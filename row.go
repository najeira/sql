package sql

import (
	"github.com/najeira/goutils/varutil"
)

type Row map[string]interface{}

func (m Row) Value(key string, args ...interface{}) interface{} {
	if v, ok := m[key]; ok {
		return v
	}
	if len(args) > 0 {
		return args[0]
	}
	return nil
}

func (m Row) String(key string, args ...interface{}) string {
	if v, ok := m[key]; ok {
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

func (m Row) Int(key string, args ...interface{}) int64 {
	if v, ok := m[key]; ok {
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

func (m Row) Float(key string, args ...interface{}) float64 {
	if v, ok := m[key]; ok {
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

func (m Row) Bool(key string, args ...interface{}) bool {
	if v, ok := m[key]; ok {
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
