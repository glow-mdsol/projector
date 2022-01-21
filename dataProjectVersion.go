package main

import (
	"github.com/jmoiron/sqlx"
	"sort"
)

type EditStatusCounts struct {
	ActiveEdits   int `db:"active_count"`
	InactiveEdits int `db:"inactive_count"`
}

// ProjectVersion represents the structure for an individual Project Version
type ProjectVersion struct {
	ProjectID          int  `db:"project_id"`
	CRFVersionID       int  `db:"crf_version_id"`
	LastVersion        bool `db:"last_version"`
	EditStatus         EditStatusCounts
	FieldEditMetrics   EditTypeMetric
	ProgramEditMetrics EditTypeMetric
	ActiveCheckCount   int
	InActiveCheckCount int
}

// Total count of edits
func (pv *ProjectVersion) getTotalEdits() int {
	return pv.FieldEditMetrics.TotalEdits + pv.ProgramEditMetrics.TotalEdits
}

// load the check counts
func (pv *ProjectVersion) getActivityCounts(db *sqlx.DB) {
	pv.EditStatus = getActivityCount(db, pv.ProjectID, pv.CRFVersionID)
}

// load the metrics
func (pv *ProjectVersion) getMetrics(db *sqlx.DB) {
	// field edits
	fieldEdits := getStudyMetricsByProjectAndCheckType(db, pv.ProjectID, pv.CRFVersionID, Field)
	programmedEdits := getStudyMetricsByProjectAndCheckType(db, pv.ProjectID, pv.CRFVersionID, Programmed)
	// impute the raw values
	fieldEdits.fixUpMetrics()
	programmedEdits.fixUpMetrics()
	// calculate the percentages
	fieldEdits.calculatePercentages()
	programmedEdits.calculatePercentages()
	// set the values
	pv.FieldEditMetrics = fieldEdits
	pv.ProgramEditMetrics = programmedEdits
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

func (pj *Project) sortVersions() []*ProjectVersion {
	crfVersion := func(v1, v2 *ProjectVersion) bool {
		return v1.CRFVersionID < v2.CRFVersionID
	}
	ByPV(crfVersion).Sort(pj.Versions)
	return pj.Versions
}

// order the project versions by CRFVersion
func orderVersions(pj []*ProjectVersion) []*ProjectVersion {
	crfVersion := func(v1, v2 *ProjectVersion) bool {
		return v1.CRFVersionID < v2.CRFVersionID
	}
	ByPV(crfVersion).Sort(pj)
	return pj
}

//func (pv *ProjectVersion) fixUpNullValues() {
//	var activeCount = 0
//	var inactiveCount = 0
//	if pv.ActiveEditsOnly != nil {
//		//log.Printf("INFO: Generating Active Edits for ProjectVersion %s(%d)", pv.ProjectName, pv.CRFVersionID)
//		if pv.ActiveEditsOnly.FieldEditMetrics.TotalEdits >= 0 {
//			activeCount += pv.ActiveEditsOnly.FieldEditMetrics.TotalEdits
//		}
//		if pv.ActiveEditsOnly.ProgramEditMetrics.TotalEdits >= 0 {
//			activeCount += pv.ActiveEditsOnly.ProgramEditMetrics.TotalEdits
//		}
//	}
//	if pv.InactiveEditsOnly != nil {
//		//log.Printf("INFO: Generating InActive Edits for ProjectVersion %s(%d)", pv.ProjectName, pv.CRFVersionID)
//		if pv.InactiveEditsOnly.FieldEditMetrics.TotalEdits >= 0 {
//			inactiveCount += pv.InactiveEditsOnly.FieldEditMetrics.TotalEdits
//		}
//		if pv.InactiveEditsOnly.ProgramEditMetrics.TotalEdits >= 0 {
//			inactiveCount += pv.InactiveEditsOnly.ProgramEditMetrics.TotalEdits
//		}
//	}
//	pv.ActiveEditCount = activeCount
//	pv.InactiveEditCount = inactiveCount
//}
