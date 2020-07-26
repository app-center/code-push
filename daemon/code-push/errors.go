package code_push

import "github.com/funnyecho/code-push/pkg/errors"

const (
	ErrInternalError = errors.Error("FA_INTERNAL_ERROR")
	ErrParamsInvalid = errors.Error("FA_PARAMS_INVALID")

	ErrBranchExisted     = errors.Error("FA_BRANCH_EXISTED")
	ErrBranchNameExisted = errors.Error("FA_BRANCH_NAME_EXISTED")
	ErrBranchNotFound    = errors.Error("FA_BRANCH_NOT_FOUND")

	ErrEnvExisted     = errors.Error("FA_ENV_EXISTED")
	ErrEnvNameExisted = errors.Error("FA_ENV_NAME_EXISTED")
	ErrEnvNotFound    = errors.Error("FA_ENV_NOT_FOUND")

	ErrVersionExisted  = errors.Error("FA_VERSION_EXISTED")
	ErrVersionNotFound = errors.Error("FA_VERSION_NOT_FOUND")
)
