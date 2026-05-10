# ─────────────────────────────────────────────────────────────────────────────
# spec-api-check.psm1 — Spec-vs-upstream API fabrication lint (S-106)
#
# Scans spec/01-app/**/*.md for Go API references (e.g. `errcore.VarTwo`,
# `corevalidator.New.Line.NotEmpty()`, `coredynamic.AllFields`) and verifies
# every referenced top-level package exists in the upstream `core-v9` clone.
#
# Designed to catch the fabrication pattern that produced 50 ❌ across
# Cycles 19-25 (C-CVS-11..59). Specifically catches:
#
#   • Non-existent packages          (e.g. `corestr.*`, `coredynamic.*`)
#   • Non-existent top-level vars    (e.g. `corevalidator.New`, `errcore.InvalidInput`)
#   • Non-existent functions/types   (e.g. `errcore.VarTwo` w/ wrong arity is NOT caught;
#                                     this is a presence lint, not a signature lint)
#
# Usage:
#   Import-Module ./scripts/spec-api-check.psm1 -Force
#   Invoke-SpecApiCheck                    # default paths
#   Invoke-SpecApiCheck -SpecDir spec/01-app -UpstreamDir /tmp/core-v9-upstream
#   Invoke-SpecApiCheck -StrictExitCode    # exit 1 on any fabrication (CI mode)
#
# Dependencies: none beyond PowerShell 7+ and a populated upstream clone.
# Out-of-scope (deliberately): signature-level lint; that requires a Go AST
# pass and is tracked separately. S-106 is the cheap-and-fast first wall.
# ─────────────────────────────────────────────────────────────────────────────

$script:SpecApiCheckVersion = '1.2.0'  # S-115 — Test-UpstreamClone helper + sentinel guard

# Pinned upstream tag — bump when core-v9 releases a new audit baseline.
$script:UpstreamCloneUrl = 'https://github.com/alimtvnetwork/core-v9'
$script:UpstreamCloneTag = 'v1.5.8'
$script:UpstreamCloneDefaultPath = '/tmp/core-v9-upstream'

function Test-UpstreamClone {
    <#
    .SYNOPSIS S-115 — Verify the core-v9 upstream clone is present and well-formed.
    .DESCRIPTION
        Audit probes (e.g. AB-residual ❓→✅/❌ classification) read source files
        directly from the upstream clone with rg/grep. If the clone is missing,
        every probe returns 0 hits and produces FALSE fabrication conclusions
        (R-CVS-01, R-CVS-02, R-CVS-03 were all the same drift). This helper
        is the single source of truth: call it BEFORE any upstream probe.
    .PARAMETER Path Defaults to /tmp/core-v9-upstream.
    .PARAMETER AutoClone If set, runs the documented git clone when missing.
    .OUTPUTS [pscustomobject]@{ Ok; Path; Reason; PackageCount }
    #>
    [CmdletBinding()]
    param(
        [string]$Path = $script:UpstreamCloneDefaultPath,
        [switch]$AutoClone
    )

    if (-not (Test-Path $Path)) {
        if ($AutoClone) {
            Write-Host "  ▶ Cloning core-v9 @ $script:UpstreamCloneTag → $Path" -ForegroundColor Yellow
            & git clone --depth 1 --branch $script:UpstreamCloneTag $script:UpstreamCloneUrl $Path 2>&1 | Out-Host
            if ($LASTEXITCODE -ne 0 -or -not (Test-Path $Path)) {
                return [pscustomobject]@{ Ok = $false; Path = $Path; Reason = 'clone-failed'; PackageCount = 0 }
            }
        } else {
            return [pscustomobject]@{ Ok = $false; Path = $Path; Reason = 'missing'; PackageCount = 0 }
        }
    }

    # Sentinel: a real core-v9 checkout always has coredata/coregeneric.
    $sentinel = Join-Path $Path 'coredata/coregeneric'
    if (-not (Test-Path $sentinel)) {
        return [pscustomobject]@{ Ok = $false; Path = $Path; Reason = 'sentinel-missing (coredata/coregeneric absent — wrong branch or shallow rename?)'; PackageCount = 0 }
    }

    $pkgCount = (Get-ChildItem -Path $Path -Recurse -Directory -ErrorAction SilentlyContinue |
                 Where-Object {
                     $_.FullName -notmatch '[\\/]\.git[\\/]' -and
                     (Get-ChildItem -Path $_.FullName -Filter '*.go' -File -ErrorAction SilentlyContinue | Select-Object -First 1)
                 }).Count

    return [pscustomobject]@{ Ok = $true; Path = $Path; Reason = 'ok'; PackageCount = $pkgCount }
}

