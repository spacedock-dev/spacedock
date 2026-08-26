---
id: enq8pfg8v1gdskqabnq3xtr1
title: gate record requires a prepared room for the open attempt it closes
status: backlog
source: "rx validation detached audit (2026-08-26): planting briefing.json beside index.json, or symlinking the room dir, flips the preparedRoomBinding discriminator to archived and gate record then closes with retained-authority SKIPPED — durable approval over tampered Briefing bytes; FO-authorized deferred risk on request-less-gate-rooms-by-default with this task as the promote path"
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
---

`gate record` closes an open attempt even when the attempt's room classifies as archived, which skips retained-authority validation. An adversarial archived-shape perturbation of a prepared room (a planted `briefing.json`, a symlinked room dir) plus a tampered Briefing yields a durable approval over tampered bytes. `Withdraw` already refuses a non-prepared room; `record` can hold the same line.

## Problem

{Ideation fills this in. Seeded: audit probes 8/9 on the rx candidate — both perturbations exit 0 and durably write the approval. Deferred-risk boundary: the trigger needs hand mutation inside the state checkout, has zero occurrences in the 700-room corpus, and the actor could already rewrite the entity binding; promote to material if any feature writes additional files into briefing-N rooms or a prepared room is observed holding briefing.json.}

## Proposed approach

{Ideation fills this in. Seeded narrow remedy from the audit: `gate record` requires `preparedRoomBinding` for the open attempt it closes, as `Withdraw` already does. Also fold the two recorded polish items: scope the new spec sentence ("Every mutating verb runs that comparison before it writes") to prepared and request-backed rooms; and consider the pre-existing `status --validate` VALID-over-deleted-Briefing scope note.}

## Risk evidence

{Backlog: the rx validation report's audit probes and disposition entry (request-less-gate-rooms-by-default, validation cycle 2) decide design should start.}

## Out of scope

The one-shape room format itself (request-less-gate-rooms-by-default ships it). q0's preflight.

## Expected surface and tolerance

Estimate: production +15 across 2 files; proof +60 across 2 files. {Backlog seed; ideation refines.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A prepared attempt whose room has been perturbed into an archived shape cannot be closed: `gate record` refuses, exit nonzero, entity bytes unchanged.**
Verified by: {ideation refines — seed: the audit's two probes as red-then-green tests (planted briefing.json; symlinked room dir); falsifying edit: drop the prepared-room requirement from record — both probes go quiet again.}

## Test plan

{Ideation fills this in.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
