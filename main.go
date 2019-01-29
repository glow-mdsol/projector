package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
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

func main() {
	var patternsArray, raveUrls arrayFlags
	flag.Var(&patternsArray, "pattern", "Supply the URL patterns")
	flag.Var(&raveUrls, "url", "Specific Rave URLs")
	dumpURLs := flag.Bool("listurls", false, "Dump the list of urls")
	hostName := flag.String("dbhost", "localhost", "Database Host")
	dbName := flag.String("dbname", "editsfive", "Database Name")
	dbUser := flag.String("user", "edits", "Database User")
	dbPass := flag.String("password", "apple01", "Database Password")
	fileName := flag.String("output", "report", "Output File Name")
	threshold := flag.Int("threshold", 10, "Threshold for Reporting")
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
	workbook := xlsx.NewFile()
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
		if !doesPatternMatch(urlPattern, dbConn) {
			log.Println("No matching URLs for", urlPattern)
			continue
		}
		log.Println("Processing URL Pattern ", urlPattern)
		// Get the Subject Counts
		log.Println("Retrieving Subject Counts")
		subjectCounts := getSubjectCounts(dbConn, urlPattern)
		// Get the unfired edits
		log.Println("Retrieving Unfired Edits")
		uselessEdits := getUselessEdits(dbConn, urlPattern)
		// Get the Study Metrics
		log.Println("Retrieving URL Metrics")
		studyMetrics := getStudyMetrics(dbConn, urlPattern)
		// Get the LastVersionData
		lastVersions := getURLLastVersionData(dbConn, studyMetrics)
		// Subject Counts
		log.Println("Writing Subject Counts")
		writeSubjectCounts(subjectCounts, workbook)
		// Useless Edits
		log.Println("Writing Unfired Edits")
		writeUselessEdits(uselessEdits, workbook)
		// Study Metrics
		log.Println("Writing Study Metrics")
		writeStudyMetrics(studyMetrics, workbook)
		// Last Project Versions
		log.Println("Writing Last Project Version Data")
		writeLastProjectVersions(lastVersions, *threshold, workbook)
	}
	// make up the prefix using the range of patterns, removing the extraneous domains
	var prefixes []string
	for _, prefix := range patternsArray {
		if strings.HasSuffix(prefix, ".mdsol.com") {
			prefixes = append(prefixes, strings.Split(prefix, ".")[0])
		} else {
			prefixes = append(prefixes, prefix)
		}
	}
	prefix := strings.Join(prefixes, "_")
	filename := fmt.Sprintf("%s_%s_%s.xlsx", prefix, *fileName, time.Now().Format("2006-01-02"))
	err = workbook.Save(filename)
	if err != nil {
		log.Fatalf("Saving file failed: %s", err)
	}
}
