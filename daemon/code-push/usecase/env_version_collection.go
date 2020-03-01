package usecase

import (
	"github.com/Workiva/go-datastructures/tree/avl"
	"github.com/funnyecho/code-push/daemon/code-push/domain/model"
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
	"github.com/funnyecho/code-push/daemon/code-push/domain/service"
	"github.com/funnyecho/code-push/daemon/code-push/usecase/errors"
	"github.com/funnyecho/code-push/pkg/semver"
	"sync"
)

type versionCompatTreeNode struct {
	compatVersion  *semver.SemVer
	appVersionList []*semver.SemVer
}

func (n *versionCompatTreeNode) Compare(compare avl.Entry) int {
	return n.compatVersion.Compare(compare.(*versionCompatTreeNode).compatVersion)
}

type envVersionCollection struct {
	envId string
	mux   *sync.RWMutex

	versionRepo    repository.IVersion
	versionService service.IVersionService

	envRepo    repository.IEnv
	envService service.IEnvService

	versionList       model.VersionList
	versionCompatTree *avl.Immutable
}

func (c *envVersionCollection) ReleaseVersion(params IVersionReleaseParams) error {
	panic("implement me")
}

func (c *envVersionCollection) RollbackVersion(appVersion string) (*Version, error) {
	panic("implement me")
}

func (c *envVersionCollection) UpdateVersion(appVersion string, params IVersionUpdateParams) error {
	panic("implement me")
}

func (c *envVersionCollection) GetVersion(appVersion string) (*Version, error) {
	panic("implement me")
}

func (c *envVersionCollection) ListVersions() (VersionList, error) {
	panic("implement me")
}

func (c *envVersionCollection) VersionCompatQuery(appVersion string) (IVersionCompatQueryResult, error) {
	panic("implement me")
}

func (c *envVersionCollection) init() error {
	defer func() {
		c.mux.Unlock()
	}()

	c.mux.Lock()

	return nil
}

func (c *envVersionCollection) resetSource() error {
	panic("implement me")
}

type envVersionCollectionConfig struct {
	EnvId          string
	VersionRepo    repository.IVersion
	VersionService service.IVersionService
	EnvRepo        repository.IEnv
	EnvService     service.IEnvService
}

func newEnvVersionCollection(config envVersionCollectionConfig) (*envVersionCollection, error) {
	if config.VersionRepo == nil ||
		config.VersionService == nil ||
		config.EnvRepo == nil ||
		config.EnvService == nil ||
		len(config.EnvId) == 0 {
		return nil, errors.ThrowVersionOperationForbiddenError(
			nil,
			"invalid version collection params",
			nil,
		)
	}

	collection := &envVersionCollection{
		envId:          config.EnvId,
		mux:            &sync.RWMutex{},
		versionRepo:    config.VersionRepo,
		versionService: config.VersionService,

		envRepo:    config.EnvRepo,
		envService: config.EnvService,
	}

	if initErr := collection.init(); initErr != nil {
		return nil, initErr
	} else {
		return collection, nil
	}
}
