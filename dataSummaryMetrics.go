package main

import "database/sql"

// Field or Programmed Edit check
type EditCheckClass int

const (
	Field EditCheckClass = iota
	Programmed
)

// Active or Inactive
type EditCheckOutcome int

const (
	OpenQuery EditCheckOutcome = iota
	WithoutOpenQuery
)

// EditMetric (per ProjectVersion, per EditCheckClass)
type EditTypeMetric struct {
	// All the edits
	RawTotalEdits sql.NullInt64 `db:"total_edits"`
	TotalEdits    int
	// All the edits
	RawTotalActiveEdits sql.NullInt64 `db:"total_active_edits"`
	TotalActiveEdits    int
	// Edits with OpenQuery
	RawTotalEditsWithOpenQuery sql.NullInt64 `db:"total_edits_with_openquery"`
	TotalEditsWithOpenQuery    int
	// All the queries
	RawTotalQueries sql.NullInt64 `db:"total_queries"`
	TotalQueries    int
	// Total queries from checks with OpenQuery
	RawTotalQueriesOpenQuery sql.NullInt64 `db:"total_queries_open_query"`
	TotalQueriesOpenQuery    int
	// All the queries
	RawTotalOpenQueries sql.NullInt64 `db:"total_open_queries"`
	TotalOpenQueries    int
	// All Edits that fired
	RawTotalEditsFired sql.NullInt64 `db:"total_edits_fired"`
	TotalEditsFired    int
	// All Edits that didn't fire
	RawTotalEditsNotFired sql.NullInt64 `db:"total_edits_not_fired"`
	TotalEditsNotFired    int
	// Edits fired with OpenQuery (duh)
	RawTotalFiredWithOpenQuery sql.NullInt64 `db:"total_edits_fired_with_openquery"`
	TotalFiredWithOpenQuery    int
	// Edits not fired with OpenQuery (duh)
	RawTotalNotFiredWithOpenQuery sql.NullInt64 `db:"total_edits_not_fired_with_openquery"`
	TotalNotFiredWithOpenQuery    int
	// Fired with change in datapoint
	RawTotalEditsFiredWithChange sql.NullInt64 `db:"total_edits_fired_with_change"`
	TotalEditsFiredWithChange    int
	// No Change in datapoint
	RawTotalEditsFiredWithNoChange sql.NullInt64 `db:"total_edits_fired_without_change"`
	TotalEditsFiredWithNoChange    int
	// Total Queries Leading to Datapoint Change
	RawTotalQueriesWithChange sql.NullInt64 `db:"total_changes_from_edits"`
	TotalQueriesWithChange    int
	// Total Open Edits
	RawTotalOpenEdits sql.NullInt64 `db:"total_open_edits"`
	TotalOpenEdits    int
	// Percentages are calculated from Values
	PercentageFired                 float64
	PercentageNotFired              float64
	PercentageFiredWithOpenQuery    float64
	PercentageNotFiredWithOpenQuery float64
	PercentageChanged               float64
	PercentageNotChanged            float64
}

