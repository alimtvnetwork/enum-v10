# `spec/05-failing-tests/` — Historical Post-Mortems

> **Scope warning:** These 25 files are **upstream `core-v9` test post-mortems** imported into `enum-v10` for institutional memory. They reference packages (`corepayload`, `coredynamic`, `corejson`, `corestr`, `errcore`, `corevalidator`, `reflectcore`, `coreonce`, `jsoninternal`, `codestack`, `lazyregex`, …) and test directories (`tests/integratedtests/<pkg>tests/`) that **do not exist inside `enum-v10`**.
>
> All file-path references that begin with `tests/integratedtests/` describe the **upstream `core-v9` consumer layout**, not this repository. Inside `enum-v10` itself, tests live exclusively under `tests/creationtests/` (per Core memory).
>
> Per the audit AH sweep, this directory is intentionally left in upstream form for forensic traceability — it is **read-only reference material**. Do not "fix" the `integratedtests/` paths to `creationtests/` because they refer to the *upstream* directory, where the post-mortem actually happened.

## Status legend

- ✅ **RESOLVED** — fix landed upstream (5 files: 01, 03, 04, 06, 09)
- 📜 **HISTORICAL** — upstream post-mortem; no status line was ever added (20 files)
- All 25 are advisory-only for `enum-v10`. No file in this directory is a current `enum-v10` action item.

## Audit history

- **Cycle 19** (2026-05-07, `spec-v0.54.0`): scoped this directory as upstream-historical reference material (Task **AA**). 25 files surveyed: 5 explicitly RESOLVED, 20 lack status lines. 14 contain `tests/integratedtests/` path references — all correctly describing the upstream `core-v9` test layout, not `enum-v10`. No drift findings. No D-CVS or C-CVS items spawned. Directory frozen against forensic-rewrite drift via this README. See `spec/07-code-vs-spec-audits/38-cycle19-failing-tests-walkthrough.md`.
