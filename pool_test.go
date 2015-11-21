package sql

import (
	"testing"
)

func TestPoolGetPut(t *testing.T) {
	s := getString()
	if s == nil {
		t.Fatalf("getString: nil")
	}

	i := getInt64()
	if i == nil {
		t.Fatalf("getInt64: nil")
	}

	f := getFloat64()
	if f == nil {
		t.Fatalf("getFloat64: nil")
	}

	b := getBool()
	if b == nil {
		t.Fatalf("getBool: nil")
	}

	s.Valid = true
	s.String = "hoge"
	poolString(s)

	i.Valid = true
	i.Int64 = 123
	poolInt64(i)

	f.Valid = true
	f.Float64 = 123.456
	poolFloat64(f)

	b.Valid = true
	b.Bool = true
	poolBool(b)

	s2 := getString()
	if s != s2 {
		t.Fatalf("getString: another one")
	}
	if s2.Valid != false {
		t.Errorf("getString: Valid is not zero")
	}
	if s2.String != "" {
		t.Errorf("getString: String is not zero")
	}

	i2 := getInt64()
	if i != i2 {
		t.Fatalf("getInt64: another one")
	}
	if i2.Valid != false {
		t.Errorf("getInt64: Valid is not zero")
	}
	if i2.Int64 != 0 {
		t.Errorf("getInt64: String is not zero")
	}

	f2 := getFloat64()
	if f != f2 {
		t.Fatalf("getFloat64: another one")
	}
	if f2.Valid != false {
		t.Errorf("getFloat64: Valid is not zero")
	}
	if f2.Float64 != 0 {
		t.Errorf("getFloat64: String is not zero")
	}

	b2 := getBool()
	if b != b2 {
		t.Fatalf("getBool: another one")
	}
	if b2.Valid != false {
		t.Errorf("getBool: Valid is not zero")
	}
	if b2.Bool != false {
		t.Errorf("getBool: String is not zero")
	}
}

func TestPoolValuesGetPut(t *testing.T) {
	v := getValues()
	if v == nil {
		t.Fatalf("getValues: nil")
	}

	s := v.String()
	if s == nil {
		t.Fatalf("values.String: nil")
	}

	i := v.Int64()
	if i == nil {
		t.Fatalf("values.Int64: nil")
	}

	f := v.Float64()
	if f == nil {
		t.Fatalf("values.Float64: nil")
	}

	b := v.Bool()
	if b == nil {
		t.Fatalf("values.Bool: nil")
	}

	s.Valid = true
	s.String = "hoge"

	i.Valid = true
	i.Int64 = 123

	f.Valid = true
	f.Float64 = 123.456

	b.Valid = true
	b.Bool = true

	v.Close()

	v2 := getValues()
	if v2 != v {
		t.Fatalf("getValues: another one")
	}

	s2 := v2.String()
	if s != s2 {
		t.Fatalf("values.String: another one")
	}
	if s2.Valid != false {
		t.Errorf("values.String: Valid is not zero")
	}
	if s2.String != "" {
		t.Errorf("values.String: String is not zero")
	}

	i2 := v2.Int64()
	if i != i2 {
		t.Fatalf("values.Int64: another one")
	}
	if i2.Valid != false {
		t.Errorf("values.Int64: Valid is not zero")
	}
	if i2.Int64 != 0 {
		t.Errorf("values.Int64: String is not zero")
	}

	f2 := v2.Float64()
	if f != f2 {
		t.Fatalf("values.Float64: another one")
	}
	if f2.Valid != false {
		t.Errorf("values.Float64: Valid is not zero")
	}
	if f2.Float64 != 0 {
		t.Errorf("values.Float64: String is not zero")
	}

	b2 := v2.Bool()
	if b != b2 {
		t.Fatalf("values.Bool: another one")
	}
	if b2.Valid != false {
		t.Errorf("values.Bool: Valid is not zero")
	}
	if b2.Bool != false {
		t.Errorf("values.Bool: String is not zero")
	}
}
