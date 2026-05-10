# CI Guards

Reusable, language-agnostic CI guards from [`spec/12 §03-reusable-ci-guards`](https://github.com/alimtvnetwork/coding-guidelines-v20/tree/main/spec/12-cicd-pipeline-workflows/03-reusable-ci-guards), implemented for `enum-v10`.

Workflow: [`.github/workflows/ci-guards.yml`](../../.github/workflows/ci-guards.yml)

## Guard 1 — Cross-File Collision Audit (Pattern 03)

Script: [`scripts/ci/check-collisions.py`](../../scripts/ci/check-collisions.py)

Detects three categories of identifier collisions **per Go package**:

| # | Category | Example |
|---|----------|---------|
| 1 | Cross-file exact | `Ranges` declared in two files in the same package |
| 2 | Case-insensitive | `FileReader` vs `fileReader` in the same package |
| 3 | Intra-file duplicates | Same name declared twice in one file |

### Behavior

- Methods (`func (r Recv) Name`) are intentionally excluded — Go allows them across receivers.
- String-literal aware: skips identifiers inside `"..."`, `` `...` ``, and `// ` / `/* */` comments.
- Tracks `const ( ... )` / `var ( ... )` blocks so grouped declarations are picked up.
- Scope is **per package directory** — `accesstype.Ranges` and `brackets.Ranges` are not collisions.

### Exit codes

| Code | Meaning |
|------|---------|
| 0 | No collisions |
| 1 | Findings present (full report on stdout) |
| 2 | No `.go` files matched |

## Guard 2 — Baseline-Diff Lint Gate (Pattern 04)

Script: [`scripts/ci/lint-baseline-diff.py`](../../scripts/ci/lint-baseline-diff.py)

Runs `golangci-lint --out-format json`, then diffs the result against a cached baseline:

| State | Behavior |
|-------|----------|
| Baseline missing (first run) | **Seeding mode** — emit `::warning::` per finding, exit 0 |
| Baseline present, no NEW findings | Pass |
| Baseline present, NEW findings | `::error::` per new finding, exit 1 |

### Finding identity

`(file, line, linter, message)` — severity, column, and source snippets are excluded so linter version bumps don't produce spurious diffs.

### Baseline storage

GitHub Actions cache, key `lint-baseline-v1.64.8-refs/heads/main`. Refreshed on every successful push to `main`. Cache miss = seeding mode (acceptable).

## Constraints

- `golangci-lint` and `govulncheck` versions are pinned.
- Guards run in a separate workflow (`ci-guards.yml`) so failures don't block the main CI's coverage gate.
- Both scripts have stable exit codes for local use (`python3 scripts/ci/check-collisions.py .`).
