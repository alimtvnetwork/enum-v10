# Cycle 19 — AB pass 1: `spec/01-app/09-converters.md` verification

**Date:** 2026-05-06
**Auditor:** AI agent (Lovable)
**Trigger:** Task **AB** — upstream `core-v9 v1.5.8` source successfully cloned from `https://github.com/alimtvnetwork/core-v9` into `/tmp/core-v9-upstream` (sandbox-local).
**Scope:** First pass of ❓→✅/⚠️/❌ promotions for the 23 unverifiable claims left by Cycle 7.

> 🧊 **Freeze interaction:** `spec/01-app/` is DRIFT-FROZEN (spec-v0.30.0). This audit is allowed under the explicit AB-promotion exception. New ❌ contradictions are recorded here but **NOT** patched into `09-converters.md` until the user lifts the freeze (or grants a one-shot AB-fix waiver).

## 0. Pre-flight discoveries

While bootstrapping the AB clone, two **Core-memory drift items** surfaced and were corrected in `mem://index.md` in this same turn:

| ID | Stale memory claim | Reality (verified in sandbox) | Action |
|---|---|---|---|
| **M-CVS-01** | "This is a Go enum library (enum-v10, module path `github.com/alimtvnetwork/enum-v10`...)" | `go.mod` declares `module github.com/alimtvnetwork/enum-v10` | Memory line 4 corrected to `enum-v10`. |
| **M-CVS-02** | "upstream `core-v9` repo's own `go.mod` STILL declares `module github.com/alimtvnetwork/core-v8`. ... `go.mod` carries `replace ...`" | Upstream `core-v9 v1.5.8` declares `module github.com/alimtvnetwork/core-v9`. `go.mod` no longer carries the `replace` bridge (cleared by Task AG). | Memory line 6 rewritten to "rename complete; bridge removed." |

Both confirm Tasks **W** and **AG** are legitimately Done. The previous memory text was a snapshot frozen from a pre-AG cycle.

## 1. Verification matrix — `09-converters.md`

Total claims sampled in Cycle 7: **23 ❓** (no verifiable subset at the time).

### 1.1 Promoted to ✅ (10 claims)

| # | Spec line(s) | Claim | Evidence in upstream |
|---|--------------|-------|----------------------|
| 1 | 21 | `converters/` provides struct-as-namespace surface | `converters/vars.go` exposes `StringTo`, `BytesTo`, `PrettyJson`, `AnyTo`, `StringsTo`, `Map`, etc. as struct singletons |
| 2 | 26 | Import path `github.com/alimtvnetwork/core-v9/converters` | upstream `go.mod` line 1: `module github.com/alimtvnetwork/core-v9` |
| 3 | 29 | `converters.StringTo.Integer(s)` returns `(int, error)` | `stringTo.go:154` `func (it stringTo) Integer(input string) (value int, err error)` |
| 4 | 33 | `StringTo.IntegerWithDefault(s, def)` returns `(int, ok)` | `stringTo.go:39` matches signature `(value int, isSuccess bool)` |
| 5 | 36 | `StringTo.Float64(s)` returns `(float64, error)` | `stringTo.go:203` matches |
| 6 | 40 | `StringTo.Byte(s)` returns `(byte, error)` | `stringTo.go:260` matches |
| 7 | 53 | `BytesTo.String(b)` direct cast | `bytesTo.go:37` `func (it bytesTo) String(...)` |
| 8 | 70 | `typesconv/` package exists | `/tmp/core-v9-upstream/typesconv/` present (5 source files) |
| 9 | 73 | Import path `github.com/alimtvnetwork/core-v9/typesconv` | confirmed by upstream `go.mod` |
| 10 | 102 | `errcore.FailedToConvertType` referenced by converters | `stringTo.go` imports `core-v9/errcore` |

### 1.2 Demoted to ❌ Contradiction (5 claims — NEW findings)

These spec lines describe APIs that **do not exist upstream**.

