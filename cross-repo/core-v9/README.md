# cross-repo/core-v9/

> ℹ️ **Historical naming — intentional. Do not rename this directory.**
>
> This folder mirrors a separate upstream repo whose name is still
> `alimtvnetwork/core-v8`. Even though the **import path used by `enum-v10`
> source code is `github.com/alimtvnetwork/core-v9`** (renamed 2026-05-05,
> tagged `v1.5.8`), the directory name here keeps the historical `core-v8`
> token so it stays in lockstep with the upstream repo it mirrors.
>
> **When editing spec text or scripts that reference *this directory*,
> always write `cross-repo/core-v9/`** — even when the surrounding sentence
> is about `core-v9` content. The mismatch is by design.
>
> Historical body references below (`enum-v1`, `core-v8`) likewise track the
> mirrored repo's vintage and must NOT be rewritten — see Core memory.

Adapted CI/CD workflows and CI-guard scripts ready to be copied into the
upstream **`alimtvnetwork/core-v8`** repository. They are **not** wired
into this `enum-v1` repo's own CI — they live here only because the
Lovable workspace can't directly edit the other GitHub repo.

## How to apply

From a local clone of `core-v8`:

```bash
git checkout -b ci/initial-pipeline

# Workflows
mkdir -p .github/workflows
cp <enum-v1>/cross-repo/core-v9/.github/workflows/*.yml .github/workflows/

# CI guard scripts (shared, repo-agnostic)
mkdir -p scripts/ci
cp <enum-v1>/scripts/ci/check-collisions.py     scripts/ci/
cp <enum-v1>/scripts/ci/lint-baseline-diff.py   scripts/ci/

# Lint baseline seed
mkdir -p .ci-baselines
cp <enum-v1>/cross-repo/core-v9/.ci-baselines/golangci-lint.json .ci-baselines/

# Repo metadata
cp <enum-v1>/cross-repo/core-v9/.golangci.yml .
cp <enum-v1>/cross-repo/core-v9/.github/dependabot.yml .github/

git add .github scripts/ci .ci-baselines .golangci.yml
git commit -m "ci: add CI/CD pipeline, vulncheck, release, guards, and Dependabot"
git push -u origin ci/initial-pipeline
```

Then open a PR. Once merged, configure the same branch-protection rules
described in `enum-v1`'s `spec/04-tooling/05-branch-protection.md`.

## What was changed vs. the `enum-v1` originals

| File | Adaptation for core-v8 |
|---|---|
| `ci.yml` | Test matrix replaced with a **single `./...` shard** — `core-v8` has 100+ packages and no natural curated split. Coverage gate kept at 60%. No `replace` directive juggling. |
| `vulncheck.yml` | Identical — the standalone weekly scan is repo-agnostic. |
| `release.yml` | Archive name changed from `enum-v1-…` to `core-v8-…`. Otherwise identical. |
| `ci-guards.yml` | Identical — references shared scripts under `scripts/ci/`. |
| `.golangci.yml` | Same baseline config; `core-v8` may want to tighten `enabled` linters over time. |
| `.ci-baselines/golangci-lint.json` | Empty seed — gate runs in warning-only mode until the first `main` push populates the cache. |
| `.github/dependabot.yml` | Identical — weekly `gomod` + `github-actions` updates, grouped by minor/patch. |

## Sensitivity / vulnerability scanning

Once `vulncheck.yml` is in place on `core-v8`:

- The weekly cron (Monday 09:00 UTC) runs `govulncheck` against the
  whole module independently of code changes — catches CVEs published
  *after* the last commit.
- The in-CI `vulncheck` job in `ci.yml` runs on every PR, with the same
  classification: third-party vulns fail the build, stdlib-only vulns
  emit a warning (since the fix requires bumping the Go toolchain).
- Reports are uploaded as artifacts (`vulncheck-report`) and retained
  for 30 days.
