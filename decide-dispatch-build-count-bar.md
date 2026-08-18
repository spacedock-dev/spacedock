---
id: 6htv3p97aq8sexrkjhrcdwt1
title: Decide whether one dispatch build per stage is the right bar for a live agent
status: backlog
source: "Captain CL, 2026-08-18, from the live-lane inventory reframe. The recorded-gate-lifecycle journey red in run 32105482382 on 'successor dispatch build attempts/successes = 2/2, want 1/1' — two SUCCESSFUL builds, not an error-then-retry. codex-live-dispatch-build-checklist-race already carries the open question and is codex-scoped; this occurrence is claude."
started:
completed:
verdict:
score:
worktree:
issue:
---

A journey grades an FO red for building a dispatch envelope twice. Nobody has decided whether that is a violation or normal live conduct, so the bar is currently an assumption.

## Problem

The recorded-gate-lifecycle journey asserts `dispatch.builds/successfulBuilds == 1/1` for a successor dispatch. Run 32105482382 produced `2/2` on claude — two builds that both succeeded, which rules out the benign error-then-retry sub-mode the sibling entity describes.

The question the bar begs has already been written down and left open. `codex-live-dispatch-build-checklist-race` asks whether `1/1` is the right bar for a live agent versus a scripted CLI replay. It is scoped to codex and describes `2/1` with a failed first attempt. Neither the host nor the sub-mode matches this occurrence.

`repair-opus-recorded-gate-lifecycle` is opus-scoped, names no mechanism, and is priority-held. It does not cover sonnet.

So the failure is real in the sense that the count exceeded the bar, and unproven in the sense that nobody has established the bar is right. A live FO that builds an envelope, reconsiders, and builds again with corrected inputs may be behaving well. Or it may be wasting a build and signalling confusion. The grader currently assumes the second without evidence.

Deciding this matters more than fixing this instance: the same bar now has two hosts and two sub-modes failing against it, and every future occurrence inherits the assumption.

## Proposed approach

{Ideation fills this in. Establish what a second successful build means from the retained streams of both occurrences — whether inputs changed between them, and whether the second superseded the first. Then either justify 1/1 with the reasoning, widen the bar with a stated tolerance, or widen the sibling entity host-neutral and close this one into it.}

## Out of scope

Changing dispatch build itself. The error-then-retry sub-mode owned by the sibling entity, except where this decision subsumes it.

## Expected surface and tolerance

Estimate net LOC change: +20 across 2 files if the bar changes, or 0 code if the decision is that 1/1 stands and the sibling entity absorbs this host. Declare which after the evidence read. Report insertions and deletions separately; do not declare a gross tolerance.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - The bar is justified by evidence rather than assumed.**
This is the measuring AC: the count of dispatch-build assertions in the live journeys whose expected value has no recorded reasoning must be ZERO. Verified by the decision recorded against both occurrences' actual build inputs, showing whether the second build differed from the first. Fails if the bar is restated without establishing what a second successful build means.

**AC-2 - A live FO that behaves acceptably grades green.**
Verified by replaying both recorded occurrences against the decided bar: whichever behavior the decision deems acceptable must pass, and the unacceptable one must fail. Fails if the outcome is a bar that no observed live run can satisfy, which would make the journey permanently red rather than informative.
