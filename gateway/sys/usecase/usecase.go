package usecase

import "github.com/funnyecho/code-push/pkg/log"

func NewUseCase(config CtorConfig, optionFns ...func(*Options)) UseCase {
	ctorOptions := &Options{
		RootUserName: "",
		RootUserPwd:  "",
	}

	for _, fn := range optionFns {
		fn(ctorOptions)
	}

	uc := &useCase{
		adapters: &adapters{codePush: config.CodePushAdapter, session: config.SessionAdapter},
		Logger:   config.Logger,
		options:  ctorOptions,
	}

	return uc
}

type CtorConfig struct {
	CodePushAdapter
	SessionAdapter
	log.Logger
}

type useCase struct {
	*adapters
	log.Logger
	options *Options
}

type adapters struct {
	codePush CodePushAdapter
	session  SessionAdapter
}

type Options struct {
	RootUserName string
	RootUserPwd  string
}
