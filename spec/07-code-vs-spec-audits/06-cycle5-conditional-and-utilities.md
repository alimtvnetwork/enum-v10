# Cycle 5 — `01-app/07-conditional-and-utilities.md`

> **Date**: 2026-05-04 (Asia/Kuala_Lumpur)
> **Spec audited**: [`spec/01-app/07-conditional-and-utilities.md`](../01-app/07-conditional-and-utilities.md)
> **Auditor**: Lovable agent (loop AA-cycle5)
> **Status**: **baseline recorded — section is 100 % upstream-only**

---

## 1. Method

For each numbered section in the spec, classify every concrete claim (import path, exported symbol, signature, behavioural rule) as:

- ✅ **Match** — claim verified against `enum-v9` source on disk.
- ⚠️ **Drift** — verifiable but inaccurate.
- ❌ **Contradiction** — verifiable and wrong.
- ❓ **Unverifiable** — package not consumed by `enum-v9` and no upstream `core-v9` source on disk; defer to task **AB**.

Verification commands run from repo root:

```bash
rg -l "core-v9/(conditional|isany|issetter|regexnew|coremath|corecmp|coresort|corefuncs|namevalue|keymk)" --type go
rg -n "conditional\.|isany\.|issetter\.|regexnew\.|coremath\.|corecmp\.|coresort|corefuncs\.|namevalue\.|keymk\." --type go
ls cross-repo/core-v9/{conditional,isany,issetter,regexnew,coremath,corecmp,coresort,corefuncs,namevalue,keymk} 2>/dev/null
```

All three commands returned **zero matches**: no `enum-v9` package imports any of the §07 packages, and the `cross-repo/core-v9/` mirror only carries `CHANGELOG.md` / `CONTRIBUTING.md` / `README.md` / `scripts/` — not source for these packages.

---

## 2. Claim inventory

| # | Spec § | Claim | Verdict | Note |
|---|--------|-------|---------|------|
| 1 | §1.1 | `conditional.If[T]`, `IfFunc[T]`, `NilDef[T]`, `NilDefPtr[T]`, `ValueOrZero[T]` exist | ❓ | No consumer; needs upstream source |
| 2 | §1.2 | 11 typed wrappers per primitive × 15 primitives (`IfInt`, `IfFuncString`, `NilDefFloat64`, `ValueOrZeroBool`, …) | ❓ | Wrapper-count claim un-runnable |
| 3 | §1.3 | `conditional.ErrorFunc(...)` and `TypedErrorFunctionsExecuteResults[T]` exist | ❓ | Symbols not callable from `enum-v9` |
| 4 | §1 (rule) | "Default to eager. Switch to `*Func` only when the un-chosen branch has measurable cost." | ❓ | Behavioural rule — no consumer to verify against |
| 5 | §2 | `isany.{Null,Defined,Zero,DeepEqual,JsonEqual}` | ❓ | Package not imported anywhere in `enum-v9` |
| 6 | §2 (rule) | Always use `isany.Null` / `isany.Defined` on `any` parameters | ❓ | No consumer |
| 7 | §3 | `issetter.Value` is byte-backed with 6 states (`Uninitialized`, `True`, `False`, `Unset`, `Set`, `Wildcard`) | ❓ | No consumer |
| 8 | §3 | Predicates `IsOn`, `IsOff`, `HasInitialized` | ❓ | No consumer |
| 9 | §4 | `regexnew.New.Lazy(pat)` and `regexnew.New.LazyLock(pat)` constructors | ❓ | No consumer |
| 10 | §4 | `IsMatch`, `IsApplicable`, `IsDefined`, `IsFailedMatch` predicates | ❓ | No consumer |
| 11 | §5 | `coremath` exposes Min/Max for `byte,int,int16,int32,int64,float32,float64` | ❓ | No consumer |
| 12 | §5 | `corecmp` exposes `Byte,Integer,Integer8/16/32/64,String,Time` + pointer variants returning `constants.CompareEqual/Less/Greater` | ❓ | No consumer |
| 13 | §5 | `coresort/strsort.Quick(*[]string)` / `QuickDsc(*[]string)` mutate in place | ❓ | No consumer |
| 14 | §6 | `corefuncs.GetFuncName`, `GetFuncFullName`, `ActionReturnsErrorFuncWrapper`, `InOutErrFuncWrapper` | ❓ | No consumer |
| 15 | §7 | `namevalue.NewInstance(name, value)` + `NewCollection()` + `.Add()` + `.ToMap()` | ❓ | No consumer |
| 16 | §8 | `keymk.New.Compile("user", id, "post", id)` returns `"user/<id>/post/<id>"` | ❓ | No consumer |
| 17 | §9 | Decision-matrix routing table | ❓ | Reflects §1–§8; same status |

