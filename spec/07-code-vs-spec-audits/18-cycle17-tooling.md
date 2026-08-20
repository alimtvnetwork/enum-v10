# Cycle 17 — `spec/04-tooling/` directory audit

> **Date:** 2026-05-06
> **Auditor:** Lovable agent
> **Spec under audit:** [`spec/04-tooling/`](../04-tooling/) (10 files, 2 553 lines)
> **Predecessor cycle:** [Cycle 16](./17-cycle16-powershell-test-run.md)
> **Significance:** Folds in residual task **AH** debt for `04-bootstrap-into-new-repo.md`. Surfaces 6 new drift findings beyond the AH-tracked one.

## 1. Method

Dual-dimension probe (same as Cycles 13–16):

1. **Code-vs-spec** — confirm referenced workflow files exist (`.github/workflows/{ci,release,vulncheck,ci-guards,python-tests}.yml`), `cross-repo/<dir>/` actually exists with the documented layout, and the staged `cross-repo/` README path is reachable.
2. **Spec-internal-consistency** — cross-refs resolve, no banned tokens (`enum-v10`, `enum-v10`, `enum-v10`, mojibake `core-v9 → core-v9`, `.lovable/user-preferences`), no contradiction with Core memory rule "the `cross-repo/core-v9/` directory intentionally keeps its `core-v8` name."

```bash
rg -nc 'integratedtests|enum-v10|enum-v10|enum-v10|core-v9 → core-v9|\.lovable/user-preferences|cross-repo/core-v9' spec/04-tooling/*.md
ls cross-repo/ .github/workflows/ scripts/ci/ 2>&1
```

**Result of the consumer probe:**
- All 5 referenced workflows exist: `.github/workflows/{ci,ci-guards,python-tests,release,vulncheck}.yml`.
- `cross-repo/` contains exactly **one** sub-directory: `cross-repo/core-v9/` — confirming the Core memory rule (the directory keeps its historical `core-v8` name even though the import path is now `core-v9`).
- `scripts/ci/` exists (referenced by `06-cross-repo-sync.md` §2 "Out of scope").
- `tests/creationtests/` exists; `tests/integratedtests/` does not exist in `enum-v10`.

## 2. Claim-by-claim table

> The 10 files together make ~85 normative claims. Below is a representative subset (30 claims) covering each file (00-overview, 01-ci-pipeline, 02a-release-pipeline, 02b-powershell-dashboard-ui, 03a-vulnerability-scanning, 03b-powershell-implementation, 04a-bootstrap, 04b-ci-guards, 05-branch-protection, 06-cross-repo-sync).

