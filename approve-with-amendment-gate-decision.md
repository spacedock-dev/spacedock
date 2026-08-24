---
id: trb9t7dv2fbzhgt64b3n35pg
title: Approve-with-amendment gate decision - record captain amendment bytes in the resolution and carry them to the successor stage
status: backlog
source: "Captain CL in chat, 2026-08-24: 'file approve+amendment protocol' - motivated by the install-sh-edge-prerelease-parity ideation gate, where a one-line captain doc amendment cost withdraw -> route-to-ensign -> re-prepare -> re-present -> second captain decision"
started:
completed:
verdict:
score:
worktree:
issue:
---

Add an approve+amendment decision to the gate protocol so a captain can approve a presented snapshot AND attach a minor amendment in one touch, instead of the full withdraw/re-prepare/re-present cycle. Live evidence 2026-08-24 (entity tdng3g6fe5, ideation): a one-line doc correction cost a withdrawal (attempt-1, digest 0541b9d8), an ensign amendment round (commit 8fb91ed06), a re-prepared room (attempt-2, digest 21d4a15f), and a second captain decision - and the second decision had the captain re-approving prose the captain authored. Byte-honesty must be preserved: the room binds a digest, and unanchored amend-after-approve would let approvals drift onto unreviewed bytes.

## Problem

{Ideation fills this in. Seeded: the recorder's decision vocabulary (approve|revise|hold) has no shape for "approve, and also fold in this minor change". Any post-presentation body change stales the briefing, forcing withdrawal and a second full cycle even when the amendment is captain-authored - the one reviewer who needs no re-presentation of their own words.}

## Proposed approach

{Ideation fills this in. Seeded design from the FO/captain discussion 2026-08-24:
1. `gate record <entity> --decision approve --amendment-file F` - amendment bytes recorded VERBATIM in the resolution, captain-attributed; the approval covers snapshot-plus-amendment, both reviewed by construction (snapshot presented; amendment is the captain's own words).
2. Consume proceeds normally; `dispatch build` forwards the amendment to the successor via the existing scope-notes/feedback-context transport - no new plumbing shape.
3. The successor worker's first act folds the amendment into the entity body; fold correctness is verified at the NEXT gate, where fold-verification naturally lives (ideation proposes, implementation applies, validation checks).
4. Materiality is the captain's call at decision time; abuse is policed downstream by the existing baseline-deviation measurement and AC cross-check against the approved baseline.
Boundaries: refuse --amendment when the gate's successor is terminal (nothing downstream verifies the fold - the withdraw/re-present path stands there); amendment chains on one gate escalate like feedback cycles so approve+amendment does not become design-by-amendment; an FO-under-conn amendment cites the grant like any delegated decision.}

## Out of scope

Advisory results on open gates (retain-open-gate-advisory-resolution, vp4) and recorder workflow-neutrality (workflow-neutral-advisory-round-recorder, zrcg) - coordinate durable storage with both, do not duplicate it. Relaxing digest authority or one-use consume semantics. Amendments authored by anyone but the deciding actor.

## Expected surface and tolerance

Estimate net LOC change: ~+200, across ~6 files (recorder Go + resolution schema, dispatch build forwarding, fo-gate-lifecycle and present-gate skill text, tests). Ideation refines with tolerance.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Seeded; ideation refines and re-anchors.

**AC-1 (value) - A minor captain amendment at a non-terminal gate costs one captain touch end-to-end.**
Verified by: a fixture/journey that replays the td-class round (open gate, captain approve+amendment, successor dispatch) and asserts exactly one captain decision event in the durable record, against the two-touch baseline the 2026-08-24 cycle produced. Fails if the flow still requires a withdrawal or a second decision for a captain-authored amendment.

**AC-2 - The resolution retains the amendment bytes verbatim with actor attribution, and the approval remains bound to the presented digest.**
Verified by: recorder fixture asserting stored amendment bytes equal input bytes and the resolution still names the bound briefing digest. Fails if amendment storage rewrites, summarizes, or detaches from the snapshot.

**AC-3 - The successor dispatch carries the amendment and the next gate can see whether it was folded.**
Verified by: dispatch build output containing the amendment transport, plus a gate-side read exposing amendment-vs-body fold state. Fails if the amendment reaches neither the worker nor the next gate's evidence.

**AC-4 - A terminal-successor gate refuses --amendment.**
Verified by: recorder test asserting refusal with a named reason and no write. Fails if an amendment can ride an approval into a terminal transition unverified.

## Test plan

{Ideation fills this in. Seeded: recorder unit tests (store/refuse/attribution), a dispatch build forwarding test, one CLI fixture journey for AC-1, skill-text updates validated by the existing smoke tests.}
