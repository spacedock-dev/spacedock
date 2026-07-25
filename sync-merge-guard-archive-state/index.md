---
id: rdjjq9hbv86skkw12z106z6q
title: Make merge-guard archive finalization durable across split-root hosts
status: validation
source: "Roborev 2146 during durable-decisions 6y final implementation review, 2026-07-24"
started: 2026-07-24T15:38:16Z
completed:
verdict:
score: "1.0"
worktree: .worktrees/spacedock-ensign-sync-merge-guard-archive-state
issue:
sprint: durable-decisions
gates:
    version: 1
    current:
        gate: gate:docs-dev:rd:ideation
    records:
        - id: gate:docs-dev:rd:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:rd-backlog-1
              briefing:
                id: briefing:docs-dev:rd:backlog:attempt-1:revision-1
                digest: sha256:a6109cdc4013c201d8a268c32dc6795d354718cfdeae85ba0fd84546ed720659
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:rd:backlog:1
                briefing: briefing:docs-dev:rd:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-24T15:38:08.065814Z"
                decision: approve
                reason: The split-root remote can retain an active merge sentinel after local archive finalization; a restart-safe supported publication path is a material durable-decisions prerelease requirement.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:docs-dev:rd:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:rd-ideation-1
              briefing:
                id: briefing:docs-dev:rd:ideation:attempt-1:revision-1
                digest: sha256:3737281543bf5be0e7b815157905b792b8adaf1ca4a756da8b531b8ae5a9c8af
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:rd:ideation:1
                briefing: briefing:docs-dev:rd:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-24T15:58:44.401791Z"
                decision: approve
                reason: The design reuses one state publisher, preserves archived read-only semantics, closes crash and duplicate-shape gaps with real-Git evidence, and passed independent staff review; implementation remains pending until 6y lands.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
---

Make a successful split-root `merge guard` finalization durable on the configured state remote, including restart after interruption, so another First Officer cannot resurrect an already archived task from the last pushed sentinel.

## Problem

`merge guard` terminalizes and commits the path-scoped archive move in the state checkout, but it does not publish that commit. The preceding `pr=pr-merge:<N>` state commit is pushed, so the remote can remain at an active merge sentinel while the local checkout alone contains the terminal archive. A fresh host then sees stale active state. Worse, after a restart the local host cannot name the active task: mutation resolution correctly reports the archived entity read-only, and `state commit` currently searches only the active/tracked source path.

This is a supported split-root durability defect exposed by 6y's terminal-gate resumption, not a reason to expand 6y's no-CLI/nine-file implementation. It is a durable-decisions prerelease blocker.

## Spike: exact archive/publish seam

On 2026-07-24, a throwaway bare origin plus two independent clones exercised the current checkout binary. The state branch began with active `task.md` carrying `pr: pr-merge:99`; both hosts used linked `.spacedock-state` worktrees on `spacedock-state/dev`.

- Origin began at sentinel commit `90610bc1ba6813cbc455ccc5a69ece1c90ffcb72`.
- Host A ran `spacedock merge guard task --workflow-dir <host-a>/docs/dev --verdict passed`; it exited 0, moved the task to `_archive/task.md`, and created local commit `016acfa2ccf419065c12b1f4b0b420625346aa46`.
- Host A was clean and exactly one commit ahead, but origin remained at `90610bc1…`; origin still contained active `task.md` and no archive.
- A rerun of `merge guard task` exited 1 with `archived entity is read-only`. `state commit task` exited 1 with `no entity "task"`. `state ready` exited 0 on both hosts but, as a pull-only boot gate, did not publish A's ahead commit.
- Host B still rendered `task` as active at `status: implementation`.

The cheapest existing public boundary is therefore `state commit <slug>`, not `state ready` and not a new verb. It already denotes one entity's commit-and-publish unit and owns the required push/rebase/HALT discipline; it needs archived-path resolution and clean-but-ahead publication. The spike fixture becomes the first implementation test.

## Proposed approach

