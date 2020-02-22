package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/model"
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
	"github.com/funnyecho/code-push/daemon/code-push/domain/service"
	"github.com/funnyecho/code-push/daemon/code-push/usecase/errors"
	"github.com/funnyecho/code-push/pkg/cache"
	"time"
)

type Version struct {
	EnvId            string
	AppVersion       string
	CompatAppVersion string
	MustUpdate       bool
	Changelog        string
	PackageUri       string
	CreateTime       time.Time
}

type VersionList = []*Version

func toVersion(ver *model.Version) *Version {
	return &Version{
		EnvId:            ver.EnvId(),
		AppVersion:       ver.AppVersion(),
		CompatAppVersion: ver.CompatAppVersion(),
		MustUpdate:       ver.MustUpdate(),
		Changelog:        ver.Changelog(),
		PackageUri:       ver.PackageUri(),
		CreateTime:       ver.CreateTime(),
	}
}

type IVersionUseCase interface {
	ReleaseVersion(params IVersionReleaseParams) error
	RollbackVersion(envId, appVersion string) (*Version, error)
	UpdateVersion(envId, appVersion string, params IVersionUpdateParams) error
	GetVersion(envId, appVersion string) (*Version, error)
	ListVersions(envId string) (VersionList, error)
	VersionCompatQuery(envId, appVersion string) (IVersionCompatQueryResult, error)
}

type versionUseCase struct {
	versionRepo    repository.IVersion
	versionService service.IVersionService

	envRepo    repository.IEnv
	envService service.IEnvService

	envVersionCollectionCache *cache.Cache
}

func (v *versionUseCase) ReleaseVersion(params IVersionReleaseParams) error {
	collection, collectionErr := v.getEnvVersionCollection(params.EnvId())

	if collectionErr != nil {
		return collectionErr
	}

	return collection.ReleaseVersion(params)
}

func (v *versionUseCase) RollbackVersion(envId, appVersion string) (*Version, error) {
	collection, collectionErr := v.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.RollbackVersion(appVersion)
}

func (v *versionUseCase) UpdateVersion(envId, appVersion string, params IVersionUpdateParams) error {
	collection, collectionErr := v.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return collectionErr
	}

	return collection.UpdateVersion(appVersion, params)
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

func (v *versionUseCase) VersionCompatQuery(envId, appVersion string) (IVersionCompatQueryResult, error) {
	collection, collectionErr := v.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.VersionCompatQuery(appVersion)
}

func (v *versionUseCase) init() error {
	v.envVersionCollectionCache = cache.New(cache.CtorConfig{
		Capacity: 10,
		AllocFunc: func(key cache.KeyType) (collection cache.ValueType, ok bool) {
			envId, isEnvIdType := key.(string)

			if !isEnvIdType {
				return nil, false
			}

			if !v.envService.IsEnvExisted(envId) {
				return nil, false
			}

			collection, collectionErr := newEnvVersionCollection(envVersionCollectionConfig{
				EnvId:          envId,
				VersionRepo:    v.versionRepo,
				VersionService: v.versionService,
				EnvRepo:        v.envRepo,
				EnvService:     v.envService,
			})

			ok = collectionErr == nil

			return
		},
	})

	return nil
}

func (v *versionUseCase) getEnvVersionCollection(envId string) (*envVersionCollection, error) {
	if !v.envService.IsEnvExisted(envId) {
		return nil, errors.ThrowEnvNotFoundError(envId, nil)
	}

	collection, hasCollection := v.envVersionCollectionCache.Get(envId)

	if hasCollection {
		return collection.(*envVersionCollection), nil
	} else {
		return nil, errors.ThrowVersionOperationForbiddenError(
			nil,
			"can not find version collection",
			errors.MetaFields{"envId": envId},
		)
	}
}

type VersionUseCaseConfig struct {
	VersionRepo    repository.IVersion
	VersionService service.IVersionService
	EnvRepo        repository.IEnv
	EnvService     service.IEnvService
}

func NewVersionUseCase(config VersionUseCaseConfig) (IVersionUseCase, error) {
	if config.VersionRepo == nil ||
		config.VersionService == nil ||
		config.EnvRepo == nil ||
		config.EnvService == nil {
		return nil, errors.ThrowVersionOperationForbiddenError(
			nil,
			"invalid version use case params",
			nil,
		)
	}

	userCase := &versionUseCase{
		versionRepo:    config.VersionRepo,
		versionService: config.VersionService,

		envRepo:    config.EnvRepo,
		envService: config.EnvService,
	}

	if initErr := userCase.init(); initErr != nil {
		return nil, initErr
	} else {
		return userCase, nil
	}
}
