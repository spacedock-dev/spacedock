---
id: rdjjq9hbv86skkw12z106z6q
title: Make merge-guard archive finalization durable across split-root hosts
status: ideation
source: "Roborev 2146 during durable-decisions 6y final implementation review, 2026-07-24"
started: 2026-07-24T15:38:16Z
completed:
verdict:
score: "1.0"
worktree:
issue:
sprint: durable-decisions
gates:
    version: 1
    current:
        gate: gate:docs-dev:rd:backlog
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

1. Extract the post-commit publisher from `internal/cli/state_sync.go` into a small internal state-sync boundary with no CLI or status dependency. It owns ordinary push, non-fast-forward `pull --rebase` and re-push, same-task conflict discovery plus `rebase --abort`, peer-commit evidence, and the no-origin result. Both callers receive one typed outcome; neither may force-push or auto-resolve.
   - Serves AC-1, AC-2, and AC-3.
   - Simplest alternative considered: duplicate the Git sequence in `merge.go`. That creates a second policy whose conflict and no-origin behavior can drift, so it is insufficient.
2. After `commitArchiveMove` succeeds, split-root `merge guard` calls that publisher before emitting its final signal. The archive commit remains the durable local boundary: a publication/network failure does not roll it back, and the error names `spacedock state commit <slug> --workflow-dir <dir>` as the supported restart. Inline workflows do not publish the code branch.
   - Serves AC-1, AC-2, and AC-5.
   - Simplest alternative considered: add `git push` to the FO contract. That is not binary-owned, cannot enforce rebase/HALT behavior, and leaves restart unresolved.
3. Extend `state commit <slug>` resolution so active form still wins, then archived flat/folder form is accepted as the same entity commit unit. If the entity has no new dirt but local HEAD is ahead of the configured state branch's remote-tracking ref, the command runs the extracted publisher instead of returning the current clean no-op. It creates no second archive commit. A genuinely clean, fully published entity remains `no-op`.
   - Serves AC-2 and AC-4.
   - Simplest alternative considered: rerun `merge guard`. Mutation resolution intentionally refuses archived entities, so widening that mutation surface would weaken read-only archive semantics.
4. Final output binds lifecycle and durability: `merge guard --json` carries one valid object whose finalized result says `pushed`, `local-only`, or `inline`; prose says the equivalent. An exit-3 conflict uses the existing HALT class and does not print a remotely durable success. No new public command or flag is added.
   - Serves AC-1, AC-3, and AC-5.
   - Simplest alternative considered: emit the current finalize object followed by a separate state-sync object. Two JSON documents are not one valid command result and could falsely announce success before publication.

## Scope boundaries

- Preserve `commitArchiveMove` and its exact live/archive path-scoped staging.
- Preserve archived entities as read-only to `status` and `merge guard`; only `state commit` gains archived resolution for publication.
- Reuse one state publisher for ordinary push, non-fast-forward integration, conflict abort/HALT, and no-origin behavior.
- Never force-push, auto-resolve, or stage unrelated state.
- Do not change gate authority, recorder/application schemas, presentation providers, PR-host behavior, worktree teardown, or 6y's lifecycle mechanics.
- Do not add a public verb or require FO-authored raw Git.

## Expected surface

Baseline: 9 files and about 520 changed lines (`+450/-70`), dominated by real-Git fixtures.

- `internal/statesync/publish.go` (new, about 140 LOC): extracted publish/rebase/HALT engine and typed outcome.
- `internal/cli/state_sync.go` (about 70 changed LOC): archived resolution, clean-but-ahead resume, and rendering through the shared publisher.
- `internal/status/merge.go` (about 45 changed LOC): publish after archive commit and bind durability into the single finalize result.
- `internal/cli/merge_state_sync_test.go` (new, about 230 LOC): two-host value, interruption, peer, conflict, isolation, and local-only fixtures.
- `internal/cli/state_commit_test.go` (about 55 added/changed LOC): archived flat/folder resolution and idempotent clean-ahead resume.
- `internal/cli/help.go`, `docs/site/reference/command-reference.md`, `skills/first-officer/references/fo-merge-core.md`, and `docs/dev/_mods/pr-merge.md` (about 30 changed lines total): supported outcome and restart wording.

