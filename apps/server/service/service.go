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
	timer := time.NewTimer(requestOpLatency + s.additionalLatency)
	for {
		select {
		case <-timer.C:
			return &pb.DoSomethingResponse{}, nil
		default:
		}
	}
}

func NewService(additionalLatency time.Duration) *Service {
	return &Service{
		additionalLatency: additionalLatency,
	}
}
