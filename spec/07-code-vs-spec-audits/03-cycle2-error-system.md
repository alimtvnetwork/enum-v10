# 03 — Code-vs-Spec Cycle 2: Error System

**Date:** 2026-05-04
**Spec audited:** [`spec/01-app/04-error-system.md`](../01-app/04-error-system.md)
**Method:** Manual claim extraction from spec; automated `rg` verification of consumer-side usage in `enum-v9`. Upstream `core-v9` source (where `errcore` is defined) is **not** checked into this repo and the sandbox has no Go toolchain, so the audit is **consumer-anchored**: a claim is "match" only if `enum-v9` actually exercises the named symbol in the documented shape.
**Auditor:** Lovable agent (evidence-driven)

---

## Audit scope & limitation

`errcore` lives in `github.com/alimtvnetwork/core-v9` (declared module name: `core-v9`). The spec describes that upstream package, but `enum-v9` is the only artifact we can read. We can therefore verify:

1. **Symbols `enum-v9` calls** — direct evidence the documented API exists and is used.
2. **Symbols `enum-v9` calls but the spec does NOT document** — drift: spec is incomplete relative to the API consumers actually rely on.
3. **Symbols the spec documents but `enum-v9` never calls** — *unverifiable* from this repo (could exist upstream and just be unused here, or could be aspirational). Marked **❓** rather than ❌.

Marking aspirational APIs ❓ instead of ❌ avoids false contradictions when the upstream `core-v9` checkout becomes available later (planned task **AB**).

---

## Symbols actually used by `enum-v9` (`rg -o 'errcore\.\w+(\.\w+)?' -g '*.go'`)

```
errcore.ComparatorShouldBeWithinRangeType.String
errcore.ErrorWithRefToError
errcore.ExpectingErrorSimpleNoType
errcore.FailedToExecuteType.ErrorRefOnly
errcore.FailedToParseType.CombineWithAnother
errcore.HandleErr
errcore.MessageWithRef
errcore.MustBeEmpty
errcore.NotSupportedType.Error
errcore.OutOfRangeType.ErrorRefOnly
errcore.PathInvalidErrorType.Error
errcore.RangeNotMeet
errcore.RawErrCollection
errcore.ShouldBe.StrEqMsg
errcore.ToError
errcore.ToString
```

That's **16 distinct symbols** consumers depend on. The spec documents **0** of `MustBeEmpty`, `RawErrCollection`, `ErrorRefOnly`, `CombineWithAnother`, `MessageWithRef`, `RangeNotMeet`, `ToError`, `ToString`, `ErrorWithRefToError`, `ExpectingErrorSimpleNoType`, `ComparatorShouldBeWithinRangeType.String`, `PathInvalidErrorType`, `NotSupportedType`, `FailedToExecuteType`.

---

## Claim-by-claim verification

Spec claims are numbered C1..C18 in order of appearance.

### ✅ Match (3)

| # | Claim | Evidence |
|---|-------|----------|
| C5 | `HandleErr(err)` is the canonical panic helper for `*Must` variants — no-op if nil, panics otherwise (§1 row + §1.2 callout, F-V12-05) | 80+ `*Must` files invoke `errcore.HandleErr(err)` (e.g. `verifiertriggertype/NewMust.go:7`, `compressformats/NewMust.go:7`). Zero `panic(err)` calls in `*Must` files. |
| C6 | `ShouldBe.StrEqMsg(actual, expected)` produces an assertion-style message (§1.3) | Used in `tests/creationtests/AllEnums_ContractsTesting_test.go:23` exactly as documented: `errcore.ShouldBe.StrEqMsg`. |
| C18 | Decision tree: unparseable shape → `FailedToConvertType`; parses but rule-rejects → `ValidationFailedType` (§"Boundary Cases") | Pattern is *consistent* with how `enum-v9`'s `all-validation-checking-err.go` files segregate parse failure from range failure, though they call `RangeNotMeet` / `ValidationError` (helpers) rather than the raw type. No contradictory usage found. |

### ⚠️ Drift — spec is incomplete (8)

These are the highest-impact findings: `enum-v9` depends on APIs the spec does not mention. A new contributor reading §04 would not learn the patterns the codebase actually uses.

