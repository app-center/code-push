package portal

import "github.com/funnyecho/code-push/pkg/errors"

const (
	ErrInternalError = errors.Error("INTERNAL_ERROR")
	ErrParamsInvalid = errors.Error("PARAMS_INVALID")
	ErrUnauthorized  = errors.Error("UNAUTHORIZED")
	ErrInvalidToken  = errors.Error("INVALID_TOKEN")
	ErrInvalidBranch = errors.Error("INVALID_BRANCH")
)
