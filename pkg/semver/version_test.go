package semver

import (
	semverErrors "github.com/funnyecho/code-push/pkg/semver/errors"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestParseVersion(t *testing.T) {
	type InvalidVersionInputs = []struct {
		message   string
		input     []string
		errorType error
	}

	type ValidVersionInputs = map[string][]string

	t.Run("majorV.minorV.patchV", func(t *testing.T) {
		invalidTests := InvalidVersionInputs{
			{
				"majorV shall be a natural number",
				[]string{
					"-1.2.3",
					"major.2.3",
					".2.3",
				},
				semverErrors.NewInvalidMajorVersionError(semverErrors.InvalidMajorVersionErrorConfig{}),
			},
			{
				"minorV shall be a natural number",
				[]string{
					"1.-2.3",
					"1.minor.3",
					"1..3",
				},
				semverErrors.NewInvalidMinorVersionError(semverErrors.InvalidMinorVersionErrorConfig{}),
			},
			{
				"patchV shall be a natural number",
				[]string{
					"1.2.-3",
					"1.2.patch",
					"1.2.",
				},
				semverErrors.NewInvalidPatchVersionError(semverErrors.InvalidPatchVersionErrorConfig{}),
			},
			{
				"0.0.0 is not meaningless",
				[]string{
					"0.0.0",
				},
				semverErrors.NewInvalidRawVersionFormatError(semverErrors.InvalidRawVersionFormatErrorConfig{}),
			},
			{
				"divider must be a dot",
				[]string{
					"1,2,3",
					"1-2-3",
					"1_2_3",
				},
				semverErrors.NewInvalidRawVersionFormatError(semverErrors.InvalidRawVersionFormatErrorConfig{}),
			},
		}

		for _, data := range invalidTests {
			for _, rawVersion := range data.input {
				_, err := ParseVersion(rawVersion)
				assert.IsType(t, data.errorType, err, data.message, rawVersion)
			}
		}

		validTests := ValidVersionInputs{
			"1.2.3.4011": []string{
				"v1.2.3",
				"1.2.3",
			},
		}

		for encodedVer, tests := range validTests {
			for _, rawVersion := range tests {
				ver, err := ParseVersion(rawVersion)
				assert.Nil(t, err, rawVersion)
				assert.Equal(t, encodedVer, ver.String())
				assert.EqualValues(t, PRStageRelease, ver.PRStage())
				assert.EqualValues(t, 1, ver.PRVersion())
				assert.EqualValues(t, 1, ver.PRBuild())
			}
		}

		t.Run("majorV.minorV.patchV-prStage.prVersion[+prBuild]", func(t *testing.T) {
			preReleaseError := semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{})
			invalidTests := InvalidVersionInputs{
				{
					`divider of patchV and pre-release is "-"`,
					[]string{
						"1.2.3+rc.1",
						"1.2.3_beta.1",
						"1.2.3?beta.1",
					},
					semverErrors.NewInvalidPatchVersionError(semverErrors.InvalidPatchVersionErrorConfig{}),
				},
				{
					"prStage in:[release,rc,beta,alpha]",
					[]string{
						"1.2.3-rel.1",
						"1.2.3-rcrc.1",
						"1.2.3-bt.1",
						"1.2.3-al.1",
						"1.2.3-0.1",
						"1.2.3--1.1",
						"1.2.3-true.1",
						"1.2.3-nil.1",
					},
					preReleaseError,
				},
				{
					"prVersion shall be a positive integer",
					[]string{
						"1.2.3-release.-1",
						"1.2.3-release.0",
						"1.2.3-rc.",
						"1.2.3-beta.prVersion",
						"1.2.3-alpha.true",
					},
					preReleaseError,
				},
				{
					"prBuild fragment shall be nil or a natural number",
					[]string{
						"1.2.3-release.1+-1",
						"1.2.3-rc.2+",
						"1.2.3-beta.3+prBuild",
						"1.2.3-alpha.5+true",
					},
					preReleaseError,
				},
				{
					`divider of prStage and prVersion is "."`,
					[]string{
						"1.2.3-release+1",
						"1.2.3-rc-1",
						"1.2.3-beta_1",
						"1.2.3-beta?1",
					},
					preReleaseError,
				},
				{
					`divider of prStage and prBuild is "+"`,
					[]string{
						"1.2.3-release.1.1",
						"1.2.3-rc.2-1",
						"1.2.3-beta.3_1",
						"1.2.3-beta.4?1",
					},
					preReleaseError,
				},
			}

			for _, data := range invalidTests {
				for _, rawVersion := range data.input {
					_, err := ParseVersion(rawVersion)
					assert.IsType(t, data.errorType, err, data.message, rawVersion)
				}
			}

			validTests := ValidVersionInputs{
				"1.2.3.4011": []string{
					"v1.2.3-release.1",
					"1.2.3-release.1",
				},
				"1.2.3.3025": []string{
					"v1.2.3-rc.2+5",
					"1.2.3-rc.2+5",
				},
				"1.2.3.2031": []string{
					"v1.2.3-beta.3",
					"1.2.3-beta.3",
				},
				"1.2.3.1046": []string{
					"v1.2.3-alpha.4+6",
					"1.2.3-alpha.4+6",
				},
			}

			for encodedVer, tests := range validTests {
				for _, rawVersion := range tests {
					ver, err := ParseVersion(rawVersion)
					assert.Nil(t, err, rawVersion)
					assert.Equal(t, encodedVer, ver.String())
				}
			}
		})

		t.Run("majorV.minorV.patchV.numericPrVersion", func(t *testing.T) {
			preReleaseError := semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{})
			invalidTests := InvalidVersionInputs{
				{
					"numericPrVersion must be a four-digit integer",
					[]string{
						"1.2.3.release.1",
					},
					preReleaseError,
				},
				{
					"numericPrVersion must be a four-digit integer",
					[]string{
						"1.2.3.1",
						"1.2.3.12",
						"1.2.3.123",
						"1.2.3.12345",
						"1.2.3.123456",
					},
					preReleaseError,
				},
				{
					"first digit of numericPrVersion <= 4, >=1",
					[]string{
						"1.2.3.0234",
						"1.2.3.5234",
					},
					preReleaseError,
				},
				{
					"the second and third digit of numericPrVersion can not be zero both",
					[]string{
						"1.2.3.1004",
					},
					preReleaseError,
				},
				{
					"last digit of numericPrVersion > 0",
					[]string{
						"1.2.3.1230",
					},
					preReleaseError,
				},
			}

			for _, data := range invalidTests {
				for _, rawVersion := range data.input {
					_, err := ParseVersion(rawVersion)
					assert.IsType(t, data.errorType, err, data.message, rawVersion)
				}
			}

			validTests := ValidVersionInputs{
				"1.2.3.4014": []string{
					"v1.2.3.4014",
					"1.2.3.4014",
				},
				"1.2.3.3023": []string{
					"v1.2.3.3023",
					"1.2.3.3023",
				},
				"1.2.3.2032": []string{
					"v1.2.3.2032",
					"1.2.3.2032",
				},
				"1.2.3.1041": []string{
					"v1.2.3.1041",
					"1.2.3.1041",
				},
			}

			for encodedVer, tests := range validTests {
				for _, rawVersion := range tests {
					ver, err := ParseVersion(rawVersion)
					assert.Nil(t, err, rawVersion)
					assert.Equal(t, encodedVer, ver.String())
				}
			}
		})
	})
}

