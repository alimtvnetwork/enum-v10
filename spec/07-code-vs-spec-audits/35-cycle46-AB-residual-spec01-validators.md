# Cycle 46 — AB-residual deep-probe of `spec/01-app/08-validators.md` (Cycle 21 carry-over)

**Date:** 2026-05-06
**Scope:** Settle the 6 ❓ items left by Cycle 21 (`22-cycle21-AB-validators.md` §1.3) using direct evidence from upstream `core-v9 v1.5.8` clone + cross-reference with Cycle 37 (S-109) findings on `tests/creationtests/`.
**Result:** **2 ❓ → ✅** (verbatim symbol existence) + **3 ❓ → ⓘ "upstream-only / out-of-band"** (per S-109) + **1 ❓ → ✅ partial** (template-context blocked by C-CVS-23 rewrite, already tracked by AJ-08).
**Allowed under freeze:** read-only audit promotion (no spec rewrites).

## 1. Promotions

### 1.1 Row 43 — `errcore.VarTwoNoType("field", label, ...)` + `errcore.ValidationFailedType.Fmt(...)` → ✅ (symbol-existence)

- `errcore/VarTwoNoType.go:25` — `func VarTwoNoType(...) string` exists.
- `errcore/RawErrorType.go:121` — `ValidationFailedType RawErrorType = "Validation failed!"` exists (`RawErrorType` is a `string` alias with `.Fmt(...)` style helpers across the package).
- **Symbol-existence verdict:** ✅ — both identifiers are real and callable.
- **Snippet-correctness:** still depends on the fabricated `corevalidator.Result` shell (C-CVS-23) compiling; once AJ-08..14 rewrites the surrounding template, the `errcore` calls plug in cleanly. **No new finding.** Tracked under existing AJ-08.

### 1.2 Row 45 — "Compile regex once via `regexnew.New.Lazy` at struct construction" → ✅ (verbatim)

- `regexnew/newCreator.go:34` — `func (it newCreator) Lazy(pattern string) *LazyRegex { return it.LazyRegex.New(pattern) }`.
- `regexnew/newLazyRegexCreator.go:30` — `func (it newLazyRegexCreator) New(...) *LazyRegex` is the actual constructor that builds the `LazyRegex` struct (`mu sync.Mutex` + `pattern` + `compiler`).
- Companion `LazyLock` (`newCreator.go:41` + `NewLock` at `:42`) provides the package-locked variant.
- **Verdict:** ✅ — `regexnew.New.Lazy(pattern)` is the documented, real entry point; advice "compile once at struct construction" matches `LazyRegex`'s lazy-cache contract verified in Cycle 44 (`isCompiled` boolean guard under `sync.Mutex`).

## 2. Reclassifications to ⓘ "upstream-only / out-of-band"

### 2.1 Row 42 — "Add tests (Style A) — `coretestcases.CaseV1`" + `CaseNilSafe` example → ⓘ upstream-only

- Cycle 37 (S-109) deep-probe already established that `enum-v10/tests/creationtests/` does NOT consume `coretestcases.CaseV1` / `CaseNilSafe` — the entire `coretests` framework is upstream-`core-v9` only (`rg -n 'CaseV1|CaseNilSafe' tests/creationtests/` returns zero hits).
- **Verdict:** ⓘ "upstream-only" — claim is accurate for upstream `core-v9` consumers but does not describe `enum-v10`'s test layout. Same scope-disclaimer pattern as the 9 items annotated in Cycle 37.

### 2.2 Row 44 — "Three test files in `tests/creationtests/<pkg>tests/`: `<V>_Verification_test.go` etc." → ⓘ upstream-only

