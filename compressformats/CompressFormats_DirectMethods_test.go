package compressformats

import "testing"

// TestDirectMethods_CompressFormats covers methods with non-nullary signatures
// that the reflective Uplift sweep cannot reach (BA: per-package coverage uplift).
func TestDirectMethods_CompressFormats(t *testing.T) {
	a, b := Zip, TarGZ

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

	if got := a.Format("%s"); got == "" {
		t.Error("Format empty")
	}
	if err := a.OnlySupportedErr(b.Name()); err == nil {
		t.Error("OnlySupportedErr should be non-nil")
	}
	if err := a.OnlySupportedMsgErr("ctx", b.Name()); err == nil {
		t.Error("OnlySupportedMsgErr should be non-nil")
	}

	_ = Ranges()
	_ = RangesInvalidErr()
	_ = Min()
	_ = Max()
	if !Is(a.Name(), a) || Is("__bogus__", a) {
		t.Error("Is helper broken")
	}
	if err := ValidationError("__bogus__", a); err == nil {
		t.Error("ValidationError should be non-nil for bogus input")
	}
	if err := ValidationError(a.Name(), a); err != nil {
		t.Errorf("ValidationError unexpected: %v", err)
	}
	StringMustBe(a.Name(), a)

	func() {
		defer func() { _ = recover() }()
		_ = NewMust(a.Name())
	}()
	func() {
		defer func() { _ = recover() }()
		_ = NewMust("__bogus__")
	}()

	jr := a.Json()
	var got Variant
	if err := got.JsonParseSelfInject(&jr); err != nil {
		t.Errorf("JsonParseSelfInject: %v", err)
	}
	if got != a {
		t.Errorf("JsonParseSelfInject round-trip: want %v, got %v", a, got)
	}
}
