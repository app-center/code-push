package semver

import (
	"fmt"
	"github.com/funnyecho/code-push/pkg/errors"
	semverErrors "github.com/funnyecho/code-push/pkg/semver/errors"
	"regexp"
	"strconv"
	"strings"
)

const (
	CompareLessFlag  = -1
	CompareEqualFlag = 0
	CompareLargeFlag = 1
)

const (
	PRStageAlpha   = 1
	PRStageBeta    = 2
	PRStageRC      = 3
	PRStageRelease = 4
)

type SemVer struct {
	// major.minor.patch-preRelease
	majorV int
	minorV int
	patchV int

	prStage   int
	prVersion int
	prBuild   int
}

func (ver *SemVer) StageSafetyCompare(ver2 interface{}) int {
	if rawVerString, isString := ver2.(string); isString {
		if semVer2, err := ParseVersion(rawVerString); err != nil {
			return ver.stageSafetyCompare(semVer2)
		} else {
			return CompareLessFlag
		}
	}

	switch i := ver2.(type) {
	case string:
		if semVer2, err := ParseVersion(i); err != nil {
			return ver.stageSafetyCompare(semVer2)
		} else {
			return CompareLessFlag
		}
	case *SemVer:
		return ver.stageSafetyCompare(i)
	case SemVer:
		return ver.stageSafetyCompare(&i)
	default:
		return CompareLessFlag
	}
}

func (ver *SemVer) Compare(ver2 interface{}) int {
	if rawVerString, isString := ver2.(string); isString {
		if semVer2, err := ParseVersion(rawVerString); err != nil {
			return ver.compare(semVer2)
		} else {
			return CompareLessFlag
		}
	}

	switch i := ver2.(type) {
	case string:
		if semVer2, err := ParseVersion(i); err != nil {
			return ver.compare(semVer2)
		} else {
			return CompareLessFlag
		}
	case *SemVer:
		return ver.compare(i)
	case SemVer:
		return ver.compare(&i)
	default:
		return CompareLessFlag
	}
}

func (ver *SemVer) stageSafetyCompare(ver2 *SemVer) int {
	mainCompare := ver.compare(ver2)

	switch {
	case mainCompare == CompareLessFlag:
		return CompareLessFlag
	case mainCompare == CompareEqualFlag:
		return ver.prCompare(ver2)
	case mainCompare == CompareLargeFlag:
		fallthrough
	default:
		if ver.prStageCompare(ver2) == CompareLessFlag {
			return CompareLessFlag
		} else {
			return CompareLargeFlag
		}
	}
}

func (ver *SemVer) compare(ver2 *SemVer) int {
	c := ver.majorV - ver2.majorV

	if c > 0 {
		return CompareLargeFlag
	} else if c < 0 {
		return CompareLessFlag
	}

	c = ver.minorV - ver2.minorV

	if c > 0 {
		return CompareLargeFlag
	} else if c < 0 {
		return CompareLessFlag
	}

	c = ver.patchV - ver2.patchV

	if c > 0 {
		return CompareLargeFlag
	} else if c < 0 {
		return CompareLessFlag
	}

	return CompareEqualFlag
}

func (ver *SemVer) prStageCompare(ver2 *SemVer) int {
	c := ver.prStage - ver2.prStage
	switch {
	case c > 0:
		return CompareLargeFlag
	case c == 0:
		return CompareEqualFlag
	case c < 0:
		return CompareLessFlag
	}

	return CompareLessFlag
}

func (ver *SemVer) prCompare(ver2 *SemVer) int {
	c := ver.prStage - ver2.prStage
	switch {
	case c > 0:
		return CompareLargeFlag
	case c < 0:
		return CompareLessFlag
	}

	c = ver.prVersion - ver2.prVersion
	switch {
	case c > 0:
		return CompareLargeFlag
	case c < 0:
		return CompareLessFlag
	}

	c = ver.prBuild - ver2.prBuild
	switch {
	case c > 0:
		return CompareLargeFlag
	case c < 0:
		return CompareLessFlag
	}

	return CompareEqualFlag
}

func (ver *SemVer) String() string {
	return fmt.Sprintf(
		"%v.%v.%v.%1d%02d%1d",
		ver.majorV,
		ver.minorV,
		ver.patchV,
		ver.prStage,
		ver.prVersion,
		ver.prBuild,
	)
}

