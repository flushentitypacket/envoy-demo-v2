package main

import (
	"context"
	pb "dummy-grpc/lib/proto/dummy"
	"flag"
	"log"
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

    statsdClient, err := statsd.New("dogstatsd:8125", statsd.WithNamespace("client"))
    if err != nil {
        log.Fatalf("failed to create statsd client: %v", err)
    }

	conn, err := grpc.NewClient(
        *addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`),
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
