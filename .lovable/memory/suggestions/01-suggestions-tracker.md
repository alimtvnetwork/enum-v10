# Suggestions Tracker

> Single file tracking all Lovable suggestions. Update in-place — do not create per-suggestion files.
> When a suggestion is completed, move it to the **Completed** section with date and notes.

---

## Convention

- **suggestionId**: `S-<NNN>` (sequential)
- **status**: `open` | `in-progress` | `done` | `rejected`
- When done: move to Completed section, add completion date and notes
- When rejected: move to Rejected section, add reason

---

## Open Suggestions

### BA: Per-package coverage uplift — restore floor on remaining library packages

- **Status:** ✅ DONE (7/7 as of v1.40.0). All `--ignore` entries for library packages removed; only structural ignores remain (`cmd/main`, `tests/creationtests`, `scripts/autofix|bracecheck|specapisig`, `cross-repo/core-v9/scripts/*`).
- **Created:** 2026-05-10 (post-remix CI failure on per-package coverage gate).
- **Source:** Lovable / CI failure.
- **Done:** ✅ `httpstatusfamily` (v1.34.0). ✅ `compressformats` (v1.35.0). ✅ `compresslevels` (v1.36.0). ✅ `osarchs` (v1.37.0). ✅ `instructiontype` (v1.38.0). ✅ `inttype` (v1.39.0). ✅ `dbdrivertype` 58.8% → ≥95% (v1.40.0, new `DbDriverType_DirectMethods_test.go`).
- **Reusable template:** `<Pkg>_DirectMethods_test.go` covering non-nullary methods (`IsEqual/IsAnyOf/IsAnyValuesEqual/IsAnyNamesOf/Format/OnlySupported*Err`), top-level helpers (`Min/Max/RangesInvalidErr/New/NewMust/Is/ValidationError/StringMustBe`), and JSON round-trip (`MarshalJSON` + `UnmarshallEnumToValue` or `JsonParseSelfInject`). Typically pushes a 60-75% package to ≥95% in 60-100 LOC.
- **Acceptance:** ✅ All library `--ignore` entries removed from `.github/workflows/ci.yml`; CI passes per-package gate at threshold 75 with bare flag.

### S-115: Harden `Get-UpstreamPackages` to recursively walk `coredata/`


- **createdAt:** 2026-05-06
- **source:** Lovable (Cycle 47 — third retraction R-CVS-03 of same drift class)
- **affectedProject:** enum-v9 tooling
- **description:** `scripts/spec-api-check.psm1` `Get-UpstreamPackages` indexes `core-v9` package directories but missed nested `coredata/coredynamic` (R-CVS-01), `coredata/corestr` (R-CVS-02), and `coredata/coreonce` (R-CVS-03) — all three caused initial false fabrication-flagging that had to be retracted.
- **rationale:** Three same-class probe failures across three audit cycles indicates the indexer needs to recursively descend at least one level into `coredata/` (and possibly other parent dirs like `internal/`).
- **proposed change:** Add a recursive walk of `coredata/*` directories in the `-LocalDir` and `-UpstreamDir` indexing paths; re-run S-106 lint after the change to confirm no regression on existing flagged items.
- **acceptance criteria:** A test-clone scenario where `coredata/coreonce` exists is correctly indexed and `coreonce.NewAnyErrorOnce` is NOT flagged as a fabrication.
- **status:** ✅ done — moved to Completed (Cycle 48). Original framing was wrong (recursion already worked). See completed entry for the actual fix.

---

### S-002: Promote `errcore.VarTwoNoType` from ❓ to ✅ in Cycle 6

- **createdAt:** 2026-05-05
- **source:** Lovable (Cycle 13, §15 audit)
- **affectedProject:** enum-v9 spec
- **description:** `VarTwoNoType` was scored ❓ but IS cross-referenced from multiple spec files.
- **rationale:** Under spec-internal-consistency dimension it qualifies as ✅.
- **proposed change:** Backport promotion when Task AC runs.
- **acceptance criteria:** Cycle 6 audit report row 16 updated.
- **status:** open (deferred — folded into Task AC consistency-dimension re-audit; will land with the §15 sweep after the freeze lifts)

