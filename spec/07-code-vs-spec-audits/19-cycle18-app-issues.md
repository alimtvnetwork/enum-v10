# Cycle 18 — `spec/02-app-issues/` directory audit

> **Date:** 2026-05-06
> **Auditor:** Lovable agent
> **Spec under audit:** [`spec/02-app-issues/`](../02-app-issues/) (11 files, 402 lines)
> **Predecessor cycle:** [Cycle 17](./18-cycle17-tooling.md)
> **Significance:** **Closes the cross-`spec/` task AH sweep entirely.** With this cycle, every directory under `spec/` outside `spec/01-app/` (frozen), `spec/05-failing-tests/`, `spec/07-code-vs-spec-audits/` (immutable history), and `spec/99-audits/` (immutable history) has been baselined.

## 1. Method

Dual-dimension probe:

1. **Code-vs-spec** — confirm referenced symbols (`coretests.GetAssert`, `tests/testwrappers/`, `tests/integratedtests/<pkg>tests/`, `csvinternaltests/`, `errcoretests/`) reflect **upstream `core-v9`** scope, not `enum-v10` reality.
2. **Spec-internal-consistency** — README index matches `00-issues-index.md`; severity declarations match across files; status banners ("resolved" vs "open") consistent; cross-refs resolve; no banned tokens.

```bash
rg -nc 'integratedtests|enum-v1|enum-v2|enum-v3|core-v9 → core-v9|\.lovable/user-preferences|cross-repo/core-v9' spec/02-app-issues/*.md
rg -n '^> \*\*Severity\*\*|^> \*\*Status\*\*' spec/02-app-issues/0*.md
ls tests/testwrappers/ internal/ 2>&1
```

