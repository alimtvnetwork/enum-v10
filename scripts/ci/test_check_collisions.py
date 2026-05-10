#!/usr/bin/env python3
"""
Unit tests for `scripts/ci/check-collisions.py`.

Run:
    python3 -m unittest scripts/ci/test_check_collisions.py -v

Covers:
  - GOOS/GOARCH build-tag sibling collapsing (`build_tag_key`)
  - Exported/unexported accessor-pair allowance (`is_exported_unexported_pair`)
  - Top-level decl parsing (`parse_decls`): func / type / var / const,
    plus `const ( ... )` and `var ( ... )` blocks
  - String-literal & comment skipping (raw strings, line comments,
    block comments, double-quoted strings containing `func`)
  - End-to-end `main()` exit codes:
      0  no collisions
      1  cross-file or case-insensitive collision detected
      2  no .go files in the search root
  - Per-package scoping (same name in different packages is OK)
  - Intra-file duplicate detection
"""
from __future__ import annotations

import importlib.util
import io
import sys
import tempfile
import unittest
from contextlib import redirect_stdout, redirect_stderr
from pathlib import Path

# Load `check-collisions.py` (hyphen in name → can't use plain `import`).
_HERE = Path(__file__).resolve().parent
_SPEC = importlib.util.spec_from_file_location(
    "check_collisions", _HERE / "check-collisions.py"
)
cc = importlib.util.module_from_spec(_SPEC)
assert _SPEC.loader is not None
_SPEC.loader.exec_module(cc)


def write(root: Path, rel: str, body: str) -> Path:
    p = root / rel
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(body, encoding="utf-8")
    return p


class TestBuildTagKey(unittest.TestCase):
    def test_strips_single_goos_suffix(self):
        self.assertEqual(
            cc.build_tag_key(Path("pkg/Foo_linux.go")),
            str(Path("pkg/Foo.go")),
        )

    def test_strips_goos_and_goarch(self):
        self.assertEqual(
            cc.build_tag_key(Path("pkg/Foo_linux_amd64.go")),
            str(Path("pkg/Foo.go")),
        )

    def test_leaves_non_build_suffix_alone(self):
        # `_test` is not a GOOS/GOARCH token.
        self.assertEqual(
            cc.build_tag_key(Path("pkg/Foo_test.go")),
            str(Path("pkg/Foo_test.go")),
        )

    def test_no_suffix_unchanged(self):
        self.assertEqual(
            cc.build_tag_key(Path("pkg/Foo.go")),
            str(Path("pkg/Foo.go")),
        )


class TestAccessorPair(unittest.TestCase):
    def test_exported_unexported_pair_allowed(self):
        self.assertTrue(cc.is_exported_unexported_pair(["Foo", "foo"]))

    def test_three_variants_not_a_pair(self):
        self.assertFalse(cc.is_exported_unexported_pair(["Foo", "foo", "FOO"]))

    def test_different_letters_not_a_pair(self):
        self.assertFalse(cc.is_exported_unexported_pair(["Foo", "bar"]))

    def test_two_exported_not_a_pair(self):
        self.assertFalse(cc.is_exported_unexported_pair(["Foo", "FOO"]))


class TestParseDecls(unittest.TestCase):
    def _decls(self, body: str):
        with tempfile.TemporaryDirectory() as td:
            p = write(Path(td), "x.go", body)
            return list(cc.parse_decls(p))

    def test_func_type_var_const(self):
        decls = self._decls(
            "package p\n"
            "func Foo() {}\n"
            "type Bar struct{}\n"
            "var Baz = 1\n"
            "const Qux = 2\n"
        )
        kinds = {(k, n) for k, n, _ in decls}
        self.assertEqual(kinds, {("func", "Foo"), ("type", "Bar"),
                                 ("var", "Baz"), ("const", "Qux")})

    def test_method_receiver_excluded(self):
        decls = self._decls(
            "package p\n"
            "type T struct{}\n"
            "func (t T) Method() {}\n"
        )
        names = {n for _, n, _ in decls}
        self.assertIn("T", names)
        self.assertNotIn("Method", names)

    def test_const_block(self):
        decls = self._decls(
            "package p\n"
            "const (\n"
            "    A = 1\n"
            "    B int = 2\n"
            "    C = \"hi\"\n"
            ")\n"
        )
        names = {n for k, n, _ in decls if k == "const"}
        self.assertEqual(names, {"A", "B", "C"})

    def test_var_block(self):
        decls = self._decls(
            "package p\n"
            "var (\n"
            "    X = 1\n"
            "    Y string = \"z\"\n"
            ")\n"
        )
        names = {n for k, n, _ in decls if k == "var"}
        self.assertEqual(names, {"X", "Y"})

    def test_skips_string_literals(self):
        # `func Inside` lives inside a double-quoted string and a raw string —
        # neither should be reported as a top-level decl.
        decls = self._decls(
            "package p\n"
            'var s = "func Inside() {}"\n'
            "var r = `func AlsoInside() {}`\n"
            "func Real() {}\n"
        )
        funcs = {n for k, n, _ in decls if k == "func"}
        self.assertEqual(funcs, {"Real"})

    def test_skips_line_and_block_comments(self):
        decls = self._decls(
            "package p\n"
            "// func Commented() {}\n"
            "/* func AlsoCommented() {} */\n"
            "/*\n"
            "func MultiLineCommented() {}\n"
            "*/\n"
            "func Real() {}\n"
        )
        funcs = {n for k, n, _ in decls if k == "func"}
        self.assertEqual(funcs, {"Real"})

    # --- Regression-lock for RCA Pattern 10 (brace-unaware decl scanner) ---
    # See `.lovable/memory/07-test-failure-rca-patterns.md` §"Pattern 10".
    # Function-local `var got Variant` declarations inside `t.Run` bodies
    # MUST be ignored — they are not top-level decls.

    def test_p10_function_local_var_inside_trun_is_ignored(self):
        decls = self._decls(
            "package p\n"
            "import \"testing\"\n"
            "func TestX(t *testing.T) {\n"
            "    t.Run(\"a\", func(t *testing.T) {\n"
            "        var got Variant\n"
            "        _ = got\n"
            "    })\n"
            "    t.Run(\"b\", func(t *testing.T) {\n"
            "        var got Variant\n"
            "        _ = got\n"
            "    })\n"
            "}\n"
            "var TopLevel = 1\n"
        )
        var_names = [n for k, n, _ in decls if k == "var"]
        self.assertNotIn("got", var_names)
        self.assertIn("TopLevel", var_names)

    def test_p10_braces_inside_strings_do_not_shift_depth(self):
        decls = self._decls(
            "package p\n"
            "var Decoy = \"{ this is not a real brace { { {\"\n"
            "var Real = 1\n"
        )
        var_names = {n for k, n, _ in decls if k == "var"}
        self.assertIn("Decoy", var_names)
        self.assertIn("Real", var_names)


