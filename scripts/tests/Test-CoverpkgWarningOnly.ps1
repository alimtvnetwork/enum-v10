# AZ smoke test (RCA Pattern P8): Test-IsCoverpkgWarningOnlyOutput
#
# Locks the contract that powers the four call sites that previously
# misclassified `-coverpkg=` warnings as build/runtime failures
# (CoverageCompileCheck.psm1:200; CoverageRunner.psm1:217,236).
#
# Runs under pwsh on Linux/Windows. No Pester required — exits non-zero
# on the first failed assertion so the CI step turns red immediately.
# See .lovable/memory/07-test-failure-rca-patterns.md (Pattern 8).

[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSCommandPath)
Import-Module (Join-Path $repoRoot 'Utilities.psm1') -Force -DisableNameChecking

$failures = New-Object System.Collections.Generic.List[string]

function Assert-Eq {
    param([string]$name, [bool]$expected, [object]$actualLines)
    $actual = Test-IsCoverpkgWarningOnlyOutput $actualLines
    if ($actual -ne $expected) {
        $failures.Add(("[FAIL] {0}: expected={1} actual={2}" -f $name, $expected, $actual))
    } else {
        Write-Host ("[PASS] {0}" -f $name)
    }
}

# 1. Pure warnings stream — the canonical false-positive that bombed
#    licensetype / onofftype / rootcmdnames in Cycle 51.
$warnOnly = @(
    'warning: no packages being tested depend on matches for pattern github.com/alimtvnetwork/enum-v10/...',
    'warning: no packages being tested depend on matches for pattern github.com/alimtvnetwork/enum-v10/...'
)
Assert-Eq 'pure-warnings' $true $warnOnly

# 2. Warnings interleaved with PASS / "ok " / blank lines — still warnings-only.
$warnWithNoise = @(
    '',
    'warning: no packages being tested depend on matches for pattern github.com/alimtvnetwork/enum-v10/...',
    'PASS',
    'ok  	github.com/alimtvnetwork/enum-v10/foo	0.123s',
    ''
)
Assert-Eq 'warnings+pass+ok+blank' $true $warnWithNoise

# 3. Warnings + a real "FAIL pkg [build failed]" marker — tolerated by
#    the helper as expected (the marker is metadata, not a diagnostic).
$warnWithBuildFailedMarker = @(
    'warning: no packages being tested depend on matches for pattern github.com/alimtvnetwork/enum-v10/...',
    'FAIL	github.com/alimtvnetwork/enum-v10/brackets	[build failed]'
)
Assert-Eq 'warnings+build-failed-marker' $true $warnWithBuildFailedMarker

# 4. Empty / null inputs — must NOT be classified as warnings-only.
Assert-Eq 'empty-array' $false @()
Assert-Eq 'null-input'  $false $null

# 5. Real compile error mixed with warnings — MUST return false so the
#    error is surfaced.
$realError = @(
    'warning: no packages being tested depend on matches for pattern github.com/alimtvnetwork/enum-v10/...',
    './foo.go:12:3: undefined: bar'
)
Assert-Eq 'real-compile-error' $false $realError

# 6. Real test failure — MUST return false (this is the regression we're
#    locking: the helper must never swallow a genuine FAIL line).
$realFail = @(
    'warning: no packages being tested depend on matches for pattern github.com/alimtvnetwork/enum-v10/...',
    '--- FAIL: TestFoo (0.00s)',
    '    foo_test.go:42: expected 1 got 2',
    'FAIL'
)
Assert-Eq 'real-test-failure' $false $realFail

# 7. Only PASS/ok markers, no warnings — must be false (no warning seen,
#    so the helper has nothing to suppress).
$pkgPassOnly = @('PASS', 'ok  	github.com/alimtvnetwork/enum-v10/foo	0.123s')
Assert-Eq 'pass-without-warnings' $false $pkgPassOnly

# 8. CRLF input — `\r` on line endings must not throw the matcher off.
$crlf = @(
    "warning: no packages being tested depend on matches for pattern github.com/alimtvnetwork/enum-v10/...`r",
    "PASS`r"
)
Assert-Eq 'crlf-line-endings' $true $crlf

if ($failures.Count -gt 0) {
    Write-Host ''
    Write-Host ('=== {0} FAILURE(S) ===' -f $failures.Count) -ForegroundColor Red
    foreach ($f in $failures) { Write-Host $f -ForegroundColor Red }
    exit 1
}

Write-Host ''
Write-Host '=== ALL CHECKS PASSED ===' -ForegroundColor Green
exit 0
