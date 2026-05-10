# 04 — Code-vs-Spec Cycle 3: Enum System

**Date:** 2026-05-04
**Spec audited:** [`spec/01-app/05-enum-system.md`](../01-app/05-enum-system.md)
**Method:** Manual claim extraction; automated `rg` verification across all 71 enum packages in this repo.
**Auditor:** Lovable agent (evidence-driven)

---

## Summary

| Metric | Value |
|---|---|
| Claims extracted | **18** |
| ✅ Match | **8** |
| ⚠️ Drift | **6** |
| ❌ Contradiction | **3** |
| ❓ Unverifiable | **1** |
| **Match rate (verifiable)** | **8 / 17 = 47.1 %** |

> **Verdict:** §05 was written assuming a "canonical 3-file layout" (`consts.go` + `vars.go` + `<Type>.go`) that **does not match a single enum package** in this repo. The actual convention is `Variant.go` (combining type + iota + predicates) plus `vars.go`, with the type **always named `Variant`** in 64 of 71 packages. The spec also mandates `Invalid` as the first iota constant, but **10 packages put a different value first** (and `inttype` uses `InvalidIndex` instead). The factory surface is also wider than documented.

---

## Claim-by-claim verification

### ✅ Match (8)

| # | Claim | Evidence |
|---|-------|----------|
| C1 | Three-layer architecture (`enuminf` interface → `enumimpl` engine → enum package) | All 71 enum packages import `coreinterface/enuminf` and `coreimpl/enumimpl` and embed a package-level `BasicEnumImpl` var |
| C2 | `enumimpl.New` is a struct-as-namespace factory hub | `rg 'enumimpl\.New\.\w+'` shows clean `enumimpl.New.BasicByte`, `.BasicInt8`, `.BasicString` access pattern |
| C3 | Backing types include `byte`, `int8`, `string` | Confirmed in enum-package type declarations |
| C4 | The `Ranges`-style name slice indexes by enum constant | All 64 packages with `vars.go` use the `IotaConst: "Name"` map-literal form |
| C5 | Domain predicates use `Is<Name>()` shape | E.g. `compressformats.IsZip()`, `iptype.IsV4()` — matches across all packages |
| C6 | Enums are value types compared with `==` | No `*Variant` comparison patterns; all `==` usage is direct |
| C7 | F-NEW-07: comprehensive `Value*()` accessors required by `BasicEnumValuer` | `rg '^func \(it \w+\) Value(Byte|Int8?|Int16|Int32|Int64|UInt16|String)\(\)' --type go .` returns dense matches across packages |
| C8 | Spelling: stdlib `UnmarshalJSON` (1 'l') vs engine helper `UnmarshallEnumToValue` (2 'l's) | Confirmed in every enum's `Variant.go`: both spellings co-exist exactly as documented |

### ⚠️ Drift — code is fine, spec prose is stale (6)

| # | Claim | Code reality | Severity | Fix path |
|---|-------|--------------|----------|----------|
| D-CVS-14 | "**Step 3 — `Status.go`** (method set)" — spec implies file is named after the type | Actual filename is **`Variant.go`** in 64 of 71 packages (the 7 exceptions: `osgroupexecution/Precedence.go`, `brackets/{Bracket,Category}.go`, `servicestate/Action.go`, `quotes/Quote.go`, `linuxservicestate/ExitCode.go`, `osarchs/Architecture.go`). All use the file's own type name. | LOW | Spec — note the convention: file is named `<TypeName>.go`, and the type is conventionally `Variant` |
| D-CVS-15 | Recipe shows separate `consts.go` and `<Type>.go` files | No package has a `consts.go`. The type, iota constants, and method set all live together in `<TypeName>.go` (typically `Variant.go`) | MED | Spec — collapse Step 1 + Step 3 into a single-file recipe matching reality |
| D-CVS-16 | Factory list (§6) only documents `UsingTypeSlice`, `Default`, `DefaultWithAliasMap`, `CreateUsingMap`, `CreateUsingMapPlusAliasMap` | Code uses 9 distinct methods incl. **undocumented**: `UsingFirstItemSliceAllCases`, `DefaultAllCases`, `DefaultWithAliasMapAllCases`, `UsingFirstItemSliceAliasMap`, `CreateUsingSlicePlusAliasMapOptions`, `CreateUsingStringersSpread`. `CreateUsingMap` is never called. | MED | Spec — add the `*AllCases` variants and document the "first item" form; remove unused `CreateUsingMap` from the recommendation table |
| D-CVS-17 | §8 says tests live in `tests/integratedtests/<pkg>tests/` with three test files per enum | Tests live in **`tests/creationtests/`** (single shared package) — same C-CVS-01 issue already fixed in §03. Per-enum test files do not exist; testing uses table-driven shared registry (`allBasicEnumsCollection.go`) | MED | Spec — rewrite §8 to match the shared-registry pattern, mirror C-CVS-01 fix from §03 |
| D-CVS-18 | Spec example uses `reflectinternal.TypeName(Invalid)` from `core-v9/internal/reflectinternal` | Same module-boundary problem as C-CVS-02: `enum-v9` packages **cannot** import `core-v9/internal/...`. Real packages call `enumimpl.New.BasicByte.UsingTypeSlice("Variant", ranges[:])` with a string literal, **not** a reflection helper. | MED | Spec — replace the `reflectinternal.TypeName(...)` example with the actual string-literal pattern, or use `DefaultAllCases(firstItem, ranges[:])` which derives the type name in the engine itself |
| D-CVS-19 | "Predicate file-split rule: split if >6 predicates OR any >20 lines" | Code keeps predicates in `Variant.go` even when the count exceeds 6 (`pathpatterntype/Variant.go` has 113 constants and dozens of predicates in one file). The "rule" is not enforced anywhere. | LOW | Spec — soften to a guideline; document the actual practice: predicates stay in `<TypeName>.go` regardless of count. |

