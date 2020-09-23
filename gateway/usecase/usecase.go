package usecase

import (
	"github.com/funnyecho/code-push/pkg/log"
)

func New(configFn func(*CtorConfig)) UseCase {
	config := &CtorConfig{
		Options: &options{},
	}

	configFn(config)

	return &useCase{
		daemon:  config.DaemonAdapter,
		Logger:  config.Logger,
		options: config.Options,
	}
}

type CtorConfig struct {
	DaemonAdapter
	log.Logger

	Options *options
}

type useCase struct {
	daemon DaemonAdapter
	log.Logger

	*options
}

type options struct {
	RootUserName string
	RootUserPwd  string
}
