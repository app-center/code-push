package model

import "time"

type Version struct {
	envId            string    `json:"env_id"`
	appVersion       string    `json:"app_version"`
	compatAppVersion string    `json:"compat_app_version"`
	mustUpdate       bool      `json:"must_update"`
	changelog        string    `json:"changelog"`
	packageUri       string    `json:"package_uri"`
	createTime       time.Time `json:"create_time"`
}

func (v Version) EnvId() string {
	return v.envId
}

func (v Version) AppVersion() string {
	return v.appVersion
}

func (v Version) CompatAppVersion() string {
	return v.compatAppVersion
}

func (v Version) MustUpdate() bool {
	return v.mustUpdate
}

func (v Version) Changelog() string {
	return v.changelog
}

func (v Version) PackageUri() string {
	return v.packageUri
}

func (v Version) CreateTime() time.Time {
	return v.createTime
}

type VersionConfig struct {
	EnvId            string
	AppVersion       string
	CompatAppVersion string
	MustUpdate       bool
	Changelog        string
	PackageUri       string
	PackageBlob      string
	CreateTime       time.Time
}

func NewVersion(config VersionConfig) *Version {
	return &Version{
		envId:            config.EnvId,
		appVersion:       config.AppVersion,
		compatAppVersion: config.CompatAppVersion,
		mustUpdate:       config.MustUpdate,
		changelog:        config.Changelog,
		packageUri:       config.PackageUri,
		createTime:       config.CreateTime,
	}
}

type VersionMap = map[string]*Version
type VersionList = []*Version
