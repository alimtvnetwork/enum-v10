# Cycle 6 — `01-app/08-validators.md`

> **Date**: 2026-05-05 (Asia/Kuala_Lumpur)
> **Spec audited**: [`spec/01-app/08-validators.md`](../01-app/08-validators.md)
> **Auditor**: Lovable agent (loop AA-cycle6)
> **Status**: **baseline recorded — section is 100 % upstream-only**

---

## 1. Method

For each numbered section in the spec, classify every concrete claim (import path, exported symbol, signature, behavioural rule, diagnostic-format contract) as:

- ✅ **Match** — claim verified against `enum-v10` source on disk.
- ⚠️ **Drift** — verifiable but inaccurate.
- ❌ **Contradiction** — verifiable and wrong.
- ❓ **Unverifiable** — package not consumed by `enum-v10` and no upstream `core-v9` source on disk; defer to task **AB**.

Verification commands run from repo root:

```bash
rg -l "core-v9/corevalidator" --type go
rg -n "corevalidator\.(New|LineValidator|SliceValidator|TextValidator|RangeValidator|StringCompareAs|Result)" --type go
rg -n "errcore\.VarTwoNoType|ValidationFailedType" --type go
rg -n "regexnew\.New\.Lazy" --type go
rg -n "coretestcases\.(CaseV1|CaseNilSafe)" --type go
ls cross-repo/core-v9/corevalidator 2>/dev/null
```

All commands returned **zero matches**: no `enum-v10` package imports `corevalidator` or any of its peers (`errcore.VarTwoNoType`, `ValidationFailedType`, `regexnew.New.Lazy`, `coretestcases.CaseV1/CaseNilSafe`), and the `cross-repo/core-v9/` mirror does not carry `corevalidator` source.

---

## 2. Claim inventory

| #  | Spec § | Claim                                                                                      | Verdict | Note |
|----|--------|--------------------------------------------------------------------------------------------|---------|------|
| 1  | §1     | `corevalidator` provides Line/Slice/Text/Range validators implementing `coreinterface/*`   | ❓ | No consumer; needs upstream source |
| 2  | §1     | Validator contract: `IsValid`, `IsSuccess`, `IsFailed`, `Message`, `Error`                 | ❓ | Contract surface unverifiable |
| 3  | §2.1   | `corevalidator.New.Line.NotEmpty().MaxLength(N).Matches(re).Build() *LineValidator`        | ❓ | No consumer |
| 4  | §2.1   | `Build()` returns pointer; value-typed storage breaks interface assertion (F-V14-01)       | ❓ | Behavioural rule — no consumer |
| 5  | §2.2   | `corevalidator.New.Slice.MinLength.MaxLength.EachItem(v).Build()`                          | ❓ | No consumer |
| 6  | §2.3   | `TextValidator` — multi-line whole-document rules                                          | ❓ | No consumer |
| 7  | §2.4   | `corevalidator.New.Range.Int.Min(N).MaxExclusive(N).Build()` semantics                     | ❓ | No consumer |
| 8  | §2.5   | `StringCompareAs` diagnostic format is regex-checked by framework tests                    | ❓ | No consumer |
| 9  | §3.1   | `errcore.MergeErrors(errs...)` aggregates validator errors                                 | ❓ | `MergeErrors` not invoked anywhere in `enum-v10` |
| 10 | §3.2   | Domain-type embedding pattern (`(r *T) Validate() error`)                                  | ❓ | Pattern not exercised in `enum-v10` |
| 11 | §3.3   | `conditional.IfFunc[error](...)` for branching validation                                  | ❓ | `conditional` already ❓ in Cycle 5 |
| 12 | §4.1–4.3 | Custom-validator template uses `regexnew.New.Lazy` + nil-receiver guard                  | ❓ | No consumer |
| 13 | §4.3   | `Result.Error()` returns `errcore.ValidationFailedType.Fmt(...)`                           | ❓ | `ValidationFailedType` symbol not imported |
| 14 | §4.4   | Tests use `coretestcases.CaseV1` + `CaseNilSafe` (Style A)                                 | ❓ | Neither symbol imported anywhere in `enum-v10` |
| 15 | §5     | Diagnostic message format: `<Label>: field=<n> value=<v> reason=<short>`                   | ❓ | Format contract has no in-repo emitter |
| 16 | §5     | `errcore.VarTwoNoType("field", label, "value", actual)` is the canonical builder          | ❓ | Symbol not imported anywhere in `enum-v10` |
| 17 | §6     | Three-test layout: `*_Verification_test.go` / `*_NilReceiver_test.go` / `*_Format_test.go` | ❓ | No validator tests under `tests/` |
| 18 | §6     | Path: `tests/integratedtests/<pkg>tests/`                                                  | ⚠️ | **Stale path** — repeats C-CVS-01 / D-CVS-17 pattern: actual layout is `tests/creationtests/`. Filed as **D-CVS-26**. |
| 19 | §7     | Common-mistakes table (7 rows)                                                             | ❓ | All rows depend on §1–§4 surface |

