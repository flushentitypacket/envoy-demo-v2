package grpc_stats

import (
	"context"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"google.golang.org/grpc"
)

func UnaryClientInterceptor(statsdClient statsd.ClientInterface) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		callOpts ...grpc.CallOption,
	) error {
		statsdClient.Incr("grpc.client.call.total", []string{}, 1.0)
		err := invoker(ctx, method, req, reply, cc, callOpts...)
		return err
	}
}

func UnaryServerInterceptor(statsdClient statsd.ClientInterface) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		statsdClient.Incr("grpc.server.call.total", []string{}, 1.0)

		t_start := time.Now()

		resp, err := handler(ctx, req)

		elapsedTime := float64(time.Since(t_start)) / float64(time.Millisecond)
		statsdClient.Distribution("grpc.server.latency", elapsedTime, []string{}, 1.0)

		return resp, err
	}
}
