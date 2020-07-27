package client

import "time"

type VersionCompatQueryResult struct {
	AppVersion          []byte
	LatestAppVersion    []byte
	CanUpdateAppVersion []byte
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
