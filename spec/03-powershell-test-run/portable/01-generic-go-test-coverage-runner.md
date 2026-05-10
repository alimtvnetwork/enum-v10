# Generic Go Test Coverage Runner — PowerShell Spec

> **Consumer-coverage note (`enum-v9`)** — this entire file is a **portable / generic** spec for building a Go-coverage runner from scratch in any Go module. References to `tests/integratedtests/` describe the **upstream `core-v9`** convention; the same runner deployed inside `enum-v9` resolves test packages from `tests/creationtests/`. See `spec/01-app/13-testing-patterns.md` §6.1 for the `enum-v9` layout.

## Purpose

This spec documents how to build a PowerShell-based Go test coverage runner
that can work with **any** Go module / repository. An AI agent or developer
can use this spec to produce a working `run.ps1` (or equivalent) for a new
Go project without prior context.

---

## 1. High-Level Architecture

```
┌───────────────────────────────────────────────────────────────────┐
│                       run.ps1 <Command>                          │
│                                                                  │
│  Commands:                                                       │
│    T   — Run all tests (verbose)                                 │
│    TC  — Run tests with coverage (parallel default)              │
│    TCP — Run coverage for a single package                       │
│    PC  — Pre-commit compile check                                │
│    TF  — Show last failing tests log                             │
│    C   — Clean build artifacts                                   │
│    H   — Help                                                    │
│                                                                  │
│  Flags:                                                          │
│    --sync   — Force sequential mode (default: parallel)          │
│    --open   — Open HTML coverage report in browser               │
└───────────────────────────────────────────────────────────────────┘
```

## 2. Directory Layout

```
<project-root>/
├── run.ps1                          # Entry point — thin dispatcher (~167 lines)
├── scripts/
│   ├── README.md                          # Module documentation
│   ├── DashboardUI.psm1                   # ANSI dashboard rendering
│   ├── Utilities.psm1                     # Common helpers
│   ├── TestLogWriter.psm1                 # Test output → log files
│   ├── TestRunner.psm1                    # Test execution + build checks
│   ├── CoverageRunner.psm1               # TC + TCP coverage pipelines
│   ├── BuildTools.psm1                    # Build, format, vet, tidy, clean
│   ├── PreCommitCheck.psm1                # PC pre-commit validation
│   ├── GoConvey.psm1                      # GoConvey launcher
│   ├── Help.psm1                          # Help display + misc commands
│   ├── bracecheck/main.go                 # Go syntax pre-checker
│   ├── autofix/main.go                    # Auto-fixer for common syntax issues
│   ├── check-safetest-boundaries.ps1      # SafeTest lint checker
│   ├── check-integrated-regressions.ps1   # API-drift regression scanner
│   └── coverage/
│       └── Export-UncoveredMethodsJson.ps1 # Coverage gap JSON exporter
├── data/
│   ├── coverage/
│   │   ├── partial/                 # Per-package .out profiles
│   │   ├── coverage.out             # Merged coverage profile
│   │   ├── coverage.html            # HTML report
│   │   ├── coverage-summary.txt     # go tool cover -func output
│   │   ├── coverage-summary.json    # Machine-readable summary
│   │   ├── per-package-coverage.txt # Per-package breakdown
│   │   ├── per-package-coverage.json
│   │   ├── blocked-packages.txt     # Compile failures (if any)
│   │   ├── blocked-packages.json
│   │   ├── coverage-previous.json   # Snapshot for regression diff
│   │   └── uncovered-method-lines.json
│   ├── test-logs/
│   │   ├── raw-output.txt           # Full go test stdout/stderr
│   │   ├── passing-tests.txt
│   │   └── failing-tests.txt
│   └── prompts/                     # AI coverage-prompt files
│       ├── coverage-prompt-1.txt
│       └── prompts-summary.json
└── tests/
    └── integratedtests/             # Test packages (one dir per source pkg)
        ├── foopkgtests/
        └── barpkgtests/
```

