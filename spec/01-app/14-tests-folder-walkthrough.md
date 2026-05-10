# 14 — Tests Folder Walkthrough

> ✅ **Done** — new file (no analog in `spec/00-llm-integration-guide.md`).
> **Status**: filled in audit Step 4 (2026-04-23, Asia/Kuala_Lumpur).
> **Audience**: AI agents and contributors needing the **physical** layout of `tests/` and the public surface of in-tree test helpers (`coretests.GetAssert`, `tests/testwrappers/*`).

> **Consumer-coverage note (`enum-v9`)**: every layout, wrapper, and helper described on this page (`tests/integratedtests/`, `tests/testwrappers/`, `coretests.GetAssert`, `coretestcases.CaseV1`, `StringsTestWrapper`, etc.) refers to **upstream `core-v9`**. `enum-v9` does not consume any of them — `rg tests/testwrappers` and `rg coretests.GetAssert` over `enum-v9` source both return zero hits. This module's tests live at `tests/creationtests/` (one shared package, Goconvey-based registry over `EnumTestWrapper`); see [`13-testing-patterns.md` §6.1](./13-testing-patterns.md#61-enum-v9-specific-layout) for that file-by-file layout. Treat §§1–5 below as the authoritative reference for upstream `core-v9`.

For the **conceptual** style matrix (when to use `CaseV1` vs `BaseTestCase` etc.), see [`13-testing-patterns.md`](./13-testing-patterns.md).

---

## 1. `tests/creationtests/` *(upstream)* — One Folder per Source Package

