---
id: vkatjs25g9a9gmk3jtvx5ce0
title: merge guard refuses a terminal transition with no preceding worker report
status: ideation
source: "Captain CL, 2026-08-18, from the live-lane inventory reframe. Failing assertion: internal/ensigncycle/shared_keep_moving_durable_test.go:103, 'first terminal transition must follow worker report', red in two consecutive claude-live runs (32092321763 attempt 2 and 32105482382) while the FO's own final messages claimed the reports had completed."
started: 2026-08-18T18:41:27Z
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:vkatjs25g9a9gmk3jtvx5ce0:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:vkatjs25g9a9gmk3jtvx5ce0-ideation-1
              briefing:
                id: briefing:vkatjs25g9a9gmk3jtvx5ce0:ideation:attempt-1:revision-1
                digest: sha256:ebc208b306ea71dbb15d03c4380e1f6a66f3038164fb205f07b050ba6d685973
                request-digest: sha256:a0482d868df68d3814925440b283f74fadb4cac34afc4aa798046bf3d9967268
                room-ref: ./merge-guard-requires-preceding-report/review/ideation/briefing-1
---

The terminal ceremony checks that the merge hook ran. It does not check that the work it is finalizing was ever reported. Make it refuse when no report commit precedes the terminal transition.

## Problem

`merge guard` drives the terminal ceremony as one ordered envelope — arm the mod-block, detect hook completion by state delta, terminalize, archive with a path-scoped commit, publish state. It owns the sequence so steps cannot be combined, skipped, or reordered.

It verifies the hook. It never verifies the report.

So an FO that trusts a worker's completion message over durable state can finalize an entity whose report never landed, and the ceremony's own atomicity hides it: one verb takes the entity from validation to archived-and-pushed with no intermediate state where the absence would show.

Observed twice in consecutive claude-live runs. The durable check that caught it walks the entity's git history and requires the terminal-transition commit to follow a commit carrying the completing stage's report.

The correct shape is already what a healthy run produces, verified on this session's own terminalizations: the worker commits its report touching no frontmatter (`+15` and `+43` lines, zero status change), then `merge guard` flips `status`, `verdict`, and `completed` in a later commit. Different writers, naturally ordered. Nothing has to change about how a good run behaves.

Two sub-shapes are possible and the evidence does not separate them: the report and the flip batched into one commit, or the report never landing at all. The second is the likelier one, because our deliberate coalescing batches FO-side state writes together and never mixes them with a worker's report.

## Sub-shape determination (settled from the retained streams)

Both failing runs' artifacts were downloaded and read (`gh run download 32092321763` / `32105482382`, artifact `runtime-live-e2e-claude-live-claude-sonnet-5`). The red scenario is `smallest-sufficient-mechanism` (the grader in `shared_keep_moving_durable_test.go` is shared): run 32092321763 failed `ready-one` AND `ready-two`; run 32105482382 failed `ready-two` only. The observed reason — line 103's `firstTerminal <= report` — structurally requires a report commit to EXIST in history (a never-landed report produces line 101's "missing worker report" instead), so "the report never landed at all" is refuted for both runs. What the streams show:

- **Run 32092321763 (both entities): the worker batched the flip with its report.** The ensign ran `spacedock status --set ready-one status=done` and it succeeded ("status: ready -> done" — a silent hand-terminal-set on a gate-less entity), then appended its report via Edit and committed BOTH in one commit: "ready-one: append stage report and advance to done", 10 insertions 1 deletion (the deletion is `status: ready`). Same shape for ready-two ("Ready Two: complete ready stage, advance to done"). The first terminal transition IS the report commit.
- **Run 32105482382 (ready-two): the worker never committed.** It appended the report via Edit and signaled done without committing — the classic signal-without-commit. `merge guard ready-two --verdict passed` then finalized against the dirty tree, and the guard's own archive commit swept the uncommitted report in with the flip. The run's final `git log` shows it plainly: `86e83b7 archive ready-two (merge guard)` directly follows `d22d79d dispatch: ready-two entering ready`, no worker commit between. The green `ready-one` in the same run is the healthy control: worker commit "Ready One: complete ready stage" (+9 lines, zero frontmatter change), then the guard's flip.

Both proximate causes converge on the same durable shape — report and terminal flip in one commit — and one ceremony-time check catches both.

## Proposed approach

Before any mutation on the finalize path (`--verdict passed` only), `merge guard` walks the entity's own git history — state it already reads and commits to — and refuses unless a commit carrying the completing stage's report strictly precedes the first commit introducing terminal state:

1. Resolve the live entity path (flat or folder form, as `captureArchiveState` already does) relative to `FindGitRoot(roots.entityDir)`.
2. `git log --follow --format=%H -- <path>` newest-first, reversed in code (`--follow` with `--reverse` silently truncates at renames — spike finding; the durable grader already reverses in code), then `git show <hash>:<path>` per commit, parsed with the existing `ParseFrontmatterData`.
3. `firstTerminal` := first commit whose blob has `status` equal to the README's terminal stage, or nonempty `completed`, or nonempty `verdict`. `firstReport` := first commit whose blob carries the `## Stage Report: {completing stage}` header (exact or its `(cycle N)` form; never a bare prefix, so stage `ready` cannot match `ready-extra`).
4. The completing stage is the entity's current status; when current status is already terminal (crash re-run, or a rogue flip), it is the status in the parent blob of `firstTerminal`.
5. Refuse — exit 1, nothing mutated; the check runs before the mod-block clear — when `firstReport` is absent ("no committed worker report for stage {stage}; commit the report path-scoped and re-run") or when `firstTerminal <= firstReport` ("terminal state does not follow the committed report; this history cannot be finalized as-is — escalate"). Otherwise proceed unchanged. The refusal message carries its own recovery at fire time, per the verb's existing signal convention.

Coverage: run-2's shape refuses with a recovery the FO can execute — commit the report, re-run (AC-1's red-to-green). Run-1's shape refuses without one: history is append-only, so a worker-committed flip has already made the journey permanently red; the guard's job there is converting silent finalization into a loud stop instead of archiving and publishing a poisoned entity. The sanctioned crash-recovery re-run (report committed, flip committed, archive missing) passes the ordering check unchanged.