func TestNew(t *testing.T) {
	newVersionConfig := func() CtorConfig {
		return CtorConfig{
			MajorV:    1,
			MinorV:    2,
			PatchV:    3,
			PRStage:   4,
			PRVersion: 5,
			PRBuild:   6,
		}
	}

	t.Run("0.0.0 is not meaningless", func(t *testing.T) {
		tErr := semverErrors.NewInvalidRawVersionFormatError(semverErrors.InvalidRawVersionFormatErrorConfig{})

		verConfig := newVersionConfig()

		verConfig.MajorV = 0
		verConfig.MinorV = 0
		verConfig.PatchV = 0

		_, err := New(verConfig)
		assert.IsType(t, tErr, err, verConfig.ToRawVersion())
	})

	t.Run("prStage in:1,2,3,4", func(t *testing.T) {
		tErr := semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{})

		verConfig := newVersionConfig()

		for _, v := range []uint8{0, 5, 6, 7} {
			verConfig.PRStage = v

			_, err := New(verConfig)
			assert.IsType(t, tErr, err, verConfig.ToRawVersion())
		}
	})

	t.Run("prVersion in:[1,99]", func(t *testing.T) {
		tErr := semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{})

		verConfig := newVersionConfig()

		for _, v := range []uint8{0, 100, 101} {
			verConfig.PRStage = v

			_, err := New(verConfig)
			assert.IsType(t, tErr, err, verConfig.ToRawVersion())
		}
	})

	t.Run("prBuild in:[1,9]", func(t *testing.T) {
		tErr := semverErrors.NewInvalidPreReleaseVersionError(semverErrors.InvalidPreReleaseVersionErrorConfig{})

		verConfig := newVersionConfig()

		for _, v := range []uint8{0, 10, 11} {
			verConfig.PRStage = v

			_, err := New(verConfig)
			assert.IsType(t, tErr, err, verConfig.ToRawVersion())
		}
	})
}

func TestCompare(t *testing.T) {
	rawSemVer := "4.3.0-rc.1+2"
	semVer, _ := ParseVersion(rawSemVer)

	tests := []struct {
		in                       interface{}
		stageSafetyCompareExpect int
		compareExpect            int
	}{
		{
			"invalid_version",
			CompareLargeFlag,
			CompareLargeFlag,
		},
		{
			"4.2.0-beta.10",
			CompareLargeFlag,
			CompareLargeFlag,
		},
		{
			"4.2.0-rc.3",
			CompareLargeFlag,
			CompareLargeFlag,
		},
		{
			"4.2.0-release.1",
			CompareLessFlag,
			CompareLargeFlag,
		},
		{
			"4.3.0-rc.1+1",
			CompareLargeFlag,
			CompareEqualFlag,
		},
		{
			"4.3.0-rc.1+2",
			CompareEqualFlag,
			CompareEqualFlag,
		},
		{
			"4.3.0-rc.1+3",
			CompareLessFlag,
			CompareEqualFlag,
		},
		{
			"4.3.0-rc.2",
			CompareLessFlag,
			CompareEqualFlag,
		},
		{
			"4.3.1-beta.10",
			CompareLargeFlag,
			CompareLessFlag,
		},
		{
			"4.3.1-rc.1",
			CompareLessFlag,
			CompareLessFlag,
		},
		{
			"100.1000.10000-alpha.1",
			CompareLargeFlag,
			CompareLessFlag,
		},
		{
			"100.1000.10000-rc.1",
			CompareLessFlag,
			CompareLessFlag,
		},
	}

	for _, entry := range tests {
		assert.Equal(t, entry.stageSafetyCompareExpect, semVer.StageSafetyCompare(entry.in), "stage safety compare, target: %s, input: %s", rawSemVer, entry.in)
		assert.Equal(t, entry.compareExpect, semVer.Compare(entry.in), "compare, target: %s, input: %s", rawSemVer, entry.in)
	}
}
