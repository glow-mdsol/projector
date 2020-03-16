package main

import (
	"github.com/tealeg/xlsx"
	"strings"
)

func writeUselessEdits(projectName string, edits []*UnusedEdit, checkOutcome EditCheckOutcome, wbk *xlsx.File) {
	headers := []string{"Project Name",
		"Edit Check Name",
		"Form OID",
		"Field OID",
		"Variable OID",
		"Times Used",
		"Custom Function?",
		"Non-conformance check?",
		"Required check?",
		"Future check?",
		"Range check?",
	}
	var tabName string
	if checkOutcome == OpenQuery {
		tabName = "Unused Edits w OpenQuery"
		//log.Println("Printing", len(edits), "edits with OpenQuery")
	} else {
		tabName = "Unused Edits wo OpenQuery"
		//log.Println("Printing", len(edits), "edits without OpenQuery")
	}
	// create the sheet
	sheet, created := getOrAddSheet(wbk, tabName)
	if created {
		// Add the headers
		colWidths := writeHeaderRow(headers, sheet)
		autoFilter := new(xlsx.AutoFilter)
		autoFilter.TopLeftCell = "A1"
		autoFilter.BottomRightCell = "M1"
		sheet.AutoFilter = autoFilter
		for idx, width := range colWidths {
			_ = sheet.SetColWidth(idx, idx, width)
		}
	}

	// minimum lengths
	projectLength := 12
	checkLength := 12
	formOIDLength := 12
	fieldOIDLength := 12
	vblOIDLength := 12
	// Hard code the upper limit
	maxLength := 70

	// Export the results
	for _, edit := range edits {
		if len(projectName) > projectLength {
			if len(projectName) < maxLength {
				projectLength = len(projectName)
			} else {
				projectLength = maxLength
			}
			_ = sheet.SetColWidth(0, 0, float64(projectLength))
		}
		if len(edit.EditCheckName) > checkLength {
			if len(edit.EditCheckName) < maxLength {
				checkLength = len(edit.EditCheckName)
			} else {
				checkLength = maxLength
			}
			_ = sheet.SetColWidth(1, 1, float64(checkLength))
		}
		if len(edit.FormOID) > formOIDLength {
			if len(edit.FormOID) < maxLength {
				formOIDLength = len(edit.FormOID)
			} else {
				formOIDLength = maxLength
			}
			_ = sheet.SetColWidth(2, 2, float64(formOIDLength))
		}
		if len(edit.FieldOID) > fieldOIDLength {
			if len(edit.FieldOID) < maxLength {
				fieldOIDLength = len(edit.FieldOID)
			} else {
				fieldOIDLength = maxLength
			}
			_ = sheet.SetColWidth(3, 3, float64(fieldOIDLength))
		}
		if len(edit.VariableOID) > vblOIDLength {
			if len(edit.VariableOID) < maxLength {
				vblOIDLength = len(edit.VariableOID)
			} else {
				vblOIDLength = maxLength
			}
			_ = sheet.SetColWidth(4, 4, float64(vblOIDLength))
		}
		var cell *xlsx.Cell
		// Rows
		row := sheet.AddRow()
		// Project Name
		cell = row.AddCell()
		cell.SetString(projectName)
		// Edit Name
		cell = row.AddCell()
		cell.SetString(edit.EditCheckName)
		// Form OID
		cell = row.AddCell()
		cell.SetString(edit.FormOID)
		cell = row.AddCell()
		cell.SetString(edit.FieldOID)
		cell = row.AddCell()
		cell.SetString(edit.VariableOID)
		cell = row.AddCell()
		cell.SetInt(edit.UsageCount)
		cell = row.AddCell()
		if edit.CustomFunction {
			cell.SetString("Y")
		} else {
			cell.SetString("N")
		}
		cell = row.AddCell()
		// Non-conformant
		if strings.HasPrefix(edit.EditCheckName, "SYS_NC_") {
			cell.SetString("Y")
		} else {
			cell.SetString("N")
		}
		cell = row.AddCell()
		// cell.SetString(edit.RequiredCheck)
		if strings.HasPrefix(edit.EditCheckName, "SYS_REQ_") {
			cell.SetString("Y")
		} else {
			cell.SetString("N")
		}
		cell = row.AddCell()
		// Future Date Check
		if strings.HasPrefix(edit.EditCheckName, "SYS_FUTURE_") {
			cell.SetString("Y")
		} else {
			cell.SetString("N")
		}
		cell = row.AddCell()
		// Range Check
		if strings.HasPrefix(edit.EditCheckName, "SYS_Q_RANGE_") {
			cell.SetString("Y")
		} else {
			cell.SetString("N")
		}
		// cell.SetString(edit.RangeCheck)
	}
	//sheet.SetColWidth(0, 0, float64(projectLength))
	//sheet.SetColWidth(1, 1, float64(checkLength))
	//sheet.SetColWidth(2, 2, float64(formOIDLength))
	//sheet.SetColWidth(3, 3, float64(fieldOIDLength))
	//sheet.SetColWidth(4, 4, float64(vblOIDLength))

}
