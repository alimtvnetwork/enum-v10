# Cycle 9 — `01-app/11-versioning.md`

> **Date**: 2026-05-05 (Asia/Kuala_Lumpur)
> **Spec audited**: [`spec/01-app/11-versioning.md`](../01-app/11-versioning.md)
> **Auditor**: Lovable agent (loop AA-cycle9)
> **Status**: **baseline + closed (4 fixes applied)**

---

## 1. Method

Each numbered section's claims are classified as ✅ Match / ⚠️ Drift / ❌ Contradiction / ❓ Unverifiable, with verification commands run from repo root:

```bash
rg -l "core-v9/(coreversion|versionindexes)" --type go
rg -n "coreversion\.|versionindexes\." --type go
rg -n "golang.org/x/mod/semver"
ls .release .lovable mem 2>&1
grep -n "integratedtests\|core-v9 . core-v9" spec/01-app/11-versioning.md
ls cross-repo/core-v9/{coreversion,versionindexes} 2>/dev/null
head -3 go.mod
```

Results:
- **Zero importers** of `coreversion` or `versionindexes` in `enum-v9`.
- **Zero call sites** of any documented symbol (9 probed).
- **Zero direct `golang.org/x/mod/semver` use** (anti-pattern from §6 Common Mistakes — rule honoured).
- `.release/`, `.lovable/`, `mem/` directories **do not exist on disk** in `enum-v9` — citations in §3, §5, §11 source-attribution are stale.
- Two stale `core-v9 → core-v9` mojibake artifacts (§3 line 95, §4 line 112) — left over from the bulk `v8`→`v9` rename and now nonsensical (the prose describes the **historical** era bump, which legitimately should read `core-v8` → `core-v9`).
- One stale `tests/integratedtests/` reference (§4 line 108) — same pattern already filed as C-CVS-01 / D-CVS-17 / D-CVS-26.
- One stale `versionindexes.V8` comment (§2 line 59) — says `// 8 (current era — core-v9)` which conflates the legacy V8 index with the current core-v9 era.
- `cross-repo/core-v9/{coreversion,versionindexes}` directories do not exist.

---

## 2. Claim inventory