## 3. TC Command — Core Flow

### Phase 1: Discovery

```powershell
# List ALL packages in the module
$allPkgs = go list ./... 2>&1 | ForEach-Object { $_.ToString() }

# Source packages = everything EXCLUDING tests/
$srcPkgs = $allPkgs | Where-Object { $_ -notmatch '/tests/' }
$covPkgList = $srcPkgs -join ","

# Test packages = only integrated test dirs, sorted for determinism
$testPkgs = go list ./tests/integratedtests/... 2>&1 |
    ForEach-Object { $_.ToString() } |
    Where-Object { $_ -and $_ -notmatch '^warning:' } |
    Sort-Object
```

### Phase 2: Pre-Coverage Compile Check

Compile each test package **before** running tests. This prevents
misleading 0% coverage from broken packages.

```powershell
# Key flags:
#   -gcflags=all=-e   — report ALL errors, not just the first 10
#   -o <tempfile>     — write binary to temp dir (cleaned up later)
#   -coverpkg=$list   — instrument source packages

foreach ($pkg in $testPkgs) {
    $out = & go test -c -gcflags=all=-e -o $tempFile "-coverpkg=$covPkgList" "$pkg" 2>&1
    $exitCode = $LASTEXITCODE

    if ($exitCode -ne 0) {
        # FALLBACK: run a second compile pass to capture more errors
        $diagOut = & go test -count=1 -run '^$' -gcflags=all=-e "$pkg" 2>&1
        $combined = Merge-UniqueOutputLines $out $diagOut

        # Filter noise from combined output:
        #   - Remove "warning: no packages being tested..." lines
        #   - Remove bare package headers without file:line references
        #   - Keep only lines matching .go:\d+ (actual compile errors)
        $blocked[$shortName] = Filter-BlockedCompileLines $combined
    }
}
```

**Why two compile passes?** `go test -c` may stop after a threshold of
errors per file. The fallback `go test -count=1 -run '^$'` pass often
surfaces additional errors. Merging both (deduplicated) gives the
fullest picture in one shot.

### Phase 3: Coverage Run

Run each compilable test package individually with `-coverprofile`:

```powershell
foreach ($pkg in $compilablePkgs) {
    $profile = Join-Path $partialDir "cover-$safeName.out"
    & go test -count=1 "-coverprofile=$profile" "-coverpkg=$covPkgList" "$pkg" 2>&1
}
```

### Phase 4: Merge Profiles

Multiple partial profiles must be merged using **MAX count** per
unique coverage line (not last-write-wins):

```powershell
# Coverage line format:
#   pkg/file.go:startLine.startCol,endLine.endCol numStatements count
#
# For each unique key (everything before the count), keep the HIGHEST count.
# This prevents a later 0-count from overwriting a covered line.

foreach ($line in $allPartialLines) {
    if ($line -match '^(\S+\.go:\d+\.\d+,\d+\.\d+\s+\d+)\s+(\d+)$') {
        $key = $Matches[1]
        $count = [int]$Matches[2]
        if (-not $map.ContainsKey($key) -or $count -gt $map[$key]) {
            $map[$key] = $count
        }
    }
}
```

### Phase 5: Reports

Generate all reports from the merged profile:

```powershell
# Text summary
go tool cover "-func=$mergedProfile" > coverage-summary.txt

# HTML report
go tool cover "-html=$mergedProfile" -o coverage.html

# Per-package JSON (parse func output, group by package)
# Blocked packages JSON (filtered errors only)
# Uncovered methods JSON (for AI prompts)
```

## 4. Console Output — 4-Section Layout

The TC command prints **exactly four boxed sections**. No per-package
progress lines, no raw test output, no debug noise.

```
Section 1: BLOCKED PACKAGES    (if any failed to compile)
Section 2: FAILING TESTS       (if any tests failed)
Section 3: COVERAGE SUMMARY    (per-package table, sorted by %)
Section 4: WRITTEN FILES       (list of generated report paths)
```

