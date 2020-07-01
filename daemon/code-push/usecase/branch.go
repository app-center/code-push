package usecase

import (
	code_push "github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/daemon/code-push/domain"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

func NewBranchUseCase(config BranchUseCaseConfig) (IBranch, error) {
	if config.BranchService == nil {
		panic("invalid branch use case params")
	}

	return &branchUseCase{
		branchService: config.BranchService,
	}, nil
}

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

type IBranch interface {
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
		return nil, errors.Wrapf(
			code_push.ErrParamsInvalid,
			"branchName: %s, branchAuthHost: %s",
			branchName,
			branchAuthHost,
		)
	}

	encToken, encTokenErr := generateBranchEncToken()
	if encTokenErr != nil {
		return nil, errors.Wrapf(
			encTokenErr,
			"generate branch enc token failed",
		)
	}

	branchToCreate := &domain.Branch{
		ID:       generateBranchId(branchName),
		Name:     branchName,
		AuthHost: branchAuthHost,
		EncToken: encToken,
	}

	createErr := b.branchService.CreateBranch(branchToCreate)
	if createErr != nil {
		return nil, errors.WithStack(createErr)
	}

	return toBranch(branchToCreate), nil
}

func (b *branchUseCase) GetBranch(branchId string) (*Branch, error) {
	if len(branchId) == 0 {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "branchId is empty")
	}

	branchEntity, fetchErr := b.branchService.Branch(branchId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
	}
	if branchEntity == nil {
		return nil, errors.WithMessagef(code_push.ErrBranchNotFound, "branchId: %v", branchId)
	}

	return toBranch(branchEntity), nil
}

func (b *branchUseCase) DeleteBranch(branchId string) error {
	if len(branchId) == 0 {
		return errors.Wrapf(code_push.ErrParamsInvalid, "branchId is empty")
	}

	deleteErr := b.branchService.DeleteBranch(branchId)
	if deleteErr != nil {
		return errors.WithStack(deleteErr)
	}

	return nil
}

func (b *branchUseCase) GetBranchEncToken(branchId string) (string, error) {
	if len(branchId) == 0 {
		return "", errors.Wrapf(code_push.ErrParamsInvalid, "branchId is empty")
	}

	branchEntity, fetchErr := b.branchService.Branch(branchId)
	if fetchErr != nil {
		return "", errors.WithStack(fetchErr)
	}

	return branchEntity.EncToken, nil
}

type BranchUseCaseConfig struct {
	BranchService domain.IBranchService
}

func generateBranchEncToken() (string, error) {
	token, err := util.RandomPass(16, 8, 0, false, true)

	return token, err
}

func generateBranchId(branchName string) string {
	return util.EncodeBase64(util.EncodeMD5(branchName + "/" + uuid.NewV4().String()))
}
