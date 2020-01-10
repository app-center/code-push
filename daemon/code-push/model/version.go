package model

import "time"

type Version struct {
	EnvId			 string	   `json:"env_id"`
	AppVersion       string    `json:"app_version"`
	CompatAppVersion string    `json:"compat_app_version"`
	MustUpdate       bool      `json:"must_update"`
	Changelog        string    `json:"changelog"`
	PackageUri       string    `json:"package_uri"`
	PackageBlob      string    `json:"package_blob"`
	CreateTime       time.Time `json:"create_time"`
}

func NewVersion(
	appVersion, compatAppVersion string,
	mustUpdate bool,
	changeLog, packageUri, packageBlob string,
	createTime time.Time,
) *Version {
	return &Version{
		AppVersion:       appVersion,
		CompatAppVersion: compatAppVersion,
		MustUpdate:       mustUpdate,
		Changelog:        changeLog,
		PackageUri:       packageUri,
		PackageBlob:      packageBlob,
		CreateTime:       createTime,
	}
}

type VersionMap = map[string]*Version
type VersionList = []*Version