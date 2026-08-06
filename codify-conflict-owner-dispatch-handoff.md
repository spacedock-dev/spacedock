---
title: Make FO dispatch of moving-target conflict owners explicit
status: backlog
source: "Captain follow-up after the 2026-08-04 conflict-owner and shared-credential diagnosis."
started:
completed:
verdict:
score: 0.93
worktree:
issue:
sprint: durable-decisions
group: fo-contract
milestone: 0.27.0
id: d8qmey415fsb5q9h6q639ngf
gates:
    version: 1
    records:
        - id: gate:d8qmey415fsb5q9h6q639ngf:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:d8qmey415fsb5q9h6q639ngf-backlog-1
              briefing:
                id: briefing:d8qmey415fsb5q9h6q639ngf:backlog:attempt-1:revision-1
                digest: sha256:e94e3a3c61d1947a0a06e071414b7cb175f4b297cf6a87f0f0c09c092d8a98e6
                request-digest: sha256:0ce185eaddb7133c17ad91e086e0d59db9adcf0ca9bc9cb64d9ffb201ef69731
                room-ref: ./codify-conflict-owner-dispatch-handoff/review/backlog/briefing-1
---

The First Officer contract names semantic ownership for an open-PR moving-target
conflict, but it does not make the owner dispatch mechanical. It also leaves
GitHub author data looking like worker identity even though all workers commit
with the Captain's shared credential. When the original worker handle is gone,
the contract gives no fresh-dispatch route to the same code/worktree owner.

## Problem

`halt.rebase-conflict` says to abort and surface paths, while the merge and
dispatch cores do not require a concrete owner-handoff package or recipient.
The FO can therefore report "the owner must reconcile" without dispatching the
owner, or send the work to the Captain because the PR author is the Captain's
credential. This is a generic First-Officer/runtime contract defect, not a
durable-decisions-only behavior and not a reason to add a resolver worker.

## Proposed approach

- Define the owner as the workflow dispatch identity plus its entity stage,
  branch, and worktree—not the GitHub author or commit credential.
- For a code-worktree moving-target conflict in an open PR, preserve the entity
  stage, `pr:#N`, merge mod-block, and pending approval; abort the local rebase;
  classify the conflict as a per-entity hold; and assemble an owner-handoff
  package containing the entity, stage, PR, branch, worktree, old/new base and
  head, exact paths, and the next owner action.
- Dispatch that package to the live owner handle when addressable. If the old
  handle is gone, fresh-dispatch the same stage's owner identity against the
  existing branch/worktree. Do not route to the Captain merely because the PR
  uses the Captain's credential, and do not spawn a generic resolver.
- Keep unrelated entities moving. A split-root state-sync conflict remains a
  global state-safety halt and is distinct from this code-worktree hold.
- When the owner produces a new exact head, invalidate old CI and gate evidence,
  rerun fresh validation/CI, and resume merge guard only from that evidence.

## Acceptance criteria

**AC-1 (VALUE) — An open-PR moving-target conflict reaches the correct worker.**
Verified by a fixture that advances the merge target after terminal approval is
pending and asserts a concrete owner-handoff dispatch to the recorded workflow
identity/branch/worktree, while `status`, `pr:#N`, merge mod-block, and pending
authority remain unchanged.

**AC-2 — Shared Git credentials cannot misidentify the owner.**
Verified by a fixture whose PR author is the Captain but whose dispatched
implementation identity is an ensign; the handoff targets the ensign identity,
not the Git author, and the fresh-dispatch fallback works after the old handle
is absent.

**AC-3 — Rewritten heads cannot reuse evidence.**
Verified by changing the PR head after a green validation result and asserting
that the old CI/Briefing is held until exact-head validation and required lanes
run again.

**AC-4 — The conflict hold is per entity.**
Verified by an event-loop fixture with one dirty PR and two independent ready
entities; the owner handoff is dispatched while the other two continue to their
next dispatch or gate action.

**AC-5 — No resolver escape hatch is created.**
Verified by detached adversarial attempts to auto-resolve, force-push, use
`ours`/`theirs`, clear authority, or spawn a generic resolver; each is rejected
while the owner-handoff path remains executable.

## Expected surface and semantic boundaries

Expected changes are limited to the generic First-Officer shared, dispatch,
merge, and runtime-adapter contract plus focused behavior fixtures. No workflow
README, stored-state schema, PR-provider policy, generic resolver worker, or
shared Git credential behavior changes.

## Test plan

Use a throwaway two-branch fixture with a Captain-authored PR and a distinct
ensign dispatch identity. Advance the base, capture the conflict, abort without
state mutation, exercise live-handle and fresh-owner dispatch paths, then
change the head and prove stale evidence rejection. Run focused dispatch/merge
tests, detached adversarial checks, `go test ./...`, `go test ./... -race`, and
format checks.

## Stage-specific test gates

- Ideation must replay the conflict and prove the credential/worker identity
  distinction before selecting contract wording.
- Implementation must preserve pending authority and emit a concrete owner
  handoff without adding a resolver or state field.
- Validation must reproduce live-owner and fresh-owner dispatch, exact-head
  freshness, per-entity keep-moving, and no-auto-resolution behavior.
