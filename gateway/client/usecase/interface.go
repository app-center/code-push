package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	sessionAdapter "github.com/funnyecho/code-push/daemon/session/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway/client"
)

type UseCase interface {
	Auth
	Version
	Metrics
}

type Auth interface {
	Auth(envId, timestamp, nonce, sign []byte) error
	SignToken(envId []byte) ([]byte, error)
	VerifyToken(token []byte) (envId []byte, err error)
}

type Version interface {
	GetVersion(envId, appVersion []byte) (*client.Version, error)
	VersionDownloadPkg(envId, appVersion []byte) ([]byte, error)
	VersionStrictCompatQuery(envId, appVersion []byte) (*client.VersionCompatQueryResult, error)
}

type Metrics interface {
	RequestDuration(path string, success bool, durationSecond float64)
}

type CodePushAdapter interface {
	GetEnvEncToken(envId []byte) ([]byte, error)
	GetVersion(envId, appVersion []byte) (*pb.VersionResponse, error)
	VersionStrictCompatQuery(envId, appVersion []byte) (*pb.VersionStrictCompatQueryResponse, error)
}

type SessionAdapter interface {
	GenerateAccessToken(issuer sessionAdapter.AccessTokenIssuer, subject string) ([]byte, error)
	VerifyAccessToken(token string) (subject []byte, err error)
}

type FilerAdapter interface {
	GetSource(fileKey []byte) ([]byte, error)
}

type MetricsAdapter interface {
	HttpRequestDuration(svr, path string, success bool, durationSecond float64)
}
