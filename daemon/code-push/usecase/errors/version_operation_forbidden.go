package errors

import "github.com/funnyecho/code-push/pkg/errors"

type VersionOperationForbiddenError error

func ThrowVersionOperationForbiddenError(err error, msg string, params MetaFields) VersionOperationForbiddenError {
	return errors.Throw(errors.CtorConfig{
		Error: err,
		Code:  FA_VERSION_OPERATION_FORBIDDEN,
		Msg:   msg,
		Meta:  params,
	})
}
