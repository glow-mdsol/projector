package main

import (
	"database/sql"
	"strings"

	"github.com/lib/pq"
)

// UnusedEdit represents the Structure for an unused Edit
type UnusedEdit struct {
	URL            string
	ProjectName    string `db:"project_name"`
	EditCheckName  string `db:"edit_check_name"`
	FormOID        string `db:"form_oids"`
	FieldOID       string `db:"field_oids"`
	VariableOID    string `db:"variable_oids"`
	UsageCount     int    `db:"edit_check_count"`
	OpenQuery      string `db:"open_query"`
	CustomFunction string `db:"custom_function"`
	NonConformant  string `db:"non_conformant"`
	RangeCheck     string `db:"range_checks"`
	RequiredCheck  string `db:"required_check"`
	FutureCheck    string `db:"future_checks"`
}

// SubjectCount represents the Structure for an Subject Count
type SubjectCount struct {
	URL          string
	ProjectName  string        `db:"project_name"`
	SubjectCount sql.NullInt64 `db:"subject_count"`
	RefreshDate  pq.NullTime   `db:"refresh_date"`
}

// Record represents the Structure for the returned stats from the DB
type Record struct {
	URL                                string
	URLID                              int           `db:"url_id"`
	ProjectName                        string        `db:"project_name"`
	CRFVersionID                       int           `db:"crf_version_id"`
	LastVersion                        bool          `db:"last_version"`
	SubjectCount                       sql.NullInt64 `db:"subject_count"`
	CheckStatus                        string        `db:"check_status"`
	RawTotalFieldEdits                 sql.NullInt64 `db:"total_edits_fld"`
	TotalFieldEdits                    int
	RawTotalProgEdits                  sql.NullInt64 `db:"total_edits_prg"`
	TotalProgEdits                     int
	RawTotalQueries                    sql.NullInt64 `db:"total_queries"`
	TotalQueries                       int
	RawTotalFieldQueries               sql.NullInt64 `db:"total_queries_fld"`
	TotalFieldQueries                  int
	RawTotalProgQueries                sql.NullInt64 `db:"total_queries_prg"`
	TotalProgQueries                   int
	RawTotalQueriesOpenQuery           sql.NullInt64 `db:"total_queries_open_query"`
	TotalQueriesOpenQuery              int
	TotalFieldEditsWithOpenQuery       int           `db:"total_edits_query_fld"`
	TotalProgEditsWithOpenQuery        int           `db:"total_edits_query_prg"`
	RawTotalFieldQueriesWithOpenQuery  sql.NullInt64 `db:"total_queries_query_fld"`
	RawTotalProgQueriesWithOpenQuery   sql.NullInt64 `db:"total_queries_query_prg"`
	TotalFieldQueriesWithOpenQuery     int
	TotalProgQueriesWithOpenQuery      int
	TotalFieldEditsFired               int           `db:"total_fired_fld"`
	TotalFieldEditsNotFired            int           `db:"total_not_fired_fld"`
	TotalProgEditsFired                int           `db:"total_fired_prg"`
	TotalProgEditsNotFired             int           `db:"total_not_fired_prg"`
	RawTotalFieldWithOpenQueryFired    sql.NullInt64 `db:"total_fired_query_fld"`
	TotalFieldWithOpenQueryFired       int
	RawTotalFieldWithOpenQueryNotFired sql.NullInt64 `db:"total_not_fired_query_fld"`
	TotalFieldWithOpenQueryNotFired    int
	RawTotalProgWithOpenQueryFired     sql.NullInt64 `db:"total_fired_query_prg"`
	TotalProgWithOpenQueryFired        int
	RawTotalProgWithOpenQueryNotFired  sql.NullInt64 `db:"total_not_fired_query_prg"`
	TotalProgWithOpenQueryNotFired     int
	TotalFieldEditsFiredWithChange     int `db:"fired_change_fld"`
	TotalFieldEditsFiredWithNoChange   int `db:"fired_no_change_fld"`
	TotalProgEditsFiredWithChange      int `db:"fired_change_prg"`
	TotalProgEditsFiredWithNoChange    int `db:"fired_no_change_prg"`
	FieldPercentageFired               float64
	FieldPercentageNotFired            float64
	FieldPercentageChanged             float64
	FieldPercentageNotChanged          float64
	ProgPercentageFired                float64
	ProgPercentageNotFired             float64
	ProgPercentageChanged              float64
	ProgPercentageNotChanged           float64
	TotalFieldEditsOpen                int `db:"open_edits_sys"`
	TotalProgEditsOpen                 int `db:"open_edits_prg"`
}

