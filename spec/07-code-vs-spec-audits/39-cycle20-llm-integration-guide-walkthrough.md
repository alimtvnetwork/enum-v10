# Cycle 20 — Walk-through Audit: `spec/00-llm-integration-guide.md`

**Date:** 2026-05-07
**Task:** AA (Cycle 20 — final unaudited spec target)
**Auditor:** AI agent
**Target:** `spec/00-llm-integration-guide.md` (2388 lines, 22 H2 sections, 63 H3 sections)
**Upstream baseline:** `/tmp/core-v9-upstream` @ tag `v1.5.8` (module `github.com/alimtvnetwork/core-v9`)

---

## 1. Scope & method

Walk-through audit of the 2388-line monolith that serves as the LLM-onboarding entry point. Each major section was checked against the upstream `core-v9 v1.5.8` source tree and against already-audited `spec/01-app/`, `spec/06-testing-guidelines/`, and `spec/05-failing-tests/` content (Cycles 15, 19, 41–48).

This is an **advisory walk-through**, not a per-claim AB sweep. The guide is a curated digest of content already individually verified in earlier cycles; this pass confirms internal consistency, link integrity, and version-identity correctness.

---

## 2. Findings

### 2.1 Module identity (lines 83–101)

✅ Correct. `module github.com/alimtvnetwork/core-v9`, `go 1.25.0`, root package name `core` — matches `/tmp/core-v9-upstream/go.mod` exactly. F-NEW-04 callout (package name vs. module path) is current.

### 2.2 Package map (lines 118–188)

✅ All 34 top-level directories listed in the package map exist in upstream `core-v9 v1.5.8` (verified by directory probe). `coreimpl/enumimpl/` path is correct (not the older flat `enumimpl/`). `coredata/` sub-packages (`coreapi`, `coredynamic`, `coregeneric`, `corejson`, `coreonce`, `corepayload`, `corerange`, `corestr`, `stringslice`) all present.

### 2.3 Import conventions (lines 192–214)

✅ All 11 import paths use `core-v9` and resolve to real upstream packages.

### 2.4 Enum system / factory reference (lines 453–632)

✅ `enumimpl.New` exposes `BasicByte`, `BasicInt8`, `BasicInt16`, `BasicInt32`, `BasicUInt16`, `BasicString` — confirmed in `/tmp/core-v9-upstream/coreimpl/enumimpl/newCreator.go` lines 26–31. `UsingTypeSlice`, `Default`, `DefaultWithAliasMap`, `CreateUsingMap`, `CreateUsingMapPlusAliasMap` factory methods all present in upstream `newBasicByteCreator.go`.

### 2.5 `tests/integratedtests/` references (lines 36, 825, 828)

✅ Three references remain. **All three are explicitly scoped as upstream-consumer examples**:
- Line 36 (Decision Matrix Style D): explicit upstream example.
- Line 825 (Test Folder Structure preamble): contains the disclaimer *"the layout below describes upstream `core-v9` consumers. `enum-v9` itself uses a single shared `tests/creationtests/` package"* with cross-links to `spec/01-app/13-testing-patterns.md §6.1` and `spec/06-testing-guidelines/01-folder-structure.md`.
- Line 828 (code-fence label): part of the upstream-scoped folder diagram.

These are **policy-compliant** under Task **AH** + PI-002 (do NOT rewrite `integratedtests/` → `creationtests/` when the surrounding context is explicitly upstream-scoped). No action.

### 2.6 Cross-spec link integrity (line 36, 79, all "Further Reading" entries)

✅ All `./01-app/`, `./02-app-issues/`, `./03-powershell-test-run/`, `./04-tooling/`, `./05-failing-tests/`, `./06-testing-guidelines/`, `./99-audits/` links resolve to extant files. Anchor fragments not deep-checked (out of scope for walk-through).

### 2.7 `coremath` deprecation note (lines 758–772)

✅ Consistent with F-NEW-01 callout authored in earlier cycle; matches upstream `coremath/*.go` `// Deprecated:` comments (sampled).

### 2.8 Testing-pattern content (lines 819–1121)

✅ Already individually verified in Cycle 15 (§06-testing-guidelines walk-through) and Cycles 41–47 (§01-app residual deep-probes). No new drift here.

---

## 3. Score

| Dimension | Result |
|-----------|--------|
| Internal consistency vs. audited spec corpus | ✅ Pass |
| Import-path / module-name correctness | ✅ Pass |
| Package-map vs. upstream `v1.5.8` directory probe | ✅ Pass (34/34) |
| Enum-factory-method existence | ✅ Pass |
| `integratedtests/` references appropriately scoped | ✅ Pass (3/3) |
| Cross-spec link existence (file-level) | ✅ Pass |
| **Findings (D-CVS / C-CVS) spawned** | **0** |

**Verdict:** ✅ READY. The guide is a faithful digest of the already-audited spec corpus. Zero drift.

---

## 4. Closure

This closes Task **AA**. All 6 Cycle-15→20 walk-through targets are now audited:

| Cycle | Target | Status |
|-------|--------|--------|
| 15 | `spec/06-testing-guidelines/` | ✅ Done |
| 16 | `spec/03-powershell-test-run/` | ✅ Done |
| 17 | `spec/04-tooling/` | ✅ Done |
| 18 | `spec/02-app-issues/` | ✅ Done |
| 19 | `spec/05-failing-tests/` | ✅ Done |
| **20** | **`spec/00-llm-integration-guide.md`** | **✅ Done (this doc)** |

No follow-up cycles required. The remaining `spec/`-scoped backlog is the **AJ rewrite queue** (~54 items, all blocked by `spec/01-app/` freeze per `spec-v0.53.0` policy).
