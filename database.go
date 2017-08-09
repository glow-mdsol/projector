package main

import (
	"log"
	"github.com/jmoiron/sqlx"
	"strings"
)

func doesPatternMatch(pattern string, db *sqlx.DB) bool {
	q := `SELECT COUNT(*) FROM rave_url
	WHERE rave_url.url LIKE '%' || $1 || '%' OR rave_url.alternate_url LIKE '%' || $1 || '%'`
	var count int
	err := db.Get(&count, q, pattern)
	if err == nil {
		return count != 0
	}
	return false
}

// export the subject counts
func getSubjectCounts(db *sqlx.DB, pattern string) ([]SubjectCount) {
	subjectCounts := []SubjectCount{}
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
  project.project_name AS project_name,
  counts.subject_count AS subject_count,
  refresh_date.refresh_date AS refresh_date
FROM counts
  JOIN project
    ON project.id = counts.project_id
  JOIN rave_url
    ON project.url_id = rave_url.id
  LEFT JOIN refresh_date
    ON refresh_date.project_id = project.id
ORDER BY  rave_url.url, project.project_name
`
	rows, err := db.Queryx(q, pattern)
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()

	// iterate over rows
	for rows.Next() {
		var r SubjectCount
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		subjectCounts = append(subjectCounts, r)
	}
	return subjectCounts
}

func getUselessEdits(db *sqlx.DB, pattern string) ([]UnusedEdit) {
	unusedEdits := []UnusedEdit{}
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
SELECT CASE WHEN rave_url.url LIKE '%hdcvc%' THEN rave_url.alternate_url ELSE rave_url.url END AS URL,
  project.project_name AS project_name,
  edit_check_name AS edit_check_name,
  edit_check_count AS edit_check_count,
  CASE WHEN open_query_count > 0 THEN 'Yes' ELSE 'No' END AS open_query
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

	// Export the results
	for rows.Next() {
		var r UnusedEdit
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		unusedEdits = append(unusedEdits, r)
	}
	return unusedEdits
}

// export the study metrics, this is a map with a key matching the URL, and the projects as values
func getStudyMetrics(db *sqlx.DB, pattern string) (map[string][]ProjectVersion) {
	q := `WITH AllData AS (SELECT
                   CASE WHEN rave_url.url LIKE '%hdcvc%'
                     THEN rave_url.alternate_url
                   ELSE rave_url.url END AS URL,
                   rave_url.id AS url_id,
                   project.project_name,
                   crf_version_id,
                   'ACTIVEONLY'          AS check_status,
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
                 WHERE rave_url.url LIKE '%' || $1 || '%' OR rave_url.alternate_url LIKE '%' || $1 || '%'
                 UNION
                 SELECT
                   CASE WHEN rave_url.url LIKE '%hdcvc%'
                     THEN rave_url.alternate_url
                   ELSE rave_url.url END AS URL,
                   rave_url.id AS url_id,
                   project.project_name,
                   crf_version_id,
                   'ALLCHECKS'           AS check_status,
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
                 WHERE rave_url.url LIKE '%' || $1 || '%' OR rave_url.alternate_url LIKE '%' || $1 || '%'
                 ORDER BY url, project_name, crf_version_id, check_status),
    SubjectData AS (SELECT
                      project_name,
                      MAX(subject_count)  AS subject_count,
                      MAX(crf_version_id) AS last_version
                    FROM edit_check chk
                      JOIN project prj
                        ON chk.project_id = prj.id
                    GROUP BY project_name)
SELECT
  AllData.URL,
  AllData.url_id,
  AllData.project_name,
  AllData.crf_version_id,
  CASE WHEN AllData.crf_version_id = SubjectData.last_version
    THEN TRUE
  ELSE FALSE END            AS last_version,
  SubjectData.subject_count AS subject_count,
  AllData.check_status AS check_status,
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
FROM AllData
  JOIN SubjectData
    ON AllData.project_name = SubjectData.project_name
ORDER BY URL, project_name, crf_version_id, check_status`
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
		// First Row
		if project_version == nil {
			project_version = createProjectVersion(r)
		} else if r.URL != project_version.URL || r.ProjectName != project_version.ProjectName || r.CRFVersionID != project_version.CRFVersionID {
			// remove the .mdsol.com
			// refresh the the inactive counts
			prefix := strings.Split(project_version.URL, ".")[0]
			// Data munging
			project_version.fixUpNullValues()
			project_version.calculateInactiveCounts()
			// store for posterity
			urls[prefix] = append(urls[prefix], *project_version)
			project_version = createProjectVersion(r)
		}
		if r.CheckStatus == "ALLCHECKS" {
			project_version.AllEdits = r
		} else {
			project_version.ActiveEditsOnly = r
		}
		//log.Println("Generated",r.CheckStatus,"for",r.ProjectName,"(",r.URL,")")

	}
	// missing last loop
	prefix := strings.Split(project_version.URL, ".")[0]
	project_version = fixUpNullValues(project_version)
	project_version = calculateInactiveCounts(project_version)
	urls[prefix] = append(urls[prefix], *project_version)
	// Log output
	log.Println("Generated metrics for ", len(urls), "URLs")
	return urls
}

// Get the last version dataset for each of the URLs
func getURLLastVersionData(db *sqlx.DB, urls map[string][]ProjectVersion)(map[string][]LastProjectVersion){
	last := make(map[string][]LastProjectVersion)
	for url, versions := range urls{
		versions := getLastVersionDataset(db, versions[0].URLID)
		if len(versions) > 0 {
			last[url] = versions
		} else {
			log.Println("No project versions found for", url)
		}
	}
	return last
}

// get the last version dataset
func getLastVersionDataset(db *sqlx.DB, url_id int) ([]LastProjectVersion){
	q := `WITH SET AS (SELECT
  project_last_version.project_id,
  COUNT(*) AS total_count,
  SUM(CASE WHEN edit_check_name LIKE 'SYS_%' THEN 1 ELSE 0 END) AS fld_total,
  SUM(CASE WHEN edit_check_name LIKE 'SYS_%' AND query_count > 0 THEN 1 ELSE 0 END) AS fld_total_fired,
  SUM(CASE WHEN edit_check_name LIKE 'SYS_%' AND query_count = 0 THEN 1 ELSE 0 END) AS fld_total_not_fired,
  SUM(CASE WHEN edit_check_name LIKE 'SYS_%' AND change_count = 0 THEN 1 ELSE 0 END) AS fld_no_change_count,
  SUM(CASE WHEN edit_check_name LIKE 'SYS_%' AND change_count > 0 THEN 1 ELSE 0 END) AS fld_change_count,
  SUM(CASE WHEN edit_check_name NOT LIKE 'SYS_%' THEN 1 ELSE 0 END) AS prg_total,
  SUM(CASE WHEN edit_check_name NOT LIKE 'SYS_%' AND query_count > 0 THEN 1 ELSE 0 END) AS prg_total_fired,
  SUM(CASE WHEN edit_check_name NOT LIKE 'SYS_%' AND query_count = 0 THEN 1 ELSE 0 END) AS prg_total_not_fired,
  SUM(CASE WHEN edit_check_name NOT LIKE 'SYS_%' AND change_count = 0 THEN 1 ELSE 0 END) AS prg_no_change_count,
  SUM(CASE WHEN edit_check_name NOT LIKE 'SYS_%' AND change_count > 0 THEN 1 ELSE 0 END) AS prg_change_count
FROM edit_check
  INNER JOIN project_last_version
    ON edit_check.project_id = project_last_version.project_id AND
       edit_check.crf_version_id = project_last_version.crf_version_id
WHERE url_id = $1 AND actions LIKE '%OpenQuery%' AND is_active = 1
GROUP BY project_last_version.project_id
)
SELECT project.project_name,
  project_last_version.crf_version_id,
  project_last_version.subject_count,
  SET.* FROM SET
  JOIN project ON project_id = project.id
  JOIN project_last_version ON project.id = project_last_version.project_id
ORDER BY project.project_name;
`
	rows, err := db.Queryx(q, url_id)
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()
	project_versions := []LastProjectVersion{}
	for rows.Next() {
		var r LastProjectVersion
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		// Add the percentages
		r.calculatePercentages()
		project_versions = append(project_versions, r)
	}
	return project_versions
}