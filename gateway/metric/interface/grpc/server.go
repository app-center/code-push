package grpc

import (
	"context"
	"github.com/funnyecho/code-push/gateway/metric/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/metric/usecase"
)

type Server interface {
	pb.RequestDurationServer
}

func New(uc usecase.UseCase) Server {
	svr := &server{
		uc: uc,
	}

	return svr
}

type server struct {
	uc usecase.UseCase
}

func (s *server) Http(ctx context.Context, request *pb.HttpRequestDurationRequest) (*pb.RequestDurationResponse, error) {
	s.uc.HttpDuration(request.GetSvr(), request.GetPath(), request.GetSuccess(), request.GetDurationSecond())
	return nil, nil
}

func (s *server) Grpc(ctx context.Context, request *pb.GrpcRequestDurationRequest) (*pb.RequestDurationResponse, error) {
	s.uc.GrpcDuration(request.GetSvr(), request.GetMethod(), request.GetSuccess(), request.GetDurationSecond())
	return nil, nil
}
