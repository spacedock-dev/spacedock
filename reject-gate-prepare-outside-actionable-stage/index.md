---
title: Reject gate prepare outside an actionable gated stage
status: validation
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: the real binary prepared and persisted a room while the ticket was in ungated implementation, then gate record refused it."
started: 2026-08-01T14:00:57Z
completed:
verdict:
score: "0.95"
worktree: .worktrees/spacedock-ensign-reject-gate-prepare-outside-actionable-stage
issue:
pr:
sprint: durable-decisions
id: hq3d00mewqrys3s0z9pf27df
gates:
    version: 1
    current:
        gate: gate:hq3d00mewqrys3s0z9pf27df:validation
    records:
        - id: gate:hq3d00mewqrys3s0z9pf27df:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:hq3d00mewqrys3s0z9pf27df-backlog-1
              briefing:
                id: briefing:hq3d00mewqrys3s0z9pf27df:backlog:attempt-1:revision-1
                digest: sha256:f38887cdc76574701bc3ae979313bf353a919f964ae293be62737d99ccd6c663
                digest-domain: canonical-bytes
                request-digest: sha256:7792cd306d70e0ede71e31eea77881d173b857c8dd91773ab69089d5c5750e38
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:hq3d00mewqrys3s0z9pf27df:backlog:1
                briefing: briefing:hq3d00mewqrys3s0z9pf27df:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-01T13:50:37.494931Z"
                decision: approve
                reason: Captain approved dispatching this durable-decisions ideation lane in parallel with wj, nth, and jc.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:hq3d00mewqrys3s0z9pf27df:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:hq3d00mewqrys3s0z9pf27df-ideation-1
              briefing:
                id: briefing:hq3d00mewqrys3s0z9pf27df:ideation:attempt-1:revision-1
                digest: sha256:85b14b832d2c93a3fde23a0efaad822f9d97e04ace242daeb907cd1f29af92b0
                digest-domain: canonical-bytes
                request-digest: sha256:eacaaa30def18073b26b7013f54c78c5f60e32d845f01ce86a373ae474b69035
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:hq3d00mewqrys3s0z9pf27df:ideation:1
                briefing: briefing:hq3d00mewqrys3s0z9pf27df:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-02T01:30:01.201939Z"
                decision: approve
                reason: 'Captain approves the bounded pre-write actionable-stage guard: 3 files, 55-85 inserted lines, under the stated 4-file/110-line tolerance, preserving valid successor replay and existing gate formats.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:hq3d00mewqrys3s0z9pf27df:validation
          stage: validation
          attempts:
            - id: gate-attempt:hq3d00mewqrys3s0z9pf27df-validation-1
              briefing:
                id: briefing:hq3d00mewqrys3s0z9pf27df:validation:attempt-1:revision-1
                digest: sha256:0fbf8ca57395569466acc734f86f167a98694e150bc80c6de5241f146580253f
                digest-domain: canonical-bytes
                request-digest: sha256:46ce53ca3f526387c1b68d078e9e242fdc6d538694094f8de87a7efcb8f36220
                room-ref: ./review/validation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-02T08:28:30.473351Z"
                reason: The validation candidate was rebased from d6958a782 onto origin/main 988163969; retire the old open Briefing before binding fresh validation authority for ab2f095d3.
            - id: gate-attempt:hq3d00mewqrys3s0z9pf27df-validation-2
              briefing:
                id: briefing:hq3d00mewqrys3s0z9pf27df:validation:attempt-2:revision-1
                digest: sha256:930a062ad3f15e884bd19853612050cb7c6d7f6b6905e7a8e19a59b0064cfee4
                digest-domain: canonical-bytes
                request-digest: sha256:85a10b27696ece59656de141b9ac89dd516d66a5192beb771d5fc8adffa0e36b
                room-ref: ./review/validation/briefing-2
---

`gate prepare` must fail before mutation unless the ticket's current workflow stage is an actionable gate that can accept a new attempt. The command currently checks only that the stage exists. A real invocation at the ungated `implementation` stage exited zero, added a gate attempt, and wrote a room that the later recorder correctly refused as non-actionable.

## Problem

Preparation is the first durable operation in the gate journey. Allowing it at an ungated or terminal stage manufactures authoritative-looking metadata that no supported decision flow can spend. It also forced a false `hold` record in the live `1w6` journey merely to retire an attempt that should never have existed.

## Proposed approach

