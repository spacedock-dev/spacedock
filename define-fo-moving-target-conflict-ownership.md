---
title: Define FO ownership and recovery for moving-target PR conflicts
status: ideation
source: "Captain follow-up after the 2026-08-03 durable-decisions conflict diagnosis."
started: 2026-08-03T16:01:22Z
completed:
verdict:
score: 0.95
worktree:
issue:
sprint: durable-decisions
group: fo-contract
milestone: 0.27.0
id: g3912c6f6jkgd0yjmyg6h7yn
gates:
    version: 1
    records:
        - id: gate:g3912c6f6jkgd0yjmyg6h7yn:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:g3912c6f6jkgd0yjmyg6h7yn-backlog-1
              briefing:
                id: briefing:g3912c6f6jkgd0yjmyg6h7yn:backlog:attempt-1:revision-1
                digest: sha256:fbe944ff1e7905c23765219faa581bc2c1cba4980f68664253038adacb2afafa
                request-digest: sha256:7c706e3f3ab394a0653b087e1007f2aab9246a08c95e8c8947bb976a60d2736a
                room-ref: ./define-fo-moving-target-conflict-ownership/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:g3912c6f6jkgd0yjmyg6h7yn:backlog:1
                briefing: briefing:g3912c6f6jkgd0yjmyg6h7yn:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-03T15:59:40.23286Z"
                decision: approve
                reason: Captain conn authorizes approval when SO concurs. SO concurs with parallel EJ/G3 ideation; keep EJ-before-G3 implementation landing because they share the FO dispatch and merge contract.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:g3912c6f6jkgd0yjmyg6h7yn:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:g3912c6f6jkgd0yjmyg6h7yn-ideation-1
              briefing:
                id: briefing:g3912c6f6jkgd0yjmyg6h7yn:ideation:attempt-1:revision-1
                digest: sha256:4d10f045deb6106fb8b30200e4e726677c24b73aa534c4cdffe7a86d3fba85ce
                request-digest: sha256:4459edfb7e1c00dca121c198716a0b3fa42dd76a97fa07d5d8479ad75ff9e69e
                room-ref: ./define-fo-moving-target-conflict-ownership/review/ideation/briefing-1
---

Define the ownership, durable hold, evidence, and keep-moving behavior when a pending PR meets a moving merge target. The task closes the contract ambiguity without adding a resolver worker or changing the workflow definition.

## Problem

The First Officer contract correctly forbids automatic conflict resolution, force-push, and discarding either side, but it does not name the owner or the durable transition when an open PR becomes dirty because the merge target moved. A First Officer can therefore mistake a filtered empty `status --next` for a session stop, ask the wrong worker to resolve a two-writer conflict, clear a pending terminal approval, or let stale validation evidence authorize a new head. The current PR and state must remain truthful while unrelated sprint work continues.

## Proposed approach

Make ownership explicit. The code/worktree owner may rebase and commit only within its owned branch when the Captain authorizes that ordinary branch maintenance; the Captain owns semantic reconciliation of two-writer conflicts and decides whether to rework, rebase, or close/reopen. The First Officer owns transport and state safety: abort any conflicted local rebase, surface the PR number, old and new base/head commits, and conflict paths, and preserve the affected task's current status, bare `pr:#N`, `mod-block: merge:pr-merge`, and pending terminal approval. No conflict path clears or consumes authority, terminalizes, auto-resolves, force-pushes, or reuses old validation evidence.

After the Captain or owning branch author produces a new exact head, the First Officer records the new head as the delivery candidate, reruns the required CI and fresh validation/gate evidence, and only then resumes the merge ceremony. The event loop holds only the affected task and keeps independently ready tasks dispatching or presenting. The contract must name this as a per-entity conflict hold, not a generic resolver route.

## Out of scope

- Automatic conflict classification, conflict resolution, `git checkout --ours/--theirs`, force-push, or force-finalization.
- A standing resolver worker, new conflict database, new frontmatter field, or workflow-specific finding taxonomy.
- Changes to the commissioned workflow definition, PR provider, branch protection, or GitHub merge policy.
- Reusing a superseded approval, old-head CI, or prior gate Briefing as proof for a rewritten candidate.

## Acceptance criteria

**AC-1 (VALUE) — A moving-target conflict cannot finalize or lose a pending decision.**
Verified by: a local Git/PR fixture that changes the merge target after an approval is pending; the observed state keeps the affected task at its current stage with `pr:#N`, `mod-block: merge:pr-merge`, and the pending terminal application, while the command log records no terminalization, consume, worker dispatch, force operation, or auto-resolution and surfaces exact conflict paths plus peer/base/head evidence.

**AC-2 — Ownership and the required transition are unambiguous.**
Verified by: a behavior fixture that exercises an FO-held conflict, Captain-authorized branch rebase, and new-head handoff; the trace attributes semantic resolution to the Captain/owning branch and transport/state guarding to the FO, aborts failed rebases, and requires a fresh exact-head validation before merge guard can finalize.

**AC-3 — Stale evidence cannot authorize a rewritten head.**
Verified by: an exact-head fixture that changes the PR head after a green validation result and asserts that the old CI/gate evidence is rejected or held until the new head's required lanes and fresh gate evidence pass.

**AC-4 — Unrelated sprint work keeps moving.**
Verified by: a multi-entity event-loop fixture with one conflict-held PR and two independent ready entities; the conflicted task remains held with its evidence while both unrelated entities reach their next dispatch or gate action in the same loop sequence.

**AC-5 — The generic contract does not create a resolver escape hatch.**
Verified by: a detached adversarial audit that attempts auto-resolution, `-X ours`/`-X theirs`, force-push, direct frontmatter clearing, and a generic resolver spawn; each attempt fails or remains outside the accepted trace, while the supported Captain-owned reconciliation path remains executable.