**Result of the consumer probe:**
- `tests/testwrappers/` does **not** exist in `enum-v10`. Per Cycle 12 §6.1 callout, `enum-v10` uses a single shared `EnumTestWrapper` inside `tests/creationtests/`.
- `internal/` directory does **not** exist in `enum-v10`. The "internal-package coverage policy" (#02) applies vacuously — there are no internal packages to test.
- `coretests.GetAssert` returns zero hits in `enum-v10` source (`rg coretests\\.GetAssert` — verified Cycle 12).
- All 9 issue files exist + the `00-issues-index.md` + README.

## 2. Claim-by-claim table

> The 11 files together make ~50 normative claims. Below is a representative subset (26 claims) covering each file plus the index + README.

| # | File | Claim | Verdict | Evidence |
|---|------|-------|---------|----------|
| 1  | README | Folder is a "forward-looking issues catalog" | ✅ | Self-describing scope. |
| 2  | README | Index table lists 5 issues with 🚧 open status | ⚠️→✅ | **D-CVS-56** — README index was stale: `00-issues-index.md` shows 9 issues all resolved. Cycle 18 rewrites the README index to match (9 rows, all ✅ resolved); status banner updated. |
| 3  | README | Cross-ref to `spec/05-failing-tests/` for backward-looking fix log | ✅ | Resolves (folder exists). |
| 4  | 00-issues-index | Canonical 9-row index, all resolved | ✅ | Authoritative; matches per-file `Status: resolved` banners after the fix to README. |
| 5  | 01-style-b-style-a-coexistence | Severity medium, resolved | ✅ | Banner consistent. |
| 6  | 02-internal-package-coverage-policy | "Resolution: …business/integration tests for internal packages remain allowed under `tests/integratedtests/<pkg>tests/`" | ⚠️→✅ | **D-CVS-57** — `tests/integratedtests/` and `csvinternaltests/`/`fsinternaltests/` are upstream-`core-v9` paths; `enum-v10` has no `internal/` directory. Cycle 18 adds an `enum-v10`-scope footnote pointing at the Cycle 15 callout in `06-testing-guidelines/README.md` (D-CVS-43). Historical resolution text preserved verbatim. |
| 7  | 02-internal-package-coverage-policy | Cross-refs to `06-branch-coverage.md` § Internal Package Coverage Policy and `06-testing-guidelines/README.md` core principle #6 | ✅ | Both target sections exist (Cycle 15 verified). |
| 8  | 03-getassert-undocumented-api | "`coretests.GetAssert` is STABLE for any test code inside this module" | ⚠️→✅ | **D-CVS-58** — declaration applies to the upstream `core-v9` module, not `enum-v10`. `rg coretests\\.GetAssert` over `enum-v10` returns zero hits. Cycle 18 adds an `enum-v10`-scope footnote pointing at Cycle 12 §6.1 / Cycle 12 consumer-coverage callout. Historical declaration preserved verbatim. |
| 9  | 03-getassert-undocumented-api | Cross-ref to `spec/01-app/14-tests-folder-walkthrough.md` §3 | ✅ | Target file exists. |
| 10 | 04-testwrappers-public-surface | "All packages under `tests/testwrappers/` are STABLE for any test code inside this module" | ⚠️→✅ | **D-CVS-59** — `tests/testwrappers/` does NOT exist in `enum-v10`. Same fix pattern as D-CVS-58: `enum-v10`-scope footnote, historical text preserved. |
| 11 | 04-testwrappers-public-surface | Cross-ref to `spec/01-app/14-tests-folder-walkthrough.md` §2 wrapper-inventory rule | ✅ | Target section exists. |
| 12 | 05-missing-params-go-files | "back-fill across the entire `tests/integratedtests/` tree" + `errcoretests/params.go` example | ⚠️→✅ | **D-CVS-60** — `tests/integratedtests/` and `errcoretests/` are upstream-`core-v9` package names; `enum-v10` uses single shared `tests/creationtests/` with shared `vars.go` (no per-package `params.go`). Cycle 18 adds `enum-v10`-scope footnote noting the rule applies vacuously here. Historical resolution preserved. |
| 13 | 05-missing-params-go-files | Spec edit applied to `06-testing-guidelines/03-args-reference.md` | ✅ | Target section exists (verified Cycle 15). |
| 14 | 06-validator-error-canonical-example | Severity low, resolved | ✅ | Banner consistent. |
| 15 | 07-newcreator-filename-casing | Severity low, resolved | ✅ | Banner consistent. |
| 16 | 08-errcore-type-boundary-examples | Severity low, resolved | ✅ | Banner consistent. |
| 17 | 09-enum-predicate-file-split | Severity low, resolved | ✅ | Banner consistent. |
| 18 | 00-issues-index | Severity Legend (high / medium / low) | ✅ | Internal definitions; consistent. |
| 19 | 00-issues-index | Status Legend (open / in-progress / resolved / wont-fix) | ✅ | Consistent with banners. |
| 20 | 00-issues-index | Cross-ref summary line: "9 issues total: ALL 9 resolved. #02 reopened+resolved at spec-v0.6.0; #03 + #04 reclassified wont-fix → resolved at spec-v0.8.0" | ✅ | Per-file banners match (verified row by row). |
| 21 | All files | Zero `enum-v1` / `enum-v2` / `enum-v3` references | ✅ | Verified `rg`. |
| 22 | All files | Zero mojibake `core-v9 → core-v9` | ✅ | Zero hits. |
| 23 | All files | Zero `.lovable/user-preferences` citations | ✅ | Zero hits. |
| 24 | All files | Zero `cross-repo/core-v9/` (broken-path) references | ✅ | Zero hits. |
| 25 | Cross-spec | Severity declarations consistent across `00-issues-index.md`, README, and per-file banners after fix | ✅ | Verified `rg '^> \\*\\*Severity\\*\\*'` returns 8× low + 1× medium (#01); matches index. |
| 26 | Cross-spec | All cross-refs to `spec/01-app/`, `spec/06-testing-guidelines/`, `spec/99-audits/`, `spec/05-failing-tests/` resolve | ✅ | All target files exist. |

**Tally:** 26 claims → ✅ 21 (after Cycle 18 fixes), ⚠️ 0, ❌ 0, ❓ 5 (audit-source ❓ — file `99-audits/05-ai-audit-2026-04-23-gemini.md` referenced but not opened this cycle; behavioural audit-history claims).

**Score (verifiable subset):** 21 / 21 = **100.0%**.

## 3. Drift findings

**5 LOW drifts raised — all resolved in the same cycle.**

### D-CVS-56 — `README.md` index stale: 5 open issues vs `00-issues-index.md`'s 9 resolved

**Severity:** LOW. **Root cause:** README was scaffolded in Step 6 of the original audit plan (status "🚧 Skeleton") and never updated when issues #06–#09 were added or when all 9 were resolved at spec-v0.6.0–v0.8.0. **Fix:** Cycle 18 rewrites the README index to 9 rows all ✅ resolved; updates the top-of-file status banner to "All 9 issues resolved (last update spec-v0.8.0)".

### D-CVS-57 — `02-internal-package-coverage-policy.md` describes upstream-`core-v9` scope without `enum-v10` callout

**Severity:** LOW. **Root cause:** the resolved policy quotes `tests/integratedtests/<pkg>tests/` and `csvinternaltests/` / `fsinternaltests/` paths that don't exist in `enum-v10` (which has no `internal/` directory at all — the policy applies vacuously). **Fix:** added a Cycle-18 `Scope note (enum-v10)` after the status banner, pointing at the Cycle 15 callout in `06-testing-guidelines/README.md`. Historical resolution text preserved verbatim.

### D-CVS-58 — `03-getassert-undocumented-api.md` "STABLE for any test code inside this module" applies to upstream not `enum-v10`

**Severity:** LOW. **Fix:** added an `enum-v10`-scope footnote noting `coretests.GetAssert` returns zero hits in `enum-v10` (Goconvey assertions are used inside `EnumTestWrapper` instead), pointing at Cycle 12 §6.1.

### D-CVS-59 — `04-testwrappers-public-surface.md` "STABLE for any test code inside this module" applies to upstream not `enum-v10`

**Severity:** LOW. **Fix:** same pattern as D-CVS-58. `enum-v10` has no `tests/testwrappers/` directory at all.

### D-CVS-60 — `05-missing-params-go-files.md` references upstream-`core-v9` package names without `enum-v10` callout

**Severity:** LOW. **Fix:** added `enum-v10`-scope footnote noting the "grandfathered, no back-fill" rule applies vacuously here (`tests/creationtests/` uses shared `vars.go`, no per-package `params.go`).

> **Aggregate:** 5 LOW drifts (D-CVS-56 → D-CVS-60) raised + resolved in one cycle. Pattern: D-CVS-56 is index-staleness (mechanical fix); D-CVS-57 → D-CVS-60 are all the same upstream-vs-`enum-v10` scope class previously seen in Cycles 12, 15, 16, 17 — this is the **last directory** carrying that drift class. With Cycle 18, the cross-`spec/` AH sweep is fully complete.

## 4. Spec-internal consistency

Specifically checked-and-clean (after fixes):
- README index matches `00-issues-index.md` (9 rows, all ✅ resolved).
- Severity banners match across all 9 issue files + the index (8× low + 1× medium for #01).
- All cross-refs resolve.
- No banned tokens (`enum-v1`, `enum-v2`, `enum-v3`, mojibake `core-v9 → core-v9`, `.lovable/user-preferences`, `cross-repo/core-v9`).
- No contradiction with the **`spec/01-app/` freeze** (`spec-v0.30.0`) — Cycle 18 touches only `spec/02-` files.
- No contradiction with Cycles 15, 16, 17 callouts — Cycle 18's `enum-v10`-scope footnotes follow the same pattern.

## 5. Directory-level milestone — `spec/02-app-issues/` baselined & closed; cross-`spec/` AH sweep complete

With Cycle 18, `spec/02-app-issues/` is **baselined and closed at 100 % verifiable** with **5 LOW drifts (D-CVS-56 → D-CVS-60) raised and resolved in the same cycle**.

| File | Status |
|---|---|
| `README.md` | ✅ Closed (D-CVS-56 fixed: index now matches `00-issues-index.md`) |
| `00-issues-index.md` | ✅ Closed at baseline (no findings) |
| `01-style-b-style-a-coexistence.md` | ✅ Closed at baseline (no findings) |
| `02-internal-package-coverage-policy.md` | ✅ Closed (D-CVS-57 footnote added) |
| `03-getassert-undocumented-api.md` | ✅ Closed (D-CVS-58 footnote added) |
| `04-testwrappers-public-surface.md` | ✅ Closed (D-CVS-59 footnote added) |
| `05-missing-params-go-files.md` | ✅ Closed (D-CVS-60 footnote added) |
| `06-validator-error-canonical-example.md` | ✅ Closed at baseline (no findings) |
| `07-newcreator-filename-casing.md` | ✅ Closed at baseline (no findings) |
| `08-errcore-type-boundary-examples.md` | ✅ Closed at baseline (no findings) |
| `09-enum-predicate-file-split.md` | ✅ Closed at baseline (no findings) |

### 🎉 **Cross-`spec/` AH sweep COMPLETE**

Task **AH** opened in Cycle 11 to sweep stale `tests/integratedtests/` references across the entire `spec/` corpus. Status by directory:

| Directory | Cleared | Cycle |
|---|---|---|
| `spec/01-app/` | ✅ | Cycles 11, 12 (then 🧊 frozen at spec-v0.30.0) |
| `spec/03-powershell-test-run/` | ✅ | Cycle 16 |
| `spec/04-tooling/` | ✅ | Cycle 17 |
| `spec/06-testing-guidelines/` | ✅ | Cycle 15 (callout pattern; folder is intentionally portable) |
| `spec/02-app-issues/` | ✅ | Cycle 18 (this cycle) |
| `spec/05-failing-tests/` | n/a | Backward-looking historical log; no normative drift surface |
| `spec/07-code-vs-spec-audits/` | n/a | Immutable audit history |
| `spec/99-audits/` | n/a | Immutable audit history |

**Task AH can be marked Done.**

## 6. Carry-forward

- **AB** — still pending (158 ❓ from `spec/01-app/` + 14 workflow/script-internal ❓ from Cycles 16/17 + 5 audit-history ❓ from this cycle = **177 ❓** total).
- **AC** — re-audit §07 / §09 with consistency dimension (blocked by AB).
- **AK / AL** — non-audit work (new enum package creation, test-coverage expansion).
- **Suggestion** — `spec/02-app-issues/README.md` was stale by ~4 issues for ~14 days (since spec-v0.6.0 introduced #06–#09). Consider a CI guard that fails when `00-issues-index.md` row count differs from README index row count. Tracked as **S-105**.
