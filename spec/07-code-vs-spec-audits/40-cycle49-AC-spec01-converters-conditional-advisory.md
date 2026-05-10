# Cycle 49 — AC dimension: spec-internal-consistency probe of §09 + §07 advisory carry-overs

**Date:** 2026-05-10
**Dimension:** **AC** — spec-internal-consistency probe of items previously deferred from AB-residual passes as "advisory / behavioural / unprobeable" (Cycles 19, 20, 43, 44).
**Scope:**
- `spec/01-app/09-converters.md` — rows 64, 99-107, 119, 130, 161, 172 (deferred by Cycles 19 + 43 to "AC contract pass").
- `spec/01-app/07-conditional-and-utilities.md` — row 142 (deferred by Cycles 20 + 44 as "advisory / behavioural").
**Result:** **3 ❓ → ✅** (with 1 mechanism-name note) + **1 ❓ → ❌ NEW C-CVS-65 (HIGH)** + **4 ❓ retained as genuine contract-pass deferrals** (rows 99-107, 119, 130, 172 — require per-method body inspection of every converter, out of scope for AC).
**Allowed under freeze:** read-only audit promotion (no `spec/01-app/` rewrites).
**Evidence root:** `/tmp/core-v9-upstream` @ `v1.5.8` (re-cloned 2026-05-10).

---

## 1. Promotions

### 1.1 §09 row 64 — "`converters.PrettyJson` namespace duplicates a subset of `corejson`" → ✅ (with mechanism note)

- `converters/vars.go:35` — `PrettyJson = jsoninternal.Pretty` (package-level binding).
- `converters/anyItemConverter.go:331-336` — `ToPrettyJson(...)` method on `anyItemConverter` (delegates to the same `jsoninternal.Pretty` surface).
- `coredata/corejson/PrettyJsonStringer.go:25` — `type PrettyJsonStringer interface { ... }` defines the JSON-pretty-printing contract that `converters.PrettyJson` implements/exposes.
- **Spec claim:** "PrettyJson namespace duplicates a subset of corejson."
- **Evidence verdict:** structurally accurate — `converters.PrettyJson` exposes the same pretty-printing surface that `corejson.PrettyJsonStringer` declares. Calling it a "duplicate" overstates: it is a **delegating re-export** via `jsoninternal.Pretty`. Both routes funnel through the same backing implementation, so the spec's intent ("the same capability is reachable from two top-level packages") is correct.
- **Mechanism note (informational, not a finding):** suggested wording tweak in a future spec pass — replace "duplicates a subset of" with "delegates to a subset of" to better reflect the `jsoninternal.Pretty` indirection. Filed as **S-002** (LOW) in the suggestions tracker; non-blocking.

### 1.2 §07 row 142 — `issetter.Value` "Pitfall: not a drop-in for `bool`" → ✅ (self-consistent advisory)

- Already promoted in Cycle 20 row 10: `issetter/Value.go:51-58` defines the 6-state byte enum (`Uninitialized=0, True=1, False=2, Default=3, Conditional=4, Wildcard=5`).
- The advisory text — "not a drop-in for `bool`" — is **structurally entailed** by the byte-enum shape: a 6-state enum is provably not assignment-compatible with the 2-state Go `bool` type, and the package exposes `IsOn()`/`IsOff()` accessor methods (`Value.go:148, 152`) precisely because callers cannot use the value directly in a boolean context.
- **Spec-internal consistency check:** the pitfall claim is consistent with the same file's own documented API (the 6 named constants + accessor methods). No internal contradiction. → ✅
- **No new finding.**

### 1.3 §07 row 173 follow-on — `LazyLock` defers + caches → already ✅ in Cycle 44 (D-CVS-66 LOW filed)

- Listed for completeness; resolved by Cycle 44 §1.1. AC has nothing further to add.

---

## 2. Demotion → NEW finding

### 2.1 §09 row 161 — "`errcore.OverflowType.Fmt(...)` exists" → ❌ NEW **C-CVS-65** (HIGH)

- **Probe:** `grep -rln "Overflow" /tmp/core-v9-upstream/errcore/` returns **zero matches**.
- **Broader probe:** `grep -rln "Overflow" /tmp/core-v9-upstream/ | grep -v _test | grep -v spec/` returns **zero matches** (only the `spec/03-powershell-test-run/07-tc-console-output.md` and `spec/testing-guidelines/02-test-case-types.md` doc references appear, both upstream-internal copies).
- **Verdict:** `errcore.OverflowType` is a **fabricated symbol** in `spec/01-app/09-converters.md` line 161. Code following the spec will not compile. This is the same fabrication-class as C-CVS-29 (`coredynamic.AllFields` etc.) and C-CVS-44 (`errcore.VarTwo` arity drift).
- **Severity:** **HIGH** — the spec instructs callers to wrap converter overflow errors via `errcore.OverflowType.Fmt(...)`; no such type or method exists.
- **Suggested replacement:** `errcore.Expected.MessageVarMap(...)` or `errcore.ShouldBe.MessageVarMap(...)` — both exist in upstream and produce the diagnostic-message shape the spec is reaching for.
- **Spawned remediation:** **AJ-45** — purge `errcore.OverflowType.Fmt(...)` from `spec/01-app/09-converters.md` line 161 (and any sibling references); replace with `errcore.Expected.MessageVarMap` example. Blocked by `spec/01-app/` freeze.

