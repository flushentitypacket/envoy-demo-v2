package service

import (
	"context"
	pb "dummy-grpc/lib/proto/dummy"
	"log"
)

type Service struct {
	pb.UnimplementedDummyServiceServer
}

func (s *Service) DoSomething(ctx context.Context, in *pb.DoSomethingRequest) (*pb.DoSomethingResponse, error) {
	log.Printf("got req")
	return &pb.DoSomethingResponse{}, nil
}
