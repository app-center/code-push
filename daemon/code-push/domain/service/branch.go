package service

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
)

type IBranchService interface {
	IsBranchNameExisted(branchName string) bool
}

type branchService struct {
	branchRepo repository.IBranch
}

func (b *branchService) IsBranchNameExisted(branchName string) bool {
	if len(branchName) == 0 {
		return false
	}

	_, err := b.branchRepo.FindBranchByName(branchName)

	return err == nil
}

type BranchServiceConfig struct {
	BranchRepo repository.IBranch
}

func NewBranchService(config BranchServiceConfig) IBranchService {
	return &branchService{
		branchRepo: config.BranchRepo,
	}
}
