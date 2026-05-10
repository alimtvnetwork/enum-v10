# 03 — Import Conventions

> ✅ **Status**: filled in audit Step 5 (2026-04-23, Asia/Kuala_Lumpur).
> **Audience**: AI agents and contributors writing code that consumes `core-v9` packages, or extending the framework itself.
> **Source**: extracted from [`/spec/00-llm-integration-guide.md`](../00-llm-integration-guide.md) §Import Conventions, §Package Map, plus `internal/` audit findings.

---

## Table of Contents

1. [Import Path Format](#1-import-path-format)
2. [Public Package Imports](#2-public-package-imports)
3. [The `internal/` Boundary](#3-the-internal-boundary)
4. [Test-Package Imports (`*tests` Suffix)](#4-test-package-imports-tests-suffix)
5. [Avoiding Cyclic Imports](#5-avoiding-cyclic-imports)
6. [Import Aliases — When and When Not](#6-import-aliases--when-and-when-not)
7. [Import Group Ordering](#7-import-group-ordering)

---

## 1. Import Path Format

The module path is `github.com/alimtvnetwork/core-v9`. Every public sub-package is reached by appending its directory:

```
github.com/alimtvnetwork/core-v9/<package-path>
```

Examples:

| Package | Import path |
|---|---|
| Root | `github.com/alimtvnetwork/core-v9` |
| `errcore` | `github.com/alimtvnetwork/core-v9/errcore` |
| `coredata/corejson` | `github.com/alimtvnetwork/core-v9/coredata/corejson` |
| `coreinterface/enuminf` | `github.com/alimtvnetwork/core-v9/coreinterface/enuminf` |
| `coreimpl/enumimpl` | `github.com/alimtvnetwork/core-v9/coreimpl/enumimpl` |

> **Rule**: Always use the **full path**. Do not introduce shortened module aliases at the `go.mod` level.

---

## 2. Public Package Imports

The canonical import block, copy-paste-ready:

```go
import (
    // Root package (generic factories)
    "github.com/alimtvnetwork/core-v9"

    // Foundation
    "github.com/alimtvnetwork/core-v9/conditional"
    "github.com/alimtvnetwork/core-v9/constants"
    "github.com/alimtvnetwork/core-v9/converters"
    "github.com/alimtvnetwork/core-v9/errcore"

    // Data structures
    "github.com/alimtvnetwork/core-v9/coredata/corejson"
    "github.com/alimtvnetwork/core-v9/coredata/corestr"
    "github.com/alimtvnetwork/core-v9/coredata/coregeneric" // optional — not every consumer needs this

    // Interfaces & implementations
    "github.com/alimtvnetwork/core-v9/coreinterface/enuminf"
    "github.com/alimtvnetwork/core-v9/coreimpl/enumimpl"

    // Predicates & setters
    "github.com/alimtvnetwork/core-v9/isany"
    "github.com/alimtvnetwork/core-v9/issetter"
)
```

For the full inventory of public packages, see [`01-package-map.md`](./01-package-map.md). Not every consumer uses every package — `enum-v9`, for example, currently uses 8 of the 11 listed canonical imports.

### Root package usage

The root `core` package contains generic slice/map factories. Import it with the bare path:

```go
import "github.com/alimtvnetwork/core-v9"

func example() {
    s := core.EmptySlice[int]()
    m := core.EmptyMapPtr[string, int]()
}
```

> The package name is `core`, not `corev9`. Even though the path ends in `core-v9`, Go uses the `package core` declaration in the source files.

---

## 3. The `internal/` Boundary

Go's compiler enforces that `internal/` packages can only be imported from packages **rooted at the `internal/` parent**. For `core-v9`, this means:

```
github.com/alimtvnetwork/core-v9/internal/...
```

is importable **only** by packages under `github.com/alimtvnetwork/core-v9/`. External consumers will get a compile error.

### Rules

| Scenario | Allowed? |
|---|---|
| Public package → public package | ✅ |
| Public package → `internal/...` package (same module) | ✅ |
| `internal/...` → public package (same module) | ✅ |
| `internal/...` → `internal/...` (same module) | ✅ |
| External consumer → any `internal/...` package | ❌ compile error |
| Test package (`*tests/`) → `internal/...` of the same module | ✅ |

### `internal/` access from tests

Test packages that live in **the same module** as an `internal/` package may import it directly. For example, in the upstream `core-v9` repo, tests under its own `tests/` tree can do:

```go
import "github.com/alimtvnetwork/core-v9/internal/reflectinternal"
```

> **Note (consumer repos like `enum-v9`):** This `internal/` access does **not** cross modules. A consumer (e.g. `enum-v9`) cannot import `core-v9/internal/...` because Go's `internal/` rule is enforced at module boundaries, not just package boundaries. Consumer-side tests should depend only on `core-v9`'s public API. See the upstream `core-v9` repo for examples of in-module `internal/` test imports.

> **Rule**: If you are tempted to add a re-export wrapper "to expose an internal helper", **don't**. Either move the helper to a public package, or accept that external consumers don't need it.

---

## 4. Test-Package Imports (`*tests` Suffix)

Test packages live under the repo's `tests/` directory in a sub-folder named after the test suite. In **`enum-v9`**, the layout is:

```
tests/
└── creationtests/
    ├── creation_test.go          // package creationtests
    ├── PathType_Creation_test.go
    ├── ScriptType_test.go
    ├── allBasicEnumsCollection.go
    ├── EnumTestWrapper.go
    └── …                          // shared fixtures + *_test.go files
```

In the upstream **`core-v9`** repo, the equivalent layout uses one folder per source package (e.g. `tests/<suite>/footests/` containing `package footests`). Both layouts share the same rules:

- Each test sub-folder is a **separate Go package** (different package name from any source package it tests).
- It imports the source package(s) under test as normal external dependencies.
- This avoids cyclic imports between source and test code.
- It can import `internal/` packages **only if the test sub-folder is in the same module** as those `internal/` packages (see §3).

### Import pattern in a `*_test.go`

```go
package errcoretests

import (
    "testing"

    // Source package under test
    "github.com/alimtvnetwork/core-v9/errcore"

    // Test framework
    "github.com/alimtvnetwork/core-v9/coretests"
    "github.com/alimtvnetwork/core-v9/coretests/args"
    "github.com/alimtvnetwork/core-v9/coretests/coretestcases"

    // Sometimes needed
    "github.com/alimtvnetwork/core-v9/issetter"
)
```

### Import pattern in a `_testcases.go`

```go
package errcoretests

import (
    // NEVER import "testing" here — testcases files hold pure data
    "github.com/alimtvnetwork/core-v9/coretests/args"
    "github.com/alimtvnetwork/core-v9/coretests/coretestcases"
)
```

> **Rule**: `_testcases.go` files **must not** import the `testing` package. Only `_test.go` files import `testing`. This keeps test data portable and serializable. See [`13-testing-patterns.md`](./13-testing-patterns.md) for the full Style A/B/C/D matrix.

### Shared test wrappers

Shared wrappers live under `tests/testwrappers/`:

```go
import "github.com/alimtvnetwork/core-v9/tests/testwrappers/stringstestwrapper"

type testWrapper = stringstestwrapper.StringsTestWrapper
```

See [`14-tests-folder-walkthrough.md` §2](./14-tests-folder-walkthrough.md#2-shared-test-wrappers) for the full wrapper inventory.

---

## 5. Avoiding Cyclic Imports

The package graph is **strictly layered** to prevent cycles. The layers, from low (no dependencies) to high (depends on everything below):

| Layer | Packages |
|---|---|
| L0 — primitives | `constants`, `isany`, `issetter` |
| L1 — foundation | `conditional`, `converters`, `coreinterface/*` |
| L2 — error system | `errcore` (depends on L0–L1) |
| L3 — data structures | `coredata/coregeneric`, `coredata/corestr` (depend on L0–L2) |
| L4 — JSON & enums | `coredata/corejson`, `coreimpl/enumimpl` (depend on L0–L3) |
| L5 — utilities | `corecmp`, `coremath`, `coresort`, `corevalidator`, `corefuncs` |
| L6 — testing | `coretests`, `coretests/args`, `coretests/coretestcases`, `coretests/results` |

### Rules

1. **A package may only import from strictly lower layers.** Same-layer imports are allowed only if there is no cycle.
2. **Interfaces (`coreinterface/*`) live in L1** so anything in L2+ can implement them without creating an upward edge.
3. **`errcore` is L2** because almost everything above it returns errors. Lower layers (`constants`, `isany`, `issetter`) must not import `errcore`.
4. **Test packages (L6) can import any source layer** but source layers must never import test packages.

### When you hit a cycle

If a new feature needs a function that would create an upward edge, do one of:

- **Move the function down** to a lower layer (most common fix).
- **Extract a shared interface** into `coreinterface/` and have both packages depend on the interface, not each other.
- **Re-shape the API** so the dependency flows in the natural direction.

Do **not** "break the cycle" with `init()` registration tricks — those hide the architectural problem.

---

## 6. Import Aliases — When and When Not

### Allowed aliases

| Case | Example | Why |
|---|---|---|
| Standard lib name collision | `gostrings "strings"` next to `corestr` | Disambiguation only |
| Shadowed identifier in scope | rare; document in a comment | Prevents bugs |

### Forbidden aliases

| Anti-pattern | Why bad |
|---|---|
| `e "github.com/alimtvnetwork/core-v9/errcore"` | Single-letter aliases destroy grep-ability |
| Renaming a package "for taste" | Inconsistent across files breaks IDE navigation |
| Aliasing to hide a long path | The path encodes the layer; hiding it hides design intent |

> **Rule**: Default to **no alias**. Only alias when Go's compiler forces you to.

---

## 7. Import Group Ordering

Standard `goimports` grouping with three groups separated by blank lines:

```go
import (
    // Group 1: Go standard library
    "fmt"
    "strings"
    "testing"

    // Group 2: Third-party (external modules)
    "github.com/some-other-org/some-lib"

    // Group 3: This module's packages
    "github.com/alimtvnetwork/core-v9"
    "github.com/alimtvnetwork/core-v9/errcore"
)
```

`gofmt -s` and `goimports` enforce this automatically. The CI lint step (see [`/spec/04-tooling/01-ci-pipeline.md`](../04-tooling/01-ci-pipeline.md)) will fail on misordered import groups.

---

## See Also

- [`01-package-map.md`](./01-package-map.md) — Authoritative inventory of all public packages
- [`02-design-philosophy.md`](./02-design-philosophy.md) — The 7 design pillars (zero-runtime-deps, struct-as-namespace, etc.)
- [`13-testing-patterns.md`](./13-testing-patterns.md) — Test-side import requirements
- [`14-tests-folder-walkthrough.md`](./14-tests-folder-walkthrough.md) — Shared wrappers and `coretests.GetAssert`
- [`/spec/00-llm-integration-guide.md` §Import Conventions](../00-llm-integration-guide.md#import-conventions) — Quick reference
