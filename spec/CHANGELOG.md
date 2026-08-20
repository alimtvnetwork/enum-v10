# Spec Changelog

> Tracks **versioned releases of the `spec/` deliverable** (this folder).
> Distinct from the `core-v9` Go module's own version (governed by `spec/01-app/11-versioning.md`).
>
> **Rule** (per `.lovable/user-preferences` line 8 + `spec/01-app/11-versioning.md` §3): every meaningful spec edit cycle bumps **at least minor**.
>
> Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and [Semantic Versioning](https://semver.org/).

---

## [spec-v0.58.0] — 2026-05-10 (Cycle 50 — Task **AJ-46**: converter behavioural contract pass — **§09 ❓ pool DRAINED**)

### Added — `spec/07-code-vs-spec-audits/41-cycle50-AJ46-converter-contract-pass.md`
- 3 ❓ → ✅ in `spec/01-app/09-converters.md`: rows 99-107 cluster (no-panic / errcore-wrapped / locale-independent), row 119 (`IntegerWithDefault` fallback contract), row 130 (`parsePagination` worked example).
- Backed by 5 new tests in `tests/contracttests/converters_test.go` (new package; `nix run nixpkgs#go -- test ./tests/contracttests/...` → all PASS in 0.005s).
- Direct upstream evidence cited: `converters/stringTo.go:47/73/154/164/203/211/262/274/281/286` and `errcore/RawErrorType.go`.
- §09 active ❓ pool: 0 remaining (only the 1 ⓘ row-172 advisory carry-over). AB-residual `spec/01-app/` active ❓ pool: 3 → **0**. Cumulative AB+AC ❌ unchanged at **54**.
- AJ-46 ✅ DONE. AJ-45 still ⏸️ blocked by `spec/01-app/` freeze (pure spec edit).

### Notes
- New `tests/contracttests/` directory is now the canonical home for cross-boundary spec-behaviour tests; future AJ-* contract tasks should add files here rather than per-package `_test` files.

## [spec-v0.57.0] — 2026-05-10 (Cycle 49 — Task **AC**: spec-internal-consistency probe of §07 + §09 advisory carry-overs — **AC CLOSED**)

### Added — `spec/07-code-vs-spec-audits/40-cycle49-AC-spec01-converters-conditional-advisory.md`
- 3 ❓ → ✅: §09 row 64 (`PrettyJson` namespace ↔ `corejson` overlap, structurally consistent — `converters/vars.go:35` + `coredata/corejson/PrettyJsonStringer.go:25`); §07 row 142 (`issetter.Value` "not a drop-in for `bool`" advisory — entailed by 6-state byte enum at `issetter/Value.go:51-58`); §07 row 173 (already promoted in Cycle 44).
- 1 ❓ → ❌ NEW **C-CVS-65 (HIGH)** — `errcore.OverflowType.Fmt(...)` on §09 line 161 is fabricated; `grep -rln Overflow /tmp/core-v9-upstream/errcore/` returns zero matches. Suggested replacement: `errcore.Expected.MessageVarMap(...)`. Spawned **AJ-45** (purge from spec, blocked by `spec/01-app/` freeze).
- 4 ❓ retained as out-of-AC-scope: rows 99-107 (no-panic / errcore-wrapped / locale-independent claims), 119, 130 (behavioural traces) consolidated under new **AJ-46** (converter behavioural contract pass). Row 172 (`*WithDefault` re-validation hides malformed input) reclassified ❓ → ⓘ "advisory by design" — not a measurable code property; removed from open ❓ pool.
- AB-residual `spec/01-app/` ❓ pool: 6 → **3** (Cycle-19 cluster + row 130 under AJ-46, plus 1 retained-advisory). Cumulative AB+AC ❌: 53 → **54** (CRITICAL unchanged at 23, HIGH +1).
- **AC umbrella complete** for all currently-loaded spec sections. AC dimension closed; future spec additions may re-open AC-class items per cycle.

### Notes
- Read-only audit promotion only. No `spec/01-app/` rewrites (under freeze). Spawned tasks (AJ-45, AJ-46, S-002) all blocked by the freeze.
- Editorial: the §09 row 64 "duplicates" → "delegates to" wording suggestion is filed as **S-002 (LOW)** in `.lovable/memory/suggestions/01-suggestions-tracker.md`.

## [spec-v0.56.0] — 2026-05-07 (Recipe-distillation pass — `00-llm-integration-guide.md` §10 lessons-learned)

### Added
- **`spec/00-llm-integration-guide.md` §10** — added "Lessons from Recipe Passes 1–3 (informative)" subsection capturing as-built conventions distilled from three end-to-end validations (`httpmethodtype`, `httpstatusfamily`, `mimetype`):
  1. `Invalid` is the LAST iota member (not first).
  2. Pattern-8 fix is mandatory for `Min`/`Max`/`MinByte`/`MaxByte` overrides.
  3. Type is named `Variant` (disambiguated by package path).
  4. Per-file split (11 files per package) — not single-file.
  5. Required pointer-receiver binding wrappers (`AsJsoner`/`AsJsonContractsBinder`/`AsJsonMarshaller`/`AsBasicByteEnumContractsBinder`/`AsBasicEnumContractsBinder`/`ToPtr`/`Json`/`JsonPtr`/`JsonParseSelfInject`).
  6. Domain-helper convention (1–2 small classifiers/constructors per package).
  7. Single-line registration in `tests/creationtests/allBasicEnumsCollection.go`.
- Added `BasicByte.UsingFirstItemSliceAllCases` to the Factory Method Reference table — the actual factory used by all three recipe-pass packages.

### Notes
- This entry edits `spec/00-llm-integration-guide.md` only. `spec/01-app/` remains frozen per spec-v0.53.0; this guide lives outside the freeze scope.
- Reverse-engineered from real packages, not new design.

## [spec-v0.55.0] — 2026-05-07 (Cycle 20 — Task **AA**: walk-through audit of `spec/00-llm-integration-guide.md` — **AA CLOSED**)

### Added
- `spec/07-code-vs-spec-audits/39-cycle20-llm-integration-guide-walkthrough.md` — final walk-through audit of the 2388-line LLM-onboarding monolith. Verified against upstream `core-v9 v1.5.8` source at `/tmp/core-v9-upstream`: module identity ✅, package map 34/34 dirs ✅, 11 import paths ✅, all 6 `enumimpl.New.Basic*` factory creators ✅, factory methods (`UsingTypeSlice`, `Default`, `CreateUsingMap`, …) ✅. The 3 remaining `tests/integratedtests/` references (lines 36, 825, 828) are explicitly upstream-consumer-scoped (line 825 carries the disclaimer + cross-links to the `enum-v10`-specific layout) and are policy-compliant under Task **AH** / PI-002 — no rewrite. Zero D-CVS / C-CVS findings.

### Notes
- **Task AA is now CLOSED.** All 6 walk-through targets are complete: Cycle 15 (`spec/06-testing-guidelines/`), Cycle 16 (`spec/03-powershell-test-run/`), Cycle 17 (`spec/04-tooling/`), Cycle 18 (`spec/02-app-issues/`), Cycle 19 (`spec/05-failing-tests/`), Cycle 20 (`spec/00-llm-integration-guide.md`).
- Remaining `spec/`-scoped backlog is exclusively the **AJ rewrite queue** (~54 items, all blocked by `spec/01-app/` freeze per `spec-v0.53.0`).
- No file inside `spec/00-llm-integration-guide.md` was modified.

---

## [spec-v0.54.0] — 2026-05-07 (Cycle 19 — Task **AA**: walk-through audit of `spec/05-failing-tests/`)

### Added
- `spec/05-failing-tests/README.md` — directory-level scope/provenance disclaimer. The 25 files in this folder are upstream `core-v9` test post-mortems imported for institutional memory and reference packages (`corepayload`, `coredynamic`, `corejson`, `corestr`, `errcore`, `corevalidator`, etc.) and test paths (`tests/integratedtests/<pkg>tests/`) that **do not exist inside `enum-v10`**. README pins the legend (5 RESOLVED / 20 historical) and instructs future AI cycles NOT to rewrite `integratedtests/` → `creationtests/` here, since those paths refer to upstream layout (per Task **AH** policy / PI-002).
- `spec/07-code-vs-spec-audits/38-cycle19-failing-tests-walkthrough.md` — full walk-through audit. 25/25 files inventoried. Every referenced package confirmed absent from `enum-v10`. Zero D-CVS / C-CVS spawned (advisory-only directory). Score: N/A (no per-claim AB sweep applicable).

### Notes
- This closes Task **AA** Cycle 19. The only remaining un-walk-through'd `spec/` target is `spec/00-llm-integration-guide.md` (2388 lines) → Cycle 20.
- No file inside `spec/05-failing-tests/` was modified.

---

## [spec-v0.53.0] — 2026-05-06 (Cycle 61 — Task **AI**: `spec/01-app/` formally FROZEN)

### Changed
- **🧊 `spec/01-app/` is now FROZEN.** All AB-residual deep-probe sweep work is complete (Cycles 19–25 + 41–47, settled in spec-v0.52.0 / Cycle 47). Zero active probe targets remain across all 14 files. Remaining work on `spec/01-app/` is purely the **AJ rewrite backlog** (~54 items, all already authored as `spec/07-code-vs-spec-audits/AJ-NN-…` companion docs and explicitly tagged "BLOCKED by `spec/01-app/` freeze" in their changelog entries).
- **Freeze policy:** Until the AJ backlog is unblocked by an explicit user "thaw" instruction, no AI cycle may modify any file under `spec/01-app/` for any reason. New findings on `spec/01-app/` content must be authored as new `AJ-NN-*.md` companion docs in `spec/07-code-vs-spec-audits/` and queued behind the freeze. Audit / research / cross-reference work that *reads* `spec/01-app/` is unaffected.
- **Why now:** Cycle 47 closed the last open ❓ pool (§10 reflection-and-dynamic) and explicitly noted "Remaining `spec/01-app/` work is the AJ rewrite backlog (52 items, all blocked by freeze)". This entry promotes that observation to a binding, top-of-CHANGELOG policy so future cycles inherit it without re-reading 7 cycles of audit context.

### Notes
- Scope of freeze: every file under `spec/01-app/` (16 files: `01-overview.md` through `16-security.md`, plus the section index).
- NOT frozen: `spec/02-app-issues/`, `spec/03-powershell-test-run/`, `spec/04-tooling/`, `spec/05-failing-tests/`, `spec/06-testing-guidelines/`, `spec/07-code-vs-spec-audits/`, `spec/99-audits/`, `spec/00-llm-integration-guide.md`, `spec/CHANGELOG.md`.
- Companion AJ items live at `spec/07-code-vs-spec-audits/AJ-*.md` (already created across spec-v0.43.0 onwards). Each carries its own "BLOCKED by freeze" header.
- Thaw instruction shape (for future user reference): _"Thaw `spec/01-app/` and apply AJ-NN through AJ-MM"_ — explicit range, explicit task letter.

---

## [spec-v0.52.0] — 2026-05-06 (Cycle 47 — AB-residual deep-probe of `spec/01-app/10-reflection-and-dynamic.md` — COMPLETES THE SWEEP)

### Added
- `spec/07-code-vs-spec-audits/36-cycle47-AB-residual-spec01-reflection.md` — settles 5 of 6 ❓ items left by Cycle 22. **1 promotion to ✅** (item 6: `isany.DeepEqual` verbatim at `isany/DeepEqual.go:29`). **2 demotions to ❌**: item 5 (NEW **C-CVS-63 HIGH** — `coreonce` exists at `coredata/coreonce/` per **R-CVS-03 retraction** but is typed memoization, NOT reflection-binding); item 2 (NEW **C-CVS-64 HIGH** — `grep '"unsafe"' coredata/corejson/ coredata/coredynamic/` returns zero; the unsafe-pointer fast-path claim is fabricated). **1 retained ❓** (item 3 "type-cache keyed on `reflect.Type`" — plausible-no-emitter, defer to Task AC). **2 out-of-band** (item 1 subjective motivation prose, item 4 benchmark claim — Task AC advisory dimension).
- Spawned **AJ-17b** (delete unsafe-pointer sentence at §4) and **AJ-19b** (rewrite §6 `coreonce` framing) — both folded into existing AJ-17/AJ-19.
- NEW suggestion **S-115**: harden `scripts/spec-api-check.psm1` `Get-UpstreamPackages` to recursively walk `coredata/*` subpackages so future audits don't repeat the R-CVS-01/02/03 "missed parent dir" mistake.
- Scoreboard top-line + per-section row updated (§10 verifiable 38.5% → 37.5% as denominator grows from demotions).

### Notes
- 🎉 **AB-residual deep-probe sweep across `spec/01-app/` is now COMPLETE.** All seven AB-cycle ❓ pools (Cycles 19, 20, 21, 22, 23, 24, 25) settled or classified out-of-band.
- AB-residual `spec/01-app/` ❓ pool drops 11 → 6 (all OOB/plausible-no-emitter; **zero active probe targets remain**).
- Cumulative AB ❌ across 7 sections: 51 → 53 (CRITICAL still 23, HIGH +2).
- Cumulative retractions: R-CVS-01 + R-CVS-02 + **R-CVS-03** — all same drift class (missed `coredata/` parent in initial probe).
- Remaining `spec/01-app/` work is the AJ rewrite backlog (52 items, all blocked by freeze).

---

## [spec-v0.51.0] — 2026-05-06 (Cycle 46 — AB-residual deep-probe of `spec/01-app/08-validators.md`)

### Added
- `spec/07-code-vs-spec-audits/35-cycle46-AB-residual-spec01-validators.md` — settles all 6 ❓ items left by Cycle 21. **2 promotions to ✅** (row 43: `errcore.VarTwoNoType` + `ValidationFailedType` symbol existence at `errcore/VarTwoNoType.go:25` + `RawErrorType.go:121`; row 45: `regexnew.New.Lazy` constructor at `regexnew/newCreator.go:34`). **3 reclassifications to ⓘ "upstream-only"** per Cycle 37 (S-109): row 42 (`coretestcases.CaseV1`/`CaseNilSafe`), row 44 (`<pkg>tests/<V>_Verification_test.go` naming — enum-v10 has flat `tests/creationtests/`, upstream `core-v9` uses `<pkg>tests/` subdirs but neither uses `_Verification_` filenames). **1 out-of-band** (row 46: aspirational diagnostic rules → Task AC advisory dimension). **NO NEW FINDINGS.**
- Scoreboard top-line + per-section row updated (§08 verifiable 33.3% → 42.9%; §08 ❓ pool fully cleared).

### Notes
- AB-residual `spec/01-app/` ❓ pool drops 17 → 11 — **only Cycle 22 (`10-reflection-and-dynamic.md`) has active probe targets left**.
- Cumulative AB ❌ across 7 sections: unchanged at 51 (CRITICAL 23).
- Cumulative ⓘ "upstream-only" annotations: 9 → 11.
- §08 is no longer the worst-drift section (§11 still worst at 18.2%).

---

## [spec-v0.50.0] — 2026-05-06 (Cycle 45 — AB-residual deep-probe of `spec/01-app/11-versioning.md`)

### Added
- `spec/07-code-vs-spec-audits/34-cycle45-AB-residual-spec01-versioning.md` — settles the single ❓ left by Cycle 23 (row 6: "`coreversion` plays well with `coregeneric.Collection`"). One demotion to ❌: `grep -rn coregeneric coreversion/` returns zero hits; `coreversion/VersionsCollection.go` is a hand-rolled `{Versions []Version}` wrapper that imports only `constants`, `coredata/corejson`, `coreinterface`. **NEW C-CVS-62 (HIGH)** — fabricated interop claim.
- Scoreboard top-line + per-section row updated (§11 ❓ pool fully cleared; verifiable score unchanged at 18.2%).

### Spawned (BLOCKED by `spec/01-app/` freeze)
- **AJ-21b** — drop `coregeneric.Collection` interop bullet from `11-versioning.md` §1; document real `VersionsCollection` surface. Folded into existing AJ-21 §1 constructor rewrite.

### Notes
- AB-residual `spec/01-app/` ❓ pool drops 18 → 17.
- Cumulative AB ❌ across 7 sections: 50 → 51 (CRITICAL still 23 — C-CVS-62 is HIGH).
- §11 fully closed for ❓ — all 11 claims now have verdicts (2 ✅ / 9 ❌ / 0 ❓).

---

## [spec-v0.49.0] — 2026-05-06 (Cycle 44 — AB-residual deep-probe of `spec/01-app/07-conditional-and-utilities.md`)

### Added
- `spec/07-code-vs-spec-audits/33-cycle44-AB-residual-spec01-conditional.md` — settles 2 of 3 ❓ items left by Cycle 20. Both promotions to ✅: row 51 (`LazyLock` defers + caches — confirmed via `regexnew/LazyRegex.go:34-110` `mu sync.Mutex` + `isCompiled` guard); row 52 (`corecmp` returns `CompareEqual/Less/Greater` = 0/-1/1 — verbatim at `constants/constants.go:336-338`). One **NEW D-CVS-66 (LOW)** — mechanism-name drift: spec implies `sync.Once`, real impl uses `sync.Mutex` + boolean guard (functionally equivalent but cached-error semantics require this design).
- Scoreboard top-line + per-section row updated (§07 verifiable 70.6% → 73.7%).

### Spawned (BLOCKED by `spec/01-app/` freeze)
- **AJ-04b** — add footnote at `07-conditional-and-utilities.md:173` clarifying `sync.Mutex` + `isCompiled` mechanism vs `sync.Once`.

### Notes
- AB-residual `spec/01-app/` ❓ pool drops 20 → 18.
- Cumulative AB ❌ across 7 sections: unchanged at 50 (23 CRITICAL).
- Out-of-band: row 50 (advisory pitfall on `issetter.Value`) deferred to Task AC.

---

## [spec-v0.48.0] — 2026-05-06 (Cycle 43 — AB-residual deep-probe of `spec/01-app/09-converters.md`)

### Added
- `spec/07-code-vs-spec-audits/32-cycle43-AB-residual-spec01-converters.md` — settles 4 of 8 ❓ items left by Cycle 19. Three promotions to ✅ (rows 57, 58, 60: `BytesTo.PrettyJsonString` intent, PrettyJson↔corejson overlap, `IntegerWithDefault` fall-back behaviour). One demotion to ❌ (row 62: **C-CVS-61 CRITICAL** — `errcore.OverflowType.Fmt(...)` fabricated; zero `Overflow` hits in upstream `errcore/`). One new drift **D-CVS-65 (LOW)** — spec line 54 should call `converters.PrettyJson.Bytes.Safe(jsonBytes)` not `BytesTo.PrettyJsonString(...)`.
- Scoreboard top-line + per-section row updated (§09 verifiable 66.7% → 68.4%).

### Spawned (BLOCKED by `spec/01-app/` freeze)
- **AJ-03b** — rewrite `09-converters.md:54` to real `PrettyJson.Bytes.Safe` call shape.
- **AJ-44** — drop `errcore.OverflowType.Fmt(...)` at `09-converters.md:161`; replace with real `errcore` builder.

### Notes
- AB-residual `spec/01-app/` ❓ pool drops 24 → 20.
- Cumulative AB ❌ across 7 sections: 49 → 50 (23 CRITICAL).
- Out-of-band: 4 remaining ❓ classified (1 unprobeable per C-CVS-11; 3 deferred to Task AC contract pass).

---

## [spec-v0.44.0] — 2026-05-06 (Cycle 37 — S-109 `tests/creationtests/` deep-probe)

### Added
- `spec/07-code-vs-spec-audits/29-cycle37-S109-creationtests-deep-probe.md` — full deep-probe report of all 14 files under `enum-v10/tests/creationtests/`. Settles the 10 ❓ items left by Cycle 15 in `spec/06-testing-guidelines/`: **1 promoted ✅** (claim 20, diff-based assertion pattern — found in `AllEnums_ContractsTesting_test.go` via `enumimpl.DynamicMap.LogShouldDiffMessage` + `So(diff, ShouldBeEmpty)`), **9 annotated ⓘ "upstream-only"** (`CaseV1`/`CaseNilSafe`/`GenericGherkins`, `args.*`, `results.*`/`InvokeWithPanicRecovery`, `ShouldBeEqual*`/`ShouldBeSafe` upstream-custom assertions distinct from GoConvey's `ShouldEqual`/`ShouldResemble`/`ShouldBeNil` family used by enum-v10, `07-diagnostics-output-standards` 5 sub-claims, `08-good-vs-bad` examples, `09-creating-custom-cases` `BaseTestCase` extension pattern). Probe commands documented for reproducibility. Closes **S-109**.
- New finding **D-CVS-64** (LOW): `02-test-case-types.md` and `05-assertion-patterns.md` don't surface the **GoConvey-only sub-pattern** that `enum-v10` itself uses (plain `So(actual, ShouldEqual, expected)` + AAA comments + plain `[]*Wrapper` registry, no `args.Map` bundling, no `BaseTestCase` extension, no `coretests.GetAssert`). Tracked as carry-forward suggestion **S-111** — cosmetic, non-blocking.
- `spec/07-code-vs-spec-audits/01-scoreboard.md` — Cycle 37 row added.

### Notes
- Cycle-15 verifiable subset grows 22/22 → 23/23 (still 100%). The spec/06 unknown ❓ pool drops **10 → 0** — the 9 ⓘ items are now *known* to be upstream-only consumer claims, not unknown.
- Cumulative AB ❌ unchanged at **49** across 7 sections (24 CRITICAL). This cycle promotes/annotates ❓ items only.
- Direct evidence files: `tests/creationtests/EnumTestWrapper.go` (local wrapper struct, NOT `BaseTestCase` extension), `AllEnums_ContractsTesting_test.go` (GoConvey + diff-pattern + `errcore.ShouldBe.StrEqMsg` header formatter), `creation_test.go` / `PathType_Creation_test.go` / `ScriptType_test.go` (AAA-commented GoConvey tests over plain registry slices/maps).

---

## [spec-v0.43.0] — 2026-05-06 (Cycle 36 — S-103 portable runner reorg)

### Changed
- Moved `spec/03-powershell-test-run/08-generic-go-test-coverage-runner.md` → `spec/03-powershell-test-run/portable/01-generic-go-test-coverage-runner.md`.
- Moved `spec/03-powershell-test-run/09-ai-agent-complete-reference.md` → `spec/03-powershell-test-run/portable/02-ai-agent-complete-reference.md`.
- Updated live cross-refs in `spec/00-llm-integration-guide.md` (line 2380) and `spec/04-tooling/03-powershell-implementation.md` (line 456) to the new paths.
- Updated internal cross-ref inside the new `portable/02-ai-agent-complete-reference.md` (table row pointing to its sibling).

### Added
- `spec/03-powershell-test-run/portable/README.md` — explains the scope split (portable vs `enum-v10`-specific), lists the two files, and lays out three editor rules to keep the portability promise intact.

### Notes
- Historical `08-`/`09-` filename references in `spec/CHANGELOG.md` (Cycle-16 entry), `spec/07-code-vs-spec-audits/17-cycle16-powershell-test-run.md`, and `spec/99-audits/01-original-11-step-plan.md` are intentionally left as-is — they document the audit history at the time those cycles ran.
- Closes **S-103**. The structural split makes it easier for future portable-runner edits to ship without touching `enum-v10`-specific files, and replaces the reliance on top-of-file consumer-coverage callouts (D-CVS-47/48 from Cycle 16) with a directory-level signal.

---

## [spec-v0.42.0] — 2026-05-06 (Cycle 27 — AB residual deep-probe of scripts + workflows)

### Added
- `spec/07-code-vs-spec-audits/28-cycle27-AB-scripts-deep-probe.md` — first dedicated "scripts deep-probe" cycle. Promoted 11 ❓ from runner-internal/workflow-internal claims (6 from §03 Cycle 16 + 5 from §04 Cycle 17) using direct evidence from `scripts/CoverageRunner.psm1`, `scripts/PreCommitCheck.psm1`, `.github/workflows/release.yml`, `.github/workflows/ci.yml`, and `run.ps1`.
- 2 NEW LOW drift findings:
  - **D-CVS-62** — `scripts/coverage/Generate-CoveragePrompts.ps1` is referenced by `scripts/CoverageRunner.psm1:313-316` and `scripts/PackageCoverage.psm1:150` (with `-BatchSize 500`) but is **MISSING** from the repo. Both call-sites guard with `if (Test-Path)` so the missing-file case silently no-ops. Tracked as suggestion **S-108**.
  - **D-CVS-63** — `spec/03-powershell-test-run/04-pre-commit-api-checker.md` JSON-schema example was missing the `source` field (actual schema in `PreCommitCheck.psm1:169` writes 7 fields, spec showed 6). Resolved this cycle by adding the `source` field to the example.

### Changed
- `spec/03-powershell-test-run/04-pre-commit-api-checker.md` — added `source` field to the JSON schema example to match `scripts/PreCommitCheck.psm1` ground truth.
- `spec/07-code-vs-spec-audits/01-scoreboard.md` — Cycle 27 row added; AB-residual ❓ tally drops 53 → 42; cumulative AB ❌ unchanged at 49.

### Notes
- Cumulative AB ❌ unchanged at **49** across 7 sections (24 CRITICAL). This cycle promotes ❓ items only.
- 3 §04 Cycle-17 ❓ deferred (out-of-band metadata + cosmetic UI tokens — non-behavioural).
- 21 ❓ remain in `spec/06-testing-guidelines/` (Cycle 15) — distinct probe technique (read `tests/creationtests/` patterns), deferred to future cycle.

---

## [spec-v0.42.1] — 2026-05-06 (Cycle 34 — S-004 scope callout)

### Changed
- `spec/00-llm-integration-guide.md` — added an inline upstream-vs-enum-v10 scope callout above the "Test Folder Structure" code fence (line 826), matching the disclaimer model already used by `spec/06-testing-guidelines/01-folder-structure.md` line 3. Cross-links `spec/01-app/13-testing-patterns.md` §6.1 and `spec/06-testing-guidelines/01-folder-structure.md`. Closes suggestion **S-004** and disambiguates the Decision-Matrix Style-D row (line 36) by the same callout — `integratedtests/` references are now correctly scoped to upstream `core-v9` consumers.

### Notes
- **S-003** verified already-resolved: `spec/06-testing-guidelines/01-folder-structure.md` line 3 already carries the same scope disclaimer.
- **PI-002** (Cross-spec stale `integratedtests/` paths) plan marked complete in `.lovable/memory/pending-issues/01-all-pending-issues.md`.

---

## [spec-v0.41.0] — 2026-05-06 (Cycle 26 — S-106 self-audit & retractions)

### Added
- `spec/07-code-vs-spec-audits/27-cycle26-S106-self-audit-retractions.md` — first run of `scripts/spec-api-check.psm1` v1.0.0; surfaces 2 audit-author errors and 1 new aggregate finding.
- 2 retraction findings:
  - **R-CVS-01** retracts **C-CVS-29 (Cycle 22)** — `coredynamic` package exists at `coredata/coredynamic/` (20+ files). The 8 specific method symbols (`AllFields`/`SetField`/`InvokeMethod`/etc.) remain confirmed sym-fabrications. Original CRITICAL pkg-fab claim re-classified as LOW path-drift + HIGH sym-fab.
  - **R-CVS-02** retracts **C-CVS-51 (Cycle 25)** — `corestr` package exists at `coredata/corestr/` (30+ files). `StringBuilder`/`IsValidUTF8`/`NewCollectionPtrUsingStrings` remain confirmed sym-fabrications.
- 1 NEW aggregate finding **C-CVS-60 (LOW)** — bundles 4 individually-low-impact sym-fabs (`coronce.New`, `enumimpl.NewBasicByte`, `errcore.OverflowType`, `errcoreinf.SomeErrorer`).
- AJ rescoping: **AJ-15** split into AJ-15a (path-qualify §2) + AJ-15b (purge 8 fabricated symbols); **AJ-36/37/38** re-scoped (keep `corestr` package, purge fabricated symbols only).

### Changed
- `spec/07-code-vs-spec-audits/01-scoreboard.md` — Cycle 26 row, headline updated (cumulative AB ❌ 50 → **49**, CRITICAL 24 → **22**, fabrication rate ~52 %), Open-findings list extended to C-CVS-11..60 with retractions noted.

### Notes
- **🛡️ S-106 v1.0 BUILT** at `scripts/spec-api-check.psm1`. Indexes 182 upstream packages and 10,216 top-level symbols; scans `spec/01-app/**/*.md` for `package.Symbol` refs; reports unresolved packages and unresolved symbols within existing packages. Local-variable tracking per-fence + `ProseLooseMode` heuristic for ≤3-char tokens eliminate most false positives. Known limits: presence-only (does NOT catch arity, return-type, or receiver-shape drift — C-CVS-44/45/49-class items remain undetected; tracked as S-106 v2 needing Go AST pass).
- **AJ-01..43 may now safely proceed** with S-106 as a guardrail against fresh fabrications during rewrites.

---

## [spec-v0.40.0] — 2026-05-06 (Cycle 25 — AB pass 7 on `01-app/16-security.md` — **completes AB sweep of `spec/01-app/`**)

### Added
- `spec/07-code-vs-spec-audits/26-cycle25-AB-security.md` — seventh AB promotion pass; 13 ❓ → **3 ✅ / 1 ⚠️ / 9 ❌ / 0 ❓** (verifiable score **66.7 %**). §16 inherits fabrication from earlier cycles: trust-boundary rules cite fabricated `coredynamic.*`, `corevalidator.New.Line/Slice` fluent, `corestr.*`.
- 9 NEW contradiction findings **C-CVS-51..59**:
  - **C-CVS-51 (CRITICAL)** — `corestr` package does not exist in upstream (zero matches for `find . -type d -name corestr*`).
  - **C-CVS-52 (CRITICAL)** — `corestr.StringBuilder` fabricated; consumers must use stdlib `strings.Builder`.
  - **C-CVS-53 (CRITICAL)** — `corestr.IsValidUTF8` fabricated; real call is stdlib `unicode/utf8.ValidString`.
  - **C-CVS-54 (CRITICAL)** — `errcore.InvalidInput.MergeError(...)` won't compile; `InvalidInput` is not an `errcore` category. `errcore/vars.go` exposes only `ShouldBe`, `Expected`, `StackEnhance` at package level.
  - **C-CVS-55 (CRITICAL)** — `coredynamic.AllFields` rule (§4 rule 4) unactionable — package doesn't exist (inherited from C-CVS-29).
  - **C-CVS-56 (CRITICAL)** — `coredynamic.SetField`/`InvokeMethod` rules (§5 rules 2-3) unactionable — same.
  - **C-CVS-57 (CRITICAL)** — Trust-boundary §6 example uses fabricated `corevalidator.New.Line.NotEmpty().MaxLength(255).Matches(re).Build()` + `result.IsFailed()`/`result.Error()` shape (inherited from C-CVS-22/23).
  - **C-CVS-58 (HIGH)** — `corevalidator.New.Slice.MaxLength(N)` cited in 4 separate locations (§4 rule 1, §6 rule 2, §7 mistake row) — all fabricated.
  - **C-CVS-59 (HIGH)** — `errcore.VarTwo("password", pwd, …)` example in §2 reproduces C-CVS-44 missing-`isIncludeType bool` defect; folded into AJ-28.
- 1 NEW drift finding **D-CVS-61 (LOW)** — `coregeneric` import path should be `coredata/coregeneric` (same drift class as C-CVS-14 for `corejson`).
- 9 new action items **AJ-35..43** — all BLOCKED on `spec/01-app/` 🧊 freeze waiver.

### Changed
- `spec/07-code-vs-spec-audits/01-scoreboard.md` — Cycle 25 row, headline updated (cumulative AB ❌ = 50; ~54 % fabrication rate; 24 CRITICAL across 7 sections), Open-findings list extended to C-CVS-11..59 + D-CVS-61.

### Notes
- **🎉 AB sweep of `spec/01-app/` is COMPLETE.** All 7 sections that previously held ≥10 ❓ (§07, §08, §09, §10, §11, §15, §16) have been promoted. Pattern across 7 cycles: 50 ❌ accumulated, ~54 % fabrication rate, 24 CRITICAL items (~48 %). **S-106 lint remains MANDATORY** before any AJ rewrite.

---

## [spec-v0.39.0] — 2026-05-06 (Cycle 24 — AB pass 6 on `01-app/15-observability.md`)

### Added
- `spec/07-code-vs-spec-audits/25-cycle24-AB-observability.md` — sixth AB promotion pass; 13 ❓ → **6 ✅ / 7 ❌ / 0 ❓** (verifiable score **74.1 %**). Section drops from a clean baseline due to fabricated function signatures and a fabricated test-failure output format.
- 7 NEW contradiction findings **C-CVS-44..50**:
  - **C-CVS-44 (CRITICAL)** — `errcore.VarTwo` example missing mandatory leading `isIncludeType bool` parameter; spec example will not compile.
  - **C-CVS-45 (CRITICAL)** — `VarTwo`/`VarTwoNoType`/`MessageVarMap` return `string`, not `error`; spec assigns to `err :=` throughout §2, misframing the entire helper family as error builders.
  - **C-CVS-46 (HIGH)** — §2.4 decision table presents `VarTwo`/`VarTwoNoType` as semantically distinct; reality: `VarTwoNoType` is literally `VarTwo(false, ...)`.
  - **C-CVS-47 (HIGH)** — `coretests/results/ResultAny.go` does not exist; real files are `Result.go`, `ResultAssert.go`, `Results.go`.
  - **C-CVS-48 (CRITICAL)** — Test-failure output format `Test #N — {scenario}: should be equal\n  expected:\n  actual:` is fabricated; zero matches in `coretests/results/`.
  - **C-CVS-49 (CRITICAL)** — `errcore.HandleErr` does NOT attach stack-enhanced wrapping — it just `panic(err.Error())`. Spec rule §3 #2 cites the wrong rationale (correct advice, wrong reason).
  - **C-CVS-50 (MEDIUM)** — `StackEnhance` documented as 2-method (`Error`, `Msg`); reality is 8 methods including the `*Skip` family used by 24 internal call-sites.
- 7 new action items **AJ-28..34** — all BLOCKED on `spec/01-app/` 🧊 freeze waiver.

### Changed
- `spec/07-code-vs-spec-audits/01-scoreboard.md` — Cycle 24 row, headline updated (cumulative AB ❌ = 41; ~52 % fabrication rate across 6 sections), Open-findings list extended to C-CVS-11..50.

### Notes
- Pattern across 6 cycles: 41 ❌ accumulated, ~52 % fabrication rate, 18 CRITICAL items. **S-106 lint remains MANDATORY** before any AJ rewrite.

---

## [spec-v0.38.0] — 2026-05-06 (Cycle 23 — AB pass 5 on `01-app/11-versioning.md` — **WORST section in project**)

### Added
- `spec/07-code-vs-spec-audits/24-cycle23-AB-versioning.md` — fifth AB promotion pass; 11 claims → **2 ✅ / 8 ❌ / 1 ❓** (verifiable score **18.2 %** — new worst score in project).
- 8 NEW contradiction findings **C-CVS-37..43**:
  - **C-CVS-37 (CRITICAL)** — `coreversion.Parse(s) (Version, error)` fabricated; real constructor is `coreversion.New.Default(s) Version` and returns no error (uses `IsInvalid` flag).
  - **C-CVS-38 (HIGH)** — typed accessors `v.Major()/Minor()/Patch()` fabricated; real surface is **public fields** `VersionMajor/VersionMinor/VersionPatch/VersionBuild int` plus `MajorString()/MinorString()/PatchString()/BuildString()` helpers.
  - **C-CVS-39 (CRITICAL)** — fluent comparison `v1.LessThan(v2)/Equal/GreaterThanOrEqual` fabricated; real comparison is package-level `coreversion.Compare(left, right *Version) corecomparator.Compare`.
  - **C-CVS-40 (LOW)** — `String()` claim "returns `\"v1.2.3\"`" inaccurate; returns the `Compiled` field which can be empty for `Empty()`/`Invalid` versions.
  - **C-CVS-41 (HIGH)** — "Wraps stdlib errors in `errcore.FailedToConvertType`" rationale fabricated; `coreversion/` package has zero `errcore` references.
  - **C-CVS-42 (HIGH)** — package path `versionindexes/` is wrong; real path is `enums/versionindexes/`.
  - **C-CVS-43 (CRITICAL)** — constants `versionindexes.V1=1, V2=2, V8=8` "version eras" entirely fabricated; real consts are `Major=0, Minor=1, Patch=2, Build=3, Invalid=4` indexing **version-component positions**, not historical eras. **Conceptual error, not just API drift** — the spec invents a different *purpose* for the package.
- 7 new action items **AJ-21..27** (rewrites for §1 constructor, §1 accessors, §1 compare, §1 `String()`, §1 errcore bullet, §2 import path, §2 entire framing) — all BLOCKED on `spec/01-app/` 🧊 freeze waiver.

### Changed
- `spec/07-code-vs-spec-audits/01-scoreboard.md` — Cycle 23 row, headline replaced (§11 now WORST at 18.2 %; cumulative AB ❌ = 34; ~55 % fabrication rate across 5 sections), Open-findings list extended to C-CVS-11..43.

### Notes
- Pattern across 5 cycles: 34 ❌ accumulated, ~55 % fabrication rate, ~44 % CRITICAL severity. **S-106 lint is now MANDATORY** before any AJ rewrite — without it, the same author pattern that produced 34 ❌ will likely produce more during rewrites.
- §11 is uniquely bad among the 5 audited sections because C-CVS-43 invents an entirely different *purpose* for `versionindexes` (eras vs. component positions). Cycles 19–22 fabricated APIs that *plausibly could exist*; this cycle fabricates a package's reason for existing. Future audits should explicitly cross-check package-purpose claims against upstream `readme.md` / `doc.go`.

---

## [spec-v0.37.0] — 2026-05-06 (Cycle 22 — AB pass 4 on `01-app/10-reflection-and-dynamic.md` — second-worst section)

### Added
- `spec/07-code-vs-spec-audits/23-cycle22-AB-reflection-and-dynamic.md` — fourth AB promotion pass; 19 claims → **5 ✅ / 8 ❌ / 6 ❓** (verifiable score **38.5 %** — second-worst score in project, after §08 33.3 %).
- 8 NEW contradiction findings **C-CVS-29..36**:
  - **C-CVS-29 (CRITICAL)** — entire `coredynamic` package documented in §2 does not exist in upstream `core-v9 v1.5.8`. `coredynamic/` directory absent; `grep -rln coredynamic` returns zero source files. `InvokeMethod`, `HasMethod`, `MethodNames`, `GetField`, `SetField`, `AllFields`, `TypeFullName`, `IsNullOrUndefined` all fabricated.
  - **C-CVS-30 (CRITICAL)** — `reflectcore.IsPointer/IsStruct/IsSlice/IsMap/IsFunc/IsChannel/IsInterface` as bare top-level predicates do not exist; real shape is aggregate `reflectcore.Is.Pointer(v)` etc.
  - **C-CVS-31 (CRITICAL)** — `reflectcore.WalkFields` fabricated; field walking lives at `reflectcore/reflectmodel.FieldProcessor`.
  - **C-CVS-32 (CRITICAL)** — `reflectcore.GetTag` fabricated.
  - **C-CVS-33 (CRITICAL)** — `reflectcore.DerefAll` fabricated.
  - **C-CVS-34 (HIGH)** — `internal/reflectinternal` "off-limits to consumers" framing misleading: `reflectcore/vars.go` publicly re-exports 15 internal symbols (`Converter, Utils, Looper, CodeStack, GetFunc, Is, TypeName, TypeNames, TypeNamesString, TypeNamesReferenceString, ReflectType, ReflectGetter, ReflectGetterUsingReflectValue, SliceConverter, MapConverter`).
  - **C-CVS-35 (HIGH)** — claimed `errcore` panic-wrapping facade does not exist in `reflectcore`.
  - **C-CVS-36 (HIGH)** — §7 mistake-row guidance about `InvokeMethod`/`HasMethod` fabricated (those methods don't exist).
- 6 new action items **AJ-15..20** (rewrites for §2, §3.1, §3.2, §3.3, §1+§4 framing, §7 table) — all BLOCKED on `spec/01-app/` 🧊 freeze waiver.

### Changed
- `spec/07-code-vs-spec-audits/01-scoreboard.md` — added Cycle 22 row, headline updated, ❓ count reduced 98 → 89.

### Notes
- **Cumulative fabrication rate is now ~45 %** across 4 audited sections (29 ✅ / 26 ❌ / 23 ❓ promoted out of 80 baseline ❓).
- **S-106 lint is now MANDATORY** before any AJ rewrite (was "strongly recommended" in spec-v0.36.0).
- `spec/01-app/` 🧊 freeze remains in force; no `spec/01-app/` files modified by this cycle (audit file lives under `spec/07-code-vs-spec-audits/`).

---

## [spec-v0.36.0] — 2026-05-06 (Cycle 21 — AB pass 3 on `01-app/08-validators.md` — worst-drift section in project)

### Added
- `spec/07-code-vs-spec-audits/22-cycle21-AB-validators.md` — third AB promotion pass; 18 ❓ → 4 ✅ / 8 ❌ / 6 ❓ (verifiable score 33.3 % — lowest in project).
- Scoreboard cycle history row 21.

### Changed
- Scoreboard top-line now leads with §08 = 33.3 % ahead of §09 = 66.7 % and §07 = 70.6 %.
- Open-drift block now covers C-CVS-11..28 (18 ❌, 6 of which are CRITICAL severity).
- Remaining-❓ count drops 110 → 98.

### Discovered (NOT yet patched into spec — blocked by freeze)
- **C-CVS-21 (CRITICAL)** Entire 5-method validator contract (`IsValid/IsSuccess/IsFailed/Message/Error`) fabricated; no `IsSuccessValidator` interface in `coreinterface/`.
- **C-CVS-22 (CRITICAL)** Entire fluent-builder API (`corevalidator.New.Line.NotEmpty().MaxLength().Build()`, slice/range equivalents) fabricated; no `New` var, no fluent methods.
- **C-CVS-23 (CRITICAL)** No `Validate(string) Result` method; real surface is `IsMatch(lineNumber, content, isCaseSensitive) bool`.
- **C-CVS-24 (HIGH)** No `RangeValidator` type; numeric-range checks live in `coremath/integer*Within.go`.
- **C-CVS-25 (HIGH)** `StringCompareAs` is a `Variant` enum in `enums/stringcompareas/`, not a "specialty validator".
- **C-CVS-26 (HIGH)** No `corevalidator.Result` type.
- **C-CVS-27 (HIGH)** "Authoring a Custom Validator" §4 template teaches the fabricated 5-method contract; uses wrong `*regexnew.Lazy` field type for what is actually `*regexnew.LazyRegex`.
- **C-CVS-28 (CRITICAL)** "PowerShell runner parses validator output" attribution pipeline does not exist; runner parses `go test` output.

### Bumped (pending action items, blocked)
- **AJ-08..14** — corresponding rewrites of §1 contract / §2.1-2.4 examples / §2.5 enum reclassification / §3.1 result-slice / §4 custom-validator template / §5 attribution pipeline / §4 type fix.

### Notes
- Cumulative AB pattern: across the 3 sections audited so far, **41 % of API claims are fabricated** (18 ❌ vs 26 ✅). Validators chapter alone contributes 8 ❌.
- **S-106 (`spec-api-check.psm1`) is now strongly recommended before lifting the freeze for any AJ rewrite** — every fabrication in this cycle would be caught by `go vet` against a synthesized stub file.

---

## [spec-v0.35.0] — 2026-05-06 (Cycle 20 — AB pass 2 on `01-app/07-conditional-and-utilities.md`)

### Added
- `spec/07-code-vs-spec-audits/21-cycle20-AB-conditional-and-utilities.md` — second AB promotion pass; 17 ❓ → 12 ✅ / 5 ❌ / 3 ❓ (verifiable score 70.6 %).
- Scoreboard cycle history row 20.

### Changed
- Scoreboard top-line now leads with §07 = 70.6 % and §09 = 66.7 % (post-AB) ahead of the 100 % closed sections.
- Open-drift block now covers C-CVS-11..20 (10 ❌, all blocked by `spec/01-app/` 🧊 freeze).
- Remaining-❓ count drops 124 → 110.

### Discovered (NOT yet patched into spec — blocked by freeze)
- **C-CVS-16 (HIGH)** `conditional.TypedErrorFunctionsExecuteResults` is a branch-selector `(isTrue, trueFns, falseFns)`, not a `(fn1, fn2)` fan-out aggregator.
- **C-CVS-17 (HIGH)** `coremath` Min/Max exists only for `Byte`/`Int`/`Float32` (3 families, not 7); `Int16/32/64` + `Float64` are absent.
- **C-CVS-18 (CRITICAL)** `namevalue.NewInstance(...)` constructor does not exist; `Instance[K,V]` is a generic struct with public `Name`/`Value` fields, not methods. No `ValueAny()`.
- **C-CVS-19 (HIGH)** `Collection.ToMap()` does not exist.
- **C-CVS-20 (CRITICAL)** Entire `keymk.New.Compile(...)` snippet fabricated; real surface is `keymk.NewKey.Create(opt, main).CompileKeys(...)`.

### Bumped (pending action items, blocked)
- **AJ-04..07** — corresponding rewrites of §1.3 / §5 / §7 / §8.

### Note
Cumulative AB pattern: across the 2 sections audited so far, ~25 % of API claims are fabricated. Suggests S-106 (`spec-api-check.psm1`) lint should land **before** any further AB-fix waiver; otherwise AJ rewrites risk introducing new fabrications.

---

## [spec-v0.34.0] — 2026-05-06 (Cycle 19 — AB pass 1 on `01-app/09-converters.md`)

### Added
- `spec/07-code-vs-spec-audits/20-cycle19-AB-converters-promotion.md` — first ❓→ground-truth promotion pass enabled by Task **AB** (upstream `core-v9 v1.5.8` cloned to `/tmp/core-v9-upstream`).
- Scoreboard cycle history row 19 (66.7 % verifiable for `09-converters.md`).

### Changed
- Scoreboard top-line now leads with §09 = 66.7 % (post-AB pass 1) ahead of the 100 % closed sections, reflecting that AB promotions can lower a section's verifiable score by replacing ❓ with ❌.
- Open-drift block now lists C-CVS-11..15 (5 ❌, all blocked by `spec/01-app/` 🧊 freeze).

### Discovered (NOT yet patched into spec — blocked by freeze)
- **C-CVS-11 (HIGH)** `converters.StringTo.Integer64` does not exist.
- **C-CVS-12 (HIGH)** `converters.StringTo.Float32` does not exist.
- **C-CVS-13 (HIGH)** `converters.StringTo.Bool` does not exist; closest equivalent lives at `typesconv.StringToBool` and returns `bool`, not `(bool, error)`.
- **C-CVS-14 (MEDIUM)** `converters.PrettyJson.String` / `.FromAny` should be `.PrettyString` / `.SafePrettyString` (or other `PrettyString*` variants).
- **C-CVS-15 (CRITICAL)** Entire `typesconv` §2 + §4.3 example fabricated — actual surface is `*Ptr`, `*PtrToSimple`, `StringToBool`, etc., not numeric widening (`IntToInt64`, `Int64ToInt32`, `Float64ToInt`).

### Bumped (pending action items, blocked)
- **AJ-01** rewrite §1.1 of `09-converters.md` against verified API.
- **AJ-02** rewrite §2 + §4.3 of `09-converters.md` against the real `typesconv` surface.
- **AJ-03** correct `PrettyJson.*` callsites at §1.3 and §4.

### Note on freeze
`spec/01-app/` remains DRIFT-FROZEN (spec-v0.30.0). AB promotions are explicitly allowed under the freeze; the AJ-01..03 rewrites that would fix the discovered ❌s are NOT — they need a one-shot waiver from the user.

---

## [spec-v0.33.0] — 2026-05-06 (Cycle 18 baseline & closed — `spec/02-app-issues/` directory at 100 % verifiable; 🎉 cross-`spec/` AH sweep COMPLETE)

### Added
- **`spec/07-code-vs-spec-audits/19-cycle18-app-issues.md`** — Cycle 18 audit covering all 11 files in `spec/02-app-issues/` (26 representative claims across 402 lines). Closes the directory at **100 % verifiable** (21 ✅ / 5 ❓ audit-history claims). Raises and resolves **D-CVS-56 → D-CVS-60** (5 LOW drifts) in the same cycle. **Marks task AH (cross-`spec/` cleanup sweep) Done** — every directory under `spec/` outside the immutable history folders has been baselined.

### Fixed
- **`spec/02-app-issues/README.md`** (D-CVS-56) — index was stale by ~4 issues for ~14 days (since spec-v0.6.0 introduced #06–#09 and resolved all 9). Rewrote the index table to 9 rows all ✅ resolved; updated the top-of-file status banner to "All 9 issues resolved (last update spec-v0.8.0)".
- **`spec/02-app-issues/02-internal-package-coverage-policy.md`** (D-CVS-57) — added a Cycle-18 `Scope note (enum-v10)` after the status banner, explaining that the cited `tests/integratedtests/<pkg>tests/` and `csvinternaltests/`/`fsinternaltests/` paths describe upstream `core-v9` (which has an `internal/` directory); `enum-v10` has no `internal/` at all, so the policy applies vacuously here. Historical resolution text preserved verbatim.
- **`spec/02-app-issues/03-getassert-undocumented-api.md`** (D-CVS-58) — added an `enum-v10`-scope footnote noting `coretests.GetAssert` returns zero hits in `enum-v10` (Goconvey assertions are used inside `EnumTestWrapper` instead). Historical "STABLE for any test code inside this module" declaration preserved verbatim — it now correctly applies to upstream `core-v9` only.
- **`spec/02-app-issues/04-testwrappers-public-surface.md`** (D-CVS-59) — added an `enum-v10`-scope footnote noting `tests/testwrappers/` does NOT exist in `enum-v10` (the project uses a single shared `EnumTestWrapper` inside `tests/creationtests/`). Historical declaration preserved verbatim.
- **`spec/02-app-issues/05-missing-params-go-files.md`** (D-CVS-60) — added an `enum-v10`-scope footnote noting `tests/integratedtests/` and `errcoretests/` are upstream-`core-v9` package names; the "grandfathered, no back-fill" rule applies vacuously here (`tests/creationtests/` uses shared `vars.go`, no per-package `params.go`).

### Changed
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — added Cycle 18 history row; lifted top-line milestone to include **🎉 Cross-`spec/` AH sweep COMPLETE**; bumped open-❓ aggregate from **172 → 177** (5 audit-history claims).

---

## [spec-v0.32.0] — 2026-05-06 (Cycle 17 baseline & closed — `spec/04-tooling/` directory at 100 % verifiable)

### Added
- **`spec/07-code-vs-spec-audits/18-cycle17-tooling.md`** — Cycle 17 audit covering all 10 files in `spec/04-tooling/` (30 representative claims across 2 553 lines). Closes the directory at **100 % verifiable** (22 ✅ / 8 ❓; the 8 ❓ are workflow/script-internal behaviours requiring direct `.github/workflows/*.yml` and `scripts/*.psm1` probes). Raises and resolves **D-CVS-49 → D-CVS-55** (7 LOW drifts) in the same cycle. Folds in residual task **AH** debt for this directory.

### Fixed
- **`spec/04-tooling/00-overview.md`** (D-CVS-49) — fixed 2 broken `cross-repo/core-v9/` paths (Map table row 06 + Maintenance §3) to `cross-repo/core-v9/` with explicit Core-memory note that the directory intentionally retains its historical `core-v8` name even though the import path is `core-v9`.
- **`spec/04-tooling/04-bootstrap-into-new-repo.md`** (D-CVS-50) — §7 decoupling row "tests/integratedtests/ mirror layout required ❌ No" now names both upstream-`core-v9` and `enum-v10` (`tests/creationtests/`) layouts as concrete examples. Closes the AH-tracked occurrence for this directory.
- **`spec/04-tooling/06-cross-repo-sync.md`** (D-CVS-51 → D-CVS-55) — fixed 5 stale-token sites: `enum-v10 → enum-v10` at lines 11, 80 (template comment), 91; `cross-repo/core-v9 → cross-repo/core-v9` at lines 19, 80 (template comment), 103 (See Also). Template comment at line 80 carried both stale tokens (D-CVS-53). Each rewrite includes a Core-memory clarification where appropriate.

### Changed
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — added Cycle 17 history row; lifted top-line milestone to include **`spec/04-tooling/` baselined & closed**; bumped open-❓ aggregate from **164 → 172** (8 from `spec/04` workflow-internal behaviours).

---

## [spec-v0.31.0] — 2026-05-06 (Cycle 16 baseline & closed — `spec/03-powershell-test-run/` directory at 100 % verifiable)

### Added
- **`spec/07-code-vs-spec-audits/17-cycle16-powershell-test-run.md`** — Cycle 16 audit covering all 9 files in `spec/03-powershell-test-run/` (28 representative claims across 2 519 lines). Closes the directory at **100 % verifiable** (22 ✅ / 6 ❓; the 6 ❓ are runner-internal behaviours — parallel-threading model, JSON-schema fidelity — requiring a direct `scripts/*.psm1` probe). Raises and resolves **D-CVS-44 → D-CVS-48** (5 LOW drifts) in the same cycle. Folds in task **AH** debt for this directory.

### Fixed
- **`spec/03-powershell-test-run/01-overview.md`** (D-CVS-44) — added a top-of-file **`Scope note (enum-v10)`** explaining `run.ps1` is layout-agnostic and example JSON paths (`tests/integratedtests/corecmptests/...`) use upstream-`core-v9` package names for illustration. The runner reads test packages from disk via `go list ./tests/...` and works on either layout.
- **`spec/03-powershell-test-run/04-pre-commit-api-checker.md`** (D-CVS-45) — added the same `Scope note (enum-v10)` callout, redirecting to `01-overview.md` and `01-app/13-testing-patterns.md` §6.1.
- **`spec/03-powershell-test-run/06-coverage-prompt-generator.md`** (D-CVS-46) — inline rewrite of the prompt template (line 71) to name both upstream-`core-v9` and `enum-v10` test layouts. Inline rewrite chosen because the line is *quoted prompt output* — a top-of-file scope note wouldn't propagate into the generated prompt.
- **`spec/03-powershell-test-run/08-generic-go-test-coverage-runner.md`** (D-CVS-47) — added a top-of-file **`Consumer-coverage note (enum-v10)`** scoping the entire portable-runner spec to upstream-`core-v9` nomenclature.
- **`spec/03-powershell-test-run/09-ai-agent-complete-reference.md`** (D-CVS-48) — added the same `Consumer-coverage note (enum-v10)` callout. Per-token rewrite of the 7 in-text `tests/integratedtests/` references avoided to preserve the file's portability promise (header line: "self-contained reference for any AI agent working on a Go project").

### Changed
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — added Cycle 16 history row; lifted top-line milestone to include **`spec/03-powershell-test-run/` baselined & closed**; bumped open-❓ aggregate from **158 → 164** (6 from `spec/03` runner-internal behaviours).

---

## [spec-v0.30.0] — 2026-05-06 (Freeze marker — `spec/01-app/` directory closed for code-vs-spec drift)

### Added
- **`spec/01-app/` — DRIFT-FROZEN marker.** With Cycle 14 closing the last numbered file (`16-security.md`) and the directory audit complete (12 files at 100 % verifiable + 2 baseline-only — §07, §09 — awaiting task **AB**), `spec/01-app/` is now declared **frozen for code-vs-spec drift work**. Future edits to files in this directory MUST fall into one of these allowed categories:
  1. **Upstream-source verification** — when task **AB** lands, ❓ claims may be promoted to ✅ (or new findings raised) for §04 (7), §05 (1), §06 (6), §07 (17), §08 (18), §09 (23), §10 (15), §11 (11), §12 (6), §13 (8), §14 (10), §15 (13), §16 (13) — **148 ❓ total**.
  2. **Re-audit of §07 / §09** under the spec-internal-consistency dimension (task **AC**).
  3. **New normative content** introduced because of an upstream `core-v9` API change (must be paired with a new audit cycle row in `spec/07-code-vs-spec-audits/01-scoreboard.md`).
  4. **Typo / formatting / link fixes** (no normative change).

  All other categories of edit (drive-by rewording, reorganisation, additional examples, etc.) are **out of scope** until the freeze is explicitly lifted in a future `spec-vX.Y.0` entry. Drift work moves to `spec/03-powershell-test-run/`, `spec/04-tooling/`, and `spec/02-app-issues/` (Cycles 16+).

### Changed
- **`spec/CHANGELOG.md`** — recorded the freeze as its own versioned entry so the marker is discoverable from the changelog (not buried in an audit cycle file). No file under `spec/01-app/` is touched by this entry; the freeze is a process declaration, not a content edit.

---

## [spec-v0.29.0] — 2026-05-06 (Cycle 15 baseline & closed — `spec/06-testing-guidelines/` directory at 100 % verifiable)

### Added
- **`spec/07-code-vs-spec-audits/16-cycle15-testing-guidelines.md`** — Cycle 15 audit covering all 10 files in `spec/06-testing-guidelines/` (32 representative claims). Closes the directory at **100 % of its verifiable subset** (22 ✅ / 10 ❓; the 10 ❓ are upstream-`core-v9` behavioural claims pending task **AB**). Introduces **D-CVS-43** (LOW) — same upstream-vs-`enum-v10` scope mismatch already resolved at `01-app/13` and `01-app/14` — and resolves it in the same cycle via the README + 01-folder-structure callout pattern.

### Fixed
- **`spec/06-testing-guidelines/README.md`** (resolves **D-CVS-43** part 1) — added a "**Consumer-coverage note (`enum-v10`)**" callout immediately after the title block, scoping the entire portable testing-guideline folder to **upstream `core-v9`** and redirecting `enum-v10` readers at `spec/01-app/13-testing-patterns.md` §6.1 and `spec/01-app/14-tests-folder-walkthrough.md` for this module's actual `tests/creationtests/` layout.
- **`spec/06-testing-guidelines/01-folder-structure.md`** (resolves **D-CVS-43** part 2) — added a `⚠️ Scope` warning at the top of the file marking the per-package `tests/integratedtests/<pkg>tests/` directory tree as upstream-only and redirecting to the same `01-app/13-` §6.1 anchor. Other in-text references in `03-args-reference.md` and `06-branch-coverage.md` are now covered by the README callout (no per-token rewrite — the spec is deliberately portable).

### Changed
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — added Cycle 15 history row; lifted top-line milestone string to include **`spec/06-testing-guidelines/` baselined**; bumped open-❓ aggregate from **148 → 158** (10 from `spec/06`).

---

## [spec-v0.28.0] — 2026-05-04 (Tooling — surface blocked-package compile diagnostics inline)

### Fixed
- **`scripts/CoverageCompileCheck.psm1`** — when a test package fails the pre-coverage compile check, the actual `go test` diagnostic is now echoed inline directly under the "✗ Blocked: …" line (capped at 25 lines per package, with a pointer to `data/coverage/blocked-packages.txt` for the full text). Previously the only on-screen signal was the package name, forcing the user to open the saved blocked-packages report to diagnose. Applies to both sync and parallel modes via a new `Write-BlockedDiagnostic` helper. **Closes operational task AF.**

## [spec-v0.27.0] — 2026-05-04 (Cycle 5 baseline — §07 conditional-and-utilities audited)

### Added
- **`spec/07-code-vs-spec-audits/06-cycle5-conditional-and-utilities.md`** — Cycle 5 baseline audit of `spec/01-app/07-conditional-and-utilities.md`. All 17 claims (covering `conditional`, `isany`, `issetter`, `regexnew`, `coremath`, `corecmp`, `coresort`, `corefuncs`, `namevalue`, `keymk`) classified as **❓ unverifiable**: zero `enum-v10` consumers and no source mirrored under `cross-repo/core-v9/`. No drift or contradiction provable from `enum-v10` alone — verification deferred to task **AB** (fetch upstream `core-v9` source). Section coverage advances **4 / 16 → 5 / 16**; verifiable-match rate unchanged.

### Changed
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — added Cycle 5 (baseline) row marked **N/A** *(no verifiable subset)*; updated section-coverage milestone to **5/16**; restated **AB** task as now spanning **17 §07 + 7 §04 + 1 §05 + 6 §06** ❓ claims.

## [spec-v0.26.0] — 2026-05-04 (Cycle 4 closed — §06 data-structures at 100 % verifiable)

### Fixed
- **`spec/01-app/06-data-structures.md`** — resolves C-CVS-06..08 + D-CVS-20..25 from Cycle 4 in a single pass. Lifts §06 from **35.7 % → 100 %** of its verifiable subset; **❌ contradiction count: 3 → 0**.
  - **§1 "Consumer-coverage note"** (resolves **D-CVS-25**): added a callout listing actual `enum-v10` consumer counts (`corejson` 80 / `corestr` 4 / `coreonce` 1; `coregeneric` and `corepayload` zero) so readers know which sub-sections are upstream-reference.
  - **§2 `coregeneric` header**: added explicit "⚠️ Upstream-only sub-package" callout.
  - **§3 `corestr` rewrite** (resolves **D-CVS-23**): replaced the unused `corestr.NewCollectionPtrUsingStrings(&values, 0)` example with the actual exported surface (`corestr.New.Hashset`, `corestr.New.SimpleSlice`, `corestr.SimpleStringOnce`); cross-referenced the upstream-only `coregeneric.New.Collection.String` for callers that genuinely need a mutex-protected string list.
  - **§4 `corejson` code block** (resolves **C-CVS-07** + **D-CVS-20** + **D-CVS-21**): replaced fictional `Serialize.ToString` / `Serialize.Raw` / `Deserialize.UsingBytes` / `Deserialize.FromTo` with the real consumer-side API: `Serialize.ToBytesErr` returning `*Result`, `Deserialize.BytesTo`, and `corejson.NewPtr(...).PrettyJsonString()`. Added a contracts subsection naming `Jsoner` / `JsonMarshaller` / `JsonContractsBinder`.
  - **§4 "Rule (with documented exceptions)"** (resolves **C-CVS-06**): rule wording softened from "**never** `encoding/json` directly" to "**should** prefer `corejson`" with an explicit table of the two legitimate exceptions: `inttype/Variant.go:440` calling `json.Marshal(it.Value())` inside `MarshalJSON`, and `inttype/all-constructors.go:75` accepting `*json.Number` as a parameter type.
  - **§5 `coreonce` rewrite** (resolves **D-CVS-22** + **D-CVS-24**): replaced the fictional `coreonce.New.String(producer)` namespace with the real top-level constructors `coreonce.NewAnyOnce` / `coreonce.NewByteOnce`; cross-referenced `corestr.SimpleStringOnce` as the string equivalent (which lives in `corestr`, not `coreonce`); softened "all common types" wording.
  - **§6 `corepayload` upstream-only callout** (resolves **C-CVS-08**): added "⚠️ Upstream-only sub-package" callout deferring `PayloadCreateInstruction` field-set verification to task **AB**.
  - **§7 decision matrix**: added a "Verified in `enum-v10`?" column (✅ vs ⚠️ upstream-only); replaced fictional rows (`corestr.Collection`, `coreonce.New.<Type>`, `corejson.Serialize`/`Deserialize` shorthand) with concrete entries pointing at real APIs.

### Changed
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — moved C-CVS-06..08 + D-CVS-20..25 to **Resolved** (9 entries); cleared the **Open drift findings** table; added Cycle 4 (closed) row at **100.0 %** verifiable on §06; updated targets to ✅ for the ≥95 % aggregate goal and ✅ for zero-contradictions.

## [spec-v0.25.0] — 2026-05-04 (Cycle 4 baseline — §06 data-structures audited)

### Added
- **`spec/07-code-vs-spec-audits/05-cycle4-data-structures.md`** — full Cycle 4 audit of `01-app/06-data-structures.md`. 20 claims: 5 ✅ / 6 ⚠️ / **3 ❌** / **6 ❓**. Verifiable score: **35.7 %** (5/14). Three contradictions: `corejson.Serialize.ToString` / `Deserialize.FromTo` examples don't compile (real API: `Serialize.ToBytesErr` / `Deserialize.BytesTo`); `coreonce.New.String(...)` namespace doesn't match real top-level constructors (`coreonce.NewAnyOnce` / `NewByteOnce`); §4's "never `encoding/json` directly" rule is violated by `inttype/Variant.go` (`json.Marshal` in `MarshalJSON`) and `inttype/all-constructors.go` (`*json.Number` parameter type). The high ❓ count (6) reflects that `coregeneric` and `corepayload` have **zero consumers** in `enum-v10`.
- **9 new drift findings** (C-CVS-06..08 + D-CVS-20..25): documented vs actual `corejson` API, `coreonce` constructor surface, `corestr` real exports (`Hashset` / `SimpleSlice` / `SimpleStringOnce` rather than a string-list collection), and the upstream-only status of `coregeneric` / `corepayload`.

### Changed
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — added Cycle 4 row, 9 new open findings, 3 new milestones; section progress **4/16**. ❌ contradiction count went 0 → 3 (all on §06).

## [spec-v0.24.0] — 2026-05-04 (Cycle 3 closed — §05 enum-system at 100 % verifiable)

### Fixed
- **`spec/01-app/05-enum-system.md`** — resolves C-CVS-03..05 + D-CVS-14..19 from Cycle 3 in a single pass. Lifts §05 from **47.1 % → 100 %** of its verifiable subset; **❌ contradiction count: 3 → 0**.
  - **§1 architecture diagram**: replaced "consts.go + vars.go + `<Type>.go`" 3-file layout with the 2-file layout (`<TypeName>.go` + `vars.go`) that every enum package actually uses; added a callout that the type is conventionally named `Variant` (64 / 71 packages) and the file is named after the type.
  - **New §4.1 "Sentinel-first rule"** (resolves **C-CVS-03** + **C-CVS-05**): reframed the "first const must be `Invalid`" rule as "**the first iota constant must occupy the zero value of the backing type**", with a sentinel-name table (`Invalid` / `Unspecified` / `Uninitialized` / `Default` / domain term) showing real packages for each form. Documented the signed-int exception (`InvalidIndex Variant = -1` in `inttype`) explicitly.
  - **§4.2 (was Step 1 + Step 3, resolves D-CVS-14 + D-CVS-15)**: collapsed into a single `Variant.go` recipe combining type + iota + full method set. Renamed every example from `Status` → `Variant` to mirror the actual repo convention.
  - **§4.3 (was Step 2, resolves C-CVS-04 + D-CVS-18)**: deleted the `core-v9/internal/reflectinternal` import and `reflectinternal.TypeName(Invalid)` call (forbidden cross-module `internal/`). Replaced with `enumimpl.New.BasicByte.DefaultAllCases(Invalid, Ranges[:])` as the recommended pattern, plus an alternate `UsingTypeSlice("Variant", Ranges[:])` fallback. Added an explicit warning callout.
  - **§4.5 "Predicate file-split guideline" (resolves D-CVS-19)**: softened the hard-rule (>6 OR >20 lines triggers split) to a guideline that matches `pathpatterntype/Variant.go` reality (113 constants, all predicates kept in one file).
  - **§6 Factory Method Reference (resolves D-CVS-16)**: expanded the table from 5 → 10 methods to cover the actually-used surface (`UsingFirstItemSliceAllCases`, `DefaultAllCases`, `DefaultWithAliasMapAllCases`, `UsingFirstItemSliceAliasMap`, `CreateUsingSlicePlusAliasMapOptions`, `CreateUsingStringersSpread`); dropped the unused `CreateUsingMap` row with an explanatory note. Updated "When to pick which" to put `DefaultAllCases` first.
  - **§8 Testing an Enum (resolves D-CVS-17)**: rewrote the "three test files per enum under `tests/integratedtests/<pkg>tests/`" pattern to the actual shared-registry approach under `tests/creationtests/` (registration in `allBasicEnumsCollection.go`, table-driven driver). Mirrors the C-CVS-01 fix already applied to §03.
  - **§9 Common Mistakes**: replaced the "First constant is not `Invalid`" row with a sentinel-aware version; added rows for the `internal/reflectinternal` mistake and the wrong test-folder mistake.

### Changed
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — moved C-CVS-03..05 + D-CVS-14..19 to **Resolved** (9 entries); cleared the **Open drift findings** table (now empty); added Cycle 3 (closed) row at **100.0 %** verifiable on §05; updated targets to ✅ for the ≥95 % aggregate goal and ✅ for zero-contradictions.

## [spec-v0.23.0] — 2026-05-04 (Cycle 3 baseline — §05 enum-system audited)

### Added
- **`spec/07-code-vs-spec-audits/04-cycle3-enum-system.md`** — full Cycle 3 audit of `01-app/05-enum-system.md`. 18 claims: 8 ✅ / 6 ⚠️ / **3 ❌** / 1 ❓. Verifiable score: **47.1 %** (8/17). Three real contradictions surfaced (C-CVS-03..05): the "first const = `Invalid`" rule is violated by 10 packages using alternate sentinel names (`Default`, `Unspecified`, `Uninitialized`, `InvalidIndex = -1`); the recipe imports `core-v9/internal/reflectinternal` which `enum-v10` cannot legally do across module boundaries; and `inttype.InvalidIndex Variant = -1` directly contradicts the "zero-value sentinel" wording.
- **6 new drift findings** (D-CVS-14..19): documented vs actual file layout (`Variant.go` everywhere, no `consts.go`), missing `*AllCases` factory family in §6 (10+ call sites), unused `CreateUsingMap` listed, stale `tests/integratedtests/` reference (mirrors C-CVS-01), unrunnable `reflectinternal.TypeName(...)` example, predicate file-split rule that's never enforced.

### Changed
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — added Cycle 3 row, 9 new open findings, 3 new milestones; section progress **3/16**. ❌ contradiction count went from 0 → 3 (all on §05).

## [spec-v0.22.0] — 2026-05-04 (Cycle 2 closed — §04 at 100% verifiable)

### Fixed
- **`spec/01-app/04-error-system.md`** — resolves D-CVS-06..13 from Cycle 2.
  - **§1 table**: added rows for `MustBeEmpty`, `RawErrCollection`, `ToError` / `ToString`, `MessageWithRef`, `RangeNotMeet`.
  - **§1.1**: extended `RawErrorType` examples with `FailedToConvertType`, `NotSupportedType`, `PathInvalidErrorType`, `FailedToExecuteType`, `ComparatorShouldBeWithinRangeType`; added footnote pointing to upstream `errcore/RawErrorType.go` for the exhaustive enumeration.
  - **§1.2**: added `ErrorRefOnly(ref)` and `CombineWithAnother(other)` constructor rows; expanded code block with examples for both, `MustBeEmpty`, and the comparison vs `HandleErr`.
  - **New §1.5 "Reference Helpers"**: documents `MessageWithRef` and `RangeNotMeet` with `onofftype/vars.go`, `promptclitype/vars.go`, and `internal/messages/messages.go` references.
  - **New §1.6 "Accumulating Errors"**: documents `RawErrCollection` accumulator pattern with `osdetect/windowsSystemDetailGenerator_windows.go` reference.
  - **New §1.7 "Conversion Helpers"**: documents `ToString` / `ToError` with `osdetect/vars.go` reference.
  - **§7 "When to Use Which API"**: added 9 new rows covering all newly-documented APIs.

### Changed
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — Cycle 2 closed: §04 verifiable score moved 27.3 → **100.0**. All 8 open drift findings (D-CVS-06..13) moved to "Resolved". 7 ❓ remain pending task **AB** (fetch upstream `core-v9` source).

## [spec-v0.21.0] — 2026-05-04 (Cycle 2 baseline — §04 error-system audited)

### Added
- **`spec/07-code-vs-spec-audits/03-cycle2-error-system.md`** — full Cycle 2 audit of `01-app/04-error-system.md`. 18 claims extracted: 3 ✅, 8 ⚠️ (spec is incomplete vs consumer usage), 7 ❓ (unverifiable without upstream `core-v9` source — sandbox lacks Go + the module is not vendored), 0 ❌. Verifiable-subset score: **27.3 %** (3/11). New ❓ bucket introduced so aspirational APIs aren't logged as contradictions before upstream source is available.
- **8 new open drift findings** (D-CVS-06..13) covering APIs `enum-v10` actively uses but the spec does not document: `MustBeEmpty`, `RawErrCollection`, `<RawErrorType>.ErrorRefOnly`, `<RawErrorType>.CombineWithAnother`, `MessageWithRef`, `RangeNotMeet`, `ToError`/`ToString`, plus 4 missing `RawErrorType` examples (`FailedToExecuteType`, `NotSupportedType`, `PathInvalidErrorType`, `ComparatorShouldBeWithinRangeType`).

### Changed
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — added Cycle 2 row, new ❓ column, 8 open findings, 2 new milestones (apply D-CVS-06..13 spec fixes; fetch upstream `core-v9` source as task **AB**). Section progress: **2/16** done.

## [spec-v0.20.0] — 2026-05-04 (Cycle 1 fully closed — §03 at 100%)

### Fixed
- **`spec/01-app/03-import-conventions.md` §3 "`internal/` access from tests"** — resolves **C-CVS-02**. Section was titled "Common `internal/` packages used by tests" and used `core-v9/internal/reflectinternal` as a live example, but `enum-v10` (this repo) imports zero `internal/` packages. Reframed as a forward-looking explanation that `internal/` access is a **same-module** rule, with the `reflectinternal` example explicitly attributed to the upstream `core-v9` repo's own tests, and a new **consumer-side note** stating that `enum-v10`-style consumers cannot import `core-v9/internal/...` because Go enforces `internal/` at module boundaries.
- **`spec/01-app/03-import-conventions.md` §4 "Test-Package Imports"** — resolves **C-CVS-01**. Old text claimed tests live at `tests/integratedtests/footests/`, which doesn't exist in this repo. Replaced with the actual `enum-v10` layout (`tests/creationtests/` — flat, single `package creationtests`, mix of `*_test.go` and shared fixture `.go` files) plus a cross-reference to the upstream `core-v9` per-suite layout (`tests/<suite>/footests/`). The four shared rules (separate package, normal imports of source, no cycles, same-module `internal/` only) apply to both layouts.
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — moved C-CVS-01 and C-CVS-02 from Open → Resolved; §03 score updated 83.3 → **100.0** (12/12). Open-findings list is now empty. Cycle 1 is closed.

### Verified
- `rg -n "integratedtests|core-v9" spec/01-app/03-import-conventions.md` → 0 hits.

---

## [spec-v0.19.0] — 2026-05-04 (Cycle 1 LOW-drift fixes applied)

### Fixed
- **`spec/01-app/03-import-conventions.md`** — applied 5 LOW-severity Cycle 1 findings:
  - **D-CVS-01** (line 4): `consumes core-v9 packages` → `consumes core-v9 packages`.
  - **D-CVS-02** (line 88): `path ends in core-v9` / `not corev8` → `core-v9` / `corev9`.
  - **D-CVS-03** (line 94): `For core-v9, this means:` → `For core-v9, this means:`.
  - **D-CVS-04** (line 121): reworded "rooted at the same `core-v9` module" to a module-generic statement that applies equally to `core-v9`, `enum-v10`, or any other consumer with its own `internal/` tree.
  - **D-CVS-05** (lines 61, 73): annotated `coredata/coregeneric` as `// optional` in the canonical import block and added a sentence noting "not every consumer uses every package — `enum-v10`, for example, currently uses 8 of the 11 listed canonical imports".
- **`spec/07-code-vs-spec-audits/01-scoreboard.md`** — moved D-CVS-01..05 from Open → Resolved; §03 score updated 41.7 → **83.3** (10/12). The two remaining `tests/integratedtests/` and `internal/reflectinternal` contradictions stay open pending a structural decision.

### Verified
- `rg -n "core-v9" spec/01-app/03-import-conventions.md` → 0 hits (clean).

---

## [spec-v0.18.0] — 2026-05-04 (open code-vs-spec audit phase)

### Added
- **`spec/07-code-vs-spec-audits/`** — NEW folder kicking off the implementation-vs-documentation drift audit phase. The spec-only cycle reached its ceiling (97.5–98.0 / 100, 0 open findings) at v0.17.1; this is the inverse: does the code match what the docs claim?
  - `00-index.md` — scope rules + methodology (extract claims → locate evidence → verify ✅/⚠️/❌ → file findings).
  - `01-scoreboard.md` — living drift scoreboard.
  - `02-cycle1-import-conventions.md` — first cycle, audited `01-app/03-import-conventions.md` (12 claims).

### Findings opened (Cycle 1)
- 5 LOW drifts (D-CVS-01..05): stale `core-v9` prose in §03 that didn't get rewritten during the v8→v9 module rename.
- 1 MED + 1 HIGH contradiction (C-CVS-01..02): §03 references a `tests/integratedtests/` directory that doesn't exist in this repo and an `internal/reflectinternal` import nobody uses.
- Initial measured drift score: **41.7%** (5/12 matches). After applying the 5 LOW spec fixes this jumps to ~91.7%.

---

## [spec-v0.17.1] — 2026-04-25 (close F-V16-01 — trust-boundary worked example)

### Added
- `spec/01-app/15-observability.md` §5.1 — End-to-end "trust-boundary handler" worked example (~75 lines including rationale table). Shows inbound untrusted form data → `corevalidator.New.Line` validation → PII-redacted struct copy → structured `slog` log call → generic error returned to caller. Cross-references `08-validators.md`, `16-security.md`, and the existing §5 logging rules.

### Resolved
- **F-V16-01** — Trust-boundary worked example missing in `15-observability.md` §5. The 5th-party mediocre-AI simulation (audit #10) flagged that a junior agent attempting an end-to-end "validate untrusted input + log with PII redaction" task had to assemble pieces from 3 files. Now unified as one canonical example. Projected score uplift: 0.3–0.5 → ~98.0.

### Status
- **0 open / 27 resolved findings.** Spec is at the 97–98 ceiling with no remaining bug-fix work. The natural next phase is **code-vs-spec** (lifts the spec-only directive).

---

## [spec-v0.17.0] — 2026-04-25 (5th-party mediocre-AI reproducibility audit)

### Added
- `spec/99-audits/10-ai-audit-2026-04-25-mediocre-sim.md` — 5th-party simulation audit (mediocre/junior-AI persona, the user's original benchmark).

### Measured
- **MEASURED score: 97.5 / 100** (5th-party mediocre-AI simulation, post spec-v0.16.0). Confirms the v0.16.0 projection range of 97.5–97.8. Score progression across 5 audits: 85.5 → 88.9 → 96.0 → **97.5** (converged inside the 97–98 ceiling).
- 4-task implementation simulation: 3 ✅ PASS + 1 ⚠️ PARTIAL.
- v0.16.0 expansion (`15-observability.md`, `16-security.md`) verified Adequate. Zero regressions.

### Opened (OPTIONAL — spec is at-ceiling regardless)
- **F-V16-01** — `spec/01-app/15-observability.md` §5: trust-boundary worked example missing. A junior agent attempting an end-to-end "validate untrusted input + log with PII redaction" task finds all the rules across `08-validators.md`, `15-observability.md`, `16-security.md` but no single 25-line example unifies them. Severity: very low. Uplift: 0.3–0.5 → ~98.0 if landed.

### Audit-doc updates
- `spec/99-audits/00-index.md` — current-quality block updated to MEASURED 97.5.
- `spec/99-audits/07-scoreboard.md` — 5th-party row added; targets table updated; decision block confirms cycle complete.

### Status
- **Spec-only cycle COMPLETE.** 5 audit passes across 3 distinct model families have converged. The natural next phase is **code-vs-spec** (lifts the spec-only directive).

---

## [spec-v0.16.0] — 2026-04-25 (close F-V14-05 — observability + security expansion)

### Added (new content — not a fix)
- **`spec/01-app/15-observability.md`** (~165 lines) — Diagnostic primitives (`VarTwo`, `VarTwoNoType`, `MessageVarMap`), stack-aware error wrapping (`StackEnhance`), test-failure output format pointer, logging boundaries, tracing/metrics scope, common mistakes. Unifies guidance previously scattered across `04-error-system.md` §5, `13-testing-patterns.md` §8, and `00-llm-integration-guide.md` §StackEnhance.
- **`spec/01-app/16-security.md`** (~175 lines) — Threat model, PII/secret handling, panic policy, allocation/resource safety, reflection safety, input-validation boundary, common mistakes. Unifies guidance previously scattered across `04-error-system.md`, `00-llm-integration-guide.md`, and `02-design-philosophy.md`.

### Edited
- `spec/01-app/README.md` — reading order extended to 17 entries; per-file status table appends rows 15 + 16; expansion note added.

### Resolved
- **F-V14-05** — observability/security split into dedicated spec files. Closes the last open audit finding.

### Audit-doc updates
- `spec/99-audits/09-ai-audit-2026-04-25-claude-sonnet-4.5.md` — banner now reads "ALL F-V14 findings RESOLVED".
- `spec/99-audits/07-scoreboard.md` — open register cleared (0 open / 26 resolved); resolved register includes F-V14-05.
- `spec/99-audits/00-index.md` — current-quality block: 26 resolved / 0 open.

### Score impact
- Projected score: 96.0 → **~97.5–97.8** (at the 97–98 ceiling). All 26 findings across 4 audit cycles now RESOLVED.
- A 5th-party measurement audit is OPTIONAL — convergence pattern shows diminishing returns past audit #09.

---

## [spec-v0.15.1] — 2026-04-25 (close 4 of 5 F-V14 findings)

### Fixed (documentation only — zero implementation changes)
- **F-V14-01** — `spec/01-app/08-validators.md` §2.1: added inline comment after the fluent-builder example explicitly stating `Build()` returns `*LineValidator` (pointer; satisfies `IsSuccessValidator`, `MessageGetter`, `ErrorGetter`). Notes that pointer storage is required because receiver methods are pointer-bound.
- **F-V14-02** — `spec/01-app/10-reflection-and-dynamic.md` §2.1: rewrote the `InvokeMethod` example to (a) state the full signature `InvokeMethod(target any, name string, args ...any) (any, error)` in a comment, (b) explicitly note "no single-return overload exists", and (c) use `_, _ = coredynamic.InvokeMethod(...)` for the discard case so the two-return shape is visible at first glance.
- **F-V14-03** — `spec/01-app/09-converters.md` §1.1: replaced the under-specified `"true"/"false"/"1"/"0"` note with the full `strconv.ParseBool` accepted-set (case-sensitive: `1`, `t`, `T`, `TRUE`, `true`, `True`, `0`, `f`, `F`, `FALSE`, `false`, `False`); explicit rejection of whitespace and `yes`/`no`/`on`/`off`.
- **F-V14-04** — `spec/01-app/12-cmd-entrypoints.md` §1: appended a paragraph to the "no `cmd/`" rule clarifying it is **PR-review-enforced** (no `go vet` lint, no CI check, no PowerShell guard exists). Pointer to `spec/02-app-issues/` for anyone wanting machine enforcement.

### Deferred
- **F-V14-05** — observability/security split into dedicated spec files. Genuinely new content (not a fix); deferred to **spec-v0.16.0+**.

### Audit-doc updates
- `spec/99-audits/09-ai-audit-2026-04-25-claude-sonnet-4.5.md` — banner marking F-V14-01..04 RESOLVED; F-V14-05 deferred.
- `spec/99-audits/07-scoreboard.md` — open register now lists only F-V14-05; resolved register expanded with the 4 new closures.
- `spec/99-audits/00-index.md` — current-quality block: 25 resolved / 1 open.

### Projected score impact
- 96.0 → **~97.5** (sum of per-finding uplifts: 0.4 + 0.3 + 0.3 + 0.2 = 1.2). At the GPT-5-stated ceiling of 97–98.
- A 5th-party measurement is **optional** — convergence pattern shows diminishing returns past audit #09.

---

## [spec-v0.15.0] — 2026-04-25 (4th-party confirmation audit by Claude Sonnet 4.5)

### Added
- `spec/99-audits/09-ai-audit-2026-04-25-claude-sonnet-4.5.md` — full 4th-party audit report.

### Measured
- **MEASURED score: 96.0 / 100** (Claude Sonnet 4.5, 4th-party, post spec-v0.14.0). Up from GPT-5's 88.9 pre-fix baseline. Just 1.5 points below the 97–98 ceiling.
- 9/9 of v0.14.0's F-V12 fixes verified Adequate. Zero regressions.

### Opened
- **F-V14-01** — `01-app/08-validators.md` §2.1: `LineValidator.Build()` return type unstated (low / S).
- **F-V14-02** — `01-app/10-reflection-and-dynamic.md` §2.1: `coredynamic.InvokeMethod` example discards both returns, contradicting `(any, error)` signature (low / S).
- **F-V14-03** — `01-app/09-converters.md` §1.1: `StringTo.Bool` accepted-input set under-specified (low / S).
- **F-V14-04** — `01-app/12-cmd-entrypoints.md` §1: `cmd/` rule has no documented enforcement mechanism (very low / S).
- **F-V14-05** — DEFERRED to v0.16.0+: no dedicated `15-observability.md` / `16-security.md` (very low / M).

### Audit-doc updates
- `spec/99-audits/00-index.md` — current-quality block updated to MEASURED 96.0.
- `spec/99-audits/07-scoreboard.md` — 4th-party row added; open register lists F-V14-01..05.

### Projected score impact (after F-V14-01..04 closed in v0.15.1)
- 96.0 → ~97.5 (at the GPT-5-stated ceiling).

---

## [spec-v0.14.0] — 2026-04-25 (close all 9 F-V12 findings)

### Fixed (documentation only — zero implementation changes)
- **F-V12-01** — `spec/00-llm-integration-guide.md` §Pattern 7: rewrote suffix-order grammar to the canonical 8-slot order **Base + Filter + Result + Type + Lock + If + Must + Ptr**. Added a worked right-vs-wrong diff example.
- **F-V12-02** — `spec/01-app/01-package-map.md` + `spec/01-app/00-repo-overview.md` §2: appended ⚠️ **LEGACY** label to `coremath` entries with cross-link to integration guide. (Closes the F-NEW-01 ⚠️ Partial gap GPT-5 surfaced.)
- **F-V12-03** — `spec/00-llm-integration-guide.md` §coregeneric: added "Iterator types — import & contract" box with explicit `import "iter"`, full Seq/Seq2 signatures for every container, and a compile-ready usage example.
- **F-V12-04** — `spec/00-llm-integration-guide.md` §Code Style Rules: refined the Constants row to permit direct `""` / `' '` comparisons in fast paths and predicate guards.
- **F-V12-05** — `spec/01-app/04-error-system.md` §1: added `HandleErr(err error)` to the public-API table with signature and behavior; added a runnable example showing the contract (no-op when nil, panics with stack-enhanced wrap when non-nil).
- **F-V12-06** — `spec/01-app/05-enum-system.md` Step 3: added an explicit spelling note distinguishing stdlib `UnmarshalJSON` (one `l`) from the enum engine's `UnmarshallEnumToValue` family (two `l`s) — this is a deliberate historical convention; do NOT "correct" it.
- **F-V12-07** — `spec/00-llm-integration-guide.md` §coregeneric: added "Return-type conventions" boxes to both `Hashset[T]` and `Hashmap[K,V]` clarifying why `*Lock` variants return the receiver while non-`Lock` `*Bool` variants return the boolean signal. Also clarified `Hashset.AddBool` returns `true` if the key ALREADY existed (counterintuitive — explicitly called out).
- **F-V12-08** — **False positive.** `spec/06-testing-guidelines/07-diagnostics-output-standards.md` already existed (78 lines, substantive content covering stack-trace rules, single/multi-line assertion format, alignment). It was simply missing from the audit bundle sent to GPT-5. Verified the link from `01-app/13-testing-patterns.md §8` resolves correctly. No file change needed.
- **F-V12-09** — `spec/01-app/05-enum-system.md` §2: added `BasicEnumValuer` to the Key Interfaces table with cross-reference to Step 3 / F-NEW-07.

### Audit-doc updates
- `spec/99-audits/08-ai-audit-2026-04-25-gpt5-confirm.md` — banner marking all 9 RESOLVED.
- `spec/99-audits/07-scoreboard.md` — open register cleared; resolved register expanded; current state = **0 open / 21 resolved**.
- `spec/99-audits/00-index.md` — current quality block updated.

### Projected score impact
- 88.9 → ~99.4 (sum of per-finding uplifts +10.5), **capped at GPT-5's stated practical ceiling ~97–98**.
- Awaits a 4th-party (Claude Opus 4.5) re-audit to *measure* the new score.

---

## [spec-v0.13.0] — 2026-04-25 (3rd-party confirmation audit)

### Added
- **`spec/99-audits/08-ai-audit-2026-04-25-gpt5-confirm.md`** — independent re-audit by GPT-5 (different model family from prior Gemini passes).

### Findings
- **Measured score: 88.9 / 100** (vs ~94.5 self-projected — overestimate corrected).
- **6 of 7 F-NEW fixes verified ✅ Adequate**; F-NEW-01 marked ⚠️ Partial because `01-app/01-package-map.md` and `01-app/00-repo-overview.md §2` still list `coremath` as active utility (the LEGACY label only landed in the integration guide).
- **9 NEW gaps opened** (F-V12-01..F-V12-09) — full register in `99-audits/07-scoreboard.md`.

### Worst axis
- **Consistency = 82/100** — `coremath` leak across files, suffix-grammar Result-slot omission, Unmarshal/Unmarshall spelling inconsistency.

### Updated
- `spec/99-audits/00-index.md` — table extended; current-quality block updated to 88.9 measured.
- `spec/99-audits/07-scoreboard.md` — score history extended; F-NEW marked RESOLVED with verification verdicts; F-V12 OPEN register added.

### Decision
- Owner directive still active: **spec-only cycle**. All 9 F-V12 fixes will be documentation edits — no Go code, no tests, no PowerShell.

---

## [spec-v0.12.0] — 2026-04-25 (close all 7 F-NEW findings)

### Fixed (documentation only — zero implementation changes)
- **F-NEW-01** — Marked `coremath` as LEGACY in `spec/00-llm-integration-guide.md`; instructs new code to use Go 1.21+ built-in `min` / `max` / `clear`.
- **F-NEW-02** — Added explicit *tokenization rule* + worked example to "Compound `*Or*` Naming in Filter Chains": `AOrB` is a single naming token occupying one slot in the suffix grammar.
- **F-NEW-03** — Added decision matrix at the top of `## coregeneric — Generic Data Structures API Reference` clarifying when to use `Collection[T]` (default, mutex-protected) vs `SimpleSlice[T]` (function-local, perf-sensitive).
- **F-NEW-04** — Elevated the `core` package name vs `core-v9` module path warning to top-level Module Identity blocks in both `spec/00-llm-integration-guide.md` and `spec/01-app/00-repo-overview.md`.
- **F-NEW-05** — Added "Serialization Asymmetry" subsection to `spec/01-app/05-enum-system.md` §7, documenting that `MarshalJSON` MUST emit string names while `UnmarshalJSON` MAY accept names/aliases/numeric strings.
- **F-NEW-06** — Added rationale block to "The CaseV1 Cast Idiom" in `spec/06-testing-guidelines/02-test-case-types.md` explaining why the cast keeps `BaseTestCase` decoupled from assertion libraries.
- **F-NEW-07** — Added comment + rationale to `spec/01-app/05-enum-system.md` §Step 3 explaining all `Value<Type>()` methods are required by `enuminf.BasicEnumValuer` and MUST NOT be deleted.

### Audit-doc updates
- `spec/99-audits/06-ai-audit-2026-04-25-gemini-2.5-pro.md` — banner added marking all 7 findings RESOLVED.
- `spec/99-audits/07-scoreboard.md` — open findings register cleared; resolved register expanded; projected score raised from 85.5 → ~94.5 pending an independent re-audit for confirmation.

### Projected score impact
- 85.5 → ~94.5 (+9.0 weighted points, matches Gemini 2.5 Pro's per-finding uplift sum).
- Subscore expected gains: Unambiguity 75 → ~88; Self-containment 80 → ~90; Consistency 80 → ~90; Completeness 85 → ~90.

---

## [spec-v0.11.0] — 2026-04-25 (audit folder reorganization + new audit)

### Added
- **`spec/99-audits/`** — new folder consolidating every spec-quality audit. All previous `spec/99-*.md` files renamed and moved with serialized prefixes:
  - `99-audit-report.md` → `99-audits/01-original-11-step-plan.md`
  - `99-step11-simulation.md` → `99-audits/02-step11-simulation-cycle1.md`
  - `99-step11-simulation-cycle2.md` → `99-audits/03-step11-simulation-cycle2.md`
  - `99-step11-simulation-cycle3.md` → `99-audits/04-step11-simulation-cycle3.md`
  - `99-ai-audit-2026-04-23.md` → `99-audits/05-ai-audit-2026-04-23-gemini.md`
  - `99-ai-audit-2026-04-25.md` → `99-audits/06-ai-audit-2026-04-25-gemini-2.5-pro.md`
- **`spec/99-audits/00-index.md`** — folder index with reading order, scope rules, and process for adding new audits.
- **`spec/99-audits/07-scoreboard.md`** — living scoreboard.

### Changed
- All cross-references throughout `spec/`, `readme.txt`, and module docs updated to the new `99-audits/NN-…` paths.

### Decision
- **Scope is spec-only.** F-NEW-01..07 will be addressed by editing documentation only.

---

## [spec-v0.5.0] — 2026-04-23 (audit retirement)

### Added
- **`spec/99-audits/`** — new folder consolidating every spec-quality audit. All previous `spec/99-*.md` files renamed and moved with serialized prefixes:
  - `99-audit-report.md` → `99-audits/01-original-11-step-plan.md`
  - `99-step11-simulation.md` → `99-audits/02-step11-simulation-cycle1.md`
  - `99-step11-simulation-cycle2.md` → `99-audits/03-step11-simulation-cycle2.md`
  - `99-step11-simulation-cycle3.md` → `99-audits/04-step11-simulation-cycle3.md`
  - `99-ai-audit-2026-04-23.md` → `99-audits/05-ai-audit-2026-04-23-gemini.md`
  - `99-ai-audit-2026-04-25.md` → `99-audits/06-ai-audit-2026-04-25-gemini-2.5-pro.md`
- **`spec/99-audits/00-index.md`** — folder index with reading order, scope rules, and process for adding new audits.
- **`spec/99-audits/07-scoreboard.md`** — living scoreboard: current score 85.5/100, open findings register (F-NEW-01..07), targets, and explicit "spec-only, no implementation" decision.
- **New independent audit (Gemini 2.5 Pro)** opened **7 new findings** (F-NEW-01..F-NEW-07), all OPEN. Honest score correction: previously self-reported ~96.8 → independently measured **85.5**.

### Changed
- All cross-references throughout `spec/`, `readme.txt`, and module docs updated to the new `99-audits/NN-…` paths.

### Decision
- **Scope is spec-only.** F-NEW-01..07 will be addressed by editing documentation (`00-llm-integration-guide.md`, `05-enum-system.md`, `02-test-case-types.md`). No Go code, tests, or PowerShell will be touched in this cycle.

---

## [spec-v0.5.0] — 2026-04-23 (audit retirement)

### Closed
- **#02** internal-package coverage policy → **wont-fix** (no maintainer; contradiction documented in both surfaces; reopening criteria preserved).
- **#03** `GetAssert` stability declaration → **wont-fix** (discoverability solved; interim "in-tree only" policy works; reopening criteria preserved).
- **#04** `testwrappers` stability declaration → **wont-fix** (inventory complete; reopening criteria preserved).

### State
- **All 9 issues now closed**: 6 resolved with spec edits, 3 wont-fix with explicit reopening criteria.
- **Audit cycle formally retired** — no remaining self-actionable or maintainer-blocked work in this cycle.
- **Spec is production-ready** at v0.5.0 with 97.17% three-cycle reproducibility average.

### Reopening any wont-fix issue
Each wont-fix issue file (#02, #03, #04) contains an explicit "Reopening criteria" section. To reopen: flip the issue file's `Status:` line back to `open`, update `00-issues-index.md`, and bump spec to next minor version.

---

## [spec-v0.4.0] — 2026-04-23 (Cycle 3 audit)

### Added
- **`spec/99-audits/04-step11-simulation-cycle3.md`** — third fresh-AI simulation, enum scenario.
- **`spec/01-app/05-enum-system.md` §4 "Predicate file-split rule"** — explicit threshold (>6 predicates or >20 lines each → separate files).
- **`spec/02-app-issues/09-enum-predicate-file-split.md`** — issue filed and resolved in same cycle.

### Reproducibility
- **97.50%** on Cycle 3 enum scenario.
- **Three-cycle moving average: 97.17%** (C1 96.25% + C2 97.75% + C3 97.50%).
- Confirms **convergence above ≥95% target** with no regression from Followup edits.

### Resolved (issues)
- **#09** Enum predicate file-split rule — resolved by §4 addition.

---

## [spec-v0.3.0] — 2026-04-23 (Cycle 2 audit)

### Added
- **`spec/99-audits/03-step11-simulation-cycle2.md`** — second fresh-AI simulation, converter scenario.
- **`spec/01-app/04-error-system.md` §Boundary Cases** — 6-row table + decision tree distinguishing `FailedToConvertType` from `ValidationFailedType`.
- **`spec/02-app-issues/08-errcore-type-boundary-examples.md`** — issue filed and resolved in same cycle.

### Reproducibility
- **97.75%** on Cycle 2 converter scenario (up +1.5pp from Cycle 1's 96.25%).
- **Combined two-cycle average ≈ 97.0%** — well above the ≥95% target.

### Resolved (issues)
- **#08** `errcore` type boundary examples — resolved by new §"Boundary Cases".

---

## [spec-v0.2.0] — 2026-04-23 (Asia/Kuala_Lumpur, UTC+8)

### Added
- **`spec/01-app/`** — 14 new atomic per-topic architectural files (~3,400 lines total) covering: repo overview, package map, design philosophy, import conventions, error system, enum system, data structures, conditionals/utilities, validators, converters, reflection/dynamic, versioning, cmd entrypoints, testing patterns, tests folder walkthrough.
- **`spec/02-app-issues/`** — issue tracker (7 issues filed; 4 resolved, 3 maintainer-blocked).
- **`spec/04-tooling/04-bootstrap-into-new-repo.md`** — bootstrap guide.
- **`spec/06-testing-guidelines/02-test-case-types.md`** — Style B section (149 lines).
- **`spec/06-testing-guidelines/03-args-reference.md`** — Style C section + empty-map gotcha + `params.go` convention.
- **`spec/99-audits/02-step11-simulation-cycle1.md`** — fresh-AI reproducibility report (96.25% PASS).
- **`spec/01-app/02-design-philosophy.md`** — explicit filename casing rule (PascalCase exported, camelCase unexported, no snake_case).
- **`spec/01-app/08-validators.md`** — canonical validator error message format with 4 worked examples.
- **`spec/01-app/00-llm-integration-guide.md`** — AI reading order + decision matrix (Step 10).

### Changed
- **Renumbered all of `spec/`** to strict `NN-name.md` convention with no duplicate prefixes.
- **Renamed** `00-audit-report.md` → `99-audits/01-original-11-step-plan.md` (parked at slot 99 as meta).
- **Renamed** `02-tooling/` → `04-tooling/`; `testing-guidelines/` → `06-testing-guidelines/`.
- **Promoted** all 14 `spec/01-app/` files from 📝 Draft → ✅ Done after Step 11 verification.

### Resolved (issues)
- **#01** Style B / Style A coexistence — documented in `06-testing-guidelines/02-test-case-types.md`.
- **#05** `params.go` convention — rule added: mandatory for new packages with > 3 cases.
- **#06** Validator error canonical example — added to `08-validators.md` §5.
- **#07** `newCreator.go` filename casing — added to `02-design-philosophy.md` Pillar 1.

### Open (maintainer-blocked)
- **#02** Internal-package coverage policy — needs maintainer governance call.
- **#03** `GetAssert` stability declaration — discoverability done; stability undecided.
- **#04** `testwrappers` stability declaration — inventory done; stability undecided.

### Reproducibility
- **96.25%** average across 4 axes (package layout, test scaffolding, error categorisation, diagnostic format) — exceeds the ≥95% target set in `99-audits/01-original-11-step-plan.md` §9 Step 11.

---

## [spec-v0.1.0] — pre-2026-04-23 (initial scaffold)

### Added
- Initial `spec/` skeleton with `00-audit-report.md` (later renamed to `99-`).
- `00-llm-integration-guide.md` baseline.
- `testing-guidelines/` with 9 numbered files (later renamed to `06-`).
- `02-tooling/` with 3 files (later renamed to `04-`).
- `03-powershell-test-run/` with 9 files.
- `05-failing-tests/` with 25 chronological fix logs.

---

## Version-bump policy (recap)

Per `spec/01-app/11-versioning.md` §3:

| Spec change category | Bump |
|---|---|
| Major restructure (folder renumbering, breaking link changes) | **major** |
| New topic file, new section, new issue resolution | **minor** |
| Typo fix, link repair, formatting | **none** required |

This release (`v0.2.0`) is a **minor** bump — no breaking changes to deep links, but substantial new content.