---

## Completed Suggestions

### S-115: Harden upstream-clone presence detection (audit-probe correctness)

- **completed:** 2026-05-06 (Cycle 48)
- **source:** Lovable (Cycle 47 — third retraction R-CVS-03 of same drift class)
- **resolution:** Empirically verified the original framing was wrong: `Get-UpstreamPackages` already walks `coredata/*` recursively via `Get-ChildItem -Recurse` (line 82). A live indexing run against `/tmp/core-v9-upstream` correctly indexed `coreonce`, `coregeneric`, `corestr`, `coredynamic` (177 pkgs total; `coreonce.NewAnyErrorOnce` resolved). The real defect was operator-side: when the upstream clone is **missing**, audit probes that read source via `rg`/`grep` directly silently return 0 hits and produce false fabrication conclusions (R-CVS-01/02/03 were all the same drift). Built two guardrails in `scripts/spec-api-check.psm1` v1.2.0: (1) **`Test-UpstreamClone` exported helper** that returns `{ Ok; Path; Reason; PackageCount }` and supports `-AutoClone` for one-shot remediation. Sentinel = `coredata/coregeneric` (always present in a real core-v9 checkout). Reasons: `missing` / `sentinel-missing` / `clone-failed` / `ok`. (2) **Sentinel-missing warning inside `Get-UpstreamPackages`** so every spec-api-check run also surfaces the drift even if the helper isn't called explicitly. Wired the helper into `scripts/CoveragePreChecks.psm1` "Spec-API Lint" phase: the old plain `Test-Path $upstreamDir` skip is replaced with `Test-UpstreamClone` so a directory that exists but lacks the sentinel (wrong branch, partial clone) now skips with a precise reason instead of running and producing false negatives. Smoke test `tests/scripts/Test-UpstreamClone.ps1` covers 4 cases (missing path / sentinel-missing / `Get-UpstreamPackages` warning / real clone) — all 7 assertions pass via `nix run nixpkgs#powershell`. `package.json` 0.17.0 → 0.18.0.
- **acceptance criteria:** ✅ `Test-UpstreamClone` exported and returns the documented shape. ✅ Real clone scenario indexes `coreonce.NewAnyErrorOnce` (verified pre-change — already worked). ✅ Sentinel-missing scenario detected. ✅ CI pre-check uses sentinel-aware skip with precise reason. ✅ `Get-UpstreamPackages` emits a warning when sentinel is absent. ✅ Smoke test in repo prevents regression.

### S-114: `PackageCoverage.psm1` (and 3 other modules) hard-codes `tests/integratedtests/` path

- **completed:** 2026-05-06 (Cycle 40)
- **source:** Lovable (discovered while scoping AL — Cycle 39 follow-up)
- **resolution:** Added shared `Resolve-TestSuiteRoot` helper to `scripts/Utilities.psm1` (~30 LOC). Probes `creationtests` first, then `integratedtests`, with optional `-Package` parameter for per-package resolution; falls back to `creationtests` default so the downstream `go test` error becomes the user-facing diagnostic. Refactored four callsites to use it: `scripts/PackageCoverage.psm1` (TCP — lines 45 + 62), `scripts/TestRunner.psm1` (TP — lines 35/42/44), `scripts/Help.psm1` (`Invoke-IntegratedTests` — lines 27/31), `scripts/PreCommitCheck.psm1:55`. `CoverageRunner.psm1` (lines 103-107) was already correct and left untouched. Smoke test at `tests/scripts/Test-ResolveTestSuiteRoot.ps1` covers 5 cases (creationtests-only / `-Package` known / `-Package` missing / legacy-only / live repo) — all pass. Module-import smoke confirms helper exported + all four refactored modules parse cleanly. `package.json` 0.9.0 → 0.10.0.
- **acceptance criteria:** ✅ TCP/TP/integrated-tests/pre-commit commands now find the real `tests/creationtests/` packages on this repo. ✅ Helper accepts either layout per Core-memory rule (never hard-coded). ✅ Smoke test in repo prevents regression. ✅ `CoverageSplitRecovery.psm1` / `CoverageProfileMerger.psm1` / `CoverageReportHtml.psm1` / `CoveragePreChecks.psm1` references to `integratedtests` left as comments/doc-strings only — no behaviour changes needed (they regex-match against either name or are display-only).

