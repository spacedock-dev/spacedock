---
id: vkatjs25g9a9gmk3jtvx5ce0
title: merge guard refuses a terminal transition with no preceding worker report
status: validation
source: "Captain CL, 2026-08-18, from the live-lane inventory reframe. Failing assertion: internal/ensigncycle/shared_keep_moving_durable_test.go:103, 'first terminal transition must follow worker report', red in two consecutive claude-live runs (32092321763 attempt 2 and 32105482382) while the FO's own final messages claimed the reports had completed."
started: 2026-08-18T18:41:27Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-merge-guard-requires-preceding-report
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
              resolution:
                type: Resolution
                id: resolution:spacedock:vkatjs25g9a9gmk3jtvx5ce0:ideation:1
                briefing: briefing:vkatjs25g9a9gmk3jtvx5ce0:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-18T19:32:26.228606Z"
                decision: approve
                reason: 'Captain approved in chat: ''approve those 4, and have them be on a pr stack.'' Accepts the ceremony-time ordering refusal, its exemptions, and the +260 surface. Tip of the stack; CI approved on the tip.'
              application:
                target-stage: implementation
                state: consumed
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

### Feedback Cycles

- Cycle 1: REJECTED — validation false-positive matrix; surface 8 files/net +527 vs estimate ≤12 files/net +260±40% (203%, 42% past the +371 ceiling); AC unchanged. Material finding: `entityBlobHistory` follows renames in the log but then reads every blob at the entity's CURRENT path, so any history crossing a rename — slug rename, flat→folder conversion, unarchive — exits 1 on `git show` and can never be finalized. Trigger is sanctioned and observed: `gate record --round` requires folder form, and this repo really did convert and merge two live entities. Separately: the ideation Spike record's six real archived histories, cited as PASS, all hard-error under the shipped read, so that spike exercised a different rule than the one that landed.

### Design-reset decision (surface past declared tolerance)

Surface is 203% of the approved estimate, past this entity's own declared ±40% tolerance, so the workflow requires a recorded design-reset decision before any further correction round.

**DECISION: RECONFIRM — Captain CL, 2026-08-18, in chat ("reconfirm and file").** The approved design stands; cycle 2 corrects the rename-crossing read inside `entityBlobHistory` and re-establishes the spike evidence against the shipped code. No re-scope, no AC change, no new mechanism.

- FO recommendation (as presented): RECONFIRM. The approved design is unchanged and validated — AC-1 verified against both observed shapes, and the coalescing question that could have invalidated the whole approach was settled negatively with evidence (no sanctioned FO fold writes terminal state, so no correct ceremony refuses). The defect is a narrow, contained bug inside one unexported helper, and the fix preserves both ACs.
- Recorded against reconfirm: the ideation's spike evidence did not reproduce against the shipped code. Confidence in "proven green against six real archived histories" was misplaced, so the correction round must re-run that evidence against the shipped read rather than re-citing the spike.

## Stage Report: implementation

- DONE: AC-1 - a passed finalize refuses when no committed stage report precedes the first terminal transition — both observed shapes refuse with nothing mutated and no archive.
  `verifyTerminalFollowsReport` in `internal/status/merge.go` (commit 1b14c09bc); `TestMergeGuardRefusesFinalizeWithNoCommittedReport` and `TestMergeGuardRefusesBatchedReportAndFlip` assert exit=1, unmutated status/verdict, no `_archive` entry for both shapes, then recovery finalizes.
- DONE: AC-2 - a healthy run is unaffected — report-only commit then flip proceeds, crash re-run proceeds, rejected verdict and non-git root stay exempt.
  `TestMergeGuardFinalizesHealthyReportThenFlip`, `TestMergeGuardFinalizesCrashRecoveryReRun`, `TestMergeGuardRejectedVerdictExemptFromReportCheck`, `TestMergeGuardNonGitRootExemptFromReportCheck` all pass; `go test ./...` and `./... -race` are green module-wide (one pre-existing, unrelated `TestCodexResolveManifestAgainstInstalledHost` failure confirmed present on unmodified main — a machine-local codex-install state check, nothing to do with this change).
