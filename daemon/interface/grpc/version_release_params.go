package grpc

import "github.com/funnyecho/code-push/daemon/interface/grpc/pb"

func NewVersionReleaseParams(request *pb.VersionReleaseRequest) *versionReleaseParams {
	return &versionReleaseParams{request}
}

type versionReleaseParams struct {
	request *pb.VersionReleaseRequest
}

func (v *versionReleaseParams) EnvId() []byte {
	return v.request.GetEnvId()
}

func (v *versionReleaseParams) AppVersion() []byte {
	return v.request.GetAppVersion()
}

func (v *versionReleaseParams) CompatAppVersion() []byte {
	return v.request.GetCompatAppVersion()
}

func (v *versionReleaseParams) Changelog() []byte {
	return v.request.GetChangelog()
}

func (v *versionReleaseParams) PackageFileKey() []byte {
	return v.request.GetPackageFileKey()
}

func (v *versionReleaseParams) MustUpdate() bool {
	return v.request.GetMustUpdate()
}
