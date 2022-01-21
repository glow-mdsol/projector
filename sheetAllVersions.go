package main

import (
	"github.com/tealeg/xlsx"
)

// Write a metric type to a row
func writeEditMetricType(editCheckTypeMetric *EditTypeMetric, row *xlsx.Row) {
	/*
		"Total Edits (fld)",
		"Total Edits Fired (fld)",
		"Total Edits Unfired (fld)",
		"Total Edits Open (fld)",
		"%ge Edits Fired (fld)",
		"%ge Edits Unfired (fld)",
		"Edits with Change (prg)",
		"Edits with No Change (prg)",
	*/
	// Total Edits
	cell := row.AddCell()
	cell.SetInt(editCheckTypeMetric.TotalEdits)
	// Total Fired
	cell = row.AddCell()
	cell.SetInt(editCheckTypeMetric.TotalFiredWithOpenQuery)
	// Total unfired
	cell = row.AddCell()
	cell.SetInt(editCheckTypeMetric.TotalNotFiredWithOpenQuery)
	// Percentage Fired
	cell = row.AddCell()
	cell.SetFloatWithFormat(editCheckTypeMetric.PercentageFiredWithOpenQuery, "0.00")
	// Percentage unfired
	cell = row.AddCell()
	cell.SetFloatWithFormat(editCheckTypeMetric.PercentageNotFiredWithOpenQuery, "0.00")
	// Edits leading to Change
	cell = row.AddCell()
	cell.SetInt(editCheckTypeMetric.TotalEditsFiredWithChange)
	// Edits not leading to Change
	cell = row.AddCell()
	cell.SetInt(editCheckTypeMetric.TotalEditsFiredWithNoChange)
	// TOTAL QUERIES
	cell = row.AddCell()
	cell.SetInt(editCheckTypeMetric.TotalQueries)
	// Total OPEN QUERIES
	cell = row.AddCell()
	cell.SetInt(editCheckTypeMetric.TotalOpenQueries)
}

func writeStudyMetricsForProject(urlName string, project *Project, wbk *xlsx.File) {
	// standard headers
	headers := []string{"Project Name",
		"CRF Version",
		"Last Version",
		"Active Edits",
		"Inactive Edits",
		"Total Edits (fld)",
		"Total Edits Fired (fld)",
		"Total Edits Unfired (fld)",
		"%ge Edits Fired (fld)",
		"%ge Edits Unfired (fld)",
		"Edits with Change (fld)",
		"Edits with No Change (fld)",
		"Total Queries (fld)",
		"Total Open Queries (fld)",
		"Total Edits (prg)",
		"Total Edits Fired (prg)",
		"Total Edits Unfired (prg)",
		"%ge Edits Fired (prg)",
		"%ge Edits Unfired (prg)",
		"Edits with Change (prg)",
		"Edits with No Change (prg)",
		"Total Queries (prg)",
		"Total Open Queries (prg)",
	}
	// log.Println("Reporting for", len(project.Versions), "versions of", project.ProjectName)
	//var colWidths []float64
	//// calculate the widths
	//for _, header := range headers {
	//	colWidths = append(colWidths, float64(len(header)))
	//}
	var sheet *xlsx.Sheet
	var created bool
	for _, projectVersion := range project.Versions {
		// create the sheet
		sheet, created = getOrAddSheet(wbk, urlName)
		// setup the fields
		if created {
			// Add the headers
			writeHeaderRow(headers, sheet)
			// filter project -> subject count
			autoFilter := new(xlsx.AutoFilter)
			autoFilter.TopLeftCell = "A1"
			autoFilter.BottomRightCell = "E1"
			sheet.AutoFilter = autoFilter
		}
		var cell *xlsx.Cell
		row := sheet.AddRow()
		// add the projectName
		cell = row.AddCell()
		cell.SetString(project.ProjectName)
		// CRF Version
		cell = row.AddCell()
		cell.SetInt(projectVersion.CRFVersionID)
		// Last Version
		cell = row.AddCell()
		if projectVersion.LastVersion {
			cell.SetString("Y")
		} else {
			cell.SetString("N")
		}
		// ActiveEdits
		cell = row.AddCell()
		cell.SetInt(projectVersion.EditStatus.ActiveEdits)
		// InActiveEdits
		cell = row.AddCell()
		cell.SetInt(projectVersion.EditStatus.InactiveEdits)
		// write the program metrics
		writeEditMetricType(&projectVersion.ProgramEditMetrics, row)
		// write the field metrics
		writeEditMetricType(&projectVersion.FieldEditMetrics, row)
	}
	autoSizeSheet(sheet)
}

