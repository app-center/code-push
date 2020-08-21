package usecase

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/sys"
	"time"
)

func (u *useCase) CreateBranch(ctx context.Context, branchName []byte) (*sys.Branch, error) {
	b, err := u.codePush.CreateBranch(ctx, branchName)
	return unmarshalBranch(b), err
}

func (u *useCase) DeleteBranch(ctx context.Context, branchId []byte) error {
	return u.codePush.DeleteBranch(ctx, branchId)
}

func unmarshalBranch(b *pb.BranchResponse) *sys.Branch {
	if b == nil {
		return nil
	}

	return &sys.Branch{
		ID:         b.GetBranchId(),
		Name:       b.GetBranchName(),
		EncToken:   b.GetBranchEncToken(),
		CreateTime: time.Unix(0, b.CreateTime),
	}
}
