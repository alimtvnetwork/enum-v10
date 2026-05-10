# core-v9 API Migration — Broken Converter & Coredynamic Calls

## Description

`enum-v9` source code still uses old `core-v8`-era function signatures that no longer exist in `core-v9 v1.5.8`. The upstream library moved from package-level functions to struct-based namespaces (e.g., `converters.AnyToValueString()` → `converters.AnyTo.ValueString()`).

## Root Cause

`core-v9` refactored its `converters` package from top-level exported functions to struct-based namespace vars (`AnyTo`, `StringTo`, `BytesTo`, etc.). The old function names were removed entirely — no backward compatibility shims.

## Steps to Reproduce

```bash
cd enum-v9
go build ./...
# Expected: clean build
# Actual: undefined: converters.AnyToValueString, undefined: coredynamic.TypeName, etc.
```

## Attempted Solutions

- [x] Identified the mapping via user's PowerShell investigation (2026-05-06)
- [ ] Apply migration fixes across all affected files (Task AM)

## Priority

High — blocks Go build after the go.mod bridge was resolved (W+AG done).

## Blocked By

Partial: awaiting `stringTo` method list from user to complete the mapping. See `.lovable/memory/06-core-v9-api-migration.md`.

## Related

- `.lovable/memory/06-core-v9-api-migration.md` (full mapping table)
- Task AM in `.lovable/plan.md`