// SummaryCounts represents the Structure for the computed stats
type SummaryCounts struct {
	Threshold            int
	RecordCount          int
	SubjectCount         int
	TotalEdits           int
	TotalFldEdits        int
	TotalFldEditsFired   int
	TotalFldEditsUnfired int
	TotalFldEditsOpen    int
	TotalFldWithChange   int
	TotalFldWithNoChange int
	TotalPrgEdits        int
	TotalPrgEditsFired   int
	TotalPrgEditsUnfired int
	TotalPrgEditsOpen    int
	TotalPrgWithChange   int
	TotalPrgWithNoChange int
}

func (pv *Record) calculatePercentages() {
	// Gate the counts
	pv.FieldPercentageFired = 0.0
	pv.FieldPercentageNotFired = 0.0
	pv.FieldPercentageChanged = 0.0
	pv.FieldPercentageNotChanged = 0.0
	pv.ProgPercentageFired = 0.0
	pv.ProgPercentageNotFired = 0.0
	pv.ProgPercentageChanged = 0.0
	pv.ProgPercentageNotChanged = 0.0

	if pv.TotalFieldEdits > 0 {
		pv.FieldPercentageFired = float64(pv.TotalFieldEditsFired) / float64(pv.TotalFieldEdits)
		pv.FieldPercentageNotFired = float64(pv.TotalFieldEditsNotFired) / float64(pv.TotalFieldEdits)
	}
	if pv.TotalFieldEditsFired > 0 {
		pv.FieldPercentageChanged = float64(pv.TotalFieldEditsFiredWithChange) / float64(pv.TotalFieldEditsFired)
		pv.FieldPercentageNotChanged = float64(pv.TotalFieldEditsFiredWithNoChange) / float64(pv.TotalFieldEditsFired)
	}
	if pv.TotalProgEdits > 0 {
		pv.ProgPercentageFired = float64(pv.TotalProgEditsFired) / float64(pv.TotalProgEdits)
		pv.ProgPercentageNotFired = float64(pv.TotalProgEditsNotFired) / float64(pv.TotalProgEdits)
	}
	if pv.TotalProgEditsFired > 0 {
		pv.ProgPercentageChanged = float64(pv.TotalProgEditsFiredWithChange) / float64(pv.TotalProgEditsFired)
		pv.ProgPercentageNotChanged = float64(pv.TotalProgEditsFiredWithNoChange) / float64(pv.TotalProgEditsFired)
	}
}

func (rec *Record) fixUpRecord() {
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
	if rec.RawTotalQueriesOpenQuery.Valid {
		rec.TotalQueriesOpenQuery = int(rec.RawTotalQueriesOpenQuery.Int64)
	} else {
		rec.TotalQueriesOpenQuery = -1
	}
}

func createRaveURL(r Record) *RaveURL {
	raveURL := new(RaveURL)
	raveURL.URL = r.URL
	raveURL.URLID = r.URLID
	raveURL.URLPrefix = strings.Split(r.URL, ".")[0]
	return raveURL
}

func createProject(r Record) *Project {
	project := new(Project)
	project.URL = r.URL
	project.ProjectName = r.ProjectName
	project.URLID = r.URLID
	if r.SubjectCount.Valid {
		project.SubjectCount = int(r.SubjectCount.Int64)
	}
	return project
}

func createProjectVersion(r Record) *ProjectVersion {
	projectVersion := new(ProjectVersion)
	projectVersion.URL = r.URL
	projectVersion.ProjectName = r.ProjectName
	projectVersion.CRFVersionID = r.CRFVersionID
	projectVersion.LastVersion = r.LastVersion
	projectVersion.URLID = r.URLID
	return projectVersion
}
