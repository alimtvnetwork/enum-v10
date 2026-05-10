# 05 — `params.go` Files Convention

> **Status**: resolved (2026-04-23)
> **Severity**: low
> **Resolved**: 2026-04-23 — adopted proposed rule and added it to `spec/06-testing-guidelines/03-args-reference.md` §"`params.go` convention" (see edit below).

> **Scope note (`enum-v10`, added Cycle 18)**: `tests/integratedtests/` and `errcoretests/` cited below are **upstream `core-v9`** package names. `enum-v10` uses a single shared `tests/creationtests/` package with shared `vars.go` (no per-package `params.go`), so the "grandfathered, no back-fill" rule applies vacuously here. The `params.go` convention itself remains relevant only for upstream consumers. Historical resolution text is preserved verbatim below.

## Decision

**`params.go` is mandatory for new test packages with > 3 test cases**, optional otherwise. Existing packages are **grandfathered** — no back-fill required.

### Rationale

- Centralising `args.Map` keys in `params.go` pays off only when keys are reused across ≥ 3 cases; below that, inline string literals are clearer.
- Mandating a back-fill across the entire `tests/integratedtests/` tree would generate noise without behaviour change.
- The spec previously implied `params.go` was universal, which contradicts visible practice (only `errcoretests` and a handful of others follow it).

### Spec edit applied

Added to `spec/06-testing-guidelines/03-args-reference.md`:

```markdown
## `params.go` Convention

A `params.go` file at the root of a test package centralises `args.Map` key constants used across multiple test cases.

**Required for new packages** with more than 3 distinct test cases sharing args keys.
**Optional** for packages with ≤ 3 cases — inline string literals are equally readable.
**Not required to back-fill** existing packages.

Example:

\`\`\`go
// errcoretests/params.go
package errcoretests

const (
    keyInput    = "input"
    keyExpected = "expected"
    keyWantErr  = "wantErr"
)
\`\`\`
```

## Original Problem

`spec/06-testing-guidelines/` recommended a `params.go` file per test package, but in practice only a few packages followed it. Spec contradicted observed reality.

## Verification

Step 11 fresh-AI simulation did not re-trigger this issue under the new rule.
