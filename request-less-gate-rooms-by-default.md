---
id: rx3daftacggfmw1pt2febw31
title: Gate rooms are request-less by default; the provider handoff is opt-in at prepare
status: ideation
source: "Captain gate-format review 2026-08-25: variance scan found 499/499 rooms carry request.json while only ~8 July-era provider results ever consumed one; q0 (subspace-r-scaffolded-gate-room, spacedock-subspace, at validation) finalizes the provider journey the file exists for"
started: 2026-08-25T17:27:46Z
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
        - id: gate:rx3daftacggfmw1pt2febw31:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:rx3daftacggfmw1pt2febw31-backlog-1
              briefing:
                id: briefing:rx3daftacggfmw1pt2febw31:backlog:attempt-1:revision-1
                digest: sha256:9a936b06ee55e932f0ef8997e0b3b79810718d78076fa69df98af75a87eae2e8
                request-digest: sha256:e9dbc6437c5a986912593c1ec9cbe6e016f561b89f7d53818221ad5ac556af86
                room-ref: ./request-less-gate-rooms-by-default/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:rx3daftacggfmw1pt2febw31:backlog:1
                briefing: briefing:rx3daftacggfmw1pt2febw31:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-25T17:27:27.228339Z"
                decision: approve
                reason: 'Captain chat 2026-08-25: ''dispatch rx'' — approved seeding into design; q0 coherence requirement baked into the seed'
              application:
                target-stage: ideation
                state: consumed
---

`gate prepare` mints `request.json` for every room. Only the subspace provider journey consumes it. Chat-presented gates are 100% of current practice, so almost every room pays for a handoff no consumer reads. Make the chat room one file, and mint the provider handoff only when the caller selects a provider presentation. Keep the minted form byte-compatible with q0's preflight.

## Problem

{Ideation fills this in. Seeded facts: prepare writes gate-briefing.json + request.json unconditionally (internal/gates/prepare.go:216, 474); request.json carries gate id, attempt id, briefing id+digest (all already in the frontmatter binding) and constant actor/approver person:captain (operation.go:530 refuses other values). The spec already documents the omission path ("Request-less and chat-only attempts may omit it") and validateGateRoomRequest no-ops on empty request-digest. The consumer that makes the file real is q0's journey: subspace-tui preflights the FROZEN request before presenting; gate record --room cross-checks resolution.by == request.approver. Also constant in practice: briefing artifacts' mediaType (501/501 text/markdown) and per-item type (501/501 absent).}

## Proposed approach

{Ideation fills this in. Seeded: presentation-channel selection at prepare. Default (chat): one-file room — gate-briefing.json only, binding carries no request-digest, riding the existing empty-RequestDigest validator path; room-shape invariant relaxes from exactly-two-files to the selected shape. Provider opt-in (flag name is ideation's): mint request.json + request-digest exactly as today, byte-compatible with q0's preflight and gate record --room. Drop the per-item type field from written manifests (unused everywhere; q0 does not read it) with back-compat reads for archived rooms. mediaType: q0's canonical five-key Result echoes artifact.mediaType, so drop it only with a recorded cross-repo agreement with the q0 contract; otherwise defer it and record why. Archived rooms and existing bindings stay readable unchanged.}

## Risk evidence

{Backlog: the 2026-08-25 variance scan (499/499 request.json; 576 attempts; ~8 provider-result rooms, July-era) and the q0 journey read (subspace-r-scaffolded-gate-room index: Intended journey steps 2 and 5) decide design should start. Ideation spikes: run q0's preflight against a room prepared by the modified binary in both shapes — the riskiest cross-repo claim is byte-compatibility of the opt-in request.}

## Out of scope

Subspace-side changes (q0 owns its preflight and Result contract). The canonical five-key Result format. Retiring the provider machinery (it is the point of keeping the opt-in). The round recorder's briefing.json/briefing.review.jsonl shapes.

## Expected surface and tolerance

Estimate net LOC change: +20, across 6 files (insertions ~+80, deletions ~-60). {Backlog seed; ideation refines with tolerance.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A chat-presented gate journey (prepare, present, record, consume) produces a room with no request.json and closes, consumes, and archives identically to today.**
Verified by: {ideation refines — seed: existing gate lifecycle tests re-pointed at the one-file default plus a fixture asserting the room's file set; baseline that fails today: every prepared room contains request.json. Falsifying edit: restore unconditional minting — the file-set assertion reds.}

**AC-2 - A provider-opted prepare mints a request.json and request-digest binding that q0's preflight and gate record --room accept unchanged.**
Verified by: {ideation refines — seed: byte-compare the opt-in request against a pre-change golden; exercise the recorder's room path against an opt-in room; if feasible, run the q0 subspace-tui preflight against the opt-in room in the spike. Falsifying edit: change any bound field — digest validation reds.}

**AC-3 - Written briefing manifests carry no per-item type field, and archived rooms with the old shape still read.**
Verified by: {ideation refines — seed: manifest writer test + a read test over a retained archived room fixture.}

**AC-4 (MEANS, serving AC-1) - The gate spec documents the request-less default and the provider opt-in.**
Verified by: {docs/specs/gate-resolution-frontmatter-contract.md diff reviewed at validation; prose follows ASD-STE100 per the workflow README.}

## Test plan

{Ideation fills this in. Seeded: internal/gates unit tests over both prepare shapes; recorder room-path test; archived-room back-compat fixture; the q0 cross-repo preflight spike recorded in the body. internal/gates is the status-mutation/guard high-stakes surface: the detached adversarial audit applies. All comments and doc text follow ASD-STE100.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

### Dispatch Retries

- Retry 1: ideation — no-completion-signal (host session limit, login restored); nudged