### S-111: Surface the GoConvey-only sub-pattern in spec/06

- **completed:** 2026-05-06 (Cycle 39)
- **source:** Lovable (Cycle 37 — D-CVS-64 carry-forward)
- **resolution:** Added a "Sub-Pattern: GoConvey-Only (Local Wrapper)" section to `spec/06-testing-guidelines/02-test-case-types.md` between the Style D footnote and the `CaseV1` heading — describes scope/when-to-use, gives a worked example from `enum-v9/tests/creationtests/AllEnums_ContractsTesting_test.go` (`Convey` + `EnumTestWrapper` registry + `LogShouldDiffMessage` + `So(diff, ShouldBeEmpty)`), and provides an equivalence table mapping upstream primitives (`CaseV1`, `coretests.GetAssert.ShouldBeEqualMap`, `args.Map`, `t.Run`, `tc.ShouldBeEqualFirst`) to GoConvey-only counterparts (local wrapper struct, diff-based assertion, module-level registry, `Convey` nested scopes, `So(actual, ShouldEqual, expected)`). Added a companion "Sub-Pattern: GoConvey-Only Diff Assertion" section to `spec/06-testing-guidelines/05-assertion-patterns.md` after the named-map-types pitfall — explains why diff-based returns empty-string-on-success / human-readable-diff-on-failure (matching `ShouldBeEqualMap` ergonomics with no `coretests` dependency), lists five companion assertions (`ShouldEqual`, `ShouldResemble`, `ShouldBeNil`, `ShouldBeTrue`/`ShouldBeFalse`, `ShouldBeEmpty`), cross-links the worked example. Cosmetic spec-only change. Treated as editorial under the `spec/01-app/` freeze (the freeze covers `spec/01-app/`; this edit lives in `spec/06-testing-guidelines/`). Spec changelog → spec-v0.45.0. `package.json` 0.8.0 → 0.9.0.
- **acceptance criteria:** ✅ Both files include the GoConvey-only callout. ✅ Cross-ref to `tests/creationtests/AllEnums_ContractsTesting_test.go` resolves. ✅ Equivalence table covers all 5 framework primitive ↔ sub-pattern mappings.

### S-112 + S-113: Truthful Git-Pull phase + remote-probe skip

- **completed:** 2026-05-06 (Cycle 38)
- **source:** Lovable (user-reported `./run.ps1 -tc` run showing `✗ git pull failed` immediately followed by `✓ Git Pull pulled from remote` in the Phase Summary box)
- **resolution:** Two related defects fixed in one pass against `scripts/TestRunnerCore.psm1` + `scripts/CoverageRunner.psm1`. **S-113 (root cause):** `Invoke-GitPull` ran `git pull` unconditionally and emitted `remote: Repository not found / fatal: repository '…/enum-v9.git/' not found` whenever the local clone had a misconfigured/private/missing `origin`. Added a two-step early probe BEFORE `git pull`: (1) `git remote get-url origin` → if missing return `Status='skip'/Message='no origin remote'`; (2) `git ls-remote --exit-code origin HEAD` → if unreachable return `Status='skip'/Message='remote unreachable'` and surface the URL. **S-112 (truthfulness):** `Invoke-GitPull` was void-return so `CoverageRunner.psm1:27-30` hard-coded `Register-Phase "Git Pull" "pass" "pulled from remote"` regardless of outcome. Refactored `Invoke-GitPull` to return `[pscustomobject]@{ Status='pass'|'warn'|'skip'; Message=... }`; `Invoke-FetchLatest` returns `@{ GitPull; Tidy }` with the same shape for `go mod tidy`; `CoverageRunner.psm1` uses both results. The dashboard already supports `skip → ⊘` and `warn → ⚠` glyphs in `scripts/DashboardPhases.psm1:45-51` so no UI changes were needed. Added regression test `tests/scripts/Test-InvokeGitPull.ps1` — both cases pass via `pwsh -File tests/scripts/Test-InvokeGitPull.ps1` (Test 1 no-origin → skip / no origin remote; Test 2 bogus unreachable → skip / remote unreachable). `package.json` 0.7.0 → 0.8.0.
- **acceptance criteria:** ✅ Phase Summary now renders `⊘ Git Pull remote unreachable` (or `⊘ Git Pull no origin remote`) instead of the lying `✓ Git Pull pulled from remote` after a soft-fail. ✅ No more confusing `remote: Repository not found` for users with missing/private `origin`. ✅ Smoke tests cover both skip paths.

