# 02 — Design Philosophy

> ✅ **Done** — extracted from `spec/00-llm-integration-guide.md` §Design Philosophy.
> **Status**: filled in audit Step 3 (2026-04-23, Asia/Kuala_Lumpur).
> **Audience**: AI agents and human contributors writing **new code** for `core-v9`. Read this *before* opening a PR.

This file expands the seven design pillars from [`00-repo-overview.md` §3](./00-repo-overview.md#3-design-pillars-tldr) with **rules**, **rationale**, and **examples**.

---

## Pillar 1 — One File per Function

**Rule.** Each public function lives in its own `.go` file, named after the function. Files stay in the **50–200 line** sweet spot. When a function naturally exceeds 200 lines, split helpers into sibling files with the same prefix.

**Rationale.**
- Diffs are scoped to the function under review.
- File-name search instantly locates a function (`ripgrep -l SplitByDelimiter` → one hit).
- Symbol moves between packages are file moves, preserving git history.

**Example.**

```
corestr/
├── Split.go               // func Split(...)
├── SplitByDelimiter.go    // func SplitByDelimiter(...)
├── SplitByDelimiterMany.go
└── SplitByRegex.go
```

**Filename casing rule.** The filename **mirrors the exported identifier exactly**:

- Exported function `SplitByDelimiter` → `SplitByDelimiter.go` (PascalCase)
- Exported method `IsSuccess` on a struct-as-namespace → `IsSuccess.go` (PascalCase)
- Internal-only helper `newCreator` → `newCreator.go` (camelCase, matches identifier)

Do **not** use `snake_case` (`split_by_delimiter.go`) or kebab-case. The rule is mechanical: copy the Go identifier verbatim, append `.go`. This keeps `ripgrep <FuncName>` aligned with file search.

**Exception.** Trivial constructors (`New`, `NewWith`) and tiny method sets on a value type may share one file when each function is < 20 lines.

---

## Pillar 2 — Struct-as-Namespace

**Rule.** Group related operations on an **unexported struct type**, then expose a **package-level `var`** of that struct. Methods become the public API surface.

**Rationale.**
- IDE auto-complete after `pkg.Verb.` shows the entire verb's vocabulary.
- Avoids global function-name collisions inside one package.
- Method receivers can carry shared state (e.g. cached configuration) without exposing it.

**Example.**

```go
// corejson/serialize.go
type serialize struct{}

var Serialize serialize

func (serialize) ToString(v any) string { /* ... */ }
func (serialize) ToBytes(v any) []byte  { /* ... */ }
func (serialize) Pretty(v any) string   { /* ... */ }
```

Call site:

```go
s := corejson.Serialize.ToString(myValue)
b := corejson.Serialize.ToBytes(myValue)
```

**Companion pattern**: `New` factory variable (Pillar 7).

---

## Pillar 3 — Interface-First

**Rule.** Define contracts in `coreinterface/` (or a `*inf` sub-package) using Go's `-er` suffix. Production code depends on interfaces, not concrete types.

**Examples of canonical names.**

| Interface | Method signature |
|---|---|
| `NameGetter` | `Name() string` |
| `Serializer` | `Serialize() ([]byte, error)` |
| `BasicEnumer` | `Value() byte; Name() string` |
| `ErrorWrapper` | `Wrap(err error) error` |

**Rationale.**
- Test doubles slot in without `gomock`.
- Cyclic-import risk drops because most packages depend on `coreinterface/*` only.
- Cross-package contracts are discoverable in one place.

---

## Pillar 4 — Zero-Nil Safety

**Rule.** Return **empty** containers, never `nil`. Pointer-receiver methods include a nil-receiver guard at the top.

**Why.** Callers can iterate or `len()` without a nil check, eliminating an entire class of panics.

**Example.**

```go
// Good
func (c *Collection) Items() []string {
    if c == nil {
        return []string{}
    }
    return c.items
}

// Bad
func (c *Collection) Items() []string {
    return c.items // panics when c is nil
}
```

**Test contract.** Every pointer-receiver method has a `NilReceiver_test.go` file in its `*tests` package. See [`/spec/06-testing-guidelines/02-test-case-types.md`](../06-testing-guidelines/02-test-case-types.md) for the `CaseNilSafe` pattern.

---

## Pillar 5 — Generics Where Clear

**Rule.** Add generic versions **alongside** backward-compatible typed wrappers. Do not retro-fit generics if it forces every caller to add type parameters.

**Example.**

```go
// Generic
func MapValues[K comparable, V any](m map[K]V) []V { /* ... */ }

// Typed wrapper preserved for callers that don't want to spell the type
func StringMapValues(m map[string]string) []string {
    return MapValues(m)
}
```

**When to skip generics.** Reflection-heavy code (`coredynamic`, `reflectcore`) — runtime types defeat compile-time parameters.

---

## Pillar 6 — Value Receivers by Default

**Rule.** Read-only methods use **value receivers**. Switch to pointer receivers only when:

1. The method **mutates** the receiver, OR
2. The receiver struct is **large** (rule of thumb: > 64 bytes), OR
3. **Interface satisfaction** requires pointer (e.g. `io.Writer`).

**Why.** Value receivers compose better, are inherently nil-safe (you can't call them on a nil pointer because the value is copied), and avoid accidental sharing across goroutines.

---

## Pillar 7 — `newCreator` / `New` Factory Variable

**Rule.** Constructor families live on a `New` package-level variable of an unexported `newCreator` struct. Clients write `pkg.New.X.Y(...)` instead of `pkg.NewXY(...)`.

**Example.**

```go
// enumimpl/new.go
type newCreator struct {
    BasicByte basicByteCreator
    BaseByte  baseByteCreator
}

var New = newCreator{}

// Call site
e := enumimpl.New.BasicByte.UsingTypeSlice(myTypes)
```

**Why.** Constructor count grows fast; the namespace structure scales without polluting the package's top-level identifiers.

---

## Cross-Cutting Conventions

These follow from the pillars above but are worth calling out:

1. **Errors are values, never panics.** Use `errcore` to construct rich errors. See [`04-error-system.md`](./04-error-system.md).
2. **No package-level mutable state.** `var` declarations must be either constants-in-disguise (the `New` factory, the `Serialize` namespace) or backed by `sync.Once` / `mutexbykey`.
3. **Test packages mirror production packages** *(upstream `core-v9`)*. A package `foo/` has tests at `tests/creationtests/footests/`. See [`/spec/06-testing-guidelines/01-folder-structure.md`](../06-testing-guidelines/01-folder-structure.md). In `enum-v9` (this module) all enum tests instead live in the shared `tests/creationtests/` registry — see [`13-testing-patterns.md` §6.1](./13-testing-patterns.md#61-enum-v9-specific-layout).
4. **Internal packages stay internal.** Code under `internal/` can be refactored without bumping a major version.

---

## Anti-Patterns (do **not** do)

| Anti-pattern | Why it's banned | Correct alternative |
|---|---|---|
| Returning `nil` slice from a getter | Forces nil checks at every call site | Return `[]T{}` |
| Free function `pkg.DoX()` for related verbs | Pollutes top-level; no IDE grouping | `pkg.Verb.DoX()` (Pillar 2) |
| `pkg.NewFooBarBaz(...)` constructors | Constructor explosion | `pkg.New.Foo.BarBaz(...)` (Pillar 7) |
| Pointer receiver on small read-only struct | Loses value-semantics safety | Value receiver (Pillar 6) |
| Two functions per file | Defeats Pillar 1 | Split into two files |
| Importing `internal/` cross-module | Compiler-blocked anyway | Add a public package or interface |

---

## Where to Go Next

- [`03-import-conventions.md`](./03-import-conventions.md) — exact import paths and the `internal/` rule
- [`13-testing-patterns.md`](./13-testing-patterns.md) — how the philosophy is enforced via tests
- [`/spec/06-testing-guidelines/08-good-vs-bad.md`](../06-testing-guidelines/08-good-vs-bad.md) — concrete good-vs-bad examples
