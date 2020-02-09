package service

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
)

type IVersionService interface {
}

type versionService struct {
	versionRepo repository.IVersion
}

type VersionServiceConfig struct {
	VersionRepo repository.IVersion
}

func NewVersionService(config VersionServiceConfig) IVersionService {
	return &versionService{
		versionRepo: config.VersionRepo,
	}
}
