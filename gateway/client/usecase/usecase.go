package usecase

import (
	"github.com/funnyecho/code-push/gateway/client"
	"github.com/funnyecho/code-push/pkg/log"
)

func NewUseCase(config *CtorConfig, optionsFns ...func(*Options)) UseCase {
	ctorOptions := &Options{}

	for _, fn := range optionsFns {
		fn(ctorOptions)
	}

	return &useCase{
		adapters: &adapters{
			codePush: config.CodePushAdapter,
			session:  config.SessionAdapter,
			filer:    config.FilerAdapter,
		},
		Logger:  config.Logger,
		options: ctorOptions,
	}
}

type CtorConfig struct {
	CodePushAdapter
	SessionAdapter
	FilerAdapter
	log.Logger
	*client.Metrics
}

type useCase struct {
	*adapters
	log.Logger
	*client.Metrics
	options *Options
}

type adapters struct {
	codePush CodePushAdapter
	session  SessionAdapter
	filer    FilerAdapter
}

type Options struct {
}
