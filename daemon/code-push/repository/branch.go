package repository

import (
	"github.com/funnyecho/code-push/daemon/code-push/model"
)

type IBranch interface {
	Find(branchId string) (*model.Branch, error)
	FindByName(branchName string) (*model.Branch, error)
	Create(branch model.Branch) (*model.Branch, error)
	Save(branch *model.Branch) error
	Delete(branch *model.Branch) error
}
