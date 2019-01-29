package main

import (
	"database/sql"

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
	CRFVersionID                       string        `db:"crf_version_id"`
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
	TotalFldWithChange   int
	TotalFldWithNoChange int
	TotalPrgEdits        int
	TotalPrgEditsFired   int
	TotalPrgEditsUnfired int
	TotalPrgWithChange   int
	TotalPrgWithNoChange int
}

// ProjectVersion represents the structure for an individual Project Version
type ProjectVersion struct {
	URL               string
	URLID             int
	ProjectName       string
	CRFVersionID      string
	ActiveEditsOnly   *Record
	InactiveEditsOnly *Record
	LastVersion       bool
	SubjectCount      int
	ActiveEditCount   int
	InactiveEditCount int
}

// LastProjectVersion represents the stats for the last Project Version
type LastProjectVersion struct {
	ProjectName               string `db:"project_name"`
	ProjectID                 int    `db:"project_id"`
	CRFVersionID              int    `db:"crf_version_id"`
	SubjectCount              int    `db:"subject_count"`
	TotalCount                int    `db:"total_count"`
	FieldTotal                int    `db:"fld_total"`
	FieldTotalFired           int    `db:"fld_total_fired"`
	FieldTotalNotFired        int    `db:"fld_total_not_fired"`
	FieldNotChanged           int    `db:"fld_no_change_count"`
	FieldChanged              int    `db:"fld_change_count"`
	ProgTotal                 int    `db:"prg_total"`
	ProgTotalFired            int    `db:"prg_total_fired"`
	ProgTotalNotFired         int    `db:"prg_total_not_fired"`
	ProgNotChanged            int    `db:"prg_no_change_count"`
	ProgChanged               int    `db:"prg_change_count"`
	FieldPercentageFired      float64
	FieldPercentageNotFired   float64
	FieldPercentageChanged    float64
	FieldPercentageNotChanged float64
	ProgPercentageFired       float64
	ProgPercentageNotFired    float64
	ProgPercentageChanged     float64
	ProgPercentageNotChanged  float64
}

func (pv *LastProjectVersion) calculatePercentages() {
	// Gate the counts
	pv.FieldPercentageFired = 0.0
	pv.FieldPercentageNotFired = 0.0
	pv.FieldPercentageChanged = 0.0
	pv.FieldPercentageNotChanged = 0.0
	pv.ProgPercentageFired = 0.0
	pv.ProgPercentageNotFired = 0.0
	pv.ProgPercentageChanged = 0.0
	pv.ProgPercentageNotChanged = 0.0

	if pv.FieldTotal > 0 {
		pv.FieldPercentageFired = float64(pv.FieldTotalFired) / float64(pv.FieldTotal)
		pv.FieldPercentageNotFired = float64(pv.FieldTotalNotFired) / float64(pv.FieldTotal)
	}
	if pv.FieldTotalFired > 0 {
		pv.FieldPercentageChanged = float64(pv.FieldChanged) / float64(pv.FieldTotalFired)
		pv.FieldPercentageNotChanged = float64(pv.FieldNotChanged) / float64(pv.FieldTotalFired)
	}
	if pv.ProgTotal > 0 {
		pv.ProgPercentageFired = float64(pv.ProgTotalFired) / float64(pv.ProgTotal)
		pv.ProgPercentageNotFired = float64(pv.ProgTotalNotFired) / float64(pv.ProgTotal)
	}
	if pv.ProgTotalFired > 0 {
		pv.ProgPercentageChanged = float64(pv.ProgChanged) / float64(pv.ProgTotalFired)
		pv.ProgPercentageNotChanged = float64(pv.ProgNotChanged) / float64(pv.ProgTotalFired)
	}
}

