package semver

import (
	"github.com/funnyecho/code-push/pkg/errors"
	semverErrors "github.com/funnyecho/code-push/pkg/semver/errors"
	"testing"
)

func TestParseVersion_InvalidInput(t *testing.T) {
	failedDataTable := []struct {
		in             []string
		errorAssertion func(err errors.IError) (hasError bool, expectType string)
	}{
		{
			[]string{
				"v2",
				"v2.ex",
				"v2.3",
				"2.3",
				"abcdefg",
				"abc.efg",
			},
			func(err errors.IError) (hasError bool, expectType string) {
				if _, isOk := err.(*semverErrors.InvalidRawVersionFormatError); !isOk {
					return true, "semverErrors.InvalidRawVersionFormatError"
				} else {
					return false, ""
				}
			},
		},
		{
			[]string{
				"-1.2.3",
				"foo.2.3",
			},
			func(err errors.IError) (hasError bool, expectType string) {
				if _, isOk := err.(*semverErrors.InvalidMajorVersionError); !isOk {
					return true, "semverErrors.InvalidMajorVersionError"
				} else {
					return false, ""
				}
			},
		},
		{
			[]string{
				"1.-2.3",
				"1.foo.3",
			},
			func(err errors.IError) (hasError bool, expectType string) {
				if _, isOk := err.(*semverErrors.InvalidMinorVersionError); !isOk {
					return true, "semverErrors.InvalidMinorVersionError"
				} else {
					return false, ""
				}
			},
		},
		{
			[]string{
				"1.2.-3",
				"1.2.foo",
			},
			func(err errors.IError) (hasError bool, expectType string) {
				if _, isOk := err.(*semverErrors.InvalidPatchVersionError); !isOk {
					return true, "semverErrors.InvalidPatchVersionError"
				} else {
					return false, ""
				}
			},
		},
		{
			[]string{
				"1.2.3.alpha1",
				"1.2.3-alpha1",
				"1.2.3-other_stage.1",
				"1.2.3-alpha.-1",
				"1.2.3.2.4",
				"1.2.3.210",
				"1.2.3.-111",
			},
			func(err errors.IError) (hasError bool, expectType string) {
				if _, isOk := err.(*semverErrors.InvalidPreReleaseVersionError); !isOk {
					return true, "semverErrors.InvalidPreReleaseVersionError"
				} else {
					return false, ""
				}
			},
		},
	}

	for _, entry := range failedDataTable {
		in := entry.in
		errorAssertion := entry.errorAssertion

		for _, v := range in {
			_, parseErr := ParseVersion(v)
			assertErr, expectType := errorAssertion(parseErr)
			if assertErr {
				t.Fatalf(`wrong error type of version#%v; expected %s; got: %T`, v, expectType, parseErr)
			}
		}
	}
}