# Packages that ALWAYS resolve (Go stdlib + project-local non-upstream pkgs).
# Anything matching is allow-listed without a directory check.
$script:AllowListedPackages = @(
    # Go stdlib (subset commonly cited in spec)
    'fmt', 'strings', 'strconv', 'errors', 'context', 'time', 'os', 'io',
    'bytes', 'sort', 'sync', 'reflect', 'regexp', 'unicode', 'utf8',
    'http', 'json', 'log', 'slog', 'testing', 'math', 'unsafe', 'runtime',
    'path', 'filepath', 'url', 'net', 'rand', 'crypto', 'sha256', 'base64',
    # Common third-party / OTel convention names
    'codes', 'tracer', 'span', 'security', 'logger', 'otel',
    # Pseudo-package names used in templates / illustrative examples
    # (these are placeholder names, not real Go packages — flagging them
    # creates noise without adding signal).
    'emailvalidator', 'corev8', 'expected', 'validator', 'downstream',
    'registry',
    # Generic placeholder words used in prose / pseudo-code
    'pkg', 'parameter', 'iter', 'core'
)

# Heuristic: in prose (outside fenced code blocks) a single-segment lowercase
# token followed by `.Foo` is more likely a local variable from a removed
# fence above than a true package reference. We still flag it inside fences.
$script:ProseLooseMode = $true

# Heuristic: longer (4-char+) local-variable names commonly appear in
# illustrative snippets where the binding is OUTSIDE the visible fence
# (the fence is indented under a numbered-list item, or the binding is
# narratively elided). Flagging `safe.PasswordHash`, `cart.Total`,
# `payload.Error` is pure noise. Allow-list the most common ones.
$script:CommonLocalVarNames = @(
    'tc', 'col', 'svc', 'cart', 'safe', 'payload', 'pattern', 'result',
    'input', 'status', 'err', 'opts', 'cfg', 'req', 'resp', 'ctx', 'val',
    'item', 'items', 'row', 'rows', 'msg', 'data', 'out', 'buf'
)

function Get-UpstreamPackages {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)][string]$UpstreamDir,
        [string]$LocalDir
    )

    if (-not (Test-Path $UpstreamDir)) {
        throw "Upstream clone not found at $UpstreamDir. Re-clone with: git clone --depth 1 --branch v1.5.8 https://github.com/alimtvnetwork/core-v9 $UpstreamDir"
    }

    # S-115 sentinel: warn loudly if the directory exists but isn't a real
    # core-v9 checkout — prevents silent under-indexing that masquerades as
    # fabrications in audit probes.
    $sentinel = Join-Path $UpstreamDir 'coredata/coregeneric'
    if (-not (Test-Path $sentinel)) {
        Write-Warning "Upstream clone at $UpstreamDir is missing sentinel 'coredata/coregeneric'. Audit probes may produce false fabrications. Re-clone with: git clone --depth 1 --branch v1.5.8 https://github.com/alimtvnetwork/core-v9 $UpstreamDir"
    }

    # Collect every directory that contains at least one .go file (= a Go package).
    # Index by the BASENAME so spec references like `coregeneric.Hashmap` resolve
    # even though the real path is `coredata/coregeneric`.
    $pkgs = @{}
    $sources = @($UpstreamDir)
    if ($LocalDir -and (Test-Path $LocalDir)) { $sources += $LocalDir }
    foreach ($src in $sources) {
        Get-ChildItem -Path $src -Recurse -Filter '*.go' -File -ErrorAction SilentlyContinue |
            Where-Object {
                $_.FullName -notmatch '[\\/]\.git[\\/]' -and
                $_.FullName -notmatch '[\\/](node_modules|cross-repo|tests|scripts|spec|src|public|data|cmd|assets|configs)[\\/]'
            } |
            ForEach-Object {
                $dir = $_.Directory.FullName
                $base = Split-Path $dir -Leaf
                if (-not $pkgs.ContainsKey($base)) {
                    $pkgs[$base] = $dir
                }
            }
    }
    return $pkgs
}

