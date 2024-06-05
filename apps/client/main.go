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
    requestsPerSecond = flag.Float64("requests_per_second", 1.0, "the number of requests per second")
    operationMillis = flag.Int64("operation_millis", 10, "how long each requested operation should take in milliseconds")
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
        grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`),
        grpc.WithChainUnaryInterceptor(grpc_stats.UnaryClientInterceptor(statsdClient, tags)),
    )
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDummyServiceClient(conn)

    tickerSeconds := (1 / *requestsPerSecond)
    tickerMicros := int64(tickerSeconds * 1e6)
	ticker := time.NewTicker(time.Microsecond * time.Duration(tickerMicros))
	defer ticker.Stop()
	ctx := context.Background()
	for {
		<-ticker.C

		go func() {
			_, err = c.DoSomething(ctx, &pb.DoSomethingRequest{
				OperationMillis: *operationMillis,
			})
			if err != nil {
				log.Fatalf("could not do something: %v", err)
			}
			statsdClient.Count("request", 1, []string{}, 1.0)
		}()
	}
}
