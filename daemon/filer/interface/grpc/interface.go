package grpc

import (
	"github.com/funnyecho/code-push/daemon/filer/usecase"
)

type Endpoints interface {
	usecase.File
	usecase.Upload
}
