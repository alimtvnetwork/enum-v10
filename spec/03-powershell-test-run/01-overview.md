# PowerShell Test Runner (`run.ps1`)

> **Scope note (`enum-v9`)** — `run.ps1` is a generic Go-coverage runner that discovers test packages from disk (`go list ./tests/...`) and works with either layout: upstream `core-v9`'s `tests/integratedtests/<pkg>tests/` or `enum-v9`'s `tests/creationtests/`. **Example output and JSON snippets in this file use upstream-`core-v9` package names (`corecmptests`, `corejsontests`, etc.) for illustration.** When the same runner is invoked inside `enum-v9`, those paths render as `tests/creationtests/...` instead. The runner never hard-codes either folder name (per Core memory rule). See `spec/01-app/13-testing-patterns.md` §6.1 for the `enum-v9` test layout.

## Overview

`run.ps1` is the primary task runner for the project. It provides short, memorable commands for running tests, building, formatting, and more.

## Quick Reference

```powershell
./run.ps1 -t              # Run all tests
./run.ps1 -tc             # Run tests with coverage (parallel, HTML + summary)
./run.ps1 -tc --sync      # Run tests with coverage (sequential mode)
./run.ps1 -h              # Show help
```

## All Commands

| Short | Flag | Long | Description |
|-------|------|------|-------------|
| `T` | `-t` | `test` | Run all tests (verbose, with log output) |
| `TP` | `-tp` | `test-pkg` | Run tests for a specific package |
| `TC` | `-tc` | `test-cover` | Run tests with coverage report |
| `TI` | `-ti` | `test-int` | Run integrated tests only |
| `TF` | `-tf` | `test-fail` | Show last failing tests log |
| `GC` | `-gc` | `goconvey` | Launch GoConvey browser test runner |
| `R` | `-r` | `run` | Run the main application |
| `B` | `-b` | `build` | Build the binary |
| `BR` | `-br` | `build-run` | Build then run |
| `F` | `-f` | `fmt` | Format all Go files |
| `L` | `-l` | `lint` | Run `go vet` |
| `V` | `-v` | `vet` | Run `go vet` |
| `TY` | `-ty` | `tidy` | Run `go mod tidy` |
| `PC` | `-pc` | `pre-commit` | Check Coverage\* files for API mismatches |
| `C` | `-c` | `clean` | Clean build artifacts + coverage |
| `H` | `-h` | `help` | Show help |

## Usage Examples

```powershell
# Run all tests
./run.ps1 T
./run.ps1 -t
./run.ps1 test

# Run a specific package
./run.ps1 TP regexnewtests
./run.ps1 -tp corestrtests

# Run tests with coverage (parallel by default, auto-opens HTML report)
./run.ps1 TC
./run.ps1 -tc
./run.ps1 -tc --sync       # sequential mode
./run.ps1 -tc --no-open    # skip auto-open
./run.ps1 -tc --sync --no-open  # both flags

# Show last failing tests
./run.ps1 TF

# Launch GoConvey on custom port
./run.ps1 GC 9090

# Show help
./run.ps1 -h
./run.ps1 help
./run.ps1              # defaults to help
```

## Test Execution Pipeline

When you run `./run.ps1 -t`, the following happens in order:

```
1. git pull               (Invoke-FetchLatest)
2. go mod tidy            (dependency sync)
3. go build ./...         (Invoke-BuildCheck — fails fast if compilation errors)
4. go test -v -count=1    (run all tests, no caching)
5. Write-TestLogs         (parse output → passing/failing logs)
6. Open-FailingTestsIfAny (auto-open failing-tests.txt if failures exist)
```

### Build Check Gate

Before running tests, the script compiles the test packages. If compilation fails:
- Tests are **skipped entirely**
- Build errors are written to `data/test-logs/failing-tests.txt`
- The failing log is auto-opened

This prevents confusing test output when the code doesn't compile.

## Test Output & Logs

All test runs produce structured log files in `data/test-logs/`:

| File | Content |
|------|---------|
| `passing-tests.txt` | List of passing test names with count and timestamp |
| `failing-tests.txt` | Summary of failed tests + full diagnostic details |
| `raw-output.txt` | Complete unprocessed `go test` output |

### Failing Tests Log Format

```
# Failing Tests — 2026-03-11 10:30:00
# Count: 3

# ── Summary ──
  - TestFoo/Case_1
  - TestBar/Case_3
  - TestBaz/Case_0

# ── Details ──
FAIL: TestFoo/Case_1
  expected: "hello"
  actual:   "world"

FAIL: TestBar/Case_3
  ...
```

The summary section lists all failed test names sorted alphabetically, followed by detailed diagnostic output for each failure.

## Coverage Reports (`-tc`)

Running `./run.ps1 -tc` produces:

