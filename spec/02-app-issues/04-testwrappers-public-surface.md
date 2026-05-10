# 04 — `tests/testwrappers/` Public-Surface Stability

> **Status**: resolved (2026-04-23, spec-v0.8.0)
> **Severity**: low
> **Resolved via**: AI-driven audit finding F03 — explicit stability declaration replaces deferred maintainer call.

> **Scope note (`enum-v9`, added Cycle 18)**: `tests/testwrappers/` is an **upstream `core-v9`** directory. `enum-v9` does NOT have a `tests/testwrappers/` directory at all — it has `tests/creationtests/` with a single shared `EnumTestWrapper` (and `PathPatternTypeCreationTestWrapper`) registered via the Goconvey-based collection in `allBasicEnumsCollection.go`. The "STABLE for any test code inside this module" declaration below applies inside the **upstream `core-v9`** module, not this one. See `spec/01-app/13-testing-patterns.md` §6.1 and `spec/01-app/14-tests-folder-walkthrough.md` consumer-coverage callout for the `enum-v9` wrapper layout. Historical declaration text preserved verbatim below.

## Stability Declaration (authoritative)

All packages under `tests/testwrappers/` are **STABLE for any test code
inside this module**.

- ✅ **MUST** be reused by in-module test packages instead of writing
  duplicate wrappers.
- ✅ **MAY** be extended with new wrappers following the inventory rules
  in `spec/01-app/14-tests-folder-walkthrough.md` §2 (every new wrapper
  package MUST be added to that inventory in the same change).
- ❌ **MUST NOT** be imported by external consumers of this Go module.
  External consumers should treat `tests/testwrappers/*` as internal —
  no backward compatibility guarantees apply across module boundaries.

This declaration is sufficient for autonomous AI agents to reuse wrappers
with full confidence. Publishing wrappers as stable external packages
remains a separate maintainer-only decision and is **out of scope** for
this issue.

## Resolution Rationale (spec-v0.8.0)

The audit (F03) found that the previous `wont-fix` status caused cautious
AI agents to write duplicate wrappers rather than risk an "unstable" API.
The fix is purely a **wording change** in this issue file plus the
explicit "stable for in-module use" rule above — no code changes required.

The discoverability work (referenced below) was already complete; only
the stability ambiguity remained, and it is now removed.

## Discoverability (already in place — reference)

`spec/01-app/14-tests-folder-walkthrough.md` §2 inventories all 4 wrapper
packages with file lists, purpose, and consumer packages. Step 11
Cycle 1 explicitly reused `corevalidatortestwrappers` rather than
authoring a duplicate.

## Reopening criteria

Reopen only when:
- A maintainer decides to publish `tests/testwrappers/*` as stable
  external packages (would expand the stability scope beyond in-module
  use), OR
- A new wrapper is added that does not get inventoried in
  `14-tests-folder-walkthrough.md` §2 in the same change.

## Original problem statement (preserved)

`tests/testwrappers/` contained shared wrappers but had no inventory,
leaving fresh agents to write duplicates.