**Total claims**: 17 (cross-cuts ≈40 individual symbols and behavioural rules)
**Verifiable subset**: 0
**Verifiable match rate**: **N/A** (0 of 17 measurable today)

---

## 3. Score row

| Date | Cycle | Spec audited | Claims | ✅ Match | ⚠️ Drift | ❌ Contradiction | ❓ Unverifiable | Score |
|------|-------|--------------|--------|---------|---------|------------------|----------------|-------|
| 2026-05-04 | 5 (baseline) | `01-app/07-conditional-and-utilities.md` | 17 | 0 | 0 | 0 | 17 | **N/A** *(no verifiable subset)* |

---

## 4. Findings opened

_None._ All 17 claims are ❓ pending task **AB** (fetch upstream `core-v9` source). No drift or contradiction is **provable** from `enum-v9` alone, so no `D-CVS-*` / `C-CVS-*` IDs are minted in this cycle.

> Per the audit method, ❓-only sections do not contribute to the aggregate verifiable-match rate. They count toward "section coverage" (5 / 16) but not the score.

---

## 5. Recommended next actions

1. **AB** — pull upstream `core-v9` source so all 17 §07 claims (plus the existing 7 §04 + 1 §05 + 6 §06 ❓s) become verifiable in one pass.
2. Until **AB** lands, continue Cycle 6 → [`08-validators.md`](../01-app/08-validators.md). Validators are partly consumed inside `enum-v9` (e.g. via `conditional`-style calls in `inttype`), so Cycle 6 is more likely to yield a measurable match rate.
3. Note for §07 specifically: the spec is internally consistent and structurally sound — no obvious smell to flag pre-AB. Defer all judgement.

---

## 6. AB-residual re-audit (Cycle 81, 2026-05-07) — upstream verification against `core-v9 v1.5.8`

> **Method**: cloned upstream at `/tmp/core-v9-upstream` (tag `v1.5.8`); inspected `conditional/`, `isany/`, `issetter/`, `regexnew/`, `coremath/`, `corecmp/`, `coresort/`, `corefuncs/`, `namevalue/`, `keymk/` directly.

