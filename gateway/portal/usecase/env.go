package usecase

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/portal"
	"time"
)

func (u *useCase) CreateEnv(ctx context.Context, branchId, envName []byte) (*portal.Env, error) {
	res, err := u.codePush.CreateEnv(ctx, branchId, envName)
	return unmarshalEnv(res), err
}

func (u *useCase) GetEnv(ctx context.Context, envId []byte) (*portal.Env, error) {
	res, err := u.codePush.GetEnv(ctx, envId)
	return unmarshalEnv(res), err
}

func (u *useCase) DeleteEnv(ctx context.Context, envId []byte) error {
	return u.codePush.DeleteEnv(ctx, envId)
}

func (u *useCase) GetEnvEncToken(ctx context.Context, envId []byte) ([]byte, error) {
	return u.codePush.GetEnvEncToken(ctx, envId)
}

func unmarshalEnv(e *pb.EnvResponse) *portal.Env {
	if e == nil {
		return nil
	}

	return &portal.Env{
		BranchId:   e.GetBranchId(),
		ID:         e.GetEnvId(),
		Name:       e.GetName(),
		EncToken:   e.GetEnvEncToken(),
		CreateTime: time.Unix(0, e.CreateTime),
	}
}
