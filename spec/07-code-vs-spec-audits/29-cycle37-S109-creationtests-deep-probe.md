# Cycle 37 — S-109 deep-probe of `tests/creationtests/`

> **Date:** 2026-05-06
> **Auditor:** Lovable agent
> **Suggestion closed:** **S-109** (Cycle-15 deep-probe of `tests/creationtests/` patterns to clear ❓ in `spec/06-testing-guidelines/`)
> **Predecessor cycle:** [Cycle 36](./28-cycle27-AB-scripts-deep-probe.md) (S-103 portable runner reorg — non-audit)
> **Significance:** Settles the 10 ❓ left over from Cycle 15 by direct inspection of `enum-v10/tests/creationtests/` (14 files). Confirms that `enum-v10` deliberately does **not** consume the upstream `coretests`/`args`/`results` framework documented in `spec/06-testing-guidelines/`, so the Cycle-15 ❓ items remain bound to upstream `core-v9` consumers and cannot be promoted from `enum-v10` source alone.

## 1. Method

Direct inspection of all 14 files under `enum-v10/tests/creationtests/`:

| File | Role |
|------|------|
| `EnumTestWrapper.go` | Local test-case struct (NOT `BaseTestCase` extension) |
| `PathPatternTypeCreationTestWrapper.go` | Second local wrapper struct |
| `allBasicEnumsCollection.go` | Imports of every enum-v10 package (registry input) |
| `allEnumGeneralTestCases.go` | Generated registry of `*EnumTestWrapper` |
| `vars.go` | Shared local fixtures |
| `simpleEnumCollectionTestCases.go` | Invalid-enumer fixtures |
| `pathPatternTypeCreationTestCases.go` | Generated PathType fixtures |
| `generateAllEnumGeneralTestCases.go` | One-shot test-case generator (writes Go to stdout) |
| `generateAllBasicEnumTestCases.go` | Companion generator |
| `generatePathPatternTestCases.go` | PathType generator |
| `creation_test.go` | `Test_Creation` over `simpleEnumCollectionTestCases` |
| `AllEnums_ContractsTesting_test.go` | `Test_AllEnums_ContractsTesting` over `allEnumGeneralTestCases` |
| `PathType_Creation_test.go` | `Test_PathType_Creation` |
| `ScriptType_test.go` | `Test_ScriptType` over a map fixture |

Probe commands (sandbox):

```bash
rg -n 'coretests\.|coretestcases\.|args\.Map|args\.One|args\.Six|args\.Holder|args\.LeftRight|CaseV1|CaseNilSafe|GenericGherkins|GetAssert|ShouldBeEqualMap|ShouldBeSafe|InvokeWithPanicRecovery|results\.Result|results\.ResultAny|results\.ExpectAnyError|BaseTestCase' tests/creationtests/
rg -n 'goconvey|Convey\(|So\(' tests/creationtests/
rg -n '^// Arrange|^// Act|^// Assert' tests/creationtests/
rg -n 'tests/integratedtests' tests/
```

**Result:** the first probe returns **zero hits**. The second confirms ubiquitous GoConvey usage. The third confirms AAA comment style. The fourth confirms `tests/integratedtests/` does not exist in `enum-v10`.

## 2. Claim-by-claim promotion table (the 10 ❓ from Cycle 15)

