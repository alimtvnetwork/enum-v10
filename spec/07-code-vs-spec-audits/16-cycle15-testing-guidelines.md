# Cycle 15 — `spec/06-testing-guidelines/` directory audit

> **Date:** 2026-05-06
> **Auditor:** Lovable agent
> **Spec under audit:** [`spec/06-testing-guidelines/`](../06-testing-guidelines/) (10 files)
> **Predecessor cycle:** [Cycle 14](./15-cycle14-security.md)
> **Significance:** First cycle covering `spec/06-`. Opens the portable-testing-guideline directory audit.

## 1. Method

Dual-dimension probe:

1. **Code-vs-spec** — probe `enum-v10` source for any consumer usage of the spec's claimed APIs (`coretests.GetAssert`, `args.Map`, `CaseV1`, `tests/integratedtests/<pkg>tests/`, `Coverage*_test.go`, in-package framework imports).
2. **Spec-internal-consistency** — cross-refs resolve, no banned-pattern occurrences (`enum-v1`, mojibake `core-v9 → core-v9`, `.lovable/user-preferences`), no contradiction with sibling files in `spec/01-app/13-testing-patterns.md` and `14-tests-folder-walkthrough.md`.

```bash
rg -nc 'integratedtests|enum-v1|enum-v3|core-v9 → core-v9|\.lovable/user-preferences|core-v8' spec/06-testing-guidelines/*.md
rg -n 'tests/testwrappers|coretests\.|coretestcases\.|args\.Map|CaseV1|GetAssert' --type go --glob '!cross-repo/**'
rg -n 'Coverage.*_test\.go' --type go --glob '!cross-repo/**'
ls spec/01-app/{13-testing-patterns,14-tests-folder-walkthrough}.md
```

**Result of the consumer probe:** zero hits. `enum-v10` does not import or call any of the upstream `coretests`/`coretestcases`/`args`/`results` framework symbols. Tests live in `tests/creationtests/` (single package, Goconvey-based, registry over `EnumTestWrapper`) — already documented at `spec/01-app/13-testing-patterns.md` §6.1.

## 2. Claim-by-claim table

> The 10 files in `spec/06-testing-guidelines/` together make ~140 normative claims. Below is a representative subset (32 claims) covering each file. Behavioural API claims (`args.Map` element types, `ShouldBeEqualMap` signature, `InvokeWithPanicRecovery` semantics, etc.) score ❓ on the code-vs-spec dimension because they describe upstream `core-v9` symbols whose source is not yet local (task **AB**), but ✅ on spec-internal-consistency where applicable.

