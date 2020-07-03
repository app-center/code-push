package filer

import "github.com/funnyecho/code-push/pkg/errors"

const (
	ErrInternalError = errors.Error("FA_INTERNAL_ERROR")
	ErrParamsInvalid = errors.Error("FA_PARAMS_INVALID")

	ErrInvalidFileKey   = errors.Error("FA_INVALID_FILE_KEY")
	ErrInvalidFileValue = errors.Error("FA_INVALID_FILE_VALUE")

	ErrFileKeyNotFound = errors.Error("FA_FILE_KEY_NOT_FOUND")
	ErrFileKeyExisted  = errors.Error("FA_FILE_KEY_EXISTED")

	ErrInvalidAliOssEndpoint        = errors.Error("FA_INVALID_ALI_OSS_ENDPOINT")
	ErrInvalidAliOssAccessKeyId     = errors.Error("FA_INVALID_ALI_OSS_ACCESS_KEY_ID")
	ErrInvalidAliOssAccessKeySecret = errors.Error("FA_INVALID_ALI_OSS_ACCESS_KEY_SECRET")
)
