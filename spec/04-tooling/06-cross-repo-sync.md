# 06 — Cross-Repo Sync (`cross-repo/`)

> ✅ **Status**: drafted at spec-v0.6.0 (2026-05-04, Asia/Kuala_Lumpur).
> **Audience**: maintainers who keep the staged copies in
> `cross-repo/<other-repo>/` aligned with the live workflows here.

---

## 1. Why this directory exists

`enum-v9` depends on `core-v9` and shares the same CI/CD philosophy with
it. Ideally both repos would consume the workflows from a single source
of truth (a reusable workflow or a `.github` org-level repo). Until that
exists, the workflows are **duplicated** with light per-repo
adaptations.

Lovable can only edit files inside this project, so adapted copies
intended for sister repos live under `cross-repo/<repo-name>/` and are
copied across by a maintainer (see `cross-repo/core-v9/README.md` for
the exact `cp` recipe — the directory intentionally retains its
`core-v8` name even though the upstream import path is now `core-v9`,
per Core memory).

These files are **not** loaded by GitHub Actions in this repo — only
files under the literal `.github/workflows/` path are. Anything inside
`cross-repo/` is inert here.

---

## 2. Layout

```
cross-repo/
└── <target-repo>/
    ├── README.md                  # Per-repo install instructions + diff log
    ├── .github/workflows/*.yml    # Adapted workflow copies
    ├── .ci-baselines/*.json       # Seed baselines (empty on first publish)
    └── .golangci.yml              # Adapted lint config (if it differs)
```

**Out of scope**: shared scripts (`scripts/ci/check-collisions.py`,
`scripts/ci/lint-baseline-diff.py`) are *not* duplicated. They are
repo-agnostic and copied verbatim from this repo's `scripts/ci/` —
documented in each `cross-repo/<repo>/README.md`.

---

## 3. Sync rules

1. **Workflows are the source of truth here.** Whenever
   `.github/workflows/*.yml` changes in this repo, re-derive the
   adapted copy under `cross-repo/<repo>/.github/workflows/` in the
   same commit.
2. **Per-repo deltas must be documented.** The target's `README.md`
   has a "What was changed vs. the originals" table. Update it whenever
   you change the adaptation.
3. **Run `actionlint` on the staged copies.** Same gate as the live
   workflows:
   ```bash
   nix run nixpkgs#actionlint -- cross-repo/<repo>/.github/workflows/*.yml
   ```
4. **Never reference Lovable-only paths** in staged files (e.g. don't
   leave `src/`, `vite.config.ts`, or any Lovable scaffolding
   references in adapted Go workflows).
5. **Treat `cross-repo/` as deliverables, not scratch.** Reviewers
   should be able to `cp -r cross-repo/<repo>/* /path/to/<repo>/` and
   open a PR with no further edits.

---

## 4. Promoting a staged file

When the target repo merges its first copy of these files, the
maintainer should:

1. Open the target-repo PR using the `README.md` recipe.
2. Once merged, leave the `cross-repo/<repo>/` directory in place here
   — it remains the canonical adaptation for future syncs.
3. Add a comment to the target-repo workflow header pointing back at
   this spec, e.g.:
   ```yaml
   # Synced from github.com/alimtvnetwork/enum-v9/cross-repo/core-v9/
   # See: spec/04-tooling/06-cross-repo-sync.md
   ```

---

## 5. End state (when this directory becomes obsolete)

Replace `cross-repo/` with one of:

- **A reusable workflow** in a shared `alimtvnetwork/.github` repo, with
  both `enum-v9` and `core-v9` calling it via `uses:`.
- **A `composite action`** packaged in its own repo for the lint /
  vulncheck / collision-audit steps.
- **A Go-module template repo** that downstream repos vendor.

Once any of those lands, delete `cross-repo/` and update
`CONTRIBUTING.md` to point at the new mechanism.

---

## See Also

- `cross-repo/core-v9/README.md` — concrete install recipe for the
  current target. (Directory keeps `core-v8` name; import path is `core-v9`.)
- [`01-ci-pipeline.md`](./01-ci-pipeline.md), [`02-release-pipeline.md`](./02-release-pipeline.md),
  [`03-vulnerability-scanning.md`](./03-vulnerability-scanning.md),
  [`04-ci-guards.md`](./04-ci-guards.md) — what the staged copies
  ultimately implement on the target side.
