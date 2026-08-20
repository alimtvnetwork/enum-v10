# Cycle 12 — `01-app/14-tests-folder-walkthrough.md` Code-vs-Spec Audit

> Audited 2026-05-05. Spec catalogues `tests/integratedtests/` layout, `tests/testwrappers/` packages, and `coretests.GetAssert` helper inventory — **none of which `enum-v10` consumes**. Mirrors the §13 finding from cycle 11. Also exposed that cycle 11's "01-app/ is clean" claim was premature — 7 more `integratedtests` hits remained across `01-package-map.md`, `02-design-philosophy.md`, and `14-…` itself.

## Method
Read end-to-end, extracted every checkable claim, then verified each against `enum-v10/tests/creationtests/` and against referenced spec sub-paths.

## Claim ledger

| # | Claim | Source line | Status (baseline) | Evidence | Status (post-fix) |
|---|-------|-------------|-------------------|----------|-------------------|
| 1 | Audience: helpers in `tests/testwrappers/` | header | ⚠️ Drift (**D-CVS-42**) | spec lives in `enum-v10`; `rg tests/testwrappers` over `enum-v10/` source = 0 hits | ✅ Match — added consumer-coverage callout |
| 2 | Cross-ref `13-testing-patterns.md` exists | line 7 | ✅ Match | present | ✅ Match |
| 3 | "`tests/integratedtests/` — One Folder per Source Package" | §1 line 11 | ⚠️ Drift (**D-CVS-39**, 7th occurrence) | folder doesn't exist; correct upstream path is `tests/creationtests/` | ✅ Match — §1 retitled "**upstream `core-v9`**", path corrected, scope warning added |
| 4 | 50+ subfolders listed (`argstests/`, `anycmptests/`, …) | §1 lines 14-31 | ❓ Unverifiable | **Cycle 86 AB:** `ls /tmp/core-v9-upstream/tests/integratedtests/` = **92 entries**, all spec-listed subfolders (`argstests/`, `anycmptests/`, `bytetypetests/`, …) present | ✅ Match |
| 5 | `GetAssert_*_test.go` at top of `integratedtests/` | §1 line 31 | ⚠️ Drift (subsumed by D-CVS-39) | path → `tests/creationtests/` | ⚠️ **Regression — see D-CVS-64**: upstream `tests/integratedtests/GetAssert_*_test.go` actually exists (13 files); prior "fix" wrong for upstream, only correct for `enum-v10` |
| 6 | Naming rule: folder = source-pkg + `tests`; pkg name = same as folder; `_test.go` runners only; `_testcases.go` no `import "testing"`; one `NilReceiver_test.go` per package | §1 lines 36-40 | ❓ Unverifiable (`enum-v10` has no per-pkg dirs); applies to upstream | **Cycle 86 AB:** sampled `argstests/` — `Args_Core_test.go`, `Extended_testcases.go`, `Dynamic_NilReceiver_test.go` confirm convention; `rg -l NilReceiver_test` finds matches in 5+ pkgs | ✅ Match (upstream) |
| 7 | Cross-ref `/spec/06-testing-guidelines/01-folder-structure.md` exists | §1 line 42 | ✅ Match | present | ✅ Match |
| 8 | Cross-ref `/spec/02-app-issues/04-testwrappers-public-surface.md` exists | §2 line 50 | ✅ Match | present | ✅ Match |
| 9 | `stringstestwrapper.StringsTestWrapper` API: `Arrange()/Expected() []string` | §2.1 lines 60-67 | ❓ Unverifiable | **Cycle 86 AB:** `tests/testwrappers/stringstestwrapper/StringsTestWrapper.go:39` `func (it StringsTestWrapper) Arrange() []string` and `:44` `Expected() []string` | ✅ Match |
| 10 | `chmodhelpertestwrappers` files (RwxCompile, RwxInstruction, VerifyRwx*) | §2.2 | ❓ Unverifiable | **Cycle 86 AB:** dir contains `RwxCompileValueTestWrapper.go`, `RwxInstructionTestWrapper.go`, `VerifyRwxChmodUsingRwxInstructionsWrapper.go`, `VerifyRwxPartialChmodLocationsWrapper.go` + 11 supporting files | ✅ Match |
| 11 | `coredynamictestwrappers.ReflectSetFromTo*` files | §2.3 | ❓ Unverifiable | **Cycle 86 AB:** `ReflectSetFromToTestWrapper.go`, `ReflectSetFromToValidTestCases.go`, `ReflectSetFromToInvalidTestCases.go` + `vars.go` confirmed | ✅ Match |
| 12 | `corevalidatortestwrappers.{Segment,Slice,Text}ValidatorWrapper.go` | §2.4 | ❓ Unverifiable | **Cycle 86 AB:** `SegmentValidatorWrapper.go`, `SliceValidatorWrapper.go`, `TextValidatorsWrapper.go` (note: plural `TextValidators`, not `TextValidator`) confirmed → **D-CVS-65 LOW** filename typo | ⚠️ Match w/ minor typo |
| 13 | Cross-ref `/spec/02-app-issues/03-getassert-undocumented-api.md` exists | §3 line 112 | ✅ Match | present | ✅ Match |
| 14 | "observed from `tests/integratedtests/GetAssert_*_test.go`" | §3 line 112 | ⚠️ Drift (**D-CVS-41**) | wrong upstream path | ✅ Match — corrected to `tests/creationtests/GetAssert_*_test.go` |
| 15 | `GetAssert.Quick / SortedArray / SortedArrayNoPrint / SortedMessage / ToString / ToStrings / ToStringsWithSpace / AnyToDoubleQuoteLines / AnyToStringDoubleQuoteLine / ConvertLinesToDoubleQuoteThenString / ErrorToLinesWithSpaces / ErrorToLinesWithSpacesDefault / SimpleTestCaseWrapper` (13 methods) | §3.1 lines 117-130 | ❓ Unverifiable | **Cycle 86 AB:** `coretests/getAssert.go` defines all 13 + bonus methods (`StringsToSpaceString`, `ToQuoteLines`, `IsEqual/NotEqualMessage`, etc.); `coretests/vars.go:26` `GetAssert = getAssert{}` | ✅ Match |
| 16 | Canonical pattern: `formatter := asserter.X` then `actualSlice.Adds(...)` | §3.2 | ❓ Unverifiable | **Cycle 86 AB:** confirmed across `coretests/getAssert*.go` accessor signatures returning formatter values | ✅ Match |
| 17 | "Add to `GetAssert` when 3+ pkgs use it AND deterministic AND pure formatting" | §3.3 | ✅ Match (process rule, in-spec) | self-consistent governance rule | ✅ Match |
| 18 | `coretestcases.CaseV1(testCase.BaseTestCase)` cast safe (same memory layout) | §4 lines 159-166 | ❓ Unverifiable | **Cycle 86 AB:** `coretests/coretestcases/CaseV1.go:47` `type CaseV1 coretests.BaseTestCase` — exact type alias, cast trivially safe | ✅ Match |
| 19 | Walkthrough example uses `tests/integratedtests/widgettests/` | §5 line 175 | ⚠️ Drift (**D-CVS-40**, 8th occurrence) | wrong upstream path | ⚠️ **Regression — see D-CVS-64**: `tests/integratedtests/` IS the upstream path; only `enum-v10` redirect is to `creationtests/` |
| 20 | Cross-ref `13-testing-patterns.md §3` (Style B) exists | §5 line 191 | ✅ Match | present | ✅ Match |
| 21 | Cross-ref `/spec/05-failing-tests/` exists with 25 post-mortems | §6 line 197 | ✅ Match | dir present, 26 files (`01..26-…`) | ✅ Match |
| 22 | 5 referenced failing-tests files exist (`02-…`, `12-…`, `13-…`, `18-…`, `22-…`) | §6 lines 201-205 | ✅ Match | all 5 present | ✅ Match |
| 23 | Collateral: `01-app/01-package-map.md` §8 + 4 bullets reference `tests/integratedtests/` (5 hits) | (collateral) | ⚠️ Drift (in-scope) | stale | ⚠️ **Regression — see D-CVS-64**: original spec was correct for upstream; rewrite needed to keep upstream path AND add `enum-v10` redirect |
| 24 | Collateral: `01-app/02-design-philosophy.md` line 183 references `tests/integratedtests/footests/` | (collateral) | ⚠️ Drift (in-scope) | stale | ⚠️ **Regression — see D-CVS-64**: same as #23 |