## Expected surface and semantic boundaries

Expected change is limited to the generic First Officer conflict/merge/dispatch contract and its behavior fixtures: `skills/first-officer/references/first-officer-shared-core.md`, `skills/first-officer/references/fo-merge-core.md`, `skills/first-officer/references/fo-dispatch-core.md`, relevant host adapter notes, and focused conflict/event-loop tests. Estimate 4–9 files and +150/-30 lines, with a 2x tolerance. Allowed semantic change: owner attribution, evidence, and per-entity hold/re-entry ordering. Stored state schema, workflow README, PR provider behavior, and merge authority remain unchanged.

## Test plan

Start with a throwaway two-branch fixture whose base advances after an approved PR is recorded. Capture the failed rebase paths and peer commit, abort cleanly, and assert byte-preserved state. Add command-log tests for the Captain/branch-owner handoff, exact-head freshness, and keep-moving matrix. Run focused status/dispatch/merge tests, a detached adversarial audit on the contract/guard surface, then `go test ./...`, `go test ./... -race`, formatting, and all required host lanes if the runtime adapters change.

## Stage-specific test gates

- Ideation must replay one moving-target conflict and record the exact pre/post state and ownership handoff before selecting wording or a helper.
- Implementation must preserve pending authority bytes, prove no auto-resolution or force path, and add no resolver worker or state field.
- Validation must reproduce the conflict/abort/new-head sequence, stale-evidence rejection, and unrelated-task continuation, plus focused/full/race/format and the required detached audit/live evidence.

## Stage Report: ideation

- DONE: Replay one moving-target PR conflict and record the exact pre/post state, conflict path, and Captain/branch-owner/FO handoff.
  Evidence: a throwaway Git/PR fixture at `/tmp/spacedock-g3-conflict.oWmCsG` used `main` as the merge target and `pr-7` as the open PR. Before the target moved, the approved-but-not-consumed state was `status: validation`, `pr: '#7'`, `mod-block: merge:pr-merge`, and `terminal-application: pending`; its exact state-file SHA-256 was `67ed11855aacfbd6e6254ed9c6629c7bdb33398516606a822402b169ca05da48`. The old base/merge-base was `1ee5119d9828811d7ce26adac08d020bc3a65640`; the old PR head was `7918cba8d1f765a1805efd7e016366b73cd2b9f3`. After `main` advanced to base `0e97600e663d56beb807067e2b175ff6f739e6f1`, `git rebase main` exited 1 with `UU contract.txt`, `REBASE_HEAD=7918cba8d1f765a1805efd7e016366b73cd2b9f3`, and conflict path `contract.txt` (the peer/base commit was `0e97600…`). The FO ran `git rebase --abort`; afterward `pr-7` was back at `7918cba…`, `main` remained `0e97600…`, the worktree was clean, no unmerged paths remained, and the state-file SHA-256 was unchanged. No terminalization, gate consume, worker dispatch, force operation, or automatic resolver command occurred. The Captain then authorized ordinary branch maintenance; the owning branch author manually reconciled `contract.txt`, producing exact new head `e54865db2dff623ab9d7bfe7ba069ec2aae6226c` on base `0e97600…` (ahead/behind `0/1`).

- DONE: Reconcile pending authority, exact-head freshness, keep-moving, and no-resolver constraints against the existing merge and dispatch contract, and name the remaining risk.
  Evidence: `fo-merge-core.md` makes the terminal-target application remain `pending` after `gate consume`; only `merge guard` may consume it after proven delivery, while retryable delivery trouble leaves the same approval pending. That preserves `pr:#7`, `mod-block: merge:pr-merge`, and the current stage during the conflict. `first-officer-shared-core.md` already requires aborting a conflicted rebase and surfacing paths/peer evidence, with no force-push, `-X ours/theirs`, or discard; `fo-merge-core.md` applies the same no-auto-resolve rule to local merge. `fo-dispatch-core.md` checks blocked entities before `status --next`, so the affected entity can be held while unrelated ready entities continue; the owned-worktree rule in the Claude adapter permits a rebase only when `drift.owned == true` and makes an unowned branch report-only. The old green evidence was tied to head `7918cba…`; after the owner handoff the exact-head check rejected it against `e54865d…`, so required CI and fresh gate evidence must run on the new head before merge guard resumes. The unresolved risk is scope ambiguity: the current generic word “HALT” can be read as stopping the whole event loop, and it does not identify the Captain/owning-branch semantic owner or the exact old/new base/head tuple. Implementation must distinguish a code-worktree PR conflict (per-entity hold, keep unrelated work moving) from a split-root state-sync conflict (global state safety halt), and re-check the PR head immediately before fresh validation and delivery.

- DONE: Select the smallest implementation direction and recommend entry to the implementation gate.
  Evidence: no helper, resolver worker, frontmatter field, workflow change, or provider change is needed. Add contract wording and focused behavior fixtures on the existing shared-core, merge-core, dispatch-core, and host-adapter surfaces; make the fixture assert byte-preserved pending authority, abort/path/peer evidence, explicit Captain/branch-owner handoff, exact-head invalidation, and independent dispatch. Recommend APPROVE ideation and enter implementation with this per-entity conflict-hold/re-entry contract, preserving the existing sole terminal consumer and merge authority.

### Summary

The real moving-target replay shows that a dirty PR must remain truthful and pending while the FO aborts and reports the conflict. Semantic reconciliation belongs to the Captain and owning branch; the FO owns transport/state safety and fresh exact-head evidence. Unrelated ready work may continue, but a rewritten head cannot reuse old validation. Proceed to implementation with the existing contract surfaces and no resolver escape hatch.
