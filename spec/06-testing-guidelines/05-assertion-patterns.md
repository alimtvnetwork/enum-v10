# 05 — Assertion Patterns

## Mandatory AAA Format

Every test function MUST have these three comments:

```go
func Test_Feature_Scenario(t *testing.T) {
    // Arrange
    // ... setup input

    // Act
    // ... call function under test

    // Assert
    // ... verify results
}
```

No exceptions. Even if Arrange is empty, include the comment:

```go
func Test_Feature_NilInput(t *testing.T) {
    tc := featureNilTestCase

    // Arrange
    // (nil input — no setup needed)

    // Act
    result := mypkg.Feature(nil)

    // Assert
    tc.ShouldBeEqualFirst(t, fmt.Sprintf("%v", result == nil))
}
```

---

## Pattern 1: String-Based Assertions (ShouldBeEqual)

For simple single-value or multi-line comparisons:

### Loop-based (slice of test cases)

```go
func Test_Sort_Quick_Verification(t *testing.T) {
    for caseIndex, tc := range sortQuickTestCases {
        // Arrange
        input := tc.ArrangeInput.(args.Map)
        items, _ := input.Get("input")
        clone := make([]int, len(items.([]int)))
        copy(clone, items.([]int))

        // Act
        result := intsort.Quick(&clone)

        // Assert
        tc.ShouldBeEqual(t, caseIndex, fmt.Sprintf("%v", *result))
    }
}
```

### Single test case (named variable)

```go
func Test_ConcatMessage_NilPassthrough(t *testing.T) {
    tc := concatMessageNilTestCase

    // Arrange
    input := tc.ArrangeInput.(args.Map)
    msg, _ := input.GetAsString("message")

    // Act
    result := errcore.ConcatMessageWithErr(msg, nil)

    // Assert
    tc.ShouldBeEqualFirst(t, fmt.Sprintf("%v", result == nil))
}
```

---

## Pattern 2: Map-Based Assertions (ShouldBeEqualMap)

For multi-field comparisons with semantic keys. **Preferred for most tests.**

### How It Works

1. `ExpectedInput` (CaseV1) or `Expected` (MapGherkins) is `args.Map` with raw typed values
2. `actual` in test runner is `args.Map` with raw typed values
3. Both are compiled to sorted `"key : value"` strings
4. Strings are compared line-by-line

### CaseV1 — Loop-based

```go
func Test_Validator_Verification(t *testing.T) {
    for caseIndex, tc := range validatorTestCases {
        // Arrange
        input := tc.ArrangeInput.(args.Map)
        name, _ := input.GetAsString("name")
        age, _ := input.GetAsInt("age")

        // Act
        result := validator.Validate(name, age)
        actual := args.Map{
            "isValid": result.IsValid,
            "errors":  len(result.Errors),
        }

        // Assert
        tc.ShouldBeEqualMap(t, caseIndex, actual)
    }
}
```

### MapGherkins — Loop-based

```go
func Test_LazyRegex_Verification(t *testing.T) {
    for caseIndex, tc := range lazyRegexTestCases {
        // Arrange
        pattern, _ := tc.Input.GetAsString("pattern")
        compareInput, _ := tc.Input.GetAsString("compareInput")

        // Act
        lazy := regexnew.New.Lazy(pattern)
        actual := args.Map{
            "isDefined":    lazy.IsDefined(),
            "isApplicable": lazy.IsApplicable(),
            "isMatch":      lazy.IsMatch(compareInput),
        }

        // Assert
        tc.ShouldBeEqualMap(t, caseIndex, actual)
    }
}
```

### Single test case

```go
func Test_Validator_EmptyName(t *testing.T) {
    tc := validatorEmptyNameTestCase

    // Arrange
    // (uses tc.ArrangeInput directly)

    // Act
    result := validator.Validate("", 25)
    actual := args.Map{
        "isValid": result.IsValid,
        "errors":  len(result.Errors),
    }

    // Assert
    tc.ShouldBeEqualMapFirst(t, actual)
}
```

---

## Pattern 3: Nil-Safety Assertions (ShouldBeSafe)

For `CaseNilSafe` — no Arrange/Act needed, the framework handles everything:

```go
func Test_MyStruct_NilReceiver(t *testing.T) {
    for caseIndex, tc := range myStructNilSafeTestCases {
        tc.ShouldBeSafe(t, caseIndex)
    }
}
```

