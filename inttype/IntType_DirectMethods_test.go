package inttype

import (
	"encoding/json"
	"testing"

	"github.com/alimtvnetwork/core-v9/corecomparator"
)

// TestDirectMethods_IntType covers top-level constructors, helpers, and
// non-nullary methods the reflective Uplift sweep cannot reach
// (BA: per-package coverage uplift).
func TestDirectMethods_IntType(t *testing.T) {
	a := New(7)
	b := New(42)

	// Comparison & arithmetic
	if !a.IsEqual(7) || a.IsEqual(8) {
		t.Error("IsEqual broken")
	}
	if !a.IsNotEqual(8) || a.IsNotEqual(7) {
		t.Error("IsNotEqual broken")
	}
	if !a.Is(a) || a.Is(b) {
		t.Error("Is broken")
	}
	if !a.IsDiff(b) || a.IsDiff(a) {
		t.Error("IsDiff broken")
	}
	if !b.IsGreater(7) || b.IsGreater(100) {
		t.Error("IsGreater broken")
	}
	if !b.IsGreaterEqual(42) || !a.IsLess(42) || !a.IsLessEqual(7) {
		t.Error("IsGreaterEqual/IsLess/IsLessEqual broken")
	}
	if !a.IsBetweenInt(0, 10) || a.IsBetweenInt(20, 30) {
		t.Error("IsBetweenInt broken")
	}
	if !a.IsBetween(New(0), New(10)) || a.IsNotBetween(New(0), New(10)) {
		t.Error("IsBetween/IsNotBetween broken")
	}
	if got := a.Add(3); got.Value() != 10 {
		t.Errorf("Add: got %v want 10", got.Value())
	}
	if got := a.Subtract(2); got.Value() != 5 {
		t.Errorf("Subtract: got %v want 5", got.Value())
	}
	if got := a.AddStringAsNumber("3"); got.Value() != 10 {
		t.Errorf("AddStringAsNumber: got %v want 10", got.Value())
	}
	_ = a.AddStringAsNumber("not-a-number") // exercises error branch

	// Predicates / name-of variants
	if !a.IsAnyOf(a, b) || a.IsAnyOf(b) {
		t.Error("IsAnyOf broken")
	}
	if !a.IsNameOf(a.Name()) || a.IsNameOf("__nope__") {
		t.Error("IsNameOf broken")
	}
	if !a.IsNameOfValues(7, 8) || a.IsNameOfValues(99) {
		t.Error("IsNameOfValues broken")
	}
	if !a.IsNameEqual(a.Name()) {
		t.Error("IsNameEqual broken")
	}
	if !a.IsAnyNamesOf(a.Name(), b.Name()) {
		t.Error("IsAnyNamesOf broken")
	}

	// Format / surfaces
	if a.Format("%d") == "" {
		t.Error("Format empty")
	}
	if err := a.OnlySupportedErr(b.Name()); err == nil {
		t.Error("OnlySupportedErr should be non-nil")
	}
	if err := a.OnlySupportedMsgErr("ctx", b.Name()); err == nil {
		t.Error("OnlySupportedMsgErr should be non-nil")
	}

	// ConvValueByte (in-range + out-of-range)
	if _, ok := a.ConvValueByte(true); !ok {
		t.Error("ConvValueByte in-range should be ok")
	}
	if _, ok := New(9999).ConvValueByte(true); ok {
		t.Error("ConvValueByte out-of-range should be !ok")
	}

	// IsCompareResult — every branch
	for _, c := range []corecomparator.Compare{
		corecomparator.Equal, corecomparator.NotEqual,
		corecomparator.LeftGreater, corecomparator.LeftGreaterEqual,
		corecomparator.LeftLess, corecomparator.LeftLessEqual,
	} {
		_ = a.IsCompareResult(7, c)
	}

	// Top-level helpers
	if got := GetSet(true, a, b); got != a {
		t.Error("GetSet true wrong")
	}
	if got := GetSet(false, a, b); got != b {
		t.Error("GetSet false wrong")
	}
	if got := GetSetVariant(true, 7, 42); got.Value() != 7 {
		t.Error("GetSetVariant true wrong")
	}
	if got := GetSetVariant(false, 7, 42); got.Value() != 42 {
		t.Error("GetSetVariant false wrong")
	}

	// Constructors family
	if v, err := NewString("123"); err != nil || v.Value() != 123 {
		t.Errorf("NewString: %v %v", v, err)
	}
	if _, err := NewString("abc"); err == nil {
		t.Error("NewString bogus should fail")
	}
	if v, err := NewUInt(uint(5)); err != nil || v.Value() != 5 {
		t.Errorf("NewUInt: %v %v", v, err)
	}
	if v, err := NewInt64(int64(99)); err != nil || v.Value() != 99 {
		t.Errorf("NewInt64: %v %v", v, err)
	}
	if v, err := NewUsingStringer(a); err != nil || v.Value() != a.Value() {
		t.Errorf("NewUsingStringer: %v %v", v, err)
	}
	if v, err := NewUsingJsoner(a); err != nil || v.Value() != a.Value() {
		t.Errorf("NewUsingJsoner: %v %v", v, err)
	}
	jr := a.Json()
	if v, err := NewUsingJsonResult(&jr); err != nil || v.Value() != a.Value() {
		t.Errorf("NewUsingJsonResult: %v %v", v, err)
	}
	num := json.Number("42")
	if v, err := NewUsingJsonNumber(&num); err != nil || v.Value() != 42 {
		t.Errorf("NewUsingJsonNumber: %v %v", v, err)
	}
}
