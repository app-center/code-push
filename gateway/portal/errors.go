package portal

import "github.com/funnyecho/code-push/pkg/errors"

const (
	ErrInternalError = errors.Error("FA_INTERNAL_ERROR")
	ErrParamsInvalid = errors.Error("FA_PARAMS_INVALID")
	ErrUnauthorized  = errors.Error("FA_UNAUTHORIZED")
	ErrInvalidToken  = errors.Error("FA_INVALID_TOKEN")
	ErrInvalidBranch = errors.Error("FA_INVALID_BRANCH")
)