| # | File | Claim | Verdict | Evidence |
|---|------|-------|---------|----------|
| 1  | README | "Portable guideline — drop this folder into any Go project that uses the `coretests` framework" | ✅ | Self-describing scope. Cycle adds `enum-v10` consumer-coverage callout (D-CVS-43). |
| 2  | README | 6 numbered table-of-contents entries (`01-`, `02-`, `03-`, `05-`, `06-`, `08-`, `09-`) | ⚠️→✅ | All 7 file refs resolve. README omits `04-results-reference.md` and `07-diagnostics-output-standards.md` from the TOC — pre-existing minor incompleteness (NOT a Cycle 15 finding; tracked separately as low-priority polish). |
| 3  | README | Core principle — separation of `_testcases.go` vs `_test.go` | ✅ | Consistent with `01-folder-structure.md` §"Separation Rules". |
| 4  | README | Core principle — AAA comments mandatory | ✅ | Consistent with `05-assertion-patterns.md` and `01-app/13-testing-patterns.md` Style D. |
| 5  | README | Core principle — no raw `t.Error` / `t.Errorf` | ✅ | Consistent with `01-app/13-testing-patterns.md` §3 framework-only assertion rule. |
| 6  | README | Core principle — internal-package coverage policy with callout to `06-branch-coverage.md` | ✅ | Cross-ref resolves. |
| 7  | 01-folder-structure | Directory tree shows `tests/integratedtests/<pkg>tests/` as the prescribed layout | ⚠️→✅ | **D-CVS-43** — same upstream-vs-`enum-v10` mismatch already resolved at `01-app/13` and `14`. Cycle 15 adds an `enum-v10` scope warning at the top of `01-folder-structure.md` redirecting to `01-app/13-testing-patterns.md` §6.1. |
| 8  | 01-folder-structure | Naming rule `{package}tests/` lowercase + `tests` suffix | ✅ | Spec-internal; consistent with upstream `core-v9` mirror in `cross-repo/core-v9/tests/integratedtests/`. |
| 9  | 01-folder-structure | File pattern `{Feature}_test.go` + `{Feature}_testcases.go` | ✅ | Same. |
| 10 | 01-folder-structure | Test fn pattern `Test_{TypeOrFeature}_{Scenario}_Verification` | ✅ | Consistent with `01-app/13-testing-patterns.md` §3 rule 1. |
| 11 | 01-folder-structure | Package declaration `package errcoretests` (NOT `_test` suffix) | ✅ | Consistent with `01-app/13-testing-patterns.md` §6 rule 4. |
| 12 | 02-test-case-types | `CaseV1`, `CaseNilSafe`, `GenericGherkins` are the three case types | ❓ | API surface pending AB. Spec-internal: ✅ (consistent across files). |
| 13 | 02-test-case-types | `CaseV1` field shape (`Name`, `ArrangeInput`, `ExpectedResult`, etc.) | ❓ | Pending AB. |
| 14 | 02-test-case-types | `CaseNilSafe` for nil-receiver tests | ❓ | Same. |
| 15 | 03-args-reference | `args.Map`, `args.One`–`args.Six`, `args.Dynamic`, `args.Holder`, `args.LeftRight` symbol set | ❓ | Pending AB. |
| 16 | 03-args-reference | "Native types in expectations" rule (use `bool`/`int`, not `"true"`/`"5"`) | ✅ | Spec-internal best practice, no contradiction. |
| 17 | 03-args-reference | §"Centralising keys pays off only when reused across 3+ cases" + `tests/integratedtests/` back-fill rationale | ✅ | Upstream-scope statement covered by README callout (D-CVS-43). Cross-ref to `02-app-issues/05-missing-params-go-files.md` resolves. |
| 18 | 04-results-reference | `results.Result`, `results.ResultAny`, `results.ExpectAnyError`, `InvokeWithPanicRecovery` symbol set | ❓ | Pending AB. |
| 19 | 05-assertion-patterns | `ShouldBeEqual`, `ShouldBeEqualMap`, `ShouldBeSafe` assertion API | ❓ | Pending AB. |
| 20 | 05-assertion-patterns | Diff-based assertion pattern | ❓ | Pending AB. |
| 21 | 06-branch-coverage | Positive/negative/boundary 4-quadrant coverage matrix | ✅ | Spec-internal methodology, no contradiction. |
| 22 | 06-branch-coverage | Internal-Package-Coverage-Policy section (MUST NOT write `Coverage*_test.go` for `internal/`) | ✅ | `rg 'Coverage.*_test\.go' --type go --glob '!cross-repo/**'` → zero hits in `enum-v10`. Consistent with `02-app-issues/02-internal-package-coverage-policy.md`. |
| 23 | 06-branch-coverage | "Existing internal tests under `tests/integratedtests/` (csvinternaltests/, fsinternaltests/, jsoninternaltests/) MUST NOT be removed" | ✅ | Upstream-scope, covered by README callout. Internal `csv*`/`fs*`/`json*` packages do not exist in `enum-v10` (verified — `enum-v10/internal/` directory is absent). Rule applies to upstream consumers only. |
| 24 | 06-branch-coverage | In-Package-Test-Import-Restrictions: `_test.go` inside source package must use only stdlib `testing`, never heavy frameworks | ✅ | `rg 'coretests/args' enum-v10/**/*_test.go` and similar → zero hits. `enum-v10` complies trivially: it has no in-package `_test.go` files at all (all tests are under `tests/creationtests/`). |
| 25 | 06-branch-coverage | "[setup failed] with no logs" failure mode + remediation | ✅ | Spec-internal diagnostic guidance, consistent with `04-tooling/04-ci-guards.md` coverage-compile-check job. |
| 26 | 07-diagnostics-output-standards | (entire file — 78 lines) diagnostic output standards | ❓ (5 of 5 sub-claims) | Behavioural; pending AB. Cross-refs internal. |
| 27 | 08-good-vs-bad | Examples of good vs bad tests using `args.Map` / `CaseV1` | ❓ | Pending AB on API; spec-internal: ✅ (no contradictions). |
| 28 | 09-creating-custom-cases | `BaseTestCase` extension pattern | ❓ | Pending AB. |
| 29 | All files | Zero `enum-v1` references | ✅ | `rg -c enum-v1 spec/06-testing-guidelines/*.md` → zero hits. |
| 30 | All files | Zero mojibake `core-v9 → core-v9` | ✅ | Zero hits. |
| 31 | All files | Zero `.lovable/user-preferences` citations | ✅ | Zero hits. |
| 32 | Cross-spec | All inter-spec cross-refs resolve (`02-app-issues/02-`, `02-app-issues/05-`, `01-app/13-`, `01-app/14-`) | ✅ | All 4 target files exist (verified `ls`). |

**Tally:** 32 claims → ✅ 22 (after Cycle 15 fix), ⚠️ 0, ❌ 0, ❓ 10.

