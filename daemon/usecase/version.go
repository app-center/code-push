package usecase

import (
	"github.com/funnyecho/code-push/daemon"
	"github.com/funnyecho/code-push/pkg/cache"
	"github.com/funnyecho/code-push/pkg/semver"
	"github.com/pkg/errors"
)

func (uc *useCase) ReleaseVersion(params VersionReleaseParams) error {
	if params == nil {
		return errors.Wrap(daemon.ErrParamsInvalid, "release params is required")
	}

	if params.EnvId() == nil {
		return errors.Wrap(daemon.ErrParamsInvalid, "envId, appVersion, compatAppVersion, packageFileKey are required")
	}

	collection, collectionErr := uc.getEnvVersionCollection(params.EnvId())

	if collectionErr != nil {
		return collectionErr
	}

	return collection.ReleaseVersion(params)
}

func (uc *useCase) GetVersion(envId, appVersion []byte) (*daemon.Version, error) {
	if envId == nil || appVersion == nil {
		return nil, errors.Wrap(daemon.ErrParamsInvalid, "envId and appVersion are required")
	}

	semAppVersion, semAppVersionErr := semver.ParseVersion(string(appVersion))
	if semAppVersionErr != nil {
		return nil, errors.WithMessagef(daemon.ErrParamsInvalid, "failed to parse version, rawAppVersion: %s", appVersion)
	}

	collection, collectionErr := uc.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.GetVersion(semAppVersion)
}

func (uc *useCase) ListVersions(envId []byte) (daemon.VersionList, error) {
	if envId == nil {
		return nil, errors.Wrapf(daemon.ErrParamsInvalid, "envId is required")
	}

	collection, collectionErr := uc.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.ListVersions()
}

func (uc *useCase) VersionStrictCompatQuery(envId, appVersion []byte) (VersionCompatQueryResult, error) {
	if envId == nil || appVersion == nil {
		return nil, errors.Wrap(daemon.ErrParamsInvalid, "envId and appVersion are required")
	}

	semAppVersion, semAppVersionErr := semver.ParseVersion(string(appVersion))
	if semAppVersionErr != nil {
		return nil, errors.Wrapf(daemon.ErrParamsInvalid, "failed to parse version, rawAppVersion: %s", appVersion)
	}

	collection, collectionErr := uc.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.VersionStrictCompatQuery(semAppVersion)
}

func (uc *useCase) initVersionUseCase() {
	uc.envVersionCollectionCache = cache.New(cache.CtorConfig{
		Capacity: 10,
		AllocFunc: func(key cache.KeyType) (collection cache.ValueType, ok bool) {
			envId, isEnvIdType := key.(string)

			if !isEnvIdType {
				return nil, false
			}

			if env, envErr := uc.domain.Env([]byte(envId)); envErr != nil || env == nil {
				return nil, false
			}

			collection, collectionErr := NewEnvVersionCollection(EnvVersionCollectionConfig{
				EnvId:         []byte(envId),
				DomainAdapter: uc.domain,
			})

			ok = collectionErr == nil

			return
		},
	})
}

func (uc *useCase) getEnvVersionCollection(envId []byte) (*EnvVersionCollection, error) {
	if env, envErr := uc.domain.Env(envId); envErr != nil || env == nil {
		return nil, errors.Wrapf(daemon.ErrEnvNotFound, "envId: %s", envId)
	}

	collection, hasCollection := uc.envVersionCollectionCache.Get(string(envId))

	if hasCollection {
		return collection.(*EnvVersionCollection), nil
	} else {
		return nil, errors.Wrapf(daemon.ErrEnvNotFound, "can not find version collection, envId: %s", envId)
	}
}
