package grpckit_interceptor_metrics

import (
	"context"
	grpckit_metrics "github.com/funnyecho/code-push/pkg/interfacekit/grpc/metrics"
	"google.golang.org/grpc"
	"time"
)

func UseUnaryInterceptor(metrics *grpckit_metrics.Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()
		method := info.FullMethod

		resp, err = handler(ctx, req)

		durationSeconds := durationToSeconds(time.Since(startTime))
		errCode := ""

		if err == nil {
			metrics.IncGRPCRequestSucceed(method)
		} else {
			errCode = err.Error()
			metrics.IncGRPCRequestFailed(method, errCode)
		}

		metrics.ObserveGRPCRequestDuration(durationSeconds, method, errCode)

		return
	}
}

func UseStreamInterceptor(metrics *grpckit_metrics.Metrics) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		startTime := time.Now()
		method := info.FullMethod

		err = handler(srv, ss)

		durationSeconds := durationToSeconds(time.Since(startTime))
		errCode := ""

		if err == nil {
			metrics.IncGRPCRequestSucceed(method)
		} else {
			errCode = err.Error()
			metrics.IncGRPCRequestFailed(method, errCode)
		}

		metrics.ObserveGRPCRequestDuration(durationSeconds, method, errCode)

		return
	}
}

func durationToSeconds(duration time.Duration) float64 {
	return float64(duration.Nanoseconds()/1000) / 1000 / 1000
}
