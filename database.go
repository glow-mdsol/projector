package main

import (
	"github.com/jmoiron/sqlx"
	"log"
)

// does the pattern return any Rave URLS by name
//func doesPatternMatch(pattern string, db *sqlx.DB) bool {
//	q := `SELECT COUNT(*) FROM rave_url
//	WHERE rave_url.url LIKE '%' || $1 || '%' OR rave_url.alternate_url LIKE '%' || $1 || '%'`
//	var count int
//	err := db.Get(&count, q, pattern)
//	if err == nil {
//		return count != 0
//	}
//	return false
//}

// get RaveURLS that match the pattern
func GetURLsThatMatch(db *sqlx.DB, pattern string) (urls []RaveURL, err error) {
	q := `SELECT id, url, alternate_url FROM rave_url 
		WHERE rave_url.url LIKE '%' || $1 || '%' 
		OR rave_url.alternate_url LIKE '%' || $1 || '%' `
	rows, err := db.Queryx(q, pattern)
	if err != nil {
		return
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// iterate over rows
	for rows.Next() {
		var r RaveURL
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		urls = append(urls, r)
	}
	return
}

// dump a list of URLs
func listURLs(db *sqlx.DB) (urls []string, err error) {
	q := `SELECT url, alternate_url FROM rave_url ORDER BY url, alternate_url`
	rows, err := db.Queryx(q)
	if err != nil {
		return
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
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

// get the Projects for a URL
func getProjects(db *sqlx.DB, urlID int) (projects []*Project) {
	// NOTE: project.id is the autogenerated value
	q := `SELECT DISTINCT prj.url_id,
		   prj.id AS project_id,
		   prj.project_name
	FROM project prj
	WHERE prj.url_id = $1`
	rows, err := db.Queryx(q, urlID)
	if err != nil {
		log.Fatal("PJ Query failed: ", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// iterate over rows
	for rows.Next() {
		var r Project
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		// log.Println("Project",r.ProjectName,"(",r.ProjectID,")")
		projects = append(projects, &r)
	}
	return
}

// Get the project versions
func getProjectVersions(db *sqlx.DB, projectID int) (projectVersions []*ProjectVersion) {
	q := `SELECT DISTINCT
                edt.project_id AS project_id,
                edt.crf_version_id AS crf_version_id,
                CASE WHEN plv.crf_version_id = edt.crf_version_id THEN 1 ELSE 0 END AS last_version
			FROM edit_check edt
				LEFT JOIN project_last_version plv ON edt.project_id = plv.project_id
			WHERE edt.project_id = $1
			`
	rows, err := db.Queryx(q, projectID)
	if err != nil {
		log.Fatal("PV Query failed: ", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// iterate over rows
	for rows.Next() {
		var r ProjectVersion
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		projectVersions = append(projectVersions, &r)
	}
	return

}

func getSubjectCounts(db *sqlx.DB, urlID int) []SubjectCount {
	var subjectCounts []SubjectCount
	q := `SELECT
	edt.url_id,
	edt.project_id AS project_id,
    pj.project_name AS project_name,
    rd.refresh_date AS refresh_date,
    MAX(subject_count) AS subject_count,
    MAX(screening_subjects) AS screening_subject_count,
    MAX(screening_failure_subjects) AS screening_failure_subject_count,
    MAX(enrolled_subjects) AS enrolled_subject_count,
    MAX(completed_subjects) AS completed_subject_count,
    MAX(enrolled_follow_up_subjects) AS follow_up_subject_count,
    MAX(early_terminated_subjects) AS early_terminated_subject_count
FROM edit_check edt
    JOIN project pj ON edt.project_id = pj.id
    JOIN refresh_date rd on pj.id = rd.project_id
  WHERE edt.url_id = $1
GROUP BY edt.url_id, edt.project_id, pj.project_name, rd.refresh_date;`
	rows, err := db.Queryx(q, urlID)
	if err != nil {
		log.Fatal("SBJS Query failed: ", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// iterate over rows
	for rows.Next() {

		var r SubjectCount
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		//  log.Println(r.ProjectName, "Loaded: ", r.SubjectCount,"for ProjectID",r.ProjectID)
		subjectCounts = append(subjectCounts, r)
	}
	return subjectCounts
}

// export the subject counts
func getProjectSubjectCount(db *sqlx.DB, urlID, projectID int) (subjectCount SubjectCount) {
	q := `SELECT
    edt.project_id AS project_id,
    pj.project_name AS project_name,
    rd.refresh_date AS refresh_date,
    MAX(subject_count) AS subject_count,
    MAX(screening_subjects) AS screening_subject_count,
    MAX(screening_failure_subjects) AS screening_failure_subject_count,
    MAX(enrolled_subjects) AS enrolled_subject_count,
    MAX(completed_subjects) AS completed_subject_count,
    MAX(enrolled_follow_up_subjects) AS follow_up_subject_count,
    MAX(early_terminated_subjects) AS early_terminated_subject_count
FROM edit_check edt
    JOIN project pj ON edt.project_id = pj.id
    JOIN refresh_date rd on pj.id = rd.project_id
  WHERE edt.url_id = $1 AND edt.project_id = $2
GROUP BY edt.project_id, pj.project_name, rd.refresh_date;`
	rows, err := db.Queryx(q, urlID, projectID)
	if err != nil {
		log.Fatal("SBJ Query failed: ", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// iterate over rows
	for rows.Next() {
		if err := rows.StructScan(subjectCount); err != nil {
			log.Fatal(err)
		}
	}
	return
}

func hasCustomFunction(db *sqlx.DB, projectID int, editCheckName string) bool {
	q := `SELECT project_id,
       edit_check_name,
       SUM(CASE WHEN actions LIKE '%CustomFunction%' THEN 1 ELSE 0 END) AS cf_count
FROM edit_check
WHERE project_id = $1 AND edit_check_name = $2
GROUP BY project_id, edit_check_name`
	rows, err := db.Queryx(q, projectID, editCheckName)
	if err != nil {
		log.Fatal("CF Query failed: ", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	type editResult struct {
		ProjectID     int    `db:"project_id"`
		EditCheckName string `db:"edit_check_name"`
		CFCount       int    `db:"cf_count"`
	}
	// Export the results
	for rows.Next() {
		var r editResult
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		return r.CFCount > 0
	}
	return false
}

func getUselessEditsForProject(db *sqlx.DB, projectID int, withOpenQueryFilter EditCheckOutcome) []*UnusedEdit {
	var unusedEdits []*UnusedEdit
	q := `SELECT project_id,
        edit_check_name AS edit_check_name,
        total_count,
       (SELECT array_to_string(array_remove(array_agg(DISTINCT chk.form_oid), NULL), '|')
        FROM edit_check chk
        WHERE chk.edit_check_name = total.edit_check_name)                                     AS form_oids,
       (SELECT array_to_string(array_remove(array_agg(DISTINCT chk.field_oid), NULL), '|')
        FROM edit_check chk
        WHERE chk.edit_check_name = total.edit_check_name)                                     AS field_oids,
       (SELECT array_to_string(array_remove(array_agg(DISTINCT chk.variable_oid), NULL), '|')
        FROM edit_check chk
        WHERE chk.edit_check_name = total.edit_check_name)                                     AS variable_oids
 FROM (SELECT project_id,
               edit_check_name,
               COUNT(*) as total_count,
               SUM(total_check_executions) as total_executions
    FROM edit_check
      WHERE project_id = $1 AND
            CASE WHEN $2 = 0 THEN
		        actions LIKE '%OpenQuery%'
	        ELSE
		        actions NOT LIKE '%OpenQuery%'
	        END
      GROUP BY edit_check_name, project_id) total
WHERE
	total.total_executions = 0
GROUP BY project_id, edit_check_name, total_count`
	rows, err := db.Queryx(q, projectID, int(withOpenQueryFilter))
	if err != nil {
		log.Fatal("BE Query failed: ", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Export the results
	for rows.Next() {
		var r UnusedEdit
		if err := rows.StructScan(&r); err != nil {
			log.Fatal(err)
		}
		if hasCustomFunction(db, r.ProjectID, r.EditCheckName) {
			r.CustomFunction = true
		} else {
			r.CustomFunction = false
		}
		unusedEdits = append(unusedEdits, &r)
	}
	return unusedEdits
}

// get the summary by type
func getStudyMetricsByProjectAndCheckType(db *sqlx.DB, projectID, crfVersionID int, checkType EditCheckClass) (metrics EditTypeMetric) {
	q := `SELECT 
		-- total edits per version
		COUNT(*) 														AS total_edits
		-- total edits with OpenQuery action (filtered to active only)
		, SUM(CASE WHEN actions LIKE '%OpenQuery%'
				  THEN 1
				ELSE 0 END)                                             AS total_edits_with_openquery
		-- total queries (across all types)
		, SUM(total_check_executions)                                   AS total_queries
		-- total open queries
		, SUM(CASE WHEN open_checks >= 0 THEN open_checks ELSE 0 END )	AS total_open_queries
		-- total edits that have fired once
		, SUM(CASE WHEN total_check_executions > 0
				THEN 1
				ELSE 0 END)                                             AS total_edits_fired
		-- total edits that have not fired once
		, SUM(CASE
				WHEN total_check_executions = 0
				  THEN 1
				ELSE 0 END)                                             AS total_edits_not_fired
		-- total checks that have fired once with OpenQuery
		, SUM(CASE
				WHEN total_check_executions > 0 AND actions LIKE '%OpenQuery%'
				  THEN 1
				ELSE 0 END)                                             AS total_edits_fired_with_openquery
		-- total checks that have not fired once with OpenQuery
		, SUM(CASE
				WHEN total_check_executions = 0 AND actions LIKE '%OpenQuery%'
				  THEN 1
				ELSE 0 END)                                             AS total_edits_not_fired_with_openquery
		-- sum of queries fired with Action OpenQuery
		, SUM(CASE
				WHEN actions LIKE '%OpenQuery%'
				THEN
				  total_check_executions
				ELSE 0 END)                                             AS total_queries_open_query
		-- count of checks that have fired and led to a change in the data
		, SUM(CASE WHEN change_count > 0 
				  THEN 1
				ELSE 0 END)                                             AS total_edits_fired_with_change
		-- count of checks that have fired, but never led to a change in the data
		, SUM(CASE
				WHEN change_count = 0 AND no_change_count > 0 AND is_active = 1
				  THEN 1
				ELSE 0 END)                                             AS total_edits_fired_without_change
		-- total count of changes 
		, SUM(CASE
				WHEN change_count > 0
				  THEN change_count
				ELSE 0 END)                                             AS total_changes_from_edits
		-- total count of edits that have open queries 
		, SUM(CASE WHEN open_checks > 0
                THEN 1
				ELSE 0 END)                                             AS total_open_edits
	FROM edit_check edt
	WHERE edt.project_id = $1 AND edt.crf_version_id = $2
		AND CASE WHEN 1 = 1 THEN
			-- pick only active checks
			edt.is_active = 1
		ELSE
			-- pick only inactive checks
			edt.is_active = 0
		END
		AND CASE WHEN $3 = 0 THEN
			-- field edits
			edit_check_name LIKE 'SYS_%'
		ELSE
			-- programmed edits
			edit_check_name NOT LIKE 'SYS_%'
		END
	GROUP BY edt.project_id, edt.crf_version_id
	`
	rows, err := db.Queryx(q, projectID, crfVersionID, int(checkType))
	if err != nil {
		log.Fatal("SM Query failed: ", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// iterate over rows
	for rows.Next() {
		if err := rows.StructScan(&metrics); err != nil {
			log.Fatal(err)
		}
	}
	return
}

// get the counts by edit check status
func getActivityCount(db *sqlx.DB, projectID, crfVersionID int) EditStatusCounts {
	var editStatusCounts EditStatusCounts
	q := `SELECT SUM(CASE WHEN is_active = 1 THEN 1 ELSE 0 END) AS active_count,
       SUM(CASE WHEN is_active = 0 THEN 1 ELSE 0 END) AS inactive_count
		FROM edit_check edt
	WHERE edt.project_id = $1 AND edt.crf_version_id = $2
	GROUP BY edt.project_id, edt.crf_version_id
	`
	rows, err := db.Queryx(q, projectID, crfVersionID)
	if err != nil {
		log.Fatal("AC Query failed: ", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	for rows.Next() {
		if err := rows.StructScan(&editStatusCounts); err != nil {
			log.Fatal(err)
		}
	}
	//log.Println("Loaded Status Counts for",
	//	projectID, "(",
	//	crfVersionID, ") with active edits = ", editStatusCounts.ActiveEdits)

	return editStatusCounts
}
