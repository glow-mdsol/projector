package main

// Aggregate Count
type AggregateCount struct {
	ConditionDescription string
	AllProjects          SummaryCounts
	GreaterThanTen       SummaryCounts
	CompletedSubjects    SummaryCounts
}

// SummaryCounts represents the Structure for the computed stats
type SummaryCounts struct {
	// what criteria are we applying
	Threshold                         int
	RecordCount                       int
	SubjectCount                      int
	TotalEdits                        int
	TotalFldEdits                     int
	TotalFldEditsFired                int
	TotalFldEditsUnfired              int
	TotalFldEditsOpen                 int
	TotalFldWithChange                int
	TotalFldWithNoChange              int
	TotalPrgEdits                     int
	TotalPrgEditsWithOpenQuery        int
	TotalPrgEditsFired                int
	TotalPrgEditsUnfired              int
	TotalPrgEditsFiredWithOpenQuery   int
	TotalPrgEditsUnfiredWithOpenQuery int
	TotalPrgEditsOpen                 int
	TotalPrgWithChange                int
	TotalPrgWithNoChange              int
}

type AverageSummaryCounts struct {
	RecordCount                int
	SubjectCount               float64
	TotalEdits                 float64
	TotalFldEdits              float64
	TotalFldEditsFired         float64
	TotalFldEditsUnfired       float64
	TotalFldEditsOpen          float64
	TotalFldWithChange         float64
	TotalFldWithNoChange       float64
	TotalPrgEdits              float64
	TotalPrgEditsWithOpenQuery float64
	TotalPrgEditsFired         float64
	TotalPrgEditsUnfired       float64
	TotalPrgEditsOpen          float64
	TotalPrgWithChange         float64
	TotalPrgWithNoChange       float64
}

// Encapsulate the calculation of averages
func (av AverageSummaryCounts) init() {
	av.RecordCount = 0
	av.SubjectCount = 0.0
	av.TotalEdits = 0.0
	av.TotalFldEdits = 0.0
	av.TotalFldEditsFired = 0.0
	av.TotalFldEditsUnfired = 0.0
	av.TotalFldEditsOpen = 0.0
	av.TotalFldWithChange = 0.0
	av.TotalFldWithNoChange = 0.0
	av.TotalPrgEdits = 0.0
	av.TotalPrgEditsWithOpenQuery = 0.0
	av.TotalPrgEditsFired = 0.0
	av.TotalPrgEditsUnfired = 0.0
	av.TotalPrgEditsOpen = 0.0
	av.TotalPrgWithChange = 0.0
	av.TotalPrgWithNoChange = 0.0
}

func (sc *SummaryCounts) getAverageCounts() (avg AverageSummaryCounts) {
	// set the counts
	avg.init()
	if sc.RecordCount > 0 {
		avg.RecordCount = sc.RecordCount
		avg.SubjectCount = float64(sc.SubjectCount) / float64(sc.RecordCount)
		avg.TotalEdits = float64(sc.TotalEdits) / float64(sc.RecordCount)
		avg.TotalFldEdits = float64(sc.TotalFldEdits) / float64(sc.RecordCount)
		avg.TotalFldEditsFired = float64(sc.TotalFldEditsFired) / float64(sc.RecordCount)
		avg.TotalFldEditsUnfired = float64(sc.TotalFldEditsUnfired) / float64(sc.RecordCount)
		avg.TotalFldEditsOpen = float64(sc.TotalFldEditsOpen) / float64(sc.RecordCount)
		avg.TotalFldWithChange = float64(sc.TotalFldWithChange) / float64(sc.RecordCount)
		avg.TotalFldWithNoChange = float64(sc.TotalFldWithNoChange) / float64(sc.RecordCount)
		avg.TotalPrgEdits = float64(sc.TotalPrgEdits) / float64(sc.RecordCount)
		avg.TotalPrgEditsWithOpenQuery = float64(sc.TotalPrgEditsWithOpenQuery) / float64(sc.RecordCount)
		avg.TotalPrgEditsFired = float64(sc.TotalPrgEditsFired) / float64(sc.RecordCount)
		avg.TotalPrgEditsUnfired = float64(sc.TotalPrgEditsFired) / float64(sc.RecordCount)
		avg.TotalPrgEditsOpen = float64(sc.TotalPrgEditsOpen) / float64(sc.RecordCount)
		avg.TotalPrgWithChange = float64(sc.TotalPrgWithChange) / float64(sc.RecordCount)
		avg.TotalPrgWithNoChange = float64(sc.TotalPrgWithNoChange) / float64(sc.RecordCount)
	} else {
		avg.RecordCount = sc.RecordCount
	}
	return
}
