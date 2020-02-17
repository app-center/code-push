package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/model"
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
	"github.com/funnyecho/code-push/daemon/code-push/domain/service"
	"github.com/funnyecho/code-push/daemon/code-push/usecase/errors"
	"github.com/funnyecho/code-push/pkg/util"
	uuid "github.com/satori/go.uuid"
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
	UpdateBranch(branchId string, params IBranchUpdateParams) error
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
		Id:         generateBranchId(branchName),
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

func (b *branchUseCase) UpdateBranch(branchId string, params IBranchUpdateParams) error {
	if len(branchId) == 0 {
		return errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "branchId is empty",
			Params: errors.MetaFields{
				"branchId": branchId,
			},
		})
	}

	updateBranchName, newBranchName := params.BranchName()
	updateBranchAuthHost, newBranchAuthHost := params.BranchAuthHost()

	if !updateBranchName && !updateBranchAuthHost {
		return errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "update params do not have valid keys",
		})
	}

	hasInvalidParamsErr := false
	invalidErrParams := errors.MetaFields{}

	if updateBranchName && len(newBranchName) == 0 {
		hasInvalidParamsErr = true
		invalidErrParams["branchName"] = newBranchName
	}

	if updateBranchAuthHost && len(newBranchAuthHost) == 0 {
		hasInvalidParamsErr = true
		invalidErrParams["branchAuthHost"] = newBranchAuthHost
	}

	if hasInvalidParamsErr {
		return errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg:    "invalid update params",
			Params: invalidErrParams,
		})
	}

	entity, findErr := b.branchRepo.FindBranch(branchId)

	if findErr != nil {
		return errors.ThrowBranchCanNotFoundError(branchId)
	}

	if updateBranchName {
		entity.SetBranchName(newBranchName)
	}

	if updateBranchAuthHost {
		entity.SetBranchAuthHost(newBranchAuthHost)
	}

	_, updateErr := b.branchRepo.SaveBranch(entity)
	if updateErr != nil {
		return errors.ThrowBranchSaveError(updateErr, params)
	}

	return nil
}

func (b *branchUseCase) GetBranch(branchId string) (*Branch, error) {
	if len(branchId) == 0 {
		return nil, errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "branchId is empty",
			Params: errors.MetaFields{
				"branchId": branchId,
			},
		})
	}

	branchEntity, fetchErr := b.branchRepo.FindBranch(branchId)
	if fetchErr != nil {
		return nil, errors.ThrowBranchCanNotFoundError(branchId)
	}

	return toBranch(branchEntity), nil
}

func (b *branchUseCase) DeleteBranch(branchId string) error {
	if len(branchId) == 0 {
		return errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "branchId is empty",
			Params: errors.MetaFields{
				"branchId": branchId,
			},
		})
	}

	deleteErr := b.branchRepo.DeleteBranchById(branchId)
	if deleteErr != nil {
		return errors.ThrowBranchDeleteFailedError(branchId, "")
	}

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

func generateBranchId(branchName string) string {
	return util.EncodeBase64(util.EncodeMD5(branchName + "/" + uuid.NewV4().String()))
}
