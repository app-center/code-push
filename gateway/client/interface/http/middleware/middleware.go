package middleware

import (
	"github.com/funnyecho/code-push/gateway/client"
	"github.com/funnyecho/code-push/gateway/client/usecase"
)

func New(uc usecase.UseCase, metrics *client.Metrics) *Middleware {
	return &Middleware{uc, metrics}
}

type Middleware struct {
	uc      usecase.UseCase
	metrics *client.Metrics
}
