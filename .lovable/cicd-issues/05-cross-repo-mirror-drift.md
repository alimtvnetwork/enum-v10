# 05 — `cross-repo/core-v9/` Mirror Drift Risk

## Symptom

CI behaviour differs between the main `enum-v10` repo and the upstream `core-v8` GitHub repo because their `.github/workflows/`, `.golangci.yml`, baselines, or CI guard scripts have drifted.

## Root Cause

`cross-repo/core-v9/` is a manually-maintained mirror. Whenever the main repo's CI surface changes, the mirror must be updated in the same commit. If the AI forgets, drift accumulates silently until the next CI run on the upstream repo fails differently from `enum-v10`'s.

## Fix / Workaround

When changing any of these files in the main repo, mirror the change to `cross-repo/core-v9/` in the same commit:

- `.github/workflows/*.yml`
- `.golangci.yml`
- `.github/dependabot.yml`
- `.ci-baselines/*`
- `scripts/ci/test_*.py`
- `CHANGELOG.md` (with `core-v9` → `core-v8` substitution in the mirrored copy)

Do NOT:
- Rename `cross-repo/core-v9/`.
- Rewrite `core-v8` → `core-v9` inside this directory.
- Rewrite `enum-v10` → `enum-v10` inside this directory.

The actual `git push` to the upstream `core-v8` GitHub repo is **Task A** (manual user action, credential-bound).

## Status

🔄 Procedural — handled per-change.

## Related

- `.lovable/memory/05-cross-repo-mirror.md`
- Task A in `.lovable/plan.md`
