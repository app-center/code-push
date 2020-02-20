package errors

import "github.com/funnyecho/code-push/pkg/errors"

type EnvInvalidEncTokenError error

func ThrowEnvInvalidEncTokenError(err error) EnvInvalidEncTokenError {
	return errors.Throw(errors.CtorConfig{
		Error: err,
		Code:  FA_ENV_INVALID_ENC_TOKEN,
	})
}