**Total claims**: 19
**Verifiable subset**: 1 (claim #18 — path-string drift, verified against on-disk `tests/` layout)
**Verifiable match rate (baseline)**: **0 / 1 = 0.0%** *(one ⚠️ Drift, no other measurable claims)*

---

## 3. Score row

| Date       | Cycle | Spec audited                      | Claims | ✅ | ⚠️ | ❌ | ❓ | Score (verifiable) |
|------------|-------|-----------------------------------|--------|----|-----|----|----|--------------------|
| 2026-05-05 | 6 (baseline) | `01-app/08-validators.md`   | 19     | 0  | 1   | 0  | 18 | **0.0%** *(1 of 19 measurable today)* |

---

## 4. Findings

### D-CVS-26 — Stale `tests/integratedtests/` reference (§6) — **LOW (path-string fix)**

Spec §6 line 347:
> Each validator package needs three test files in `tests/integratedtests/<pkg>tests/`:

Repository reality: `tests/integratedtests/` does not exist; tests live under `tests/creationtests/` (this is the same drift pattern already filed as **C-CVS-01** for §03 and **D-CVS-17** for §05, both resolved with identical s/integratedtests/creationtests/ rewrites).

**Fix**: rewrite §6 path to `tests/creationtests/<pkg>tests/` and cross-reference the established convention in `13-testing-patterns.md` (matching the C-CVS-01 / D-CVS-17 resolutions).

**No other contradictions found.** All 17 ❓ rows defer to task **AB** (fetch upstream `core-v9/corevalidator` source).

---

## 5. Next actions

1. Apply **D-CVS-26** as a one-line spec fix (closes §08's only verifiable drift).
2. Update the scoreboard (`01-scoreboard.md`) with the Cycle 6 baseline + closed rows.
3. Defer the remaining 17 ❓ to **AB** alongside Cycle 5's 17 ❓ — both cycles share the same blocker (no upstream `corevalidator` / `conditional` / `isany` / `regexnew` source on disk).

---

## 6. AC re-audit (Cycle 79, 2026-05-07) — upstream verification against `core-v9 v1.5.8`

> **Method**: cloned upstream at `/tmp/core-v9-upstream` (tag `v1.5.8`), inspected `corevalidator/`, `errcore/`, `regexnew/`, `coretests/` directly. Promoted ❓ rows.

| #  | Original verdict | New verdict | Evidence |
|----|------------------|-------------|----------|
| 1  | ❓ | ⚠️ | `corevalidator/` exposes `LineValidator`, `TextValidator`, `SliceValidator` and `RangeSegmentsValidator`; **no `RangeValidator`** by that name (upstream uses `RangeSegmentsValidator`). Filed as **D-CVS-27**. |
| 2  | ❓ | ❌ | Upstream `LineValidator` exposes `IsMatch/IsMatchMany/VerifyError/VerifyMany/VerifyFirstError/AllVerifyError` — **no `IsValid`/`IsSuccess`/`IsFailed`/`Message`/`Error` contract**. Filed as **D-CVS-28**. |
| 3  | ❓ | ❌ | **No `corevalidator.New.Line.NotEmpty().MaxLength(N).Matches(re).Build()` fluent builder exists** in upstream. Constructors are top-level funcs (`NewSliceValidatorUsingErr`, etc.) and validators are constructed as plain structs (`TextValidator{Search, SearchAs, Condition}`). Filed as **D-CVS-29** (entire fluent-API surface fictitious). |
| 4  | ❓ | ❌ | No `Build()` method exists; F-V14-01 rule cannot apply. Subsumed by D-CVS-29. |
| 5  | ❓ | ❌ | No `New.Slice.MinLength.MaxLength.EachItem(v).Build()` builder. Subsumed by D-CVS-29. |
| 6  | ❓ | ✅ | `TextValidator` exists and supports multi-line via `verifyDetailErrorUsingLineProcessing`. |
| 7  | ❓ | ❌ | No `New.Range.Int.Min(N).MaxExclusive(N).Build()`. Upstream provides `RangeSegmentsValidator` with a different surface. Subsumed by D-CVS-29. |
| 8  | ❓ | ✅ | `stringcompareas` enum exists (imported by `corevalidator/vars.go`). |
| 9  | ❓ | ✅ | `errcore.MergeErrors(errs ...error) error` exists at `errcore/MergeErrors.go:32`. |
| 10 | ❓ | ❓ | Pattern still not exercised in `enum-v10`; no upstream change required. |
| 11 | ❓ | ✅ | `conditional/` package exists with `ErrorFunc`, `Functions`, etc. |
| 12 | ❓ | ✅ | `regexnew/CreateLock.go` + `LazyRegex.go` provide the lazy regex surface. |
| 13 | ❓ | ✅ | `errcore.ValidationFailedType` exists at `errcore/RawErrorType.go:121` and is used by `corevalidator/LinesValidators.go:188,221`. |
| 14 | ❓ | ⚠️ | `coretests/` (note: **`coretests`**, not `coretestcases`) provides `BaseTestCase` family. Spec's `CaseV1` / `CaseNilSafe` symbol names not found upstream. Filed as **D-CVS-30**. |
| 15 | ❓ | ⚠️ | Diagnostic format `<Label>: field=<n> value=<v> reason=<short>` not emitted verbatim by upstream `VarTwoNoType`. Filed as **D-CVS-31**. |
| 16 | ❓ | ✅ | `errcore.VarTwoNoType(n1, v1, n2, v2)` exists at `errcore/VarTwoNoType.go:25`. |
| 17 | ❓ | ❓ | Behavioural; no validator tests in `enum-v10` to verify against. |
| 18 | ⚠️ | ⚠️ | (D-CVS-26 already filed) — unchanged. |
| 19 | ❓ | ⚠️ | All rows depend on §1–§4 surface; with D-CVS-29 confirming the fluent API is fictitious, table needs ground-up rewrite. Filed as **D-CVS-32**. |

### Updated score row

| Date       | Cycle | Spec audited                      | Claims | ✅ | ⚠️ | ❌ | ❓ | Score (verifiable) |
|------------|-------|-----------------------------------|--------|----|-----|----|----|--------------------|
| 2026-05-07 | 79 (AC re-audit) | `01-app/08-validators.md` | 19 | 6  | 5   | 5  | 3  | **6 / 16 = 37.5%** |

### New findings opened (Cycle 79)

- **D-CVS-27 (LOW)** — Spec claims `RangeValidator`; upstream provides `RangeSegmentsValidator`.
- **D-CVS-28 (HIGH)** — Validator contract surface (`IsValid`/`IsSuccess`/`IsFailed`/`Message`/`Error`) does not exist upstream. Real surface is `IsMatch*`/`Verify*Error`.
- **D-CVS-29 (CRITICAL)** — Entire `corevalidator.New.<X>.…Build()` fluent builder API in §2.x is fictitious. Upstream uses `NewSliceValidatorUsingErr`/`NewSliceValidatorUsingAny` and plain struct construction.
- **D-CVS-30 (LOW)** — Spec uses `coretestcases.CaseV1`/`CaseNilSafe`; upstream package is `coretests` with a `BaseTestCase` family — neither symbol exists by name.
- **D-CVS-31 (LOW)** — §5 diagnostic format string is aspirational; not the literal output of `VarTwoNoType`.
- **D-CVS-32 (MEDIUM)** — §7 common-mistakes table assumes the fictitious fluent API; needs rewrite once §2.x is reworked against real upstream surface.

