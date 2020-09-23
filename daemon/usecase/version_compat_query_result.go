package usecase

import "github.com/funnyecho/code-push/pkg/semver"

type VersionCompatQueryResultConfig struct {
	AppVersion          *semver.SemVer
	LatestAppVersion    *semver.SemVer
	CanUpdateAppVersion *semver.SemVer
	MustUpdate          bool
}

func NewVersionCompatQueryResult(config VersionCompatQueryResultConfig) *versionCompatQueryResult {
	var appVersion, latestAppVersion, canUpdateAppVersion []byte

	if config.AppVersion != nil {
		appVersion = []byte(config.AppVersion.String())
	}

	if config.LatestAppVersion != nil {
		latestAppVersion = []byte(config.LatestAppVersion.String())
	}

	if config.CanUpdateAppVersion != nil {
		canUpdateAppVersion = []byte(config.CanUpdateAppVersion.String())
	}

	return &versionCompatQueryResult{
		appVersion:          appVersion,
		latestAppVersion:    latestAppVersion,
		canUpdateAppVersion: canUpdateAppVersion,
		mustUpdate:          config.MustUpdate,
	}
}

type versionCompatQueryResult struct {
	appVersion          []byte
	latestAppVersion    []byte
	canUpdateAppVersion []byte
	mustUpdate          bool
}

func (v *versionCompatQueryResult) AppVersion() []byte {
	return v.appVersion
}

func (v *versionCompatQueryResult) LatestAppVersion() []byte {
	return v.latestAppVersion
}

func (v *versionCompatQueryResult) CanUpdateAppVersion() []byte {
	return v.canUpdateAppVersion
}

func (v *versionCompatQueryResult) MustUpdate() bool {
	return v.mustUpdate
}
