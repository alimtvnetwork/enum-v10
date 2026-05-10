# Release Pipeline Specification

## Overview

Library release pipeline for `enum-v10`. Triggers on `release/**` branches and `v*` tags. Produces source archives + checksums + a GitHub Release with an extracted changelog section.

This is a **Go library** — no cross-compiled binaries, install scripts, or code signing.

## Trigger

```yaml
on:
  push:
    branches: ["release/**"]
    tags:     ["v*"]
```

## Concurrency & Permissions

```yaml
concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false   # never cancel a release run

permissions:
  contents: write             # required to create the GitHub Release
```

## Jobs

### 1. validate

- `actions/checkout@v4`
- `actions/setup-go@v5` with `go-version-file: go.mod`
- `go vet ./...`
- `go build ./...`
- `go test -count=1 ./...`

### 2. release

| Step | Purpose |
|------|---------|
| Resolve version | `refs/tags/vX.Y.Z` or `refs/heads/release/vX.Y.Z` → `vX.Y.Z` |
| Package source archives | `dist/enum-v10-<ver>-source.tar.gz`, `.zip`, `checksums.txt` (SHA256) |
| Extract changelog | `awk` over `CHANGELOG.md` for the matching `## [vX.Y.Z]` block; conventional-commits fallback when missing |
| Assemble release body | Changelog + release-info table + `go get` install snippet + checksums block |
| Publish | `softprops/action-gh-release@v2` with auto pre-release detection (`-` in version → prerelease + not latest) |

## Release Body Structure

```
<Changelog section / commit-derived bullets>

---

### Release info
| Field | Value |
| Version | `vX.Y.Z` |
| Commit  | `<short-sha>` |
| Ref     | `<github-ref>` |
| Build Date | YYYY-MM-DD HH:MM:SS UTC |
| Module  | `github.com/alimtvnetwork/enum-v10` |

### Install
```bash
go get github.com/alimtvnetwork/enum-v10@vX.Y.Z
```

### Checksums (SHA256)
<contents of checksums.txt>
```

## Pre-release Detection

| Tag | Type | `prerelease` | `make_latest` |
|-----|------|--------------|---------------|
| `v1.2.3` | Stable | `false` | `true` |
| `v1.2.3-rc.1` | Pre-release | `true` | `false` |

## Constraints

- All actions pinned to exact tags — no `@latest` / `@main`.
- Release concurrency never cancels (`cancel-in-progress: false`).
- Source archives exclude `.git/`, `node_modules/`, `dist/`, `data/`, `tmp/`.
- The release body is assembled from `CHANGELOG.md` first, commits second — never both.

## Cross-References

- [`spec/12 — shared conventions`](https://github.com/alimtvnetwork/coding-guidelines-v20/blob/main/spec/12-cicd-pipeline-workflows/01-shared-conventions.md)
- [`spec/12 — GitHub release standard`](https://github.com/alimtvnetwork/coding-guidelines-v20/blob/main/spec/12-cicd-pipeline-workflows/02-github-release-standard.md)
- [`./01-ci-pipeline.md`](./01-ci-pipeline.md)
- [`./03-vulnerability-scanning.md`](./03-vulnerability-scanning.md)
