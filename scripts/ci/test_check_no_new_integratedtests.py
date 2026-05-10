#!/usr/bin/env python3
"""Self-tests for BC guard `check-no-new-integratedtests.py`."""
from __future__ import annotations

import importlib.util
import os
import sys
import tempfile
from pathlib import Path

HERE = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location(
    "bc_guard", HERE / "check-no-new-integratedtests.py"
)
assert spec and spec.loader
mod = importlib.util.module_from_spec(spec)
spec.loader.exec_module(mod)  # type: ignore[union-attr]


def make_repo(files: dict[str, str]) -> Path:
    root = Path(tempfile.mkdtemp(prefix="bcguard-"))
    for rel, body in files.items():
        p = root / rel
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(body, encoding="utf-8")
    return root


def case(name: str, ok: bool):
    print(f"  {'PASS' if ok else 'FAIL'}: {name}")
    return 0 if ok else 1


def main() -> int:
    failures = 0

    # Case 1: clean repo (no occurrences) → exit 0
    root = make_repo({"scripts/foo.ps1": "Write-Host 'hi'\n"})
    rc = mod.check(root)
    failures += case("clean repo passes", rc == 0)

    # Case 2: new (non-allowlisted) file mentions the token → exit 1
    root = make_repo({"scripts/new.ps1": "$x = 'tests/integratedtests/foo'\n"})
    rc = mod.check(root)
    failures += case("new file with token fails", rc == 1)

    # Case 3: allowlisted file at baseline → exit 0
    # Use Utilities.psm1 (baseline 4) with exactly 4 occurrences.
    body = "\n".join([f"# integratedtests line {i}" for i in range(4)]) + "\n"
    root = make_repo({"scripts/Utilities.psm1": body})
    rc = mod.check(root)
    failures += case("allowlisted at baseline passes", rc == 0)

    # Case 4: allowlisted file ABOVE baseline → exit 1
    body = "\n".join([f"# integratedtests line {i}" for i in range(5)]) + "\n"
    root = make_repo({"scripts/Utilities.psm1": body})
    rc = mod.check(root)
    failures += case("allowlisted above baseline fails", rc == 1)

    # Case 5: allowlisted file BELOW baseline → exit 0 (shrinkage allowed)
    body = "# integratedtests once\n"
    root = make_repo({"scripts/Utilities.psm1": body})
    rc = mod.check(root)
    failures += case("allowlisted below baseline passes", rc == 0)

    # Case 6: token in spec/.lovable/cross-repo is ignored
    root = make_repo({
        "spec/foo.md": "tests/integratedtests/everywhere\n",
        ".lovable/bar.md": "tests/integratedtests/everywhere\n",
        "cross-repo/baz.go": "// tests/integratedtests/foo\n",
    })
    rc = mod.check(root)
    failures += case("docs/cross-repo are ignored", rc == 0)

    # Case 7: real repo (current state) → exit 0 (sanity)
    rc = mod.check()
    failures += case("real repo passes at HEAD", rc == 0)

    if failures:
        print(f"\n{failures} test(s) failed", file=sys.stderr)
        return 1
    print("\nAll BC guard self-tests passed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
