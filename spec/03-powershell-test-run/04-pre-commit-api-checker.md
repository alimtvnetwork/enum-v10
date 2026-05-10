# Pre-Commit API Mismatch Checker

> **Scope note (`enum-v9`)** — the discovery glob below (`tests/integratedtests/`) is the upstream-`core-v9` layout. The actual `run.ps1 PC` implementation reads the test root from disk and works equally on `enum-v9`'s `tests/creationtests/`. JSON example paths use upstream package names (`corejsontests`) for illustration. See [`01-overview.md`](./01-overview.md) scope note and `spec/01-app/13-testing-patterns.md` §6.1.

## Problem

Coverage test files (`Coverage*_test.go`) are written against assumed method signatures of production source packages. When signatures drift (renamed methods, changed parameters, field-vs-method access), the test files fail to compile. Because the runner uses a **Build Failure Cascade** model, a single bad test package blocks the entire coverage pipeline.

These mismatches are only caught at `go test -c` time, which may be minutes into a coverage run. A pre-commit check catches them in seconds.

## Solution: `./run.ps1 PC` (Pre-Commit Check)

A targeted compile check that validates **only** `Coverage*` test files against their source packages, producing a structured JSON report of mismatches.

### Usage

```powershell
./run.ps1 PC                # Check all test packages (parallel)
./run.ps1 -pc               # Same
./run.ps1 pre-commit        # Same
./run.ps1 PC --sync         # Sequential mode
./run.ps1 PC corejsontests  # Check a single package
```

### What It Does

1. **Discovers** all packages under `tests/integratedtests/` containing `Coverage*` files.
2. **Compiles** each package via `go test -c` (no execution — build-only).
3. **Parses** compiler errors to extract:
   - File and line number
   - Error message (missing method, wrong arg count, type mismatch)
   - Affected symbol (function/method name)
4. **Outputs** summary-only console output (no per-package ✓/✗ lines — only boxed summary sections).
5. **Writes** `data/precommit/api-check.json` with structured results.

### JSON Output Schema

```json
{
  "timestamp": "2026-03-14T10:30:00Z",
  "passed": false,
  "checkedCount": 18,
  "passedCount": 16,
  "failedCount": 2,
  "source": "PreCommitCheck.psm1 → Invoke-PreCommitCheck",
  "failures": [
    {
      "package": "corejsontests",
      "errorCount": 2,
      "errors": [
        {
          "file": "Coverage2_test.go",
          "line": 14,
          "message": "too many arguments in call to result.Clone",
          "raw": "tests/integratedtests/corejsontests/Coverage2_test.go:14:2: too many arguments in call to result.Clone"
        }
      ]
    }
  ]
}
```

### Error Classification

The checker categorizes compile errors into API mismatch types:

| Pattern | Category | Typical Cause |
|---------|----------|---------------|
| `too many arguments` | `arg-count` | Method signature changed (params added/removed) |
| `not enough arguments` | `arg-count` | Missing required parameter |
| `undefined:` | `undefined` | Method/function renamed or removed |
| `cannot use .* as` | `type-mismatch` | Parameter type changed |
| `has no field or method` | `missing-member` | Field renamed or receiver changed |
| `cannot call non-function` | `field-vs-method` | Field accessed as method or vice versa |

### Console Output

The PC command produces **summary-only** console output — no per-package `✓`/`✗` lines are printed during the compile check. Only boxed summary sections appear:

**All passed:**
```
  ┌─────────────────────────────────────────────────
  │ ✓ ALL 18 PACKAGES PASSED API CHECK
  └─────────────────────────────────────────────────
```

**Failures detected:**
```
  ┌─────────────────────────────────────────────────
  │ ✗ 2 PACKAGE(S) HAVE API MISMATCHES
  │
  │   ✗ corejsontests (3 error(s))
  │   ✗ corestrtests (1 error(s))
  │
  │ Fix these before committing Coverage* files.
  └─────────────────────────────────────────────────

  ── corejsontests ──
    Coverage2_test.go:14 [arg-count] too many arguments in call to result.Clone
    Coverage2_test.go:28 [undefined] undefined: corejson.BadFunc
    Coverage3_test.go:9 [missing-member] result.Items has no field or method Len
```

### Integration Points

- **Pre-commit hook**: Run `./run.ps1 PC` before committing Coverage files.
- **CI pipeline**: Add as a gate before coverage runs.
- **AI agent**: Parse `api-check.json` to auto-fix signature mismatches.
- **Coverage runner**: `./run.ps1 TC` already runs compile checks; `PC` is a fast, focused subset.

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All packages compiled successfully |
| 1 | One or more packages have API mismatches |

### Relationship to TC Compile Check

`PC` is **not** a replacement for the TC pre-coverage compile check. It is a **fast, focused** check meant to run before commits. Key differences:

| Aspect | `PC` (Pre-Commit) | `TC` (Coverage Compile) |
|--------|--------------------|-------------------------|
| Scope | Coverage* files only | All test files |
| Speed | Seconds | Minutes |
| Output | JSON + boxed summary | TXT + boxed summary |
| `-coverpkg` | No | Yes (instruments source) |
| When | Before commit | Before coverage run |
