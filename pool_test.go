package sql

import (
	"testing"
)

func TestPoolGetPut(t *testing.T) {
	s := getString()
	if s == nil {
		t.Fatalf("getString: returns nil")
	}

	i := getInt64()
	if i == nil {
		t.Fatalf("getInt64: returns nil")
	}

	f := getFloat64()
	if f == nil {
		t.Fatalf("getFloat64: returns nil")
	}

	b := getBool()
	if b == nil {
		t.Fatalf("getBool: returns nil")
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
		t.Fatalf("getString: returns another one")
	}
	if s2.Valid != false {
		t.Errorf("getString: Valid is not zero")
	}
	if s2.String != "" {
		t.Errorf("getString: String is not zero")
	}

	i2 := getInt64()
	if i != i2 {
		t.Fatalf("getInt64: returns another one")
	}
	if i2.Valid != false {
		t.Errorf("getInt64: Valid is not zero")
	}
	if i2.Int64 != 0 {
		t.Errorf("getInt64: String is not zero")
	}

	f2 := getFloat64()
	if f != f2 {
		t.Fatalf("getFloat64: returns another one")
	}
	if f2.Valid != false {
		t.Errorf("getFloat64: Valid is not zero")
	}
	if f2.Float64 != 0 {
		t.Errorf("getFloat64: String is not zero")
	}

	b2 := getBool()
	if b != b2 {
		t.Fatalf("getBool: returns another one")
	}
	if b2.Valid != false {
		t.Errorf("getBool: Valid is not zero")
	}
	if b2.Bool != false {
		t.Errorf("getBool: String is not zero")
	}
}
