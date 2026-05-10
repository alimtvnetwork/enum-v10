#!/usr/bin/env python3
"""
check-stringer-recursion.py — static guard for RCA Pattern P9.

Pattern P9 (see .lovable/memory/07-test-failure-rca-patterns.md):
  A `func (it T) String() string` body that calls `converters.AnyTo.ValueString`
  (which falls through to `fmt.Sprintf("%v", ...)`) is an infinite-recursion
  bomb — `%v` re-invokes the type's own `String()` method.

This script scans Go source (excluding `*_test.go`) and fails non-zero if any
`String() string` method body contains a call to `converters.AnyTo.ValueString`.
It uses brace-depth tracking (cf. RCA Pattern P10) so nested blocks don't fool
the scanner.

Usage:
    check-stringer-recursion.py <root> [<root> ...]

Exit codes:
    0  no violations
    1  one or more violations found
    2  invalid arguments / no Go files scanned
"""
from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

# `func (recv T) String() string {`  — also accepts pointer receivers and
# whitespace variations. We deliberately do NOT match generic type params
# because `String() string` is never generic.
STRING_METHOD_RE = re.compile(
    r"^\s*func\s*\(\s*\w+\s+\*?\w[\w.]*\s*\)\s*String\s*\(\s*\)\s*string\s*\{",
)
FORBIDDEN_CALL = "converters.AnyTo.ValueString"


def strip_strings_and_comments(line: str) -> str:
    """Cheap pass — drop // line comments and string contents so braces inside
    them don't shift depth. Sufficient for brace-depth bookkeeping; not a full
    Go tokenizer."""
    out = []
    i = 0
    n = len(line)
    while i < n:
        c = line[i]
        if c == "/" and i + 1 < n and line[i + 1] == "/":
            break
        if c == '"' or c == "`":
            quote = c
            i += 1
            while i < n and line[i] != quote:
                if quote == '"' and line[i] == "\\" and i + 1 < n:
                    i += 2
                    continue
                i += 1
            i += 1
            continue
        out.append(c)
        i += 1
    return "".join(out)


def scan_file(path: Path) -> list[tuple[int, str]]:
    """Return [(lineno, snippet), ...] for each violation inside a String() body."""
    try:
        lines = path.read_text(errors="replace").splitlines()
    except OSError:
        return []

    violations: list[tuple[int, str]] = []
    in_string_method = False
    method_depth = 0
    brace_depth = 0

    for idx, raw in enumerate(lines, 1):
        clean = strip_strings_and_comments(raw)

        if not in_string_method and STRING_METHOD_RE.match(raw):
            in_string_method = True
            method_depth = brace_depth + 1  # the `{` on this line opens the body
            opens = clean.count("{")
            closes = clean.count("}")
            brace_depth += opens - closes
            continue

        if in_string_method and FORBIDDEN_CALL in raw:
            violations.append((idx, raw.strip()))

        opens = clean.count("{")
        closes = clean.count("}")
        brace_depth += opens - closes

        if in_string_method and brace_depth < method_depth:
            in_string_method = False
            method_depth = 0

    return violations


def iter_go_files(root: Path):
    for p in root.rglob("*.go"):
        if p.name.endswith("_test.go"):
            continue
        # Skip vendored / cross-repo upstream mirrors.
        parts = set(p.parts)
        if "vendor" in parts or "cross-repo" in parts:
            continue
        yield p


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("roots", nargs="+", type=Path)
    args = parser.parse_args(argv)

    scanned = 0
    total_violations: list[tuple[Path, int, str]] = []

    for root in args.roots:
        if not root.exists():
            print(f"::error::{root}: does not exist", file=sys.stderr)
            return 2
        for f in iter_go_files(root):
            scanned += 1
            for ln, snippet in scan_file(f):
                total_violations.append((f, ln, snippet))

    if scanned == 0:
        print("::error::no Go files scanned", file=sys.stderr)
        return 2

    print(f"check-stringer-recursion: scanned {scanned} Go source files.")

    if total_violations:
        print(f"\n::error::Pattern P9 violation — Stringer recursion bomb detected:")
        print("  A String() string method body must NOT call converters.AnyTo.ValueString.")
        print("  See .lovable/memory/07-test-failure-rca-patterns.md (Pattern P9).\n")
        for f, ln, snippet in total_violations:
            print(f"  {f}:{ln}: {snippet}")
        return 1

    print("✅ No Pattern P9 violations.")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
