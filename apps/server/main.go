package main

import (
	"dummy-grpc/apps/server/service"
	"dummy-grpc/lib/grpc_stats"
	pb "dummy-grpc/lib/proto/dummy"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/DataDog/datadog-go/v5/statsd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

    statsdClient, err := statsd.New("dogstatsd:8125", statsd.WithNamespace("ron.demo"))
    if err != nil {
        log.Fatalf("failed to create statsd client: %v", err)
    }

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpc_stats.UnaryServerInterceptor(statsdClient)),
	)
	pb.RegisterDummyServiceServer(s, &service.Service{})

	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
