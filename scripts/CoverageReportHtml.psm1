# ─────────────────────────────────────────────────────────────────────────────
# CoverageReportHtml.psm1 — HTML report + AI button + console summary
#
# Dependencies: CoverageReportTxt.psm1 (Get-LowCoverageFunctions),
#               DashboardCoverage.psm1 (Write-CoverageComparison, Load/Save-CoverageSnapshot)
# ─────────────────────────────────────────────────────────────────────────────

function Write-CoverageHtmlWithAiButton {
    <# .SYNOPSIS Inject a "Copy for AI" button into the Go HTML coverage report. #>
    [CmdletBinding()]
    param([string]$CoverHtml, [string]$CoverProfile, [string[]]$FuncOutput, [hashtable]$SrcPkgStmts)

    $htmlArgs = @("-html=$CoverProfile", "-o=$CoverHtml")
    $htmlErr = & go tool cover $htmlArgs 2>&1
    $htmlExitCode = $LASTEXITCODE

    if ($htmlExitCode -ne 0 -or -not (Test-Path $CoverHtml)) {
        $s = Get-CallerSource
        Write-Host "  ⚠ Failed to generate HTML report via 'go tool cover -html' (exit: $htmlExitCode) (source: $s)" -ForegroundColor Red
        if ($htmlErr) { Write-Host "  Error: $htmlErr" -ForegroundColor Red }
        $fallbackHtml = '<!DOCTYPE html><html><head><meta charset="utf-8"><title>Coverage Report</title>' +
            '<style>body{font-family:monospace;padding:20px;background:#1e1e2e;color:#cdd6f4}' +
            'pre{white-space:pre-wrap}</style></head><body>' +
            "<h1>Coverage Report</h1><pre>$($FuncOutput -join "`n")</pre></body></html>"
        Set-Content -Path $CoverHtml -Value $fallbackHtml -Encoding UTF8
        Write-Host "  Generated fallback HTML report" -ForegroundColor Yellow
    }

    # Build AI-friendly text
    $aiTextLines = [System.Collections.Generic.List[string]]::new()
    $aiTextLines.Add("# Coverage Report — $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"); $aiTextLines.Add("")
    $aiTextLines.Add("## Goal: Improve test coverage for the packages listed below.")
    $aiTextLines.Add("Please write tests for uncovered functions, following the project's AAA pattern."); $aiTextLines.Add("")

    $totalLine = $FuncOutput | Where-Object { $_ -match "^total:" } | Select-Object -Last 1
    if ($totalLine) { $aiTextLines.Add("## Total Coverage"); $aiTextLines.Add($totalLine); $aiTextLines.Add("") }

    if ($SrcPkgStmts.Count -gt 0) {
        $aiTextLines.Add("## Per-Source-Package Coverage")
        $computed = $SrcPkgStmts.GetEnumerator() | ForEach-Object {
            $pct = if ($_.Value.Stmts -gt 0) { [math]::Round(($_.Value.Covered / $_.Value.Stmts) * 100, 1) } else { 0 }
            [pscustomobject]@{ Name = $_.Key; Pct = $pct; Stmts = $_.Value.Stmts; Covered = $_.Value.Covered }
        } | Sort-Object Pct
        foreach ($e in $computed) { $aiTextLines.Add("  $($e.Pct)%  $($e.Name)  ($($e.Covered)/$($e.Stmts) stmts)") }
        $aiTextLines.Add("")
    }

    $lowCovFuncs = Get-LowCoverageFunctions $FuncOutput
    if ($lowCovFuncs.Count -gt 0) {
        $aiTextLines.Add("## Uncovered/Low-Coverage Functions (< 50%)")
        $aiTextLines.Add("Count: $($lowCovFuncs.Count)"); $aiTextLines.Add("")
        foreach ($f in $lowCovFuncs) { $aiTextLines.Add($f.TrimStart()) }; $aiTextLines.Add("")
    }
    $aiTextLines.Add("## Instructions")
    $aiTextLines.Add("- Tests go in tests/creationtests/{pkg}tests/ (this repo) or tests/integratedtests/{pkg}tests/ (upstream core-v9 layout)")
    $aiTextLines.Add("- Use CaseV1 table-driven pattern with AAA comments")
    $aiTextLines.Add("- Focus on the lowest coverage packages first")

    $aiTextEscaped = ($aiTextLines -join "`n") -replace '\\', '\\\\' -replace "'", "\\\'" -replace "`n", '\n' -replace "`r", '' -replace '"', '\"'

    if (Test-Path $CoverHtml) {
        $htmlContent = Get-Content -Path $CoverHtml -Raw
        $buttonHtml = '<div id="ai-copy-panel" style="position:fixed;top:12px;right:12px;z-index:9999;font-family:system-ui,sans-serif;">' +
            '<button onclick="copyForAI()" style="' +
            'background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;border:none;' +
            'padding:10px 20px;border-radius:8px;font-size:14px;font-weight:600;' +
            'cursor:pointer;box-shadow:0 4px 12px rgba(99,102,241,0.4);' +
            'display:flex;align-items:center;gap:6px;transition:all 0.2s;' +
            '" onmouseover="this.style.transform=''scale(1.05)''" onmouseout="this.style.transform=''scale(1)''">' +
            '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>' +
            '  Copy for AI' +
            '</button>' +
            '<span id="ai-copy-status" style="display:none;color:#22c55e;font-size:13px;margin-top:4px;text-align:center;">Copied!</span>' +
            '</div>' +
            '<script>' +
            "var __aiCoverageText ='"
        $scriptEnd = "'" + ';' +
            'function copyForAI(){' +
            '  try {' +
            '    var ta = document.createElement("textarea");' +
            '    ta.value = __aiCoverageText;' +
            '    ta.style.position = "fixed";' +
            '    ta.style.left = "-9999px";' +
            '    document.body.appendChild(ta);' +
            '    ta.select();' +
            '    document.execCommand("copy");' +
            '    document.body.removeChild(ta);' +
            '    var s = document.getElementById("ai-copy-status");' +
            '    s.style.display = "block";' +
            '    setTimeout(function(){ s.style.display = "none"; }, 2000);' +
            '  } catch(e) {' +
            '    alert("Copy failed: " + e.message);' +
            '  }' +
            '}' +
            '</script>'
        $injectedHtml = $buttonHtml + $aiTextEscaped + $scriptEnd
        $htmlContent = $htmlContent -replace '</body>', ($injectedHtml + "`n</body>")
        Set-Content -Path $CoverHtml -Value $htmlContent -Encoding UTF8
        Write-Host "  ✓ Injected 'Copy for AI' button into HTML report" -ForegroundColor Green
    }
}

function Write-CoverageConsoleSummary {
    <# .SYNOPSIS Print coverage summary and comparison to console. #>
    [CmdletBinding()]
    param([hashtable]$SrcPkgStmts, [string[]]$FuncOutput)

    $totalLine = $FuncOutput | Where-Object { $_ -match "^total:" } | Select-Object -Last 1
    $lowCovFuncs = Get-LowCoverageFunctions $FuncOutput

    if ($SrcPkgStmts.Count -gt 0) {
        Write-Host ""
        Write-Host "  ┌─────────────────────────────────────────────────" -ForegroundColor Cyan
        Write-Host "  │ COVERAGE SUMMARY" -ForegroundColor Cyan
        Write-Host "  │" -ForegroundColor Cyan
        $sorted = $SrcPkgStmts.GetEnumerator() | ForEach-Object {
            $pct = if ($_.Value.Stmts -gt 0) { [math]::Round(($_.Value.Covered / $_.Value.Stmts) * 100, 1) } else { 0 }
            [pscustomobject]@{ Name = $_.Key; Pct = $pct }
        } | Sort-Object Pct -Descending
        foreach ($entry in $sorted) {
            $color = if ($entry.Pct -ge 100) { "Green" } elseif ($entry.Pct -ge 80) { "Yellow" } else { "Red" }
            Write-Host "  │  $($entry.Pct)%`t$($entry.Name)" -ForegroundColor $color
        }
        Write-Host "  │" -ForegroundColor Cyan
        if ($totalLine -and $totalLine -match "(\d+\.\d+)%") { Write-Host "  │  total: (statements)  $($Matches[1])%" -ForegroundColor Cyan }
        if ($lowCovFuncs.Count -gt 0) { Write-Host "  │  ⚠ $($lowCovFuncs.Count) function(s) below 50% coverage" -ForegroundColor Yellow }
        Write-Host "  └─────────────────────────────────────────────────" -ForegroundColor Cyan
    }

    if (Get-Command Write-CoverageComparison -ErrorAction SilentlyContinue) {
        $currentCovData = @()
        if ($SrcPkgStmts.Count -gt 0) {
            $currentCovData = @($SrcPkgStmts.GetEnumerator() | ForEach-Object {
                $pct = if ($_.Value.Stmts -gt 0) { [math]::Round(($_.Value.Covered / $_.Value.Stmts) * 100, 1) } else { 0 }
                @{ Package = $_.Key; Coverage = $pct }
            })
        }
        if ($currentCovData.Count -gt 0) {
            $previousCovData = $null
            if (Get-Command Load-CoverageSnapshot -ErrorAction SilentlyContinue) { $previousCovData = Load-CoverageSnapshot }
            Write-Host ""
            Write-CoverageComparison -Current $currentCovData -Previous $previousCovData
            if (Get-Command Save-CoverageSnapshot -ErrorAction SilentlyContinue) { Save-CoverageSnapshot -CoverageData $currentCovData }
        }
    }
}

Export-ModuleMember -Function @(
    'Write-CoverageHtmlWithAiButton', 'Write-CoverageConsoleSummary'
)
