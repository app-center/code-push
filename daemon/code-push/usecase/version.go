package usecase

import (
	code_push "github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/daemon/code-push/domain"
	"github.com/funnyecho/code-push/pkg/cache"
	"github.com/pkg/errors"
	"time"
)

func NewVersionUseCase(config VersionUseCaseConfig) IVersion {
	if config.VersionService == nil ||
		config.EnvService == nil {
		panic("invalid version use case params")
	}

	userCase := &versionUseCase{
		versionService: config.VersionService,
		envService:     config.EnvService,
	}

	userCase.init()
	return userCase
}

type Version struct {
	EnvId            string
	AppVersion       string
	CompatAppVersion string
	MustUpdate       bool
	Changelog        string
	PackageFileKey   string
	CreateTime       time.Time
}

type VersionList = []*Version

func toVersion(ver *domain.Version) *Version {
	return &Version{
		EnvId:            ver.EnvId,
		AppVersion:       ver.AppVersion,
		CompatAppVersion: ver.CompatAppVersion,
		MustUpdate:       ver.MustUpdate,
		Changelog:        ver.Changelog,
		PackageFileKey:   ver.PackageFileKey,
		CreateTime:       ver.CreateTime,
	}
}

type IVersion interface {
	ReleaseVersion(params IVersionReleaseParams) error
	GetVersion(envId, appVersion string) (*Version, error)
	ListVersions(envId string) (VersionList, error)
	VersionStrictCompatQuery(envId, appVersion string) (IVersionCompatQueryResult, error)
}

type versionUseCase struct {
	envService                domain.IEnvService
	versionService            domain.IVersionService
	envVersionCollectionCache *cache.Cache
}

func (v *versionUseCase) ReleaseVersion(params IVersionReleaseParams) error {
	collection, collectionErr := v.getEnvVersionCollection(params.EnvId())

	if collectionErr != nil {
		return collectionErr
	}

	return collection.ReleaseVersion(params)
}

func (v *versionUseCase) GetVersion(envId, appVersion string) (*Version, error) {
	collection, collectionErr := v.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.GetVersion(appVersion)
}

func (v *versionUseCase) ListVersions(envId string) (VersionList, error) {
	collection, collectionErr := v.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.ListVersions()
}

func (v *versionUseCase) VersionStrictCompatQuery(envId, appVersion string) (IVersionCompatQueryResult, error) {
	collection, collectionErr := v.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.VersionStrictCompatQuery(appVersion)
}

func (v *versionUseCase) init() {
	v.envVersionCollectionCache = cache.New(cache.CtorConfig{
		Capacity: 10,
		AllocFunc: func(key cache.KeyType) (collection cache.ValueType, ok bool) {
			envId, isEnvIdType := key.(string)

			if !isEnvIdType {
				return nil, false
			}

			if env, envErr := v.envService.Env(envId); envErr != nil || env == nil {
				return nil, false
			}

			collection, collectionErr := newEnvVersionCollection(envVersionCollectionConfig{
				EnvId:          envId,
				VersionService: v.versionService,
				EnvService:     v.envService,
			})

			ok = collectionErr == nil

			return
		},
	})
}

func (v *versionUseCase) getEnvVersionCollection(envId string) (*envVersionCollection, error) {
	if env, envErr := v.envService.Env(envId); envErr != nil || env == nil {
		return nil, errors.Wrapf(code_push.ErrEnvNotFound, "envId: %s", envId)
	}

	collection, hasCollection := v.envVersionCollectionCache.Get(envId)

	if hasCollection {
		return collection.(*envVersionCollection), nil
	} else {
		return nil, errors.Wrapf(code_push.ErrEnvNotFound, "can not find version collection, envId: %s", envId)
	}
}

type VersionUseCaseConfig struct {
	VersionService domain.IVersionService
	EnvService     domain.IEnvService
}
