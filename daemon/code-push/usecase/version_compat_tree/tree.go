package version_compat_tree

import (
	"container/list"
	"github.com/funnyecho/code-push/pkg/semver"
)

type versionRangeTree struct {
	tree *list.List
}

func (v *versionRangeTree) Add(entries ...IEntry) {
	for i := 0; i < len(entries); {
		entry := entries[i]

		if v.tree.Len() == 0 {
			v.tree.PushFront(entry)
			i++
		} else {
			for e := v.tree.Front(); e != nil; {
				entryInList := e.Value.(IEntry)

				switch entry.CompatVersion().StageSafetyLooseCompare(entryInList.CompatVersion()) {
				case semver.CompareLargeFlag:
					e := e.Next()
					if e == nil {
						v.tree.PushBack(entry)
						i++
					}
				case semver.CompareEqualFlag:
					strictCompare := entry.Version().StageSafetyStrictCompare(entryInList.Version())
					looseCompare := entry.Version().StageSafetyLooseCompare(entryInList.Version())

					switch strictCompare {
					case semver.CompareEqualFlag:
						e = nil
						i++
					case semver.CompareLessFlag:
						if looseCompare == semver.CompareLargeFlag {
							e = e.Next()
							if e == nil {
								v.tree.PushBack(entry)
								i++
							}
						} else if looseCompare == semver.CompareLessFlag {
							e = nil
							i++
						}
					case semver.CompareLargeFlag:
						fallthrough
					default:
						if looseCompare == semver.CompareLargeFlag {
							nextElement := e.Next()

							v.tree.Remove(e)

							e = nextElement
							if e == nil {
								v.tree.PushBack(entry)
								i++
							}
						} else if looseCompare == semver.CompareLessFlag {
							v.tree.InsertBefore(entry, e)
							e = nil
							i++
						}
					}
				case semver.CompareLessFlag:
					fallthrough
				default:
					v.tree.InsertBefore(entry, e)
					e = nil
					i++
				}
			}
		}
	}
}

func (v *versionRangeTree) StrictCompat(anchor ICompatQueryAnchor) ICompatQueryResult {
	queryResult := &compatQueryResult{}

	version := anchor.Version()

	for e := v.tree.Back(); e != nil; e = e.Prev() {
		entryInList := e.Value.(IEntry)

		compatVersionCompare := entryInList.CompatVersion().StageSafetyLooseCompare(version)
		versionCompare := entryInList.Version().StageSafetyStrictCompare(version)

		if versionCompare == semver.CompareLargeFlag && compatVersionCompare != semver.CompareLargeFlag {
			queryResult.canUpdateVersion = entryInList
			break
		}
	}

	return queryResult
}

func (v *versionRangeTree) LooseCompat(anchor ICompatQueryAnchor) ICompatQueryResult {
	queryResult := &compatQueryResult{}

	version := anchor.Version()

	for e := v.tree.Back(); e != nil; e = e.Prev() {
		entryInList := e.Value.(IEntry)

		compatVersionCompare := entryInList.CompatVersion().StageSafetyLooseCompare(version)
		versionCompare := entryInList.Version().StageSafetyLooseCompare(version)

		if versionCompare == semver.CompareLessFlag {
			if queryResult.latestVersion == nil {
				queryResult.latestVersion = entryInList
			}
			continue
		}

		if compatVersionCompare == semver.CompareLargeFlag {
			continue
		}

		queryResult.canUpdateVersion = entryInList
		if queryResult.latestVersion == nil {
			queryResult.latestVersion = entryInList
		}

		break
	}

	return queryResult
}

func NewVersionCompatTree(entries []IEntry) ITree {
	tree := &versionRangeTree{
		tree: list.New(),
	}

	if entries != nil {
		tree.Add(entries...)
	}

	return tree
}

type compatQueryResult struct {
	latestVersion    IEntry
	canUpdateVersion IEntry
}

func (r *compatQueryResult) LatestVersion() IEntry {
	return r.latestVersion
}

func (r *compatQueryResult) CanUpdateVersion() IEntry {
	return r.canUpdateVersion
}
