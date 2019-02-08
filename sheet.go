package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/tealeg/xlsx"
)

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
	for idx, datum := range data {
		cell := row.AddCell()
		cell.SetStyle(titleFace)
		cell.SetString(datum)
		sheet.SetColWidth(idx, idx, float64(len(datum)))
	}
}

// Write a Project Version to the Sheet
func writeProjectVersion(subjectCount int, prj *ProjectVersion, row *xlsx.Row) {
	sheet := row.Sheet
	var cell *xlsx.Cell
	// Study URL
	cell = row.AddCell()
	cell.Value = prj.URL
	sheet.SetColWidth(0, 0, float64(len(prj.URL)))
	// Project Name
	cell = row.AddCell()
	cell.Value = prj.ProjectName
	if float64(len(prj.ProjectName)) > sheet.Col(1).Width {
		sheet.SetColWidth(1, 1, float64(len(prj.ProjectName)))
	}
	// CRF Version
	cell = row.AddCell()
	cell.SetInt(prj.CRFVersionID)
	// Last Version?
	cell = row.AddCell()
	if prj.LastVersion {
		cell.SetString("Y")
	} else {
		cell.SetString("N")
	}
	// Current Subject Count
	cell = row.AddCell()
	cell.SetInt(subjectCount)
	// Active Edits
	cell = row.AddCell()
	cell.SetInt(prj.ActiveEditCount)
	// Inactive Edits
	cell = row.AddCell()
	cell.SetInt(prj.InactiveEditCount)
	// Write the metrics out
	writeMetrics(prj.ActiveEditsOnly, row)
	// Inactive Edits - Removed
	// writeMetrics(prj.InactiveEditsOnly, row)
}

// write a set of metrics to the file
func writeMetrics(rec *Record, row *xlsx.Row) {
	var cell *xlsx.Cell
	// Total Edits (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEdits)
	// Total Edits Fired (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEditsFired)
	// Total Edits Unfired (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEditsNotFired)
	// Total Edits Open (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEditsOpen)
	// %ge Edits Fired (fld)
	cell = row.AddCell()
	if rec.TotalFieldEdits == 0 {
		cell.SetFloat(0.0)
	} else {
		cell.SetFloatWithFormat(float64(rec.TotalFieldEditsFired)/float64(rec.TotalFieldEdits)*100.0, "#,##0.00;(#,##0.00)")
	}
	// %ge Edits Unfired (fld)
	cell = row.AddCell()
	if rec.TotalFieldEdits == 0 {
		cell.SetFloat(0.0)
	} else {
		cell.SetFloatWithFormat(float64(rec.TotalFieldEditsNotFired)/float64(rec.TotalFieldEdits)*100.0, "#,##0.00;(#,##0.00)")
	}
	// Total Edits (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEdits)
	// Total Edits With OpenQuery (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEditsWithOpenQuery)
	// Total Edits Fired (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEditsFired)
	// Total Edits Unfired (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEditsNotFired)
	// Total Edits Open (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEditsOpen)
	// %ge Edits Fired (prg)
	cell = row.AddCell()
	if rec.TotalProgEditsWithOpenQuery == 0 {
		cell.SetFloat(0.0)
	} else {
		cell.SetFloatWithFormat(float64(rec.TotalProgEditsFired)/float64(rec.TotalProgEdits)*100.0, "#,##0.00;(#,##0.00)")
	}
	// %ge Edits Unfired (prg)
	cell = row.AddCell()
	if rec.TotalProgEditsWithOpenQuery == 0 {
		cell.SetFloat(0.0)
	} else {
		cell.SetFloatWithFormat(float64(rec.TotalProgEditsNotFired)/float64(rec.TotalProgEdits)*100.0, "#,##0.00;(#,##0.00)")
	}
	// Total Queries (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldQueries)
	// Total Queries (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgQueries)
	// Total Queries With OpenQuery (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalQueriesOpenQuery)
}

