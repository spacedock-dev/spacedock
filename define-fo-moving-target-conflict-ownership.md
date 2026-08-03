---
title: Define FO ownership and recovery for moving-target PR conflicts
status: backlog
source: "Captain follow-up after the 2026-08-03 durable-decisions conflict diagnosis."
started:
completed:
verdict:
score: 0.95
worktree:
issue:
sprint: durable-decisions
group: fo-contract
milestone: 0.27.0
id: g3912c6f6jkgd0yjmyg6h7yn
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
