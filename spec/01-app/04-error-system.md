# 04 — Error System

> ✅ **Status**: filled in audit Step 5 (2026-04-23, Asia/Kuala_Lumpur).
> **Audience**: anyone returning, wrapping, or asserting on errors in `core-v9` code.
> **Source**: extracted from [`/spec/00-llm-integration-guide.md`](../00-llm-integration-guide.md) §errcore + audit cross-references.

---

## Table of Contents

1. [Public API — `errcore`](#1-public-api--errcore)
2. [Internal Contract — `errcoreinf`](#2-internal-contract--errcoreinf)
3. [Error Wrapping & Merging](#3-error-wrapping--merging)
4. [Nil-Receiver Behaviour](#4-nil-receiver-behaviour)
5. [Diagnostic Output Format](#5-diagnostic-output-format)
6. [Function Type Aliases](#6-function-type-aliases)
7. [When to Use Which API](#7-when-to-use-which-api)

---

## 1. Public API — `errcore`

The `errcore` package is the canonical way to construct errors. It exposes:

| Element | Kind | Purpose |
|---|---|---|
| `RawErrorType` | `string` type | 80+ predefined error categories with constructor methods |
| `ShouldBe` | struct-as-namespace var | Assertion-style messages and errors |
| `Expected` | struct-as-namespace var | "Expected X but got Y" formatting |
| `StackEnhance` | struct-as-namespace var | Stack-trace-aware wrapping |
| `MergeErrors`, `ManyErrorToSingle`, `SliceToError` | functions | Combining multiple errors |
| `HandleErr(err error)` | function | **Panic helper used exclusively by `*Must` variants.** No-op when `err == nil`; panics with a stack-enhanced wrapping when `err != nil`. **Always prefer `errcore.HandleErr(err)` over bare `panic(err)` in `*Must` methods.** *(F-V12-05)* |
| `MustBeEmpty(err error)` | function | **Sister of `HandleErr`** used by non-`*Must` call sites that still want fail-fast semantics — typical pattern is package-internal helpers that genuinely cannot recover (e.g. `compressformats/all-validation-checking-err.go`, `dbdrivertype/connectionStringCompiler.go`). Same nil-safety: no-op if `err == nil`, panics otherwise. Pick `HandleErr` inside `*Must` constructors; pick `MustBeEmpty` in invariant-asserting helpers. |
| `RawErrCollection` | struct type | **Accumulator** for batched validation. Embed as a field (e.g. `osdetect/windowsSystemDetailGenerator_windows.go: rawErrCollection errcore.RawErrCollection`) and append errors as you discover them; flush as a single merged error at the end. Use when one operation can produce many independent failures. |
| `ToError(...)` / `ToString(err error) string` | functions | **Conversion helpers**. `ToString` safely renders an `error` (returns `""` if nil) for log/JSON fields — see `osdetect/vars.go:111`. `ToError` is the inverse for serialised payloads. |
| `MessageWithRef(name string, ref any) string` | function | Returns a `string` (not an `error`) of the form `"<name> = <ref>"` for embedding inside a parent error's message. Used heavily in `*/vars.go` to attach a reference table to the package-level error template. |
| `RangeNotMeet(label string, min, max, ranges any) string` | function | Domain-specific range-violation message builder. See `internal/messages/messages.go` for the canonical use. |
| `VarTwo`, `VarTwoNoType`, `MessageVarMap` | functions | Variable-context formatting |
| `ErrFunc`, `ErrBytesFunc`, `ErrStringsFunc`, `ErrStringFunc`, `ErrAnyFunc` | type aliases | Common error-returning function signatures |

### 1.1 `RawErrorType` — Typed Error Categories

`RawErrorType` is a `string` type defined in `errcore/RawErrorType.go`. The package exports 80+ predefined values; each is itself a value of that type with attached methods.

Common categories:

```go
errcore.InvalidValueType                       // "Invalid : value cannot process it."
errcore.CannotBeNilOrEmptyType                 // "Values or value cannot be nil or null or empty."
errcore.NotFound                               // "not found"
errcore.FailedToParseType                      // "Failed : request failed to parse!"
errcore.ValidationFailedType                   // "Validation failed!"
errcore.UnMarshallingFailedType                // "Failed to unmarshal or deserialize."
errcore.OutOfRangeType                         // "Out of range : given value, cannot process it."
errcore.FailedToConvertType                    // "Failed to convert : input shape cannot be parsed."

// Additional categories exercised by enum-v10 (see audit cycle 2):
errcore.NotSupportedType                       // "Operation not supported on this variant."
errcore.PathInvalidErrorType                   // "Path invalid: cannot resolve / open."
errcore.FailedToExecuteType                    // "Execution failed."
errcore.ComparatorShouldBeWithinRangeType      // "Comparator value out of allowed range."
```

> The full enumeration (80+ values) lives in `errcore/RawErrorType.go` upstream. The list above shows the categories most often used by `enum-v10`; consult the upstream file for the exhaustive set.

### 1.2 Constructor Methods on `RawErrorType`

| Method | Signature | Use |
|---|---|---|
| `Error(name string, ref any) error` | with reference | When you have both a field name and a value to include |
| `ErrorNoRefs(name string) error` | name only | Simple categorical errors with no value context |
| `Fmt(format string, args ...any) error` | printf-style | Custom message formatting |
| `FmtIf(cond bool, format string, args ...any) error` | conditional | Returns nil if `cond` is false |
| `MergeError(other error) error` | wrap one | Combines this category with a downstream error |
| `MergeErrorWithMessage(other error, msg string) error` | wrap + label | Combine with extra context string |
| `ErrorRefOnly(ref any) error` | value only | When the value itself is the entire context — no field name needed (e.g. `errcore.OutOfRangeType.ErrorRefOnly(badIndex)`). Used heavily by enum constructors and `panic(errcore.NotSupportedType.ErrorRefOnly(it))`-style guards. |
| `CombineWithAnother(other error) error` | wrap one | Alias of `MergeError` preserved for older call sites (e.g. `errcore.FailedToParseType.CombineWithAnother(downstreamErr)`). New code should prefer `MergeError`. |

```go
err := errcore.InvalidValueType.Error("field name", someRef)
err := errcore.FailedToParseType.Fmt("cannot parse %q as date", input)
err := errcore.ValidationFailedType.FmtIf(len(name) == 0, "name is required")
err := errcore.NotFound.ErrorNoRefs("user with id 42")
err := errcore.FailedToConvertType.MergeError(originalErr)
err := errcore.FailedToConvertType.MergeErrorWithMessage(originalErr, "while converting X")
err := errcore.OutOfRangeType.ErrorRefOnly(givenIndex)         // value-only form
err := errcore.FailedToParseType.CombineWithAnother(parseErr)  // legacy alias

// HandleErr — the canonical panic helper for *Must variants:
//   func HandleErr(err error)
//   - if err == nil → no-op
//   - if err != nil → panics with a stack-enhanced wrapping of err
errcore.HandleErr(err) // never bare panic(err); see method-writing pattern §*Must

// MustBeEmpty — sister of HandleErr for non-*Must invariant assertions:
//   func MustBeEmpty(err error)
//   Same nil-safety. Use inside helpers that genuinely cannot recover.
errcore.MustBeEmpty(err) // see compressformats/all-validation-checking-err.go
```

### 1.3 Struct-as-Namespace Entry Points

`errcore` follows the **struct-as-namespace** pillar (see [`02-design-philosophy.md` §3](./02-design-philosophy.md)). Functionality is grouped under unexported structs exposed as package-level vars:

```go
// Assertion-style messages and errors
msg := errcore.ShouldBe.StrEqMsg("actual", "expected")
err := errcore.ShouldBe.AnyEqErr(got, want)

// Expectation comparison (with type info)
err := errcore.Expected.But("config", "production", "staging")
err := errcore.Expected.ButUsingType("field", 42, "not a number")

// Stack trace enhancement
err := errcore.StackEnhance.Error(originalErr)
msg := errcore.StackEnhance.Msg("something went wrong")
```

> **Rule**: When you need a new error pattern that doesn't fit `RawErrorType`, check if it belongs under one of these namespaces before adding a new top-level function.

### 1.4 Variable-Context Formatting

When the error message must show "what values were involved":

```go
// Two-variable context (with types)
msg := errcore.VarTwo("src", srcVal, "dst", dstVal)
// → "(src [t:string], dst[t:int]) = (hello, 42)"

// Without type tags
msg := errcore.VarTwoNoType("left", 5, "right", 10)
// → "(left, right) = (5, 10)"

// Message + variable map
msg := errcore.MessageVarMap("validation failed", map[string]any{
    "field":  "email",
    "reason": "invalid",
})
```

These produce **strings**, not errors. Wrap with `errors.New(msg)` or pass to `RawErrorType.Fmt` if you need an error value.

### 1.5 Reference Helpers — `MessageWithRef` and `RangeNotMeet`

Two additional **string** producers cover the common "attach a reference table to a package-level error template" pattern used throughout `*/vars.go`:

```go
// MessageWithRef(name, ref) → "name = <ref>"  (string, not error)
mapReferenceMessage := errcore.MessageWithRef(
    "mapping list",
    isSetterWithVariantMap)
// → see onofftype/vars.go:86, promptclitype/vars.go:112

// RangeNotMeet(label, min, max, ranges) → range-violation message
errMsg := errcore.RangeNotMeet(
    errcore.ComparatorShouldBeWithinRangeType.String(),
    corecomparator.Min(),
    corecomparator.Max(),
    corecomparator.Ranges())
// → see internal/messages/messages.go
```

Use these to build a **package-level constant message** once at init time, then reference it from every constructor / validator that needs the same wording. This keeps error text consistent across all enums in a package.

### 1.6 Accumulating Errors — `RawErrCollection`

When one operation can produce **many independent failures** (e.g. probing several Windows registry keys), use `errcore.RawErrCollection` as a struct field instead of returning early on the first error:

```go
type windowsSystemDetailGenerator struct {
    rawErrCollection errcore.RawErrCollection // see osdetect/windowsSystemDetailGenerator_windows.go:16
    rootRegistryKey  registry.Key
}

// Append errors as you discover them; flush once at the end.
g.rawErrCollection.Add(err1)
g.rawErrCollection.Add(err2)
finalErr := g.rawErrCollection.ToError() // nil if no errors were appended
```

`RawErrCollection` follows the same nil-safety rules as the merge functions (§3.2): flushing an empty collection yields `nil`, never a non-nil wrapper around zero errors.

### 1.7 Conversion Helpers — `ToString` and `ToError`

For log/JSON serialisation where a `nil` error must render as an empty value:

```go
// ToString(err) → "" if err == nil, err.Error() otherwise
errStr := errcore.ToString(err)            // see osdetect/vars.go:111
payload.Error = errcore.ToString(err)      // safe to call without a nil check

// ToError(s) — inverse for deserialised payloads
err := errcore.ToError(payload.Error)      // returns nil if s == ""
```

Use `ToString` instead of `if err != nil { ... } else { "" }` boilerplate at every JSON boundary.

---

## 2. Internal Contract — `errcoreinf`

`coreinterface/errcoreinf/` defines the interface contracts that error producers and consumers can implement without depending on the concrete `errcore` package. This lets lower layers (e.g. `coredata`) accept `errcoreinf.SomeErrorer` instead of importing `errcore` directly, preventing upward dependencies (see [`03-import-conventions.md` §5](./03-import-conventions.md#5-avoiding-cyclic-imports)).

### Common interfaces

| Interface | Method | Purpose |
|---|---|---|
| `MessageGetter` | `Message() string` | Read the error message text |
| `IsErrorChecker` | `IsError() bool` | Boolean error-state predicate |
| `ErrorWrapper` | `WrapError(error) error` | Compose two errors |
| `ShouldBeer` | `ShouldBe...` family | Assertion-style helpers |

> **Rule**: When writing a function that "needs an error formatter", accept the smallest `errcoreinf.*` interface that does the job. Do not parameter-type with `*errcore.RawErrorType`.

---

## 3. Error Wrapping & Merging

### 3.1 The three combiners

```go
combined  := errcore.MergeErrors(err1, err2, err3)              // variadic
single    := errcore.ManyErrorToSingle(errorSlice)              // []error → error
fromLines := errcore.SliceToError([]string{"issue 1", "issue 2"}) // []string → error
```

### 3.2 Behaviour

- All three return `nil` if every input is nil/empty — **never** return a non-nil error wrapping zero underlying errors.
- `MergeErrors` and `ManyErrorToSingle` preserve the order of inputs in the resulting message.
- `SliceToError` joins lines with the OS-aware line ending from `constants.DefaultLine` (always `"\n"` regardless of platform).

### 3.3 Wrapping pattern

When you catch a downstream error and want to add context:

```go
result, err := downstream.DoWork(input)
if err != nil {
    return errcore.FailedToParseType.MergeErrorWithMessage(err, "during do-work step")
}
```

The merged error's `.Error()` output contains both the category prefix and the original message — useful for log scanning.

> **Pitfall**: Do **not** `fmt.Errorf("... %w", err)` and lose the `RawErrorType` category. Use `MergeError` / `MergeErrorWithMessage` so downstream consumers can still pattern-match the category.

---

## 4. Nil-Receiver Behaviour

This is non-negotiable across the framework — see [`02-design-philosophy.md` §4 (Zero-Nil Safety)](./02-design-philosophy.md).

### Rules for error producers

1. **A nil pointer receiver method must not panic.** If the wrapped value is nil, return `nil` (no error), an empty string, or a zero value as appropriate.
2. **Constructor methods on `RawErrorType` always return a non-nil `error`** when called — they cannot themselves be nil because `RawErrorType` is a value type, not a pointer.
3. **Merge/combine functions return nil if all inputs are nil.** Never return a non-nil error that wraps zero real errors.

### Rules for error consumers

1. Always check `if err != nil` — do not assume a function returns nil only on success without checking.
2. When asserting on error categories, use `errors.Is` semantics if available, or compare against the `RawErrorType` value directly.

### Test enforcement

The framework has a dedicated test style — `CaseNilSafe` — to verify nil-receiver behaviour. Every pointer-receiver method on every public type should have at least one nil-receiver test. See [`/spec/06-testing-guidelines/02-test-case-types.md` §CaseNilSafe](../06-testing-guidelines/02-test-case-types.md#casenilsafe).

---

## 5. Diagnostic Output Format

When an error reaches the user (or a log), it should follow a predictable shape so log parsers and AI agents can extract the relevant fields.

### Single-error format

```
<RawErrorType prefix> : <user message> (var-context if any)
```

Examples:

```
Invalid : value cannot process it. : field "age" got -1
Validation failed! : (email [t:string], reason [t:string]) = (foo@bar, missing @)
```

### Merged-error format

```
<merged-error-prefix>
  <child error 1>
  <child error 2>
  ...
```

Two-space indent per nested level; one error per line.

### Why this matters

The PowerShell test runner (`run.ps1`) parses error output to attribute failures to specific source files. See [`/spec/04-tooling/03-powershell-implementation.md` §8 Error Attribution System](../04-tooling/03-powershell-implementation.md#8-error-attribution-system). Drift in error format breaks attribution — the format is part of the contract, not an implementation detail.

---

## 6. Function Type Aliases

Common function signatures are aliased so APIs that accept "any error-returning function" can be precise without `func() (X, error)` repetition:

```go
errcore.ErrFunc          // func() error
errcore.ErrBytesFunc     // func() ([]byte, error)
errcore.ErrStringsFunc   // func() ([]string, error)
errcore.ErrStringFunc    // func() (string, error)
errcore.ErrAnyFunc       // func() (any, error)
```

Use these in parameter lists and struct fields instead of inline function types.

---

## 7. When to Use Which API

| You want to… | Use |
|---|---|
| Categorical error with a value reference | `RawErrorType.Error(name, ref)` |
| Categorical error, no value | `RawErrorType.ErrorNoRefs(name)` |
| Custom message under a category | `RawErrorType.Fmt(format, args...)` |
| Conditional error (return nil if false) | `RawErrorType.FmtIf(cond, format, args...)` |
| Add context to an existing error | `RawErrorType.MergeError(err)` or `MergeErrorWithMessage` |
| Combine multiple errors | `MergeErrors`, `ManyErrorToSingle`, or `SliceToError` |
| Assertion-style "should be equal" | `ShouldBe.StrEqMsg` / `AnyEqErr` |
| "Expected X but got Y" | `Expected.But` / `ButUsingType` |
| Stack-trace-enhanced wrapping | `StackEnhance.Error` / `StackEnhance.Msg` |
| Two-variable context formatting | `VarTwo` / `VarTwoNoType` |
| Format a `map[string]any` of variables | `MessageVarMap` |
| Value-only error (no field name) | `RawErrorType.ErrorRefOnly(ref)` |
| Legacy alias for `MergeError` | `RawErrorType.CombineWithAnother(err)` (prefer `MergeError` in new code) |
| Panic-on-error inside a `*Must` constructor | `errcore.HandleErr(err)` |
| Panic-on-error inside an invariant helper | `errcore.MustBeEmpty(err)` |
| Build a "name = ref" message fragment | `errcore.MessageWithRef(name, ref)` |
| Build a range-violation message | `errcore.RangeNotMeet(label, min, max, ranges)` |
| Accumulate many errors, flush once | `errcore.RawErrCollection` |
| Render an error to string for JSON / logs | `errcore.ToString(err)` |
| Parse a string back into an error | `errcore.ToError(s)` |

---

## Boundary Cases — `FailedToConvertType` vs `ValidationFailedType`

These two categories are easy to confuse. The rule is **about the input
shape, not about what went wrong**:

- **`FailedToConvertType`** — the input *cannot be parsed* into the target
  type. The shape itself is wrong.
- **`ValidationFailedType`** — the input *parses fine*, but the resulting
  value fails a business or range rule.

| Input | Operation | Use type | Rationale |
|---|---|---|---|
| `"abc"` | `converters.StringTo.Integer` | `FailedToConvertType` | Unparseable shape — never was a number |
| `"1500"` | `converters.StringTo.DurationMillis` | `FailedToConvertType` | Missing required unit suffix — shape rule |
| `"-5"` | `PositiveIntegerValidator.IsSuccess` | `ValidationFailedType` | Parses as integer, fails range rule |
| `""` | `NonEmptyLineValidator.IsSuccess` | `ValidationFailedType` | Format-valid string, business-rejected |
| `"99999999999999999999"` | `converters.StringTo.Int32` | `FailedToConvertType` | Overflow — int32 cannot hold value |
| `42` (out of allowed enum) | `EnumValidator.IsSuccess` | `ValidationFailedType` | Type matches; value not in allow-list |

**Decision tree**:

```
Did the input parse into the target Go type at all?
├── No  → FailedToConvertType
└── Yes → Did it pass business / range / allow-list rules?
          ├── No  → ValidationFailedType
          └── Yes → no error
```

> Resolves [`/spec/02-app-issues/08-errcore-type-boundary-examples.md`](../02-app-issues/08-errcore-type-boundary-examples.md).


## See Also

- [`02-design-philosophy.md`](./02-design-philosophy.md) — Zero-nil safety pillar (§4) and struct-as-namespace pillar (§3)
- [`03-import-conventions.md`](./03-import-conventions.md) — Why `errcore` is L2 in the layer graph
- [`/spec/06-testing-guidelines/02-test-case-types.md`](../06-testing-guidelines/02-test-case-types.md) — `CaseNilSafe` for nil-receiver tests
- [`/spec/04-tooling/03-powershell-implementation.md`](../04-tooling/03-powershell-implementation.md) — Error attribution system (depends on diagnostic format)
- [`/spec/00-llm-integration-guide.md` §errcore](../00-llm-integration-guide.md#errcore--error-construction) — Quick reference