| # | File | Claim | Verdict | Evidence |
|---|------|-------|---------|----------|
| 1  | 00-overview | spec-v0.7.0 (2026-05-04) status header | ❓ | Out-of-band metadata. |
| 2  | 00-overview | Map table — 9 sub-spec rows with companion-code paths | ⚠️→✅ | **D-CVS-49** — row 06 cited `cross-repo/core-v9/` (broken — actual dir is `cross-repo/core-v9/`). Fixed inline this cycle with explicit Core-memory note. |
| 3  | 00-overview | Duplicate `02-` / `03-` / `04-` prefixes are intentional (paired specs) | ✅ | Note at line 28 explains; renumbering would break cross-refs. |
| 4  | 00-overview | "Maintenance" §3 references `cross-repo/core-v9/` | ⚠️→✅ | **D-CVS-49** part 2 — same cross-repo path drift fixed inline. |
| 5  | 00-overview | All 5 reading-path cross-refs resolve | ✅ | Verified — every `[NN — Title](./NN-...)` link targets a present file. |
| 6  | 01-ci-pipeline | `.github/workflows/ci.yml` referenced as primary CI workflow | ✅ | `ls .github/workflows/ci.yml` → present. |
| 7  | 01-ci-pipeline | Stage descriptions (lint, build, test, coverage) | ❓ | Behavioural; not probed against actual workflow YAML this cycle. |
| 8  | 02-release-pipeline | "Library release pipeline for `enum-v10`" | ✅ | Post-rename confirmed (`enum-v10` token correct). |
| 9  | 02-release-pipeline | Triggers on `release/**` branches and `v*` tags | ❓ | Not probed against `release.yml` this cycle. |
| 10 | 02-release-pipeline | Produces source archives + checksums + GitHub Release | ❓ | Same. |
| 11 | 02-powershell-dashboard-ui | (1244 lines) UI design tokens, phase rendering, dashboard layout | ❓ (sampled 3 sub-claims) | Behavioural / cosmetic; not a code-vs-spec audit target this cycle. Spec-internal: ✅ no banned tokens. |
| 12 | 03-vulnerability-scanning | Two-tier scanning per upstream `coding-guidelines-v20/spec/12 §03` | ✅ | Cross-ref resolves to upstream URL (well-formed). |
| 13 | 03-vulnerability-scanning | `.github/workflows/vulncheck.yml` is the implementation | ✅ | `ls` → present. |
| 14 | 03-powershell-implementation | (466 lines) "For AI agents working on the PowerShell tooling itself" | ✅ | Self-described scope. Banned-token-clean. |
| 15 | 03-powershell-implementation | Module loading / error-guarding / phase tracking architecture | ❓ | Behavioural; not probed against `scripts/*.psm1` this cycle. |
| 16 | 04-bootstrap-into-new-repo | "Generic — no `core-v9` assumptions baked in" | ✅ | §7 decoupling table makes this verifiable claim-by-claim. |
| 17 | 04-bootstrap-into-new-repo | §7 row "`tests/integratedtests/` mirror layout required ❌ No" | ⚠️→✅ | **D-CVS-50** — semantically correct but didn't name the `enum-v10` precedent (`tests/creationtests/`). Fixed inline this cycle to name both upstream-`core-v9` and `enum-v10` layouts. (This is the AH-tracked occurrence.) |
| 18 | 04-bootstrap-into-new-repo | §7 row "Module path hard-coded ❌ No — read from `go.mod`" | ✅ | Consistent with `01-app/03-import-conventions.md` §1. |
| 19 | 04-bootstrap-into-new-repo | §7 row "`coretests` framework required ❌ No" | ✅ | Consistent with `06-testing-guidelines/README.md` Cycle 15 callout (the framework is upstream-only). |
| 20 | 04-ci-guards | References `.github/workflows/ci-guards.yml` | ✅ | `ls` → present. |
| 21 | 04-ci-guards | Cross-ref to upstream `coding-guidelines-v20/spec/12 §03-reusable-ci-guards` | ✅ | URL well-formed. |
| 22 | 05-branch-protection | Repo-admin guidance for branch protection rules | ✅ | Process spec; no API surface to verify. |
| 23 | 06-cross-repo-sync | "`enum-v10` depends on `core-v9`" (line 11) | ⚠️→✅ | **D-CVS-51** — `enum-v10` is stale (project is now `enum-v10` after two renames). Fixed inline. |
| 24 | 06-cross-repo-sync | "`cross-repo/core-v9/README.md`" (line 19) | ⚠️→✅ | **D-CVS-52** — broken path (actual dir is `cross-repo/core-v9/`). Fixed inline with Core-memory clarification. |
| 25 | 06-cross-repo-sync | Comment template "Synced from github.com/alimtvnetwork/enum-v10/cross-repo/core-v9/" (line 80) | ⚠️→✅ | **D-CVS-53** — combines both stale tokens (`enum-v10` + `cross-repo/core-v9`). Fixed inline to `enum-v10/cross-repo/core-v9/`. |
| 26 | 06-cross-repo-sync | "both `enum-v10` and `core-v9` calling it via `uses:`" (line 91) | ⚠️→✅ | **D-CVS-54** — `enum-v10` stale. Fixed inline to `enum-v10`. |
| 27 | 06-cross-repo-sync | "See Also: `cross-repo/core-v9/README.md`" (line 103) | ⚠️→✅ | **D-CVS-55** — broken path. Fixed inline with Core-memory note. |
| 28 | 06-cross-repo-sync | §3 sync rules (workflows are source of truth, deltas documented, `actionlint` gate) | ✅ | Spec-internal best practice; consistent with `04-tooling/01-ci-pipeline.md`. |
| 29 | All files | Zero mojibake `core-v9 → core-v9` | ✅ | Zero hits. |
| 30 | All files | Zero `.lovable/user-preferences` citations | ✅ | Zero hits. |

**Tally:** 30 claims → ✅ 22 (after Cycle 17 fixes), ⚠️ 0, ❌ 0, ❓ 8.

**Score (verifiable subset):** 22 / 22 = **100.0%**.

## 3. Drift findings

**7 LOW drifts raised — all resolved in the same cycle.**

### D-CVS-49 — `00-overview.md` cites `cross-repo/core-v9/` but directory is `core-v8/`

**Severity:** LOW (broken link in landing index). **Locations:** lines 26 (Map table) + 80 (Maintenance §3). **Fix:** inline rewrite to `cross-repo/core-v9/` with explicit Core-memory note explaining the directory keeps its historical name even though the import path is `core-v9`.

### D-CVS-50 — `04-bootstrap-into-new-repo.md` §7 doesn't name the `enum-v10` test-layout precedent

**Severity:** LOW (already semantically correct — claims `tests/integratedtests/` is **not** required). **Location:** line 242. **Fix:** inline rewrite to name both upstream-`core-v9` (`tests/integratedtests/<pkg>tests/`) and `enum-v10` (`tests/creationtests/`) layouts as concrete examples. Closes the AH-tracked occurrence for this directory.

### D-CVS-51 — `06-cross-repo-sync.md` line 11 cites stale `enum-v10`

