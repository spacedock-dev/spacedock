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