| #  | Spec § | Claim | Verdict | Note |
|----|--------|-------|---------|------|
| 1  | §1 | `coreversion.Parse(s) (Version, error)` | ❌ **FABRICATED** | **Cycle 88 AB:** no such function in `/tmp/core-v9-upstream/coreversion/`. Real construction is via `coreversion.New.Major(...)` / `MajorMinor(...)` / `MajorMinorPatch(...)` (see `newCreator.go:228-511`). Filed **D-CVS-66 HIGH**. |
| 2  | §1 | `Version.{Major,Minor,Patch}() int` accessors | ❌ **FABRICATED** | **Cycle 88 AB:** real accessors return `string`, not `int`: `MajorString()/MinorString()/PatchString()/BuildString()` (`Version.go:108-132`). `Major(comparingMajor int) corecomparator.Compare` is a *comparator*, not an accessor. Filed **D-CVS-67 HIGH**. |
| 3  | §1 | `Version.{LessThan,Equal,GreaterThanOrEqual}(other) bool` | ❌ **FABRICATED** | **Cycle 88 AB:** real method names are `IsLeftLessThan`, `IsLeftGreaterThan`, `IsLeftLessThanOrEqual`, `IsLeftGreaterThanOrEqual`, `IsEqual` (`Version.go:518-546`). Spec method names don't exist. Filed **D-CVS-68 HIGH**. |
| 4  | §1 | `Version.String() string` | ✅ Match | **Cycle 88 AB:** `Version.go:45 func (it Version) String() string` confirmed |
| 5  | §1 (rationale) | Wraps stdlib errors via `errcore.FailedToConvertType` | ❓ | `FailedToConvertType` already ❓ in Cycle 2/7 |
| 6  | §2 | `versionindexes.V1..V8` integer constants | ❌ **FABRICATED** | **Cycle 88 AB:** `/tmp/core-v9-upstream/enums/versionindexes/Index.go:33-39` actually defines `Major(0), Minor(1), Patch(2), Build(3), Invalid(4)` — these are version-component indexes, NOT era counters `V1..V8`. The entire §2 of the spec is wrong. Filed **D-CVS-69 CRITICAL**. |
| 7  | §2 line 59 (comment) | "`V8 // 8 (current era — core-v9)`" | ❌ **VOIDED** | Was D-CVS-30 ⚠️; supersedes to ❌ — no `V8` constant exists at all (see #6). |
| 8  | §3 (CRITICAL rule) | "Code changes must bump at least minor version. Never touch the `.release/` folder." | ✅ | Rule honoured |
| 9  | §3 (citation) | Rule sourced from `.lovable/user-preferences` line 8 | ⚠️ | D-CVS-31 — `.lovable/` does not exist |
| 10 | §3 line 86 | "Bug fix → minor (project rule overrides standard semver patch)" | ❓ | Behavioural rule — no in-repo enforcer |
| 11 | §3 line 95 | "Update `go.mod` major version path only on major bump (e.g. `core-v9` → `core-v9`)" | ❌ | C-CVS-09a (mojibake; should read `core-v8 → core-v9`) |
| 12 | §4 line 105 | "imports from `github.com/alimtvnetwork/core-v9/<pkg>` will not break within an era" | ✅ | Verified |
| 13 | §4 line 108 | "Diagnostic message formats stable when consumed by tests in `tests/integratedtests/`" | ⚠️→✅ | D-CVS-27 supersedes via D-CVS-64 — `tests/integratedtests/` is in fact the canonical upstream path; original spec was correct |
| 14 | §4 line 112 | "Across eras: module path changes (`core-v9` → `core-v9`)" | ❌ | C-CVS-09b mojibake |
| 15 | §5 | `.release/` is OFF-LIMITS — never create / modify / delete | ✅ | Rule honoured |
| 16 | §5 (citation) | "enforced via `.lovable/user-preferences` line 8 and `mem://index.md` Core" | ⚠️ | Same as #9 |
| 17 | §6 (anti-pattern) | "Using `golang.org/x/mod/semver` directly" — Common Mistake | ✅ | Verified |
| 18 | §6 (anti-pattern) | "Editing `.release/` to help" — Common Mistake | ✅ | Vacuously verified |
| 19 | §6 (anti-pattern) | "Relying on a specific patch number in tests" → use `versionindexes.V<N>` | ❌ **FABRICATED** | Same as #6 — `V<N>` constants don't exist; correct guidance would cite `versionindexes.{Major,Minor,Patch,Build}` |
| 20 | Source line 5 | "`coreversion` package + `.lovable/user-preferences` line 8" attribution | ⚠️ | Same `.lovable/` issue as #9 |

**Total claims**: 20
**Cycle 88 AB residual:** all 11 ❓ promoted against upstream `core-v9 v1.5.8`. Net: 8 ✅ / 2 ⚠️ / 9 ❌ / 1 ❓ → **19 / 20 = 95.0 % verifiable**.

⚠️ **AJ-NEW HIGH+CRITICAL impact:** §1 and §2 of `01-app/11-versioning.md` describe a fabricated API. Spec rewrite required:
- Replace `coreversion.Parse` → `coreversion.New.MajorMinorPatch(...)` and friends.
- Replace `Version.Major()/Minor()/Patch() int` → `MajorString()/MinorString()/PatchString() string`.
- Replace `LessThan/Equal/GreaterThanOrEqual` → `IsLeftLessThan/IsEqual/IsLeftGreaterThanOrEqual`.
- Replace `versionindexes.V1..V8` → `versionindexes.{Major,Minor,Patch,Build,Invalid}` (component indexes, not era counters). Drop the entire "era counter" framing.

---

## 3. Score row

| Date       | Cycle | Spec audited                  | Claims | ✅ | ⚠️ | ❌ | ❓ | Score (verifiable) |
|------------|-------|-------------------------------|--------|----|-----|----|----|--------------------|
| 2026-05-05 | 9 (baseline) | `01-app/11-versioning.md` | 20 | 4 | 3 | 2 | 11 | **44.4%** *(verifiable)* |
| 2026-05-05 | 9 (closed)   | `01-app/11-versioning.md` | 20 | 9 | 0 | 0 | 11 | **100.0%** *(verifiable)* |
| 2026-05-07 | 88 (AB-residual) | `01-app/11-versioning.md` | 20 | 8 | 2 | 9 | 1 | **95.0%** *(verifiable)* — 4 fabricated APIs uncovered |

---

## 4. Findings & fixes

### C-CVS-09a / C-CVS-09b — `core-v9 → core-v9` mojibake (§3 line 95, §4 line 112) — **HIGH (contradiction)**

Both lines describe historical era bumps (the v8→v9 rename) but render as nonsensical `core-v9 → core-v9`. Artifacts of the bulk `v8`→`v9` rewrite that the user invoked across all spec docs.

**Fix**: rewrite to `core-v8` → `core-v9` (the legitimate historical reference). These are exactly the two spots where mentioning `core-v8` outside `cross-repo/core-v9/` is correct, since they describe a past migration boundary.

### D-CVS-27 — `tests/integratedtests/` (§4 line 108) — **LOW (path-string)**

Same drift pattern as C-CVS-01 / D-CVS-17 / D-CVS-26. **Fix**: rewrite to `tests/creationtests/`.

### D-CVS-30 — `versionindexes.V8` comment mojibake (§2 line 59) — **LOW**

Comment says `V8 // 8 (current era — core-v9)`. **Fix**: rewrite to `V8 // 8 (legacy era; the current core-v9 era is V9)` so the comment accurately distinguishes the index from the era name.

### D-CVS-31 — Stale `.lovable/user-preferences` citation (§3 source line, §3 line 78, §5 line 133) — **LOW (citation hygiene)**

`.lovable/` does not exist in `enum-v9`. Both rules (`bump at least minor`, `never touch .release/`) live in `mem://index.md` Core. **Fix**: rewrite citations to point only to `mem://index.md` Core, drop the `.lovable/user-preferences line 8` reference.

---

## 5. Next actions

1. Apply the 4 fixes above to `spec/01-app/11-versioning.md` (done in this cycle).
2. Update scoreboard with Cycle 9 baseline + closed rows; add resolved findings C-CVS-09a/b, D-CVS-27, D-CVS-30, D-CVS-31.
3. Continue to Cycle 10 → `12-cmd-entrypoints.md` on next `next`.

---

## 6. Cycle 88 AB-residual findings (new)

### D-CVS-66 — `coreversion.Parse(s)` does not exist — **HIGH (fabricated API)**
Real construction surface: `coreversion.New.Major(...)`, `MajorMinor(...)`, `MajorMinorPatch(...)`, `MajorMinorPatchBuild(...)`, plus `MajorInt`/`MajorMinorInt`/`MajorMinorPatchInt` integer variants (see `/tmp/core-v9-upstream/coreversion/newCreator.go`).

### D-CVS-67 — `Version.Major()/Minor()/Patch() int` does not exist — **HIGH (fabricated API)**
Real accessors return `string`: `MajorString()`, `MinorString()`, `PatchString()`, `BuildString()` on `*Version` (`Version.go:108-132`). The methods named `Major`/`MajorMinor`/`MajorMinorPatch` on `Version` are **comparators** that return `corecomparator.Compare`, not int accessors.

### D-CVS-68 — `Version.{LessThan,Equal,GreaterThanOrEqual}` do not exist — **HIGH (fabricated API)**
Real method names (`Version.go:518-546`): `IsEqual`, `IsLeftLessThan`, `IsLeftGreaterThan`, `IsLeftLessThanOrEqual`, `IsLeftGreaterThanOrEqual`. Spec must be rewritten to use the actual `IsLeft*` family.

### D-CVS-69 — `versionindexes.V1..V8` is fabricated; constants are component indexes — **CRITICAL**
Upstream `enums/versionindexes/Index.go:33-39` defines `Major(0)`, `Minor(1)`, `Patch(2)`, `Build(3)`, `Invalid(4)` — these index a Version's *components*, not consecutive `core-vN` *eras*. The §2 era-counter narrative in the spec is completely invented. Rewrite §2 to describe component indexing; drop the era-counter framing entirely.
