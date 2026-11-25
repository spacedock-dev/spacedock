---
id: 301s961pr7b2cnywnkz88jec
title: A gate that a second surface can read and decide safely
status: backlog
source: "MetaFO friction list 2026-11-24, sent on the word of CL, /tmp/spacedock-frictions-2026-11-24.md, group A items 1-5"
group: A-second-surface-gate
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:301s961pr7b2cnywnkz88jec:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:301s961pr7b2cnywnkz88jec-backlog-1
              briefing:
                id: briefing:301s961pr7b2cnywnkz88jec:backlog:attempt-1:revision-1
                digest: sha256:52a3c58a278bc4f991b7cffcae827ada3f87f1b68dfcc51e81e09a86a3d605a7
                room-ref: '@review/backlog/briefing-1'
---

A captain can decide a gate from a pane, a Spacedock Review page or another agent's message, and the decision lands on the briefing the captain read, with its full text and its origin recorded.

## Problem

A second surface is anything other than the first officer's own chat. Checked against origin/main `a85a26a35`:

1. A chat decision cannot name the briefing it answers. `gate record` closes whatever attempt is current (`internal/gates/operation.go:25` `RecordInput` has no expected briefing or digest; `:248` `currentStageAttempt`). A late decision can land on a briefing the captain did not read. The same fault was seen on the Review page: a decision went to version 1 while the agent waited on version 2.
2. `gate read --json` returns only the `gates:` block (`internal/status/gate_history.go:13-52`). It does not give the entity path and status, each attempt's state (`attemptState`, `internal/gates/model.go:403`), which attempt is current, or the resolved room (`ResolveRoomRef`, `internal/gates/prepare.go:786`). A pane parses the entity file itself.
3. The briefing the captain reads is not stored with the gate. `--summary` has no length limit (`internal/gates/prepare.go:85-89`: valid UTF-8 and nonblank only), so a short room is a first-officer choice. But the rendered gate review (the `present-gate` text) exists only in chat; nothing stores it in the room.
4. A captain Resolution has no field for where it came from: `by`, `at`, `decision`, `reason` (`internal/gates/model.go:69-78`), and conn fields are refused for `person:captain` (`internal/gates/operation.go:246`).
5. When another agent starts a first officer's turn, the gate presentation did not reach the captain's chat. A manager-prompt line now corrects it in one project; the product question stays open.

## Proposed approach

{Ideation fills this in. Group A changes what a gate decision means, so the design goes to CL at the ideation gate.}

Related backlog to fold or keep separate at ideation: `gate-legibility-summary-versus-excerpt`, `subspace-gate-attempt-presentation-integration`, `multi-artifact-gate-prepare`, `reliable-exact-digest-in-gate-review`, `gate-room-request-ordering-trap`.

## Risk evidence

Backlog: items 1, 2 and 4 are confirmed in the source above. Item 3 is narrower than reported: the room can hold a long summary today, but there is no place for the rendered review. Evidence of use: the Ready Room pane (`/Users/clkao/git/picard/ready-room-next/`, `/Users/clkao/git/picard/openspec/FRICTION.md`) and the first outside Review run (`/Users/clkao/git/picard/review-skill-findings-2026-11-23.md`).

## Out of scope

The pane itself; Subspace's Review page; an "act done" decision kind (task `gate-act-done-record`).

## Expected surface and tolerance

Estimate net LOC change: {ideation}, across {M} files.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A decision that names a briefing other than the current open attempt is refused, and changes nothing.**
Verified by: a CLI test that prepares attempt 1, re-prepares attempt 2, then records with attempt 1's briefing id and digest: exit non-zero, entity bytes unchanged; the same record with attempt 2's id and digest closes attempt 2. Dropping the check makes the first record close attempt 2 and fail the test.

## Test plan

{Ideation fills this in.}

### Feedback Cycles
