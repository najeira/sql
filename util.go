package sql

import (
	"database/sql"

	sqlxtypes "github.com/jmoiron/sqlx/types"
)

type (
	JSONText     = sqlxtypes.JSONText
	NullJSONText = sqlxtypes.NullJSONText
	Rows         = sql.Rows
	Row          = sql.Row
	Result       = sql.Result
	NullInt64    = sql.NullInt64
	NullString   = sql.NullString
)

func StringPtr(v string) *string {
	if len(v) <= 0 {
		return nil
	}
	return &v
}

func NullStringOf(v string) NullString {
	return NullString{
		String: v,
		Valid:  len(v) > 0,
	}
}

func StringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func NullInt64Of(v int64) NullInt64 {
	return NullInt64{
		Int64: v,
		Valid: v != 0,
	}
}
