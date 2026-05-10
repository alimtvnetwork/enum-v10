package compresslevels

import "testing"

// TestDirectMethods_CompressLevels covers methods with non-nullary signatures
// the reflective Uplift sweep cannot reach (BA: per-package coverage uplift).
func TestDirectMethods_CompressLevels(t *testing.T) {
	a, b := Default, Best

	if !a.IsEqual(a) || a.IsEqual(b) {
		t.Error("IsEqual broken")
	}
	if !a.IsValueEqual(a.ValueByte()) {
		t.Error("IsValueEqual broken")
	}
	if !a.IsByteValueEqual(a.ValueByte()) {
		t.Error("IsByteValueEqual broken")
	}
	if !a.IsInteger8ValueEqual(a.ValueInt8()) {
		t.Error("IsInteger8ValueEqual broken")
	}
	if !a.IsNameEqual(a.Name()) || a.IsNameEqual("__nope__") {
		t.Error("IsNameEqual broken")
	}
	_ = a.IsAboveOrEqual(b)
	_ = a.IsLowerOrEqual(b)
	if !a.IsAnyOf(a, b) || a.IsAnyOf(b) {
		t.Error("IsAnyOf broken")
	}
	if !a.IsAnyValuesEqual(a.ValueInt8(), b.ValueInt8()) {
		t.Error("IsAnyValuesEqual broken")
	}
	if !a.IsAnyNamesOf(a.Name(), b.Name()) {
		t.Error("IsAnyNamesOf broken")
	}

	if got := a.Format("%s"); got == "" {
		t.Error("Format empty")
	}
	if a.ToEnumString(a.ValueInt8()) == "" || a.ToInt8EnumString(a.ValueInt8()) == "" {
		t.Error("ToEnumString broken")
	}
	if err := a.OnlySupportedErr(b.Name()); err == nil {
		t.Error("OnlySupportedErr should be non-nil")
	}
	if err := a.OnlySupportedMsgErr("ctx", b.Name()); err == nil {
		t.Error("OnlySupportedMsgErr should be non-nil")
	}

	_ = Min()
	_ = Max()

	jr := a.Json()
	var got Variant
	if err := got.JsonParseSelfInject(&jr); err != nil {
		t.Errorf("JsonParseSelfInject: %v", err)
	}
	if got != a {
		t.Errorf("JsonParseSelfInject round-trip: want %v, got %v", a, got)
	}

	if _, err := a.UnmarshallEnumToValue(jr.SecondaryData); err != nil {
		t.Logf("UnmarshallEnumToValue: %v", err)
	}
	if _, err := a.UnmarshallEnumToValueInt8(jr.SecondaryData); err != nil {
		t.Logf("UnmarshallEnumToValueInt8: %v", err)
	}
}
