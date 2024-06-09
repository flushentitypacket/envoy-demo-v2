package grpc_stats

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"google.golang.org/grpc"
)

func UnaryClientInterceptor(statsdClient statsd.ClientInterface, tags []string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		callOpts ...grpc.CallOption,
	) error {
		statsdClient.Incr("grpc.client.call.total", tags, 1.0)
		err := invoker(ctx, method, req, reply, cc, callOpts...)
		return err
	}
}

func UnaryServerInterceptor(statsdClient statsd.ClientInterface, tags []string) grpc.UnaryServerInterceptor {
	var inFlightCount atomic.Int64
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		statsdClient.Incr("grpc.server.call.total", tags, 1.0)
		inFlightCount.Add(1)
		statsdClient.Gauge("grpc.server.call.inflight", float64(inFlightCount.Load()), tags, 1.0)
		defer inFlightCount.Add(-1)

		t_start := time.Now()

		resp, err := handler(ctx, req)

		elapsedTime := float64(time.Since(t_start)) / float64(time.Millisecond)
		statsdClient.Distribution("grpc.server.latency", elapsedTime, tags, 1.0)

		return resp, err
	}
}
