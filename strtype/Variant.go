package strtype

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/core-v9/bytetype"
	"github.com/alimtvnetwork/core-v9/constants"
	"github.com/alimtvnetwork/core-v9/converters"
	"github.com/alimtvnetwork/core-v9/corecsv"
	"github.com/alimtvnetwork/core-v9/coredata/corestr"
	"github.com/alimtvnetwork/core-v9/coredata/stringslice"
	"github.com/alimtvnetwork/core-v9/coreimpl/enumimpl"
	"github.com/alimtvnetwork/core-v9/coreimpl/enumimpl/enumtype"
	"github.com/alimtvnetwork/core-v9/coreinterface/enuminf"
	"github.com/alimtvnetwork/core-v9/coremath"
	"github.com/alimtvnetwork/core-v9/coreutils/stringutil"
	"github.com/alimtvnetwork/core-v9/coreversion"
	"github.com/alimtvnetwork/core-v9/errcore"
	"github.com/alimtvnetwork/enum-v10/inttype"
)

type Variant string

func (it Variant) ValueUInt16() uint16 {
	return 0
}

func (it Variant) AllNameValues() []string {
	return []string{}
}

func (it Variant) OnlySupportedErr(names ...string) error {
	panic("not implemented for generic string enum")
}

func (it Variant) OnlySupportedMsgErr(message string, names ...string) error {
	panic("not implemented for generic string enum")
}

func (it Variant) IntegerEnumRanges() []int {
	return []int{}
}

func (it Variant) MinMaxAny() (min, max interface{}) {
	return "", ""
}

func (it Variant) MinValueString() string {
	return ""
}

