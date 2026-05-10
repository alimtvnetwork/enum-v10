#!/usr/bin/env python3
"""Tests for BC guard `check-no-new-integratedtests.py` (Task BC, Cycle 56)."""
from __future__ import annotations

import importlib.util
import sys
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parent
SPEC = importlib.util.spec_from_file_location(
    "bc_guard", ROOT / "check-no-new-integratedtests.py"
)
assert SPEC and SPEC.loader
bc = importlib.util.module_from_spec(SPEC)
sys.modules["bc_guard"] = bc
SPEC.loader.exec_module(bc)


def make_repo(files: dict[str, str]) -> Path:
    root = Path(tempfile.mkdtemp(prefix="bcguard-"))
    for rel, body in files.items():
        p = root / rel
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(body, encoding="utf-8")
    return root


class BcGuardTests(unittest.TestCase):
    def test_clean_repo_passes(self):
        root = make_repo({"scripts/foo.ps1": "Write-Host 'hi'\n"})
        self.assertEqual(bc.check(root), 0)

    def test_new_file_with_token_fails(self):
        root = make_repo({"scripts/new.ps1": "$x = 'tests/integratedtests/foo'\n"})
        self.assertEqual(bc.check(root), 1)

    def test_allowlisted_at_baseline_passes(self):
        body = "\n".join([f"# integratedtests line {i}" for i in range(4)]) + "\n"
        root = make_repo({"scripts/Utilities.psm1": body})
        self.assertEqual(bc.check(root), 0)

    def test_allowlisted_above_baseline_fails(self):
        body = "\n".join([f"# integratedtests line {i}" for i in range(5)]) + "\n"
        root = make_repo({"scripts/Utilities.psm1": body})
        self.assertEqual(bc.check(root), 1)

    def test_allowlisted_below_baseline_passes(self):
        # Shrinkage allowed — baseline is a ceiling, not equality.
        root = make_repo({"scripts/Utilities.psm1": "# integratedtests once\n"})
        self.assertEqual(bc.check(root), 0)

    def test_docs_and_cross_repo_ignored(self):
        root = make_repo({
            "spec/foo.md": "tests/integratedtests/everywhere\n",
            ".lovable/bar.md": "tests/integratedtests/everywhere\n",
            "cross-repo/baz.go": "// tests/integratedtests/foo\n",
        })
        self.assertEqual(bc.check(root), 0)

    def test_real_repo_passes_at_head(self):
        # Sanity: the live repo must be clean against its own baseline.
        self.assertEqual(bc.check(), 0)


if __name__ == "__main__":
    unittest.main()
