package usecase

import (
	"github.com/funnyecho/code-push/daemon"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *useCase) CreateEnv(branchId, envId, envName, envEncToken []byte) (*daemon.Env, error) {
	if branchId == nil || envName == nil {
		return nil, errors.Wrapf(daemon.ErrParamsInvalid, "branchId or envName can't not be empty")
	}

	isBranchAvailable := uc.domain.IsBranchAvailable(branchId)
	if !isBranchAvailable {
		return nil, errors.Wrapf(daemon.ErrBranchNotFound, "branchId:%s", branchId)
	}

	isEnvNameExisted, envNameExistedErr := uc.domain.IsEnvNameExisted(branchId, envName)
	if envNameExistedErr != nil {
		return nil, errors.Wrap(envNameExistedErr, "failed to check whether env name existed")
	}
	if isEnvNameExisted {
		return nil, errors.Wrapf(daemon.ErrEnvNameExisted, "branchId:%s, envName:%s", branchId, envName)
	}

	id := ""

	if envId != nil {
		if uc.domain.IsEnvAvailable(envId) {
			return nil, errors.Wrapf(daemon.ErrEnvExisted, "envId:%s existed", string(envId))
		} else {
			id = string(envId)
		}
	} else {
		id = generateEnvId(string(branchId))
	}

	encToken := ""

	if envEncToken == nil || len(envEncToken) == 0 {
		token, encTokenErr := generateEnvEncToken()

		if encTokenErr != nil {
			return nil, errors.Wrapf(encTokenErr, "generate env enc token failed")
		}

		encToken = token
	} else {
		encToken = string(envEncToken)
	}

	envToCreate := &daemon.Env{
		BranchId: string(branchId),
		ID:       id,
		Name:     string(envName),
		EncToken: encToken,
	}

	createErr := uc.domain.CreateEnv(envToCreate)
	if createErr != nil {
		return nil, errors.WithStack(createErr)
	}

	return envToCreate, nil
}

func (uc *useCase) GetEnv(envId []byte) (*daemon.Env, error) {
	if envId == nil {
		return nil, errors.Wrapf(daemon.ErrParamsInvalid, "envId is empty")
	}

	envEntity, fetchErr := uc.domain.Env(envId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
	}
	if envEntity == nil {
		return nil, nil
	}

	return envEntity, nil
}

func (uc *useCase) GetEnvsWithBranchId(branchId string) ([]*daemon.Env, error) {
	if branchId == "" {
		return nil, errors.WithMessage(daemon.ErrParamsInvalid, "branchId required")
	}

	return uc.domain.GetEnvsWithBranchId(branchId)
}

func (uc *useCase) DeleteEnv(envId []byte) error {
	if envId == nil {
		return errors.Wrapf(daemon.ErrParamsInvalid, "envId is empty")
	}

	if !uc.domain.IsEnvAvailable(envId) {
		return errors.Wrapf(daemon.ErrEnvNotFound, "envId:%s", envId)
	}

	deleteErr := uc.domain.DeleteEnv(envId)
	if deleteErr != nil {
		return errors.WithStack(deleteErr)
	}

	return nil
}

func (uc *useCase) GetEnvEncToken(envId []byte) ([]byte, error) {
	if envId == nil {
		return nil, errors.Wrapf(daemon.ErrParamsInvalid, "envId is empty")
	}

	envEntity, fetchErr := uc.domain.Env(envId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
	}
	if envEntity == nil {
		return nil, errors.Wrapf(daemon.ErrEnvNotFound, "envId: %s", envId)
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
