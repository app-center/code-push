package wip_parser

import (
	"fmt"
	semverErrors "github.com/funnyecho/code-push/pkg/semver/errors"
)

type segReporter func(err error, chars []rune) error

func noopSegReporter (err error, chars []rune) error {
	if err != nil {
		return err
	}

	return nil
}

func noopFailedSegReporter(err error, chars []rune) error {
	if err != nil {
		return err
	}

	return fmt.Errorf("auto failed segment")
}

func invalidRawVersionErrorReporterFactory(rawVer string) segReporter {
	return func(err error, chars []rune) error {
		if err == nil {
			return nil
		} else {
			return semverErrors.NewInvalidRawVersionFormatError(semverErrors.InvalidRawVersionFormatErrorConfig{RawVersion: rawVer})
		}
	}
}

func invalidPreReleaseVersionErrorReporterFactory(rawVer string) segReporter {
	return func(err error, chars []rune) error {
		if err == nil {
			return nil
		} else {
			return semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{RawVersion: rawVer})
		}
	}
}