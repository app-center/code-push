package gateway

import "time"

type Branch struct {
	ID         string
	Name       string
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

type FileSource struct {
	Key        string
	Value      string
	Desc       string
	CreateTime time.Time
	FileMD5    string
	FileSize   int64
}
