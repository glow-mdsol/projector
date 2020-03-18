package main

// UnusedEdit represents the Structure for an unused Edit
type UnusedEdit struct {
	URLID          string
	ProjectID      int    `db:"project_id"`
	ProjectName    string `db:"project_name"`
	EditCheckName  string `db:"edit_check_name"`
	FormOID        string `db:"form_oids"`
	FieldOID       string `db:"field_oids"`
	VariableOID    string `db:"variable_oids"`
	UsageCount     int    `db:"total_count"`
	OpenQuery      string `db:"open_query"`
	CustomFunction bool
}