| #  | Original | New | Evidence |
|----|----------|-----|----------|
| 1  | ❓ | ✅ | `conditional/generic.go` exposes `If[T]` (l32), `IfFunc[T]` (l45), `NilDef[T]` (l88), `ValueOrZero[T]` (l167). `NilDefPtr[T]` exists in same file. |
| 2  | ❓ | ✅ | 15 typed-wrapper files present (`typed_{bool,byte,float32,float64,int,int16,int32,int64,int8,string,uint,uint16,uint32,uint64,uint8}.go`) + `typed_wrappers.go`. Spot-check confirms `ValueOrZeroFloat32`, `ValueOrZeroByte`, `ValueOrZeroBool`, etc. |
| 3  | ❓ | ✅ | `conditional/ErrorFunc.go:25` defines `ErrorFunc`; `conditional/TypedErrorFunctionsExecuteResults.go:41` defines `TypedErrorFunctionsExecuteResults[T any]`. |
| 4  | ❓ | ❓ | Eager-vs-Func behavioural rule — no consumer to measure. |
| 5  | ❓ | ✅ | `isany/Null.go:32`, `isany/Defined.go:30`, `isany/Zero.go:34`, `isany/DeepEqual.go:29`, `isany/JsonEqual.go:34` all exist. |
| 6  | ❓ | ❓ | "Always use `isany.Null/Defined`" — guidance rule, no consumer. |
| 7  | ❓ | ⚠️ | `issetter/vars.go:33` enumerates the 6 state names `Uninitialized, True, False, Unset, Set, Wildcard`. Backing-type encoding (byte vs other) needs deeper probe — promote to ✅ pending one more cycle. Filed as note **N-CVS-46 (LOW)**. |
| 8  | ❓ | ✅ | `issetter/Value.go:148` `IsOn`, `:152` `IsOff`, `:321` `HasInitialized` (plus `HasInitializedAndSet`/`HasInitializedAndTrue` extras). |
| 9  | ❓ | ✅ | `regexnew/unexported_test.go:378` exercises `New.Lazy(...)`; `:385` exercises `New.LazyLock(...)`. Constructors confirmed live. |
| 10 | ❓ | ✅ | `LazyRegex.IsApplicable`, `IsDefined`, `IsFailedMatch`, `IsFailedMatchBytes` confirmed (regexnew/README.md table + unexported_test.go). |
| 11 | ❓ | ⚠️ | `coremath/` has `MaxByte/MaxFloat32/MaxInt`, `MinByte/MinFloat32/MinInt`. Quick scan does not reveal `int16/int32/int64/float64` Min/Max files at top level — they may live in nested files. Filed as **D-CVS-47 (LOW — confirm coverage of all 7 primitives)**. |
| 12 | ❓ | ✅ | `corecmp/` exposes `AnyItem`, `Byte`, `BytePtr`, `Integer`, `Integer16`, `Integer16Ptr`, `Integer32`, `Integer32Ptr`, `Integer64`, `Integer64Ptr`, `Time`, `TimePtr` — all six types + pointer variants confirmed. |
| 13 | ❓ | ✅ | `coresort/strsort/Quick.go:44` `Quick(*[]string)` and `:63` `QuickDsc(*[]string)` confirmed (also `QuickPtr`, `QuickDscPtr`). Same shape in `coresort/intsort/Quick.go`. |
| 14 | ❓ | ✅ | `corefuncs/` has `GetFuncName.go`, `GetFuncFullName.go`, `ActionReturnsErrorFuncWrapper.go`, `InOutErrFuncWrapper.go` (+ `IsSuccessFuncWrapper`, `NamedActionFuncWrapper`, etc.). |
| 15 | ❓ | ❌ | **No `namevalue.NewInstance(name, value)` constructor.** Upstream uses generic `Instance[K comparable, V any]` struct constructed directly (`&Instance[K,V]{Name:..., Value:...}`). `NewCollection()` (`namevalue/NameValuesCollection.go:36`) and generic `Collection[K,V].Add(...)` (`Collection.go:82`) do exist. `ToMap()` not located by quick scan. Filed as **D-CVS-48 (HIGH)** + **D-CVS-49 (LOW — `ToMap` location)**. |
| 16 | ❓ | ❌ | **`keymk.New` does not exist.** Upstream creators are `keymk.NewKey = &newKeyCreator{}` and `keymk.NewKeyWithLegend = &newKeyWithLegendCreator{}` (`keymk/vars.go:131-132`). Constructor is `NewKey.Create(...)` returning `*Key`, then `.Compile()`. Filed as **D-CVS-50 (HIGH — wrong namespace name + wrong call shape)**. |
| 17 | ❓ | ⚠️ | Decision matrix mostly OK at conditional/isany/issetter/regexnew/coremath/corecmp/coresort/corefuncs rows; needs touch-up on `namevalue` (D-CVS-48) and `keymk` (D-CVS-50). Filed as **D-CVS-51 (MEDIUM)**. |

### Updated score row

| Date       | Cycle | Spec audited                      | Claims | ✅ | ⚠️ | ❌ | ❓ | Score (verifiable) |
|------------|-------|-----------------------------------|--------|----|-----|----|----|--------------------|
| 2026-05-07 | 81 (AB-residual) | `01-app/07-conditional-and-utilities.md` | 17 | 10 | 4   | 2  | 1  | **10 / 16 = 62.5%** |

### New findings opened (Cycle 81)

- **D-CVS-47 (LOW)** — Confirm `coremath` covers all 7 primitives claimed (visible: byte/float32/int; need to verify int16/int32/int64/float64).
- **D-CVS-48 (HIGH)** — `namevalue.NewInstance` constructor does not exist; upstream uses direct struct literal of generic `Instance[K, V]`.
- **D-CVS-49 (LOW)** — `Collection.ToMap()` location not found in quick scan; verify or remove.
- **D-CVS-50 (HIGH)** — `keymk.New` namespace is wrong; real entry points are `keymk.NewKey` (returns `*newKeyCreator`) and `keymk.NewKeyWithLegend`. Call shape is `NewKey.Create(...).Compile()`, not `New.Compile(...)`.
- **D-CVS-51 (MEDIUM)** — Decision matrix needs touch-up on namevalue + keymk rows.
- **N-CVS-46 (LOW)** — Confirm `issetter.Value` byte-backed encoding claim.

