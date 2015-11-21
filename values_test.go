package sql

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type nullable struct {
	StringNVal NullString
	StringVal  string

	Int64NVal NullInt64
	Int64Val  int64

	Float64NVal NullFloat64
	Float64Val  float64

	BoolNVal NullBool
	BoolVal  bool
}

func dummy() *nullable {
	return &nullable{
		StringNVal: String("string_n_val"),
		StringVal:  "string_val",

		Int64NVal: Int64(123),

		Int64Val: int64(123),

		Float64NVal: Float64(12),
		Float64Val:  float64(12),

		BoolNVal: Bool(true),
		BoolVal:  true,
	}
}

func TestMarshal(t *testing.T) {
	nset := dummy()

	allSetRes, err := json.Marshal(nset)
	if err != nil {
		t.Fatalf("err while marshaling: %s", err.Error())
	}

	// test all-set variables marshaling
	allSetExpectedResString := `{"StringNVal":"string_n_val","StringVal":"string_val","Int64NVal":123,"Int64Val":123,"Float64NVal":12,"Float64Val":12,"BoolNVal":true,"BoolVal":true}`
	if allSetExpectedResString != string(allSetRes) {
		t.Fatalf("Marshal err: expected: %s, got: %s", allSetExpectedResString, string(allSetRes))
	}

	// test not-set variables marshalling
	nnonset := &nullable{}
	nonSetRes, err := json.Marshal(nnonset)
	if err != nil {
		t.Fatalf("err while marshaling:%s", err.Error())
	}

	nonSetExpectedResString := `{"StringNVal":null,"StringVal":"","Int64NVal":null,"Int64Val":0,"Float64NVal":null,"Float64Val":0,"BoolNVal":null,"BoolVal":false}`
	if nonSetExpectedResString != string(nonSetRes) {
		t.Fatalf("Marshal err: expected: %s, got: %s", nonSetExpectedResString, string(nonSetRes))
	}
}

func TestUnMarshal(t *testing.T) {
	nset := dummy()

	allSetRes, err := json.Marshal(nset)
	if err != nil {
		t.Fatalf("err while marshaling: %s", err.Error())
	}

	// test not-set variables marshalling
	nnonset := &nullable{}
	nonSetRes, err := json.Marshal(nnonset)
	if err != nil {
		t.Fatalf("err while marshaling:%s", err.Error())
	}

	// test all set variables unmarshalling
	nset2 := &nullable{}
	if err := json.Unmarshal(allSetRes, nset2); err != nil {
		t.Fatalf("Unmarshal err: %s", err.Error())
	}

	if !reflect.DeepEqual(nset, nset2) {
		t.Fatalf("not same: \nn:%#v,\nn2:%#v", nset, nset2)
	}

	// test not-set variables unmarshaling
	nnonset2 := &nullable{}
	if err := json.Unmarshal(nonSetRes, nnonset2); err != nil {
		t.Fatalf("Unmarshal err: %s", err.Error())
	}

	testGetNonNil(t, nset)
	testGetNil(t, nnonset)

	testGetNonNil(t, nset2)
	testGetNil(t, nnonset2)
}

func TestMarshalNullBool(t *testing.T) {
	nb := &NullBool{}
	err := nb.UnmarshalJSON([]byte("null1"))
	if err == nil {
		t.Fatal("null1 is not a valid bool")
	}
}

func testGetNil(t *testing.T, n *nullable) {
	testF(t, "n.StringNVal.Get() == nil", n.StringNVal.Get() == nil)
	testF(t, "n.Int64NVal.Get() == nil", n.Int64NVal.Get() == nil)
	testF(t, "n.Float64NVal.Get() == nil", n.Float64NVal.Get() == nil)
	testF(t, "n.BoolNVal.Get() == nil", n.BoolNVal.Get() == nil)
}

func testGetNonNil(t *testing.T, n *nullable) {
	testF(t, "n.StringNVal.Get() != nil", n.StringNVal.Get() != nil)
	testF(t, "n.Int64NVal.Get() != nil", n.Int64NVal.Get() != nil)
	testF(t, "n.Float64NVal.Get() != nil", n.Float64NVal.Get() != nil)
	testF(t, "n.BoolNVal.Get() != nil", n.BoolNVal.Get() != nil)
}

func testF(tb testing.TB, msg string, res bool) {
	if !res {
		fmt.Printf("exp: %s\n", msg)
		tb.Fail()
	}
}
