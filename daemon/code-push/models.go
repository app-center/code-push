package code_push

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
	PackageFileKey   string
	CreateTime       time.Time
}

type BranchList = []*Branch
type EnvList = []*Env
type VersionList = []*Version
