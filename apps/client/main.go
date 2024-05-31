package main

import (
	"context"
	pb "dummy-grpc/lib/proto/dummy"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "server:50051", "the address to connect to")
)

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
    }
}
