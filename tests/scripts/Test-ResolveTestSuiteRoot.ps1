$ErrorActionPreference = 'Stop'
Import-Module ./scripts/Utilities.psm1 -Force

# Setup: temp project root with both possible layouts
$tmpRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("rtsr-" + [guid]::NewGuid().ToString('N').Substring(0,8))
New-Item -ItemType Directory -Path (Join-Path $tmpRoot 'tests/creationtests/osdetect') -Force | Out-Null

# Test 1: prefers creationtests when only it exists
$r1 = Resolve-TestSuiteRoot -ProjectRoot $tmpRoot
if ($r1 -ne 'creationtests') { throw "Test1 FAIL: expected 'creationtests', got '$r1'" }
Write-Host "Test 1 PASS — only creationtests exists → $r1"

# Test 2: with -Package 'osdetect' → still creationtests (pkg exists there)
$r2 = Resolve-TestSuiteRoot -ProjectRoot $tmpRoot -Package 'osdetect'
if ($r2 -ne 'creationtests') { throw "Test2 FAIL: expected 'creationtests', got '$r2'" }
Write-Host "Test 2 PASS — package osdetect resolves under creationtests → $r2"

# Test 3: with -Package 'doesnotexist' → falls back to default 'creationtests'
$r3 = Resolve-TestSuiteRoot -ProjectRoot $tmpRoot -Package 'doesnotexist'
if ($r3 -ne 'creationtests') { throw "Test3 FAIL: expected 'creationtests' (fallback), got '$r3'" }
Write-Host "Test 3 PASS — missing package falls back to default → $r3"

# Test 4: legacy-only layout → returns integratedtests
$tmpLegacy = Join-Path ([System.IO.Path]::GetTempPath()) ("rtsr-legacy-" + [guid]::NewGuid().ToString('N').Substring(0,8))
New-Item -ItemType Directory -Path (Join-Path $tmpLegacy 'tests/integratedtests/foo') -Force | Out-Null
$r4 = Resolve-TestSuiteRoot -ProjectRoot $tmpLegacy
if ($r4 -ne 'integratedtests') { throw "Test4 FAIL: expected 'integratedtests', got '$r4'" }
Write-Host "Test 4 PASS — legacy-only layout → $r4"

# Test 5: real enum-v9 repo → creationtests
$r5 = Resolve-TestSuiteRoot -ProjectRoot (Get-Location).Path
if ($r5 -ne 'creationtests') { throw "Test5 FAIL: expected 'creationtests' for live repo, got '$r5'" }
Write-Host "Test 5 PASS — live repo resolves → $r5"

Remove-Item -Recurse -Force $tmpRoot, $tmpLegacy -ErrorAction SilentlyContinue
Write-Host "All 5 tests passed."
