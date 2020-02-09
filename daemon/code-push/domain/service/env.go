package service

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
)

type IEnvService interface {
}

type envService struct {
	envRepo repository.IEnv
}

type EnvServiceConfig struct {
	EnvRepo repository.IEnv
}

func NewEnvService(config EnvServiceConfig) IEnvService {
	return &envService{
		envRepo: config.EnvRepo,
	}
}
