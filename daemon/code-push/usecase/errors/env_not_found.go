package errors

import "github.com/funnyecho/code-push/pkg/errors"

type EnvNotFoundError error

func ThrowEnvNotFoundError(envId string, err error) EnvNotFoundError {
	return errors.Throw(errors.CtorConfig{
		Error: err,
		Code:  FA_ENV_NOT_FOUND,
		Msg:   "env not found",
		Meta:  errors.MetaFields{"envId": envId},
	})
}
