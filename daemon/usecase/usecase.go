package usecase

import (
	"github.com/funnyecho/code-push/pkg/cache"
	"github.com/funnyecho/code-push/pkg/log"
	go_cache "github.com/patrickmn/go-cache"
)

func NewUseCase(configFn func(*CtorConfig)) UseCase {
	config := &CtorConfig{}
	configFn(config)

	instance := &useCase{
		domain: config.DomainAdapter,
		aliOss: config.AliOssAdapter,
		logger: config.Logger,
	}

	instance.initVersionUseCase()
	instance.initAccessTokenUseCase()

	return instance
}

type useCase struct {
	domain                    DomainAdapter
	aliOss					  AliOssAdapter
	logger                    log.Logger
	accessTokenCache		  *go_cache.Cache
	envVersionCollectionCache *cache.Cache
}

type CtorConfig struct {
	DomainAdapter
	AliOssAdapter
	log.Logger
}