1. Extract the post-commit publisher and rebase inspection from `internal/cli/state_sync.go` into a standard-library-only `internal/statesync` package. It has no CLI or status dependency, so both `internal/cli` and `internal/status` import it without a cycle. It owns ordinary push, non-fast-forward `pull --rebase` and re-push, conflict discovery plus `rebase --abort`, peer-commit evidence, and the no-origin result. Both callers receive one typed outcome; neither may force-push, auto-resolve, or autostash.
   - Serves AC-1, AC-2, and AC-3.
   - Simplest alternative considered: duplicate the Git sequence in `merge.go`. That creates a second policy whose conflict and no-origin behavior can drift, so it is insufficient.
2. Immediately after each caller resolves its checkout/roots, and before entity resolution, staging, or mutation, call a shared preflight. If a rebase is already in progress, the preflight captures unmerged paths and the peer commit first, aborts the rebase, and returns the existing exit-3 HALT outcome. `merge guard` must not terminalize an entity and `state commit` must not stage it before this check.
   - Serves AC-3 by making restart from the exact interrupted rebase state deterministic.
   - Simplest alternative considered: let the publisher notice only when it starts a new rebase. That misses a process restart with Git already mid-rebase and allows mutation against an unmerged index.
3. After `commitArchiveMove` succeeds, split-root `merge guard` calls the publisher before emitting its final signal. The archive commit remains the durable local boundary: a publication/network failure does not roll it back, and the error names `spacedock state commit <slug> --workflow-dir <dir>` as the supported restart. Inline workflows do not publish the code branch.
   - Serves AC-1, AC-2, and AC-5.
   - Simplest alternative considered: add `git push` to the FO contract. That is not binary-owned, cannot enforce rebase/HALT behavior, and leaves restart unresolved.
4. Make `state commit <slug>` resolution return both the canonical path set and `active|archived` scope. Before any movement, fail closed if the slug exists in both active and archive scope or if archive scope contains both flat and folder forms; “active wins” is forbidden for these invalid shapes. Active scope retains today's path-scoped stage/commit behavior. Archived scope is publish-only: require its active-source and archived-destination pathspecs to be clean in both index and worktree, never run `git add` or `git commit`, and publish only existing ahead history. A dirty archived task refuses with HEAD and origin unchanged; a clean, fully published archive is `no-op`.
   - Serves AC-2, AC-4, and AC-6.
   - Simplest alternative considered: rerun `merge guard`. Mutation resolution intentionally refuses archived entities, so widening that mutation surface would weaken read-only archive semantics.
5. Final output binds lifecycle and durability: `merge guard --json` carries exactly one JSON value followed by EOF, whose finalized result says `pushed`, `local-only`, or `inline`; prose says the equivalent. An exit-3 conflict uses the existing HALT class and does not print a remotely durable success. No new public command or flag is added.
   - Serves AC-1, AC-3, and AC-5.
   - Simplest alternative considered: emit the current finalize object followed by a separate state-sync object. Two JSON documents are not one valid command result and could falsely announce success before publication.

## Scope boundaries

- Preserve `commitArchiveMove` and its exact live/archive path-scoped staging.
- Preserve archived entities as read-only to `status` and `merge guard`; `state commit` gains only clean, publish-only archived resolution and never commits archived dirt.
- Reuse one state publisher for ordinary push, non-fast-forward integration, conflict abort/HALT, and no-origin behavior.
- Fail closed before movement on active/archive identity collisions and archived flat/folder shape collisions.
- Never force-push, auto-resolve, autostash, or stage unrelated state.
- Do not change gate authority, recorder/application schemas, presentation providers, PR-host behavior, worktree teardown, or 6y's lifecycle mechanics.
- Do not add a public verb or require FO-authored raw Git.

## Deferred trigger: sibling dirt plus non-fast-forward

