package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
	"strings"
	"database/sql"
	"sort"
)

type Record struct {
	URL                               string
	ProjectName                       string    `db:"project_name"`
	CRFVersionID                      string `db:"crf_version_id"`
	CheckStatus                       string `db:"check_status"`
	RawTotalFieldEdits                sql.NullInt64    `db:"total_edits_fld"`
	TotalFieldEdits                   int
	RawTotalProgEdits                 sql.NullInt64    `db:"total_edits_prg"`
	TotalProgEdits                    int
	RawTotalFieldQueries              sql.NullInt64    `db:"total_queries_fld"`
	TotalFieldQueries                 int
	RawTotalProgQueries               sql.NullInt64    `db:"total_queries_prg"`
	TotalProgQueries                  int
	TotalFieldEditsWithOpenQuery      int    `db:"total_edits_query_fld"`
	TotalProgEditsWithOpenQuery       int    `db:"total_edits_query_prg"`
	RawTotalFieldQueriesWithOpenQuery sql.NullInt64    `db:"total_queries_query_fld"`
	RawTotalProgQueriesWithOpenQuery  sql.NullInt64    `db:"total_queries_query_prg"`
	TotalFieldQueriesWithOpenQuery    int
	TotalProgQueriesWithOpenQuery     int
	TotalFieldEditsFired              int    `db:"total_fired_fld"`
	TotalProgEditsFired               int    `db:"total_fired_prg"`
	RawTotalFieldWithOpenQueryFired sql.NullInt64    `db:"total_fired_query_fld"`
	RawTotalProgWithOpenQueryFired  sql.NullInt64    `db:"total_fired_query_prg"`
	TotalFieldWithOpenQueryFired int
	TotalProgWithOpenQueryFired  int
	TotalFieldEditsFiredWithNoChange  int    `db:"fired_no_change_fld"`
	TotalProgEditsFiredWithNoChange   int    `db:"fired_no_change_prg"`
}

type ProjectVersion struct {
	URL               string
	ProjectName       string
	CRFVersionID      string
	AllEdits          Record
	ActiveEditsOnly   Record
	InactiveEditsOnly Record
}

