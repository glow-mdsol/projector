package main

// ProjectVersion represents the structure for an individual Project Version
type ProjectVersion struct {
	URL               string
	URLID             int
	ProjectName       string
	CRFVersionID      int
	ActiveEditsOnly   *Record
	InactiveEditsOnly *Record
	LastVersion       bool
	ActiveEditCount   int
	InactiveEditCount int
}

func (pv *ProjectVersion) fixUpNullValues() {
	var activeCount = 0
	var inactiveCount = 0
	if pv.ActiveEditsOnly != nil {
		//log.Printf("INFO: Generating Active Edits for ProjectVersion %s(%d)", pv.ProjectName, pv.CRFVersionID)
		pv.ActiveEditsOnly.fixUpRecord()
		if pv.ActiveEditsOnly.TotalFieldEdits >= 0 {
			activeCount += pv.ActiveEditsOnly.TotalFieldEdits
		}
		if pv.ActiveEditsOnly.TotalProgEdits >= 0 {
			activeCount += pv.ActiveEditsOnly.TotalProgEdits
		}
	}
	if pv.InactiveEditsOnly != nil {
		//log.Printf("INFO: Generating InActive Edits for ProjectVersion %s(%d)", pv.ProjectName, pv.CRFVersionID)
		pv.InactiveEditsOnly.fixUpRecord()
		if pv.InactiveEditsOnly.TotalProgEdits >= 0 {
			// skip the null values
			inactiveCount += pv.InactiveEditsOnly.TotalProgEdits
		}
		if pv.InactiveEditsOnly.TotalFieldEdits >= 0 {
			// skip the null values
			inactiveCount += pv.InactiveEditsOnly.TotalFieldEdits
		}
	}
	pv.ActiveEditCount = activeCount
	pv.InactiveEditCount = inactiveCount
}
