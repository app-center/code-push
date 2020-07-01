package usecase

import (
	code_push "github.com/funnyecho/code-push/daemon/code-push"
	"github.com/funnyecho/code-push/daemon/code-push/domain"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

func NewEnvUseCase(config EnvUseCaseConfig) (IEnv, error) {
	if config.BranchService == nil ||
		config.EnvService == nil {
		panic("invalid env use case params")
	}

	return &envUseCase{
		envService:    config.EnvService,
		branchService: config.BranchService,
	}, nil
}

type Env struct {
	BranchId   string
	EnvId      string
	Name       string
	CreateTime time.Time
}

func toEnv(env *domain.Env) *Env {
	return &Env{
		BranchId:   env.BranchId,
		EnvId:      env.ID,
		Name:       env.Name,
		CreateTime: env.CreateTime,
	}
}

type IEnv interface {
	CreateEnv(branchId, envName string) (*Env, error)
	GetEnv(envId string) (*Env, error)
	DeleteEnv(envId string) error
	GetEnvEncToken(envId string) (string, error)
	GetEnvAuthHost(envId string) (string, error)
}

type envUseCase struct {
	branchService domain.IBranchService
	envService    domain.IEnvService
}

func (e envUseCase) CreateEnv(branchId, envName string) (*Env, error) {
	if len(branchId) == 0 || len(envName) == 0 {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "branchId or envName can't not be empty")
	}

	envId := generateEnvId(branchId)
	encToken, encTokenErr := generateEnvEncToken()

	if encTokenErr != nil {
		return nil, errors.Wrapf(encTokenErr, "generate env enc token failed")
	}

	envToCreate := &domain.Env{
		BranchId: branchId,
		ID:       envId,
		Name:     envName,
		EncToken: encToken,
	}

	createErr := e.envService.CreateEnv(envToCreate)
	if createErr != nil {
		return nil, errors.WithStack(createErr)
	}

	return toEnv(envToCreate), nil
}

func (e envUseCase) GetEnv(envId string) (*Env, error) {
	if len(envId) == 0 {
		return nil, errors.Wrapf(code_push.ErrParamsInvalid, "envId is empty")
	}

	envEntity, fetchErr := e.envService.Env(envId)
	if fetchErr != nil {
		return nil, errors.WithStack(fetchErr)
	}
	if envEntity == nil {
		return nil, errors.WithMessagef(code_push.ErrEnvNotFound, "envId: %s", envId)
	}

	return toEnv(envEntity), nil
}

func (e envUseCase) DeleteEnv(envId string) error {
	if len(envId) == 0 {
		return errors.Wrapf(code_push.ErrParamsInvalid, "envId is empty")
	}

	deleteErr := e.envService.DeleteEnv(envId)
	if deleteErr != nil {
		return errors.WithStack(deleteErr)
	}

	return nil
}

func (e envUseCase) GetEnvEncToken(envId string) (string, error) {
	if len(envId) == 0 {
		return "", errors.Wrapf(code_push.ErrParamsInvalid, "envId is empty")
	}

	envEntity, fetchErr := e.envService.Env(envId)
	if fetchErr != nil {
		return "", errors.WithStack(fetchErr)
	}
	if envEntity == nil {
		return "", errors.Wrapf(code_push.ErrEnvNotFound, "envId: %s", envId)
	}

	return envEntity.EncToken, nil
}

func (e envUseCase) GetEnvAuthHost(envId string) (string, error) {
	if len(envId) == 0 {
		return "", errors.Wrapf(code_push.ErrParamsInvalid, "envId is empty")
	}

	envEntity, fetchErr := e.envService.Env(envId)
	if fetchErr != nil {
		return "", errors.WithStack(fetchErr)
	}
	if envEntity == nil {
		return "", errors.Wrapf(code_push.ErrEnvNotFound, "envId: %s", envId)
	}

	branchEntity, branchFetchErr := e.branchService.Branch(envEntity.BranchId)
	if branchFetchErr != nil {
		return "", errors.WithStack(fetchErr)
	}
	if branchEntity == nil {
		return "", errors.Wrapf(code_push.ErrBranchNotFound, "envId: %s, branchId: %s", envId, envEntity.BranchId)
	}

	return branchEntity.EncToken, nil
}

type EnvUseCaseConfig struct {
	BranchService domain.IBranchService
	EnvService    domain.IEnvService
}

func generateEnvEncToken() (string, error) {
	token, err := util.RandomPass(16, 8, 0, false, true)

	return token, err
}

func generateEnvId(branchId string) string {
	return util.EncodeMD5(branchId + "/" + uuid.NewV4().String())
}