// write the Subject Counts
func writeSubjectCounts(subjectCounts []SubjectCount, wbk *xlsx.File) {
	tabName := "Subject Counts"
	headers := []string{"Rave URL", "Project Name", "Subject Count", "Date Updated"}

	// create the sheet
	sheet, created := getOrAddSheet(wbk, tabName)
	if created {
		// Add the headers
		writeHeaderRow(headers, sheet)
	}
	autoFilter := new(xlsx.AutoFilter)
	autoFilter.TopLeftCell = "A1"
	autoFilter.BottomRightCell = "D1"
	sheet.AutoFilter = autoFilter

	// default widths
	urlLength := 12
	projectLength := 12
	dateLength := 18.5

	for _, subjectCount := range subjectCounts {
		var cell *xlsx.Cell
		// Rows
		row := sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(subjectCount.URL)
		cell = row.AddCell()
		cell.SetString(subjectCount.ProjectName)
		cell = row.AddCell()
		if subjectCount.SubjectCount.Valid {

			cell.SetInt(int(subjectCount.SubjectCount.Int64))
		} else {
			cell.SetString("-")
		}
		cell = row.AddCell()
		if subjectCount.RefreshDate.Valid {
			cell.SetDateTime(subjectCount.RefreshDate.Time)

		} else {
			cell.SetString("-")
		}
		if len(subjectCount.URL) > urlLength {
			urlLength = len(subjectCount.URL)
		}
		if len(subjectCount.ProjectName) > projectLength {
			projectLength = len(subjectCount.ProjectName)
		}

	}
	err := sheet.SetColWidth(0, 0, float64(urlLength))
	if err != nil {
		log.Printf("Error setting column width: %s", err)
	}
	err = sheet.SetColWidth(1, 1, float64(projectLength))
	if err != nil {
		log.Printf("Error setting column width: %s", err)
	}
	err = sheet.SetColWidth(3, 3, float64(dateLength))
	if err != nil {
		log.Printf("Error setting column width: %s", err)
	}
}

// write study metrics
func writeStudyMetrics(data map[string]*RaveURL, wbk *xlsx.File) {
	headers := []string{"Study URL",
		"Project Name",
		"CRF Version",
		"Last Version",
		"Subject Count",
		"Active Edits",
		"Inactive Edits",
		"Total Edits (fld)",
		"Total Edits Fired (fld)",
		"Total Edits Unfired (fld)",
		"Total Edits Open (fld)",
		"%ge Edits Fired (fld)",
		"%ge Edits Unfired (fld)",
		"Total Edits (prg)",
		"Total Edits With OpenQuery (prg)",
		"Total Edits Fired (prg)",
		"Total Edits Unfired (prg)",
		"Total Edits Open (prg)",
		"%ge Edits Fired (prg)",
		"%ge Edits Unfired (prg)",
		"Total Queries (fld)",
		"Total Queries (prg)",
		"Total Queries with OpenQuery",
		//"Total Queries With OpenQuery (prg)",
		// "Inactive - Total Edits (fld)",
		// "Inactive - Total Edits Fired (fld)",
		// "Inactive - Total Edits Unfired (fld)",
		// "Inactive - %ge Edits Fired (fld)",
		// "Inactive - %ge Edits Unfired (fld)",
		// "Inactive - Total Edits (prg)",
		// "Inactive - Total Edits With OpenQuery (prg)",
		// "Inactive - Total Edits Fired (prg)",
		// "Inactive - Total Edits Unfired (prg)",
		// "Inactive - %ge Edits Fired (prg)",
		// "Inactive - %ge Edits Unfired (prg)",
		// "Inactive - Total Queries (fld)",
		// "Inactive - Total Queries (prg)",
		// "Inactive - Total Queries With OpenQuery (prg)",
		//"Total Edits Fired With No Change (fld)",
		//"Total Edits Fired With No Change (prg)"
	}
	var urls []string
	for k := range data {
		// put the urls out in order
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
		autoFilter := new(xlsx.AutoFilter)
		autoFilter.TopLeftCell = "A1"
		autoFilter.BottomRightCell = "D1"
		sheet.AutoFilter = autoFilter
		// Get the URL instance
		raveURL, ok := data[url]
		if ok == false {
			log.Fatalf("Unable to locate %s", url)
		}
		//log.Println("Created Sheet for URL ", url)
		projects := raveURL.getProjects()
		for _, project := range projects {
			projectVersions := project.getVersions()
			for _, projectVersion := range projectVersions {
				//// Add the row for Checks
				row := sheet.AddRow()
				writeProjectVersion(project.SubjectCount, projectVersion, row)
			}
		}
	}
}

