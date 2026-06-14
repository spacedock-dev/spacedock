---
id: gwttp4rhr07bqtejzebn7hha
title: Reconcile Class D/E must never destroy committed or unpushed work — rethink repo-hygiene, consider an ownership lease
status: implementation
source: captain + FO (2026-06-14, this session) — the reconcile sweep flagged local main (ahead by 1, the unpushed README cut) as Class E "stale local main" and prescribed `git reset --hard origin/main`, which would have deleted the commit. Sibling to #369 (separated reconcile team-management A/B/C from repo-hygiene D/E, fixed trunk detection) and #370 (pr-merge trunk refit off pre-flip `next`); this addresses the D/E remedy-safety gap those left.
started: 2026-06-14T21:10:17Z
completed:
verdict:
score: 0.42
worktree: .worktrees/spacedock-ensign-reconcile-de-ownership-safety
issue:
sprint: 0203-fo-efficiency
---

The reconcile sweep's repo-hygiene classes (D: stale worktree branch; E: stale local main) can destroy work. Rethink them so reconcile never destroys committed or unpushed work and never mutates a worktree the running session does not own. Keep it simple; if needed, model ownership as a lease.

## Problem

The event-loop reconcile sweep (contract: `claude-fo-dispatch.md` Event Loop step 0; helper: `internal/dispatch/reconcile.go`) classifies repo-hygiene drift into Class D (worktree branch behind `origin/{trunk}` → `pull --rebase`) and Class E (local main "stale" → `git fetch && git reset --hard origin/{trunk} && rebuild`).

