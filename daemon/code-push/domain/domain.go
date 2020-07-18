package domain

import "github.com/funnyecho/code-push/daemon/code-push"

type Service struct {
	BranchService
	EnvService
	VersionService
}

type BranchService interface {
	Branch(branchId []byte) (*code_push.Branch, error)
	CreateBranch(branch *code_push.Branch) error
	DeleteBranch(branchId []byte) error

	IsBranchAvailable(branchId []byte) bool
}

type EnvService interface {
	Env(envId []byte) (*code_push.Env, error)
	CreateEnv(env *code_push.Env) error
	DeleteEnv(envId []byte) error
	IsEnvAvailable(envId []byte) bool
}

type VersionService interface {
	Version(envId, appVersion []byte) (*code_push.Version, error)
	VersionsWithEnvId(envId []byte) (code_push.VersionList, error)
	CreateVersion(version *code_push.Version) error
}
