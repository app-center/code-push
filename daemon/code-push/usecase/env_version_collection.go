package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain"
	"github.com/funnyecho/code-push/daemon/code-push/usecase/errors"
	"github.com/funnyecho/code-push/pkg/semver"
	"github.com/funnyecho/code-push/pkg/versionCompatTree"
)

type envVersionCollection struct {
	envId string

	versionService domain.IVersionService
	envService     domain.IEnvService

	versionList       domain.VersionList
	versionCompatTree versionCompatTree.ITree
}

func (c *envVersionCollection) ReleaseVersion(params IVersionReleaseParams) error {
	rawAppVersion := params.AppVersion()
	rawCompatAppVersion := params.CompatAppVersion()
	changelog := params.Changelog()
	packageUri := params.PackageUri()
	mustUpdate := params.MustUpdate()

	if len(rawAppVersion) == 0 ||
		len(rawCompatAppVersion) == 0 ||
		len(packageUri) == 0 {
		return errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "invalid release params",
			Params: errors.MetaFields{
				"appVersion":       rawAppVersion,
				"compatAppVersion": rawCompatAppVersion,
				"packageUri":       packageUri,
			},
		})
	}

	appVersion, appVersionErr := semver.ParseVersion(rawAppVersion)
	if appVersionErr != nil {
		return errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Err: appVersionErr,
			Msg: "invalid app version",
			Params: errors.MetaFields{
				"appVersion": rawAppVersion,
			},
		})
	}

	compatAppVersion, compatVersionErr := semver.ParseVersion(rawCompatAppVersion)
	if compatVersionErr != nil {
		return errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Err: appVersionErr,
			Msg: "invalid compat app version",
			Params: errors.MetaFields{
				"compatAppVersion": compatAppVersion,
			},
		})
	}

	newVersion := &domain.Version{
		EnvId:            c.envId,
		AppVersion:       appVersion.String(),
		CompatAppVersion: compatAppVersion.String(),
		MustUpdate:       mustUpdate,
		Changelog:        changelog,
		PackageUri:       packageUri,
	}

	releaseErr := c.versionService.CreateVersion(newVersion)
	if releaseErr != nil {
		return errors.ThrowVersionReleaseFailedError(releaseErr, errors.FA_VERSION_RELEASE_FAILED, BoxingVersionReleaseParams(params))
	}

	c.versionCompatTree.Publish(newVersionCompatEntry(newVersion))

	return nil
}

func (c *envVersionCollection) GetVersion(rawAppVersion string) (*Version, error) {
	appVersion, appVersionErr := semver.ParseVersion(rawAppVersion)
	if appVersionErr != nil {
		return nil, errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Err:    appVersionErr,
			Msg:    "invalid appVersion",
			Params: errors.MetaFields{"envId": c.envId, "appVersion": rawAppVersion},
		})
	}

	version, versionErr := c.versionService.Version(c.envId, appVersion.String())
	if versionErr != nil {
		return nil, errors.ThrowVersionNotFoundError(c.envId, appVersion.String(), versionErr)
	}

	return toVersion(version), nil
}

func (c *envVersionCollection) ListVersions() (VersionList, error) {
	versionList, versionListErr := c.versionService.VersionsWithEnvId(c.envId)
	if versionListErr != nil {
		return nil, errors.ThrowVersionOperationForbiddenError(versionListErr, "fetch version list failed", errors.MetaFields{"envId": c.envId})
	}

	versionListOutput := make(VersionList, len(versionList))
	for i, v := range versionList {
		versionListOutput[i] = toVersion(v)
	}

	return versionListOutput, nil
}

func (c *envVersionCollection) VersionStrictCompatQuery(rawAppVersion string) (IVersionCompatQueryResult, error) {
	appVersion, appVersionErr := semver.ParseVersion(rawAppVersion)
	if appVersionErr != nil {
		return nil, errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Err:    appVersionErr,
			Msg:    "invalid appVersion",
			Params: errors.MetaFields{"envId": c.envId, "appVersion": rawAppVersion},
		})
	}

	r := c.versionCompatTree.StrictCompat(newVersionCompatQueryAnchor(appVersion))

	latestAppVersionEntry := r.LatestVersion()
	canUpdateAppVersionEntry := r.CanUpdateVersion()

	var latestAppVersion, canUpdateAppVersion *semver.SemVer
	var mustUpdate bool

	if latestAppVersionEntry != nil {
		latestAppVersion = latestAppVersionEntry.Version()
	}

	if canUpdateAppVersionEntry != nil {
		canUpdateAppVersion = canUpdateAppVersionEntry.Version()
		canUpdateAppVersionModel, canUpdateAppVersionModelErr := c.GetVersion(canUpdateAppVersion.String())
		if canUpdateAppVersionModelErr != nil {
			return nil, errors.ThrowVersionOperationForbiddenError(
				canUpdateAppVersionModelErr,
				"failed to get canUpdateAppVersion scheme",
				errors.MetaFields{"canUpdateAppVersion": canUpdateAppVersion.String()},
			)
		}

		mustUpdate = canUpdateAppVersionModel.MustUpdate
	}

	queryResult := NewVersionCompatQueryResult(VersionCompatQueryResultConfig{
		AppVersion:          appVersion,
		LatestAppVersion:    latestAppVersion,
		CanUpdateAppVersion: canUpdateAppVersion,
		MustUpdate:          mustUpdate,
	})

	return queryResult, nil
}

func (c *envVersionCollection) init() error {
	versionList, fetchErr := c.versionService.VersionsWithEnvId(c.envId)

	if fetchErr != nil {
		return fetchErr
	}

	c.versionList = versionList

	treeEntries := make([]versionCompatTree.IEntry, len(versionList))
	for i, version := range versionList {
		treeEntries[i] = newVersionCompatEntry(version)
	}

	c.versionCompatTree.Publish(treeEntries...)

	return nil
}

func (c *envVersionCollection) resetSource() error {
	panic("implement me")
}

type envVersionCollectionConfig struct {
	EnvId          string
	EnvService     domain.IEnvService
	VersionService domain.IVersionService
}

func newEnvVersionCollection(config envVersionCollectionConfig) (*envVersionCollection, error) {
	if config.VersionService == nil ||
		config.EnvService == nil ||
		len(config.EnvId) == 0 {
		return nil, errors.ThrowVersionOperationForbiddenError(
			nil,
			"invalid version collection params",
			nil,
		)
	}

	collection := &envVersionCollection{
		envId:             config.EnvId,
		versionService:    config.VersionService,
		envService:        config.EnvService,
		versionCompatTree: versionCompatTree.NewVersionCompatTree(),
	}

	if initErr := collection.init(); initErr != nil {
		return nil, initErr
	} else {
		return collection, nil
	}
}

type versionCompatEntry struct {
	version *domain.Version
}

func newVersionCompatEntry(version *domain.Version) *versionCompatEntry {
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
