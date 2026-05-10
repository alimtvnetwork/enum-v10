# Active Plan — enum-v8

> Single source of truth for the project roadmap. Letter IDs are stable across sessions.
> Last updated: 2026-05-06 (Cycle 19 — AB pass 1 done; AJ-01..03 spawned, blocked by `spec/01-app/` freeze).

---

## Phase 1: Unblock Build (CRITICAL)

### W. Upstream `core-v9` `go.mod` rename + tag `v1.5.8`

- **Status:** ✅ Done (2026-05-05)

### AG. Drop the `replace` bridge and pin clean `core-v9 v1.5.8`

- **Status:** ✅ Done (2026-05-05)
- **Result:** `replace` directive removed, `require core-v9 v1.5.8` pinned.

### AM. Fix broken `core-v9` API calls (converter + coredynamic migration)

- **Status:** ✅ Done for reported blocker — `tests/creationtests` compile check passes in sandbox
- **Objective:** Update all `enum-v8` source files that use old `core-v8`-era function signatures (`converters.AnyToValueString`, `coredynamic.TypeName`, etc.) to the new struct-namespace API (`converters.AnyTo.ValueString`, `coredynamic.SafeTypeName`, etc.).
- **Dependencies:** None for the reported `creationtests` compile blocker.
- **Expected outputs:** All affected `.go` files updated, `go build ./...` passes.
- **Acceptance criteria:** `go build ./...` succeeds with `core-v9 v1.5.8`.
- **Reference:** `.lovable/memory/06-core-v9-api-migration.md`

---

## Phase 2: Spec Accuracy (HIGH)

### AA. Continue spec-audit cycles

- **Status:** ✅ DONE (2026-05-07, Cycle 20, `spec-v0.55.0` / enum-v8 v0.70.0)
- **Result:** All 6 walk-through targets audited. Final cycle (20) verified `spec/00-llm-integration-guide.md` against upstream `core-v9 v1.5.8` — zero drift, zero D-CVS/C-CVS findings.
- **Cycle plan (closed):**

| Cycle | Target | Status |
|-------|--------|--------|
| 15 | `spec/06-testing-guidelines/` (9 files) | ✅ Done |
| 16 | `spec/03-powershell-test-run/` (9 files) | ✅ Done |
| 17 | `spec/04-tooling/` (10 files) | ✅ Done |
| 18 | `spec/02-app-issues/` (10 files) | ✅ Done |
| 19 | `spec/05-failing-tests/` (25 files) | ✅ Done |
| 20 | `spec/00-llm-integration-guide.md` | ✅ Done |

### AH. Cross-`spec/` cleanup sweep

- **Status:** ✅ Done (2026-05-07, audit during v0.98.0)
- **Objective:** Replace all stale `tests/integratedtests/` references outside `cross-repo/`.
- **Result:** All 5 plan-listed files (`spec/06-testing-guidelines/01-folder-structure.md`, `spec/03-powershell-test-run/01-overview.md` + 3 sibling/portable files, `spec/04-tooling/04-bootstrap-into-new-repo.md`, `spec/02-app-issues/02-internal-package-coverage-policy.md`, `spec/00-llm-integration-guide.md`) now carry explicit upstream-`core-v9` scope notes. Remaining `integratedtests/` mentions are intentional upstream-layout callouts (verified line-by-line); `spec/05-failing-tests/` is intentionally frozen as upstream-historical reference (PI-002 res. + Cycle 19). No source-tree path drift remains.
- **Acceptance:** `rg integratedtests spec/` returns only callout / scope-note / upstream-reference lines.


### AB. Fetch upstream `core-v9` source

- **Status:** 🔄 In Progress (per-spec-directory residual passes)
- **Objective:** Get upstream `core-v9` source into a workspace path so auditor can verify ❓ claims.
- **Dependencies:** Fetch access to upstream repo ✅ (clone at `/tmp/core-v9-upstream`, tag `v1.5.8`)
- **Expected outputs:** ❓ claims promoted to ✅ or ❌, per audited directory.
- **Acceptance criteria:** ❓ count drops to <10 across all audited directories.
- **Progress:** §06 (`spec/06-testing-guidelines/`) — Cycle 48 (this session, 2026-05-06): 10 ❓ → 10 ✅, zero new findings.

