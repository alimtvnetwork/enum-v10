package conntrackstate

import (
	"encoding/json"
	"testing"
)

// AO pass 2 — comprehensive Variant accessor sweep.

func TestConnTrackState_Uplift(t *testing.T) {
	v, err := Create("Established")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if v.Name() != "Established" {
		t.Errorf("Name: %q", v.Name())
	}
	_ = CreateMust("Established")
	if _, err := Create("__bogus__"); err == nil {
		t.Error("bogus should err")
	}

	if Max() == Invalid {
		t.Error("Max")
	}
	_ = RangesInvalidErr()

	if len(v.AllNameValues()) == 0 || len(v.IntegerEnumRanges()) == 0 || len(v.RangesByte()) == 0 {
		t.Error("ranges")
	}
	if v.RangeNamesCsv() == "" || v.TypeName() == "" || len(v.RangesDynamicMap()) == 0 {
		t.Error("range strings/map")
	}
	if v.String() == "" || v.ValueString() == "" || v.ToNumberString() == "" || v.NameValue() == "" || v.Format("%s") == "" || v.MinValueString() == "" || v.MaxValueString() == "" {
		t.Error("string accessors")
	}
	if mn, mx := v.MinMaxAny(); mn == nil || mx == nil {
		t.Error("MinMaxAny")
	}

	_ = v.MinByte()
	_ = v.MaxByte()
	_ = v.ValueByte()
	_ = v.MinInt()
	_ = v.MaxInt()
	_ = v.ValueInt()
	_ = v.ValueInt8()
	_ = v.ValueInt16()
	_ = v.ValueInt32()
	_ = v.Value()
	_ = v.ValueUInt16()

	if !v.IsValid() || v.IsInvalid() {
		t.Error("validity")
	}
	if !v.IsAnyOf(v, Invalid) {
		t.Error("IsAnyOf")
	}
	if !v.IsNameEqual("Established") || !v.IsValueEqual(v.Value()) || !v.IsByteValueEqual(v.Value()) || !v.IsAnyValuesEqual(v.Value()) || !v.IsAnyNamesOf("Established") {
		t.Error("predicates")
	}
	_ = v.IsEqual(v)
	_ = v.IsAboveOrEqual(v)
	_ = v.IsLowerOrEqual(v)

	_ = v.OnlySupportedErr("Established")
	_ = v.OnlySupportedMsgErr("ctx", "Established")

	_ = v.EnumType()
	_ = v.AsJsoner()
	_ = v.AsJsonContractsBinder()
	_ = v.AsJsonMarshaller()
	_ = v.AsBasicByteEnumContractsBinder()
	_ = v.AsBasicEnumContractsBinder()
	_ = v.ToPtr()

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var got Variant
	if err := json.Unmarshal(data, &got); err != nil || got != v {
		t.Errorf("roundtrip: %v %v", got, err)
	}
	_ = v.Json()
	_ = v.JsonPtr()
	if _, err := v.UnmarshallEnumToValue(data); err != nil {
		t.Errorf("UnmarshallEnumToValue: %v", err)
	}


	for _, n := range v.IntegerEnumRanges() {
		x := Variant(byte(n))
		_ = x.Name()
		_ = x.Format("%s")
		_ = x.IsValid()
	}
}
