package usecase

import (
	"github.com/funnyecho/code-push/daemon"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *useCase) CreateBranch(branchName []byte) (*daemon.Branch, error) {
	if branchName == nil {
		return nil, errors.Wrapf(
			daemon.ErrParamsInvalid,
			"branchName: %s",
			branchName,
		)
	}

	if nameExisted, nameExistedErr := uc.domain.IsBranchNameExisted(branchName); nameExistedErr != nil {
		return nil, errors.Wrapf(
			nameExistedErr,
			"failed to check branch name was existed",
		)
	} else if nameExisted {
		return nil, errors.Wrapf(
			daemon.ErrBranchNameExisted,
			"branchName: %s",
			branchName,
		)
	}

	encToken, encTokenErr := generateBranchEncToken()
	if encTokenErr != nil {
		return nil, errors.Wrapf(
			encTokenErr,
			"generate branch enc token failed",
		)
	}

	branchToCreate := &daemon.Branch{
		ID:       generateBranchId(string(branchName)),
		Name:     string(branchName),
		EncToken: encToken,
	}

	createErr := uc.domain.CreateBranch(branchToCreate)
	if createErr != nil {
		return nil, errors.WithStack(createErr)
	}

	return branchToCreate, nil
}

func (uc *useCase) GetBranch(branchId []byte) (*daemon.Branch, error) {
	if branchId == nil {
		return nil, errors.Wrapf(daemon.ErrParamsInvalid, "branchId is empty")
	}

	branchEntity, fetchErr := uc.domain.Branch(branchId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
	}
	if branchEntity == nil {
		return nil, nil
	}

	return branchEntity, nil
}

func (uc *useCase) DeleteBranch(branchId []byte) error {
	if branchId == nil {
		return errors.Wrapf(daemon.ErrParamsInvalid, "branchId is empty")
	}

	deleteErr := uc.domain.DeleteBranch(branchId)
	if deleteErr != nil {
		return errors.WithStack(deleteErr)
	}

	return nil
}

func (uc *useCase) GetBranchEncToken(branchId []byte) ([]byte, error) {
	if branchId == nil {
		return nil, errors.Wrapf(daemon.ErrParamsInvalid, "branchId is empty")
	}

	branchEntity, fetchErr := uc.domain.Branch(branchId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
	}
	if branchEntity == nil {
		return nil, errors.Wrapf(daemon.ErrBranchNotFound, "branchId:%s", branchId)
	}

	return []byte(branchEntity.EncToken), nil
}

func generateBranchEncToken() (string, error) {
	token, err := util.RandomPass(16, 8, 0, false, true)

	return token, err
}

func generateBranchId(branchName string) string {
	return util.EncodeBase64(util.EncodeMD5(branchName + "/" + uuid.NewV4().String()))
}