### AC. Re-audit §07 and §09

- **Status:** ⏳ Pending — waits on AB
- **Objective:** Apply spec-internal-consistency dimension to the two baseline-only sections.
- **Dependencies:** AB
- **Expected outputs:** Updated audit reports, possible spec fixes.

### AI. Mark `spec/01-app/` as frozen

- **Status:** ✅ Done (2026-05-06, Cycle 61, `spec-v0.53.0`)
- **Result:** Freeze policy entry added to `spec/CHANGELOG.md`. All future `spec/01-app/` changes blocked behind explicit user thaw.

---

## Phase 3: Implementation

### AJ. Implement spec fixes from Cycle 15 findings

- **Status:** 📋 Planned
- **Objective:** Apply all drift/contradiction fixes found in Cycle 15 (`spec/06-testing-guidelines/`).
- **Dependencies:** AA (Cycle 15)
- **Expected outputs:** Updated spec files, scoreboard, changelog bump.

### AK. New enum package creation (template validation)

- **Status:** ✅ DONE (2026-05-07, v0.71.0)
- **Result:** Created `httpmethodtype/` package end-to-end (11 files) following `spec/00-llm-integration-guide.md` §10. Registered in `tests/creationtests/allBasicEnumsCollection.go`. `go build`, `go test`, and creationtests build all pass with go1.25.7 + `core-v9 v1.5.8`. Recipe verified complete; zero spec gaps encountered.

### AL. Test coverage expansion (umbrella)

- **Status:** 🔄 In Progress (broken into AL-01..AL-08 below)
- **Objective:** Lift `tests/creationtests/` coverage from baseline **15.5%** toward **≥60%** total statements.
- **Baseline (2026-05-06):** 15.5% total, all 73 packages individually <25%, 4154 functions <50%.
- **Strategy:** The existing `Test_AllEnums_ContractsTesting` only exercises ~6 methods per Variant. Each `Variant.go` exposes ~40–50 methods (`Json`, `JsonPtr`, `Format`, `IsAnyOf`, `IsAnyNamesOf`, `IsAnyValuesEqual`, `MarshalJSON`/`UnmarshalJSON`, `ValueInt8/16/32`, `MinByte`/`MaxByte`, `RangesByte`, `NameValue`, `ToPtr`, `AsJsoner`, etc.). Adding one shared-loop test per method family unlocks coverage across **all 73 packages simultaneously**.
- **Dependencies:** AG (clean Go build) ✅
- **Expected outputs:** New `*_test.go` files in `tests/creationtests/`, coverage delta recorded each pass.

#### AL-01. Json round-trip suite (highest leverage) ✅ DONE (2026-05-06, Cycle 49)
- **Target:** `MarshalJSON`, `UnmarshalJSON` (and `BasicEnumImpl.ToEnumJsonBytes` / `UnmarshallToValue`).
- **Approach:** Loop `allBasicEnumsCollection`; for each, marshal → re-unmarshal into the same pointer → assert `Name()` + `ValueString()` round-trip.
- **Result:** **15.5% → 21.6% total (+6.1pp)** with one new file `tests/creationtests/AllEnums_JsonRoundTrip_test.go`.
- **Findings:** PI-005 (sqliteconnpathtype double-quoted MarshalJSON) — type skipped via `jsonRoundTripSkipTypeNames` map.
- **Acceptance:** `./run.ps1 tc` green, total ≥ 21%. ✅

