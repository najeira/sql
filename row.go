package sql

import (
	"github.com/najeira/goutils/varutil"
)

// Row is the result of query.
type Row map[string]interface{}

// Value returns the value for the column in the row.
// It returns nil if not found.
func (m Row) Value(column string, args ...interface{}) interface{} {
	if v, ok := m[column]; ok {
		return v
	}
	if len(args) > 0 {
		return args[0]
	}
	return nil
}

// String returns the string value for the column in the row.
// It returns empty string if not found or it is not string.
func (m Row) String(column string, args ...interface{}) string {
	if v, ok := m[column]; ok {
		if ns, ok := v.(NullString); ok && ns.Valid {
			return ns.String
		} else if ns, ok := v.(*NullString); ok && ns.Valid {
			return ns.String
		} else if s, ok := varutil.String(v); ok {
			return s
		}
	}
	if len(args) > 0 {
		return varutil.AsString(args[0])
	}
	return ""
}

// Int returns the integer value for the column in the row.
// It returns 0 if not found or it is not integer.
func (m Row) Int(column string, args ...interface{}) int64 {
	if v, ok := m[column]; ok {
		if ns, ok := v.(NullInt64); ok && ns.Valid {
			return ns.Int64
		} else if ns, ok := v.(*NullInt64); ok && ns.Valid {
			return ns.Int64
		} else if s, ok := varutil.Int(v); ok {
			return s
		}
	}
	if len(args) > 0 {
		return varutil.AsInt(args[0])
	}
	return 0
}

// Float returns the float value for the column in the row.
// It returns 0 if not found or it is not float.
func (m Row) Float(column string, args ...interface{}) float64 {
	if v, ok := m[column]; ok {
		if ns, ok := v.(NullFloat64); ok && ns.Valid {
			return ns.Float64
		} else if ns, ok := v.(*NullFloat64); ok && ns.Valid {
			return ns.Float64
		} else if s, ok := varutil.Float(v); ok {
			return s
		}
	}
	if len(args) > 0 {
		return varutil.AsFloat(args[0])
	}
	return 0
}

// Bool returns the boolean value for the column in the row.
// It returns false if not found or it is not boolean.
func (m Row) Bool(column string, args ...interface{}) bool {
	if v, ok := m[column]; ok {
		if ns, ok := v.(NullBool); ok && ns.Valid {
			return ns.Bool
		} else if ns, ok := v.(*NullBool); ok && ns.Valid {
			return ns.Bool
		} else if s, ok := varutil.Bool(v); ok {
			return s
		}
	}
	if len(args) > 0 {
		return varutil.AsBool(args[0])
	}
	return false
}
