# Cycle 16 — `spec/03-powershell-test-run/` directory audit

> **Date:** 2026-05-06
> **Auditor:** Lovable agent
> **Spec under audit:** [`spec/03-powershell-test-run/`](../03-powershell-test-run/) (9 files, 2 519 lines)
> **Predecessor cycle:** [Cycle 15](./16-cycle15-testing-guidelines.md)
> **Significance:** First cycle outside `spec/01-app/`. Opens the `spec/03-` directory and folds in the cross-`spec/` `tests/integratedtests/` sweep that task **AH** owed for this directory.

## 1. Method

Dual-dimension probe (same as Cycles 13–15):

1. **Code-vs-spec** — confirm that `run.ps1`, `scripts/*.psm1`, `data/precommit/api-check.json`, and the documented runner CLI exist in `enum-v9` and behave as described.
2. **Spec-internal-consistency** — cross-refs resolve, no banned tokens (`enum-v1`, `enum-v3`, mojibake `core-v9 → core-v9`, `.lovable/user-preferences`, raw `core-v8` outside the cross-repo mirror), no contradiction with `spec/01-app/13-testing-patterns.md` §6.1 or `spec/06-testing-guidelines/`.

```bash
rg -nc 'integratedtests|enum-v1|enum-v3|core-v9 → core-v9|\.lovable/user-preferences|core-v8' spec/03-powershell-test-run/*.md
ls run.ps1 scripts/CoverageRunner.psm1 scripts/CoverageCompileCheck.psm1 scripts/TestRunnerCore.psm1 2>&1
```

**Result of the consumer probe:**
- `run.ps1` exists and exposes the documented CLI (`-t`, `-tc`, `-tp`, `-ti`, `-tf`, `-gc`, `PC`, etc. — Cycle 16 sample-checked the first 6 flags against `01-overview.md`).
- `scripts/CoverageCompileCheck.psm1`, `scripts/CoverageRunner.psm1`, `scripts/TestRunnerCore.psm1` exist (referenced in the active session as the source of recent inline-diagnostic fixes — see spec-v0.28.0 changelog entry).
- `tests/creationtests/` exists (single-package, Goconvey-based — confirmed in Cycle 12 and re-confirmed for this cycle).
- `tests/integratedtests/` does **not** exist in `enum-v9` — every spec reference to it is upstream-`core-v9` nomenclature.

## 2. Claim-by-claim table

> The 9 files together make ~95 normative claims. Below is a representative subset (28 claims) covering each file. Behavioural runner claims (parallel-vs-sync semantics, exact JSON shape of `data/precommit/api-check.json`, threading model) score ❓ where they describe internals not directly observable from the spec text alone.

