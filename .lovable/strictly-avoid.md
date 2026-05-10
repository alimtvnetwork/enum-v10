# Strictly Avoid — Hard Prohibitions

This is the project's "never do this" list. Every entry is enforced across all sessions.

## Communication & delivery

- **Email-based notifications of any kind.** Never send emails, configure Dependabot recipients, add SMTP code, propose CI email alerts, or suggest "we'll notify you by email." User explicitly rejects all email flows. See: `mem://constraints/no-email-notifications.md`.

## Module paths & renames

- **Reintroducing `github.com/alimtvnetwork/core-v8` outside `cross-repo/core-v9/`.** All source imports use `core-v9`. The mirror directory keeps `core-v8` intentionally and tracks a different upstream module.
- **Rewriting `enum-v1` → `enum-v10` inside `cross-repo/core-v9/`.** That directory tracks a different module and must not be touched by the bulk-rename.
- **Pseudo-version `v1.5.6-0.<date>-<sha>` for the `replace` bridge.** It needs a `v1.5.5` predecessor tag that doesn't exist; Go module proxy will reject it.

## Spec content

- **Writing `tests/integratedtests/` in new spec content.** Tests live under `tests/creationtests/`. Audit cycles 1, 3, 6, 8, 9, 10, 11, 12 (D-CVS-17/26/27/32/33/36/37/39/40/41) all corrected this pattern. Anti-pattern callouts that explicitly mark the string as wrong (e.g. `05-enum-system.md:417`) are the only exception.
- **Mojibake `core-v9 → core-v9`** (broken historical references). The historical migration was `core-v8` → `core-v9`. See cycle 9 fixes C-CVS-09a/b.
- **Citing `.lovable/user-preferences line 8`** as a source of authority — the file did not exist when cycle 9 ran (D-CVS-31). The current loop is creating `.lovable/` for the first time; even now, prefer citing `mem://index.md` Core for cross-cutting rules.

## Tooling probes

- **Hard-coding `tests/creationtests/` OR `tests/integratedtests/` in test-discovery tooling.** Tooling must accept either name or read from disk — never branch on a single hard-coded value.

## Process

- **Skipping the remaining-task list at the end of a "Next" reply.** Breaks the user's mental tracking. See: `mem://preferences/workflow.md`.
- **Batching multiple lettered tasks into one "Next" turn** unless the user explicitly asks for it.
- **Renaming `cross-repo/core-v9/`** — it intentionally tracks a different upstream repo.
