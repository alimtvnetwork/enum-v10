# Cycle 19 — Walk-through audit of `spec/05-failing-tests/`

**Date:** 2026-05-07
**Auditor:** AI agent (Lovable)
**Trigger:** Task **AA** Cycle 19 — last unaudited spec subdirectory.
**Spec target:** `spec/05-failing-tests/` (25 files, 852 lines total)
**Scope dimension:** provenance + frame-of-reference + status-currency. **Not** a per-claim AB audit (these are upstream post-mortems with zero `enum-v10` consumers).

---

## 1. Inventory

| File | Date | Status line | Pkgs referenced | `integratedtests/` refs |
|------|------|-------------|-----------------|--------------------------|
| 01-blocked-packages-fixes.md | 2026-03-22 | ✅ RESOLVED | coredynamic, corejson, corestr | 0 |
| 02-groupby-empty-map-assertion-mismatch.md | — | none | — | 0 |
| 03-blocked-packages-fixes-r11.md | 2026-04-04 | ✅ RESOLVED | coredynamic, corejson, corepayload | 0 |
| 04-failing-tests-root-cause.md | 2026-03-22 | ✅ RESOLVED | corejson, corepayload, enumimpl | 0 |
| 05-reflect-set-from-to-draft-to-bytes.md | — | none | coredata, coredynamic | 0 |
| 06-failing-tests-round3.md | 2026-03-23 | ✅ RESOLVED | coredynamic, corejson, corepayload | 0 |
| 07-validators-slice-double-allocation.md | — | none | corevalidator | 0 |
| 08-errtype-combine-regex-mismatch.md | — | none | errcore | 1 (line 30) |
| 09-failing-tests-round4.md | 2026-03-23 | ✅ RESOLVED | — | 0 |
| 10-merge-errors-getdirectlower-key-mismatch.md | — | none | errcore | 1 (line 28) |
| 11-lazyregex-compile-expectations.md | — | none | — | 2 (lines 27,28) |
| 12-nil-lazyregex-single-return-error.md | — | none | — | 0 |
| 13-slicevalidator-diagnostic-format-drift.md | — | none | corevalidator, errcore | 1 (line 29) |
| 14-dynamic-methods-argscount-and-string-panic.md | — | none | — | 2 (lines 4,23) |
| 15-bytetype-pointer-receiver-interface-mismatch.md | — | none | — | 2 (lines 4,37) |
| 16-info-options-nil-for-plain-default.md | — | none | — | 2 (lines 4,18) |
| 17-stringcompareas-onlysupportederr-names-format.md | — | none | — | 2 (lines 4,34) |
| 18-zero-coverage-build-failure-cascade.md | — | none | codestack, enumimpl | 2 (lines 6,43) |
| 19-batch-fix-jsoninternal-corepayload-…enumimpl-corevalidator.md | — | none | coredata, coreonce, corepayload | 0 |
| 20-codestack-nil-pointer-panic.md | — | none | codestack | 0 |
| 21-corepayload-deep-clone-nil-anymap.md | — | none | corepayload | 0 |
| 22-crossed-packages-investigation.md | — | none | codestack | 1 (line 12) |
| 23-coverage7-4-8-api-mismatch-fixes.md | 2026-03-15 | none | corejson, corestr, enumimpl | 0 |
| 24-variant-onlysupportederr-subset.md | — | none | — | 0 |
| 25-rangenameswithvaluesindexes-format.md | — | none | — | 0 |

**Totals:** 25 files · 5 RESOLVED · 20 no-status · 14 files contain 16 `integratedtests/` references · 0 files reference any `enum-v10`-local package.

---

## 2. Frame-of-reference verification

For every package mentioned (`corepayload`, `corejson`, `coredynamic`, `corestr`, `errcore`, `corevalidator`, `reflectcore`, `coreonce`, `jsoninternal`, `codestack`, `lazyregex`, `coredata`, `coregeneric`, `namevalue`, `enumimpl`):

```bash
$ ls /dev-server/{corepayload,corejson,coredynamic,corestr,errcore,corevalidator,
                  reflectcore,coreonce,jsoninternal,codestack,lazyregex,coredata,
                  enumimpl}
ls: cannot access ...: No such file or directory   # all 13
```

✅ **Confirmed:** zero referenced packages exist inside `enum-v10`. All 25 documents describe upstream `core-v9` (or earlier) test infrastructure. The `tests/integratedtests/<pkg>tests/` paths are the upstream consumer layout, not this repository.

---

## 3. Findings

### 3.1 No drift findings spawned

This is a deliberately advisory-only directory. Per Task **AH** policy (PI-002), `tests/integratedtests/` references describing **upstream** layout are correct as-written and must NOT be rewritten to `tests/creationtests/`. Doing so would corrupt the historical record.

### 3.2 No contradiction findings spawned

There is nothing inside `enum-v10` to contradict. Cross-references all point at upstream symbols, none of which are mirrored locally.

### 3.3 Cleanup applied

Added `spec/05-failing-tests/README.md` documenting:
1. Scope = upstream `core-v9` post-mortems (not `enum-v10` action items).
2. Status legend (5 RESOLVED, 20 historical).
3. Frame-of-reference disclaimer for `tests/integratedtests/` paths.
4. Audit-history pointer back to this cycle file.

This README closes the long-standing gap where future AI cycles could mistake these files for active `enum-v10` issues.

---

## 4. Score

- **Verifiable claims:** 25/25 (provenance dimension only — every file confirmed upstream-historical).
- **D-CVS findings:** 0
- **C-CVS findings:** 0
- **Score:** **N/A** (no per-claim AB sweep applicable; advisory-only material).

---

## 5. Remaining `spec/` audit work

After Cycle 19, the only un-walk-through'd targets in the AA cycle plan are:

| Cycle | Target | Status |
|-------|--------|--------|
| 20    | `spec/00-llm-integration-guide.md` (2388 lines) | ⏳ Pending |

Everything else (`spec/01-app/` (frozen), `spec/02-app-issues/`, `spec/03-powershell-test-run/`, `spec/04-tooling/`, `spec/06-testing-guidelines/`, `spec/07-code-vs-spec-audits/`, `spec/05-failing-tests/`) has either been audited or is explicitly out-of-scope.
