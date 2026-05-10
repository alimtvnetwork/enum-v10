package packageinstallmethod

import "testing"

func TestPackageInstallMethod_Uplift(t *testing.T) {
	for _, v := range []Variant{Invalid, Url, OsPackages, AdvanceScript} {
		_ = v.IsInvalid()
		_ = v.IsValid()
		_ = v.IsUrl()
		_ = v.IsOsPackages()
		_ = v.IsAdvanceScript()
		ptr := v.ToPtr()
		if ptr.ToSimple() != v {
			t.Errorf("ToPtr/ToSimple roundtrip mismatch for %v", v)
		}
		_ = v.Json()
		_ = v.JsonPtr()
		_ = ptr.AsJsonContractsBinder()
		_ = ptr.AsJsoner()
		_ = ptr.AsJsonMarshaller()
		_ = ptr.AsBasicByteEnumContractsBinder()
		_ = v.AsBasicEnumContractsBinder()

		jr := v.Json()
		var v2 Variant
		if err := v2.JsonParseSelfInject(&jr); err != nil {
			t.Errorf("JsonParseSelfInject(%v): %v", v, err)
		}
		if v2 != v {
			t.Errorf("self-inject mismatch %v != %v", v2, v)
		}
	}
	var nilPtr *Variant
	if nilPtr.ToSimple() != Invalid {
		t.Error("nil ToSimple should be Invalid")
	}
}
