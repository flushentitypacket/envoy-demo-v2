package service

import (
	"context"
	pb "dummy-grpc/lib/proto/dummy"
)

type Service struct {
	pb.UnimplementedDummyServiceServer
}

func (s *Service) DoSomething(ctx context.Context, in *pb.DoSomethingRequest) (*pb.DoSomethingResponse, error) {
	return &pb.DoSomethingResponse{}, nil
}
