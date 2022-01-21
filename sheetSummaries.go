package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

// Write the aggregated averages, broken down by the threshold
func writeAggregatedCounts(agg AggregateCount, sheet *xlsx.Sheet) {
	// write the averages
	headers := []string{"Criteria",
		"Aggregate",
		"Threshold",
		"Sample Count",
		"Subject Count",
		"Total Checks",
		"Total Checks (fld)",
		"Total Checks Fired (fld)",
		"Total Checks Not Fired (fld)",
		"Total Checks Open (fld)",
		"%ge Checks Fired (fld)",
		"%ge Checks Not Fired (fld)",
		"Checks with Change (fld)",
		"Checks with No Change (fld)",
		"%ge Checks with Change (fld)",
		"%ge Checks with No Change (fld)",
		"Total Checks (prg)",
		"Total Checks Fired (prg)",
		"Total Checks Not Fired (prg)",
		"Total Checks Open (prg)",
		"%ge Checks Fired (prg)",
		"%ge Checks Not Fired (prg)",
		"Checks with Change (prg)",
		"Checks with No Change (prg)",
		"%ge Checks with Change (prg)",
		"%ge Checks with No Change (prg)",
	}
	writeHeaderRow(headers, sheet)
	var summary SummaryCounts
	// All Projects
	summary = agg.AllProjects
	// no studies above the threshold
	writeAggregates("All Projects", sheet, summary)
	// Greater than 10 subjects
	summary = agg.GreaterThanTen
	// no studies above the threshold
	writeAggregates("Subject Count", sheet, summary)
	// Completed Subjects
	summary = agg.CompletedSubjects
	// no studies above the threshold
	writeAggregates("Completed Subjects", sheet, summary)
}

func writeAggregates(description string, sheet *xlsx.Sheet, summary SummaryCounts) {
	var cell *xlsx.Cell
	// check if there are any records
	if summary.RecordCount > 0 {
		// Add a row
		row := sheet.AddRow()
		// Criteria
		cell = row.AddCell()
		cell.SetString(description)
		// Aggregation => Sum
		cell = row.AddCell()
		cell.SetString("Sum")
		// Threshold
		cell = row.AddCell()
		if summary.Threshold > 0 {
			cell.SetString(fmt.Sprintf("> %d", summary.Threshold))
		} else {
			cell.SetString("ALL")
		}
		writeSumSummaryCounts(row, summary)
	}
	avg := summary.getAverageCounts()
	// check if there are any records
	if avg.RecordCount > 0 {
		// Aggregation => Average
		// Add a row
		row := sheet.AddRow()
		// Criteria
		cell = row.AddCell()
		cell.SetString(description)
		// Aggregation
		cell = row.AddCell()
		cell.SetString("Average")
		// Threshold
		cell = row.AddCell()
		if summary.Threshold > 0 {
			cell.SetString(fmt.Sprintf("> %d", summary.Threshold))
		} else {
			cell.SetString("ALL")
		}

		writeAvgSummaryCounts(row, summary.getAverageCounts())
	}
}