Each section uses box-drawing characters:

```powershell
Write-Host "  ┌─────────────────────────────────────────────────" -ForegroundColor Red
Write-Host "  │ BLOCKED PACKAGES ($count failed to compile)"     -ForegroundColor Red
Write-Host "  │   ✗ $packageName"                                -ForegroundColor Red
Write-Host "  └─────────────────────────────────────────────────" -ForegroundColor Red
```

Color coding for coverage: ≥100% green, ≥80% yellow, <80% red.

## 5. Blocked Packages — Error Filtering Rules

Raw `go test -c` output contains noise. Apply these filters:

| Pattern | Action |
|---------|--------|
| `warning: no packages being tested depend on matches for pattern ...` | **Remove** |
| `# github.com/org/repo [...]` (no `.go:\d+`) | **Remove** (bare header) |
| `github.com/...` or `gitlab.com/...` line without `.go:\d+` | **Remove** (package path) |
| `FAIL\t...` without file reference | **Keep** (build-failed summary) |
| `*.go:123:45: error message` | **Keep** (actionable error) |

```powershell
function Filter-BlockedCompileLines([string[]]$lines) {
    $filtered = [System.Collections.Generic.List[string]]::new()
    foreach ($raw in $lines) {
        $trimmed = $raw.ToString().Trim()
        if (-not $trimmed) { continue }
        if ($trimmed -match '^\s*warning:\s*no packages being tested') { continue }
        if ($trimmed -match '^#\s+\S+' -and $trimmed -notmatch '\.go:\d+') { continue }
        if ($trimmed -match '^(github\.com|gitlab\.com)/\S+' -and $trimmed -notmatch '\.go:\d+') { continue }
        $filtered.Add($raw)
    }
    return $filtered.ToArray()
}
```

## 6. Parallel Execution

Both compile-check and coverage-run support parallel mode using
`ForEach-Object -Parallel -ThrottleLimit $N`:

```powershell
$throttle = [Math]::Min($pkgCount, [Environment]::ProcessorCount * 2)

$results = $packages | ForEach-Object -ThrottleLimit $throttle -Parallel {
    $pkg = $_
    $covPkgs = $using:covPkgList
    $tempDir = $using:compileTemp
    # ... run go test ...
    [pscustomobject]@{ Pkg = $pkg; ExitCode = $ec; Output = $out }
}

# CRITICAL: Sort results by Pkg for deterministic output
foreach ($r in ($results | Sort-Object Pkg)) { ... }
```

**Rules:**
- Use package-derived sanitized filenames: `$pkg -replace '[^a-zA-Z0-9\.-]', '_'`
- Always sort results by package name before processing
- `--sync` flag falls back to sequential `foreach`

## 7. Test Log Separation

All raw `go test` output goes to `data/test-logs/raw-output.txt`.
Parsed results split into:

- **passing-tests.txt** — list of `--- PASS:` test names
- **failing-tests.txt** — structured report with:
  - Summary section (list of failed test names)
  - Details section (per-test diagnostic blocks)
  - Compilation error fallback (if no `=== RUN` lines found)

## 8. Pre-Commit Check (PC Command)

Compile-only check for test files matching `Coverage*.go`:

```powershell
# 1. Find all integratedtest dirs containing Coverage*.go files
# 2. go test -c -gcflags=all=-e each one
# 3. Parse errors with ParseCompileErrors (categorizes: arg-count,
#    undefined, type-mismatch, missing-member, field-vs-method)
# 4. Print boxed summary + per-package error details
# 5. Write api-check.json report
```

## 9. Regression Guard Script

`scripts/check-integrated-regressions.ps1` scans `Coverage*.go` files
for known API-drift patterns **before** compilation:

