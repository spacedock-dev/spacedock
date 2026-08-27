---
id: ntnywe6wfk1g5sersjbe5yt7
title: Verify the Pi spawn skill name resolves - bare "ensign" versus "spacedock:ensign"
status: ideation
source: "Pi/GLM FO field report 2026-08-26: dispatched ensigns repeatedly produced broken DONE formatting, and the FO's diagnosis was that the dispatch artifact passes skill \"ensign\" while the Pi skill loader wants the exact name spacedock:ensign — if true, Pi ensigns spawn without the ensign contract at all"
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
        - id: gate:ntnywe6wfk1g5sersjbe5yt7:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ntnywe6wfk1g5sersjbe5yt7-backlog-1
              briefing:
                id: briefing:ntnywe6wfk1g5sersjbe5yt7:backlog:attempt-1:revision-1
                digest: sha256:184ce70a29b675fee2ea8cd40983596f1c2d5c7a73796c1d9221d726a9e1fe9e
                request-digest: sha256:0bcc9236a024db1353e8c8ca9f770965634d5634082625db1503064403d6bcfc
                room-ref: ./pi-spawn-skill-name-resolution/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ntnywe6wfk1g5sersjbe5yt7:backlog:1
                briefing: briefing:ntnywe6wfk1g5sersjbe5yt7:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-27T00:11:28.952704Z"
                decision: approve
                reason: 'Captain approve: enter ideation to flesh out the approach and test plan'
              application:
                target-stage: ideation
                state: consumed
---

`dispatch build --host pi` emits the bare skill name `ensign` by design (internal/dispatch/build.go piSpawnSkill), on a documented assumption: "pi-subagents resolves agents and skills by directory basename only." A field report says the loader needs `spacedock:ensign`. If the assumption rotted, every Pi ensign runs without its contract, and the observed broken stage-report formatting is a symptom, not a separate defect.

## Problem

{Ideation fills this in. Seeded: the discriminator is one transcript read — open a dispatched Pi ensign's session and check whether the ensign skill content loaded near the top. If it did not, the fix is the spawn constant or a loader-version-aware form; if it did, the report-format problem belongs wholly to the stage-report protocol embed (t4) and this task closes as verified-working. Either outcome is durable evidence. Related: embed-stage-report-protocol-in-dispatch (t4, in flight, Pi-only embed) treats the symptom surface; dispatch-build-flag-form-version-skew (jha) is the same rotted-assumption class on another flag.}

## Proposed approach

{Ideation fills this in. Verify first against a real Pi host and the current pi-subagents loader; then either correct piSpawnSkill (or make it loader-version aware) or record the assumption as re-verified with the loader version pinned in the comment.}

## Risk evidence

{Backlog: the field report plus the source comment's unverified assumption decide design should start. The repeated DONE-format failures give the symptom baseline.}

## Out of scope

The stage-report grammar and its messaging (9x) and the protocol embed (t4).

## Expected surface and tolerance

Estimate: production +5 across 1 file; proof +20 across 1 file. {Backlog seed; ideation refines.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A Pi-dispatched ensign demonstrably loads the ensign contract: a live or recorded Pi dispatch transcript shows the skill content resident before stage work begins.**
Verified by: {ideation refines — seed: one transcript check on the current loader; falsifying condition: the skill name in the artifact fails to resolve and the transcript shows no contract load.}

## Test plan

{Ideation fills this in. Pi surface: the pi-live lane's relevance per the path-to-lane rule if the spawn constants change.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