#### AL-02. Format & string conversion suite ✅ DONE (2026-05-06, Cycle 50)
- **Target:** `Format`, `Name`, `String`, `ValueString`, `ToNumberString`, `RangeNamesCsv`, `NameValue`, `MinValueString`, `MaxValueString`, `AllNameValues`.
- **Result:** **21.6% → 26.1% (+4.5pp)** with `tests/creationtests/AllEnums_Format_test.go`. Cumulative AL-01+02: **15.5% → 26.1% (+10.6pp)**.
- **Findings:** PI-006 (sqliteconnpathtype `NameValue` wrong fmt verb + empty `MinValueString`); strtype.Variant correctly excluded from min/max checks (free-form string enum).

#### AL-03. Comparison & predicate suite ✅ DONE (2026-05-06, Cycle 51)
- **Target:** `IsValid`, `IsInvalid`, `IsNameEqual`, `IsAnyNamesOf`, `ValueByte`, `ValueInt`, `ValueInt8/16/32`, `ValueUInt16` (numeric-width consistency).
- **Result:** **26.1% → 33.8% (+7.7pp)** with `tests/creationtests/AllEnums_Predicates_test.go`. Cumulative AL-01+02+03: **15.5% → 33.8% (+18.3pp)**.
- **Findings:** PI-007 (sqliteconnpathtype `IsAnyNamesOf()` empty-args returns true); strtype.Variant correctly excluded from numeric-width block (string-backed; `ValueByte` panics).

#### AL-04. Numeric width & range suite ✅ DONE (2026-05-06, Cycle 52)
- **Target:** `MinInt`, `MaxInt`, `MinMaxAny`, `MinValueString`, `MaxValueString`, `RangesDynamicMap`, `AllNameValues`, `IntegerEnumRanges`.
- **Result:** Test created at `tests/creationtests/AllEnums_NumericRange_test.go`. Coverage delta TBD on next `./run.ps1 -tc` run; expected +4–6pp.
- **Findings:** Reuses PI-006 skip for sqliteconnpathtype `MinValueString` (empty); strtype skipped from `MinInt <= MaxInt` invariant only (string-backed semantics differ).

#### AL-05. Constructor suite (`New`, `NewMust`, `RangesInvalidErr`, `Max`, `Min`) ✅ DONE pass-1 (2026-05-06, Cycle 53)
- **Pass 1 (this cycle):** 4 per-package suites shipped — `accesstype`, `certaction`, `completionstate`, `compressformats`. Each test calls `New(name)` for every known constant, asserts round-trip via `Name()`, asserts `New("__bogus__")` returns `(Invalid, non-nil err)`, runs `NewMust(name)` for every constant, and exercises `Max()`/`Min()`/`RangesInvalidErr()` where present. compressformats has unusual iota ordering (`Invalid = 5`) — its `Min<=Max` invariant is intentionally not asserted.
- **Files:** `accesstype/AccessType_Constructor_test.go`, `certaction/CertAction_Constructor_test.go`, `completionstate/CompletionState_Constructor_test.go`, `compressformats/CompressFormats_Constructor_test.go`.
- **Result:** Coverage delta TBD on next `./run.ps1 -tc`; pass-1 expected lift ≈ **+1.5–2.5pp** (4 of ~10 planned packages).
- **Pass 2 ✅ DONE (Cycle 54):** 6 additional per-package suites — `dbaction`, `envtype`, `iptype`, `onofftype`, `overwritetype`, `timeunit`. envtype uses `Uninitialized` as its zero-value sentinel (not `Invalid`); onofftype has a `newOtherWays` shorthand fallback (`yes`/`y`/`1` → On, etc.) that is also exercised. Pass-2 expected lift ≈ **+2–3pp**. Cumulative AL-05 (passes 1+2) covers 10 low-coverage packages.
#### AL-06. `quotes/` and `brackets/` dedicated suites ✅ DONE (2026-05-06, Cycle 55)
- **Why:** Both currently 7–12%, neither in `allBasicEnumsCollection`. Bespoke tests required.
- **Files:** `quotes/Quotes_WrapUnwrap_test.go` + `brackets/Brackets_WrapUnwrap_test.go`. Coverage targets exercised: `WrapWith`, `UnWrapWith`, `HasBothWrappedWith`, `WhichQuote`/`WhichBracket`, `getQuoteStatus`/`getSingleBracketStatus`, `Quote.Wrap`/`SelfWrap`/`IsEqual`/`GetOther`/`WrapAny`/`WrapString`/`WrapSkipOnExist`/`WrapRegardless`/`WrapFmtString`/`IsWrapped`/`UnWrap`/`WrapWithOptions`/`WrapAnySkipOnExist`, plus `Bracket.IsStart`/`IsEnd`/`IsParenthesis`/`IsCurly`/`IsSquare` and category/start/end variants, `Bracket.Pair`/`Category`/`OtherBracket`/`Value`, `Pair.Wrap`/`SelfWrap`. Boundary cases (empty, single-char, mismatched, left-only, right-only, plain).
- **Expected lift:** +1–2% total; lifts both packages from 7–12% into the 50–70% band.

