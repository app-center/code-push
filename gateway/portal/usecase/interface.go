package usecase

import (
	"github.com/funnyecho/code-push/gateway/portal"
	"io"
)

type UseCase interface {
	Auth
	Env
	Version
	Uploader
}

type Auth interface {
	Auth(branchId, timestamp, nonce, sign []byte) error
	SignToken(branchId []byte) ([]byte, error)
	VerifyToken(token []byte) (branchId []byte, err error)
}

type Env interface {
	CreateEnv(branchId, envName []byte) (*portal.Env, error)
	GetEnv(envId []byte) (*portal.Env, error)
	DeleteEnv(envId []byte) error
	GetEnvEncToken(envId []byte) ([]byte, error)
}

type Version interface {
	ReleaseVersion(params *portal.VersionReleaseParams) error
	GetVersion(envId, appVersion []byte) (*portal.Version, error)
	VersionStrictCompatQuery(envId, appVersion []byte) (*portal.VersionCompatQueryResult, error)
}

type Uploader interface {
	UploadPkg(stream io.Reader) (fileKey []byte, err error)
}

type CodePushAdapter interface {
	GetBranchEncToken(branchId []byte) ([]byte, error)
	CreateEnv(branchId, envName []byte) (*portal.Env, error)
	GetEnv(envId []byte) (*portal.Env, error)
	DeleteEnv(envId []byte) error
	GetEnvEncToken(envId []byte) ([]byte, error)
	ReleaseVersion(params *portal.VersionReleaseParams) error
	GetVersion(envId, appVersion []byte) (*portal.Version, error)
	VersionStrictCompatQuery(envId, appVersion []byte) (*portal.VersionCompatQueryResult, error)
}

type SessionAdapter interface {
	GenerateAccessToken(subject string) ([]byte, error)
	VerifyAccessToken(token string) (subject []byte, err error)
}

type FilerAdapter interface {
	UploadPkg(source io.Reader) (fileKey []byte, err error)
}
