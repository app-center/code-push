package portal

import "time"

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

type VersionReleaseParams struct {
	EnvId            []byte
	AppVersion       []byte
	CompatAppVersion []byte
	Changelog        []byte
	PackageFileKey   []byte
	MustUpdate       bool
}

type VersionCompatQueryResult struct {
	AppVersion          []byte
	LatestAppVersion    []byte
	CanUpdateAppVersion []byte
	MustUpdate          bool
}