// func calculateInactiveCounts(pv *ProjectVersion) *ProjectVersion {
// 	rec := new(Record)
// 	rec.URL = pv.URL
// 	rec.ProjectName = pv.ProjectName
// 	rec.CRFVersionID = pv.CRFVersionID
// 	rec.CheckStatus = "INACTIVE"
// 	rec.TotalFieldEdits = pv.AllEdits.TotalFieldEdits - pv.ActiveEditsOnly.TotalFieldEdits
// 	rec.TotalProgEdits = pv.AllEdits.TotalProgEdits - pv.ActiveEditsOnly.TotalProgEdits
// 	rec.TotalFieldQueries = pv.AllEdits.TotalFieldQueries - pv.ActiveEditsOnly.TotalFieldQueries
// 	rec.TotalProgQueries = pv.AllEdits.TotalProgQueries - pv.ActiveEditsOnly.TotalProgQueries
// 	rec.TotalFieldEditsWithOpenQuery = pv.AllEdits.TotalFieldEditsWithOpenQuery - pv.ActiveEditsOnly.TotalFieldEditsWithOpenQuery
// 	rec.TotalProgEditsWithOpenQuery = pv.AllEdits.TotalProgEditsWithOpenQuery - pv.ActiveEditsOnly.TotalProgEditsWithOpenQuery
// 	rec.TotalFieldQueriesWithOpenQuery = pv.AllEdits.TotalFieldQueriesWithOpenQuery - pv.ActiveEditsOnly.TotalFieldQueriesWithOpenQuery
// 	rec.TotalProgQueriesWithOpenQuery = pv.AllEdits.TotalProgQueriesWithOpenQuery - pv.ActiveEditsOnly.TotalProgQueriesWithOpenQuery
// 	rec.TotalFieldEditsFired = pv.AllEdits.TotalFieldEditsFired - pv.ActiveEditsOnly.TotalFieldEditsFired
// 	rec.TotalProgEditsFired = pv.AllEdits.TotalProgEditsFired - pv.ActiveEditsOnly.TotalProgEditsFired
// 	rec.TotalFieldEditsFired = pv.AllEdits.TotalFieldEditsFired - pv.ActiveEditsOnly.TotalFieldEditsFired
// 	rec.TotalFieldWithOpenQueryFired = pv.AllEdits.TotalFieldWithOpenQueryFired - pv.ActiveEditsOnly.TotalFieldWithOpenQueryFired
// 	rec.TotalProgWithOpenQueryFired = pv.AllEdits.TotalProgWithOpenQueryFired - pv.ActiveEditsOnly.TotalProgWithOpenQueryFired
// 	rec.TotalFieldEditsFiredWithNoChange = pv.AllEdits.TotalFieldEditsFiredWithNoChange - pv.ActiveEditsOnly.TotalFieldEditsFiredWithNoChange
// 	rec.TotalProgEditsFiredWithNoChange = pv.AllEdits.TotalProgEditsFiredWithNoChange - pv.ActiveEditsOnly.TotalProgEditsFiredWithNoChange
// 	pv.InactiveEditsOnly = *rec
// 	return pv
// }

// func (pv *ProjectVersion) calculateInactiveCounts() {
// 	rec := new(Record)
// 	rec.URL = pv.URL
// 	rec.ProjectName = pv.ProjectName
// 	rec.CRFVersionID = pv.CRFVersionID
// 	rec.CheckStatus = "INACTIVE"
// 	rec.TotalFieldEdits = pv.AllEdits.TotalFieldEdits - pv.ActiveEditsOnly.TotalFieldEdits
// 	rec.TotalProgEdits = pv.AllEdits.TotalProgEdits - pv.ActiveEditsOnly.TotalProgEdits
// 	rec.TotalFieldQueries = pv.AllEdits.TotalFieldQueries - pv.ActiveEditsOnly.TotalFieldQueries
// 	rec.TotalProgQueries = pv.AllEdits.TotalProgQueries - pv.ActiveEditsOnly.TotalProgQueries
// 	rec.TotalFieldEditsWithOpenQuery = pv.AllEdits.TotalFieldEditsWithOpenQuery - pv.ActiveEditsOnly.TotalFieldEditsWithOpenQuery
// 	rec.TotalProgEditsWithOpenQuery = pv.AllEdits.TotalProgEditsWithOpenQuery - pv.ActiveEditsOnly.TotalProgEditsWithOpenQuery
// 	rec.TotalFieldQueriesWithOpenQuery = pv.AllEdits.TotalFieldQueriesWithOpenQuery - pv.ActiveEditsOnly.TotalFieldQueriesWithOpenQuery
// 	rec.TotalProgQueriesWithOpenQuery = pv.AllEdits.TotalProgQueriesWithOpenQuery - pv.ActiveEditsOnly.TotalProgQueriesWithOpenQuery
// 	rec.TotalFieldEditsFired = pv.AllEdits.TotalFieldEditsFired - pv.ActiveEditsOnly.TotalFieldEditsFired
// 	rec.TotalProgEditsFired = pv.AllEdits.TotalProgEditsFired - pv.ActiveEditsOnly.TotalProgEditsFired
// 	rec.TotalFieldEditsFired = pv.AllEdits.TotalFieldEditsFired - pv.ActiveEditsOnly.TotalFieldEditsFired
// 	rec.TotalFieldWithOpenQueryFired = pv.AllEdits.TotalFieldWithOpenQueryFired - pv.ActiveEditsOnly.TotalFieldWithOpenQueryFired
// 	rec.TotalProgWithOpenQueryFired = pv.AllEdits.TotalProgWithOpenQueryFired - pv.ActiveEditsOnly.TotalProgWithOpenQueryFired
// 	rec.TotalFieldEditsFiredWithNoChange = pv.AllEdits.TotalFieldEditsFiredWithNoChange - pv.ActiveEditsOnly.TotalFieldEditsFiredWithNoChange
// 	rec.TotalProgEditsFiredWithNoChange = pv.AllEdits.TotalProgEditsFiredWithNoChange - pv.ActiveEditsOnly.TotalProgEditsFiredWithNoChange
// 	pv.InactiveEditsOnly = *rec
// }

