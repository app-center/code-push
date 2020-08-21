package usecase

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	filerpb "github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/client"
	"github.com/pkg/errors"
	"time"
)

func (uc *useCase) VersionStrictCompatQuery(ctx context.Context, envId, appVersion []byte) (*client.VersionCompatQueryResult, error) {
	res, err := uc.codePush.VersionStrictCompatQuery(ctx, envId, appVersion)
	return unmarshalVersionCompatQueryResult(res), err
}

func (uc *useCase) GetVersion(ctx context.Context, envId, appVersion []byte) (*client.Version, error) {
	res, err := uc.codePush.GetVersion(ctx, envId, appVersion)
	return unmarshalVersion(res), err
}

func (uc *useCase) VersionPkgSource(ctx context.Context, envId, appVersion string) (*client.FileSource, error) {
	ver, verErr := uc.codePush.GetVersion(ctx, []byte(envId), []byte(appVersion))
	if verErr != nil {
		return nil, verErr
	}
	if ver == nil {
		return nil, errors.Wrap(client.ErrVersionNotFound, "version not existed")
	}

	fileKey := ver.PackageFileKey

	source, sourceErr := uc.filer.GetSource(ctx, []byte(fileKey))
	if sourceErr != nil {
		return nil, sourceErr
	}

	return unmarshalFileSource(source), nil
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
		AppVersion:          string(r.GetAppVersion()),
		LatestAppVersion:    string(r.GetLatestAppVersion()),
		CanUpdateAppVersion: string(r.GetCanUpdateAppVersion()),
		MustUpdate:          r.GetMustUpdate(),
	}
}

func unmarshalFileSource(v *filerpb.FileSource) *client.FileSource {
	if v == nil {
		return nil
	}

	return &client.FileSource{
		Key:        v.GetKey(),
		Value:      v.GetValue(),
		Desc:       v.GetDesc(),
		CreateTime: time.Unix(0, v.GetCreateTime()),
		FileMD5:    v.GetFileMD5(),
		FileSize:   v.GetFileSize(),
	}
}
