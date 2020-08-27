package middleware

import (
	"github.com/funnyecho/code-push/gateway/client/usecase"
)

func New(uc usecase.UseCase) *Middleware {
	return &Middleware{uc}
}

type Middleware struct {
	uc usecase.UseCase
}
