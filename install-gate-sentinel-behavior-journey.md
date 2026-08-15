---
id: xwbfqydr844jr6wd4ty5tghy
title: Behavior coverage for the install-gate sentinel loop bound
status: backlog
source: "Captain ruling at the retire-prose-grep-contract-tests ideation gate, 2026-08-15: accept checkless prose, file the sentinel-bound behavior test separately"
started:
completed:
verdict:
score:
worktree:
issue:
---

After the prose-grep retirement, the install-gate's sentinel one-attempt loop bound (test -f, create-before-run, per-runtime identity keys, the rm recovery message) and the uname -s OS-detection escape have zero committed coverage, and no live journey reaches them: the gate fires only when the binary is absent, which never happens in the live lanes. Design a BEHAVIOR check that exercises the shipped machinery - a binary-absent journey arm with a captive PATH, or an FO-driven fixture - where the oracle is observed behavior, never prose wording. This is a new standing check: captain-approved at the gate ruling above.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in - the simplest substrate that makes the sentinel loop observable.}

## Out of scope

Prose pins of any kind; re-creating the deleted shell-mirror harness as a second implementation.

## Expected surface and tolerance

Estimate net LOC change: +NN, test-side only.

## Acceptance criteria

**AC-1 - The sentinel one-attempt bound is proven by exercised behavior: a second install attempt in the same identity is refused, and the rm recovery path re-arms it.**
Verified by: the journey/fixture run's observed side effects; a falsifying edit (dropping the sentinel check) reds it.

**AC-2 - No assertion reads instruction-file wording.**
Verified by: the test's sources are captive filesystem state and command output only.

**AC-3 - The suite stays green.**
Verified by: the owning package plain and -race.

## Test plan

{Ideation decides the substrate and cost; live-lane budget consciousness required.}
