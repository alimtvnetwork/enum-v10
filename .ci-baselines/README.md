# CI Baselines

Cached snapshots used by the CI guards in `.github/workflows/ci-guards.yml`.

## `golangci-lint.json`

Baseline of accepted `golangci-lint` findings. The lint gate
(`scripts/ci/lint-baseline-diff.py`) compares a fresh report against
this file and fails only on **new** findings.

## `package-coverage.json`

Per-package coverage snapshot `{importPath: percent}` used by the
trend gate (`scripts/ci/check-package-coverage-trend.py`). Warning-only:
emits one `::warning::` per package whose coverage drops more than
`--tolerance` percentage points (default 1.0) vs the baseline. Refreshed
automatically on `main` pushes by the `Refresh per-package coverage
baseline` step in `.github/workflows/ci.yml`. Seed value is `{}` —
the script auto-detects seeding mode and only prints the current
snapshot until a real baseline is committed.

## Seed-Then-Gate Workflow

The baseline operates in two modes — driven entirely by whether
`.ci-baselines/golangci-lint.json` is missing or empty.

### Mode 1 — Seeding (onboarding)

**When:** First time enabling the gate on a repo with pre-existing
lint debt, OR after a deliberate baseline reset.

**State:** `golangci-lint.json` is missing, empty, or contains
`{"Issues": []}`.

**Behavior:**

- `lint-baseline-diff.py` emits one `::warning::` per finding.
- The script exits `0` — the gate **never blocks** in this mode.
- The `update-baseline` job on `main` writes the first real
  baseline at the next push.

**How to seed manually** (if you want a baseline now without
waiting for the next `main` push):

```bash
golangci-lint run --out-format=json ./... > .ci-baselines/golangci-lint.json
git add .ci-baselines/golangci-lint.json
git commit -m "ci: seed golangci-lint baseline"
```

### Mode 2 — Gating (steady state)

**When:** A non-empty baseline has been committed.

**State:** `golangci-lint.json` contains one or more issues.

**Behavior:**

- `lint-baseline-diff.py` computes three sets:
  - **NEW** — present in current report, absent from baseline → `::error::` + exit `1`
  - **FIXED** — present in baseline, absent from current report → informational
  - **UNCHANGED** — in both → silently accepted
- A PR fails iff it introduces any NEW finding.
- The `update-baseline` job on `main` regenerates the baseline
  after each merge so FIXED entries drop off and the debt monotonically
  shrinks (or holds flat).

### Mode transition

To move from seeding to gating:

1. Land any PR to `main` — the `update-baseline` job populates
   `.ci-baselines/golangci-lint.json` automatically, OR
2. Run the manual seed command above and commit the result.

To reset back to seeding (rare — only when the baseline becomes
unmaintainable, e.g. after a major refactor):

```bash
echo '{"Issues": []}' > .ci-baselines/golangci-lint.json
git commit -am "ci: reset lint baseline"
```

## Reviewing Baseline Changes

The `update-baseline` job on `main` regenerates the baseline
automatically. Reviewers should inspect the diff:

- **Net decrease** ✅ — fixes landed, debt shrinking.
- **Net increase** ⚠️ — likely means someone bypassed the gate
  (force-merge, admin override, or seeding-mode PR slipped through).
  Investigate whether the change should be reverted or the underlying
  issues fixed.

## Finding Identity

A finding is keyed by the 4-tuple `(file, line, linter, message)`.
Column and source snippet are intentionally excluded so cosmetic
re-formatting doesn't churn the baseline. Linter version bumps that
only change column positions or snippet text are also tolerated.

## Related

- Gate script: [`scripts/ci/lint-baseline-diff.py`](../scripts/ci/lint-baseline-diff.py)
- Tests: [`scripts/ci/test_lint_baseline_diff.py`](../scripts/ci/test_lint_baseline_diff.py)
- Workflow: [`.github/workflows/ci-guards.yml`](../.github/workflows/ci-guards.yml)
- Spec: [`spec/04-tooling/04-ci-guards.md`](../spec/04-tooling/04-ci-guards.md)
