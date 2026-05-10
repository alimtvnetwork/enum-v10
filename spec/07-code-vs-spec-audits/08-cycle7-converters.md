# Cycle 7 — `01-app/09-converters.md`

> **Date**: 2026-05-05 (Asia/Kuala_Lumpur)
> **Spec audited**: [`spec/01-app/09-converters.md`](../01-app/09-converters.md)
> **Auditor**: Lovable agent (loop AA-cycle7)
> **Status**: **baseline recorded — section is 100 % upstream-only**

---

## 1. Method

For each numbered section in the spec, classify every concrete claim (import path, exported symbol, signature, behavioural rule, error-category contract) as:

- ✅ **Match** — claim verified against `enum-v9` source on disk.
- ⚠️ **Drift** — verifiable but inaccurate.
- ❌ **Contradiction** — verifiable and wrong.
- ❓ **Unverifiable** — package not consumed by `enum-v9` and no upstream `core-v9` source on disk; defer to task **AB**.

Verification commands run from repo root:

```bash
rg -l "core-v9/(converters|typesconv)" --type go
rg -n "converters\.(StringTo|BytesTo|PrettyJson)" --type go
rg -n "typesconv\." --type go
rg -n "errcore\.(FailedToConvertType|OverflowType)" --type go
rg -n "strconv\.(Atoi|ParseInt|ParseFloat|ParseBool)" --type go
ls cross-repo/core-v9/{converters,typesconv} 2>/dev/null
```

All commands returned **zero matches**: no `enum-v9` package imports `converters` or `typesconv`, none of the documented symbols (`StringTo.Integer`, `BytesTo.String`, `PrettyJson.FromAny`, `IntToInt64`, `Int64ToInt32`, `Float64ToInt`, etc.) appear in source, and the `cross-repo/core-v9/` mirror does not carry either package. The `strconv.Atoi/ParseBool` anti-pattern from §5 is also absent — there is nothing to violate the rule against.

---

## 2. Claim inventory

| #  | Spec § | Claim                                                                                                  | Verdict | Note |
|----|--------|--------------------------------------------------------------------------------------------------------|---------|------|
| 1  | §1     | `converters/` package exists with struct-as-namespace pattern                                          | ❓ | No consumer; needs upstream source |
| 2  | §1.1   | `converters.StringTo.Integer(s) (int, error)`                                                          | ❓ | No consumer |
| 3  | §1.1   | `converters.StringTo.Integer64(s) (int64, error)`                                                      | ❓ | No consumer |
| 4  | §1.1   | `converters.StringTo.IntegerWithDefault(s, def) (int, bool)` — `(value, ok)` mode                      | ❓ | No consumer |
| 5  | §1.1   | `converters.StringTo.Float64(s)` / `Float32(s)`                                                        | ❓ | No consumer |
| 6  | §1.1   | `converters.StringTo.Byte(s) (byte, error)`                                                            | ❓ | No consumer |
| 7  | §1.1   | `converters.StringTo.Bool(s)` accepts only `strconv.ParseBool` set; rejects whitespace; no yes/no/on/off (F-V14-03) | ❓ | No consumer to exercise the contract |
| 8  | §1.2   | `converters.BytesTo.String([]byte) string` is direct cast (zero-copy)                                  | ❓ | No consumer |
| 9  | §1.2   | `converters.BytesTo.PrettyJsonString(jsonBytes) string`                                                | ❓ | No consumer |
| 10 | §1.3   | `converters.PrettyJson.String(jsonBytes)` / `PrettyJson.FromAny(any)`                                  | ❓ | No consumer |
| 11 | §1.3 (rule) | `PrettyJson` is legacy; prefer `corejson.NewPtr(x).PrettyJsonString()` in new code                  | ❓ | `corejson.NewPtr` itself ❓ in Cycle 4 (no consumer) |
| 12 | §2     | `typesconv/` package — non-string ↔ non-string numeric conversions                                     | ❓ | No consumer |
| 13 | §2     | `typesconv.IntToInt64(int) int64` (always-safe widening, no error)                                     | ❓ | No consumer |
| 14 | §2     | `typesconv.Int64ToInt32(int64) (int32, bool)` (narrowing returns ok flag)                              | ❓ | No consumer |
| 15 | §2     | `typesconv.Float64ToInt(float64) (int, bool)` truncates; `ok=false` on NaN/Inf                         | ❓ | No consumer |
| 16 | §2     | "When to use `typesconv` vs `converters`" decision matrix                                              | ❓ | Reflects §1+§2 surface; same status |
| 17 | §3     | Two-return-mode contract: `(value, error)` for log-on-failure / `(value, bool)` for fallback           | ❓ | Behavioural contract — no consumer |
| 18 | §3     | "No panics on bad input — always returns zero value + failure signal"                                  | ❓ | Behavioural rule — no consumer |
| 19 | §3     | Errors wrapped via `errcore.FailedToConvertType`                                                       | ❓ | `FailedToConvertType` not invoked anywhere in `enum-v9` (see Cycle 2 ❓) |
| 20 | §3     | "Locale-independent — `.` decimal separator, no thousand grouping"                                     | ❓ | Behavioural rule — no consumer |
| 21 | §3     | "Truncation is silent in `*WithDefault` variants"                                                      | ❓ | Behavioural rule — no consumer |
| 22 | §4.3   | `errcore.OverflowType.Fmt(...)` for narrowing overflow                                                 | ❓ | `OverflowType` not invoked anywhere in `enum-v9` |
| 23 | §5     | Common-mistakes table (5 rows: prefer `converters` over `strconv.Atoi`, etc.)                          | ❓ | Anti-pattern absent (no `strconv.Atoi/ParseBool` calls in `enum-v9`) — no rule-violation to flag, but consumer-side enforcement also unverifiable |

