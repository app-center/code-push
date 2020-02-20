package repository

import "github.com/funnyecho/code-push/daemon/code-push/domain/model"

type IVersion interface {
	FindVersion(envId, appVersion string) (*model.Version, error)
	FindVersionsWithEnvId(envId string) (model.VersionList, error)
	SaveVersion(version model.Version) (*model.Version, error)
	DeleteVersion(envId, appVersion string) error
}