AC-4 promises immediate publication with unrelated sibling dirt only when the ordinary push succeeds. If that push is rejected and Git refuses `pull --rebase` because unrelated tracked sibling dirt remains, the supported behavior is a non-HALT exit 1: do not autostash, do not move either ref, retain the local archive commit, and allow `state commit <archived-slug>` to resume after the sibling dirt is settled. The real-Git suite characterizes that recoverability, but this task does not promise publication through unrelated dirt. Promote the case only if immediate publication under both conditions becomes a product promise.

## Expected surface

Baseline: 10 files and about 560 changed lines (`+480/-80`), dominated by real-Git fixtures.

- `internal/statesync/publish.go` (new, about 155 LOC): standard-library-only publish/rebase/preflight/HALT engine and typed outcome.
- `internal/cli/state_sync.go` (about 85 changed LOC): scoped resolver, corruption guards, publish-only archived resume, and shared rendering.
- `internal/status/merge.go` (about 50 changed LOC): preflight, publish after archive commit, and one durability-bound finalize result.
- `internal/cli/merge_state_sync_test.go` (new, about 185 LOC): two-host value, interruption, peer, pre-existing rebase, isolation, and local-only fixtures.
- `internal/cli/state_commit_test.go` (about 55 added/changed LOC): publish-only archive, dirty refusal, duplicate-shape refusal, and idempotent resume.
- `internal/status/merge_guard_test.go` (about 20 added/changed LOC): exact one-value-plus-EOF JSON and inline/local-only result assertions.
- `internal/cli/help.go`, `docs/site/reference/command-reference.md`, `skills/first-officer/references/fo-merge-core.md`, and `docs/dev/_mods/pr-merge.md` (about 30 changed lines total): supported outcome and restart wording.

Tolerance: the implementation should remain within 520–600 changed lines across these 10 files; at most 13 files and 900 changed lines are allowed for review-driven fixture or structural adjustments. Exceeding the upper tolerance, adding a public verb/flag, adding a second Git harness, introducing a non-standard-library `internal/statesync` dependency, or touching any excluded gate/provider/PR-host/lifecycle surface requires a return to ideation.

## Acceptance criteria

**AC-1 (VALUE) - Successful split-root finalization is visible as the same archived terminal task from a fresh host.**
Verified by: a real two-clone fixture pushes the merge sentinel, runs the public `merge guard` command on host A, runs `state ready` on host B, and observes exactly one archived terminal task at origin's archive commit with no active sentinel row. Removing the post-archive publisher makes host B reproduce the stale active task.

**AC-2 - Interruption after the local archive commit but before publication has a supported idempotent resume.**
Verified by: a pre-push failure fixture stops after the archive commit, removes the injected failure, and reruns documented `state commit <slug>` against the archived slug. Origin advances to the existing archive commit, archive-commit count remains exactly one, a second resume is a no-op, and both worktree and index are clean. Separate staged, unstaged, and untracked archived-dirt legs refuse before publication and prove both HEAD and origin byte-identical. Any archived `git add`/commit path, or reverting clean-ahead publication, makes a leg fail.

**AC-3 - Peer synchronization retains the existing conflict safety.**
Verified by: one two-host leg pushes a disjoint peer entity before host A publishes and observes a linear rebase/re-push with both entities present. The conflict leg creates a local archive-rename commit, pushes a conflicting remote edit to the active task, manually starts the real `pull --rebase` until Git is mid-conflict, then invokes the supported command. Preflight exits 3, names the captured unmerged path and peer commit, aborts the rebase, leaves the remote peer edit intact, and restores the local archive commit as recoverable HEAD. Moving preflight after entity resolution/staging, omitting abort, or force-pushing fails the fixture.

**AC-4 - Archive publication remains path-scoped and does not sweep sibling dirt.**
Verified by: on the ordinary direct-push path, a sibling tracked entity is dirtied before finalization; the archive commit's name-only tree delta is exactly `{slug}.md` plus `_archive/{slug}.md` (or the two folder roots), the sibling remains dirty and absent from the pushed commit, and archived resume creates no commit. Replacing path-scoped staging with broad add makes the fixture fail. The separately recorded sibling-dirt/non-fast-forward characterization proves recoverability after dirt settles without promising or implementing autostash.