| # (Cycle 15) | File | Claim | Cycle-15 verdict | **Cycle-37 deep-probe verdict** | Evidence |
|---|------|-------|---|---|---|
| 12 | `02-test-case-types` | `CaseV1`, `CaseNilSafe`, `GenericGherkins` are the three case types | ❓ | **❓→ ⓘ upstream-only** | `rg 'CaseV1\|CaseNilSafe\|GenericGherkins' tests/creationtests/` → zero hits. `enum-v10` uses two **local** wrapper structs (`EnumTestWrapper`, `PathPatternTypeCreationTestWrapper`) — neither extends nor mentions `CaseV1`. Confirms claim is upstream-`core-v9`-consumer-only; no `enum-v10` evidence available. |
| 13 | `02-test-case-types` | `CaseV1` field shape (`Name`, `ArrangeInput`, `ExpectedResult`, etc.) | ❓ | **❓→ ⓘ upstream-only** | Same — `CaseV1` is not present. The local `EnumTestWrapper` field shape (`Header`, `InitialBasicEnumer`, `TypeName`, `ExpectedEnumType`, `ExpectedMapValues`, `ExpectedInvalidName`, `ExpectedInvalidValueString`, `ExpectedRangesNamesCsv`, `IntegerMinMax`, `StringMin`, `StringMax`) is independent of the spec's `CaseV1` schema. |
| 14 | `02-test-case-types` | `CaseNilSafe` for nil-receiver tests | ❓ | **❓→ ⓘ upstream-only** | `rg 'CaseNilSafe\|nil-receiver' tests/creationtests/` → zero hits. `enum-v10` has no nil-receiver tests in `tests/creationtests/`. |
| 15 | `03-args-reference` | `args.Map`, `args.One`–`args.Six`, `args.Dynamic`, `args.Holder`, `args.LeftRight` symbol set | ❓ | **❓→ ⓘ upstream-only** | `rg '\bargs\.' tests/creationtests/` → zero hits. `enum-v10` does not import any `args` package — fixtures are plain Go slices/maps (`[]*EnumTestWrapper`, `map[string]ScriptType`). |
| 18 | `04-results-reference` | `results.Result`, `results.ResultAny`, `results.ExpectAnyError`, `InvokeWithPanicRecovery` symbol set | ❓ | **❓→ ⓘ upstream-only** | `rg '\bresults\.\|InvokeWithPanicRecovery' tests/creationtests/` → zero hits. `enum-v10` uses GoConvey `So(err, ShouldBeNil)` directly (e.g. `ScriptType_test.go:21`), not `results.ExpectAnyError`. |
| 19 | `05-assertion-patterns` | `ShouldBeEqual`, `ShouldBeEqualMap`, `ShouldBeSafe` assertion API | ❓ | **❓→ ⓘ upstream-only (partial — GoConvey `ShouldEqual` IS used)** | `rg 'ShouldBeEqual\|ShouldBeEqualMap\|ShouldBeSafe' tests/creationtests/` → zero hits. `enum-v10` instead uses GoConvey-stdlib assertions: `ShouldEqual` (PathType_Creation_test.go:23, ScriptType_test.go:23), `ShouldResemble` (AllEnums_ContractsTesting_test.go:34-37), `ShouldBeTrue`/`ShouldBeFalse` (creation_test.go:33,36), `ShouldBeEmpty` (AllEnums_ContractsTesting_test.go:48), `ShouldBeNil` (ScriptType_test.go:22). The `Should*` family in spec/06 is **upstream-`coretests`-specific** custom assertions distinct from these. One bridge: `errcore.ShouldBe.StrEqMsg` is used in `AllEnums_ContractsTesting_test.go:23` but only as a *header formatter*, not an assertion. |
| 20 | `05-assertion-patterns` | Diff-based assertion pattern | ❓ | **✅ partially evidenced in enum-v10** | `AllEnums_ContractsTesting_test.go:42-47` calls `actualEnumDynamicMap.LogShouldDiffMessage(true, typeName+" - type mismatch", testCase.ExpectedMapValues)` and asserts `So(diffMessage, ShouldBeEmpty)`. This is the **diff-based pattern** the spec describes (compute diff, assert empty diff string), implemented via `enumimpl.DynamicMap.LogShouldDiffMessage` rather than `coretests.ShouldBeEqualMap`. Promotes the *pattern* claim ✅ via behavioural equivalence; the *specific symbol* remains upstream-only. |
| 26 | `07-diagnostics-output-standards` | (5 sub-claims) diagnostic output standards | ❓ | **❓→ ⓘ upstream-only** | None of the 5 sub-claims surface in `tests/creationtests/`. `Test_AllEnums_ContractsTesting` uses GoConvey's stdout reporter; failure messages composed via `errcore.ShouldBe.StrEqMsg(expected, actual)` for the `Convey` header (line 23, 36) — not via `coretests` diagnostics. |
| 27 | `08-good-vs-bad` | Examples of good vs bad tests using `args.Map` / `CaseV1` | ❓ | **❓→ ⓘ upstream-only** | All examples reference `args.Map`/`CaseV1` which `enum-v10` does not consume. Spec-internal: ✅ (still consistent across files). |
| 28 | `09-creating-custom-cases` | `BaseTestCase` extension pattern | ❓ | **❓→ ⓘ upstream-only (alternative pattern in enum-v10)** | `rg 'BaseTestCase' tests/creationtests/` → zero hits. `enum-v10` uses **plain struct registries** instead — `EnumTestWrapper` has no embedded `BaseTestCase`, fixtures live in module-level `var allEnumGeneralTestCases = []*EnumTestWrapper{...}` (generated by `generateAllEnumGeneralTestCases.go`). This is a *different* extension idiom — both work, but the spec's pattern is upstream-only. |

