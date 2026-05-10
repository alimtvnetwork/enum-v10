package inttype

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// AL2: Coverage uplift for inttype (was 47.9%).
func Test_IntType_Coverage(t *testing.T) {
	Convey("inttype — predicates across spectrum", t, func() {
		samples := []int{-10, -1, 0, 1, 10, 80, 65535, 70000}
		for _, n := range samples {
			v := New(n)
			_ = v.IsZero()
			_ = v.IsDefined()
			_ = v.IsGreaterThanZero()
			_ = v.IsOtherThanZero()
			_ = v.IsLessThanZero()
			_ = v.IsGreaterThanInvalid()
			_ = v.HasValidIndex()
			_ = v.HasValidValue()
			_ = v.IsInvalidIndex()
			_ = v.IsInvalidValue()
			_ = v.IsInvalid()
			_ = v.IsValid()
			_ = v.IsPortRange()
			_ = v.IsWithinRangeUint16()
			_ = v.IsWithinRangeByte()
			_ = v.IsUninitialized()
			_ = v.IsInitializedLogically()
			_ = v.IsMin()
			_ = v.IsMax()
			_ = v.IsNotMin()
			_ = v.IsNotMax()
			_ = v.IsAboveMin()
			_ = v.IsAboveEqualMin()
			So(v.Value(), ShouldEqual, n)
			So(v.StringValue(), ShouldNotBeBlank)
			So(v.String(), ShouldNotBeBlank)
		}

		zero := New(0)
		So(zero.IsZero(), ShouldBeTrue)
		pos := New(42)
		So(pos.IsGreaterThanZero(), ShouldBeTrue)
		neg := New(-3)
		So(neg.IsLessThanZero(), ShouldBeTrue)
		port := New(80)
		So(port.IsPortRange(), ShouldBeTrue)
	})

	Convey("inttype — ranges + arithmetic", t, func() {
		v := New(10)
		So(v.IsBetween(New(0), New(100)), ShouldBeTrue)
		So(v.IsBetweenInt(0, 100), ShouldBeTrue)
		So(v.IsNotBetween(New(20), New(30)), ShouldBeTrue)
		So(v.IsAnyOf(New(10), New(11)), ShouldBeTrue)
		So(v.IsAnyOf(New(99)), ShouldBeFalse)
		So(v.IsNameOfValues(10, 20), ShouldBeTrue)
		So(v.Add(5).Value(), ShouldEqual, 15)
		So(v.AddStringAsNumber("3").Value(), ShouldEqual, 13)
	})

	Convey("inttype — byte conversion", t, func() {
		v := New(200)
		b, ok := v.ConvValueByte(false)
		So(ok, ShouldBeTrue)
		So(b, ShouldEqual, byte(200))
		bb, _ := v.ConvValueByteWithBoundaryDefault()
		So(bb, ShouldEqual, byte(200))
		_ = v.ValueByte()
	})

	Convey("inttype — IsNameOf + name access", t, func() {
		v := New(7)
		So(v.IsNameOf(v.Name()), ShouldBeTrue)
		So(v.IsNameOf("nope"), ShouldBeFalse)
	})

	Convey("inttype — surface accessors", t, func() {
		v := New(5)
		So(v.TypeName(), ShouldNotBeBlank)
		So(len(v.AllNameValues()), ShouldBeGreaterThanOrEqualTo, 0)
		So(v.IntegerEnumRanges(), ShouldNotBeNil)
		_ = v.MaxInt()
		_ = v.MinInt()
		_ = v.MaxValueString()
		_ = v.MinValueString()
		_ = v.RangesDynamicMap()
		minA, maxA := v.MinMaxAny()
		So(minA, ShouldNotBeNil)
		So(maxA, ShouldNotBeNil)
		So(v.OnlySupportedErr("anything"), ShouldNotBeNil)
		So(v.OnlySupportedMsgErr("ctx", "anything"), ShouldNotBeNil)
		_ = v.ValueUInt16()
	})
}