**AC-5 - Inline and no-origin workflows report truthful local behavior.**
Verified by: every `merge guard --json` leg decodes one JSON value and requires a second decoder call to return EOF. The origin-backed leg reports `pushed` only when the remote ref equals the local archive commit; a split-root checkout with no `origin` keeps its local archive commit and reports `local-only`; an inline fixture reports `inline` and proves its code-branch remote ref never moved. A second JSON value, a mismatched remote ref, or any inline push fails.

**AC-6 - Invalid active/archive identity or archive-shape collisions cannot move local or remote state.**
Verified by: real-Git fixtures create (a) the same slug in active and archive scope and (b) both `_archive/<slug>.md` and `_archive/<slug>/index.md`. `state commit <slug>` refuses before staging or publication in each case, with HEAD, index, worktree, and origin unchanged. Restoring “active wins”, folder preference, or any push-before-validation ordering makes the test fail.

## Test plan

- Add red tests by reusing `twoHostStateWorkflow`, `runStateCommitCmd`, and `runMergeCLI`; do not create a second Git harness. The existing real bare-origin/two-clone substrate serves AC-1 through AC-6. A mock or second bespoke harness was simpler locally, but would either miss real refs/indexes or duplicate fixture semantics.
- Drive `run(...)` for public CLI behavior. The repository-local failing pre-push hook serves AC-2 by leaving a real archive commit at the exact network boundary. An injected Go callback was simpler, but it would not prove the supported binary survives a real Git publication failure.
- Add table legs for archived staged/unstaged/untracked dirt and both duplicate shapes, asserting zero HEAD/origin movement. These serve AC-2 and AC-6; checking only stderr was simpler but would not prove a refusal occurred before staging or push.
- Create the AC-3 pre-existing conflict with real Git: local archive rename, remote active-task edit, explicit `pull --rebase` to the conflict, then invoke the binary. A synthetic rebase marker was simpler but cannot prove path capture, abort restoration, or peer preservation.
- Characterize sibling dirt plus non-fast-forward as exit 1 with no autostash and a recoverable archive commit; settle the dirt and prove the same archived-slug command then publishes. This guards the deferred boundary without promoting immediate publication to an AC.
- Decode JSON with `json.Decoder`, assert the first value's result against refs, and require the second decode to return `io.EOF`. Substring checks were simpler but accept concatenated JSON and false durability claims.
- Run `gofmt -w ./cmd ./internal`, focused `go test ./internal/statesync ./internal/cli ./internal/status ./internal/contractlint`, then `go test ./...` and `go test ./... -race`.
- Because the diff changes both a status mutation boundary and host-neutral FO contract text, run the detached adversarial audit plus every required host live lane from the diff-to-lane policy. The audit must delete/bypass the post-archive publish call and confirm the two-host value test turns red.

## Documentation diff

- `spacedock merge --help`: change “terminalize and archive” to “terminalize, archive with a path-scoped commit, and publish split-root state”; add “after an interrupted publication, `spacedock state commit <slug>` resumes from the archived slug without creating another archive commit.”
- Command reference `state commit`: change “one canonical top-level entity slug” to “one canonical active or clean archived top-level entity slug”; state that archived scope is publish-only, dirty/colliding archive shapes refuse before movement, a clean unpublished commit is pushed, and fully published state remains a no-op.
- FO merge core `effect`: change “including the path-scoped archive commit” to “including the path-scoped archive commit and split-root publication”; define done as remote-durable when an origin exists and explicitly local-only otherwise.
- Development `pr-merge` mod: change each “terminalizes, archives, and commits” description to “terminalizes, archives, commits path-scoped, and publishes through the shared state-sync discipline.” No raw Git instruction is added.

## Stage Report: ideation

- DONE: Reproduce the pushed-sentinel/local-archive split with a real two-clone fixture, including interruption after archive commit and the current inability to resume by active slug.
  Origin stayed at `90610bc1…` while host A was clean at archive `016acfa2…`; guard, state-commit, and state-ready resume attempts left host B on the active sentinel.
