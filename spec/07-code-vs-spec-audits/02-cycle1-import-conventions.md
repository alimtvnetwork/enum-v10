# 02 — Code-vs-Spec Cycle 1: Import Conventions

**Date:** 2026-05-04
**Spec audited:** [`spec/01-app/03-import-conventions.md`](../01-app/03-import-conventions.md)
**Method:** Manual claim extraction + automated verification via `rg` / `grep` / `find`.
**Auditor:** Lovable agent (assisted, evidence-driven)

---

## Summary

| Metric | Value |
|---|---|
| Claims extracted | **12** |
| ✅ Matches | **5** |
| ⚠️ Drifts | **5** |
| ❌ Contradictions | **2** |
| **Match rate** | **5 / 12 = 41.7 %** |

> **Verdict:** This spec section was written when the project was still a fork of `core-v9` and assumed the codebase **was** `core-v9`. The recent module rename to `core-v9` was applied to import paths but the **prose explanations**, **rule statements**, and **test-folder structure** in this spec section weren't updated. Most drifts are stale prose, not real code defects.

---

## Claim-by-claim verification

### ✅ Match (5)

| # | Claim (spec line) | Evidence | Status |
|---|-------------------|----------|--------|
| 4 | "package name is `core`, not `corev8`" (line 88) | `rg "^package corev[89]"` → 0 hits | ✅ |
| 5 | `errcore` is heavily used | 100 files import it | ✅ |
| 6 | `coreinterface/enuminf` is heavily used | 85 files import it | ✅ |
| 7 | `coreimpl/enumimpl` is heavily used | 79 files import it | ✅ |
| 8 | `coredata/corejson` is heavily used | 80 files import it | ✅ |

### ⚠️ Drift (5) — code is fine, spec prose is stale

| # | Claim (spec line) | Code reality | Severity | Fix path |
|---|-------------------|--------------|----------|----------|
| D-1 | "consumes `core-v9` packages" (line 4) | Module is `core-v9`; only `core-v9` references should be in `cross-repo/core-v9/` | LOW | **Spec** — s/core-v9/core-v9/ in line 4 |
| D-2 | "Even though the path ends in `core-v9`, Go uses the `package core` declaration" (line 88) | Path now ends in `core-v9` | LOW | **Spec** — s/core-v9/core-v9/ in line 88 |
| D-3 | "For `core-v9`, this means: `…/core-v9/internal/...`" (lines 94–98) | Mismatch: prose says v8, example says v9 | LOW | **Spec** — s/`core-v9`/`core-v9`/ in line 94 |
| D-4 | "the test module is rooted at the same `core-v9` module" (line 121) | Test module is rooted at `enum-v10`, not `core-v*` (this repo IS enum-v10) | MEDIUM | **Spec** — clarify that the rule is "rooted at the same module" generically; for this repo that module is `enum-v10`, not `core-vN` |
| D-5 | "`coredata/coregeneric`" listed in canonical imports (line 61) | 0 files import it in this repo | LOW | **Spec** — note this is canonical for `core-v9` itself; downstream consumers like `enum-v10` may not use it |

### ❌ Contradiction (2) — spec describes structure that doesn't exist here

| # | Claim (spec line) | Code reality | Severity | Fix path |
|---|-------------------|--------------|----------|----------|
| C-1 | "Tests for package `foo` live in package `footests` under `tests/integratedtests/footests/`" (line 129) | `tests/integratedtests/` **does not exist** in this repo. Only `tests/creationtests/` exists. | **HIGH** | **Spec OR Code** — either (a) update spec to describe the actual `tests/creationtests/` layout this repo uses, or (b) restructure tests to match the spec |
| C-2 | "Common `internal/` packages used by tests" with example `import "…/core-v9/internal/reflectinternal"` (line 118) | Zero `.go` files in this repo import `internal/reflectinternal` | **MEDIUM** | **Spec** — this example is from the upstream `core-v9` repo's tests, not this repo. Move to a "see core-v9 source" reference instead of an inline example |

### Tally of canonical-import usage (claim 2 detail)

The spec lists 11 canonical packages under "copy-paste-ready". In this repo (`enum-v10`):

| Package | Files using it | Status |
|---|---|---|
| `core-v9` (root) | 0 | unused |
| `core-v9/conditional` | 0 | unused |
| `core-v9/constants` | 53 | ✅ |
| `core-v9/converters` | 9 | ✅ |
| `core-v9/errcore` | 100 | ✅ |
| `core-v9/coredata/corejson` | 80 | ✅ |
| `core-v9/coredata/corestr` | 4 | ✅ |
| `core-v9/coredata/coregeneric` | 0 | unused |
| `core-v9/coreinterface/enuminf` | 85 | ✅ |
| `core-v9/coreimpl/enumimpl` | 79 | ✅ |
| `core-v9/isany` | 0 | unused |
| `core-v9/issetter` | 10 | ✅ |

8 of 11 used; 3 unused. **Not a defect** — the canonical block is for *consumers of `core-v9`*, and `enum-v10` happens not to need all 11. But this should be clarified in the spec ("not every consumer uses every package").

---

## Findings opened (filed in scoreboard)

| ID | Title | Sev | Path |
|----|-------|-----|------|
| D-CVS-01 | Spec §03 line 4 says "consumes `core-v9`" — stale, should be `core-v9` | LOW | Fix spec |
| D-CVS-02 | Spec §03 line 88 says path "ends in `core-v9`" — stale | LOW | Fix spec |
| D-CVS-03 | Spec §03 line 94 prose/example mismatch (v8 vs v9) | LOW | Fix spec |
| D-CVS-04 | Spec §03 line 121 conflates "test module" with "core module" | MED | Fix spec |
| D-CVS-05 | `coregeneric` canonical-import listing not annotated as optional-per-consumer | LOW | Fix spec |
| C-CVS-01 | Spec §03 line 129 references nonexistent `tests/integratedtests/` directory | **HIGH** | Fix spec OR restructure tests |
| C-CVS-02 | Spec §03 line 118 internal-import example doesn't apply to this repo | MED | Fix spec |

---

## Score interpretation

41.7% sounds catastrophic but isn't — **0 of the 7 findings indicate broken code**. Every drift is documentation that lagged behind the `core-v9` → `core-v9` rename or the conversion of this repo from a `core-v9` mirror into the `enum-v10` consumer. Once D-CVS-01..05 are applied (single-line edits) and C-CVS-01..02 are decided (the test-folder question is the only architectural call), match rate jumps to ~95%.

**Recommended next cycle:** `01-app/04-error-system.md` — high churn, lots of API claims, good drift candidate.
