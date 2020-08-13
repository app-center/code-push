package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/client"
	"github.com/pkg/errors"
	"time"
)

func (uc *useCase) VersionStrictCompatQuery(envId, appVersion []byte) (*client.VersionCompatQueryResult, error) {
	res, err := uc.codePush.VersionStrictCompatQuery(envId, appVersion)
	return unmarshalVersionCompatQueryResult(res), err
}

func (uc *useCase) GetVersion(envId, appVersion []byte) (*client.Version, error) {
	res, err := uc.codePush.GetVersion(envId, appVersion)
	return unmarshalVersion(res), err
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

func unmarshalVersion(v *pb.VersionResponse) *client.Version {
	if v == nil {
		return nil
	}

	return &client.Version{
		EnvId:            v.GetEnvId(),
		AppVersion:       v.GetAppVersion(),
		CompatAppVersion: v.GetCompatAppVersion(),
		MustUpdate:       v.GetMustUpdate(),
		Changelog:        v.GetChangelog(),
		PackageFileKey:   v.GetPackageFileKey(),
		CreateTime:       time.Unix(0, v.GetCreateTime()),
	}
}

func unmarshalVersionCompatQueryResult(r *pb.VersionStrictCompatQueryResponse) *client.VersionCompatQueryResult {
	if r == nil {
		return nil
	}

	return &client.VersionCompatQueryResult{
		AppVersion:          r.GetAppVersion(),
		LatestAppVersion:    r.GetLatestAppVersion(),
		CanUpdateAppVersion: r.GetCanUpdateAppVersion(),
		MustUpdate:          r.GetMustUpdate(),
	}
}
