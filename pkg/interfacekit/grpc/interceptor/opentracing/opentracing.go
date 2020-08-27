package grpckit_interceptor_opentracing

import (
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func UseUnaryInterceptor() grpc.UnaryServerInterceptor {
	return grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer()))
}

func UseStreamInterceptor() grpc.StreamServerInterceptor {
	return grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer()))
}
