package service

import (
	"context"
	pb "dummy-grpc/lib/proto/dummy"
	"time"
)

type Service struct {
	pb.UnimplementedDummyServiceServer
	additionalLatency time.Duration
}

func (s *Service) DoSomething(ctx context.Context, in *pb.DoSomethingRequest) (*pb.DoSomethingResponse, error) {
	requestOpLatency := time.Duration(in.GetOperationMillis()) * time.Millisecond
	time.Sleep(requestOpLatency + s.additionalLatency)
	return &pb.DoSomethingResponse{}, nil
}

func NewService(additionalLatency time.Duration) *Service {
	return &Service{
		additionalLatency: additionalLatency,
	}
}
