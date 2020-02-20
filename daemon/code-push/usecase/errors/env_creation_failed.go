package errors

import "github.com/funnyecho/code-push/pkg/errors"

type EnvCreationFailedError error

func ThrowEnvCreationFailedError(err error, code string, env interface{}) EnvCreationFailedError {
	if len(code) == 0 {
		code = FA_ENV_CREATION_FAILED
	}

	return errors.Throw(errors.CtorConfig{
		Code: code,
		Msg:  "env creation failed",
		Meta: errors.MetaFields{"err": err, "env": env},
	})
}
