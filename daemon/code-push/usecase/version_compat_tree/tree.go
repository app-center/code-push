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

				switch entry.CompatVersion().Compare(entryInList.CompatVersion()) {
				case semver.CompareLargeFlag:
					e := e.Next()
					if e == nil {
						v.tree.PushBack(entry)
						i++
					}
				case semver.CompareEqualFlag:
					switch entry.Version().Compare(entryInList.Version()) {
					case semver.CompareEqualFlag:
						e = nil
						i++
					case semver.CompareLessFlag:
						e = nil
						i++
					case semver.CompareLargeFlag:
						fallthrough
					default:
						nextElement := e.Next()

						v.tree.Remove(e)

						e = nextElement
						if e == nil {
							v.tree.PushBack(entry)
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
	panic("implement me")
}

func (v *versionRangeTree) LooseCompat(anchor ICompatQueryAnchor) ICompatQueryResult {
	panic("implement me")
}

func New(entries []IEntry) ITree {
	tree := &versionRangeTree{
		tree: list.New(),
	}

	tree.Add(entries...)

	return tree
}
