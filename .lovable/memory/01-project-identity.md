# Project Identity

## Module

- **Module path:** `github.com/alimtvnetwork/enum-v9`
- **Previous name:** `enum-v1` (renamed; do NOT reintroduce in source).
- **Type:** Go enum library + PowerShell + Python tooling.
- **Frontend shell:** A React/Vite project lives in `src/` for the Lovable preview. It is **incidental** — real work is in Go packages, PowerShell scripts, and Python CI guards.

## Core dependency

- **Import path used in source:** `github.com/alimtvnetwork/core-v9`
- **Previous name:** `core-v8` (renamed; allowed only inside `cross-repo/core-v9/`).

## Repo layout (top-level highlights)

- `accesstype/`, `enumimpl/`, `cmd/main/`, `tests/creationtests/` — Go.
- `scripts/ci/` — Python CI guard scripts + their `test_*.py` regression tests.
- `scripts/CoveragePreChecks.psm1`, `run.ps1` — PowerShell.
- `cross-repo/core-v9/` — Mirror of an external upstream repo. Different module, different rename rules.
- `spec/` — Documentation. `spec/01-app/` is the audited section.
- `spec/07-code-vs-spec-audits/` — Audit cycle reports and the living scoreboard.
