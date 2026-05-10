package osdetect_test

import (
	"encoding/json"
	"testing"

	"github.com/alimtvnetwork/enum-v9/inttype"
	"github.com/alimtvnetwork/enum-v9/osarchs"
	"github.com/alimtvnetwork/enum-v9/osdetect"
	"github.com/alimtvnetwork/enum-v9/strtype"
)

// AL2-08 follow-up uplift for osdetect using only `testing`. Targets
// every Variant accessor + JSON/contract surface, both `IsX` predicates
// for OperatingSystemDetail and WindowsSystemDetail across populated /
// empty / nil receivers, and the constructor / version-comparison paths.

func TestOsDetect_VariantFullSurface(t *testing.T) {
	all := []osdetect.Variant{
		osdetect.Invalid, osdetect.AnyOs, osdetect.Windows, osdetect.Unix,
		osdetect.Linux, osdetect.MacOs, osdetect.Ubuntu, osdetect.Debian,
		osdetect.ArchLinux, osdetect.FreeBsd, osdetect.Centos,
		osdetect.RedHatEnterpriseLinux, osdetect.Docker, osdetect.Android,
	}
	for _, v := range all {
		_ = v.Value()
		_ = v.ValueByte()
		_ = v.ValueInt()
		_ = v.ValueInt8()
		_ = v.ValueInt16()
		_ = v.ValueInt32()
		_ = v.ValueUInt16()
		_ = v.ValueString()
		_ = v.Index()
		_ = v.Name()
		_ = v.NameLower()
		_ = v.ProductName()
		_ = v.RawProductName()
		_ = v.ToNumberString()
		_ = v.Format("%s")
		_ = v.MaxByte()
		_ = v.MinByte()
		_ = v.MaxInt()
		_ = v.MinInt()
		_ = v.MinValueString()
		_ = v.MaxValueString()
		_ = v.AllNameValues()
		_ = v.IntegerEnumRanges()
		_ = v.RangesDynamicMap()
		_ = v.RangesByte()
		_, _ = v.MinMaxAny()
		_ = v.EnumType()
		_ = v.DefaultCmdProcessName()

		_ = v.IsValid()
		_ = v.IsInvalid()
		_ = v.IsUninitialized()
		_ = v.IsAnyOs()
		_ = v.IsAllOs()
		_ = v.IsAnyOsLogically()
		_ = v.IsAllLogically()
		_ = v.IsWindows()
		_ = v.IsUnix()
		_ = v.IsUnixLogically()
		_ = v.IsLinux()
		_ = v.IsRedHatEnterpriseLinux()
		_ = v.IsMacOs()
		_ = v.IsUbuntu()
		_ = v.IsCentos()
		_ = v.IsDebian()
		_ = v.IsArchLinux()
		_ = v.IsDocker()
		_ = v.IsAndroid()

		_ = v.IsAnyNamesOf("Linux", "Windows")
		_ = v.IsAnyValuesEqual(byte(osdetect.Linux), byte(osdetect.Windows))
		_ = v.IsByteValueEqual(v.ValueByte())
		_ = v.IsValueEqual(v.ValueByte())
		_ = v.IsNameEqual(v.Name())
		_ = v.IsAnyOf(osdetect.Linux, osdetect.Windows)
		_ = v.IsNameOf("Linux", "Windows")

		_ = v.OnlySupportedErr("Linux")
		_ = v.OnlySupportedMsgErr("ctx", "Linux")

		// JSON round-trip
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("Marshal(%v): %v", v, err)
		}
		var got osdetect.Variant
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("Unmarshal(%s): %v", string(data), err)
		}
	}

	if osdetect.Min().IsValid() && osdetect.Max().IsInvalid() {
		// just to touch Min/Max
		t.Logf("Min=%v Max=%v", osdetect.Min(), osdetect.Max())
	}
	if osdetect.RangesInvalidErr() == nil {
		t.Fatal("RangesInvalidErr nil")
	}

	// nil ToSimple already exercised; touch ToPtr round-trip
	p := osdetect.Ubuntu.ToPtr()
	if p == nil || p.ToSimple() != osdetect.Ubuntu {
		t.Fatal("ToPtr/ToSimple wrong")
	}
}

