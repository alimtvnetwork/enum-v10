"""Tests for check-stringer-recursion.py — guards RCA Pattern P9."""
from __future__ import annotations

import importlib.util
import subprocess
import sys
import tempfile
import textwrap
import unittest
from pathlib import Path

HERE = Path(__file__).parent
SCRIPT = HERE / "check-stringer-recursion.py"


def _load():
    spec = importlib.util.spec_from_file_location("check_stringer", SCRIPT)
    mod = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(mod)
    return mod


def _write_go(tmp: Path, name: str, body: str) -> Path:
    p = tmp / name
    p.write_text(textwrap.dedent(body).lstrip())
    return p


class CheckStringerRecursionTests(unittest.TestCase):
    def setUp(self):
        self._tmpdir = tempfile.TemporaryDirectory()
        self.tmp = Path(self._tmpdir.name)

    def tearDown(self):
        self._tmpdir.cleanup()

    def _run(self, *roots: Path):
        return subprocess.run(
            [sys.executable, str(SCRIPT), *map(str, roots)],
            capture_output=True, text=True,
        )

    def test_clean_source_passes(self):
        _write_go(self.tmp, "ok.go", """
            package x
            func (it Foo) String() string {
                return fmt.Sprintf("Foo(%s)", it.name)
            }
        """)
        rc = self._run(self.tmp)
        self.assertEqual(rc.returncode, 0, rc.stdout + rc.stderr)
        self.assertIn("No Pattern P9 violations", rc.stdout)

    def test_recursion_bomb_detected(self):
        _write_go(self.tmp, "bad.go", """
            package x
            func (it Foo) String() string {
                return converters.AnyTo.ValueString(it)
            }
        """)
        rc = self._run(self.tmp)
        self.assertEqual(rc.returncode, 1)
        self.assertIn("Pattern P9 violation", rc.stdout)
        self.assertIn("bad.go", rc.stdout)

    def test_test_files_are_excluded(self):
        # _test.go file with the bomb — should be ignored.
        _write_go(self.tmp, "fine_test.go", """
            package x
            func (it Foo) String() string {
                return converters.AnyTo.ValueString(it)
            }
        """)
        # Plus a clean non-test file so we have something to scan.
        _write_go(self.tmp, "clean.go", """
            package x
            func Hello() string { return "hi" }
        """)
        rc = self._run(self.tmp)
        self.assertEqual(rc.returncode, 0, rc.stdout + rc.stderr)

    def test_call_outside_string_method_is_ok(self):
        _write_go(self.tmp, "elsewhere.go", """
            package x
            func (it Foo) Render() string {
                return converters.AnyTo.ValueString(it.payload)
            }
            func (it Foo) String() string {
                return it.name
            }
        """)
        rc = self._run(self.tmp)
        self.assertEqual(rc.returncode, 0, rc.stdout + rc.stderr)

    def test_pointer_receiver_detected(self):
        _write_go(self.tmp, "ptr.go", """
            package x
            func (it *Foo) String() string {
                if it == nil { return "" }
                return converters.AnyTo.ValueString(it)
            }
        """)
        rc = self._run(self.tmp)
        self.assertEqual(rc.returncode, 1)
        self.assertIn("ptr.go", rc.stdout)

    def test_brace_in_string_literal_does_not_confuse_scanner(self):
        _write_go(self.tmp, "tricky.go", """
            package x
            func (it Foo) Other() string {
                return "}{ converters.AnyTo.ValueString"
            }
            func (it Foo) String() string {
                return it.name
            }
        """)
        rc = self._run(self.tmp)
        self.assertEqual(rc.returncode, 0, rc.stdout + rc.stderr)

    def test_missing_root(self):
        rc = self._run(self.tmp / "nope")
        self.assertEqual(rc.returncode, 2)

    def test_no_go_files(self):
        rc = self._run(self.tmp)
        self.assertEqual(rc.returncode, 2)


if __name__ == "__main__":
    unittest.main()
