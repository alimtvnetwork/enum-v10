# 05 — Code-vs-Spec Cycle 4: Data Structures

**Date:** 2026-05-04
**Spec audited:** [`spec/01-app/06-data-structures.md`](../01-app/06-data-structures.md)
**Method:** Manual claim extraction; automated `grep` verification across all 83 `enum-v9` files that import `coredata/*`.
**Auditor:** Lovable agent (evidence-driven)

---

## Summary

| Metric | Value |
|---|---|
| Claims extracted | **20** |
| ✅ Match | **5** |
| ⚠️ Drift | **6** |
| ❌ Contradiction | **3** |
| ❓ Unverifiable | **6** |
| **Match rate (verifiable)** | **5 / 14 = 35.7 %** |

> **Verdict:** §06 documents the **upstream `core-v9` API surface** as if `enum-v9` consumed all of it. Reality from a `enum-v9`-side scan: only **3 of 5 sub-packages** have any consumers in this repo (`corejson` 80 files, `corestr` 4 files, `coreonce` 1 file); **`coregeneric` and `corepayload` have zero consumers**. Worse, several documented API entry points (`corejson.Serialize.ToString`, `Serialize.Raw`, `Deserialize.UsingBytes`, `Deserialize.FromTo`, `corestr.NewCollectionPtrUsingStrings`, `coreonce.New.String`) are **never called** — actual call sites use different names (`corejson.Serialize.ToBytesErr`, `Deserialize.BytesTo`, `coreonce.NewByteOnce`/`NewAnyOnce`). Two files import `encoding/json` directly, contradicting the §4 "never" rule.

---

## Claim-by-claim verification

### ✅ Match (5)

| # | Claim | Evidence |
|---|-------|----------|
| C1 | `coredata/` is a parent directory grouping data-structure sub-packages | Confirmed by import paths (`coredata/corejson`, `coredata/corestr`, `coredata/coreonce`) in the 83 consumer files |
| C2 | `corejson` is the canonical JSON entry point with a fluent `New` / `NewPtr` / `Serialize` / `Deserialize` surface | `rg corejson\.` finds `corejson.New`, `corejson.NewPtr`, `corejson.NewResult.Ptr`, `corejson.Serialize.ToBytesErr`, `corejson.Deserialize.BytesTo` across the codebase |
| C3 | `corejson.NewPtr(myStruct).PrettyJsonString()` is the pretty-print idiom | Confirmed at `cmd/main/main.go` (`osDetail.PrettyJsonString()`, `json.PrettyJsonString()`) and other consumers |
| C4 | `corejson` exposes contract interfaces (`Jsoner`, `JsonMarshaller`, `JsonContractsBinder`) | All three appear as types in consumer signatures (e.g. `inttype/all-constructors.go: NewUsingJsoner(jsoner corejson.Jsoner)`) |
| C5 | `corestr` provides `Hashset`, `SimpleSlice`, plus a `SimpleStringOnce` lazy value | `corestr.New.Hashset`, `corestr.New.SimpleSlice`, `corestr.SimpleSlice`, `corestr.SimpleStringOnce` all confirmed in `enum-v9` consumers |

### ⚠️ Drift — code is fine, spec prose is stale (6)

| # | Claim | Code reality | Severity | Fix path |
|---|-------|--------------|----------|----------|
| D-CVS-20 | "Serialize: `corejson.Serialize.ToString(myStruct)` / `Serialize.Raw(myStruct)`" | Zero call sites for either method. Real serializer used: **`corejson.Serialize.ToBytesErr(...)`** (returns `[]byte, error`). | MED | Spec — replace `ToString` / `Raw` examples with `ToBytesErr`; document the actual return type |
| D-CVS-21 | "Deserialize: `corejson.Deserialize.UsingBytes(jsonBytes, &target)` / `Deserialize.FromTo(source, &target)`" | Zero call sites for either. Real deserializer used: **`corejson.Deserialize.BytesTo(...)`**. | MED | Spec — replace examples with `Deserialize.BytesTo`; remove `FromTo` claim or add evidence link |
| D-CVS-22 | `coreonce.New.String(producer)` is the construction pattern | Zero call sites. Real API: **`coreonce.NewAnyOnce(...)`** and **`coreonce.NewByteOnce(...)`** (top-level functions, not a `New.<Type>` namespace). | MED | Spec — rewrite §5 to show the actual top-level constructors; remove the `New.<Type>` namespace claim or document it as upstream-only |
| D-CVS-23 | `corestr` is "a thread-safe list of strings" + `corestr.NewCollectionPtrUsingStrings(&values, 0)` | Zero call sites for `NewCollectionPtrUsingStrings`. Real `corestr` use is **`Hashset`**, **`SimpleSlice`**, and **`SimpleStringOnce`** — *not* a string-list collection. | MED | Spec — rewrite §3 around the actual exported surface (`Hashset`, `SimpleSlice`, `SimpleStringOnce`); drop the `Collection`-style example |
| D-CVS-24 | `coreonce` covers "all common types" via a uniform `New.<Type>` namespace | Only `AnyOnce` and `ByteOnce` constructors appear in consumer call sites. | LOW | Spec — soften "all common types" to "common typed wrappers (`AnyOnce`, `ByteOnce`, …)" or fetch upstream listing (task **AB**) |
| D-CVS-25 | `coregeneric` and `corepayload` are part of the `enum-v9` consumer story | **Zero** consumers in `enum-v9` for either. Spec presents them as first-class citizens; from this repo they are aspirational. | LOW | Spec — add a note in §1 that `coregeneric` and `corepayload` are documented for upstream completeness; `enum-v9` does not consume them |

