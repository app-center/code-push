package usecase

func NewUseCase(config CtorConfig, optionFns ...func(*options)) UseCase {
	ctorOptions := &options{
		RootUserName: nil,
		RootUserPwd:  nil,
	}

	for _, fn := range optionFns {
		fn(ctorOptions)
	}

	return &useCase{
		adapters: &adapters{codePush: config.CodePushAdapter},
		options:  ctorOptions,
	}
}

type CtorConfig struct {
	CodePushAdapter
}

type useCase struct {
	*adapters
	*options
}

type adapters struct {
	codePush CodePushAdapter
}

type options struct {
	RootUserName string
	RootUserPwd  string
}
