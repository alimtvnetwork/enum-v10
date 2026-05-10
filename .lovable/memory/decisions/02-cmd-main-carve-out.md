# Decision: `cmd/main/` Smoke-Test Carve-Out

## Date

2026-05-05 (Cycle 10, fix C-CVS-10)

## Context

`spec/01-app/12-cmd-entrypoints.md` originally asserted: "no `cmd/` directory anywhere; `core-v9` is a pure library." But reality:

- Upstream `core-v9` truly has zero `cmd/` directories.
- `enum-v10` ships exactly one: `cmd/main/main.go` (`package main` + `func main()`), used as a smoke-test harness.

So the spec was wrong about `enum-v10`'s reality.

## Decision

Rewrite §1 of `12-cmd-entrypoints.md` as a **library-first, smoke-test allowed** policy:

1. Upstream `core-v9` MUST have zero `cmd/` entrypoints (unchanged).
2. `enum-v10` MAY ship exactly ONE `cmd/main/` smoke-test harness (the existing one).
3. Any additional `cmd/<name>/` is forbidden.

Cross-link `cmd/README.md` documents the smoke-test usage.

## Why this matters

- Distinguishes the two modules cleanly.
- Allows the smoke-test harness without weakening the "library-first" doctrine for upstream.
- Sets a precedent for future "module-specific carve-out" decisions: name the rule, name the exception, name the boundary.

## Don't

- Do NOT add a second `cmd/<name>/` directory in `enum-v10`.
- Do NOT add ANY `cmd/` directory in upstream `core-v9`.
