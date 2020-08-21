package usecase

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	filerpb "github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	sessionAdapter "github.com/funnyecho/code-push/daemon/session/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway/client"
)

type UseCase interface {
	Auth
	Version
	Filer
}

type Auth interface {
	Auth(ctx context.Context, envId, timestamp, nonce, sign []byte) error
	SignToken(ctx context.Context, envId []byte) ([]byte, error)
	VerifyToken(ctx context.Context, token []byte) (envId []byte, err error)
}

type Version interface {
	GetVersion(ctx context.Context, envId, appVersion []byte) (*client.Version, error)
	VersionPkgSource(ctx context.Context, envId, appVersion string) (*client.FileSource, error)
	VersionStrictCompatQuery(ctx context.Context, envId, appVersion []byte) (*client.VersionCompatQueryResult, error)
}

type Filer interface {
	FileDownload(ctx context.Context, fileId []byte) ([]byte, error)
}

type CodePushAdapter interface {
	GetEnvEncToken(ctx context.Context, envId []byte) ([]byte, error)
	GetVersion(ctx context.Context, envId, appVersion []byte) (*pb.VersionResponse, error)
	VersionStrictCompatQuery(ctx context.Context, envId, appVersion []byte) (*pb.VersionStrictCompatQueryResponse, error)
}

type SessionAdapter interface {
	GenerateAccessToken(ctx context.Context, issuer sessionAdapter.AccessTokenIssuer, subject string) ([]byte, error)
	VerifyAccessToken(ctx context.Context, token string) (subject []byte, err error)
}

type FilerAdapter interface {
	GetSource(ctx context.Context, fileKey []byte) (*filerpb.FileSource, error)
}
