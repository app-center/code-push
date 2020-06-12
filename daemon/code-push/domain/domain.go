package domain

import "time"

type Branch struct {
	ID         string
	Name       string
	AuthHost   string
	EncToken   string
	CreateTime time.Time
}

type Env struct {
	BranchId   string
	ID         string
	Name       string
	EncToken   string
	CreateTime time.Time
}

type Version struct {
	EnvId            string
	AppVersion       string
	CompatAppVersion string
	MustUpdate       bool
	Changelog        string
	PackageUri       string
	CreateTime       time.Time
}

type BranchList = []*Branch
type EnvList = []*Env
type VersionList = []*Version

type IBranchService interface {
	Branch(branchId string) (*Branch, error)
	CreateBranch(branch *Branch) error
	DeleteBranch(branchId string) error

	IsBranchAvailable(branchId string) bool
}

type IEnvService interface {
	Env(envId string) (*Env, error)
	CreateEnv(env *Env) error
	DeleteEnv(envId string) error
	IsEnvAvailable(envId string) bool
}

type IVersionService interface {
	Version(envId, appVersion string) (*Version, error)
	VersionsWithEnvId(envId string) (VersionList, error)
	CreateVersion(version *Version) error
}
