# Per-package coverage gate (75% floor)

## Symptom
Total-coverage gate alone allows individual packages to silently regress below 75%, eroding the AO uplift gains over time.

## Root Cause
`.github/workflows/ci.yml` `test-summary` job's "Coverage gate (75%)" step parsed only the `total:` line of `go tool cover -func`. A PR could bring `osdetect` from 77% → 50% as long as the totals stayed above the floor.

## Fix / Workaround
Cycle 53 (v1.30.0, Task AS):

- Added `scripts/ci/check-package-coverage.py` — parses raw `coverage.out` (the `mode: atomic` profile, NOT the `-func` summary), aggregates statement counts per package, fails non-zero if any non-ignored package falls below `--threshold` (default 75.0).
- Added `scripts/ci/test_check_package_coverage.py` — 6 unittest cases (parsing, pass, fail, ignore, missing input, malformed input). Wired automatically into existing `unittest discover -s scripts/ci -p 'test_*.py'` harness.
- Added `Coverage gate per-package (75%)` step in `test-summary` job, runs immediately after the existing total-coverage gate.

## Status
✅ Working as designed. Locks in AO uplift permanently. Use `--ignore <import-path>` for opt-out packages (none currently needed).

## Related
- Task AO (Cycle ≤49) — original per-package uplift to ≥75%.
- Task AS (Cycle 53) — this gate.
- `.lovable/cicd-issues/07-coverage-gate-60.md` — older total-coverage gate (still in place at 75%).
