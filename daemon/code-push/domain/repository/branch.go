package repository

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/model"
)

type IBranch interface {
	FindBranch(branchId string) (*model.Branch, error)
	FindBranchByName(branchName string) (*model.Branch, error)
	SaveBranch(branch model.Branch) (*model.Branch, error)
	DeleteBranchById(branchId string) error
}