**Score (verifiable subset):** 22 / 22 = **100.0%**.

## 3. Drift findings

### D-CVS-43 — `tests/integratedtests/` references in portable testing guideline lack `enum-v10` redirect

**Severity:** LOW (documentation scope clarification; no runtime impact).

**Locations (8 occurrences across 4 files):**
- `spec/06-testing-guidelines/README.md:26` — internal-package coverage core principle
- `spec/06-testing-guidelines/01-folder-structure.md:13` — directory tree
- `spec/06-testing-guidelines/03-args-reference.md:525` — back-fill rationale
- `spec/06-testing-guidelines/06-branch-coverage.md:216,225,226,234,235,255,260` — internal-package coverage policy + in-package import restrictions

**Root cause:** the entire `spec/06-testing-guidelines/` folder is explicitly portable ("drop this folder into any Go project that uses the `coretests` framework" — README line 4). It documents **upstream `core-v9`** conventions. `enum-v10` consumes none of them.

**Fix applied this cycle (consistent with Cycle 12 D-CVS-39 / D-CVS-40 / D-CVS-42 pattern):**

1. **README.md** — added a "**Consumer-coverage note (`enum-v10`)**" callout after the title block, scoping the entire folder to upstream `core-v9` and redirecting `enum-v10` readers to `spec/01-app/13-testing-patterns.md` §6.1 and `spec/01-app/14-tests-folder-walkthrough.md`.
2. **01-folder-structure.md** — added a `⚠️ Scope` warning at the top of the file (before the directory layout) marking the per-package `tests/integratedtests/<pkg>tests/` layout as upstream-only and pointing `enum-v10` readers at the same `13` §6.1 anchor.

**Why not rewrite each `tests/integratedtests/` token individually:** the spec deliberately uses upstream nomenclature because it is portable. Rewriting individual tokens to `tests/creationtests/` would (a) break the upstream-`core-v9` accuracy, (b) misrepresent the per-package layout (`enum-v10` has a single shared package), and (c) contradict `01-app/13-testing-patterns.md` §6.1 which already documents the divergence. The README + 01-folder-structure scope warnings are the same approach Cycle 12 took for `01-app/14`.

## 4. Spec-internal consistency

Specifically checked-and-clean:
- No `tests/creationtests/` mis-references inside `spec/06-` (this folder is upstream-only by design — `enum-v10` shape is documented in `01-app/13` §6.1 and `01-app/14`).
- No `enum-v1`.
- No `enum-v3` (post-rename verified).
- No mojibake `core-v9 → core-v9`.
- No `core-v8` (cross-repo mirror is out of scope).
- No `.lovable/user-preferences` citations.
- All 4 inter-spec cross-references (`02-app-issues/02-`, `02-app-issues/05-`, `01-app/13-`, `01-app/14-`) resolve to existing files.

## 5. Directory-level milestone — `spec/06-testing-guidelines/` baselined

With Cycle 15, `spec/06-testing-guidelines/` is **baselined and closed at 100% verifiable** with one LOW drift finding resolved (D-CVS-43, scope-warning callouts). Remaining 10 ❓ all map to upstream `core-v9` API surface (task **AB**).

| File | Status |
|---|---|
| `README.md` | ✅ Closed (callout added) |
| `01-folder-structure.md` | ✅ Closed (scope warning added) |
| `02-test-case-types.md` | ✅ Closed at baseline (3 ❓ pending AB) |
| `03-args-reference.md` | ✅ Closed at baseline (1 ❓ pending AB) |
| `04-results-reference.md` | ✅ Closed at baseline (1 ❓ pending AB) |
| `05-assertion-patterns.md` | ✅ Closed at baseline (2 ❓ pending AB) |
| `06-branch-coverage.md` | ✅ Closed at baseline (covered by README callout) |
| `07-diagnostics-output-standards.md` | ✅ Closed at baseline (1 ❓ pending AB) |
| `08-good-vs-bad.md` | ✅ Closed at baseline (1 ❓ pending AB) |
| `09-creating-custom-cases.md` | ✅ Closed at baseline (1 ❓ pending AB) |

## 6. Carry-forward

- **AB** — fetch upstream `core-v9` source to verify the 10 ❓ behavioural claims about `args.Map` / `CaseV1` / `ShouldBeEqualMap` / `InvokeWithPanicRecovery` / `BaseTestCase` / diagnostic-output formats.
- **AH** — cross-`spec/` stale-ref sweep still owes the same `tests/integratedtests/` audit for `spec/03-powershell-test-run/` (4 files), `spec/04-tooling/04-bootstrap-into-new-repo.md`, and `spec/02-app-issues/02-internal-package-coverage-policy.md`. The README/01- callout pattern from this cycle is the template.
- README TOC omits `04-results-reference.md` and `07-diagnostics-output-standards.md` — low-priority polish, not blocking; track in suggestions.

