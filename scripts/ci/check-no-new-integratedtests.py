#!/usr/bin/env python3
"""
BC guard: lint that fails when new hardcoded `tests/integratedtests/` (or bare
`integratedtests`) strings appear in source/tooling outside the documented
allowlist.

Locks Task BB so future PRs cannot reintroduce hardcoded paths in PowerShell /
Go / YAML / Python sources. Spec docs, audit memory, and the cross-repo mirror
are intentionally out of scope (those legitimately discuss the upstream layout).

Allowlisted files have a per-file occurrence baseline. The check fails if:
  * a NEW file (outside the allowlist) contains the token, OR
  * an allowlisted file's occurrence count GROWS beyond its baseline.

Counts may shrink (e.g. when BE removes the stale CoverageReportHtml string) —
shrinkage just means the baseline can be tightened in a follow-up PR.

Run locally:
  python3 scripts/ci/check-no-new-integratedtests.py
Run self-tests:
  python3 scripts/ci/test_check_no_new_integratedtests.py
"""
from __future__ import annotations

import os
import re
import sys
from pathlib import Path

# Per-file baseline: max allowed occurrences of the token "integratedtests".
# Each entry documents WHY the file legitimately mentions the legacy layout.
ALLOWLIST: dict[str, int] = {
    # Dual-layout regex matchers / fallback arrays (intentional):
    "scripts/CoverageCompileCheck.psm1": 2,    # comment + regex strip
    "scripts/CoverageProfileMerger.psm1": 1,   # regex strip on profile paths
    "scripts/CoverageRunner.psm1": 4,          # comment + dual-root array + 2 regex strips
    "scripts/CoverageSplitRecovery.psm1": 4,   # legacy split-recovery path probes
    # Helper that defines the dual-layout abstraction:
    "scripts/Utilities.psm1": 4,               # Resolve-TestSuiteRoot doc + candidate list
    # User-facing diagnostics / doc strings (display-only):
    "scripts/CoveragePreChecks.psm1": 1,       # Write-Fail message text
    # CoverageReportHtml.psm1: stale instruction string removed in BE — must stay at 0.
    "scripts/Help.psm1": 3,                    # `it` cmdlet help text + legacy ./integratedtests/... fallback
    "scripts/PackageCoverage.psm1": 3,         # comment + warning message text
    "scripts/PreCommitCheck.psm1": 3,          # Core-memory comment + 2 regex strips
    "scripts/TestRunner.psm1": 1,              # cmdlet help text
    # Smoke test that exercises the legacy-only branch of Resolve-TestSuiteRoot:
    "tests/scripts/Test-ResolveTestSuiteRoot.ps1": 3,
    # The guard itself + its self-tests must mention the token:
    "scripts/ci/check-no-new-integratedtests.py": 13,
    "scripts/ci/test_check_no_new_integratedtests.py": 9,
}

# Directories scanned (relative to repo root). `cross-repo/` mirrors a separate
# upstream repo; `spec/`, `.lovable/`, `.release/` are documentation/memory that
# legitimately discuss the upstream layout.
SCAN_ROOTS = ("scripts", "tests", "osdetect", ".github", "Cmd")
SKIP_PATH_FRAGMENTS = (
    f"{os.sep}.git{os.sep}",
    f"{os.sep}node_modules{os.sep}",
    f"{os.sep}__pycache__{os.sep}",
)
SCAN_EXTENSIONS = {
    ".ps1", ".psm1", ".psd1",
    ".go",
    ".yml", ".yaml",
    ".py",
    ".json",
    ".sh", ".bash",
    ".md",
    ".txt",
}

TOKEN = re.compile(r"integratedtests")


def repo_root() -> Path:
    return Path(__file__).resolve().parents[2]


def iter_source_files(root: Path):
    for top in SCAN_ROOTS:
        base = root / top
        if not base.is_dir():
            continue
        for dirpath, _dirnames, filenames in os.walk(base):
            if any(frag in dirpath + os.sep for frag in SKIP_PATH_FRAGMENTS):
                continue
            for name in filenames:
                p = Path(dirpath) / name
                if p.suffix.lower() in SCAN_EXTENSIONS:
                    yield p


def count_token(path: Path) -> tuple[int, list[tuple[int, str]]]:
    try:
        text = path.read_text(encoding="utf-8", errors="replace")
    except OSError:
        return 0, []
    matches: list[tuple[int, str]] = []
    for lineno, line in enumerate(text.splitlines(), start=1):
        if TOKEN.search(line):
            matches.append((lineno, line.rstrip()))
    return len(matches), matches


def check(root: Path | None = None) -> int:
    root = root or repo_root()
    violations: list[str] = []
    for path in iter_source_files(root):
        rel = path.relative_to(root).as_posix()
        count, matches = count_token(path)
        if count == 0:
            continue
        baseline = ALLOWLIST.get(rel)
        if baseline is None:
            violations.append(
                f"NEW file with hardcoded `integratedtests` (not in allowlist): {rel} ({count} occurrences)"
            )
            for ln, txt in matches:
                violations.append(f"    {rel}:{ln}: {txt}")
        elif count > baseline:
            violations.append(
                f"Occurrences GREW in {rel}: baseline={baseline}, now={count}"
            )
            for ln, txt in matches:
                violations.append(f"    {rel}:{ln}: {txt}")
    if violations:
        print("BC guard: hardcoded `integratedtests` violations detected:\n", file=sys.stderr)
        for v in violations:
            print(v, file=sys.stderr)
        print(
            "\nIf the new reference is intentional dual-layout handling, add it to "
            "ALLOWLIST in scripts/ci/check-no-new-integratedtests.py with a one-line "
            "rationale. Otherwise, refactor to use Resolve-TestSuiteRoot from "
            "scripts/Utilities.psm1 (see Task S-114 / BB).",
            file=sys.stderr,
        )
        return 1
    print("BC guard: OK — no new hardcoded `integratedtests` references.")
    return 0


if __name__ == "__main__":
    sys.exit(check())