### S-109: Cycle-15 deep-probe of `tests/creationtests/` patterns to clear 21 ❓

- **completed:** 2026-05-06 (Cycle 37)
- **source:** Lovable (Cycle 27 — carry-forward)
- **resolution:** Read all 14 files under `enum-v9/tests/creationtests/` and ran the symbol-set probe `rg -n 'coretests\.|args\.|results\.|CaseV1|CaseNilSafe|GenericGherkins|GetAssert|ShouldBeEqualMap|ShouldBeSafe|InvokeWithPanicRecovery|BaseTestCase' tests/creationtests/` → **zero hits**. Confirmed `enum-v9` deliberately does NOT consume the upstream `coretests`/`args`/`results` framework; instead uses GoConvey + 2 local wrapper structs + module-level slice/map registries + AAA comments. Settled the 10 Cycle-15 ❓ items: **1 promoted ✅** (claim 20, diff-based assertion pattern via `enumimpl.DynamicMap.LogShouldDiffMessage` + `So(diff, ShouldBeEmpty)` in `AllEnums_ContractsTesting_test.go`), **9 annotated ⓘ "upstream-only"** (no `enum-v9` evidence available; remain blocked by Task AB for upstream-clone promotion). Cycle-15 verifiable subset 22/22 → 23/23 (still 100%); spec/06 unknown ❓ pool **10 → 0**. New finding D-CVS-64 (LOW) carried forward as **S-111**. Audit file: `spec/07-code-vs-spec-audits/29-cycle37-S109-creationtests-deep-probe.md`. Spec changelog → spec-v0.44.0. `package.json` 0.6.0 → 0.7.0.
- **acceptance criteria:** ✅ New audit cycle entry on scoreboard (`Cycle 37`). ✅ Cycle-15 ❓ pool reduced from 10 → 0 unknown (1 promoted, 9 annotated upstream-only). ✅ Direct source evidence cited for every promotion/annotation.

### S-004: Fix `spec/00-llm-integration-guide.md` stale test path reference

- **completed:** 2026-05-06 (Cycle 34)
- **source:** Lovable (reliability report, 2026-05-05)
- **resolution:** Re-framing applied (not path replacement). The `tests/integratedtests/` references on lines 36 (Decision Matrix Style D row) and 826 (Test Folder Structure section) correctly document the **upstream `core-v9` consumer layout**; `enum-v9` itself uses the single shared `tests/creationtests/` package. Added an inline scope callout above the line-826 code fence that explains the upstream-vs-enum-v9 split and cross-links both `spec/01-app/13-testing-patterns.md` §6.1 and `spec/06-testing-guidelines/01-folder-structure.md` (which already carries the same disclaimer on its own line 3). Pattern matches the resolution model used in Cycles 12/15/17/18 for sibling `integratedtests` references.
- **acceptance criteria:** ✅ `tests/integratedtests/` references in `spec/00-llm-integration-guide.md` are now scope-disambiguated. ✅ No misdirection risk for AI readers — the callout lands before the code block.

