package main

import (
	"database/sql"
	"github.com/lib/pq"
)

type UnusedEdit struct {
	URL           string
	ProjectName   string        `db:"project_name"`
	EditCheckName string        `db:"edit_check_name"`
	UsageCount    int            `db:"edit_check_count"`
	OpenQuery     string            `db:"open_query"`
}

type SubjectCount struct {
	URL          string
	ProjectName  string        `db:"project_name"`
	SubjectCount sql.NullInt64 `db:"subject_count"`
	RefreshDate  pq.NullTime    `db:"refresh_date"`
}

type Record struct {
	URL                               string
	ProjectName                       string        `db:"project_name"`
	CRFVersionID                      string        `db:"crf_version_id"`
	LastVersion                       bool        	`db:"last_version"`
	SubjectCount                      sql.NullInt64 `db:"subject_count"`
	CheckStatus                       string        `db:"check_status"`
	RawTotalFieldEdits                sql.NullInt64 `db:"total_edits_fld"`
	TotalFieldEdits                   int
	RawTotalProgEdits                 sql.NullInt64 `db:"total_edits_prg"`
	TotalProgEdits                    int
	RawTotalFieldQueries              sql.NullInt64 `db:"total_queries_fld"`
	TotalFieldQueries                 int
	RawTotalProgQueries               sql.NullInt64 `db:"total_queries_prg"`
	TotalProgQueries                  int
	TotalFieldEditsWithOpenQuery      int           `db:"total_edits_query_fld"`
	TotalProgEditsWithOpenQuery       int           `db:"total_edits_query_prg"`
	RawTotalFieldQueriesWithOpenQuery sql.NullInt64 `db:"total_queries_query_fld"`
	RawTotalProgQueriesWithOpenQuery  sql.NullInt64 `db:"total_queries_query_prg"`
	TotalFieldQueriesWithOpenQuery    int
	TotalProgQueriesWithOpenQuery     int
	TotalFieldEditsFired              int           `db:"total_fired_fld"`
	TotalProgEditsFired               int           `db:"total_fired_prg"`
	RawTotalFieldWithOpenQueryFired   sql.NullInt64 `db:"total_fired_query_fld"`
	RawTotalProgWithOpenQueryFired    sql.NullInt64 `db:"total_fired_query_prg"`
	TotalFieldWithOpenQueryFired      int
	TotalProgWithOpenQueryFired       int
	TotalFieldEditsFiredWithNoChange  int 			`db:"fired_no_change_fld"`
	TotalProgEditsFiredWithNoChange   int 			`db:"fired_no_change_prg"`
}

type ProjectVersion struct {
	URL               string
	ProjectName       string
	CRFVersionID      string
	AllEdits          Record
	ActiveEditsOnly   Record
	InactiveEditsOnly Record
	LastVersion       bool
	SubjectCount      int
}


