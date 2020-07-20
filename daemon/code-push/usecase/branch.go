package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (c *useCase) CreateBranch(branchName []byte) (*code_push.Branch, error) {
	if branchName == nil {
		return nil, errors.Wrapf(
			code_push.ErrParamsInvalid,
			"branchName: %s",
			branchName,
		)
	}

	if nameExisted, nameExistedErr := c.domain.IsBranchNameExisted(branchName); nameExistedErr != nil {
		return nil, errors.Wrapf(
			nameExistedErr,
			"failed to check branch name was existed",
		)
	} else if nameExisted {
		return nil, errors.Wrapf(
			code_push.ErrBranchNameExisted,
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

	branchToCreate := &code_push.Branch{
		ID:       generateBranchId(string(branchName)),
		Name:     string(branchName),
		EncToken: encToken,
	}

	createErr := c.domain.CreateBranch(branchToCreate)
	if createErr != nil {
		return nil, errors.WithStack(createErr)
	}

	return branchToCreate, nil
}

func (c *useCase) GetBranch(branchId []byte) (*code_push.Branch, error) {
	if branchId == nil {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "branchId is empty")
	}

	branchEntity, fetchErr := c.domain.Branch(branchId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
	}
	if branchEntity == nil {
		return nil, nil
	}

	return branchEntity, nil
}

func (c *useCase) DeleteBranch(branchId []byte) error {
	if branchId == nil {
		return errors.Wrapf(code_push.ErrParamsInvalid, "branchId is empty")
	}

	deleteErr := c.domain.DeleteBranch(branchId)
	if deleteErr != nil {
		return errors.WithStack(deleteErr)
	}

	return nil
}

func (c *useCase) GetBranchEncToken(branchId []byte) ([]byte, error) {
	if branchId == nil {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "branchId is empty")
	}

	branchEntity, fetchErr := c.domain.Branch(branchId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
	}
	if branchEntity == nil {
		return nil, errors.Wrapf(code_push.ErrBranchNotFound, "branchId:%s", branchId)
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
