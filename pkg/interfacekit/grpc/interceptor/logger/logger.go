package grpckit_interceptor_logger

import (
	"context"
	"github.com/funnyecho/code-push/pkg/log"
	"google.golang.org/grpc"
)

func UseUnaryInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		method := info.FullMethod

		resp, err = handler(ctx, req)

		if err != nil {
			logger.Debug(
				"grpc-request failed",
				"method", method,
				"errorMessage", err.Error(),
				"req", req,
			)
		} else {
			logger.Debug(
				"grpc-request completed",
				"method", method,
			)
		}

		return
	}
}

func UseStreamInterceptor(logger log.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		method := info.FullMethod

		err = handler(srv, ss)

		if err != nil {
			logger.Debug(
				"grpc-request failed",
				"method", method,
				"errorMessage", err.Error(),
			)
		} else {
			logger.Debug(
				"grpc-request completed",
				"method", method,
			)
		}

		return
	}
}
