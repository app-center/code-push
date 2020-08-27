package grpckit_interceptor_error

import (
	"context"
	"github.com/funnyecho/code-push/pkg/errors"
	stdErrors "github.com/pkg/errors"
	"google.golang.org/grpc"
)

func UseUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		if err != nil {
			cpErr := errors.Error("INTERNAL_ERROR")
			stdErrors.As(err, &cpErr)
			err = &grpcErr{
				err:   cpErr,
				cause: err,
			}
		}

		return
	}
}

func UseStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		err = handler(srv, ss)

		if err != nil {
			cpErr := errors.Error("INTERNAL_ERROR")

			stdErrors.As(err, &cpErr)
			err = &grpcErr{
				err:   cpErr,
				cause: err,
			}
		}

		return
	}
}

type grpcErr struct {
	err   errors.Error
	cause error
}

func (e *grpcErr) Unwrap() error {
	return e.cause
}

func (e *grpcErr) Cause() error {
	return e.cause
}

func (e *grpcErr) Error() string {
	return e.err.Error()
}
