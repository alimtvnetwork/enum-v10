# 13 — Testing Patterns

> ✅ **Done** — extracted from `spec/00-llm-integration-guide.md` §Testing Patterns + audit §4 (Style Inventory).
> **Status**: filled in audit Step 4 (2026-04-23, Asia/Kuala_Lumpur).
> **Audience**: AI agents and contributors writing tests for `core-v9` *(upstream)* or for `enum-v10` *(this module)*.

> **Consumer-coverage note (`enum-v10`)**: every code example on this page (`coretestcases.CaseV1`, `args.Map`, `coretests.BaseTestCase`, `testWrapper`, `coretests.GetAssert`) refers to **upstream `core-v9`** test infrastructure that is **not consumed by `enum-v10/`**. This module's own tests live at `tests/creationtests/` and use a different shape entirely — a Goconvey-based registry runner over `EnumTestWrapper` (see [`14-tests-folder-walkthrough.md`](./14-tests-folder-walkthrough.md) §"`tests/creationtests/` shape"). Treat the Styles A–D below as the authoritative reference for upstream `core-v9` tests; for `enum-v10` enums, mirror the `creationtests` registry pattern instead.

This page is the **architectural overview** of how tests are organized. For the full pattern reference, drill into [`/spec/06-testing-guidelines/`](../06-testing-guidelines/) (9 files). For folder-level structure, see [`14-tests-folder-walkthrough.md`](./14-tests-folder-walkthrough.md).

---

## 1. The Four Test Styles (Decision Matrix)

The codebase has **four equally-supported test styles**. Pick by data shape, not by package age.

| Style | When to use | Case type | Assertion style | Wrapper file? | Example package |
|---|---|---|---|---|---|
| **A** | Single-shape input → single-shape output. Most common. | `coretestcases.CaseV1` + `args.Map` | `tc.ShouldBeEqualMap(t, idx, actual)` | No | `errcoretests` |
| **B** | Hetero typed-slice input (e.g. `[]args.TwoAny`) → `[]string` output | `coretests.BaseTestCase` wrapped in a `testWrapper` slice (`[]testWrapper`) | Cast to `coretestcases.CaseV1(tc.BaseTestCase)` then `ShouldBeEqual(t, idx, lines...)` | Yes — `testWrapper.go` (often a type alias to a shared wrapper) | `anycmptests`, `corecmptests` |
| **C** | One-off micro-assertion, no loop, no shared cases | None — `args.Map` literal in the test body | `expected.ShouldBeEqual(t, idx, title, actual)` (called on the literal map) | No | `argstests` |
| **D** | Formatted-output verification using the in-tree assertion harness | Style A or B, plus `coretests.GetAssert.*` helpers in the Act phase | Same as A or B | No | upstream `core-v9` `tests/creationtests/GetAssert_*_test.go` (see D-CVS-37) |

> The existing tradeoffs and unresolved questions about Style B vs A migration live in [`/spec/02-app-issues/01-style-b-style-a-coexistence.md`](../02-app-issues/01-style-b-style-a-coexistence.md).

---

## 2. Style A — `CaseV1` + `args.Map` (workhorse)

Defined in [`/spec/06-testing-guidelines/02-test-case-types.md`](../06-testing-guidelines/02-test-case-types.md).

**Skeleton:**

```go
// _testcases.go
var validateEmailTestCases = []coretestcases.CaseV1{
    {
        Title: "ValidateEmail returns valid -- given well-formed email",
        ArrangeInput: args.Map{"email": "user@example.com"},
        ExpectedInput: args.Map{"isValid": true, "errorCount": 0},
    },
}

// _test.go
func Test_ValidateEmail_Verification(t *testing.T) {
    for caseIndex, tc := range validateEmailTestCases {
        // Arrange
        input := tc.ArrangeInput.(args.Map)
        email, _ := input.GetAsString("email")

        // Act
        result := validator.ValidateEmail(email)
        actual := args.Map{
            "isValid":    result.IsValid,
            "errorCount": len(result.Errors),
        }

        // Assert
        tc.ShouldBeEqualMap(t, caseIndex, actual)
    }
}
```

**Key rules.**
- One `_test.go` runner per public function.
- All test data lives in `_testcases.go` — never inline in `_test.go`.
- Always include `// Arrange`, `// Act`, `// Assert` comments.
- Title format: `"FuncName returns X -- given Y input"`.
- Native types in `args.Map` (use `true`, not `"true"`).

For **single-case** variants, use a named variable + `*First` methods (`ShouldBeEqualMapFirst`, `ShouldBeEqualFirst`).

---

