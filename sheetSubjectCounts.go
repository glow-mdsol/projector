package main

import (
	"github.com/tealeg/xlsx"
)

// write the Subject Counts
func writeSubjectCount(urlName string, projects []*Project, wbk *xlsx.File) {
	tabName := "Subject Counts"
	headers := []string{"Rave URL",
		"Project Name",
		"Subject Count",
		"Screening Subject Count",
		"Screening Failure Count",
		"Enrolled Count",
		"Early Terminated Count",
		"Completed Count",
		"Enrolled in Follow Up",
		"Date Updated"}
	// maxWidth is an array of column widths
	var maxWidth = initColumns(headers)
	// create the sheet
	sheet, created := getOrAddSheet(wbk, tabName)
	if created {
		// Add the headers if it's newly created
		writeHeaderRow(headers, sheet)
	}
	autoFilter := new(xlsx.AutoFilter)
	autoFilter.TopLeftCell = "A1"
	autoFilter.BottomRightCell = "I1"
	sheet.AutoFilter = autoFilter

	// default widths
	dateLength := 18.5
	// hard code the width of the date
	maxWidth[9] = dateLength
	for _, project := range projects {
		var cell *xlsx.Cell
		// Rows
		row := sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(urlName)
		cell = row.AddCell()
		cell.SetString(project.ProjectName)
		//log.Println(project.ProjectName,"(",project.ProjectID,")","->",project.SubjectCount.SubjectCount)
		subjectCount := project.SubjectCount
		cell = row.AddCell()
		if subjectCount.SubjectCount > 0 {
			cell.SetInt(int(project.SubjectCount.SubjectCount))
		} else {
			cell.SetString("-")
		}
		// Screening Count
		cell = row.AddCell()
		if subjectCount.ScreeningCount.Valid {
			if subjectCount.ScreeningCount.Int64 > 0 {
				cell.SetInt64(subjectCount.ScreeningCount.Int64)
			} else {
				cell.SetString("-")
			}
		} else {
			cell.SetString("-")
		}
		// Screening Failure Count
		cell = row.AddCell()
		if subjectCount.ScreeningFailureCount.Valid {
			if subjectCount.ScreeningFailureCount.Int64 > 0 {
				cell.SetInt64(subjectCount.ScreeningFailureCount.Int64)
			} else {
				cell.SetString("-")
			}
		} else {
			cell.SetString("-")
		}
		// Enrolled Count
		cell = row.AddCell()
		if subjectCount.EnrolledCount.Valid {
			if subjectCount.EnrolledCount.Int64 > 0 {
				cell.SetInt64(subjectCount.EnrolledCount.Int64)
			} else {
				cell.SetString("-")
			}
		} else {
			cell.SetString("-")
		}
		// Early Terminated Count
		cell = row.AddCell()
		if subjectCount.EarlyTerminatedCount.Valid {
			if subjectCount.EarlyTerminatedCount.Int64 > 0 {
				cell.SetInt64(subjectCount.EarlyTerminatedCount.Int64)
			} else {
				cell.SetString("-")
			}
		} else {
			cell.SetString("-")
		}
		// Completed Count
		cell = row.AddCell()
		if subjectCount.CompletedCount.Valid {
			if subjectCount.CompletedCount.Int64 > 0 {
				cell.SetInt64(subjectCount.CompletedCount.Int64)
			} else {
				cell.SetString("-")
			}
		} else {
			cell.SetString("-")
		}
		// Enrolled in Follow Up Count
		cell = row.AddCell()
		if subjectCount.FollowUpCount.Valid {
			if subjectCount.FollowUpCount.Int64 > 0 {
				cell.SetInt64(subjectCount.FollowUpCount.Int64)
			} else {
				cell.SetString("-")
			}
		} else {
			cell.SetString("-")
		}
		// Refresh Date
		cell = row.AddCell()
		if subjectCount.RefreshDate.Valid {
			cell.SetDateTime(project.SubjectCount.RefreshDate.Time)

		} else {
			cell.SetString("-")
		}

	}
	// resize the sheet
	autoSizeSheet(sheet)
}
