package usecase

import (
	"github.com/funnyecho/code-push/pkg/semver"
)

type IVersionUpdateParams interface {
	SetCompatAppVersion(compatAppVersion string) IVersionUpdateParams
	SetChangelog(changelog string) IVersionUpdateParams
	SetPackageUri(packageUri string) IVersionUpdateParams
	SetMustUpdate(mustUpdate bool) IVersionUpdateParams

	CompatAppVersion() (set bool, val string)
	CompatAppSemVersion() (set bool, val *semver.SemVer, err error)
	Changelog() (set bool, val string)
	PackageUri() (set bool, val string)
	MustUpdate() (set, val bool)
}

type versionUpdateParams struct {
	compatAppVersionSet bool
	compatAppVersion    string

	changelogSet bool
	changelog    string

	packageUriSet bool
	packageUri    string

	mustUpdateSet bool
	mustUpdate    bool
}

func (v *versionUpdateParams) SetCompatAppVersion(compatAppVersion string) IVersionUpdateParams {
	v.compatAppVersionSet = true
	v.compatAppVersion = compatAppVersion

	return v
}

func (v *versionUpdateParams) SetChangelog(changelog string) IVersionUpdateParams {
	v.changelogSet = true
	v.changelog = changelog

	return v
}

func (v *versionUpdateParams) SetPackageUri(packageUri string) IVersionUpdateParams {
	v.packageUriSet = true
	v.packageUri = packageUri

	return v
}

func (v *versionUpdateParams) SetMustUpdate(mustUpdate bool) IVersionUpdateParams {
	v.mustUpdateSet = true
	v.mustUpdate = mustUpdate

	return v
}

func (v *versionUpdateParams) CompatAppVersion() (set bool, val string) {
	return v.compatAppVersionSet, v.compatAppVersion
}

func (v *versionUpdateParams) CompatAppSemVersion() (set bool, val *semver.SemVer, err error) {
	set = v.compatAppVersionSet

	if set {
		val, err = semver.ParseVersion(v.compatAppVersion)
	}

	return
}

func (v *versionUpdateParams) Changelog() (set bool, val string) {
	return v.changelogSet, v.changelog
}

func (v *versionUpdateParams) PackageUri() (set bool, val string) {
	return v.packageUriSet, v.packageUri
}

func (v *versionUpdateParams) MustUpdate() (set, val bool) {
	return v.mustUpdateSet, v.mustUpdate
}

func NewVersionUpdateParams() IVersionUpdateParams {
	return &versionUpdateParams{}
}
