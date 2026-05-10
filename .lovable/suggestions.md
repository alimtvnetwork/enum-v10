# Suggestions

## Active Suggestions

### AP — Brief `spec/01-app/` thaw to clear AJ-45 + S-002

- **Status:** Pending — needs explicit user authorisation to thaw the freeze.
- **Priority:** Medium (eliminates the only HIGH-severity finding still on the scoreboard, C-CVS-65).
- **Description:** Thaw `spec/01-app/` for one cycle, apply (a) AJ-45 — purge fabricated `errcore.OverflowType.Fmt` from `spec/01-app/09-converters.md:161`, replace with `errcore.Expected.MessageVarMap`; (b) S-002 — change "duplicates" → "delegates to" in §09 row 64. Re-freeze immediately. Both are verbatim spec-text edits with zero source/code impact.
- **Added:** 2026-05-10 (Cycle 49).

### AQ — Lift `osdetect` past 77.4% with Linux-branch probes

- **Status:** Pending.
- **Priority:** Low (currently above the per-package CI floor; this is opportunistic).
- **Description:** AL-08 deliberately skipped platform-specific paths in `osdetect/`. Add `runtime.GOOS == "linux"`-guarded test probes for `linux.go`'s `defaultLinuxReleaseFile` parsing branches (lines 312/358/366) to push coverage from 77.4% toward ≥90%.
- **Added:** 2026-05-10 (Cycle 53).

### AU — Sweep §09 ⓘ advisory rows into `tests/contracttests/`

- **Status:** Pending.
- **Priority:** Low (remaining items are advisory, not active findings).
- **Description:** Row-172 carry-over and any other drained-but-untested behavioural claims in §09 should get pinned by contract tests in `tests/contracttests/`, mirroring the AJ-46 pattern. Locks the spec narrative against silent upstream behavioural drift.
- **Added:** 2026-05-10 (Cycle 50).

### AV — End-to-end test for `Test-UpstreamClone` (S-115)

- **Status:** Pending.
- **Priority:** Low (S-115 has unit smoke test; this would add an integration scenario).
- **Description:** Drive `Test-UpstreamClone` against a deliberately broken clone (sentinel-missing) inside CI to confirm the documented `{Ok, Path, Reason, PackageCount}` shape and that `Reason == "sentinel-missing"` is surfaced precisely.
- **Added:** 2026-05-10 (Cycle 53).

### AW — Per-package coverage trend warning

- **Status:** Pending.
- **Priority:** Medium (improves the AS gate without adding hard failures).
- **Description:** Companion to AS. Compare current per-package coverage to the previous green main build's profile and emit a CI warning (not failure) if any package drops by >5pp even when still ≥75%. Catches gradual erosion before it crosses the floor.
- **Added:** 2026-05-10 (Cycle 53).

### AY — Regression-lock test for RCA Pattern P10 (brace-tracking)

- **Status:** Pending.
- **Priority:** Medium (completes the static-guard trifecta with AS+AX).
- **Description:** Add a synthetic fixture under `scripts/ci/test_check_collisions.py` containing `var got` inside `t.Run(...)` body. Assert no false positive. Locks in the brace-depth tracking fix from RCA Pattern P10 (v1.15.0).
- **Added:** 2026-05-10 (Cycle 54).

### AZ — Unittest mitigation lock for RCA Pattern P8

- **Status:** Pending.
- **Priority:** Medium (completes the static-guard trifecta with AS+AX+AY).
- **Description:** Add Pester or a PowerShell smoke test that feeds `Test-IsCoverpkgWarningOnlyOutput` a representative warning-only blob and asserts `$true`. Locks in the AT-verified mitigation against silent removal.
- **Added:** 2026-05-10 (Cycle 54).

### Pin Go toolchain to 1.22 as a stopgap for Task W

- **Status:** Obsolete — Task W is now ✅ Done. This stopgap is no longer needed.
- **Priority:** N/A
- **Description:** Go 1.25 rejects the dual-path `replace` bridge between `core-v9` (import path) and `core-v8` (cached module path) with `used for two different module paths`. Pinning the toolchain to Go 1.22 in `go.mod` would unblock builds while waiting for the upstream `core-v9` `go.mod` rename + `v1.5.8` tag (Task W). Trade-off: locks the project to an older toolchain and silently masks the underlying issue.
- **Added:** 2026-05-05 (Cycle 13 turn).
- **Superseded:** 2026-05-05 — Task W completed, bridge removed, suggestion no longer applicable.

### Promote `errcore.VarTwoNoType` from ❓ to ✅ in Cycle 6's §08 audit

- **Status:** Pending (cross-cycle promotion — depends on Cycle 13's new "spec-internal-consistency" dimension).
- **Priority:** Low
- **Description:** Cycle 13 introduced spec-internal-consistency as an explicit verifiability dimension. `errcore.VarTwoNoType` was scored ❓ in Cycle 6 row 16 because no `enum-v8` consumer imports it, but it IS cross-referenced from `04-error-system.md:131`, `08-validators.md:240,307,329`, and `15-observability.md`. Under the new dimension it would promote to ✅. Backport mechanically when Task AC runs.
- **Added:** 2026-05-05 (Cycle 13 §15 audit, table row 4).

## Implemented Suggestions

### Add `cmd/main/` smoke-test policy carve-out to spec §12

- **Implemented:** 2026-05-05 (Cycle 10, fix C-CVS-10).
- **Notes:** Spec §12 used to assert "no `cmd/` directory anywhere"; reality: `enum-v8/cmd/main/main.go` is a single permitted smoke-test harness. §12 §1 rewritten as a "library-first, smoke-test allowed" policy. Cross-link to `cmd/README.md` added.

### Rewrite §06 around `SimpleSlice`/`Hashset`/`SimpleStringOnce` instead of fictional `coreonce.New.String`

- **Implemented:** 2026-05-04 (Cycle 4, fixes D-CVS-22/23/24).
- **Notes:** Original §06 §3 / §5 referenced symbols that don't exist. Rewrote around the actual top-level constructors used by `enum-v8`.

### Add consumer-coverage callouts wherever the spec describes upstream-only API

- **Implemented:** 2026-05-04 → 2026-05-05 (D-CVS-25 in §06, D-CVS-38 in §13, D-CVS-42 in §14).
- **Notes:** Three sections now carry explicit "this surface has no `enum-v8` consumer; verify via Task AB" callouts so future readers don't assume verified ✅ status incorrectly.