#### AL-07. `strtype` / `inttype` constructor & GetSet suites ✅ DONE (2026-05-06, Cycle 56)
- **Why:** strtype 4.3%, inttype 10.3% — used by every other package via `IntType()`.
- **Files:** `inttype/IntType_Constructor_test.go` (covers `New`, `NewString`, `NewUInt`, `NewInt64`, `NewUsingJsonNumber`, `GetSet`, `GetSetVariant`, full `IsCompareResult` switch); `strtype/StrType_Constructor_test.go` (covers `New`, `NewUsingInteger`, `NewFileReader` smoke, `GetSet`, `GetSetVariant`).
- **Expected lift:** +2–4pp.

#### AL-08. `osdetect` selective coverage (cross-platform safe parts only) ✅ DONE (2026-05-06, Cycle 58)
- **Why:** 3% but bulk is platform-specific (`windows.go`, `linux.go`).
- **Target:** Only `IsRunningInDockerContainer` stub, `Variant`, `OperatingSystemDetail` JSON, `CurrentOsType` smoke.
- **File:** `osdetect/OsDetect_CrossPlatform_test.go` — Variant predicates + name/value/byte accessors, IsAnyOf, ToPtr round-trip, DefaultCmdProcessName branches, New/NewMust round-trip + error path, Variant JSON round-trip, OperatingSystemDetail JSON helpers smoke (zero value), CurrentOsType / CurrentOsMixTypes / CurrentOsTypesMap / IsCurrentOsTypesContains smoke, IsRunningInDockerContainer stability check.
- **Expected lift:** +1pp (osdetect goes from ~3% to ~25–35%; pure-Variant logic is now well covered).

**Combined target:** 15.5% → ~60% across AL-01..AL-08.

---

### AL2. Per-package coverage uplift, Phase 2 (NEW — top priority)

- **Status:** 📋 Planned — execute one sub-task per `next` command
- **Motivation:** After AL-01..AL-08, total coverage rose, but most individual packages still report **<50–60%**. The shared-loop matrix exercises `Variant` accessors but skips package-private surface: `New`/`NewMust` for the ~40 packages still without a constructor suite, `all-is-checkers.go` (`IsX()` predicate set per Variant), `all-validation-checking-err.go` (per-Variant validation error helpers), `Min`/`Max`/`RangesInvalidErr`, and package-bespoke files (`dbdrivertype/Connection*.go`, `osdetect/*` platform branches).
- **Strategy:** Group remaining ~40 packages into 6 batches of ~6–7 packages each. Per batch, ship one `<Pkg>_Coverage_test.go` per package using a small shared template (constructor round-trip for every Variant, every `IsX()` predicate, every `*Err()` helper, `Min`/`Max`/`RangesInvalidErr`, JSON round-trip if not already in matrix). Plus 2 bespoke sub-tasks for packages with non-template surface.
- **Dependencies:** AG ✅, AL-01..AL-08 ✅
- **Acceptance:** Each touched package reaches ≥60% statement coverage; total rises toward **≥70%**.
- **Sequencing (one per `next`):**

