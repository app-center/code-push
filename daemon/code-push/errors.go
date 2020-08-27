package code_push

import "github.com/funnyecho/code-push/pkg/errors"

const (
	ErrInternalError = errors.Error("INTERNAL_ERROR")
	ErrParamsInvalid = errors.Error("PARAMS_INVALID")

	ErrBranchExisted     = errors.Error("BRANCH_EXISTED")
	ErrBranchNameExisted = errors.Error("BRANCH_NAME_EXISTED")
	ErrBranchNotFound    = errors.Error("BRANCH_NOT_FOUND")

	ErrEnvExisted     = errors.Error("ENV_EXISTED")
	ErrEnvNameExisted = errors.Error("ENV_NAME_EXISTED")
	ErrEnvNotFound    = errors.Error("ENV_NOT_FOUND")

	ErrVersionExisted  = errors.Error("VERSION_EXISTED")
	ErrVersionNotFound = errors.Error("VERSION_NOT_FOUND")
)
