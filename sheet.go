package main

import (
	"log"

	"github.com/tealeg/xlsx"
)

const MaxWidth float64 = 70.0

// initialise the set of columns
func initColumns(headers []string) []float64 {
	var colMax []float64
	for _, header := range headers {
		colMax = append(colMax, float64(len(header)))
	}
	return colMax
}

// wrap the adding of a sheet
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
func writeHeaderRow(data []string, sheet *xlsx.Sheet) (colWidth []float64) {
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
	for _, datum := range data {
		cell := row.AddCell()
		cell.SetStyle(titleFace)
		cell.SetString(datum)
		colWidth = append(colWidth, float64(len(datum)))
		//err := sheet.SetColWidth(idx, idx, float64(len(datum)))
		//if err != nil {
		//	panic("Unable to set width of column")
		//}
	}
	return
}

// Resize a sheet automatically
func autoSizeSheet(sheet *xlsx.Sheet) {
	var targetWidths []float64
	for _, row := range sheet.Rows {
		for idx, cell := range row.Cells {
			// initialise the array for cell
			if len(targetWidths) < idx+1 {
				targetWidths = append(targetWidths, 0.0)
			}
			contentWidth := float64(len(cell.Value))
			if contentWidth > targetWidths[idx] {
				if contentWidth < MaxWidth {
					targetWidths[idx] = contentWidth
				} else {
					targetWidths[idx] = MaxWidth
				}
			}
		}
	}
	for idx, tWidth := range targetWidths {
		_ = sheet.SetColWidth(idx, idx, tWidth)
	}
}
