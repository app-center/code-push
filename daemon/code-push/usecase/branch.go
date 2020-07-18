package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (c *UseCase) CreateBranch(branchName, branchAuthHost []byte) (*code_push.Branch, error) {
	if branchName == nil || branchAuthHost == nil {
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

	branchToCreate := &code_push.Branch{
		ID:       generateBranchId(string(branchName)),
		Name:     string(branchName),
		AuthHost: string(branchAuthHost),
		EncToken: encToken,
	}

	createErr := c.domain.CreateBranch(branchToCreate)
	if createErr != nil {
		return nil, errors.WithStack(createErr)
	}

	return branchToCreate, nil
}

func (c *UseCase) GetBranch(branchId []byte) (*code_push.Branch, error) {
	if branchId == nil {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "branchId is empty")
	}

	branchEntity, fetchErr := c.domain.Branch(branchId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
	}
	if branchEntity == nil {
		return nil, errors.Wrapf(code_push.ErrBranchNotFound, "branchId: %v", branchId)
	}

	return branchEntity, nil
}

func (c *UseCase) DeleteBranch(branchId []byte) error {
	if branchId == nil {
		return errors.Wrapf(code_push.ErrParamsInvalid, "branchId is empty")
	}

	deleteErr := c.domain.DeleteBranch(branchId)
	if deleteErr != nil {
		return errors.WithStack(deleteErr)
	}

	return nil
}

func (c *UseCase) GetBranchEncToken(branchId []byte) ([]byte, error) {
	if branchId == nil {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "branchId is empty")
	}

	branchEntity, fetchErr := c.domain.Branch(branchId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
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
