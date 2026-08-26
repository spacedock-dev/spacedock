---
id: 9xnaq83nryb38fyt18mh0gbt
title: The completion-guard error names the failing sub-check
status: ideation
source: "External field report 2026-08-26 (/tmp/spacedock-durability-guard-defect.md, WhisperLiveKit sandbox session): three human round-trips to diagnose a dirty entity file because one generic error covers four distinct failures; the reporter misread the guard as requiring a remote push — refuted by source read and a live no-remote repro on 0.28.0-pre0"
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
        - id: gate:9xnaq83nryb38fyt18mh0gbt:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:9xnaq83nryb38fyt18mh0gbt-backlog-1
              briefing:
                id: briefing:9xnaq83nryb38fyt18mh0gbt:backlog:attempt-1:revision-1
                digest: sha256:484cdf1ea7ca83bc9d95c2302368bb5ca70c64fe31e2220a3b6e37e4449115f5
                request-digest: sha256:65459f1836c9fb15d4ea01d5a55184a3581fd59098ffdf0833bbec2b0dbb7473
                room-ref: ./completion-guard-error-names-failing-subcheck/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:9xnaq83nryb38fyt18mh0gbt:backlog:1
                briefing: briefing:9xnaq83nryb38fyt18mh0gbt:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-26T19:50:43.055728Z"
                decision: approve
                reason: The seed has a bounded diagnostic outcome, preserves the existing guard invariants, and names falsifiable ideation proof for each failure class.
              application:
                target-stage: ideation
                state: consumed
---

The stage-advance guard emits one message — "cannot change status away from entered stage until a durable, complete ## Stage Report is committed" — for four distinct failures: the stage-report heading is not found (exact-token match), the checklist is incomplete (missing bullets, FAILED items, or no Summary), the entity file is untracked, or the entity file differs from local HEAD. An operator who cannot tell which failed diagnoses by round-trip; one field session needed three.

## Problem

{Ideation fills this in. Seeded: internal/status/entered_stage.go hasCompleteCommittedStageReport returns one bool from hasCompleteStageReport && entityPathCleanInHEAD; handlers.go renders one string. The guard is fully local (two git commands against local HEAD) — the field report's remote-push theory arose only because the message gave no clue the failure was a dirty working tree. The exact-heading match (a parenthetical suffix defeats it) is also invisible in the message. Folded by captain direction 2026-08-26, same error-vocabulary family: (1) merge guard's "ineligible" refusal on an ungated terminal transition names no path forward — it should say the workflow's terminal transition is ungated and give the status --set form; (2) a checklist bullet written as `- DONE (annotation):` silently drops out of the guard's checklist (the colon must follow the status token), so a report can pass or fail for invisible reasons — the message should name unrecognized near-miss bullets.}

## Proposed approach

{Ideation fills this in. Seeded: return which sub-check failed and say it plainly — "no ## Stage Report: {stage} heading found (the match is exact)", "checklist incomplete: {first failing item or missing Summary}", "entity file is not tracked in {git root}", "entity file differs from local HEAD — commit it in {git root}". No behavior change to the guard itself. Note the release-machinery proof posture does not apply; this is the status-guard surface.}

## Risk evidence

{Backlog: the field report's three-round-trip timeline plus this session's own live falsification (no-remote fixture repo: local commit alone passes the guard) decide design should start.}

## Out of scope

Relaxing the exact-heading match or the clean-vs-HEAD requirement (fences, not bugs). Any remote/sync coupling (none exists).

## Expected surface and tolerance

Estimate: production +25 across 2 files; proof +40 across 1 file. {Backlog seed; ideation refines.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A blocked stage advance names the failing sub-check, and each of the four failure classes produces a distinct message.**
Verified by: {ideation refines — seed: a table test driving --set against four fixture states (missing heading, incomplete checklist, untracked file, dirty file) asserting four distinct stderr messages; falsifying edit: collapse them back to the generic string — the distinctness assertions red. Baseline that fails today: all four produce the identical message.}

## Test plan

{Ideation fills this in.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
