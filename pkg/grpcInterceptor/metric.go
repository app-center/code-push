package grpcInterceptor

import (
	"context"
	"github.com/funnyecho/code-push/pkg/grpcInterceptor/internal"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"time"
)

func UnaryServerMetricInterceptor(logger internal.InterceptorLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()

		toLog := []interface{}{
			"method", info.FullMethod,
			"request", req,
		}

		defer func() {
			toLog = append(toLog, "duration.time_ms", durationToMilliseconds(time.Since(startTime)))
			logger.Info(
				"grpc unary server",
				toLog...,
			)
		}()

		var errStack error

		ctx = context.WithValue(ctx, KeyErrorStack, &errStack)

		resp, err = handler(ctx, req)

		if err != nil {
			toLog = append(toLog, "err", err)
			errStack := ctx.Value(KeyErrorStack)
			if errStack != nil {
				toLog = append(toLog, "errStack", errStack)
			}
		} else {
			toLog = append(toLog, "response", resp)
		}

		return
	}
}

func StreamServerMetricInterceptor(logger internal.InterceptorLogger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		startTime := time.Now()

		toLog := []interface{}{
			"method", info.FullMethod,
		}

		defer func() {
			toLog = append(toLog, "duration.time_ms", durationToMilliseconds(time.Since(startTime)))
			logger.Info(
				"grpc stream server",
				toLog...,
			)
		}()

		var errStack error

		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = context.WithValue(ss.Context(), KeyErrorStack, &errStack)

		err = handler(srv, wrapped)
		if err != nil {
			toLog = append(toLog, "err", err)
			errStack := ss.Context().Value(KeyErrorStack)
			if errStack != nil {
				toLog = append(toLog, "errStack", errStack)
			}
		}

		return
	}
}

func UnaryClientMetricInterceptor(logger internal.InterceptorLogger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		startTime := time.Now()

		toLog := []interface{}{
			"method", method,
			"request", req,
		}

		defer func() {
			toLog = append(toLog, "duration.time_ms", durationToMilliseconds(time.Since(startTime)))
			logger.Info(
				"grpc unary client",
				toLog...,
			)
		}()

		err = invoker(ctx, method, req, reply, cc, opts...)

		if err != nil {
			toLog = append(toLog, "err", err)
		} else {
			toLog = append(toLog, "response", reply)
		}

		return
	}
}

func StreamClientMetricInterceptor(logger internal.InterceptorLogger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (clientStream grpc.ClientStream, err error) {
		startTime := time.Now()

		toLog := []interface{}{
			"method", method,
		}

		defer func() {
			toLog = append(toLog, "duration.time_ms", durationToMilliseconds(time.Since(startTime)))
			logger.Info(
				"grpc stream client",
				toLog...,
			)
		}()

		clientStream, err = streamer(ctx, desc, cc, method, opts...)

		if err != nil {
			toLog = append(toLog, "err", err)
		}

		return
	}
}

func durationToMilliseconds(duration time.Duration) float32 {
	return float32(duration.Nanoseconds()/1000) / 1000
}