> ⚠️ **Scope:** the layout below applies to **upstream `core-v9`**. `enum-v9` uses a single shared `tests/creationtests/` package — see [`13-testing-patterns.md` §6.1](./13-testing-patterns.md#61-enum-v9-specific-layout). Prior fixes that removed stale `tests/integratedtests/` references from the spec corpus: C-CVS-01 / D-CVS-17 / D-CVS-26 / D-CVS-27 / D-CVS-32 / D-CVS-36 — this section is the 7th occurrence (D-CVS-39).

```
tests/creationtests/
├── argstests/              # tests for coretests/args
├── anycmptests/            # tests for anycmp
├── bytetypetests/          # tests for bytetype
├── chmodhelpertests/
├── codestacktests/
├── conditionaltests/
├── corecmptests/
├── coredynamictests/
├── corejsontests/
├── corestrtests/
├── csvinternaltests/       # ← internal-package coverage (see issue 02-app-issues/02)
├── errcoretests/
├── isanytests/
├── jsoninternaltests/
├── regexnewtests/
├── ...                     # ≥ 50 packages total
└── GetAssert_*_test.go     # top-level Style D coverage of the GetAssert harness
```

### Naming rules

- **Folder name** = source package name + `tests` suffix. `corestr` → `corestrtests`.
- **Package name inside the folder** = same as folder. `package corestrtests` (NOT `corestrtests_test`).
- **`_test.go`** files contain runner functions only.
- **`_testcases.go`** files contain expectations only — no `import "testing"`.
- **`NilReceiver_test.go`** — at most one per package, runs all `CaseNilSafe` entries for that package's pointer-receiver methods.

Full naming spec: [`/spec/06-testing-guidelines/01-folder-structure.md`](../06-testing-guidelines/01-folder-structure.md).

---

## 2. `tests/testwrappers/` — Shared Wrapper Reference

Four wrapper packages live under `tests/testwrappers/`. Each provides a typed `BaseTestCase` extension that test packages alias via `type testWrapper = ...` (the Style B pattern).

> **Stability note**: the public surface of these wrappers is currently **undeclared** — see [`/spec/02-app-issues/04-testwrappers-public-surface.md`](../02-app-issues/04-testwrappers-public-surface.md). Treat them as stable for in-tree use; do not import from outside this module.

### 2.1 `stringstestwrapper`

| File | Purpose |
|---|---|
| `StringsTestWrapper.go` | Wraps `BaseTestCase` with `Arrange() []string` and `Expected() []string` accessors |

**API surface:**

```go
type StringsTestWrapper struct {
    coretests.BaseTestCase
}

func (it StringsTestWrapper) Arrange() []string  { /* ArrangeInput.([]string) */ }
func (it StringsTestWrapper) Expected() []string { /* ExpectedInput.([]string) */ }
```

**Used by**: `anycmptests`, `corecmptests`, and any package whose Style B output is a `[]string` line list.

### 2.2 `chmodhelpertestwrappers`

| File | Purpose |
|---|---|
| `RwxCompileValueTestWrapper.go` | Wraps cases for `chmodhelper.RwxCompile*` |
| `RwxInstructionTestWrapper.go` | Wraps cases for `chmodhelper.RwxInstruction*` |
| `VerifyRwxChmodUsingRwxInstructionsWrapper.go` | Wraps the verify-chmod scenarios |
| `VerifyRwxPartialChmodLocationsWrapper.go` | Wraps partial-chmod-location scenarios |
| `consts.go`, `vars.go`, `defaultValues.go`, `SimpleLocations.go` | Shared fixture data |
| `*TestCases.go` (multiple) | Shared case slices used by multiple `*_test.go` runners |

**Used by**: `chmodhelpertests`, `chmodclasstypetests`, `chmodinstests`.

### 2.3 `coredynamictestwrappers`

| File | Purpose |
|---|---|
| `ReflectSetFromToTestWrapper.go` | Wraps cases for `coredynamic.ReflectSetFromTo` |
| `ReflectSetFromToValidTestCases.go` | Shared valid-input cases |
| `ReflectSetFromToInvalidTestCases.go` | Shared invalid-input cases |
| `vars.go` | Shared fixture vars |

**Used by**: `coredynamictests`.

### 2.4 `corevalidatortestwrappers`

| File | Purpose |
|---|---|
| `SegmentValidatorWrapper.go` | Wraps cases for segment validators |
| `SliceValidatorWrapper.go` | Wraps cases for slice validators |
| `TextValidatorsWrapper.go` | Wraps cases for text validators |
| `SegmentValidatorTestCases.go`, `TextValidatorTestCases.go` | Shared case slices |

**Used by**: `corevalidatortests`.

---

## 3. `coretests.GetAssert` — Helper Inventory

`coretests.GetAssert` is a struct-as-namespace helper exposing **formatters** (turn raw values into stringified lines) and **shaping helpers** (sort, double-quote, error-line-format) that Style D tests call inside the Act phase.

> **Stability note**: not officially documented as a public API — see [`/spec/02-app-issues/03-getassert-undocumented-api.md`](../02-app-issues/03-getassert-undocumented-api.md). The methods below are **observed** from upstream `tests/creationtests/GetAssert_*_test.go` and `CoreTests_*_test.go` files (D-CVS-41).

### 3.1 Discovered methods (from coverage tests in this repo)

| Method | Signature shape | Purpose |
|---|---|---|
| `Quick` | `Quick(when, actual, expected string, counter int) string` | One-call AAA-formatted summary line |
| `SortedArray` | `SortedArray(isPrint, isSort bool, message string) []string` | Sort + format a string array, optionally print to stdout |
| `SortedArrayNoPrint` | `SortedArrayNoPrint(isSort bool, message string) []string` | Variant that never prints |
| `SortedMessage` | `SortedMessage(...) []string` | Sorted message lines |
| `ToString` | `ToString(any) string` | Convert any value to its diagnostic string form |
| `ToStrings` | `ToStrings(any) []string` | Multi-line variant |
| `ToStringsWithSpace` | `ToStringsWithSpace(any) []string` | With spacing/indent |
| `AnyToDoubleQuoteLines` | `AnyToDoubleQuoteLines(any) []string` | Each line wrapped in `"…"` |
| `AnyToStringDoubleQuoteLine` | `AnyToStringDoubleQuoteLine(any) string` | Single double-quoted line |
| `ConvertLinesToDoubleQuoteThenString` | `ConvertLinesToDoubleQuoteThenString([]string) string` | Lines → quoted-joined string |
| `ErrorToLinesWithSpaces` | `ErrorToLinesWithSpaces(error) []string` | Format an error as indented lines |
| `ErrorToLinesWithSpacesDefault` | `ErrorToLinesWithSpacesDefault(error) []string` | Default indent variant |
| `SimpleTestCaseWrapper` | `SimpleTestCaseWrapper(...)` | Wraps a simple case for batch invocation |

### 3.2 Canonical usage pattern

```go
asserter := coretests.GetAssert
formatter := asserter.SortedArray   // or .Quick, .ToStrings, etc.

outputs := formatter(arg1, arg2, arg3)
actualSlice.Adds(outputs...)

testCase.ShouldBeEqual(t, caseIndex, actualSlice.Strings()...)
```

The pattern is always: **bind the helper into a local variable** (improves readability and lets the IDE rename consistently).

### 3.3 When to add a new GetAssert helper

Add a method to `GetAssert` when:
1. The same formatting block appears in **3+ test packages**, AND
2. The output is **deterministic** (same inputs → same output across OSes), AND
3. The block is **purely formatting** (no production logic).

Otherwise keep the formatting inline in the test.

---

## 4. The `coretestcases.CaseV1(BaseTestCase)` Cast Idiom

In Style B, the runner needs to call assertion methods that live on `CaseV1`, not on `BaseTestCase`. Because `CaseV1` is structurally a `BaseTestCase` (it embeds it), Go allows a direct conversion:

```go
finalTestCase := coretestcases.CaseV1(testCase.BaseTestCase)
finalTestCase.ShouldBeEqual(t, caseIndex, finalActLines...)
```

This is **idiomatic** and **safe** — both types share the same memory layout. Do not "fix" this by introducing wrapper methods on `BaseTestCase`; that would duplicate the assertion API.

---

## 5. End-to-End Walkthrough — Adding `widgettests/`

Suppose you are adding a new public package `widget/` with three pointer-receiver methods on `*Widget` and one package-level function `BuildWidget`. The minimum test package looks like:

```
tests/creationtests/widgettests/    # ← upstream core-v9 layout; in enum-v9 register the enum in tests/creationtests/allBasicEnumsCollection.go instead
├── params.go                          # args.Map key constants (Style A)
├── BuildWidget_testcases.go           # Style A cases for BuildWidget
├── BuildWidget_test.go                # Style A runner
├── Widget_NilReceiver_testcases.go    # CaseNilSafe entries for *Widget methods
└── NilReceiver_test.go                # nil-safety runner
```

If `BuildWidget` returns a multi-line diagnostic string and you want to assert against sorted lines, add Style D in `BuildWidget_test.go`:

```go
asserter := coretests.GetAssert
sortedLines := asserter.SortedArray(false, true, output)
tc.ShouldBeEqual(t, caseIndex, sortedLines...)
```

If the input is a hetero typed slice (e.g. `[]args.TwoAny`), switch to Style B and add `testWrapper.go` + `testCases.go` per [`13-testing-patterns.md` §3](./13-testing-patterns.md#3-style-b--basetestcase--testwrapper-typed-slice--hetero-input).

---

## 6. Pitfalls (cross-link to `05-failing-tests/`)

The 25 failing-test post-mortems in [`/spec/05-failing-tests/`](../05-failing-tests/) are **mandatory skim reading**. The most relevant landmines for new test authors:

| Failing-tests file | Landmine |
|---|---|
| [`02-groupby-empty-map-assertion-mismatch.md`](../05-failing-tests/02-groupby-empty-map-assertion-mismatch.md) | Use `args.Map{}` for empty expectations, not `""` |
| [`12-nil-lazyregex-single-return-error.md`](../05-failing-tests/12-nil-lazyregex-single-return-error.md) | Nil-receiver methods may need single-error-return signature |
| [`13-slicevalidator-diagnostic-format-drift.md`](../05-failing-tests/13-slicevalidator-diagnostic-format-drift.md) | Diagnostic format is checked literally — no whitespace drift |
| [`18-zero-coverage-build-failure-cascade.md`](../05-failing-tests/18-zero-coverage-build-failure-cascade.md) | Compile failure in one test package cascades to coverage 0% — fix build first |
| [`22-crossed-packages-investigation.md`](../05-failing-tests/22-crossed-packages-investigation.md) | Watch for package-name collisions (e.g. two `package foo` in different folders) |
