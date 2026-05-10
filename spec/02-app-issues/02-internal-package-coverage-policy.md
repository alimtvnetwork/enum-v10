# 02 — Internal-Package Coverage Policy

> **Status**: resolved (2026-04-23, spec-v0.6.0)
> **Severity**: low
> **Resolution**: AI-driven audit (`spec/99-audits/05-ai-audit-2026-04-23-gemini.md` finding F01) ranked this contradiction as the spec's #1 reproducibility blocker (-5 pts, dropped `internal_consistency` axis to 70). Fixed by clarifying that the rule forbids **coverage-motivated** tests, while **business/integration** tests for internal packages remain allowed under `tests/integratedtests/<pkg>tests/`. Both surfaces (`06-branch-coverage.md` § Internal Package Coverage Policy, and `06-testing-guidelines/README.md` core principle #6) now state the same nuanced rule with a decision tree.

> **Scope note (`enum-v9`, added Cycle 18)**: the `tests/integratedtests/<pkg>tests/` and `csvinternaltests/` / `fsinternaltests/` paths cited in the resolution and original problem statement below describe **upstream `core-v9`** layout. `enum-v9` itself has neither an `internal/` directory nor any `*internaltests/` packages — the entire policy applies to upstream consumers only. See `spec/06-testing-guidelines/README.md` consumer-coverage callout (D-CVS-43, Cycle 15) and `spec/01-app/13-testing-patterns.md` §6.1 for the `enum-v9` `tests/creationtests/` layout. The historical text is preserved verbatim below.

## Resolution summary

| Before | After |
|---|---|
| "Never create tests for `internal/` packages" — contradicted by existing `csvinternaltests/`, `fsinternaltests/`, etc. | "MUST NOT create coverage-motivated tests; MAY add business/integration tests under `tests/integratedtests/`." |
| Two surfaces conflicted silently | Both surfaces aligned + decision tree + cross-link |

## Original problem statement (preserved)

`spec/06-testing-guidelines/06-branch-coverage.md` states "no coverage tests
for `internal/` packages", but the codebase contains many internal-test
packages (`csvinternaltests/`, `fsinternaltests/`, `jsoninternaltests/`,
…). Either the rule is grandfathered or it is being silently violated.

## Reopening criteria

Reopen this issue only if:
- A maintainer rules that **all** internal-package tests must be deleted (would require a separate cleanup task), OR
- A future audit shows the new wording is itself ambiguous to AI agents.
