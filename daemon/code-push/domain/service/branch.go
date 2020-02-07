package service

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
	"sync"
)

type IBranchService interface {
}

type branchService struct {
	mtx        sync.RWMutex
	branchRepo repository.IBranch
}

type BranchServiceConfig struct {
	BranchRepo repository.IBranch
}

func NewBranchService(config BranchServiceConfig) IBranchService {
	return &branchService{
		branchRepo: config.BranchRepo,
	}
}
