# 07 — Coverage Gate ≥ 60%

## Symptom

`.github/workflows/ci.yml` fails when total Go test coverage drops below 60%.

## Root Cause

Intentional gate. Coverage threshold is enforced at the workflow level after `go test -cover`.

## Fix / Workaround

Working as designed. To raise coverage:
- Add tests in `tests/creationtests/` (NOT `tests/integratedtests/` — that path doesn't exist in `enum-v9`).
- Use the Goconvey + `EnumTestWrapper` registry pattern documented in `spec/01-app/13-testing-patterns.md` §6.1.

## Status

✅ Working as designed.

## Related

- `mem://features/ci-tooling.md` § "Coverage gate"
- `.github/workflows/ci.yml`
- `.lovable/memory/04-test-layout.md`
