package taskcategory

import (
	"encoding/json"
	"testing"
)

// AO pass 2 — comprehensive Variant accessor sweep.

func TestTaskCategory_Uplift(t *testing.T) {
	v, err := New("Help")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if v.Name() != "Help" {
		t.Errorf("Name: %q", v.Name())
	}
	_ = NewMust("Help")
	if _, err := New("__bogus__"); err == nil {
		t.Error("bogus should err")
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
	_ = v.ValueByte()
	_ = v.ValueUInt16()

	if !v.IsValid() || v.IsInvalid() {
		t.Error("validity")
	}
	if !v.IsAnyOf(v, Invalid) {
		t.Error("IsAnyOf")
	}
	if !v.IsNameEqual("Help") || !v.IsValueEqual(v.ValueByte()) || !v.IsByteValueEqual(v.ValueByte()) || !v.IsAnyValuesEqual(v.ValueByte()) || !v.IsAnyNamesOf("Help") {
		t.Error("predicates")
	}

	_ = v.OnlySupportedErr("Help")
	_ = v.OnlySupportedMsgErr("ctx", "Help")

	_ = v.EnumType()
	_ = v.AsJsoner()
	_ = v.AsJsonContractsBinder()
	_ = v.AsJsonMarshaller()
	_ = v.AsBasicByteEnumContractsBinder()
	_ = v.AsBasicEnumContractsBinder()

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