## 3. Style B — `BaseTestCase` + `[]testWrapper` (typed-slice / hetero input)

Used when the **input is a slice of complex types** (`[]args.TwoAny`, `[]args.Map`, etc.) but the **output is a stringified line list** suitable for `ShouldBeEqual`. The wrapper centralises typed accessors so multiple `_test.go` files can share one `testCases.go` slice.

**Three required files** (in addition to `_test.go`):

### 3.1 `testWrapper.go` — type alias to a shared wrapper

```go
package anycmptests

import (
    "github.com/alimtvnetwork/core-v9/tests/testwrappers/stringstestwrapper"
)

type testWrapper = stringstestwrapper.StringsTestWrapper
```

> **Why a type alias?** It lets every package keep its own `testWrapper` identifier in `testCases.go` while pointing at a shared implementation under `tests/testwrappers/`. Swap implementations by changing one line.

### 3.2 `testCases.go` — slice of wrappers + optional `VerifyTypeOf`

```go
var (
    arrangeTypeVerification = &coretests.VerifyTypeOf{
        ArrangeInput:  reflect.TypeOf([]args.TwoAny{}),
        ActualInput:   reflect.TypeOf([]string{}),
        ExpectedInput: reflect.TypeOf([]string{}),
    }

    testCases = []testWrapper{
        {
            BaseTestCase: coretests.BaseTestCase{
                Title:        "left and right is null checking ...",
                ArrangeInput: []args.TwoAny{ /* ... */ },
                ExpectedInput: []string{ /* ... */ },
                IsEnable:     issetter.True,
            },
        },
    }
)
```

> `VerifyTypeOf` documents the expected type shapes — the runner uses it to fail fast on type mismatches.
> `IsEnable: issetter.True` is **required boilerplate** — `issetter.False` skips the case.

### 3.3 `*_test.go` — runner with the `CaseV1(tc.BaseTestCase)` cast idiom

```go
func Test_Cmp_Verification(t *testing.T) {
    for caseIndex, testCase := range testCases {
        // Arrange
        inputs := testCase.ArrangeInput.([]args.TwoAny)
        actualSlice := corestr.New.SimpleSlice.Cap(len(inputs))

        // Act
        for i, parameter := range inputs {
            actualSlice.AppendFmt(
                "%d : %s (%T, %T)",
                i,
                anycmp.Cmp(parameter.First, parameter.Second).String(),
                parameter.First,
                parameter.Second,
            )
        }
        finalActLines := actualSlice.Strings()

        // The cast idiom — recover CaseV1 to call its assertion methods
        finalTestCase := coretestcases.CaseV1(testCase.BaseTestCase)

        // Assert
        finalTestCase.ShouldBeEqual(t, caseIndex, finalActLines...)
    }
}
```

The cast `coretestcases.CaseV1(testCase.BaseTestCase)` works because `CaseV1` has the same field set as `BaseTestCase` (it embeds it). This is **idiomatic** in Style B and not a code smell.

---

## 4. Style C — Standalone `args.Map.ShouldBeEqual` (micro-tests)

For one-off assertions where the test function is the case (no shared cases, no loop).

```go
func Test_Map_Length_FromMapLength(t *testing.T) {
    // Arrange
    m := args.Map{"a": 1, "b": 2}

    // Act
    actual := args.Map{"length": m.Length()}

    // Assert
    expected := args.Map{"length": 2}
    expected.ShouldBeEqual(t, 0, "Map.Length returns 2 -- two entries", actual)
}
```

Notes:
- The receiver `expected` carries the expected value; `actual` is the third positional arg.
- The third arg (a string title) replaces the `tc.Title` you'd have in Style A.
- Use this style when adding a case to a shared slice would be more boilerplate than the test itself (typically < 5 lines of arrange).

> See [`/spec/06-testing-guidelines/03-args-reference.md`](../06-testing-guidelines/03-args-reference.md) for the full `args.Map` API. (Style C documentation is being added in audit Step 8.)

---

## 5. Style D — `coretests.GetAssert.*` Formatted Output

When the actual output is a multi-line formatted string (sorted lists, error tables, double-quoted lines), the harness `coretests.GetAssert` provides shared formatters so the test only owns the input and the expected text.

```go
asserter := coretests.GetAssert
sortedArray := asserter.SortedArray

outputs := sortedArray(
    isPrint,            // bool — also write to stdout
    isSort,             // bool — sort before formatting
    message,            // string — header line
)
```

