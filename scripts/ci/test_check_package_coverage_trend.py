"""Tests for check-package-coverage-trend (AW, Cycle 53)."""
from __future__ import annotations

import importlib.util
import io
import json
import sys
import tempfile
import unittest
from contextlib import redirect_stdout, redirect_stderr
from pathlib import Path

ROOT = Path(__file__).resolve().parent
SPEC = importlib.util.spec_from_file_location(
    "check_package_coverage_trend",
    ROOT / "check-package-coverage-trend.py",
)
assert SPEC and SPEC.loader
trend = importlib.util.module_from_spec(SPEC)
sys.modules["check_package_coverage_trend"] = trend
SPEC.loader.exec_module(trend)


def write_profile(tmp: Path, lines: list[str]) -> Path:
    p = tmp / "coverage.out"
    p.write_text("mode: atomic\n" + "\n".join(lines) + "\n")
    return p


def write_baseline(tmp: Path, data: dict[str, float]) -> Path:
    p = tmp / "baseline.json"
    p.write_text(json.dumps(data))
    return p


# Profile lines covering 4 statements: covered/total counts let us
# control the percentage exactly. Format:
#   <pkg>/<file>.go:<sl>.<sc>,<el>.<ec> <numStmt> <count>
def line(pkg: str, num_stmt: int, count: int, idx: int = 1) -> str:
    return f"{pkg}/file.go:{idx}.1,{idx}.99 {num_stmt} {count}"


class TestSnapshotFromProfile(unittest.TestCase):
    def test_percentage_is_rounded_to_one_dp(self):
        with tempfile.TemporaryDirectory() as td:
            tmp = Path(td)
            # 1 of 3 covered = 33.333...% -> 33.3
            prof = write_profile(tmp, [line("foo", 1, 0, 1), line("foo", 2, 1, 2)])
            snap = trend.snapshot_from_profile(prof)
            self.assertEqual(snap, {"foo": round(2 / 3 * 100, 1)})

    def test_zero_statement_packages_skipped(self):
        with tempfile.TemporaryDirectory() as td:
            tmp = Path(td)
            prof = write_profile(tmp, [line("zero", 0, 0, 1), line("good", 1, 1, 2)])
            snap = trend.snapshot_from_profile(prof)
            self.assertNotIn("zero", snap)
            self.assertEqual(snap["good"], 100.0)


class TestLoadBaseline(unittest.TestCase):
    def test_missing_baseline_returns_empty(self):
        self.assertEqual(trend.load_baseline(Path("/nonexistent/path.json")), {})

    def test_empty_file_returns_empty(self):
        with tempfile.TemporaryDirectory() as td:
            p = Path(td) / "empty.json"
            p.write_text("")
            self.assertEqual(trend.load_baseline(p), {})

    def test_invalid_json_exits_2(self):
        with tempfile.TemporaryDirectory() as td:
            p = Path(td) / "bad.json"
            p.write_text("{not json")
            with self.assertRaises(SystemExit) as cm, redirect_stderr(io.StringIO()):
                trend.load_baseline(p)
            self.assertEqual(cm.exception.code, 2)

    def test_non_object_exits_2(self):
        with tempfile.TemporaryDirectory() as td:
            p = Path(td) / "list.json"
            p.write_text("[]")
            with self.assertRaises(SystemExit) as cm, redirect_stderr(io.StringIO()):
                trend.load_baseline(p)
            self.assertEqual(cm.exception.code, 2)

    def test_skips_unparseable_entries(self):
        with tempfile.TemporaryDirectory() as td:
            p = Path(td) / "mixed.json"
            p.write_text(json.dumps({"good": 80.5, "bad": "not-a-number"}))
            with redirect_stderr(io.StringIO()):
                loaded = trend.load_baseline(p)
            self.assertEqual(loaded, {"good": 80.5})


class TestDiffSnapshots(unittest.TestCase):
    def test_drop_within_tolerance_is_not_regression(self):
        regs, new, rm = trend.diff_snapshots(
            current={"a": 79.5}, baseline={"a": 80.0}, tolerance_pp=1.0
        )
        self.assertEqual(regs, [])
        self.assertEqual(new, [])
        self.assertEqual(rm, [])

    def test_drop_beyond_tolerance_is_regression(self):
        regs, _, _ = trend.diff_snapshots(
            current={"a": 70.0}, baseline={"a": 90.0}, tolerance_pp=1.0
        )
        self.assertEqual(len(regs), 1)
        pkg, base, cur, drop = regs[0]
        self.assertEqual(pkg, "a")
        self.assertEqual(base, 90.0)
        self.assertEqual(cur, 70.0)
        self.assertAlmostEqual(drop, 20.0, places=6)

    def test_improvements_never_warn(self):
        regs, _, _ = trend.diff_snapshots(
            current={"a": 95.0}, baseline={"a": 80.0}, tolerance_pp=1.0
        )
        self.assertEqual(regs, [])

    def test_new_and_removed_packages_classified(self):
        regs, new, rm = trend.diff_snapshots(
            current={"a": 80.0, "new": 50.0},
            baseline={"a": 80.0, "removed": 90.0},
            tolerance_pp=1.0,
        )
        self.assertEqual(regs, [])
        self.assertEqual(new, ["new"])
        self.assertEqual(rm, ["removed"])

    def test_regressions_sorted_worst_first(self):
        regs, _, _ = trend.diff_snapshots(
            current={"small": 79.0, "big": 50.0},
            baseline={"small": 80.0, "big": 90.0},
            tolerance_pp=0.0,
        )
        self.assertEqual([r[0] for r in regs], ["big", "small"])


class TestMainEndToEnd(unittest.TestCase):
    def test_seeding_mode_exits_zero_and_does_not_warn(self):
        with tempfile.TemporaryDirectory() as td:
            tmp = Path(td)
            prof = write_profile(tmp, [line("pkg", 1, 1, 1)])
            buf = io.StringIO()
            with redirect_stdout(buf):
                rc = trend.main([str(prof), "--baseline", str(tmp / "missing.json")])
            self.assertEqual(rc, 0)
            self.assertIn("SEEDING", buf.getvalue())
            self.assertNotIn("::warning::", buf.getvalue())

    def test_gating_mode_emits_warning_for_regression(self):
        with tempfile.TemporaryDirectory() as td:
            tmp = Path(td)
            # 2 of 4 covered -> 50.0%
            prof = write_profile(
                tmp, [line("pkg", 2, 1, 1), line("pkg", 2, 0, 2)]
            )
            base = write_baseline(tmp, {"pkg": 90.0})
            buf = io.StringIO()
            with redirect_stdout(buf):
                rc = trend.main([str(prof), "--baseline", str(base)])
            self.assertEqual(rc, 0)  # warning-only — never fail
            out = buf.getvalue()
            self.assertIn("::warning::coverage trend regression", out)
            self.assertIn("pkg: 90.0% → 50.0%", out)

    def test_write_persists_current_snapshot(self):
        with tempfile.TemporaryDirectory() as td:
            tmp = Path(td)
            prof = write_profile(tmp, [line("pkg", 1, 1, 1)])
            base = tmp / "baseline.json"
            with redirect_stdout(io.StringIO()):
                rc = trend.main([str(prof), "--baseline", str(base), "--write"])
            self.assertEqual(rc, 0)
            self.assertEqual(json.loads(base.read_text()), {"pkg": 100.0})

    def test_missing_profile_exits_2(self):
        buf = io.StringIO()
        with redirect_stderr(buf):
            rc = trend.main(["/nonexistent/coverage.out"])
        self.assertEqual(rc, 2)


if __name__ == "__main__":
    unittest.main()
