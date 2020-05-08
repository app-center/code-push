package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/usecase/errors"
	"github.com/funnyecho/code-push/pkg/semver"
)

type IVersionReleaseParams interface {
	EnvId() string
	AppVersion() string
	CompatAppVersion() string
	Changelog() string
	PackageUri() string
	MustUpdate() bool
}

type versionReleaseParams struct {
	envId            string
	appVersion       string
	compatAppVersion string
	changelog        string
	packageUri       string
	mustUpdate       bool
}

func (v versionReleaseParams) EnvId() string {
	return v.envId
}

func (v versionReleaseParams) AppVersion() string {
	return v.appVersion
}

func (v versionReleaseParams) AppSemVersion() (*semver.SemVer, error) {
	return semver.ParseVersion(v.appVersion)
}

func (v versionReleaseParams) CompatAppVersion() string {
	return v.compatAppVersion
}

func (v versionReleaseParams) CompatAppSemVersion() (*semver.SemVer, error) {
	return semver.ParseVersion(v.compatAppVersion)
}

func (v versionReleaseParams) Changelog() string {
	return v.changelog
}

func (v versionReleaseParams) PackageUri() string {
	return v.packageUri
}

func (v versionReleaseParams) MustUpdate() bool {
	return v.mustUpdate
}

type VersionReleaseParamsConfig struct {
	EnvId            string
	AppVersion       string
	CompatAppVersion string
	Changelog        string
	PackageUri       string
	MustUpdate       bool
}

func NewVersionReleaseParams(config VersionReleaseParamsConfig) (IVersionReleaseParams, error) {
	hasMissField := false
	missFields := errors.MetaFields{}

	if len(config.EnvId) == 0 {
		hasMissField = true
		missFields["envId"] = config.EnvId
	}

	if len(config.AppVersion) == 0 {
		hasMissField = true
		missFields["appVersion"] = config.AppVersion
	}

	if len(config.Changelog) == 0 {
		hasMissField = true
		missFields["changelog"] = config.Changelog
	}

	if len(config.PackageUri) == 0 {
		hasMissField = true
		missFields["packageUri"] = config.PackageUri
	}

	if hasMissField {
		return nil, errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg:    "required release params was omitted",
			Params: missFields,
		})
	}

	if len(config.CompatAppVersion) == 0 {
		config.CompatAppVersion = config.AppVersion
	}

	return &versionReleaseParams{
		envId:            config.EnvId,
		appVersion:       config.AppVersion,
		compatAppVersion: config.CompatAppVersion,
		changelog:        config.Changelog,
		packageUri:       config.PackageUri,
		mustUpdate:       config.MustUpdate,
	}, nil
}

func BoxingVersionReleaseParams(params IVersionReleaseParams) map[string]interface{} {
	return map[string]interface{}{
		"envId":            params.EnvId(),
		"appVersion":       params.AppVersion(),
		"compatAppVersion": params.CompatAppVersion(),
		"changelog":        params.Changelog(),
		"packageUri":       params.PackageUri(),
		"mustUpdate":       params.MustUpdate(),
	}
}
