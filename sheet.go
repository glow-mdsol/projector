package main

import (
	"fmt"
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
		if subject_count.SubjectCount.Valid {

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

func writeUselessEdits(edits []UnusedEdit, wbk *xlsx.File) {
	headers := []string{"Study URL",
		"Project Name",
		"Edit Check Name",
		"Times Used",
		"OpenQuery Check?",
		"Custom Function?",
		"Non-conformance check?",
		"Required check?",
		"Future check?",
		"Range check?",
	}

	tab_name := "Unused Edits"

	// create the sheet
	sheet, created := getOrAddSheet(wbk, tab_name)
	if created {
		// Add the headers
		writeHeaderRow(headers, sheet)
	}
	auto_filter := new(xlsx.AutoFilter)
	auto_filter.TopLeftCell = "A1"
	auto_filter.BottomRightCell = "J1"
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
	sheet.SetColWidth(0, 0, float64(url_length))
	sheet.SetColWidth(1, 1, float64(project_length))
	sheet.SetColWidth(2, 2, float64(check_length))
}

// write out the data for the LastProjectVersion
func writeLastProjectVersions(lpv map[string][]LastProjectVersion, threshold int, wbk *xlsx.File) {
	for url, lpvs := range lpv {
		if len(lpvs) > 0 {
			writeLastProjectVersion(url, threshold, lpvs, wbk)
		}
	}
}

// write the Subject Counts
func writeLastProjectVersion(url string, threshold int, projectVersions []LastProjectVersion, wbk *xlsx.File) {

	tabName := fmt.Sprintf("Last - %s", url)
	headers := []string{"Project Name", "CRF Version ID", "Subject Count",
		"Total Checks",
		"Total Checks (Field)",
		"Total Checks Fired (Field)", "Total Checks Not Fired (Field)",
		"%ge Checks Fired (Field)", "%ge Checks Not Fired (Field)",
		"Checks with Change (Field)", "Checks with No Change (Field)",
		"%ge Checks with Change (Field)", "%ge Checks with No Change (Field)",
		"Total Checks (Prog)",
		"Total Checks Fired (Prog)", "Total Checks Not Fired (Prog)",
		"%ge Checks Fired (Prog)", "%ge Checks Not Fired (Prog)",
		"Checks with Change (Prog)", "Checks with No Change (Prog)",
		"%ge Checks with Change (Prog)", "%ge Checks with No Change (Prog)",
	}

	// create the sheet
	sheet, created := getOrAddSheet(wbk, tabName)
	if created {
		// Add the headers
		writeHeaderRow(headers, sheet)
	}

	for _, lastProjectVersion := range projectVersions {
		// Build the summary counts
		var cell *xlsx.Cell
		// Rows
		row := sheet.AddRow()
		// Project Name
		cell = row.AddCell()
		cell.SetString(lastProjectVersion.ProjectName)
		// CRF Version ID
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.CRFVersionID)
		// Subject Count
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.SubjectCount)
		// Total Checks
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.TotalCount)
		// Total Field Checks
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.FieldTotal)
		// Total Field Checks Fired
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.FieldTotalFired)
		// Total Field Checks Not Fired
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.FieldTotalNotFired)
		// Percentage Field Checks Fired
		cell = row.AddCell()
		cell.SetFloatWithFormat(lastProjectVersion.FieldPercentageFired, "0.00%")
		// Percentage Field Checks Not Fired
		cell = row.AddCell()
		cell.SetFloatWithFormat(lastProjectVersion.FieldPercentageNotFired, "0.00%")
		// Total Field Checks with Change
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.FieldChanged)
		// Total Field Checks with No Change
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.FieldNotChanged)
		// Percentage Field Checks Leading to Data Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(lastProjectVersion.FieldPercentageChanged, "0.00%")
		// Percentage Field Checks Leading to No Data Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(lastProjectVersion.FieldPercentageNotChanged, "0.00%")
		// Total Prog Checks
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.ProgTotal)
		// Total Prog Checks Fired
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.ProgTotalFired)
		// Total Prog Checks Not Fired
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.ProgTotalNotFired)
		// Percentage Prog Checks Fired
		cell = row.AddCell()
		cell.SetFloatWithFormat(lastProjectVersion.ProgPercentageFired, "0.00%")
		// Percentage Prog Checks Not Fired
		cell = row.AddCell()
		cell.SetFloatWithFormat(lastProjectVersion.ProgPercentageNotFired, "0.00%")
		// Total Prog Checks with Change
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.ProgChanged)
		// Total Prog Checks with No Change
		cell = row.AddCell()
		cell.SetInt(lastProjectVersion.ProgNotChanged)
		// Percentage Prog Checks Leading to Data Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(lastProjectVersion.ProgPercentageChanged, "0.00%")
		// Percentage Prog Checks Leading to No Data Change
		cell = row.AddCell()
		cell.SetFloatWithFormat(lastProjectVersion.ProgPercentageNotChanged, "0.00%")
	}
	// write the summary counts
	writeSummaryCounts(threshold, projectVersions, sheet)

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
		"%ge Checks Fired (Field)",
		"%ge Checks Not Fired (Field)",
		"Checks with Change (Field)",
		"Checks with No Change (Field)",
		"%ge Checks with Change (Field)",
		"%ge Checks with No Change (Field)",
		"Total Checks (Prog)",
		"Total Checks Fired (Prog)",
		"Total Checks Not Fired (Prog)",
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
		for i := 0; i < 22; i++ {
			cell := row.AddCell()
			if i == 0 {
				if summary.Threshold > 0 {
					cell.SetString(fmt.Sprintf("> %d", summary.Threshold))
				} else {
					cell.SetString("ALL")
				}
			} else if i == 1 {
				cell.SetInt(summary.RecordCount)
			} else if i == 2 {
				// Average Subject Count
				cell.SetFloatWithFormat(float64(summary.SubjectCount)/float64(summary.RecordCount), "0.00")
			} else if i == 3 {
				// Total Edit Count
				cell.SetInt(summary.TotalEdits)
			} else if i == 4 {
				// Total Field Edit Count
				cell.SetInt(summary.TotalFldEdits)
			} else if i == 5 {
				// Total Field Edit Fired Count
				cell.SetInt(summary.TotalFldEditsFired)
			} else if i == 6 {
				// Total Field Edit Not Fired Count
				cell.SetInt(summary.TotalFldEditsUnfired)
			} else if i == 7 {
				if summary.TotalFldEdits > 0 {
					// Percentage Field Edit Fired Count
					cell.SetFloatWithFormat(float64(summary.TotalFldEditsFired)/float64(summary.TotalFldEdits),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 8 {
				if summary.TotalFldEdits > 0 {
					// Percentage Field Edit Not Fired Count
					cell.SetFloatWithFormat(float64(summary.TotalFldEditsUnfired)/float64(summary.TotalFldEdits),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 9 {
				// Field Edit Fired with Change
				cell.SetInt(summary.TotalFldWithChange)
			} else if i == 10 {
				// Field Edit Fired with No Change
				cell.SetInt(summary.TotalFldWithNoChange)
			} else if i == 11 {
				if summary.TotalFldEditsFired > 0 {
					// Percentage Field Edit Fired Leading to Change
					cell.SetFloatWithFormat(float64(summary.TotalFldWithChange)/float64(summary.TotalFldEditsFired),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 12 {
				if summary.TotalFldEditsFired > 0 {
					// Percentage Field Edit Fired Leading to No Change
					cell.SetFloatWithFormat(float64(summary.TotalFldWithNoChange)/float64(summary.TotalFldEditsFired),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 13 {
				// Total Prg Edit Count
				cell.SetInt(summary.TotalPrgEdits)
			} else if i == 14 {
				// Total Prg Edit Fired Count
				cell.SetInt(summary.TotalPrgEditsFired)
			} else if i == 15 {
				// Total Prg Edit Not Fired Count
				cell.SetInt(summary.TotalPrgEditsUnfired)
			} else if i == 16 {
				if summary.TotalPrgEdits > 0 {
					// Percentage Prog Edit Fired
					cell.SetFloatWithFormat(float64(summary.TotalPrgEditsFired)/float64(summary.TotalPrgEdits),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 17 {
				if summary.TotalPrgEdits > 0 {
					// Percentage Prog Edit Not Fired
					cell.SetFloatWithFormat(float64(summary.TotalPrgEditsUnfired)/float64(summary.TotalPrgEdits),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 18 {
				// Prg Edit Fired with Change
				cell.SetInt(summary.TotalPrgWithChange)
			} else if i == 19 {
				// Prg Edit Fired with No Change
				cell.SetInt(summary.TotalPrgWithNoChange)
			} else if i == 20 {
				if summary.TotalPrgEditsFired > 0 {
					// Percentage Prog Edit Fired Leading to Change
					cell.SetFloatWithFormat(float64(summary.TotalPrgWithChange)/float64(summary.TotalPrgEditsFired),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 21 {
				if summary.TotalPrgEditsFired > 0 {
					// Percentage Prog Edit Fired Leading to No Change
					cell.SetFloatWithFormat(float64(summary.TotalPrgWithNoChange)/float64(summary.TotalPrgEditsFired),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else {
				cell.SetString("")
			}
		}
	}
}

func writeAverageSummaryCounts(summaryCounts []SummaryCounts, sheet *xlsx.Sheet) {
	// write the averages
	headers := []string{"Threshold", "Sample Count", "Avg. Subject Count",
		"Avg. Total Checks",
		"Avg. Total Checks (Field)",
		"Avg. Total Checks Fired (Field)",
		"Avg. Total Checks Not Fired (Field)",
		"%ge Checks Fired (Field)",
		"%ge Checks Not Fired (Field)",
		"Avg. Checks with Change (Field)",
		"Avg. Checks with No Change (Field)",
		"%ge Checks with Change (Field)",
		"%ge Checks with No Change (Field)",
		"Avg. Total Checks (Prog)",
		"Avg. Total Checks Fired (Prog)",
		"Avg. Total Checks Not Fired (Prog)",
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
		for i := 0; i < 22; i++ {
			cell := row.AddCell()
			if i == 0 {
				if summary.Threshold > 0 {
					cell.SetString(fmt.Sprintf("> %d", summary.Threshold))
				} else {
					cell.SetString("ALL")
				}
			} else if i == 1 {
				cell.SetInt(summary.RecordCount)
			} else if i == 2 {
				// Average Subject Count
				cell.SetFloatWithFormat(float64(summary.SubjectCount)/float64(summary.RecordCount), "0.00")
			} else if i == 3 {
				// Average Total Edit Count
				cell.SetFloatWithFormat(float64(summary.TotalEdits)/float64(summary.RecordCount), "0.00")
			} else if i == 4 {
				// Average Total Field Edit Count
				cell.SetFloatWithFormat(float64(summary.TotalFldEdits)/float64(summary.RecordCount), "0.00")
			} else if i == 5 {
				// Average Total Field Edit Fired Count
				cell.SetFloatWithFormat(float64(summary.TotalFldEditsFired)/float64(summary.RecordCount), "0.00")
			} else if i == 6 {
				// Average Total Field Edit Not Fired Count
				cell.SetFloatWithFormat(float64(summary.TotalFldEditsUnfired)/float64(summary.RecordCount), "0.00")
			} else if i == 7 {
				if summary.TotalFldEdits > 0 {
					// Percentage Field Edit Fired Count
					cell.SetFloatWithFormat(float64(summary.TotalFldEditsFired)/float64(summary.TotalFldEdits),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 8 {
				if summary.TotalFldEdits > 0 {
					// Percentage Field Edit Not Fired Count
					cell.SetFloatWithFormat(float64(summary.TotalFldEditsUnfired)/float64(summary.TotalFldEdits),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 9 {
				// Average Field Edit Fired with Change
				cell.SetFloatWithFormat(float64(summary.TotalFldWithChange)/float64(summary.RecordCount),
					"0.00")
			} else if i == 10 {
				// Average Field Edit Fired with No Change
				cell.SetFloatWithFormat(float64(summary.TotalFldWithNoChange)/float64(summary.RecordCount),
					"0.00")
			} else if i == 11 {
				if summary.TotalFldEditsFired > 0 {
					// Percentage Field Edit Fired Leading to Change
					cell.SetFloatWithFormat(float64(summary.TotalFldWithChange)/float64(summary.TotalFldEditsFired),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 12 {
				if summary.TotalFldEditsFired > 0 {
					// Percentage Field Edit Fired Leading to No Change
					cell.SetFloatWithFormat(float64(summary.TotalFldWithNoChange)/float64(summary.TotalFldEditsFired),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 13 {
				// Average Total Prg Edit Count
				cell.SetFloatWithFormat(float64(summary.TotalPrgEdits)/float64(summary.RecordCount), "0.00")
			} else if i == 14 {
				// Average Total Prg Edit Fired Count
				cell.SetFloatWithFormat(float64(summary.TotalPrgEditsFired)/float64(summary.RecordCount), "0.00")
			} else if i == 15 {
				// Average Total Prg Edit Not Fired Count
				cell.SetFloatWithFormat(float64(summary.TotalPrgEditsUnfired)/float64(summary.RecordCount), "0.00")
			} else if i == 16 {
				if summary.TotalPrgEdits > 0 {
					// Percentage Prog Edit Fired
					cell.SetFloatWithFormat(float64(summary.TotalPrgEditsFired)/float64(summary.TotalPrgEdits),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 17 {
				if summary.TotalPrgEdits > 0 {
					// Percentage Prog Edit Not Fired
					cell.SetFloatWithFormat(float64(summary.TotalPrgEditsUnfired)/float64(summary.TotalPrgEdits),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 18 {
				// Average Prg Edit Fired with Change
				cell.SetFloatWithFormat(float64(summary.TotalPrgWithChange)/float64(summary.RecordCount), "0.00")
			} else if i == 19 {
				// Average Prg Edit Fired with No Change
				cell.SetFloatWithFormat(float64(summary.TotalPrgWithNoChange)/float64(summary.RecordCount), "0.00")
			} else if i == 20 {
				if summary.TotalPrgEditsFired > 0 {
					// Percentage Prog Edit Fired Leading to Change
					cell.SetFloatWithFormat(float64(summary.TotalPrgWithChange)/float64(summary.TotalPrgEditsFired),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else if i == 21 {
				if summary.TotalPrgEditsFired > 0 {
					// Percentage Prog Edit Fired Leading to No Change
					cell.SetFloatWithFormat(float64(summary.TotalPrgWithNoChange)/float64(summary.TotalPrgEditsFired),
						"0.00%")
				} else {
					cell.SetInt(0)
				}
			} else {
				cell.SetString("")
			}
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

func writeSummaryCounts(thresholdLevel int, projectVersions []LastProjectVersion, sheet *xlsx.Sheet) {
	// No matching versions
	if len(projectVersions) == 0 {
		return
	}
	// Count holders
	var summaryCounts []SummaryCounts
	thresholds := []int{thresholdLevel, 0}
	// just those over the threshold
	for i, threshold := range thresholds {
		// add a value for
		summaryCounts = append(summaryCounts, SummaryCounts{})
		for _, lastProjectVersion := range projectVersions {
			// filtered set of counts
			if lastProjectVersion.SubjectCount > threshold {
				//log.Println("Adding counts for ", last_project_version.ProjectName,"with count",last_project_version.SubjectCount, "with threshold",threshold)
				summaryCounts[i].RecordCount += 1
				summaryCounts[i].Threshold = threshold
				summaryCounts[i].SubjectCount += lastProjectVersion.SubjectCount
				summaryCounts[i].TotalEdits += lastProjectVersion.TotalCount
				summaryCounts[i].TotalFldEdits += lastProjectVersion.FieldTotal
				summaryCounts[i].TotalFldEditsFired += lastProjectVersion.FieldTotalFired
				summaryCounts[i].TotalFldEditsUnfired += lastProjectVersion.FieldTotalNotFired
				summaryCounts[i].TotalFldWithChange += lastProjectVersion.FieldChanged
				summaryCounts[i].TotalFldWithNoChange += lastProjectVersion.FieldNotChanged
				summaryCounts[i].TotalPrgEdits += lastProjectVersion.ProgTotal
				summaryCounts[i].TotalPrgEditsFired += lastProjectVersion.ProgTotalFired
				summaryCounts[i].TotalPrgEditsUnfired += lastProjectVersion.ProgTotalNotFired
				summaryCounts[i].TotalPrgWithChange += lastProjectVersion.ProgChanged
				summaryCounts[i].TotalPrgWithNoChange += lastProjectVersion.ProgNotChanged
			}
		}
	}
	// write the counts out
	writeTotalSummaryCounts(summaryCounts, sheet)
	writeAverageSummaryCounts(summaryCounts, sheet)
	writeNotes(sheet)
}
