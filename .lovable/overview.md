# Project Overview — enum-v10

## Identity

- **Module path:** `github.com/alimtvnetwork/enum-v10` (renamed from `enum-v10`)
- **Type:** Go enum library with PowerShell + Python tooling
- **Frontend:** A React/Vite shell exists but is incidental — all real work happens in Go packages, PowerShell scripts, and Python CI guards.

## Core dependency

- **Import path used in source:** `github.com/alimtvnetwork/core-v9` (renamed from `core-v8`)
- **Upstream status:** ✅ Resolved. `core-v9` `go.mod` now declares `module github.com/alimtvnetwork/core-v9` (tag `v1.5.8`). The `replace` directive has been removed. `go.mod` pins `require core-v9 v1.5.8` cleanly.
- **API migration:** `core-v9` refactored `converters` from package-level functions to struct namespaces (`AnyTo`, `StringTo`, etc.). `coredynamic.TypeName` → `coredynamic.SafeTypeName`. See `.lovable/memory/06-core-v9-api-migration.md`.

## Directories the AI must respect

- `cross-repo/core-v9/` — mirror of a separate upstream repo. **Never** rename `core-v8` → `core-v9` here. **Never** rewrite `enum-v10` → `enum-v10` here.
- `tests/creationtests/` — actual test layout. The string `tests/integratedtests/` in any spec file is **stale** (audit finding C-CVS-01).

## Hard prohibitions

- **No emails.** Never configure Dependabot recipients, SMTP, or any email-based notification flow. User explicitly rejects all email flows.
- **No `core-v8` reintroduction** outside `cross-repo/core-v9/`.
- **No `tests/integratedtests/`** in new spec content.

## Workflow

- User drives via incremental "Next" commands.
- Tasks are tracked with letter IDs (A, B, C, …, AI, AJ, …) carried across turns.
- Every reply MUST end with the remaining-task list.