func (ver *SemVer) PRStage() int {
	return ver.prStage
}

func ParseVersion(rawVer string) (semVer *SemVer, parseErr errors.IError) {
	rawVer = strings.TrimPrefix(rawVer, "v")

	var majorV, minorV, patchV, prStage, prVersion, prBuild int

	re := regexp.MustCompile(`[.-]`)

	matches := re.Split(rawVer, 4)

	if len(matches) < 3 || len(matches) > 4 {
		return nil, &semverErrors.InvalidRawVersionFormatError{RawVersion: rawVer}
	}

	var v int64
	var err error

	v, err = strconv.ParseInt(matches[0], 10, 64)
	if err != nil {
		return nil, &semverErrors.InvalidMajorVersionError{
			RawVersion:   rawVer,
			MajorVersion: matches[0],
		}
	} else {
		majorV = int(v)
	}

	v, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, &semverErrors.InvalidMinorVersionError{
			Err:          err,
			RawVersion:   rawVer,
			MinorVersion: matches[1],
		}
	} else {
		minorV = int(v)
	}

	v, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return nil, &semverErrors.InvalidPatchVersionError{
			Err:          err,
			RawVersion:   rawVer,
			PatchVersion: matches[2],
		}
	} else {
		patchV = int(v)
	}

	if len(matches) == 3 {
		prStage = PRStageRelease
		prVersion = 1
		prBuild = 1
	} else {
		rawPR := strings.ToLower(matches[3])
		if i := strings.Index(rawPR, "."); i > -1 {
			rawPRStage := rawPR[0:i]
			rawPRVersion := rawPR[i+1:]

			switch rawPRStage {
			case "alpha":
				prStage = PRStageAlpha
			case "beta":
				prStage = PRStageBeta
			case "rc":
				prStage = PRStageRC
			case "release":
				prStage = PRStageRelease
			default:
				return nil, &semverErrors.InvalidPreReleaseVersionError{
					RawVersion: rawVer,
					PRStage:    rawPRStage,
				}
			}

			v, err = strconv.ParseInt(rawPRVersion, 10, 32)
			if err != nil {
				return nil, &semverErrors.InvalidPreReleaseVersionError{
					Err:        err,
					RawVersion: rawVer,
					PRVersion:  rawPRVersion,
				}
			} else {
				prVersion = int(v)
				prBuild = 1
			}
		} else {
			v, err = strconv.ParseInt(rawPR, 10, 32)
			if err != nil {
				return nil, &semverErrors.InvalidPreReleaseVersionError{
					Err:        err,
					RawVersion: rawVer,
				}
			}

			prStage = int(v / 1000)
			v = v % 1000

			prVersion = int(v / 10)
			v = v % 10

			prBuild = int(v)
		}
	}

	return New(CtorConfig{
		majorV:    majorV,
		minorV:    minorV,
		patchV:    patchV,
		prStage:   prStage,
		prVersion: prVersion,
		prBuild:   prBuild,
	})
}

type CtorConfig struct {
	majorV int
	minorV int
	patchV int

	prStage   int
	prVersion int
	prBuild   int
}

func (config CtorConfig) ToRawVersion() string {
	return fmt.Sprintf(
		"%v.%v.%v-%v.%v.%v",
		config.majorV,
		config.minorV,
		config.patchV,
		config.prStage,
		config.prVersion,
		config.prBuild,
	)
}

func New(config CtorConfig) (semVer *SemVer, err errors.IError) {
	if config.majorV < 0 ||
		config.minorV < 0 ||
		config.patchV < 0 ||
		(config.prStage < PRStageAlpha || config.prStage > PRStageRelease) ||
		config.prVersion < 0 ||
		config.prBuild < 0 {
		return nil, &semverErrors.InvalidRawVersionFormatError{
			RawVersion: config.ToRawVersion(),
		}
	}

	return &SemVer{
		majorV:    config.majorV,
		minorV:    config.minorV,
		patchV:    config.patchV,
		prStage:   config.prStage,
		prVersion: config.prVersion,
		prBuild:   config.prBuild,
	}, nil
}