**Total claims**: 23
**Verifiable subset**: 0
**Verifiable match rate (baseline)**: **N/A** *(0 of 23 measurable today — same shape as Cycle 5)*

---

## 3. Score row

| Date       | Cycle | Spec audited                | Claims | ✅ | ⚠️ | ❌ | ❓ | Score (verifiable) |
|------------|-------|-----------------------------|--------|----|-----|----|----|--------------------|
| 2026-05-05 | 7 (baseline) | `01-app/09-converters.md` | 23     | 0  | 0   | 0  | 23 | **N/A** *(no verifiable subset)* |

---

## 4. Findings

**No drifts, no contradictions found.** Unlike Cycle 6 (which had the `tests/integratedtests/` path string), §09 contains no on-disk path references that could be measured against the repo. All 23 claims defer to task **AB**.

### Cross-cycle observation

§09 is the second consecutive section (after §07) with **zero verifiable subset** and **zero on-disk drift signals**. Both packages are documented entirely from upstream knowledge with no `enum-v9` adoption. This raises the running ❓ count on the scoreboard significantly without changing the closed-section count.

---

## 5. Next actions

1. Update the scoreboard (`01-scoreboard.md`) with the Cycle 7 baseline row and bump the §07/§08/§04/§05/§06 + §09 ❓ tally for task **AB**.
2. Mark Cycle 7 as **baseline-only closed** (no fixes needed; nothing to re-verify).
3. Continue to Cycle 8 → `10-reflection-and-dynamic.md` on next `next`.

---

## 6. AB-residual re-audit (Cycle 80, 2026-05-07) — upstream verification against `core-v9 v1.5.8`

> **Method**: cloned upstream at `/tmp/core-v9-upstream` (tag `v1.5.8`); inspected `converters/`, `typesconv/`, `internal/jsoninternal/`, and `errcore/` directly.

### Critical structural discrepancy

Upstream `typesconv/` does **not** provide the narrowing/widening conversion helpers documented in §2 (`IntToInt64`, `Int64ToInt32`, `Float64ToInt`). Real `typesconv/` is a pointer-utility package: `IntPtr`, `IntPtrToSimple`, `IntPtrToSimpleDef`, `IntPtrToDefPtr`, `IntPtrDefValFunc`, plus the same Ptr-helper family for `Float`/`Byte`/`Bool`/`String`. The `string.go` file additionally provides `StringToBool`/`StringPointerToBool` top-level funcs. Filed as **D-CVS-37 (CRITICAL)**.

