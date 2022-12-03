package sql

import (
	sqld "database/sql"

	sqlxtypes "github.com/jmoiron/sqlx/types"
)

type (
	JSONText     = sqlxtypes.JSONText
	NullJSONText = sqlxtypes.NullJSONText
	Rows         = sqld.Rows
	Row          = sqld.Row
	Result       = sqld.Result
	NullInt64    = sqld.NullInt64
	NullString   = sqld.NullString
)
