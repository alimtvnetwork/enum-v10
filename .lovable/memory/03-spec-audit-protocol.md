# Spec Audit Protocol

## Purpose

Drive `spec/` to factual truth by auditing each section against `enum-v9`'s actual code, then closing drift via spec edits (preferred) or escalating to the user.

## Cycle workflow

For each spec file (e.g. `spec/01-app/15-observability.md`):

1. **Extract claims** — every normative assertion gets a row in a claim-by-claim table.
2. **Probe code** — run `rg` against `enum-v9` Go source (excluding `cross-repo/`, `tests/`, `cmd/`) for each named symbol or pattern.
3. **Score** — assign one of:
   - ✅ Match (verified against code or spec-internal cross-ref)
   - ⚠️ Drift (spec says X, code says Y; fix the spec)
   - ❌ Contradiction (spec contradicts itself or another spec file; HIGH priority)
   - ❓ Unverifiable (claim is about upstream `core-v9`; defer to Task AB)
4. **Write the cycle report** at `spec/07-code-vs-spec-audits/NN-cycleN-<topic>.md`.
5. **Apply spec fixes** for ⚠️ and ❌ findings (catalogue as D-CVS-XX / C-CVS-XX).
6. **Update the scoreboard** at `spec/07-code-vs-spec-audits/01-scoreboard.md`.

## Verifiability dimensions

- **Code-vs-spec** — does `enum-v9` source confirm or contradict the claim?
- **Spec-internal consistency** *(introduced Cycle 13)* — does the cross-reference resolve? Does any sibling spec file contradict this one? Is the rule self-consistent?
- **Banned-pattern absence** — is the file free of `tests/integratedtests/`, `enum-v1` (outside the mirror), `core-v9 → core-v9` mojibake, `.lovable/user-preferences line 8` citations?

A claim that is ❓ on code-vs-spec can still be ✅ on spec-internal consistency. Track both.

## Scoring formula

```
Score = ✅ / (✅ + ⚠️ + ❌)        # ❓ excluded from denominator (verifiable subset)
```

## Closure criteria

- A section is **closed at 100% verifiable** when ⚠️ = 0 and ❌ = 0 (regardless of ❓ count).
- A section is **baseline-only** when verifiable subset is empty (all claims are ❓). §07 and §09 currently sit here.
- A section is **closed at baseline with zero corrective edits** when first-pass scoring already shows ⚠️ = 0, ❌ = 0, and at least one ✅. §15 and §16 achieved this.

## Finding ID schemes

- `C-CVS-XX` — Contradiction-class findings (HIGH).
- `D-CVS-XX` — Drift-class findings (MED/LOW).
- IDs are global across all cycles, monotonically increasing.

## Key recurring patterns

- **`tests/integratedtests/`** appeared 9 times across `01-app/` (cycles 1, 3, 6, 8, 9, 10, 11, 12). Always corrected to `tests/creationtests/` (the real `enum-v9` layout) with cross-link to predecessors.
- **Upstream-only API surfaces** (`coregeneric`, `corepayload`, `corevalidator`, `coretests`) get explicit consumer-coverage callouts (D-CVS-25, D-CVS-38, D-CVS-42).
- **Mojibake `core-v9 → core-v9`** appeared in §11 §3:95 and §4:112 — fixed in Cycle 9 (C-CVS-09a/b) by rewriting as the historical `core-v8` → `core-v9` migration.
