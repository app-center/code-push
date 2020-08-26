package grpc_interceptor

import (
	"context"
	cpErrors "github.com/funnyecho/code-push/pkg/errors"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"strings"
)

const (
	KeyErrorStack = "error_stack"
)

func UnaryServerErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		if err != nil {
			cpErr := cpErrors.Error("FA_INTERNAL_ERROR")

			errors.As(err, &cpErr)
			errStack, _ := ctx.Value(KeyErrorStack).(*error)
			*errStack = err

			err = cpErr
		}

		return
	}
}

func StreamServerErrorInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		err = handler(srv, ss)

		if err != nil {
			cpErr := cpErrors.Error("FA_INTERNAL_ERROR")

			errors.As(err, &cpErr)
			errStack, _ := ss.Context().Value(KeyErrorStack).(*error)
			*errStack = err

			err = cpErr
		}

		return
	}
}

func UnaryClientErrorInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		err = invoker(ctx, method, req, reply, cc, opts...)

		if err != nil {
			s, ok := status.FromError(err)

			cpErr := cpErrors.Error("FA_INTERNAL_ERROR")

			if ok {
				message := s.Message()
				if strings.HasPrefix(message, "FA_") {
					cpErr = cpErrors.Error(message)
				}
			}

			err = cpErr
		}

		return
	}
}

func StreamClientErrorInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (clientStream grpc.ClientStream, err error) {

		clientStream, err = streamer(ctx, desc, cc, method, opts...)

		if err != nil {
			s, ok := status.FromError(err)

			cpErr := cpErrors.Error("FA_INTERNAL_ERROR")

			if ok {
				message := s.Message()
				if strings.HasPrefix(message, "FA_") {
					cpErr = cpErrors.Error(message)
				}
			}

			err = cpErr
		}

		return
	}
}
