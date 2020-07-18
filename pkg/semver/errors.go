package semver

import "github.com/funnyecho/code-push/pkg/errors"

const (
	ErrInvalidMajorVersion      = errors.Error("FA_INVALID_MAJOR_VERSION")
	ErrInvalidMinorVersion      = errors.Error("FA_INVALID_MINOR_VERSION")
	ErrInvalidPatchVersion      = errors.Error("FA_INVALID_PATCH_VERSION")
	ErrInvalidPreReleaseVersion = errors.Error("FA_INVALID_PRERELEASE_VERSION")
	ErrInvalidVersionFormat     = errors.Error("FA_INVALID_VERSION_FORMAT")
)