Make the existing workflow-stage definition the sole authority. Strengthen `validatePreparedStage`, which `Prepare` already calls before reading or constructing gate state, so it resolves the current stage definition and accepts it only when `gate: true` and `terminal` is false. Return a stable error naming the current stage and stating that it is not an actionable gate. Only after that guard succeeds may the existing retained-authority and attempt-lifecycle checks decide whether an attempt can be replayed or succeeded.

This serves AC-1 by placing policy at the existing pre-write boundary and serves AC-2 by leaving `prepareTarget` and successor allocation unchanged. The simplest alternative was to rely on `gate record`'s current actionable-stage refusal. It cannot deliver the value because preparation has already mutated the entity and created a room by then. A second cleanup/withdrawal path is also insufficient: it legitimizes impossible attempts and expands lifecycle semantics instead of preventing them.

## Common journey effect

- At `backlog`, `ideation`, or `validation` in this workflow, an otherwise eligible gate can be prepared normally.
- At ungated `implementation`, preparation exits nonzero, names `implementation`, and leaves the ticket and review tree byte-for-byte and entry-for-entry unchanged.
- At terminal `done`, preparation exits nonzero, names `done`, and has the same zero-mutation guarantee even if a malformed workflow also marks that stage `gate: true`.
- At a gated stage with an unresolved open attempt, the existing replay/divergence rules still apply; this guard does not authorize a second open attempt.
- At a gated stage whose prior attempt was legitimately retired, the existing successor path remains available.
- A legitimate stale open attempt is retired with the existing `gate withdraw` lifecycle owned by `0m6`; invalid-stage preparation is not another withdrawal case.

## Expected surface and semantic scope

Expected implementation surface: 3 files, approximately 55-85 inserted lines and fewer than 10 changed/deleted lines.

- `internal/gates/prepare.go`: strengthen the existing stage guard (about 8-15 inserted/changed lines).
- `internal/gates/prepare_test.go`: add table-driven package coverage for gated, ungated, terminal, and contradictory terminal-plus-gate definitions, plus valid successor re-entry (about 35-55 inserted lines).
- `docs/site/reference/command-reference.md`: document the precondition and zero-mutation refusal (about 12-15 inserted/changed lines).

Tolerance: up to 4 files and 110 inserted lines if a small reusable snapshot helper is needed in `internal/gates/prepare_test.go` or CLI-level coverage belongs in `internal/cli/gate_test.go`. Exceeding either bound, changing a stored schema, or touching another gate verb requires returning to the ideation gate.

Declared semantic change: runtime behavior and user-visible diagnostics for the existing `gate prepare` command change. An ungated or terminal current stage changes from success-with-mutation to nonzero refusal before mutation. Command grammar, successful output, workflow authority, gate/attempt/Briefing/Resolution stored formats, room format, valid-stage attempt lifecycle, and all other gate verbs remain unchanged.

## Out of scope

Do not change decision recording, terminal consumption, withdrawal semantics, provider presentation, or workflow stage declarations. Do not add compatibility parsing for invalid pilot attempts.

## Acceptance criteria

**AC-1 - Gate preparation creates zero durable state outside an actionable gated stage.**
Verified by: a focused table-driven `gates.Prepare` behavior test runs against one gated stage, one ungated stage, one terminal stage, and a contradictory `gate: true` plus `terminal: true` stage. The gated case succeeds. Every refusal exits through an error naming the stage and containing `is not an actionable gate`; exact entity bytes and a recursive review-tree snapshot (relative paths, entry types, and file bytes) are identical before and after. Removing the guard or checking only stage existence makes at least the ungated case create a gate record and room; checking only `gate: true` makes the contradictory terminal case mutate.

**AC-2 - Existing valid re-entry remains possible.**
Verified by: the smallest package fixture prepares at a gated stage, retires the first attempt through the existing lifecycle fixture state, and prepares successor attempt 2 at the same stage. It asserts the successor room and attempt identity while retaining attempt 1. A guard that rejects all historical records or bypasses existing successor allocation makes this fail.

**AC-3 - Operators receive the actionable-stage precondition on the command surface.**
Verified by: review of the rendered command-reference row together with the behavior test's stable diagnostic assertion; omitting either the gated/nonterminal precondition or the zero-mutation refusal fails the documentation review.

## Test plan

Add focused package tests around the existing pre-write guard, reusing `prepareFixture`. Snapshot the entity and review subtree before each rejected call and compare after the call. This is fixture-level behavior coverage with low complexity; no live workflow, transcript parser, provider fixture, or new CLI harness is needed because CLI routing delegates the policy directly to `gates.Prepare` and preserves its error as a nonzero result.

