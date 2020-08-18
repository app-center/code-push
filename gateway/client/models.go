package client

import "time"

type VersionCompatQueryResult struct {
	AppVersion          string
	LatestAppVersion    string
	CanUpdateAppVersion string
	MustUpdate          bool
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

type FileSource struct {
	Key        string
	Value      string
	Desc       string
	CreateTime time.Time
	FileMD5    string
	FileSize   int64
}