| #  | Original | New | Evidence |
|----|----------|-----|----------|
| 1  | ❓ | ✅ | `converters/vars.go` exposes `StringTo`, `BytesTo`, `StringsTo`, `AnyTo`, `Map`, `PrettyJson`, `JsonString`, `Integers`, `KeyValuesTo`, `CodeFormatter` as struct-as-namespace vars. |
| 2  | ❓ | ✅ | `(it stringTo) Integer(input string) (value int, err error)` at `converters/stringTo.go:154`. |
| 3  | ❓ | ❌ | **No `Integer64` method.** Upstream surface: `Integer`, `IntegerMust`, `IntegerDefault`, `IntegerWithDefault`, `IntegersWithDefaults`, `IntegersConditional`. Filed as **D-CVS-38**. |
| 4  | ❓ | ✅ | `(it stringTo) IntegerWithDefault(...)` at `converters/stringTo.go:39` — returns `(int, bool)` ok-mode. |
| 5  | ❓ | ⚠️ | `Float64`/`Float64Must`/`Float64Default`/`Float64Conditional` exist; **no `Float32` method on `stringTo`**. Filed as **D-CVS-39 (LOW)**. |
| 6  | ❓ | ✅ | `(it stringTo) Byte(input string) (byte, error)` at `converters/stringTo.go:260`. |
| 7  | ❓ | ❌ | **No `Bool` method on `stringTo`.** String→bool lives at `typesconv/string.go:73` as the top-level `StringToBool` func — wrong package + wrong signature. F-V14-03 contract (rejects whitespace, no yes/no/on/off) cannot be verified at the documented location. Filed as **D-CVS-40 (HIGH)**. |
| 8  | ❓ | ✅ | `(it bytesTo) String([]byte) string` at `converters/bytesTo.go:37` (zero-copy via internal helper). |
| 9  | ❓ | ❌ | **No `BytesTo.PrettyJsonString` method.** `bytesTo` only has `PtrString`/`String`/`PointerToBytes`. Closest is `converters.PrettyJson.*` (separate namespace). Filed as **D-CVS-41**. |
| 10 | ❓ | ⚠️ | `converters.PrettyJson` is a re-export of `internal/jsoninternal.Pretty` (`prettyConverter` struct). Concrete method names (`String`, `FromAny`) not enumerated by quick scan — needs deeper probe to fully verify. Promote to ✅ once methods listed. |
| 11 | ❓ | ⚠️ | "PrettyJson is legacy; prefer `corejson.NewPtr(x).PrettyJsonString()`" — `corejson` exists upstream but `PrettyJson` is actively re-exported (not deprecated in vars.go). Guidance is **aspirational**, not enforced. |
| 12 | ❓ | ⚠️ | `typesconv/` exists, but its purpose is **pointer-helper utilities**, not non-string ↔ non-string numeric conversions. Filed as **D-CVS-42**. |
| 13 | ❓ | ❌ | **No `IntToInt64` func.** Closest is `IntPtr`/`IntPtrToSimple`. Subsumed by D-CVS-37. |
| 14 | ❓ | ❌ | **No `Int64ToInt32` func.** Subsumed by D-CVS-37. |
| 15 | ❓ | ❌ | **No `Float64ToInt` func.** Subsumed by D-CVS-37. |
| 16 | ❓ | ⚠️ | Decision matrix depends on §1+§2; with §2 fictitious (D-CVS-37), matrix needs rewrite. Filed as **D-CVS-43 (MEDIUM)**. |
| 17 | ❓ | ⚠️ | Two-return-mode contract partially holds: `Integer` returns `(int, error)`, `IntegerWithDefault` returns `(int, bool)`. ✅ for those two. Other claimed funcs (Bool, Integer64, Float32) don't exist so contract unverifiable for them. |
| 18 | ❓ | ❓ | "No panics on bad input" — not exhaustively verified across all upstream funcs. |
| 19 | ❓ | ✅ | `errcore.FailedToConvertType` exists at `errcore/RawErrorType.go:63` and is invoked from `errcore/RawErrCollection.go:283,289`. |
| 20 | ❓ | ❓ | Locale-independence — implicit via `strconv.ParseInt`/`ParseFloat`, but not documented in upstream code. |
| 21 | ❓ | ❓ | "Truncation silent in `*WithDefault`" — partially verifiable; not exhaustively checked. |
| 22 | ❓ | ❌ | **`errcore.OverflowType` does not exist** in upstream `errcore/RawErrorType.go`. Filed as **D-CVS-44**. |
| 23 | ❓ | ⚠️ | Common-mistakes table assumes the §1+§2 surface; needs review once §2 is reworked. Filed as **D-CVS-45 (MEDIUM)**. |

### Updated score row

| Date       | Cycle | Spec audited                | Claims | ✅ | ⚠️ | ❌ | ❓ | Score (verifiable) |
|------------|-------|-----------------------------|--------|----|-----|----|----|--------------------|
| 2026-05-07 | 80 (AB-residual) | `01-app/09-converters.md` | 23 | 6  | 8   | 7  | 2  | **6 / 21 = 28.6%** |

### New findings opened (Cycle 80)

- **D-CVS-37 (CRITICAL)** — Entire §2 `typesconv` surface fictitious. Real `typesconv/` is a pointer-utility package, not a numeric conversion package.
- **D-CVS-38 (LOW)** — `StringTo.Integer64` does not exist.
- **D-CVS-39 (LOW)** — `StringTo.Float32` does not exist.
- **D-CVS-40 (HIGH)** — `StringTo.Bool` does not exist on `converters/stringTo`. String→bool conversion lives in `typesconv/string.go` as a top-level func with different signature.
- **D-CVS-41 (LOW)** — `BytesTo.PrettyJsonString` does not exist; use `converters.PrettyJson.*`.
- **D-CVS-42 (MEDIUM)** — `typesconv/` package purpose mis-stated.
- **D-CVS-43 (MEDIUM)** — Decision matrix needs rewrite once §2 is corrected.
- **D-CVS-44 (LOW)** — `errcore.OverflowType` symbol does not exist.
- **D-CVS-45 (MEDIUM)** — Common-mistakes table needs review against real surface.

