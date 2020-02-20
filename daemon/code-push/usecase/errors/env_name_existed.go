package errors

import "github.com/funnyecho/code-push/pkg/errors"

type EnvNameExistedError error

func ThrowEnvNameExistedError(branchId, envName string) EnvNameExistedError {
	return errors.Throw(errors.CtorConfig{
		Code: FA_ENV_NAME_EXISTED,
		Msg:  "branch name existed",
		Meta: errors.MetaFields{"branchId": branchId, "envName": envName},
	})
}
