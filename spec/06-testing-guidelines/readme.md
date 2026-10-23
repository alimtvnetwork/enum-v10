# Integrated Testing Guidelines for Go

> **Portable guideline** — drop this folder into any Go project that uses the `coretests` framework.
> It covers folder structure, test case types, input holders, result assertions, branch coverage strategy, and real-world examples.

> **Consumer-coverage note (`enum-v10`)**: this folder is the **portable upstream `core-v9`** testing guideline. Every reference below to `tests/integratedtests/<pkg>tests/`, `coretests.GetAssert`, `args.Map`, `CaseV1`, `coretestcases`, the per-package directory layout, and the internal-package coverage policy describes **upstream `core-v9`** conventions. `enum-v10` does **not** consume any of them — `rg tests/integratedtests` and `rg coretests\\.` over `enum-v10` source both return zero hits. This module's tests live at `tests/creationtests/` (one shared package, Goconvey-based registry over `EnumTestWrapper`); see [`spec/01-app/13-testing-patterns.md` §6.1](../01-app/13-testing-patterns.md#61-enum-v10-specific-layout) and [`spec/01-app/14-tests-folder-walkthrough.md`](../01-app/14-tests-folder-walkthrough.md) for the `enum-v10`-specific layout. Treat the files in this folder as the authoritative reference for **upstream** projects only.

## Table of Contents

| File | Purpose |
|------|---------|
| [01-folder-structure.md](01-folder-structure.md) | Directory layout, naming conventions, file separation rules |
| [02-test-case-types.md](02-test-case-types.md) | `CaseV1`, `CaseNilSafe`, `GenericGherkins` — when to use each |
| [03-args-reference.md](03-args-reference.md) | `args.Map`, `args.One`–`args.Six`, `args.Dynamic`, `args.Holder`, `args.LeftRight` |
| [04-results-reference.md](04-results-reference.md) | `results.Result`, `results.ResultAny`, `results.ExpectAnyError`, `InvokeWithPanicRecovery` |
| [05-assertion-patterns.md](05-assertion-patterns.md) | `ShouldBeEqual`, `ShouldBeEqualMap`, `ShouldBeSafe`, diff-based assertions |
| [06-branch-coverage.md](06-branch-coverage.md) | Positive/negative paths, if/else, nil guards, boundary conditions |
| [08-good-vs-bad.md](08-good-vs-bad.md) | Concrete examples of good tests vs. bad tests |
| [09-creating-custom-cases.md](09-creating-custom-cases.md) | How to create your own case type using `BaseTestCase` |

## Core Principles

1. **Separation of data and logic** — test cases (`_testcases.go`) never contain assertions; test runners (`_test.go`) never contain expected values
2. **AAA comments are mandatory** — every test function has `// Arrange`, `// Act`, `// Assert` comments
3. **No raw `t.Error` / `t.Errorf`** — always use framework assertions (`ShouldBeEqual`, `ShouldBeEqualMap`, etc.)
4. **Native types in expectations** — use `bool`, `int`, `string` in `args.Map`, not `"true"`, `"5"`
5. **One test function per logical scenario** — no branching logic inside test bodies
6. **No coverage-motivated tests for `internal/` packages** — MUST NOT write `Coverage*_test.go` files or include `internal/` packages in coverage iteration plans. Business/integration tests for internal packages are allowed under `tests/integratedtests/<pkg>tests/` (e.g. `csvinternaltests/`, `fsinternaltests/`). See [06-branch-coverage.md § Internal Package Coverage Policy](06-branch-coverage.md#internal-package-coverage-policy).
