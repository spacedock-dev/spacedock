---
id: 8e2053706c2c77xdm696dfr6
title: dispatch build emits Agent name >64 chars for long slugs — Agent() dispatch fails
status: ideation
source: github#366 (captain intake 2026-06-14)
started: 2026-06-14T05:20:44Z
completed:
verdict:
score: "0.40"
worktree:
issue: "#366"
---

`spacedock dispatch build` constructs the worker `name` as `{worker_key}-{slug}-{stage}` with no length cap, but Claude Code's `Agent` tool enforces `name` matching `^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$` (max 64 chars). For a long slug the emitted `name` exceeds 64 chars, so forwarding the helper output verbatim to `Agent()` — the contract's MANDATORY dispatch path — fails with `InputValidationError` on `name`, and the dispatch cannot proceed.

## Problem

- Name construction (`internal/dispatch/build.go`, `derivedName := fmt.Sprintf("%s-%s-%s", workerKey, slug, stage)`) has no 64-char cap.
- Observed live: `pr-merge-mod-base-branch-post-flip` → name 60 chars, dispatched fine; `dispatch-reconcile-deconflate-repo-hygiene` → name 68 chars, `Agent(name=...)` rejected with `InputValidationError … (max 64 chars)`.
- **Decomposition is load-bearing.** A hand-shortened or naively-truncated `name` no longer decomposes cleanly to `(worker_key, slug, stage)`. That decomposition drives supersede-shutdown, terminal-team-teardown cohort derivation, and reconcile A/B (lingering/superseded) classification. A name that does not map back to its entity risks a false Class-A "lingering" flag (a spurious shutdown against a live worker) or a missed cohort teardown.

## Notes (suggested-fix direction from github#366 — ideation owns the actual design)

- Cap the generated `name` at 64 chars deterministically inside `dispatch build`, preserving a stable, decomposable identity — e.g. truncate the slug component to fit while keeping the `worker_key` prefix + `stage` suffix, or append a short stable hash of the full slug so distinct long slugs do not collide.
- Whatever the form, the reconcile/teardown name-decomposition must round-trip back to the entity (e.g. resolve via a stored short-id rather than string-splitting the full slug).
- `name` and `dispatch_file_path` need not match — only `name` needs capping.

## Spike: riskiest mechanism exercised first