**Severity:** LOW. **Fix:** `enum-v10` → `enum-v10` (project went through two renames: `v1 → v2 → v3 → v4`; this occurrence was missed in earlier sweeps).

### D-CVS-52 — `06-cross-repo-sync.md` line 19 cites broken `cross-repo/core-v9/README.md`

**Severity:** LOW. **Fix:** inline rewrite to `cross-repo/core-v9/README.md` with Core-memory clarification.

### D-CVS-53 — `06-cross-repo-sync.md` line 80 comment template combines two stale tokens

**Severity:** LOW (the template is meant to be copy-pasted into other repos, so propagates the drift). **Fix:** inline rewrite of both `enum-v10 → enum-v10` and `cross-repo/core-v9 → cross-repo/core-v9`.

### D-CVS-54 — `06-cross-repo-sync.md` line 91 cites stale `enum-v10`

**Severity:** LOW. **Fix:** `enum-v10` → `enum-v10`.

### D-CVS-55 — `06-cross-repo-sync.md` line 103 "See Also" cites broken `cross-repo/core-v9/`

**Severity:** LOW. **Fix:** inline rewrite to `cross-repo/core-v9/` with Core-memory clarification.

> **Aggregate:** 7 LOW drifts (D-CVS-49 → D-CVS-55) raised + resolved in one cycle. No HIGH or MEDIUM drift, no contradictions. The 5 cross-repo-path drifts (D-CVS-49, -52, -53, -55) all stem from the same root cause: the spec text was written assuming `cross-repo/<dirname>` would mirror the import path (`core-v9`), but the actual convention (per Core memory) keeps the historical `core-v8` directory name. Cycle 17 makes the convention explicit at every cite site.

## 4. Spec-internal consistency

Specifically checked-and-clean (after fixes):
- No `enum-v10` / `enum-v10` / `enum-v10` references remain (post-rename verified).
- No mojibake `core-v9 → core-v9`.
- No `.lovable/user-preferences` citations.
- All `cross-repo/...` paths resolve (only `cross-repo/core-v9/` exists).
- All `.github/workflows/*.yml` paths resolve.
- All inter-spec cross-refs resolve.
- No contradiction with the **`spec/01-app/` freeze** (`spec-v0.30.0`) — Cycle 17 touches only `spec/04-` files.
- No contradiction with Cycle 15 (`spec/06-testing-guidelines/`) or Cycle 16 (`spec/03-powershell-test-run/`) callouts.

## 5. Directory-level milestone — `spec/04-tooling/` baselined & closed

With Cycle 17, `spec/04-tooling/` is **baselined and closed at 100 % verifiable** with **7 LOW drifts (D-CVS-49 → D-CVS-55) raised and resolved in the same cycle**. Remaining 8 ❓ are workflow-internal behaviours (release-trigger conditions, dashboard-rendering details, PowerShell module-loading semantics) that need direct probes of `.github/workflows/*.yml` and `scripts/*.psm1` — not blocking the directory closure.

| File | Status |
|---|---|
| `00-overview.md` | ✅ Closed (D-CVS-49 fixed at 2 sites) |
| `01-ci-pipeline.md` | ✅ Closed at baseline (1 ❓ behavioural) |
| `02-release-pipeline.md` | ✅ Closed at baseline (2 ❓ behavioural) |
| `02-powershell-dashboard-ui.md` | ✅ Closed at baseline (3 ❓ cosmetic/behavioural) |
| `03-vulnerability-scanning.md` | ✅ Closed at baseline (no findings) |
| `03-powershell-implementation.md` | ✅ Closed at baseline (1 ❓ behavioural) |
| `04-bootstrap-into-new-repo.md` | ✅ Closed (D-CVS-50 fixed) |
| `04-ci-guards.md` | ✅ Closed at baseline (no findings) |
| `05-branch-protection.md` | ✅ Closed at baseline (no findings) |
| `06-cross-repo-sync.md` | ✅ Closed (D-CVS-51 → D-CVS-55 fixed at 5 sites) |

## 6. Carry-forward

- **AH** — `spec/04-tooling/` debt cleared this cycle. Final residual AH item: `spec/02-app-issues/02-internal-package-coverage-policy.md` (folds into Cycle 18).
- **AB** — 8 ❓ workflow-internal claims need direct probes of `.github/workflows/*.yml` (release-pipeline triggers, dashboard rendering) and `scripts/*.psm1` (module loading). Not blocking; can fold into a future "workflows + scripts audit" cycle.
- **Suggestion** — the 5 `cross-repo/core-v9 → cross-repo/core-v9` drifts (D-CVS-49, -52, -53, -55) point at a teaching-friction in the spec corpus: readers who know the import path is `core-v9` will instinctively type the wrong directory name. Consider adding a `cross-repo/core-v9/README.md` top-of-file note explaining the historical naming. Tracked in suggestions as **S-104**.