func TestParseVersion_ValidInput(t *testing.T) {
	validDataInput := []struct {
		in         []string
		assertFunc func(semVer *SemVer, err errors.IError)
	}{
		{
			[]string{
				"v1.2.3",
				"1.2.3",
				"1.2.3.4011",
				"1.2.3-release.1",
			},
			func(semVer *SemVer, err errors.IError) {
				if err != nil {
					t.Fatalf("expect no error; got %v", err.Error())
				}

				if semVer.MajorVersion() != 1 {
					t.Fatalf("expect %d; got %v", 1, semVer.MajorVersion())
				}

				if semVer.MinorVersion() != 2 {
					t.Fatalf("expect %d; got %v", 2, semVer.MinorVersion())
				}

				if semVer.PatchVersion() != 3 {
					t.Fatalf("expect %d; got %v", 3, semVer.PatchVersion())
				}

				if semVer.PRStage() != PRStageRelease {
					t.Fatalf("expect %d; got %v", PRStageRelease, semVer.PRStage())
				}

				if semVer.PRVersion() != 1 {
					t.Fatalf("expect %d; got %v", 1, semVer.PRVersion())
				}

				if semVer.PRBuild() != 1 {
					t.Fatalf("expect %d; got %v", 1, semVer.PRBuild())
				}

				if semVer.String() != "1.2.3.4011" {
					t.Fatalf("expect %s; got %v", "1.2.3.4011", semVer.String())
				}
			},
		},
		{
			[]string{
				"v1.2.3-alpha.2",
				"1.2.3.1021",
			},
			func(semVer *SemVer, err errors.IError) {
				if err != nil {
					t.Fatalf("expect no error; got %v", err.Error())
				}

				if semVer.MajorVersion() != 1 {
					t.Fatalf("expect %d; got %v", 1, semVer.MajorVersion())
				}

				if semVer.MinorVersion() != 2 {
					t.Fatalf("expect %d; got %v", 2, semVer.MinorVersion())
				}

				if semVer.PatchVersion() != 3 {
					t.Fatalf("expect %d; got %v", 3, semVer.PatchVersion())
				}

				if semVer.PRStage() != PRStageAlpha {
					t.Fatalf("expect %d; got %v", PRStageAlpha, semVer.PRStage())
				}

				if semVer.PRVersion() != 2 {
					t.Fatalf("expect %d; got %v", 2, semVer.PRVersion())
				}

				if semVer.PRBuild() != 1 {
					t.Fatalf("expect %d; got %v", 1, semVer.PRBuild())
				}

				if semVer.String() != "1.2.3.1021" {
					t.Fatalf("expect %s; got %v", "1.2.3.1021", semVer.String())
				}
			},
		},
		{
			[]string{
				"v1.2.3-beta.3",
				"1.2.3.2031",
			},
			func(semVer *SemVer, err errors.IError) {
				if err != nil {
					t.Fatalf("expect no error; got %v", err.Error())
				}

				if semVer.MajorVersion() != 1 {
					t.Fatalf("expect %d; got %v", 1, semVer.MajorVersion())
				}

				if semVer.MinorVersion() != 2 {
					t.Fatalf("expect %d; got %v", 2, semVer.MinorVersion())
				}

				if semVer.PatchVersion() != 3 {
					t.Fatalf("expect %d; got %v", 3, semVer.PatchVersion())
				}

				if semVer.PRStage() != PRStageBeta {
					t.Fatalf("expect %d; got %v", PRStageBeta, semVer.PRStage())
				}

				if semVer.PRVersion() != 3 {
					t.Fatalf("expect %d; got %v", 3, semVer.PRVersion())
				}

				if semVer.PRBuild() != 1 {
					t.Fatalf("expect %d; got %v", 1, semVer.PRBuild())
				}

				if semVer.String() != "1.2.3.2031" {
					t.Fatalf("expect %s; got %v", "1.2.3.2031", semVer.String())
				}
			},
		},
		{
			[]string{
				"v1.2.3-rc.4",
				"1.2.3.3041",
			},
			func(semVer *SemVer, err errors.IError) {
				if err != nil {
					t.Fatalf("expect no error; got %v", err.Error())
				}

				if semVer.MajorVersion() != 1 {
					t.Fatalf("expect %d; got %v", 1, semVer.MajorVersion())
				}

				if semVer.MinorVersion() != 2 {
					t.Fatalf("expect %d; got %v", 2, semVer.MinorVersion())
				}

				if semVer.PatchVersion() != 3 {
					t.Fatalf("expect %d; got %v", 3, semVer.PatchVersion())
				}

				if semVer.PRStage() != PRStageRC {
					t.Fatalf("expect %d; got %v", PRStageRC, semVer.PRStage())
				}

				if semVer.PRVersion() != 4 {
					t.Fatalf("expect %d; got %v", 4, semVer.PRVersion())
				}

				if semVer.PRBuild() != 1 {
					t.Fatalf("expect %d; got %v", 1, semVer.PRBuild())
				}

				if semVer.String() != "1.2.3.3041" {
					t.Fatalf("expect %s; got %v", "1.2.3.3041", semVer.String())
				}
			},
		},
	}

	for _, entry := range validDataInput {
		in := entry.in
		assertFunc := entry.assertFunc

		for _, v := range in {
			semVer, parseErr := ParseVersion(v)
			assertFunc(semVer, parseErr)
		}
	}
}

func TestSemVer_Compare(t *testing.T) {
	semVer, _ := ParseVersion("4.3.0-rc.1")

	dataTable := []struct {
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

	for _, entry := range dataTable {
		if compareResult := semVer.StageSafetyCompare(entry.in); compareResult != entry.stageSafetyCompareExpect {
			t.Fatalf("StageSafetyCompare input: %v; expecte: %d; got: %d", entry.in, entry.stageSafetyCompareExpect, compareResult)
		}

		if compareResult := semVer.Compare(entry.in); compareResult != entry.compareExpect {
			t.Fatalf("Compare input: %v; expecte: %d; got: %d", entry.in, entry.compareExpect, compareResult)
		}
	}
}
