# 04 — Bootstrap Into a New Repo

> **Purpose**: Step-by-step instructions for installing the `run.ps1` PowerShell toolchain
> and the GitHub Actions CI pipeline into a **fresh Go repository**. Generic — no `core-v9`
> assumptions baked in.
>
> **Audience**: Anyone (human or AI agent) standing up a new project that wants the same
> dashboard-style local test runner and CI gates documented in this folder.
>
> **Prerequisites**: Go 1.21+, PowerShell 7+, Git, a GitHub repo (for CI). Tested on
> Windows, macOS, and Linux PowerShell.

---

## Table of Contents

1. [Files to Copy](#1-files-to-copy)
2. [Parameters to Customize](#2-parameters-to-customize)
3. [Step-by-Step Install](#3-step-by-step-install)
4. [Verifying the Install](#4-verifying-the-install)
5. [Common Pitfalls](#5-common-pitfalls)
6. [Optional Add-Ons](#6-optional-add-ons)
7. [Decoupling from `core-v9` Assumptions](#7-decoupling-from-core-v9-assumptions)

---

## 1. Files to Copy

The toolchain has three layers. Copy all of them from the source repo into the destination repo at the **same relative paths**:

### 1.1 Entry-point dispatcher

| Source | Destination | Purpose |
|---|---|---|
| `run.ps1` | `<repo>/run.ps1` | Thin (~167 lines) dispatcher that imports the modules and routes commands |

### 1.2 PowerShell modules (`scripts/*.psm1`)

Copy the **entire `scripts/` folder**. The 9 modules and their dependency graph are documented in [`03-powershell-implementation.md` §1](./03-powershell-implementation.md#1-module-architecture):

| Module | Required? | Notes |
|---|---|---|
| `DashboardUI.psm1` | ✅ required | ANSI rendering, phase tracking |
| `Utilities.psm1` | ✅ required | Console + error extraction |
| `TestLogWriter.psm1` | ✅ required | `go test` output → log files |
| `TestRunner.psm1` | ✅ required | Test execution, build checks |
| `CoverageRunner.psm1` | ✅ required | TC + TCP coverage pipelines |
| `BuildTools.psm1` | ✅ required | Build/format/vet/tidy/clean |
| `PreCommitCheck.psm1` | ✅ required | PC validation pipeline |
| `Help.psm1` | ✅ required | Help + fail log viewer |
| `GoConvey.psm1` | ⚪ optional | Browser test runner; drop if unused |

> **Rule**: Copy them as a unit. The import order in `run.ps1` must match the dependency graph (DashboardUI → Utilities → TestLogWriter → TestRunner → CoverageRunner → BuildTools → GoConvey → PreCommitCheck → Help). Do not reorder.

### 1.3 CI workflow

| Source | Destination | Purpose |
|---|---|---|
| `.github/workflows/ci.yml` | `<repo>/.github/workflows/ci.yml` | Runs lint + vet + test + coverage + govulncheck on every push/PR. See [`01-ci-pipeline.md`](./01-ci-pipeline.md) |

### 1.4 Optional design-token includes

If you want the dashboard banner colors and box-drawing to match exactly, also copy any helper PowerShell scripts referenced in [`02-powershell-dashboard-ui.md` §2.3](./02-powershell-dashboard-ui.md#23-variable-definitions-copy-paste-block). Most installs do not need this — the defaults in `DashboardUI.psm1` already include the full palette.

---

## 2. Parameters to Customize

After copying, **only the following are project-specific** and need editing. Everything else is generic.

| Parameter | Where it lives | What to change |
|---|---|---|
| **Module path** | `go.mod` (you create this normally) | The toolchain reads it via `go list -m`; nothing in `scripts/` is hard-coded to `github.com/alimtvnetwork/core-v9` |
| **Package list for coverage** | Auto-discovered via `go list ./...` | No edit required — works for any module layout |
| **Coverage threshold** | `.github/workflows/ci.yml` (env block) | Default is `60`. Adjust to your standard |
| **Lint config** | `.golangci.yml` (project root) | Bring your own; the CI workflow only references it by name |
| **Test parallelism** | `.github/workflows/ci.yml` test step | Default `-parallel=4`. Lower for memory-constrained runners |
| **Repository name in dashboard banner** | `scripts/DashboardUI.psm1` (header banner function) | Optional cosmetic; defaults to current directory name |
| **Phase registry** | `scripts/CoverageRunner.psm1` & `PreCommitCheck.psm1` | Add/remove phases if your project has extra steps (rare). See [`02-powershell-dashboard-ui.md` §12.1](./02-powershell-dashboard-ui.md#121-phase-registry) |

> **Anti-pattern**: Do **not** edit module names, function names, or the import order in `run.ps1`. The modules are coupled by exported function names; renaming breaks the dispatcher.

---

## 3. Step-by-Step Install

Run these from the destination repo root. All commands are PowerShell 7 (`pwsh`) compatible.

### Step 1 — Confirm prerequisites

```powershell
go version              # need 1.21+
pwsh --version          # need 7.0+
git --version           # any modern version
```

### Step 2 — Initialize Go module (if not already)

```powershell
cd <new-repo-root>
go mod init github.com/<owner>/<repo>
```

### Step 3 — Copy the toolchain

From the source repo (`core-v9` or any repo with this toolchain installed):

```powershell
# Adjust $src to point at the source repo
$src = "C:\path\to\core-v9"
$dst = (Get-Location).Path

Copy-Item "$src\run.ps1"                       "$dst\run.ps1"
Copy-Item "$src\scripts"                       "$dst\scripts"          -Recurse
New-Item  "$dst\.github\workflows" -ItemType Directory -Force | Out-Null
Copy-Item "$src\.github\workflows\ci.yml"      "$dst\.github\workflows\ci.yml"
```

On macOS / Linux:

```bash
SRC=/path/to/core-v9
cp    "$SRC/run.ps1"                       ./run.ps1
cp -r "$SRC/scripts"                       ./scripts
mkdir -p .github/workflows
cp    "$SRC/.github/workflows/ci.yml"      ./.github/workflows/ci.yml
```

### Step 4 — (Optional) Bring your linter config

```powershell
Copy-Item "$src\.golangci.yml" "$dst\.golangci.yml"
```

If you skip this, the CI lint step will fail until you add a `.golangci.yml`.

### Step 5 — Set the execution policy (Windows only, once per machine)

```powershell
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned
```

### Step 6 — First-run smoke test

```powershell
pwsh ./run.ps1 help
```

You should see the dashboard-styled help banner and command list. If you do, the modules loaded successfully.

### Step 7 — Run the build + test pipeline locally

```powershell
pwsh ./run.ps1 tc          # full Test+Coverage pipeline
pwsh ./run.ps1 pc          # Pre-Commit validation pipeline
```

### Step 8 — Commit and push

```powershell
git add run.ps1 scripts/ .github/workflows/ci.yml .golangci.yml
git commit -m "chore: install run.ps1 toolchain + CI"
git push
```

CI should run automatically on the push. Watch the Actions tab.

---

## 4. Verifying the Install

After Step 8, confirm each layer is wired correctly:

| Check | Expected |
|---|---|
| `pwsh ./run.ps1 help` | Dashboard banner + command list (no module-load errors) |
| `pwsh ./run.ps1 tc` | All 10 TC phases render with `[OK]`, `[SKIP]`, or `[FAIL]` markers |
| `pwsh ./run.ps1 pc` | All 5 PC phases render |
| `pwsh ./run.ps1 build` | `go build ./...` succeeds |
| GitHub Actions tab | First push triggers `lint`, `vet`, `test`, `vulncheck` jobs |
| Coverage % in dashboard | Matches `go test -coverprofile=coverage.out ./...` standalone |

If `tc` fails immediately with "module not found", the import order in `run.ps1` was broken — re-copy from source.

---

## 5. Common Pitfalls

### 5.1 PowerShell 5 vs 7

The toolchain uses PS7-only syntax (`$using:`, ternary `?:`, ANSI escape sequences via `` `e ``). Windows PowerShell 5.1 will throw parse errors. **Always run with `pwsh`**, never `powershell.exe`.

### 5.2 ExecutionPolicy on Windows

If `./run.ps1` errors with "running scripts is disabled on this system":

```powershell
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned
```

### 5.3 ANSI rendering on legacy terminals

`cmd.exe` and Git Bash mintty render the dashboard correctly. Older Windows Terminal versions (< 1.10) garble the box-drawing characters. Upgrade Terminal or use VS Code's integrated terminal.

### 5.4 Skipping required modules

You cannot drop `DashboardUI.psm1` "to keep things minimal" — `Utilities.psm1` falls back to it for ANSI. The only safely-droppable module is `GoConvey.psm1`.

### 5.5 Hard-coded coverage threshold

The CI threshold lives in `ci.yml`, not in any `.psm1`. Local `tc` runs do not enforce it; only CI does. Adjust per project.

### 5.6 Missing `.golangci.yml`

If you skip Step 4 and don't add a config, CI lint will fail. Either add a config or comment out the lint job in `ci.yml`.

---

## 6. Optional Add-Ons

These are common follow-ups but **not required** for the toolchain to work:

| Add-on | Purpose | Where |
|---|---|---|
| Pre-commit git hook calling `pwsh ./run.ps1 pc` | Block commits that fail `pc` | `.git/hooks/pre-commit` |
| `Makefile` aliases (`make tc`, `make pc`) | Friendly to non-PowerShell users | repo root |
| GoConvey browser runner | Interactive coverage browser | enable `GoConvey.psm1` + `pwsh ./run.ps1 convey` |
| Codecov / Codecov Cloud upload | External coverage tracking | extra step in `ci.yml` after `coverage.out` is produced |
| Cache `~/go/pkg/mod` in CI | Faster CI runs | `actions/cache@v4` step in `ci.yml` |

---

## 7. Decoupling from `core-v9` Assumptions

The toolchain was extracted from `core-v9` but is **not coupled** to it. Specifically:

| Assumption you might worry about | Reality |
|---|---|
| Module path hard-coded | ❌ No — read from `go.mod` via `go list -m` |
| Package list hard-coded | ❌ No — discovered via `go list ./...` |
| `coretests` framework required | ❌ No — toolchain only invokes `go test`, doesn't care which framework you use |
| `tests/integratedtests/` mirror layout required | ❌ No — works with any test layout `go test ./...` accepts (e.g. `enum-v9` uses `tests/creationtests/`; upstream `core-v9` uses `tests/integratedtests/<pkg>tests/`) |
| Specific Go version required | ⚠️ 1.21+ recommended for generics support in your code; toolchain itself is version-agnostic |
| `args.Map` / `BaseTestCase` patterns required | ❌ No — those are `core-v9` test-style choices, completely independent of the runner |
| Phase names must match `core-v9` | ⚪ No, but renaming requires updating `02-powershell-dashboard-ui.md` §12.1 in your fork |

> **Rule of thumb**: If you can run `go test ./... -coverprofile=...` in your repo, the toolchain will work.

---

## See Also

- [`01-ci-pipeline.md`](./01-ci-pipeline.md) — what the CI workflow does, stage by stage
- [`02-powershell-dashboard-ui.md`](./02-powershell-dashboard-ui.md) — visual design tokens and phase rendering
- [`03-powershell-implementation.md`](./03-powershell-implementation.md) — internal module architecture (read this if you need to extend or debug)
- [`/spec/03-powershell-test-run/01-overview.md`](../03-powershell-test-run/01-overview.md) — user-facing command reference for `run.ps1`
