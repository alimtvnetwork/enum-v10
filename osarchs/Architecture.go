package osarchs

import (
	"github.com/alimtvnetwork/core-v9/coreinterface/enuminf"
	"github.com/alimtvnetwork/enum-v9/internal/consts"
)

type Architecture byte

const (
	X32 Architecture = iota
	X64
	Invalid
)

func (it Architecture) ValueUInt16() uint16 {
	return uint16(it)
}

func (it Architecture) AllNameValues() []string {
	return BasicEnumImpl.AllNameValues()
}

func (it Architecture) OnlySupportedErr(names ...string) error {
	return BasicEnumImpl.OnlySupportedErr(names...)
}

func (it Architecture) OnlySupportedMsgErr(message string, names ...string) error {
	return BasicEnumImpl.OnlySupportedMsgErr(message, names...)
}

func (it Architecture) IntegerEnumRanges() []int {
	return BasicEnumImpl.IntegerEnumRanges()
}

func (it Architecture) MinMaxAny() (min, max interface{}) {
	return BasicEnumImpl.MinMaxAny()
}

func (it Architecture) MinValueString() string {
	return BasicEnumImpl.MinValueString()
}

func (it Architecture) MaxValueString() string {
	return BasicEnumImpl.MaxValueString()
}

func (it Architecture) MaxInt() int {
	return BasicEnumImpl.MaxInt()
}

func (it Architecture) MinInt() int {
	return BasicEnumImpl.MinInt()
}

func (it Architecture) RangesDynamicMap() map[string]interface{} {
	return BasicEnumImpl.RangesDynamicMap()
}

func (it Architecture) IsX32() bool {
	return it == X32
}

func (it Architecture) IsX64() bool {
	return it == X64
}

func (it Architecture) IsUnknown() bool {
	return it == Invalid
}

func (it Architecture) Value() byte {
	return byte(it)
}

func (it Architecture) ValueInt() int {
	return int(it)
}

func (it Architecture) ValueByte() byte {
	return byte(it)
}

func (it Architecture) IsAnyNamesOf(names ...string) bool {
	return BasicEnumImpl.IsAnyNamesOf(it.ValueByte(), names...)
}

func (it Architecture) IsAnyValuesEqual(anyByteValues ...byte) bool {
	return BasicEnumImpl.IsAnyOf(it.ValueByte(), anyByteValues...)
}

func (it Architecture) IsByteValueEqual(value byte) bool {
	return it.ValueByte() == value
}

func (it Architecture) IsNameEqual(name string) bool {
	return it.Name() == name
}

func (it Architecture) IsValueEqual(value byte) bool {
	return it.ValueByte() == value
}

func (it Architecture) ValueInt8() int8 {
	return int8(it)
}

func (it Architecture) ValueInt16() int16 {
	return int16(it)
}

func (it Architecture) ValueInt32() int32 {
	return int32(it)
}

func (it Architecture) ValueString() string {
	return it.ToNumberString()
}

func (it Architecture) Format(format string) (compiled string) {
	return BasicEnumImpl.Format(format, it.ValueByte())
}

func (it Architecture) EnumType() enuminf.EnumTyper {
	return BasicEnumImpl.EnumType()
}

func (it Architecture) Name() string {
	return architectures[it]
}

func (it Architecture) NameValue() string {
	return BasicEnumImpl.NameWithValue(it)
}

func (it *Architecture) ToNumberString() string {
	return BasicEnumImpl.ToNumberString(it.ValueByte())
}

func (it Architecture) IsInvalid() bool {
	return it == Invalid
}

func (it Architecture) IsValid() bool {
	return it != Invalid
}

func (it Architecture) IsAnyOf(anyOfItems ...Architecture) bool {
	for _, item := range anyOfItems {
		if item == it {
			return true
		}
	}

	return false
}

func (it Architecture) IsNameOf(anyNames ...string) bool {
	for _, name := range anyNames {
		if name == it.Name() {
			return true
		}
	}

	return false
}

func (it *Architecture) UnmarshallEnumToValue(
	jsonUnmarshallingValue []byte,
) (byte, error) {
	return BasicEnumImpl.UnmarshallToValue(
		consts.IsMappedToDefault,
		jsonUnmarshallingValue)
}

func (it Architecture) String() string {
	return architectures[it]
}

func (it Architecture) RangeNamesCsv() string {
	return BasicEnumImpl.RangeNamesCsv()
}

func (it Architecture) TypeName() string {
	return BasicEnumImpl.TypeName()
}

func (it Architecture) MarshalJSON() ([]byte, error) {
	return BasicEnumImpl.ToEnumJsonBytes(it.Value())
}

func (it *Architecture) UnmarshalJSON(data []byte) error {
	rawScriptType, err := it.UnmarshallEnumToValue(
		data)

	if err == nil {
		*it = Architecture(rawScriptType)
	}

	return err
}

func (it *Architecture) AsBasicEnumContractsBinder() enuminf.BasicEnumContractsBinder {
	return it
}

func (it Architecture) MaxByte() byte {
	// Invalid is the trailing sentinel; BasicEnumImpl.Max() would return it.
	// Return the highest real architecture (X64) instead. Mirrors Max().
	return byte(X64)
}

func (it Architecture) MinByte() byte {
	return byte(X32)
}

func (it Architecture) RangesByte() []byte {
	return BasicEnumImpl.Ranges()
}

func (it Architecture) AsBasicByteEnumContractsBinder() enuminf.BasicByteEnumContractsBinder {
	return &it
}

func (it Architecture) ToPtr() *Architecture {
	return &it
}