Run `go test ./internal/gates -run 'TestPrepare.*(Actionable|Successor)'`, then the repository-required `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`. The refusal test proves AC-1 by observing errors and on-disk state, not by matching implementation text. The successor fixture proves AC-2. Documentation review proves AC-3 and is paired with AC-1's behavioral evidence.

## Risk spike / mechanism evidence

No throwaway spike is needed. The riskiest boundary is already exercised and inspectable in shipped code: `Prepare` obtains the entity lock, reads `status`, and calls `validatePreparedStage` before `Read`, `prepareTarget`, `preparedRoomPath`, or any room/entity write; `applicationStages` already decodes both `gate` and `terminal` flags and is used by the existing gate lifecycle. The implementation's first focused test will seed the missing negative case at that proven seam. The simpler recorder-only alternative was disproved by the source incident: recording refused the non-actionable stage only after preparation had persisted impossible authority.

## Required documentation diff

In `docs/site/reference/command-reference.md`, change the `gate prepare` table entry from:

> Derive and bind a recorder-ready room for folder or flat form. Immediately after preparation the room contains exactly `gate-briefing.json` and `request.json`, with no copied sources or association.

to:

> At an actionable current workflow stage (`gate: true` and nonterminal), derive and bind a recorder-ready room for folder or flat form. At an ungated or terminal stage, preparation exits nonzero before changing the entity or review tree and names the rejected stage. Immediately after successful preparation the room contains exactly `gate-briefing.json` and `request.json`, with no copied sources or association.

The remainder of that row, including selected-source authority and success output, stays unchanged.

## Stage Report: ideation

- DONE: Produce a concrete, bounded problem statement, proposed approach, common scenarios, and expected surface for rejecting gate preparation outside actionable stages.
  The body defines the guard, five journey cases, a 3-file/55-85-line baseline, a 4-file/110-line tolerance, and explicit semantic boundaries.
- DONE: Define independently falsifiable end-state acceptance criteria with a reproducible test plan, declared semantic scope, tolerance, and required documentation diff.
  AC-1 observes error plus byte/tree identity, AC-2 observes successor attempt 2, and AC-3 binds the exact command-reference wording.
- DONE: Exercise or explicitly record the riskiest unverified mechanism and explain why the simplest alternative cannot deliver the value.
  The risk section records the already-proven pre-write call order and parser flags; recorder-only enforcement is too late because preparation has already persisted state.

### Summary

The ideation now bounds the fix to strengthening the existing pre-write stage guard and proving zero mutation at ungated and terminal stages. It preserves valid replay and successor behavior, declares the precise runtime-only semantic change, and includes an implementation-ready test plan and concrete documentation wording.

## Stage Report: implementation

- DONE: Strengthen the existing pre-write gate-prepare stage guard to require gate=true and terminal=false.
  Commit `d6958a782` rejects a resolved stage unless it is gated and nonterminal, before gate-state or room reads and writes.
- DONE: Add focused tests for actionable, ungated, terminal, contradictory, and valid successor-reentry cases with zero-mutation snapshots.
  `TestPrepareRequiresActionableCurrentStage` fails if any invalid stage mutates entity/review bytes or entry types, or if a held actionable attempt cannot allocate and retain successor attempt 2.
- DONE: Apply the approved command-reference wording and remain within the 3-file/55-85-line estimate and 4-file/110-line tolerance.
  The command row names the precondition and zero-mutation refusal; the committed surface is 3 files, 82 insertions, and 3 deletions, with `go test ./...` and `go test ./... -race` passing.

### Summary

Gate preparation now refuses ungated and terminal current stages before durable mutation while preserving actionable preparation and same-stage successor re-entry. The implementation updates the existing guard, focused package proof, and command reference only, and all required formatting, normal, and race checks pass.

## Stage Report: validation

- DONE: Independently verify AC-1 through AC-3 against the pre-write guard implementation and approved semantic boundary.
  AC-1: `TestPrepareRequiresActionableCurrentStage` accepts gated `validation`; refuses ungated `implementation`, terminal `done`, and contradictory gated-terminal `contradictory` with the stable diagnostic while exact entity bytes and recursive review-tree paths, types, and bytes remain unchanged. AC-2: the same fixture retires attempt 1 with the existing hold lifecycle and proves retained attempt 1 plus successor attempt/room 2. AC-3: the command row states `gate: true` plus nonterminal, nonzero refusal, rejected-stage naming, and zero entity/review-tree mutation. Removing the guard, ignoring `terminal`, rejecting historical records, or dropping either documentation promise makes this evidence fail.
