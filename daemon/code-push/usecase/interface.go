package usecase

import "github.com/funnyecho/code-push/daemon/code-push"

type Branch interface {
	CreateBranch(branchName, branchAuthHost []byte) (*code_push.Branch, error)
	GetBranch(branchId []byte) (*code_push.Branch, error)
	DeleteBranch(branchId []byte) error
	GetBranchEncToken(branchId []byte) ([]byte, error)
}

type Env interface {
	CreateEnv(branchId, envName []byte) (*code_push.Env, error)
	GetEnv(envId []byte) (*code_push.Env, error)
	DeleteEnv(envId []byte) error
	GetEnvEncToken(envId []byte) ([]byte, error)
	GetEnvAuthHost(envId []byte) ([]byte, error)
}

type Version interface {
	ReleaseVersion(params VersionReleaseParams) error
	GetVersion(envId, appVersion []byte) (*code_push.Version, error)
	ListVersions(envId []byte) (code_push.VersionList, error)
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

type DomainAdapter interface {
	Branch(branchId []byte) (*code_push.Branch, error)
	CreateBranch(branch *code_push.Branch) error
	DeleteBranch(branchId []byte) error
	IsBranchAvailable(branchId []byte) bool

	Env(envId []byte) (*code_push.Env, error)
	CreateEnv(env *code_push.Env) error
	DeleteEnv(envId []byte) error
	IsEnvAvailable(envId []byte) bool

	Version(envId, appVersion []byte) (*code_push.Version, error)
	VersionsWithEnvId(envId []byte) (code_push.VersionList, error)
	CreateVersion(version *code_push.Version) error
}
