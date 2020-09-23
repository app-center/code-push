package usecase

import (
	"context"
	"github.com/funnyecho/code-push/daemon/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway"
	"github.com/pkg/errors"
	"time"
)

func (uc *useCase) ReleaseVersion(ctx context.Context, params *gateway.VersionReleaseParams) error {
	return uc.daemon.ReleaseVersion(ctx, marshalVersionReleaseParams(params))
}

func (uc *useCase) GetVersion(ctx context.Context, envId, appVersion []byte) (*gateway.Version, error) {
	res, err := uc.daemon.GetVersion(ctx, envId, appVersion)
	return unmarshalVersion(res), err
}

func (uc *useCase) VersionStrictCompatQuery(ctx context.Context, envId, appVersion []byte) (*gateway.VersionCompatQueryResult, error) {
	res, err := uc.daemon.VersionStrictCompatQuery(ctx, envId, appVersion)
	return unmarshalVersionCompatQueryResult(res), err
}

func (uc *useCase) VersionPkgSource(ctx context.Context, envId, appVersion string) (*gateway.FileSource, error) {
	ver, verErr := uc.daemon.GetVersion(ctx, []byte(envId), []byte(appVersion))
	if verErr != nil {
		return nil, verErr
	}
	if ver == nil {
		return nil, errors.Wrap(gateway.ErrVersionNotFound, "version not existed")
	}

	fileKey := ver.PackageFileKey

	source, sourceErr := uc.daemon.GetSource(ctx, []byte(fileKey))
	if sourceErr != nil {
		return nil, sourceErr
	}

	return unmarshalFileSource(source), nil
}

func unmarshalVersion(v *pb.VersionResponse) *gateway.Version {
	if v == nil {
		return nil
	}

	return &gateway.Version{
		EnvId:            v.GetEnvId(),
		AppVersion:       v.GetAppVersion(),
		CompatAppVersion: v.GetCompatAppVersion(),
		MustUpdate:       v.GetMustUpdate(),
		Changelog:        v.GetChangelog(),
		PackageFileKey:   v.GetPackageFileKey(),
		CreateTime:       time.Unix(0, v.GetCreateTime()),
	}
}

func unmarshalVersionCompatQueryResult(r *pb.VersionStrictCompatQueryResponse) *gateway.VersionCompatQueryResult {
	if r == nil {
		return nil
	}

	return &gateway.VersionCompatQueryResult{
		AppVersion:          r.GetAppVersion(),
		LatestAppVersion:    r.GetLatestAppVersion(),
		CanUpdateAppVersion: r.GetCanUpdateAppVersion(),
		MustUpdate:          r.GetMustUpdate(),
	}
}

func marshalVersionReleaseParams(p *gateway.VersionReleaseParams) *pb.VersionReleaseRequest {
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

