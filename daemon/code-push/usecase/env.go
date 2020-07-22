package usecase

import (
	"github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (c *useCase) CreateEnv(branchId, envName []byte) (*code_push.Env, error) {
	if branchId == nil || envName == nil {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "branchId or envName can't not be empty")
	}

	isBranchAvailable := c.domain.IsBranchAvailable(branchId)
	if !isBranchAvailable {
		return nil, errors.Wrapf(code_push.ErrBranchNotFound, "branchId:%s", branchId)
	}

	isEnvNameExisted, envNameExistedErr := c.domain.IsEnvNameExisted(branchId, envName)
	if envNameExistedErr != nil {
		return nil, errors.Wrap(envNameExistedErr, "failed to check whether env name existed")
	}
	if isEnvNameExisted {
		return nil, errors.Wrapf(code_push.ErrEnvNameExisted, "branchId:%s, envName:%s", branchId, envName)
	}

	envId := generateEnvId(string(branchId))
	encToken, encTokenErr := generateEnvEncToken()

	if encTokenErr != nil {
		return nil, errors.Wrapf(encTokenErr, "generate env enc token failed")
	}

	envToCreate := &code_push.Env{
		BranchId: string(branchId),
		ID:       envId,
		Name:     string(envName),
		EncToken: encToken,
	}

	createErr := c.domain.CreateEnv(envToCreate)
	if createErr != nil {
		return nil, errors.WithStack(createErr)
	}

	return envToCreate, nil
}

func (c *useCase) GetEnv(envId []byte) (*code_push.Env, error) {
	if envId == nil {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "envId is empty")
	}

	envEntity, fetchErr := c.domain.Env(envId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
	}
	if envEntity == nil {
		return nil, nil
	}

	return envEntity, nil
}

func (c *useCase) DeleteEnv(envId []byte) error {
	if envId == nil {
		return errors.Wrapf(code_push.ErrParamsInvalid, "envId is empty")
	}

	if !c.domain.IsEnvAvailable(envId) {
		return errors.Wrapf(code_push.ErrEnvNotFound, "envId:%s", envId)
	}

	deleteErr := c.domain.DeleteEnv(envId)
	if deleteErr != nil {
		return errors.WithStack(deleteErr)
	}

	return nil
}

func (c *useCase) GetEnvEncToken(envId []byte) ([]byte, error) {
	if envId == nil {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "envId is empty")
	}

	envEntity, fetchErr := c.domain.Env(envId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
	}
	if envEntity == nil {
		return nil, errors.Wrapf(code_push.ErrEnvNotFound, "envId: %s", envId)
	}

	return []byte(envEntity.EncToken), nil
}

func generateEnvEncToken() (string, error) {
	token, err := util.RandomPass(16, 8, 0, false, true)

	return token, err
}

func generateEnvId(branchId string) string {
	return util.EncodeMD5(branchId + "/" + uuid.NewV4().String())
}