The full helper inventory lives in [`14-tests-folder-walkthrough.md` §3](./14-tests-folder-walkthrough.md#3-coretestsgetassert-helper-inventory).

---

## 6. Per-Package File Layout (canonical, **upstream `core-v9`**)

> ⚠️ **Scope:** the layout below applies to **upstream `core-v9`** packages. `enum-v10` does **not** have per-package `*tests/` folders — all enum contract tests live in the shared `tests/creationtests/` registry runner. See the consumer-coverage note at the top of this file and §6.1 below.

For a public package `foo/`, its tests live at `tests/creationtests/footests/` *(upstream)*:

```
footests/
├── params.go                          # (Style A) shared args.Map key constants
├── testCases.go                       # (Style B) shared []testWrapper
├── testWrapper.go                     # (Style B) type alias to a shared wrapper
├── helpers.go                         # shared test-only types/utilities
│
├── MyFunc_testcases.go                # (Style A) per-function test data
├── MyFunc_test.go                     # (Style A) per-function runner
│
├── MyStruct_NilReceiver_testcases.go  # nil-safety data
├── NilReceiver_test.go                # nil-safety runner (one per package)
│
└── MicroAssertion_test.go             # (Style C) standalone test functions
```

Folder-naming and import rules: [`/spec/06-testing-guidelines/01-folder-structure.md`](../06-testing-guidelines/01-folder-structure.md).
Prior fixes that removed stale `tests/integratedtests/` references from the spec corpus: C-CVS-01 / D-CVS-17 / D-CVS-26 / D-CVS-27 / D-CVS-32 — this section is the 6th and final occurrence (see D-CVS-36).

### 6.1 `enum-v10`-specific layout

```
tests/creationtests/
├── EnumTestWrapper.go                       # shared registry record (Header, InitialBasicEnumer, Expected*…)
├── allBasicEnumsCollection.go               # registry: every enum's `.Invalid` instance for is-valid checks
├── allEnumGeneralTestCases.go               # registry: per-enum expected MapValues / Min / Max / CSV
├── simpleEnumCollectionTestCases.go         # registry slice consumed by Test_Creation
├── pathPatternTypeCreationTestCases.go      # bespoke registry for pathpatterntype
├── generateAllBasicEnumTestCases.go         # populates the registries at init
├── generateAllEnumGeneralTestCases.go       #   ″
├── generatePathPatternTestCases.go          #   ″
├── creation_test.go                         # Goconvey runner over simpleEnumCollectionTestCases
├── AllEnums_ContractsTesting_test.go        # Goconvey runner over allEnumGeneralTestCases
├── PathType_Creation_test.go                # Goconvey runner over pathPatternTypeCreationTestCases
└── ScriptType_test.go                       # bespoke per-enum runner
```

Key differences vs the upstream layout above: (a) one shared package, not per-enum sub-packages; (b) Goconvey (`. "github.com/smartystreets/goconvey/convey"`) instead of `coretestcases`/`coretests`; (c) registry-driven — adding an enum means adding one entry, not creating a new folder; (d) no `args.Map`, no `testWrapper` type alias, no `params.go`.

---

## 7. When to Add a `params.go`

`params.go` centralises `args.Map` keys into typed constants so a typo becomes a compile error.

```go
var params = struct {
    pattern      string
    compareInput string
    isMatch      string
}{
    pattern:      "pattern",
    compareInput: "compareInput",
    isMatch:      "isMatch",
}
```

**Rule of thumb (provisional, see [`02-app-issues/05-missing-params-go-files.md`](../02-app-issues/05-missing-params-go-files.md)):**
- Mandatory when the package has > 3 test cases sharing keys.
- Optional otherwise.

---

## 8. Coverage Expectations

- Every public function: at least one positive case in `*_testcases.go`.
- Every public function with a non-trivial branch: at least one negative case.
- Every pointer-receiver method: a `NilReceiver_test.go` entry using `CaseNilSafe`.
- Boundary inputs (empty string, zero, nil) covered explicitly per [`/spec/06-testing-guidelines/06-branch-coverage.md`](../06-testing-guidelines/06-branch-coverage.md).
- Diagnostic output formatted per [`/spec/06-testing-guidelines/07-diagnostics-output-standards.md`](../06-testing-guidelines/07-diagnostics-output-standards.md).

---

## 9. Mandatory Reading Before Writing Tests

1. This page (decision matrix).
2. [`14-tests-folder-walkthrough.md`](./14-tests-folder-walkthrough.md) — testwrappers + GetAssert.
3. [`/spec/06-testing-guidelines/`](../06-testing-guidelines/) — the 9-file pattern reference.
4. [`/spec/05-failing-tests/`](../05-failing-tests/) — **skim every file** (these are the landmines you must not re-trigger).
