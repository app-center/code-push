package version_compat_tree

type versionRangeTree struct {
	compatVersionBranchMap
	*latestVersionBranch
	*strictCompatQueryBranch
}

func (v *versionRangeTree) Publish(entries ...IEntry) {
	for _, entry := range entries {
		evictEntryList, ignored := v.compatVersionBranchMap.Enqueue(entry)
		if evictEntryList != nil {
			for _, entry := range evictEntryList {
				v.strictCompatQueryBranch.Evict(entry)
			}
		}

		if !ignored {
			v.latestVersionBranch.Enqueue(entry)
			v.strictCompatQueryBranch.Insert(entry)
		}
	}
}

func (v *versionRangeTree) StrictCompat(anchor ICompatQueryAnchor) ICompatQueryResult {
	queryResult := &compatQueryResult{}

	canUpdateEntry := v.strictCompatQueryBranch.Query(anchor)
	strictLatest := v.latestVersionBranch.StrictLatest(anchor)

	queryResult.canUpdateVersion = canUpdateEntry
	queryResult.latestVersion = strictLatest

	return queryResult
}

func NewVersionCompatTree() ITree {
	tree := &versionRangeTree{
		compatVersionBranchMap:  newCompatVersionBranchMap(),
		latestVersionBranch:     newLatestVersionBranch(),
		strictCompatQueryBranch: newStrictCompatQueryBranch(),
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