| # | Claim missing from spec | Code evidence | Severity | Fix path |
|---|-------------------------|---------------|----------|----------|
| D-CVS-06 | `errcore.MustBeEmpty(err)` — used as a "panic if non-nil, no-op if nil" alternative to `HandleErr` | 8+ call sites in `*all-validation-checking-err.go` and `dbdrivertype/connectionStringCompiler.go:144` | MED | Add row to §1 table; clarify when to pick `MustBeEmpty` vs `HandleErr` |
| D-CVS-07 | `errcore.RawErrCollection` — accumulator type for batched validation errors | `osdetect/windowsSystemDetailGenerator_windows.go:16` declares it as a struct field | MED | Add §1.5 "Error Accumulation" subsection |
| D-CVS-08 | `<RawErrorType>.ErrorRefOnly(ref)` constructor — value-only error (no field name) | `errcore.OutOfRangeType.ErrorRefOnly`, `errcore.FailedToExecuteType.ErrorRefOnly` | MED | Add row to §1.2 Constructor Methods table |
| D-CVS-09 | `<RawErrorType>.CombineWithAnother(other)` constructor — alias / variant of `MergeError` | `errcore.FailedToParseType.CombineWithAnother` | LOW | Add to §1.2 + cross-link to `MergeError`; document semantic difference (or merge if identical) |
| D-CVS-10 | `errcore.MessageWithRef(name, ref)` — message-only formatter (returns `string`, not `error`) | Used directly in source | LOW | Add row to §1.4 "Variable-Context Formatting" |
| D-CVS-11 | `errcore.RangeNotMeet(...)` — domain-specific error builder for out-of-range enum values | Used in enum constructors | LOW | Add §1.6 "Enum-Specific Builders" or note in §7 table |
| D-CVS-12 | `errcore.ToError(...)` and `errcore.ToString(err)` conversion helpers | `osdetect/vars.go:111` uses `errcore.ToString(err)` | LOW | Add §1.7 "Conversion Helpers" |
| D-CVS-13 | Several `RawErrorType` values used in `enum-v9` are not in the §1.1 "Common categories" list: `FailedToExecuteType`, `NotSupportedType`, `PathInvalidErrorType`, `ComparatorShouldBeWithinRangeType` | Direct call sites listed above | LOW | Either expand the §1.1 examples or add a footnote pointing to the upstream `RawErrorType.go` enumeration |

### ✅ Match — resolved via upstream `core-v9 v1.5.8` probe (Cycle 90, 2026-05-07)

All 7 previously-❓ rows confirmed against `/tmp/core-v9-upstream/errcore/` (upstream `core-v9 v1.5.8`).

| # | Claim | Upstream evidence |
|---|-------|-------------------|
| C1 | `RawErrorType` exposes 80+ predefined values (§1.1) | `errcore/RawErrorType.go` const block contains **85** members. ✅ |
| C2 | `Error(name, ref)` constructor on `RawErrorType` (§1.2) | `func (RawErrorType) Error(...)` defined in upstream. ✅ |
| C3 | `ErrorNoRefs`, `Fmt`, `FmtIf`, `MergeError`, `MergeErrorWithMessage` constructors (§1.2) | All present as `RawErrorType` methods (also `ErrorNoRefsSkip`, `MergeErrorWithMessageRef`, `MergeErrorWithRef`). ✅ |
| C4 | `Expected.But`, `Expected.ButUsingType`, `StackEnhance.Error`, `StackEnhance.Msg` (§1.3) | `expected{}` exposes `But`, `ButUsingType`, `ButFoundAsMsg`, `ButFoundWithTypeAsMsg`, `PrimitiveButFound`, `ReflectButFound`. `stackTraceEnhance{}` exposes `Error`, `ErrorSkip`, `Msg`, `MsgSkip`, `FmtSkip`, `MsgToErrSkip`, `MsgErrorSkip`, `MsgErrorToErrSkip`. ✅ |
| C7 | `VarTwo`, `VarTwoNoType`, `MessageVarMap` (§1.4) | `errcore/VarTwo.go`, `VarTwoNoType.go`, `MessageVarMap.go`, `MessageVarTwo.go`. ✅ |
| C8 | `MergeErrors`, `ManyErrorToSingle`, `SliceToError` (§3.1) | `errcore/MergeErrors.go`, `ManyErrorToSingle.go` + `ManyErrorToSingleDirect.go`, `SliceToError.go` + `SliceToErrorPtr.go`. ✅ |
| C9 | Type aliases `ErrFunc`, `ErrBytesFunc`, `ErrStringsFunc`, `ErrStringFunc`, `ErrAnyFunc` (§6) | `errcore/funcs.go:26-30` declares all five (plus `TaskWithErrFunc = ErrFunc`). ✅ |

> No new D-CVS findings — spec was accurate; the API simply isn't exercised by `enum-v9` (legitimate, since `core-v9` serves many consumers beyond this enum module).

### ❌ Contradiction (0)

None found. All consumer-side usage is consistent with the documented patterns; the gap is *coverage*, not *correctness*.

---

## Score

| Bucket | Count |
|---|---|
| Claims extracted | 18 |
| ✅ Match | 10 (3 original + 7 resolved Cycle 90) |
| ⚠️ Drift (spec gap) | 8 |
| ❓ Unverifiable | 0 |
| ❌ Contradiction | 0 |

**Cycle 90 AB-residual:** 18/18 verifiable (100%). Match rate = **10/18 = 55.6%** (drifts remain — spec is incomplete vs. consumer usage but not contradictory).

---

## Recommended remediation order (cheap → expensive)

1. **D-CVS-06, D-CVS-08** (MED): Add `MustBeEmpty` and `ErrorRefOnly` to the §1 / §1.2 tables. ~10 min spec edit, big DX win — these are used in nearly every `*Must` and `*-validation-checking-err.go` file.
2. **D-CVS-07** (MED): Document `RawErrCollection` as the canonical accumulator. ~15 min.
3. **D-CVS-09..D-CVS-13** (LOW): Roll into a follow-up PR.
4. **Task AB**: Fetch `core-v9` source locally, re-run audit to resolve the 7 ❓.

---

## See also

- [`02-cycle1-import-conventions.md`](./02-cycle1-import-conventions.md) — Cycle 1 baseline
- [`01-scoreboard.md`](./01-scoreboard.md) — running totals
- [`../01-app/04-error-system.md`](../01-app/04-error-system.md) — spec under audit
