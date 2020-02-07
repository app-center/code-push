package service

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
	"sync"
)

type IVersionService interface {
}

type versionService struct {
	mtx         sync.RWMutex
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
