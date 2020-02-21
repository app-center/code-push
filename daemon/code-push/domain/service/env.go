package service

import (
	"github.com/funnyecho/code-push/daemon/code-push/domain/repository"
)

type IEnvService interface {
	IsEnvExisted(envId string) bool
	IsEnvNameExisted(branchId, envName string) bool
}

type envService struct {
	envRepo repository.IEnv
}

func (e envService) IsEnvExisted(envId string) bool {
	env, err := e.envRepo.FirstEnv(envId)

	return err != nil && env != nil
}

func (e envService) IsEnvNameExisted(branchId, envName string) bool {
	env, err := e.envRepo.FirstEnvWithBranchIdAndEnvName(branchId, envName)

	return err != nil && env != nil
}

type EnvServiceConfig struct {
	EnvRepo repository.IEnv
}

func NewEnvService(config EnvServiceConfig) IEnvService {
	return &envService{
		envRepo: config.EnvRepo,
	}
}
