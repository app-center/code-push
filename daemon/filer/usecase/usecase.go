package usecase

func NewUseCase(config CtorConfig) *UseCase {
	return &UseCase{adapters{
		domain: config.DomainAdapter,
		aliOss: config.AliOssAdapter,
	}}
}

type UseCase struct {
	adapters
}

type CtorConfig struct {
	DomainAdapter
	AliOssAdapter
}

type adapters struct {
	domain DomainAdapter
	aliOss AliOssAdapter
}
