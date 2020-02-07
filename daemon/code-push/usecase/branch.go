package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
	"github.com/funnyecho/code-push/daemon/code-push/domain/service"
)

type BranchUseCase interface {
}

type branchUseCase struct {
	branchRepo    repository.IBranch
	branchService service.IBranchService
}

type BranchUseCaseConfig struct {
	branchRepo    repository.IBranch
	branchService service.IBranchService
}

func NewBranchUseCase(config BranchUseCaseConfig) BranchUseCase {
	return &branchUseCase{
		branchRepo:    config.branchRepo,
		branchService: config.branchService,
	}
}
