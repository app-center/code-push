package version_compat_tree

import "github.com/funnyecho/code-push/pkg/semver"

type IEntry interface {
	CompatVersion() *semver.SemVer
	Version() *semver.SemVer
}

type Entries []IEntry

type ICompatQueryAnchor interface {
	Version() *semver.SemVer
}

type ICompatQueryResult interface {
	LatestVersion() IEntry
	CanUpdateVersion() IEntry
}

type ITree interface {
	Add(entries ...IEntry)
	StrictCompat(anchor ICompatQueryAnchor) ICompatQueryResult
	LooseCompat(anchor ICompatQueryAnchor) ICompatQueryResult
}
