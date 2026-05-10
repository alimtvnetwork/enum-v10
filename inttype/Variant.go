package inttype

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/alimtvnetwork/core-v9/constants"
	"github.com/alimtvnetwork/core-v9/converters"
	"github.com/alimtvnetwork/core-v9/corecomparator"
	"github.com/alimtvnetwork/core-v9/corecsv"
	"github.com/alimtvnetwork/core-v9/coredata/corerange"
	"github.com/alimtvnetwork/core-v9/coreimpl/enumimpl"
	"github.com/alimtvnetwork/core-v9/coreimpl/enumimpl/enumtype"
	"github.com/alimtvnetwork/core-v9/coreinterface/enuminf"
	"github.com/alimtvnetwork/core-v9/errcore"
)

type Variant int

const (
	InvalidIndex Variant = -1
	Invalid      Variant = -1
	InvalidValue Variant = -1
	Zero         Variant = 0
	One          Variant = 1
	Two          Variant = 2
	Three        Variant = 3
	Min                  = Variant(constants.MinInt)
	Max                  = Variant(constants.MaxInt)
)

func (it Variant) ValueUInt16() uint16 {
	return uint16(it)
}

func (it Variant) AllNameValues() []string {
	return []string{
		Min.StringValue(),
		Invalid.StringValue(),
		Zero.StringValue(),
		Max.StringValue(),
	}
}

func (it Variant) OnlySupportedErr(names ...string) error {
	if len(names) == 0 {
		return errors.New("inttype: only supported names: <none specified>")
	}
	return fmt.Errorf("inttype: only supported names: %v", names)
}

func (it Variant) OnlySupportedMsgErr(message string, names ...string) error {
	return fmt.Errorf("inttype: %s — only supported names: %v", message, names)
}

func (it Variant) IntegerEnumRanges() []int {
	return []int{
		Min.Value(),
		Invalid.Value(),
		Zero.Value(),
		Max.Value(),
	}
}

func (it Variant) MinMaxAny() (min, max interface{}) {
	return Min, Max
}

func (it Variant) MinValueString() string {
	return Min.String()
}

func (it Variant) MaxValueString() string {
	return Max.String()
}

func (it Variant) MaxInt() int {
	return constants.MinInt
}

func (it Variant) MinInt() int {
	return constants.MinInt
}

func (it Variant) RangesDynamicMap() map[string]interface{} {
	return map[string]interface{}{}
}

func (it Variant) TypeName() string {
	return typeName
}

func (it Variant) Value() int {
	return int(it)
}

func (it Variant) StringValue() string {
	return strconv.Itoa(it.Value())
}

func (it Variant) ValueByte() byte {
	b, _ := corerange.Within.RangeByte(
		true,
		it.Value())

	return b
}

func (it Variant) ConvValueByte(isUseBoundaryDefaults bool) (byte byte, isInRange bool) {
	return corerange.Within.RangeByte(
		isUseBoundaryDefaults,
		it.Value())
}

func (it Variant) ConvValueByteWithBoundaryDefault() (byte byte, isInRange bool) {
	return corerange.Within.RangeByte(
		true,
		it.Value())
}

func (it Variant) String() string {
	return it.StringValue()
}

func (it Variant) Name() string {
	return it.StringValue()
}

// IsUninitialized
//
//	returns true if <= 0
func (it Variant) IsUninitialized() bool {
	return it <= Zero
}

func (it Variant) IsInitializedLogically() bool {
	return it.IsDefined()
}

func (it Variant) IsMin() bool {
	return it == Min
}

func (it Variant) IsAboveMin() bool {
	return it > Min
}

func (it Variant) IsAboveEqualMin() bool {
	return it >= Min
}

func (it Variant) IsMax() bool {
	return it == Max
}

func (it Variant) IsNotMin() bool {
	return it != Min
}

func (it Variant) IsNotMax() bool {
	return it != Max
}

func (it Variant) IsBetween(
	startIncluding, endIncluding Variant,
) bool {
	return it.IsBetweenInt(
		startIncluding.Value(),
		endIncluding.Value())
}

func (it Variant) IsBetweenInt(
	startIncluding, endIncluding int,
) bool {
	start := startIncluding
	end := endIncluding
	curr := it.Value()

	return curr >= start && curr <= end
}

func (it Variant) IsNotBetween(
	startIncluding, endIncluding Variant,
) bool {
	return !it.IsBetweenInt(
		startIncluding.Value(),
		endIncluding.Value())
}

func (it Variant) IsZero() bool {
	return it == Zero
}

// IsDefined
//
//	Greater than zero
func (it Variant) IsDefined() bool {
	return it > Zero
}

func (it Variant) IsGreaterThanZero() bool {
	return it > Zero
}

func (it Variant) IsOtherThanZero() bool {
	return it != Zero
}

func (it Variant) IsLessThanZero() bool {
	return it < Zero
}

func (it Variant) IsGreaterThanInvalid() bool {
	return it > Invalid
}

func (it Variant) HasValidIndex() bool {
	return it > Invalid
}

func (it Variant) HasValidValue() bool {
	return it > Invalid
}

func (it Variant) IsInvalidIndex() bool {
	return it == Invalid
}

func (it Variant) IsInvalidValue() bool {
	return it == Invalid
}

