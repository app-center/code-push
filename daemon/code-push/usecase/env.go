package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
	"github.com/funnyecho/code-push/daemon/code-push/domain/service"
)

type EnvUseCase interface {
}

type envUseCase struct {
	envRepo    repository.IEnv
	envService service.IEnvService
}

type EnvUseCaseConfig struct {
	EnvRepo    repository.IEnv
	EnvService service.IEnvService
}

func NewEnvUseCase(config EnvUseCaseConfig) EnvUseCase {
	return &envUseCase{
		envRepo:    config.EnvRepo,
		envService: config.EnvService,
	}
}
