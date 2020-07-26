package grpc

import "github.com/funnyecho/code-push/daemon/session/usecase"

type Endpoints interface {
	usecase.AccessToken
}