| Rule | Detects |
|------|---------|
| corejson-result-err | `.Err` instead of `.Error` on corejson.Result |
| legacy-casev1-field | Old field names (Name/Input/Expected/Actual) |
| hashmap-invalid-add | `.Add()` instead of `.AddOrUpdate()` |
| simpleslice-renamed | Deprecated method names (SortedAsc→Sort, etc.) |
| simpleslice-strings-variadic | Strings() used with variadic args instead of []string |

## 10. Generic Sample Script

Below is a **minimal, self-contained** PowerShell script that works
with any Go module. Copy it, adjust `$TestPkgPattern`, and run.

```powershell
#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Generic Go test coverage runner.
.DESCRIPTION
    Usage: ./run-coverage.ps1 [-Sync] [-Open]
    Runs all test packages, generates merged coverage profile + reports.
#>
param(
    [switch]$Sync,
    [switch]$Open
)

$ErrorActionPreference = "Stop"

# ── CONFIGURATION (adjust for your project) ──────────────────────
$TestPkgPattern  = "./tests/..."          # go list pattern for test packages
$ExcludePattern  = '/vendor/'             # regex to exclude from source pkgs
$OutputDir       = Join-Path $PSScriptRoot "data" "coverage"
$TestLogDir      = Join-Path $PSScriptRoot "data" "test-logs"
# ─────────────────────────────────────────────────────────────────

New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null
New-Item -ItemType Directory -Path $TestLogDir -Force | Out-Null
$partialDir = Join-Path $OutputDir "partial"
if (Test-Path $partialDir) { Remove-Item -Recurse -Force $partialDir }
New-Item -ItemType Directory -Path $partialDir -Force | Out-Null

# ── Phase 1: Discovery ───────────────────────────────────────────
Write-Host "`n=== Discovering packages ===" -ForegroundColor Cyan

$allPkgs  = go list ./... 2>&1 | ForEach-Object { $_.ToString() }
$srcPkgs  = $allPkgs | Where-Object { $_ -notmatch '/tests/' -and $_ -notmatch $ExcludePattern }
$covPkgs  = $srcPkgs -join ","
$testPkgs = go list $TestPkgPattern 2>&1 |
    ForEach-Object { $_.ToString() } |
    Where-Object { $_ -and $_ -notmatch '^warning:' } |
    Sort-Object

Write-Host "  Source packages : $($srcPkgs.Count)"
Write-Host "  Test packages   : $($testPkgs.Count)"

# ── Phase 2: Compile Check ───────────────────────────────────────
Write-Host "`n=== Pre-coverage compile check ===" -ForegroundColor Cyan

$blocked = [ordered]@{}
$compilable = [System.Collections.Generic.List[string]]::new()
$compileTemp = Join-Path $OutputDir "compile-check"
if (Test-Path $compileTemp) { Remove-Item -Recurse -Force $compileTemp }
New-Item -ItemType Directory -Path $compileTemp -Force | Out-Null

function Get-SafeName([string]$pkg) {
    return $pkg -replace '[^a-zA-Z0-9\.-]', '_'
}

function Filter-Noise([string[]]$lines) {
    return @($lines | Where-Object {
        $t = $_.Trim()
        $t -and
        $t -notmatch '^\s*warning:\s*no packages being tested' -and
        -not ($t -match '^#\s+\S+' -and $t -notmatch '\.go:\d+')
    })
}

foreach ($pkg in $testPkgs) {
    $safe = Get-SafeName $pkg
    $outFile = Join-Path $compileTemp "compile-$safe.test"

    $ErrorActionPreference = "Continue"
    $out = & go test -c -gcflags=all=-e -o $outFile "-coverpkg=$covPkgs" "$pkg" 2>&1 |
        ForEach-Object { $_.ToString() }
    $ec = $LASTEXITCODE
    $ErrorActionPreference = "Stop"

    if ($ec -eq 0) {
        $compilable.Add($pkg)
    } else {
        # Second pass for more errors
        $ErrorActionPreference = "Continue"
        $diag = & go test -count=1 -run '^$' -gcflags=all=-e "$pkg" 2>&1 |
            ForEach-Object { $_.ToString() }
        $ErrorActionPreference = "Stop"

        # Merge unique lines
        $seen = [System.Collections.Generic.HashSet[string]]::new()
        $merged = @(($out + $diag) | Where-Object { $_ -and $seen.Add($_.Trim()) })

        $short = $pkg -replace '.*tests/?', ''
        $blocked[$short] = Filter-Noise $merged
    }
}

