# ─────────────────────────────────────────────────────────────────────────────
# PreCommitCheck.psm1 — Pre-commit API mismatch checker for Coverage* files
#
# Now reuses Invoke-CoveragePreChecks for the autofix/bracecheck pipeline,
# keeping only the regression guard and compile-check logic.
#
# Usage:
#   Import-Module ./scripts/PreCommitCheck.psm1 -Force
#
# Dependencies:
#   - Utilities.psm1 (Write-Header, Write-Success, Write-Fail, Merge-UniqueOutputLines, ParseCompileErrors)
#   - CoveragePreChecks.psm1 (Invoke-CoveragePreChecks)
#   - DashboardUI/DashboardPhases (Register-Phase, Reset-Phases, Write-PhaseSummaryBox)
# ─────────────────────────────────────────────────────────────────────────────

function Invoke-PreCommitCheck {
    <#
    .SYNOPSIS
        Run the full pre-commit validation pipeline.
    .DESCRIPTION
        Phases: regression guard → safetest/autofix/bracecheck (via CoveragePreChecks)
        → compile-check Coverage* packages → JSON report → dashboard summary.
    .PARAMETER singlePkg
        Optional package name to check only one package.
    #>
    [CmdletBinding()]
    param([string]$singlePkg)

    if (Get-Command Reset-Phases -ErrorAction SilentlyContinue) { Reset-Phases }
    Write-Header "Pre-commit API mismatch checker"

    $isSyncMode = $ExtraArgs -and ($ExtraArgs -contains "--sync")

    # ── Regression guard ──
    $regressionScript = Join-Path $global:ProjectRoot "scripts" "check-integrated-regressions.ps1"
    if (-not (Test-Path $regressionScript)) {
        $s = Get-CallerSource; Write-Fail "Regression guard script not found: $regressionScript (source: $s)"
        exit 1
    }
    Write-Host "  Running regression guard scan..." -ForegroundColor Yellow
    if ($singlePkg) { & $regressionScript -SinglePackage $singlePkg } else { & $regressionScript }
    if ($LASTEXITCODE -ne 0) {
        if (Get-Command Register-Phase -ErrorAction SilentlyContinue) { Register-Phase "Regression Guard" "fail" "regressions detected" }
        $s = Get-CallerSource; Write-Fail "Regression guard failed. Fix reported issues before PC. (source: $s)"
        exit 1
    }
    if (Get-Command Register-Phase -ErrorAction SilentlyContinue) { Register-Phase "Regression Guard" "pass" "no regressions" }

    # ── Pre-checks (safetest + autofix + bracecheck) via shared module ──
    $coverDir = Join-Path $global:DataDir "coverage"
    $preCheckOk = Invoke-CoveragePreChecks -ScriptRoot $global:ProjectRoot -ExtraArgs $ExtraArgs -CoverDir $coverDir
    if (-not $preCheckOk) { exit 1 }

    # ── Discover Coverage* packages ──
    # BB (Cycle 53): Resolve-TestSuiteRoot honours the Core-memory rule
    # that tests live under tests/creationtests/, NOT tests/integratedtests/.
    # Per-package: prefer the root that actually contains $singlePkg so a
    # mixed-layout repo still works.
    $suiteRoot   = if ($singlePkg) { Resolve-TestSuiteRoot -Package $singlePkg } else { Resolve-TestSuiteRoot }
    $testBaseDir = Join-Path (Join-Path $global:ProjectRoot 'tests') $suiteRoot
    if ($singlePkg) {
        $targetDirs = @(Join-Path $testBaseDir $singlePkg)
        if (-not (Test-Path $targetDirs[0])) { $s = Get-CallerSource; Write-Fail "Package not found: $singlePkg under tests/$suiteRoot/ (source: $s)"; return }
    } else {
        $targetDirs = @(Get-ChildItem -Path $testBaseDir -Directory -ErrorAction SilentlyContinue | ForEach-Object { $_.FullName })
    }

    $pkgsWithCoverage = [System.Collections.Generic.List[string]]::new()
    foreach ($dir in $targetDirs) {
        $coverFiles = Get-ChildItem -Path $dir -Filter "Coverage*" -File -ErrorAction SilentlyContinue
        if ($coverFiles -and $coverFiles.Count -gt 0) { $pkgsWithCoverage.Add($dir) }
    }

    if ($pkgsWithCoverage.Count -eq 0) { Write-Success "No Coverage* files found to check"; return }

    $modeLabel = if ($isSyncMode) { "sync" } else { "parallel" }
    Write-Host "  Checking $($pkgsWithCoverage.Count) packages with Coverage* files ($modeLabel)..." -ForegroundColor Yellow
    Write-Host ""

    $goTestPkgs = [System.Collections.Generic.List[string]]::new()
    foreach ($dir in $pkgsWithCoverage) {
        $relPath = $dir -replace [regex]::Escape($global:ProjectRoot), '' -replace '^[\\/]', '' -replace '\\', '/'
        $goTestPkgs.Add("github.com/alimtvnetwork/core-v9/$relPath")
    }

    # ── Compile check ──
    $compileTemp = Join-Path $global:DataDir "precommit"
    if (Test-Path $compileTemp) { Remove-Item -Recurse -Force $compileTemp }
    New-Item -ItemType Directory -Path $compileTemp -Force | Out-Null

    $failures = [System.Collections.Generic.List[object]]::new()
    $passedCount = 0

    if ($isSyncMode) {
        foreach ($pkg in $goTestPkgs) {
            $shortName = $pkg -replace '.*integratedtests/?', ''
            $safeName = $pkg -replace '[^a-zA-Z0-9\.-]', '_'
            $outFile = Join-Path $compileTemp "check-$safeName.test"
            $prevPref = $ErrorActionPreference; $ErrorActionPreference = "Continue"
            $compOut = & go test -c -gcflags=all=-e -o $outFile "$pkg" 2>&1 | ForEach-Object { $_.ToString() }
            $ec = $LASTEXITCODE; $ErrorActionPreference = $prevPref
            if ($ec -eq 0) { $passedCount++ }
            else {
                $prevPref = $ErrorActionPreference; $ErrorActionPreference = "Continue"
                $diagOut = & go test -count=1 -run '^$' -gcflags=all=-e "$pkg" 2>&1 | ForEach-Object { $_.ToString() }
                $ErrorActionPreference = $prevPref
                $compOut = Merge-UniqueOutputLines $compOut $diagOut
                $parsedErrors = ParseCompileErrors $compOut
                $failures.Add(@{ package = $shortName; errorCount = $parsedErrors.Count; errors = $parsedErrors })
            }
        }
    } else {
        $throttle = [Math]::Min($goTestPkgs.Count, [Environment]::ProcessorCount * 2)
        $results = $goTestPkgs | ForEach-Object -ThrottleLimit $throttle -Parallel {
            $pkg = $_; $tempDir = $using:compileTemp
            $safeName = $pkg -replace '[^a-zA-Z0-9\.-]', '_'
            $outFile = Join-Path $tempDir "check-$safeName.test"
            $ErrorActionPreference = "Continue"
            $rawOut = & go test -c -gcflags=all=-e -o $outFile "$pkg" 2>&1
            $ec = $LASTEXITCODE
            $out = @($rawOut | ForEach-Object { $_.ToString() })
            if ($ec -ne 0) {
                $diagRaw = & go test -count=1 -run '^$' -gcflags=all=-e "$pkg" 2>&1
                $diagOut = @($diagRaw | ForEach-Object { $_.ToString() })
                $seen = [System.Collections.Generic.HashSet[string]]::new([System.StringComparer]::Ordinal)
                $merged = [System.Collections.Generic.List[string]]::new()
                foreach ($line in @($out + $diagOut)) {
                    if ($null -eq $line) { continue }
                    $normalized = $line.ToString().TrimEnd("`r")
                    if (-not $normalized) { continue }
                    if ($seen.Add($normalized)) { $merged.Add($normalized) | Out-Null }
                }
                $out = $merged.ToArray()
            }
            [pscustomobject]@{ Pkg = $pkg; ExitCode = $ec; Output = $out }
        }
        foreach ($r in ($results | Sort-Object Pkg)) {
            $shortName = $r.Pkg -replace '.*integratedtests/?', ''
            if ($r.ExitCode -eq 0) { $passedCount++ }
            else {
                $parsedErrors = ParseCompileErrors $r.Output
                $failures.Add(@{ package = $shortName; errorCount = $parsedErrors.Count; errors = $parsedErrors })
            }
        }
    }

    Get-ChildItem -Path $compileTemp -Filter "*.test" -ErrorAction SilentlyContinue | Remove-Item -Force

    # ── Summary ──
    $allPassed = $failures.Count -eq 0
    Write-Host ""
    if ($allPassed) {
        Write-Host "  ┌─────────────────────────────────────────────────" -ForegroundColor Green
        Write-Host "  │ ✓ ALL $passedCount PACKAGES PASSED API CHECK" -ForegroundColor Green
        Write-Host "  └─────────────────────────────────────────────────" -ForegroundColor Green
    } else {
        Write-Host "  ┌─────────────────────────────────────────────────" -ForegroundColor Red
        Write-Host "  │ ✗ $($failures.Count) PACKAGE(S) HAVE API MISMATCHES" -ForegroundColor Red
        Write-Host "  │" -ForegroundColor Red
        foreach ($f in $failures) { Write-Host "  │   ✗ $($f.package) ($($f.errorCount) error(s))" -ForegroundColor Red }
        Write-Host "  │" -ForegroundColor Red
        Write-Host "  │ Fix these before committing Coverage* files." -ForegroundColor Yellow
        Write-Host "  └─────────────────────────────────────────────────" -ForegroundColor Red
        Write-Host ""
        foreach ($f in $failures) {
            Write-Host "  ── $($f.package) ──" -ForegroundColor Red
            foreach ($e in $f.errors) { Write-Host "    $($e.file):$($e.line) [$($e.category)] $($e.message)" -ForegroundColor Yellow }
            Write-Host ""
        }
    }

    # JSON report
    $pcSource = Get-CallerSource
    $jsonReport = @{ timestamp = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"); passed = $allPassed; checkedCount = $goTestPkgs.Count; passedCount = $passedCount; failedCount = $failures.Count; source = $pcSource; failures = $failures.ToArray() }
    $jsonPath = Join-Path $compileTemp "api-check.json"
    $jsonReport | ConvertTo-Json -Depth 5 | Set-Content -Path $jsonPath -Encoding UTF8
    Write-Host "  Report → $jsonPath" -ForegroundColor Gray

    if (Get-Command Register-Phase -ErrorAction SilentlyContinue) {
        if ($allPassed) { Register-Phase "API Compile Check" "pass" "$passedCount/$($goTestPkgs.Count) passed" }
        else { Register-Phase "API Compile Check" "fail" "$($failures.Count) failed, $passedCount passed" }
    }
    if (Get-Command Write-PhaseSummaryBox -ErrorAction SilentlyContinue) { Write-Host ""; Write-PhaseSummaryBox }
    if (-not $allPassed) { exit 1 }
}

Export-ModuleMember -Function @('Invoke-PreCommitCheck')
