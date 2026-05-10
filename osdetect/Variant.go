package osdetect

import (
	"github.com/alimtvnetwork/core-v9/coredata/corejson"
	"github.com/alimtvnetwork/core-v9/coreinterface/enuminf"
	"github.com/alimtvnetwork/enum-v10/linuxvendortype"
)

type Variant byte

const (
	Invalid Variant = iota
	AnyOs           // refers to Default, or All Os
	Windows
	Unix
	Linux
	MacOs
	Ubuntu
	Debian
	ArchLinux
	FreeBsd
	Centos
	RedHatEnterpriseLinux
	Docker
	Android
)

func (it Variant) ValueUInt16() uint16 {
	return uint16(it)
}

func (it Variant) AllNameValues() []string {
	return BasicEnumImpl.AllNameValues()
}

func (it Variant) OnlySupportedErr(names ...string) error {
	return BasicEnumImpl.OnlySupportedErr(names...)
}

func (it Variant) OnlySupportedMsgErr(message string, names ...string) error {
	return BasicEnumImpl.OnlySupportedMsgErr(message, names...)
}

func (it Variant) IntegerEnumRanges() []int {
	return BasicEnumImpl.IntegerEnumRanges()
}

func (it Variant) MinMaxAny() (min, max interface{}) {
	return BasicEnumImpl.MinMaxAny()
}

func (it Variant) MinValueString() string {
	return BasicEnumImpl.MinValueString()
}

func (it Variant) MaxValueString() string {
	return BasicEnumImpl.MaxValueString()
}

func (it Variant) MaxInt() int {
	return BasicEnumImpl.MaxInt()
}

func (it Variant) MinInt() int {
	return BasicEnumImpl.MinInt()
}

func (it Variant) RangesDynamicMap() map[string]interface{} {
	return BasicEnumImpl.RangesDynamicMap()
}

func (it Variant) DefaultCmdProcessName() string {
	if it.IsWindows() {
		return powershell
	}
	
	return bash
}

func (it Variant) IsAnyOs() bool {
	return it == AnyOs
}

func (it Variant) IsAllOs() bool {
	return it == AnyOs
}

func (it Variant) IsWindows() bool {
	return it == Windows
}

func (it Variant) IsUnix() bool {
	return it == Unix
}

func (it Variant) IsUnixLogically() bool {
	return it != Windows
}

func (it Variant) IsLinux() bool {
	return it == Linux
}

func (it Variant) IsRedHatEnterpriseLinux() bool {
	return it == RedHatEnterpriseLinux
}

func (it Variant) IsMacOs() bool {
	return it == MacOs
}

func (it Variant) IsUbuntu() bool {
	return it == Ubuntu
}

func (it Variant) IsCentos() bool {
	return it == Centos
}

func (it Variant) IsDebian() bool {
	return it == Debian
}

func (it Variant) IsArchLinux() bool {
	return it == ArchLinux
}

func (it Variant) IsDocker() bool {
	return it == Docker
}

func (it Variant) IsAndroid() bool {
	return it == Android
}

func (it Variant) IsAnyOsLogically() bool {
	return it != Invalid
}

func (it Variant) IsAllLogically() bool {
	return it != Invalid
}

func (it Variant) IsAnyNamesOf(names ...string) bool {
	return BasicEnumImpl.IsAnyNamesOf(it.ValueByte(), names...)
}

func (it Variant) ValueInt() int {
	return int(it)
}

func (it Variant) Index() int {
	return int(it)
}

func (it Variant) IsAnyValuesEqual(anyByteValues ...byte) bool {
	return BasicEnumImpl.IsAnyOf(it.ValueByte(), anyByteValues...)
}

func (it Variant) IsByteValueEqual(value byte) bool {
	return it.ValueByte() == value
}

func (it Variant) IsNameEqual(name string) bool {
	return it.Name() == name
}

