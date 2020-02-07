package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
	"github.com/funnyecho/code-push/daemon/code-push/domain/service"
)

type VersionUseCase interface {
}

type versionUseCase struct {
	versionRepo    repository.IVersion
	versionService service.IVersionService
}

type VersionUseCaseConfig struct {
	VersionRepo    repository.IVersion
	VersionService service.IVersionService
}

func NewVersionUseCase(config VersionUseCaseConfig) VersionUseCase {
	return &versionUseCase{
		versionRepo:    config.VersionRepo,
		versionService: config.VersionService,
	}
}