class TestMainEndToEnd(unittest.TestCase):
    def _run(self, root: Path) -> tuple[int, str, str]:
        out, err = io.StringIO(), io.StringIO()
        with redirect_stdout(out), redirect_stderr(err):
            rc = cc.main(["check-collisions.py", str(root)])
        return rc, out.getvalue(), err.getvalue()

    def test_no_go_files_returns_2(self):
        with tempfile.TemporaryDirectory() as td:
            rc, _, err = self._run(Path(td))
            self.assertEqual(rc, 2)
            self.assertIn("No .go files matched", err)

    def test_clean_repo_returns_0(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            write(root, "pkg/a.go", "package pkg\nfunc Alpha() {}\n")
            write(root, "pkg/b.go", "package pkg\nfunc Beta() {}\n")
            rc, out, _ = self._run(root)
            self.assertEqual(rc, 0, msg=out)
            self.assertIn("No collisions", out)

    def test_cross_file_collision_returns_1(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            write(root, "pkg/a.go", "package pkg\nfunc Dup() {}\n")
            write(root, "pkg/b.go", "package pkg\nfunc Dup() {}\n")
            rc, out, _ = self._run(root)
            self.assertEqual(rc, 1)
            self.assertIn("CROSS-FILE EXACT COLLISIONS", out)
            self.assertIn("Dup", out)

    def test_build_tag_siblings_are_not_collisions(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            write(root, "pkg/Get_linux.go", "package pkg\nfunc Get() {}\n")
            write(root, "pkg/Get_darwin.go", "package pkg\nfunc Get() {}\n")
            write(root, "pkg/Get_windows.go", "package pkg\nfunc Get() {}\n")
            rc, out, _ = self._run(root)
            self.assertEqual(rc, 0, msg=out)

    def test_accessor_pair_is_not_a_case_collision(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            write(root, "pkg/a.go", "package pkg\nvar Foo = foo\n")
            write(root, "pkg/b.go", "package pkg\nvar foo = 1\n")
            rc, out, _ = self._run(root)
            self.assertEqual(rc, 0, msg=out)

    def test_case_insensitive_collision_returns_1(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            # Two distinct exported names differing only in case → not the
            # legit accessor pattern.
            write(root, "pkg/a.go", "package pkg\nfunc Handler() {}\n")
            write(root, "pkg/b.go", "package pkg\nfunc HANDLER() {}\n")
            rc, out, _ = self._run(root)
            self.assertEqual(rc, 1)
            self.assertIn("CASE-INSENSITIVE COLLISIONS", out)

    def test_intra_file_duplicate_returns_1(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            write(root, "pkg/a.go",
                  "package pkg\nfunc Dup() {}\nfunc Dup() {}\n")
            rc, out, _ = self._run(root)
            self.assertEqual(rc, 1)
            self.assertIn("INTRA-FILE DUPLICATES", out)

    def test_collisions_scoped_per_package(self):
        # Same name in different package directories is allowed.
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            write(root, "pkga/a.go", "package pkga\nfunc Shared() {}\n")
            write(root, "pkgb/b.go", "package pkgb\nfunc Shared() {}\n")
            rc, out, _ = self._run(root)
            self.assertEqual(rc, 0, msg=out)


if __name__ == "__main__":
    unittest.main()
