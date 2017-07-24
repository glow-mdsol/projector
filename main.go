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
)

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
		if !doesPatternMatch(url_pattern, dbConn) {
			log.Println("No matching URLs for", url_pattern)
			continue
		}
		log.Println("Processing URL Pattern ", url_pattern)
		// Get the Subject Counts
		log.Println("Retrieving Subject Counts")
		subject_counts := getSubjectCounts(dbConn, url_pattern)
		// Get the unfired edits
		log.Println("Retrieving Unfired Edits")
		useless_edits := getUselessEdits(dbConn, url_pattern)
		// Get the Study Metrics
		log.Println("Retrieving URL Metrics")
		study_metrics := getStudyMetrics(dbConn, url_pattern)
		// Subject Counts
		log.Println("Writing Subject Counts")
		writeSubjectCounts(subject_counts, workbook)
		// Useless Edits
		log.Println("Writing Unfired Edits")
		writeUselessEdits(useless_edits, workbook)
		// Study Metrics
		log.Println("Writing Study Metrics")
		writeStudyMetrics(study_metrics, workbook)
	}
	// make up the prefix using the range of patterns
	prefix := strings.Join(patternsArray, "_")
	filename := fmt.Sprintf("%s_%s_%s.xlsx", prefix, *file_name, time.Now().Format("2006-01-02"))
	workbook.Save(filename)
}
