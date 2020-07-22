package usecase

import "github.com/funnyecho/code-push/gateway/sys"

func (u *useCase) CreateBranch(branchName []byte) (*sys.Branch, error) {
	return u.codePush.CreateBranch(branchName)
}

func (u *useCase) DeleteBranch(branchId []byte) error {
	return u.codePush.DeleteBranch(branchId)
}