- DONE: Report actual net surface against the approved estimate of +260 across up to 12 files, and name the fixture fallout count.
  Actual: +534/-7 net across 8 files (`git diff --numstat` vs `merge-base main`) — over the declared ±40% tolerance (max ~+371). `merge.go` +154/-3 (est. ~+85); new `merge_guard_report_order_test.go` +281 (est. ~+150, covers all 7 approved test-plan items). Fixture fallout: one testdata fix (`130-missing-mod-with-sentinel.md`, a stray pre-set `verdict` field that pre-empted the real terminal flip), plus a committed-report step added ahead of finalize in 1 `internal/status` test and 9 `internal/cli` tests across 3 files (`merge_test.go`, `merge_state_sync_test.go`, `terminal_consume_test.go`) — fallout the ideation flagged as firming up at implementation, outside the +260 estimate's scope (internal/cli wasn't named).

### Summary

Implemented the git-history ordering check in `finalize()` per the Proposed approach: completing stage = live status, or (already terminal) the parent-of-first-terminal-commit's status, falling back to the latest committed status when the flip itself is uncommitted. `--verdict rejected`, `--rework`, and non-git roots stay exempt. All 7 test-plan items covered in a new test file; `go test ./...` and `./... -race` are green module-wide. Net surface exceeds the declared tolerance, driven by the approved 7-item test plan plus internal/cli fallout the estimate didn't scope for — flagging for the gate Cycle line rather than trimming test coverage or leaving existing tests broken.

## Review-finding disposition

