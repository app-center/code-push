package usecase

import (
	"github.com/funnyecho/code-push/daemon"
	"io"
)

type UseCase interface {
	Branch
	Env
	Version
	File
	Upload
	AccessToken
}

type Branch interface {
	CreateBranch(branchName []byte) (*daemon.Branch, error)
	GetBranch(branchId []byte) (*daemon.Branch, error)
	DeleteBranch(branchId []byte) error
	GetBranchEncToken(branchId []byte) ([]byte, error)
}

type Env interface {
	CreateEnv(branchId, envName []byte) (*daemon.Env, error)
	GetEnv(envId []byte) (*daemon.Env, error)
	DeleteEnv(envId []byte) error
	GetEnvEncToken(envId []byte) ([]byte, error)
}

type Version interface {
	ReleaseVersion(params VersionReleaseParams) error
	GetVersion(envId, appVersion []byte) (*daemon.Version, error)
	ListVersions(envId []byte) (daemon.VersionList, error)
	VersionStrictCompatQuery(envId, appVersion []byte) (VersionCompatQueryResult, error)
}

type VersionReleaseParams interface {
	EnvId() []byte
	AppVersion() []byte
	CompatAppVersion() []byte
	Changelog() []byte
	PackageFileKey() []byte
	MustUpdate() bool
}

type VersionCompatQueryResult interface {
	AppVersion() []byte
	LatestAppVersion() []byte
	CanUpdateAppVersion() []byte
	MustUpdate() bool
}

type File interface {
	GetSource(key string) (*daemon.File, error)
	InsertSource(value, desc, fileMD5 string, fileSize int64) (daemon.FileKey, error)
}

type Upload interface {
	UploadToAliOss(stream io.Reader) (daemon.FileKey, error)
}

type AccessToken interface {
	GenerateAccessToken(claims *daemon.AccessTokenClaims) ([]byte, error)
	VerifyAccessToken(token []byte) (*daemon.AccessTokenClaims, error)
}
