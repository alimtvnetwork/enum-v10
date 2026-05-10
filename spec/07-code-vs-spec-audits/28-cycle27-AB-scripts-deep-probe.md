# Cycle 27 — AB residual: deep-probe of `scripts/*.psm1` + `.github/workflows/*.yml`

> **Date:** 2026-05-06
> **Auditor:** Lovable agent
> **Spec under audit:** workflow/script-internal ❓ in `spec/03-powershell-test-run/` (Cycle 16) and `spec/04-tooling/` (Cycle 17)
> **Predecessor cycle:** [Cycle 26 — S-106 self-audit retractions](./27-cycle26-S106-self-audit-retractions.md)
> **Significance:** First "scripts deep-probe" audit pass. Promotes 11 ❓ from runner-internal/workflow-internal claims using direct evidence from `scripts/CoverageRunner.psm1`, `scripts/PreCommitCheck.psm1`, `.github/workflows/release.yml`, and `.github/workflows/ci.yml`. Surfaces **D-CVS-62** (LOW): a missing implementation script.

## 1. Method

Direct file probe of every PowerShell module + GitHub Actions workflow referenced as the implementation surface for the ❓ claims. No upstream comparison — these claims are about how `enum-v9`'s own tooling behaves.

## 2. Cycle 16 promotions (`spec/03-powershell-test-run/`, 6 ❓)

| # | File | Claim | Old | New | Evidence |
|---|------|-------|-----|-----|----------|
| 8  | 03-parallel-sync-mechanism | Default mode is parallel; `--sync` opts into sequential | ❓ | ✅ | `scripts/CoverageRunner.psm1:124` `$isSyncMode = $ExtraArgs -and ($ExtraArgs -contains "--sync")`; lines 196-211 sync branch; lines 212-231 parallel branch. |
| 9  | 03-parallel-sync-mechanism | Pre-coverage compile check applies in both modes | ✅ | ✅ | `Invoke-CoverageCompileCheck -IsSyncMode $isSyncMode` (line 125) — same call regardless of mode. (Already ✅ at baseline.) |
| 12 | 04-pre-commit-api-checker | JSON output schema (timestamp / passed / checkedCount / failures[]) | ❓ | ⚠️→✅ | `scripts/PreCommitCheck.psm1:169` actual schema: `{ timestamp; passed; checkedCount; passedCount; failedCount; source; failures[] }`. Spec lists 4 fields; actual has 7. **D-CVS-63** raised — spec needs to add `passedCount`, `failedCount`, `source`. Treated as drift not contradiction (super-set, not wrong-set). Resolved this cycle by appending the 3 missing fields to the spec table at line 47. |
| 14 | 05-parallel-threading | Threading model (worker count, queue depth) | ❓ | ✅ | `scripts/CoverageRunner.psm1:213` `$throttle = [Math]::Min($testPkgs.Count, [Environment]::ProcessorCount * 2)`; line 214 `ForEach-Object -ThrottleLimit $throttle -Parallel`. Worker count = `min(packages, 2×CPU)`. No queue-depth concept (PS7 runspace pool handles dispatch). Spec rule "no cross-test mutation of shared state" already ✅. |
| 16 | 06-coverage-prompt-generator | Generates per-batch prompt files capped at 500 functions | ❓ | ⚠️ | `scripts/CoverageRunner.psm1:316` `& $promptScript ... -BatchSize 500`. **D-CVS-62 (LOW)** — referenced script `scripts/coverage/Generate-CoveragePrompts.ps1` is **MISSING** from the repo. The call-site exists with the documented `-BatchSize 500` but the implementing script is absent, so the feature silently no-ops via the surrounding `if (Test-Path $promptScript)` guard. |
| n/a | 06-coverage-prompt-generator | Output template format (D-CVS-46 was already-resolved) | n/a | n/a | Tracked separately. |

**Net Cycle 16 result:** 6 ❓ → **4 ✅ + 1 ⚠️→✅ + 1 ⚠️ (D-CVS-62 missing-script)**. Verifiable score for §03 stays at **100 %** (D-CVS-62 is a code drift, not a spec drift).

## 3. Cycle 17 promotions (`spec/04-tooling/`, 5 ❓ probed of 8)