func writeUselessEdits(edits []UnusedEdit, wbk *xlsx.File) {
	headers := []string{"Study URL",
		"Project Name",
		"Edit Check Name",
		"Form OID",
		"Field OID",
		"Variable OID",
		"Times Used",
		"OpenQuery Check?",
		"Custom Function?",
		"Non-conformance check?",
		"Required check?",
		"Future check?",
		"Range check?",
	}

	tabName := "Unused Edits"

	// create the sheet
	sheet, created := getOrAddSheet(wbk, tabName)
	if created {
		// Add the headers
		writeHeaderRow(headers, sheet)
	}
	autoFilter := new(xlsx.AutoFilter)
	autoFilter.TopLeftCell = "A1"
	autoFilter.BottomRightCell = "M1"
	sheet.AutoFilter = autoFilter

	urlLength := 12
	projectLength := 12
	checkLength := 12
	formOIDLength := 12
	fieldOIDLength := 12
	vblOIDLength := 12
	maxLength := 70
	// Export the results
	for _, edit := range edits {
		if len(edit.URL) > urlLength {
			urlLength = len(edit.URL)
		}
		if len(edit.ProjectName) > projectLength {
			projectLength = len(edit.ProjectName)
		}
		if len(edit.EditCheckName) > checkLength {
			if len(edit.EditCheckName) < maxLength {
				checkLength = len(edit.EditCheckName)
			} else {
				checkLength = maxLength
			}
		}
		if len(edit.FormOID) > formOIDLength {
			if len(edit.FormOID) < maxLength {
				formOIDLength = len(edit.FormOID)
			} else {
				formOIDLength = maxLength
			}
		}
		if len(edit.FieldOID) > fieldOIDLength {
			if len(edit.FieldOID) < maxLength {
				fieldOIDLength = len(edit.FieldOID)
			} else {
				fieldOIDLength = maxLength
			}

		}
		if len(edit.VariableOID) > vblOIDLength {
			if len(edit.VariableOID) < maxLength {
				vblOIDLength = len(edit.VariableOID)
			} else {
				vblOIDLength = maxLength
			}
		}
		var cell *xlsx.Cell
		// Rows
		row := sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(edit.URL)
		cell = row.AddCell()
		cell.SetString(edit.ProjectName)
		cell = row.AddCell()
		cell.SetString(edit.EditCheckName)
		cell = row.AddCell()
		cell.SetString(edit.FormOID)
		cell = row.AddCell()
		cell.SetString(edit.FieldOID)
		cell = row.AddCell()
		cell.SetString(edit.VariableOID)
		cell = row.AddCell()
		cell.SetInt(edit.UsageCount)
		cell = row.AddCell()
		cell.SetString(edit.OpenQuery)
		cell = row.AddCell()
		cell.SetString(edit.CustomFunction)
		cell = row.AddCell()
		cell.SetString(edit.NonConformant)
		cell = row.AddCell()
		cell.SetString(edit.RequiredCheck)
		cell = row.AddCell()
		cell.SetString(edit.FutureCheck)
		cell = row.AddCell()
		cell.SetString(edit.RangeCheck)
	}
	err := sheet.SetColWidth(0, 0, float64(urlLength))
	if err != nil {
		fmt.Printf("Error setting the Column: %s", err)
	}
	err = sheet.SetColWidth(1, 1, float64(projectLength))
	if err != nil {
		fmt.Printf("Error setting the Column: %s", err)
	}
	err = sheet.SetColWidth(2, 2, float64(checkLength))
	if err != nil {
		fmt.Printf("Error setting the Column: %s", err)
	}
	err = sheet.SetColWidth(3, 3, float64(formOIDLength))
	if err != nil {
		fmt.Printf("Error setting the Column: %s", err)
	}
	err = sheet.SetColWidth(4, 4, float64(fieldOIDLength))
	if err != nil {
		fmt.Printf("Error setting the Column: %s", err)
	}
	err = sheet.SetColWidth(5, 5, float64(vblOIDLength))
	if err != nil {
		fmt.Printf("Error setting the Column: %s", err)
	}

}