function Get-UpstreamTopLevelSymbols {
    <#
    .SYNOPSIS Build {pkgName -> Set[symbolName]} index of top-level Go identifiers.
    .DESCRIPTION
        Greps `func ` / `type ` / `var ` / `const ` declarations in each package's
        .go files (excluding _test.go). Cheap regex pass — no AST.
    #>
    [CmdletBinding()]
    param([Parameter(Mandatory)][hashtable]$PackageMap)

    $symbols = @{}
    foreach ($pkg in $PackageMap.Keys) {
        $set = [System.Collections.Generic.HashSet[string]]::new()
        $dir = $PackageMap[$pkg]
        Get-ChildItem -Path $dir -Filter '*.go' -File -ErrorAction SilentlyContinue |
            Where-Object { $_.Name -notmatch '_test\.go$' } |
            ForEach-Object {
                # Top-level decls: line starts with func/type/var/const followed by an Exported identifier.
                # For func receivers, capture both the receiver-method name and bare functions.
                Select-String -Path $_.FullName -Pattern '^(func\s+(?:\([^)]+\)\s+)?|type\s+|var\s+|const\s+)([A-Z][A-Za-z0-9_]*)' |
                    ForEach-Object {
                        if ($_.Matches.Count -gt 0) {
                            [void]$set.Add($_.Matches[0].Groups[2].Value)
                        }
                    }
                # Also pick up names declared inside `var ( ... )` / `const ( ... )` / `type ( ... )` blocks.
                Select-String -Path $_.FullName -Pattern '^\s+([A-Z][A-Za-z0-9_]*)\s*(=|\s+[A-Za-z\[\*])' |
                    ForEach-Object {
                        if ($_.Matches.Count -gt 0) {
                            [void]$set.Add($_.Matches[0].Groups[1].Value)
                        }
                    }
            }
        $symbols[$pkg] = $set
    }
    return $symbols
}