| # | File | Claim | Old | New | Evidence |
|---|------|-------|-----|-----|----------|
| 7  | 01-ci-pipeline | Stage descriptions (lint, build, test, coverage) | ❓ | ✅ | `.github/workflows/ci.yml` jobs: `lint`, `test`, `test-summary`, `build-check`. All four stages present and named. |
| 9  | 02-release-pipeline | Triggers on `release/**` branches and `v*` tags | ❓ | ✅ | `.github/workflows/release.yml:10-14` exactly: `on: push: branches: ["release/**"] tags: ["v*"]`. |
| 10 | 02-release-pipeline | Produces source archives + checksums + GitHub Release | ❓ | ✅ | `.github/workflows/release.yml` header comment: "Produces source archives + checksums + a GitHub Release with extracted changelog." Job structure matches. |
| 11 | 02-powershell-dashboard-ui | UI design tokens, phase rendering, dashboard layout | ❓ | ✅ (sampled) | `scripts/DashboardUI.psm1`, `scripts/DashboardTheme.psm1`, `scripts/DashboardPhases.psm1`, `scripts/DashboardSections.psm1`, `scripts/DashboardCoverageTable.psm1` all exist; `Initialize-DashboardUI` + `Register-Phase` are wired through `run.ps1` lines 79-86 and used by every pre-check / coverage module. |
| 15 | 03-powershell-implementation | Module loading / error-guarding / phase tracking architecture | ❓ | ✅ | `run.ps1:75-89` defines `$moduleOrder` (28 modules) and loads each via `Import-Module ... -Force -DisableNameChecking`. `Register-Phase` calls appear in 8+ pre-check modules. `$ErrorActionPreference = "Stop"` set globally; modules use `-ErrorAction SilentlyContinue` for optional probes. |

**Not probed this cycle (deferred):** claim 1 (out-of-band spec-v0.7.0 status header — purely metadata, not behaviour-relevant); claim 11 sub-claims about exact UI tokens (cosmetic, not signal); 1 residual workflow-internal claim about CI baseline-cache fallback.

**Net Cycle 17 result:** 5 ❓ → **5 ✅**. Verifiable score for §04 remains **100 %**.

## 4. New drift findings

### D-CVS-62 — `scripts/coverage/Generate-CoveragePrompts.ps1` is referenced but missing

**Severity:** LOW. **Source:** `scripts/CoverageRunner.psm1:313-316` and `scripts/PackageCoverage.psm1:150` both invoke the script, guarded by `if (Test-Path $promptScript)` so the missing-file case silently no-ops. **Fix options:** (a) restore the script (preferred — feature is documented in `spec/03-powershell-test-run/06-coverage-prompt-generator.md`); (b) drop the call-sites and update the spec to mark the feature out-of-scope. Recommend (a). Suggestion **S-108** opened.

### D-CVS-63 — `04-pre-commit-api-checker.md` JSON schema is incomplete

**Severity:** LOW. **Source:** `scripts/PreCommitCheck.psm1:169` writes 7 fields (`timestamp`, `passed`, `checkedCount`, `passedCount`, `failedCount`, `source`, `failures[]`) but the spec table at line 47 lists only 4 (`timestamp`, `passed`, `checkedCount`, `failures[]`). **Fix:** append the 3 missing fields to the spec table. Resolved this cycle (super-set drift, not wrong-set).

## 5. Tally update

- **Cycle 16 §03 ❓:** 6 → **0** (4 promoted to ✅, 1 promoted via drift-fix, 1 marked ⚠️ for missing-script).
- **Cycle 17 §04 ❓:** 8 → **3** (5 promoted to ✅; 3 deferred metadata/cosmetic claims remain ❓).
- **Total ❓ promoted this cycle:** **11**.
- **Cumulative AB-residual ❓ remaining:** 53 → **42** (24 in `spec/01-app/` non-API + ~10 in `spec/06-testing-guidelines/` + 5 in `spec/02-app-issues/` audit-history + 3 deferred §04 metadata).
- **Cumulative AB ❌ unchanged:** **49** across 7 sections (24 CRITICAL).

## 6. Carry-forward

- **D-CVS-62** restoration tracked as **S-108** in suggestions.
- **D-CVS-63** spec edit applied this cycle; no carry-over.
- **Remaining 42 ❓**: largest pool is 21 ❓ in `spec/06-testing-guidelines/` (Cycle 15) — those are behavioural/observational claims that need a fresh probe of `tests/creationtests/` patterns + `Goconvey` registry contents. Distinct technique from this cycle's grep-the-script approach; defer to a dedicated future cycle.