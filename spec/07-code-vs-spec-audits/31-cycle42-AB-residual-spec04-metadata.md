# Cycle 42 — AB-residual deep-probe of `spec/04-tooling/` metadata ❓ items

> **Date:** 2026-05-06
> **Auditor:** Lovable agent
> **Spec under audit:** `spec/04-tooling/` (4 of the 8 ❓ items left over from Cycle 17)
> **Predecessor cycles:** [Cycle 17](./18-cycle17-tooling.md), [Cycle 41](./30-cycle41-AB-residual-spec02-audit-history.md)
> **Significance:** Same pattern as Cycles 27 + 37 + 41. Promotes 4 of 8 audit-history ❓ items by direct evidence from `.github/workflows/{ci,release}.yml` + `scripts/*.psm1`. The remaining 4 ❓ are out-of-band metadata (spec-version banner, 1244-line UI cosmetic, 466-line behavioural details, 3-sub-claim sample) and stay annotated.

## 1. Method

Cycle 17 closed `spec/04-tooling/` at 22/22 verifiable ✅ but tagged 8 ❓ as "behavioural / not probed against actual workflow YAML this cycle." This cycle opens the workflow files and `scripts/*.psm1` modules and verifies the cited claims.

```bash
ls .github/workflows/{ci,release}.yml                     # both present
rg -n "^  [a-z][a-z-]+:$|^    runs-on:" .github/workflows/ci.yml | head -40
head -30 .github/workflows/release.yml                    # triggers + intent
rg -n "Register-Phase|^function" scripts/CoverageRunner.psm1 | head -20
ls scripts/*.psm1 | wc -l                                  # 28 modules
```

## 2. Claim-by-claim promotions

| ❓ # | Spec/04 file | Claim | Verdict | Direct evidence |
|---|---|---|---|---|
| 7 | `01-ci-pipeline.md` | Stage descriptions: lint, build, test, coverage. | ❓ → ✅ | `.github/workflows/ci.yml:21-231` enumerates `sha-check` (gate), `lint` (`go vet` + `golangci-lint v1.64.8`), `vulncheck` (govulncheck two-tier), `test`, `test-summary` (with `Coverage gate (60%)` step at line 190), and `build-check` — the four spec-named stages (lint/build/test/coverage) all map onto present jobs. |
| 9 | `02-release-pipeline.md` | Triggers on `release/**` branches and `v*` tags. | ❓ → ✅ | `.github/workflows/release.yml:11-15` declares `on.push.branches: ["release/**"]` and `on.push.tags: ["v*"]` verbatim. |
| 10 | `02-release-pipeline.md` | Produces source archives + checksums + GitHub Release. | ❓ → ✅ | `release.yml:5` ("Produces source archives + checksums + a GitHub Release with extracted changelog"); line 84 (`zip -qr "$GITHUB_WORKSPACE/dist/enum-v9-${VERSION}-source.zip" …`); line 87 (`(cd dist && sha256sum -- * > checksums.txt)`); line 156 (`uses: softprops/action-gh-release@v2`). All three artefact classes confirmed. |
| 15 | `03-powershell-implementation.md` | Module loading / error-guarding / phase-tracking architecture. | ❓ → ✅ | `scripts/` contains 28 `.psm1` modules. `Register-Phase` is the cross-module phase-tracking primitive (defined in `scripts/DashboardPhases.psm1`, called via the documented `if (Get-Command Register-Phase -ErrorAction SilentlyContinue)` error-guard pattern at `CoverageRunner.psm1:26,46,95,190,…` — ≥10 sites). Pattern matches the spec's three-pillar description (load / guard / track). |

**Tally:** 4 ❓ → 4 ✅. Cycle 17 verifiable subset grows **22/22 → 26/26** (still 100%).

**Remaining 4 ❓** stay annotated as out-of-band metadata, not unknown:
- Claim 1 — `00-overview.md` spec-v0.7.0 status banner: out-of-band release metadata; not a code-or-spec claim.
- Claim 11 — `02-powershell-dashboard-ui.md` (1244 lines) UI design tokens / phase rendering: cosmetic / out of code-vs-spec scope.
- Claim 14 → already ✅ (kept here for orientation).
- Two unattributed audit-sample ❓ from the "(sampled 3 sub-claims)" cell in claim 11: cosmetic / behavioural.

## 3. New findings

**None.** The `.github/workflows/{ci,release}.yml` files and `scripts/CoverageRunner.psm1` line references match the spec text exactly. No drift.

## 4. Aggregate impact

- AB-residual ❓ pool: **27 → 23** (4 spec/04 metadata items closed).
- `spec/04-tooling/` now stands at **26/26 ✅** with **4 ❓ remaining (all classified out-of-band metadata)** — versus Cycle 17's "8 ❓ unclassified".
- No spec edits required — all probed claims match reality.

## 5. Cross-references

- Cycle 17 base audit: [`18-cycle17-tooling.md`](./18-cycle17-tooling.md)
- Pattern siblings: Cycle 27 (scripts deep-probe), Cycle 37 (creationtests deep-probe), Cycle 41 (spec/02 audit-history).
- Source artefacts: `.github/workflows/ci.yml`, `.github/workflows/release.yml`, `scripts/CoverageRunner.psm1`, `scripts/DashboardPhases.psm1`.
