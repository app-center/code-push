package endpoint

import "github.com/funnyecho/code-push/gateway/usecase"

func New(uc usecase.UseCase) *Endpoints {
	return &Endpoints{uc}
}

type Endpoints struct {
	uc usecase.UseCase
}