| File | Description |
|------|-------------|
| `data/coverage/coverage.out` | Raw Go coverage profile |
| `data/coverage/coverage.html` | Visual HTML report (auto-opens in browser) |
| `data/coverage/coverage-summary.txt` | Text summary with per-package and low-coverage highlights |
| `data/coverage/coverage-summary.json` | Machine-readable JSON (see schema below) |
| `data/coverage/per-package-coverage.txt` | Per-package coverage table |
| `data/coverage/per-package-coverage.json` | Per-package JSON |
| `data/coverage/blocked-packages.txt` | Text report of packages that failed to compile (full errors + stack traces) |
| `data/coverage/blocked-packages.json` | Machine-readable JSON (see schema below) |

### Console Output Sections

The TC command prints exactly **four sections** to the console (no per-package test rows):

1. **Build Failure Packages** — boxed list of packages that failed `go test -c`, or a single "all compiled" message
2. **Failing Test Summary** — boxed list of test functions that produced `--- FAIL:`, with pointer to `failing-tests.txt` for details
3. **Coverage Summary** — boxed per-source-package coverage table sorted by % descending, total line, and low-coverage warning
4. **Written Files Summary** — boxed list of all generated report file paths

Individual package compile/test results are **not** printed to the console — they are captured in log files only.

### Coverage Summary Contents

1. **Total Coverage** — aggregate percentage
2. **Per-Package Coverage** — breakdown by test package
3. **Low Coverage Functions (< 50%)** — functions needing attention
4. **Report file paths**

### JSON Schemas

#### `coverage-summary.json`

```json
{
  "timestamp": "2026-03-14T10:30:00Z",
  "totalCoverage": 62.5,
  "packageCount": 18,
  "packages": [
    {
      "package": "corestr",
      "coverage": 34.2,
      "statements": 120,
      "covered": 41,
      "uncovered": 79
    }
  ],
  "lowCoverageFuncCount": 5,
  "lowCoverageFunctions": [
    {
      "file": "github.com/user/core/corestr/Split.go",
      "function": "SplitByDelimiter",
      "coverage": 12.5
    }
  ],
  "blockedPackages": ["corecmptests", "isanytests"],
  "reports": {
    "profile": "data/coverage/coverage.out",
    "html": "data/coverage/coverage.html",
    "summary": "data/coverage/coverage-summary.txt",
    "json": "data/coverage/coverage-summary.json"
  }
}
```

Packages are sorted by coverage **ascending** so AI agents can iterate top-down to fix the worst gaps first.

#### `blocked-packages.json`

```json
{
  "timestamp": "2026-03-14T10:30:00Z",
  "blockedCount": 2,
  "compiledCount": 16,
  "totalCount": 18,
  "blockedPackages": [
    {
      "package": "corecmptests",
      "errorCount": 3,
      "errors": [
        "tests/integratedtests/corecmptests/Coverage5_test.go:14:2: too many arguments"
      ]
    }
  ]
}
```

Pass `--no-open` to skip auto-opening the HTML report:
```powershell
./run.ps1 -tc --no-open
```

## Command Dispatch

`run.ps1` is a thin dispatcher (~167 lines) that imports specialized `.psm1` modules from `scripts/` and routes commands to the appropriate function.

All three forms are equivalent and case-insensitive:

```powershell
./run.ps1 T       # uppercase shorthand
./run.ps1 -t      # hyphen-lowercase flag
./run.ps1 test    # long name
```

The dispatch table normalizes all forms via `$Command.ToLower()` and matches against a set of aliases.

### Module Architecture

| Module | Functions | Description |
|--------|-----------|-------------|
| `DashboardUI.psm1` | Phase tracking, coverage tables, box rendering | ANSI dashboard UI |
| `Utilities.psm1` | `Write-Header`, `Write-Success`, `ParseCompileErrors`, etc. | Common helpers |
| `TestLogWriter.psm1` | `Write-TestLogs` | Go test output → log files |
| `TestRunner.psm1` | `Invoke-AllTests`, `Invoke-PackageTests`, `Invoke-BuildCheck` | Test execution |
| `CoverageRunner.psm1` | `Invoke-TestCoverage`, `Invoke-PackageTestCoverage` | TC + TCP pipelines |
| `BuildTools.psm1` | `Invoke-Build`, `Invoke-Format`, `Invoke-Vet`, etc. | Build commands |
| `PreCommitCheck.psm1` | `Invoke-PreCommitCheck` | PC pipeline |
| `GoConvey.psm1` | `Invoke-GoConvey` | GoConvey launcher |
| `Help.psm1` | `Show-Help`, `Invoke-ShowFailLog`, `Invoke-IntegratedTests` | Help + misc |

See [`scripts/README.md`](/scripts/README.md) for full dependency graph and module documentation.

## Cleanup

```powershell
./run.ps1 -c      # removes: build/, tests/coverage.out, data/coverage/
```

## Related Docs

- [Repo Overview](/spec/01-app/00-repo-overview.md) — 🚧 skeleton (to be filled in audit Step 3)
- [CMD Entrypoints](/spec/01-app/12-cmd-entrypoints.md) — 🚧 skeleton (to be filled in audit Step 5)
- [Testing Patterns](/spec/01-app/13-testing-patterns.md) — 🚧 skeleton (to be filled in audit Step 4)
- [Spec Audit Report](/spec/99-audits/01-original-11-step-plan.md) — full gap analysis and step plan
