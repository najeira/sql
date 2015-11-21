package sql

import (
	"testing"
)

func TestRow(t *testing.T) {
	r := make(Row)
	if len(r) != 0 {
		t.Errorf("Row len: expected %d, got %d", 0, len(r))
	}

	r["name"] = "Test"
	r["age"] = 30
	r["score"] = 123.456
	r["flag"] = true

	if r.String("name") != "Test" {
		t.Errorf("Row.String: expected %s, got %s", "Test", r.String("name"))
	}
	if r.String("age") != "" {
		t.Errorf("Row.String: expected %s, got %s", "", r.String("age"))
	}
	if r.String("score") != "" {
		t.Errorf("Row.String: expected %s, got %s", "", r.String("score"))
	}
	if r.String("flag") != "" {
		t.Errorf("Row.String: expected %s, got %s", "", r.String("flag"))
	}

	if r.Int("name") != 0 {
		t.Errorf("Row.Int: expected %s, got %s", 0, r.String("name"))
	}
	if r.Int("age") != 30 {
		t.Errorf("Row.Int: expected %s, got %s", 30, r.String("age"))
	}
	if r.Int("score") != 0 {
		t.Errorf("Row.Int: expected %s, got %s", 0, r.String("score"))
	}
	if r.Int("flag") != 0 {
		t.Errorf("Row.Int: expected %s, got %s", 0, r.String("flag"))
	}

	if r.Float("name") != 0 {
		t.Errorf("Row.Float: expected %s, got %s", 0, r.String("name"))
	}
	if r.Float("age") != 0 {
		t.Errorf("Row.Float: expected %s, got %s", 0, r.String("age"))
	}
	if r.Float("score") != 123.456 {
		t.Errorf("Row.Float: expected %s, got %s", 123.456, r.String("score"))
	}
	if r.Float("flag") != 0 {
		t.Errorf("Row.Float: expected %s, got %s", 0, r.String("flag"))
	}

	if r.Bool("name") != false {
		t.Errorf("Row.Bool: expected %s, got %s", false, r.Bool("name"))
	}
	if r.Bool("age") != false {
		t.Errorf("Row.Bool: expected %s, got %s", false, r.Bool("age"))
	}
	if r.Bool("score") != false {
		t.Errorf("Row.Bool: expected %s, got %s", false, r.Bool("score"))
	}
	if r.Bool("flag") != true {
		t.Errorf("Row.Bool: expected %s, got %s", true, r.Bool("flag"))
	}
}

func TestRowNullValue(t *testing.T) {
	r := make(Row)
	if len(r) != 0 {
		t.Errorf("Row len: expected %d, got %d", 0, len(r))
	}

	r["name"] = String("Test")
	r["age"] = Int64(30)
	r["score"] = Float64(123.456)
	r["flag"] = Bool(true)

	if r.String("name") != "Test" {
		t.Errorf("Row.String: expected %s, got %s", "Test", r.String("name"))
	}
	if r.String("age") != "" {
		t.Errorf("Row.String: expected %s, got %s", "", r.String("age"))
	}
	if r.String("score") != "" {
		t.Errorf("Row.String: expected %s, got %s", "", r.String("score"))
	}
	if r.String("flag") != "" {
		t.Errorf("Row.String: expected %s, got %s", "", r.String("flag"))
	}

	if r.Int("name") != 0 {
		t.Errorf("Row.Int: expected %s, got %s", 0, r.String("name"))
	}
	if r.Int("age") != 30 {
		t.Errorf("Row.Int: expected %s, got %s", 30, r.String("age"))
	}
	if r.Int("score") != 0 {
		t.Errorf("Row.Int: expected %s, got %s", 0, r.String("score"))
	}
	if r.Int("flag") != 0 {
		t.Errorf("Row.Int: expected %s, got %s", 0, r.String("flag"))
	}

	if r.Float("name") != 0 {
		t.Errorf("Row.Float: expected %s, got %s", 0, r.String("name"))
	}
	if r.Float("age") != 0 {
		t.Errorf("Row.Float: expected %s, got %s", 0, r.String("age"))
	}
	if r.Float("score") != 123.456 {
		t.Errorf("Row.Float: expected %s, got %s", 123.456, r.String("score"))
	}
	if r.Float("flag") != 0 {
		t.Errorf("Row.Float: expected %s, got %s", 0, r.String("flag"))
	}

	if r.Bool("name") != false {
		t.Errorf("Row.Bool: expected %s, got %s", false, r.Bool("name"))
	}
	if r.Bool("age") != false {
		t.Errorf("Row.Bool: expected %s, got %s", false, r.Bool("age"))
	}
	if r.Bool("score") != false {
		t.Errorf("Row.Bool: expected %s, got %s", false, r.Bool("score"))
	}
	if r.Bool("flag") != true {
		t.Errorf("Row.Bool: expected %s, got %s", true, r.Bool("flag"))
	}
}

func TestRowNullValuePointer(t *testing.T) {
	r := make(Row)
	if len(r) != 0 {
		t.Errorf("Row len: expected %d, got %d", 0, len(r))
	}

	s := String("Test")
	i := Int64(30)
	f := Float64(123.456)
	b := Bool(true)
	r["name"] = &s
	r["age"] = &i
	r["score"] = &f
	r["flag"] = &b

	if r.String("name") != "Test" {
		t.Errorf("Row.String: expected %s, got %s", "Test", r.String("name"))
	}
	if r.String("age") != "" {
		t.Errorf("Row.String: expected %s, got %s", "", r.String("age"))
	}
	if r.String("score") != "" {
		t.Errorf("Row.String: expected %s, got %s", "", r.String("score"))
	}
	if r.String("flag") != "" {
		t.Errorf("Row.String: expected %s, got %s", "", r.String("flag"))
	}

	if r.Int("name") != 0 {
		t.Errorf("Row.Int: expected %s, got %s", 0, r.String("name"))
	}
	if r.Int("age") != 30 {
		t.Errorf("Row.Int: expected %s, got %s", 30, r.String("age"))
	}
	if r.Int("score") != 0 {
		t.Errorf("Row.Int: expected %s, got %s", 0, r.String("score"))
	}
	if r.Int("flag") != 0 {
		t.Errorf("Row.Int: expected %s, got %s", 0, r.String("flag"))
	}

	if r.Float("name") != 0 {
		t.Errorf("Row.Float: expected %s, got %s", 0, r.String("name"))
	}
	if r.Float("age") != 0 {
		t.Errorf("Row.Float: expected %s, got %s", 0, r.String("age"))
	}
	if r.Float("score") != 123.456 {
		t.Errorf("Row.Float: expected %s, got %s", 123.456, r.String("score"))
	}
	if r.Float("flag") != 0 {
		t.Errorf("Row.Float: expected %s, got %s", 0, r.String("flag"))
	}

	if r.Bool("name") != false {
		t.Errorf("Row.Bool: expected %s, got %s", false, r.Bool("name"))
	}
	if r.Bool("age") != false {
		t.Errorf("Row.Bool: expected %s, got %s", false, r.Bool("age"))
	}
	if r.Bool("score") != false {
		t.Errorf("Row.Bool: expected %s, got %s", false, r.Bool("score"))
	}
	if r.Bool("flag") != true {
		t.Errorf("Row.Bool: expected %s, got %s", true, r.Bool("flag"))
	}
}
