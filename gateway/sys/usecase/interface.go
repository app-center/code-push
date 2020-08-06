package usecase

import (
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
	GenerateAccessToken(subject string) ([]byte, error)
	VerifyAccessToken(token string) (subject []byte, err error)
}
