package osdetect

import (
	"github.com/alimtvnetwork/core-v9/corecomparator"
	"github.com/alimtvnetwork/core-v9/coredata/corejson"
	"github.com/alimtvnetwork/core-v9/coreversion"
	"github.com/alimtvnetwork/core-v9/ostype"
	"github.com/alimtvnetwork/enum-v9/linuxvendortype"
	"github.com/alimtvnetwork/enum-v9/osarchs"
	"github.com/alimtvnetwork/enum-v9/strtype"
)

// OperatingSystemDetail
//
// References:
// - Sample from other packages  : https://t.ly/Au3Q
// - Enum all os detect examples : github.com/alimtvnetwork/enum-v9/-/issues/4
type OperatingSystemDetail struct {
	OsMixType       Variant
	Name            strtype.Variant         `json:"Name,omitempty"`            // eg. "CentOS Linux 7 (Core)",
	ProductName     strtype.Variant         `json:"ProductName,omitempty"`     // eg. "CentOS Linux 7 (Core)",
	Vendor          strtype.Variant         `json:"Vendor,omitempty"`          // eg. ubuntu, centos
	LinuxVendorType linuxvendortype.Variant `json:"LinuxVendorType,omitempty"` // eg. ubuntu, centos
	Version         strtype.Variant         `json:"Version,omitempty"`         // eg. "7"
	Release         strtype.Variant         `json:"Release,omitempty"`         // eg. "7.2.1511"
	Architecture    osarchs.Architecture    `json:"Architecture,omitempty"`    // eg. "amd64"
	WindowsDetail   *WindowsSystemDetail    `json:"WindowsDetail,omitempty"`
	IsLinux         bool                    `json:"IsLinux,omitempty"`
	IsMacOs         bool                    `json:"IsMacOs,omitempty"`
	IsDocker        bool                    `json:"IsDocker,omitempty"` // TODO VM detect
	Group           ostype.Group
	releaseVersion  *coreversion.Version
}

func (it OperatingSystemDetail) AllSysTypes() []Variant {
	return CurrentOsMixTypes()
}

func (it OperatingSystemDetail) AllSysTypesMap() map[Variant]bool {
	return CurrentOsTypesMap()
}

func (it OperatingSystemDetail) IsName(name string) bool {
	if it.IsNull() {
		return false
	}
	
	return it.Name.IsEqual(name)
}

func (it OperatingSystemDetail) IsNameContains(name string) bool {
	if it.IsNull() {
		return false
	}
	
	return it.Name.IsContains(name)
}

func (it OperatingSystemDetail) IsNameStartsWith(name string) bool {
	if it.IsNull() {
		return false
	}
	
	return it.Name.IsStartsWith(name)
}

func (it OperatingSystemDetail) IsNameEndsWith(name string) bool {
	if it.IsNull() {
		return false
	}
	
	return it.Name.IsEndsWith(name)
}

func (it OperatingSystemDetail) IsArch(arch osarchs.Architecture) bool {
	if it.IsNull() {
		return false
	}
	
	return it.Architecture == arch
}

func (it OperatingSystemDetail) Is32BitArch() bool {
	return it.IsArch(osarchs.X32)
}

func (it OperatingSystemDetail) Is64BitArch() bool {
	return it.IsArch(osarchs.X64)
}

func (it OperatingSystemDetail) IsMajorVersionAtLeast(
	major int,
) bool {
	return it.
		ReleaseVersion().
		IsMajorAtLeast(major)
}

func (it OperatingSystemDetail) IsMajorVersion(
	major int,
) bool {
	comparison := it.
		ReleaseVersion().
		Major(major)
	
	return comparison.IsEqual()
}

func (it OperatingSystemDetail) IsVersion(
	versionCompare string,
) bool {
	return it.
		ReleaseVersion().
		IsExpectedComparisonRawVersion(
			corecomparator.Equal,
			versionCompare,
		)
}

func (it OperatingSystemDetail) IsVersionAtLeast(
	versionCompare string,
) bool {
	return it.ReleaseVersion().
		IsExpectedComparisonRawVersion(
			corecomparator.LeftGreaterEqual,
			versionCompare,
		)
}

func (it *OperatingSystemDetail) HasWindowsDetail() bool {
	return it != nil && it.WindowsDetail != nil
}