- **Reviewer (validation, 2026-08-18) — Material: `entityBlobHistory` pairs a rename-following log with a fixed-path blob read, so any entity renamed mid-history refuses permanently.** Released user and normal workflow: the First Officer running `spacedock merge guard <slug> --verdict passed` — the sole sanctioned terminal ceremony — on an entity converted flat→folder for a review room, or unarchived for repair. Observable harm: `merge.go:512` walks history with `git log --follow` (which crosses renames by design) but `merge.go:520` reads each blob as `git show <hash>:<CURRENT rel>`, which exits 128 for every pre-rename commit; `verifyTerminalFollowsReport` converts that into `merge guard: failed to read git history for {slug}`, exit 1, no terminalization, no archive, no recovery named in the message, and no bypass flag. Git history is append-only, so the entity can never be finalized through the ceremony. Affected authority: `value-ac[AC-2]` — "A healthy run is unaffected … Fails if the guard refuses a correctly-ordered run, which would block every normal merge." Both real entities below had correctly-ordered report-then-flip histories. Trigger evidence: sanctioned and observed in this repo's own state checkout. Flat→folder conversion is structurally required, since `gate record --round` refuses flat entities; commit `2c9f205d8` converted `fix-tautological-output-grep-tests` and `merge-guard-arm-not-a-stopping-point` while live, and `merge guard` later archived both (`f030e037e`, `4c43d2c1a`). Replaying the shipped read at the commit just before each archive: 8 of 17 and 7 of 26 commits unresolvable by `git show` — both merges would have hard-refused. Unarchive is recorded twice for `git-root-review-v1-materialization` (`231a35b91`, `24baede38`). One live entity is already in this state today: `stakes-declaration-read-through/index.md` (status `backlog`), unresolvable from `73170770a` back; 1 of 185 live entities scanned. Reproduced as three failing Go tests on a throwaway clone (slug rename, flat→folder, unarchive), all exiting 1 on the "failed to read git history" path. Defect kind: outcome defect. Narrow fix, contained inside the one new function and requiring no AC or surface change: resolve the path per commit instead of reusing the current one — `git log --follow --format=%H --name-only` already emits the path as of each commit (verified: it yields `stakes-declaration-read-through.md` for the pre-rename commits and `.../index.md` after).
- **Reviewer (validation, 2026-08-18) — Deferred risk: `hasStageReport` scans the whole blob, so a fenced protocol example satisfies the report check.** Released user and normal workflow: an entity whose body documents the stage-report protocol with a column-0 `## Stage Report: {stage}` line inside a fenced block. Observable harm: the guard passes vacuously from the entity's first commit — a false negative that weakens the check; it can never produce a false refusal, since a smaller `firstReport` only makes `firstTerminal <= firstReport` less likely. Affected authority: `value-ac[AC-1]` — not violated; AC-1's measuring property (zero terminal transitions with no preceding report commit) holds for every real entity today. Trigger evidence: none observed — a fence-aware scan over all 185 live entities in this state checkout found zero column-0 `## Stage Report:` headings inside fenced blocks. Promotes to material if an entity template, spec entity, or the report protocol itself starts embedding the heading at column 0; the narrow fix is to skip fenced spans when scanning, the same way the fence-safe `status --read` reader already does.
- **Reviewer (validation, 2026-08-18) — Deferred risk: a pre-terminal `verdict` or `completed` counts as terminal state, and the implementation edited a fixture rather than narrowing the detector.** Released user and normal workflow: any writer that sets `verdict` on a nonterminal entity. Observable harm: `verifyTerminalFollowsReport` treats a nonempty `verdict`/`completed` as terminal, so an entity carrying `status: implementation` + `verdict: PASSED` from its first commit refuses with "terminal state does not follow the committed report … escalate" even after a clean report commit (reproduced on the throwaway clone). This is why `testdata/merge-pr-workflow/130-missing-mod-with-sentinel.md` lost its `verdict: passed` line in `1b14c09bc`. Affected authority: `value-ac[AC-2]` — not violated; no sanctioned writer produces the shape. `merge guard` owns the terminal write, and `--set verdict=` is documented only for archiving an already-terminal entity. Trigger evidence: none observed — no live entity in this state checkout carries a nonempty `verdict` at a nonterminal status. The fixture edit does not weaken `TestMergeGuardFinalizesMissingMergeModWithSentinel`, which asserts only exit 0, `finalized`, and archived `status: done`. Promotes to material if a hand `--set verdict=` on a nonterminal entity becomes sanctioned; the narrow fix is to key terminal detection on `status == terminal` plus nonempty `completed` and drop `verdict`.

## Stage Report: validation

- DONE: Attack the refusal for false positives — build the matrix the implementation did not: folder-form entities, entities renamed or archived mid-history, feedback-cycle entities carrying "## Stage Report: {stage} (cycle N)" as the only report, an entity whose completing stage differs from its live status, and a report commit that also touches unrelated paths. Each must proceed or refuse correctly, with the reason stated.
  Matrix run as Go tests on a throwaway clone, never the implementation worktree. Correctly PROCEED: folder-form healthy; report commit also touching unrelated paths; multi-stage entity whose completing stage is its live status; live status terminal but flip uncommitted (falls back to newest committed status); `(cycle N)`-only report. Correctly REFUSE: no committed report; batched report+flip; same-prefix stage `implementation-extra`. FALSE POSITIVE, material: renamed mid-history — slug rename, flat→folder conversion, and unarchive all exit 1 on "failed to read git history", including two entities this repo really did merge (finding 1).
