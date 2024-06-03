package main

import (
	"context"
	"dummy-grpc/lib/grpc_stats"
	pb "dummy-grpc/lib/proto/dummy"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "server:50051", "the address to connect to")
)

func main() {
	flag.Parse()

    statsdClient, err := statsd.New("dogstatsd:8125", statsd.WithNamespace("ron.demo"))
    if err != nil {
        log.Fatalf("failed to create statsd client: %v", err)
    }

	hostname, err := os.Hostname()
    if err != nil {
        log.Fatalf("failed to retrieve hostname: %v", err)
    }
	tags := []string{fmt.Sprintf("hostname:%s", hostname)}

	conn, err := grpc.NewClient(
        *addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "random": {} } ]}`),
        grpc.WithChainUnaryInterceptor(grpc_stats.UnaryClientInterceptor(statsdClient, tags)),
    )
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDummyServiceClient(conn)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	ctx := context.Background()
	for {
		<-ticker.C

		_, err = c.DoSomething(ctx, &pb.DoSomethingRequest{})
		if err != nil {
			log.Fatalf("could not do something: %v", err)
		}
		log.Printf("did something")
        statsdClient.Count("request", 1, []string{}, 1.0)
	}
}
