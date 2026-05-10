package pathpatterntype

import (
	"encoding/json"
	"testing"
)

// AO pass-3 uplift: cover JSON binders, ToSimple-nil, panic helpers,
// CompileTemplateReplaceOption / CompileCurlyTemplateReplace, findUsingInternalMapping.

func TestPathPatternType_Uplift(t *testing.T) {
	v := App
	ptr := v.ToPtr()
	if ptr.ToSimple() != v {
		t.Errorf("ToSimple roundtrip")
	}
	var nilPtr *Variant
	if nilPtr.ToSimple() != Invalid {
		t.Error("nil ToSimple should be Invalid")
	}

	_ = v.Json()
	_ = v.JsonPtr()
	_ = ptr.AsJsonContractsBinder()
	_ = ptr.AsJsoner()
	_ = ptr.AsJsonMarshaller()
	_ = v.AsBasicByteEnumContractsBinder()
	_ = v.AsBasicEnumContractsBinder()

	jr := v.Json()
	var v2 Variant
	if err := v2.JsonParseSelfInject(&jr); err != nil {
		t.Errorf("self-inject: %v", err)
	}
	if v2 != v {
		t.Errorf("self-inject mismatch")
	}

	// UnmarshalJSON
	b, _ := json.Marshal(v)
	var v3 Variant
	if err := json.Unmarshal(b, &v3); err != nil {
		t.Errorf("UnmarshalJSON: %v", err)
	}

	// ToSimple via valid pointer
	_ = ptr.ToSimple()

	// CompileTemplateReplaceOption / CompileCurlyTemplateReplace
	rep := map[string]string{"app": "myapp", "{app}": "myapp"}
	_ = PrefixAppRelativeIdFile.CompileTemplateReplaceOption(true, rep)
	_ = PrefixAppRelativeIdFile.CompileTemplateReplaceOption(false, rep)
	_ = PrefixAppRelativeIdFile.CompileCurlyTemplateReplace(rep)

	// findUsingInternalMapping (via New): empty + curly + bare + bogus
	if _, err := New(""); err == nil {
		t.Error("empty New should error")
	}
	if _, err := New("{app}"); err != nil {
		t.Errorf("curly New failed: %v", err)
	}
	if _, err := New("app"); err != nil {
		t.Errorf("bare New failed: %v", err)
	}
	if _, err := New("__bogus__"); err == nil {
		t.Error("bogus New should error")
	}

	// notFoundInListMessage path via panicOnNotFoundInList — wrap in recover.
	func() {
		defer func() { _ = recover() }()
		Variant(254).PathFullName()
	}()
	func() {
		defer func() { _ = recover() }()
		Variant(254).ExpandedAssociatedVariants()
	}()
}