### S-003: Fix `spec/06-testing-guidelines/01-folder-structure.md` stale path

- **completed:** 2026-05-06 (Cycle 34) — already-resolved discovery
- **source:** Lovable (reliability report, 2026-05-05)
- **resolution:** Verified obsolete. `spec/06-testing-guidelines/01-folder-structure.md` line 3 already carries the upstream-scope disclaimer added in earlier audit cycles: *"⚠️ Scope: the layout below describes upstream `core-v9`. enum-v9 uses a single shared `tests/creationtests/` package — see [spec/01-app/13-testing-patterns.md §6.1]…"*. The `integratedtests` references in the file body intentionally document the upstream consumer layout, not enum-v9. Closing as **resolved-by-prior-cycle** with no edit needed.
- **acceptance criteria:** ✅ Scope disclaimer present (line 3). ✅ No misdirection risk.

### S-001: Pin Go toolchain to 1.22 as stopgap for Task W

- **completed:** 2026-05-06 (Cycle 34) — closed as **obsolete**
- **source:** Lovable (Cycle 13, 2026-05-05)
- **resolution:** Stopgap no longer relevant. The underlying problem (Go 1.25 rejecting the dual-path `replace` bridge to upstream `core-v9`) was fixed at the source by Tasks W + AG (2026-05-05): upstream `core-v9` was renamed to `module github.com/alimtvnetwork/core-v9` and tagged `v1.5.8`; `enum-v9/go.mod` now pins `core-v9 v1.5.8` directly with no `replace` bridge (per Core memory). Pinning to Go 1.22 would now mask the working modern setup. **Do not** add `toolchain go1.22.0` to `go.mod`.
- **acceptance criteria:** ✅ `replace` bridge removed (Task AG). ✅ `enum-v9` builds against current toolchain without the stopgap. ✅ Risk noted in original suggestion (toolchain lock-in) avoided.

### S-107: Goconvey-failure summarizer for `failing-tests.txt`

- **completed:** 2026-05-06 (Cycle 33)
- **source:** Lovable (carry-forward from Cycle 24 — `TestLogWriter.psm1` Pass-2 collector buries goconvey failures under `N total assertions` lines)
- **resolution:** Added `Get-GoconveyFailureSummary` to `scripts/TestLogWriter.psm1` (also exported). For each failed test's captured block, walks the lines and pairs every `Expected:` with its nearest following `Actual:` / `(Line N)` / `Message:` (8-line look-ahead window, stops when the next `Expected:` starts). `Write-TestLogs` now prepends a `── Failure summary (N) ──` mini-section to each `FAIL:` block in `failing-tests.txt`, listing each triplet as `#i Expected: ... | Actual: ... (Line N) [Message]` BEFORE the noisy raw block. The original raw block is preserved underneath for full context. Smoke-tested via `nix run nixpkgs#powershell` on a synthetic 2-failure goconvey block: extracted both triplets correctly (incl. optional Message field), empty blocks return 0, non-goconvey blocks (e.g. plain `panic:`) return 0 with no false positives. `package.json` 0.2.0 → 0.3.0.
- **acceptance criteria:** ✅ Failed-test blocks now lead with the goconvey Expected/Actual triplets. ✅ Existing raw-block content is preserved (no information loss). ✅ Helper handles empty + non-goconvey blocks safely.

### S-110: Restore standalone coverage utilities documented alongside S-108

- **completed:** 2026-05-06 (Cycle 32)
- **source:** Lovable (Cycle 31 follow-up to S-108 completion)
- **resolution:** Created the three standalone utilities listed in `spec/03-powershell-test-run/06-coverage-prompt-generator.md`:
  - `scripts/coverage/Get-UncoveredLines.ps1` — emits a single-row object `{SourceFile, UncoveredCount, Ranges}` with collapsed line ranges (e.g. `L15-L17, L22`).
  - `scripts/coverage/Get-FunctionCoverage.ps1` — parses `go tool cover -func` output, filters below `-Threshold` (default 100.0), returns objects sorted ascending by coverage.
  - `scripts/coverage/Get-PackageCoverageReport.ps1` — combines both for a single `-Package`; supports `-Format text|markdown|json` and optional `-OutputFile`.
  All three reuse the same regexes and range-collapse helper as the S-108 main script and avoid the `$Input` automatic-variable shadowing pitfall (use `$Source`). Smoke-tested via `nix run nixpkgs#powershell` against synthetic 6-block `coverage.out` + 3-line func output: each utility produces exactly the spec-shaped output (correct sort order, correct ranges, all 3 formats). `package.json` 0.1.0 → 0.2.0.
