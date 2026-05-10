# 01 — Package Map

> ✅ **Done** — extracted from `spec/00-llm-integration-guide.md` §Package Map.
> **Status**: filled in audit Step 3 (2026-04-23, Asia/Kuala_Lumpur).
> **Audience**: AI agents needing an authoritative inventory before importing or referencing a package.

---

## 1. Public Packages (alphabetical, one-line purpose)

| Package | Path suffix | Purpose |
|---|---|---|
| Root | `core-v9` | Generic slice/map factories (`core.go`, `generic.go`) |
| `bytetype` | `/bytetype` | Byte-type utilities |
| `chmodhelper` | `/chmodhelper` | File-permission parsing & verification |
| `conditional` | `/conditional` | Ternary helpers: `If[T]`, `IfFunc[T]`, `NilDef[T]` |
| `constants` | `/constants` | 400+ named constants (strings, bytes, runes, numbers) |
| `converters` | `/converters` | String↔bytes, maps, JSON formatting |
| `coreappend` | `/coreappend` | Append/prepend with nil-skip |
| `corecmp` | `/corecmp` | Typed comparison helpers |
| `corefuncs` | `/corefuncs` | Function-type wrappers (`ErrFunc`, `InOutErrFuncWrapper`) |
| `coremath` | `/coremath` | ⚠️ **LEGACY** — Min/Max for all numeric types. Do NOT use in new code; use Go 1.21+ built-in `min` / `max` / `clear`. Kept for backward compatibility. *(F-V12-02)* |
| `coresort` | `/coresort` | Quick sort (strings, integers) |
| `coreunique` | `/coreunique` | Uniqueness helpers |
| `corevalidator` | `/corevalidator` | Line, slice, text, range validators |
| `coreversion` | `/coreversion` | Semantic versioning |
| `defaultcapacity` | `/defaultcapacity` | Default capacity constants |
| `defaulterr` | `/defaulterr` | Default error types |
| `errcore` | `/errcore` | Rich error construction + stack traces |
| `filemode` | `/filemode` | File-mode types |
| `isany` | `/isany` | Type-checking predicates (`Null`, `Zero`, `DeepEqual`) |
| `issetter` | `/issetter` | Six-valued boolean |
| `keymk` | `/keymk` | Key compilation with legends/templates |
| `mutexbykey` | `/mutexbykey` | Per-key mutex locking |
| `namevalue` | `/namevalue` | Name-value pair types |
| `osconsts` | `/osconsts` | OS-specific constants |
| `ostype` | `/ostype` | OS type detection |
| `pagingutil` | `/pagingutil` | Paging/pagination utilities |
| `reflectcore` | `/reflectcore` | Reflection utilities |
| `regexnew` | `/regexnew` | Lazy-compiled regex with thread-safe caching |
| `typesconv` | `/typesconv` | Pointer↔value conversions |

---

## 2. Umbrella: `coredata/` — Data Structures

| Sub-package | Purpose |
|---|---|
| `coredata/coreapi` | Typed API request/response models |
| `coredata/coredynamic` | Reflection-based dynamic data |
| `coredata/coregeneric` | Generic `Collection[T]`, `Hashset[T]`, `Hashmap[K,V]` |
| `coredata/corejson` | JSON pipeline: `Serialize.*`, `Deserialize.*` |
| `coredata/coreonce` | Compute-once lazy values |
| `coredata/corepayload` | `PayloadWrapper` — structured data transport |
| `coredata/corerange` | Range types (int, byte) |
| `coredata/corestr` | String `Collection`, `Hashset`, `Hashmap` |
| `coredata/stringslice` | 80+ pure `[]string` manipulation functions |

---

## 3. Umbrella: `coreinterface/` — Contracts

| Sub-package | Purpose |
|---|---|
| `coreinterface` (root) | 100+ canonical interface contracts |
| `coreinterface/enuminf` | Enum interfaces (`BasicEnumer`, `BaseEnumer`, etc.) |
| `coreinterface/errcoreinf` | Error-wrapper interfaces |
| `coreinterface/loggerinf` | Logger interfaces |
| `coreinterface/serializerinf` | Serialization contracts |
| `coreinterface/entityinf` | Entity-level interfaces |
| `coreinterface/payloadinf` | Payload interfaces |

