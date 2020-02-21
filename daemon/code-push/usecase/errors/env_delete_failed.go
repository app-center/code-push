package errors

import "github.com/funnyecho/code-push/pkg/errors"

type EnvDeleteFailedError error

func ThrowEnvDeleteFailedError(err error, envId string) EnvDeleteFailedError {
	return errors.Throw(errors.CtorConfig{
		Code: FA_ENV_DELETE_FAILED,
		Meta: errors.MetaFields{"err": err, "envId": envId},
	})
}
