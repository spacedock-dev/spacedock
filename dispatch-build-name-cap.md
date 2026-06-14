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

The design's soundness rests on one unverified mechanism: that a length-capped `name` still round-trips through `decompose()` (`internal/dispatch/reconcile.go`) back to the correct entity, so the cap does not introduce a false reconcile Class-A (a spurious shutdown against a live worker) or a missed terminal teardown. That mechanism was exercised before committing to the design, with a throwaway test driving the real `dispatch build` surface and the real `decompose`/reconcile code. Results:

- **Bug reproduced over real `dispatch build`.** A 57-char slug at `backlog` produced an 82-char `name` (`spacedock-ensign-{slug}-backlog`) — well past the Agent-tool 64-char ceiling, confirming the failure `Agent(name=...)` would reject.
- **Naive head-truncation collides.** Two distinct 56-char slugs sharing a 42-char common prefix truncate to the identical 32-char head, so head-only truncation merges two entities' cohorts — exactly the collision the github#366 note warns about.
- **Truncated-head + slug-hash fits and stays distinct.** Replacing the slug component with `{truncated-head}-{6-hex-of-sha256(full-slug)}` produced two 64-char names that pass `namePattern` and differ in the hash segment (`…-deconf-41535c-implementation` vs `…-deconf-31f6ce-implementation`) despite the shared head. No collision.
- **Capped name still decomposes its stage/cycle.** Because the cap rewrites only the slug head and leaves the `-{stage}` suffix (and any FO-appended `-cycleN`) untouched, `decompose()` peels `stage` and `cycle` from a capped name unchanged. The embedded `{head}-{hash}` token is a stable, collision-resistant cohort key for Class-B / supersede grouping.
- **Round-trip to the real slug works via the known-slug set.** `decompose()` alone cannot string-split a capped name back to the on-disk slug (`active[slug]` is an exact-match map lookup). Resolving the capped token against the reconcile sweep's existing `active`/`archived` slug maps — exact match first, else recompute the capped form per known slug and compare — recovered the correct full slug for both long slugs and left short (uncapped) names resolving by exact match.

The spike is throwaway (removed before commit); what it teaches seeds the implementation's first test (see Test plan).

## Proposed approach

