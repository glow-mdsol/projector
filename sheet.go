package main

import (
	"github.com/tealeg/xlsx"
	"log"
	"sort"
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
		sheet.SetColWidth(idx, idx, float64(len(datum)))
	}
}

// Write a Project Version to the Sheet
func writeProjectVersion(prj ProjectVersion, row *xlsx.Row) {
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
	cell.Value = prj.CRFVersionID
	// Last Version?
	cell = row.AddCell()
	if prj.LastVersion {
		cell.SetString("Y")
	} else {
		cell.SetString("N")
	}
	// Current Subject Count
	cell = row.AddCell()
	cell.SetInt(prj.SubjectCount)
	// Active Edits
	writeMetrics(prj.ActiveEditsOnly, row)
	// Inactive Edits
	writeMetrics(prj.InactiveEditsOnly, row)
}

// Write a metric Row
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

// write a set of metrics to the file
func writeMetrics(rec Record, row *xlsx.Row) {
	var cell *xlsx.Cell
	// Total Edits (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEdits)
	// Total Edits Fired (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEditsFired)
	// Total Edits Unfired (fld)
	cell = row.AddCell()
	cell.SetInt(rec.TotalFieldEdits - rec.TotalFieldEditsFired)
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
	// Total Edits Unfired (prg)
	cell = row.AddCell()
	cell.SetInt(rec.TotalProgEditsWithOpenQuery - rec.TotalProgWithOpenQueryFired)
	// %ge Edits Fired (prg)
	cell = row.AddCell()
	if rec.TotalProgEditsWithOpenQuery == 0 {
		cell.SetFloat(0.0)
	} else {
		cell.SetFloatWithFormat(float64(rec.TotalProgWithOpenQueryFired)/float64(rec.TotalProgEditsWithOpenQuery)*100.0, "#,##0.00;(#,##0.00)")
	}
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
}

// write the Subject Counts
func writeSubjectCounts(subject_counts []SubjectCount, wbk *xlsx.File) {
	tab_name := "Subject Counts"
	headers := []string{"Rave URL", "Project Name", "Subject Count", "Date Updated"}

	// create the sheet
	sheet, created := getOrAddSheet(wbk, tab_name)
	if created {
		// Add the headers
		writeHeaderRow(headers, sheet)
	}
	auto_filter := new(xlsx.AutoFilter)
	auto_filter.TopLeftCell = "A1"
	auto_filter.BottomRightCell = "D1"
	sheet.AutoFilter = auto_filter

	// default widths
	url_length := 12
	project_length := 12

	for _, subject_count := range subject_counts {
		var cell *xlsx.Cell
		// Rows
		row := sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(subject_count.URL)
		cell = row.AddCell()
		cell.SetString(subject_count.ProjectName)
		cell = row.AddCell()
		if subject_count.SubjectCount.Valid{

			cell.SetInt(int(subject_count.SubjectCount.Int64))
		} else {
			cell.SetString("-")
		}
		cell = row.AddCell()
		if subject_count.RefreshDate.Valid {
			cell.SetDateTime(subject_count.RefreshDate.Time)
		} else {
			cell.SetString("-")
		}
		if len(subject_count.URL) > url_length {
			url_length = len(subject_count.URL)
		}
		if len(subject_count.ProjectName) > project_length {
			project_length = len(subject_count.ProjectName)
		}

	}
	sheet.SetColWidth(0, 0, float64(url_length))
	sheet.SetColWidth(1, 1, float64(project_length))
}

// write study metrics
func writeStudyMetrics(data map[string][]ProjectVersion, wbk *xlsx.File) {
	headers := []string{"Study URL",
						"Project Name",
						"CRF Version",
						"Last Version",
						"Subject Count",
						"Active - Total Edits (fld)",
						"Active - Total Edits Fired (fld)",
						"Active - Total Edits Unfired (fld)",
						"Active - %ge Edits Fired (fld)",
						"Active - %ge Edits Unfired (fld)",
						"Active - Total Edits (prg)",
						"Active - Total Edits With OpenQuery (prg)",
						"Active - Total Edits Fired (prg)",
						"Active - Total Edits Unfired (prg)",
						"Active - %ge Edits Fired (prg)",
						"Active - %ge Edits Unfired (prg)",
						"Active - Total Queries (fld)",
						"Active - Total Queries (prg)",
						"Active - Total Queries With OpenQuery (prg)",
						"Inactive - Total Edits (fld)",
						"Inactive - Total Edits Fired (fld)",
						"Inactive - Total Edits Unfired (fld)",
						"Inactive - %ge Edits Fired (fld)",
						"Inactive - %ge Edits Unfired (fld)",
						"Inactive - Total Edits (prg)",
						"Inactive - Total Edits With OpenQuery (prg)",
						"Inactive - Total Edits Fired (prg)",
						"Inactive - Total Edits Unfired (prg)",
						"Inactive - %ge Edits Fired (prg)",
						"Inactive - %ge Edits Unfired (prg)",
						"Inactive - Total Queries (fld)",
						"Inactive - Total Queries (prg)",
						"Inactive - Total Queries With OpenQuery (prg)",
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
		auto_filter := new(xlsx.AutoFilter)
		auto_filter.TopLeftCell = "A1"
		auto_filter.BottomRightCell = "AG1"
		sheet.AutoFilter = auto_filter

		//log.Println("Created Sheet for URL ", url)

		for _, project_version := range data[url] {
			// Add the row for Checks
			row := sheet.AddRow()
			writeProjectVersion(project_version, row)
		}
	}
}

// write study metrics
func writeStudyMetricsPartitioned(data map[string][]ProjectVersion, wbk *xlsx.File) {
	headers := []string{"Study URL",
						"Project Name",
						"CRF Version",
						"Status",
						"Total Edits (fld)",
						"Total Edits Fired (fld)",
						"Total Edits Unfired (fld)",
						"%ge Edits Fired (fld)",
						"%ge Edits Unfired (fld)",
						"Total Edits (prg)",
						"Total Edits With OpenQuery (prg)",
						"Total Edits Fired (prg)",
						"Total Edits Unfired (prg)",
						"%ge Edits Fired (prg)",
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

func writeUselessEdits(edits []UnusedEdit, wbk *xlsx.File) {
	headers := []string{"Study URL",
						"Project Name",
						"Edit Check Name",
						"Times Used",
						"OpenQuery Check?"}

	tab_name := "Unused Edits"

	// create the sheet
	sheet, created := getOrAddSheet(wbk, tab_name)
	if created {
		// Add the headers
		writeHeaderRow(headers, sheet)
	}
	auto_filter := new(xlsx.AutoFilter)
	auto_filter.TopLeftCell = "A1"
	auto_filter.BottomRightCell = "E1"
	sheet.AutoFilter = auto_filter

	url_length := 12
	project_length := 12
	check_length := 12

	// Export the results
	for _, edit := range edits {
		if len(edit.URL) > url_length {
			url_length = len(edit.URL)
		}
		if len(edit.ProjectName) > project_length {
			project_length = len(edit.ProjectName)
		}
		if len(edit.EditCheckName) > check_length {
			check_length = len(edit.EditCheckName)
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
		cell.SetInt(edit.UsageCount)
		cell = row.AddCell()
		cell.SetString(edit.OpenQuery)
	}
	sheet.SetColWidth(0, 0, float64(url_length))
	sheet.SetColWidth(1, 1, float64(project_length))
	sheet.SetColWidth(2, 2, float64(check_length))
}
