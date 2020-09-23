package grpc

import (
	"context"
	"github.com/funnyecho/code-push/daemon"
	"github.com/funnyecho/code-push/daemon/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/usecase"
	grpc_streamer "github.com/funnyecho/code-push/pkg/grpc-streamer"
	"github.com/funnyecho/code-push/pkg/log"
	"io"
)

func NewServer(configFn func(*ServerConfig)) *server {
	config := &ServerConfig{}
	configFn(config)

	return &server{uc: config.UseCase, Logger: config.Logger}
}

type server struct {
	uc usecase.UseCase
	log.Logger
}

type ServerConfig struct {
	usecase.UseCase
	log.Logger
}

func (s *server) CreateBranch(ctx context.Context, request *pb.CreateBranchRequest) (*pb.BranchResponse, error) {
	res, err := s.uc.CreateBranch(request.GetBranchName())

	return MarshalBranchResponse(res), err
}

func (s *server) GetBranch(ctx context.Context, request *pb.GetBranchRequest) (*pb.BranchResponse, error) {
	res, err := s.uc.GetBranch(request.BranchId)

	return MarshalBranchResponse(res), err
}

func (s *server) DeleteBranch(ctx context.Context, request *pb.DeleteBranchRequest) (*pb.PlainResponse, error) {
	err := s.uc.DeleteBranch(request.BranchId)
	return nil, err
}

func (s *server) GetBranchEncToken(ctx context.Context, request *pb.GetBranchEncTokenRequest) (*pb.StringResponse, error) {
	res, err := s.uc.GetBranchEncToken(request.BranchId)
	return MarshalBytesToStringResponse(res), err
}

func (s *server) CreateEnv(ctx context.Context, request *pb.CreateEnvRequest) (*pb.EnvResponse, error) {
	res, err := s.uc.CreateEnv(request.BranchId, request.EnvName)
	return MarshalEnvResponse(res), err
}

func (s *server) GetEnv(ctx context.Context, request *pb.EnvIdRequest) (*pb.EnvResponse, error) {
	res, err := s.uc.GetEnv(request.EnvId)
	return MarshalEnvResponse(res), err
}

func (s *server) DeleteEnv(ctx context.Context, request *pb.EnvIdRequest) (*pb.PlainResponse, error) {
	err := s.uc.DeleteEnv(request.EnvId)
	return nil, err
}

func (s *server) GetEnvEncToken(ctx context.Context, request *pb.EnvIdRequest) (*pb.StringResponse, error) {
	res, err := s.uc.GetEnvEncToken(request.EnvId)
	return MarshalBytesToStringResponse(res), err
}

func (s *server) ReleaseVersion(ctx context.Context, request *pb.VersionReleaseRequest) (*pb.PlainResponse, error) {
	err := s.uc.ReleaseVersion(UnmarshalVersionReleaseParams(request))
	return nil, err
}

func (s *server) GetVersion(ctx context.Context, request *pb.GetVersionRequest) (*pb.VersionResponse, error) {
	res, err := s.uc.GetVersion(request.EnvId, request.AppVersion)
	return MarshalVersionResponse(res), err
}

func (s *server) ListVersions(ctx context.Context, request *pb.ListVersionsRequest) (*pb.VersionListResponse, error) {
	res, err := s.uc.ListVersions(request.EnvId)
	return MarshalVersionList(res), err
}

func (s *server) VersionStrictCompatQuery(ctx context.Context, request *pb.VersionStrictCompatQueryRequest) (*pb.VersionStrictCompatQueryResponse, error) {
	res, err := s.uc.VersionStrictCompatQuery(request.EnvId, request.AppVersion)
	return MarshalVersionCompatQueryResultResponse(res), err
}

func (s *server) GenerateAccessToken(ctx context.Context, request *pb.GenerateAccessTokenRequest) (*pb.StringResponse, error) {
	res, err := s.uc.GenerateAccessToken(unmarshalAccessTokenClaims(request.GetClaims()))

	return marshalBytesToStringResponse(res), err
}

func (s *server) VerifyAccessToken(ctx context.Context, request *pb.VerifyAccessTokenRequest) (*pb.VerifyAccessTokenResponse, error) {
	res, err := s.uc.VerifyAccessToken([]byte(request.GetToken()))
	return &pb.VerifyAccessTokenResponse{
		Claims: marshalAccessTokenClaims(res),
	}, err
}

func (f *server) UploadToAliOss(stream pb.Upload_UploadToAliOssServer) error {
	fileKey, err := f.uc.UploadToAliOss(grpc_streamer.NewStreamReader(grpc_streamer.StreamReaderConfig{
		RecvByte: func() (b byte, err error) {
			chunk, recvErr := stream.Recv()
			if recvErr != nil {
				err = recvErr
				return
			}

			if chunk == nil {
				err = io.EOF
				return
			}

			return byte(chunk.Data), nil
		},
	}))

	if err != nil {
		return err
	}

	return stream.SendAndClose(marshalBytesToStringResponse(fileKey))
}

func (f *server) GetSource(ctx context.Context, request *pb.GetSourceRequest) (*pb.FileSource, error) {
	source, err := f.uc.GetSource(request.GetKey())

	return marshalFileSource(source), err
}

func (f *server) InsertSource(ctx context.Context, request *pb.InsertSourceRequest) (*pb.StringResponse, error) {
	key, err := f.uc.InsertSource(request.GetValue(), request.GetDesc(), request.GetFileMD5(), request.GetFileSize())
	return marshalBytesToStringResponse(key), err
}

func MarshalBranchResponse(b *daemon.Branch) *pb.BranchResponse {
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

func MarshalEnvResponse(e *daemon.Env) *pb.EnvResponse {
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

func MarshalVersionResponse(v *daemon.Version) *pb.VersionResponse {
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

func MarshalVersionList(l daemon.VersionList) *pb.VersionListResponse {
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

	func marshalAccessTokenClaims(claims *daemon.AccessTokenClaims) *pb.AccessTokenClaims {
	if claims == nil {
		return nil
	}

	return &pb.AccessTokenClaims{
		Issuer:   pb.AccessTokenIssuer(claims.Issuer),
		Subject:  claims.Subject,
		Audience: claims.Audience,
	}
}

func marshalBytesToStringResponse(bytes []byte) *pb.StringResponse {
	if bytes == nil {
		return nil
	}

	return &pb.StringResponse{Data: string(bytes)}
}

func unmarshalAccessTokenClaims(request *pb.AccessTokenClaims) *daemon.AccessTokenClaims {
	if request == nil {
		return nil
	}

	return &daemon.AccessTokenClaims{
		Issuer:   daemon.AccessTokenIssuer(request.Issuer),
		Subject:  request.Subject,
		Audience: request.Audience,
	}
}

func marshalFileSource(file *daemon.File) *pb.FileSource {
	if file == nil {
		return nil
	}

	return &pb.FileSource{
		Key:        file.Key,
		Value:      file.Value,
		Desc:       file.Desc,
		CreateTime: file.CreateTime.UnixNano(),
		FileMD5:    file.FileMD5,
		FileSize:   file.FileSize,
	}
}