func TestOsDetect_OperatingSystemDetail_FullSurface(t *testing.T) {
	d := &osdetect.OperatingSystemDetail{
		OsMixType:    osdetect.Ubuntu,
		Name:         strtype.Variant("Ubuntu 22.04 LTS"),
		Vendor:       strtype.Variant("ubuntu"),
		Version:      strtype.Variant("22"),
		Release:      strtype.Variant("22.04"),
		Architecture: osarchs.X64,
		IsLinux:      true,
	}
	// Type / arch / name predicates
	_ = d.IsName("Ubuntu 22.04 LTS")
	_ = d.IsNameContains("22.04")
	_ = d.IsNameStartsWith("Ubuntu")
	_ = d.IsNameEndsWith("LTS")
	_ = d.IsArch(osarchs.X64)
	_ = d.Is32BitArch()
	_ = d.Is64BitArch()
	_ = d.IsType(osdetect.Ubuntu)
	_ = d.IsAnyOfTypes(osdetect.Ubuntu, osdetect.Debian)
	_ = d.IsTypePlusRunningInDocker(osdetect.Ubuntu)
	_ = d.IsUbuntu()
	_ = d.IsCentos()
	_ = d.IsWindows()
	_ = d.IsUnix()
	_ = d.IsAndroid()
	_ = d.IsValid()
	_ = d.IsInvalid()
	_ = d.HasWindowsDetail()
	_ = d.IsEmptyWindowsDetail()
	_ = d.HasWindowsDetails()
	_ = d.HasAnyItem()
	_ = d.IsEmpty()
	_ = d.IsNull()
	_ = d.AllSysTypes()
	_ = d.AllSysTypesMap()

	// Version comparison paths (ReleaseVersion populated)
	_ = d.IsMajorVersion(22)
	_ = d.IsMajorVersionAtLeast(20)
	_ = d.IsVersion("22.04")
	_ = d.IsVersionAtLeast("20.04")
	_ = d.IsUbuntuAtLeast(20)
	_ = d.IsCentOsAtLeast(7)
	_ = d.IsDebianAtLeast(11)
	_ = d.IsWindowsAtLeast(10)
	_ = d.IsMacOsAtLeast(11)
	_ = d.IsTypePlusMajorAtLeast(osdetect.Ubuntu, 20)

	// JSON / serialization
	if _, err := d.Serialize(); err != nil {
		t.Errorf("Serialize: %v", err)
	}
	_ = d.SerializeMust()
	_ = d.PrettyJsonString()
	_ = d.Json()
	if d.JsonPtr() == nil {
		t.Error("JsonPtr nil")
	}
	var dest osdetect.OperatingSystemDetail
	if err := d.Deserialize(&dest); err != nil {
		t.Errorf("Deserialize: %v", err)
	}

	// Empty Release path -> ReleaseVersion returns nil branch
	dEmpty := &osdetect.OperatingSystemDetail{OsMixType: osdetect.Linux}
	if dEmpty.ReleaseVersion() != nil {
		t.Error("empty Release should yield nil version")
	}

	// Nil receiver paths
	var n *osdetect.OperatingSystemDetail
	_ = n.IsNull()
	_ = n.IsEmpty()
	_ = n.HasAnyItem()
	_ = n.HasWindowsDetail()
	_ = n.IsEmptyWindowsDetail()
	_ = n.HasWindowsDetails()
	_ = n.IsValid()
	_ = n.IsInvalid()
	if n.ReleaseVersion() != nil {
		t.Error("nil ReleaseVersion should be nil")
	}
}

