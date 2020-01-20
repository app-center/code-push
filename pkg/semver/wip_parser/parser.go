package wip_parser

import (
	"github.com/funnyecho/code-push/pkg/semver"
	semverErrors "github.com/funnyecho/code-push/pkg/semver/errors"
	"strconv"
)

func parseVersion(rawVer string) (semVer *semver.SemVer, parseErr error) {
	var majorV, minorV, patchV uint
	var prStage, prVersion, prBuild uint8

	var v uint64

	invalidRawVersionErrorReporter := invalidRawVersionErrorReporterFactory(rawVer)
	invalidPreReleaseVersionErrorReporter := invalidPreReleaseVersionErrorReporterFactory(rawVer)

	semverSegment := newChainSegment(chainSegmentConfig{
		reporter: invalidRawVersionErrorReporter,
		chains: []iSegment{
			newSegment(segmentConfig{
				walker: charWalkerFactory('v'),
				reporter: func(err error, chars []rune) error {
					return nil
				},
			}),
			newSegment(segmentConfig{
				walker: digitWalker,
				reporter: func(err error, chars []rune) error {
					if err == nil {
						v, err = strconv.ParseUint(string(chars), 10, 64)
					}

					if err != nil {
						return semverErrors.NewInvalidMajorVersionError(semverErrors.InvalidMajorVersionErrorConfig{
							Err:          err,
							RawVersion:   rawVer,
							MajorVersion: chars,
						})
					} else {
						majorV = uint(v)
						return nil
					}
				},
			}),
			newSegment(segmentConfig{
				walker: lengthLimitedWalkerFactory(1, dotWalker),
				reporter: func(err error, chars []rune) error {
					if err != nil || (len(chars) != 1 || !dotWalker(chars[0])) {
						return semverErrors.NewInvalidRawVersionFormatError(semverErrors.InvalidRawVersionFormatErrorConfig{RawVersion: rawVer})
					} else {
						return nil
					}
				},
			}),
			newSegment(segmentConfig{
				walker: digitWalker,
				reporter: func(err error, chars []rune) error {
					if err == nil {
						v, err = strconv.ParseUint(string(chars), 10, 64)
					}

					if err != nil {
						return semverErrors.NewInvalidMinorVersionError(semverErrors.InvalidMinorVersionErrorConfig{
							Err:          err,
							RawVersion:   rawVer,
							MinorVersion: chars,
						})
					} else {
						minorV = uint(v)
						return nil
					}
				},
			}),
			newSegment(segmentConfig{
				walker: lengthLimitedWalkerFactory(1, dotWalker),
				reporter: func(err error, chars []rune) error {
					if err != nil || (len(chars) != 1 || !dotWalker(chars[0])) {
						return semverErrors.NewInvalidRawVersionFormatError(semverErrors.InvalidRawVersionFormatErrorConfig{RawVersion: rawVer})
					} else {
						return nil
					}
				},
			}),
			newSegment(segmentConfig{
				walker: digitWalker,
				reporter: func(err error, chars []rune) error {
					if err == nil {
						v, err = strconv.ParseUint(string(chars), 10, 64)
					}

					if err != nil {
						return semverErrors.NewInvalidPatchVersionError(semverErrors.InvalidPatchVersionErrorConfig{
							Err:          err,
							RawVersion:   rawVer,
							PatchVersion: chars,
						})
					} else {
						patchV = uint(v)
						return nil
					}
				},
			}),
			newSwitchSegment(switchSegmentConfig{
				cases: func(c rune) iSegment {
					switch c {
					case '.':
						return newChainSegment(chainSegmentConfig{
							reporter: invalidPreReleaseVersionErrorReporter,
							chains: []iSegment{
								newSegment(segmentConfig{
									walker: digitWithLenWalkerFactory(1),
									reporter: func(err error, chars []rune) error {
										if err != nil {
											return invalidRawVersionErrorReporter(err, chars)
										} else {
											return nil
										}
									},
								}),
								newSegment(segmentConfig{
									walker: digitWithLenWalkerFactory(2),
									reporter: func(err error, chars []rune) error {
										if err != nil {
											return invalidRawVersionErrorReporter(err, chars)
										} else {
											return nil
										}
									},
								}),
								newSegment(segmentConfig{
									walker: digitWithLenWalkerFactory(1),
									reporter: func(err error, chars []rune) error {
										if err != nil {
											return invalidRawVersionErrorReporter(err, chars)
										} else {
											return nil
										}
									},
								}),
							},
						})
					case '-':
						return newChainSegment(chainSegmentConfig{
							reporter: invalidPreReleaseVersionErrorReporter,
							chains: []iSegment{
								newSegment(segmentConfig{
									walker: stringCandidateWalkerFactory([]string{"release", "rc", "beta", "alpha"}),
									reporter: func(err error, chars []rune) error {
										if err != nil {
											return invalidPreReleaseVersionErrorReporter(err, chars)
										} else {
											return nil
										}
									},
								}),
								newSegment(segmentConfig{
									walker: dotWalker,
									reporter: func(err error, chars []rune) error {
										if err != nil {
											return invalidPreReleaseVersionErrorReporter(err, chars)
										} else {
											return nil
										}
									},
								}),
								newSegment(segmentConfig{
									walker: digitWithLenWalkerFactory(2),
									reporter: func(err error, chars []rune) error {
										if err != nil {
											return invalidPreReleaseVersionErrorReporter(err, chars)
										} else {
											return nil
										}
									},
								}),
								newSwitchSegment(switchSegmentConfig{
									cases: func(c rune) iSegment {
										switch c {
										case 0x00:
											return newImmediateSegment(immediateSegmentConfig{
												reporter: func(err error, chars []rune) error {
													return nil
												},
											})
										case '+':
											return newSegment(segmentConfig{
												walker: digitWithLenWalkerFactory(1),
												reporter: func(err error, chars []rune) error {
													if err != nil {
														return invalidPreReleaseVersionErrorReporter(err, chars)
													} else {
														return nil
													}
												},
											})
										default:
											return nil
										}
									},
								}),
							},
						})
					case 0x00:
						return newImmediateSegment(immediateSegmentConfig{
							reporter: func(err error, chars []rune) error {
								prStage = semver.PRStageRelease
								prVersion = 1
								prBuild = 1

								return nil
							},
						})
					default:
						return nil
					}

				},
				reporter: invalidPreReleaseVersionErrorReporter,
			}),
		},
	})

	for _, c := range rawVer {
		hit := semverSegment.walk(c)
		if !hit {
			reportErr := semverSegment.report()
			return nil, reportErr
		}
	}
	semverSegment.walk(0x00)

	if reportErr := semverSegment.report(); reportErr != nil {
		return nil, reportErr
	}

	return semver.New(semver.CtorConfig{
		MajorV:    majorV,
		MinorV:    minorV,
		PatchV:    patchV,
		PRStage:   prStage,
		PRVersion: prVersion,
		PRBuild:   prBuild,
	})

	return
}