function Get-SpecApiReferences {
    <#
    .SYNOPSIS Extract `package.Symbol` references from a markdown spec file.
    .DESCRIPTION
        Returns a list of [pscustomobject]@{ Package; Symbol; Line; Context }.
        Ignores: code-block fences, table-pipe lines that are pure prose,
        comment-only lines inside Go fences. Also detects local variables
        declared via `:=` inside the same fenced code block and treats
        their names as allow-listed (so `tc := …` followed by `tc.Foo()`
        does not flag `tc` as a fabricated package).
    #>
    [CmdletBinding()]
    param([Parameter(Mandatory)][string]$Path)

    $refs = New-Object System.Collections.Generic.List[object]
    $lines = Get-Content -Path $Path
    $inFence = $false
    $fenceLang = ''
    $localVars = [System.Collections.Generic.HashSet[string]]::new()

    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i]

        # Track fenced code blocks. Fences may be INDENTED (e.g. inside a
        # numbered-list item) — match leading whitespace too (S-106 v1.1).
        if ($line -match '^\s*```(\w*)') {
            if ($inFence) { $inFence = $false; $fenceLang = ''; [void]$localVars.Clear() }
            else { $inFence = $true; $fenceLang = $Matches[1].ToLower() }
            continue
        }

        # Inside a fence, track local Go variable bindings: `name :=`, `name, _ :=`, `var name `
        if ($inFence) {
            foreach ($vm in [regex]::Matches($line, '(?:^|\s|\(|,)([a-z][a-zA-Z0-9_]*)\s*(?:,\s*[a-zA-Z0-9_]+\s*)?:=')) {
                [void]$localVars.Add($vm.Groups[1].Value)
            }
            foreach ($vm in [regex]::Matches($line, '\bvar\s+([a-z][a-zA-Z0-9_]*)\s+\S')) {
                [void]$localVars.Add($vm.Groups[1].Value)
            }
            foreach ($vm in [regex]::Matches($line, '\bfor\s+([a-z][a-zA-Z0-9_]*)\s*(?:,\s*[a-zA-Z0-9_]+\s*)?:?=')) {
                [void]$localVars.Add($vm.Groups[1].Value)
            }
            # Function parameters: `func Foo(name Type, …)`
            foreach ($vm in [regex]::Matches($line, '(?:func[^(]*\(|,\s*)([a-z][a-zA-Z0-9_]*)\s+\*?[A-Za-z]')) {
                [void]$localVars.Add($vm.Groups[1].Value)
            }
            # Method receiver binding: `func (it Variant) Name() ...`
            foreach ($vm in [regex]::Matches($line, '\bfunc\s+\(\s*([a-z][a-zA-Z0-9_]*)\s+\*?[A-Z]')) {
                [void]$localVars.Add($vm.Groups[1].Value)
            }
        }

        # Match `package.Symbol` — lowercase package, Capitalised symbol.
        # Allow inline-code wrapping: `errcore.VarTwo` and bare errcore.VarTwo.
        $regex = '(?<!\w)([a-z][a-z0-9]+)\.([A-Z][A-Za-z0-9_]*)'
        $matches = [regex]::Matches($line, $regex)
        foreach ($m in $matches) {
            $pkg = $m.Groups[1].Value
            $sym = $m.Groups[2].Value
            # Skip allow-listed aliases / local variables that look package-like.
            if ($script:AllowListedPackages -contains $pkg) { continue }
            if ($localVars.Contains($pkg)) { continue }
            # Skip common local-variable names that frequently appear with an
            # elided binding (S-106 v1.1).
            if ($script:CommonLocalVarNames -contains $pkg) { continue }
            # Skip identifiers that look like a versioned local var (`v1`, `v2`, ...).
            if ($pkg -match '^v\d+$') { continue }
            # Skip references inside markdown links / URLs (file paths with .md).
            if ($line -match '\[.+\]\([^)]*' + [regex]::Escape($pkg) + '\.[A-Z]') { continue }
            # Outside fences, a token of length ≤ 3 is most likely a placeholder
            # variable from prose (e.g. `tc.Method`, `v1.Equal`). Skip in loose mode.
            if ($script:ProseLooseMode -and -not $inFence -and $pkg.Length -le 3) { continue }
            $refs.Add([pscustomobject]@{
                Package = $pkg
                Symbol  = $sym
                Line    = $i + 1
                Context = $line.Trim()
                Fenced  = $inFence
            })
        }
    }
    return $refs
}