Tolerance: at most 12 files and 900 changed lines. Exceeding either bound, adding a public verb/flag, or touching any excluded gate/provider/PR-host/lifecycle surface requires a return to ideation before implementation continues.

## Acceptance criteria

**AC-1 (VALUE) - Successful split-root finalization is visible as the same archived terminal task from a fresh host.**
Verified by: a real two-clone fixture pushes the merge sentinel, runs the public `merge guard` command on host A, runs `state ready` on host B, and observes exactly one archived terminal task at origin's archive commit with no active sentinel row. Removing the post-archive publisher makes host B reproduce the stale active task.

**AC-2 - Interruption after the local archive commit but before publication has a supported idempotent resume.**
Verified by: a pre-push failure fixture stops after the archive commit, removes the injected failure, and reruns documented `state commit <slug>` against the archived slug. Origin advances to the existing archive commit, archive-commit count remains exactly one, a second resume is a no-op, and both worktree and index are clean. Reverting archived resolution or clean-but-ahead publication makes the test fail.

**AC-3 - Peer synchronization retains the existing conflict safety.**
Verified by: one two-host leg pushes a disjoint peer entity before host A publishes and observes a linear rebase/re-push with both entities present. A second leg changes the same active sentinel remotely and observes exit 3, named conflicting paths and peer commit, an aborted rebase, no force push, the remote peer version intact, and the local archive commit still recoverable. Removing abort/no-force behavior or bypassing the shared publisher fails the corresponding assertion.

**AC-4 - Archive publication remains path-scoped and does not sweep sibling dirt.**
Verified by: a sibling tracked entity is dirtied before finalization; the archive commit's name-only tree delta is exactly `{slug}.md` plus `_archive/{slug}.md` (or the two folder roots), the sibling remains dirty and absent from the pushed commit, and the resume creates no commit. Replacing path-scoped staging with broad add makes the fixture fail.

**AC-5 - Inline and no-origin workflows report truthful local behavior.**
Verified by: an inline fixture archives without invoking state publication and reports `inline`; a split-root checkout with no `origin` keeps its local archive commit and reports `local-only`. The origin-backed fixture reports `pushed`. Claiming pushed durability without a matching remote ref, or pushing the inline code branch, fails.

## Test plan

- Add red tests first by adapting `twoHostStateWorkflow`; the real bare-origin/two-clone harness serves AC-1 through AC-4. A mocked publisher was the simpler alternative, but it cannot expose the actual ref, worktree, rebase, index, or sibling-dirt states those criteria measure.
- Drive `run(...)` for public CLI behavior. The repository-local failing pre-push hook serves AC-2 by leaving a real archive commit at the exact network boundary. An injected Go callback was simpler, but it would not prove the supported binary survives a real Git publication failure.
- Keep one focused state-sync test per branch: direct push, disjoint non-fast-forward, same-task conflict/HALT, no origin, and already-published no-op. The change that would falsify each is named in the AC above.
- Run `gofmt -w ./cmd ./internal`, focused `go test ./internal/statesync ./internal/cli ./internal/status ./internal/contractlint`, then `go test ./...` and `go test ./... -race`.
- Because the diff changes both a status mutation boundary and host-neutral FO contract text, run the detached adversarial audit plus every required host live lane from the diff-to-lane policy. The audit must delete/bypass the post-archive publish call and confirm the two-host value test turns red.

## Documentation diff

- `spacedock merge --help`: change “terminalize and archive” to “terminalize, archive with a path-scoped commit, and publish split-root state”; add “after an interrupted publication, `spacedock state commit <slug>` resumes from the archived slug without creating another archive commit.”
- Command reference `state commit`: change “one canonical top-level entity slug” to “one canonical active or archived top-level entity slug”; add that a clean-but-unpublished local commit is pushed, while fully published state remains a no-op.
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
