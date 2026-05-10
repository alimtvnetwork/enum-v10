# ─────────────────────────────────────────────────────────────────────────────
# TestRunner.psm1 — High-level test commands (TA, TP)
#
# Dependencies: TestRunnerCore.psm1, Utilities.psm1, TestLogWriter.psm1
# ─────────────────────────────────────────────────────────────────────────────

function Invoke-AllTests {
    <# .SYNOPSIS Run all Go tests with verbose output, build-check first. #>
    [CmdletBinding()]
    param()
    Write-Header "Running all tests"
    Invoke-FetchLatest
    Push-Location tests
    try {
        if (-not (Invoke-BuildCheck "./...")) { return }
        $prevPref = $ErrorActionPreference; $ErrorActionPreference = "Continue"
        $output = & go test -v -count=1 ./... 2>&1 | ForEach-Object { $_.ToString() }
        $exitCode = $LASTEXITCODE; $ErrorActionPreference = $prevPref
        Filter-TestWarnings $output | ForEach-Object { Write-Host $_ }
        Write-TestLogs $output
        if ($exitCode -eq 0) { Write-Success "All tests passed" }
        else { $s = Get-CallerSource; Write-Fail "Some tests failed (exit code: $exitCode) (source: $s)" }
    }
    finally { Pop-Location }
    Open-FailingTestsIfAny
}

function Invoke-PackageTests {
    <# .SYNOPSIS Run Go tests for a single package under the active test suite root
        (creationtests preferred, integratedtests legacy fallback). #>
    [CmdletBinding()]
    param([string]$pkg)
    if (-not $pkg) {
        $s = Get-CallerSource; Write-Fail "Package name required. Usage: ./run.ps1 TP <package> (source: $s)"
        # BB (Cycle 53): list packages from whichever suite root actually exists.
        $listRoot = Resolve-TestSuiteRoot
        $listDir  = Join-Path 'tests' $listRoot
        Write-Host "  Available packages (tests/$listRoot):" -ForegroundColor Yellow
        Get-ChildItem -Path $listDir -Directory -ErrorAction SilentlyContinue |
            ForEach-Object { Write-Host "    - $($_.Name)" -ForegroundColor Gray }
        return
    }
    Write-Header "Running tests for package: $pkg"
    Invoke-FetchLatest
    Push-Location tests
    try {
        # BB: per-package resolution — prefer the root that actually contains $pkg.
        $suiteRoot = Resolve-TestSuiteRoot -Package $pkg
        if (-not (Invoke-BuildCheck "./$suiteRoot/$pkg/...")) { return }
        $prevPref = $ErrorActionPreference; $ErrorActionPreference = "Continue"
        $output = & go test -v -count=1 "./$suiteRoot/$pkg/..." 2>&1 | ForEach-Object { $_.ToString() }
        $exitCode = $LASTEXITCODE; $ErrorActionPreference = $prevPref
        Filter-TestWarnings $output | ForEach-Object { Write-Host $_ }
        Write-TestLogs $output
        if ($exitCode -eq 0) { Write-Success "Package tests passed" }
        else { $s = Get-CallerSource; Write-Fail "Package tests failed (exit code: $exitCode) (source: $s)" }
    }
    finally { Pop-Location }
    Open-FailingTestsIfAny
}

Export-ModuleMember -Function @('Invoke-AllTests', 'Invoke-PackageTests')
