package service

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
	"sync"
)

type IEnvService interface {
}

type envService struct {
	mtx sync.RWMutex

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