Remove-Item -Recurse -Force $compileTemp -ErrorAction SilentlyContinue

# Print blocked summary
if ($blocked.Count -gt 0) {
    Write-Host "`n  ┌─────────────────────────────────────────────────" -ForegroundColor Red
    Write-Host "  │ BLOCKED PACKAGES ($($blocked.Count) failed to compile)" -ForegroundColor Red
    foreach ($name in $blocked.Keys) {
        Write-Host "  │   ✗ $name" -ForegroundColor Red
    }
    Write-Host "  │ Fix build errors to include them." -ForegroundColor Yellow
    Write-Host "  └─────────────────────────────────────────────────" -ForegroundColor Red

    # Write blocked-packages.json
    $items = @($blocked.GetEnumerator() | ForEach-Object {
        @{ package = $_.Key; errorCount = $_.Value.Count; errors = $_.Value }
    })
    @{
        timestamp     = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")
        blockedCount  = $blocked.Count
        compiledCount = $compilable.Count
        totalCount    = $testPkgs.Count
        blockedPackages = $items
    } | ConvertTo-Json -Depth 4 |
        Set-Content (Join-Path $OutputDir "blocked-packages.json") -Encoding UTF8
} else {
    Write-Host "  ✓ All $($compilable.Count) packages compiled" -ForegroundColor Green
}

if ($compilable.Count -eq 0) {
    Write-Host "  ✗ No packages compiled — aborting" -ForegroundColor Red
    exit 1
}

# ── Phase 3: Coverage Run ────────────────────────────────────────
Write-Host "`n=== Running $($compilable.Count) test packages ===" -ForegroundColor Yellow

$allOutput = [System.Collections.Generic.List[string]]::new()

if ($Sync) {
    foreach ($pkg in $compilable) {
        $safe = Get-SafeName $pkg
        $profile = Join-Path $partialDir "cover-$safe.out"
        $ErrorActionPreference = "Continue"
        $out = & go test -count=1 "-coverprofile=$profile" "-coverpkg=$covPkgs" "$pkg" 2>&1 |
            ForEach-Object { $_.ToString() }
        $ErrorActionPreference = "Stop"
        foreach ($line in $out) { $allOutput.Add($line) }
    }
} else {
    $throttle = [Math]::Min($compilable.Count, [Environment]::ProcessorCount * 2)
    $results = $compilable | ForEach-Object -ThrottleLimit $throttle -Parallel {
        $pkg = $_
        $covList = $using:covPkgs
        $pDir = $using:partialDir
        $safe = $pkg -replace '[^a-zA-Z0-9\.-]', '_'
        $profile = Join-Path $pDir "cover-$safe.out"
        $ErrorActionPreference = "Continue"
        $out = & go test -count=1 "-coverprofile=$profile" "-coverpkg=$covList" "$pkg" 2>&1 |
            ForEach-Object { $_.ToString() }
        [pscustomobject]@{ Pkg = $pkg; Output = $out; ExitCode = $LASTEXITCODE }
    }
    foreach ($r in ($results | Sort-Object Pkg)) {
        foreach ($line in $r.Output) { $allOutput.Add($line) }
    }
}

# Write raw test output
Set-Content (Join-Path $TestLogDir "raw-output.txt") ($allOutput -join "`n") -Encoding UTF8