func (it Variant) MaxValueString() string {
	return ""
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

func (it Variant) NameValue() string {
	return string(it)
}

func (it Variant) FileReader() FileReader {
	return &fileReader{
		it.String(),
	}
}

func (it Variant) IsNameEqual(name string) bool {
	return it.String() == name
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
	return it.String()
}

func (it Variant) ToByteUsingMap(
	givenMap map[string]byte,
) (val byte, isApplicable bool) {
	if len(givenMap) == 0 {
		return 0, false
	}

	val, has := givenMap[it.String()]

	if has {
		return val, has
	}

	return 0, false
}

func (it Variant) ToByteUsingMapValidationErr(
	givenMap map[string]byte,
) (val byte, err error) {
	if len(givenMap) == 0 {
		return 0, errcore.
			InvalidEmptyValueType.
			ErrorNoRefs("empty map given to convert string type to byte")
	}

	val, has := givenMap[it.String()]

	if has {
		return val, nil
	}

	return 0, errcore.
		InvalidEmptyValueType.
		Error(
			"map doesn't contain string to byte convert key. Key:"+it.String(),
			givenMap)
}

func (it Variant) ToIntUsingMapValidationErr(
	givenMap map[string]int,
) (val int, err error) {
	if len(givenMap) == 0 {
		return constants.InvalidIndex, errcore.
			InvalidEmptyValueType.
			ErrorNoRefs("empty map given to convert string type to int")
	}

	val, has := givenMap[it.String()]

	if has {
		return val, nil
	}

	return constants.InvalidIndex, errcore.
		InvalidEmptyValueType.
		Error(
			"map doesn't contain string to int convert key. Key:"+it.String(),
			givenMap)
}

func (it Variant) ValueInt() int {
	val, _ := converters.StringTo.IntegerWithDefault(
		it.String(),
		constants.InvalidIndex)

	return val
}

func (it Variant) ValueInt8() int8 {
	val, _ := converters.StringTo.IntegerWithDefault(
		it.String(),
		constants.InvalidIndex)

	if coremath.IsOutOfRange.Integer.ToInt8(val) {
		return constants.InvalidIndex
	}

	return int8(val)
}

func (it Variant) ValueInt16() int16 {
	val, _ := converters.StringTo.IntegerWithDefault(
		it.String(),
		constants.InvalidIndex)

	if coremath.IsOutOfRange.Integer.ToInt16(val) {
		return constants.InvalidIndex
	}

	return int16(val)
}

func (it Variant) ValueInt32() int32 {
	val, _ := converters.StringTo.IntegerWithDefault(
		it.String(),
		constants.InvalidIndex)

	if coremath.IsOutOfRange.Integer.ToInt32(val) {
		return constants.InvalidIndex
	}

	return int32(val)
}

func (it Variant) ValueString() string {
	return it.String()
}

func (it Variant) RangeNamesCsv() string {
	return it.String()
}

func (it Variant) Format(format string) (compiled string) {
	return enumimpl.FormatUsingFmt(it, format)
}

func (it Variant) EnumType() enuminf.EnumTyper {
	return enumtype.String
}

func (it Variant) TypeName() string {
	return typeName
}

func (it Variant) ValueByte() byte {
	b, err := converters.StringTo.Byte(it.String())
	errcore.MustBeEmpty(err)

	return b
}

func (it Variant) ByteType() (val bytetype.Variant, isValid bool) {
	b, err := converters.StringTo.Byte(it.String())

	if err != nil {
		return 0, false
	}

	return bytetype.New(b), true
}

func (it Variant) IsInvalid() bool {
	return invalidMaps[it.StringValue()]
}

func (it Variant) IsValid() bool {
	return !it.IsInvalid()
}

func (it Variant) Value() string {
	return string(it)
}

func (it Variant) Length() int {
	return len(string(it))
}

func (it Variant) Size() int {
	return len(string(it))
}

func (it Variant) Count() int {
	return len(string(it))
}

func (it Variant) RunesLength() (length int, allRunes []rune) {
	allRunes = []rune(it)

	return len(allRunes), allRunes
}

func (it Variant) AllChars() []byte {
	return []byte(it)
}

func (it Variant) AllRunes() []rune {
	return []rune(it)
}

func (it Variant) TitleQuotation(
	title string,
) string {
	return fmt.Sprintf(
		TitleQuotationWrapFormat,
		title,
		string(it))
}

func (it Variant) TitleCurly(
	title string,
) string {
	return fmt.Sprintf(
		TitleCurlyWrapFormat,
		title,
		string(it))
}

func (it Variant) TitleSquare(
	title string,
) string {
	return fmt.Sprintf(
		TitleBracketWrapFormat,
		title, string(it))
}

func (it Variant) TitleQuotationReferenceStrings(
	title string,
	csvItems ...string,
) string {
	return fmt.Sprintf(
		TitleValueQuotationParenthesisRefWrapReferenceFormat,
		title,
		string(it),
		corecsv.DefaultCsv(csvItems...))
}

func (it Variant) TitleQuotationRefs(
	title string,
	csvItems ...interface{},
) string {
	return fmt.Sprintf(
		TitleValueQuotationParenthesisRefWrapReferenceFormat,
		title,
		string(it),
		corecsv.DefaultAnyCsv(csvItems...))
}

func (it Variant) QuotationWrap() string {
	return fmt.Sprintf(
		constants.SprintDoubleQuoteFormat,
		string(it))
}

func (it Variant) CurlyWrap() string {
	return fmt.Sprintf(
		CurlyStringWrapFormat,
		string(it))
}

func (it Variant) SquareWrap() string {
	return fmt.Sprintf(
		BracketStringWrapFormat,
		string(it))
}

func (it Variant) StringValue() string {
	return string(it)
}

func (it Variant) IsEmpty() bool {
	return string(it) == ""
}

func (it Variant) IsDefined() bool {
	return string(it) != ""
}

func (it Variant) IsWhitespace() bool {
	return strings.TrimSpace(string(it)) != ""
}

func (it Variant) Trim() Variant {
	return Variant(strings.TrimSpace(string(it)))
}

func (it Variant) IsEqualTrim(right string) bool {
	return strings.TrimSpace(string(it)) !=
		strings.TrimSpace(right)
}

func (it Variant) Replace(
	oldText, newText string,
) Variant {
	replaced := strings.ReplaceAll(
		it.String(),
		oldText,
		newText)

	return Variant(replaced)
}

func (it Variant) ReplaceUsingMapCurly(
	replacingMap map[string]string,
) Variant {
	replaced := stringutil.ReplaceTemplate.UsingMapOptions(
		true,
		it.String(),
		replacingMap)

	return Variant(replaced)
}

func (it Variant) ReplaceUsingMapDirect(
	replacingMap map[string]string,
) Variant {
	replaced := stringutil.
		ReplaceTemplate.
		UsingMapOptions(
			false,
			it.String(),
			replacingMap)

	return Variant(replaced)
}

func (it Variant) ReplaceUsingMapOption(
	isWrapKeysWithCurly bool,
	replacingMap map[string]string,
) Variant {
	replaced := stringutil.
		ReplaceTemplate.
		UsingMapOptions(
			isWrapKeysWithCurly,
			it.String(),
			replacingMap)

	return Variant(replaced)
}

func (it Variant) Remove(
	removeText string,
) Variant {
	replaced := strings.ReplaceAll(
		it.String(),
		removeText,
		constants.EmptyString)

	return Variant(replaced)
}

func (it Variant) RemoveMany(
	removeTexts ...string,
) Variant {
	replaced := stringutil.RemoveMany(
		it.String(),
		removeTexts...)

	return Variant(replaced)
}

func (it Variant) RemoveManyBySplitting(
	splitsBy string,
	removeTexts ...string,
) []string {
	return stringutil.RemoveManyBySplitting(
		it.String(),
		splitsBy,
		removeTexts...)
}

func (it Variant) SplitBy(
	splitsBy string,
) []string {
	return strings.Split(
		it.String(),
		splitsBy)
}

func (it Variant) SplitKeyVal(
	splitsBy string,
) (key, val string) {
	return stringutil.SplitLeftRight(
		it.String(),
		splitsBy)
}

func (it Variant) SplitKeyValTrim(
	splitsBy string,
) (keyTrim, valTrim string) {
	return stringutil.SplitLeftRightTrimmed(
		it.String(),
		splitsBy)
}

func (it Variant) AddSuffixOnMissing(
	suffixAdd string,
) (compiled string) {
	if it.HasSuffix(suffixAdd) {
		return it.String()
	}

	return it.String() + suffixAdd
}

func (it Variant) AddPrefixOnMissing(
	prefix string,
) (compiled string) {
	if it.HasPrefix(prefix) {
		return it.String()
	}

	return prefix + it.String()
}

func (it Variant) SplitTrimmedNonEmpty(
	splitsBy string,
) []string {
	return stringslice.SplitTrimmedNonEmptyAll(
		it.String(),
		splitsBy)
}

func (it Variant) SplitByWhitespace() []string {
	return stringslice.SplitContentsByWhitespace(
		it.String())
}

func (it Variant) SimpleStringOnce(
	isInitialized bool,
) corestr.SimpleStringOnce {
	return corestr.
		New.
		SimpleStringOnce.
		Create(string(it), isInitialized)
}

func (it Variant) SafeSubStringEnd(
	end int,
) Variant {
	return it.SafeSubString(
		0, end)
}

func (it Variant) SafeSubStringStart(
	start int,
) Variant {
	return it.SafeSubString(
		start, it.Length())
}

func (it Variant) SafeSplit(
	midPoint int,
) (left, right Variant) {
	left = it.SafeSubStringEnd(midPoint)
	right = it.SafeSubStringStart(midPoint)

	return left, right
}

func (it Variant) SplitKeyValue(
	splitter string,
) (left, right string) {
	return stringutil.SplitLeftRight(
		it.String(),
		splitter,
	)
}

func (it Variant) SplitKeyValueTrim(
	splitter string,
) (left, right string) {
	return stringutil.SplitLeftRightTrimmed(
		it.String(),
		splitter,
	)
}

func (it Variant) SplitKeyValueAsType(
	splitter string,
) (left, right Variant) {
	l, r := stringutil.SplitLeftRight(
		it.String(),
		splitter,
	)

	return Variant(l), Variant(r)
}

func (it Variant) SafeSubString(
	start, end int,
) Variant {
	s := it.String()
	length := len(s)

	if s == "" || start > length {
		return ""
	}

	if start < 0 {
		start = 0
	}

	if length < end {
		end = length
	}

	return Variant(s[start:end])
}

func (it Variant) String() string {
	return string(it)
}

func (it Variant) ConvInteger() (int, error) {
	return converters.StringTo.Integer(it.String())
}

func (it Variant) Integer() int {
	return converters.StringTo.IntegerDefault(it.String())
}

func (it Variant) IntType() inttype.Variant {
	return inttype.Variant(it.Integer())
}

func (it Variant) Version() *coreversion.Version {
	return coreversion.New.DefaultPtr(it.String())
}

func (it Variant) IntegerDefaultVal(defaultVal int) (valueInt int, isSuccess bool) {
	return converters.StringTo.IntegerWithDefault(it.String(), defaultVal)
}

func (it Variant) Name() string {
	return string(it)
}

func (it Variant) AddAnother(n Variant) Variant {
	return Variant(it.Value() + n.Value())
}

func (it Variant) Append(n Variant) Variant {
	return Variant(it.Value() + n.Value())
}

func (it Variant) Prepend(n Variant) Variant {
	return Variant(n.Value() + it.Value())
}

func (it Variant) PrependString(n string) Variant {
	return Variant(n + it.Value())
}

func (it Variant) Join(elements ...Variant) Variant {
	slice := make([]string, len(elements))
	for i, element := range elements {
		slice[i] = element.String()
	}

	return it.JoinStrings(slice...)
}

func (it Variant) JoinStrings(elements ...string) Variant {
	return Variant(strings.Join(elements, it.String()))
}

func (it Variant) PrependStringIf(isPrepend bool, n string) Variant {
	if isPrepend {
		return Variant(n + it.Value())
	}

	return it
}

func (it Variant) AppendIf(isAppend bool, n Variant) Variant {
	if isAppend {
		return Variant(it.Value() + n.Value())
	}

	return it
}

func (it Variant) AppendStringIf(isAppend bool, n string) Variant {
	if isAppend {
		return Variant(it.Value() + n)
	}

	return it
}

func (it Variant) PrependIf(isPrepend bool, n Variant) Variant {
	if isPrepend {
		return Variant(n.Value() + it.Value())
	}

	return it
}

func (it Variant) HasAnyItem() bool {
	return it != ""
}

func (it Variant) OrEmpty(n Variant) bool {
	return it.IsEmpty() || n.IsEmpty()
}

func (it Variant) OrHasElement(n Variant) bool {
	return it.HasAnyItem() || n.HasAnyItem()
}

func (it Variant) AndHasElement(n Variant) bool {
	return it.HasAnyItem() && n.HasAnyItem()
}

func (it Variant) AndIsEmpty(n Variant) bool {
	return it.IsEmpty() && n.IsEmpty()
}

func (it Variant) ToErr() error {
	if it == "" {
		return nil
	}

	return errors.New(it.String())
}

func (it Variant) HasInAliasMap(givenMap map[string]Variant, checkingElement Variant) bool {
	_, has := givenMap[checkingElement.String()]

	return has
}

// Add v + n
func (it Variant) Add(n string) Variant {
	return Variant(it.Value() + n)
}

func (it Variant) Is(n Variant) bool {
	return it.Value() == n.Value()
}

func (it Variant) IsContains(n string) bool {
	return strings.Contains(it.Value(), n)
}

func (it Variant) IsStartsWith(n string) bool {
	return strings.HasPrefix(it.Value(), n)
}

func (it Variant) IsEndsWith(n string) bool {
	return strings.HasSuffix(it.Value(), n)
}

func (it Variant) HasPrefix(n string) bool {
	return strings.HasPrefix(it.Value(), n)
}

func (it Variant) HasSuffix(n string) bool {
	return strings.HasSuffix(it.Value(), n)
}

func (it Variant) Index(n string) int {
	return strings.Index(it.Value(), n)
}

func (it Variant) LastIndexOf(n string) int {
	return strings.LastIndex(it.Value(), n)
}

func (it Variant) IsEqual(n string) bool {
	return it.Value() == n
}

func (it Variant) IsEqualAnother(n Variant) bool {
	return it.Value() == n.Value()
}

// IsGreater v.Value() > n
func (it Variant) IsGreater(n string) bool {
	return it.Value() > n
}

// IsGreaterEqual v.Value() >= n
func (it Variant) IsGreaterEqual(n string) bool {
	return it.Value() >= n
}

// IsLess v.Value() < n
func (it Variant) IsLess(n string) bool {
	return it.Value() < n
}

// IsLessEqual v.Value() <= n
func (it Variant) IsLessEqual(n string) bool {
	return it.Value() <= n
}

func (it Variant) NameUsingMap(
	nameRanges map[Variant]string,
) string {
	return nameRanges[it]
}

func (it Variant) MarshalJSON() ([]byte, error) {
	return bytesSerializer(it.String())
}

func (it *Variant) UnmarshalJSON(data []byte) error {
	dataConv, err := stringDeserializer(data)

	if err == nil {
		*it = Variant(dataConv)
	}

	return err
}

func (it Variant) AsBasicEnumer() enuminf.BasicEnumer {
	return &it
}

func (it Variant) ToPtr() *Variant {
	return &it
}