func (it Variant) IsInvalid() bool {
	return it == Invalid
}

func (it Variant) IsValid() bool {
	return it != Invalid
}

// IsPortRange
//
//	Refers to be under math.MaxUint16 and above Zero
func (it Variant) IsPortRange() bool {
	val := it.Value()

	return val >= 0 && val <= math.MaxUint16
}

// IsWithinRangeUint16
//
//	Refers to be under math.MaxUint16 and above Zero
func (it Variant) IsWithinRangeUint16() bool {
	val := it.Value()

	return val >= 0 && val <= math.MaxUint16
}

// IsWithinRangeByte
//
//	Refers to be under math.MaxUint16 and above Zero
func (it Variant) IsWithinRangeByte() bool {
	val := it.Value()

	return val >= 0 && val <= math.MaxUint8
}

func (it Variant) IsAnyOf(anyOfItems ...Variant) bool {
	for _, item := range anyOfItems {
		if item == it {
			return true
		}
	}

	return false
}

func (it Variant) IsNameOf(anyNames ...string) bool {
	for _, name := range anyNames {
		if name == it.StringValue() {
			return true
		}
	}

	return false
}

func (it Variant) IsNameOfValues(anyValues ...int) bool {
	for _, val := range anyValues {
		if val == it.Value() {
			return true
		}
	}

	return false
}

func (it Variant) AddStringAsNumber(n string) Variant {
	convInt, isSuccess := converters.StringTo.IntegerWithDefault(
		n, 0)

	if isSuccess {
		return Variant(it.Value() + convInt)
	}

	return it
}

// Add v + n
func (it Variant) Add(n int) Variant {
	return Variant(it.Value() + n)
}

// Subtract v - n
func (it Variant) Subtract(n int) Variant {
	return Variant(it.Value() - n)
}

func (it Variant) Is(n Variant) bool {
	return it.Value() == n.Value()
}

func (it Variant) IsDiff(n Variant) bool {
	return it.Value() != n.Value()
}

func (it Variant) IsCmp(
	compare corecomparator.Compare, n Variant,
) bool {
	switch compare {
	case corecomparator.Equal:
		return it.Is(n)
	case corecomparator.LeftLess:
		return it.IsLess(n.Value())
	case corecomparator.LeftLessEqual:
		return it.IsLessEqual(n.Value())
	case corecomparator.LeftGreater:
		return it.IsGreater(n.Value())
	case corecomparator.LeftGreaterEqual:
		return it.IsLessEqual(n.Value())
	case corecomparator.NotEqual:
		return it.IsDiff(n)
	}

	panic(errcore.OutOfRangeType.ErrorRefOnly(corecomparator.Inconclusive))
}

func (it Variant) IsEqual(n int) bool {
	return it.Value() == n
}

func (it Variant) IsNotEqual(n int) bool {
	return it.Value() != n
}

// IsGreater v.Value() > n
func (it Variant) IsGreater(n int) bool {
	return it.Value() > n
}

// IsGreaterEqual v.Value() >= n
func (it Variant) IsGreaterEqual(n int) bool {
	return it.Value() >= n
}

// IsLess v.Value() < n
func (it Variant) IsLess(n int) bool {
	return it.Value() < n
}

// IsLessEqual v.Value() <= n
func (it Variant) IsLessEqual(n int) bool {
	return it.Value() <= n
}

func (it Variant) NameValue() string {
	return enumimpl.NameWithValue(it)
}

func (it Variant) IsNameEqual(name string) bool {
	return it.Name() == name
}

func (it Variant) IsAnyNamesOf(names ...string) bool {
	for _, name := range names {
		if it.IsNameEqual(name) {
			return true
		}
	}

	return false
}

func (it Variant) ToNumberString() string {
	return strconv.Itoa(it.Value())
}

func (it Variant) ValueInt() int {
	return int(it)
}

func (it Variant) ValueInt8() int8 {
	return int8(it)
}

func (it Variant) ValueInt16() int16 {
	return int16(it)
}

func (it Variant) ValueInt32() int32 {
	return int32(it)
}

func (it Variant) ValueFloat() float32 {
	return float32(it)
}

func (it Variant) ValueFloat32() float32 {
	return float32(it)
}

func (it Variant) ValueFloat64() float64 {
	return float64(it)
}

func (it Variant) ValueString() string {
	return strconv.Itoa(
		it.Value())
}

func (it Variant) RangeNamesCsv() string {
	return corecsv.RangeNamesWithValuesIndexesCsvString(
		"Invalid",
		"Zero",
		"One",
		"Two",
		"Min",
		"Max")
}

func (it Variant) MarshalJSON() ([]byte, error) {
	return json.Marshal(it.Value())
}

func (it *Variant) UnmarshalJSON(
	data []byte,
) error {
	toInt, err := bytesToDeserializer.Integer(data)

	if err == nil {
		*it = Variant(toInt)
	}

	return err
}

func (it Variant) Format(format string) (compiled string) {
	return enumimpl.FormatUsingFmt(it, format)
}

func (it Variant) EnumType() enuminf.EnumTyper {
	return enumtype.Integer
}

func (it Variant) AsBasicEnumer() enuminf.BasicEnumer {
	return &it
}

func (it Variant) ToPtr() *Variant {
	return &it
}
