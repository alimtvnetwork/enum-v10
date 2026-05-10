# Stale `tests/integratedtests/` References Outside `spec/01-app/`

## Description

Cycles 1–12 cleared every occurrence of the stale path `tests/integratedtests/` from `spec/01-app/`. The same drift pattern still exists in three other directories.

## Root Cause

Bulk authoring carried the upstream-only `core-v9` test layout into spec content that documents `enum-v9` (which actually uses `tests/creationtests/`). Same root cause as audit findings C-CVS-01 / D-CVS-17 / D-CVS-26 / D-CVS-27 / D-CVS-32 / D-CVS-33 / D-CVS-36 / D-CVS-37 / D-CVS-39 / D-CVS-40 / D-CVS-41 / D-CVS-43.

## Steps to Reproduce

```bash
rg -n 'tests/integratedtests' spec/03-powershell-test-run/ \
                              spec/04-tooling/04-bootstrap-into-new-repo.md \
                              spec/02-app-issues/02-internal-package-coverage-policy.md
```

Expected: zero hits OR explicit anti-pattern callouts only.
Actual: stale prescriptive references remain.

## Attempted Solutions

- [ ] To be folded into the upcoming directory audits (Cycle 15+) under Task AH.

## Priority

Medium — does not break builds, but degrades spec accuracy and traps future readers.

## Blocked By

Nothing — can be done at any time. Recommended approach: combine with the next directory audit (Task AA / Cycle 15 candidate `spec/06-testing-guidelines/`).

## Related

- `.lovable/memory/04-test-layout.md`
- Task AH in `.lovable/plan.md`
