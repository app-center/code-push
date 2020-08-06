package usecase

import "github.com/funnyecho/code-push/pkg/log"

func NewUseCase(config CtorConfig) *UseCase {
	return &UseCase{
		adapters{
			domain: config.DomainAdapter,
			aliOss: config.AliOssAdapter,
		},
		config.Logger,
	}
}

type UseCase struct {
	adapters
	log.Logger
}

type CtorConfig struct {
	DomainAdapter
	AliOssAdapter
	log.Logger
}

type adapters struct {
	domain DomainAdapter
	aliOss AliOssAdapter
}
