package usecase

import (
	"context"
	"github.com/funnyecho/code-push/daemon/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway"
	"time"
)

func (uc *useCase) CreateBranch(ctx context.Context, branchName []byte) (*gateway.Branch, error) {
	b, err := uc.daemon.CreateBranch(ctx, branchName)
	return unmarshalBranch(b), err
}

func (uc *useCase) DeleteBranch(ctx context.Context, branchId []byte) error {
	return uc.daemon.DeleteBranch(ctx, branchId)
}

func (uc *useCase) GetBranch(ctx context.Context, branchId string) (*gateway.Branch, error) {
	b, err := uc.daemon.GetBranch(ctx, branchId)
	return unmarshalBranch(b), err
}

func unmarshalBranch(b *pb.BranchResponse) *gateway.Branch {
	if b == nil {
		return nil
	}

	return &gateway.Branch{
		ID:         b.GetBranchId(),
		Name:       b.GetBranchName(),
		EncToken:   b.GetBranchEncToken(),
		CreateTime: time.Unix(0, b.CreateTime),
	}
}
