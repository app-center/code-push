package usecase

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	sessionAdapter "github.com/funnyecho/code-push/daemon/session/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway/portal"
	"mime/multipart"
)

type UseCase interface {
	Auth
	Env
	Version
	Uploader
}

type Auth interface {
	Auth(ctx context.Context, branchId, timestamp, nonce, sign []byte) error
	SignToken(ctx context.Context, branchId []byte) ([]byte, error)
	VerifyToken(ctx context.Context, token []byte) (branchId []byte, err error)
}

type Env interface {
	CreateEnv(ctx context.Context, branchId, envName []byte) (*portal.Env, error)
	GetEnv(ctx context.Context, envId []byte) (*portal.Env, error)
	DeleteEnv(ctx context.Context, envId []byte) error
	GetEnvEncToken(ctx context.Context, envId []byte) ([]byte, error)
}

type Version interface {
	ReleaseVersion(ctx context.Context, params *portal.VersionReleaseParams) error
	GetVersion(ctx context.Context, envId, appVersion []byte) (*portal.Version, error)
	VersionStrictCompatQuery(ctx context.Context, envId, appVersion []byte) (*portal.VersionCompatQueryResult, error)
}

type Uploader interface {
	UploadPkg(ctx context.Context, stream multipart.File) (fileKey []byte, err error)
}

type CodePushAdapter interface {
	GetBranchEncToken(ctx context.Context, branchId []byte) ([]byte, error)
	CreateEnv(ctx context.Context, branchId, envName []byte) (*pb.EnvResponse, error)
	GetEnv(ctx context.Context, envId []byte) (*pb.EnvResponse, error)
	DeleteEnv(ctx context.Context, envId []byte) error
	GetEnvEncToken(ctx context.Context, envId []byte) ([]byte, error)
	ReleaseVersion(ctx context.Context, params *pb.VersionReleaseRequest) error
	GetVersion(ctx context.Context, envId, appVersion []byte) (*pb.VersionResponse, error)
	VersionStrictCompatQuery(ctx context.Context, envId, appVersion []byte) (*pb.VersionStrictCompatQueryResponse, error)
}

type SessionAdapter interface {
	GenerateAccessToken(ctx context.Context, issuer sessionAdapter.AccessTokenIssuer, subject string) ([]byte, error)
	VerifyAccessToken(ctx context.Context, token string) (subject []byte, err error)
}

type FilerAdapter interface {
	UploadPkg(ctx context.Context, source multipart.File) (fileKey []byte, err error)
}