func TestOsDetect_WindowsSystemDetail_FullSurface(t *testing.T) {
	// Client side
	c := &osdetect.WindowsSystemDetail{
		WindowsVersion: inttype.Variant(11),
		IsClient:       true,
	}
	_ = c.IsNull()
	_ = c.IsDefined()
	_ = c.IsNullOr(false)
	_ = c.IsDefinedPlus(true)
	_ = c.IsWindows7()
	_ = c.IsWindows8()
	_ = c.IsWindows10()
	_ = c.IsWindows11()
	_ = c.IsWindows11OrAbove()
	_ = c.IsWindowsGreaterEqual(10)
	_ = c.IsWindowsEqual(11)
	_ = c.IsWindowsSever()
	_ = c.IsWindowsSever2016()
	_ = c.IsWindowsSever2019()
	_ = c.IsWindowsSeverGreaterEqual2016()
	_ = c.IsWindowsSeverGreaterEqual2019()
	_ = c.IsWindowsServerGreaterEqual(2016)
	_ = c.IsWindowsServerEqual(2019)
	_ = c.WinVer()

	// Server side
	s := &osdetect.WindowsSystemDetail{
		ServerVersion: inttype.Variant(2022),
		IsServer:      true,
	}
	_ = s.IsWindowsSever()
	_ = s.IsWindowsSeverGreaterEqual2016()
	_ = s.IsWindowsSeverGreaterEqual2019()
	_ = s.IsWindowsServerGreaterEqual(2016)
	_ = s.IsWindowsServerEqual(2022)
	_ = s.WinVer()

	// Null / undefined receiver — only methods marked with pointer receivers
	// can take nil safely.
	var n *osdetect.WindowsSystemDetail
	_ = n.IsNull()
	_ = n.IsDefined()
	_ = n.IsNullOr(true)
	_ = n.IsDefinedPlus(true)
	// WinVer dereferences fields, so it's not nil-safe — skip on nil receiver.
}

func TestOsDetect_OsDetailWithErr_FullSurface(t *testing.T) {
	o := &osdetect.OsDetailWithErr{
		OperatingSystemDetail: &osdetect.OperatingSystemDetail{
			OsMixType: osdetect.Linux,
			Name:      strtype.Variant("Linux"),
		},
	}
	if o.String() == "" || o.PrettyJsonString() == "" {
		t.Error("string helpers empty on populated detail")
	}
	if jr := o.Json(); jr.HasError() {
		t.Error("Json should not have error")
	}
	if o.JsonPtr() == nil {
		t.Error("JsonPtr nil")
	}
	// AsJsonContractsBinder: take by value
	_ = (*o).AsJsonContractsBinder()

	// JsonParseSelfInject round-trip
	jr := o.Json()
	var sink osdetect.OsDetailWithErr
	if err := sink.JsonParseSelfInject(&jr); err != nil {
		t.Errorf("JsonParseSelfInject: %v", err)
	}

	// Nil receiver -> empty strings, not panic
	var n *osdetect.OsDetailWithErr
	if n.String() != "" || n.PrettyJsonString() != "" {
		t.Error("nil receiver should return empty strings")
	}
}

func TestOsDetect_HostSmoke(t *testing.T) {
	// Just touch every host-dependent helper; never assert their values.
	_ = osdetect.IsWindows10()
	_ = osdetect.IsWindows11()
	_ = osdetect.IsRunningInDockerContainer()
	_ = osdetect.CurrentOsType()
	_ = osdetect.CurrentOsMixTypes()
	_ = osdetect.CurrentOsTypesMap()
	_ = osdetect.IsCurrentOsTypesContains(osdetect.Linux)
	_ = osdetect.IsCurrentOsTypesContains(osdetect.Invalid)
	_ = osdetect.CurrentOsTypesNotContainsError(osdetect.Linux)
	_, _ = osdetect.GetCurrentOsDetail()
}