Three faults, found 2026-06-14 (after the #369 trunk fix was already live):

- **Class E detects the wrong direction.** `classE` counts `git rev-list --count origin/{trunk}..main` — commits local main has that origin lacks, i.e. AHEAD / unpushed — stores it in an `Ahead` field, then emits `reset main->origin/{trunk}`. So Class E fires precisely when the FO holds unpushed commits on main (an FO-direct README/state edit awaiting a push), and the remedy discards them. The class is labelled "stale local main" but detects unpushed-ahead, not behind.
- **The remedy is destructive.** `reset --hard` against a branch that legitimately carries local commits (FO write scope commits the workflow README and state transitions to main, then pushes on the captain's word) can lose committed work. The contract should never prescribe a hard reset that can destroy committed work.
- **Class D mutates un-owned worktrees.** It auto-rebases any active entity's worktree behind trunk, including a worktree another session (the Commander) is actively driving — a cross-session collision on someone else's branch.

## Spike (riskiest path, exercised first)

Three throwaway spikes ran against the real helper + real git fixtures before design was fixed. They seed the implementation's first tests.

- **Spike 1 — confirm the destructive-reset bug (the riskiest path).** A reconcile fixture (real `git init`, `origin/next` set via `update-ref`, then a commit on `main` so main is AHEAD of `origin/next` by 1) run through `Reconcile(... include={E} ...)`. Result: one Class-E item with `reason: "local main carries 1 commits not on origin/next; reset main->origin/next"`. The current helper emits a destructive `reset` remedy for an unpushed-ahead main — bug confirmed by exercising, not by reading. This fixture becomes AC-1's failing test.
- **Spike 2 — confirm behind/diverged are distinguishable with the same git runner.** From the same fixture, advancing `origin/next` by one commit (helper `repoBumpOriginNext`) yields `origin/next..main`=0, `main..origin/next`=1 (behind-only); then a further commit on main yields `origin/next..main`=1, `main..origin/next`=1 (diverged). Both counts come from `gitRunnerExec` with no new dependency. This proves the new classE can branch on the two counts.
- **Spike 3 — confirm the ownership signal is already in hand.** `decompose("spacedock-ensign-yaml-parser-migration-implementation", stages)` returns `slug="yaml-parser-migration", ok=true`. So "does the current trusted roster own this worktree's entity?" is answerable by decomposing each trusted-roster ensign name and matching its slug to the entity slug — the exact resolution classes A/B/C already perform. **No ownership-lease primitive is needed** (captain steer: lease only if necessary). Ownership = the entity's ensign is a member of the current session's trusted roster.

All three ran green/red-as-expected; the spike scripts were removed (throwaway). Mechanisms relied on: `gitRunnerExec` two-count comparison, `decompose`/trusted-roster membership, real-git fixtures matching `reconcile_test.go` helpers.

## Proposed approach

Two independent changes in `internal/dispatch/reconcile.go`, each turning a destructive/un-owned mutation into a direction-aware or ownership-gated remedy. The helper still only *detects and prescribes* — the FO acts — but it must never prescribe a remedy that can lose work or touch a peer's branch.

### Class E — direction-aware, never reset

`classE` currently counts only `origin/{trunk}..main` (AHEAD) and unconditionally prescribes `reset main->origin/{trunk}`. Replace with a two-count classification:

- `ahead = rev-list --count origin/{trunk}..main` (commits main has, origin lacks — unpushed)
- `behind = rev-list --count main..origin/{trunk}` (commits origin has, main lacks)

Then:

| ahead | behind | situation | remedy (non-destructive) | drift `reason` |
|------|--------|-----------|--------------------------|----------------|
| 0 | 0 | in sync | no item | — |
| 0 | >0 | behind only | `merge --ff-only origin/{trunk}` + rebuild | `local main behind origin/{trunk} by {behind}; ff-merge main<-origin/{trunk}` |
| >0 | 0 | ahead/unpushed | **report only** (FO pushes on captain's word) | `local main ahead of origin/{trunk} by {ahead} (unpushed); push when ready — no auto-reset` |
| >0 | >0 | diverged | **report only** | `local main diverged from origin/{trunk} (ahead {ahead}, behind {behind}); manual reconcile — no auto-reset` |

The `reason` string is the contract-visible remedy. The word `reset` never appears for ahead>0; only the behind-only case carries an automated remedy and it is `ff-merge` (which fails rather than discards if it cannot fast-forward). The driftItem keeps `Ahead` and now also populates `Behind` for E.

The behind-only `ff-merge` is safe because `merge --ff-only` refuses (non-zero exit, no mutation) the moment local main has diverged — git itself enforces non-destructiveness, so even a race between detection and action cannot lose a commit.

### Class D — ownership-gated, report un-owned

`classD` currently rebases any `active` entity's worktree behind trunk. Gate the *remedy* on ownership without dropping the *report*:

- Compute `owned` once per sweep: decompose each `spacedock:ensign` member of the **trusted** roster (`rosterTrusted == true`) and collect the set of slugs they resolve to (via `decompose` + `resolveSlugToken`, the A/B/C path).
- In classD, when a worktree is behind trunk:
  - if the entity's slug ∈ `owned` → emit the rebase remedy as today (`reason: "branch behind origin/{trunk} by {n}; pull --rebase (owned)"`).
  - else (slug not owned, OR roster untrusted) → emit a **report-only** item: same `Class:"D"`, `Behind`, `Trunk`, but `reason: "branch behind origin/{trunk} by {n}; peer-owned or un-owned — reporting only, not rebasing"`. A new boolean field `Owned bool` (json `owned,omitempty`) lets the FO branch deterministically rather than string-matching the reason.

This reuses the exact ownership signal A/B/C use; no lease file, no new state. When the roster is untrusted (bare git-only reconcile), *every* class-D worktree is report-only — correct, because without a trusted roster we cannot prove we own anything.

`classD`'s signature gains the `owned` set: `classD(active, repoRoot, trunk, owned, git)`. The `owned` set is built in `Reconcile` from `ensigns` + `stageNames` + the active/archived maps (already in scope), passed in alongside `trunk`.

### Contract doc diff (`skills/first-officer/references/claude-fo-dispatch.md`)

The action table at lines 228-229 prescribes the remedies the FO runs; it must change in lockstep. Unified diff:

```diff
-   - **D (stale branch)** → `git -C {worktree} pull --rebase origin {drift.trunk}`; halt on conflict per the rebase-conflict halt rule.
-   - **E (stale local main)** → `git -C {repo} fetch origin {drift.trunk} && git -C {repo} reset --hard origin/{drift.trunk} && cd {repo} && go build -o spacedock ./cmd/spacedock`.
+   - **D (stale branch)** → only when `drift.owned == true`: `git -C {worktree} pull --rebase origin {drift.trunk}`; halt on conflict per the rebase-conflict halt rule. When `drift.owned` is false the item is report-only — surface it to the captain; do NOT rebase a worktree the current session does not own.
+   - **E (local main drift)** → behind only (`drift.behind > 0 && drift.ahead == 0`): `git -C {repo} fetch origin {drift.trunk} && git -C {repo} merge --ff-only origin/{drift.trunk} && cd {repo} && go build -o spacedock ./cmd/spacedock`. Ahead/unpushed or diverged (`drift.ahead > 0`): report-only — surface `drift.reason` to the captain and NEVER `reset --hard`; the captain decides push vs. manual reconcile.
```

The reconcile-sweep prose at line 225 needs no change (the JSON-shape and class enumeration are unchanged). The one-line drift report at line 231 is unaffected (still `A=.. B=.. C=.. D=.. E=..`).

## Out of scope

- The A/B/C team-management classes and the trunk-detection fix (already shipped by #369/sr — `resolveTrunk` and the roster-trust gating are reused as-is, not re-touched).
- Introducing a standalone ownership-lease file or any new on-disk ownership state. The captain steer is "lease only if necessary"; spike 3 proved the trusted-roster membership signal already answers ownership, so this stays out.
- The FO actually *executing* the remedies (the helper detects and prescribes; the FO's event loop runs them). This entity changes what is prescribed, not the FO's execution machinery beyond the doc-diff wording.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 — Reconcile never prescribes a reset that can lose committed work.** For a local main carrying unpushed commits (ahead>0, behind 0) the Class-E remedy is report-only and its `reason` contains no `reset`; for a diverged main (ahead>0, behind>0) it is likewise report-only with no `reset`; for a behind-only main the remedy is `ff-merge`, not `reset`.
Verified by: a Go unit test over `internal/dispatch/reconcile.go` driving `Reconcile(... include={E} ...)` against three seeded real-git fixtures — ahead-only, behind-only, diverged — asserting (independent oracle = the seeded git state): the ahead-only and diverged drift items carry no substring `reset` (case-insensitive) in `reason` and set the expected `Ahead`/`Behind` counts; the behind-only item's `reason` prescribes `ff-merge`. Pinned fixture: spike-1's ahead-only graph (`git init` → `update-ref origin/next HEAD` → one commit on main) is the ahead-only case verbatim; behind-only adds `repoBumpOriginNext`; diverged adds a further `repoMakeCommit` on main (spike-2 proved all three count-pairs).

**AC-2 — Reconcile does not prescribe mutating a worktree the running session does not own.** A behind worktree whose entity's ensign is in the current trusted roster yields `owned:true` + a `pull --rebase` remedy; a behind worktree whose entity's ensign is absent from the trusted roster (or when the roster is untrusted) yields `owned:false` + a report-only remedy and no `rebase`/`pull` verb is prescribed for it.
Verified by: a Go unit test seeding two behind worktrees — one whose slug an in-roster `spacedock-ensign-{slug}-implementation` member decomposes to, one whose slug no member resolves to — and asserting the first drift item has `Owned==true` with a `rebase` remedy and the second has `Owned==false` with a report-only `reason` containing no `pull`/`rebase` verb. A second sub-case runs the same fixture with `--team-name` omitted and no session match (untrusted roster) and asserts BOTH class-D items are `Owned==false`/report-only. The team-config fixture follows `teamConfigJSONWithSession` + `makeFixtureWorktree` conventions; the oracle is the seeded roster+worktree state, not a prose grep.

## Test plan

- **Where:** new Go unit tests in `internal/dispatch/` (e.g. `reconcile_de_safety_test.go`), alongside the existing `reconcile_e_test.go` / `reconcile_trunk_test.go`. Reuse helpers: `repoGitInit`, `repoSetOriginNext`, `repoMakeCommit`, `repoBumpOriginNext`, `makeFixtureWorktree`, `teamConfigJSONWithSession`, `mustGit`, `formatDrift`.
- **AC-1 fixtures (Class E):** three real-git graphs (ahead-only, behind-only, diverged) driven through `Reconcile`, asserting on `Ahead`/`Behind` and the absence/presence of `reset`/`ff-merge` in `reason`. Cost: low — three tmpdir git repos, sub-second each; pattern proven by spikes 1+2.
- **AC-2 fixtures (Class D):** one repo with two behind worktrees + a team config carrying one matching and zero-matching ensign members, run twice (trusted roster via `--team-name`; untrusted via omitted name + non-matching session). Asserts `Owned` and remedy verbs. Cost: low; pattern proven by `reconcile_trunk_test.go` + spike 3.
- **Existing-test impact:** `reconcile_e_test.go`'s `TestReconcileEDetectsAndResetAdvancesMain` currently asserts the reset path on an ahead-main fixture. Under the new contract that fixture (ahead-only) becomes report-only, so this test must be updated to assert the report-only `reason` (no `reset`) — flagged for implementation; it is a contract change, not a regression. The `go build`-rebuild half (`TestReconcileEGoBuildIsRunnable`) stays as-is (the behind-only remedy still rebuilds).
- **Doc-diff verification:** the contract wording change is verified by the behavior it points at (the `reason` strings the tests assert), not by a prose grep. The doc diff in this body is applied verbatim by implementation.
- **No live workflow test needed:** the change is entirely in helper detection + prescribed `reason` strings; unit fixtures over seeded git/roster state are the independent oracle.

## Stage Report: ideation

- DONE: Design the Class-E safety rethink (distinguish behind from ahead/unpushed; report diverged; never reset)
  Two-count classification (`origin/{trunk}..main` ahead, `main..origin/{trunk}` behind): behind-only → `ff-merge`+rebuild; ahead/diverged → report-only, `reason` never contains `reset`. See ## Proposed approach Class E table.
- DONE: Fix the classE detection bug (counts ahead yet labels "stale local main" + emits reset)
  Root cause confirmed by spike 1 (real-git ahead fixture → current helper emits `reset main->origin/next`). New design replaces the single-count + unconditional-reset with the direction table.
- DONE: Exercise the riskiest path first (spike): ahead-main fixture asserts NO reset remedy emitted
  Spike 1 ran the real helper over a real-git ahead fixture and FAILED on the current code (reset emitted) — bug proven by exercising. Spike fixture becomes AC-1's ahead-only case. Spikes 2 (behind/diverged counts) and 3 (decompose ownership) ran green.
- DONE: Design the Class-D ownership safety (auto-rebase only an owned worktree; report peer-owned)
  Ownership = entity's ensign ∈ current trusted roster (decompose + resolveSlugToken, the A/B/C signal). `owned` set built in Reconcile, passed to classD; un-owned/untrusted → report-only `Owned:false`. No lease primitive (captain steer satisfied).
- DONE: ACs bound to reconcile UNIT TESTS over seeded repos (oracle = seeded git/worktree state)
  AC-1 over three real-git E fixtures; AC-2 over a two-worktree + roster fixture, run trusted + untrusted. No prose-grep oracle. Fixtures pinned in body + Test plan.
- DONE: Pin the fixtures in ideation
  AC-1 fixtures = spike-1 ahead graph + `repoBumpOriginNext` (behind) + extra `repoMakeCommit` (diverged). AC-2 = `makeFixtureWorktree` ×2 + `teamConfigJSONWithSession`. Named in Test plan.
- DONE: Out-of-scope recorded (A/B/C + trunk detection, already #369/sr; no standalone lease)
  See ## Out of scope.
- DONE: Doc diff for user-visible contract change (FO action table)
  Unified diff against `claude-fo-dispatch.md` lines 228-229 in ## Proposed approach → Contract doc diff. Implementation applies verbatim.

### Summary

Designed both safety rethinks grounded in the real `reconcile.go`. The classE bug was confirmed by an exercised spike (ahead-main → destructive `reset` emitted on current code), not by reading; the fix is a two-count direction table where `reset` never appears for ahead>0 and the only automated remedy (`ff-merge`) is git-enforced non-destructive. Class-D ownership reuses the existing trusted-roster decompose signal — no lease primitive needed, satisfying the captain's "simple" steer. ACs bind to Go unit tests over seeded real-git/roster fixtures (independent oracle), with fixtures pinned and the existing `reconcile_e_test.go` ahead-fixture flagged for update to the report-only contract.

## Stage Report: implementation

- DONE: Implement the Class-E two-count direction table in internal/dispatch/reconcile.go (behind→ff-merge; ahead/diverged→report-only, reason never contains `reset`); AC-1 unit tests over three seeded real-git fixtures
  `classE` now counts both `origin/{trunk}..main` and `main..origin/{trunk}`; reasons reworded to avoid the `reset` substring (spec's "no auto-reset" wording would itself fail AC-1 — used "no auto-rewrite"). Tests in reconcile_de_safety_test.go: ahead-only, behind-only (ff-merge), diverged. Commit 3c0408fa.
- DONE: Implement Class-D ownership safety (rebase only an in-trusted-roster worktree; un-owned/untrusted → report-only Owned:false, no pull/rebase verb); AC-2 unit tests
  New `ownedSlugs` builds the owned set from the trusted roster via the A/B/C decompose+resolveSlugToken path; classD gains an `owned` param. Tests assert in-roster owned→rebase, absent→report-only, and untrusted-roster→both report-only. No lease primitive. Commit 3c0408fa.
- DONE: Apply the contract doc-diff to claude-fo-dispatch.md; UPDATE TestReconcileEDetectsAndResetAdvancesMain to the report-only contract; keep behind-only rebuild; `go test ./...` green
  FO step-0 D/E action table updated verbatim. Both ahead-fixture tests (reconcile_e_test.go + reconcile_test.go's TestReconcileFiveClasses) now assert no `reset` substring. TestReconcileEGoBuildIsRunnable kept as-is. `go build ./...` + `go test ./...` all green.

### Summary

Class E was rebuilt as a two-count direction table that never prescribes a destructive reset: behind-only carries a git-enforced-safe `ff-merge`+rebuild, ahead/unpushed and diverged are report-only. Class D now gates its `pull --rebase` on the worktree entity's slug being in the current trusted roster (reusing the existing decompose signal, no new lease primitive), reporting un-owned/untrusted worktrees with `Owned:false` and no mutation verb. One deviation from the ideation's literal reason strings: the spec's "no auto-reset" wording contains the `reset` substring AC-1 forbids, so the report-only reasons say "no auto-rewrite" instead — the binding safety property (no `reset` in reason) won over the example prose. Full `go test ./...` is green.
