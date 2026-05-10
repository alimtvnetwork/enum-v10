# Cycle 21 — AB pass 3: `spec/01-app/08-validators.md` verification

**Date:** 2026-05-06
**Auditor:** AI agent (Lovable)
**Trigger:** Task **AB** pass 3 — upstream `core-v9 v1.5.8` cloned at `/tmp/core-v9-upstream`.
**Scope:** Promotion pass for the 18 unverifiable claims left by Cycle 6 (`08-validators.md`).

> 🧊 **Freeze interaction:** `spec/01-app/` remains DRIFT-FROZEN. The findings below are recorded but **NOT patched** until the user lifts the freeze.

> ⚠️ **This is the worst-drift section audited so far.** Almost the entire chapter describes an API that does not exist in upstream `core-v9`. The fluent-builder pattern (`corevalidator.New.Line.NotEmpty().MaxLength().Build()`), the `Result` type, the `IsValid()/IsSuccess()/IsFailed()/Message()/Error()` contract, the `Validate(input)` method, and the `RangeValidator`/`StringCompareAs` validators are all fabricated. Real upstream `corevalidator` exposes `LineValidator{LineNumber, TextValidator}` with `IsMatch(lineNumber, content, isCaseSensitive) bool` and side-channel constructors like `NewSliceValidatorUsingErr(...)`.

## 1. Verification matrix

Cycle 6 baseline: 18 ❓, no verifiable subset.

### 1.1 Promoted to ✅ (4 claims)

| # | Spec line(s) | Claim | Evidence in upstream |
|---|--------------|-------|----------------------|
| 1 | 23 | `corevalidator` exists with line/slice/text validators | `corevalidator/` directory with `LineValidator.go`, `SliceValidator.go`, `TextValidator.go` (23 files) |
| 2 | 72 | Import path `core-v9/corevalidator` | upstream `go.mod` |
| 3 | 150 | `errcore.MergeErrors(errs...)` exists | `errcore/MergeErrors.go:32` (uses `errors.Join`) |
| 4 | 196, 216 | `regexnew.New.Lazy(pattern)` | already verified in Cycle 20 (`regexnew/vars.go:44 New = newCreator{}`, `newCreator.go:34 Lazy(...)`) |

### 1.2 Demoted to ❌ Contradiction (8 claims — NEW findings, all stem from one root cause)