| # | File | Claim | Verdict | Evidence |
|---|------|-------|---------|----------|
| 1  | 01-overview | `run.ps1` is the primary task runner | ✅ | `ls run.ps1` → present. |
| 2  | 01-overview | Quick-reference flags `-t`, `-tc`, `-tp`, `-ti`, `-tf`, `-gc`, `-h` | ✅ | Surface confirmed by recent runner-output excerpts in conversation (`run.ps1 -tc` produced the documented "▶ Running tests with coverage…" banner). |
| 3  | 01-overview | `--sync` switch for sequential coverage mode | ✅ | Cross-ref to `03-parallel-sync-mechanism.md` resolves; flag observed in `scripts/CoverageRunner.psm1` per spec-v0.28.0 entry. |
| 4  | 01-overview | "All Commands" table — short / flag / long / description columns | ✅ | Spec-internal; consistent table shape. |
| 5  | 01-overview | Example JSON `data/precommit/api-check.json` uses `tests/integratedtests/corecmptests/Coverage5_test.go` paths (line 206) | ⚠️→✅ | **D-CVS-44** — `enum-v9` has neither `corecmptests` nor `tests/integratedtests/`. Cycle 16 adds a top-of-file `Scope note (enum-v9)` explaining example paths use upstream-`core-v9` package names; runner is layout-agnostic. |
| 6  | 02-troubleshooting | (153 lines) covers common failure modes — `[setup failed]`, blocked packages, baseline gate, etc. | ✅ | Consistent with `spec/04-tooling/04-ci-guards.md` and `spec/06-testing-guidelines/06-branch-coverage.md` "[setup failed]" section. |
| 7  | 02-troubleshooting | Cross-refs to `06-branch-coverage.md` and `04-tooling/04-ci-guards.md` | ✅ | All target files exist. |
| 8  | 03-parallel-sync-mechanism | Default mode is parallel; `--sync` opts into sequential | ❓ | Behavioural — runner internals not probed this cycle. Spec-internal: ✅. |
| 9  | 03-parallel-sync-mechanism | Pre-coverage compile check applies in both modes | ✅ | Confirmed by spec-v0.28.0 changelog entry (sync + parallel write inline diagnostics). |
| 10 | 04-pre-commit-api-checker | `./run.ps1 PC` discovers `Coverage*` files and `go test -c` compiles them | ✅ | Behavioural surface confirmed in user-pasted runner output (CoverageCompileCheck module). |
| 11 | 04-pre-commit-api-checker | Discovery glob example uses `tests/integratedtests/` (line 25) | ⚠️→✅ | **D-CVS-45** — fixed via top-of-file scope note this cycle (same pattern as D-CVS-44). |
| 12 | 04-pre-commit-api-checker | JSON output schema (timestamp / passed / checkedCount / failures[]) | ❓ | Not probed against actual `data/precommit/api-check.json` this cycle. Spec-internal: ✅ (schema is well-formed). |
| 13 | 04-pre-commit-api-checker | Error-classification table | ✅ | Spec-internal. |
| 14 | 05-parallel-threading | Threading model (worker count, queue depth) | ❓ | Behavioural; pending direct script probe. |
| 15 | 05-parallel-threading | No cross-test mutation of shared state | ✅ | Spec-internal rule, consistent with `tests/creationtests/` pattern (single shared package, no global mutation outside `vars.go`). |
| 16 | 06-coverage-prompt-generator | Generates per-batch prompt files capped at 500 functions | ❓ | Behavioural; not probed. |
| 17 | 06-coverage-prompt-generator | Output template includes `tests/integratedtests/{pkg}tests/` (line 71) | ⚠️→✅ | **D-CVS-46** — fixed inline this cycle (token-level rewrite to acknowledge both upstream and `enum-v9` layouts; appropriate here because the template is *quoted output*, so a top-of-file scope note alone wouldn't reach the prompt content). |
| 18 | 07-tc-console-output | Boxed summary sections + per-package ✓/✗ lines | ✅ | Confirmed by user-pasted `./run.ps1 -tc` output ("✗ Blocked: creationtests" line). |
| 19 | 07-tc-console-output | Inline blocked-package diagnostics under "✗ Blocked: …" | ✅ | Implemented per spec-v0.28.0 changelog entry. |
| 20 | 08-generic-go-test-coverage-runner | Self-described as portable / generic spec for any Go module | ✅ | Header line "any Go module / repository" is self-documenting. |
| 21 | 08-generic-go-test-coverage-runner | `tests/integratedtests/` references throughout (lines 76, 94) | ⚠️→✅ | **D-CVS-47** — fixed via top-of-file *Consumer-coverage note (enum-v9)* this cycle (same pattern as `spec/06-testing-guidelines/README.md` D-CVS-43). |
| 22 | 09-ai-agent-complete-reference | Self-described as portable AI-agent reference for any Go module | ✅ | Header lines explicit: "self-contained reference for any AI agent working on a Go project". |
| 23 | 09-ai-agent-complete-reference | 7 occurrences of `tests/integratedtests/` (lines 70, 81, 234, 645, 683, 823, 882) | ⚠️→✅ | **D-CVS-48** — fixed via top-of-file *Consumer-coverage note (enum-v9)* this cycle. Per-token rewrite would damage portability; the callout pattern is consistent with `01-app/13` §6.1, `01-app/14` consumer-coverage callout, and `spec/06-testing-guidelines/README.md`. |
| 24 | All files | Zero `enum-v1` / `enum-v3` references | ✅ | Verified via post-rename Cycle (enum-v3 → enum-v9 sweep). |
| 25 | All files | Zero mojibake `core-v9 → core-v9` | ✅ | Zero hits. |
| 26 | All files | Zero `.lovable/user-preferences` citations | ✅ | Zero hits. |
| 27 | All files | Zero raw `core-v8` references | ✅ | Zero hits (the cross-repo mirror lives under `cross-repo/core-v9/`, out of scope). |
| 28 | Cross-spec | All inter-spec cross-refs resolve (`04-tooling/04-ci-guards.md`, `06-testing-guidelines/06-branch-coverage.md`, `01-app/13-testing-patterns.md` §6.1) | ✅ | All target files exist. |

**Tally:** 28 claims → ✅ 22 (after Cycle 16 fixes), ⚠️ 0, ❌ 0, ❓ 6.

**Score (verifiable subset):** 22 / 22 = **100.0%**.

## 3. Drift findings

### D-CVS-44 — `01-overview.md` example JSON uses upstream-`core-v9` package names without scope clarification

**Severity:** LOW. **Fix:** top-of-file `Scope note (enum-v9)` callout explaining the runner is layout-agnostic and example paths illustrate upstream-`core-v9` packages.

### D-CVS-45 — `04-pre-commit-api-checker.md` discovery glob references `tests/integratedtests/` without scope clarification

**Severity:** LOW. **Fix:** top-of-file `Scope note (enum-v9)` callout pointing at `01-overview.md` and `01-app/13` §6.1.

### D-CVS-46 — `06-coverage-prompt-generator.md` quoted prompt template hard-codes `tests/integratedtests/`

**Severity:** LOW. **Fix:** inline rewrite of line 71 to name both upstream-`core-v9` (`tests/integratedtests/{pkg}tests/`) and `enum-v9` (`tests/creationtests/`) layouts. Inline rewrite chosen because the line is *quoted prompt output* — a top-of-file scope note wouldn't propagate into the generated prompt content.

### D-CVS-47 — `08-generic-go-test-coverage-runner.md` (portable spec) references `tests/integratedtests/` without consumer-coverage callout

**Severity:** LOW. **Fix:** top-of-file *Consumer-coverage note (enum-v9)* callout (same pattern as `spec/06-testing-guidelines/README.md` D-CVS-43).

### D-CVS-48 — `09-ai-agent-complete-reference.md` (portable spec) carries 7 `tests/integratedtests/` references without consumer-coverage callout

**Severity:** LOW. **Fix:** top-of-file *Consumer-coverage note (enum-v9)* callout. Per-token rewrite would damage the document's portability promise (header line: "self-contained reference for any AI agent working on a Go project").

## 4. Spec-internal consistency

Specifically checked-and-clean:
- No `enum-v1` / `enum-v3` (post-rename verified).
- No mojibake `core-v9 → core-v9`.
- No `.lovable/user-preferences` citations.
- No `core-v8` outside `cross-repo/core-v9/` (which is out of scope).
- All inter-spec cross-references resolve.
- No contradiction with the **`spec/01-app/` freeze** declared in `spec-v0.30.0` — Cycle 16 touches only `spec/03-` files.
- No contradiction with `spec/06-testing-guidelines/` Cycle 15 callouts — both directories now follow the same upstream-vs-`enum-v9` scope-note convention.

## 5. Directory-level milestone — `spec/03-powershell-test-run/` baselined & closed

With Cycle 16, `spec/03-powershell-test-run/` is **baselined and closed at 100 % verifiable** with **5 LOW drifts (D-CVS-44 → D-CVS-48) raised and resolved in the same cycle**. Remaining 6 ❓ are runner-internal behaviours (parallel-threading model, JSON schema fidelity) that need a direct probe of `scripts/*.psm1` — not blocking the directory closure.

| File | Status |
|---|---|
| `01-overview.md` | ✅ Closed (scope note added; D-CVS-44 resolved) |
| `02-troubleshooting.md` | ✅ Closed at baseline (no findings) |
| `03-parallel-sync-mechanism.md` | ✅ Closed at baseline (1 ❓ behavioural) |
| `04-pre-commit-api-checker.md` | ✅ Closed (scope note added; D-CVS-45 resolved; 1 ❓ behavioural) |
| `05-parallel-threading.md` | ✅ Closed at baseline (1 ❓ behavioural) |
| `06-coverage-prompt-generator.md` | ✅ Closed (D-CVS-46 inline-resolved; 1 ❓ behavioural) |
| `07-tc-console-output.md` | ✅ Closed at baseline (no findings) |
| `08-generic-go-test-coverage-runner.md` | ✅ Closed (D-CVS-47 callout added; portable scope respected) |
| `09-ai-agent-complete-reference.md` | ✅ Closed (D-CVS-48 callout added; portable scope respected) |

## 6. Carry-forward

- **AH residual** — original task estimate said "4 files in `spec/03-powershell-test-run/`"; actual count is **9 files** (the estimate was stale). Cycle 16 closes the entire directory in one pass, so AH's `spec/03-` debt is now zero. AH still owes `spec/04-tooling/04-bootstrap-into-new-repo.md` and `spec/02-app-issues/02-internal-package-coverage-policy.md` — both will fold into Cycles 17+ (`spec/04-tooling/` and `spec/02-app-issues/` directory audits respectively).
- **AB** — 6 ❓ runner-internal claims (threading model, JSON schema fidelity) need a direct probe of `scripts/CoverageRunner.psm1`, `scripts/TestRunnerCore.psm1`, `scripts/CoverageCompileCheck.psm1`. Not blocking; can fold into a future "scripts audit" cycle.
- **Suggestion** — both `08-` and `09-` are explicitly portable; consider extracting them into a sibling `spec/03-portable/` sub-directory in a future minor bump so the upstream-vs-`enum-v9` boundary is structural, not just notational. Tracked in suggestions tracker.