func (it *OperatingSystemDetail) IsEmptyWindowsDetail() bool {
	return it == nil || it.WindowsDetail == nil
}

func (it OperatingSystemDetail) IsWindows() bool {
	return it.Group.IsWindows()
}

func (it OperatingSystemDetail) IsUnix() bool {
	return it.Group.IsUnix()
}

func (it OperatingSystemDetail) IsAndroid() bool {
	return it.Group.IsAndroid()
}

func (it *OperatingSystemDetail) IsInvalid() bool {
	return it.IsEmpty() || it.Group.IsInvalid()
}

func (it *OperatingSystemDetail) IsValid() bool {
	return !it.IsInvalid()
}

func (it OperatingSystemDetail) IsUbuntu() bool {
	return it.OsMixType == Ubuntu
}

func (it OperatingSystemDetail) IsCentos() bool {
	return it.OsMixType == Centos
}

func (it OperatingSystemDetail) IsType(mixType Variant) bool {
	return it.OsMixType == mixType
}

func (it OperatingSystemDetail) IsTypePlusMajorAtLeast(
	mixType Variant,
	majorVersion int,
) bool {
	return it.OsMixType == mixType &&
		it.ReleaseVersion().
			IsMajorAtLeast(majorVersion)
}

func (it OperatingSystemDetail) IsUbuntuAtLeast(
	majorVersion int,
) bool {
	return it.IsTypePlusMajorAtLeast(
		Ubuntu, majorVersion)
}

func (it OperatingSystemDetail) IsCentOsAtLeast(
	majorVersion int,
) bool {
	return it.IsTypePlusMajorAtLeast(
		Centos, majorVersion)
}

func (it OperatingSystemDetail) IsDebianAtLeast(
	majorVersion int,
) bool {
	return it.IsTypePlusMajorAtLeast(
		Debian, majorVersion)
}

func (it OperatingSystemDetail) IsWindowsAtLeast(
	majorVersion int,
) bool {
	return it.IsTypePlusMajorAtLeast(
		Windows, majorVersion)
}

func (it OperatingSystemDetail) IsMacOsAtLeast(
	majorVersion int,
) bool {
	return it.IsTypePlusMajorAtLeast(
		MacOs, majorVersion)
}

func (it OperatingSystemDetail) IsAnyOfTypes(
	mixTypes ...Variant,
) bool {
	return it.OsMixType.IsAnyOf(mixTypes...)
}

func (it OperatingSystemDetail) IsTypePlusRunningInDocker(
	mixType Variant,
) bool {
	return it.OsMixType == mixType && it.IsDocker
}

func (it OperatingSystemDetail) Serialize() ([]byte, error) {
	return it.JsonPtr().Raw()
}

func (it OperatingSystemDetail) SerializeMust() []byte {
	return it.JsonPtr().RawMust()
}

func (it OperatingSystemDetail) Deserialize(toPtr interface{}) error {
	return it.JsonPtr().Deserialize(toPtr)
}

func (it OperatingSystemDetail) PrettyJsonString() string {
	return it.JsonPtr().PrettyJsonString()
}

func (it OperatingSystemDetail) Json() corejson.Result {
	return corejson.New(it)
}

func (it OperatingSystemDetail) JsonPtr() *corejson.Result {
	return corejson.NewPtr(it)
}

func (it *OperatingSystemDetail) IsNull() bool {
	return it == nil
}

func (it *OperatingSystemDetail) IsEmpty() bool {
	return it == nil ||
		it.Name == "" &&
			it.Vendor == "" &&
			it.Release == ""
}

func (it *OperatingSystemDetail) HasAnyItem() bool {
	return it == nil ||
		it.Name == "" &&
			it.Vendor == "" &&
			it.Release == ""
}

func (it *OperatingSystemDetail) HasWindowsDetails() bool {
	return it != nil &&
		it.WindowsDetail != nil
}

func (it *OperatingSystemDetail) ReleaseVersion() *coreversion.Version {
	if it == nil || it.Release.IsEmpty() {
		return nil
	}
	
	if it.releaseVersion != nil {
		return it.releaseVersion
	}
	
	it.releaseVersion = coreversion.New.DefaultPtr(
		it.Release.String())
	
	return it.releaseVersion
}