func calculateInactiveCounts(pv *ProjectVersion) *ProjectVersion {
	rec := new(Record)
	rec.URL = pv.URL
	rec.ProjectName = pv.ProjectName
	rec.CRFVersionID = pv.CRFVersionID
	rec.CheckStatus = "INACTIVE"
	rec.TotalFieldEdits = pv.AllEdits.TotalFieldEdits - pv.ActiveEditsOnly.TotalFieldEdits
	rec.TotalProgEdits = pv.AllEdits.TotalProgEdits - pv.ActiveEditsOnly.TotalProgEdits
	rec.TotalFieldQueries = pv.AllEdits.TotalFieldQueries - pv.ActiveEditsOnly.TotalFieldQueries
	rec.TotalProgQueries = pv.AllEdits.TotalProgQueries - pv.ActiveEditsOnly.TotalProgQueries
	rec.TotalFieldEditsWithOpenQuery = pv.AllEdits.TotalFieldEditsWithOpenQuery - pv.ActiveEditsOnly.TotalFieldEditsWithOpenQuery
	rec.TotalProgEditsWithOpenQuery = pv.AllEdits.TotalProgEditsWithOpenQuery - pv.ActiveEditsOnly.TotalProgEditsWithOpenQuery
	rec.TotalFieldQueriesWithOpenQuery = pv.AllEdits.TotalFieldQueriesWithOpenQuery - pv.ActiveEditsOnly.TotalFieldQueriesWithOpenQuery
	rec.TotalProgQueriesWithOpenQuery = pv.AllEdits.TotalProgQueriesWithOpenQuery - pv.ActiveEditsOnly.TotalProgQueriesWithOpenQuery
	rec.TotalFieldEditsFired = pv.AllEdits.TotalFieldEditsFired - pv.ActiveEditsOnly.TotalFieldEditsFired
	rec.TotalProgEditsFired = pv.AllEdits.TotalProgEditsFired - pv.ActiveEditsOnly.TotalProgEditsFired
	rec.TotalFieldEditsFired = pv.AllEdits.TotalFieldEditsFired - pv.ActiveEditsOnly.TotalFieldEditsFired
	rec.TotalFieldWithOpenQueryFired = pv.AllEdits.TotalFieldWithOpenQueryFired - pv.ActiveEditsOnly.TotalFieldWithOpenQueryFired
	rec.TotalProgWithOpenQueryFired = pv.AllEdits.TotalProgWithOpenQueryFired - pv.ActiveEditsOnly.TotalProgWithOpenQueryFired
	rec.TotalFieldEditsFiredWithNoChange = pv.AllEdits.TotalFieldEditsFiredWithNoChange - pv.ActiveEditsOnly.TotalFieldEditsFiredWithNoChange
	rec.TotalProgEditsFiredWithNoChange = pv.AllEdits.TotalProgEditsFiredWithNoChange - pv.ActiveEditsOnly.TotalProgEditsFiredWithNoChange
	pv.InactiveEditsOnly = *rec
	return pv
}

func (pv *ProjectVersion) calculateInactiveCounts(){
	rec := new(Record)
	rec.URL = pv.URL
	rec.ProjectName = pv.ProjectName
	rec.CRFVersionID = pv.CRFVersionID
	rec.CheckStatus = "INACTIVE"
	rec.TotalFieldEdits = pv.AllEdits.TotalFieldEdits - pv.ActiveEditsOnly.TotalFieldEdits
	rec.TotalProgEdits = pv.AllEdits.TotalProgEdits - pv.ActiveEditsOnly.TotalProgEdits
	rec.TotalFieldQueries = pv.AllEdits.TotalFieldQueries - pv.ActiveEditsOnly.TotalFieldQueries
	rec.TotalProgQueries = pv.AllEdits.TotalProgQueries - pv.ActiveEditsOnly.TotalProgQueries
	rec.TotalFieldEditsWithOpenQuery = pv.AllEdits.TotalFieldEditsWithOpenQuery - pv.ActiveEditsOnly.TotalFieldEditsWithOpenQuery
	rec.TotalProgEditsWithOpenQuery = pv.AllEdits.TotalProgEditsWithOpenQuery - pv.ActiveEditsOnly.TotalProgEditsWithOpenQuery
	rec.TotalFieldQueriesWithOpenQuery = pv.AllEdits.TotalFieldQueriesWithOpenQuery - pv.ActiveEditsOnly.TotalFieldQueriesWithOpenQuery
	rec.TotalProgQueriesWithOpenQuery = pv.AllEdits.TotalProgQueriesWithOpenQuery - pv.ActiveEditsOnly.TotalProgQueriesWithOpenQuery
	rec.TotalFieldEditsFired = pv.AllEdits.TotalFieldEditsFired - pv.ActiveEditsOnly.TotalFieldEditsFired
	rec.TotalProgEditsFired = pv.AllEdits.TotalProgEditsFired - pv.ActiveEditsOnly.TotalProgEditsFired
	rec.TotalFieldEditsFired = pv.AllEdits.TotalFieldEditsFired - pv.ActiveEditsOnly.TotalFieldEditsFired
	rec.TotalFieldWithOpenQueryFired = pv.AllEdits.TotalFieldWithOpenQueryFired - pv.ActiveEditsOnly.TotalFieldWithOpenQueryFired
	rec.TotalProgWithOpenQueryFired = pv.AllEdits.TotalProgWithOpenQueryFired - pv.ActiveEditsOnly.TotalProgWithOpenQueryFired
	rec.TotalFieldEditsFiredWithNoChange = pv.AllEdits.TotalFieldEditsFiredWithNoChange - pv.ActiveEditsOnly.TotalFieldEditsFiredWithNoChange
	rec.TotalProgEditsFiredWithNoChange = pv.AllEdits.TotalProgEditsFiredWithNoChange - pv.ActiveEditsOnly.TotalProgEditsFiredWithNoChange
	pv.InactiveEditsOnly = *rec
}

