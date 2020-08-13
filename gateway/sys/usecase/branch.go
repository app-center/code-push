package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/sys"
	"time"
)

func (u *useCase) CreateBranch(branchName []byte) (*sys.Branch, error) {
	b, err := u.codePush.CreateBranch(branchName)
	return unmarshalBranch(b), err
}

func (u *useCase) DeleteBranch(branchId []byte) error {
	return u.codePush.DeleteBranch(branchId)
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
