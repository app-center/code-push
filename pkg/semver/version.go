package semver

import (
	"fmt"
	"github.com/funnyecho/code-push/pkg/safemath"
	semverErrors "github.com/funnyecho/code-push/pkg/semver/errors"
	"strconv"
	"strings"
)

const (
	CompareLessFlag  = safemath.CompareLessFlag
	CompareEqualFlag = safemath.CompareEqualFlag
	CompareLargeFlag = safemath.CompareLargeFlag
)

const (
	PRStageAlpha   = 1
	PRStageBeta    = 2
	PRStageRC      = 3
	PRStageRelease = 4
)

type SemVer struct {
	majorV uint
	minorV uint
	patchV uint

	prStage   uint8
	prVersion uint8
	prBuild   uint8
}

func (ver *SemVer) StageSafetyStrictCompare(ver2 interface{}) int {
	switch i := ver2.(type) {
	case string:
		if semVer2, err := ParseVersion(i); err == nil {
			return ver.stageSafetyStrictCompare(semVer2)
		} else {
			return CompareLargeFlag
		}
	case *SemVer:
		return ver.stageSafetyStrictCompare(i)
	case SemVer:
		return ver.stageSafetyStrictCompare(&i)
	default:
		return CompareLargeFlag
	}
}

func (ver *SemVer) StageSafetyLooseCompare(ver2 interface{}) int {
	switch i := ver2.(type) {
	case string:
		if semVer2, err := ParseVersion(i); err == nil {
			return ver.stageSafetyStrictCompare(semVer2)
		} else {
			return CompareLargeFlag
		}
	case *SemVer:
		return ver.stageSafetyStrictCompare(i)
	case SemVer:
		return ver.stageSafetyStrictCompare(&i)
	default:
		return CompareLargeFlag
	}
}

func (ver *SemVer) Compare(ver2 interface{}) int {
	switch i := ver2.(type) {
	case string:
		if semVer2, err := ParseVersion(i); err == nil {
			return ver.compare(semVer2)
		} else {
			return CompareLargeFlag
		}
	case *SemVer:
		return ver.compare(i)
	case SemVer:
		return ver.compare(&i)
	default:
		return CompareLargeFlag
	}
}