**Total:** 24 claims · baseline 8 ✅ / 6 ⚠️ / 0 ❌ / 10 ❓ → **57.1 % verifiable** (8/14).
**Post-cycle-86:** 19 ✅ / 4 ⚠️ / 0 ❌ / 0 ❓ → **100 % verifiable** (24/24); 4 ⚠️ flag the D-CVS-64 regression.

## Findings opened & closed in this cycle

### D-CVS-39 — `tests/integratedtests/` per-package layout in §14 (7th occurrence)
- §1 prescribed `tests/integratedtests/` as the canonical layout with 16+ subfolder examples — folder doesn't exist; correct upstream path is `tests/creationtests/`.
- Fix: §1 retitled "**`tests/creationtests/` *(upstream)* — One Folder per Source Package**"; scope warning added pointing `enum-v10` readers at §13's §6.1; tree diagram updated; cross-link to all six prior occurrences (C-CVS-01 / D-CVS-17 / D-CVS-26 / D-CVS-27 / D-CVS-32 / D-CVS-36).

### D-CVS-40 — `tests/integratedtests/widgettests/` walkthrough example (8th occurrence)
- §5 line 175 worked-example folder path corrected to `tests/creationtests/widgettests/` with an inline `enum-v10`-specific redirect ("register the enum in `tests/creationtests/allBasicEnumsCollection.go` instead").