---

## 4. Umbrella: `coreimpl/` — Implementations

| Sub-package | Purpose |
|---|---|
| `coreimpl/enumimpl` | Enum implementation engine |
| `coreimpl/enumimpl/enumtype` | Enum-type metadata (`Variant`) |

---

## 5. Domain Enums

| Package | Purpose |
|---|---|
| `enums/stringcompareas` | String-comparison-area enum |
| `enums/versionindexes` | Version-index enum |

---

## 6. Test Framework (in-tree)

| Package | Purpose |
|---|---|
| `coretests/args` | Test argument types (`FuncWrap`, `Map`) |
| `coretests/coretestcases` | `CaseV1` test-case definitions |
| `coretests` (helpers root) | `GetAssert` and friends — see [`14-tests-folder-walkthrough.md`](./14-tests-folder-walkthrough.md) |

---

## 7. Internal (NOT importable externally)

| Package | Purpose |
|---|---|
| `internal/convertinternal` | Conversion helpers used by `converters/` |
| `internal/reflectinternal` | Reflection helpers used by `reflectcore/` |
| `internal/strutilinternal` | String helpers used by `corestr/` and friends |

> **Hard rule**: external code must never import any path under `internal/`. The Go compiler enforces this for code outside the module; this rule reminds us not to create thin re-export packages that would defeat the boundary.

---

## 8. Test Packages (under `tests/creationtests/`)

> ⚠️ **Scope:** the per-package mirror layout below is **upstream `core-v9`**. `enum-v10` uses one shared `tests/creationtests/` package — see [`13-testing-patterns.md` §6.1](./13-testing-patterns.md#61-enum-v10-specific-layout).

There is a **mirror** of public packages under upstream `tests/creationtests/` with the suffix `tests`. Examples:

- `tests/creationtests/argstests/`
- `tests/creationtests/errcoretests/`
- `tests/creationtests/anycmptests/`
- `tests/creationtests/isanytests/`
- … (≥ 50 packages total)

Naming rules and file layout for these packages are documented in [`/spec/06-testing-guidelines/01-folder-structure.md`](../06-testing-guidelines/01-folder-structure.md) and [`14-tests-folder-walkthrough.md`](./14-tests-folder-walkthrough.md).

Shared wrappers live under `tests/testwrappers/`. Their public surface is currently undocumented — tracked as [`02-app-issues/04-testwrappers-public-surface.md`](../02-app-issues/04-testwrappers-public-surface.md).

---

## 9. Cross-Reference: Package → Owning Topic File

| Package family | Detail file |
|---|---|
| Root, conditional, constants, isany, issetter, regexnew | [`07-conditional-and-utilities.md`](./07-conditional-and-utilities.md) |
| `errcore`, `coreinterface/errcoreinf` | [`04-error-system.md`](./04-error-system.md) |
| `coreinterface/enuminf`, `coreimpl/enumimpl`, `enums/*` | [`05-enum-system.md`](./05-enum-system.md) |
| `coredata/*` umbrella | [`06-data-structures.md`](./06-data-structures.md) |
| `corevalidator` | [`08-validators.md`](./08-validators.md) |
| `converters`, `typesconv` | [`09-converters.md`](./09-converters.md) |
| `coredynamic`, `reflectcore` | [`10-reflection-and-dynamic.md`](./10-reflection-and-dynamic.md) |
| `coreversion`, `versionindexes` | [`11-versioning.md`](./11-versioning.md) |
| `cmd/*` binaries | [`12-cmd-entrypoints.md`](./12-cmd-entrypoints.md) |
| `coretests/*`, `tests/*` | [`13-testing-patterns.md`](./13-testing-patterns.md) + [`14-tests-folder-walkthrough.md`](./14-tests-folder-walkthrough.md) |