func fixUpRecord(rec Record) Record {
	if rec.RawTotalFieldQueriesWithOpenQuery.Valid {
		rec.TotalFieldQueriesWithOpenQuery = int(rec.RawTotalFieldQueriesWithOpenQuery.Int64)
	} else {
		rec.TotalFieldQueriesWithOpenQuery = -1
	}
	if rec.RawTotalProgQueriesWithOpenQuery.Valid {
		rec.TotalProgQueriesWithOpenQuery = int(rec.RawTotalProgQueriesWithOpenQuery.Int64)
	} else {
		rec.TotalProgQueriesWithOpenQuery = -1
	}
	if rec.RawTotalFieldEdits.Valid {
		rec.TotalFieldEdits = int(rec.RawTotalFieldEdits.Int64)
	} else {
		rec.TotalFieldEdits = -1
	}
	if rec.RawTotalProgEdits.Valid {
		rec.TotalProgEdits = int(rec.RawTotalProgEdits.Int64)
	} else {
		rec.TotalProgEdits = -1
	}
	if rec.RawTotalFieldQueries.Valid {
		rec.TotalFieldQueries = int(rec.RawTotalFieldQueries.Int64)
	} else {
		rec.TotalFieldQueries = -1
	}
	if rec.RawTotalProgQueries.Valid {
		rec.TotalProgQueries = int(rec.RawTotalProgQueries.Int64)
	} else {
		rec.TotalProgQueries = -1
	}
	if rec.RawTotalFieldWithOpenQueryFired.Valid {
		rec.TotalFieldWithOpenQueryFired = int(rec.RawTotalFieldWithOpenQueryFired.Int64)
	} else {
		rec.TotalFieldWithOpenQueryFired = -1
	}
	if rec.RawTotalProgWithOpenQueryFired.Valid {
		rec.TotalProgWithOpenQueryFired = int(rec.RawTotalProgWithOpenQueryFired.Int64)
	} else {
		rec.TotalProgWithOpenQueryFired = -1
	}
	return rec
}

func fixUpNullValues(pv *ProjectVersion) *ProjectVersion {
	pv.ActiveEditsOnly = fixUpRecord(pv.ActiveEditsOnly)
	pv.AllEdits = fixUpRecord(pv.AllEdits)
	return pv
}

func (pv *ProjectVersion) fixUpNullValues(){
	pv.ActiveEditsOnly = fixUpRecord(pv.ActiveEditsOnly)
	pv.AllEdits = fixUpRecord(pv.AllEdits)
}

func createProjectVersion(r Record)(*ProjectVersion){
	project_version := new(ProjectVersion)
	project_version.URL = r.URL
	project_version.ProjectName = r.ProjectName
	project_version.CRFVersionID = r.CRFVersionID
	project_version.LastVersion = r.LastVersion
	if r.SubjectCount.Valid{
		project_version.SubjectCount = int(r.SubjectCount.Int64)
	}
	return project_version
}