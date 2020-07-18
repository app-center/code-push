package usecase

import (
	code_push "github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/daemon/code-push/domain"
	"github.com/funnyecho/code-push/pkg/cache"
	"github.com/pkg/errors"
)

func (c *UseCase) ReleaseVersion(params VersionReleaseParams) error {
	collection, collectionErr := c.getEnvVersionCollection(params.EnvId())

	if collectionErr != nil {
		return collectionErr
	}

	return collection.ReleaseVersion(params)
}

func (c *UseCase) GetVersion(envId, appVersion []byte) (*code_push.Version, error) {
	collection, collectionErr := c.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.GetVersion(appVersion)
}

func (c *UseCase) ListVersions(envId []byte) (code_push.VersionList, error) {
	collection, collectionErr := c.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.ListVersions()
}

func (c *UseCase) VersionStrictCompatQuery(envId, appVersion []byte) (VersionCompatQueryResult, error) {
	collection, collectionErr := c.getEnvVersionCollection(envId)

	if collectionErr != nil {
		return nil, collectionErr
	}

	return collection.VersionStrictCompatQuery(appVersion)
}

func (c *UseCase) initVersionUseCase() {
	c.envVersionCollectionCache = cache.New(cache.CtorConfig{
		Capacity: 10,
		AllocFunc: func(key cache.KeyType) (collection cache.ValueType, ok bool) {
			envId, isEnvIdType := key.([]byte)

			if !isEnvIdType {
				return nil, false
			}

			if env, envErr := c.domain.Env(envId); envErr != nil || env == nil {
				return nil, false
			}

			collection, collectionErr := NewEnvVersionCollection(EnvVersionCollectionConfig{
				EnvId:         envId,
				DomainAdapter: c.domain,
			})

			ok = collectionErr == nil

			return
		},
	})
}

func (c *UseCase) getEnvVersionCollection(envId []byte) (*EnvVersionCollection, error) {
	if env, envErr := c.domain.Env(envId); envErr != nil || env == nil {
		return nil, errors.Wrapf(code_push.ErrEnvNotFound, "envId: %s", envId)
	}

	collection, hasCollection := c.envVersionCollectionCache.Get(envId)

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
