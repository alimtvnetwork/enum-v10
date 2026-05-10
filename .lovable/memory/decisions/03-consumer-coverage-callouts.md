# Decision: Consumer-Coverage Callouts for Upstream-Only API

## Date

2026-05-04 (D-CVS-25, §06 / Cycle 4)
2026-05-05 (D-CVS-38, §13 / Cycle 11; D-CVS-42, §14 / Cycle 12)

## Context

Several `core-v9` API surfaces (e.g. `coregeneric`, `corepayload`, `corevalidator`, `coretests`, `tests/testwrappers`) are documented at length in `spec/01-app/` but have **zero `enum-v10` consumers**. Without explicit callouts, future readers (human or AI) would assume the API was verified against `enum-v10` code — but it wasn't, and can't be until Task AB lands.

## Decision

Wherever a section describes an upstream-only API surface, add an explicit **consumer-coverage callout** at the section header:

```markdown
> ⚠️ **Consumer coverage:** The symbols in this section have no `enum-v10` consumer.
> Verification deferred to Task AB (fetch upstream `core-v9` source).
> `enum-v10` readers should redirect to <pointer to actual `enum-v10` equivalent>.
```

The redirect pointer is critical — it tells `enum-v10` readers where to look for the **actually-used** equivalent (e.g. §13 `tests/integratedtests/` callout redirects to §6.1 documenting the real `tests/creationtests/` Goconvey + `EnumTestWrapper` registry).

## Effect

- §06 §1, §2, §6 + §7 — three sub-callouts (D-CVS-25).
- §13 — header-level callout naming `coretestcases.CaseV1`, `args.Map`, `coretests.BaseTestCase`, `testWrapper`, `coretests.GetAssert` (D-CVS-38).
- §14 — header-level callout naming `tests/testwrappers/`, `coretests.GetAssert`, `coretestcases.CaseV1`, `StringsTestWrapper` (D-CVS-42).

## Rule going forward

When a future audit cycle encounters a high-❓ section:
1. Identify which API surfaces have zero `enum-v10` consumers (`rg --type go` returns empty).
2. Add a consumer-coverage callout naming each upstream-only symbol.
3. Add a redirect pointer if `enum-v10` has an equivalent.
4. Score the documented-API claims as ❓ (not ⚠️) — they're unverifiable, not wrong.
