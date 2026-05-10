# Finding Registry

> All C-CVS-XX (contradictions) and D-CVS-XX (drifts) findings with current status.
> Authoritative source: `spec/07-code-vs-spec-audits/01-scoreboard.md`. This file is the human-browsable index.

## Contradictions (C-CVS-XX) — HIGH severity

| ID | Cycle | Section | Status | Title |
|---|---|---|---|---|
| C-CVS-01 | 1  | §03 | ✅ Resolved | `tests/integratedtests/` doesn't exist; correct path is `tests/creationtests/` |
| C-CVS-02 | 1  | §03 | ✅ Resolved | Conflicting import-order rules |
| C-CVS-03 | 3  | §05 | ✅ Resolved | Enum factory contradictions |
| C-CVS-04 | 3  | §05 | ✅ Resolved | (see cycle 3 report) |
| C-CVS-05 | 3  | §05 | ✅ Resolved | (see cycle 3 report) |
| C-CVS-06 | 4  | §06 | ✅ Resolved | "Never `encoding/json` directly" rule violated by `inttype` — documented exceptions added |
| C-CVS-07 | 4  | §06 | ✅ Resolved | `corejson.Serialize.ToString` example didn't compile |
| C-CVS-08 | 4  | §06 | ✅ Resolved | `corepayload.New.PayloadWrapper.UsingInstruction(...)` example unverifiable |
| C-CVS-09a | 9  | §11 | ✅ Resolved | Mojibake `core-v9 → core-v9` at §3:95 |
| C-CVS-09b | 9  | §11 | ✅ Resolved | Mojibake `core-v9 → core-v9` at §4:112 |
| C-CVS-10 | 10 | §12 | ✅ Resolved | "No `cmd/` directory" assertion contradicted by `enum-v9/cmd/main/main.go`; carved out smoke-test policy |

## Drifts (D-CVS-XX) — MED/LOW severity

| ID | Cycle | Section | Status | Pattern |
|---|---|---|---|---|
| D-CVS-01..05 | 1  | §03 | ✅ | Various import-convention nits |
| D-CVS-06..13 | 2  | §04 | ✅ | Error-system surface drift |
| D-CVS-14..19 | 3  | §05 | ✅ | Enum factory drift; D-CVS-17 = `tests/integratedtests/` |
| D-CVS-20..25 | 4  | §06 | ✅ | Data-structure surface drift; D-CVS-25 = consumer-coverage callout |
| D-CVS-26 | 6  | §08 | ✅ | `tests/integratedtests/` (3rd occurrence) |
| D-CVS-27 | 9  | §11 | ✅ | `tests/integratedtests/` at §4:108 (4th occurrence) |
| D-CVS-30 | 9  | §11 | ✅ | `versionindexes.V8 // 8 (current era — core-v9)` self-contradiction |
| D-CVS-31 | 9  | §11 | ✅ | 4 stale `.lovable/user-preferences line 8` citations |
| D-CVS-32 | 10 | §12 | ✅ | `tests/integratedtests/` at §3:71 (5th occurrence) |
| D-CVS-33 | 10 | §12 | ✅ | `go test ./tests/integratedtests/coregenerictests/...` example |
| D-CVS-35 | 11 | §12 | ❌ RETRACTED | False claim that `04-bootstrap-into-new-repo.md` was missing |
| D-CVS-36 | 11 | §13 | ✅ | `tests/integratedtests/footests/` (6th occurrence) — added §6.1 documenting real layout |
| D-CVS-37 | 11 | §13 | ✅ | Style D path correction |
| D-CVS-38 | 11 | §13 | ✅ | NEW consumer-coverage callout for upstream-only test API |
| D-CVS-39 | 12 | §14 | ✅ | `tests/integratedtests/` per-package layout (7th occurrence) |
| D-CVS-40 | 12 | §14 | ✅ | `tests/integratedtests/widgettests/` walkthrough (8th occurrence) |
| D-CVS-41 | 12 | §14 | ✅ | GetAssert observation-source path (9th occurrence) |
| D-CVS-42 | 12 | §14 | ✅ | NEW consumer-coverage callout (mirrors D-CVS-38) |
| D-CVS-43 | 12 | collateral | ✅ | Collateral: `01-package-map.md` §8 (5 hits) + `02-design-philosophy.md:183` |

## Open

_None._

## Reserved / skipped IDs

- D-CVS-14, 15, 16, 18, 19, 20, 21, 22, 23, 24, 28, 29, 34 — assigned across cycles; details in per-cycle reports. (Numbers 28, 29, 34 are gaps in the public scoreboard text — reserved for findings that were merged into adjacent IDs during cycle authoring.)
