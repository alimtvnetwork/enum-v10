package dbdrivertype

import "testing"

// TestDirectMethods_DbDriverType — final BA cycle. Covers Variant-level
// helpers and top-level functions the reflective Uplift sweep cannot reach
// (the Connection / compiler suite is in DbDriverType_Connection_Coverage_test.go).
func TestDirectMethods_DbDriverType(t *testing.T) {
	all := []Variant{
		Invalid, Sqlite, Redis, MySql, MariaDb, PostgreSql,
		MicrosoftSqlExpress, MicrosoftSqlServer, MicrosoftSqlCompact,
		MicrosoftAccess, Oracle, Firebird, MongoDb, CouchDb,
		AmazonDynamoDb, HSqlDb, Text, Json, Yaml, Protobuf,
	}

	// Sweep every variant through the DefaultPort family, IsSql / IsNoSql,
	// and the .Connection() factory. This trips the map-lookup branches
	// for both present and absent entries.
	for _, v := range all {
		_ = v.DefaultPort()
		_, _ = v.DefaultPortStatus()
		_, _ = v.DefaultPortStatusInteger()
		_ = v.DefaultPortInteger()
		_ = v.IsSqlDb()
		_ = v.IsNoSql()
		c := v.Connection()
		_ = c.Format()
		_ = c.AllDbFormat()
	}

	// Positive predicates
	if !MySql.IsSqlDb() || MySql.IsNoSql() {
		t.Error("MySql IsSqlDb/IsNoSql wrong")
	}
	if !MongoDb.IsNoSql() || MongoDb.IsSqlDb() {
		t.Error("MongoDb IsNoSql/IsSqlDb wrong")
	}
	if p, ok := MySql.DefaultPortStatus(); !ok || p != 3306 {
		t.Errorf("MySql DefaultPortStatus: %d %v", p, ok)
	}
	if _, ok := Invalid.DefaultPortStatus(); ok {
		t.Error("Invalid DefaultPortStatus should be !ok")
	}

	// Non-nullary methods
	a, b := MySql, PostgreSql
	if !a.IsAnyOf(a, b) || a.IsAnyOf(b) {
		t.Error("IsAnyOf broken")
	}
	if !a.IsNameOf(a.Name()) {
		t.Error("IsNameOf broken")
	}
	if !a.IsNameEqual(a.Name()) || a.IsNameEqual("__nope__") {
		t.Error("IsNameEqual broken")
	}
	if !a.IsValueEqual(a.ValueByte()) || a.IsValueEqual(255) {
		t.Error("IsValueEqual broken")
	}
	if !a.IsByteValueEqual(a.ValueByte()) {
		t.Error("IsByteValueEqual broken")
	}
	if !a.IsAnyValuesEqual(a.ValueByte(), b.ValueByte()) {
		t.Error("IsAnyValuesEqual broken")
	}
	if !a.IsAnyNamesOf(a.Name(), b.Name()) {
		t.Error("IsAnyNamesOf broken")
	}
	if a.Format("%s") == "" {
		t.Error("Format empty")
	}
	if err := a.OnlySupportedErr(b.Name()); err == nil {
		t.Error("OnlySupportedErr should be non-nil")
	}
	if err := a.OnlySupportedMsgErr("ctx", b.Name()); err == nil {
		t.Error("OnlySupportedMsgErr should be non-nil")
	}

	// Top-level helpers
	_ = Min()
	_ = Max()
	_ = RangesInvalidErr()
	if v, err := New(a.Name()); err != nil || v != a {
		t.Errorf("New: %v %v", v, err)
	}
	if _, err := New("__bogus__"); err == nil {
		t.Error("New bogus should fail")
	}
	func() {
		defer func() { _ = recover() }()
		_ = NewMust(a.Name())
	}()
	func() {
		defer func() { _ = recover() }()
		_ = NewMust("__bogus__")
	}()
	if !Is(a.Name(), a) || Is("__bogus__", a) {
		t.Error("Is helper broken")
	}
	if err := ValidationError("__bogus__", a); err == nil {
		t.Error("ValidationError should be non-nil for bogus input")
	}
	if err := ValidationError(a.Name(), a); err != nil {
		t.Errorf("ValidationError unexpected: %v", err)
	}
	StringMustBe(a.Name(), a)

	// UnmarshallEnumToValue happy path
	bs, _ := a.MarshalJSON()
	if _, err := a.UnmarshallEnumToValue(bs); err != nil {
		t.Errorf("UnmarshallEnumToValue: %v", err)
	}
}
