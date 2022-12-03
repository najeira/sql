package sql

import (
	"database/sql"
	"reflect"
)

func RowsAffected(i interface{}) int {
	if i == nil {
		return 0
	}

	if _, ok := i.(*sql.Rows); ok {
		return 0
	}

	if _, ok := i.(*sql.Row); ok {
		return 0
	}

	if res, ok := i.(sql.Result); ok {
		if res != nil {
			n, _ := res.RowsAffected()
			return int(n)
		}
		return 0
	}

	return destCount(i)
}

func destCount(i interface{}) int {
	if i == nil {
		return 0
	}

	// nolint:exhaustive
	switch i.(type) {
	case []byte, *[]byte, string, *string:
		return 1
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return 0
		}
		v = reflect.Indirect(v)
	}

	// nolint:exhaustive
	switch v.Kind() {
	case reflect.Array:
		return v.Len()
	case reflect.Slice:
		return v.Len()
	}
	return 1
}
