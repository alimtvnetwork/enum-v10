# 01 — `go.mod` Rename Bridge Breaks Go 1.25 Builds

## Symptom

```
go: downloading github.com/alimtvnetwork/core-v8 v0.0.0-20260504132614-f978d4d81b5e
go: github.com/alimtvnetwork/enum-v9/accesstype imports
        github.com/alimtvnetwork/core-v9/coredata/coredynamic:
        github.com/alimtvnetwork/core-v8@... used for two different module paths
```

`./run.ps1 -tc` reports `STATUS ✗ BLOCKED` after `Pulling latest from remote / Fetching latest dependencies`.

## Root Cause

`enum-v9` source imports `github.com/alimtvnetwork/core-v9/...`. Upstream `core-v9` repo's own `go.mod` STILL declares `module github.com/alimtvnetwork/core-v8`. The `replace` directive in `enum-v9/go.mod` bridges the path but Go 1.25 rejects using one cached module artifact for two different declared paths.

## Fix / Workaround

- **Real fix (Task W):** Edit upstream `core-v9` `go.mod` to declare `module github.com/alimtvnetwork/core-v9`, sweep its source imports, tag `v1.5.8`. Then drop the `replace` bridge in `enum-v9/go.mod` (Task AG).
- **Stopgap (offered, not accepted):** Pin Go toolchain to 1.22 in `enum-v9/go.mod`. Older Go is more permissive about dual-path module use. Trade-off: locks the project to an older toolchain.

## Status

🚫 Blocked on Task W (manual upstream action by user).

## Related

- `.lovable/memory/02-go-mod-bridge.md`
- `.lovable/cicd-issues/02-internal-package-rule.md` (downstream symptom)
- `.lovable/cicd-issues/03-pseudo-version-rejected.md` (failed workaround)