- DONE: Run focused preparation tests, full and race suites, formatting and diff checks, and inspect zero-mutation/successor evidence.
  `go test ./internal/gates -run 'TestPrepare.*(Actionable|Successor)' -count=1 -v`, `go test ./...`, and `go test ./... -race` passed; recursive `gofmt -d` over `cmd` and `internal` produced no diff, and `git diff --check` passed. Source tracing confirms the guard follows status resolution but precedes gate-state reads, target/room allocation, source inspection, and writes; HEAD and the candidate worktree remained unchanged at `d6958a782`.
- DONE: Confirm the 3-file/82-insertion/3-deletion implementation remains within the 4-file/110-insertion tolerance and docs match behavior.
  `git diff --numstat "$(git merge-base main HEAD)"..HEAD` reports only `internal/gates/prepare.go`, `internal/gates/prepare_test.go`, and `docs/site/reference/command-reference.md`, totaling 82 insertions and 3 deletions; no schema, other gate verb, grammar, success output, or lifecycle surface changed.

### Summary

Recommendation: PASSED. The actionable/ungated/terminal/contradictory matrix preserves one invariant—only a gated nonterminal stage may reach durable preparation—and valid same-stage successor allocation remains intact; no scaling hazard is introduced because the added stage lookup reuses the existing bounded workflow-stage parse and linear lookup before I/O-heavy preparation. All promised ACs have independently reproduced command/state evidence, no material finding or deferred risk was found, and the implementation stays within the approved scope and tolerance.

## Stage Report: validation

- DONE: Validate AC-1 through AC-3 on exact rebased hq head `ab2f095d3` against `origin/main` `988163969`, including the composed docs semantics and actionable-stage guard.
  AC-1: `TestPrepareRequiresActionableCurrentStage` passed for gated `validation` and refused ungated `implementation`, terminal `done`, and contradictory gated-terminal `contradictory`; each refusal names its stage, contains `is not an actionable gate`, and preserves exact entity bytes plus recursive review-tree paths, entry types, and bytes. AC-2: the gated case records a legitimate hold, retains attempt 1, and prepares successor attempt/room 2. AC-3: the command-reference row states the gated/nonterminal precondition, nonzero refusal, rejected-stage naming, and zero entity/review-tree mutation. Source tracing places the guard after locked status resolution and before gate-state reads, target/room allocation, source inspection, or writes.
- DONE: Run focused gate-prepare tests, full and race suites, formatting and diff checks, and the detached adversarial evidence required by the workflow; triage failures from this head only.
  `go test ./internal/gates -run 'TestPrepare.*(Actionable|Successor)' -count=1 -v`, `go test ./...`, and `go test ./... -race` passed. Recursive `gofmt -d` over `cmd` and `internal` was empty, and `git diff --check origin/main..HEAD` passed. In a detached throwaway checkout, disabling the actionable guard turned the ungated, terminal, and contradictory cases red; checking only `gate` turned the contradictory case red; removing successor-number advancement turned the valid re-entry case red. The throwaway checkout was removed and candidate HEAD/status remained unchanged.
- DONE: Reconcile the 3-file/+82/-3 scope and record a fresh validation report and AC scan; the prior `d6958a782` validation authority is stale.
  `git merge-base origin/main HEAD` is exactly `988163969`; HEAD is exactly `ab2f095d3`; `git diff --numstat origin/main..HEAD` reports only `internal/gates/prepare.go`, `internal/gates/prepare_test.go`, and `docs/site/reference/command-reference.md`, totaling 82 insertions and 3 deletions. A fresh `status --read ... --stage validation --ac-scan` reports AC-1 with 2 citations, AC-2 with 1, and AC-3 with 2, all `unevidenced=false`. This remains within the approved 4-file/110-insertion tolerance and changes no schema, command grammar, successful output, other gate verb, or attempt lifecycle.

### Summary

Recommendation: PASSED on `ab2f095d3`. All three ACs have fresh behavioral or shipped-surface evidence, the detached audit proved the tests fail under three claim-breaking edits, and no material finding, deferred risk, or scaling hazard was found. The sole invariant is preserved across the adjacent-stage matrix: only a gated, nonterminal current stage may reach durable preparation.
