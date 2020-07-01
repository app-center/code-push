package grpc

import (
	"context"
	code_push "github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/code-push/usecase"
)

type codePushServer struct {
	branchUseCase  usecase.IBranch
	envUseCase     usecase.IEnv
	versionUseCase usecase.IVersion
}

func (s *codePushServer) CreateBranch(ctx context.Context, request *pb.CreateBranchRequest) (*pb.CreateBranchResponse, error) {
	res, err := s.branchUseCase.CreateBranch(request.BranchName, request.BranchAuthHost)

	return &pb.CreateBranchResponse{
		Code: MarshalErrorCode(err),
		Data: MarshalBranch(res),
	}, nil
}

func (s *codePushServer) GetBranch(ctx context.Context, request *pb.GetBranchRequest) (*pb.GetBranchResponse, error) {
	res, err := s.branchUseCase.GetBranch(request.BranchId)

	return &pb.GetBranchResponse{
		Code: MarshalErrorCode(err),
		Data: MarshalBranch(res),
	}, nil
}

func (s *codePushServer) DeleteBranch(ctx context.Context, request *pb.DeleteBranchRequest) (*pb.PlainResponse, error) {
	err := s.branchUseCase.DeleteBranch(request.BranchId)
	return &pb.PlainResponse{
		Code: MarshalErrorCode(err),
	}, nil
}

func (s *codePushServer) GetBranchEncToken(ctx context.Context, request *pb.GetBranchEncTokenRequest) (*pb.StringResponse, error) {
	res, err := s.branchUseCase.GetBranchEncToken(request.BranchId)
	return &pb.StringResponse{
		Code: MarshalErrorCode(err),
		Data: res,
	}, nil
}

func (s *codePushServer) CreateEnv(ctx context.Context, request *pb.CreateEnvRequest) (*pb.EnvResponse, error) {
	res, err := s.envUseCase.CreateEnv(request.BranchId, request.EnvName)
	return &pb.EnvResponse{
		Code: MarshalErrorCode(err),
		Data: MarshalEnv(res),
	}, nil
}

func (s *codePushServer) GetEnv(ctx context.Context, request *pb.EnvIdRequest) (*pb.EnvResponse, error) {
	res, err := s.envUseCase.GetEnv(request.EnvId)
	return &pb.EnvResponse{
		Code: MarshalErrorCode(err),
		Data: MarshalEnv(res),
	}, nil
}

func (s *codePushServer) DeleteEnv(ctx context.Context, request *pb.EnvIdRequest) (*pb.PlainResponse, error) {
	err := s.envUseCase.DeleteEnv(request.EnvId)
	return &pb.PlainResponse{Code: MarshalErrorCode(err)}, nil
}

func (s *codePushServer) GetEnvEncToken(ctx context.Context, request *pb.EnvIdRequest) (*pb.StringResponse, error) {
	res, err := s.envUseCase.GetEnvEncToken(request.EnvId)
	return &pb.StringResponse{
		Code: MarshalErrorCode(err),
		Data: res,
	}, nil
}

func (s *codePushServer) GetEnvAuthHost(ctx context.Context, request *pb.EnvIdRequest) (*pb.StringResponse, error) {
	res, err := s.envUseCase.GetEnvAuthHost(request.EnvId)
	return &pb.StringResponse{
		Code: MarshalErrorCode(err),
		Data: res,
	}, nil
}

func (s *codePushServer) ReleaseVersion(ctx context.Context, request *pb.VersionReleaseRequest) (*pb.PlainResponse, error) {
	params, err := usecase.NewVersionReleaseParams(usecase.VersionReleaseParamsConfig{
		EnvId:            request.EnvId,
		AppVersion:       request.AppVersion,
		CompatAppVersion: request.CompatAppVersion,
		Changelog:        request.Changelog,
		PackageFileKey:   request.PackageFileKey,
		MustUpdate:       request.MustUpdate,
	})

	if err != nil {
		return &pb.PlainResponse{Code: MarshalErrorCode(err)}, nil
	}

	err = s.versionUseCase.ReleaseVersion(params)
	return &pb.PlainResponse{Code: MarshalErrorCode(err)}, nil
}

func (s *codePushServer) GetVersion(ctx context.Context, request *pb.GetVersionRequest) (*pb.VersionResponse, error) {
	res, err := s.versionUseCase.GetVersion(request.EnvId, request.AppVersion)
	return &pb.VersionResponse{
		Code: MarshalErrorCode(err),
		Data: MarshalVersion(res),
	}, nil
}

func (s *codePushServer) ListVersions(ctx context.Context, request *pb.ListVersionsRequest) (*pb.VersionListResponse, error) {
	res, err := s.versionUseCase.ListVersions(request.EnvId)
	return &pb.VersionListResponse{
		Code: MarshalErrorCode(err),
		Data: MarshalVersionList(res),
	}, nil
}

func (s *codePushServer) VersionStrictCompatQuery(ctx context.Context, request *pb.VersionStrictCompatQueryRequest) (*pb.VersionStrictCompatQueryResponse, error) {
	res, err := s.versionUseCase.VersionStrictCompatQuery(request.EnvId, request.AppVersion)
	return &pb.VersionStrictCompatQueryResponse{
		Code: MarshalErrorCode(err),
		Data: MarshalVersionCompatQueryResult(res),
	}, nil
}

func MarshalErrorCode(err error) string {
	if err != nil {
		return "S_OK"
	}

	if codePushErr := err.(code_push.Error); codePushErr == "" {
		return code_push.ErrInternalError.Error()
	} else {
		return codePushErr.Error()
	}
}

func MarshalBranch(b *usecase.Branch) *pb.Branch {
	if b == nil {
		return nil
	}

	return &pb.Branch{
		BranchId:   b.BranchId,
		BranchName: b.BranchId,
		CreateTime: b.CreateTime.UnixNano(),
	}
}

func MarshalEnv(e *usecase.Env) *pb.Env {
	if e == nil {
		return nil
	}

	return &pb.Env{
		BranchId:   e.BranchId,
		EnvId:      e.EnvId,
		Name:       e.Name,
		CreateTime: e.CreateTime.UnixNano(),
	}
}

func MarshalVersion(v *usecase.Version) *pb.Version {
	if v == nil {
		return nil
	}

	return &pb.Version{
		EnvId:            v.EnvId,
		AppVersion:       v.AppVersion,
		CompatAppVersion: v.CompatAppVersion,
		MustUpdate:       v.MustUpdate,
		Changelog:        v.Changelog,
		PackageFileKey:   v.PackageFileKey,
		CreateTime:       v.CreateTime.UnixNano(),
	}
}

func MarshalVersionList(l usecase.VersionList) []*pb.Version {
	if l == nil {
		return nil
	}

	var v []*pb.Version

	for _, ver := range l {
		v = append(v, MarshalVersion(ver))
	}

	return v
}

func MarshalVersionCompatQueryResult(r usecase.IVersionCompatQueryResult) *pb.VersionCompatQueryResult {
	if r == nil {
		return nil
	}

	return &pb.VersionCompatQueryResult{
		AppVersion:          r.AppVersion(),
		LatestAppVersion:    r.LatestAppVersion(),
		CanUpdateAppVersion: r.CanUpdateAppVersion(),
		MustUpdate:          r.MustUpdate(),
	}
}
