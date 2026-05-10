# Coverage Prompt Generator

## Overview

After `./run.ps1 TC` completes, the runner automatically generates AI-friendly prompt files listing all functions below 100% coverage with their specific uncovered line ranges.

## Scripts

```
scripts/coverage/
├── Generate-CoveragePrompts.ps1      # Main generator (called by run.ps1)
├── Get-UncoveredLines.ps1            # Standalone: uncovered lines for one file
├── Get-FunctionCoverage.ps1          # Standalone: filter functions by threshold
└── Get-PackageCoverageReport.ps1     # Combined: detailed report per package
```

### Generate-CoveragePrompts.ps1

Main generator invoked automatically at the end of `Invoke-TestCoverage`. Parses both `coverage.out` and `go tool cover -func` output to produce batched prompt files with uncovered line ranges.

### Get-UncoveredLines.ps1

Standalone utility that extracts uncovered line ranges for a single source file from a coverage profile. Useful for debugging why a specific file has gaps.

| Parameter | Required | Description |
|-----------|----------|-------------|
| `CoverProfile` | Yes | Path to `coverage.out` |
| `SourceFile` | Yes | Full module-qualified file path (e.g., `github.com/alimtvnetwork/core-v9/errcore/ErrorNew.go`) |

### Get-FunctionCoverage.ps1

Standalone utility that filters `go tool cover -func` output to list functions below a given coverage threshold, sorted ascending by coverage.

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `FuncOutput` | Yes | — | Lines from `go tool cover -func` |
| `Threshold` | No | `100.0` | Coverage percentage threshold |

### Get-PackageCoverageReport.ps1

Combined utility that merges function filtering and uncovered-line extraction into a single detailed, color-coded report for one package. Shows each function's coverage percentage and specific uncovered line ranges.

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `CoverProfile` | Yes | — | Path to `coverage.out` |
| `FuncOutput` | Yes | — | Lines from `go tool cover -func` |
| `Package` | Yes | — | Package path (e.g., `errcore`, `internal/strutilinternal`) |
| `Threshold` | No | `100.0` | Coverage percentage threshold |
| `OutputFile` | No | — | Path to write a copy of the report |
| `Format` | No | `text` | Output format: `text` (color-coded), `markdown` (table), or `json` (machine-readable) |

## Output

```
data/prompts/
├── coverage-prompt-1.txt          # Functions 1-500 (sorted by coverage ascending)
├── coverage-prompt-2.txt          # Functions 501-1000
├── ...
└── prompts-summary.json           # Metadata (counts, batch info)
```

## Prompt File Format

```text
# Coverage Improvement Prompt — Batch 1/3
# Generated: 2026-03-15 12:00:00
# Functions: 500 (of 1200 total below 100%)

Please improve the code coverage to 100% for these functions.
Each function lists its current coverage and the specific uncovered lines.
Write tests in tests/integratedtests/{pkg}tests/ (upstream `core-v9`) — or `tests/creationtests/` for `enum-v10` — using the AAA pattern with args.Map + ShouldBeEqual.

─────────────────────────────────────────────────────────────

## NewError
   File:     errcore/ErrorNew.go
   Package:  errcore
   Coverage: 66.7%
   Uncovered lines: L15-L17, L22

## SplitLeftRight
   File:     internal/strutilinternal/all-left-right-splits.go
   Package:  internal/strutilinternal
   Coverage: 40.0%
   Uncovered lines: L8-L12, L18
```

## Parameters (Generate-CoveragePrompts.ps1)

| Parameter | Default | Description |
|-----------|---------|-------------|
| `CoverProfile` | (required) | Path to merged coverage.out |
| `FuncOutput` | (required) | Lines from `go tool cover -func` |
| `OutputDir` | `data/prompts` | Where to write prompt files |
| `BatchSize` | 500 | Functions per file |
| `ProjectRoot` | (auto-detect) | Project root for path resolution |

## Standalone Usage

```powershell
# Get uncovered lines for a specific file
./scripts/coverage/Get-UncoveredLines.ps1 `
  -CoverProfile data/coverage/coverage.out `
  -SourceFile "github.com/alimtvnetwork/core-v9/errcore/ErrorNew.go"

# Get all functions below 80% coverage
$funcLines = go tool cover -func=data/coverage/coverage.out
./scripts/coverage/Get-FunctionCoverage.ps1 -FuncOutput $funcLines -Threshold 80

# Get all functions below 100% (default threshold)
./scripts/coverage/Get-FunctionCoverage.ps1 -FuncOutput $funcLines

# Detailed report for a specific package (combines both utilities)
$funcLines = go tool cover -func=data/coverage/coverage.out
./scripts/coverage/Get-PackageCoverageReport.ps1 `
  -CoverProfile data/coverage/coverage.out `
  -FuncOutput $funcLines `
  -Package "errcore"

# Same, but only functions below 80%
./scripts/coverage/Get-PackageCoverageReport.ps1 `
  -CoverProfile data/coverage/coverage.out `
  -FuncOutput $funcLines `
  -Package "internal/strutilinternal" `
  -Threshold 80
```

## Integration Point

Called automatically at end of `Invoke-TestCoverage` in `run.ps1` (line ~1062):

```powershell
$promptScript = Join-Path $PSScriptRoot "scripts" "coverage" "Generate-CoveragePrompts.ps1"
if (Test-Path $promptScript) {
    & $promptScript -CoverProfile $coverProfile -FuncOutput $funcOutput `
      -OutputDir $promptsDir -BatchSize 500 -ProjectRoot $PSScriptRoot
}
```

## Related Docs

- [PowerShell Test Runner Overview](01-overview.md)
- [Parallel Threading Strategy](05-parallel-threading.md)
- [Coverage Prompts Spec](../13-app-issues/testing/09-coverage-prompts-from-powershell.md)