func fixUpRecord(rec *Record) *Record {
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

// func fixUpNullValues(pv *ProjectVersion) *ProjectVersion {
// 	var activeCount = 0
// 	var inactiveCount = 0
// 	pv.ActiveEditsOnly = fixUpRecord(pv.ActiveEditsOnly)
// 	if pv.InactiveEditsOnly != nil {
// 		pv.InactiveEditsOnly = fixUpRecord(pv.InactiveEditsOnly)
// 		if pv.InactiveEditsOnly.TotalProgEdits >= 0 {
// 			// skip the null values
// 			inactiveCount += pv.InactiveEditsOnly.TotalProgEdits
// 		}
// 		if pv.InactiveEditsOnly.TotalFieldEdits >= 0 {
// 			// skip the null values
// 			inactiveCount += pv.InactiveEditsOnly.TotalFieldEdits
// 		}
// 	}
// 	if pv.ActiveEditsOnly.TotalFieldEdits >= 0 {
// 		activeCount += pv.ActiveEditsOnly.TotalFieldEdits
// 	}
// 	if pv.ActiveEditsOnly.TotalProgEdits >= 0 {
// 		activeCount += pv.ActiveEditsOnly.TotalProgEdits
// 	}
// 	pv.ActiveEditCount = activeCount
// 	pv.InactiveEditCount = inactiveCount
// 	return pv
// }

func (pv *ProjectVersion) fixUpNullValues() {
	var activeCount = 0
	var inactiveCount = 0
	pv.ActiveEditsOnly = fixUpRecord(pv.ActiveEditsOnly)
	if pv.InactiveEditsOnly != nil {
		pv.InactiveEditsOnly = fixUpRecord(pv.InactiveEditsOnly)
		if pv.InactiveEditsOnly.TotalProgEdits >= 0 {
			// skip the null values
			inactiveCount += pv.InactiveEditsOnly.TotalProgEdits
		}
		if pv.InactiveEditsOnly.TotalFieldEdits >= 0 {
			// skip the null values
			inactiveCount += pv.InactiveEditsOnly.TotalFieldEdits
		}
	}
	if pv.ActiveEditsOnly.TotalFieldEdits >= 0 {
		activeCount += pv.ActiveEditsOnly.TotalFieldEdits
	}
	if pv.ActiveEditsOnly.TotalProgEdits >= 0 {
		activeCount += pv.ActiveEditsOnly.TotalProgEdits
	}
	pv.ActiveEditCount = activeCount
	pv.InactiveEditCount = inactiveCount
}

func createProjectVersion(r Record) *ProjectVersion {
	projectVersion := new(ProjectVersion)
	projectVersion.URL = r.URL
	projectVersion.ProjectName = r.ProjectName
	projectVersion.CRFVersionID = r.CRFVersionID
	projectVersion.LastVersion = r.LastVersion
	projectVersion.URLID = r.URLID
	if r.SubjectCount.Valid {
		projectVersion.SubjectCount = int(r.SubjectCount.Int64)
	}
	return projectVersion
}
