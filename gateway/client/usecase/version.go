package usecase

import (
	"github.com/funnyecho/code-push/gateway/client"
	"github.com/pkg/errors"
)

func (u useCase) VersionStrictCompatQuery(envId, appVersion []byte) (*client.VersionCompatQueryResult, error) {
	return u.codePush.VersionStrictCompatQuery(envId, appVersion)
}

func (u *useCase) GetVersion(envId, appVersion []byte) (*client.Version, error) {
	return u.codePush.GetVersion(envId, appVersion)
}

func (u *useCase) VersionDownloadPkg(envId, appVersion []byte) ([]byte, error) {
	ver, verErr := u.codePush.GetVersion(envId, appVersion)
	if verErr != nil {
		return nil, verErr
	}
	if ver == nil {
		return nil, errors.Wrap(client.ErrVersionNotFound, "version not existed")
	}

	fileKey := ver.PackageFileKey
	return u.filer.GetSource([]byte(fileKey))
}
