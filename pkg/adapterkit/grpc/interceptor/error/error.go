package adapterkit_grpc_interceptor_error

import (
	"context"
	"github.com/funnyecho/code-push/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"strings"
)

func WithUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		err = invoker(ctx, method, req, reply, cc, opts...)

		if err != nil {
			s, ok := status.FromError(err)

			cpErr := errors.Error("INTERNAL_ERROR")

			if ok {
				message := s.Message()
				if strings.HasPrefix(message, "FA") {
					cpErr = errors.Error(message)
				}
			}

			err = cpErr
		}

		return
	}
}

func WithStreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (clientStream grpc.ClientStream, err error) {
		clientStream, err = streamer(ctx, desc, cc, method, opts...)

		if err != nil {
			s, ok := status.FromError(err)

			cpErr := errors.Error("INTERNAL_ERROR")

			if ok {
				message := s.Message()
				if strings.HasPrefix(message, "FA") {
					cpErr = errors.Error(message)
				}
			}

			err = cpErr
		}

		return
	}
}