#### AL2-01. Batch A — simple state enums (6 pkgs) ✅ DONE (2026-05-06, Cycle 64)
- **Targets:** `compresslevels`, `configfilestate`, `conntrackstate`, `osgroupexecution`, `servicestate`, `sitestatetype`.
- **Surface:** `Create`/`CreateMust` (or `New`/`NewMust`), `all-is-checkers.go`, `all-validation-checking-err.go`, `Min`/`Max`/`RangesInvalidErr`.
- **Result:** 6 new test files (`<Pkg>_Coverage_test.go` per package) covering constructors + every `IsX` predicate + full Variant accessor sweep.

#### AL2-01b. Batch A test fix-up ✅ DONE (2026-05-06, Cycle 66)
- **Result:** Repaired 8 failing tests caused by misuse of diagnostic methods (`OnlySupportedErr`, `RangesInvalidErr` always non-nil) and brittle pinning of fuzzy-matched shorthand inputs (`onofftype`). Added `compresslevels.Variant` to `RangesDynamicMap` skip list. Documented as RCA Patterns 5 & 6.
- **Notes:** `compresslevels` and `configfilestate` have no `New`/`NewMust` (Variant-only); `servicestate.Min()` returns `Status` not `Invalid` (verified).

#### AL2-02. Batch B — DB family (6 pkgs)
- **Targets:** `dbexposetype`, `dbuserprivillegetype`, `sqljointype`, `sqliteconnpathtype`, `querymethodtype`, `resauthtype`.
- **Surface:** Same template; sqliteconnpathtype already partly covered — focus on remaining `IsX`/`*Err` helpers.
- **Expected lift:** ~+3pp.

#### AL2-03. Batch C — networking / IP (5 pkgs)
- **Targets:** `inputiptype`, `protocoltype`, `nginxlogtype`, `pathpatterntype`, `verifiertriggertype`.
- **Surface:** Template + any `pathpatterntype` regex/match helpers.
- **Expected lift:** ~+2.5pp.

#### AL2-04. Batch D — Linux / OS (6 pkgs)
- **Targets:** `linuxservicestate`, `linuxtype`, `linuxvendortype`, `osarchs`, `packageinstallmethod`, `runtype`.
- **Surface:** Template; verify `osarchs` width helpers.
- **Expected lift:** ~+3pp.

#### AL2-05. Batch E — misc value enums (7 pkgs)
- **Targets:** `eventtype`, `instructiontype`, `leveltype`, `licensetype`, `linescomparetype`, `logtype`, `revokereason`.
- **Surface:** Template; many are pure label enums — coverage should jump fast.
- **Expected lift:** ~+3pp.

#### AL2-06. Batch F — task / script / prompt (5 pkgs)
- **Targets:** `taskcategory`, `taskpriority`, `scripttype`, `promptclitype`, `cmdenumtypes`.
- **Surface:** Template + scripttype any path helpers.
- **Expected lift:** ~+2.5pp.

#### AL2-07. Bespoke — `dbdrivertype` connection-string suite
- **Why:** Has `Connection.go`, `ConnectionOptions.go`, `connectionStringCompiler.go` — ~3 source files of business logic the template won't reach.
- **Approach:** Build representative `Connection` for each driver Variant; assert compiled DSN string matches expected pattern; round-trip options.
- **Expected lift:** +1–2pp; package 0% → 60%+.

#### AL2-08. Bespoke — `osdetect` Linux/Windows guarded branches
- **Status:** ✅ Done (2026-05-07, v0.99.0). Pure-`testing` uplift sweep added at `osdetect/OsDetect_Uplift_test.go` covers every Variant accessor/predicate/JSON method, OperatingSystemDetail populated/empty/nil paths (incl. `ReleaseVersion` cached / empty / nil branches and the `Is*AtLeast` family), WindowsSystemDetail client + server + nil-receiver paths, OsDetailWithErr round-trip, and host-smoke calls (no host-specific assertions). Cross-platform safe — runs identically on Linux/macOS/Windows.

