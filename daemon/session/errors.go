package session

import "github.com/funnyecho/code-push/pkg/errors"

const (
	ErrInternalError = errors.Error("FA_INTERNAL_ERROR")
	ErrParamsInvalid = errors.Error("FA_PARAMS_INVALID")

	ErrAccessTokenInvalid = errors.Error("FA_ACCESS_TOKEN_INVALID")
)
