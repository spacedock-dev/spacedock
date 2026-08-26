---
id: 5xpdtrr18q1fc5sbrg8z872h
title: internal/ensigncycle sits on the 10-minute test-timeout edge
status: backlog
source: "Preserved finding from patch-release-line-support implementation cycle 2 (2026-08-25): at the base commit with a busted cache the package took 570.9s of its 600s budget; one branch run measured 600.4s and timed out; the branch is untouched and strictly reduces suite time"
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
---

The `internal/ensigncycle` package runs within 29 seconds of the default 10-minute `go test` timeout on a cold cache, before any candidate change exists. Any machine slower than the measuring one, or any added scenario, flakes the whole suite with a timeout that names no test.

## Problem

{Ideation fills this in. Seeded: measured 570.9s at base with cache busted; 490.5s and 600.4s on two runs of an unrelated branch. A timeout failure is a suite-wide red that masquerades as the touched branch's fault — the exact false-attribution class the workflow's triage rules exist to prevent.}

## Proposed approach

{Ideation fills this in. Candidates: raise the package timeout explicitly in CI and document the local flag; split the package's scenarios; profile and cut the slowest fixture setup. The simplest sufficient one wins.}

## Risk evidence

{Backlog: the three measured runs above decide design should start; they are recorded in patch-release-line-support's cycle-2 stage report.}

## Out of scope

The release-machinery changes that surfaced the finding (patch-release-line-support owns them).

## Expected surface and tolerance

Estimate: production +10 across 2 files; proof +20 across 1 file. {Backlog seed; ideation refines.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A cold-cache full-suite run completes with measured headroom of at least 3 minutes on the package's budget, or the package's failure names the slow test instead of timing out whole.**
Verified by: {ideation refines — seed: timed cold-cache runs before and after, recorded in the report; falsifying edit: revert the change and the headroom collapses to the measured 29s.}

## Test plan

{Ideation fills this in.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