- **acceptance criteria:** ✅ All four scripts now exist in `scripts/coverage/`. ✅ Spec drift D-CVS-62 is **fully** closed (was partially closed by S-108).

### S-108: Restore missing `scripts/coverage/Generate-CoveragePrompts.ps1`

- **completed:** 2026-05-06 (Cycle 31)
- **source:** Lovable (Cycle 27 — AB scripts deep-probe, surfaced D-CVS-62)
- **resolution:** Created `scripts/coverage/Generate-CoveragePrompts.ps1` per the contract in `spec/03-powershell-test-run/06-coverage-prompt-generator.md`. Parses `coverage.out` (statement blocks with `count==0`) + `go tool cover -func` output, sorts functions ascending by coverage, batches at `-BatchSize` (default 500) into `data/prompts/coverage-prompt-N.txt`, emits `prompts-summary.json`. Smoke-tested end-to-end via `nix run nixpkgs#powershell` against synthetic `coverage.out` + func output: produced exactly the spec sample format including range-collapsed uncovered-line ranges (`L15-L17, L22`) and ascending coverage sort. Both call-sites in `scripts/CoverageRunner.psm1:313-316` and `scripts/PackageCoverage.psm1:150` now resolve. Standalone utilities (`Get-UncoveredLines.ps1` etc.) tracked separately as **S-110**.
- **acceptance criteria:** ✅ `Test-Path scripts/coverage/Generate-CoveragePrompts.ps1` → True. ✅ Smoke run produces `coverage-prompt-1.txt` matching spec format and `prompts-summary.json`.

### S-100: Add `cmd/main/` smoke-test policy carve-out to spec §12

- **completedAt:** 2026-05-05 (Cycle 10, fix C-CVS-10)
- **notes:** Spec §12 rewritten as "library-first, smoke-test allowed" policy.

### S-101: Rewrite §06 around `SimpleSlice`/`Hashset`/`SimpleStringOnce`

- **completedAt:** 2026-05-04 (Cycle 4, fixes D-CVS-22/23/24)
- **notes:** Replaced fictional `coreonce.New.String` with actual constructors.

### S-102: Add consumer-coverage callouts for upstream-only API

- **completedAt:** 2026-05-04 → 2026-05-05 (D-CVS-25, D-CVS-38, D-CVS-42)
- **notes:** Three sections now carry explicit upstream-only callouts.

### S-103: Extract portable runner specs into `spec/03-portable/` sub-directory

- **completed:** 2026-05-06 (Cycle 36)
- **source:** Lovable (Cycle 16 carry-forward)
- **resolution:** Moved `spec/03-powershell-test-run/08-generic-go-test-coverage-runner.md` → `spec/03-powershell-test-run/portable/01-generic-go-test-coverage-runner.md` and `09-ai-agent-complete-reference.md` → `portable/02-ai-agent-complete-reference.md`. Added `spec/03-powershell-test-run/portable/README.md` with the scope-split explanation, the 2-file table, and 3 editor rules. Updated 2 live cross-refs (`spec/00-llm-integration-guide.md:2380`, `spec/04-tooling/03-powershell-implementation.md:456`) and 1 internal cross-ref inside the moved `02-ai-agent-complete-reference.md`. Historical references in `spec/CHANGELOG.md`, `spec/07-code-vs-spec-audits/17-cycle16-powershell-test-run.md`, and `spec/99-audits/01-original-11-step-plan.md` left as-is (audit history). Spec changelog → spec-v0.43.0. `package.json` 0.5.0 → 0.6.0.
- **acceptance criteria:** ✅ `rg -n 'spec/03-powershell-test-run/(08|09)-' spec/ --glob '!spec/CHANGELOG.md' --glob '!spec/07-code-vs-spec-audits/**' --glob '!spec/99-audits/**'` returns zero hits. ✅ `spec/03-powershell-test-run/portable/README.md` exists and documents the split.

