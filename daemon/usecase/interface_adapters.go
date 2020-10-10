package usecase

import (
	"github.com/funnyecho/code-push/daemon"
	"io"
)

type DomainAdapter interface {
	Branch(branchId []byte) (*daemon.Branch, error)
	CreateBranch(branch *daemon.Branch) error
	DeleteBranch(branchId []byte) error
	IsBranchAvailable(branchId []byte) bool
	IsBranchNameExisted(branchName []byte) (bool, error)

	Env(envId []byte) (*daemon.Env, error)
	GetEnvsWithBranchId(branchId string) ([]*daemon.Env, error)
	CreateEnv(env *daemon.Env) error
	DeleteEnv(envId []byte) error
	IsEnvAvailable(envId []byte) bool
	IsEnvNameExisted(branchId, envName []byte) (bool, error)

	Version(envId, appVersion []byte) (*daemon.Version, error)
	VersionsWithEnvId(envId []byte) (daemon.VersionList, error)
	CreateVersion(version *daemon.Version) error
	IsVersionAvailable(envId, appVersion []byte) (bool, error)

	File(fileKey string) (*daemon.File, error)
	InsertFile(file *daemon.File) error
	IsFileKeyExisted(fileKey string) bool
}

type AliOssAdapter interface {
	SignFetchURL(key []byte) ([]byte, error)
	Upload(stream io.Reader) ([]byte, error)
}
