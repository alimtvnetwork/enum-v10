# Stringer-recursion static guard (RCA Pattern P9)

## Symptom
A `func (it T) String() string` body that calls `converters.AnyTo.ValueString` is an infinite-recursion bomb: that helper falls through to `fmt.Sprintf("%v", ...)`, which re-invokes the type's own `String()` method → stack overflow at runtime, often surfaced only by tests in the affected package.

## Root Cause
The pattern is documented in `.lovable/memory/07-test-failure-rca-patterns.md` (Pattern P9, originally added 2026-05-07 after the `brackets/Pair.go` and `brackets/BothBrackets.go` crashes). Until Cycle 54 there was no automated check preventing the pattern from creeping back in.

Cycle 54 audit (Task AT) verified the live tree was clean via an `rg` sweep. But human discipline alone is not durable — Cycle 54 (Task AX) immediately surfaced one real new violation in `scripttype/ScriptDefault.go` that had landed without a sweep noticing.

## Fix / Workaround
Cycle 54 (v1.31.0, Task AX):

- Added `scripts/ci/check-stringer-recursion.py` — scans every non-test, non-vendor, non-`cross-repo/` `.go` file. Brace-depth tracks each `String() string` method body (cf. RCA Pattern P10) and fails non-zero if any contains `converters.AnyTo.ValueString`. Sufficient string/comment stripping so braces inside `"..."` or `// ...` don't shift depth.
- Added `scripts/ci/test_check_stringer_recursion.py` — 8 unittest cases: clean source passes, value-receiver bomb detected, pointer-receiver bomb detected, `_test.go` files excluded, calls outside `String()` allowed, brace-in-string-literal robustness, missing root, no-Go-files. All 8 pass; full `unittest discover` count: 48 → **56 OK**.
- Added `stringer-recursion-guard` job in `.github/workflows/ci-guards.yml` (runs after `python-tests`).
- Fixed real violation in `scripttype/ScriptDefault.go`: replaced `converters.AnyTo.ValueString(*it)` with explicit per-field `fmt.Sprintf` and added a `nil`-receiver guard. `go build ./scripttype/...` and `go test ./scripttype/... -count=1` both pass.

## Status
✅ Working as designed. Live-tree scan: scanned 672 Go source files → ✅ no violations.

## Related
- `.lovable/memory/07-test-failure-rca-patterns.md` Pattern P9 — original bug pattern.
- `.lovable/memory/07-test-failure-rca-patterns.md` Pattern P10 — brace-depth tracking technique reused here.
- Task AT (Cycle 52) — manual audit that motivated the static guard.
- Task AX (Cycle 54) — this guard.
