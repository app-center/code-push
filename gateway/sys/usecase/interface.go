package usecase

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	sessionAdapter "github.com/funnyecho/code-push/daemon/session/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway/sys"
)

type UseCase interface {
	Auth
	Branch
}

type Auth interface {
	Auth(ctx context.Context, name, pwd string) error
	SignToken(ctx context.Context) ([]byte, error)
	VerifyToken(ctx context.Context, token []byte) error
}

type Branch interface {
	CreateBranch(ctx context.Context, branchName []byte) (*sys.Branch, error)
	DeleteBranch(ctx context.Context, branchId []byte) error
}

type CodePushAdapter interface {
	CreateBranch(ctx context.Context, branchName []byte) (*pb.BranchResponse, error)
	DeleteBranch(ctx context.Context, branchId []byte) error
}

type SessionAdapter interface {
	GenerateAccessToken(ctx context.Context, issuer sessionAdapter.AccessTokenIssuer, subject string) ([]byte, error)
	VerifyAccessToken(ctx context.Context, token string) (subject []byte, err error)
}
