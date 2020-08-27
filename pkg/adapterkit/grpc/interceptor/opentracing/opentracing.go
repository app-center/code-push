package adapterkit_grpc_interceptor_opentracing

import (
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func WithUnaryInterceptor() grpc.UnaryClientInterceptor {
	return grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer()))
}

func WithStreamInterceptor() grpc.StreamClientInterceptor {
	return grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer()))
}
