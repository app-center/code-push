package version_compat_tree

type versionRangeTree struct {
}

func (v *versionRangeTree) Add(entries ...IEntry) {
	panic("implement me")
}

func (v *versionRangeTree) StrictCompat(anchor ICompatQueryAnchor) ICompatQueryResult {
	panic("implement me")
}

func (v *versionRangeTree) LooseCompat(anchor ICompatQueryAnchor) ICompatQueryResult {
	panic("implement me")
}

func New() ITree {
	return &versionRangeTree{}
}