func (ver *SemVer) stageSafetyStrictCompare(ver2 *SemVer) int {
	mainCompare := ver.compare(ver2)

	switch {
	case mainCompare == CompareLessFlag:
		if ver.prStageCompare(ver2) == CompareLargeFlag {
			return CompareLargeFlag
		} else {
			return CompareLessFlag
		}
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

func (ver *SemVer) stageSafetyLooseCompare(ver2 *SemVer) int {
	mainCompare := ver.compare(ver2)

	switch {
	case mainCompare == CompareLessFlag:
		return CompareLessFlag
	case mainCompare == CompareEqualFlag:
		return ver.prCompare(ver2)
	case mainCompare == CompareLargeFlag:
		fallthrough
	default:
		return CompareLargeFlag
	}
}

func (ver *SemVer) compare(ver2 *SemVer) int {
	c := safemath.UintCompare(ver.majorV, ver2.majorV)

	if c > 0 {
		return CompareLargeFlag
	} else if c < 0 {
		return CompareLessFlag
	}

	c = safemath.UintCompare(ver.minorV, ver2.minorV)

	if c > 0 {
		return CompareLargeFlag
	} else if c < 0 {
		return CompareLessFlag
	}

	c = safemath.UintCompare(ver.patchV, ver2.patchV)

	if c > 0 {
		return CompareLargeFlag
	} else if c < 0 {
		return CompareLessFlag
	}

	return CompareEqualFlag
}

func (ver *SemVer) prStageCompare(ver2 *SemVer) int {
	c := safemath.Uint8Compare(ver.prStage, ver2.prStage)
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
	c := safemath.Uint8Compare(ver.prStage, ver2.prStage)
	switch {
	case c > 0:
		return CompareLargeFlag
	case c < 0:
		return CompareLessFlag
	}

	c = safemath.Uint8Compare(ver.prVersion, ver2.prVersion)
	switch {
	case c > 0:
		return CompareLargeFlag
	case c < 0:
		return CompareLessFlag
	}

	c = safemath.Uint8Compare(ver.prBuild, ver2.prBuild)
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

func (ver *SemVer) MajorVersion() uint {
	return ver.majorV
}

func (ver *SemVer) MinorVersion() uint {
	return ver.minorV
}

func (ver *SemVer) PatchVersion() uint {
	return ver.patchV
}

func (ver *SemVer) PRStage() uint8 {
	return ver.prStage
}

func (ver *SemVer) PRVersion() uint8 {
	return ver.prVersion
}

func (ver *SemVer) PRBuild() uint8 {
	return ver.prBuild
}

func ParseVersion(rawVer string) (semVer *SemVer, parseErr error) {
	rawVer = strings.TrimPrefix(rawVer, "v")

	var majorV, minorV, patchV uint
	var prStage, prVersion, prBuild uint8

	var rawMajorV, rawMinorV, rawPatchV, rawPR string

	parts := strings.SplitN(rawVer, ".", 3)

	if len(parts) < 3 {
		return nil, semverErrors.NewInvalidRawVersionFormatError(semverErrors.InvalidRawVersionFormatErrorConfig{RawVersion: rawVer})
	}

	rawMajorV = parts[0]
	rawMinorV = parts[1]

	var v uint64
	var err error

	v, err = strconv.ParseUint(rawMajorV, 10, 64)
	if err != nil {
		return nil, semverErrors.NewInvalidMajorVersionError(semverErrors.InvalidMajorVersionErrorConfig{
			Err:          err,
			RawVersion:   rawVer,
			MajorVersion: rawMajorV,
		})
	} else {
		majorV = uint(v)
	}

	v, err = strconv.ParseUint(rawMinorV, 10, 64)
	if err != nil {
		return nil, semverErrors.NewInvalidMinorVersionError(semverErrors.InvalidMinorVersionErrorConfig{
			Err:          err,
			RawVersion:   rawVer,
			MinorVersion: rawMinorV,
		})
	} else {
		minorV = uint(v)
	}

	rawPRBound := strings.IndexAny(parts[2], ".-")

	if rawPRBound > 0 {
		rawPatchV = parts[2][0:rawPRBound]
		rawPR = parts[2][rawPRBound:]
	} else {
		rawPatchV = parts[2]
		rawPR = ""
	}

	v, err = strconv.ParseUint(rawPatchV, 10, 64)
	if err != nil {
		return nil, semverErrors.NewInvalidPatchVersionError(semverErrors.InvalidPatchVersionErrorConfig{
			Err:          err,
			RawVersion:   rawVer,
			PatchVersion: rawPatchV,
		})
	} else {
		patchV = uint(v)
	}

	if len(rawPR) > 0 {
		switch rawPR[:1] {
		case ".":
			numericRawPR := rawPR[1:]
			v, err = strconv.ParseUint(numericRawPR, 10, 32)
			if err != nil {
				return nil, semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{
					Err:        err,
					RawVersion: rawVer,
					RawPR:      rawPR,
				})
			}

			prStage = uint8(v / 1000)
			v = v % 1000

			prVersion = uint8(v / 10)
			v = v % 10

			prBuild = uint8(v)
		case "-":
			strRawPR := rawPR[1:]

			versionBound := strings.Index(strRawPR, ".")

			if versionBound <= 0 || versionBound == len(strRawPR)-1 {
				return nil, semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{
					Err:        err,
					RawVersion: rawVer,
					RawPR:      rawPR,
				})
			}

			rawPRStage := strRawPR[0:versionBound]

			switch strings.ToLower(rawPRStage) {
			case "alpha":
				prStage = PRStageAlpha
			case "beta":
				prStage = PRStageBeta
			case "rc":
				prStage = PRStageRC
			case "release":
				prStage = PRStageRelease
			default:
				return nil, semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{
					RawVersion: rawVer,
					RawPR:      strRawPR,
					PRStage:    rawPRStage,
				})
			}

			rawPRLeft := strRawPR[versionBound+1:]
			buildBound := strings.Index(rawPRLeft, "+")

			var rawPRVersion string

			if buildBound >= 0 {
				rawPRVersion = rawPRLeft[0:buildBound]
				rawPRBuild := rawPRLeft[buildBound+1:]

				v, err = strconv.ParseUint(rawPRBuild, 10, 32)
				if err != nil {
					return nil, semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{
						Err:        err,
						RawVersion: rawVer,
						PRVersion:  rawPRVersion,
					})
				} else {
					prBuild = uint8(v)
				}

			} else {
				rawPRVersion = rawPRLeft
				prBuild = 1
			}

			v, err = strconv.ParseUint(rawPRVersion, 10, 32)
			if err != nil {
				return nil, semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{
					Err:        err,
					RawVersion: rawVer,
					PRVersion:  rawPRVersion,
				})
			} else {
				prVersion = uint8(v)
			}
		default:
			return nil, semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{
				RawVersion: rawVer,
				RawPR:      rawPR,
			})
		}
	} else {
		prStage = PRStageRelease
		prVersion = 1
		prBuild = 1
	}

	return New(CtorConfig{
		MajorV:    majorV,
		MinorV:    minorV,
		PatchV:    patchV,
		PRStage:   prStage,
		PRVersion: prVersion,
		PRBuild:   prBuild,
	})
}

