package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/portal"
	"time"
)

func (u *useCase) ReleaseVersion(params *portal.VersionReleaseParams) error {
	return u.codePush.ReleaseVersion(marshalVersionReleaseParams(params))
}

func (u *useCase) GetVersion(envId, appVersion []byte) (*portal.Version, error) {
	res, err := u.codePush.GetVersion(envId, appVersion)
	return unmarshalVersion(res), err
}

func (u *useCase) VersionStrictCompatQuery(envId, appVersion []byte) (*portal.VersionCompatQueryResult, error) {
	res, err := u.codePush.VersionStrictCompatQuery(envId, appVersion)
	return unmarshalVersionCompatQueryResult(res), err
}

func unmarshalVersion(v *pb.VersionResponse) *portal.Version {
	if v == nil {
		return nil
	}

	return &portal.Version{
		EnvId:            v.GetEnvId(),
		AppVersion:       v.GetAppVersion(),
		CompatAppVersion: v.GetCompatAppVersion(),
		MustUpdate:       v.GetMustUpdate(),
		Changelog:        v.GetChangelog(),
		PackageFileKey:   v.GetPackageFileKey(),
		CreateTime:       time.Unix(0, v.GetCreateTime()),
	}
}

func unmarshalVersionCompatQueryResult(r *pb.VersionStrictCompatQueryResponse) *portal.VersionCompatQueryResult {
	if r == nil {
		return nil
	}

	return &portal.VersionCompatQueryResult{
		AppVersion:          r.GetAppVersion(),
		LatestAppVersion:    r.GetLatestAppVersion(),
		CanUpdateAppVersion: r.GetCanUpdateAppVersion(),
		MustUpdate:          r.GetMustUpdate(),
	}
}

func marshalVersionReleaseParams(p *portal.VersionReleaseParams) *pb.VersionReleaseRequest {
	if p == nil {
		return nil
	}

	return &pb.VersionReleaseRequest{
		EnvId:            p.EnvId,
		AppVersion:       p.AppVersion,
		CompatAppVersion: p.CompatAppVersion,
		Changelog:        p.Changelog,
		PackageFileKey:   p.PackageFileKey,
		MustUpdate:       p.MustUpdate,
	}
}
