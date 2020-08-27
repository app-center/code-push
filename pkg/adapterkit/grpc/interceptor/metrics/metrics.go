package adapterkit_grpc_interceptor_metrics

import "google.golang.org/grpc"

func WithUnaryInterceptor() grpc.UnaryClientInterceptor {
	return nil
}

func WithStreamInterceptor() grpc.StreamClientInterceptor {
	return nil
}