// Just for the last version
func writeLastStudyMetricsForProject(urlName string, project *Project, wbk *xlsx.File) {
	tabName := urlName + " - Last"
	// standard headers
	headers := []string{"Project Name",
		"CRF Version",
		"Subject Count",
		"Active Edits",
		"Inactive Edits",
		"Total Edits (fld)",
		"Total Edits Fired (fld)",
		"Total Edits Unfired (fld)",
		"%ge Edits Fired (fld)",
		"%ge Edits Unfired (fld)",
		"Edits with Change (fld)",
		"Edits with No Change (fld)",
		"Total Queries (fld)",
		"Total Open Queries (fld)",
		"Total Edits (prg)",
		"Total Edits Fired (prg)",
		"Total Edits Unfired (prg)",
		"%ge Edits Fired (prg)",
		"%ge Edits Unfired (prg)",
		"Edits with Change (prg)",
		"Edits with No Change (prg)",
		"Total Queries (prg)",
		"Total Open Queries (prg)",
	}
	// log.Println("Reporting for", len(project.Versions), "versions of", project.ProjectName)
	//const maxWidth = 70
	//var colWidths []float64
	//for _, header := range headers {
	//	colWidths = append(colWidths, float64(len(header)))
	//}
	//
	var sheet *xlsx.Sheet
	var created bool
	for _, projectVersion := range project.Versions {
		if !projectVersion.LastVersion {
			continue
		}
		// create the sheet
		sheet, created = getOrAddSheet(wbk, tabName)
		// setup the fields
		if created {
			// Add the headers
			writeHeaderRow(headers, sheet)
			// filter project -> subject count
			autoFilter := new(xlsx.AutoFilter)
			autoFilter.TopLeftCell = "A1"
			autoFilter.BottomRightCell = "E1"
			sheet.AutoFilter = autoFilter
		}
		var cell *xlsx.Cell
		row := sheet.AddRow()
		// add the projectName
		cell = row.AddCell()
		//if float64(len(project.ProjectName)) > colWidths[0]{
		//	if len(project.ProjectName) > maxWidth{
		//		colWidths[0] = float64(maxWidth)
		//	} else {
		//		colWidths[0] = float64(len(project.ProjectName))
		//	}
		//	_ = sheet.SetColWidth(0, 0, colWidths[0])
		//}
		cell.SetString(project.ProjectName)
		// CRF Version
		cell = row.AddCell()
		cell.SetInt(projectVersion.CRFVersionID)
		// Subject Count
		cell = row.AddCell()
		cell.SetInt(project.SubjectCount.SubjectCount)
		//log.Println("Generating version",
		//	projectVersion.ProjectID, "(",
		//	projectVersion.CRFVersionID, ") ->",
		//	projectVersion.EditStatus.ActiveEdits)
		// ActiveEdits
		cell = row.AddCell()
		cell.SetInt(projectVersion.EditStatus.ActiveEdits)
		// InActiveEdits
		cell = row.AddCell()
		cell.SetInt(projectVersion.EditStatus.InactiveEdits)
		// write the program metrics
		writeEditMetricType(&projectVersion.ProgramEditMetrics, row)
		// write the field metrics
		writeEditMetricType(&projectVersion.FieldEditMetrics, row)
	}
	autoSizeSheet(sheet)
}
