package repository

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/model"
)

type IEnv interface {
	Find(envId string) (*model.Env, error)
	FindAllWithBranchId(branchId string) (model.EnvList, error)
	Create(env model.Env) error
	Save(env *model.Env) error
	Delete(env *model.Env) error
}
