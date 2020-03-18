package main

import (
	"github.com/jmoiron/sqlx"
	"sort"
)

// Project Structure
type Project struct {
	URLID               int    `db:"url_id"`
	ProjectID           int    `db:"project_id"`
	ProjectName         string `db:"project_name"`
	SubjectCount        SubjectCount
	Versions            []*ProjectVersion
	UnusedWithOpenQuery []*UnusedEdit
	Unused              []*UnusedEdit
}

// load the subject count for a Project
func (pv *Project) loadSubjectCount(db *sqlx.DB, urlID int) {
	pv.SubjectCount = getProjectSubjectCount(db, urlID, pv.ProjectID)
}

// load the useless edits
func (pj *Project) loadUnusedQueries(db *sqlx.DB) {
	// with OpenQuery
	openQuery := getUselessEditsForProject(db, pj.ProjectID, OpenQuery)
	noOpenQuery := getUselessEditsForProject(db, pj.ProjectID, WithoutOpenQuery)
	pj.Unused = noOpenQuery
	pj.UnusedWithOpenQuery = openQuery
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
type ByP func(v1, v2 *Project) bool

func (by ByP) Sort(projects []*Project) {
	vs := &projectSorter{
		projects: projects,
		by:       by,
	}
	sort.Sort(vs)
}

type projectSorter struct {
	projects []*Project
	by       func(v1, v2 *Project) bool
}

// Len is part of sort.Interface.
func (s *projectSorter) Len() int {
	return len(s.projects)
}

// Swap is part of sort.Interface.
func (s *projectSorter) Swap(i, j int) {
	s.projects[i], s.projects[j] = s.projects[j], s.projects[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *projectSorter) Less(i, j int) bool {
	return s.by(s.projects[i], s.projects[j])
}

func orderProjects(prj []*Project) []*Project {
	name := func(p1, p2 *Project) bool {
		return p1.ProjectName < p2.ProjectName
	}
	ByP(name).Sort(prj)
	return prj
}
