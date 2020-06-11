package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain"
	"github.com/funnyecho/code-push/daemon/code-push/usecase/errors"
	"github.com/funnyecho/code-push/pkg/util"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Branch struct {
	BranchId   string
	BranchName string
	CreateTime time.Time
}

func toBranch(branch *domain.Branch) *Branch {
	return &Branch{
		BranchId:   branch.ID,
		BranchName: branch.Name,
		CreateTime: branch.CreateTime,
	}
}

type IBranchUseCase interface {
	CreateBranch(branchName, branchAuthHost string) (*Branch, error)
	GetBranch(branchId string) (*Branch, error)
	DeleteBranch(branchId string) error
	GetBranchEncToken(branchId string) (string, error)
}

type branchUseCase struct {
	branchService domain.IBranchService
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

	encToken, encTokenErr := generateBranchEncToken()
	if encTokenErr != nil {
		return nil, errors.ThrowBranchInvalidEncTokenError(encTokenErr)
	}

	//branchToCreate := model.NewBranch(model.BranchConfig{
	//	Id:         generateBranchId(branchName),
	//	Name:       branchName,
	//	AuthHost:   branchAuthHost,
	//	EncToken:   encToken,
	//	CreateTime: time.Now(),
	//})

	branchToCreate := &domain.Branch{
		ID:         generateBranchId(branchName),
		Name:       branchName,
		AuthHost:   branchAuthHost,
		EncToken:   encToken,
		CreateTime: time.Now(),
	}

	createErr := b.branchService.CreateBranch(branchToCreate)
	if createErr != nil {
		return nil, errors.ThrowBranchSaveError(createErr, branchToCreate)
	}

	return toBranch(branchToCreate), nil
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

	branchEntity, fetchErr := b.branchService.Branch(branchId)
	if fetchErr != nil {
		return nil, errors.ThrowBranchNotFoundError(branchId, fetchErr)
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

	deleteErr := b.branchService.DeleteBranch(branchId)
	if deleteErr != nil {
		return errors.ThrowBranchDeleteFailedError(branchId, "")
	}

	return nil
}

func (b *branchUseCase) GetBranchEncToken(branchId string) (string, error) {
	if len(branchId) == 0 {
		return "", errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "branchId is empty",
			Params: errors.MetaFields{
				"branchId": branchId,
			},
		})
	}

	branchEntity, fetchErr := b.branchService.Branch(branchId)
	if fetchErr != nil {
		return "", errors.ThrowBranchNotFoundError(branchId, fetchErr)
	}

	return branchEntity.EncToken, nil
}

type BranchUseCaseConfig struct {
	BranchService domain.IBranchService
}

func NewBranchUseCase(config BranchUseCaseConfig) (IBranchUseCase, error) {
	if config.BranchService == nil {
		panic("invalid branch use case params")
	}

	return &branchUseCase{
		branchService: config.BranchService,
	}, nil
}

func generateBranchEncToken() (string, error) {
	token, err := util.RandomPass(16, 8, 0, false, true)

	return token, err
}

func generateBranchId(branchName string) string {
	return util.EncodeBase64(util.EncodeMD5(branchName + "/" + uuid.NewV4().String()))
}
