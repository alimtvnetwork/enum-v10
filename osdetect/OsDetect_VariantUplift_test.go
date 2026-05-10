package osdetect_test

import (
	"testing"

	"github.com/alimtvnetwork/enum-v8/osdetect"
)

// AO pass-3 uplift: cover trivial accessors and JSON binders on Variant.
func TestOsDetect_VariantUplift(t *testing.T) {
	v := osdetect.Ubuntu
	_ = osdetect.Max()
	_ = v.RangeNamesCsv()
	_ = v.TypeName()
	_ = v.NameValue()
	_ = v.Json()
	_ = v.JsonPtr()
	ptr := &v
	_ = ptr.AsJsonContractsBinder()
	_ = ptr.AsJsoner()
	_ = ptr.AsJsonMarshaller()
	_ = ptr.AsBasicByteEnumContractsBinder()
	_ = v.AsBasicEnumContractsBinder()

	jr := v.Json()
	var v2 osdetect.Variant
	if err := v2.JsonParseSelfInject(&jr); err != nil {
		t.Errorf("self-inject: %v", err)
	}
	if v2 != v {
		t.Errorf("self-inject mismatch")
	}

	_ = v.AllSysMatchingTypes()
	_ = v.AllSysMatchingTypesMap()
	_ = v.IsCurrentOs()
	_ = v.LinuxVendorType()
	_ = osdetect.Windows.LinuxVendorType()
	_ = v.OsDetail()
	_ = v.IsMajorAtLeast(0)
}
