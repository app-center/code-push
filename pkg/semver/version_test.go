package semver_test

import (
	"github.com/funnyecho/code-push/pkg/semver"
	"github.com/pkg/errors"
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
				semver.ErrInvalidMajorVersion,
			},
			{
				"minorV shall be a natural number",
				[]string{
					"1.-2.3",
					"1.minor.3",
					"1..3",
				},
				semver.ErrInvalidMinorVersion,
			},
			{
				"patchV shall be a natural number",
				[]string{
					"1.2.-3",
					"1.2.patch",
					"1.2.",
				},
				semver.ErrInvalidPatchVersion,
			},
			{
				"0.0.0 is not meaningless",
				[]string{
					"0.0.0",
				},
				semver.ErrInvalidVersionFormat,
			},
			{
				"divider must be a dot",
				[]string{
					"1,2,3",
					"1-2-3",
					"1_2_3",
				},
				semver.ErrInvalidVersionFormat,
			},
		}

		for _, data := range invalidTests {
			for _, rawVersion := range data.input {
				_, err := semver.ParseVersion(rawVersion)
				assert.True(t, errors.Is(err, data.errorType), data.message, rawVersion)
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
				ver, err := semver.ParseVersion(rawVersion)
				assert.Nil(t, err, rawVersion)
				assert.Equal(t, encodedVer, ver.String())
				assert.EqualValues(t, semver.PRStageRelease, ver.PRStage())
				assert.EqualValues(t, 1, ver.PRVersion())
				assert.EqualValues(t, 1, ver.PRBuild())
			}
		}

		t.Run("majorV.minorV.patchV-prStage.prVersion[+prBuild]", func(t *testing.T) {
			preReleaseError := semver.ErrInvalidPreReleaseVersion
			invalidTests := InvalidVersionInputs{
				{
					`divider of patchV and pre-release is "-"`,
					[]string{
						"1.2.3+rc.1",
						"1.2.3_beta.1",
						"1.2.3?beta.1",
					},
					semver.ErrInvalidPatchVersion,
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
					_, err := semver.ParseVersion(rawVersion)
					assert.True(t, errors.Is(err, data.errorType), data.message, rawVersion)
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
					ver, err := semver.ParseVersion(rawVersion)
					assert.Nil(t, err, rawVersion)
					assert.Equal(t, encodedVer, ver.String())
				}
			}
		})

		t.Run("majorV.minorV.patchV.numericPrVersion", func(t *testing.T) {
			preReleaseError := semver.ErrInvalidPreReleaseVersion
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
					_, err := semver.ParseVersion(rawVersion)
					assert.True(t, errors.Is(err, data.errorType), data.message, rawVersion)
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
					ver, err := semver.ParseVersion(rawVersion)
					assert.Nil(t, err, rawVersion)
					assert.Equal(t, encodedVer, ver.String())
				}
			}
		})
	})
}

func TestNew(t *testing.T) {
	newVersionConfig := func() semver.CtorConfig {
		return semver.CtorConfig{
			MajorV:    1,
			MinorV:    2,
			PatchV:    3,
			PRStage:   4,
			PRVersion: 5,
			PRBuild:   6,
		}
	}

	t.Run("0.0.0 is not meaningless", func(t *testing.T) {
		tErr := semver.ErrInvalidVersionFormat

		verConfig := newVersionConfig()

		verConfig.MajorV = 0
		verConfig.MinorV = 0
		verConfig.PatchV = 0

		_, err := semver.New(verConfig)
		assert.True(t, errors.Is(err, tErr), verConfig.ToRawVersion())
	})

	t.Run("prStage in:1,2,3,4", func(t *testing.T) {
		tErr := semver.ErrInvalidPreReleaseVersion

		verConfig := newVersionConfig()

		for _, v := range []uint8{0, 5, 6, 7} {
			verConfig.PRStage = v

			_, err := semver.New(verConfig)
			assert.True(t, errors.Is(err, tErr), verConfig.ToRawVersion())
		}
	})

	t.Run("prVersion in:[1,99]", func(t *testing.T) {
		tErr := semver.ErrInvalidPreReleaseVersion

		verConfig := newVersionConfig()

		for _, v := range []uint8{0, 100, 101} {
			verConfig.PRStage = v

			_, err := semver.New(verConfig)
			assert.True(t, errors.Is(err, tErr), verConfig.ToRawVersion())
		}
	})

	t.Run("prBuild in:[1,9]", func(t *testing.T) {
		tErr := semver.ErrInvalidPreReleaseVersion

		verConfig := newVersionConfig()

		for _, v := range []uint8{0, 10, 11} {
			verConfig.PRStage = v

			_, err := semver.New(verConfig)
			assert.True(t, errors.Is(err, tErr), verConfig.ToRawVersion())
		}
	})
}

func TestCompare(t *testing.T) {
	rawSemVer := "4.3.0-rc.1+2"
	semVer, _ := semver.ParseVersion(rawSemVer)

	tests := []struct {
		in                       interface{}
		stageSafetyCompareExpect int
		compareExpect            int
	}{
		{
			"invalid_version",
			semver.CompareLargeFlag,
			semver.CompareLargeFlag,
		},
		{
			"4.2.0-beta.10",
			semver.CompareLargeFlag,
			semver.CompareLargeFlag,
		},
		{
			"4.2.0-rc.3",
			semver.CompareLargeFlag,
			semver.CompareLargeFlag,
		},
		{
			"4.2.0-release.1",
			semver.CompareLessFlag,
			semver.CompareLargeFlag,
		},
		{
			"4.3.0-rc.1+1",
			semver.CompareLargeFlag,
			semver.CompareEqualFlag,
		},
		{
			"4.3.0-rc.1+2",
			semver.CompareEqualFlag,
			semver.CompareEqualFlag,
		},
		{
			"4.3.0-rc.1+3",
			semver.CompareLessFlag,
			semver.CompareEqualFlag,
		},
		{
			"4.3.0-rc.2",
			semver.CompareLessFlag,
			semver.CompareEqualFlag,
		},
		{
			"4.3.1-beta.10",
			semver.CompareLargeFlag,
			semver.CompareLessFlag,
		},
		{
			"4.3.1-rc.1",
			semver.CompareLessFlag,
			semver.CompareLessFlag,
		},
		{
			"100.1000.10000-alpha.1",
			semver.CompareLargeFlag,
			semver.CompareLessFlag,
		},
		{
			"100.1000.10000-rc.1",
			semver.CompareLessFlag,
			semver.CompareLessFlag,
		},
	}

	for _, entry := range tests {
		assert.Equal(t, entry.stageSafetyCompareExpect, semVer.StageSafetyStrictCompare(entry.in), "stage safety compare, target: %s, input: %s", rawSemVer, entry.in)
		assert.Equal(t, entry.compareExpect, semVer.Compare(entry.in), "compare, target: %s, input: %s", rawSemVer, entry.in)
	}
}
