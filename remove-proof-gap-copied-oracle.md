---
id: pph11xwsa2s3z73ts983wd0k
title: Remove the copied proof-gap oracle from the reconciliation lint
status: backlog
source: "Captain directive 2026-08-15 after ffacf584a reintroduced the pattern zvk9 deleted; zv ensign composition report flagged it"
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:pph11xwsa2s3z73ts983wd0k:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:pph11xwsa2s3z73ts983wd0k-backlog-1
              briefing:
                id: briefing:pph11xwsa2s3z73ts983wd0k:backlog:attempt-1:revision-1
                digest: sha256:f3a9c595b6653640f2bfa8e58198d69c5b65e5a73ddd9456b5d3fef6ccb89f8c
                request-digest: sha256:55cd1712c5b2a8150909aad6747819de2dcefeb3886ee9b2c7d17c5ddf231828
                room-ref: ./remove-proof-gap-copied-oracle/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:pph11xwsa2s3z73ts983wd0k:backlog:1
                briefing: briefing:pph11xwsa2s3z73ts983wd0k:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T21:24:34.622076Z"
                decision: approve
                reason: 'Captain directive 2026-08-15: dispatch all five onto the stack tip'
              application:
                target-stage: ideation
                state: pending
---

ffacf584a added `wantProofGaps`, a hand-written map of proof-gap bindings, DeepEqual-checked against `readRuntimeProofGaps` - an AST reader that derives the same data from the live test sources three lines away. The copy's only failure mode is "you forgot to update the copy": it catches no real defect and reintroduces the two-file lockstep tax for XFAILs on live-proof tests, one abstraction above the wantGaps map the workflow just removed for the same reason.

Keep-boundary: `readRuntimeProofGaps` itself, the `parseLiveGap` shape validation it feeds, and the owner-join extension in TestRuntimeLiveTODOOwnersAreActive are load-bearing and stay.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

The runner-file liveXFail vocabulary and the registry doc.

## Expected surface and tolerance

Estimate net LOC change: -NN, 1 file.

## Acceptance criteria

**AC-1 - An XFAIL change on a live-proof test is a one-file edit that leaves contractlint green.**
Verified by: the one-file probe edit that fails at HEAD (DeepEqual mismatch) passes after; probe reverted.

**AC-2 - The change removes more lines than it adds; owner-join and shape validation survive.**
Verified by: negative delta; the owner-join and parseLiveGap tests still pass.

**AC-3 - The suite stays green.**
Verified by: go test ./internal/contractlint/ plain and -race.

## Test plan

Deletion plus the existing suite; the one-file probe both sides as the value proof.
