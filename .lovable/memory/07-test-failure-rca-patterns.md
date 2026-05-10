# Test Failure RCA Patterns

> **Last audit:** 2026-05-10 (Cycle 52, Task AT). All 3 active patterns (P8/P9/P10) verified still mitigated in current source — see audit log at bottom. Patterns 1–7 are historical (recoverable from git history); see `.lovable/plan.md` cycles 49–63 for original context.

## Pattern 10 — Brace-unaware Go declaration scanner (added 2026-05-07)

**Symptom:** `scripts/ci/check-collisions.py` (or any similar tooling) reports cross-file "var collisions" for symbols like `got`, `n`, `data`, `sink` across `*_Coverage_test.go` and `*_Uplift_test.go` pairs in many packages, plus intra-file dupes for the same identifier appearing in multiple test functions.

**Root cause:** The scanner matched `^\s*var\s+NAME` line-by-line without tracking `{`/`}` brace depth. Function-local `var got Variant` declarations inside `t.Run` bodies were classified as package-level vars, colliding with identically-named locals in sibling test files.

**Fix:** Maintain a `brace_depth` counter; compute `at_top = (brace_depth == 0)` at the START of each line; only emit a declaration when `at_top` is true. Then `brace_depth += line.count("{") - line.count("}")`. String/comment stripping must happen first so braces inside `"..."` or `// ...` don't shift depth.

**Detection sweep:** Run the collision checker against a known-clean tree; any hit on a common local name (`got`, `err`, `n`, `data`, `sink`, `v`) is almost certainly a brace-tracking bug, not a real collision.

**Manifested as:** 16 cross-file false positives + 1 intra-file false positive across the test suite (v1.15.0 fix).

---

## Pattern 9 — `Stringer` infinite recursion via `converters.AnyTo.ValueString(self)` (added 2026-05-07)

**Symptom:** A package's tests crash with `fatal error: stack overflow`, the trace alternating between `runtime/fmt/print.go` frames and one specific source line, e.g. `brackets/Pair.go:114` or `brackets/BothBrackets.go:100`. Go test runner reports it as a runtime failure for the package even though no `t.Error` was called.

**Root cause:** A `func (it T) String() string` implementation that returns `converters.AnyTo.ValueString(it)`. The helper falls through to `fmt.Sprintf("%v", it)`, which invokes the type's own `String()` method again — infinite recursion → stack overflow.

**Fix:** Never call `converters.AnyTo.ValueString` (or any `%v`/`fmt.Sprint` formatter) on the receiver inside a `String()` method. Always format the struct's fields explicitly:

```go
func (it Pair) String() string {
    return fmt.Sprintf("{Start:%s End:%s Category:%s}",
        it.Start.String(), it.End.String(), it.Category.String())
}
```

**Detection sweep:** `rg "converters\.AnyTo\.ValueString" --type go | rg -v _test.go` — review every hit; any one that passes `it`, `*it`, or the receiver's value to a `String()`/`MarshalJSON` method is a recursion bomb.

**Manifested as:** `brackets` package reported as RUNTIME FAILURE with stack overflow originating in `Pair.String` and `BothBrackets.String` (v1.14.0 fix).

---

## Pattern 8 — `-coverpkg` warning-only false-positive (added 2026-05-07)

**Symptom:** Packages reported as `Blocked` in pre-coverage compile check, or as `RUNTIME FAILURE` after the coverage run, but the captured diagnostic contains ONLY repeated lines like:

    warning: no packages being tested depend on matches for pattern github.com/alimtvnetwork/enum-v9/...

These warnings are emitted by Go when a `-coverpkg=` glob includes packages the test binary doesn't transitively import. They are harmless. Under parallel runspace load, the build cache contention sometimes causes `go test -c` to also exit non-zero transiently, masking the warning-only nature.

**Fix:** `scripts/Utilities.psm1` exports `Test-IsCoverpkgWarningOnlyOutput`. Call it as a final guard in:
- `CoverageCompileCheck.psm1` parallel runspace (set `$confirmed = $false` when warnings-only) and in the post-parallel re-confirmation pass.
- `CoverageRunner.psm1` sync + parallel coverage loops, so warning-only output never reaches `Add-RuntimeFailuresForPackage` / `Add-BuildErrorsForPackage`.

Manifested as: `licensetype` 58.9%, `onofftype` 32.1%, `rootcmdnames` 58.1% all blocked while passing tests cleanly elsewhere; `brackets` falsely listed under RUNTIME FAILURES despite "✓ No failing tests".

---

(prior patterns 1–7 retained — see git history)

---

## Audit Log

### 2026-05-10 (Cycle 52, Task AT) — all patterns ✅ still mitigated

- **P8 (`-coverpkg` warning-only false-positive):** `Test-IsCoverpkgWarningOnlyOutput` exported from `scripts/Utilities.psm1:118` and wired into the four call sites that previously misclassified warnings as failures (`CoverageRunner.psm1:217,236`; `CoverageCompileCheck.psm1:200`). Mitigation intact.
- **P9 (`Stringer` infinite recursion via `converters.AnyTo.ValueString(self)`):** Detection sweep `rg "converters\.AnyTo\.ValueString" --type go -g '!*_test.go'` returns **zero hits**. The recursion bombs that motivated this pattern have not regressed.
- **P10 (Brace-unaware Go declaration scanner):** `scripts/ci/check-collisions.py` retains `brace_depth` accounting (line 103) with `at_top = (brace_depth == 0)` gate (line 147) and per-branch `brace_depth += opens - closes` updates (lines 155, 160, 165, 176). Mitigation intact.

No new patterns surfaced this cycle. Index entry in `mem://index.md` says "4 patterns" — that referred to the original Cycles 57/60/63 set (P1–P4). The doc has since rotated to the 3 most recent active patterns (P8–P10) with P1–P7 retained in git history. Index wording left unchanged for historical accuracy.
