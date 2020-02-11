package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/model"
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
	"github.com/funnyecho/code-push/daemon/code-push/domain/service"
	"github.com/funnyecho/code-push/daemon/code-push/usecase/errors"
	"github.com/funnyecho/code-push/pkg/util"
	"time"
)

type Branch struct {
	BranchName string
}

func toBranch(branch model.Branch) *Branch {
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
	if len(branchName) == 0 || len(branchAuthHost) == 0 {
		return nil, errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Params: errors.MetaFields{
				"branchName":     branchName,
				"branchAuthHost": branchAuthHost,
			},
		})
	}

	if b.branchService.IsBranchNameExisted(branchName) {
		return nil, errors.ThrowBranchNameExistedError(branchName)
	}

	encToken, encTokenErr := generateBranchEncToken()
	if encTokenErr != nil {
		return nil, errors.ThrowBranchInvalidEncTokenError(encTokenErr)
	}

	branchToCreate := model.NewBranch(model.BranchConfig{
		Name:       branchName,
		AuthHost:   branchAuthHost,
		EncToken:   encToken,
		CreateTime: time.Now(),
	})

	branchCreated, createErr := b.branchRepo.SaveBranch(branchToCreate)
	if createErr != nil {
		return nil, errors.ThrowBranchSaveError(createErr, branchToCreate)
	}

	return toBranch(branchCreated), nil
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

func generateBranchEncToken() (string, error) {
	token, err := util.RandomPass(16, 8, 0, false, true)

	return token, err
}