### ❌ Contradiction — spec rule is actively violated by code (3)

| # | Claim | Reality | Severity | Resolution |
|---|-------|---------|----------|------------|
| C-CVS-06 | §4 rule: "Any package that touches JSON should import `corejson` and **never** `encoding/json` directly" | Two files import `encoding/json`: **`inttype/Variant.go`** (calls `json.Marshal(it.Value())` in `MarshalJSON`) and **`inttype/all-constructors.go`** (uses `*json.Number` as a parameter type in `NewUsingJsonNumber`). | HIGH | **Either** (a) document the legitimate exceptions — `MarshalJSON` may delegate to `json.Marshal` for primitive value emission; `*json.Number` is the only stdlib type that meaningfully crosses the API boundary — **or** (b) refactor `inttype` to route through `corejson`. (a) matches author intent and is recommended. |
| C-CVS-07 | "Constructors return non-nil values" + spec example `corejson.Serialize.ToString(myStruct)` returns `(string, error)` | The actually-used `corejson.Serialize.ToBytesErr` returns `(*Result, error)` (a non-nil result wrapper) — but the specced `ToString` / `Raw` signatures *do not exist* and so cannot satisfy any rule. The spec example would not compile. | HIGH | Replace the §4 code block with a runnable snippet using `Serialize.ToBytesErr` / `Deserialize.BytesTo` |
| C-CVS-08 | §6 example: `corepayload.New.PayloadWrapper.UsingInstruction(&corepayload.PayloadCreateInstruction{...})` | Zero usages of `corepayload` in `enum-v9`. The example cannot be verified against any consumer call site in this repo. (Spec also forces the reader to instantiate a `PayloadCreateInstruction` literal whose fields cannot be cross-checked here.) | MED | Either move §6 to an "upstream-only" appendix, or fetch upstream `core-v9` source (task **AB**) and verify the literal field set |

### ❓ Unverifiable (6)

| # | Claim | Why unverifiable |
|---|-------|------------------|
| Q1 | Full `coregeneric.Collection[T]` mutation/query/iterator method tables | Zero `coregeneric` consumers in `enum-v9`. Need upstream `core-v9` source (task **AB**). |
| Q2 | `coregeneric.Hashset` / `Hashmap` / `Pair` / `Triple` / `LinkedList` API surfaces | Same — no `enum-v9` consumers. |
| Q3 | "Constructors return non-nil values" as a global guarantee across all 5 sub-packages | Verifiable for the 3 sub-packages with consumers (`corejson`, `corestr`, `coreonce`) — no `nil`-result patterns observed — but cannot be confirmed for `coregeneric` / `corepayload`. |
| Q4 | "All containers expose a `*Lock` family for thread-safe access" | No `*Lock` methods are called from `enum-v9` consumers. Cannot confirm the family exists or is uniform. |
| Q5 | "Mutation methods return the receiver (fluent style)" | No mutation chains observed in consumer code. |
| Q6 | `corejson` "wraps stdlib errors in `errcore.UnMarshallingFailedType` / `FailedToConvertType`" | Error types are produced by the upstream package; consumers receive but do not introspect them. |

---

## Score

**Verifiable subset:** 14 claims (5 ✅ + 6 ⚠️ + 3 ❌) → **5/14 = 35.7 %**

The 3 ❌ contradictions matter most: one is a real code/spec mismatch (`encoding/json` exception in `inttype`), one is an unrunnable code example (`Serialize.ToString` doesn't exist), and one is an unverifiable example (`corepayload` literal).

---

## Recommended remediation order (cheap → expensive)

1. **D-CVS-20 + D-CVS-21 + C-CVS-07** (HIGH/MED): Replace §4 code block with the actually-used `corejson.Serialize.ToBytesErr` + `corejson.Deserialize.BytesTo` examples. Add a note that the upstream `Result` wrapper provides `PrettyJsonString()`. ~10 min. Bumps score 35.7 → 57.1.
2. **D-CVS-22 + D-CVS-24** (MED+LOW): Rewrite §5 around the real top-level constructors (`coreonce.NewAnyOnce`, `coreonce.NewByteOnce`); soften "all common types" claim. ~5 min. → 71.4.
3. **D-CVS-23** (MED): Rewrite §3 around `corestr.{Hashset,SimpleSlice,SimpleStringOnce}`; drop the unused `NewCollectionPtrUsingStrings` example. ~10 min. → 78.6.
4. **C-CVS-06** (HIGH): Document the `inttype` exception to the "never `encoding/json` directly" rule (`MarshalJSON` delegation + `*json.Number` parameter type). ~5 min. → 92.9.
5. **C-CVS-08 + D-CVS-25** (MED+LOW): Add an "upstream-only" callout to §1 noting that `coregeneric` and `corepayload` have no `enum-v9` consumers; defer full verification to task **AB**. ~5 min. → 100.0 verifiable.

Total remediation budget: **~35 minutes** to lift §06 from 35.7 % → 100 % verifiable. Six ❓ remain pending task **AB** (upstream `core-v9` source fetch).

---

## See also

- [`02-cycle1-import-conventions.md`](./02-cycle1-import-conventions.md) — Cycle 1
- [`03-cycle2-error-system.md`](./03-cycle2-error-system.md) — Cycle 2
- [`04-cycle3-enum-system.md`](./04-cycle3-enum-system.md) — Cycle 3
- [`01-scoreboard.md`](./01-scoreboard.md) — running totals
- [`../01-app/06-data-structures.md`](../01-app/06-data-structures.md) — spec under audit
