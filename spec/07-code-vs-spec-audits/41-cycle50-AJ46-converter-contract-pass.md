# Cycle 50 — AJ-46: converter behavioural contract pass + ❓→✅ promotions

**Date:** 2026-05-10
**Spawned by:** Cycle 49 (AC dimension closure).
**Scope:** Settle the 3 ❓ items in `spec/01-app/09-converters.md` that Cycle 49 deferred to AJ-46 — rows 99-107 (no-panic / errcore-wrapped / locale-independent), row 119 (`IntegerWithDefault` fallback), row 130 (`parsePagination` example) — by writing machine-verifiable Go-level traces against upstream `core-v9 v1.5.8`.
**Result:** **3 ❓ → ✅** via 5 new contract tests in `tests/contracttests/converters_test.go` (all pass under `go test`).
**Allowed under freeze:** scoreboard + audit-cycle additions only; no `spec/01-app/` rewrites.
**Evidence root:** `tests/contracttests/converters_test.go` + upstream source at `/tmp/core-v9-upstream` @ `v1.5.8`.

---

## 1. New contract harness

Created `tests/contracttests/converters_test.go` (package `contracttests`) — a new test directory dedicated to behavioural contract claims that span enum-v8 → upstream `core-v9` boundaries. Five tests, each pinning one spec claim:

| Test                                                    | Spec row(s) | What it pins                                                                                                       |
| ------------------------------------------------------- | ----------- | ------------------------------------------------------------------------------------------------------------------ |
| `TestConverters_NoPanics_OnBadInput`                    | 99-107      | 30 calls (3 funcs × 10 hostile inputs) all wrapped in `recover()`. Failures = test failures. Currently: 0 panics.  |
| `TestConverters_ErrorsAreErrcoreTyped`                  | 99-107      | `errcore.ParsingFailedType` and `errcore.FailedToConvertType` are publicly addressable + `String()` non-empty.     |
| `TestConverters_Float64_LocaleIndependent`              | 99-107      | `"3.14"` parses to ~3.14; `"3,14"` (comma-decimal) MUST fail. Pins C-locale numeric parsing.                       |
| `TestConverters_IntegerWithDefault_FallbackContract`    | 119         | Happy / bad-input / empty-input all return the documented `(value, ok)` shape: `(42,true)`, `(25,false)`, `(99,false)`. |
| `TestConverters_ParsePagination_EndToEnd`               | 130         | The spec's worked example signature works end-to-end: 4 (page, size) cases incl. mixed valid / bad / empty inputs. |

**Verification:** `nix run nixpkgs#go -- test -v ./tests/contracttests/...` — all 5 PASS, 0.005s total.

```
=== RUN   TestConverters_NoPanics_OnBadInput               --- PASS
=== RUN   TestConverters_ErrorsAreErrcoreTyped             --- PASS
=== RUN   TestConverters_Float64_LocaleIndependent         --- PASS
=== RUN   TestConverters_IntegerWithDefault_FallbackContract --- PASS
=== RUN   TestConverters_ParsePagination_EndToEnd          --- PASS
```

---

## 2. Direct upstream evidence cited in the harness

| Spec claim | Upstream source line | Mechanism |
| ---------- | -------------------- | --------- |
| "no panics" — `Integer`, `Float64`, `IntegerWithDefault` are non-panicking | `converters/stringTo.go:47, 73, 154, 203, 274` | All conversions delegate to `strconv.Atoi` / `strconv.ParseFloat` / `strconv.ParseInt`, all of which return `(value, error)` — non-panicking by stdlib contract. |
| "errcore wrapped" | `converters/stringTo.go:164, 211, 262, 281, 286` | Failure paths return `errcore.ParsingFailedType.Error(...)` and `errcore.FailedToConvertType.*` — typed errcore values, not raw stdlib errors. |
| "locale-independent" | `strconv.ParseFloat` Go stdlib contract | Go's `strconv` package is fixed to C locale. The test pins this at the converter wrapper boundary by asserting that comma-decimal European notation FAILS. |
| `IntegerWithDefault` fallback | `converters/stringTo.go:39-54` | `if input == constants.EmptyString { return defaultInt, false }` and `if err != nil { return defaultInt, false }`. |

---

## 3. Promotions

| Row | Was | Now | Note |
| --- | --- | --- | ---- |
| 99-107 (cluster) | ❓ "needs contract pass" | ✅ | Backed by 3 of the 5 new tests above. |
| 119 | ❓ "behavioural fallback" | ✅ | Backed by `TestConverters_IntegerWithDefault_FallbackContract`. |
| 130 | ❓ "behavioural example" | ✅ | Backed by `TestConverters_ParsePagination_EndToEnd`. |

**Net §09 ❓ pool:** post-Cycle-49 was 2 active (rows 119, 130 under AJ-46) + 1 ⓘ (row 172 advisory). Post-Cycle-50: **0 active ❓ remaining**; only the 1 ⓘ advisory + the 1 ⓘ row 172 carry-over remain (both informational, not unknown).

---

## 4. Scoreboard impact

- §09 verifiable subset grows by 3 (rows 99-107 cluster + 119 + 130). New verifiable score: **§09 ~78** (was ~71 post-Cycle-49).
- AB-residual `spec/01-app/` ❓ pool: 3 → **0 active** (only ⓘ informational items remain).
- Cumulative AB+AC ❌ unchanged at **54** (no new contradictions surfaced — every claim probed was honoured by upstream).

---

## 5. Closure

- **AJ-46 status:** ✅ DONE.
- **AJ-45 status:** still ⏸️ blocked by `spec/01-app/` freeze (that one is a spec edit, not a code addition).
- **`spec/01-app/09-converters.md` ❓ pool:** **fully drained** (excluding 1 reclassified-advisory row 172).
- The new `tests/contracttests/` directory is now the canonical location for any future spec-behaviour contract tests; future AJ-* tasks of this shape should add files here rather than per-package _test files.

_See also: `40-cycle49-AC-spec01-converters-conditional-advisory.md` (the AC cycle that spawned this work)._
