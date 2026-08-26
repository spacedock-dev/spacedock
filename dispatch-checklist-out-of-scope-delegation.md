---
id: yv3w8rhxrjyqywadmy666nph
title: Dispatch checklists must not delegate status advancement out of worker scope
status: backlog
source: "claude-live smallest-sufficient-mechanism red, PR #762 run attempt 2, diagnosed from artifacts 2026-08-26: the live FO's checklist told both ensigns to 'advance status to done' — a frontmatter mutation the ensign contract forbids and documents no syntax for"
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
        - id: gate:yv3w8rhxrjyqywadmy666nph:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:yv3w8rhxrjyqywadmy666nph-backlog-1
              briefing:
                id: briefing:yv3w8rhxrjyqywadmy666nph:backlog:attempt-1:revision-1
                digest: sha256:fdeded6ab1ba8db3e1f57aa7a922ad2f54b7f62ea9904566b6bcf27a2b7c3f8f
                request-digest: sha256:a282e23565a3dd60c9d2b1c30a74446d283f395997d4be887147d801441421f0
                room-ref: ./dispatch-checklist-out-of-scope-delegation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:yv3w8rhxrjyqywadmy666nph:backlog:1
                briefing: briefing:yv3w8rhxrjyqywadmy666nph:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-26T20:04:00.456365Z"
                decision: approve
                reason: The seed defines a clear authority boundary, excludes unrelated repairs, and requires the smallest mechanism plus live proof of tool-mediated, correctly attributed transitions.
              application:
                target-stage: ideation
                state: pending
---

A live FO wrote "advance status to done" into two ensign checklists. One ensign complied by silently hand-editing frontmatter (the write-scope ban is unenforced for workers). The other honored the ban and improvised: five failed `--set` attempts, `strings` on the binary, and an FO-only skill loaded to find the syntax. The durable-journey grader then failed the run on the hand-edited transition. The FO's own state work was clean throughout.

## Problem

{Ideation fills this in. Seeded from the investigation record (this session's conversation and CI artifacts): (1) the ensign contract forbids frontmatter writes but gives no lawful path when a checklist demands one; (2) the fo-write-core mutation gate binds only FO-authored writes, so a worker's violation is silent; (3) the stage-def's passive phrasing ("the entity advances to done") leaves the actor unassigned and the FO resolved it onto the worker; (4) secondary: stream analyses must split on parent_tool_use_id or they mis-attribute ensign actions to the FO — this session's own FO made that mis-attribution before the investigator corrected it.}

## Proposed approach

{Ideation fills this in. Candidate levers, smallest first: (a) the FO contract's checklist guidance states that status advancement is never a checklist item — the FO flips status on the completion signal; (b) stage-def templates name the actor for stage transitions; (c) optionally, the ensign core documents the exact `status --set <id> field=value` form IF workers are ever meant to mutate state, or an enforcement check makes a worker frontmatter edit fail loudly. Pick the smallest set that removes the ambiguity.}

## Risk evidence

{Backlog: the CI red's artifact-level causal timeline (ensign-1 never-tried-then-hand-edited; ensign-2 four documented CLI failures then success) decides design should start.}

## Out of scope

The scenario's grader. The per-host smallest-sufficient repair tasks (h30, bf). Retroactive transcript tooling.

## Expected surface and tolerance

Estimate: production +20 across 3 files (contract/skill prose); proof +30 across 1 file. {Backlog seed; ideation refines.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A dispatched worker is never instructed to mutate entity status, and the live smallest-sufficient-mechanism journey produces a fully tool-mediated, correctly-attributed durable history.**
Verified by: {ideation refines — seed: the shipped checklist guidance names the rule; one live journey run (or the existing scenario's next scheduled run) shows the FO performing every status transition; falsifying edit: restore "advance status to done" to a checklist — the durable-journey grader reds again.}

## Test plan

{Ideation fills this in. Live-lane relevance: the shipped FO/ensign contract changes make claude-live REQUIRED per the path-to-lane rule.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
