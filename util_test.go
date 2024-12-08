package sql

import (
	"testing"
)

func TestStringPtr(t *testing.T) {
	if v := StringPtr("foo"); v == nil {
		t.Error("nil")
	}
	if v := StringPtr(""); v != nil {
		t.Error("not nil")
	}
}

func TestNullStringOf(t *testing.T) {
	if v := NullStringOf("foo"); !v.Valid {
		t.Error(v)
	} else if v.String != "foo" {
		t.Error(v)
	}

	if v := NullStringOf(""); v.Valid {
		t.Error(v)
	} else if v.String != "" {
		t.Error(v)
	}
}

func TestStringValue(t *testing.T) {
	s := "foo"
	if v := StringValue(&s); v != s {
		t.Error(v)
	}
	if v := StringValue(nil); v != "" {
		t.Error(v)
	}
}
