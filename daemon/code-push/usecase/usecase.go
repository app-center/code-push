package usecase

import (
	"github.com/funnyecho/code-push/pkg/cache"
	"github.com/funnyecho/code-push/pkg/log"
)

func NewUseCase(config CtorConfig) UseCase {
	instance := &useCase{
		adapters: adapters{domain: config.DomainAdapter},
		Logger:   config.Logger,
	}

	instance.initVersionUseCase()

	return instance
}

type useCase struct {
	adapters
	log.Logger
	envVersionCollectionCache *cache.Cache
}

type CtorConfig struct {
	DomainAdapter
	Logger log.Logger
}

type adapters struct {
	domain DomainAdapter
}
