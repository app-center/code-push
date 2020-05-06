package version_compat_tree

import (
	"github.com/funnyecho/code-push/pkg/semver"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testEntry struct {
	compatVer *semver.SemVer
	appVer    *semver.SemVer
}

func (t *testEntry) CompatVersion() *semver.SemVer {
	return t.compatVer
}

func (t *testEntry) Version() *semver.SemVer {
	return t.appVer
}

func (t *testEntry) RawVersion() string {
	return t.appVer.String()
}

type testCompatQueryAnchor struct {
	ver *semver.SemVer
}

func (t *testCompatQueryAnchor) Version() *semver.SemVer {
	return t.ver
}

func parseVersion(rawVer string) *semver.SemVer {
	ver, _ := semver.ParseVersion(rawVer)
	return ver
}

func rawVersionOfEntry(entry IEntry) (raw string) {
	if entry.Version() != nil {
		raw = entry.Version().String()
	}

	return
}

func TestSingleCompatRange(t *testing.T) {
	tree := NewVersionCompatTree(Entries{
		&testEntry{
			compatVer: parseVersion("v1.0.0"),
			appVer:    parseVersion("v1.0.0"),
		},
		&testEntry{
			compatVer: parseVersion("v1.0.0"),
			appVer:    parseVersion("v1.0.2-release.1"),
		},
		&testEntry{
			compatVer: parseVersion("v1.0.0"),
			appVer:    parseVersion("v1.0.4-rc.1"),
		},
		&testEntry{
			compatVer: parseVersion("v1.0.0"),
			appVer:    parseVersion("v1.0.6-beta.1"),
		},
		&testEntry{
			compatVer: parseVersion("v1.0.0"),
			appVer:    parseVersion("v1.0.8-alpha.1"),
		},
	})

	anchor := &testCompatQueryAnchor{ver: parseVersion("1.0.1-alpha.3")}
	strictCompatResult := tree.StrictCompat(anchor)

	assert.Equal(t, "1.0.8.1011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
	assert.Equal(t, "1.0.8.1011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

	anchor = &testCompatQueryAnchor{ver: parseVersion("1.0.1-beta.3")}
	strictCompatResult = tree.StrictCompat(anchor)

	assert.Equal(t, "1.0.6.2011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
	assert.Equal(t, "1.0.6.2011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

	anchor = &testCompatQueryAnchor{ver: parseVersion("1.0.1-rc.3")}
	strictCompatResult = tree.StrictCompat(anchor)

	assert.Equal(t, "1.0.4.3011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
	assert.Equal(t, "1.0.4.3011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

	anchor = &testCompatQueryAnchor{ver: parseVersion("1.0.1-release.3")}
	strictCompatResult = tree.StrictCompat(anchor)

	assert.Equal(t, "1.0.2.4011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
	assert.Equal(t, "1.0.2.4011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

	anchor = &testCompatQueryAnchor{ver: parseVersion("1.0.6-release.3")}
	strictCompatResult = tree.StrictCompat(anchor)

	assert.Nil(t, strictCompatResult.CanUpdateVersion())
	assert.Nil(t, strictCompatResult.LatestVersion())

	anchor = &testCompatQueryAnchor{ver: parseVersion("1.0.10")}
	strictCompatResult = tree.StrictCompat(anchor)

	assert.Nil(t, strictCompatResult.CanUpdateVersion())
	assert.Nil(t, strictCompatResult.LatestVersion())
}

func TestMultiCompatRanges(t *testing.T) {
	/*
		Intersect ranges:
			|----|
			   |---|
				   |---|
	*/
	t.Run("intersect", func(t *testing.T) {
		tree := NewVersionCompatTree(nil)

		tree.Publish(Entries{
			&testEntry{
				compatVer: parseVersion("v1.0.0"),
				appVer:    parseVersion("v1.0.0"),
			},
			&testEntry{
				compatVer: parseVersion("v1.0.0"),
				appVer:    parseVersion("v1.0.2"),
			},
			&testEntry{
				compatVer: parseVersion("v1.0.0"),
				appVer:    parseVersion("v1.0.4"),
			},
			&testEntry{
				compatVer: parseVersion("v1.0.0"),
				appVer:    parseVersion("v1.0.6"),
			},
		}...)

		tree.Publish(Entries{
			&testEntry{
				compatVer: parseVersion("v1.0.4"),
				appVer:    parseVersion("v1.1.0"),
			},
			&testEntry{
				compatVer: parseVersion("v1.0.4"),
				appVer:    parseVersion("v1.1.2"),
			},
			&testEntry{
				compatVer: parseVersion("v1.0.4"),
				appVer:    parseVersion("v1.1.4"),
			},
		}...)

		tree.Publish(Entries{
			&testEntry{
				compatVer: parseVersion("v1.1.4"),
				appVer:    parseVersion("v1.2.0"),
			},
			&testEntry{
				compatVer: parseVersion("v1.1.4"),
				appVer:    parseVersion("v1.2.2"),
			},
			&testEntry{
				compatVer: parseVersion("v1.1.4"),
				appVer:    parseVersion("v1.2.4"),
			},
		}...)

		anchor := &testCompatQueryAnchor{ver: parseVersion("1.0.3")}
		strictCompatResult := tree.StrictCompat(anchor)

		assert.Equal(t, "1.0.6.4011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
		assert.Equal(t, "1.2.4.4011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

		anchor = &testCompatQueryAnchor{ver: parseVersion("1.0.5")}
		strictCompatResult = tree.StrictCompat(anchor)

		assert.Equal(t, "1.1.4.4011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
		assert.Equal(t, "1.2.4.4011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

		anchor = &testCompatQueryAnchor{ver: parseVersion("1.1.0")}
		strictCompatResult = tree.StrictCompat(anchor)

		assert.Equal(t, "1.1.4.4011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
		assert.Equal(t, "1.2.4.4011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

		anchor = &testCompatQueryAnchor{ver: parseVersion("1.1.4")}
		strictCompatResult = tree.StrictCompat(anchor)

		assert.Equal(t, "1.2.4.4011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
		assert.Equal(t, "1.2.4.4011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

		anchor = &testCompatQueryAnchor{ver: parseVersion("1.2.5")}
		strictCompatResult = tree.StrictCompat(anchor)

		assert.Nil(t, strictCompatResult.CanUpdateVersion())
		assert.Nil(t, strictCompatResult.LatestVersion())
	})

	/*
		Contain ranges
			|------|
			  |--|
	*/
	t.Run("contain", func(t *testing.T) {
		tree := NewVersionCompatTree(nil)

		tree.Publish(Entries{
			&testEntry{
				compatVer: parseVersion("v1.0.0"),
				appVer:    parseVersion("v1.0.8"),
			},
		}...)

		tree.Publish(Entries{
			&testEntry{
				compatVer: parseVersion("v1.0.4"),
				appVer:    parseVersion("v1.0.6"),
			},
		}...)

		anchor := &testCompatQueryAnchor{ver: parseVersion("1.0.3")}
		strictCompatResult := tree.StrictCompat(anchor)

		assert.Equal(t, "1.0.8.4011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
		assert.Equal(t, "1.0.8.4011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

		anchor = &testCompatQueryAnchor{ver: parseVersion("1.0.5")}
		strictCompatResult = tree.StrictCompat(anchor)

		assert.Equal(t, "1.0.8.4011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
		assert.Equal(t, "1.0.8.4011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

		anchor = &testCompatQueryAnchor{ver: parseVersion("1.0.7")}
		strictCompatResult = tree.StrictCompat(anchor)

		assert.Equal(t, "1.0.8.4011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
		assert.Equal(t, "1.0.8.4011", rawVersionOfEntry(strictCompatResult.LatestVersion()))
	})

	/*
		Disjoint ranges
			|---|
				  |---|
	*/
	t.Run("disjoint", func(t *testing.T) {
		tree := NewVersionCompatTree(nil)

		tree.Publish(Entries{
			&testEntry{
				compatVer: parseVersion("v1.0.0"),
				appVer:    parseVersion("v1.0.8"),
			},
		}...)

		tree.Publish(Entries{
			&testEntry{
				compatVer: parseVersion("v1.0.10"),
				appVer:    parseVersion("v1.0.20"),
			},
		}...)

		anchor := &testCompatQueryAnchor{ver: parseVersion("1.0.6")}
		strictCompatResult := tree.StrictCompat(anchor)

		assert.Equal(t, "1.0.8.4011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
		assert.Equal(t, "1.0.20.4011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

		anchor = &testCompatQueryAnchor{ver: parseVersion("1.0.9")}
		strictCompatResult = tree.StrictCompat(anchor)

		assert.Nil(t, strictCompatResult.CanUpdateVersion())
		assert.Equal(t, "1.0.20.4011", rawVersionOfEntry(strictCompatResult.LatestVersion()))

		anchor = &testCompatQueryAnchor{ver: parseVersion("1.0.16")}
		strictCompatResult = tree.StrictCompat(anchor)

		assert.Equal(t, "1.0.20.4011", rawVersionOfEntry(strictCompatResult.CanUpdateVersion()))
		assert.Equal(t, "1.0.20.4011", rawVersionOfEntry(strictCompatResult.LatestVersion()))
	})
}
