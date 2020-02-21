package repository

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/model"
)

type IBranch interface {
	FirstBranch(branchId string) (*model.Branch, error)
	FirstBranchByName(branchName string) (*model.Branch, error)
	SaveBranch(branch model.Branch) (*model.Branch, error)
	DeleteBranchById(branchId string) error
}
