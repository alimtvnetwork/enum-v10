# 06 — Data Structures

> ✅ **Status**: filled in audit Step 5b (2026-04-23, Asia/Kuala_Lumpur).
> **Audience**: anyone using or extending `coredata/*` packages.
> **Source**: extracted from [`/spec/00-llm-integration-guide.md`](../00-llm-integration-guide.md) §coredata + §coregeneric.

---

## Table of Contents

1. [`coredata` Umbrella](#1-coredata-umbrella)
2. [`coregeneric` — Generic Containers](#2-coregeneric--generic-containers)
3. [`corestr` — String Collections](#3-corestr--string-collections)
4. [`corejson` — JSON Pipeline](#4-corejson--json-pipeline)
5. [`coreonce` — Compute-Once Values](#5-coreonce--compute-once-values)
6. [`corepayload` — `PayloadWrapper`](#6-corepayload--payloadwrapper)
7. [Choosing the Right Container](#7-choosing-the-right-container)

---

## 1. `coredata` Umbrella

`coredata/` is a parent directory, not a package. It groups data-structure sub-packages that share the same design rules:

| Sub-package | Purpose |
|---|---|
| [`coregeneric`](#2-coregeneric--generic-containers) | Generic `Collection[T]`, `Hashset[T]`, `Hashmap[K,V]`, `SimpleSlice[T]`, `LinkedList[T]`, `Pair`, `Triple` |
| [`corestr`](#3-corestr--string-collections) | String-typed collections with thread-safe operations |
| [`corejson`](#4-corejson--json-pipeline) | JSON serialize / deserialize / pretty-print |
| [`coreonce`](#5-coreonce--compute-once-values) | Lazy compute-once cached values |
| [`corepayload`](#6-corepayload--payloadwrapper) | Wire-format payload envelopes |

> ℹ️ **Consumer-coverage note** *(audit Cycle 4)*: from `enum-v9`, only `corejson` (80 files), `corestr` (4 files), and `coreonce` (1 file) are actually imported. `coregeneric` and `corepayload` have **zero `enum-v9` consumers** — they are documented here for upstream `core-v9` completeness, but the API surfaces in §2 and §6 cannot be verified against this repo. Treat those two sub-sections as *upstream-reference* until task **AB** fetches `core-v9` source.

### Shared rules

1. **Constructors return non-nil values.** An "empty" collection is still a usable receiver — see Zero-Nil Safety in [`02-design-philosophy.md`](./02-design-philosophy.md).
2. **All containers expose a `*Lock` family for thread-safe access.** Pick the lockless variant when you own the only goroutine touching the value.
3. **Mutation methods return the receiver** (fluent style) so calls can chain.
4. **Read methods never mutate, never lock unless `*Lock`.**

---

## 2. `coregeneric` — Generic Containers

Located at `coredata/coregeneric/`. Generic-first, built on Go 1.18+ type parameters.

> ⚠️ **Upstream-only sub-package** *(audit Cycle 4)*: `coregeneric` has **zero consumers in `enum-v9`**. The API tables below reflect upstream `core-v9` documentation but cannot be cross-checked against any call site in this repo. Treat as upstream-reference until task **AB** verifies against `core-v9` source.

### 2.1 `Collection[T any]`

Slice-backed collection with embedded `sync.Mutex`.

#### Construction

```go
import "github.com/alimtvnetwork/core-v9/coredata/coregeneric"

// Via the New struct-as-namespace
col := coregeneric.New.Collection.String.Cap(10)
col := coregeneric.New.Collection.Int.Items(1, 2, 3)
col := coregeneric.New.Collection.Float64.Empty()

// Via package-level functions
col := coregeneric.EmptyCollection[string]()
col := coregeneric.NewCollection[MyStruct](20)
col := coregeneric.CollectionFrom(existingSlice)   // no copy
col := coregeneric.CollectionClone(existingSlice)  // copies
```

#### Mutation (returns `*Collection[T]`)

| Method | Description |
|---|---|
| `Add(item)` / `AddLock(item)` | Append one |
| `Adds(items...)` / `AddsLock(items...)` | Append variadic |
| `AddSlice([]T)` | Append from slice |
| `AddIf(bool, item)` / `AddIfMany(bool, items...)` | Conditional append |
| `AddFunc(func() T)` | Append function result |
| `AddCollection(*Collection[T])` / `AddCollections(...*Collection[T])` | Merge |
| `RemoveAt(index) bool` | Remove by index |
| `SortFunc(less)` | In-place sort |
| `Reverse()` | In-place reverse |
| `ConcatNew(items...)` | New collection = this + items |
| `Clone()` | Deep copy |

#### Query

| Method | Returns | Description |
|---|---|---|
| `Length()` / `Count()` / `LengthLock()` | `int` | Size |
| `IsEmpty()` / `IsEmptyLock()` | `bool` | Empty check |
| `HasItems()` / `HasAnyItem()` | `bool` | Non-empty |
| `HasIndex(i)` | `bool` | Bounds check |
| `Items()` | `[]T` | Underlying slice |
| `First()` / `Last()` | `T` | **Panics** if empty |
| `FirstOrDefault()` / `LastOrDefault()` | `T` | Zero-value if empty |
| `SafeAt(i)` | `T` | Zero-value if OOB |
| `Skip(n)` / `Take(n)` | `[]T` | Slice ops |
| `Filter(pred)` | `*Collection[T]` | New filtered |
| `CountFunc(pred)` | `int` | Count matching |
| `ForEach(fn)` / `ForEachBreak(fn)` | — | Iterate |

#### Iterators (Go 1.23+)

```go
for i, item := range col.All()    { ... }   // iter.Seq2[int, T]
for item    := range col.Values() { ... }   // iter.Seq[T]
```

#### Package-level generic helpers

| Function | Constraint | Description |
|---|---|---|
| `MapCollection(src, fn)` | `T→U` | Transform `Collection[T]` → `Collection[U]` |
| `FlatMapCollection(src, fn)` | `T→[]U` | Flatten-transform |
| `ReduceCollection(src, init, fn)` | `T→U` | Fold |
| `GroupByCollection(src, keyFn)` | `K comparable` | `map[K]*Collection[T]` |
| `ContainsFunc` / `ContainsItem` | predicate / `comparable` | Search |
| `IndexOfFunc` / `IndexOfItem` | — | Find index |
| `Distinct(src)` | `comparable` | Deduplicate |
| `ContainsAll` / `ContainsAny` | `comparable` | Bulk membership |
| `RemoveItem` / `RemoveAllItems` | `comparable` | Remove |
| `ToHashset(src)` | `comparable` | Convert |

### 2.2 `Hashset[T comparable]`

Set backed by `map[T]bool` with embedded `sync.Mutex`.

```go
hs := coregeneric.New.Hashset.String.Empty()
hs := coregeneric.New.Hashset.Int.Cap(100)
hs := coregeneric.HashsetFrom([]string{"a", "b"})
hs := coregeneric.HashsetFromMap(existingMap)

hs.Add("x")
hs.Contains("x")
hs.RemoveItem("x")
hs.ToSlice()      // []T
hs.Length()
```

### 2.3 `Hashmap[K comparable, V any]`

Map wrapper with the same locking convention as `Collection`.

```go
hm := coregeneric.New.Hashmap.StringInt.Empty()
hm.Set("a", 1)
v, ok := hm.Get("a")
hm.Delete("a")
hm.Length()
hm.Keys()    // []K
hm.Values()  // []V
```

### 2.4 `SimpleSlice[T any]`

Lighter than `Collection[T]` — no mutex, no lock-suffixed methods. Use when you own the slice exclusively.

```go
ss := coregeneric.New.SimpleSlice.String.Cap(8)
ss.Append("a")
ss.AppendFmt("%d : %s", i, name)
ss.Strings()   // []string (for string-typed only)
```

> Used heavily inside Style B test runners — see [`13-testing-patterns.md` §3](./13-testing-patterns.md#3-style-b---basetestcase--testwrapper-typed-slice--hetero-input).

### 2.5 `LinkedList[T any]`, `Pair[L, R]`, `Triple[A, B, C]`

Lower-volume containers; same construction pattern via `coregeneric.New.LinkedList`, `coregeneric.New.Pair`, `coregeneric.New.Triple`. Use `Pair` / `Triple` when you need a typed multi-return without declaring a struct.

---

## 3. `corestr` — String Collections

Located at `coredata/corestr/`. The exported surface actually consumed from `enum-v9` is three string-typed helpers (audit Cycle 4):

| Type / constructor | Purpose | Example call site |
|---|---|---|
| `corestr.New.Hashset` | String set with O(1) membership | `New.Hashset.Empty()` / `New.Hashset.UsingItems("a", "b")` |
| `corestr.New.SimpleSlice` | Lightweight string slice (no mutex) | `New.SimpleSlice.Cap(8)` |
| `corestr.SimpleSlice` | Same type used directly as a parameter | `func(ss *corestr.SimpleSlice)` signatures |
| `corestr.SimpleStringOnce` | Lazy compute-once string value | package-level cached labels |

```go
import "github.com/alimtvnetwork/core-v9/coredata/corestr"

// String set
seen := corestr.New.Hashset.Empty()
seen.Add("alpha")
exists := seen.Contains("alpha")

// Lightweight string slice (no mutex)
lines := corestr.New.SimpleSlice.Cap(8)
lines.Append("first line")

// Lazy-cached string
label := corestr.SimpleStringOnce{Producer: func() string { return expensiveLabel() }}
fmt.Println(label.Value())
```

> **Note (audit Cycle 4):** earlier drafts of this spec showed `corestr.NewCollectionPtrUsingStrings(&values, 0)` as the canonical entry point. That constructor is **not used anywhere in `enum-v9`**, and a thread-safe "string list collection" is not part of the consumer-side surface today. If you need a string list with mutex-protected mutation, use `coregeneric.New.Collection.String` (upstream-only — see §1 callout) or — for the lock-free case — `corestr.New.SimpleSlice`.

---

## 4. `corejson` — JSON Pipeline

Located at `coredata/corejson/`. The single canonical entry point for JSON in the framework.

```go
import "github.com/alimtvnetwork/core-v9/coredata/corejson"

// Serialize → returns *corejson.Result (a non-nil wrapper carrying the bytes + error)
result, err := corejson.Serialize.ToBytesErr(myStruct)
if err != nil { return err }
jsonBytes := result.Bytes()

// Deserialize from bytes into a target value
err = corejson.Deserialize.BytesTo(jsonBytes, &target)

// Wrap an existing value for downstream JSON contracts
js := corejson.NewPtr(myStruct)        // *Result
pretty := js.PrettyJsonString()        // pretty-printed string

// Implementing the JSON contracts in your own type
type MyType struct{ ... }
func (it MyType) JsonPtr() *corejson.Result { return corejson.New(it) }
// satisfies corejson.Jsoner / JsonMarshaller / JsonContractsBinder
```

### Why `corejson` and not `encoding/json`?

1. **Consistent error categories** — wraps stdlib errors in `errcore.UnMarshallingFailedType` / `FailedToConvertType` so log scanners can attribute failures.
2. **Result-wrapper ergonomics** — `*Result` carries the bytes, the original error, and helpers like `PrettyJsonString()` so downstream code doesn't re-marshal.
3. **Contract interfaces** — `corejson.Jsoner`, `JsonMarshaller`, `JsonContractsBinder` let consumers accept "anything JSON-able" without committing to a concrete type.

### Rule (with documented exceptions)

Any package that touches JSON **should** import `corejson` rather than `encoding/json` directly. Two **legitimate exceptions** apply (audit Cycle 4 — these are the only direct `encoding/json` imports in `enum-v9`):

| File | Use of `encoding/json` | Why it's allowed |
|---|---|---|
| `inttype/Variant.go` | `json.Marshal(it.Value())` inside `MarshalJSON` | The receiver IS the JSON marshaller; delegating its primitive value to `json.Marshal` is the canonical Go pattern and `corejson` would call the same code with extra wrapping. |
| `inttype/all-constructors.go` | `*json.Number` as a parameter type in `NewUsingJsonNumber(jsonNumber *json.Number)` | `json.Number` is a stdlib value type with no `corejson` equivalent; consumer code that already holds one should be able to pass it through. |

Outside of these two patterns (a `MarshalJSON` body emitting a primitive, or a `*json.Number` parameter), prefer `corejson`.

---

## 5. `coreonce` — Compute-Once Values

Located at `coredata/coreonce/`. Lazy-evaluated cached values built on `sync.Once`. The exported surface used from `enum-v9` is a small set of typed top-level constructors (audit Cycle 4):

| Constructor | Cached value type | Use case |
|---|---|---|
| `coreonce.NewAnyOnce(producer func() any)` | `any` | Heterogeneous cached value (struct, pointer, etc.) |
| `coreonce.NewByteOnce(producer func() byte)` | `byte` | Numeric scalar cache |

```go
import "github.com/alimtvnetwork/core-v9/coredata/coreonce"

// Generic any-typed cache
cfg := coreonce.NewAnyOnce(func() any { return loadConfig() })
v := cfg.Value()           // first call runs the producer
v = cfg.Value()            // later calls return cached

// Byte-typed cache
firstByte := coreonce.NewByteOnce(func() byte { return readMagicByte() })
```

Use for expensive package-level computations (regex compilation, file reads, network handshakes) that should run at most once.

> **Note (audit Cycle 4):** earlier drafts described a uniform `coreonce.New.<Type>(...)` namespace covering "all common types". The actual consumer-side surface is the two top-level functions above. Additional typed wrappers may exist upstream in `core-v9` (e.g. `corestr.SimpleStringOnce` is the string equivalent and lives in `corestr` rather than `coreonce`); pending task **AB** to enumerate the full upstream listing.

---

## 6. `corepayload` — `PayloadWrapper`

Located at `coredata/corepayload/`. Wire-format envelope for payload-bearing messages.

> ⚠️ **Upstream-only sub-package** *(audit Cycle 4)*: `corepayload` has **zero consumers in `enum-v9`**. The example below reflects the documented upstream API but cannot be cross-checked against any call site in this repo. Treat it as upstream-reference; the exact field names of `PayloadCreateInstruction` are pending verification under task **AB**.

```go
import "github.com/alimtvnetwork/core-v9/coredata/corepayload"

// Empty placeholder
payload := corepayload.New.PayloadWrapper.Empty()

// From an instruction (field set per upstream docs — verify under task AB)
payload = corepayload.New.PayloadWrapper.UsingInstruction(&corepayload.PayloadCreateInstruction{
    Name:       "user-create",
    Identifier: "usr-123",
    Payloads:   myStruct,
})
```

### When to use

- Inter-service messages where the payload type is dynamic.
- Storing heterogeneous data in a single field that a downstream consumer will type-assert.
- Building command/event buses with structured metadata (`Name`, `Identifier`).

### When **not** to use

- Plain in-process function returns — use a typed struct.
- Cases where the consumer knows the concrete payload type — use that type directly.

---

## 7. Choosing the Right Container

Decision matrix. ⚠️ marks rows whose target sub-package has no `enum-v9` consumers (see §1 callout) — usable upstream but cannot be verified from this repo.

| Need | Use | Verified in `enum-v9`? |
|---|---|---|
| Generic, mutex-protected, slice-backed | `coregeneric.New.Collection.<Type>` | ⚠️ upstream-only |
| Generic, mutex-protected, set semantics | `coregeneric.New.Hashset.<Type>` | ⚠️ upstream-only |
| Generic, mutex-protected, key→value | `coregeneric.New.Hashmap.<KV>` | ⚠️ upstream-only |
| Lightweight generic slice, no concurrency | `coregeneric.New.SimpleSlice.<Type>` | ⚠️ upstream-only |
| String set | `corestr.New.Hashset` | ✅ |
| Lightweight string slice, no concurrency | `corestr.New.SimpleSlice` | ✅ |
| Lazy compute-once string | `corestr.SimpleStringOnce` | ✅ |
| Two- or three-tuple return | `coregeneric.New.Pair` / `Triple` | ⚠️ upstream-only |
| Linked-list semantics (cheap insert) | `coregeneric.New.LinkedList` | ⚠️ upstream-only |
| Cached lazy value (any) | `coreonce.NewAnyOnce(producer)` | ✅ |
| Cached lazy value (byte) | `coreonce.NewByteOnce(producer)` | ✅ |
| Wire-format envelope | `corepayload.New.PayloadWrapper` | ⚠️ upstream-only |
| JSON serialize | `corejson.Serialize.ToBytesErr(value)` → `*Result` | ✅ |
| JSON deserialize | `corejson.Deserialize.BytesTo(bytes, &target)` | ✅ |
| JSON wrap-and-pretty-print | `corejson.NewPtr(value).PrettyJsonString()` | ✅ |

> Generic-first rule: when both a generic and a string-typed option exist, prefer the generic version unless a string-specific helper is essential.

---

## See Also

- [`02-design-philosophy.md`](./02-design-philosophy.md) — `newCreator` pattern (§5), zero-nil safety (§4)
- [`05-enum-system.md`](./05-enum-system.md) — `enumimpl` builds on the same `New` namespace pattern
- [`07-conditional-and-utilities.md`](./07-conditional-and-utilities.md) — `regexnew` and `coreonce` pair well for cached compiled regex
- [`13-testing-patterns.md`](./13-testing-patterns.md) — `SimpleSlice` is heavily used in Style B test runners
- [`/spec/00-llm-integration-guide.md` §coregeneric](../00-llm-integration-guide.md#coregeneric--generic-data-structures-api-reference) — Full API tables