# ── Phase 4: Merge Profiles (MAX count) ──────────────────────────
$coverMap = [System.Collections.Generic.Dictionary[string, int]]::new()
foreach ($pf in (Get-ChildItem $partialDir -Filter "cover-*.out")) {
    foreach ($line in (Get-Content $pf.FullName)) {
        if (-not $line -or $line -match '^mode:') { continue }
        if ($line -match '^(\S+\.go:\d+\.\d+,\d+\.\d+\s+\d+)\s+(\d+)\s*$') {
            $key = $Matches[1]; $count = [int]$Matches[2]
            if (-not $coverMap.ContainsKey($key) -or $count -gt $coverMap[$key]) {
                $coverMap[$key] = $count
            }
        }
    }
}

$mergedProfile = Join-Path $OutputDir "coverage.out"
$mergedLines = @("mode: set") + @($coverMap.GetEnumerator() |
    ForEach-Object { "$($_.Key) $($_.Value)" })
Set-Content $mergedProfile ($mergedLines -join "`n") -Encoding UTF8

# ── Phase 5: Generate Reports ────────────────────────────────────
$htmlFile = Join-Path $OutputDir "coverage.html"
$summaryFile = Join-Path $OutputDir "coverage-summary.txt"

$funcOutput = & go tool cover "-func=$mergedProfile" 2>&1 |
    ForEach-Object { $_.ToString() }
Set-Content $summaryFile ($funcOutput -join "`n") -Encoding UTF8

& go tool cover "-html=$mergedProfile" -o $htmlFile 2>$null

# Extract total coverage
$totalLine = $funcOutput | Where-Object { $_ -match 'total:' } | Select-Object -Last 1
$totalPct = "0.0%"
if ($totalLine -match '(\d+\.\d+)%') { $totalPct = "$($Matches[1])%" }

# ── Phase 6: Console Summary ─────────────────────────────────────
Write-Host "`n  ┌─────────────────────────────────────────────────" -ForegroundColor Cyan
Write-Host "  │ COVERAGE SUMMARY" -ForegroundColor Cyan
Write-Host "  │  total: $totalPct" -ForegroundColor Cyan
Write-Host "  └─────────────────────────────────────────────────" -ForegroundColor Cyan

Write-Host "`n  ┌─────────────────────────────────────────────────" -ForegroundColor Gray
Write-Host "  │ WRITTEN FILES" -ForegroundColor Gray
Write-Host "  │  $mergedProfile" -ForegroundColor Gray
Write-Host "  │  $htmlFile" -ForegroundColor Gray
Write-Host "  │  $summaryFile" -ForegroundColor Gray
Write-Host "  └─────────────────────────────────────────────────" -ForegroundColor Gray

if ($Open) { Start-Process $htmlFile }

Write-Host "`nDone." -ForegroundColor Green
```

### Adapting for Your Project

| Setting | What to change |
|---------|---------------|
| `$TestPkgPattern` | Your test package go-list pattern (e.g., `./test/...`, `./...`) |
| `$ExcludePattern` | Regex for packages to exclude from `-coverpkg` |
| `$OutputDir` | Where reports go |
| Parallel throttle | Adjust multiplier for your CI runner |

## 11. Key Lessons and Anti-Patterns

| Do | Don't |
|----|-------|
| Use `-gcflags=all=-e` to get ALL errors | Stop at first 10 errors (Go default) |
| Run two compile passes and merge | Assume single pass catches everything |
| Filter warning/header noise from blocked output | Dump raw `go test -c` output to reports |
| Merge profiles with MAX count per line | Use last-write-wins (overwrites covered lines) |
| Sort parallel results by package name | Stream results as they arrive (non-deterministic) |
| Use package-derived sanitized filenames | Use sequential counters (`cover-1.out`) for parallel |
| Print 4-section boxed summary only | Print per-package progress lines |
| Separate test logs (pass/fail/raw) | Mix everything into one file |
