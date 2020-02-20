package service

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
)

type IBranchService interface {
	IsBranchExisted(branchId string) bool
	IsBranchNameExisted(branchName string) bool
}

type branchService struct {
	branchRepo repository.IBranch
}

func (b *branchService) IsBranchExisted(branchId string) bool {
	branch, err := b.branchRepo.FindBranch(branchId)

	return err != nil && branch != nil
}

func (b *branchService) IsBranchNameExisted(branchName string) bool {
	if len(branchName) == 0 {
		return false
	}

	branch, err := b.branchRepo.FindBranchByName(branchName)

	return err == nil && branch != nil
}

type BranchServiceConfig struct {
	BranchRepo repository.IBranch
}

func NewBranchService(config BranchServiceConfig) IBranchService {
	return &branchService{
		branchRepo: config.BranchRepo,
	}
}
