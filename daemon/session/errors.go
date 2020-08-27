package session

import "github.com/funnyecho/code-push/pkg/errors"

const (
	ErrInternalError = errors.Error("INTERNAL_ERROR")
	ErrParamsInvalid = errors.Error("PARAMS_INVALID")

	ErrAccessTokenInvalid = errors.Error("ACCESS_TOKEN_INVALID")
)
