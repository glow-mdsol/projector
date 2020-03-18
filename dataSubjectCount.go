package main

import (
	"database/sql"
	"github.com/lib/pq"
)

// SubjectCount represents the Structure for an Subject Count
type SubjectCount struct {
	URLID                 int           `db:"url_id"`
	ProjectID             int           `db:"project_id"`
	ProjectName           string        `db:"project_name"`
	RefreshDate           pq.NullTime   `db:"refresh_date"`
	SubjectCount          int           `db:"subject_count"`
	ScreeningCount        sql.NullInt64 `db:"screening_subject_count"`
	ScreeningFailureCount sql.NullInt64 `db:"screening_failure_subject_count"`
	EnrolledCount         sql.NullInt64 `db:"enrolled_subject_count"`
	EarlyTerminatedCount  sql.NullInt64 `db:"early_terminated_subject_count"`
	CompletedCount        sql.NullInt64 `db:"completed_subject_count"`
	FollowUpCount         sql.NullInt64 `db:"follow_up_subject_count"`
}
