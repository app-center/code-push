package usecase

import (
	"context"
	"github.com/funnyecho/code-push/daemon/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway"
	"time"
)

func (uc *useCase) CreateEnv(ctx context.Context, branchId, envName []byte) (*gateway.Env, error) {
	res, err := uc.daemon.CreateEnv(ctx, branchId, envName)
	return unmarshalEnv(res), err
}

func (uc *useCase) GetEnv(ctx context.Context, envId []byte) (*gateway.Env, error) {
	res, err := uc.daemon.GetEnv(ctx, envId)
	return unmarshalEnv(res), err
}

func (uc *useCase) DeleteEnv(ctx context.Context, envId []byte) error {
	return uc.daemon.DeleteEnv(ctx, envId)
}

func (uc *useCase) GetEnvEncToken(ctx context.Context, envId []byte) ([]byte, error) {
	return uc.daemon.GetEnvEncToken(ctx, envId)
}

func unmarshalEnv(e *pb.EnvResponse) *gateway.Env {
	if e == nil {
		return nil
	}

	return &gateway.Env{
		BranchId:   e.GetBranchId(),
		ID:         e.GetEnvId(),
		Name:       e.GetName(),
		EncToken:   e.GetEnvEncToken(),
		CreateTime: time.Unix(0, e.CreateTime),
	}
}
