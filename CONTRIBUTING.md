# Contributing to enum-v9

## Local Development

Use the bundled runner for a workflow that mirrors CI:

| Task | macOS / Linux | Windows |
|------|---------------|---------|
| All tests | `./run.sh t` | `./run.ps1 T` |
| Coverage (HTML report) | `./run.sh tc` | `./run.ps1 TC` |
| Single package | `./run.sh tp <pkg>` | `./run.ps1 TP <pkg>` |
| Show last failures | `./run.sh tf` | `./run.ps1 TF` |
| Format | `./run.sh f` | `./run.ps1 F` |
| Vet | `./run.sh v` | `./run.ps1 V` |
| `go mod tidy` | `./run.sh ty` | `./run.ps1 TY` |

Coverage HTML lands at `data/coverage/coverage.html`.

## Pre-Push Checklist

Before opening a PR, verify each item below. The PR template (`.github/PULL_REQUEST_TEMPLATE.md`) mirrors this list — tick boxes there as you go.

- [ ] `go vet ./...` clean
- [ ] `golangci-lint run --timeout=5m` clean (or baseline updated — see [`.ci-baselines/README.md`](.ci-baselines/README.md))
- [ ] `./run.sh tc` (or `./run.ps1 TC`) passes with **≥ 60% coverage**
- [ ] `python3 scripts/ci/check-collisions.py .` clean
- [ ] `CHANGELOG.md` `[Unreleased]` updated for user-visible changes
- [ ] Spec docs updated if behaviour or tooling changed

The CI workflow (`.github/workflows/ci.yml`) enforces these — run them locally first to avoid round-trips:

```bash
go vet ./...
golangci-lint run --timeout=5m
./run.sh tc                                 # ≥ 60% coverage required
python3 scripts/ci/check-collisions.py .    # cross-file identifier audit
```

## Commit Conventions

Use **Conventional Commits**. The release pipeline categorises commits when no `CHANGELOG.md` entry exists for the version:

| Prefix | Section |
|--------|---------|
| `feat:` | Features |
| `fix:` | Bug Fixes |
| `refactor:`, `chore:`, `docs:`, `ci:`, `build:`, `test:`, `perf:`, `style:` | Maintenance |

## Changelog

User-visible changes go under `## [Unreleased]` in `CHANGELOG.md`. At release time, rename the section to `## [vX.Y.Z]` so the release pipeline can extract it.

## Releasing

1. Update `CHANGELOG.md` — move `[Unreleased]` items under `## [vX.Y.Z]`.
2. Push a `release/vX.Y.Z` branch, **or** tag the commit `vX.Y.Z`.
3. The release workflow validates, packages source archives, generates checksums, and publishes the GitHub Release.

Pre-release versions (anything containing `-`, e.g. `v1.0.0-rc.1`) are auto-flagged as pre-releases.

## Workflows

| File | Purpose |
|------|---------|
| `.github/workflows/ci.yml` | Lint, vulncheck, 4-suite test matrix, coverage gate, build check |
| `.github/workflows/ci-guards.yml` | Cross-file collision audit + baseline-diff lint gate |
| `.github/workflows/vulncheck.yml` | Standalone weekly `govulncheck` |
| `.github/workflows/release.yml` | Source archives + checksums + GitHub Release |

All actions are pinned to exact version tags — **never** use `@latest` or `@main`.

## Spec References

- [`spec/04-tooling/00-overview.md`](spec/04-tooling/00-overview.md) — tooling pipeline overview (start here)
- [`spec/04-tooling/01-ci-pipeline.md`](spec/04-tooling/01-ci-pipeline.md)
- [`spec/04-tooling/02-release-pipeline.md`](spec/04-tooling/02-release-pipeline.md)
- [`spec/04-tooling/03-vulnerability-scanning.md`](spec/04-tooling/03-vulnerability-scanning.md)
- [`spec/04-tooling/04-ci-guards.md`](spec/04-tooling/04-ci-guards.md)
- [`spec/04-tooling/05-branch-protection.md`](spec/04-tooling/05-branch-protection.md)
- [`spec/04-tooling/06-cross-repo-sync.md`](spec/04-tooling/06-cross-repo-sync.md)
