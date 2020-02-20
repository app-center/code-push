package usecase

type IEnvUpdateParams interface {
	SetEnvName(name string) IEnvUpdateParams

	EnvName() (set bool, val string)
}

type envUpdateParams struct {
	envName    string
	envNameSet bool
}

func (e *envUpdateParams) EnvName() (set bool, val string) {
	return e.envNameSet, e.envName
}

func (e *envUpdateParams) SetEnvName(envName string) IEnvUpdateParams {
	e.envNameSet = true
	e.envName = envName

	return e
}

func NewEnvUpdateParams() IEnvUpdateParams {
	return &envUpdateParams{}
}
