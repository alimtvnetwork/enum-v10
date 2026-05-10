# Upstream `core-v9` `go.mod` Module Path Mismatch

## Description

`enum-v10` source imports `github.com/alimtvnetwork/core-v9/...` but the upstream `core-v9` repository's own `go.mod` declared `module github.com/alimtvnetwork/core-v8`. A `replace` bridge in `enum-v10/go.mod` redirected resolution to the cached `core-v8` artifact, but Go 1.25 rejected this with `used for two different module paths`, and Go's `internal/` rule rejected any `enum-v10` consumer reaching into `core-v9/.../internal/...` because it saw the cached path as `core-v8`.

## Root Cause

Upstream `core-v9` repo had not yet had its `go.mod` rewritten to `module github.com/alimtvnetwork/core-v9`, nor had a release been tagged under the new path.

## Solution

Upstream `core-v9` repo was updated: `go.mod` now declares `module github.com/alimtvnetwork/core-v9` and was tagged `v1.5.8`. The `replace` directive in `enum-v10/go.mod` was removed. `go.mod` now carries a clean `require github.com/alimtvnetwork/core-v9 v1.5.8`.

## Iteration Count

5+ attempts (pseudo-version, replace bridges, toolchain pinning) before the real upstream fix landed.

## Learning

Go module path mismatches cannot be reliably bridged with `replace` directives in Go 1.25+. The only real fix is updating the upstream `go.mod` declaration.

## What NOT to Repeat

- Do not attempt pseudo-version `v1.5.6-0.<date>-<sha>` — needs a non-existent predecessor tag.
- Do not attempt `replace` directive bridges for modules with `internal/` consumers — Go enforces internal/ against the cached module's declared path.

## Related

- `.lovable/cicd-issues/01-go-mod-rename-bridge.md`
- `.lovable/cicd-issues/02-internal-package-rule.md`
- `.lovable/cicd-issues/03-pseudo-version-rejected.md`
- `.lovable/memory/02-go-mod-bridge.md`