//// write out the data for the LastProjectVersion
func writeLastProjectVersions(urls map[string]*RaveURL, threshold int, wbk *xlsx.File) {
	for url, raveUrl := range urls {
		writeLastProjectVersion(url, threshold, raveUrl.getProjects(), wbk)
	}
}

// write the Subject Counts
func writeLastProjectVersion(url string, threshold int, projects []*Project, wbk *xlsx.File) {

	tabName := fmt.Sprintf("Last - %s", url)
	headers := []string{"Project Name",
		"CRF Version ID",
		"Subject Count",
		"Total Checks",
		"Total Checks (Field)",
		"Total Checks Fired (Field)",
		"Total Checks Not Fired (Field)",
		"Total Checks Open (Field)",
		"%ge Checks Fired (Field)",
		"%ge Checks Not Fired (Field)",
		"Checks with Change (Field)",
		"Checks with No Change (Field)",
		"%ge Checks with Change (Field)",
		"%ge Checks with No Change (Field)",
		"Total Checks (Prog)",
		"Total Checks Fired (Prog)",
		"Total Checks Not Fired (Prog)",
		"Total Checks Open (Prog)",
		"%ge Checks Fired (Prog)",
		"%ge Checks Not Fired (Prog)",
		"Checks with Change (Prog)",
		"Checks with No Change (Prog)",
		"%ge Checks with Change (Prog)",
		"%ge Checks with No Change (Prog)",
	}

	// create the sheet
	sheet, created := getOrAddSheet(wbk, tabName)
	if created {
		// Add the headers
		writeHeaderRow(headers, sheet)
	}
	maxWidth := 70
	projectWidth := 12
	var projectVersions []*Record
	for _, project := range projects {
		lpv := project.getLastVersion().ActiveEditsOnly
		lpv.calculatePercentages()
		projectVersions = append(projectVersions, lpv)
		// Build the summary counts
		var cell *xlsx.Cell
		// Rows
		row := sheet.AddRow()
		// Project Name
		cell = row.AddCell()
		cell.SetString(project.ProjectName)
		if len(project.ProjectName) > projectWidth {
			if len(project.ProjectName) < maxWidth {
				projectWidth = len(project.ProjectName)
			} else {
				projectWidth = maxWidth
			}
		}
		// CRF Version ID
		cell = row.AddCell()
		cell.SetInt(lpv.CRFVersionID)
		// Subject Count
		cell = row.AddCell()
		cell.SetInt(project.SubjectCount)
		// Total Checks
		cell = row.AddCell()
		cell.SetInt(lpv.TotalProgEdits + lpv.TotalFieldEdits)
		// Total Field Checks
		cell = row.AddCell()
		cell.SetInt(lpv.TotalFieldEdits)
		// Total Field Checks Fired
		cell = row.AddCell()
		cell.SetInt(lpv.TotalFieldEditsFired)
		// Total Field Checks Not Fired
		cell = row.AddCell()
		cell.SetInt(lpv.TotalFieldEditsNotFired)
		// Total Field Checks Open
		cell = row.AddCell()
		cell.SetInt(lpv.TotalFieldEditsOpen)
		// Percentage Field Checks Fired
		cell = row.AddCell()
		cell.SetFloatWithFormat(lpv.FieldPercentageFired, "0.00%")
		// Percentage Field Checks Not Fired
		cell = row.AddCell()
		cell.SetFloatWithFormat(lpv.FieldPercentageNotFired, "0.00%")
		// Total Field Checks with Change
		cell = row.AddCell()
		cell.SetInt(lpv.TotalFieldEditsFiredWithChange)
		// Total Field Checks with No Change
		cell = row.AddCell()
		cell.SetInt(lpv.TotalFieldEditsFiredWithNoChange)
		// Percentage Field Checks Leading to Data Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(lpv.FieldPercentageChanged, "0.00%")
		// Percentage Field Checks Leading to No Data Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(lpv.FieldPercentageNotChanged, "0.00%")
		// Total Prog Checks
		cell = row.AddCell()
		cell.SetInt(lpv.TotalProgEdits)
		// Total Prog Checks Fired
		cell = row.AddCell()
		cell.SetInt(lpv.TotalProgEditsFired)
		// Total Prog Checks Not Fired
		cell = row.AddCell()
		cell.SetInt(lpv.TotalProgEditsNotFired)
		// Total Prog Checks Open
		cell = row.AddCell()
		cell.SetInt(lpv.TotalProgEditsOpen)
		// Percentage Prog Checks Fired
		cell = row.AddCell()
		cell.SetFloatWithFormat(lpv.ProgPercentageFired, "0.00%")
		// Percentage Prog Checks Not Fired
		cell = row.AddCell()
		cell.SetFloatWithFormat(lpv.ProgPercentageNotFired, "0.00%")
		// Total Prog Checks with Change
		cell = row.AddCell()
		cell.SetInt(lpv.TotalProgEditsFiredWithChange)
		// Total Prog Checks with No Change
		cell = row.AddCell()
		cell.SetInt(lpv.TotalProgEditsFiredWithNoChange)
		// Percentage Prog Checks Leading to Data Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(lpv.ProgPercentageChanged, "0.00%")
		// Percentage Prog Checks Leading to No Data Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(lpv.ProgPercentageNotChanged, "0.00%")
	}
	err := sheet.SetColWidth(0, 0, float64(projectWidth))
	if err != nil {
		fmt.Printf("Error setting the Project Name Column width: %s", err)
	}

	// write the summary counts
	writeSummaryCounts(threshold, projects, sheet)

}

