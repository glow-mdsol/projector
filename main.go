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
	URL                              string
	ProjectName                      string	`db:"project_name"`
	CRFVersionID                     string `db:"crf_version_id"`
	CheckStatus                      string `db:"check_status"`
	TotalFieldEdits                  int	`db:"total_edits_fld"`
	TotalProgEdits                   int	`db:"total_edits_prg"`
	TotalFieldQueries                int	`db:"total_queries_fld"`
	TotalProgQueries                 int	`db:"total_queries_prg"`
	TotalFieldEditsWithOpenQuery     int	`db:"total_edits_query_fld"`
	TotalProgEditsWithOpenQuery      int	`db:"total_edits_query_prg"`
	RawTotalFieldQueriesWithOpenQuery   sql.NullInt64	`db:"total_queries_query_fld"`
	RawTotalProgQueriesWithOpenQuery    sql.NullInt64	`db:"total_queries_query_prg"`
	TotalFieldQueriesWithOpenQuery   int
	TotalProgQueriesWithOpenQuery    int
	TotalFieldEditsFired             int	`db:"total_fired_fld"`
	TotalProgEditsFired              int	`db:"total_fired_prg"`
	TotalFieldEditsFiredWithNoChange int	`db:"fired_no_change_fld"`
	TotalProgEditsFiredWithNoChange  int	`db:"fired_no_change_prg"`
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
	rec.TotalFieldEditsFiredWithNoChange = pv.AllEdits.TotalFieldEditsFiredWithNoChange - pv.ActiveEditsOnly.TotalFieldEditsFiredWithNoChange
	rec.TotalProgEditsFiredWithNoChange = pv.AllEdits.TotalProgEditsFiredWithNoChange - pv.ActiveEditsOnly.TotalProgEditsFiredWithNoChange
	pv.InactiveEditsOnly = *rec
	return pv
}

func fixUpNullValues(pv *ProjectVersion) (*ProjectVersion)  {
	if pv.AllEdits.RawTotalFieldQueriesWithOpenQuery.Valid {
		pv.AllEdits.TotalFieldQueriesWithOpenQuery = int(pv.AllEdits.RawTotalFieldQueriesWithOpenQuery.Int64)
	} else {
		pv.AllEdits.TotalFieldQueriesWithOpenQuery = -1
	}
	if pv.AllEdits.RawTotalProgQueriesWithOpenQuery.Valid {
		pv.AllEdits.TotalProgQueriesWithOpenQuery = int(pv.AllEdits.RawTotalProgQueriesWithOpenQuery.Int64)
	} else {
		pv.AllEdits.TotalProgQueriesWithOpenQuery = -1
	}
	if pv.ActiveEditsOnly.RawTotalFieldQueriesWithOpenQuery.Valid {
		pv.ActiveEditsOnly.TotalFieldQueriesWithOpenQuery = int(pv.ActiveEditsOnly.RawTotalFieldQueriesWithOpenQuery.Int64)
	} else {
		pv.ActiveEditsOnly.TotalFieldQueriesWithOpenQuery = -1
	}
	if pv.ActiveEditsOnly.RawTotalProgQueriesWithOpenQuery.Valid {
		pv.ActiveEditsOnly.TotalProgQueriesWithOpenQuery = int(pv.ActiveEditsOnly.RawTotalProgQueriesWithOpenQuery.Int64)
	} else {
		pv.ActiveEditsOnly.TotalProgQueriesWithOpenQuery = -1
	}
	return pv
}

// Write the Header row
func writeHeaderRow(data []string, sheet *xlsx.Sheet) {
	boldface := *xlsx.NewFont(12, "Verdana")
	boldface.Bold = true
	centerHalign := *xlsx.DefaultAlignment()
	centerHalign.Horizontal = "center"
	titleFace := xlsx.NewStyle()
	titleFace.Font = boldface
	titleFace.Alignment = centerHalign
	titleFace.ApplyAlignment = true
	titleFace.ApplyFont = true
	row := sheet.AddRow()
	for _, datum := range data {
		cell := row.AddCell()
		cell.SetStyle(titleFace)
		cell.SetString(datum)
	}
}