### ❌ Contradiction — spec rule is actively violated by code (3)

| # | Claim | Reality | Severity | Resolution |
|---|-------|---------|----------|------------|
| C-CVS-03 | "First constant **must be `Invalid`** (or equivalent zero-value) so an unset variable is detectable" | **10 packages** put a non-Invalid value first: `compressformats` (`Zip`), `compresslevels` (`Default`), `envtype` (`Uninitialized`), `inttype` (`InvalidIndex` with explicit `= -1`), `logtype` (`Silent`), `revokereason` (`Unspecified`), `scripttype` (`Default`), `sqljointype` (`Default`), `strtype`, `taskpriority` (`Default`) | HIGH | **Either** (a) accept the de-facto rule "first constant is the **sentinel** — may be named `Invalid`, `Unspecified`, `Default`, `Uninitialized`, etc., but must occupy the zero value" and rewrite the spec, **or** (b) file 10 code-side issues. (a) is far cheaper and matches author intent — recommended. |
| C-CVS-04 | Recipe Step 2 imports `core-v9/internal/reflectinternal` | Cross-module `internal/` import is forbidden by Go and zero packages do this | HIGH | Same fix as D-CVS-18 — the spec is unrunnable as written |
| C-CVS-05 | "First constant must be `Invalid` (or equivalent zero-value)" — `inttype` declares `InvalidIndex Variant = -1` | `-1` is **not** the zero value of any unsigned/iota Go type; this directly contradicts "zero-value sentinel" | HIGH | Spec must either exclude signed-int enums from the rule, or document the alternate "`InvalidIndex = -1` for signed types" pattern explicitly |

### ❓ Unverifiable (1)

| # | Claim | Why unverifiable |
|---|-------|------------------|
| C9 | Asymmetric JSON: `MarshalJSON` always emits string name; `UnmarshalJSON` accepts string OR numeric-string | Verifying the *behavior* requires running the engine. The interface methods exist on every enum; the runtime contract is enforced inside `coreimpl/enumimpl` (upstream `core-v9`). Pending task **AB**. |

---

## Score

**Verifiable subset:** 17 claims (8 ✅ + 6 ⚠️ + 3 ❌) → **8/17 = 47.1 %**

The 3 ❌ contradictions are a stronger signal than the drift bucket and should be fixed first.

---

## Recommended remediation order (cheap → expensive)

1. **C-CVS-03** + **C-CVS-05** (HIGH): Reframe the "Invalid first" rule as "**sentinel first**" and add a sub-paragraph documenting the alternate names (`Default`, `Unspecified`, `Uninitialized`, `InvalidIndex`) and the `-1` form for signed types. ~10 min, fixes 2 of 3 contradictions. Bumps score 47.1 → 58.8.
2. **C-CVS-04** + **D-CVS-18** (HIGH+MED): Delete the `reflectinternal.TypeName(...)` example. Replace with the actually-used pattern: either string-literal type name (`UsingTypeSlice("Variant", ranges[:])`) or `DefaultAllCases(firstItem, ranges[:])`. ~10 min. Bumps score 58.8 → 70.6.
3. **D-CVS-14** + **D-CVS-15** (LOW+MED): Rewrite the recipe as a 2-file pattern (`Variant.go` + `vars.go`). Note the conventional type name `Variant`. ~15 min. Bumps score 70.6 → 82.4.
4. **D-CVS-16** (MED): Expand §6 factory table with the `*AllCases` family, drop `CreateUsingMap`. ~10 min. → 88.2.
5. **D-CVS-17** (MED): Rewrite §8 tests section to point at `tests/creationtests/` with the shared-registry pattern (mirror C-CVS-01 §03 fix). ~10 min. → 94.1.
6. **D-CVS-19** (LOW): Soften the predicate file-split rule to match `pathpatterntype` reality. ~5 min. → 100.0 verifiable.

Total remediation budget: ~60 minutes to lift §05 from 47.1 % → 100 % verifiable.

---

## See also

- [`02-cycle1-import-conventions.md`](./02-cycle1-import-conventions.md) — Cycle 1
- [`03-cycle2-error-system.md`](./03-cycle2-error-system.md) — Cycle 2
- [`01-scoreboard.md`](./01-scoreboard.md) — running totals
- [`../01-app/05-enum-system.md`](../01-app/05-enum-system.md) — spec under audit