Scope decisions:
- `--verdict rejected` is exempt: rejection is the sanctioned no-work terminal; killing a never-worked entity must not require a fabricated report.
- `--rework` determination (answering Out of scope's open question): the gap does NOT apply — rework performs no terminal transition; it supersedes and routes to the nonterminal `feedback-to`.
- A non-git entity root skips the check, matching `commitArchiveMove`'s existing carve-out — no durable history exists to verify.
- A workflow whose completing stage never wrote a stage report will now refuse a passed finalize; the recovery is the stage-report protocol every completing stage already owes.

## Mechanism justification

New mechanism: one git-history read inside finalize. Value AC served: AC-1. Alternatives considered:
1. Require the report text in the entity file at guard time (no history read) — insufficient: both failing runs HAD the report on disk when the guard ran; the violated contract is durable ORDER, visible only in history.
2. Refuse on a dirty entity file — insufficient: run 1's tree was clean at guard time, and unrelated dirt would misfire.
3. A new tracking field (report-commit SHA in frontmatter) — a new tracking mechanism, and forgeable by the same weak FO the ceremony guards against.

## Spike record

Riskiest mechanism (the history-walk decision rule) exercised first, red and green both proven:
- Replayed all four shapes as git fixtures and ran the exact reads plus rule: healthy replay (dispatch commit, then a worker report commit touching no frontmatter) PROCEEDs; crash re-run PROCEEDs; run-1's batch REFUSEs; run-2's dirty-tree REFUSEs, then PROCEEDs after a path-scoped report commit.
- Ran the same rule over six real archived entities in this repo's live state checkout (pre0-cut-idempotent-on-rerun, remove-tautological-workflow-tests, collapse-duplicate-edge-marketplace-routes, prerelease-ships-stable-stamped-default-artifact, run-rejection-journey-in-team-mode, red-auto-continue-gate-bypass): all PASS — the report commit strictly precedes the first terminal commit in every real healthy history.
- Folded-in finding: `git log --follow --reverse` truncates at the archive rename; collect newest-first, reverse in code.

## Out of scope

Changing what the merge hook does, the archive step, or the ordering the envelope already enforces. The `--rework` path unless the same gap applies there.

## Expected surface and tolerance

Estimate net LOC change: +260 across up to 12 files (revised at ideation after reading the code: `internal/status/merge.go` ~+85; a new merge-guard report test file ~+150; the 6–8 finalize-path fixture entities under `internal/status/testdata/merge-*-workflow/` gain a committed stage report, ~+5 each — none carries one today, so the new refusal correctly fires on them; `docs/site/reference/command-reference.md` +2/−1). Insertions ~+265, deletions ~−5. Tolerance: net ±40%, files ±4 (the fixture fallout count firms up at implementation). Semantics changed: `merge guard` gains exactly one refusal condition — a previously-accepted `--verdict passed` finalize can now exit 1 with nothing mutated. No command grammar, stored-format, or authority changes; `rejected`, `--rework`, arm, and blocked behavior unchanged.

## Documentation diff

`docs/site/reference/command-reference.md`, the `merge guard <slug> --verdict passed|rejected` row.
Before: "…and the `pr` merge sentinel is retained through archive as durable delivery proof."
After: "…and the `pr` merge sentinel is retained through archive as durable delivery proof. A `passed` finalize refuses when the entity's own history carries no commit with the completing stage's report preceding the first terminal transition — commit the worker's report path-scoped and re-run."
No FO-skill text change: the refusal message names its own recovery at fire time (the verb's carried-at-fire-time convention), so no resident prose is added and no skill smoke tests are triggered.

## Test plan

Go CLI-level tests driving `MergeGuard` against git fixtures (extending the existing `driveMergeGuard` harness in `internal/status`), each with the change that would fail it:
1. Uncommitted-report refusal and recovery (AC-1): exit 1 naming the missing committed report, entity unmutated, no archive; commit the report path-scoped; re-run finalizes and archives. Fails on today's binary, which finalizes either way (observed live, run 32105482382).
2. A single worker commit batching the report with a status flip refuses; no archive. Fails on today's binary (observed live, run 32092321763).
3. Healthy replay — dispatch commit, then a report-only commit touching no frontmatter — finalizes and archives (AC-2). Fails if the guard misreads history order or the header match.
4. Crash re-run (report committed, flip committed, archive missing) still finalizes. Fails if the check demands "no terminal state anywhere" instead of ordering.
5. `--verdict rejected` with no report still finalizes. Fails if the exemption is dropped.
6. A non-git entity root still finalizes. Fails if the carve-out is dropped.
7. `## Stage Report: {stage} (cycle 2)` counts as the report; `## Stage Report: ready-extra` does not count for stage `ready`. Fails on prefix matching.
Cost: unit-level only; no live run needed — the two retained streams are the live evidence, and the durable grader already encodes the contract this guard enforces at ceremony time. Existing arm/blocked/rejected merge-guard tests stay untouched.

## Follow-up candidates (out of this task's surface)

- A hand `status --set <slug> status=<terminal>` on a gate-less entity succeeds silently with no verdict (run 32092321763 proved it live); the pending-approval refusal only covers gated entities. Closing that hole is a separate entity.
- Worker-side prevention (the ensign contract already forbids frontmatter mutation; the binary permitted it) is upstream discipline, separately owned.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - No entity reaches a terminal stage without its completing stage's report preceding the transition in durable history.**
This is the measuring AC: the count of terminal transitions with no preceding report commit must be ZERO. Verified by driving `merge guard` against a fixture entity whose report was never committed and observing a refusal, then committing the report and observing it proceed — and against the second observed shape, a single worker commit batching the report with a status flip, observing a refusal with nothing mutated. Fails on today's binary, which finalizes either way.

**AC-2 - A healthy run is unaffected.**
Verified by replaying this session's real terminalization shape — worker report commit touching no frontmatter, then the terminal flip — and observing `merge guard` proceed unchanged. Fails if the guard refuses a correctly-ordered run, which would block every normal merge.

## Stage Report: ideation

- DONE: Settle the sub-shape before designing: from the two failing runs' retained streams, determine whether the report and terminal flip were batched into one commit, or whether the report never landed at all. The fix differs and the evidence has not been read.
  Both runs downloaded and read; both sub-shapes occurred and both end batched: run 32092321763's workers hand-flipped status=done and committed flip+report together; run 32105482382's ready-two worker never committed and the guard's archive commit swept the dirty report in (final log: `86e83b7 archive ready-two (merge guard)` directly after `d22d79d dispatch:`). "Never landed at all" refuted — line 103's reason structurally requires a report in history.
- DONE: Build the check from state merge guard already reads — the entity's own history in the state checkout — and do not add a new tracking mechanism.
  Design walks the entity's own git history via the package's existing `runGitCmd`/`ParseFrontmatterData`/`terminalStageName`, refusing a passed finalize unless a completing-stage report commit strictly precedes the first terminal-state commit; mechanism-justification section records the rejected new-tracking-field alternative.
- DONE: Prove a healthy run is unaffected by replaying the real shape this session produced: a worker report commit touching no frontmatter, then the terminal flip.
  Replayed as a git fixture (dispatch commit → report-only commit → flip): PROCEED before and after the flip; the same rule run over six real archived entities in this state checkout all PASS; run 32105482382's own green ready-one is the live healthy control.

### Summary

Read both failing runs' retained streams and settled the open sub-shape question: both are the batched shape in durable history, reached two different ways (worker hand-flip + batch commit; worker signal-without-commit + guard sweep), so one ceremony-time ordering check covers both. Designed the refusal into `finalize()` (`--verdict passed` only, before any mutation), spiked the history-walk rule red-and-green against replayed failing shapes and six real archived histories, and recorded the `--follow --reverse` git quirk, the rejected/`--rework`/non-git exemptions, a revised surface estimate (+260 net, up to 12 files with fixture fallout), the concrete command-reference doc diff, and two follow-up candidates (hand terminal `--set` on gate-less entities; worker-side frontmatter discipline).