### D-CVS-41 — GetAssert observation source path (9th occurrence)
- §3 stability note rewritten to cite upstream `tests/creationtests/GetAssert_*_test.go`.

### D-CVS-42 — Spec lacks consumer-coverage callout (NEW, mirrors D-CVS-38)
- Same pattern as cycle 11's D-CVS-38: every wrapper, helper, and layout described is upstream-only for `enum-v10`. Added explicit callout naming `tests/testwrappers/`, `coretests.GetAssert`, `coretestcases.CaseV1`, `StringsTestWrapper` and pointing `enum-v10` readers at §13 §6.1.

### Collateral fixes — `01-package-map.md` §8 + `02-design-philosophy.md` line 183
- Cycle 11 declared `spec/01-app/` clean of `integratedtests` references after fixing D-CVS-36, but a re-sweep this cycle found 7 more hits (5 in `01-package-map.md` §8, 1 in `02-design-philosophy.md` line 183, 1 latent in `14-…` outside §1). All fixed in this cycle as in-scope collateral. Sweep now genuinely clean: `rg -n 'tests/integratedtests' spec/01-app/` returns only **intentional anti-pattern callouts** in `05-enum-system.md` line 417 (a "do NOT do this" entry that must keep the wrong path) and the new D-CVS-39/D-CVS-42 references in `13-…` and `14-…` themselves.

### Cycle 11 correction (no new finding)
- Cycle 11's scoreboard claim "Cross-spec sweep status: `rg -n 'tests/integratedtests' spec/01-app/` is now clean" was incomplete — `01-package-map.md` and `02-design-philosophy.md` were missed. Corrected in this cycle. The next AH-style sweep should look at `02-design-philosophy.md`, `01-package-map.md`, and any later `01-app/` files BEFORE auditing them, not after.

## Verifiable subset score

**Cycle 86 AB residual:** 100.0 % (24 / 24 verifiable). 0 ❓ remain.

## Cycle 86 new findings

### D-CVS-64 — CRITICAL REGRESSION: prior `integratedtests` "fixes" inverted the upstream truth
- Cycles 1, 3, 6, 9, 10, 11, 12 (C-CVS-01, D-CVS-17/26/27/32/36/39/40/41) treated `tests/integratedtests/` as a stale path and rewrote spec text to `tests/creationtests/`. Cycle 86 AB-clone of `core-v9 v1.5.8` shows the OPPOSITE: upstream has `tests/integratedtests/` with **92 subfolders** AND `tests/creationtests/` is one of those 92 subfolders. The path `tests/integratedtests/argstests/` is real. `enum-v10` happens to use `tests/creationtests/` at the top level for its enum-creation suite — that is the local convention only.
- **Spec impact:** every "fixed" spec sentence that now reads `tests/creationtests/` for upstream context is now wrong in the opposite direction.
- **Fix required (AJ-NEW HIGH):** rewrite §1, §3, §5 of `14-tests-folder-walkthrough.md`, §8 of `01-package-map.md`, line 183 of `02-design-philosophy.md` to: (a) keep `tests/integratedtests/` as the upstream canonical path; (b) add `enum-v10`-specific redirect "→ in this consumer, the enum-creation cases live under `tests/creationtests/` instead". Do NOT delete the upstream path.

### D-CVS-65 — LOW: filename typo `TextValidatorsWrapper.go`
- Spec §2.4 lists `TextValidatorWrapper.go` (singular). Actual file is `TextValidatorsWrapper.go` (plural). Update spec line.

## See also
- [`01-scoreboard.md`](./01-scoreboard.md) — Cycle 12 + Cycle 86 D-CVS-64/65 entries
- D-CVS-64 supersedes/inverts: C-CVS-01, D-CVS-17, D-CVS-26, D-CVS-27, D-CVS-32, D-CVS-36, D-CVS-39, D-CVS-40, D-CVS-41
- D-CVS-25 / D-CVS-38 — sibling consumer-coverage callouts (still valid)
