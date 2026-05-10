package servicestate

import (
	"encoding/json"
	"testing"
)

// AO pass 2 — Action accessor sweep.

func TestServiceState_Uplift(t *testing.T) {
	a, err := New("status")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if a.Name() != "status" {
		t.Errorf("Name: %q", a.Name())
	}
	_ = NewMust("status")
	if _, err := New("__bogus__"); err == nil {
		t.Error("bogus should err")
	}

	if Max() == Invalid {
		t.Error("Max")
	}
	_ = RangesInvalidErr()

	if len(a.AllNameValues()) == 0 || len(a.IntegerEnumRanges()) == 0 || len(a.RangesByte()) == 0 {
		t.Error("ranges")
	}
	if a.RangeNamesCsv() == "" || a.TypeName() == "" || len(a.RangesDynamicMap()) == 0 {
		t.Error("range strings/map")
	}
	if a.String() == "" || a.ValueString() == "" || a.ToNumberString() == "" || a.NameValue() == "" || a.Format("%s") == "" || a.MinValueString() == "" || a.MaxValueString() == "" {
		t.Error("string accessors")
	}
	if mn, mx := a.MinMaxAny(); mn == nil || mx == nil {
		t.Error("MinMaxAny")
	}

	_ = a.MinByte()
	_ = a.MaxByte()
	_ = a.ValueByte()
	_ = a.MinInt()
	_ = a.MaxInt()
	_ = a.ValueInt()
	_ = a.ValueInt8()
	_ = a.ValueInt16()
	_ = a.ValueInt32()
	_ = a.Value()
	_ = a.ValueUInt16()

	if !a.IsValid() || a.IsInvalid() {
		t.Error("validity")
	}
	if !a.IsAnyOf(a, Invalid) {
		t.Error("IsAnyOf")
	}
	if !a.IsNameEqual("status") || !a.IsValueEqual(a.Value()) || !a.IsByteValueEqual(a.Value()) || !a.IsAnyValuesEqual(a.Value()) || !a.IsAnyNamesOf("status") {
		t.Error("predicates")
	}

	_ = a.OnlySupportedErr("status")
	_ = a.OnlySupportedMsgErr("ctx", "status")

	_ = a.EnumType()
	_ = a.AsActionTyper()
	_ = a.AsBasicByteEnumContractsBinder()
	_ = a.AsBasicEnumContractsBinder()

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var got Action
	if err := json.Unmarshal(data, &got); err != nil || got != a {
		t.Errorf("roundtrip: %v %v", got, err)
	}
	if _, err := a.UnmarshallEnumToValue(data); err != nil {
		t.Errorf("UnmarshallEnumToValue: %v", err)
	}

	for _, n := range a.IntegerEnumRanges() {
		x := Action(byte(n))
		_ = x.Name()
		_ = x.Format("%s")
		_ = x.IsValid()
	}
}