- DONE: Design the smallest binary-owned publication/resume path by reusing the existing state-sync conflict discipline, with exact files/LOC and no raw-Git FO workaround.
  The design reuses `state commit <slug>`, extracts one internal publisher, and declares a 9-file/~520-line baseline with a 12-file/900-line reset threshold.
- DONE: Define falsifiable two-host, interruption, non-fast-forward/conflict, sibling-isolation, and local-only evidence that measures durable finalization.
  AC-1; AC-2; AC-3; AC-4; and AC-5 name real-Git observations and the specific removed or weakened mechanism that turns each test red.

### Summary

The spike proved the remote-resurrection defect and that no current entry point publishes the already-committed archive. The design adds no public verb: `merge guard` uses the extracted state publisher, while `state commit <slug>` becomes the idempotent archived-slug restart path under the existing conflict discipline.

## Stage Report: ideation (cycle 2)

- DONE: Reproduce the pushed-sentinel/local-archive split with a real two-clone fixture, including interruption after archive commit and the current inability to resume by active slug.
  The accepted spike remains unchanged; cycle 2 adds the real mid-rebase restart shape that begins with the same local archive commit and remote active edit.
- DONE: Design the smallest binary-owned publication/resume path by reusing the existing state-sync conflict discipline, with exact files/LOC and no raw-Git FO workaround.
  The design keeps `state commit <slug>`, makes archive scope publish-only, adds pre-mutation corruption/rebase guards, and resets the baseline to 10 files/~560 LOC.
- DONE: Define falsifiable two-host, interruption, non-fast-forward/conflict, sibling-isolation, and local-only evidence that measures durable finalization.
  AC-1; AC-2; AC-3; AC-4; AC-5; and AC-6 now bind refs, HEAD/index/worktree state, one-value JSON EOF, duplicate refusals, and the no-autostash deferred boundary.

### Summary

Cycle 2 preserves the accepted no-new-verb seam while closing the unsafe archived-commit and restart gaps. Archived state is clean-and-publish-only, invalid identities halt before movement, and shared preflight handles a rebase already in progress before either caller can resolve or mutate an entity.

### Feedback Cycles

- Cycle 1 (Roborev 2304): CHANGES REQUESTED — invalid-branch ref safety, continuation wording, verified no-op, HALT evidence, and JSON behavior were corrected. Actual surface was 10 files/~899 LOC versus 10/~560 estimated; acceptance criteria were unchanged.
- Cycle 2 (Roborev 2323): CHANGES REQUESTED — success-output compatibility, archive-folder dirt coverage, and diagnostics were corrected. Actual surface was 11 files/~1,006 LOC; the accepted dirty-sibling/non-fast-forward boundary remained no-autostash/recover-after-settlement.
- Cycle 3 (Roborev 2328): CHANGES REQUESTED — exact checkout/branch preflight, clean active retry, and post-archive cleanup guidance were corrected. Actual surface was 11 files/~1,131 LOC; acceptance criteria were unchanged.
- Cycle 4 (Roborev 2347): CHANGES REQUESTED — wrong/unknown-branch rebases are now left untouched and abort failures report possible residual conflict. Same configured-branch abort was retained because AC-3 requires archive-HEAD restoration.
- Cycle 5 (Roborev 2367): CHANGES REQUESTED — abort diagnostics, neutral ready errors, clean peer-only narration, finalized HALT evidence, and stale comments were corrected. Archived dirt remains refusal-only because AC-2 forbids committing it.
- Cycle 6 (Roborev 2383): CHANGES REQUESTED — remote-ref no-op verification, root/branch remedies, action-specific recovery prose, archived Git-error classification, exit documentation, and FO recovery contract assertions were corrected.
- Cycle 7 (Roborev 2395): CHANGES REQUESTED — the ineffective negative assertion, HALT-renderer comment, and archived `-m` documentation were corrected. Active/archive collision refusal and same-branch rebase abort were retained as explicit AC-6 and AC-3 requirements.
- Design-reset decision: no product-design reset. Final surface is 15 files/1,443 changed lines versus 10/~560 estimated (150% of files, 258% of LOC); the excess is real-Git counterexample coverage and review-driven safety/diagnostics, with no new verb, flag, dependency, raw-Git FO step, lifecycle authority, or second publisher.