### Summary of promotions

- **1 ❓ → ✅** (claim 20, diff-based assertion pattern — behaviourally equivalent local implementation found).
- **9 ❓ → ⓘ "upstream-only" (annotated)** — confirmed not consumed by `enum-v10`; cannot be promoted from local source. They remain ❓ in the *code-vs-spec* dimension when scoped to upstream `core-v9`, and are blocked by Task **AB** for the upstream clone-based promotion.

### New cumulative score for Cycle 15

- Before: 22 ✅ / 10 ❓ → 22/22 = **100%** (verifiable subset)
- After: 23 ✅ / 9 ❓ → 23/23 = **100%** (verifiable subset; ⓘ items stay flagged as upstream-only but no longer "unknown")

The `100%` ratio is unchanged because the verifiable subset grows by 1 alongside the promotion. The win is the **9 ⓘ annotations** which convert "unknown" into "known-upstream-only", removing them from the open-question pool.

## 3. New finding: D-CVS-64 (LOW)

**Severity:** LOW (documentation completeness; no runtime impact).

**Location:** `spec/06-testing-guidelines/05-assertion-patterns.md` (whole file) and `spec/06-testing-guidelines/02-test-case-types.md`.

**Finding:** Neither file mentions the **GoConvey-only path** that `enum-v10` itself uses (plain `So(actual, ShouldEqual, expected)` + AAA comments + plain struct registries). The portable spec presents the `coretests`-framework path as the only option, when in fact GoConvey-only is a valid sub-pattern that consumers can adopt when they don't need `args.Map` argument bundling or `BaseTestCase` extension. `enum-v10` is itself the proof.

**Proposed fix (NOT applied this cycle — cosmetic, falls under Task AC):** Add a top-of-file note to `02-test-case-types.md` and `05-assertion-patterns.md`:

> **Alternative path — GoConvey-only.** Consumers that don't need `args.Map` bundling can use plain `So(actual, ShouldEqual, expected)` over a module-local `[]*Wrapper` registry. `enum-v10`'s `tests/creationtests/` package is a worked example. The diff-based pattern (§5.X) maps cleanly onto `enumimpl.DynamicMap.LogShouldDiffMessage(...) + So(diff, ShouldBeEmpty)`.

Tracked as carry-forward suggestion **S-111** below.

## 4. Spec-internal-consistency checks (re-run)

```bash
rg -nc 'integratedtests|enum-v1|enum-v3|core-v9 → core-v9|\.lovable/user-preferences|core-v8' spec/06-testing-guidelines/*.md
```

Result identical to Cycle 15 — README + `01-folder-structure.md` callouts still in place. No regression.

## 5. Closing

- **S-109** closed.
- **D-CVS-64** raised LOW; carried forward as **S-111** (cosmetic; non-blocking).
- 10 Cycle-15 ❓ items now resolved as: **1 ✅** + **9 ⓘ upstream-only annotated**.
- Open-question pool for `spec/06-testing-guidelines/` shrinks from 10 → 0 unknown (9 known-upstream-only retained for AB upstream-clone promotion).