- DONE: Test the deliberate-coalescing case specifically. The First Officer intentionally batches some state writes. Determine whether any SANCTIONED coalesced shape now refuses — if the FO's own normal ceremony can produce a single commit carrying both the report and the flip, this guard breaks delivery and that is a material finding.
  No sanctioned coalesced shape refuses. The terminal write is exclusively `merge guard`'s (`finalize` → `emitSet` / `gates.FinalizeTerminalApproval`) and the new check runs ahead of it at `internal/status/merge.go:441-447`, before the mod-block clear; the two documented FO folds (`gate record --decision approve --consume`, `dispatch build --stamp`) write only nonterminal `status`/`started`, so a report swept into either lands strictly before the flip — proven by the multi-stage and unrelated-paths probes, which both proceed. `gate consume` on a terminal-target approval writes no status at all, so it cannot batch; the real split-root path is green in `internal/cli/terminal_consume_test.go`. The one refusing batched shape is a hand `--set status=<terminal>` carrying the report, which is run 32092321763's defect and an unclosed follow-up hole by the ideation's own record, not ceremony.
- DONE: Test the overrun claim: net +527 against a declared +260 with +/-40% tolerance (max ~371). The implementation attributes it to the approved 7-item test plan plus internal/cli fallout the estimate never scoped. Verify no scope beyond the approved design landed, and confirm the 7 test-plan items are each genuinely covered rather than padded.
  Overrun confirmed at +534/-7 = +527 net across 8 files (`git diff --numstat "$(git merge-base main HEAD)"..HEAD`), 42% past the +371 ceiling, honestly self-reported. No scope beyond the approved design landed: one refusal condition in `finalize`, three new unexported helpers, a `command-reference.md` line character-identical to the ideation's "After:" text, and test/fixture fallout — no new command grammar, flag, field, or stored format. The estimate mispredicted where, not just how much: `seedLegacyCompletedStages` already seeds reports into the `internal/status` merge fixtures, so only 1 testdata line changed instead of the budgeted 6-8 entities, while 9 `internal/cli` tests (+87) went unscoped and the test file landed at +281 against +150. All 7 test-plan items are genuinely covered, one named test each, no padding.

### Checks run

Throwaway clone of `1b14c09bc`; `go test ./...` green except `TestCodexResolveManifestAgainstInstalledHost`, which fails identically on unmodified main at the merge-base `a108559c` (machine-local codex install state). `go test -race` on `internal/status`, `internal/cli`, `internal/gates`: no data races, same single pre-existing failure. One full `internal/cli` run also reddened `TestUpgradeFromStaleMovesToGreen` and `TestFreshBoxInstallSucceeds`; both pass in isolation and a second full run was clean — flaky, not a regression. The 7 test-plan tests prove three distinct claims, each with the change that reddens it: the refusal fires on both observed shapes (`RefusesFinalizeWithNoCommittedReport`, `RefusesBatchedReportAndFlip` — red if the guard call at `merge.go:445` is dropped); healthy and exempt paths still finalize (`FinalizesHealthyReportThenFlip`, `FinalizesCrashRecoveryReRun`, `RejectedVerdictExempt`, `NonGitRootExempt` — red if the rule demands "no terminal state anywhere" instead of ordering, or if either exemption is dropped); the header match is token-exact (`ReportMatchIsTokenExactNotPrefix` — red on prefix matching). Stage names are kebab-only per `validateWorkflowStageNames`, so the `\S+` token capture cannot truncate a legal stage name. Cost measured: the busiest live entity has 71 commits, so the walk is 1 `git log` plus N `git show` spawns — sub-second, no scaling concern. Behavioral note, not a defect: on a `merge:` workflow the refusal lands at the finalize phase only — run 1 still auto-arms and writes `mod-block` (observed with the built branch binary against a copy of `merge-pr-workflow`), so the loud stop happens after the merge hook has run. That matches the design intent and leaves AC-1 intact, but "nothing mutated" is scoped to the finalize call, not the whole ceremony.

### Acceptance criteria

- **AC-1** VERIFIED. Both observed shapes refuse with exit 1, status unmutated, no `_archive` entry, and the uncommitted-report shape recovers: commit the report path-scoped, re-run, finalize and archive with `status: done`. Cited evidence reproduces exactly.
- **AC-2** NOT SATISFIED. The guard refuses correctly-ordered runs whose history crosses a rename. The ideation's supporting evidence does not reproduce against the shipped code: all six archived entities the Spike record cites as PASS hard-error on `git show` under the shipped read, so that spike exercised a different rule than the one that landed.