For a single nil-safety case:

```go
func Test_MyStruct_IsValid_NilReceiver(t *testing.T) {
    tc := myStructIsValidNilCase
    tc.ShouldBeSafeFirst(t)
}
```

---

## Pattern 4: Multi-Line String Assertions

For functions that return complex string output:

```go
var funcWrapTestCases = []coretestcases.CaseV1{
    {
        Title: "someFunctionV1 with valid params",
        ArrangeInput: args.ThreeFuncAny{
            First:    "f1",
            Second:   "f2",
            Third:    "f3",
            WorkFunc: someFunctionV1,
        },
        ExpectedInput: []string{
            "someFunctionV1 => called with (f1, f2, f3)",
        },
    },
    {
        Title: "nil func returns error",
        ArrangeInput: args.ThreeFuncAny{
            First:    "f1",
            WorkFunc: nil,
        },
        ExpectedInput: []string{
            "error : ",
            "  func-wrap is invalid:",
            "      given type: <nil>",
        },
    },
}
```

Test runner uses variadic `ShouldBeEqual`:

```go
func Test_FuncWrap_Verification(t *testing.T) {
    for caseIndex, tc := range funcWrapTestCases {
        // Arrange
        input := tc.ArrangeInput.(args.ThreeFuncAny)

        // Act
        output, err := input.InvokeWithValidArgs()
        lines := toStrings(output)
        if err != nil {
            lines = append(lines, "error : ", err.Error())
        }

        // Assert
        tc.ShouldBeEqual(t, caseIndex, lines...)
    }
}
```

---

## Native Types in Expectations

**Always use native Go types** in `args.Map` — the framework handles conversion:

```go
// ✅ GOOD — native types:
ExpectedInput: args.Map{
    "isValid": true,        // bool
    "count":   5,           // int
    "name":    "Alice",     // string
}

// ❌ BAD — stringified values:
ExpectedInput: args.Map{
    "isValid": "true",      // string — extra conversion needed
    "count":   "5",         // string — loses type safety
}

// ✅ GOOD — actual map uses native types:
actual := args.Map{
    "isValid": result.IsValid(),    // returns bool
    "count":   result.Count(),      // returns int
    "name":    result.Name(),       // returns string
}
```

---

## Test Params Pattern (`params.go`)

Every test package that uses `args.Map` keys MUST define a local `params.go` file
with a `params` struct holding reusable key constants. This eliminates magic strings,
prevents case-sensitivity bugs, and ensures consistency across test cases and runners.

### Structure

```go
// params.go
package mypkgtests

var params = struct {
	// Common assertion keys
	isDefined    string
	isApplicable string
	isMatch      string
	isFailedMatch string
	hasError     string

	// Package-specific keys
	pattern      string
	compareInput string
}{
	isDefined:    "isDefined",
	isApplicable: "isApplicable",
	isMatch:      "isMatch",
	isFailedMatch: "isFailedMatch",
	hasError:     "hasError",

	pattern:      "pattern",
	compareInput: "compareInput",
}
```

### Common Keys Reference

These keys appear across many test packages and SHOULD use consistent names:

| Key | Type | Meaning |
|-----|------|---------|
| `isDefined` | `bool` | Whether the object is non-nil / initialized |
| `isApplicable` | `bool` | Whether the object compiled / is usable |
| `isMatch` | `bool` | Whether the input matched the pattern |
| `isFailedMatch` | `bool` | Inverse of isMatch |
| `hasError` | `bool` | Whether an error was returned |
| `isValid` | `bool` | Whether validation passed |
| `isNil` | `bool` | Whether the result is nil |
| `pattern` | `string` | Regex or search pattern (input) |
| `compareInput` | `string` | Text to match against (input) |

### Usage in Test Cases

