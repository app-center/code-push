package versionCompatTree

import (
	"container/list"
	"github.com/funnyecho/code-push/pkg/semver"
	"sync"
)

type branchEntryList struct {
	mut           *sync.RWMutex
	list          *list.List
	entryElements map[IEntry]*list.Element
}

func (l *branchEntryList) Dequeue(entry IEntry) (ok bool) {
	l.mut.Lock()
	defer l.mut.Unlock()

	ok = l.dequeue(entry)
	return
}

func (l *branchEntryList) dequeue(entry IEntry) (ok bool) {
	if versionElement, versionExisted := l.entryElements[entry]; versionExisted {
		l.list.Remove(versionElement)
		delete(l.entryElements, entry)
		ok = true
	}

	return
}

func (l *branchEntryList) Insert(entry IEntry) {
	l.mut.Lock()
	defer l.mut.Unlock()

	l.insert(entry)
}

func (l *branchEntryList) insert(entry IEntry) {
	for element := l.list.Front(); element != nil; element = element.Next() {
		e := element.Value.(IEntry)

		switch e.Version().StageSafetyLooseCompare(entry.Version()) {
		case semver.CompareEqualFlag:
			return
		case semver.CompareLessFlag:
			continue
		case semver.CompareLargeFlag:
			ele := l.list.InsertBefore(entry, element)
			l.entryElements[entry] = ele
			return
		}
	}

	ele := l.list.PushBack(entry)
	l.entryElements[entry] = ele
}

func (l *branchEntryList) Enqueue(entry IEntry) (evictEntries []IEntry, ignored bool) {
	l.mut.Lock()
	defer l.mut.Unlock()

	evictEntries, ignored = l.enqueue(entry)
	return
}

func (l *branchEntryList) enqueue(entry IEntry) (evictEntries []IEntry, ignored bool) {
	ignored = true

	for element := l.list.Front(); element != nil; {
		e := element.Value.(IEntry)

		looseCompare := e.Version().StageSafetyLooseCompare(entry.Version())
		strictCompare := e.Version().StageSafetyStrictCompare(entry.Version())

		if looseCompare == semver.CompareLessFlag && strictCompare == semver.CompareLessFlag {
			evictEntries = append(evictEntries, e)
			l.dequeue(e)
			element = element.Next()
			continue
		} else if looseCompare == semver.CompareLargeFlag && strictCompare == semver.CompareLessFlag {
			ele := l.list.InsertBefore(entry, element)
			l.entryElements[entry] = ele
			ignored = false
			return
		} else if looseCompare == semver.CompareLessFlag && strictCompare == semver.CompareLargeFlag {
			element = element.Next()
			continue
		} else if looseCompare == semver.CompareLargeFlag && strictCompare == semver.CompareLargeFlag {
			return
		} else {
			return
		}
	}

	ele := l.list.PushBack(entry)
	l.entryElements[entry] = ele
	ignored = false

	return
}

func (l *branchEntryList) StrictCompatQuery(ver *semver.SemVer) IEntry {
	l.mut.RLock()
	defer l.mut.RUnlock()

	return l.strictCompatQuery(ver)
}

func (l *branchEntryList) strictCompatQuery(ver *semver.SemVer) IEntry {
	for element := l.list.Back(); element != nil; element = element.Prev() {
		e := element.Value.(IEntry)

		hiStrictCompare := e.Version().StageSafetyStrictCompare(ver)
		hiLooseCompare := e.Version().StageSafetyLooseCompare(ver)
		loLooseCompare := e.CompatVersion().StageSafetyLooseCompare(ver)
		if hiLooseCompare == semver.CompareLargeFlag {
			if hiStrictCompare == semver.CompareLessFlag {
				continue
			} else if loLooseCompare == semver.CompareLessFlag || loLooseCompare == semver.CompareEqualFlag {
				return e
			} else {
				continue
			}
		} else {
			return nil
		}
	}

	return nil
}

func (l *branchEntryList) StrictLatestQuery(ver *semver.SemVer) IEntry {
	l.mut.RLock()
	defer l.mut.RUnlock()

	return l.strictLatestQuery(ver)
}

func (l *branchEntryList) strictLatestQuery(ver *semver.SemVer) IEntry {
	for element := l.list.Back(); element != nil; element = element.Prev() {
		e := element.Value.(IEntry)

		hiStrictCompare := e.Version().StageSafetyStrictCompare(ver)
		hiLooseCompare := e.Version().StageSafetyLooseCompare(ver)

		if hiLooseCompare == semver.CompareLargeFlag {
			if hiStrictCompare == semver.CompareLessFlag {
				continue
			} else {
				return e
			}
		} else {
			return nil
		}
	}

	return nil
}

func newBranchEntryList() *branchEntryList {
	return &branchEntryList{
		mut:           &sync.RWMutex{},
		list:          list.New(),
		entryElements: map[IEntry]*list.Element{},
	}
}

type compatVersionBranch struct {
	entryList *branchEntryList
}

func newCompatVersionBranch() *compatVersionBranch {
	return &compatVersionBranch{
		entryList: newBranchEntryList(),
	}
}

func (b *compatVersionBranch) Enqueue(entry IEntry) (evictEntries []IEntry, ignored bool) {
	return b.entryList.Enqueue(entry)
}

type compatVersionBranchMap map[string]*compatVersionBranch

func newCompatVersionBranchMap() compatVersionBranchMap {
	return make(compatVersionBranchMap)
}

func (m compatVersionBranchMap) Enqueue(entry IEntry) (evictEntries []IEntry, ignored bool) {
	rawVersion := entry.CompatVersion().String()

	branch, hasBranch := m[rawVersion]

	if !hasBranch {
		branch = newCompatVersionBranch()
		m[rawVersion] = branch
	}

	return branch.Enqueue(entry)
}

type latestVersionBranch struct {
	entryList *branchEntryList
}

func newLatestVersionBranch() *latestVersionBranch {
	return &latestVersionBranch{
		entryList: newBranchEntryList(),
	}
}

func (b *latestVersionBranch) Enqueue(entry IEntry) {
	b.entryList.Enqueue(entry)
}

func (b *latestVersionBranch) StrictLatest(anchor ICompatQueryAnchor) IEntry {
	return b.entryList.StrictLatestQuery(anchor.Version())
}

type strictCompatQueryBranch struct {
	entryList *branchEntryList
}

func newStrictCompatQueryBranch() *strictCompatQueryBranch {
	return &strictCompatQueryBranch{
		entryList: newBranchEntryList(),
	}
}

func (b *strictCompatQueryBranch) Insert(entry IEntry) {
	b.entryList.Insert(entry)
}

func (b *strictCompatQueryBranch) Evict(entry IEntry) (ok bool) {
	return b.entryList.Dequeue(entry)
}

func (b *strictCompatQueryBranch) Query(anchor ICompatQueryAnchor) IEntry {
	return b.entryList.StrictCompatQuery(anchor.Version())
}