**1. Deterministic cap inside `dispatch build`.** Introduce a capping helper applied to `derivedName` (`internal/dispatch/build.go:435`) before the Rule-7 length/safety checks. When `{workerKey}-{slug}-{stage}` already fits the cap, emit it unchanged (the common case — short slugs are untouched, preserving today's readable names and all existing golden fixtures). When it would exceed the cap, replace the **slug component only** with `{truncated-head}-{hash}` where:
   - `hash` = first 6 lowercase-hex chars of `sha256(full-slug)` — a stable, deterministic function of the full slug, so distinct slugs almost never collide and the same slug always caps identically across runs and across FOs.
   - `truncated-head` = the leading bytes of the slug that fit the remaining budget after reserving `workerKey`, the two joining hyphens, `hash`, and `-{stage}`, with any trailing hyphen trimmed so the result still matches `namePattern` (`^[a-z0-9][a-z0-9-]*[a-z0-9]$`) and carries no `--`.
   The `workerKey` prefix and `-{stage}` suffix are preserved verbatim, so `decompose()`'s prefix-strip and stage-peel are undisturbed.

**2. Cap target reserves cycle headroom.** The FO appends `-cycleN` to the base name downstream of `dispatch build` (the supersede / fresh-redispatch roster convention; cycle 3 is the feedback-flow escalation ceiling). A base name capped to exactly 64 overflows once `-cycle2` is appended (71 chars). So the cap **target** is `64 − len("-cycle9") = 57`, leaving room for any single-digit cycle suffix. The emitted base name is therefore ≤57; `{base}-cycle3` ≤64. (`dispatch build` itself never emits a cycle — it has no cycle field — so the headroom is reserved, not consumed, in the build output.)

**3. Resolve capped names back to the entity in `decompose`/reconcile.** Make the reconcile decomposition slug-aware so a capped name resolves to its real entity slug. Concretely, thread the known active+archived slug set (which `Reconcile` already loads — `reconcile.go:201-202`) into the name→slug resolution: for each member name, try exact decomposition first (uncapped names and the common case), and when the decomposed slug token does not match any known slug, recompute the capped form for each known slug at the decomposed stage and match. The matched known slug becomes `d.slug` for the Class-A entity lookup and Class-B cohort key. Names that match no known slug stay `ok=false` (unchanged behavior — an unknown member is not classified). This keeps Class-B grouping correct (capped token is a consistent cohort key even without resolution) and Class-A/C lookups correct (resolved real slug hits the entity map).

**4. FO break-glass template note.** The FO break-glass manual-dispatch template (`skills/first-officer/references/claude-first-officer-runtime.md:109`, `name="{worker_key}-{slug}-{stage}"`) is the one contract-prose touchpoint that also constructs a name and would also overflow for a long slug. It is a degraded fallback used only when `dispatch build` is unavailable, and it already omits much of the production helper's output. The doc change is a one-line caveat on that template noting the name must be capped to ≤64 the same way `dispatch build` caps it (pointing at the helper as the authority), so the fallback path is not silently broken. This is the only doc diff (see below); the `name` is an internal roster handle and appears in no docs-site page or CLI output.

### Doc diff

`skills/first-officer/references/claude-first-officer-runtime.md`, break-glass template (around line 109):

- Before: `    name="{worker_key}-{slug}-{stage}",`
- After: `    name="{worker_key}-{slug}-{stage}",  // if this exceeds 64 chars, cap it the way `spacedock dispatch build` does: keep the {worker_key} prefix and -{stage} suffix, truncate the slug head, and append a short stable hash of the full slug`

(Exact wording to be finalized in implementation; the load-bearing content is the ≤64 cap instruction plus the keep-prefix/keep-stage/hash-the-slug recipe.)

## Out of scope

- **Stored short-id resolution / a persisted name↔slug registry.** The known-slug-set recompute (approach §3) resolves capped names with no new on-disk state, so a stored short-id table is unnecessary (YAGNI). Rejected in favor of the stateless recompute.
- **Changing `dispatch_file_path` keying.** The dispatch filename has its own separate `dispatchFileNameMaxLen` (251) cap and is keyed on `team_name + derivedName`; per github#366, only `name` needs the 64-char cap. The file path is left as-is.
- **Capping `name` to be byte-identical to a hand-shortened or human-chosen name.** The cap is deterministic and machine-derived; no human-override path is introduced.
- **Re-running `dispatch build` with a cycle parameter to re-cap cycled names.** Considered and rejected: it would add a cycle field to the build schema and a build round-trip on every feedback cycle. Reserving cycle headroom in the base cap (§2) is simpler and keeps the FO's append-only cycle convention.
- **Hash-collision recovery beyond detection.** A 6-hex (24-bit) slug hash makes a collision between two long slugs sharing a truncated head astronomically unlikely; this task does not add a collision-retry/widen-hash mechanism. If a real collision is ever observed, widening the hash is a follow-up.

## Acceptance criteria

- **AC1 — `dispatch build` emits a `name` ≤64 chars for a long slug.** For an entity whose `{workerKey}-{slug}-{stage}` would exceed 64, the real `dispatch build` output's `name` field is ≤64 chars and matches `namePattern`.
  - Verified by: a Go test that drives `dispatch build` in-process (the `runNative` harness) over a fixture entity with a long slug and asserts `len(name) <= 64` and `namePattern.MatchString(name)` on the parsed stdout JSON. Fails today (name is 82 chars).
- **AC2 — distinct long slugs produce distinct names (no cohort collision).** Two distinct entities whose slugs share a long common prefix produce two different `name` values from real `dispatch build`.
  - Verified by: a Go test building both fixtures through `dispatch build` and asserting the two emitted `name` values differ. Fails under naive head-truncation (identical names).
- **AC3 — a capped name resolves back to the correct entity slug in reconcile.** Given the capped `name` emitted by `dispatch build` for a long-slug entity that is archived, the reconcile sweep classifies it Class-A against the **correct** slug, and does NOT mis-resolve it to a different long-slug entity that is still active.
  - Verified by: a Go test that emits the capped name via `dispatch build`, seeds an `active` + `_archive` entity set (one archived long-slug entity, one active long-slug entity sharing the prefix), runs the real reconcile classification over a roster carrying the capped name, and asserts exactly one Class-A entry pointing at the archived slug and none against the active one.
- **AC4 — capped names do not produce a false Class-A against a live entity.** A capped name whose entity is still active (status ≠ terminal, not archived) produces NO Class-A drift entry.
  - Verified by: a Go test running real reconcile over a roster with the capped name of an active long-slug entity and asserting zero Class-A entries for it (guards the false-shutdown hazard called out in the Problem section).
- **AC5 — short slugs are unchanged.** An entity whose `{workerKey}-{slug}-{stage}` already fits 64 emits exactly that name (no hash, no truncation), and existing build golden/parity fixtures still pass.
  - Verified by: a Go test asserting the uncapped form for a short slug, plus the existing `go test ./internal/dispatch/` golden/parity suite passing unchanged.
- **AC6 — base name reserves cycle headroom.** A capped base name plus a realistic FO cycle suffix (`-cycle3`) stays ≤64.
  - Verified by: a Go test that takes the capped `name` from `dispatch build` for a long slug and asserts `len(name + "-cycle3") <= 64`.

Every AC's proof is a Go test over real `dispatch build` output and the real `decompose`/reconcile code — none is a substring/regex match over an instruction file. AC3/AC4 bind the expectation to an independent source (the seeded on-disk entity set), so a wrong resolution fails the check.

## Test plan

- **Where:** `internal/dispatch/` (alongside `build_*_test.go` and `reconcile_*_test.go`), reusing the existing in-process harness (`runNative`, `readmeWorktree`, `entityFM`, `writeFile`, `gitInit`) and the real `decompose`/`Reconcile` functions. No new test infrastructure.
- **AC1/AC2/AC5/AC6 (build-side):** fixture-driven Go unit tests that build through `dispatch build` and assert on the parsed `name`. Cheap (sub-second each), no live workflow needed — the spike already showed these run in the harness.
- **AC3/AC4 (reconcile round-trip):** Go tests that (1) emit a capped name via `dispatch build`, (2) seed an active dir + `_archive` dir with the colliding-prefix long-slug entities, (3) build a `reconcileOpts` with a stub roster carrying the capped name and a fixture team identity (mirroring the existing `reconcile_test.go` / `reconcile_session_test.go` stub-roster pattern), (4) run the real `Reconcile` and assert the drift classes. Stub only the team roster and gh/git runners (as the existing reconcile tests do); the entity-map read and decomposition run against real code and real on-disk fixtures — no mocked behavior under test.
- **Cost/complexity:** low. All tests are Go unit/fixture tests in one package; estimated under a second total. No CLI-subprocess or live-workflow tests needed — the claim is build-output bytes plus reconcile classification, both exercisable in-process. The cross-file change (build cap + reconcile resolve) is the only moderate-complexity piece and is covered by AC3/AC4 end-to-end.
- **Regression guard:** the existing `go test ./internal/dispatch/` golden/parity suite (195 tests today) must stay green — AC5 leans on it to prove short-slug names are byte-unchanged.

## Stage Report: ideation

- DONE: Design a deterministic ≤64-char worker `name` that stays DECOMPOSABLE; exercise the riskiest path FIRST — trace the reconcile A/B + supersede/terminal-teardown name→entity decomposition (internal/dispatch/reconcile.go + internal/dispatch/build.go) and confirm the proposed capping form round-trips back to the correct entity (distinct long slugs must not collide), or record an auditable "no spike needed" naming the proven mechanism.
  Spike ran (throwaway, removed before commit): real `dispatch build` reproduced the 82-char name; truncated-head+sha256-6hex capping form produced two 64-char non-colliding names; `decompose()` still peeled stage/cycle; known-slug-set recompute round-tripped both long slugs back to their real on-disk slugs. Traced both call sites (build.go:435 derivedName; reconcile.go:366 decompose → active[slug] exact-match lookup). Results recorded in the "Spike" section of the body.
- DONE: Acceptance criteria each proven by a Go test over REAL `dispatch build` output — name ≤64 chars for a long slug AND it still resolves to the correct entity/cohort (no false reconcile Class-A, no missed teardown) — never a prose/code grep that merely spells the 64-char constraint.
  Six ACs written; each verified by a Go test over real `dispatch build` output + real `decompose`/`Reconcile` code. AC3/AC4 bind the expectation to a seeded on-disk entity set (independent source), proving correct resolution / no false Class-A. No AC is a string match over an instruction file.

### Summary

Designed a deterministic cap: when `{workerKey}-{slug}-{stage}` exceeds the budget, replace the slug component with `{truncated-head}-{6hex-sha256(full-slug)}`, preserving the prefix and `-{stage}` suffix so `decompose()` is undisturbed; reconcile resolves the capped token back to the real slug via the known active+archived slug set it already loads (no new on-disk state — stored short-id rejected as YAGNI). The cap target is 57 to reserve `-cycleN` headroom the FO appends downstream. The one doc touchpoint is the FO break-glass template (a one-line cap caveat); the `name` is an internal roster handle in no CLI/docs output. Riskiest mechanism (capped-name → entity round-trip) was exercised first via a throwaway spike against real code before committing to the design.
