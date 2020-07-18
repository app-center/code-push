package grpc

import (
	"github.com/funnyecho/code-push/daemon/code-push/usecase"
)

type Endpoints interface {
	usecase.Branch
	usecase.Env
	usecase.Version
}