func writeSumSummaryCounts(row *xlsx.Row, summary SummaryCounts) {
	var cell *xlsx.Cell
	// Record Count
	cell = row.AddCell()
	cell.SetInt(summary.RecordCount)
	// Sum Subject Count
	cell = row.AddCell()
	cell.SetInt(summary.SubjectCount)
	//  Total Edit Count
	cell = row.AddCell()
	cell.SetInt(summary.TotalEdits)
	//  Total Field Edit Count
	cell = row.AddCell()
	cell.SetInt(summary.TotalFldEdits)
	//  Total Field Edit Fired Count
	cell = row.AddCell()
	cell.SetInt(summary.TotalFldEditsFired)
	//  Total Field Edit Not Fired Count
	cell = row.AddCell()
	cell.SetInt(summary.TotalFldEditsUnfired)
	// Total Field Edit Open Count
	cell = row.AddCell()
	cell.SetInt(summary.TotalFldEditsOpen)
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
	// Field Edit Fired with Change
	cell = row.AddCell()
	cell.SetInt(summary.TotalFldWithChange)
	// Field Edit Fired with No Change
	cell = row.AddCell()
	cell.SetInt(summary.TotalFldWithNoChange)
	// Percentage Field Edit Fired Leading to Change
	cell = row.AddCell()
	if summary.TotalFldEditsFired > 0 {
		cell.SetFloatWithFormat(float64(summary.TotalFldWithChange)/(float64(summary.TotalFldWithChange)+float64(summary.TotalFldWithNoChange)),
			"0.00%")
	} else {
		cell.SetInt(0)
	}
	// Percentage Field Edit Fired Leading to No Change
	cell = row.AddCell()
	if summary.TotalFldEditsFired > 0 {
		cell.SetFloatWithFormat(float64(summary.TotalFldWithNoChange)/(float64(summary.TotalFldWithChange)+float64(summary.TotalFldWithNoChange)),
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
	//  Total Prg Edit Not Fired Count
	cell = row.AddCell()
	cell.SetInt(summary.TotalPrgEditsUnfired)
	//  Total Prg Edit Open Count
	cell = row.AddCell()
	cell.SetInt(summary.TotalPrgEditsOpen)
	// Percentage Prog Edit Fired
	cell = row.AddCell()
	if summary.TotalPrgEdits > 0 {
		cell.SetFloatWithFormat(float64(summary.TotalPrgEditsFired)/float64(summary.TotalPrgEditsWithOpenQuery),
			"0.00%")
	} else {
		cell.SetInt(0)
	}
	// Percentage Prog Edit Not Fired
	cell = row.AddCell()
	if summary.TotalPrgEdits > 0 {
		cell.SetFloatWithFormat(float64(summary.TotalPrgEditsUnfired)/float64(summary.TotalPrgEditsWithOpenQuery),
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
		cell.SetFloatWithFormat(float64(summary.TotalPrgWithChange)/(float64(summary.TotalPrgWithChange)+float64(summary.TotalPrgWithNoChange)),
			"0.00%")
	} else {
		cell.SetInt(0)
	}
	// Percentage Prog Edit Fired Leading to No Change
	cell = row.AddCell()
	if summary.TotalPrgEditsFired > 0 {
		cell.SetFloatWithFormat(float64(summary.TotalPrgWithNoChange)/(float64(summary.TotalPrgWithChange)+float64(summary.TotalPrgWithNoChange)),
			"0.00%")
	} else {
		cell.SetInt(0)
	}
}

func writeAvgSummaryCounts(row *xlsx.Row, summary AverageSummaryCounts) {
	var cell *xlsx.Cell
	// Record Count
	cell = row.AddCell()
	cell.SetInt(summary.RecordCount)
	// Average Subject Count
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.SubjectCount, "0.00")
	// Average Total Edit Count
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalEdits, "0.00")
	// Average Total Field Edit Count
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalFldEdits, "0.00")
	// Average Total Field Edit Fired Count
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalFldEditsFired, "0.00")
	// Average Total Field Edit Not Fired Count
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalFldEditsUnfired, "0.00")
	// Average Total Field Edit Open Count
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalFldEditsOpen, "0.00")
	// Percentage Field Edit Fired Count
	cell = row.AddCell()
	if summary.TotalFldEdits > 0.0 {
		// Percentage Field Edit Fired Count
		cell.SetFloatWithFormat(summary.TotalFldEditsFired/summary.TotalFldEdits,
			"0.00%")
	} else {
		cell.SetInt(0)
	}
	// Percentage Field Edit Not Fired Count
	cell = row.AddCell()
	if summary.TotalFldEdits > 0 {
		cell.SetFloatWithFormat(summary.TotalFldEditsUnfired/summary.TotalFldEdits,
			"0.00%")
	} else {
		cell.SetInt(0)
	}
	// Average Field Edit Fired with Change
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalFldWithChange,
		"0.00")
	// Average Field Edit Fired with No Change
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalFldWithNoChange,
		"0.00")
	// Percentage Field Edit Fired Leading to Change
	cell = row.AddCell()
	if summary.TotalFldEditsFired > 0.0 {
		cell.SetFloatWithFormat(summary.TotalFldWithChange/summary.TotalFldEditsFired,
			"0.00%")
	} else {
		cell.SetInt(0)
	}
	// Percentage Field Edit Fired Leading to No Change
	cell = row.AddCell()
	if summary.TotalFldEditsFired > 0.0 {
		cell.SetFloatWithFormat(summary.TotalFldWithNoChange/summary.TotalFldEditsFired,
			"0.00%")
	} else {
		cell.SetInt(0)
	}
	// Average Total Prg Edit Count
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalPrgEdits, "0.00")
	// Average Total Prg Edit Fired Count
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalPrgEditsFired, "0.00")
	// Average Total Prg Edit Not Fired Count
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalPrgEditsUnfired, "0.00")
	// Average Total Prg Edit Open Count
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalPrgEditsOpen, "0.00")
	// Percentage Prog Edit Fired
	cell = row.AddCell()
	if summary.TotalPrgEditsWithOpenQuery > 0.0 {
		cell.SetFloatWithFormat(summary.TotalPrgEditsFired/summary.TotalPrgEditsWithOpenQuery,
			"0.00%")
	} else {
		cell.SetInt(0)
	}
	// Percentage Prog Edit Not Fired
	cell = row.AddCell()
	if summary.TotalPrgEditsWithOpenQuery > 0.0 {
		cell.SetFloatWithFormat(summary.TotalPrgEditsUnfired/summary.TotalPrgEditsWithOpenQuery,
			"0.00%")
	} else {
		cell.SetInt(0)
	}
	// Average Prg Edit Fired with Change
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalPrgWithChange, "0.00")
	// Average Prg Edit Fired with No Change
	cell = row.AddCell()
	cell.SetFloatWithFormat(summary.TotalPrgWithNoChange, "0.00")
	// Percentage Prog Edit Fired Leading to Change
	cell = row.AddCell()
	if summary.TotalPrgEditsFired > 0.0 {
		cell.SetFloatWithFormat(summary.TotalPrgWithChange/(summary.TotalPrgWithChange+summary.TotalPrgWithNoChange),
			"0.00%")
	} else {
		cell.SetInt(0)
	}
	// Percentage Prog Edit Fired Leading to No Change
	cell = row.AddCell()
	if summary.TotalPrgEditsFired > 0.0 {
		cell.SetFloatWithFormat(summary.TotalPrgWithNoChange/(summary.TotalPrgWithChange+summary.TotalPrgWithNoChange),
			"0.00%")
	} else {
		cell.SetInt(0)
	}
}