func writeTotalSummaryCounts(summaryCounts []SummaryCounts, sheet *xlsx.Sheet) {
	// write the totals
	headers := []string{"Threshold",
		"Sample Count",
		"Avg. Subject Count",
		"Total Checks",
		"Total Checks (Field)",
		"Total Checks Fired (Field)",
		"Total Checks Not Fired (Field)",
		"Total Checks Open (Field)",
		"%ge Checks Fired (Field)",
		"%ge Checks Not Fired (Field)",
		"Checks with Change (Field)",
		"Checks with No Change (Field)",
		"%ge Checks with Change (Field)",
		"%ge Checks with No Change (Field)",
		"Total Checks (Prog)",
		"Total Checks Fired (Prog)",
		"Total Checks Not Fired (Prog)",
		"Total Checks Open (Prog)",
		"%ge Checks Fired (Prog)",
		"%ge Checks Not Fired (Prog)",
		"Checks with Change (Prog)",
		"Checks with No Change (Prog)",
		"%ge Checks with Change (Prog)",
		"%ge Checks with No Change (Prog)",
	}
	writeHeaderRow(headers, sheet)
	for _, summary := range summaryCounts {
		// no studies above the threshold
		if summary.RecordCount == 0 {
			continue
		}
		// Add a row
		row := sheet.AddRow()
		// Threshold
		cell := row.AddCell()
		if summary.Threshold > 0 {
			cell.SetString(fmt.Sprintf("> %d", summary.Threshold))
		} else {
			cell.SetString("ALL")
		}
		// Record Count
		cell = row.AddCell()
		cell.SetInt(summary.RecordCount)
		// Average Subject Count
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.SubjectCount)/float64(summary.RecordCount), "0.00")
		// Total Edit Count
		cell = row.AddCell()
		cell.SetInt(summary.TotalEdits)
		// Total Field Edit Count
		cell = row.AddCell()
		cell.SetInt(summary.TotalFldEdits)
		// Total Field Edit Fired Count
		cell = row.AddCell()
		cell.SetInt(summary.TotalFldEditsFired)
		// Total Field Edit Not Fired Count
		cell = row.AddCell()
		cell.SetInt(summary.TotalFldEditsUnfired)
		// Total Field Edit Open Count
		cell = row.AddCell()
		cell.SetInt(summary.TotalFldEditsOpen)
		// Percentage Field Edit Fired Count
		cell = row.AddCell()
		if summary.TotalFldEdits > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalFldEditsFired)/float64(summary.TotalFldEdits),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Percentage Field Edit Not Fired Count
		cell = row.AddCell()
		if summary.TotalFldEdits > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalFldEditsUnfired)/float64(summary.TotalFldEdits),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Field Edit Fired with Change
		cell = row.AddCell()
		cell.SetInt(summary.TotalFldWithChange)
		// Field Edit Fired with No Change
		cell = row.AddCell()
		cell.SetInt(summary.TotalFldWithNoChange)
		// Percentage Field Edit Fired Leading to Change
		cell = row.AddCell()
		if summary.TotalFldEditsFired > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalFldWithChange)/float64(summary.TotalFldEditsFired),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Percentage Field Edit Fired Leading to No Change
		cell = row.AddCell()
		if summary.TotalFldEditsFired > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalFldWithNoChange)/float64(summary.TotalFldEditsFired),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Total Prg Edit Count
		cell = row.AddCell()
		cell.SetInt(summary.TotalPrgEdits)
		// Total Prg Edit Fired Count
		cell = row.AddCell()
		cell.SetInt(summary.TotalPrgEditsFired)
		// Total Prg Edit Not Fired Count
		cell = row.AddCell()
		cell.SetInt(summary.TotalPrgEditsUnfired)
		// Total Prg Edit Open Count
		cell = row.AddCell()
		cell.SetInt(summary.TotalPrgEditsOpen)
		// Percentage Prog Edit Fired
		cell = row.AddCell()
		if summary.TotalPrgEdits > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalPrgEditsFired)/float64(summary.TotalPrgEdits),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Percentage Prog Edit Not Fired
		cell = row.AddCell()
		if summary.TotalPrgEdits > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalPrgEditsUnfired)/float64(summary.TotalPrgEdits),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Prg Edit Fired with Change
		cell = row.AddCell()
		cell.SetInt(summary.TotalPrgWithChange)
		// Prg Edit Fired with No Change
		cell = row.AddCell()
		cell.SetInt(summary.TotalPrgWithNoChange)
		// Percentage Prog Edit Fired Leading to Change
		cell = row.AddCell()
		if summary.TotalPrgEditsFired > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalPrgWithChange)/float64(summary.TotalPrgEditsFired),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Percentage Prog Edit Fired Leading to No Change
		cell = row.AddCell()
		if summary.TotalPrgEditsFired > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalPrgWithNoChange)/float64(summary.TotalPrgEditsFired),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
	}
}

