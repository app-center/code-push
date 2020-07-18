package usecase

import (
	"github.com/funnyecho/code-push/pkg/cache"
)

func NewUseCase(config CtorConfig) *UseCase {
	instance := &UseCase{adapters: adapters{domain: config.DomainAdapter}}

	instance.initVersionUseCase()

	return instance
}

type UseCase struct {
	adapters
	envVersionCollectionCache *cache.Cache
}

type CtorConfig struct {
	DomainAdapter
}

type adapters struct {
	domain DomainAdapter
}