function Invoke-SpecApiCheck {
    <#
    .SYNOPSIS S-106 fabrication lint — verifies every spec API reference resolves.
    .PARAMETER SpecDir Default: spec/01-app
    .PARAMETER UpstreamDir Default: /tmp/core-v9-upstream
    .PARAMETER StrictExitCode If set, exits 1 when any fabrication is found.
    .PARAMETER OnlyFile Limit scan to one spec file (relative to SpecDir).
    #>
    [CmdletBinding()]
    param(
        [string]$SpecDir = 'spec/01-app',
        [string]$UpstreamDir = '/tmp/core-v9-upstream',
        [string]$LocalDir = '.',
        [switch]$StrictExitCode,
        [string]$OnlyFile
    )

    Write-Host ''
    Write-Host "  ▶ S-106 spec-api-check v$script:SpecApiCheckVersion" -ForegroundColor Cyan
    Write-Host "    SpecDir:     $SpecDir"
    Write-Host "    UpstreamDir: $UpstreamDir"
    Write-Host ''

    if (-not (Test-Path $SpecDir)) { throw "SpecDir not found: $SpecDir" }

    Write-Host '  ▶ Indexing upstream packages...' -ForegroundColor Yellow
    $pkgMap = Get-UpstreamPackages -UpstreamDir $UpstreamDir -LocalDir $LocalDir
    Write-Host "    Found $($pkgMap.Count) Go packages (upstream + local enum-v10)"

    Write-Host '  ▶ Indexing upstream top-level symbols...' -ForegroundColor Yellow
    $symMap = Get-UpstreamTopLevelSymbols -PackageMap $pkgMap
    $totalSyms = ($symMap.Values | ForEach-Object { $_.Count } | Measure-Object -Sum).Sum
    Write-Host "    Indexed $totalSyms top-level identifiers"
    Write-Host ''

    # Discover spec files
    $files = if ($OnlyFile) {
        @(Join-Path $SpecDir $OnlyFile)
    } else {
        Get-ChildItem -Path $SpecDir -Recurse -Filter '*.md' -File | ForEach-Object { $_.FullName }
    }

    $totalRefs = 0
    $pkgFabrications = New-Object System.Collections.Generic.List[object]
    $symFabrications = New-Object System.Collections.Generic.List[object]
    $okRefs = 0

    foreach ($file in $files) {
        $rel = $file -replace [regex]::Escape((Resolve-Path .).Path), '.'
        $refs = Get-SpecApiReferences -Path $file
        $totalRefs += $refs.Count

        foreach ($ref in $refs) {
            if (-not $pkgMap.ContainsKey($ref.Package)) {
                $pkgFabrications.Add([pscustomobject]@{
                    File    = $rel
                    Line    = $ref.Line
                    Package = $ref.Package
                    Symbol  = $ref.Symbol
                    Context = $ref.Context
                })
                continue
            }
            if (-not $symMap[$ref.Package].Contains($ref.Symbol)) {
                $symFabrications.Add([pscustomobject]@{
                    File    = $rel
                    Line    = $ref.Line
                    Package = $ref.Package
                    Symbol  = $ref.Symbol
                    Context = $ref.Context
                })
                continue
            }
            $okRefs++
        }
    }

    # ── Report ──
    Write-Host ('  ─' * 30) -ForegroundColor DarkGray
    Write-Host '  ▶ Results' -ForegroundColor Cyan
    Write-Host "    Files scanned:   $($files.Count)"
    Write-Host "    Total refs:      $totalRefs"
    Write-Host "    Resolved:        $okRefs" -ForegroundColor Green
    Write-Host "    Pkg fabrications:$($pkgFabrications.Count)" -ForegroundColor $(if ($pkgFabrications.Count) { 'Red' } else { 'Green' })
    Write-Host "    Sym fabrications:$($symFabrications.Count)" -ForegroundColor $(if ($symFabrications.Count) { 'Red' } else { 'Green' })
    Write-Host ''

    if ($pkgFabrications.Count -gt 0) {
        Write-Host '  ✗ FABRICATED PACKAGES (no such directory in upstream):' -ForegroundColor Red
        $pkgFabrications | Group-Object Package | Sort-Object Count -Descending | ForEach-Object {
            Write-Host "    • $($_.Name)  ($($_.Count) refs)" -ForegroundColor Red
            $_.Group | Select-Object -First 3 | ForEach-Object {
                Write-Host "        $($_.File):$($_.Line)  →  $($_.Package).$($_.Symbol)" -ForegroundColor DarkRed
            }
            if ($_.Count -gt 3) {
                Write-Host "        ... ($($_.Count - 3) more)" -ForegroundColor DarkGray
            }
        }
        Write-Host ''
    }

    if ($symFabrications.Count -gt 0) {
        Write-Host '  ✗ FABRICATED SYMBOLS (package exists but symbol does not):' -ForegroundColor Red
        $symFabrications | Group-Object { "$($_.Package).$($_.Symbol)" } | Sort-Object Count -Descending | ForEach-Object {
            Write-Host "    • $($_.Name)  ($($_.Count) refs)" -ForegroundColor Red
            $_.Group | Select-Object -First 2 | ForEach-Object {
                Write-Host "        $($_.File):$($_.Line)" -ForegroundColor DarkRed
            }
        }
        Write-Host ''
    }

    if ($pkgFabrications.Count -eq 0 -and $symFabrications.Count -eq 0) {
        Write-Host '  ✓ No fabrications detected' -ForegroundColor Green
    }

    Write-Host ''

    if ($StrictExitCode -and ($pkgFabrications.Count + $symFabrications.Count) -gt 0) {
        exit 1
    }

    return [pscustomobject]@{
        FilesScanned     = $files.Count
        TotalRefs        = $totalRefs
        Resolved         = $okRefs
        PkgFabrications  = $pkgFabrications
        SymFabrications  = $symFabrications
    }
}

Export-ModuleMember -Function Invoke-SpecApiCheck, Get-UpstreamPackages, Get-UpstreamTopLevelSymbols, Get-SpecApiReferences, Test-UpstreamClone
