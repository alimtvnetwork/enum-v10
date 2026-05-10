package httpstatusfamily

import "testing"

// TestDirectMethods_HttpStatusFamily covers methods with non-nullary signatures
// that the reflective Uplift sweep cannot reach (BA: per-package coverage uplift).
func TestDirectMethods_HttpStatusFamily(t *testing.T) {
	a, b := Successful, ClientError

	// Equality / comparison
	if !a.IsEqual(a) || a.IsEqual(b) {
		t.Error("IsEqual broken")
	}
	if !a.IsValueEqual(a.ValueByte()) || a.IsValueEqual(255) {
		t.Error("IsValueEqual broken")
	}
	if !a.IsByteValueEqual(a.ValueByte()) {
		t.Error("IsByteValueEqual broken")
	}
	if !a.IsNameEqual(a.Name()) || a.IsNameEqual("__nope__") {
		t.Error("IsNameEqual broken")
	}
	if !b.IsAboveOrEqual(a) || !a.IsLowerOrEqual(b) {
		t.Error("comparison broken")
	}
	if !a.IsAnyOf(a, b) || a.IsAnyOf(b) {
		t.Error("IsAnyOf broken")
	}
	if !a.IsAnyValuesEqual(a.ValueByte(), b.ValueByte()) {
		t.Error("IsAnyValuesEqual broken")
	}
	if !a.IsAnyNamesOf(a.Name(), b.Name()) {
		t.Error("IsAnyNamesOf broken")
	}

	// Diagnostic / formatting
	if got := a.Format("%s"); got == "" {
		t.Error("Format empty")
	}
	if err := a.OnlySupportedErr(b.Name()); err == nil {
		t.Error("OnlySupportedErr should be non-nil")
	}
	if err := a.OnlySupportedMsgErr("ctx", b.Name()); err == nil {
		t.Error("OnlySupportedMsgErr should be non-nil")
	}

	// Top-level helpers
	_ = Ranges()
	_ = RangesInvalidErr()
	if !Is(a.Name(), a) || Is("__bogus__", a) {
		t.Error("Is helper broken")
	}
	if err := ValidationError("__bogus__", a); err == nil {
		t.Error("ValidationError should be non-nil for bogus input")
	}
	if err := ValidationError(a.Name(), a); err != nil {
		t.Errorf("ValidationError unexpected: %v", err)
	}
	StringMustBe(a.Name(), a) // must not panic on the happy path

	// JsonParseSelfInject round-trip
	jr := a.Json()
	var got Variant
	if err := got.JsonParseSelfInject(&jr); err != nil {
		t.Errorf("JsonParseSelfInject: %v", err)
	}
	if got != a {
		t.Errorf("JsonParseSelfInject round-trip: want %v, got %v", a, got)
	}
}