```go
// _testcases.go
var testCases = []coretestcases.MapGherkins{
    {
        Title: "Lazy regex matches word pattern",
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

### Usage in Test Runners

```go
// _test.go
func Test_LazyRegex_Verification(t *testing.T) {
    for caseIndex, tc := range testCases {
        // Arrange
        pattern, _ := tc.Input.GetAsString(params.pattern)
        compareInput, _ := tc.Input.GetAsString(params.compareInput)

        // Act
        lazy := regexnew.New.Lazy(pattern)
        actual := args.Map{
            params.isDefined:    lazy.IsDefined(),
            params.isApplicable: lazy.IsApplicable(),
            params.isMatch:      lazy.IsMatch(compareInput),
            params.isFailedMatch: lazy.IsFailedMatch(compareInput),
        }

        // Assert
        tc.ShouldBeEqualMap(t, caseIndex, actual)
    }
}
```

### Rules

1. **One `params.go` per test package** — each package defines its own params struct
2. **Package-private** — `params` is unexported (lowercase), scoped to the test package
3. **Common keys use standard names** — `isDefined`, `isMatch`, `hasError`, etc.
4. **Package-specific keys extend the struct** — add fields unique to the domain
5. **Never use raw string literals** for map keys in test cases or runners

### ❌ Bad — magic strings

```go
Expected: args.Map{
    "isDefined": true,    // raw string — typo-prone
    "isMatch":   true,
}

hasError, _ := input.GetAsBool("hasError")  // magic string
```

### ✅ Good — params constants

```go
Expected: args.Map{
    params.isDefined: true,    // compile-time safe
    params.isMatch:   true,
}

hasError := input.GetAsBoolDefault(params.hasError, false)  // safe accessor + constant
```

---

## Safe Accessors in Test Runners

Use typed getters to avoid panics:

```go
// ✅ GOOD — safe accessor with ok check:
msg, ok := input.GetAsString("message")
if !ok {
    t.Fatal("missing 'message' in ArrangeInput")
}

// ✅ GOOD — safe accessor with default:
count := input.GetAsIntDefault("count", 0)

// ❌ BAD — raw map access (panic on missing key):
msg := input["message"].(string)
```

---

## Diff-Based Error Output

All assertions ultimately delegate to `errcore.AssertDiffOnMismatch`, which produces structured diff output on failure:

```
=== FAIL: Test_Validator_Verification (case 2)
    Title: "Validate empty name returns error"
    
    Expected:
      errors : 1
      isValid : false
    
    Actual:
      errors : 0
      isValid : true
    
    Diff:
      - errors : 1
      + errors : 0
      - isValid : false
      + isValid : true
```

This is why `args.Map` with semantic keys is preferred — failures show exactly which field mismatched.

---

## Known Pitfall: Named Map Types and convertinternal.AnyTo.Strings

`args.Map` is a **named type** (`type Map map[string]any`). Go's type switch does NOT
match named types against their underlying type. This means `args.Map` will NOT match
a `case map[string]any:` branch and will fall through to the `default` case (PrettyJSON).

**Rule**: When passing `args.Map` as `ExpectedInput`, always convert it to `[]string`
via `CompileToStrings()` before it reaches `ExpectedLines()`. The `ShouldBeEqualMap`
method handles this automatically — never set `ExpectedInput` to a raw `args.Map`
and then call `ShouldBe` directly.

---

## Sub-Pattern: GoConvey-Only Diff Assertion

> **Scope:** Companion to the [GoConvey-Only sub-pattern in §02](./02-test-case-types.md#sub-pattern-goconvey-only-local-wrapper). Use when a downstream package (e.g. `enum-v10`) opts out of `coretests` entirely.

The GoConvey-only equivalent of `ShouldBeEqualMap` is the **diff-based assertion**:

```go
// Arrange
actualEnumDynamicMap := dynamicEnumMapOf(tc.Enum)

// Act
diff := actualEnumDynamicMap.LogShouldDiffMessage(true, header, expected)

// Assert
So(diff, ShouldBeEmpty)
```

**Why diff-based:** Returns the empty string on success and a human-readable field-by-field diff on failure. `So(diff, ShouldBeEmpty)` then surfaces the diff in the GoConvey report — same failure ergonomics as `ShouldBeEqualMap`, no `coretests` dependency.

**Companion assertions** routinely paired with this pattern:

| Assertion | Purpose |
|---|---|
| `So(value, ShouldEqual, expected)` | Scalar comparison |
| `So(value, ShouldResemble, expected)` | Deep struct/slice comparison |
| `So(value, ShouldBeNil)` | Nil-safety check |
| `So(value, ShouldBeTrue)` / `ShouldBeFalse` | Boolean invariants |
| `So(diff, ShouldBeEmpty)` | Diff-based map/struct equality (above) |

**Worked example:** `github.com/alimtvnetwork/enum-v10/tests/creationtests/AllEnums_ContractsTesting_test.go`.
