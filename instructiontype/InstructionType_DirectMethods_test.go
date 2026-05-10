package instructiontype

import "testing"

// TestDirectMethods_InstructionType covers top-level helpers and methods
// the reflective Uplift sweep cannot reach (BA: per-package coverage uplift).
func TestDirectMethods_InstructionType(t *testing.T) {
	a, b := Scoping, Nginx

	if !a.IsAnyOf(a, b) || a.IsAnyOf(b) {
		t.Error("IsAnyOf broken")
	}
	if !a.IsNameOf(a.Name()) || a.IsNameOf("__nope__") {
		t.Error("IsNameOf broken")
	}
	if err := a.OnlySupportedErr(b.Name()); err == nil {
		t.Error("OnlySupportedErr should be non-nil")
	}
	if err := a.OnlySupportedMsgErr("ctx", b.Name()); err == nil {
		t.Error("OnlySupportedMsgErr should be non-nil")
	}

	_ = Min()
	_ = Max()
	_ = RangesInvalidErr()

	if v, err := New(a.Name()); err != nil || v != a {
		t.Errorf("New: got (%v,%v)", v, err)
	}
	if _, err := New("__bogus__"); err == nil {
		t.Error("New bogus should fail")
	}

	func() {
		defer func() { _ = recover() }()
		_ = NewMust(a.Name())
	}()
	func() {
		defer func() { _ = recover() }()
		_ = NewMust("__bogus__")
	}()

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

	jr := a.Json()
	var got Variant
	if err := got.JsonParseSelfInject(&jr); err != nil {
		t.Errorf("JsonParseSelfInject: %v", err)
	}
	if got != a {
		t.Errorf("round-trip: want %v got %v", a, got)
	}
}
