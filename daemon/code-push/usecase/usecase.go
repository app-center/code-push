package usecase

import (
	"github.com/funnyecho/code-push/pkg/cache"
)

func NewUseCase(config CtorConfig) UseCase {
	instance := &useCase{adapters: adapters{domain: config.DomainAdapter}}

	instance.initVersionUseCase()

	return instance
}

type useCase struct {
	adapters
	envVersionCollectionCache *cache.Cache
}

type CtorConfig struct {
	DomainAdapter
}

type adapters struct {
	domain DomainAdapter
}
