# 02 — Go's `internal/` Rule Rejects `enum-v9` Consumers

## Symptom

Build fails with errors of the form:

```
use of internal package github.com/alimtvnetwork/core-v8/.../internal/X
not allowed
```

…even though `enum-v9` source imports the package as `github.com/alimtvnetwork/core-v9/.../internal/X`.

## Root Cause

Go's `internal/` rule is enforced against the **cached module's declared path** (`core-v8`), NOT the import path the consumer used (`core-v9`). To the toolchain, `enum-v9` looks like a foreign module reaching into a different module's `internal/`.

This is a **direct downstream consequence of issue #01**.

## Fix / Workaround

Same as #01: only Task W (upstream `go.mod` rename + `v1.5.8` tag) actually fixes this. The `replace` bridge cannot.

## Status

🚫 Blocked on Task W. Root cause = #01.

## Related

- `.lovable/cicd-issues/01-go-mod-rename-bridge.md` (root cause)
- `.lovable/memory/02-go-mod-bridge.md` § "Why the bridge is INSUFFICIENT"
