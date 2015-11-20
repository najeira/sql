package sql

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

var (
	nullString []byte = []byte("null")
	escapeHTML bool   = false
)

// from https://github.com/cihangir/nisql
// Copyright (c) 2015 Cihangir

// String creates a valid NullString
func String(s string) NullString {
	return NullString{
		sql.NullString{
			String: s,
			Valid:  true,
		},
	}
}

// NullString is a type that can be null or a string
type NullString struct {
	sql.NullString
}

// MarshalJSON implements the json.Marshaler interface.
func (n *NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return nullString, nil
	}
	body, err := json.Marshal(n.String)
	if err != nil {
		return nil, err
	}
	if escapeHTML {
		var buf bytes.Buffer
		json.HTMLEscape(&buf, body)
		return buf.Bytes(), nil
	} else {
		return body, nil
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *NullString) UnmarshalJSON(b []byte) error {
	return unmarshal(n, b)
}

// Get returns nil or underlying string value
func (n *NullString) Get() *string {
	if !n.Valid {
		return nil
	}
	return &n.String
}

// Float64 creates a valid NullFloat64
func Float64(f float64) NullFloat64 {
	return NullFloat64{
		sql.NullFloat64{
			Float64: f,
			Valid:   true,
		},
	}
}

// NullFloat64 is a type that can be null or a float64
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON implements the json.Marshaler interface.
func (n *NullFloat64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return nullString, nil
	}
	return json.Marshal(n.Float64)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *NullFloat64) UnmarshalJSON(b []byte) error {
	return unmarshal(n, b)
}

// Get returns nil or underlying float64 value
func (n *NullFloat64) Get() *float64 {
	if !n.Valid {
		return nil
	}
	return &n.Float64
}

// Int64 creates a valid NullInt64
func Int64(i int64) NullInt64 {
	return NullInt64{
		sql.NullInt64{
			Int64: i,
			Valid: true,
		},
	}
}

// NullInt64 is a type that can be null or an int
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON implements the json.Marshaler interface.
func (n *NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return nullString, nil
	}
	return json.Marshal(n.Int64)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *NullInt64) UnmarshalJSON(b []byte) error {
	return unmarshal(n, b)
}

// Get returns nil or underlying int64 value
func (n *NullInt64) Get() *int64 {
	if !n.Valid {
		return nil
	}
	return &n.Int64
}

// Bool creates a valid NullBool
func Bool(b bool) NullBool {
	return NullBool{
		sql.NullBool{
			Bool:  b,
			Valid: true,
		},
	}
}

// NullBool is a type that can be null or a bool
type NullBool struct {
	sql.NullBool
}

// MarshalJSON implements the json.Marshaler interface.
func (n *NullBool) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return nullString, nil
	}
	return json.Marshal(n.Bool)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *NullBool) UnmarshalJSON(b []byte) error {
	return unmarshal(n, b)
}

// Get returns nil or underlying bool value
func (n *NullBool) Get() *bool {
	if !n.Valid {
		return nil
	}
	return &n.Bool
}

func unmarshal(s sql.Scanner, b []byte) error {
	var d interface{}
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	return s.Scan(d)
}