| ID | Spec line(s) | Spec claim | Upstream reality | Severity |
|----|--------------|------------|------------------|----------|
| **C-CVS-21** | 49-61 (the entire validator-contract table) | Every validator exposes `IsValid() bool`, `IsSuccess() bool`, `IsFailed() bool`, `Message() string`, `Error() error` ("`IsSuccessValidator` composition") | **None of these methods exist on `LineValidator`, `SliceValidator`, `TextValidator`.** Real surface is `IsMatch(...) bool` and assertion helpers. There is no `IsSuccessValidator` interface in `coreinterface/`. | **CRITICAL** — the central abstraction of the chapter does not exist. |
| **C-CVS-22** | 74-78, 98-102, 116-119 (all three fluent-builder snippets) | `corevalidator.New.Line.NotEmpty().MaxLength(80).Matches(...).Build()`; `corevalidator.New.Slice.MinLength(1).MaxLength(10).EachItem(lv).Build()`; `corevalidator.New.Range.Int.Min(0).MaxExclusive(100).Build()` | **No `corevalidator.New` var exists** (`grep "^var New\b" corevalidator/` → 0 hits). **No `NotEmpty/MaxLength/MinLength/Matches/Build/EachItem/Min/Max/MaxExclusive` methods exist anywhere in the package.** Real `LineValidator` is a struct with embedded `LineNumber` + `TextValidator` fields populated directly. | **CRITICAL** — three fabricated examples; zero of the methods exist. |
| **C-CVS-23** | 84-90 | `v.Validate(input)` returns a `Result` whose `IsValid()` and `Message()` work as shown | **No `Validate(string)` method on any validator.** Closest is `LineValidator.IsMatch(lineNumber, content, isCaseSensitive) bool`. There is no `Result` type in `corevalidator/`. | **CRITICAL** — return shape and method name both wrong. |
| **C-CVS-24** | 111-113, 116-123 | `RangeValidator` with `New.Range.Int.Min(0).MaxExclusive(100)` | **No `RangeValidator` type.** Closest names are `RangeSegmentsValidator` (segment-checking, not numeric Min/Max) and `RangesSegment`. Numeric range validation lives in `coremath/integerWithin.go`, `integer16Within.go`, etc. — not in `corevalidator/`. | **HIGH** — wrong package + fabricated type. |
| **C-CVS-25** | 125-127 | `StringCompareAs` is a "specialty validator" inside `corevalidator` | **No `StringCompareAs` type in `corevalidator/`.** There is a separate `enums/stringcompareas/` package containing the `Variant` enum used as a *parameter* by `SliceValidator` constructors — it's an enum, not a validator. | **HIGH** — wrong abstraction (parameter enum vs validator). |
| **C-CVS-26** | 138 | `[]corevalidator.Result{ ... }` | `Result` type does not exist in `corevalidator/` (already noted in C-CVS-23). | **HIGH** — depends on C-CVS-23. |
| **C-CVS-27** | 199-203, 233-261 | "Authoring a Custom Validator" template uses `*regexnew.Lazy` field type, embeds the `IsValid/IsSuccess/IsFailed/Message/Error` contract | The type is **`*regexnew.LazyRegex`**, not `*regexnew.Lazy` (`grep "regexnew.Lazy\b" core-v9-upstream/` → 0 hits; `LazyRegex` is the actual struct name in `regexnew/LazyRegex.go:34`). Plus the "implement the contract" section reproduces C-CVS-21's fabricated 5-method interface. | **HIGH** — published template would not compile and would teach a fabricated interface. |
| **C-CVS-28** | 312-327 | "Canonical error message format `<ValidatorLabel>: field=<name> value=<actual> reason=<short-reason>`" with verbatim worked examples (`LineValidator: field=username value="" reason=empty-not-allowed`, etc.) and the claim that "the PowerShell test runner parses validator output to attribute failures to source files" | No code in `corevalidator/` produces this format (validators don't emit prose at all — they return `bool` from `IsMatch`). The "PowerShell parses validator output" pipeline does not exist in `scripts/*.psm1` either (the runner parses `go test` output, not validator output). The 4 worked examples are aspirational, not extracted "verbatim from existing validators" as the spec claims. | **HIGH** — invents an attribution pipeline that doesn't exist; misleads anyone trying to debug an attribution failure. |

### 1.3 Remain ❓ (6 claims — out of scope or behavioural)

| Spec line | Claim | Why still ❓ |
|-----------|-------|--------------|
| 263-297 | "Add tests (Style A) — `coretestcases.CaseV1`" + `CaseNilSafe` example | Behavioural test pattern — would need `coretests` upstream verification (covered by future Cycle 24/25). |
| 240, 259 | `errcore.VarTwoNoType("field", label, ...)`, `errcore.ValidationFailedType.Fmt(...)` | `VarTwoNoType` exists (`errcore/VarTwoNoType.go:25`), `ValidationFailedType` exists (`RawErrorType.go:121`) — the *snippet's correctness* depends on C-CVS-23's fabricated `Result` type compiling first; treat as ❓ until the surrounding template is rewritten. |
| 347-353 | "Three test files in `tests/creationtests/<pkg>tests/`: `<V>_Verification_test.go` etc." | Naming convention claim; partially verifiable from `enum-v10/tests/creationtests/` shape. Defer. |
| 367 | "Compile regex once via `regexnew.New.Lazy` at struct construction" | Behavioural advice — verified that `regexnew.New.Lazy` exists; the *advice* is sound but trivially advisory. |
| 308-310 | Diagnostic rules: "Message starts with field label", "No trailing punctuation", "No interpolated newlines" | All three are aspirational rules; without an emitting code path (C-CVS-28) they cannot be validated against reality. |

## 2. Updated cycle 6 scoreboard line

Cycle 6 was: 0 ✅ / 1 ⚠️→0 / 0 ❌ / 18 ❓ — 0 % / 100 % verifiable swing.

After Cycle 21 pass 3: **4 ✅ / 0 ⚠️ / 8 ❌ / 6 ❓** → verifiable score = 4 / 12 = **33.3%** *(verifiable)*. **Lowest verifiable score in the project.**

## 3. Action items spawned (BLOCKED by freeze)

- **AJ-08** rewrite §1 (overview + contract table) against the real `LineValidator{LineNumber, TextValidator}` + `IsMatch(...) bool` shape; remove fabricated `IsSuccessValidator` reference.
- **AJ-09** rewrite §2.1, §2.2, §2.4 (LineValidator/SliceValidator/RangeValidator examples) — likely requires *re-categorizing* RangeValidator as `coremath.integerWithin` (different package).
- **AJ-10** rewrite §2.5 to clarify `StringCompareAs` is the `enums/stringcompareas.Variant` enum, used as a *parameter* by `SliceValidator` constructors.
- **AJ-11** rewrite §3.1 to drop the fabricated `corevalidator.Result` slice pattern.
- **AJ-12** rewrite §4 ("Authoring a Custom Validator") — template currently teaches a fabricated 5-method contract; either (a) rewrite to teach the real `IsMatch(...) bool` contract, or (b) move the prescriptive template to `spec/06-testing-guidelines/` and clearly label it as "this project's preferred validator shape, not the upstream `corevalidator/` API".
- **AJ-13** rewrite §5 to remove the fabricated "PowerShell runner parses validator output" pipeline OR back it with a real script reference.
- **AJ-14** typo fix throughout: `*regexnew.Lazy` → `*regexnew.LazyRegex` (same drift as Cycle 20 §C-CVS-?? but distinct here).

## 4. Cumulative AB-pass running totals

| Cycle | Section | ❓ before | ✅ after | ❌ after | ❓ after | Score |
|-------|---------|-----------|----------|----------|----------|-------|
| 19 | `09-converters.md` | 23 | 10 | 5 | 8 | 66.7 % |
| 20 | `07-conditional-and-utilities.md` | 17 | 12 | 5 | 3 | 70.6 % |
| 21 | `08-validators.md` | 18 | 4 | **8** | 6 | **33.3 %** |
| **Σ** | (3 sections) | **58** | **26** | **18** | **17** | — |

**Pattern update:** Across 3 sections, fabrication rate is now **18/(26+18) = 41 %** of all *verifiable* claims, not the 25 % projected after Cycle 20. Validators chapter alone contributes 8 ❌ — far above the 5 / cycle running average — because the entire chapter was written in a "wishful API" style describing what an idiomatic Go validator package *should* look like, not what `corevalidator/` actually is.

**Revised projection:** Expect 5–10 more ❌ across the remaining 4 sections (§10, §11, §15, §16) — total project ❌ likely lands in the **25–30 range** by AB pass 7.

## 5. Suggestion link

S-106 (`spec-api-check.psm1`) would have caught **all 8** ❌ in this cycle automatically — every fabricated symbol would fail `go vet`. **Strongly re-iterating priority of S-106 before any AJ rewrites land**, otherwise AJ-08..14 risk inventing a fresh layer of fabrications.

---

_Audit file: `spec/07-code-vs-spec-audits/22-cycle21-AB-validators.md`_
_See also: `21-cycle20-AB-conditional-and-utilities.md` (pass-2 precedent), `07-cycle6-validators.md` (the cycle whose ❓ are being promoted here)._