func (mtx *EditTypeMetric) fixUpMetrics() {
	if mtx.RawTotalEditsWithOpenQuery.Valid {
		mtx.TotalEditsWithOpenQuery = int(mtx.RawTotalEditsWithOpenQuery.Int64)
	} else {
		mtx.TotalEditsWithOpenQuery = -1
	}
	if mtx.RawTotalEdits.Valid {
		mtx.TotalEdits = int(mtx.RawTotalEdits.Int64)
	} else {
		mtx.TotalEdits = -1
	}
	if mtx.RawTotalEdits.Valid {
		mtx.TotalEdits = int(mtx.RawTotalEdits.Int64)
	} else {
		mtx.TotalEdits = -1
	}
	if mtx.RawTotalQueries.Valid {
		mtx.TotalQueries = int(mtx.RawTotalQueries.Int64)
	} else {
		mtx.TotalQueries = -1
	}
	if mtx.RawTotalQueries.Valid {
		mtx.TotalQueries = int(mtx.RawTotalQueries.Int64)
	} else {
		mtx.TotalQueries = -1
	}
	if mtx.RawTotalFiredWithOpenQuery.Valid {
		mtx.TotalFiredWithOpenQuery = int(mtx.RawTotalFiredWithOpenQuery.Int64)
	} else {
		mtx.TotalFiredWithOpenQuery = -1
	}
	if mtx.RawTotalFiredWithOpenQuery.Valid {
		mtx.TotalFiredWithOpenQuery = int(mtx.RawTotalFiredWithOpenQuery.Int64)
	} else {
		mtx.TotalFiredWithOpenQuery = -1
	}
	if mtx.RawTotalNotFiredWithOpenQuery.Valid {
		mtx.TotalNotFiredWithOpenQuery = int(mtx.RawTotalNotFiredWithOpenQuery.Int64)
	} else {
		mtx.TotalNotFiredWithOpenQuery = -1
	}
	if mtx.RawTotalEditsFiredWithChange.Valid {
		mtx.TotalEditsFiredWithChange = int(mtx.RawTotalEditsFiredWithChange.Int64)
	} else {
		mtx.TotalEditsFiredWithChange = -1
	}
	if mtx.RawTotalEditsFiredWithNoChange.Valid {
		mtx.TotalEditsFiredWithNoChange = int(mtx.RawTotalEditsFiredWithNoChange.Int64)
	} else {
		mtx.TotalEditsFiredWithNoChange = -1
	}
	if mtx.RawTotalQueriesOpenQuery.Valid {
		mtx.TotalQueriesOpenQuery = int(mtx.RawTotalQueriesOpenQuery.Int64)
	} else {
		mtx.TotalQueriesOpenQuery = -1
	}
	if mtx.RawTotalOpenEdits.Valid {
		mtx.TotalOpenEdits = int(mtx.RawTotalOpenEdits.Int64)
	} else {
		mtx.TotalOpenEdits = -1
	}
	if mtx.RawTotalOpenQueries.Valid {
		mtx.TotalOpenQueries = int(mtx.RawTotalOpenQueries.Int64)
	} else {
		mtx.TotalOpenQueries = -1
	}
}

func (mtx *EditTypeMetric) calculatePercentages() {
	// Gate the counts

	// Percentages fired versus not
	mtx.PercentageFired = 0.0
	mtx.PercentageNotFired = 0.0
	if mtx.TotalEdits > 0 {
		mtx.PercentageFired = 100.0 * float64(mtx.TotalEditsFired) / float64(mtx.TotalEdits)
		mtx.PercentageNotFired = 100.0 * float64(mtx.TotalEditsNotFired) / float64(mtx.TotalEdits)
	}

	// Leaf
	mtx.PercentageChanged = 0.0
	mtx.PercentageNotChanged = 0.0
	if mtx.TotalEditsFired > 0 {
		mtx.PercentageChanged = 100.0 * float64(mtx.TotalEditsFiredWithChange) / float64(mtx.TotalEditsFired)
		mtx.PercentageNotChanged = 100.0 * float64(mtx.TotalEditsFiredWithNoChange) / float64(mtx.TotalEditsFired)
	}
	// OpenQuery related
	mtx.PercentageFiredWithOpenQuery = 0.0
	mtx.PercentageNotFiredWithOpenQuery = 0.0
	if mtx.TotalEditsWithOpenQuery > 0 {
		mtx.PercentageFiredWithOpenQuery = 100.0 * float64(mtx.TotalFiredWithOpenQuery) / float64(mtx.TotalEditsWithOpenQuery)
		mtx.PercentageNotFiredWithOpenQuery = 100.0 * float64(mtx.TotalNotFiredWithOpenQuery) / float64(mtx.TotalEditsWithOpenQuery)
	}

}
