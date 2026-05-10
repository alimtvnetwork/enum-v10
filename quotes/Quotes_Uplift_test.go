package quotes

import (
	"testing"
)

// AO pass 2 — Quote enum accessor sweep (predicates / numeric / range / binders).

func TestQuotes_Uplift(t *testing.T) {
	q := Double

	if len(q.AllNameValues()) == 0 || len(q.IntegerEnumRanges()) == 0 {
		t.Error("ranges empty")
	}
	if len(q.RangesDynamicMap()) == 0 {
		t.Error("RangesDynamicMap empty")
	}
	if q.MinValueString() == "" || q.MaxValueString() == "" {
		t.Error("min/max strings")
	}
	if mn, mx := q.MinMaxAny(); mn == nil || mx == nil {
		t.Error("MinMaxAny")
	}

	_ = q.MaxInt()
	_ = q.MinInt()
	_ = q.ValueByte()
	_ = q.ValueUInt16()
	_ = q.ValueInt()

	if !q.IsValueEqual(q.Value()) {
		t.Error("IsValueEqual")
	}
	if !q.IsByteValueEqual(q.Value()) {
		t.Error("IsByteValueEqual")
	}
	if !q.IsAnyValuesEqual(q.Value()) {
		t.Error("IsAnyValuesEqual")
	}
	if !q.IsNameEqual(q.Name()) {
		t.Error("IsNameEqual")
	}
	if !q.IsAnyNamesOf(q.Name()) {
		t.Error("IsAnyNamesOf")
	}
	if q.IsEqual(0) {
		t.Error("IsEqual zero should false")
	}
	if !q.IsEqual(q.Value()) {
		t.Error("IsEqual self")
	}

	_ = q.OnlySupportedErr(q.Name())
	_ = q.OnlySupportedMsgErr("ctx", q.Name())

	_ = q.ValueInt8()
	_ = q.ValueInt16()
	_ = q.ValueInt32()
	_ = q.ValueString()
	_ = q.Format("%s")
	_ = q.EnumType()
	_ = q.NameValue()
	_ = q.ToNumberString()
	if q.IsInvalid() || !q.IsValid() {
		t.Error("validity")
	}
	if !q.IsAnyOf(q, Invalid) {
		t.Error("IsAnyOf")
	}
	_ = q.IsNameOf(q.Name())
	_ = q.MaxByte()
	_ = q.MinByte()
	_ = q.RangeNamesCsv()
	_ = q.TypeName()
	_ = q.AsBasicEnumContractsBinder()
	data, err := q.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON: %v", err)
	}
	var got Quote
	if err := got.UnmarshalJSON(data); err != nil {
		t.Errorf("UnmarshalJSON: %v", err)
	}
	if _, err := q.UnmarshallEnumToValue(data); err != nil {
		t.Errorf("UnmarshallEnumToValue: %v", err)
	}

	for _, candidate := range []Quote{Double, Single, Backtick} {
		_ = candidate.Name()
		_ = candidate.SelfWrap()
		_ = candidate.Wrap("hi")
		_ = candidate.WrapAny(123)
		_ = candidate.WrapString("hi")
		_ = candidate.IsWrapped("\"hi\"")
		_ = candidate.UnWrap("\"hi\"")
		_ = candidate.GetOther()
		_ = candidate.WrapSkipOnExist("\"hi\"")
		_ = candidate.WrapRegardless("hi")
		_ = candidate.WrapFmtString("prefix {wrapped} suffix", "hi")
		_ = candidate.WrapAnySkipOnExist(456)
	}
}
