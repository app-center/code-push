package bolt

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain"
)

var _ domain.IVersionService = &VersionService{}

type VersionService struct {
	client *Client
}

func (v *VersionService) Version(envId, appVersion string) (*domain.Version, error) {
	panic("implement me")
}

func (v *VersionService) VersionsWithEnvId(envId string) (domain.VersionList, error) {
	panic("implement me")
}

func (v *VersionService) CreateVersion(version *domain.Version) error {
	panic("implement me")
}