**Combined AL2 target:** total coverage **~35–40% → ≥65%**, every per-package ≥60% except deliberately skipped platform-only files.

---

## Phase 4: Manual / Parked

### A. Manual `cross-repo/core-v9/` push

- **Status:** ⏭️ Manual user action (credential-bound)
- **Trigger:** When main-repo CI files change, mirror then user pushes.

---

## Next Task Selection

**Recommended next task:** Pick from this list (in order):

1. **AO** — Coverage probe + uplift any package still <75% (gate threshold) ⭐ **NEXT**
2. **AN-D** — Validate `0 issues` on next CI run (awaiting fresh push) ⏳
3. **AC** — Re-audit §07 / §09 spec dimensions
4. **AB residual** — Continue ❓ promotion in §03/§04 (PI-003, LOW)
5. **A** — Manual `cross-repo/core-v9/` push ⏭️

**Done from this list (reconciled 2026-05-10):**
- AL2-02 ✅ (Cycles 68–69, v0.40.0/v0.41.0 — DB family Coverage_test files in all 6 packages)
- AL2-03 ✅ (Coverage_test files in all 5 networking packages)
- AL2-04 ✅ (Coverage_test + Uplift_test in all 6 Linux/OS packages)
- AL2-05 ✅ (Coverage_test + Uplift_test in all 7 misc enum packages)
- AL2-06 ✅ (Coverage_test in all task/script/prompt packages; cmdenumtypes Coverage_test in all 27 sub-packages)
- AL2-07 ✅ (`dbdrivertype/DbDriverType_Connection_Coverage_test.go` shipped)
- AL2-08 ✅ (Cycle, v0.99.0 — `osdetect/OsDetect_Uplift_test.go`)
- AL-01..AL-08 ✅ (Cycles 49–58 — full AL umbrella)
- Cycle 57 ✅ (test fixes for 6 failures)
- PI-008 ✅ (Cycle 59, off-by-one fix in quotes/brackets unwrap helpers)
- PI-005 + PI-006 + PI-007 ✅ (Cycle 60, sqliteconnpathtype cluster — local overrides for upstream core-v9 defects; 4 skip-list entries removed)
- AI ✅ (Cycle 61, `spec/01-app/` formally FROZEN in spec-v0.53.0)
- AA / Cycle 15 ✅ (audited `spec/06-testing-guidelines/`, baselined at 100% verifiable)
- AB residual for §06 ✅ (Cycle 62, 10 deferred ❓ → 10 ✅ via `/tmp/core-v9-upstream` v1.5.8; zero new findings)
- Cycle 63 ✅ (3 test failures fixed: osdetect lowerCaseNames gap, sqliteconnpathtype StringMin fixture drift, sqliteconnpathtype RangesDynamicMap upstream lazy-init defect; 4th "failure" was Goconvey log-conflation phantom; RCA pattern catalogue saved to `.lovable/memory/07-test-failure-rca-patterns.md`)

### AO. Coverage probe + uplift remaining sub-75% packages (NEW)

- **Status:** 📋 Planned
- **Objective:** Run `./run.ps1 -tc` (or `go test ./... -coverprofile=coverage.out`), identify every package below the 75% CI gate, ship targeted uplift tests.
- **Dependencies:** None.
- **Approach:** (1) Run coverage. (2) Sort `go tool cover -func=coverage.out` ascending. (3) For each package <75%, write a focused `*_Uplift_test.go` covering the uncovered statements. (4) Re-run, confirm all ≥75%.
- **Acceptance:** Coverage gate passes cleanly with no per-package outlier <75%.

