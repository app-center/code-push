package filer

import "github.com/funnyecho/code-push/pkg/errors"

const (
	ErrInternalError = errors.Error("INTERNAL_ERROR")
	ErrParamsInvalid = errors.Error("PARAMS_INVALID")

	ErrInvalidFileKey   = errors.Error("INVALID_FILE_KEY")
	ErrInvalidFileValue = errors.Error("INVALID_FILE_VALUE")

	ErrFileKeyNotFound = errors.Error("FILE_KEY_NOT_FOUND")
	ErrFileKeyExisted  = errors.Error("FILE_KEY_EXISTED")

	ErrInvalidAliOssEndpoint        = errors.Error("INVALID_ALI_OSS_ENDPOINT")
	ErrInvalidAliOssAccessKeyId     = errors.Error("INVALID_ALI_OSS_ACCESS_KEY_ID")
	ErrInvalidAliOssAccessKeySecret = errors.Error("INVALID_ALI_OSS_ACCESS_KEY_SECRET")
)
