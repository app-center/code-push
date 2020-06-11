package bolt

import "github.com/funnyecho/code-push/daemon/code-push/domain"

var _ domain.IBranchService = &BranchService{}

type BranchService struct {
	client *Client
}

func (b *BranchService) Branch(branchId string) (*domain.Branch, error) {
	panic("implement me")
}

func (b *BranchService) CreateBranch(branch *domain.Branch) error {
	panic("implement me")
}

func (b *BranchService) DeleteBranch(branchId string) error {
	panic("implement me")
}