func writeAverageSummaryCounts(summaryCounts []SummaryCounts, sheet *xlsx.Sheet) {
	// write the averages
	headers := []string{"Threshold",
		"Sample Count",
		"Avg. Subject Count",
		"Avg. Total Checks",
		"Avg. Total Checks (Field)",
		"Avg. Total Checks Fired (Field)",
		"Avg. Total Checks Not Fired (Field)",
		"Avg. Total Checks Open (Field)",
		"%ge Checks Fired (Field)",
		"%ge Checks Not Fired (Field)",
		"Avg. Checks with Change (Field)",
		"Avg. Checks with No Change (Field)",
		"%ge Checks with Change (Field)",
		"%ge Checks with No Change (Field)",
		"Avg. Total Checks (Prog)",
		"Avg. Total Checks Fired (Prog)",
		"Avg. Total Checks Not Fired (Prog)",
		"Avg. Total Checks Open (Prog)",
		"%ge Checks Fired (Prog)",
		"%ge Checks Not Fired (Prog)",
		"Avg. Checks with Change (Prog)",
		"Avg. Checks with No Change (Prog)",
		"%ge Checks with Change (Prog)",
		"%ge Checks with No Change (Prog)",
	}
	writeHeaderRow(headers, sheet)
	for _, summary := range summaryCounts {
		// no studies above the threshold
		if summary.RecordCount == 0 {
			continue
		}
		// Add a row
		row := sheet.AddRow()
		var cell *xlsx.Cell
		// Threshold
		cell = row.AddCell()
		if summary.Threshold > 0 {
			cell.SetString(fmt.Sprintf("> %d", summary.Threshold))
		} else {
			cell.SetString("ALL")
		}
		// Record Count
		cell = row.AddCell()
		cell.SetInt(summary.RecordCount)
		// Average Subject Count
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.SubjectCount)/float64(summary.RecordCount), "0.00")
		// Average Total Edit Count
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalEdits)/float64(summary.RecordCount), "0.00")
		// Average Total Field Edit Count
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalFldEdits)/float64(summary.RecordCount), "0.00")
		// Average Total Field Edit Fired Count
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalFldEditsFired)/float64(summary.RecordCount), "0.00")
		// Average Total Field Edit Not Fired Count
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalFldEditsUnfired)/float64(summary.RecordCount), "0.00")
		// Average Total Field Edit Open Count
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalFldEditsOpen)/float64(summary.RecordCount), "0.00")
		// Percentage Field Edit Fired Count
		cell = row.AddCell()
		if summary.TotalFldEdits > 0 {
			// Percentage Field Edit Fired Count
			cell.SetFloatWithFormat(float64(summary.TotalFldEditsFired)/float64(summary.TotalFldEdits),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Percentage Field Edit Not Fired Count
		cell = row.AddCell()
		if summary.TotalFldEdits > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalFldEditsUnfired)/float64(summary.TotalFldEdits),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Average Field Edit Fired with Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalFldWithChange)/float64(summary.RecordCount),
			"0.00")
		// Average Field Edit Fired with No Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalFldWithNoChange)/float64(summary.RecordCount),
			"0.00")
		// Percentage Field Edit Fired Leading to Change
		cell = row.AddCell()
		if summary.TotalFldEditsFired > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalFldWithChange)/float64(summary.TotalFldEditsFired),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Percentage Field Edit Fired Leading to No Change
		cell = row.AddCell()
		if summary.TotalFldEditsFired > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalFldWithNoChange)/float64(summary.TotalFldEditsFired),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Average Total Prg Edit Count
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalPrgEdits)/float64(summary.RecordCount), "0.00")
		// Average Total Prg Edit Fired Count
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalPrgEditsFired)/float64(summary.RecordCount), "0.00")
		// Average Total Prg Edit Not Fired Count
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalPrgEditsUnfired)/float64(summary.RecordCount), "0.00")
		// Average Total Prg Edit Open Count
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalPrgEditsOpen)/float64(summary.RecordCount), "0.00")
		// Percentage Prog Edit Fired
		cell = row.AddCell()
		if summary.TotalPrgEdits > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalPrgEditsFired)/float64(summary.TotalPrgEdits),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Percentage Prog Edit Not Fired
		cell = row.AddCell()
		if summary.TotalPrgEdits > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalPrgEditsUnfired)/float64(summary.TotalPrgEdits),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Average Prg Edit Fired with Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalPrgWithChange)/float64(summary.RecordCount), "0.00")
		// Average Prg Edit Fired with No Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(summary.TotalPrgWithNoChange)/float64(summary.RecordCount), "0.00")
		// Percentage Prog Edit Fired Leading to Change
		cell = row.AddCell()
		if summary.TotalPrgEditsFired > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalPrgWithChange)/float64(summary.TotalPrgEditsFired),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
		// Percentage Prog Edit Fired Leading to No Change
		cell = row.AddCell()
		if summary.TotalPrgEditsFired > 0 {
			cell.SetFloatWithFormat(float64(summary.TotalPrgWithNoChange)/float64(summary.TotalPrgEditsFired),
				"0.00%")
		} else {
			cell.SetInt(0)
		}
	}
}

