package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/portal"
	"time"
)

func (u *useCase) CreateEnv(branchId, envName []byte) (*portal.Env, error) {
	res, err := u.codePush.CreateEnv(branchId, envName)
	return unmarshalEnv(res), err
}

func (u *useCase) GetEnv(envId []byte) (*portal.Env, error) {
	res, err := u.codePush.GetEnv(envId)
	return unmarshalEnv(res), err
}

func (u *useCase) DeleteEnv(envId []byte) error {
	return u.codePush.DeleteEnv(envId)
}

func (u *useCase) GetEnvEncToken(envId []byte) ([]byte, error) {
	return u.codePush.GetEnvEncToken(envId)
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
