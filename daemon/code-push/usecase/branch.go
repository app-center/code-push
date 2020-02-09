package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/model"
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
	"github.com/funnyecho/code-push/daemon/code-push/domain/service"
)

type Branch struct {
	BranchName string
}

func toBranch(branch *model.Branch) *Branch {
	return &Branch{
		BranchName: branch.BranchName(),
	}
}

type IBranchUseCase interface {
	CreateBranch(branchName, branchAuthHost string) (*Branch, error)
	UpdateBranchName(branchId, branchName string) error
	GetBranch(branchId string) (*Branch, error)
	DeleteBranch(branchId string) error
}

type branchUseCase struct {
	branchRepo    repository.IBranch
	branchService service.IBranchService
}

func (b *branchUseCase) CreateBranch(branchName, branchAuthHost string) (*Branch, error) {
	return nil, nil
}

func (b *branchUseCase) UpdateBranchName(branchId, branchName string) error {
	return nil
}

func (b *branchUseCase) GetBranch(branchId string) (*Branch, error) {
	return nil, nil
}

func (b *branchUseCase) DeleteBranch(branchId string) error {
	return nil
}

type BranchUseCaseConfig struct {
	branchRepo    repository.IBranch
	branchService service.IBranchService
}

func NewBranchUseCase(config BranchUseCaseConfig) IBranchUseCase {
	return &branchUseCase{
		branchRepo:    config.branchRepo,
		branchService: config.branchService,
	}
}
