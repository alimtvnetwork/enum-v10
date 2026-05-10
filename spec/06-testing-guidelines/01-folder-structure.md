# 01 — Folder Structure & Naming Conventions

> ⚠️ **Scope:** the layout below describes **upstream `core-v9`**. `enum-v10` uses a single shared `tests/creationtests/` package — see [`spec/01-app/13-testing-patterns.md` §6.1](../01-app/13-testing-patterns.md#61-enum-v10-specific-layout). Per-package `*tests/` directories under `tests/integratedtests/` apply to upstream consumers only.

## Directory Layout

```
project-root/
├── coretests/                      # Framework code (shared across all tests)
│   ├── args/                       # Input holders (Map, One–Six, Dynamic, Holder, LeftRight)
│   ├── coretestcases/              # Case types (CaseV1, CaseNilSafe, GenericGherkins)
│   └── results/                    # Result types (Result, ResultAny, InvokeWithPanicRecovery)
│
├── tests/
│   └── integratedtests/
│       ├── mypkgtests/             # One directory per source package
│       │   ├── MyFunc_test.go      # Test runner — assertions only
│       │   ├── MyFunc_testcases.go # Test data — expectations only
│       │   ├── MyStruct_NilReceiver_testcases.go  # Nil-safety test data
│       │   ├── NilReceiver_test.go # Nil-safety test runner
│       │   └── helpers.go          # Shared test-only structs/utilities
│       └── anotherpkgtests/
│           └── ...
```

## Naming Rules

### Directories

| Pattern | Example | Rule |
|---------|---------|------|
| `{package}tests/` | `errcoretests/` | Always lowercase, always suffixed with `tests` |

### Files

| Pattern | Purpose | Example |
|---------|---------|---------|
| `{Feature}_test.go` | Test runner (assertions) | `Sort_test.go` |
| `{Feature}_testcases.go` | Test data (expectations) | `Sort_testcases.go` |
| `{Type}_NilReceiver_testcases.go` | Nil-safety test data | `WrappedErr_NilReceiver_testcases.go` |
| `NilReceiver_test.go` | Nil-safety test runner (one per package) | `NilReceiver_test.go` |
| `helpers.go` | Shared test-only structs | `helpers.go` |

### Test Functions

```go
// Pattern: Test_{TypeOrFeature}_{Scenario}_Verification
func Test_IntSort_Quick_Verification(t *testing.T) { ... }

// Pattern: Test_{Type}_NilReceiver
func Test_WrappedErr_NilReceiver(t *testing.T) { ... }

// For single (non-loop) tests, drop "Verification":
func Test_ConcatMessageWithErr_NilPassthrough(t *testing.T) { ... }
```

### Test Case Variables

```go
// Slice for loop tests:
var intSortQuickTestCases = []coretestcases.CaseV1{ ... }

// Standalone for single tests:
var concatMessageNilPassthroughTestCase = coretestcases.CaseV1{ ... }

// Nil-safety slice:
var wrappedErrNilSafeTestCases = []coretestcases.CaseNilSafe{ ... }
```

## Separation Rules

### `_testcases.go` — ONLY contains:
- `var` declarations of test case slices/structs
- Helper functions for test data construction (e.g., `errFromString`)
- **NO** `import "testing"`, **NO** `func Test_`, **NO** assertions

### `_test.go` — ONLY contains:
- `func Test_` functions
- `// Arrange`, `// Act`, `// Assert` comment blocks
- Calls to framework assertion methods
- **NO** expected values hardcoded in the function body

### `helpers.go` — ONLY contains:
- Shared structs used across multiple test files (e.g., `exampleStruct`)
- Shared utility functions (e.g., `errFromString`)
- **NO** test functions, **NO** test data

## Package Declaration

All files in a test directory use the same package name:

```go
package errcoretests  // NOT package errcoretests_test
```

This allows access to unexported test helpers within the same directory.
