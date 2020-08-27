package client

import "github.com/funnyecho/code-push/pkg/errors"

const (
	ErrInternalError        = errors.Error("INTERNAL_ERROR")
	ErrParamsInvalid        = errors.Error("PARAMS_INVALID")
	ErrUnauthorized         = errors.Error("UNAUTHORIZED")
	ErrInvalidToken         = errors.Error("INVALID_TOKEN")
	ErrInvalidEnv           = errors.Error("INVALID_ENV")
	ErrVersionNotFound      = errors.Error("VERSION_NOT_FOUND")
	ErrVersionNotUpgradable = errors.Error("VERSION_NOT_UPGRADABLE")
)
