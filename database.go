package main

import (
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

func listURLs(db *sqlx.DB) (urls []string, err error) {
	q := `SELECT url, alternate_url FROM rave_url ORDER BY url, alternate_url`
	rows, err := db.Queryx(q)
	if err != nil {
		return
	}
	defer rows.Close()
	// iterate over rows
	for rows.Next() {
		var (
			mainURL string
			altURL  string
		)
		if err = rows.Scan(&mainURL, &altURL); err != nil {
			return
		}
		if altURL != "" {
			urls = append(urls, altURL)
		} else {
			urls = append(urls, mainURL)
		}
	}
	return
}

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
func getSubjectCounts(db *sqlx.DB, pattern string) []SubjectCount {
	subjectCounts := []SubjectCount{}
	q := `WITH counts AS (SELECT
                  edt.project_id         AS project_id,
                  MAX(subject_count) AS subject_count
                FROM edit_check edt
                  JOIN project prj
                    ON edt.project_id = prj.id
                  JOIN rave_url url
                    ON edt.url_id = url.id
                WHERE url.url LIKE '%' || $1 || '%' OR url.alternate_url LIKE '%' || $1 || '%'
                GROUP BY edt.project_id)
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

func getUselessEdits(db *sqlx.DB, pattern string) []UnusedEdit {
	unusedEdits := []UnusedEdit{}
	q := `WITH total AS (
    SELECT
		rave_url.id as url_id,
      	project_id as project_id,
      	edit_check_name,
        COUNT(*) AS edit_check_count,
      	SUM(CASE WHEN Actions LIKE '%OpenQuery%' THEN 1 ELSE 0 END) AS open_query_count,
      	SUM(CASE WHEN Actions LIKE '%CustomFunction%' THEN 1 ELSE 0 END) AS custom_function_count,
      	SUM(total_check_executions) AS total_count
          FROM edit_check
      JOIN rave_url ON edit_check.url_id = rave_url.id
    WHERE rave_url.url LIKE '%' || $1 || '%' OR rave_url.alternate_url LIKE '%' || $1 || '%'
     AND edit_check.is_active = 1
  GROUP BY rave_url.id, project_id, edit_check_name
)
SELECT CASE WHEN rave_url.url LIKE '%hdcvc%' THEN rave_url.alternate_url ELSE rave_url.url END AS URL,
  project.project_name AS project_name,
  edit_check_name AS edit_check_name,
       (SELECT array_to_string(array_remove(array_agg(DISTINCT chk.form_oid), NULL), '|')
        FROM edit_check chk
        WHERE chk.edit_check_name = total.edit_check_name)                                     AS form_oids,
       (SELECT array_to_string(array_remove(array_agg(DISTINCT chk.field_oid), NULL), '|')
        FROM edit_check chk
        WHERE chk.edit_check_name = total.edit_check_name)                                     AS field_oids,
       (SELECT array_to_string(array_remove(array_agg(DISTINCT chk.variable_oid), NULL), '|')
        FROM edit_check chk
        WHERE chk.edit_check_name = total.edit_check_name)                                     AS variable_oids,
  edit_check_count AS edit_check_count,
  CASE WHEN open_query_count > 0 THEN 'Yes' ELSE 'No' END 					AS open_query,
  CASE WHEN custom_function_count > 0 THEN 'Yes' ELSE 'No' END 				AS custom_function,
  CASE WHEN edit_check_name LIKE 'SYS_NC_%' THEN 'Yes' ELSE 'No' END 		AS non_conformant,
  CASE WHEN edit_check_name LIKE 'SYS_Q_RANGE_%' THEN 'Yes' ELSE 'No' END 	AS range_checks,
  CASE WHEN edit_check_name LIKE 'SYS_REQ_%' THEN 'Yes' ELSE 'No' END 		AS required_check,
  CASE WHEN edit_check_name LIKE 'SYS_FUTURE_%' THEN 'Yes' ELSE 'No' END 	AS future_checks
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
func getStudyMetrics(db *sqlx.DB, pattern string) map[string][]ProjectVersion {
	q := `WITH URL AS (
  SELECT id AS url_id
  FROM rave_url
  WHERE (rave_url.url LIKE '%' || $1 || '%' OR rave_url.alternate_url LIKE '%' || $1 || '%')),
     AllData AS (
       SELECT edt.url_id
            ,edt.project_id
            ,edt.crf_version_id
            -- Flag active vs inactive checks
            ,CASE WHEN is_active = 1 THEN 'ACTIVEONLY' ELSE 'INACTIVEONLY' END  AS check_status
            -- total field checks
            , SUM(CASE WHEN edit_check_name LIKE 'SYS_%' THEN 1 ELSE 0 END) AS total_edits_fld
            -- total field queries that have fired once
            , SUM(CASE
                    WHEN edit_check_name LIKE 'SYS_%' AND query_count > 0
                      THEN 1
                    ELSE 0 END)                                             AS total_fired_fld
            -- total field queries that have not fired once
            , SUM(CASE
                    WHEN edit_check_name LIKE 'SYS_%' AND total_check_executions = 0
                      THEN 1
                    ELSE 0 END)                                             AS total_not_fired_fld
            -- total programmed checks
            , SUM(CASE
                    WHEN edit_check_name NOT LIKE 'SYS_%'
                      THEN 1
                    ELSE 0 END)                                             AS total_edits_prg
            -- total programmed queries that have fired once
            , SUM(CASE
                    WHEN edit_check_name NOT LIKE 'SYS_%' AND total_check_executions > 0
                      THEN 1
                    ELSE 0 END)                                             AS total_fired_prg
            -- total programmed queries that have not fired once (note, the default int value is -1,
            --  so if we get a null they won't be included in the count
            , SUM(CASE
                    WHEN edit_check_name NOT LIKE 'SYS_%' AND total_check_executions = 0
                      THEN 1
                    ELSE 0 END)                                             AS total_not_fired_prg
            -- total queries (across all types)
            , SUM(total_check_executions)                                   AS total_queries
            -- total field check queries
            , SUM(CASE
                    WHEN edit_check_name LIKE 'SYS_%' THEN
                      total_check_executions
                    ELSE 0 END)                                             AS total_queries_fld
            -- total programmed check queries
            , SUM(CASE
                    WHEN edit_check_name NOT LIKE 'SYS_%' THEN
                      total_check_executions
                    ELSE 0 END)                                             AS total_queries_prg
            -- total field checks with OpenQuery Action
            , SUM(CASE
                    WHEN edit_check_name LIKE 'SYS_%' AND actions LIKE '%OpenQuery%'
                      THEN 1
                    ELSE 0 END)                                             AS total_edits_query_fld
            -- total field queries with OpenQuery Action
            , SUM(CASE
                    WHEN edit_check_name LIKE 'SYS_%' AND actions LIKE '%OpenQuery%'
                      THEN total_check_executions
                    ELSE 0 END)                                             AS total_queries_query_fld
            -- total field queries with OpenQuery Action that have fired once
            , SUM(CASE
                    WHEN edit_check_name LIKE 'SYS_%' AND total_check_executions > 0
                      THEN 1
                    ELSE 0 END)                                             AS total_fired_query_fld
            -- total field queries with OpenQuery Action that have not fired once
            , SUM(CASE
                    WHEN edit_check_name LIKE 'SYS_%' AND total_check_executions = 0
                      THEN 1
                    ELSE 0 END)                                             AS total_not_fired_query_fld
            -- count of field checks that have fired, but never led to a change in the data
            , SUM(CASE
                    WHEN edit_check_name LIKE 'SYS_%' AND change_count = 0 AND no_change_count > 0
                      THEN 1
                    ELSE 0 END)                                             AS fired_no_change_fld
            -- count of field checks that have fired, and led to a change in the data
            , SUM(CASE
                    WHEN edit_check_name LIKE 'SYS_%' AND change_count > 0
                      THEN 1
                    ELSE 0 END)                                             AS fired_change_fld
            -- total programmed checks with OpenQuery Action
            , SUM(CASE
                    WHEN edit_check_name NOT LIKE 'SYS_%' AND actions LIKE '%OpenQuery%'
                      THEN 1
                    ELSE 0 END)                                             AS total_edits_query_prg
            -- total programmed queries with OpenQuery Action
            , SUM(CASE
                    WHEN edit_check_name NOT LIKE 'SYS_%' AND actions LIKE '%OpenQuery%'
                      THEN total_check_executions
                    ELSE 0 END)                                             AS total_queries_query_prg
            -- total programmed queries with OpenQuery Action that have fired once
            , SUM(CASE
                    WHEN edit_check_name NOT LIKE 'SYS_%' AND actions LIKE '%OpenQuery%'
                           AND total_check_executions > 0
                      THEN 1
                    ELSE 0 END)                                             AS total_fired_query_prg
            -- total programmed queries with OpenQuery Action that have not fired once
            , SUM(CASE
                    WHEN edit_check_name NOT LIKE 'SYS_%' AND actions LIKE '%OpenQuery%'
                           AND total_check_executions = 0
                      THEN 1
                    ELSE 0 END)                                             AS total_not_fired_query_prg
            -- count of programmed checks that have fired, but not led to a change in the data
            , SUM(CASE
                    WHEN edit_check_name NOT LIKE 'SYS_%' AND no_change_count > 0
                           AND change_count = 0
                      THEN 1
                    ELSE 0 END)                                             AS fired_no_change_prg
            -- count of programmed checks that have fired and led to a change in the data
            , SUM(CASE
                    WHEN edit_check_name NOT LIKE 'SYS_%' AND change_count > 0
                      THEN 1
                    ELSE 0 END)                                             AS fired_change_prg
       FROM edit_check edt
              JOIN URL
                   ON URL.url_id = edt.url_id
       GROUP BY edt.url_id, edt.project_id, edt.crf_version_id, edt.is_active
     ),
     SubjectData AS (SELECT chk.project_id,
                            MAX(subject_count)  AS subject_count,
                            MAX(crf_version_id) AS last_version
                     FROM edit_check chk
                     GROUP BY chk.project_id)
SELECT (SELECT url from rave_url where rave_url.id = adt.url_id)            AS url,
       adt.url_id,
       (SELECT project_name from project where project.id = adt.project_id) AS project_name,
       adt.crf_version_id,
       CASE
         WHEN adt.crf_version_id = sdt.last_version
           THEN TRUE
         ELSE FALSE END                                                     AS last_version,
       sdt.subject_count                                                    AS subject_count,
       adt.check_status                                                     AS check_status,
       total_edits_fld,
       total_fired_fld,
       fired_change_fld,
       fired_no_change_fld,
       total_not_fired_fld,
       total_edits_query_fld,
       total_fired_query_fld,
       total_not_fired_query_fld,
       total_edits_prg,
       total_fired_prg,
       total_not_fired_prg,
       fired_change_prg,
       fired_no_change_prg,
       total_not_fired_prg,
       total_edits_query_prg,
       total_fired_query_prg,
       total_not_fired_query_prg,
       total_queries,
       total_queries_fld,
       total_queries_prg
FROM AllData adt
       JOIN SubjectData sdt
            ON adt.project_id = sdt.project_id
ORDER BY URL, project_name, crf_version_id, check_status`
	rows, err := db.Queryx(q, pattern)
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()

	// Status variable
	urls := make(map[string][]ProjectVersion)
	var projectVersion *ProjectVersion
	for rows.Next() {
		var r Record
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		// First Row
		if projectVersion == nil {
			projectVersion = createProjectVersion(r)
		} else if r.URL != projectVersion.URL || r.ProjectName != projectVersion.ProjectName || r.CRFVersionID != projectVersion.CRFVersionID {
			// remove the .mdsol.com
			// refresh the the inactive counts
			prefix := strings.Split(projectVersion.URL, ".")[0]
			// Data munging
			projectVersion.fixUpNullValues()
			// projectVersion.calculateInactiveCounts()
			// store for posterity
			urls[prefix] = append(urls[prefix], *projectVersion)
			projectVersion = createProjectVersion(r)
		}
		if r.CheckStatus == "ACTIVEONLY" {
			projectVersion.ActiveEditsOnly = &r
		} else {
			projectVersion.InactiveEditsOnly = &r
		}
		//log.Println("Generated",r.CheckStatus,"for",r.ProjectName,"(",r.URL,")")
	}
	// missing last loop
	prefix := strings.Split(projectVersion.URL, ".")[0]
	//projectVersion = fixUpNullValues(projectVersion)
	//projectVersion = calculateInactiveCounts(projectVersion)
	urls[prefix] = append(urls[prefix], *projectVersion)
	// Log output
	log.Println("Generated metrics for ", len(urls), "URLs")
	return urls
}

// Get the last version dataset for each of the URLs
func getURLLastVersionData(db *sqlx.DB, urls map[string][]ProjectVersion) map[string][]LastProjectVersion {
	last := make(map[string][]LastProjectVersion)
	for url, versions := range urls {
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
func getLastVersionDataset(db *sqlx.DB, URLID int) []LastProjectVersion {
	// Note we only pull out the checks that can register (ie with OpenQuery)
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
  JOIN project ON SET.project_id = project.id
  JOIN project_last_version ON project.id = project_last_version.project_id
ORDER BY project.project_name;
`
	rows, err := db.Queryx(q, URLID)
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()
	projectVersions := []LastProjectVersion{}
	for rows.Next() {
		var r LastProjectVersion
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		// Add the percentages
		r.calculatePercentages()
		projectVersions = append(projectVersions, r)
	}
	return projectVersions
}
