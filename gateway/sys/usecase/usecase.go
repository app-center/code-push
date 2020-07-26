package usecase

func NewUseCase(config CtorConfig, optionFns ...func(*Options)) UseCase {
	ctorOptions := &Options{
		RootUserName: "",
		RootUserPwd:  "",
	}

	for _, fn := range optionFns {
		fn(ctorOptions)
	}

	return &useCase{
		&adapters{codePush: config.CodePushAdapter},
		ctorOptions,
	}
}

type CtorConfig struct {
	CodePushAdapter
	SessionAdapter
}

type useCase struct {
	*adapters
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
