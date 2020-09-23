package usecase_test

import (
	"github.com/funnyecho/code-push/daemon"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	t.Run("invalid params", func(t *testing.T) {
		{
			err := useCase.ReleaseVersion(nil)
			assert.True(t, errors.Is(err, daemon.ErrParamsInvalid))
		}
		{
			ver, err := useCase.GetVersion(nil, nil)
			assert.True(t, errors.Is(err, daemon.ErrParamsInvalid))
			assert.Nil(t, ver)

			_, env := createRandomEnv()

			ver, err = useCase.GetVersion([]byte(env.ID), nil)
			assert.True(t, errors.Is(err, daemon.ErrParamsInvalid))
			assert.Nil(t, ver)

			ver, err = useCase.GetVersion([]byte(env.ID), []byte("not.a.version"))
			assert.True(t, errors.Is(err, daemon.ErrParamsInvalid))
			assert.Nil(t, ver)
		}
		{
			list, err := useCase.ListVersions(nil)
			assert.True(t, errors.Is(err, daemon.ErrParamsInvalid))
			assert.Nil(t, list)
		}
		{
			result, err := useCase.VersionStrictCompatQuery(nil, nil)
			assert.True(t, errors.Is(err, daemon.ErrParamsInvalid))
			assert.Nil(t, result)

			result, err = useCase.VersionStrictCompatQuery([]byte("foo"), nil)
			assert.True(t, errors.Is(err, daemon.ErrParamsInvalid))
			assert.Nil(t, result)
		}
	})

	t.Run("release version on invalid env", func(t *testing.T) {
		t.Run("failed with nil envId", func(t *testing.T) {
			err := useCase.ReleaseVersion(&versionReleaseParams{
				envId:            nil,
				appVersion:       []byte("1.1.0"),
				compatAppVersion: nil,
				changelog:        nil,
				packageFileKey:   []byte("dfjkd"),
				mustUpdate:       false,
			})
			assert.True(t, errors.Is(err, daemon.ErrParamsInvalid))
		})

		t.Run("on non existed env", func(t *testing.T) {
			err := useCase.ReleaseVersion(&versionReleaseParams{
				envId:            []byte("non existed env"),
				appVersion:       []byte("1.1.0"),
				compatAppVersion: nil,
				changelog:        nil,
				packageFileKey:   []byte("dfjkd"),
				mustUpdate:       false,
			})
			assert.True(t, errors.Is(err, daemon.ErrEnvNotFound))
		})
	})

	t.Run("release version on valid env", func(t *testing.T) {
		branch, env := createRandomEnv()
		envIdSlice := []byte(env.ID)

		t.Run("failed with invalid params", func(t *testing.T) {
			err := useCase.ReleaseVersion(&versionReleaseParams{
				envId:            envIdSlice,
				appVersion:       nil,
				compatAppVersion: nil,
				changelog:        nil,
				packageFileKey:   []byte("dfjkd"),
				mustUpdate:       false,
			})
			assert.True(t, errors.Is(err, daemon.ErrParamsInvalid), "appVersion required")

			err = useCase.ReleaseVersion(&versionReleaseParams{
				envId:            envIdSlice,
				appVersion:       []byte("1.1.0"),
				compatAppVersion: nil,
				changelog:        nil,
				packageFileKey:   nil,
				mustUpdate:       false,
			})
			assert.True(t, errors.Is(err, daemon.ErrParamsInvalid), "packageFileKey required")
		})

		t.Run("release without compatAppVersion specified", func(t *testing.T) {
			err := useCase.ReleaseVersion(&versionReleaseParams{
				envId:            envIdSlice,
				appVersion:       []byte("1.2.0"),
				compatAppVersion: nil,
				changelog:        nil,
				packageFileKey:   []byte("dfjkd"),
				mustUpdate:       false,
			})
			assert.NoError(t, err)

			fetchVer, fetchErr := useCase.GetVersion(envIdSlice, []byte("1.2.0"))
			assert.NoError(t, fetchErr)
			assert.Equal(t, "1.2.0.4011", fetchVer.AppVersion)
			assert.Equal(t, fetchVer.AppVersion, fetchVer.CompatAppVersion, "compatAppVersion equal to appVersion when compatAppVersion is missing")
		})

		t.Run("release with supported format of version", func(t *testing.T) {
			env, _ := useCase.CreateEnv([]byte(branch.ID), []byte(t.Name()))
			envIdSlice := []byte(env.ID)

			data := [][2]string{
				{"v1.1.0-alpha.1", "1.1.0.1011"},
				{"1.1.0-beta.1", "1.1.0.2011"},
				{"v1.1.0-rc.1", "1.1.0.3011"},
				{"v1.1.0-release.1", "1.1.0.4011"},
				{"1.1.1", "1.1.1.4011"},

				{"1.2.0.1011", "1.2.0.1011"},
				{"v1.2.0.2011", "1.2.0.2011"},
				{"1.2.0.3011", "1.2.0.3011"},
				{"v1.2.0.4011", "1.2.0.4011"},
				{"1.2.1", "1.2.1.4011"},
			}

			for _, item := range data {
				input := item[0]

				err := useCase.ReleaseVersion(&versionReleaseParams{
					envId:            envIdSlice,
					appVersion:       []byte(input),
					compatAppVersion: nil,
					changelog:        nil,
					packageFileKey:   []byte("dfjkd"),
					mustUpdate:       false,
				})
				assert.NoError(t, err)
			}

			t.Run("fetch version", func(t *testing.T) {
				for _, item := range data {
					input := item[0]
					output := item[1]

					fetchVer, fetchErr := useCase.GetVersion(envIdSlice, []byte(input))
					assert.NoError(t, fetchErr)
					assert.Equal(t, output, fetchVer.AppVersion)
				}
			})

			t.Run("fetch version list", func(t *testing.T) {
				list, listErr := useCase.ListVersions(envIdSlice)
				assert.NoError(t, listErr)
				assert.Equal(t, len(data), len(list))

				for _, item := range data {
					output := item[1]

					var match *daemon.Version

					for _, i := range list {
						if i.AppVersion == output {
							match = i
							break
						}
					}

					assert.NotNil(t, match)
				}
			})
		})

		t.Run("release with existed appVersion", func(t *testing.T) {
			env, _ := useCase.CreateEnv([]byte(branch.ID), []byte(t.Name()))
			envIdSlice := []byte(env.ID)

			err := useCase.ReleaseVersion(&versionReleaseParams{
				envId:            envIdSlice,
				appVersion:       []byte("1.1.0"),
				compatAppVersion: nil,
				changelog:        nil,
				packageFileKey:   []byte("dfjkd"),
				mustUpdate:       false,
			})
			assert.NoError(t, err)

			err = useCase.ReleaseVersion(&versionReleaseParams{
				envId:            envIdSlice,
				appVersion:       []byte("1.1.0"),
				compatAppVersion: nil,
				changelog:        nil,
				packageFileKey:   []byte("dfjkd"),
				mustUpdate:       false,
			})
			assert.True(t, errors.Is(err, daemon.ErrVersionExisted), "failed to release existed app version")
		})

		t.Run("release when compatAppVersion is large than appVersion", func(t *testing.T) {
			data := [][2]string{
				{"1.1.0-alpha.1", "1.1.0-alpha.1+2"},
				{"1.1.0-alpha.1", "1.1.0-alpha.2"},
				{"1.1.0-alpha.1", "1.1.0-beta.2"},
				{"1.1.0", "1.2.0-alpha.1"},
			}

			for _, item := range data {
				err := useCase.ReleaseVersion(&versionReleaseParams{
					envId:            envIdSlice,
					appVersion:       []byte(item[0]),
					compatAppVersion: []byte(item[1]),
					changelog:        nil,
					packageFileKey:   []byte("dfjkd"),
					mustUpdate:       false,
				})
				assert.True(t, errors.Is(err, daemon.ErrParamsInvalid), "compatAppVersion <= appVersion")
			}
		})
	})

	t.Run("fetch non existed version", func(t *testing.T) {
		t.Log("if version not existed, no error occur, and return empty data")

		_, env := createRandomEnv()

		ver, err := useCase.GetVersion([]byte(env.ID), []byte("1.0.0"))
		assert.NoError(t, err)
		assert.Nil(t, ver)
	})

	t.Run("fetch empty version list", func(t *testing.T) {
		t.Log("if env has no version, no error occur, and return empty list")

		_, env := createRandomEnv()
		list, listErr := useCase.ListVersions([]byte(env.ID))
		assert.NoError(t, listErr)
		assert.NotNil(t, list)
		assert.Equal(t, 0, len(list))
	})

	t.Run("fetch on non existed env", func(t *testing.T) {
		ver, getErr := useCase.GetVersion([]byte("non existed env"), []byte("1.1.0"))
		assert.True(t, errors.Is(getErr, daemon.ErrEnvNotFound))
		assert.Nil(t, ver)

		list, listErr := useCase.ListVersions([]byte("not existed env"))
		assert.True(t, errors.Is(listErr, daemon.ErrEnvNotFound))
		assert.Nil(t, list)
	})
}

