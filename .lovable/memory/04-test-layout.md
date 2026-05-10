# Test Layout

## Reality

Tests live under `tests/creationtests/` in `enum-v10`. The structure is a flat Goconvey + `EnumTestWrapper` registry, NOT a per-package layout.

## Stale spec pattern

The string `tests/integratedtests/` appears in spec docs as an upstream-only convention from `core-v9`. It is **not** the layout used by `enum-v10`. Audit findings:

- C-CVS-01 (Cycle 1)
- D-CVS-17 (Cycle 3, §05)
- D-CVS-26 (Cycle 6, §08)
- D-CVS-27 (Cycle 9, §11)
- D-CVS-32, D-CVS-33 (Cycle 10, §12)
- D-CVS-36, D-CVS-37 (Cycle 11, §13)
- D-CVS-39, D-CVS-40, D-CVS-41 (Cycle 12, §14)
- D-CVS-43 (Cycle 12 collateral cleanup of §01-package-map.md and §02-design-philosophy.md)

All instances in `spec/01-app/` are now corrected. Remaining hits live in:
- `spec/03-powershell-test-run/` (4 files)
- `spec/04-tooling/04-bootstrap-into-new-repo.md`
- `spec/02-app-issues/02-internal-package-coverage-policy.md`

These are scheduled for cleanup as part of Task AH (folded into upcoming directory audits).

## Tooling rule

Test-discovery tooling MUST accept either folder name (or read from disk). Hard-coding one name will silently break either upstream or `enum-v10`.

## Anti-pattern callout

`spec/01-app/05-enum-system.md:417` intentionally references `tests/integratedtests/` as an example of what NOT to do. This single occurrence should NEVER be "corrected" — it's the canonical anti-pattern reference.
