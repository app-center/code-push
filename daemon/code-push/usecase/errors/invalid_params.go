package errors

import "github.com/funnyecho/code-push/pkg/errors"

type InvalidParamsError error

type InvalidParamsErrorConfig struct {
	Err    error
	Params MetaFields
}

func ThrowInvalidParamsError(config InvalidParamsErrorConfig) InvalidParamsError {
	return errors.Throw(errors.CtorConfig{
		Error: config.Err,
		Code:  FA_INVALID_PARAMS,
		Meta:  config.Params,
	})
}
