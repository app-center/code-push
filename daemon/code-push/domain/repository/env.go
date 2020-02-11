package repository

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/model"
)

type IEnv interface {
	FindEnv(envId string) (model.Env, error)
	FindEnvsWithBranchId(branchId string) (model.EnvList, error)
	SaveEnv(env model.Env) (model.Env, error)
	DeleteEnv(envId string) error
}
