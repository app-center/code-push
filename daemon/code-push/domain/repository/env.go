package repository

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/model"
)

type IEnv interface {
	FirstEnv(envId string) (*model.Env, error)
	FindEnvWithBranchId(branchId string) (model.EnvList, error)
	FirstEnvWithBranchIdAndEnvName(branchId, envName string) (*model.Env, error)
	SaveEnv(env model.Env) (*model.Env, error)
	DeleteEnv(envId string) error
}