### S-104: Add `cross-repo/core-v9/README.md` historical-naming top-of-file note

- **completed:** 2026-05-06 (Cycle 35)
- **source:** Lovable (Cycle 17 carry-forward)
- **resolution:** Added a prominent top-of-file callout to `cross-repo/core-v9/README.md` covering: (1) the directory name is historical and intentional (mirrors a separate upstream repo), (2) the actual import path used by `enum-v9` source is `github.com/alimtvnetwork/core-v9` (renamed 2026-05-05, tagged `v1.5.8`), (3) editors must always write `cross-repo/core-v9/` even when the surrounding sentence is about `core-v9` content, (4) the historical `enum-v1` / `core-v8` body references must NOT be rewritten (Core-memory rule). Body untouched. Closes the Cycle-17 root cause of D-CVS-49/52/53/55 at the point of truth.
- **acceptance criteria:** ✅ Head section includes the 4-point explanation. ✅ Body content preserved.

### S-105: CI guard — `spec/02-app-issues/` index-drift detector

- **completed:** 2026-05-06 (Cycle 35)
- **source:** Lovable (Cycle 18 carry-forward)
- **resolution:** Added `scripts/ci/check-issues-index-drift.py` + `scripts/ci/test_check_issues_index_drift.py` (5 unittest cases, all pass) and a new `issues-index-drift` job in `.github/workflows/ci-guards.yml` (depends on `python-tests`, mirrors the `collision-audit` job). Script extracts every `| NN |`-prefixed row from `00-issues-index.md` (canonical) and `README.md` (human), compares **row count AND id-set**, exits 1 with `Missing from README` / `Missing from index` diffs on drift, exits 2 on missing files. Live repo reports `OK: spec/02-app-issues index in sync (9 rows).` — adopted at clean baseline.
- **acceptance criteria:** ✅ Removing a row from README → CI fails (`test_row_count_mismatch_fails`). ✅ Both in sync → CI passes (`test_in_sync_passes` + live-repo smoke). ✅ Same count, mismatched ids → CI fails (`test_id_set_mismatch_same_count_fails`).

### S-106: `spec-api-check.psm1` — automate code-fence vs API verification

- **createdAt:** 2026-05-06
- **source:** Lovable (Cycle 19 carry-forward)
- **affectedProject:** enum-v9
- **description:** Cycle 19 (AB pass 1) found 5 ❌ contradictions in `spec/01-app/09-converters.md` whose root cause is identical: **the spec was authored against an internal mental model of the API, not against grep output**. C-CVS-15 in particular fabricated an entire `typesconv` numeric-widening surface that does not exist.
- **rationale:** The same author pattern is almost certainly present in the other §07/§08/§10/§11/§15/§16 ❓ pools (124 ❓ remaining). A one-shot lint can prevent the next 50+ ❌s from ever shipping.
- **proposed change:** Add `scripts/spec-api-check.psm1` that, for every Go code-fenced block in `spec/01-app/`, extracts `<pkg>.<symbol>` references and runs `go vet` (or `go list -f '{{.Imports}}'`) against a synthesised stub `.go` file in a tmp module rooted at `core-v9 v1.5.8`. Failures should print the spec file + line.
- **acceptance criteria:** Removing `Integer` from `converters.stringTo` makes the lint fail; current spec text after AJ-01..03 land makes it pass; CI-friendly.
- **status:** open
- **risk:** Tooling addition only — no change to spec or production code.

---

## Rejected Suggestions

_(none yet)_
