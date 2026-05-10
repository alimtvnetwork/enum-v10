package quotes

import (
	"strings"

	"github.com/alimtvnetwork/core-v9/converters"
	"github.com/alimtvnetwork/core-v9/coreinterface/enuminf"
	"github.com/alimtvnetwork/enum-v10/internal/consts"
)

type Quote byte

const (
	Invalid  Quote = iota
	Double   Quote = '"'
	Single   Quote = '\''
	Backtick Quote = '`'
)

func (it Quote) ValueUInt16() uint16 {
	return uint16(it)
}

func (it Quote) AllNameValues() []string {
	return BasicEnumImpl.AllNameValues()
}

func (it Quote) OnlySupportedErr(names ...string) error {
	return BasicEnumImpl.OnlySupportedErr(names...)
}

func (it Quote) OnlySupportedMsgErr(message string, names ...string) error {
	return BasicEnumImpl.OnlySupportedMsgErr(message, names...)
}

func (it Quote) IntegerEnumRanges() []int {
	return BasicEnumImpl.IntegerEnumRanges()
}

func (it Quote) MinMaxAny() (min, max interface{}) {
	return BasicEnumImpl.MinMaxAny()
}

func (it Quote) MinValueString() string {
	return BasicEnumImpl.MinValueString()
}

func (it Quote) MaxValueString() string {
	return BasicEnumImpl.MaxValueString()
}

func (it Quote) MaxInt() int {
	return BasicEnumImpl.MaxInt()
}

func (it Quote) MinInt() int {
	return BasicEnumImpl.MinInt()
}

func (it Quote) RangesDynamicMap() map[string]interface{} {
	return BasicEnumImpl.RangesDynamicMap()
}

func (it Quote) IsEqual(char uint8) bool {
	return it.Value() == char
}

func (it Quote) Value() byte {
	return uint8(it)
}

func (it Quote) ValueByte() byte {
	return uint8(it)
}

func (it Quote) SelfWrap() string {
	return selfWrap[it]
}

func (it Quote) Wrap(str string) string {
	if str == "" {
		return it.String() + it.String()
	}

	return it.String() +
		str +
		it.String()
}

func (it Quote) GetOther() Quote {
	other, _ := otherQuoteMaps[it]

	return other
}

func (it Quote) WrapAny(source interface{}) string {
	toString := converters.AnyTo.ValueString(source)

	return WrapWith(
		toString,
		it,
		false)
}

func (it Quote) WrapAnySkipOnExist(source interface{}) string {
	toString := converters.AnyTo.ValueString(source)

	return WrapWith(
		toString,
		it,
		true)
}

func (it Quote) WrapString(
	sourceString string,
) string {
	return WrapWith(
		sourceString,
		it,
		false)
}

// WrapFmtString
//
//	{wrapped} will be replaced in the
//	format by the wrapped string.
func (it Quote) WrapFmtString(
	format, sourceString string,
) string {
	wrapped := it.Wrap(sourceString)

	return strings.ReplaceAll(
		format,
		consts.WrappedFormat,
		wrapped)
}

func (it Quote) IsWrapped(
	source string,
) bool {
	return HasBothWrappedWith(
		source,
		it)
}

func (it Quote) UnWrap(
	source string,
) string {
	return UnWrapWith(
		source,
		it,
	)
}

func (it Quote) WrapWithOptions(
	isSkipOnExist bool,
	source string,
) string {
	return WrapWith(source, it, isSkipOnExist)
}

func (it Quote) WrapSkipOnExist(
	source string,
) string {
	return WrapWith(
		source,
		it,
		true)
}

func (it Quote) WrapRegardless(
	source string,
) string {
	return WrapWith(
		source,
		it,
		false)
}

func (it Quote) IsAnyNamesOf(names ...string) bool {
	return BasicEnumImpl.IsAnyNamesOf(it.ValueByte(), names...)
}

func (it Quote) ValueInt() int {
	return int(it)
}

func (it Quote) IsAnyValuesEqual(anyByteValues ...byte) bool {
	return BasicEnumImpl.IsAnyOf(it.ValueByte(), anyByteValues...)
}

func (it Quote) IsByteValueEqual(value byte) bool {
	return it.ValueByte() == value
}

func (it Quote) IsNameEqual(name string) bool {
	return it.Name() == name
}

func (it Quote) IsValueEqual(value byte) bool {
	return it.ValueByte() == value
}

func (it Quote) ValueInt8() int8 {
	return int8(it)
}

func (it Quote) ValueInt16() int16 {
	return int16(it)
}

func (it Quote) ValueInt32() int32 {
	return int32(it)
}

func (it Quote) ValueString() string {
	return it.ToNumberString()
}

func (it Quote) Format(format string) (compiled string) {
	return BasicEnumImpl.Format(format, it.ValueByte())
}

func (it Quote) EnumType() enuminf.EnumTyper {
	return BasicEnumImpl.EnumType()
}

func (it Quote) Name() string {
	return BasicEnumImpl.ToEnumString(it.Value())
}

func (it Quote) NameValue() string {
	return BasicEnumImpl.NameWithValue(it)
}

func (it *Quote) ToNumberString() string {
	return BasicEnumImpl.ToNumberString(it.ValueByte())
}

func (it Quote) IsInvalid() bool {
	return it == Invalid
}

func (it Quote) IsValid() bool {
	return it != Invalid
}

func (it Quote) IsAnyOf(anyOfItems ...Quote) bool {
	for _, item := range anyOfItems {
		if item == it {
			return true
		}
	}

	return false
}

func (it Quote) IsNameOf(anyNames ...string) bool {
	for _, name := range anyNames {
		if name == it.Name() {
			return true
		}
	}

	return false
}

func (it *Quote) UnmarshallEnumToValue(
	jsonUnmarshallingValue []byte,
) (byte, error) {
	return BasicEnumImpl.UnmarshallToValue(
		isMappedToDefault,
		jsonUnmarshallingValue)
}

func (it Quote) String() string {
	return BasicEnumImpl.ToEnumString(it.Value())
}

func (it Quote) RangeNamesCsv() string {
	return BasicEnumImpl.RangeNamesCsv()
}

func (it Quote) TypeName() string {
	return BasicEnumImpl.TypeName()
}

func (it Quote) MarshalJSON() ([]byte, error) {
	return BasicEnumImpl.ToEnumJsonBytes(it.ValueByte())
}

func (it *Quote) UnmarshalJSON(data []byte) error {
	rawScriptType, err := it.UnmarshallEnumToValue(
		data)

	if err == nil {
		*it = Quote(rawScriptType)
	}

	return err
}

func (it *Quote) AsBasicEnumContractsBinder() enuminf.BasicEnumContractsBinder {
	return it
}

func (it Quote) MaxByte() byte {
	return BasicEnumImpl.Max()
}

func (it Quote) MinByte() byte {
	return BasicEnumImpl.Min()
}

func (it Quote) RangesByte() []byte {
	return BasicEnumImpl.Ranges()
}

func (it Quote) AsBasicByteEnumContractsBinder() enuminf.BasicByteEnumContractsBinder {
	return &it
}

func (it Quote) ToPtr() *Quote {
	return &it
}
