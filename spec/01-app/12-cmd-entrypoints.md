# 12 — CMD Entrypoints

> ✅ **Status**: filled in audit Step 5c (2026-04-23, Asia/Kuala_Lumpur).
> **Audience**: anyone looking for a `main` function or a CLI binary in this module.
> **Short answer**: **there isn't one** — this module is a pure library.

---

## Table of Contents

1. [No `cmd/` Directory — On Purpose](#1-no-cmd-directory--on-purpose)
2. [How to Use the Library](#2-how-to-use-the-library)
3. [Where to Find Runnable Examples](#3-where-to-find-runnable-examples)
4. [If You Need a CLI](#4-if-you-need-a-cli)
5. [Tooling vs Entrypoints](#5-tooling-vs-entrypoints)

---

## 1. `cmd/` Policy — Library-First, Smoke-Test Allowed

Both `core-v9` (the upstream library this module depends on) and `enum-v9` (this module) are **pure Go enum/primitive libraries**. Their job is to provide reusable types that downstream applications import — not to ship binaries.

**Upstream `core-v9`** intentionally has:

- No `cmd/` directory.
- No `main` package.
- No produced binary artefacts.

**This module (`enum-v9`)** ships a single, narrowly-scoped exception: `cmd/main/main.go` (see [`/cmd/README.md`](../../cmd/README.md)). It exists purely as a developer **smoke-test harness** — invoked locally via `make` to compile to `bin/main` and exercise a handful of enums (`brackets`, `dbaction`, `instructiontype`, `osdetect`, `strtype`) end-to-end. It is **not** a shipped CLI, is not referenced by any consumer, and is not produced by CI release artefacts.

Putting a real `main` package in either module would:

- Add deps that consumers don't need.
- Encourage tight coupling between the library and a specific deployment shape.
- Conflate the package map (see [`01-package-map.md`](./01-package-map.md)) with binary inventories.

> **Rule**: The only `main` package permitted in this repo is the smoke-test harness at `cmd/main/`. Do not add additional `cmd/<name>/` entrypoints, and do not introduce any `cmd/` package in the upstream `core-v9` module. New consumer-facing entrypoints belong in the consuming application's repo.
>
> **Enforcement** *(F-V14-04 fix)*: this rule is **PR-review-enforced**, not machine-enforced — there is no `go vet` lint, no CI check, and no PowerShell guard in `spec/04-tooling/` that blocks an additional `cmd/<name>/` directory from being added. If you need machine enforcement, open a ticket in [`spec/02-app-issues/`](../02-app-issues/) describing the desired guard (suggested implementations: a `pre-commit` shell check counting `cmd/*/` entries, or a custom `go vet` analyzer). Until such a ticket lands, code review is the single line of defence.

---

## 2. How to Use the Library

```go
import (
    "github.com/alimtvnetwork/core-v9/coredata/coregeneric"
    "github.com/alimtvnetwork/core-v9/conditional"
    "github.com/alimtvnetwork/core-v9/errcore"
)

func main() {
    col := coregeneric.New.Collection.String.Items("a", "b", "c")
    label := conditional.IfString(col.IsEmpty(), "empty", "non-empty")
    if err := doWork(); err != nil {
        return errcore.FailedType.Fmt("work failed: %v", err)
    }
}
```

See the package-by-package usage notes in:

- [`04-error-system.md`](./04-error-system.md)
- [`05-enum-system.md`](./05-enum-system.md)
- [`06-data-structures.md`](./06-data-structures.md)
- [`07-conditional-and-utilities.md`](./07-conditional-and-utilities.md)
- [`08-validators.md`](./08-validators.md)
- [`09-converters.md`](./09-converters.md)
- [`10-reflection-and-dynamic.md`](./10-reflection-and-dynamic.md)

---

## 3. Where to Find Runnable Examples

The closest thing to executable code in this module is the **shared creation-test suite** at `tests/creationtests/` (plus per-package `*_test.go` files). Every enum participates in the shared registry-driven contract tests; see C-CVS-01 / D-CVS-17 / D-CVS-26 / D-CVS-27 for prior fixes of this same `integratedtests` → `creationtests` drift.

```bash
# Run all tests
go test ./...

# Run the shared creation-test suite
go test ./tests/creationtests/...
```

The `cmd/main/` smoke-test harness (see §1) is a second, looser way to exercise the library by hand: `make` from the repo root compiles it to `bin/main`.

See [`14-tests-folder-walkthrough.md`](./14-tests-folder-walkthrough.md) for the layout and [`13-testing-patterns.md`](./13-testing-patterns.md) for how the tests are structured.

---

## 4. If You Need a CLI

Build it in your own repository as a thin wrapper:

```go
// myapp/cmd/myapp/main.go
package main

import (
    "os"
    "github.com/alimtvnetwork/core-v9/coredata/coregeneric"
    "github.com/alimtvnetwork/core-v9/conditional"
)

func main() {
    args := coregeneric.New.Collection.String.Items(os.Args[1:]...)
    cmd  := conditional.IfFuncString(
        args.IsEmpty(),
        func() string { return "help" },
        func() string { return args.First() },
    )
    dispatch(cmd, args)
}
```

This keeps the library import-only and lets each consumer choose its own CLI framework, flag parser, and configuration loader.

---

## 5. Tooling vs Entrypoints

The repo does include **PowerShell scripts** under [`/spec/04-tooling/`](../04-tooling/) (test runner, repo bootstrap). These are *tooling*, not Go entrypoints:

| Asset | Type | Location | Purpose |
|---|---|---|---|
| PowerShell test runner | `.ps1` script | `/spec/04-tooling/03-powershell-implementation.md` | Runs `go test`, parses output, attributes failures |
| Repo bootstrap | `.ps1` script | `/spec/04-tooling/04-bootstrap-into-new-repo.md` | Scaffolds a new consumer repo |
| Go test binaries | implicit | produced by `go test` | One per `_test.go` package |

None of these are `main` packages in this module. The PowerShell scripts are external orchestration, and the `go test` binaries are produced by the toolchain on demand.

---

## See Also

- [`00-repo-overview.md`](./00-repo-overview.md) — Module identity and top-level layout
- [`01-package-map.md`](./01-package-map.md) — Every package, none of which are `main`
- [`14-tests-folder-walkthrough.md`](./14-tests-folder-walkthrough.md) — Runnable code lives in `tests/creationtests/`
- [`/cmd/README.md`](../../cmd/README.md) — The single permitted `main` package (smoke-test harness)
- [`/spec/04-tooling/`](../04-tooling/) — PowerShell scripts (external tooling, not Go entrypoints)
- [`/spec/00-llm-integration-guide.md` §Quick Start](../00-llm-integration-guide.md) — Example library usage
