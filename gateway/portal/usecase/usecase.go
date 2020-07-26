package usecase

func NewUseCase(config CtorConfig, optionsFns ...func(*Options)) UseCase {
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
		options: ctorOptions,
	}
}

type CtorConfig struct {
	CodePushAdapter
	SessionAdapter
	FilerAdapter
}

type useCase struct {
	*adapters
	options *Options
}

type adapters struct {
	codePush CodePushAdapter
	session  SessionAdapter
	filer    FilerAdapter
}

type Options struct {
}