func writeNotes(sheet *xlsx.Sheet) {
	// Add the notes to the sheet
	row := sheet.AddRow()
	cell := row.AddCell()
	boldface := *xlsx.NewFont(12, "Verdana")
	boldface.Bold = true
	centerHalign := *xlsx.DefaultAlignment()
	centerHalign.Horizontal = "center"
	titleFace := xlsx.NewStyle()
	titleFace.Font = boldface
	titleFace.Alignment = centerHalign
	titleFace.ApplyAlignment = true
	titleFace.ApplyFont = true
	cell.SetStyle(titleFace)
	cell.SetString("Notes")
	notes := []string{"Threshold represents the lower limit on number of participants",
		"Edit counts are restricted to those including CheckAction of OpenQuery"}
	for _, note := range notes {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(note)
	}

}

func writeSummaryCounts(thresholdLevel int, projects []*Project, sheet *xlsx.Sheet) {
	// Count holders
	var summaryCounts []SummaryCounts
	thresholds := []int{thresholdLevel, 0}
	// just those over the threshold
	for i, threshold := range thresholds {
		// add a value for
		summaryCounts = append(summaryCounts, SummaryCounts{})
		for _, project := range projects {
			// filtered set of counts
			if project.SubjectCount > threshold {
				lastProjectVersion := project.getLastVersion().ActiveEditsOnly
				//log.Println("Adding counts for ", last_project_version.ProjectName,"with count",last_project_version.SubjectCount, "with threshold",threshold)
				summaryCounts[i].RecordCount++
				summaryCounts[i].Threshold = threshold
				summaryCounts[i].SubjectCount += project.SubjectCount
				summaryCounts[i].TotalEdits += lastProjectVersion.TotalProgEdits + lastProjectVersion.TotalFieldEdits
				summaryCounts[i].TotalFldEdits += lastProjectVersion.TotalFieldEdits
				summaryCounts[i].TotalFldEditsFired += lastProjectVersion.TotalFieldEditsFired
				summaryCounts[i].TotalFldEditsUnfired += lastProjectVersion.TotalFieldEditsNotFired
				summaryCounts[i].TotalFldEditsOpen += lastProjectVersion.TotalFieldEditsOpen
				summaryCounts[i].TotalFldWithChange += lastProjectVersion.TotalFieldEditsFiredWithChange
				summaryCounts[i].TotalFldWithNoChange += lastProjectVersion.TotalFieldEditsFiredWithNoChange
				summaryCounts[i].TotalPrgEdits += lastProjectVersion.TotalProgEdits
				summaryCounts[i].TotalPrgEditsFired += lastProjectVersion.TotalProgEditsFired
				summaryCounts[i].TotalPrgEditsUnfired += lastProjectVersion.TotalProgEditsNotFired
				summaryCounts[i].TotalPrgEditsOpen += lastProjectVersion.TotalProgEditsOpen
				summaryCounts[i].TotalPrgWithChange += lastProjectVersion.TotalProgEditsFiredWithChange
				summaryCounts[i].TotalPrgWithNoChange += lastProjectVersion.TotalProgEditsFiredWithNoChange
			}
		}
	}
	// write the counts out
	writeTotalSummaryCounts(summaryCounts, sheet)
	writeAverageSummaryCounts(summaryCounts, sheet)
	writeNotes(sheet)
}