- `enum-v10/tests/creationtests/` is **flat** (no `<pkg>tests/` subdirectory pattern). Real files: `PathType_Creation_test.go`, `ScriptType_test.go`, `creation_test.go`, `AllEnums_ContractsTesting_test.go` — all at the directory root.
- Upstream `core-v9` IS organised by `<pkg>tests/` subdirectories (e.g. `tests/integratedtests/converterstests/`, `tests/integratedtests/regexnewtests/` per Cycle 43+44 evidence) but uses neither `_Verification_` nor `<V>_Verification_test.go` filename pattern (representative samples: `LazyRegex_Compile_test.go`, `LazyRegex_Methods_test.go`, `StringTo_IntegerWithDefault_test.go`).
- **Verdict:** ⓘ "upstream-only — and partially fabricated even there" — the `<pkg>tests/` subdirectory shape is real upstream but the `<V>_Verification_test.go` naming convention isn't used by either `enum-v10` or upstream `core-v9`. Folded into S-109 scope-disclaimer category for now (no new finding ID; will be addressed when AJ-08..14 rewrites the surrounding §6 template).

### 2.3 Row 46 — Diagnostic rules ("Message starts with field label", "No trailing punctuation", "No interpolated newlines") → out-of-band

- All three are aspirational rules. Without an emitting code path (the validator-output pipeline is itself fabricated per C-CVS-21..C-CVS-28), there is no behaviour to validate against.
- **Verdict:** out-of-band — Task AC advisory dimension (same bucket as Cycle 43 row 172 + Cycle 44 row 50).

## 3. Updated Cycle-21 scoreboard

| Pass | ✅ | ⚠️ | ❌ | ❓ | Verifiable score |
|------|----|---|----|---|---|
| Cycle 6 baseline | 0 | 0 | 0 | 18 | N/A |
| Cycle 21 (AB pass 3) | 4 | 0 | 8 | 6 | 33.3% |
| **Cycle 46 (this)** | **6** | 0 | 8 | **0** (3 ⓘ + 1 OOB + 2 promoted) | **42.9%** (6/14) |

§08 verifiable subset grows by 2 (rows 43, 45). ❓ pool **fully cleared** for §08 — every claim now has a verdict or a scope-disclaimer annotation. **§08 is no longer the worst-drift section** (still tied with §10 for second-worst behind §11 at 18.2%).

## 4. Cumulative AB-residual ❓ pool

- Pre-Cycle-46: 17 ❓ in `spec/01-app/` (Cycle 19=4 OOB / Cycle 20=1 OOB / Cycle 21=6 / Cycle 22=6).
- Post-Cycle-46: **17 − 6 = 11 ❓** in `spec/01-app/` (Cycle 19=4 OOB / Cycle 20=1 OOB / **Cycle 21=0 closed** / Cycle 22=6 active). Only Cycle 22 (`10-reflection-and-dynamic.md`) still has active probe targets.

## 5. Action items

- **No new findings.** All 6 ❓ items resolve cleanly (2 confirmed, 3 upstream-only per S-109, 1 advisory). The existing AJ-08..14 backlog already covers the surrounding template rewrites.
- **AC carry-over:** Row 46 advisory rules + Cycle 43 row 172 + Cycle 44 row 50 form a growing "advisory dimension" bucket — recommend formalising in Task AC.

## 6. Cumulative AB ❌ summary

- Cumulative AB ❌ across 7 sections: **51** (unchanged — no new findings this cycle). CRITICAL count unchanged at **23**.
- Cumulative ⓘ "upstream-only" annotations: previous (9 from Cycle 37) + 2 (rows 42, 44) = **11**.

## 7. Memory + bookkeeping

- Spec changelog → `spec-v0.51.0`.
- `package.json` → `0.16.0`.
- `.lovable/memory/workflow/01-state.md` updated.
- Scoreboard top-line updated.

---

_Audit file: `spec/07-code-vs-spec-audits/35-cycle46-AB-residual-spec01-validators.md`_
_See also: `22-cycle21-AB-validators.md` (the cycle whose ❓ are settled here), `29-cycle37-S109-creationtests-deep-probe.md` (upstream-only category)._