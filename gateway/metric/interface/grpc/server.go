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

func (s *server) Gateway(ctx context.Context, request *pb.GatewayRequestDurationRequest) (*pb.RequestDurationResponse, error) {
	s.uc.GatewayDuration(request.GetSvr(), request.GetProto(), request.GetPath(), request.GetSuccess(), request.GetDurationSecond())
	return nil, nil
}

func (s *server) Daemon(ctx context.Context, request *pb.DaemonRequestDurationRequest) (*pb.RequestDurationResponse, error) {
	s.uc.DaemonDuration(request.GetSvr(), request.GetProto(), request.GetMethod(), request.GetSuccess(), request.GetDurationSecond())
	return nil, nil
}