func calculateInactiveCounts(pv *ProjectVersion) (*ProjectVersion) {
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

func fixUpRecord(rec Record) (Record) {
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

func fixUpNullValues(pv *ProjectVersion) (*ProjectVersion) {
	pv.ActiveEditsOnly = fixUpRecord(pv.ActiveEditsOnly)
	pv.AllEdits = fixUpRecord(pv.AllEdits)
	return pv
}

// Write the Header row
func writeHeaderRow(data []string, sheet *xlsx.Sheet) {
	boldface := *xlsx.NewFont(12, "Verdana")
	boldface.Bold = true
	centerHalign := *xlsx.DefaultAlignment()
	centerHalign.Horizontal = "center"
	align := xlsx.Alignment{
		ShrinkToFit: true,
	}
	titleFace := xlsx.NewStyle()
	titleFace.Font = boldface
	titleFace.Alignment = centerHalign
	titleFace.ApplyAlignment = true
	titleFace.ApplyFont = true
	titleFace.Alignment = align
	row := sheet.AddRow()
	for idx, datum := range data {
		cell := row.AddCell()
		cell.SetStyle(titleFace)
		cell.SetString(datum)
		sheet.SetColWidth(idx, idx, float64(len(datum)) )
	}
}

// export the subject counts
func exportSubjectCounts(wbk *xlsx.File, db *sqlx.DB, pattern string) {
	tab_name := "Subject Counts"
	headers := []string{"Rave URL", "Project Name", "Subject Count"}
	// TODO: Add the Date of the Subject Count
	q := `WITH counts AS (SELECT
                  project_id         AS project_id,
                  MAX(subject_count) AS subject_count
                FROM edit_check
                  JOIN project
                    ON edit_check.project_id = project.id
                  JOIN rave_url
                    ON edit_check.url_id = rave_url.id
                WHERE rave_url.url LIKE '%' || $1 || '%' OR rave_url.alternate_url LIKE '%' || $1 || '%'
                GROUP BY project_id)
SELECT
  CASE WHEN rave_url.url LIKE '%hdcvc%'
    THEN rave_url.alternate_url
  ELSE rave_url.url END AS URL,
  project.project_name,
  counts.subject_count
FROM counts
  JOIN project
    ON project.id = counts.project_id
  JOIN rave_url
    ON project.url_id = rave_url.id
ORDER BY  rave_url.url, project.project_name
`
	rows, err := db.Queryx(q, pattern)
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()

	// create the sheet
	sheet, created := getOrAddSheet(wbk, tab_name)
	if created {
		// Add the headers
		writeHeaderRow(headers, sheet)
	}

	// iterate over rows
	for rows.Next() {
		var study_url string
		var project_name string
		var subject_count int
		err := rows.Scan(&study_url, &project_name, &subject_count)
		if err != nil {
			log.Fatal("Error processing data for subject counts ", err)
		}
		var cell *xlsx.Cell
		// Rows
		row := sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(study_url)
		cell = row.AddCell()
		cell.SetString(project_name)
		cell = row.AddCell()
		cell.SetInt(subject_count)
	}
}

func writeMetricsRow(rec Record, row *xlsx.Row) {
	var cell *xlsx.Cell
	// Study URL
	cell = row.AddCell()
	cell.Value = rec.URL
	// Project Name
	cell = row.AddCell()
	cell.Value = rec.ProjectName
	// CRF Version
	cell = row.AddCell()
	cell.Value = rec.CRFVersionID
	// Status
	cell = row.AddCell()
	cell.Value = rec.CheckStatus
	// Total Edits (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEdits)
	// Total Edits Fired (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEditsFired)
	// %ge Edits Fired (fld)
	cell = row.AddCell()
	if rec.TotalFieldEdits == 0 {
		cell.SetFloat(0.0)
	} else {
		cell.SetFloatWithFormat(float64(rec.TotalFieldEditsFired)/float64(rec.TotalFieldEdits)*100.0, "#,##0.00;(#,##0.00)")
	}
	// Total Edits Unfired (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEdits - rec.TotalFieldEditsFired)
	// %ge Edits Unfired (fld)
	cell = row.AddCell()
	if rec.TotalFieldEdits == 0 {
		cell.SetFloat(0.0)
	} else {
		cell.SetFloatWithFormat(float64(rec.TotalFieldEdits-rec.TotalFieldEditsFired)/float64(rec.TotalFieldEdits)*100.0, "#,##0.00;(#,##0.00)")
	}
	// Total Edits (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEdits)
	// Total Edits With OpenQuery (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEditsWithOpenQuery)
	// Total Edits Fired (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgWithOpenQueryFired)
	// %ge Edits Fired (prg)
	cell = row.AddCell()
	if rec.TotalProgEditsWithOpenQuery == 0 {
		cell.SetFloat(0.0)
	} else {
		cell.SetFloatWithFormat(float64(rec.TotalProgWithOpenQueryFired)/float64(rec.TotalProgEditsWithOpenQuery)*100.0, "#,##0.00;(#,##0.00)")
	}
	// Total Edits Unfired (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEditsWithOpenQuery - rec.TotalProgWithOpenQueryFired)
	// %ge Edits Unfired (prg)
	cell = row.AddCell()
	if rec.TotalProgEditsWithOpenQuery == 0 {
		cell.SetFloat(0.0)
	} else {
		cell.SetFloatWithFormat(float64(rec.TotalProgEditsWithOpenQuery-rec.TotalProgWithOpenQueryFired)/float64(rec.TotalProgEditsWithOpenQuery)*100.0, "#,##0.00;(#,##0.00)")
	}
	// Total Queries (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldQueries)
	// Total Queries (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgQueries)
	// Total Queries With OpenQuery (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgQueriesWithOpenQuery)
	//cell = row.AddCell()
	//cell.SetInt(rec.TotalFieldEditsFiredWithNoChange)
	//cell = row.AddCell()
	//cell.SetInt(rec.TotalProgEditsFiredWithNoChange)
}

// write study metrics
func writeStudyMetrics(data map[string][]ProjectVersion, wbk *xlsx.File) {
	headers := []string{"Study URL",
						"Project Name",
						"CRF Version",
						"Status",
						"Total Edits (fld)",
						"Total Edits Fired (fld)",
						"%ge Edits Fired (fld)",
						"Total Edits Unfired (fld)",
						"%ge Edits Unfired (fld)",
						"Total Edits (prg)",
						"Total Edits With OpenQuery (prg)",
						"Total Edits Fired (prg)",
						"%ge Edits Fired (prg)",
						"Total Edits Unfired (prg)",
						"%ge Edits Unfired (prg)",
						"Total Queries (fld)",
						"Total Queries (prg)",
						"Total Queries With OpenQuery (prg)",
		//"Total Edits Fired With No Change (fld)",
		//"Total Edits Fired With No Change (prg)"
	}
	var urls []string
	for k := range data {
		urls = append(urls, k)
	}
	sort.Strings(urls)
	for _, url := range urls {

		// create the sheet
		sheet, created := getOrAddSheet(wbk, url)
		if created {
			// Add the headers
			writeHeaderRow(headers, sheet)
		}
		//log.Println("Created Sheet for URL ", url)

		for _, project_version := range data[url] {
			// Add the row for Active Checks
			var row *xlsx.Row
			active_only := project_version.ActiveEditsOnly
			//
			if active_only.URL != "" {
				row = sheet.AddRow()
				writeMetricsRow(active_only, row)
				row = sheet.AddRow()
				inactive_only := project_version.InactiveEditsOnly
				writeMetricsRow(inactive_only, row)
			} else {
				// predate the active check
				row = sheet.AddRow()
				writeMetricsRow(project_version.AllEdits, row)
			}
		}
	}
}

// export the study metrics
func exportStudyMetrics(wbk *xlsx.File, db *sqlx.DB, pattern string) {
	q := `SELECT
	CASE WHEN rave_url.url LIKE '%hdcvc%' THEN rave_url.alternate_url ELSE rave_url.url END as URL,
  project.project_name,
  crf_version_id,
  'ACTIVEONLY' AS check_status,
  total_edits_fld,
  total_edits_prg,
  total_queries_fld,
  total_queries_prg,
  total_edits_query_fld,
  total_edits_query_prg,
  total_queries_query_fld,
  total_queries_query_prg,
  total_fired_fld,
  total_fired_prg,
  total_fired_query_fld,
  total_fired_query_prg,
  fired_no_change_fld,
  fired_no_change_prg
FROM project_version_active_view
  JOIN project ON project_version_active_view.project_id = project.id
  JOIN rave_url ON project.url_id = rave_url.id
WHERE rave_url.url LIKE '%' || $1 || '%' OR rave_url.alternate_url like '%' || $1 || '%'
  UNION
SELECT
	CASE WHEN rave_url.url LIKE '%hdcvc%' THEN rave_url.alternate_url ELSE rave_url.url END as URL,
  project.project_name,
  crf_version_id,
  'ALLCHECKS' AS check_status,
  total_edits_fld,
  total_edits_prg,
  total_queries_fld,
  total_queries_prg,
  total_edits_query_fld,
  total_edits_query_prg,
  total_queries_query_fld,
  total_queries_query_prg,
  total_fired_fld,
  total_fired_prg,
  total_fired_query_fld,
  total_fired_query_prg,
  fired_no_change_fld,
  fired_no_change_prg
FROM project_version_view
  JOIN project ON project_version_view.project_id = project.id
  JOIN rave_url ON project.url_id = rave_url.id
WHERE rave_url.url LIKE '%' || $1 || '%' OR rave_url.alternate_url like '%' || $1 || '%'
ORDER BY url, project_name, crf_version_id, check_status;
`
	rows, err := db.Queryx(q, pattern)
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()

	// Status variable
	urls := make(map[string][]ProjectVersion)
	var project_version *ProjectVersion
	for rows.Next() {
		var r Record
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		// First Rune
		if project_version == nil {
			project_version = new(ProjectVersion)
			project_version.URL = r.URL
			project_version.ProjectName = r.ProjectName
			project_version.CRFVersionID = r.CRFVersionID
		} else if r.URL != project_version.URL || r.ProjectName != project_version.ProjectName || r.CRFVersionID != project_version.CRFVersionID {
			// remove the .mdsol.com
			// refresh the the inactive counts
			prefix := strings.Split(project_version.URL, ".")[0]
			project_version = fixUpNullValues(project_version)
			project_version = calculateInactiveCounts(project_version)
			urls[prefix] = append(urls[prefix], *project_version)
			project_version = new(ProjectVersion)
			project_version.URL = r.URL
			project_version.ProjectName = r.ProjectName
			project_version.CRFVersionID = r.CRFVersionID
		}
		if r.CheckStatus == "ALLCHECKS" {
			project_version.AllEdits = r
		} else {
			project_version.ActiveEditsOnly = r
		}
	}
	log.Println("Generated metrics for ", len(urls), "URLs")
	writeStudyMetrics(urls, wbk)
}

func getOrAddSheet(wbk *xlsx.File, name string) (*xlsx.Sheet, bool) {
	for _, sheet := range wbk.Sheets {
		if sheet.Name == name {
			return sheet, false
		}
	}
	sheet, err := wbk.AddSheet(name)
	if err != nil {
		log.Fatal("Unable to create sheet", name, ": ", err)
	}
	return sheet, true
}

func exportUselessEdits(wbk *xlsx.File, db *sqlx.DB, pattern string) {
	headers := []string{"Study URL",
						"Project Name",
						"Edit Check Name",
						"Times Used",
						"OpenQuery Check?"}

	tab_name := "Unused Edits"

	q := `WITH total AS (
    SELECT
		rave_url.id as url_id,
      	project_id as project_id,
      	edit_check_name,
	    COUNT(*) AS edit_check_count,
      	SUM(CASE WHEN Actions LIKE '%OpenQuery%' THEN 1 ELSE 0 END) AS open_query_count,
      	SUM(query_count) AS total_count
          FROM edit_check
      JOIN rave_url ON edit_check.url_id = rave_url.id
    WHERE rave_url.url LIKE '%' || $1 || '%' OR rave_url.alternate_url LIKE '%' || $1 || '%'
     AND edit_check.is_active = 1
  GROUP BY rave_url.id, project_id, edit_check_name
)
SELECT CASE WHEN rave_url.url LIKE '%hdcvc%' THEN rave_url.alternate_url ELSE rave_url.url END as URL,
  project.project_name,
  edit_check_name,
  edit_check_count,
  CASE WHEN open_query_count > 0 THEN 'Yes' ELSE 'No' END
  FROM total
  JOIN rave_url ON total.url_id = rave_url.id
  JOIN project ON total.project_id = project.id
WHERE total_count = 0;
`
	rows, err := db.Queryx(q, pattern)
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()

	// create the sheet
	sheet, created := getOrAddSheet(wbk, tab_name)
	if created {
		// Add the headers
		writeHeaderRow(headers, sheet)
	}
	// Export the results
	for rows.Next() {
		var study_url string
		var project_name string
		var edit_check_name string
		var usage_count int
		var open_query string
		err := rows.Scan(&study_url, &project_name, &edit_check_name, &usage_count, &open_query)
		if err != nil {
			log.Fatal("Error processing data for unused edits ", err)
		}
		var cell *xlsx.Cell
		// Rows
		row := sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(study_url)
		cell = row.AddCell()
		cell.SetString(project_name)
		cell = row.AddCell()
		cell.SetString(edit_check_name)
		cell = row.AddCell()
		cell.SetInt(usage_count)
		cell = row.AddCell()
		cell.SetString(open_query)
	}
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var patternsArray arrayFlags
	flag.Var(&patternsArray, "pattern", "Supply the URL patterns")
	host_name := flag.String("dbhost", "localhost", "Database Host")
	db_name := flag.String("dbname", "editsfive", "Database Name")
	db_user := flag.String("user", "edits", "Database User")
	db_pass := flag.String("password", "apple01", "Database Password")
	file_name := flag.String("output", "report", "Output File Name")
	flag.Parse()
	if len(patternsArray) == 0 {
		log.Fatal("Need to specify the patterns")
	}
	var data_source_name = fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		*host_name,
		*db_user,
		*db_name,
		*db_pass)
	var dbConn *sqlx.DB
	// make the database connection
	dbConn, err := sqlx.Open("postgres", string(data_source_name))
	if err != nil {
		log.Fatal(err)
	}
	workbook := xlsx.NewFile()
	for _, url_pattern := range patternsArray {
		log.Println("Processing URL Pattern ", url_pattern)
		log.Printf("Inserting Subject Counts")
		exportSubjectCounts(workbook, dbConn, url_pattern)
		log.Print("Inserting Unfired Edits")
		exportUselessEdits(workbook, dbConn, url_pattern)
		log.Print("Inserting URL Metrics")
		exportStudyMetrics(workbook, dbConn, url_pattern)
	}
	// make up the prefix using the range of patterns
	prefix := strings.Join(patternsArray, "_")
	filename := fmt.Sprintf("%s_%s_%s.xlsx", prefix, *file_name, time.Now().Format("2006-01-02"))
	workbook.Save(filename)
}