func (it Variant) IsValueEqual(value byte) bool {
	return it.ValueByte() == value
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

func (it Variant) ValueString() string {
	return it.ToNumberString()
}

func (it Variant) Format(format string) (compiled string) {
	return BasicEnumImpl.Format(format, it.ValueByte())
}

func (it Variant) EnumType() enuminf.EnumTyper {
	return BasicEnumImpl.EnumType()
}

func (it Variant) IsUninitialized() bool {
	return it == Invalid
}

func (it Variant) IsInvalid() bool {
	return it == Invalid
}

func (it Variant) IsValid() bool {
	return it != Invalid
}

func (it Variant) Name() string {
	return BasicEnumImpl.ToEnumString(it.Value())
}

func (it Variant) ProductName() string {
	osDetail, err := GetCurrentOsDetail()
	
	if err != nil {
		return ""
	}
	
	return osDetail.Name.String()
}

func (it Variant) RawProductName() string {
	osDetail, err := GetCurrentOsDetail()
	
	if err != nil {
		return ""
	}
	
	return osDetail.ProductName.String()
}

func (it Variant) NameLower() string {
	return lowerCaseNames[it]
}

func (it Variant) ToNumberString() string {
	return BasicEnumImpl.ToNumberString(it.Value())
}

func (it Variant) UnmarshallEnumToValue(
	jsonUnmarshallingValue []byte,
) (byte, error) {
	return BasicEnumImpl.UnmarshallToValue(
		true,
		jsonUnmarshallingValue)
}

func (it Variant) MaxByte() byte {
	return BasicEnumImpl.Max()
}

func (it Variant) MinByte() byte {
	return BasicEnumImpl.Min()
}

func (it Variant) ValueByte() byte {
	return it.Value()
}

func (it Variant) RangesByte() []byte {
	return BasicEnumImpl.Ranges()
}

func (it Variant) Value() byte {
	return byte(it)
}

func (it *Variant) UnmarshalJSON(data []byte) error {
	dataConv, err := it.UnmarshallEnumToValue(data)
	
	if err == nil {
		*it = Variant(dataConv)
	}
	
	return err
}

func (it Variant) ToPtr() *Variant {
	return &it
}

func (it *Variant) ToSimple() Variant {
	if it == nil {
		return Invalid
	}
	
	return *it
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
		if name == it.Name() {
			return true
		}
	}
	
	return false
}

func (it Variant) MarshalJSON() ([]byte, error) {
	return BasicEnumImpl.ToEnumJsonBytes(it.Value())
}

func (it Variant) RangeNamesCsv() string {
	return BasicEnumImpl.RangeNamesCsv()
}

func (it Variant) TypeName() string {
	return BasicEnumImpl.TypeName()
}

func (it Variant) NameValue() string {
	return BasicEnumImpl.NameWithValue(it)
}

func (it Variant) String() string {
	return BasicEnumImpl.ToEnumString(it.Value())
}

func (it *Variant) JsonParseSelfInject(jsonResult *corejson.Result) error {
	err := jsonResult.Unmarshal(it)
	
	return err
}

func (it Variant) Json() corejson.Result {
	return corejson.New(it)
}

func (it Variant) JsonPtr() *corejson.Result {
	return corejson.NewPtr(it)
}

func (it *Variant) AsJsonContractsBinder() corejson.JsonContractsBinder {
	return it
}

func (it *Variant) AsJsoner() corejson.Jsoner {
	return it
}

func (it *Variant) AsJsonMarshaller() corejson.JsonMarshaller {
	return it
}

func (it Variant) AsBasicByteEnumContractsBinder() enuminf.BasicByteEnumContractsBinder {
	return &it
}

func (it Variant) AsBasicEnumContractsBinder() enuminf.BasicEnumContractsBinder {
	return &it
}

func (it Variant) AllSysMatchingTypes() []Variant {
	return CurrentOsMixTypes()
}

func (it Variant) AllSysMatchingTypesMap() map[Variant]bool {
	return CurrentOsTypesMap()
}

func (it Variant) IsCurrentOs() bool {
	allMap := it.AllSysMatchingTypesMap()
	
	return allMap[it]
}

func (it Variant) LinuxVendorType() linuxvendortype.Variant {
	osDetail, err := GetCurrentOsDetail()
	
	if err != nil || osDetail == nil {
		return linuxvendortype.Invalid
	}
	
	if osDetail.IsLinux {
		return osDetail.LinuxVendorType
	}
	
	return linuxvendortype.Invalid
}

func (it Variant) OsDetail() *OperatingSystemDetail {
	osDetail, err := GetCurrentOsDetail()
	
	if err != nil || osDetail == nil {
		return nil
	}
	
	return osDetail
}

func (it Variant) OsDetailWithError() (*OperatingSystemDetail, error) {
	return GetCurrentOsDetail()
}

// IsMajorAtLeast this will not yield right result for windows 11
func (it Variant) IsMajorAtLeast(majorVersion int) bool {
	osDetail, err := GetCurrentOsDetail()
	
	if err != nil {
		return false
	}
	
	return osDetail.IsTypePlusMajorAtLeast(
		it,
		majorVersion)
}

func (it Variant) IsWindows11() bool {
	osDetail, err := GetCurrentOsDetail()
	
	// Nil-guard: GetCurrentOsDetail may return (non-nil, nil) on a non-Windows
	// host where WindowsDetail is unpopulated. Without this guard the receiver
	// dereference panics with SIGSEGV (regression caught 2026-05-07).
	if err != nil ||
		osDetail == nil ||
		osDetail.IsEmptyWindowsDetail() {
		return false
	}
	
	return osDetail.IsWindows() &&
		osDetail.WindowsDetail.IsWindows11()
}

