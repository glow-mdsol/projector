package main

import (
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func getURLs(db *sqlx.DB) {
	urls, err := listURLs(db)
	if err != nil {
		log.Fatal("Unable to list URLs: ", err)
	}
	for _, url := range urls {
		fmt.Println(url)
	}
}

// load the queries, versions, etc
func loadProject(db *sqlx.DB, urlID int, project *Project) {
	// load in the UnusedQueries
	project.loadUnusedQueries(db)
	// get the versions
	projectVersions := getProjectVersions(db, project.ProjectID)
	for _, projectVersion := range projectVersions {
		projectVersion.getActivityCounts(db)
		projectVersion.getMetrics(db)
	}
	// ensure the versions are ordered appropriately
	project.Versions = orderVersions(projectVersions)
}

// fluff out a project definition
func expandProject(db *sqlx.DB, raveURL RaveURL, project *Project, subjectCounts []SubjectCount) {
	log.Println("Expanding ", project.ProjectName)
	// TODO: Concurrency
	loadProject(db, raveURL.URLID, project)
	for _, counts := range subjectCounts {
		if counts.ProjectID == project.ProjectID {
			project.SubjectCount = counts
		}
	}
}

// process a RaveURL dataset
func processRaveURL(db *sqlx.DB, raveURL RaveURL) {
	workbook := xlsx.NewFile()
	//if !doesPatternMatch(urlPattern, dbConn) {
	//	log.Println("No matching URLs for", urlPattern)
	//	continue
	//}
	log.Println("Processing Rave URL ", raveURL.URL())
	// get the projects
	projects := getProjects(db, raveURL.URLID)
	log.Println("Loaded", len(projects), "Projects")
	// sort the projects
	projects = orderProjects(projects)
	// load the subjectCounts
	subjectCounts := getSubjectCounts(db, raveURL.URLID)
	// Get the project versions
	for _, project := range projects {
		// can we parallelise this?
		expandProject(db, raveURL, project, subjectCounts)
	}
	// WRITE OUT THE SUBJECT COUNTS
	writeSubjectCount(raveURL.URL(), projects, workbook)
	// Process useless edits project by project
	for _, project := range projects {
		// OpenQuery
		writeUselessEdits(project.ProjectName, project.UnusedWithOpenQuery, OpenQuery, workbook)
		// Not OpenQuery
		writeUselessEdits(project.ProjectName, project.Unused, WithoutOpenQuery, workbook)
		// versions
		writeStudyMetricsForProject(raveURL.URLPrefix(), project, workbook)
		// last version
		writeLastStudyMetricsForProject(raveURL.URLPrefix(), project, workbook)

	}
	// aggregated counts
	writeSummaryCounts(projects, workbook)

	// write to disk
	filename := fmt.Sprintf("%s_%s.xlsx", raveURL.URLPrefix(), time.Now().Format("2006-01-02"))

	err := workbook.Save(filename)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func main() {
	var patternsArray, raveUrls arrayFlags
	flag.Var(&patternsArray, "pattern", "Supply the URL patterns")
	flag.Var(&raveUrls, "url", "Specific Rave URLs")
	dumpURLs := flag.Bool("listurls", false, "Dump the list of urls")
	hostName := flag.String("dbhost", "localhost", "Database Host")
	dbName := flag.String("dbname", "editsfive", "Database Name")
	dbUser := flag.String("user", "edits", "Database User")
	dbPass := flag.String("password", "apple01", "Database Password")
	//fileName := flag.String("output", "report", "Output File Name")
	//threshold := flag.Int("threshold", 10, "Threshold for Reporting")
	flag.Parse()
	if *dumpURLs == false && (len(patternsArray) == 0 && len(raveUrls) == 0) {
		log.Fatal("Need to specify the patterns or url")
	}
	var dataSourceName = fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		*hostName,
		*dbUser,
		*dbName,
		*dbPass)
	var dbConn *sqlx.DB
	// make the database connection
	dbConn, err := sqlx.Open("postgres", string(dataSourceName))
	if err != nil {
		log.Fatal(err)
	}
	if *dumpURLs == true {
		getURLs(dbConn)
		os.Exit(0)
	}
	if len(raveUrls) != 0 {
		for _, raveURL := range raveUrls {
			if !strings.HasSuffix(raveURL, ".mdsol.com") {
				// if we don't end with mdsol.com, then set it
				patternsArray.Set(fmt.Sprintf("%s.mdsol.com", raveURL))
			} else {
				patternsArray.Set(raveURL)
			}

		}
	}
	for _, urlPattern := range patternsArray {
		matchingURLs, err := GetURLsThatMatch(dbConn, urlPattern)
		if err != nil {
			continue
		}
		if len(matchingURLs) == 0 {
			log.Println("No matching URLs for", urlPattern)
			continue
		}

		for _, raveURL := range matchingURLs {
			processRaveURL(dbConn, raveURL)
		}

	}
}
