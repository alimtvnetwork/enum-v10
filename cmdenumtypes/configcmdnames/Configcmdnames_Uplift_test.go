package configcmdnames

import (
	"encoding/json"
	"testing"
)

// AO uplift sweep — comprehensive Variant accessor exercise.

func TestCmdNamesUplift_configcmdnames(t *testing.T) {
	v, err := New("Help")
	if err != nil {
		t.Fatalf("New(Help): %v", err)
	}
	if v.Name() != "Help" {
		t.Errorf("Name: %q", v.Name())
	}
	_ = NewMust("Help")
	if _, err := New("__bogus__"); err == nil {
		t.Error("bogus name should error")
	}

	if Min() != Invalid {
		t.Errorf("Min: %v", Min())
	}
	if Max() == Invalid {
		t.Errorf("Max invalid")
	}

	if len(v.AllNameValues()) == 0 {
		t.Error("AllNameValues empty")
	}
	if len(v.IntegerEnumRanges()) == 0 {
		t.Error("IntegerEnumRanges empty")
	}
	if len(v.RangesByte()) == 0 {
		t.Error("RangesByte empty")
	}
	if v.RangeNamesCsv() == "" {
		t.Error("RangeNamesCsv empty")
	}
	if v.TypeName() == "" {
		t.Error("TypeName empty")
	}
	if len(v.RangesDynamicMap()) == 0 {
		t.Error("RangesDynamicMap empty")
	}

	if v.String() == "" {
		t.Error("String empty")
	}
	if v.ValueString() == "" {
		t.Error("ValueString empty")
	}
	if v.ToNumberString() == "" {
		t.Error("ToNumberString empty")
	}
	if v.NameValue() == "" {
		t.Error("NameValue empty")
	}
	if v.Format("%s") == "" {
		t.Error("Format empty")
	}
	if v.MinValueString() == "" {
		t.Error("MinValueString empty")
	}
	if v.MaxValueString() == "" {
		t.Error("MaxValueString empty")
	}
	_ = v.FullName()
	_ = v.HyphenName()
	_ = v.ToNameLower()

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
	if mn, mx := v.MinMaxAny(); mn == nil || mx == nil {
		t.Error("MinMaxAny nil")
	}

	if !v.IsValid() {
		t.Error("IsValid")
	}
	if v.IsInvalid() {
		t.Error("IsInvalid")
	}
	if !v.Is(v) {
		t.Error("Is self")
	}
	if !v.IsAnyOf(v, Invalid) {
		t.Error("IsAnyOf")
	}
	if !v.IsNameEqual("Help") {
		t.Error("IsNameEqual")
	}
	if !v.IsValueEqual(v.Value()) {
		t.Error("IsValueEqual")
	}
	if !v.IsByteValueEqual(v.Value()) {
		t.Error("IsByteValueEqual")
	}
	if !v.IsAnyValuesEqual(v.Value()) {
		t.Error("IsAnyValuesEqual")
	}
	if !v.IsAnyNamesOf("Help") {
		t.Error("IsAnyNamesOf")
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
	if err := json.Unmarshal(data, &got); err != nil {
		t.Errorf("Unmarshal: %v", err)
	}
	if got != v {
		t.Errorf("round-trip: %v != %v", got, v)
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