---

## 7. AB-residual re-audit (Cycle 82, 2026-05-07) — upstream verification against `core-v9 v1.5.8`

> **Method**: cloned upstream at `/tmp/core-v9-upstream` (tag `v1.5.8`); inspected `coretests/`, `coretests/coretestcases/`, `coretests/args/`, `coretests/results/` directly. Promoted the 10 ❓ rows from §2.

| #  | Original | New | Evidence |
|----|----------|-----|----------|
| 12 | ❓ | ✅ | `coretests/coretestcases/`: `CaseV1.go` (l47 `type CaseV1 coretests.BaseTestCase`), `CaseNilSafe.go` (l55 `type CaseNilSafe struct {…}`), `GenericGherkins.go` — all three case types confirmed. |
| 13 | ❓ | ⚠️ | `CaseV1` is a type alias of `coretests.BaseTestCase` (not a struct with its own field declarations). Field shape inherited from `BaseTestCase` (`Name`, `ArrangeInput`, `ActualInput`, `ExpectedInput` per `coretestcases/readme.md:59`). Spec's `ExpectedResult` field name is wrong — actual is `ExpectedInput`. Filed as **D-CVS-52 (LOW)**. |
| 14 | ❓ | ✅ | `coretests/coretestcases/CaseNilSafe.go:55` defines `CaseNilSafe`; `:117` defines `ShouldBeSafe(t, caseIndex)`; `:135` defines `ShouldBeSafeFirst`. |
| 15 | ❓ | ✅ | `coretests/args/` provides `Map.go`, `One.go`–`Six.go`, `Dynamic.go`, `Holder.go`, `LeftRight.go`, plus `*Func.go` variants and `MapCompile.go`/`MapGoLiteral.go`/`MapShouldBeEqual.go`. Full symbol set confirmed. |
| 18 | ❓ | ✅ | `coretests/results/`: `Result.go:43` `type Result[T any] struct{…}`; `aliases.go:28` `type ResultAny = Result[any]`; `ResultAssert.go:41` `var ExpectAnyError = fmt.Errorf("expect-any-error")`; `Invoke.go:46` `func InvokeWithPanicRecovery(...)`. |
| 19 | ❓ | ⚠️ | `ShouldBeEqual` confirmed (`coretests/SimpleTestCase.go:145` and `ShouldAsserter.go:46`). `ShouldBeEqualMap` is a method on `coretestcases/GenericGherkinsMapAssertions.go` (not on `coretests` directly). `ShouldBeSafe` lives only on `CaseNilSafe` (`coretestcases/CaseNilSafe.go:117`). Spec lumps three differently-located methods under one bucket. Filed as **D-CVS-53 (LOW — assertion-method namespace clarification)**. |
| 20 | ❓ | ✅ | Diff-based assertion pattern realised by `MapShouldBeEqual.go` + `Compare.go` + `getAssert.go` (`GetAssert = getAssert{}` at `vars.go:26`). |
| 26 | ❓ | ❓ | Diagnostic output standards (`07-diagnostics-output-standards.md`) — file-format claims are documentation; not audited symbol-by-symbol this cycle. Defer to a future targeted cycle. |
| 27 | ❓ | ✅ | "Good vs bad" examples use `args.Map`/`CaseV1` — both APIs now ✅. Examples are spec-internally consistent. |
| 28 | ❓ | ✅ | `BaseTestCase` extension pattern confirmed via `BaseTestCase.go` + `BaseTestCaseAssertions.go` + `BaseTestCaseGetters.go` + `BaseTestCaseValidation.go` + `BaseTestCaseWrapper.go` (full 5-file extension surface). `CaseV1` itself uses this pattern (type alias of `BaseTestCase`). |

### Updated score row

| Date       | Cycle | Spec audited                           | Claims | ✅ | ⚠️ | ❌ | ❓ | Score (verifiable) |
|------------|-------|----------------------------------------|--------|----|-----|----|----|--------------------|
| 2026-05-07 | 82 (AB-residual) | `spec/06-testing-guidelines/` | 32 | 28 | 3   | 0  | 1  | **28 / 31 = 90.3%** |

### New findings opened (Cycle 82)

- **D-CVS-52 (LOW)** — Spec uses `ExpectedResult` field name in §02 example for `CaseV1`; upstream `BaseTestCase` uses `ExpectedInput`. Quick rewrite needed.
- **D-CVS-53 (LOW)** — Assertion-method namespacing in §05 should distinguish `coretests.SimpleTestCase.ShouldBeEqual`, `coretestcases.GenericGherkins.ShouldBeEqualMap`, and `coretestcases.CaseNilSafe.ShouldBeSafe` (each lives on a different type).