// export the subject counts
func exportSubjectCounts(wbk *xlsx.File, db *sqlx.DB, pattern *string) {
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

	// get a sheet
	sheet, err := wbk.AddSheet(tab_name)
	if err != nil {
		log.Fatal("Unable to create sheet: ", tab_name)
	}

	// Add the headers
	writeHeaderRow(headers, sheet)

	// iterate over rows
	for rows.Next() {
		var study_url string
		var project_name string
		var subject_count int
		err := rows.Scan(&study_url, &project_name, &subject_count)
		if err != nil{
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
	cell = row.AddCell()
	cell.Value = rec.URL
	cell = row.AddCell()
	cell.Value = rec.ProjectName
	cell = row.AddCell()
	cell.Value = rec.CRFVersionID
	cell = row.AddCell()
	cell.Value = rec.CheckStatus
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEdits)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEdits)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldQueries)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgQueries)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEditsWithOpenQuery)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgQueriesWithOpenQuery)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEditsFired)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEditsFired)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEditsFiredWithNoChange)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEditsFiredWithNoChange)
}

// write study metrics
func writeStudyMetrics(data map[string][]ProjectVersion, wbk *xlsx.File) {
	headers := []string{"Study URL",
		"Project Name",
		"CRF Version",
		"Status",
		"Total Edits (fld)",
		"Total Edits (prg)",
		"Total Queries (fld)",
		"Total Queries (prg)",
		"Total Edits With OpenQuery (prg)",
		"Total Queries With OpenQuery (prg)",
		"Total Edits Fired (fld)",
		"Total Edits Fired (prg)",
		"Total Edits Fired With No Change (fld)",
		"Total Edits Fired With No Change (prg)"}
	var urls []string
	for k := range data {
		urls = append(urls, k)
	}
	sort.Strings(urls)
	for _, url := range urls {

		// Create a new sheet
		sheet, err := wbk.AddSheet(url)
		if err != nil {
			log.Fatal("Error creating new sheet: ", err)
		}
		writeHeaderRow(headers, sheet)
		//log.Println("Created Sheet for URL ", url)

		for _, project_version := range data[url] {
			// Add the row for Active Checks
			var row *xlsx.Row
			row = sheet.AddRow()
			active_only := project_version.ActiveEditsOnly
			writeMetricsRow(active_only, row)
			row = sheet.AddRow()
			inactive_only := project_version.InactiveEditsOnly
			writeMetricsRow(inactive_only, row)
		}
	}
}

// export the study metrics
func exportStudyMetrics(wbk *xlsx.File, db *sqlx.DB, pattern *string) {
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
		} else if r.URL != project_version.URL || r.ProjectName != project_version.ProjectName ||  r.CRFVersionID != project_version.CRFVersionID {
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

func exportUselessEdits(wbk *xlsx.File, db *sqlx.DB, pattern *string) {
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
	sheet, err := wbk.AddSheet(tab_name)
	if err != nil {
		log.Fatal("Unable to create sheet: ", tab_name)
	}

	// Add the headers
	writeHeaderRow(headers, sheet)

	// Export the results
	for rows.Next() {
		var study_url string
		var project_name string
		var edit_check_name string
		var usage_count int
		var open_query string
		err := rows.Scan(&study_url, &project_name, &edit_check_name, &usage_count, &open_query)
		if err != nil{
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

func main() {
	url_pattern := flag.String("pattern", "", "Supply the URL pattern")
	host_name := flag.String("dbhost", "localhost", "Database Host")
	db_name := flag.String("dbname", "editsfive", "Database Name")
	db_user := flag.String("user", "edits", "Database User")
	db_pass := flag.String("password", "apple01", "Database Password")
	file_name := flag.String("output", "report", "Output File Name")
	flag.Parse()
	if *url_pattern == "" {
		log.Fatal("Need to specify the pattern")
	}
	var data_source_name = fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		*host_name,
		*db_user,
		*db_name,
		*db_pass)
	var db *sqlx.DB
	// make the database connection
	db, err := sqlx.Open("postgres", string(data_source_name))
	if err != nil {
		log.Fatal(err)
	}
	workbook := xlsx.NewFile()
	log.Printf("Inserting Subject Counts")
	exportSubjectCounts(workbook, db, url_pattern)
	log.Print("Inserting Unfired Edits")
	exportUselessEdits(workbook, db, url_pattern)
	log.Print("Inserting URL Metrics")
	exportStudyMetrics(workbook, db, url_pattern)
	filename := fmt.Sprintf("%s_%s_%s.xlsx", *url_pattern, *file_name, time.Now().Format("2006-01-02"))
	workbook.Save(filename)
}
