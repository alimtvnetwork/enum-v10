#!/usr/bin/env python3
"""
check-package-coverage-trend.py — warn-only gate that flags per-package
coverage regressions vs a committed baseline snapshot.

Pairs with `check-package-coverage.py` (the hard floor at 75%):
  - The floor catches *absolute* drops below threshold.
  - This script catches *relative* drops between PRs even when both
    sides remain above the floor — e.g. a package sliding 95% → 80%
    is invisible to the floor gate but is a real signal.

Usage:
    check-package-coverage-trend.py <coverage.out> [--baseline FILE]
                                                   [--tolerance PP]
                                                   [--write]

Modes
-----

Seeding (baseline missing or empty):
    Prints the current per-package snapshot. If --write is passed,
    persists `{pkg: pct}` to the baseline path. Always exits 0.

Gating (baseline present and non-empty):
    For every package present in BOTH current + baseline, computes
    drop = baseline_pct - current_pct. If drop > tolerance, emits one
    `::warning::` per regressed package. New packages are noted but
    do not warn. Removed packages are noted but do not warn.

Exit codes
----------
    0  always (warning-only by design — this is a trend signal,
       not a hard gate; the floor gate is the hard gate).
    2  invalid input (missing profile, malformed JSON baseline).

Why warning-only
----------------
A trend warning that fails the build creates pressure to *delete*
coverage rather than fix root causes. Keeping it visible-but-non-blocking
lets reviewers triage without weaponising the signal.
"""
from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

# Reuse the profile parser from the floor gate so both scripts agree on
# what a "package" is (directory portion of the file path).
SCRIPT_DIR = Path(__file__).resolve().parent
sys.path.insert(0, str(SCRIPT_DIR))
from importlib import import_module  # noqa: E402

_floor = import_module("check-package-coverage")  # type: ignore[assignment]
parse_profile = _floor.parse_profile


def snapshot_from_profile(profile: Path) -> dict[str, float]:
    """Return {pkg: percent} rounded to 1dp — matches what reviewers see."""
    out: dict[str, float] = {}
    for pkg, (covered, total) in parse_profile(profile).items():
        if total <= 0:
            continue
        out[pkg] = round(covered / total * 100.0, 1)
    return out


def load_baseline(path: Path) -> dict[str, float]:
    if not path.exists() or path.stat().st_size == 0:
        return {}
    try:
        raw = json.loads(path.read_text())
    except json.JSONDecodeError as exc:
        print(f"::error::{path}: invalid JSON: {exc}", file=sys.stderr)
        sys.exit(2)
    if not isinstance(raw, dict):
        print(f"::error::{path}: baseline must be a JSON object", file=sys.stderr)
        sys.exit(2)
    out: dict[str, float] = {}
    for k, v in raw.items():
        try:
            out[str(k)] = float(v)
        except (TypeError, ValueError):
            print(f"::warning::skipping baseline entry {k}={v!r}", file=sys.stderr)
    return out


def diff_snapshots(
    current: dict[str, float],
    baseline: dict[str, float],
    tolerance_pp: float,
) -> tuple[list[tuple[str, float, float, float]], list[str], list[str]]:
    """Return (regressions, new_pkgs, removed_pkgs).

    `regressions` is a list of (pkg, baseline_pct, current_pct, drop_pp)
    sorted by drop descending (worst first).
    """
    regressions: list[tuple[str, float, float, float]] = []
    new_pkgs = sorted(set(current) - set(baseline))
    removed_pkgs = sorted(set(baseline) - set(current))

    for pkg in sorted(set(current) & set(baseline)):
        cur = current[pkg]
        base = baseline[pkg]
        drop = base - cur
        if drop > tolerance_pp + 1e-9:
            regressions.append((pkg, base, cur, drop))

    regressions.sort(key=lambda r: r[3], reverse=True)
    return regressions, new_pkgs, removed_pkgs


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("profile", type=Path)
    parser.add_argument(
        "--baseline",
        type=Path,
        default=Path(".ci-baselines/package-coverage.json"),
    )
    parser.add_argument(
        "--tolerance",
        type=float,
        default=1.0,
        help="Allowed drop in percentage points before a warning is emitted (default 1.0).",
    )
    parser.add_argument(
        "--write",
        action="store_true",
        help="After diff, overwrite baseline with the current snapshot.",
    )
    args = parser.parse_args(argv)

    if not args.profile.exists() or args.profile.stat().st_size == 0:
        print(f"::error::{args.profile} missing or empty", file=sys.stderr)
        return 2

    current = snapshot_from_profile(args.profile)
    baseline = load_baseline(args.baseline)

    print(f"Per-package coverage trend (tolerance ±{args.tolerance:.1f}pp):")
    print(f"  current snapshot: {len(current)} packages")
    print(f"  baseline:         {len(baseline)} packages ({args.baseline})")

    if not baseline:
        print("  mode: SEEDING (baseline empty or missing)")
        print("  no trend comparison — current snapshot:")
        for pkg in sorted(current):
            print(f"    {pkg}: {current[pkg]:.1f}%")
    else:
        regressions, new_pkgs, removed = diff_snapshots(
            current, baseline, args.tolerance
        )
        if new_pkgs:
            print(f"  new packages ({len(new_pkgs)}, no baseline):")
            for p in new_pkgs:
                print(f"    + {p}: {current[p]:.1f}%")
        if removed:
            print(f"  removed packages ({len(removed)}, in baseline only):")
            for p in removed:
                print(f"    - {p}: was {baseline[p]:.1f}%")
        if regressions:
            print(f"  regressions ({len(regressions)}, drop > {args.tolerance:.1f}pp):")
            for pkg, base, cur, drop in regressions:
                msg = (
                    f"{pkg}: {base:.1f}% → {cur:.1f}% (-{drop:.1f}pp)"
                )
                # GitHub Actions warning + plain log line so it shows up
                # in both the annotations panel and the raw step output.
                print(f"::warning::coverage trend regression: {msg}")
                print(f"    ! {msg}")
        else:
            print("  ✅ no regressions beyond tolerance.")

    if args.write:
        args.baseline.parent.mkdir(parents=True, exist_ok=True)
        args.baseline.write_text(
            json.dumps(current, indent=2, sort_keys=True) + "\n"
        )
        print(f"  wrote new baseline → {args.baseline}")

    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
