package usecase

import "github.com/funnyecho/code-push/gateway/portal"

func (u *useCase) CreateEnv(branchId, envName []byte) (*portal.Env, error) {
	return u.codePush.CreateEnv(branchId, envName)
}

func (u *useCase) GetEnv(envId []byte) (*portal.Env, error) {
	return u.codePush.GetEnv(envId)
}

func (u *useCase) DeleteEnv(envId []byte) error {
	return u.codePush.DeleteEnv(envId)
}

func (u *useCase) GetEnvEncToken(envId []byte) ([]byte, error) {
	return u.codePush.GetEnvEncToken(envId)
}