**Done — full AL umbrella:**
- AL-01 ✅ (Cycle 49, 15.5% → 21.6%, +6.1pp)
- AL-02 ✅ (Cycle 50, 21.6% → 26.1%, +4.5pp)
- AL-03 ✅ (Cycle 51, 26.1% → 33.8%, +7.7pp; cumulative +18.3pp)
- AL-04 ✅ (Cycle 52, +4–6pp expected; pending local `./run.ps1 -tc` confirmation)
- AL-05 pass-1 ✅ (Cycle 53, 4 packages: accesstype/certaction/completionstate/compressformats)
- AL-05 pass-2 ✅ (Cycle 54, 6 packages: dbaction/envtype/iptype/onofftype/overwritetype/timeunit)
- AL-06 ✅ (Cycle 55, quotes/+brackets/ bespoke wrap-unwrap suites)
- AL-07 ✅ (Cycle 56, strtype + inttype)
- AL-08 ✅ (Cycle 58, osdetect)

**Phase AL² — second-pass coverage push (post-AK) ✅ COMPLETE:**
- AL²-04 ✅ (NumericRange / Format / JSON shared-loop suites)
- AL²-05 ✅ (Format / OnlySupportedErr / MarshalJSON-dispatch sweep, v0.72.0)
- AL²-07 ✅ (concrete-Variant byte-predicate cluster: IsValueEqual / IsByteValueEqual / IsAnyValuesEqual via reflect, v0.73.0)
- AL²-06 ✅ (pointer-receiver binding wrappers via reflect, v0.74.0)
- AL²-08 ✅ (range-edge fuzz on MinByte/MaxByte/RangesByte via reflect, v0.75.0)

**Recipe-validation passes ✅:**
- pass-1 (AK) `httpmethodtype/` (v0.71.0)
- pass-2 `httpstatusfamily/` (v0.76.0)
- pass-3 `mimetype/` (v0.77.0) — recipe confirmed repeatable three times without refinement.

---

## Phase 5: CI Lint RCA (AN)

### AN. Recurring CI lint "errcheck/ineffassign in osdetect/generate.go" failures

- **Status:** 🔄 In Progress
- **Symptom:** CI reports 4 issues at `osdetect/generate.go:98,312,358,366` across ≥7 consecutive runs; identical byte-for-byte each time.
- **Sub-tasks:**

#### AN-A. Confirm CI source mismatch ✅ DONE (2026-05-10)
- Local `osdetect/generate.go` SHA-256 = `42027f133fb3e751d9f0139e34e5223ec5860e9fec9021762f6572873ed7550f`.
- Lines 98 / 312 carry the `_ =` errcheck ignores; lines 358 / 366 sit inside `if rel := …; rel != ""` shadow-blocks.
- Reproduced CI environment locally: **exact** `golangci-lint-2.5.0-linux-amd64` binary + Go 1.25.7 → `0 issues.`
- **Conclusion:** current source is clean. CI is linting an older SHA — most likely a **re-run of an old workflow run** (GitHub Actions re-runs replay the original commit, not HEAD).

#### AN-B. Force CI to expose the linted SHA + eliminate stale-cache vector ✅ DONE (2026-05-10, v1.22.0)
- Diagnostic step `Source fingerprint before lint` already in place (v1.21.0).
- **Refined RCA:** `golangci/golangci-lint-action@v7` caches analysis results keyed on `go.sum` + linter version. Since `go.sum` was unchanged across all 7 failing runs, the action was replaying stale cached findings — exactly matching the byte-identical 4-error output.
- **Fix shipped:** `skip-cache: true` + `skip-save-cache: true` on the action. Forces fresh analysis every run.

#### AN-C. Apply the verified fix ✅ DONE (folded into AN-B above)

#### AN-D. Validate ⏳ Pending
- Re-run CI and confirm `0 issues`.

---

## Completed

### Cycles 1–14 — `spec/01-app/` directory audit

- **Completed:** 2026-05-04 → 2026-05-05
- **Result:** All 14 numbered files audited. 12 at 100% verifiable. 2 baseline-only.
- **Findings:** C-CVS-01..10 (10 contradictions), D-CVS-01..43 (43 drifts), all resolved.

### W + AG — Unblock Build

- **Completed:** 2026-05-05
- **Result:** Upstream `core-v9` `go.mod` fixed, `replace` bridge removed, `require core-v9 v1.5.8` pinned.
