package usecase

import (
	"context"
	"github.com/funnyecho/code-push/gateway"
	"mime/multipart"
)

type UseCase interface {
	Auth
	Branch
	Env
	Version
	Uploader
	Filer
}

type Auth interface {
	AuthRootUser(ctx context.Context, name, pwd string) error
	SignTokenForRootUser(ctx context.Context) ([]byte, error)
	VerifyTokenForRootUser(ctx context.Context, token []byte) error

	AuthBranch(ctx context.Context, branchId, timestamp, nonce, sign []byte) error
	AuthBranchWithJWT(ctx context.Context, token string) (branchId []byte, err error)
	SignTokenForBranch(ctx context.Context, branchId []byte) ([]byte, error)
	VerifyTokenForBranch(ctx context.Context, token []byte) (branchId []byte, err error)

	AuthEnv(ctx context.Context, envId, timestamp, nonce, sign []byte) error
	SignTokenForEnv(ctx context.Context, envId []byte) ([]byte, error)
	VerifyTokenForEnv(ctx context.Context, token []byte) (envId []byte, err error)

	EvictToken(ctx context.Context, token []byte) error
}

type Branch interface {
	CreateBranch(ctx context.Context, branchName []byte) (*gateway.Branch, error)
	DeleteBranch(ctx context.Context, branchId []byte) error
	GetBranch(ctx context.Context, branchId string) (*gateway.Branch, error)
}

type Env interface {
	CreateEnv(ctx context.Context, branchId, envName, envEncToken []byte) (*gateway.Env, error)
	GetEnv(ctx context.Context, envId []byte) (*gateway.Env, error)
	DeleteEnv(ctx context.Context, envId []byte) error
	GetEnvEncToken(ctx context.Context, envId []byte) ([]byte, error)
	GetEnvsWithBranchId(ctx context.Context, branchId string) ([]*gateway.Env, error)
}

type Version interface {
	ReleaseVersion(ctx context.Context, params *gateway.VersionReleaseParams) error
	GetVersion(ctx context.Context, envId, appVersion []byte) (*gateway.Version, error)
	GetVersionList(ctx context.Context, envId []byte) ([]*gateway.Version, error)
	VersionStrictCompatQuery(ctx context.Context, envId, appVersion []byte) (*gateway.VersionCompatQueryResult, error)
	VersionPkgSource(ctx context.Context, envId, appVersion string) (*gateway.FileSource, error)
}

type Filer interface {
	FileDownload(ctx context.Context, fileId []byte) ([]byte, error)
}

type Uploader interface {
	UploadPkg(ctx context.Context, stream multipart.File) (fileKey []byte, err error)
}
