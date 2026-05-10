#!/usr/bin/env python3
"""
check-package-coverage.py — enforce a per-package statement-coverage floor.

Reads a Go cover profile (the file produced by `go test -coverprofile=...`),
aggregates statement counts per package, and exits non-zero if any package
falls below the threshold.

Usage:
    check-package-coverage.py <coverage.out> [--threshold N] [--ignore PKG ...]

Exit codes:
    0  all packages meet the floor (or only ignored packages were below)
    1  one or more non-ignored packages below the floor
    2  invalid / missing input

Profile format (each non-header line):
    <import-path>/<file>.go:<startLine>.<startCol>,<endLine>.<endCol>
        <numStmt> <count>

The package key is the directory portion of the file path
(everything before the last `/`).

Locks in the gains from Task AO so future PRs cannot silently regress
any package below 75% statement coverage.
"""
from __future__ import annotations

import argparse
import collections
import sys
from pathlib import Path


def parse_profile(path: Path) -> dict[str, tuple[int, int]]:
    """Return {package_import_path: (covered_statements, total_statements)}."""
    totals: dict[str, list[int]] = collections.defaultdict(lambda: [0, 0])

    with path.open() as fh:
        first = fh.readline()
        if not first.startswith("mode:"):
            print(
                f"::error::{path}: missing 'mode:' header — not a Go cover profile",
                file=sys.stderr,
            )
            sys.exit(2)
        for raw in fh:
            line = raw.strip()
            if not line:
                continue
            try:
                location, num_stmt_s, count_s = line.rsplit(" ", 2)
                num_stmt = int(num_stmt_s)
                count = int(count_s)
            except ValueError:
                print(f"::warning::skipping malformed profile line: {line}", file=sys.stderr)
                continue

            file_part = location.split(":", 1)[0]
            slash = file_part.rfind("/")
            pkg = file_part[:slash] if slash != -1 else file_part

            totals[pkg][1] += num_stmt
            if count > 0:
                totals[pkg][0] += num_stmt

    return {p: (c, t) for p, (c, t) in totals.items()}


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("profile", type=Path)
    parser.add_argument("--threshold", type=float, default=75.0)
    parser.add_argument(
        "--ignore",
        action="append",
        default=[],
        help="Import path (or suffix) to skip. May be repeated.",
    )
    args = parser.parse_args(argv)

    if not args.profile.exists() or args.profile.stat().st_size == 0:
        print(f"::error::{args.profile} missing or empty", file=sys.stderr)
        return 2

    pkgs = parse_profile(args.profile)
    if not pkgs:
        print("::error::no packages found in profile", file=sys.stderr)
        return 2

    failures: list[tuple[str, float, int, int]] = []
    skipped: list[str] = []
    rows: list[tuple[str, float, int, int]] = []

    for pkg in sorted(pkgs):
        covered, total = pkgs[pkg]
        pct = (covered / total * 100.0) if total > 0 else 100.0
        rows.append((pkg, pct, covered, total))
        if any(pkg.endswith(ig) or pkg == ig for ig in args.ignore):
            if pct < args.threshold:
                skipped.append(f"{pkg} ({pct:.1f}%)")
            continue
        if total == 0:
            continue  # data-only package, no statements to cover
        if pct + 1e-9 < args.threshold:
            failures.append((pkg, pct, covered, total))

    print(f"Per-package coverage gate (threshold {args.threshold:.1f}%):")
    print(f"  packages scanned: {len(pkgs)}")
    print(f"  failures:         {len(failures)}")
    print(f"  ignored<floor:    {len(skipped)}")
    if skipped:
        print("  ignored low packages:")
        for s in skipped:
            print(f"    - {s}")

    if failures:
        print("\n::error::Per-package coverage gate failed:")
        for pkg, pct, c, t in failures:
            print(f"  - {pkg}: {pct:.1f}%  ({c}/{t})")
        return 1

    print("✅ All non-ignored packages meet the floor.")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