type CtorConfig struct {
	MajorV uint
	MinorV uint
	PatchV uint

	PRStage   uint8
	PRVersion uint8
	PRBuild   uint8
}

func (config CtorConfig) ToRawVersion() string {
	return fmt.Sprintf(
		"%v.%v.%v-%v.%v+%v",
		config.MajorV,
		config.MinorV,
		config.PatchV,
		config.PRStage,
		config.PRVersion,
		config.PRBuild,
	)
}

func New(config CtorConfig) (semVer *SemVer, err error) {
	rawVersion := config.ToRawVersion()

	switch true {
	case config.MajorV < 0:
		err = semverErrors.NewInvalidMajorVersionError(semverErrors.InvalidMajorVersionErrorConfig{
			Err:          nil,
			RawVersion:   rawVersion,
			MajorVersion: config.MajorV,
		})
	case config.MinorV < 0:
		err = semverErrors.NewInvalidMinorVersionError(semverErrors.InvalidMinorVersionErrorConfig{
			Err:          nil,
			RawVersion:   rawVersion,
			MinorVersion: config.MinorV,
		})
	case config.PatchV < 0:
		err = semverErrors.NewInvalidPatchVersionError(semverErrors.InvalidPatchVersionErrorConfig{
			Err:          nil,
			RawVersion:   rawVersion,
			PatchVersion: config.PatchV,
		})
	case config.MajorV == 0 && config.MinorV == 0 && config.PatchV == 0:
		err = semverErrors.NewInvalidRawVersionFormatError(semverErrors.InvalidRawVersionFormatErrorConfig{
			RawVersion: rawVersion,
		})
	case (config.PRStage < PRStageAlpha || config.PRStage > PRStageRelease) ||
		(config.PRVersion < 1 || config.PRVersion > 99) ||
		(config.PRBuild < 1 || config.PRBuild > 9):
		err = semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{
			Err:        nil,
			RawVersion: rawVersion,
			PRStage:    config.PRStage,
			PRVersion:  config.PRVersion,
			PRBuild:    config.PRBuild,
		})
	default:
		semVer = &SemVer{
			majorV:    config.MajorV,
			minorV:    config.MinorV,
			patchV:    config.PatchV,
			prStage:   config.PRStage,
			prVersion: config.PRVersion,
			prBuild:   config.PRBuild,
		}
	}

	return
}
