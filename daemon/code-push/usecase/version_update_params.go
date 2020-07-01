package usecase

type IVersionUpdateParams interface {
	SetChangelog(changelog string) IVersionUpdateParams
	SetPackageFileKey(packageFileKey string) IVersionUpdateParams
	SetMustUpdate(mustUpdate bool) IVersionUpdateParams

	Changelog() (set bool, val string)
	PackageFileKey() (set bool, val string)
	MustUpdate() (set, val bool)
}

type versionUpdateParams struct {
	changelogSet bool
	changelog    string

	packageFileKeySet bool
	packageFileKey    string

	mustUpdateSet bool
	mustUpdate    bool
}

func (v *versionUpdateParams) SetChangelog(changelog string) IVersionUpdateParams {
	v.changelogSet = true
	v.changelog = changelog

	return v
}

func (v *versionUpdateParams) SetPackageFileKey(packageFileKey string) IVersionUpdateParams {
	v.packageFileKeySet = true
	v.packageFileKey = packageFileKey

	return v
}

func (v *versionUpdateParams) SetMustUpdate(mustUpdate bool) IVersionUpdateParams {
	v.mustUpdateSet = true
	v.mustUpdate = mustUpdate

	return v
}

func (v *versionUpdateParams) Changelog() (set bool, val string) {
	return v.changelogSet, v.changelog
}

func (v *versionUpdateParams) PackageFileKey() (set bool, val string) {
	return v.packageFileKeySet, v.packageFileKey
}

func (v *versionUpdateParams) MustUpdate() (set, val bool) {
	return v.mustUpdateSet, v.mustUpdate
}

func NewVersionUpdateParams() IVersionUpdateParams {
	return &versionUpdateParams{}
}