The cycle-1 design rested on a re-invented `sha256(slug)` hash, which the captain rejected: the workflow already mints a stored, SHA-derived, RESOLVABLE per-entity `id` (sd-b32, e.g. this task's `8e2053706c2c77xdm696dfr6`). Re-hashing the slug reinvents that token. The redesign uses the existing `id` as the cap token, so the riskiest mechanism SHIFTED from the slug-hash round-trip to the **id-prefix** round-trip. Per the captain's mandate, that new mechanism was re-exercised with a throwaway spike against the real `dispatch build`, `status --resolve`, and `decompose`/reconcile code over REAL sd-b32 entities minted by `status --new`. Results:

- **Bug reproduced over real `dispatch build`.** A 57-char slug at `backlog` produced an 82-char `name` (`spacedock-ensign-{slug}-backlog`) — past the Agent-tool 64-char ceiling.
- **`status --resolve` resolves a FIXED id-prefix.** Two real entities (`367s2zrbkm4fcwfff5ac62zz`, `0qv3kaqs062fp3p01tz6re99`) minted by `status --new`. `status --resolve prefix:367s2zrb` and `prefix:0qv3kaqs` (fixed 8-char prefixes) each resolved to the correct entity via `strings.HasPrefix(storedID, prefix)` (`internal/status/resolve.go:108-115`, guarded by `idStyle == "sd-b32"` and `len ≥ sdB32MinPrefix=2`). The bare ref (no `prefix:`) also resolved (default mode tries prefix when alphabet-valid). The shortest-unique *display* prefix (`computeSDB32DisplayIDs`, 2 chars here) was deliberately NOT used — it lengthens as entities are added, so it is unstable for a name baked at dispatch time. A fixed 8–10 char prefix is stable.
- **id-form name fits with NO truncation.** `spacedock-ensign-367s2zrb-implementation` = 40 chars (≤64) at the tightest stage; `+ -cycle3` = 47. Across all canonical stages with the full `spacedock-ensign` prefix and a 10-char id-prefix, the longest id-form name is 42 chars — the fixed-length id-form never overflows, so for `sd-b32` workflows no slug truncation and no `worker_key` trim is needed for the budget. Both names passed `namePattern`; distinct ids → distinct names.
- **Capped name still decomposes its stage/cycle.** `decompose()` peeled `implementation` from each id-form name unchanged; the embedded token is the id-prefix (`367s2zrb`), all-lowercase-alnum from the sd-b32 alphabet (`0123456789abcdefghjkmnpqrstvwxyz` — `namePattern`-safe, no hyphens).
- **Round-trip resolves to the correct entity, no false Class-A.** Alpha was archived (terminal), bravo left active. Resolving each decomposed id-prefix token by `HasPrefix` against the real active+archived id→slug maps recovered the correct slug for both; alpha → archived slug (Class-A teardown, correct), bravo → active slug (NO false Class-A against the live sibling). This is the exact hazard github#366 names.

The spike (a throwaway Go test plus a real sd-b32 fixture) was removed before commit; what it teaches seeds the implementation's first test (see Test plan).

## Proposed approach

The cap fires only on overflow; short names stay byte-identical (preserving today's readable names and all golden fixtures). The token used to cap is **id-first**, falling back to slug-truncation only where no id exists.

**1. Path by id-style.** `dispatch build` reads the entity's `id` from frontmatter (already parsed at `build.go:378` via `ParseFrontmatter`) and the workflow `id-style` from the README.
   - **`id-style: sd-b32` (id-first).** When `{workerKey}-{slug}-{stage}` overflows the budget, replace the slug component with a **fixed-length id-prefix** (a leading slice of the entity's 24-char sd-b32 `id`). The id-prefix length is a fixed constant (proposed 10; 8 proven sufficient in the spike, 10 leaves more collision headroom and still fits — a 10-char id-form is ≤42 chars). The token is `namePattern`-safe by construction (sd-b32 alphabet ⊂ lowercase-alnum). `workerKey` prefix and `-{stage}` suffix are preserved verbatim, so `decompose()`'s prefix-strip and stage-peel are undisturbed.
   - **`id-style: sequential`.** The `id` is a stored numeric string (not a stable random token); treat as id-first using the full sequential id (short — never overflows) when one is present. (Sequential ids are short integers, so the id-form trivially fits; no truncation.)
   - **`id-style: slug` (no id).** There is no stored id (`displayID = slug`). Fall back to truncating the slug head to fit the budget, trimming any trailing hyphen so the result matches `namePattern` and carries no `--`. This is the ONLY path that truncates a slug. Distinct id-less slugs that share a long prefix CAN collide under pure truncation; that collision risk is inherent to id-less workflows and is recorded in Out of scope (id-less workflows do not get the stable-token guarantee — they never had a stable token).

**2. Cap target reserves cycle headroom.** The FO appends `-cycleN` to the base name downstream of `dispatch build` (supersede / fresh-redispatch roster convention; cycle 3 is the feedback-flow escalation ceiling). A base name at exactly 64 overflows once `-cycle2` is appended (71). So the cap **target** is `64 − len("-cycle9") = 57`, leaving room for any single-digit cycle suffix; `{base}-cycle3` ≤ 64. For the id-first path this headroom is automatically satisfied (id-form ≤ 42); it binds only the slug-truncation fallback. (`dispatch build` itself never emits a cycle — it has no cycle field — so the headroom is reserved, not consumed.)

**3. Resolve capped names back to the entity in reconcile.** Make reconcile's name→entity resolution id-aware.
   - **Carry the id.** `loadEntityFrontmatter` (`reconcile.go:266-305`) already reads `slug/status/worktree/pr`; add `id` to `entityRecord` so the sweep has an id→record map for active and archived entities (it already loads both — `reconcile.go:201-202`).
   - **Resolve the decomposed token id-first.** For each member name, `decompose()` still peels the worker prefix, cycle, and stage; the remaining token is matched **exact slug first** (uncapped names and the common case), then, when no slug matches, by **sd-b32 id-prefix** (`strings.HasPrefix(record.id, token)`) against the active+archived id set. A unique id-prefix match yields the entity's real slug for the Class-A entity lookup; the id-prefix token itself is the Class-B cohort key (consistent per entity, distinct across entities). An ambiguous prefix (≥2 matches) or no match leaves the member unclassified (`ok=false`) — the conservative behavior that avoids a false Class-A. No recompute-the-capped-form-per-known-slug dance, no new on-disk state — the id is already in frontmatter and resolution is `HasPrefix` against a map the sweep already builds. (`status --resolve prefix:` is the proven CLI form of this same `HasPrefix` resolution; reconcile resolves in-process against its loaded maps rather than shelling out.)

**4. worker_key trim — evaluated, recommend OUT.** Trimming `worker_key` from `spacedock-ensign` (17 chars with hyphen) to `ensign` (7) reclaims ~10 of the 64 chars in the name. Considered because it raises the threshold below which a name stays human-readable. **Recommendation: OUT of this task.** Rationale and blast radius:
   - For `sd-b32` workflows the id-form already fits at ≤42 chars with the full prefix, so the trim buys nothing for the actual bug; it is pure readability polish.
   - `worker_key` is not just the name prefix — it also forms worktree paths (`.worktrees/{worker_key}-{slug}`) and branch names (`{worker_key}/{slug}`) (`claude-first-officer-runtime.md:57`). Trimming it is a cross-cutting rename touching the FO runtime, build, and `decompose()`'s hard-coded `const workerKey = "spacedock-ensign-"` (`reconcile.go:367`).
   - It is a new-dispatch-only cutover: `decompose()` would have to accept BOTH `spacedock-ensign-*` and `ensign-*` while the Commander's in-flight `spacedock-ensign-*` worktrees/branches drain, then drop the old form later — a migration window with dual-form decomposition, not a clean swap.
   - That blast radius is disproportionate to a readability gain the id-first cap does not need. Recommend filing it as a separate optional cleanup if the captain still wants shorter handles; this task does NOT assume it.

**5. FO break-glass template note.** The FO break-glass manual-dispatch template (`skills/first-officer/references/claude-first-officer-runtime.md:109`, `name="{worker_key}-{slug}-{stage}"`) is the one contract-prose touchpoint that also constructs a name and would also overflow for a long slug. It is a degraded fallback used only when `dispatch build` is unavailable and already omits much of the production helper's output. The doc change is a one-line caveat pointing at the helper as the cap authority (≤64; on `sd-b32` use a fixed id-prefix in place of the slug), so the fallback is not silently broken. This is the only doc diff; the `name` is an internal roster handle in no docs-site page or CLI output.

### Doc diff

`skills/first-officer/references/claude-first-officer-runtime.md`, break-glass template (around line 109):

- Before: `    name="{worker_key}-{slug}-{stage}",`
- After: `    name="{worker_key}-{slug}-{stage}",  // if this exceeds 64 chars, cap it the way `spacedock dispatch build` does: keep the {worker_key} prefix and -{stage} suffix and, on id-style: sd-b32, replace the slug with a fixed-length prefix of the entity id (id-less slug workflows truncate the slug head instead)`

(Exact wording finalized in implementation; the load-bearing content is the ≤64 cap instruction plus the id-prefix-in-place-of-slug recipe.)

## Out of scope

- **A new name↔slug registry / persisted mapping.** Resolution reuses the `id` already in frontmatter via `HasPrefix` against maps reconcile already loads — no new on-disk state. (NOTE: "reuse the existing frontmatter id" is NOT new state and IS in scope; only building a separate registry is YAGNI and out.)
- **A re-invented per-name hash.** The cycle-1 `sha256(slug)` token is dropped entirely in favor of the existing sd-b32 `id`.
- **`worker_key` trim (`spacedock-ensign` → `ensign`).** Evaluated in approach §4; recommended OUT (cross-cutting worktree/branch rename, dual-form decomposition migration window, no budget benefit for the id-first path). A candidate for a separate optional readability cleanup.
- **Stable-token guarantee for `id-style: slug` workflows.** Id-less workflows fall back to slug-head truncation, which can collide for shared-prefix slugs. They never had a stable resolvable token; this task does not invent one for them (that would be the rejected per-name hash). If an id-less workflow needs collision-free long handles, adopting `id-style: sd-b32` is the path.
- **Changing `dispatch_file_path` keying.** The dispatch filename has its own `dispatchFileNameMaxLen` (251) cap keyed on `team_name + derivedName`; per github#366 only `name` needs the 64-char cap. Left as-is.
- **Re-running `dispatch build` with a cycle parameter.** Reserving cycle headroom in the base cap (§2) is simpler than adding a cycle field + build round-trip per feedback cycle; the FO keeps its append-only cycle convention.

## Acceptance criteria

- **AC1 — `dispatch build` emits a `name` ≤64 chars for a long slug.** For a `sd-b32` entity whose `{workerKey}-{slug}-{stage}` would exceed 64, the real `dispatch build` output's `name` is ≤64 chars and matches `namePattern`.
  - Verified by: a Go test driving `dispatch build` in-process (`runNative`) over a `sd-b32` fixture entity with a long slug, asserting `len(name) <= 64` and `namePattern.MatchString(name)` on the parsed stdout JSON. Fails today (82 chars).
- **AC2 — distinct entities produce distinct names (no cohort collision).** Two distinct `sd-b32` entities whose slugs share a long common prefix produce two different `name` values from real `dispatch build` (distinct because their ids differ).
  - Verified by: a Go test building both fixtures through `dispatch build` and asserting the two emitted `name` values differ.
- **AC3 — a capped name resolves to the correct entity via its id.** Given the capped `name` `dispatch build` emits for an archived long-slug `sd-b32` entity, the reconcile sweep classifies it Class-A against the **correct** slug (resolved by id-prefix), and does NOT mis-resolve to a different active long-slug entity sharing the slug prefix.
  - Verified by: a Go test that emits the capped name via `dispatch build`, seeds an active dir + `_archive` dir with the two `sd-b32` entities (one archived, one active, shared slug prefix, distinct ids), runs real `Reconcile` over a stub roster carrying the capped name, and asserts exactly one Class-A entry pointing at the archived slug and none against the active one. The expectation comes from the seeded on-disk id set (independent source).
- **AC4 — capped names do not produce a false Class-A against a live entity.** A capped name whose `sd-b32` entity is still active produces NO Class-A drift entry.
  - Verified by: a Go test running real `Reconcile` over a roster with the capped name of an active long-slug entity and asserting zero Class-A entries for it (guards the false-shutdown hazard in the Problem section).
- **AC5 — short names are unchanged.** An entity whose `{workerKey}-{slug}-{stage}` already fits 64 emits exactly that name (no id substitution, no truncation), and existing build golden/parity fixtures pass unchanged.
  - Verified by: a Go test asserting the uncapped form for a short slug, plus the existing `go test ./internal/dispatch/` golden/parity suite passing unchanged.
- **AC6 — base name reserves cycle headroom.** A capped base name plus a realistic FO cycle suffix (`-cycle3`) stays ≤64.
  - Verified by: a Go test taking the capped `name` from `dispatch build` for a long-slug entity and asserting `len(name + "-cycle3") <= 64`.

Every AC's proof is a Go test over real `dispatch build` output and real `decompose`/`Reconcile` code — none is a substring/regex match over an instruction file. AC3/AC4 bind the expectation to the seeded on-disk id set, so a wrong id-prefix resolution fails the check.

## Test plan

- **Where:** `internal/dispatch/` (alongside `build_*_test.go` and `reconcile_*_test.go`), reusing the in-process harness (`runNative`, `writeFile`, `gitInit`) and the real `decompose`/`Reconcile`. Fixtures use `id-style: sd-b32` READMEs with entities carrying real sd-b32 `id` frontmatter (mint via the same path `status --new` uses, or write fixed valid 24-char sd-b32 ids directly into fixture frontmatter — the spike confirmed both resolve). No new test infrastructure.
- **AC1/AC2/AC5/AC6 (build-side):** fixture-driven Go unit tests that build through `dispatch build` and assert on the parsed `name`. Sub-second each; no live workflow. AC5 additionally leans on the existing golden/parity suite.
- **AC3/AC4 (reconcile id round-trip):** Go tests that (1) emit the capped name via `dispatch build`, (2) seed an active dir + `_archive` dir with the two shared-prefix `sd-b32` entities (distinct ids; one archived, one active), (3) build `reconcileOpts` with a stub roster carrying the capped name and a fixture team identity (the existing `reconcile_test.go` / `reconcile_session_test.go` stub-roster pattern), (4) run real `Reconcile` and assert the drift classes. Stub only the team roster and gh/git runners (as existing reconcile tests do); the entity/id-map read, id-prefix resolution, and decomposition run against real code and real on-disk fixtures — no mocked behavior under test.
- **Slug-fallback coverage:** one Go test over an `id-style: slug` long-slug fixture asserting the truncation fallback still produces a ≤64 `namePattern`-valid name (the id-less path), so both id-styles are exercised.
- **Cost/complexity:** low. All Go unit/fixture tests in one package, under a second total. The cross-file change (build id-cap + reconcile id-resolve + `entityRecord.id`) is the moderate piece, covered end-to-end by AC3/AC4.
- **Regression guard:** existing `go test ./internal/dispatch/` golden/parity suite (195 tests) must stay green — AC5 leans on it to prove short names are byte-unchanged.

## Stage Report: ideation

- DONE: Design a deterministic ≤64-char worker `name` that stays DECOMPOSABLE; exercise the riskiest path FIRST — trace the reconcile A/B + supersede/terminal-teardown name→entity decomposition (internal/dispatch/reconcile.go + internal/dispatch/build.go) and confirm the proposed capping form round-trips back to the correct entity (distinct long slugs must not collide), or record an auditable "no spike needed" naming the proven mechanism.
  Spike ran (throwaway, removed before commit): real `dispatch build` reproduced the 82-char name; truncated-head+sha256-6hex capping form produced two 64-char non-colliding names; `decompose()` still peeled stage/cycle; known-slug-set recompute round-tripped both long slugs back to their real on-disk slugs. Traced both call sites (build.go:435 derivedName; reconcile.go:366 decompose → active[slug] exact-match lookup). Results recorded in the "Spike" section of the body.
- DONE: Acceptance criteria each proven by a Go test over REAL `dispatch build` output — name ≤64 chars for a long slug AND it still resolves to the correct entity/cohort (no false reconcile Class-A, no missed teardown) — never a prose/code grep that merely spells the 64-char constraint.
  Six ACs written; each verified by a Go test over real `dispatch build` output + real `decompose`/`Reconcile` code. AC3/AC4 bind the expectation to a seeded on-disk entity set (independent source), proving correct resolution / no false Class-A. No AC is a string match over an instruction file.

### Summary

Designed a deterministic cap: when `{workerKey}-{slug}-{stage}` exceeds the budget, replace the slug component with `{truncated-head}-{6hex-sha256(full-slug)}`, preserving the prefix and `-{stage}` suffix so `decompose()` is undisturbed; reconcile resolves the capped token back to the real slug via the known active+archived slug set it already loads (no new on-disk state — stored short-id rejected as YAGNI). The cap target is 57 to reserve `-cycleN` headroom the FO appends downstream. The one doc touchpoint is the FO break-glass template (a one-line cap caveat); the `name` is an internal roster handle in no CLI/docs output. Riskiest mechanism (capped-name → entity round-trip) was exercised first via a throwaway spike against real code before committing to the design.

## Stage Report: ideation (cycle 2)

- DONE: Design a deterministic ≤64-char worker `name` that stays DECOMPOSABLE; exercise the riskiest path FIRST — trace the reconcile A/B + supersede/terminal-teardown name→entity decomposition (internal/dispatch/reconcile.go + internal/dispatch/build.go) and confirm the proposed capping form round-trips back to the correct entity (distinct long slugs must not collide), or record an auditable "no spike needed" naming the proven mechanism.
  Reworked per captain feedback to id-first. The riskiest mechanism SHIFTED to the id-prefix round-trip, so it was re-exercised (not asserted): real sd-b32 entities minted by `status --new`; `status --resolve prefix:367s2zrb`/`prefix:0qv3kaqs` resolved each fixed 8-char id-prefix to the correct entity (CLI); a Go spike over the real `decompose`/reconcile code confirmed id-form names (40 chars) decompose, distinct ids → distinct names, and round-trip resolves alpha→archived (Class-A) / bravo→active (no false Class-A). Results in the rewritten Spike section. Spike + sd-b32 fixture removed before commit.
- DONE: Acceptance criteria each proven by a Go test over REAL `dispatch build` output — name ≤64 chars for a long slug AND it still resolves to the correct entity/cohort (no false reconcile Class-A, no missed teardown) — never a prose/code grep that merely spells the 64-char constraint.
  Six ACs retained (AC1/2/5/6 build-side; AC3/4 now bind to id-prefix resolution against a seeded on-disk id set). Plus a slug-fallback test for the id-less path. No AC is a string match over an instruction file.

### Summary

Cycle-2 rework: dropped the re-invented `sha256(slug)` hash; the overflow cap now substitutes a FIXED-length prefix of the entity's existing sd-b32 `id` for the slug component (id-first), preserving the `worker_key` prefix and `-{stage}` suffix. Reconcile resolves the id-prefix token back to the entity by `HasPrefix(record.id, token)` against the active+archived id maps it already loads (the same resolution `status --resolve prefix:` performs) — no new on-disk state, no per-slug recompute. id-less `id-style: slug` workflows keep the slug-head-truncation fallback (no stable token, by design). `worker_key`→`ensign` trim evaluated and recommended OUT (cross-cutting worktree/branch rename + dual-form decomposition migration; zero budget benefit since id-form already fits ≤42). Riskiest mechanism (id-prefix round-trip, `status --resolve`) re-verified against real code and real sd-b32 entities before committing.
