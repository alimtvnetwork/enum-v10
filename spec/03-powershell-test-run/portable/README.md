# Portable Runner Specs

> **Scope split (S-103, 2026-05-06).** Files in this sub-directory are **portable** — they describe a generic Go-test coverage runner and AI-agent reference applicable to **any Go module / repository**, not just `enum-v9`.
>
> Files in the parent `spec/03-powershell-test-run/` directory are **enum-v9-specific** (use the dashboard UI, `run.ps1` flag set, repo-local conventions).

## Why a sub-directory?

Cycle 16 of the spec audit raised D-CVS-47 and D-CVS-48 because both files mixed a portable promise ("any Go module / repository", "self-contained reference for any AI agent") with `enum-v9`-specific test paths in the body. Top-of-file consumer-coverage callouts patched the misdirection risk, but a structural split is more discoverable: future portable-runner edits ship without touching enum-v9-specific files, and consumers know at a glance which scope they're reading.

## Files

| File | Purpose | Lines |
|------|---------|-------|
| [`01-generic-go-test-coverage-runner.md`](./01-generic-go-test-coverage-runner.md) | Generic runner architecture with sample script. Was `08-generic-go-test-coverage-runner.md`. |
| [`02-ai-agent-complete-reference.md`](./02-ai-agent-complete-reference.md) | Self-contained reference for any AI agent working on a Go project. Was `09-ai-agent-complete-reference.md`. |

## Editor rules

1. **Do NOT add `enum-v9`-specific paths or flags here.** If a runner behaviour is enum-v9-specific, it belongs in the parent directory.
2. **`tests/integratedtests/` references inside these files describe the upstream `core-v9`-style consumer layout.** `enum-v9` itself uses `tests/creationtests/` — that detail belongs in `spec/06-testing-guidelines/01-folder-structure.md` and `spec/01-app/13-testing-patterns.md` §6.1, not here.
3. **Keep the portability promise explicit** in each file's opening paragraph.
