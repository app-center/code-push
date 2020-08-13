package usecase

import (
	sessionAdapter "github.com/funnyecho/code-push/daemon/session/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway/sys"
)

type UseCase interface {
	Auth
	Branch
}

type Auth interface {
	Auth(name, pwd string) error
	SignToken() ([]byte, error)
	VerifyToken(token []byte) error
}

type Branch interface {
	CreateBranch(branchName []byte) (*sys.Branch, error)
	DeleteBranch(branchId []byte) error
}

type CodePushAdapter interface {
	CreateBranch(branchName []byte) (*sys.Branch, error)
	DeleteBranch(branchId []byte) error
}

type SessionAdapter interface {
	GenerateAccessToken(issuer sessionAdapter.AccessTokenIssuer, subject string) ([]byte, error)
	VerifyAccessToken(token string) (subject []byte, err error)
}