func TestVersionCompatQuery(t *testing.T) {
	t.Run("SimpleCompatRange", func(t *testing.T) {
		_, env := createRandomEnv()

		envIdSlice := []byte(env.ID)

		data := [][2]string{
			{"v1.0.0", "1.0.0"},
			{"v1.0.2-release.1", "1.0.0"},
			{"v1.0.4-rc.1", "1.0.0"},
			{"v1.0.6-beta.1", "1.0.0"},
			{"v1.0.8-alpha.1", "1.0.0"},
		}

		for _, item := range data {
			appVer := item[0]
			compatAppVer := item[1]

			err := useCase.ReleaseVersion(&versionReleaseParams{
				envId:            envIdSlice,
				appVersion:       []byte(appVer),
				compatAppVersion: []byte(compatAppVer),
				changelog:        nil,
				packageFileKey:   []byte("file-key"),
				mustUpdate:       false,
			})
			assert.NoError(t, err)
		}

		result, err := useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.0.0-alpha.3"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Nil(t, result.CanUpdateAppVersion())
		assert.Equal(t, "1.0.8.1011", string(result.LatestAppVersion()))

		result, err = useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.0.0"))
		assert.NoError(t, err)
		assert.Equal(t, "1.0.0.4011", string(result.AppVersion()))
		assert.Equal(t, "1.0.2.4011", string(result.CanUpdateAppVersion()))
		assert.Equal(t, "1.0.2.4011", string(result.LatestAppVersion()))

		result, err = useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.0.1-alpha.3"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Equal(t, "1.0.8.1011", string(result.CanUpdateAppVersion()))
		assert.Equal(t, "1.0.8.1011", string(result.LatestAppVersion()))

		result, err = useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.0.1-beta.3"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Equal(t, "1.0.6.2011", string(result.CanUpdateAppVersion()))
		assert.Equal(t, "1.0.6.2011", string(result.LatestAppVersion()))

		result, err = useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.0.1-rc.3"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Equal(t, "1.0.4.3011", string(result.CanUpdateAppVersion()))
		assert.Equal(t, "1.0.4.3011", string(result.LatestAppVersion()))

		result, err = useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.0.1-release.3"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Equal(t, "1.0.2.4011", string(result.CanUpdateAppVersion()))
		assert.Equal(t, "1.0.2.4011", string(result.LatestAppVersion()))

		result, err = useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.0.6-release.3"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Nil(t, result.CanUpdateAppVersion())
		assert.Nil(t, result.LatestAppVersion())

		result, err = useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.0.10"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Nil(t, result.CanUpdateAppVersion())
		assert.Nil(t, result.LatestAppVersion())
	})

	t.Run("MultiCompatRange", func(t *testing.T) {
		_, env := createRandomEnv()

		envIdSlice := []byte(env.ID)

		data := [][2]string{
			{"v1.0.0", "1.0.0"},
			{"v1.0.2", "1.0.0"},
			{"v1.0.4", "1.0.0"},
			{"v1.0.6", "1.0.0"},

			{"1.1.0", "1.0.4"},
			{"1.1.2", "1.0.4"},
			{"1.1.4", "1.0.4"},

			{"1.2.0", "1.1.4"},
			{"1.2.2", "1.1.4"},
			{"1.2.4", "1.1.4"},
		}

		for _, item := range data {
			appVer := item[0]
			compatAppVer := item[1]

			err := useCase.ReleaseVersion(&versionReleaseParams{
				envId:            envIdSlice,
				appVersion:       []byte(appVer),
				compatAppVersion: []byte(compatAppVer),
				changelog:        nil,
				packageFileKey:   []byte("file-key"),
				mustUpdate:       false,
			})
			assert.NoError(t, err)
		}

		result, err := useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.0.3"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Equal(t, "1.0.6.4011", string(result.CanUpdateAppVersion()))
		assert.Equal(t, "1.2.4.4011", string(result.LatestAppVersion()))

		result, err = useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.0.5"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Equal(t, "1.1.4.4011", string(result.CanUpdateAppVersion()))
		assert.Equal(t, "1.2.4.4011", string(result.LatestAppVersion()))

		result, err = useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.1.0"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Equal(t, "1.1.4.4011", string(result.CanUpdateAppVersion()))
		assert.Equal(t, "1.2.4.4011", string(result.LatestAppVersion()))

		result, err = useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.1.4"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Equal(t, "1.2.4.4011", string(result.CanUpdateAppVersion()))
		assert.Equal(t, "1.2.4.4011", string(result.LatestAppVersion()))

		result, err = useCase.VersionStrictCompatQuery(envIdSlice, []byte("1.2.5"))
		assert.NoError(t, err)
		assert.Nil(t, result.AppVersion(), "nil when query version not existed")
		assert.Nil(t, result.CanUpdateAppVersion())
		assert.Nil(t, result.LatestAppVersion())
	})
}

type versionReleaseParams struct {
	envId            []byte
	appVersion       []byte
	compatAppVersion []byte
	changelog        []byte
	packageFileKey   []byte
	mustUpdate       bool
}

func (v *versionReleaseParams) EnvId() []byte {
	return v.envId
}

func (v *versionReleaseParams) AppVersion() []byte {
	return v.appVersion
}

func (v *versionReleaseParams) CompatAppVersion() []byte {
	return v.compatAppVersion
}

func (v *versionReleaseParams) Changelog() []byte {
	return v.changelog
}

func (v *versionReleaseParams) PackageFileKey() []byte {
	return v.packageFileKey
}

func (v *versionReleaseParams) MustUpdate() bool {
	return v.mustUpdate
}
