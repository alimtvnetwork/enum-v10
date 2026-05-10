# 02 — Test Case Types

## Overview

| Type | Use When | Input | Expected |
|------|----------|-------|----------|
| `CaseV1` (Style A) | Testing package functions or methods with explicit Act step | `ArrangeInput` | `ExpectedInput` |
| `BaseTestCase` + `testWrapper` (Style B) | Hetero typed-slice input → `[]string` output, shared across multiple `_test.go` files | `ArrangeInput` (typed slice) | `ExpectedInput` (`[]string`) |
| `args.Map` literal (Style C) | One-off micro-assertion, no loop, no shared cases | `args.Map` in test body | `args.Map` literal as receiver |
| `CaseNilSafe` | Testing nil-receiver safety of pointer receiver methods | `Func` (method ref) | `Expected` (ResultAny) |
| `GenericGherkins[TInput, TExpect]` | BDD-style scenarios with typed input/output | `Input` | `Expected` / `ExtraArgs` |

> Style D is not a separate case type — it is Style A or B **plus** the `coretests.GetAssert.*` formatter helpers in the Act phase. See [`/spec/01-app/14-tests-folder-walkthrough.md` §3](../01-app/14-tests-folder-walkthrough.md#3-coretestsgetassert-helper-inventory).

---

## Sub-Pattern: GoConvey-Only (Local Wrapper)

> **Scope:** This sub-pattern is a documented variation of Styles A/B for downstream consumers (e.g. `enum-v10`) that intentionally do NOT depend on the upstream `coretests` / `args` / `results` framework. It is a worked example, not a replacement.

Some downstream packages (notably `github.com/alimtvnetwork/enum-v10`) implement test cases using **GoConvey only** — `Convey` / `So` with assertions like `ShouldEqual`, `ShouldResemble`, `ShouldBeNil`, `ShouldBeTrue`, `ShouldBeEmpty` — combined with **local wrapper structs** (e.g. `EnumTestWrapper`, `PathPatternTypeCreationTestWrapper`) and module-level slice/map registries.

### When to use

- Downstream package whose test surface only needs construction + invariant assertions (no `args.Map`/`results.*` typed-slice ergonomics).
- You want zero dependency on `coretests` so the test build stays minimal.
- You still want AAA discipline (`// Arrange` / `// Act` / `// Assert` comments are mandatory — see §05).

### Worked example: `enum-v10/tests/creationtests/`

```go
// EnumTestWrapper.go (local — NOT from coretests)
type EnumTestWrapper struct {
    Title    string
    Enum     enuminf.BasicEnumer
    Expected string
}

// AllEnums_ContractsTesting_test.go
func Test_AllEnums_Contracts(t *testing.T) {
    Convey("All enums satisfy contract", t, func() {
        for _, tc := range allEnumGeneralTestCases {
            // Arrange
            actualEnumDynamicMap := dynamicEnumMapOf(tc.Enum)

            // Act
            diff := actualEnumDynamicMap.LogShouldDiffMessage(true, header, expected)

            // Assert
            So(diff, ShouldBeEmpty)
        }
    })
}
```

The diff-based assertion (`LogShouldDiffMessage` + `So(diff, ShouldBeEmpty)`) is the GoConvey-only equivalent of `ShouldBeEqualMap` — it reports the exact field mismatch on failure.

### Equivalence table

| Upstream framework (Styles A–D) | GoConvey-only sub-pattern |
|---|---|
| `CaseV1` struct | Local wrapper struct (e.g. `EnumTestWrapper`) |
| `coretests.GetAssert.ShouldBeEqualMap` | `LogShouldDiffMessage` + `So(diff, ShouldBeEmpty)` |
| `args.Map` semantic-key input | Module-level registry slice/map |
| `t.Run` sub-tests | `Convey` nested scopes |
| `tc.ShouldBeEqualFirst(t, ...)` | `So(actual, ShouldEqual, expected)` |

### Constraints

- AAA comments are still mandatory.
- This sub-pattern does NOT exempt a package from the framework-limitations section below if it later adopts `coretests`.
- Cross-link `spec/01-app/13-testing-patterns.md` §6.1 when introducing this in a new downstream consumer.

---

## CaseV1

The primary workhorse. Use for any test where you explicitly control Arrange → Act → Assert.

### Structure

```go
type CaseV1 struct {
    Title         string         // Test case name / scenario description
    ArrangeInput  any            // Input data (args.Map, args.One, etc.)
    ActualInput   any            // Set dynamically after Act phase
    ExpectedInput any            // Expected output (args.Map, string, []string, bool, etc.)
    VerifyTypeOf  *VerifyTypeOf  // Optional: type verification
    Parameters    *args.HolderAny // Optional: extra parameters
}
```

### Key Fields

| Field | Type | Purpose |
|-------|------|---------|
| `Title` | `string` | Displayed in test output on failure |
| `ArrangeInput` | `any` | Holds input data. Usually `args.Map` or positional types |
| `ExpectedInput` | `any` | Holds expected output. Must be `args.Map` for map assertions, `string`/`[]string` for line assertions |

### ExpectedInput Auto-Normalization

`ExpectedInput` is normalized to `[]string` via `ExpectedLines()`. Supported types:

| Type | Conversion |
|------|-----------|
| `string` | `[]string{s}` |
| `[]string` | as-is |
| `int` | `[]string{strconv.Itoa(v)}` |
| `bool` | `[]string{"true"}` or `[]string{"false"}` |
| `args.Map` | sorted `"key : value"` lines |
| other | `PrettyJSON` fallback |

### Assertion Methods

| Method | Use When |
|--------|----------|
| `ShouldBeEqual(t, caseIndex, actual...)` | Loop-based, exact string match |
| `ShouldBeEqualFirst(t, actual...)` | Single test case (caseIndex=0) |
| `ShouldBeEqualMap(t, caseIndex, actualMap)` | Map-based comparison |
| `ShouldBeEqualMapFirst(t, actualMap)` | Single test case, map comparison |
| `ShouldContains(t, caseIndex, actual...)` | Substring match |
| `ShouldStartsWith(t, caseIndex, actual...)` | Prefix match |
| `ShouldEndsWith(t, caseIndex, actual...)` | Suffix match |
| `ShouldBeNotEqual(t, caseIndex, actual...)` | Inverse match |
| `ShouldBeTrimEqual(t, caseIndex, actual...)` | Trimmed comparison |
| `ShouldBeSortedEqual(t, caseIndex, actual...)` | Sorted + trimmed comparison |
| `ShouldBeRegex(t, caseIndex, actual...)` | Regex match |

### Complete Example

**`_testcases.go`:**
```go
package mypkgtests

import (
    "myproject/coretests/args"
    "myproject/coretests/coretestcases"
)

// =============================================================================
// MyFunc positive path
// =============================================================================

var myFuncPositiveTestCases = []coretestcases.CaseV1{
    {
        Title: "MyFunc returns sum of two integers",
        ArrangeInput: args.Map{
            "a": 3,
            "b": 5,
        },
        ExpectedInput: args.Map{
            "result":  8,
            "isValid": true,
        },
    },
    {
        Title: "MyFunc handles zero values",
        ArrangeInput: args.Map{
            "a": 0,
            "b": 0,
        },
        ExpectedInput: args.Map{
            "result":  0,
            "isValid": true,
        },
    },
}

// =============================================================================
// MyFunc negative path — nil input
// =============================================================================

var myFuncNilInputTestCase = coretestcases.CaseV1{
    Title: "MyFunc with nil returns error",
    ArrangeInput: args.Map{
        "input": nil,
    },
    ExpectedInput: args.Map{
        "hasError": true,
        "result":   0,
    },
}
```

**`_test.go`:**
```go
package mypkgtests

import (
    "testing"

    "myproject/coretests/args"
    "myproject/mypkg"
)

// ==========================================
// MyFunc — positive path
// ==========================================

func Test_MyFunc_Positive_Verification(t *testing.T) {
    for caseIndex, tc := range myFuncPositiveTestCases {
        // Arrange
        input := tc.ArrangeInput.(args.Map)
        a, _ := input.GetAsInt("a")
        b, _ := input.GetAsInt("b")

        // Act
        result, err := mypkg.MyFunc(a, b)
        actual := args.Map{
            "result":  result,
            "isValid": err == nil,
        }

        // Assert
        tc.ShouldBeEqualMap(t, caseIndex, actual)
    }
}

// ==========================================
// MyFunc — negative path (nil input)
// ==========================================

func Test_MyFunc_NilInput(t *testing.T) {
    tc := myFuncNilInputTestCase

    // Arrange
    // (nil input — no setup needed)

    // Act
    result, err := mypkg.MyFunc(0, 0)
    actual := args.Map{
        "hasError": err != nil,
        "result":   result,
    }

    // Assert
    tc.ShouldBeEqualMapFirst(t, actual)
}
```

---

## Style B — `BaseTestCase` + `testWrapper`

Use when the **input is a typed slice** (`[]args.TwoAny`, `[]args.Map`, `[]MyStruct`, …) and the **output is a stringified line list** suitable for `ShouldBeEqual` (not `ShouldBeEqualMap`). Style B centralises typed accessors in a wrapper so multiple `_test.go` files in the same package can share one `testCases` slice.

### Structure

`BaseTestCase` lives in `coretests` and is field-compatible with `CaseV1`:

```go
// coretests
type BaseTestCase struct {
    Title         string
    ArrangeInput  any
    ActualInput   any
    ExpectedInput any
    VerifyTypeOf  *VerifyTypeOf
    Parameters    *args.HolderAny
    IsEnable      issetter.Value  // required: issetter.True to run, issetter.False to skip
}
```

A package-level `testWrapper` is a **type alias** to a shared wrapper under `tests/testwrappers/`:

```go
// testWrapper.go
package anycmptests

import "github.com/alimtvnetwork/core-v9/tests/testwrappers/stringstestwrapper"

type testWrapper = stringstestwrapper.StringsTestWrapper
```

The wrapper embeds `coretests.BaseTestCase` and adds typed accessors (e.g. `ArrangeAsTwoAnySlice()`).

### Required Files

| File | Holds |
|------|-------|
| `testWrapper.go` | One-line type alias to a shared wrapper |
| `testCases.go` | `testCases = []testWrapper{ ... }` plus optional `arrangeTypeVerification = &coretests.VerifyTypeOf{...}` |
| `*_test.go` | Runner that loops `testCases`, casts via the **`CaseV1(tc.BaseTestCase)` idiom**, then calls `ShouldBeEqual(t, idx, lines...)` |

### The CaseV1 Cast Idiom

Because `CaseV1` and `BaseTestCase` share the same field layout, you cast inside the loop to recover `CaseV1`'s assertion methods:

```go
finalTestCase := coretestcases.CaseV1(testCase.BaseTestCase)
finalTestCase.ShouldBeEqual(t, caseIndex, finalActLines...)
```

This is **idiomatic Style B**, not a code smell. Do not refactor it away.

> 🧠 **Why this cast is intentional** *(F-NEW-06 fix)*
> The cast keeps `BaseTestCase` a **pure data container** — decoupled from any specific assertion library — while still giving the runner access to `CaseV1`'s assertion methods at the call site. Refactoring it away (e.g. by adding `ShouldBeEqual` directly to `BaseTestCase`) would couple data and assertions together and break Style B's separation guarantee. **Keep the cast.**

### Required Boilerplate

- `IsEnable: issetter.True` on every case — `issetter.False` (or zero value) skips the case silently.
- `VerifyTypeOf` is recommended when `ArrangeInput` is a typed slice; the runner fails fast on shape drift.

### Complete Example

```go
// testWrapper.go
package anycmptests

import "github.com/alimtvnetwork/core-v9/tests/testwrappers/stringstestwrapper"

type testWrapper = stringstestwrapper.StringsTestWrapper
```

```go
// testCases.go
package anycmptests

import (
    "reflect"

    "github.com/alimtvnetwork/core-v9/coretests"
    "github.com/alimtvnetwork/core-v9/coretests/args"
    "github.com/alimtvnetwork/core-v9/issetter"
)

var (
    arrangeTypeVerification = &coretests.VerifyTypeOf{
        ArrangeInput:  reflect.TypeOf([]args.TwoAny{}),
        ActualInput:   reflect.TypeOf([]string{}),
        ExpectedInput: reflect.TypeOf([]string{}),
    }

    testCases = []testWrapper{
        {
            BaseTestCase: coretests.BaseTestCase{
                Title: "Cmp returns equal for matching primitives",
                ArrangeInput: []args.TwoAny{
                    {First: 1, Second: 1},
                    {First: "a", Second: "a"},
                },
                ExpectedInput: []string{
                    "0 : equal (int, int)",
                    "1 : equal (string, string)",
                },
                VerifyTypeOf: arrangeTypeVerification,
                IsEnable:     issetter.True,
            },
        },
    }
)
```

```go
// Cmp_test.go
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

        // Assert (cast idiom)
        finalTestCase := coretestcases.CaseV1(testCase.BaseTestCase)
        finalTestCase.ShouldBeEqual(t, caseIndex, finalActLines...)
    }
}
```

### When to Use Style B vs Style A

| Choose | If |
|--------|----|
| **A** (`CaseV1` + `args.Map`) | One input → one structured output; expected fields are semantic (`isValid`, `count`) |
| **B** (`BaseTestCase` + wrapper) | Input is a slice processed in a sub-loop; expected output is a list of formatted lines; same cases reused across multiple runners |

If both fit, prefer **A** — it is more self-documenting. Style B's tradeoffs are tracked in [`/spec/02-app-issues/01-style-b-style-a-coexistence.md`](../02-app-issues/01-style-b-style-a-coexistence.md).

### Pattern Abuse Warning

❌ Do **not** use Style B when the output is a structured map — use Style A's `ShouldBeEqualMap` instead.
❌ Do **not** inline `BaseTestCase` literals in `_test.go` — they belong in `testCases.go`.
❌ Do **not** drop `IsEnable: issetter.True` — the case will be silently skipped.

---

## CaseNilSafe

Designed exclusively for testing nil-receiver safety of **pointer receiver methods**.

### Structure

```go
type CaseNilSafe struct {
    Title         string          // Scenario name
    Func          any             // Direct method reference: (*Type).Method
    Args          []any           // Optional arguments for the method call
    Expected      results.ResultAny  // Expected outcome
    CompareFields []string        // Override auto-derived field comparison
}
```

### When to Use

✅ Use for: pointer receiver methods that must not panic on nil  
❌ Do NOT use for: package-level functions (use CaseV1 instead)

### How Func Works

The `Func` field accepts a **method expression** — a direct reference to a method:

```go
// Zero-arg method — use method expression directly:
Func: (*MyStruct).IsValid

// Method with arguments — wrap in a function literal:
Func: func(m *MyStruct) bool {
    return m.HasKey("someKey")
}

// Void method — wrap to suppress no-return:
Func: func(m *MyStruct) {
    m.SetName("x")
}
```

### Expected Fields (auto-derived)

| Field | Auto-compared when... | Meaning |
|-------|----------------------|---------|
| `Panicked` | always | Whether a panic occurred |
| `Value` | `Expected.Value != nil` | The stringified return value |
| `Error` | `Expected.Error != nil` | Whether an error was returned |
| `ReturnCount` | `Expected.ReturnCount != 0` | Number of return values |

### CompareFields Override

When auto-derivation isn't sufficient, explicitly specify which fields to compare:

```go
{
    Title: "SetName on nil does not panic",
    Func: func(m *MyStruct) {
        m.SetName("x")
    },
    Expected: results.ResultAny{
        Panicked: false,
    },
    // Void method has no "value" — only compare panicked + returnCount
    CompareFields: []string{"panicked", "returnCount"},
}
```

### Complete Example

**`_NilReceiver_testcases.go`:**
```go
package mypkgtests

import (
    "myproject/coretests/coretestcases"
    "myproject/coretests/results"
    "myproject/mypkg"
)

var myStructNilSafeTestCases = []coretestcases.CaseNilSafe{
    {
        Title: "IsValid on nil returns false",
        Func:  (*mypkg.MyStruct).IsValid,
        Expected: results.ResultAny{
            Value:    "false",
            Panicked: false,
        },
    },
    {
        Title: "Name on nil returns empty",
        Func:  (*mypkg.MyStruct).Name,
        Expected: results.ResultAny{
            Value:    "",
            Panicked: false,
        },
    },
    {
        Title: "HasKey on nil returns false",
        Func: func(m *mypkg.MyStruct) bool {
            return m.HasKey("anything")
        },
        Expected: results.ResultAny{
            Value:    "false",
            Panicked: false,
        },
    },
    {
        Title: "ClonePtr on nil returns nil",
        Func: func(m *mypkg.MyStruct) bool {
            return m.ClonePtr() == nil
        },
        Expected: results.ResultAny{
            Value:    "true",
            Panicked: false,
        },
    },
    {
        Title: "Clear on nil does not panic",
        Func:  (*mypkg.MyStruct).Clear,
        Expected: results.ResultAny{
            Panicked: false,
        },
        CompareFields: []string{"panicked", "returnCount"},
    },
}
```

**`NilReceiver_test.go`:**
```go
package mypkgtests

import "testing"

func Test_MyStruct_NilReceiver(t *testing.T) {
    for caseIndex, tc := range myStructNilSafeTestCases {
        tc.ShouldBeSafe(t, caseIndex)
    }
}
```

### Pattern Abuse Warning

**Never** use `CaseNilSafe` for package-level functions. If `ConcatMessageWithErr` is `func(string, error) error` (not a method), use `CaseV1`:

```go
// ❌ BAD — pattern abuse
var badTestCase = coretestcases.CaseNilSafe{
    Func: func(_ *struct{}) bool {
        return errcore.ConcatMessageWithErr("msg", nil) == nil
    },
}

// ✅ GOOD — use CaseV1
var goodTestCase = coretestcases.CaseV1{
    Title: "ConcatMessageWithErr nil error returns nil",
    ArrangeInput: args.Map{"message": "should not appear"},
    ExpectedInput: args.Map{"isNil": true},
}
```

---

## GenericGherkins[TInput, TExpect]

BDD-style test case with typed fields for input and expectations.

### Structure

```go
type GenericGherkins[TInput, TExpect any] struct {
    Title         string
    Feature       string
    Given         string
    When          string
    Then          string
    Input         TInput
    Expected      TExpect
    Actual        TExpect      // Set after Act
    IsMatching    bool
    ExpectedLines []string
    ExtraArgs     args.Map     // Overflow key-value pairs
}
```

### Common Aliases

```go
type AnyGherkins      = GenericGherkins[any, any]
type StringGherkins   = GenericGherkins[string, string]
type StringBoolGherkins = GenericGherkins[string, bool]
type MapGherkins      = GenericGherkins[args.Map, args.Map]
```

### When to Use Which Alias

| Alias | Input | Expected | Use When |
|-------|-------|----------|----------|
| `StringGherkins` | `string` | `string` | Single string input → single string result |
| `StringBoolGherkins` | `string` | `bool` | String input → boolean result (e.g., IsMatch) |
| `MapGherkins` | `args.Map` | `args.Map` | Multi-field input → multi-field result |
| `AnyGherkins` | `any` | `any` | Heterogeneous or unknown types |

### Field Responsibility

| Field | Purpose | Use When |
|-------|---------|----------|
| `Input` | Primary typed input data | Always — holds the main test input |
| `Expected` | Primary typed expected result | Always — holds what the test asserts against |
| `ExtraArgs` | Overflow key-value pairs | Only when `Input` cannot hold all arrange data |
| `ExpectedLines` | Legacy raw string assertions | **Deprecated for new tests** — use `Expected` with `args.Map` instead |
| `IsMatching` | Boolean match flag | Validation/matching tests (e.g., regex, search) |

### MapGherkins — Preferred for Multi-Field Tests

When a test has multiple inputs (e.g., pattern + compareInput) and multiple expected
results (e.g., isDefined, isApplicable, isMatch), use `MapGherkins`:

- **Input** (`args.Map`): holds all arrange data with semantic keys
- **Expected** (`args.Map`): holds all assertion data with semantic keys
- **ExtraArgs**: only used if needed beyond `Input`

This replaces the opaque `ExpectedLines: []string{"true", "false", "true"}` pattern
where each line's meaning is unknowable without reading the test runner.

#### Example — MapGherkins Test Case

```go
// params.go
package regexnewtests

var params = struct {
    pattern      string
    compareInput string
    isDefined    string
    isApplicable string
    isMatch      string
    isFailedMatch string
}{
    pattern:      "pattern",
    compareInput: "compareInput",
    isDefined:    "isDefined",
    isApplicable: "isApplicable",
    isMatch:      "isMatch",
    isFailedMatch: "isFailedMatch",
}
```

```go
// _testcases.go
var lazyRegexTestCases = []coretestcases.MapGherkins{
    {
        Title: "New.Lazy matches word pattern",
        When:  "given a simple word pattern",
        Input: args.Map{
            params.pattern:      "hello",
            params.compareInput: "hello world",
        },
        Expected: args.Map{
            params.isDefined:    true,
            params.isApplicable: true,
            params.isMatch:      true,
            params.isFailedMatch: false,
        },
    },
}
```

```go
// _test.go
func Test_LazyRegex_New_Verification(t *testing.T) {
    for caseIndex, tc := range lazyRegexTestCases {
        // Arrange
        pattern, _ := tc.Input.GetAsString(params.pattern)
        compareInput, _ := tc.Input.GetAsString(params.compareInput)

        // Act
        lazyRegex := regexnew.New.Lazy(pattern)
        actual := args.Map{
            params.isDefined:    lazyRegex.IsDefined(),
            params.isApplicable: lazyRegex.IsApplicable(),
            params.isMatch:      lazyRegex.IsMatch(compareInput),
            params.isFailedMatch: lazyRegex.IsFailedMatch(compareInput),
        }

        // Assert
        tc.ShouldBeEqualMap(t, caseIndex, actual)
    }
}
```

#### Why MapGherkins Over ExpectedLines

❌ **Bad — opaque ExpectedLines:**
```go
{
    Title:      "New.Lazy with simple word pattern",
    Input:      "hello",
    ExtraArgs:  map[string]any{"compareInput": "hello world"},
    ExpectedLines: []string{
        "hello",   // what is this?
        "true",    // isDefined? isApplicable? isMatch?
        "true",    // ???
        "true",    // ???
        "false",   // ???
    },
}
```

✅ **Good — self-documenting MapGherkins:**
```go
{
    Title: "New.Lazy matches word pattern",
    Input: args.Map{
        "pattern":      "hello",
        "compareInput": "hello world",
    },
    Expected: args.Map{
        "isDefined":    true,
        "isApplicable": true,
        "isMatch":      true,
        "isFailedMatch": false,
    },
}
```

### Assertion Methods — GenericGherkins

| Method | Use When |
|--------|----------|
| `ShouldBeEqual(t, caseIndex, actLines, expLines)` | Raw string line comparison |
| `ShouldBeEqualFirst(t, actLines, expLines)` | Single case, string lines |
| `ShouldBeEqualArgs(t, caseIndex, lines...)` | Variadic string args vs ExpectedLines |
| `ShouldBeEqualArgsFirst(t, lines...)` | Single case, variadic args |
| `ShouldBeEqualUsingExpected(t, caseIndex, actLines)` | Act lines vs struct's ExpectedLines |
| `ShouldBeEqualUsingExpectedFirst(t, actLines)` | Single case, vs ExpectedLines |
| `ShouldBeEqualMap(t, caseIndex, actual)` | Map-based comparison (MapGherkins) |
| `ShouldBeEqualMapFirst(t, actual)` | Single case, map comparison |

### Legacy Example — StringGherkins (still valid for simple cases)

```go
var regexMatchTestCases = []coretestcases.StringGherkins{
    {
        Title:      "Email pattern matches valid email",
        Given:      "a compiled email regex",
        When:       "matching against user@example.com",
        Input:      "user@example.com",
        IsMatching: true,
        ExtraArgs: args.Map{
            "pattern": `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
        },
    },
}
```

---

## Known Framework Limitations

> **Why this section exists** (AI-audit F05): These are **inherent limitations** of the test framework — not bugs. An AI that does not know about them will write tests that look correct but fail in non-obvious ways. Read this section before writing tests that involve interface return types, error-returning methods, or methods with more than 2 return values.

### L1 — `VerifyOutArgs` cannot match interface return types

**Symptom**: A test calling `VerifyOutArgs` against a method whose return type is an interface (most commonly `error`) reports a type mismatch even when the runtime value is correct.

**Root cause**: `VerifyOutArgs` compares types via `reflect.TypeOf(actualValue)`. `reflect.TypeOf` returns the **concrete dynamic type** (e.g. `*errors.errorString`), never the **declared interface type** (e.g. `error`). The two cannot be equal in a reflect comparison, so the verification always fails when the declared return is an interface.

**MUST-do workaround**:
- For methods returning `error` only — assert on the error value directly using `errcore` helpers; do **not** use `VerifyOutArgs` for the error position.
- For mixed-return methods — split the assertion: use `VerifyOutArgs` for concrete-typed positions and assert the interface-typed position separately.
- When deliberately covering this limitation in a test, set the expectation to `ok: false, noErr: false` so the failure is acknowledged rather than treated as a bug.

**Source**: `spec/05-failing-tests/09-failing-tests-round4.md` → `Test_I13_VerifyOutArgs_Success`.

### L2 — `InvokeFirstAndError` only works with 2-return methods

**Symptom**: Calling `InvokeFirstAndError` on a method with 3+ return values panics with a type-assertion failure on `results[1].(error)`.

**Root cause**: The helper is hard-coded to treat `results[1]` as `error`. For a method like `func MultiReturn() (int, string, error)`, index 1 is `string` — the assertion panics.

**MUST-do workaround**:
- Use `InvokeFirstAndError` **only** for methods with the exact signature `func(...) (T, error)`.
- For 3+ return values, invoke via `coredynamic` directly and pull the error from the correct index manually.
- When deliberately covering this limitation in a test, wrap the call in `recover()` and assert that a panic occurred.

**Source**: `spec/05-failing-tests/09-failing-tests-round4.md` → `Test_I14_InvokeFirstAndError_MultiReturn`.

### L3 — `t.Run` sub-test failures always propagate to the parent

**Symptom**: A test that intentionally exercises a failure path (to cover the "what happens on failure" branch of an assertion helper) marks the whole parent test as failed even though the failure is expected.

**Root cause**: Go's `testing.T.Run` propagates child failures to the parent unconditionally. There is no flag to mark a sub-test as "expected to fail".

**MUST-do workaround**: For deliberate-failure coverage tests, pass a fresh `&testing.T{}` value as the parent rather than the real `*testing.T`. The synthetic `T` absorbs the failure without polluting the actual test outcome.

**Source**: `spec/05-failing-tests/09-failing-tests-round4.md` → Learning #1.

### Reporting a new limitation

If a future failing-tests round identifies another framework quirk that an AI cannot discover from the public API alone:

1. Add a new `L<N>` entry above with the same four fields: **Symptom / Root cause / MUST-do workaround / Source**.
2. Cross-link the originating `spec/05-failing-tests/<round>.md` entry.
3. Bump the spec minor version and add a `readme.txt` milestone.

This section — not the `05-failing-tests/` corpus — is the **canonical** discovery point for framework limitations going forward.
