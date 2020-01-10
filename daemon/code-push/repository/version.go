package repository

import "github.com/funnyecho/code-push/daemon/code-push/model"

type IVersion interface {
	Find(envId, appVersion string) (*model.Version, error)
	FindAllWithEnvId(envId string) (*model.VersionList, error)
	Create(version model.Version) (*model.Version, error)
	Save(env *model.Env) error
	Delete(env *model.Env) error
}