package usecase

import "github.com/funnyecho/code-push/pkg/semver"

type IVersionCompatQueryResult interface {
	AppVersion() string
	LatestAppVersion() string
	CanUpdateAppVersion() string
	MustUpdate() bool
}

type versionCompatQueryResult struct {
	appVersion          string
	latestAppVersion    string
	canUpdateAppVersion string
	mustUpdate          bool
}

func (v *versionCompatQueryResult) AppVersion() string {
	return v.appVersion
}

func (v *versionCompatQueryResult) LatestAppVersion() string {
	return v.latestAppVersion
}

func (v *versionCompatQueryResult) CanUpdateAppVersion() string {
	return v.canUpdateAppVersion
}

func (v *versionCompatQueryResult) MustUpdate() bool {
	return v.mustUpdate
}

type VersionCompatQueryResultConfig struct {
	AppVersion          *semver.SemVer
	LatestAppVersion    *semver.SemVer
	CanUpdateAppVersion *semver.SemVer
	MustUpdate          bool
}

func NewVersionCompatQueryResult(config VersionCompatQueryResultConfig) IVersionCompatQueryResult {
	var appVersion, latestAppVersion, canUpdateAppVersion string

	if config.AppVersion != nil {
		appVersion = config.AppVersion.String()
	}

	if config.LatestAppVersion != nil {
		latestAppVersion = config.LatestAppVersion.String()
	}

	if config.CanUpdateAppVersion != nil {
		canUpdateAppVersion = config.CanUpdateAppVersion.String()
	}

	return &versionCompatQueryResult{
		appVersion:          appVersion,
		latestAppVersion:    latestAppVersion,
		canUpdateAppVersion: canUpdateAppVersion,
		mustUpdate:          config.MustUpdate,
	}
}