### Summary

Ran the false-positive matrix the assignment asked for on a detached throwaway clone. Nine adjacent shapes behave correctly, but three do not, all from one root cause: `entityBlobHistory` follows renames in the log and then reads every blob at the current path, so a renamed, converted, or unarchived entity dies on `git show` and can never be finalized. The trigger is sanctioned and observed — `gate record --round` requires folder form, and this repo really did convert two live entities and then merge them both. Settled the coalescing question negatively with evidence: no sanctioned FO fold writes terminal state, so no correct ceremony refuses. Overrun and test-plan coverage check out: the surface is 42% past tolerance but entirely within the approved design, and all 7 items are real tests with real falsifying changes.

**Recommendation: REJECTED** — one material finding (the rename-crossing history read), with a narrow, AC-preserving fix inside `entityBlobHistory`. Two deferred risks recorded with their promote-to-material conditions.

## Stage Report: implementation (cycle 2)

- DONE: Fix the material finding — `entityBlobHistory` read every historical commit at today's live path instead of the path each commit actually had.
  `entityPathsAtCommits` (commit `27c69c0ba`) replays git's own `--follow` rename-retargeting by parsing `--name-status` R-lines, so each `git show` reads the commit's real path. No new mechanism, no AC change — only `entityBlobHistory`'s internals changed.
- DONE: Re-run the ideation spike against the shipped read (not just re-cite it) and report each of the six named entities by name.
  All six now PASS on the fixed code: `pre0-cut-idempotent-on-rerun`, `remove-tautological-workflow-tests`, `collapse-duplicate-edge-marketplace-routes`, `prerelease-ships-stable-stamped-default-artifact`, `run-rejection-journey-in-team-mode`, `red-auto-continue-gate-bypass` (all `_archive/{slug}.md`, terminal stage `done`). Confirmed the same six all hard-error with "exit status 128" on the pre-fix code (`1b14c09bc`), matching validation's finding exactly — verified via a throwaway test run against this machine's real state checkout both before and after the fix, not committed (machine-path-dependent).
- DONE: Add regression coverage for the three rename shapes so this false positive cannot come back.
  `TestMergeGuardFinalizesAfterSlugRename`, `TestMergeGuardFinalizesAfterFlatToFolderConversion`, `TestMergeGuardFinalizesAfterUnarchive` — each report-then-rename-then-finalize, and each verified red (exact `failed to read git history` / exit 128 symptom) against `1b14c09bc` before the fix, green after.
- DONE: Report cycle-2 surface separately from the cumulative total.
  Cycle-2: +158/-11, net +147, across the same 2 files already touched (`merge.go`, `merge_guard_report_order_test.go`) — no new files. Cumulative vs `main`: +681/-7, net +674, across 8 files (`git diff --numstat "$(git merge-base main HEAD)"..HEAD`) — still past the approved +260 ±40% tolerance; unchanged file count (8, within "up to 12").

### Summary

Fixed the one material finding narrowly inside `entityBlobHistory`/the new `entityPathsAtCommits` helper: each commit is now read at the path it actually had, resolved by replaying git's own `--follow` rename chain from `--name-status` output, rather than assuming today's path throughout. Re-ran the ideation's six-entity spike against the shipped read as required — confirmed it hard-errors pre-fix and passes post-fix for every named entity, so the fix is evidenced against production history, not just synthetic fixtures. Added one regression test per rename shape (slug rename, flat→folder, unarchive), each proven red-then-green against the actual defect commit. AC-1 and AC-2 untouched; design not reopened. Cycle-2 surface (+147 net, 2 files) stayed tight; cumulative overrun is unchanged in kind from cycle 1 and reported for the gate's Cycle line.
