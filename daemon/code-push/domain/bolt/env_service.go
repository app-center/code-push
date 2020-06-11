package bolt

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain"
)

var _ domain.IEnvService = &EnvService{}

type EnvService struct {
	client *Client
}

func (e *EnvService) Env(envId string) (*domain.Env, error) {
	panic("implement me")
}

func (e *EnvService) EnvWithBranchIdAndEnvName(branchId, envName string) (*domain.Env, error) {
	panic("implement me")
}

func (e *EnvService) CreateEnv(env *domain.Env) error {
	panic("implement me")
}

func (e *EnvService) SetEnvName(envId, newEnvName string) error {
	panic("implement me")
}

func (e *EnvService) DeleteEnv(envId string) error {
	panic("implement me")
}
