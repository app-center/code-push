package usecase

type IVersionUpdateParams interface {
	SetChangelog(changelog string) IVersionUpdateParams
	SetPackageUri(packageUri string) IVersionUpdateParams
	SetMustUpdate(mustUpdate bool) IVersionUpdateParams

	Changelog() (set bool, val string)
	PackageUri() (set bool, val string)
	MustUpdate() (set, val bool)
}

type versionUpdateParams struct {
	changelogSet bool
	changelog    string

	packageUriSet bool
	packageUri    string

	mustUpdateSet bool
	mustUpdate    bool
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
