package grpc

import (
	"context"
	code_push "github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/code-push/usecase"
	"github.com/funnyecho/code-push/pkg/log"
)

func NewCodePushServer(endpoints Endpoints, logger log.Logger) *codePushServer {
	return &codePushServer{endpoints: endpoints, Logger: logger}
}

type codePushServer struct {
	endpoints Endpoints
	log.Logger
}

func (s *codePushServer) CreateBranch(ctx context.Context, request *pb.CreateBranchRequest) (*pb.BranchResponse, error) {
	res, err := s.endpoints.CreateBranch(request.GetBranchName())

	return MarshalBranchResponse(res), err
}

func (s *codePushServer) GetBranch(ctx context.Context, request *pb.GetBranchRequest) (*pb.BranchResponse, error) {
	res, err := s.endpoints.GetBranch(request.BranchId)

	return MarshalBranchResponse(res), err
}

func (s *codePushServer) DeleteBranch(ctx context.Context, request *pb.DeleteBranchRequest) (*pb.PlainResponse, error) {
	err := s.endpoints.DeleteBranch(request.BranchId)
	return nil, err
}

func (s *codePushServer) GetBranchEncToken(ctx context.Context, request *pb.GetBranchEncTokenRequest) (*pb.StringResponse, error) {
	res, err := s.endpoints.GetBranchEncToken(request.BranchId)
	return MarshalBytesToStringResponse(res), err
}

func (s *codePushServer) CreateEnv(ctx context.Context, request *pb.CreateEnvRequest) (*pb.EnvResponse, error) {
	res, err := s.endpoints.CreateEnv(request.BranchId, request.EnvName)
	return MarshalEnvResponse(res), err
}

func (s *codePushServer) GetEnv(ctx context.Context, request *pb.EnvIdRequest) (*pb.EnvResponse, error) {
	res, err := s.endpoints.GetEnv(request.EnvId)
	return MarshalEnvResponse(res), err
}

func (s *codePushServer) DeleteEnv(ctx context.Context, request *pb.EnvIdRequest) (*pb.PlainResponse, error) {
	err := s.endpoints.DeleteEnv(request.EnvId)
	return nil, err
}

func (s *codePushServer) GetEnvEncToken(ctx context.Context, request *pb.EnvIdRequest) (*pb.StringResponse, error) {
	res, err := s.endpoints.GetEnvEncToken(request.EnvId)
	return MarshalBytesToStringResponse(res), err
}

func (s *codePushServer) ReleaseVersion(ctx context.Context, request *pb.VersionReleaseRequest) (*pb.PlainResponse, error) {
	err := s.endpoints.ReleaseVersion(UnmarshalVersionReleaseParams(request))
	return nil, err
}

func (s *codePushServer) GetVersion(ctx context.Context, request *pb.GetVersionRequest) (*pb.VersionResponse, error) {
	res, err := s.endpoints.GetVersion(request.EnvId, request.AppVersion)
	return MarshalVersionResponse(res), err
}

func (s *codePushServer) ListVersions(ctx context.Context, request *pb.ListVersionsRequest) (*pb.VersionListResponse, error) {
	res, err := s.endpoints.ListVersions(request.EnvId)
	return MarshalVersionList(res), err
}

func (s *codePushServer) VersionStrictCompatQuery(ctx context.Context, request *pb.VersionStrictCompatQueryRequest) (*pb.VersionStrictCompatQueryResponse, error) {
	res, err := s.endpoints.VersionStrictCompatQuery(request.EnvId, request.AppVersion)
	return MarshalVersionCompatQueryResultResponse(res), err
}

func MarshalBranchResponse(b *code_push.Branch) *pb.BranchResponse {
	if b == nil {
		return nil
	}

	return &pb.BranchResponse{
		BranchId:       b.ID,
		BranchName:     b.Name,
		BranchEncToken: b.EncToken,
		CreateTime:     b.CreateTime.UnixNano(),
	}
}

func MarshalEnvResponse(e *code_push.Env) *pb.EnvResponse {
	if e == nil {
		return nil
	}

	return &pb.EnvResponse{
		BranchId:    e.BranchId,
		EnvId:       e.ID,
		Name:        e.Name,
		EnvEncToken: e.EncToken,
		CreateTime:  e.CreateTime.UnixNano(),
	}
}

func MarshalVersionResponse(v *code_push.Version) *pb.VersionResponse {
	if v == nil {
		return nil
	}

	return &pb.VersionResponse{
		EnvId:            v.EnvId,
		AppVersion:       v.AppVersion,
		CompatAppVersion: v.CompatAppVersion,
		MustUpdate:       v.MustUpdate,
		Changelog:        v.Changelog,
		PackageFileKey:   v.PackageFileKey,
		CreateTime:       v.CreateTime.UnixNano(),
	}
}

func MarshalVersionList(l code_push.VersionList) *pb.VersionListResponse {
	if l == nil {
		return nil
	}

	var v []*pb.VersionResponse

	for _, ver := range l {
		v = append(v, MarshalVersionResponse(ver))
	}

	return &pb.VersionListResponse{List: v}
}

func MarshalVersionCompatQueryResultResponse(r usecase.VersionCompatQueryResult) *pb.VersionStrictCompatQueryResponse {
	if r == nil {
		return nil
	}

	return &pb.VersionStrictCompatQueryResponse{
		AppVersion:          r.AppVersion(),
		LatestAppVersion:    r.LatestAppVersion(),
		CanUpdateAppVersion: r.CanUpdateAppVersion(),
		MustUpdate:          r.MustUpdate(),
	}
}

func MarshalBytesResponse(bytes []byte) *pb.BytesResponse {
	if bytes == nil {
		return nil
	}

	return &pb.BytesResponse{Data: bytes}
}

func MarshalBytesToStringResponse(bytes []byte) *pb.StringResponse {
	if bytes == nil {
		return nil
	}

	return &pb.StringResponse{Data: string(bytes)}
}

func UnmarshalVersionReleaseParams(request *pb.VersionReleaseRequest) usecase.VersionReleaseParams {
	return NewVersionReleaseParams(request)
}
