"""Tests for check-package-coverage.py — keeps the AS gate honest."""
from __future__ import annotations

import importlib.util
import subprocess
import sys
import textwrap
from pathlib import Path

HERE = Path(__file__).parent
SCRIPT = HERE / "check-package-coverage.py"


def _load():
    spec = importlib.util.spec_from_file_location("check_pkg_cov", SCRIPT)
    mod = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(mod)
    return mod


def _profile(tmp_path: Path, body: str) -> Path:
    p = tmp_path / "coverage.out"
    p.write_text("mode: atomic\n" + textwrap.dedent(body).lstrip())
    return p


def test_parses_per_package_totals(tmp_path):
    mod = _load()
    p = _profile(
        tmp_path,
        """
        github.com/x/enum/foo/file.go:1.1,2.1 3 1
        github.com/x/enum/foo/file.go:2.1,3.1 1 0
        github.com/x/enum/bar/file.go:1.1,2.1 2 2
        """,
    )
    pkgs = mod.parse_profile(p)
    assert pkgs["github.com/x/enum/foo"] == (3, 4)
    assert pkgs["github.com/x/enum/bar"] == (2, 2)


def test_passes_when_all_above_threshold(tmp_path):
    p = _profile(
        tmp_path,
        """
        github.com/x/enum/ok/file.go:1.1,2.1 4 1
        """,
    )
    rc = subprocess.run(
        [sys.executable, str(SCRIPT), str(p), "--threshold", "75"],
        capture_output=True, text=True,
    )
    assert rc.returncode == 0, rc.stdout + rc.stderr


def test_fails_when_any_below_threshold(tmp_path):
    p = _profile(
        tmp_path,
        """
        github.com/x/enum/ok/file.go:1.1,2.1 4 1
        github.com/x/enum/low/file.go:1.1,2.1 4 1
        github.com/x/enum/low/file.go:2.1,3.1 6 0
        """,
    )
    rc = subprocess.run(
        [sys.executable, str(SCRIPT), str(p), "--threshold", "75"],
        capture_output=True, text=True,
    )
    assert rc.returncode == 1
    assert "github.com/x/enum/low" in rc.stdout


def test_ignore_flag_excludes_package(tmp_path):
    p = _profile(
        tmp_path,
        """
        github.com/x/enum/skipme/file.go:1.1,2.1 10 0
        github.com/x/enum/ok/file.go:1.1,2.1 4 1
        """,
    )
    rc = subprocess.run(
        [sys.executable, str(SCRIPT), str(p), "--threshold", "75",
         "--ignore", "skipme"],
        capture_output=True, text=True,
    )
    assert rc.returncode == 0, rc.stdout + rc.stderr
    assert "ignored low packages" in rc.stdout


def test_missing_profile(tmp_path):
    rc = subprocess.run(
        [sys.executable, str(SCRIPT), str(tmp_path / "nope.out")],
        capture_output=True, text=True,
    )
    assert rc.returncode == 2


def test_rejects_non_profile(tmp_path):
    p = tmp_path / "bad.out"
    p.write_text("not a profile\n")
    rc = subprocess.run(
        [sys.executable, str(SCRIPT), str(p)],
        capture_output=True, text=True,
    )
    assert rc.returncode == 2
