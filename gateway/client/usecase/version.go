package usecase

import (
	"github.com/funnyecho/code-push/gateway/client"
	"github.com/pkg/errors"
)

func (uc *useCase) VersionStrictCompatQuery(envId, appVersion []byte) (*client.VersionCompatQueryResult, error) {
	return uc.codePush.VersionStrictCompatQuery(envId, appVersion)
}

func (uc *useCase) GetVersion(envId, appVersion []byte) (*client.Version, error) {
	return uc.codePush.GetVersion(envId, appVersion)
}

func (uc *useCase) VersionDownloadPkg(envId, appVersion []byte) ([]byte, error) {
	ver, verErr := uc.codePush.GetVersion(envId, appVersion)
	if verErr != nil {
		return nil, verErr
	}
	if ver == nil {
		return nil, errors.Wrap(client.ErrVersionNotFound, "version not existed")
	}

	fileKey := ver.PackageFileKey
	return uc.filer.GetSource([]byte(fileKey))
}
