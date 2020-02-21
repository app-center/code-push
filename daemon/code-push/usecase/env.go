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

type Env struct {
	BranchId   string
	EnvId      string
	Name       string
	CreateTime time.Time
}

func toEnv(env *model.Env) *Env {
	return &Env{
		BranchId:   env.BranchId(),
		EnvId:      env.Id(),
		Name:       env.Name(),
		CreateTime: env.CreateTime(),
	}
}

type IEnvUserCase interface {
	CreateEnv(branchId, envName string) (*Env, error)
	UpdateEnv(envId string, params IEnvUpdateParams) error
	GetEnv(envId string) (*Env, error)
	DeleteEnv(envId string) error
	GetEnvEncToken(envId string) (string, error)
	GetEnvAuthHost(envId string) (string, error)
}

type envUseCase struct {
	envRepo    repository.IEnv
	envService service.IEnvService

	branchRepo    repository.IBranch
	branchService service.IBranchService
}

func (e envUseCase) CreateEnv(branchId, envName string) (*Env, error) {
	if len(branchId) == 0 || len(envName) == 0 {
		return nil, errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Params: errors.MetaFields{
				"branchId": branchId,
				"envName":  envName,
			},
		})
	}

	if e.branchService.IsBranchExisted(branchId) {
		return nil, errors.ThrowBranchNotFoundError(branchId, nil)
	}

	if e.envService.IsEnvNameExisted(branchId, envName) {
		return nil, errors.ThrowEnvNameExistedError(branchId, envName)
	}

	envId := generateEnvId(branchId)
	encToken, encTokenErr := generateEnvEncToken()

	if encTokenErr != nil {
		return nil, errors.ThrowEnvInvalidEncTokenError(encTokenErr)
	}

	if e.envService.IsEnvExisted(envId) {
		return nil, errors.ThrowEnvCreationFailedError(nil, errors.FA_ENV_EXISTED, errors.MetaFields{
			"branchId": branchId,
			"envName":  envName,
		})
	}

	envToCreate := model.NewEnv(model.EnvConfig{
		BranchId:   branchId,
		Id:         envId,
		Name:       envName,
		EncToken:   encToken,
		CreateTime: time.Now(),
	})

	envCreated, createErr := e.envRepo.SaveEnv(envToCreate)
	if createErr != nil {
		return nil, errors.ThrowEnvCreationFailedError(createErr, errors.FA_ENV_CREATION_FAILED, errors.MetaFields{
			"branchId": branchId,
			"envName":  envName,
		})
	}

	return toEnv(envCreated), nil
}

func (e envUseCase) UpdateEnv(envId string, params IEnvUpdateParams) error {
	if len(envId) == 0 {
		return errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "envId is empty",
			Params: errors.MetaFields{
				"envId":  envId,
				"params": params,
			},
		})
	}

	updateEnvName, newEnvName := params.EnvName()

	if !updateEnvName {
		return errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "update params do not have valid keys",
		})
	}

	hasInvalidParamsErr := false
	invalidErrParams := errors.MetaFields{}

	if updateEnvName && len(newEnvName) == 0 {
		hasInvalidParamsErr = true
		invalidErrParams["envName"] = newEnvName
	}

	if hasInvalidParamsErr {
		return errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg:    "invalid update params",
			Params: invalidErrParams,
		})
	}

	entity, findErr := e.envRepo.FirstEnv(envId)

	if findErr != nil {
		return errors.ThrowEnvNotFoundError(envId, findErr)
	}

	if e.envService.IsEnvNameExisted(entity.BranchId(), newEnvName) {
		return errors.ThrowEnvNameExistedError(entity.BranchId(), newEnvName)
	}

	if updateEnvName {
		entity.SetName(newEnvName)
	}

	_, updateErr := e.envRepo.SaveEnv(*entity)
	if updateErr != nil {
		return errors.ThrowEnvSaveError(updateErr, params)
	}

	return nil
}

func (e envUseCase) GetEnv(envId string) (*Env, error) {
	if len(envId) == 0 {
		return nil, errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "envId is empty",
			Params: errors.MetaFields{
				"envId": envId,
			},
		})
	}

	envEntity, fetchErr := e.envRepo.FirstEnv(envId)
	if fetchErr != nil {
		return nil, errors.ThrowEnvNotFoundError(envId, fetchErr)
	}

	return toEnv(envEntity), nil
}

func (e envUseCase) DeleteEnv(envId string) error {
	if len(envId) == 0 {
		return errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "envId is empty",
			Params: errors.MetaFields{
				"envId": envId,
			},
		})
	}

	deleteErr := e.envRepo.DeleteEnv(envId)
	if deleteErr != nil {
		return errors.ThrowEnvDeleteFailedError(deleteErr, envId)
	}

	return nil
}

func (e envUseCase) GetEnvEncToken(envId string) (string, error) {
	if len(envId) == 0 {
		return "", errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "envId is empty",
			Params: errors.MetaFields{
				"envId": envId,
			},
		})
	}

	envEntity, fetchErr := e.envRepo.FirstEnv(envId)
	if fetchErr != nil {
		return "", errors.ThrowEnvNotFoundError(envId, fetchErr)
	}

	return envEntity.EncToken(), nil
}

func (e envUseCase) GetEnvAuthHost(envId string) (string, error) {
	if len(envId) == 0 {
		return "", errors.ThrowInvalidParamsError(errors.InvalidParamsErrorConfig{
			Msg: "envId is empty",
			Params: errors.MetaFields{
				"envId": envId,
			},
		})
	}

	envEntity, fetchErr := e.envRepo.FirstEnv(envId)
	if fetchErr != nil {
		return "", errors.ThrowEnvNotFoundError(envId, fetchErr)
	}

	branchEntity, branchFetchErr := e.branchRepo.FirstBranch(envEntity.BranchId())
	if branchFetchErr != nil {
		return "", errors.ThrowBranchNotFoundError(envEntity.BranchId(), branchFetchErr)
	}

	return branchEntity.BranchAuthHost(), nil
}

type EnvUseCaseConfig struct {
	EnvRepo    repository.IEnv
	EnvService service.IEnvService

	BranchRepo    repository.IBranch
	BranchService service.IBranchService
}

func NewEnvUseCase(config EnvUseCaseConfig) IEnvUserCase {
	return &envUseCase{
		envRepo:    config.EnvRepo,
		envService: config.EnvService,

		branchRepo:    config.BranchRepo,
		branchService: config.BranchService,
	}
}

func generateEnvEncToken() (string, error) {
	token, err := util.RandomPass(16, 8, 0, false, true)

	return token, err
}

func generateEnvId(branchId string) string {
	return util.EncodeMD5(branchId + "/" + uuid.NewV4().String())
}
