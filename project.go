package main

import "sort"

// Project Structure
type Project struct {
	URL          string
	URLID        int
	ProjectName  string
	SubjectCount int
	Versions     []*ProjectVersion
}

// retrieve a project Version by CRF Version
func (pj *Project) getVersionByID(crfVersion int) *ProjectVersion {
	for _, version := range pj.Versions {
		if version.CRFVersionID == crfVersion {
			return version
		}
	}
	return nil
}

// Get the Last Version
func (pj *Project) getLastVersion() *ProjectVersion {
	for _, version := range pj.Versions {
		if version.LastVersion {
			return version
		}
	}
	return nil
}

// Sorter
type ByPV func(v1, v2 *ProjectVersion) bool

func (by ByPV) Sort(versions []*ProjectVersion) {
	vs := &versionSorter{
		versions: versions,
		by:       by,
	}
	sort.Sort(vs)
}

type versionSorter struct {
	versions []*ProjectVersion
	by       func(v1, v2 *ProjectVersion) bool
}

// Len is part of sort.Interface.
func (s *versionSorter) Len() int {
	return len(s.versions)
}

// Swap is part of sort.Interface.
func (s *versionSorter) Swap(i, j int) {
	s.versions[i], s.versions[j] = s.versions[j], s.versions[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *versionSorter) Less(i, j int) bool {
	return s.by(s.versions[i], s.versions[j])
}

// order the project versions by CRFVersion
func (pj *Project) getVersions() []*ProjectVersion {
	crfVersion := func(v1, v2 *ProjectVersion) bool {
		return v1.CRFVersionID < v2.CRFVersionID
	}
	ByPV(crfVersion).Sort(pj.Versions)
	return pj.Versions
}
