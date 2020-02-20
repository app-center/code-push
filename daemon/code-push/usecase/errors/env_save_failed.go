package errors

import "github.com/funnyecho/code-push/pkg/errors"

type EnvSaveError error

func ThrowEnvSaveError(err error, env interface{}) EnvSaveError {
	return errors.Throw(errors.CtorConfig{
		Error: err,
		Code:  FA_ENV_SAVE_FAILED,
		Meta:  MetaFields{"env": env},
	})
}
