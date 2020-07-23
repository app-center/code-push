package usecase

import (
	"github.com/funnyecho/code-push/pkg/jwt"
	"time"
)

func NewUseCase(config CtorConfig, optionFns ...func(*Options)) UseCase {
	ctorOptions := &Options{
		RootUserName: "",
		RootUserPwd:  "",
		JwtSecret:    "",
		JwtIssuer:    "",
		JwtLifetime:  0,
	}

	for _, fn := range optionFns {
		fn(ctorOptions)
	}

	sysJwt := jwt.NewJwt(func(options *jwt.Options) {
		options.Secret = ctorOptions.JwtSecret
		options.Issuer = ctorOptions.JwtIssuer
		options.Lifetime = ctorOptions.JwtLifetime
	})

	return &useCase{
		&adapters{codePush: config.CodePushAdapter},
		ctorOptions,
		sysJwt,
	}
}

type CtorConfig struct {
	CodePushAdapter
}

type useCase struct {
	*adapters
	options *Options
	jwt     *jwt.Jwt
}

type adapters struct {
	codePush CodePushAdapter
}

type Options struct {
	RootUserName string
	RootUserPwd  string
	JwtSecret    string
	JwtIssuer    string
	JwtLifetime  time.Duration
}
