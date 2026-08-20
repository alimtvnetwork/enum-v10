# `cross-repo/core-v9/` Mirror

## Purpose

Mirrors CI surface to a separate upstream repo (`core-v8` GitHub repo). Whenever the main `enum-v10` repo's CI changes, the corresponding files must be updated here too.

## Files mirrored

- `.github/workflows/`
- `.golangci.yml`
- `.github/dependabot.yml`
- `.ci-baselines/`
- `scripts/ci/test_*.py`
- `CHANGELOG.md` (with `core-v9` → `core-v8` substitution)

## Hard rules

- **Do NOT rename `cross-repo/core-v9/`.** It tracks a different module; the name is part of the contract.
- **Do NOT rewrite `core-v8` → `core-v9` inside this directory.** This directory intentionally keeps the old name.
- **Do NOT rewrite `enum-v10` → `enum-v10` inside this directory.** It tracks a different module which is still on `enum-v10`.

## Push policy (Task A)

Pushing the mirror to the actual `core-v8` GitHub repo is a **manual user action** (credential-bound). The AI should:
1. Update files in `cross-repo/core-v9/` to mirror main-repo CI changes.
2. Commit them to `enum-v10`.
3. Mark Task A as ⏭️ (manual user action) — never attempt the actual push.