## Stage Report: implementation

- DONE: Split-root `merge guard` now publishes the exact archive commit before reporting remote durability, while inline and no-origin workflows remain truthful.
  Evidence: code commit `3d034e07`; `TestMergeGuardPublishesArchiveVisibleToFreshHost`, `TestMergeGuardNoOriginReportsLocalOnly`, and `TestMergeGuardInlineJSONReportsOneDurabilityBoundResult`.
  Falsifier: removing `publishMergeArchive` leaves the fresh host on the active sentinel and makes the origin-ref equality assertion fail.
- DONE: A clean archived slug resumes an interrupted publish without staging archive dirt or creating a second archive commit.
  Evidence: `TestArchivedStateCommitResumesInterruptedPublication` proves one archive commit, origin equality, clean index/worktree, and a second no-op.
  Evidence: `TestStateCommitRefusesDirtyArchivedEntityBeforePublication` covers staged, unstaged, and untracked folder dirt with unchanged HEAD and origin.
  Falsifier: adding an archived `git add`/commit path makes the dirty refusal and one-archive-commit assertions fail.
- DONE: The shared standard-library publisher preserves multi-host linear history, conflict evidence, and local recovery without force-push, autostash, or unrelated staging.
  Evidence: `TestMergeGuardAndStateCommitPreflightInterruptedArchiveRebase`, `TestMergeGuardSiblingDirtNonFastForwardRemainsRecoverable`, and `TestStateCommitMultiWriterHappyPath`.
  Evidence: `TestPublishVerifiesNoOpWithoutRemoteTrackingRef`, wrong-root/branch preflight tests, and clean retry/peer-only narration tests cover review-discovered boundaries.
- DONE: Invalid identity shapes fail closed before mutation or publication.
  Evidence: `TestStateCommitRefusesActiveArchiveAndArchiveShapeCollisions` proves active/archive and archived flat/folder collisions preserve HEAD, index, worktree, and origin.
  Review disposition: “active wins” was declined because it directly invalidates AC-6's same-slug counterexample.
- DONE: Rebase ownership behavior matches the accepted restart contract without touching unrelated operator rebases.
  Evidence: a rebase whose `head-name` is the configured state branch is captured and aborted to restore the recoverable archive HEAD required by AC-3.
  Evidence: `TestPreflightLeavesWrongBranchRebaseUntouched` proves a different/unknown branch is refused without abort; an abort failure returns `failed` and warns that conflict may remain.
  Review disposition: preserving every same-branch manual rebase was declined because removing the abort fails AC-3's required HEAD restoration.
- DONE: Documentation and FO recovery text bind cleanup to durable publication.
  Evidence: merge help, command reference, PR-merge mod, and `fo-merge-core.md` name archived `state commit` resume, exit-3 HALT, exit-1 recovery, local-only behavior, and publish-only archived `-m`.
  Evidence: `TestFOFunctionRequiredCallSites` pins both recovery branches in the declarative FO contract.
- DONE: Required verification is green on commit `3d034e07`.
  Evidence: `gofmt -w ./cmd ./internal`; `git diff --check`; focused statesync/CLI/status/ensigncycle/contractlint packages; `go test ./...`; and `go test ./... -race`.
  Review evidence: Roborev jobs 2304, 2323, 2328, 2347, 2367, 2383, and 2395 were triaged; demonstrated defects were fixed, while only the two acceptance-contract reversals above were declined.

### Summary

Merge-guard finalization is now remote-durable across split-root hosts, and an interrupted archive publication has one supported idempotent restart through the existing `state commit <slug>` command. The implementation keeps archived scope publish-only, preserves conflict and sibling-dirt recovery, reports durability in one result, and introduces no new public command surface.