| ID | Spec line | Spec claim | Upstream reality | Severity |
|----|-----------|------------|------------------|----------|
| **C-CVS-11** | 30 | `converters.StringTo.Integer64("...")` | No `Integer64` method on `stringTo`. Only `Integer`, `IntegerMust`, `IntegerDefault`, `IntegerWithDefault`, `IntegersWithDefaults`, `IntegersConditional`. | **HIGH** — code that follows the spec will not compile. |
| **C-CVS-12** | 37 | `converters.StringTo.Float32("3.14")` | No `Float32` method. Only `Float64`, `Float64Must`, `Float64Default`, `Float64Conditional`. | **HIGH** — same. |
| **C-CVS-13** | 47 | `converters.StringTo.Bool("true")` (with the lengthy stdlib-`strconv.ParseBool` contract paragraph 42-46) | No `Bool` method on `stringTo`. The stdlib-equivalent surface lives at `typesconv.StringToBool` / `typesconv.StringPointerToBool` / `typesconv.StringToBoolPtr` instead — and those return `bool`, not `(bool, error)`. | **HIGH** — both the API path *and* return shape are wrong. |
| **C-CVS-14** | 60-61 | `converters.PrettyJson.String(jsonBytes)`, `converters.PrettyJson.FromAny(myStruct)` | `PrettyJson` is bound to `jsoninternal.Pretty` whose surface is `PrettyString`, `PrettyStringIndent`, `SafePrettyString`, `PrettyStringDefault`, `PrettyStringDefaultMust` — no `.String` and no `.FromAny`. | **MEDIUM** — wrong method names; corrected names exist nearby. |
| **C-CVS-15** | 76, 79, 82, 159 | `typesconv.IntToInt64(i)`, `typesconv.Int64ToInt32(v)`, `typesconv.Float64ToInt(f)` | The actual `typesconv` surface is `*Ptr`, `*PtrToSimple`, `*PtrToSimpleDef`, `*PtrToDefPtr`, `*PtrDefValFunc` — pointer-tagging/defaulting helpers, NOT numeric widening/narrowing. **Entire §2 example block + §4.3 example are fabricated.** | **CRITICAL** — every call in the example will fail to compile and the whole "when to use `typesconv` vs `converters`" guidance table at §2 is built on this fiction. |

### 1.3 Remain ❓ (8 claims — pending deeper probe)

| Spec line | Claim | Why still ❓ |
|-----------|-------|--------------|
| 31 | `Integer64("9223372036854775807")` overflow behaviour | Method does not exist; no replacement to probe. |
| 54 | `BytesTo.PrettyJsonString(jsonBytes)` | Need to grep `bytesTo.go` for that method (next pass). |
| 64 | "PrettyJson namespace duplicates a subset of `corejson`" | Requires reading `corejson` to confirm overlap claim. |
| 99-107 | Conversion safety contract: "no panics", "errcore wrapped", "locale-independent" | Requires reading every converter method body — defer to focused contract pass. |
| 119 | `IntegerWithDefault(queryParam, 25)` returns `(int, _)` | Already promoted #4 above for the signature; the *behavioural* claim "fall back to default" needs a unit-style trace. |
| 130 | `parsePagination` example signature | Behavioural; depends on contract pass. |
| 161 | `errcore.OverflowType.Fmt(...)` exists | Need to grep `errcore` package. |
| 172 | "Using `*WithDefault` then re-validating hides the malformed input" | Behavioural / advisory — never strictly verifiable from code. |

## 2. Updated cycle 7 scoreboard line

Cycle 7 (`09-converters.md`) was: 0 ✅ / 0 ⚠️ / 0 ❌ / 23 ❓ — N/A score.

After Cycle 19 pass 1: **10 ✅ / 0 ⚠️ / 5 ❌ / 8 ❓** → verifiable score = 10 / 15 = **66.7%** *(verifiable)*.

This is the **first non-100 % verifiable score in the project**, and it lands directly because we now have ground truth instead of "no verifiable subset."

## 3. Action items spawned

- **AJ-01 (NEW, BLOCKED by freeze):** rewrite `09-converters.md` §1.1 to drop `Integer64`, `Float32`, `Bool` and replace with the actually-existing methods. Estimated 30 lines.
- **AJ-02 (NEW, BLOCKED by freeze):** rewrite §2 + §4.3 around the real `typesconv` `*Ptr*` surface (or relocate the numeric-widening guidance into a new spec section that references stdlib `int(x)` casts).
- **AJ-03 (NEW, BLOCKED by freeze):** correct `PrettyJson.String` → `PrettyString`, `PrettyJson.FromAny` → `PrettyString` (or `SafePrettyString` for the panic-free variant) at §1.3 + §4 callsites.
- **AC carryover:** Cycle 7 (this file) and Cycle 9 (consistency dimension) should be re-run after AJ-01..03 land.
- **AB-pass-2 (planned):** apply the same promotion-and-grep pass to Cycles 5 (`07-conditional-and-utilities.md`), 6 (`08-validators.md`), 8 (`10-reflection-and-dynamic.md`), and 9 (`11-versioning.md`). Each is expected to surface a similar number of ❌ findings because all four were authored against the pre-rename `core-v8` surface.

## 4. Suggestion (S-106 NEW)

The five ❌ findings in this file all share a single root cause: **the spec was authored against an internal mental model of the API, not against grep output**. Recommend adding an opt-in `scripts/spec-api-check.psm1` runner that, for every `converters.*`, `typesconv.*`, etc. code-fenced block in `spec/01-app/`, executes `go vet` against a synthesised test file. This would have caught all five ❌ findings automatically. Logged in `.lovable/memory/suggestions/01-suggestions-tracker.md`.

---

_Audit file: `spec/07-code-vs-spec-audits/20-cycle19-AB-converters-promotion.md`_
_See also: `15-cycle14-security.md` (last `spec/01-app/` cycle before freeze), `08-cycle7-converters.md` (the cycle whose ❓ are being promoted here)._