//func writeNotes(sheet *xlsx.Sheet) {
//	// Add the notes to the sheet
//	row := sheet.AddRow()
//	cell := row.AddCell()
//	boldface := *xlsx.NewFont(12, "Verdana")
//	boldface.Bold = true
//	centerHalign := *xlsx.DefaultAlignment()
//	centerHalign.Horizontal = "center"
//	titleFace := xlsx.NewStyle()
//	titleFace.Font = boldface
//	titleFace.Alignment = centerHalign
//	titleFace.ApplyAlignment = true
//	titleFace.ApplyFont = true
//	cell.SetStyle(titleFace)
//	cell.SetString("Notes")
//	notes := []string{"Threshold represents the lower limit on number of participants",
//		"Edit counts are restricted to those including CheckAction of OpenQuery"}
//	for _, note := range notes {
//		row = sheet.AddRow()
//		cell = row.AddCell()
//		cell.SetString(note)
//	}
//
//}

// Write the summary counts (Average and Sum) for a Last Project Version Sheet
func writeSummaryCounts(projects []*Project, wbk *xlsx.File) {
	const subjectThreshold = 10
	const completedCount = 1
	// Count holders
	var aggregateCount AggregateCount
	// initiate values
	aggregateCount.GreaterThanTen.RecordCount = 0
	aggregateCount.AllProjects.RecordCount = 0
	aggregateCount.CompletedSubjects.RecordCount = 0

	// Scan the projects
	for _, project := range projects {
		lastProjectVersion := project.getLastVersion()
		// filtered set of counts
		if project.SubjectCount.SubjectCount > subjectThreshold {
			//log.Println("Adding counts for ", last_project_version.ProjectName,"with count",last_project_version.SubjectCount, "with threshold",threshold)
			aggregateCount.GreaterThanTen.RecordCount++
			aggregateCount.GreaterThanTen.Threshold = subjectThreshold
			aggregateCount.GreaterThanTen.SubjectCount += project.SubjectCount.SubjectCount
			aggregateCount.GreaterThanTen.TotalEdits += lastProjectVersion.getTotalEdits()
			aggregateCount.GreaterThanTen.TotalFldEdits += lastProjectVersion.FieldEditMetrics.TotalEdits
			aggregateCount.GreaterThanTen.TotalFldEditsFired += lastProjectVersion.FieldEditMetrics.TotalFiredWithOpenQuery
			aggregateCount.GreaterThanTen.TotalFldEditsUnfired += lastProjectVersion.FieldEditMetrics.TotalNotFiredWithOpenQuery
			aggregateCount.GreaterThanTen.TotalFldEditsOpen += lastProjectVersion.FieldEditMetrics.TotalOpenQueries
			aggregateCount.GreaterThanTen.TotalFldWithChange += lastProjectVersion.FieldEditMetrics.TotalEditsFiredWithChange
			aggregateCount.GreaterThanTen.TotalFldWithNoChange += lastProjectVersion.FieldEditMetrics.TotalEditsFiredWithNoChange
			aggregateCount.GreaterThanTen.TotalPrgEdits += lastProjectVersion.ProgramEditMetrics.TotalEdits
			aggregateCount.GreaterThanTen.TotalPrgEditsWithOpenQuery += lastProjectVersion.ProgramEditMetrics.TotalEditsWithOpenQuery
			aggregateCount.GreaterThanTen.TotalPrgEditsFired += lastProjectVersion.ProgramEditMetrics.TotalFiredWithOpenQuery
			aggregateCount.GreaterThanTen.TotalPrgEditsUnfired += lastProjectVersion.ProgramEditMetrics.TotalNotFiredWithOpenQuery
			aggregateCount.GreaterThanTen.TotalPrgEditsOpen += lastProjectVersion.ProgramEditMetrics.TotalOpenQueries
			aggregateCount.GreaterThanTen.TotalPrgWithChange += lastProjectVersion.ProgramEditMetrics.TotalEditsFiredWithChange
			aggregateCount.GreaterThanTen.TotalPrgWithNoChange += lastProjectVersion.ProgramEditMetrics.TotalEditsFiredWithNoChange
		}
		// Check for completedSubjects
		if project.SubjectCount.CompletedCount.Valid {
			if project.SubjectCount.CompletedCount.Int64 > completedCount {
				aggregateCount.CompletedSubjects.RecordCount++
				aggregateCount.CompletedSubjects.Threshold = subjectThreshold
				aggregateCount.CompletedSubjects.SubjectCount += project.SubjectCount.SubjectCount
				aggregateCount.CompletedSubjects.TotalEdits += lastProjectVersion.getTotalEdits()
				aggregateCount.CompletedSubjects.TotalFldEdits += lastProjectVersion.FieldEditMetrics.TotalEdits
				aggregateCount.CompletedSubjects.TotalFldEditsFired += lastProjectVersion.FieldEditMetrics.TotalFiredWithOpenQuery
				aggregateCount.CompletedSubjects.TotalFldEditsUnfired += lastProjectVersion.FieldEditMetrics.TotalNotFiredWithOpenQuery
				aggregateCount.CompletedSubjects.TotalFldEditsOpen += lastProjectVersion.FieldEditMetrics.TotalOpenQueries
				aggregateCount.CompletedSubjects.TotalFldWithChange += lastProjectVersion.FieldEditMetrics.TotalEditsFiredWithChange
				aggregateCount.CompletedSubjects.TotalFldWithNoChange += lastProjectVersion.FieldEditMetrics.TotalEditsFiredWithNoChange
				aggregateCount.CompletedSubjects.TotalPrgEdits += lastProjectVersion.ProgramEditMetrics.TotalEdits
				aggregateCount.CompletedSubjects.TotalPrgEditsWithOpenQuery += lastProjectVersion.ProgramEditMetrics.TotalEditsWithOpenQuery
				aggregateCount.CompletedSubjects.TotalPrgEditsFired += lastProjectVersion.ProgramEditMetrics.TotalFiredWithOpenQuery
				aggregateCount.CompletedSubjects.TotalPrgEditsUnfired += lastProjectVersion.ProgramEditMetrics.TotalNotFiredWithOpenQuery
				aggregateCount.CompletedSubjects.TotalPrgEditsOpen += lastProjectVersion.ProgramEditMetrics.TotalOpenQueries
				aggregateCount.CompletedSubjects.TotalPrgWithChange += lastProjectVersion.ProgramEditMetrics.TotalEditsFiredWithChange
				aggregateCount.CompletedSubjects.TotalPrgWithNoChange += lastProjectVersion.ProgramEditMetrics.TotalEditsFiredWithNoChange
			}
		}
		// All Subjects
		aggregateCount.AllProjects.RecordCount++
		aggregateCount.AllProjects.Threshold = subjectThreshold
		aggregateCount.AllProjects.SubjectCount += project.SubjectCount.SubjectCount
		aggregateCount.AllProjects.TotalEdits += lastProjectVersion.getTotalEdits()
		aggregateCount.AllProjects.TotalFldEdits += lastProjectVersion.FieldEditMetrics.TotalEdits
		aggregateCount.AllProjects.TotalFldEditsFired += lastProjectVersion.FieldEditMetrics.TotalFiredWithOpenQuery
		aggregateCount.AllProjects.TotalFldEditsUnfired += lastProjectVersion.FieldEditMetrics.TotalNotFiredWithOpenQuery
		aggregateCount.AllProjects.TotalFldEditsOpen += lastProjectVersion.FieldEditMetrics.TotalOpenQueries
		aggregateCount.AllProjects.TotalFldWithChange += lastProjectVersion.FieldEditMetrics.TotalEditsFiredWithChange
		aggregateCount.AllProjects.TotalFldWithNoChange += lastProjectVersion.FieldEditMetrics.TotalEditsFiredWithNoChange
		aggregateCount.AllProjects.TotalPrgEdits += lastProjectVersion.ProgramEditMetrics.TotalEdits
		aggregateCount.AllProjects.TotalPrgEditsWithOpenQuery += lastProjectVersion.ProgramEditMetrics.TotalEditsWithOpenQuery
		aggregateCount.AllProjects.TotalPrgEditsFired += lastProjectVersion.ProgramEditMetrics.TotalFiredWithOpenQuery
		aggregateCount.AllProjects.TotalPrgEditsUnfired += lastProjectVersion.ProgramEditMetrics.TotalNotFiredWithOpenQuery
		aggregateCount.AllProjects.TotalPrgEditsOpen += lastProjectVersion.ProgramEditMetrics.TotalOpenQueries
		aggregateCount.AllProjects.TotalPrgWithChange += lastProjectVersion.ProgramEditMetrics.TotalEditsFiredWithChange
		aggregateCount.AllProjects.TotalPrgWithNoChange += lastProjectVersion.ProgramEditMetrics.TotalEditsFiredWithNoChange

	}

	//headers := []string{
	//	"Criteria",
	//	"Threshold",
	//	"Aggregate Type",
	//	"Sample Size",
	//	"Subject Count",
	//	"Total Checks",
	//	"Programmed Checks",
	//	"Fired Checks",
	//	"Not Fired Checks",
	//	"Checks Leading to Change",
	//	"Checks Not Leading to Change",
	//	"Field Checks",
	//	"Fired Checks",
	//	"Not Fired Checks",
	//	"Checks Leading to Change",
	//	"Checks Not Leading to Change",
	//}
	sheet, _ := getOrAddSheet(wbk, "Summary Counts")
	// write the counts out
	writeAggregatedCounts(aggregateCount, sheet)
	//	writeNotes(sheet)
	// filter project -> subject count
	autoFilter := new(xlsx.AutoFilter)
	autoFilter.TopLeftCell = "A1"
	autoFilter.BottomRightCell = "E1"
	sheet.AutoFilter = autoFilter
	autoSizeSheet(sheet)
}
