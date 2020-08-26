package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/pkg/semver"
	"github.com/funnyecho/code-push/pkg/version-compat-tree"
	"github.com/pkg/errors"
)

func NewEnvVersionCollection(config EnvVersionCollectionConfig) (*EnvVersionCollection, error) {
	if config.DomainAdapter == nil ||
		config.EnvId == nil {
		return nil, errors.Wrap(code_push.ErrParamsInvalid, "invalid version collection params")
	}

	collection := &EnvVersionCollection{
		envId:             string(config.EnvId),
		domain:            config.DomainAdapter,
		versionCompatTree: version_compat_tree.NewVersionCompatTree(),
	}

	if initErr := collection.init(); initErr != nil {
		return nil, initErr
	} else {
		return collection, nil
	}
}

type EnvVersionCollectionConfig struct {
	EnvId         []byte
	DomainAdapter DomainAdapter
}

type EnvVersionCollection struct {
	envId             string
	domain            DomainAdapter
	versionList       code_push.VersionList
	versionCompatTree version_compat_tree.ITree
}

func (e *EnvVersionCollection) ReleaseVersion(params VersionReleaseParams) error {
	rawAppVersion := params.AppVersion()
	rawCompatAppVersion := params.CompatAppVersion()
	changelog := params.Changelog()
	packageFileKey := params.PackageFileKey()
	mustUpdate := params.MustUpdate()

	if rawAppVersion == nil ||
		packageFileKey == nil {
		return errors.Wrapf(code_push.ErrParamsInvalid, "appVersion, compatAppVersion nor packageFileKey can't be empty")
	}

	if rawCompatAppVersion == nil {
		rawCompatAppVersion = make([]byte, len(rawAppVersion))
		copy(rawCompatAppVersion, rawAppVersion)
	}

	appVersion, appVersionErr := semver.ParseVersion(string(rawAppVersion))
	if appVersionErr != nil {
		return errors.Wrapf(code_push.ErrParamsInvalid, "parse appVersion failed")
	}

	compatAppVersion, compatVersionErr := semver.ParseVersion(string(rawCompatAppVersion))
	if compatVersionErr != nil {
		return errors.Wrapf(code_push.ErrParamsInvalid, "parse compatAppVersion failed")
	}

	if compatAppVersion.StageSafetyLooseCompare(appVersion) == semver.CompareLargeFlag {
		return errors.Wrapf(code_push.ErrParamsInvalid, "compatAppVersion shall not larger than appVersion")
	}

	if versionAvailable, versionAvailableErr := e.domain.IsVersionAvailable([]byte(e.envId), []byte(appVersion.String())); versionAvailableErr != nil {
		return errors.Wrapf(versionAvailableErr, "failed to check app version is available")
	} else if versionAvailable {
		return errors.Wrapf(code_push.ErrVersionExisted, "envId:%s, appVersion:%s", e.envId, appVersion.String())
	}

	newVersion := &code_push.Version{
		EnvId:            e.envId,
		AppVersion:       appVersion.String(),
		CompatAppVersion: compatAppVersion.String(),
		MustUpdate:       mustUpdate,
		Changelog:        string(changelog),
		PackageFileKey:   string(packageFileKey),
	}

	releaseErr := e.domain.CreateVersion(newVersion)
	if releaseErr != nil {
		return errors.WithStack(releaseErr)
	}

	e.versionCompatTree.Publish(newVersionCompatEntry(newVersion))

	return nil
}

func (e *EnvVersionCollection) GetVersion(appVersion *semver.SemVer) (*code_push.Version, error) {
	version, versionErr := e.domain.Version([]byte(e.envId), []byte(appVersion.String()))
	if versionErr != nil {
		return nil, errors.WithStack(versionErr)
	}
	if version == nil {
		return nil, nil
	}

	return version, nil
}

func (e *EnvVersionCollection) ListVersions() (code_push.VersionList, error) {
	versionList, versionListErr := e.domain.VersionsWithEnvId([]byte(e.envId))
	if versionListErr != nil {
		return nil, errors.Wrapf(versionListErr, "fetch version list failed, envId: %s", e.envId)
	}

	if versionList == nil {
		versionList = make(code_push.VersionList, 0)
	}

	return versionList, nil
}

func (e *EnvVersionCollection) VersionStrictCompatQuery(appVersion *semver.SemVer) (VersionCompatQueryResult, error) {
	isAppVersionAvailable, appVersionAvailableErr := e.domain.IsVersionAvailable([]byte(e.envId), []byte(appVersion.String()))
	var resultAppVersion *semver.SemVer
	if appVersionAvailableErr == nil && isAppVersionAvailable {
		resultAppVersion = appVersion
	}

	r := e.versionCompatTree.StrictCompat(newVersionCompatQueryAnchor(appVersion))

	latestAppVersionEntry := r.LatestVersion()
	canUpdateAppVersionEntry := r.CanUpdateVersion()

	var latestAppVersion, canUpdateAppVersion *semver.SemVer
	var mustUpdate bool

	if latestAppVersionEntry != nil {
		latestAppVersion = latestAppVersionEntry.Version()
	}

	if canUpdateAppVersionEntry != nil {
		canUpdateAppVersion = canUpdateAppVersionEntry.Version()
		canUpdateAppVersionModel, canUpdateAppVersionModelErr := e.GetVersion(canUpdateAppVersion)
		if canUpdateAppVersionModelErr != nil {
			return nil, errors.Wrapf(canUpdateAppVersionModelErr, "failed to get canUpdateAppVersion: %s", canUpdateAppVersion.String())
		}

		mustUpdate = canUpdateAppVersionModel.MustUpdate
	}

	queryResult := NewVersionCompatQueryResult(VersionCompatQueryResultConfig{
		AppVersion:          resultAppVersion,
		LatestAppVersion:    latestAppVersion,
		CanUpdateAppVersion: canUpdateAppVersion,
		MustUpdate:          mustUpdate,
	})

	return queryResult, nil
}

func (e *EnvVersionCollection) init() error {
	versionList, fetchErr := e.domain.VersionsWithEnvId([]byte(e.envId))

	if fetchErr != nil {
		return fetchErr
	}

	e.versionList = versionList

	treeEntries := make([]version_compat_tree.IEntry, len(versionList))
	for i, version := range versionList {
		treeEntries[i] = newVersionCompatEntry(version)
	}

	e.versionCompatTree.Publish(treeEntries...)

	return nil
}

func (e *EnvVersionCollection) resetSource() error {
	panic("implement me")
}

type versionCompatEntry struct {
	version *code_push.Version
}

func newVersionCompatEntry(version *code_push.Version) *versionCompatEntry {
	return &versionCompatEntry{
		version: version,
	}
}

func (e *versionCompatEntry) CompatVersion() *semver.SemVer {
	ver, _ := semver.ParseVersion(e.version.CompatAppVersion)

	return ver
}

func (e *versionCompatEntry) Version() *semver.SemVer {
	ver, _ := semver.ParseVersion(e.version.AppVersion)

	return ver
}

type versionCompatQueryAnchor struct {
	appVersion *semver.SemVer
}

func newVersionCompatQueryAnchor(appVersion *semver.SemVer) *versionCompatQueryAnchor {
	return &versionCompatQueryAnchor{
		appVersion: appVersion,
	}
}

func (v *versionCompatQueryAnchor) Version() *semver.SemVer {
	return v.appVersion
}
