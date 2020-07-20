package usecase

import (
	code_push "github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/daemon/code-push/domain"
	"github.com/funnyecho/code-push/pkg/cache"
	"github.com/funnyecho/code-push/pkg/semver"
	"github.com/pkg/errors"
)

func (c *useCase) ReleaseVersion(params VersionReleaseParams) error {
	if params == nil {
		return errors.Wrap(code_push.ErrParamsInvalid, "release params is required")
	}

	if params.EnvId() == nil {
		return errors.Wrap(code_push.ErrParamsInvalid, "envId, appVersion, compatAppVersion, packageFileKey are required")
	}

	collection, collectionErr := c.getEnvVersionCollection(params.EnvId())

	if collectionErr != nil {
		return collectionErr
	}

	return collection.ReleaseVersion(params)
}

func (c *useCase) GetVersion(envId, appVersion []byte) (*code_push.Version, error) {
	if envId == nil || appVersion == nil {
		return nil, errors.Wrap(code_push.ErrParamsInvalid, "envId and appVersion are required")
	}

	semAppVersion, semAppVersionErr := semver.ParseVersion(string(appVersion))
	if semAppVersionErr != nil {
		return nil, errors.WithMessagef(code_push.ErrParamsInvalid, "failed to parse version, rawAppVersion: %s", appVersion)
	}

	collection, collectionErr := c.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.GetVersion(semAppVersion)
}

func (c *useCase) ListVersions(envId []byte) (code_push.VersionList, error) {
	if envId == nil {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "envId is required")
	}

	collection, collectionErr := c.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.ListVersions()
}

func (c *useCase) VersionStrictCompatQuery(envId, appVersion []byte) (VersionCompatQueryResult, error) {
	if envId == nil || appVersion == nil {
		return nil, errors.Wrap(code_push.ErrParamsInvalid, "envId and appVersion are required")
	}

	semAppVersion, semAppVersionErr := semver.ParseVersion(string(appVersion))
	if semAppVersionErr != nil {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "failed to parse version, rawAppVersion: %s", appVersion)
	}

	collection, collectionErr := c.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.VersionStrictCompatQuery(semAppVersion)
}

func (c *useCase) initVersionUseCase() {
	c.envVersionCollectionCache = cache.New(cache.CtorConfig{
		Capacity: 10,
		AllocFunc: func(key cache.KeyType) (collection cache.ValueType, ok bool) {
			envId, isEnvIdType := key.(string)

			if !isEnvIdType {
				return nil, false
			}

			if env, envErr := c.domain.Env([]byte(envId)); envErr != nil || env == nil {
				return nil, false
			}

			collection, collectionErr := NewEnvVersionCollection(EnvVersionCollectionConfig{
				EnvId:         []byte(envId),
				DomainAdapter: c.domain,
			})

			ok = collectionErr == nil

			return
		},
	})
}

func (c *useCase) getEnvVersionCollection(envId []byte) (*EnvVersionCollection, error) {
	if env, envErr := c.domain.Env(envId); envErr != nil || env == nil {
		return nil, errors.Wrapf(code_push.ErrEnvNotFound, "envId: %s", envId)
	}

	collection, hasCollection := c.envVersionCollectionCache.Get(string(envId))

	if hasCollection {
		return collection.(*EnvVersionCollection), nil
	} else {
		return nil, errors.Wrapf(code_push.ErrEnvNotFound, "can not find version collection, envId: %s", envId)
	}
}

type VersionUseCaseConfig struct {
	VersionService domain.VersionService
	EnvService     domain.EnvService
}
