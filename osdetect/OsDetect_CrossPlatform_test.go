package osdetect_test

import (
	"encoding/json"
	"testing"

	"github.com/alimtvnetwork/enum-v10/osdetect"
	. "github.com/smartystreets/goconvey/convey"
)

// AL-08: Cross-platform safe coverage for the osdetect package.
// Avoids any host-specific assertions (no IsWindows()/IsLinux() == true);
// only exercises pure logic, JSON round-trips, and smoke calls that work on
// any OS.
func TestOsDetect_CrossPlatformSafe(t *testing.T) {
	Convey("osdetect cross-platform safe surface", t, func() {

		Convey("Variant basic enum surface", func() {
			Convey("Invalid is uninitialized / invalid", func() {
				So(osdetect.Invalid.IsInvalid(), ShouldBeTrue)
				So(osdetect.Invalid.IsUninitialized(), ShouldBeTrue)
				So(osdetect.Invalid.IsValid(), ShouldBeFalse)
			})

			Convey("Known variants are valid and have non-empty names", func() {
				known := []osdetect.Variant{
					osdetect.AnyOs, osdetect.Windows, osdetect.Unix,
					osdetect.Linux, osdetect.MacOs, osdetect.Ubuntu,
					osdetect.Debian, osdetect.ArchLinux, osdetect.FreeBsd,
					osdetect.Centos, osdetect.RedHatEnterpriseLinux,
					osdetect.Docker, osdetect.Android,
				}
				for _, v := range known {
					So(v.IsValid(), ShouldBeTrue)
					So(v.IsInvalid(), ShouldBeFalse)
					So(v.Name(), ShouldNotBeBlank)
					So(v.NameLower(), ShouldNotBeBlank)
					So(v.ValueByte(), ShouldEqual, byte(v))
					So(v.ValueInt(), ShouldEqual, int(v))
				}
			})

			Convey("Type-predicate methods are mutually consistent", func() {
				So(osdetect.Windows.IsWindows(), ShouldBeTrue)
				So(osdetect.Linux.IsLinux(), ShouldBeTrue)
				So(osdetect.Unix.IsUnix(), ShouldBeTrue)
				So(osdetect.MacOs.IsMacOs(), ShouldBeTrue)
				So(osdetect.Ubuntu.IsUbuntu(), ShouldBeTrue)
				So(osdetect.Debian.IsDebian(), ShouldBeTrue)
				So(osdetect.Centos.IsCentos(), ShouldBeTrue)
				So(osdetect.ArchLinux.IsArchLinux(), ShouldBeTrue)
				So(osdetect.Docker.IsDocker(), ShouldBeTrue)
				So(osdetect.Android.IsAndroid(), ShouldBeTrue)
				So(osdetect.RedHatEnterpriseLinux.IsRedHatEnterpriseLinux(), ShouldBeTrue)
				So(osdetect.AnyOs.IsAnyOs(), ShouldBeTrue)
				So(osdetect.AnyOs.IsAllOs(), ShouldBeTrue)
			})

			Convey("Cross-predicate negatives", func() {
				So(osdetect.Windows.IsLinux(), ShouldBeFalse)
				So(osdetect.Linux.IsWindows(), ShouldBeFalse)
				So(osdetect.Windows.IsUnixLogically(), ShouldBeFalse)
				So(osdetect.Linux.IsUnixLogically(), ShouldBeTrue)
			})

			Convey("IsAnyOf and IsAnyValuesEqual", func() {
				So(osdetect.Linux.IsAnyOf(osdetect.Windows, osdetect.Linux), ShouldBeTrue)
				So(osdetect.Linux.IsAnyOf(osdetect.Windows, osdetect.MacOs), ShouldBeFalse)
				So(osdetect.Linux.IsAnyValuesEqual(byte(osdetect.Windows), byte(osdetect.Linux)), ShouldBeTrue)
				So(osdetect.Linux.IsByteValueEqual(byte(osdetect.Linux)), ShouldBeTrue)
				So(osdetect.Linux.IsValueEqual(byte(osdetect.Linux)), ShouldBeTrue)
				So(osdetect.Linux.IsNameEqual(osdetect.Linux.Name()), ShouldBeTrue)
			})

			Convey("Pointer round-trip", func() {
				v := osdetect.Ubuntu
				ptr := v.ToPtr()
				So(ptr, ShouldNotBeNil)
				So(ptr.ToSimple(), ShouldEqual, osdetect.Ubuntu)

				var nilPtr *osdetect.Variant
				So(nilPtr.ToSimple(), ShouldEqual, osdetect.Invalid)
			})

			Convey("DefaultCmdProcessName branches by OS family", func() {
				So(osdetect.Windows.DefaultCmdProcessName(), ShouldNotBeBlank)
				So(osdetect.Linux.DefaultCmdProcessName(), ShouldNotBeBlank)
				So(osdetect.Windows.DefaultCmdProcessName(),
					ShouldNotEqual,
					osdetect.Linux.DefaultCmdProcessName())
			})
		})

		Convey("New / NewMust constructors", func() {
			Convey("New round-trips a known name", func() {
				v, err := osdetect.New(osdetect.Linux.Name())
				So(err, ShouldBeNil)
				So(v, ShouldEqual, osdetect.Linux)
			})

			Convey("New rejects an unknown name", func() {
				_, err := osdetect.New("not-a-real-os-xyz")
				So(err, ShouldNotBeNil)
			})

			Convey("NewMust succeeds on a valid name", func() {
				So(func() {
					v := osdetect.NewMust(osdetect.Ubuntu.Name())
					So(v, ShouldEqual, osdetect.Ubuntu)
				}, ShouldNotPanic)
			})
		})

		Convey("Variant JSON round-trip", func() {
			v := osdetect.Ubuntu
			data, err := json.Marshal(v)
			So(err, ShouldBeNil)
			So(len(data), ShouldBeGreaterThan, 0)

			var decoded osdetect.Variant
			So(json.Unmarshal(data, &decoded), ShouldBeNil)
			So(decoded, ShouldEqual, osdetect.Ubuntu)
		})

		Convey("OperatingSystemDetail JSON helpers on empty value", func() {
			d := osdetect.OperatingSystemDetail{}
			So(d.IsNull(), ShouldBeFalse) // value receiver; pointer reports differently
			// Pretty/raw JSON should not panic on a zero value
			So(func() { _ = d.PrettyJsonString() }, ShouldNotPanic)
			So(func() { _ = d.Json() }, ShouldNotPanic)
		})

		Convey("CurrentOsType smoke (host-agnostic)", func() {
			// CurrentOsType() is the singular detected variant which on
			// minimally-provisioned Linux hosts may legitimately be Invalid.
			// Use the mix-types slice (always populated) for the IsValid probe.
			mix := osdetect.CurrentOsMixTypes()
			So(mix, ShouldNotBeNil)
			So(osdetect.CurrentOsTypesMap(), ShouldNotBeNil)
			if len(mix) > 0 {
				probe := mix[0]
				So(probe.IsValid(), ShouldBeTrue)
				So(probe.Name(), ShouldNotBeBlank)
				So(osdetect.IsCurrentOsTypesContains(probe), ShouldBeTrue)
			}
			So(osdetect.IsCurrentOsTypesContains(osdetect.Invalid), ShouldBeFalse)
		})

		Convey("IsRunningInDockerContainer is a stable bool", func() {
			// Just verify it returns without panic. The actual value is
			// host-dependent; calling it twice must be consistent.
			a := osdetect.IsRunningInDockerContainer()
			b := osdetect.IsRunningInDockerContainer()
			So(a, ShouldEqual, b)
		})
	})
}
