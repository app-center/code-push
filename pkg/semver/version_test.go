package semver

import (
	semverErrors "github.com/funnyecho/code-push/pkg/semver/errors"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestParseVersion(t *testing.T) {
	type InvalidVersionInputs = []struct {
		name      string
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
				new(semverErrors.InvalidMajorVersionError),
			},
			{
				"minorV shall be a natural number",
				[]string{
					"1.-2.3",
					"1.minor.3",
					"1..3",
				},
				new(semverErrors.InvalidMinorVersionError),
			},
			{
				"patchV shall be a natural number",
				[]string{
					"1.2.-3",
					"1.2.patch",
					"1.2.",
				},
				new(semverErrors.InvalidPatchVersionError),
			},
			{
				"0.0.0 is not meaningless",
				[]string{
					"0.0.0",
				},
				new(semverErrors.InvalidRawVersionFormatError),
			},
			{
				"divider must be a dot",
				[]string{
					"1,2,3",
					"1-2-3",
					"1_2_3",
				},
				new(semverErrors.InvalidRawVersionFormatError),
			},
		}

		for _, data := range invalidTests {
			t.Run(data.name, func(t *testing.T) {
				for _, rawVersion := range data.input {
					_, err := ParseVersion(rawVersion)
					assert.IsType(t, data.errorType, err)
				}
			})
		}

		t.Run("majorV.minorV.patchV-prStage.prVersion", func(t *testing.T) {
			preReleaseError := new(semverErrors.InvalidPreReleaseVersionError)
			invalidTests := InvalidVersionInputs{
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
					"prVersion shall be a natural number",
					[]string{
						"1.2.3-release.-1",
						"1.2.3-rc.",
						"1.2.3-beta.prVersion",
						"1.2.3-alpha.true",
					},
					preReleaseError,
				},
				{
					`divider of pre-release and patchV is "-"`,
					[]string{
						"1.2.3.release.1",
						"1.2.3+rc.1",
						"1.2.3_beta.1",
						"1.2.3?beta.1",
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
			}

			for _, data := range invalidTests {
				t.Run(data.name, func(t *testing.T) {
					for _, rawVersion := range data.input {
						_, err := ParseVersion(rawVersion)
						assert.IsType(t, data.errorType, err, rawVersion)
					}
				})
			}
		})
		t.Run("majorV.minorV.patchV.numericPrVersion", func(t *testing.T) {
			preReleaseError := new(semverErrors.InvalidPreReleaseVersionError)
			invalidTests := InvalidVersionInputs{
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
					"last digit of numericPrVersion > 0",
					[]string{
						"1.2.3.1230",
					},
					preReleaseError,
				},
			}

			for _, data := range invalidTests {
				t.Run(data.name, func(t *testing.T) {
					for _, rawVersion := range data.input {
						_, err := ParseVersion(rawVersion)
						assert.IsType(t, data.errorType, err, rawVersion)
					}
				})
			}
		})
	})
}