---

## 3. Items retained as ❓ (genuine contract-pass deferrals — beyond AC scope)

| Spec row | Claim | Why still deferred |
|----------|-------|--------------------|
| 99-107 | Conversion safety contract: "no panics", "errcore wrapped", "locale-independent" | Requires reading every method body in `converters/stringTo.go`, `converters/bytesTo.go`, `converters/anyItemConverter.go`, `converters/stringsTo.go`, `converters/unsafeBytesTo.go` — ≈ 60+ methods. Behavioural contract pass; not spec-internal-consistency. **Spawned as AJ-46.** |
| 119 | `IntegerWithDefault(queryParam, 25)` "fall back to default" | Behavioural trace — requires writing a `_test.go` exercise. Tracked under AJ-46. |
| 130 | `parsePagination` example signature | Same — behavioural. Tracked under AJ-46. |
| 172 | "Using `*WithDefault` then re-validating hides the malformed input" | Pure advisory; never strictly verifiable from code (it is a developer-experience warning, not a code property). **Retained ⓘ as "advisory — out-of-band by design."** No further AC pass needed. |

The `99-107 / 119 / 130` cluster is now consolidated under one task (**AJ-46 — Converter behavioural contract pass**) instead of three separate carry-overs. Row 172 is reclassified from "❓" to "ⓘ advisory" and removed from the open ❓ pool.

---

## 4. Updated scoreboard impact

| Section | Pre-Cycle-49 ❓ | Δ | Post-Cycle-49 ❓ |
|---------|----------------:|--:|----------------:|
| §09 (converters)            | 4 (rows 64, 119, 130, 172) | −2 (row 64 → ✅, row 172 → ⓘ) | 2 (rows 119, 130 → AJ-46) |
| §07 (conditional+utilities) | 1 (row 142)                | −1 (row 142 → ✅)              | 0 |

- **Verifiable subset growth:** §09 +1 (row 64); §07 +1 (row 142). §09 ❌ count grows by 1 (C-CVS-65).
- **AB-residual `spec/01-app/` ❓ pool:** post-Cycle-47 was 6 OOB. Post-Cycle-49: **6 − 3 = 3 OOB** (Cycle-19 rows 99-107 cluster + row 130 = 2 entries under AJ-46, plus 1 retained-advisory = 3 active deferrals). 
- **Cumulative AB+AC ❌:** 53 → **54** (CRITICAL unchanged at 23, HIGH +1).

---

## 5. Spawned tasks

- **AJ-45 (HIGH)** — Purge `errcore.OverflowType.Fmt(...)` from `spec/01-app/09-converters.md` line 161; replace with `errcore.Expected.MessageVarMap` example. Blocked by `spec/01-app/` freeze.
- **AJ-46 (MEDIUM)** — Converter behavioural contract pass: write Go-level traces for "no panics / errcore wrapped / locale-independent" claims (rows 99-107) and the `IntegerWithDefault` fallback + `parsePagination` example (rows 119, 130). Output goes into a new `tests/contracttests/converters_test.go` (or equivalent), then promote rows to ✅ in a follow-on AB cycle.
- **S-002 (LOW)** — Suggest tweaking §09 line 64 wording from "duplicates a subset of" to "delegates to a subset of" (see §1.1 mechanism note). Recorded in suggestions tracker; blocked by `spec/01-app/` freeze.

---

## 6. AC dimension status after this cycle

- **Phase 1 of AC (`spec/01-app/` advisory carry-overs):** ✅ DONE for §07 and §09. No remaining AC-class ❓ in either section.
- **Phase 2 of AC (`spec/02-app-issues`, `spec/03-powershell-test-run`, `spec/04-tooling`, `spec/06-testing-guidelines`):** the AB-residual passes (Cycles 41, 42, 27/48) already burned through these advisory pools. AC has nothing further to probe in those sections.
- **Conclusion:** AC umbrella is **complete** for all currently-loaded spec sections. Future spec additions may re-open AC-class items; track per cycle.

_See also: `20-cycle19-AB-converters-promotion.md`, `21-cycle20-AB-conditional-and-utilities.md`, `32-cycle43-AB-residual-spec01-converters.md`, `33-cycle44-AB-residual-spec01-conditional.md`._
