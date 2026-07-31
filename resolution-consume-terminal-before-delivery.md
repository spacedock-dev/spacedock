---
id: 1w62z8c5fq5g5cmhzf5sd79w
title: "Resolution-consume semantic hole: gate approval spends authority into terminal status before delivery proof, leaving no send-back path when delivery later fails"
status: ideation
source: "FO self-incident, 2026-07-31, in session fielding the fo-boot-install-hint-linux-direct-sandbox PR-merge ceremony. Captain caught the model–reality mismatch: status=done + verdict/completed unset + mod-block=merge:pr-merge, all after the validation approval was consumed into done. Only the credential delay kept the failure mode from being ratified."
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    current:
        gate: gate:1w62z8c5fq5g5cmhzf5sd79w:backlog
    records:
        - id: gate:1w62z8c5fq5g5cmhzf5sd79w:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:1w62z8c5fq5g5cmhzf5sd79w-backlog-1
              briefing:
                id: briefing:1w62z8c5fq5g5cmhzf5sd79w:backlog:attempt-1:revision-1
                digest: sha256:45c7f88e4a02a5eea0d3febe7431f09b3425a4e963fe0eff7bef7d9e5398bb84
                digest-domain: canonical-bytes
                request-digest: sha256:f808d3e34848fa57e0df5e17b9e278ccce23b1c7d8a77548f944bda23ce7a343
                room-ref: ./resolution-consume-terminal-before-delivery/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:1w62z8c5fq5g5cmhzf5sd79w:backlog:1
                briefing: briefing:1w62z8c5fq5g5cmhzf5sd79w:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-07-31T14:25:14.857596Z"
                decision: approve
                reason: 'Captain''s explicit order in this session (2026-07-31): dispatch ideation for 1w62 and discuss the proposed solution with him — accepts the direction that the consume-into-done semantic hole (caught at z3''s merge boundary) must be designed before further terminal-consume ceremonies ratify the pattern.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

A gate approval on a terminal-target stage (`validation → done`) is *immediately consumable into terminal status*: the binding approval's authority is marked spent (`consumed=true`) and the entity's status flips to `done` at consume time, while delivery (PR push, CI, merge) and workflow terminal fields (`verdict`, `completed`) remain pending. Three desyncs follow:

1. **Authority is spent before the spend is justified.** The consumed approval is irrevocable by design ("never double-applied"), yet the deliverable it ratified is only *pending* — there is no delivery proof at the moment the approval is spent.
2. **Two status classes per entity.** `status=done` (decision monoid) with empty `verdict`/`completed`/`pr` (delivery incomplete). The mod-block guard withholds terminal fields deliberately, but the status field itself already claims terminal. CL feedback: "right now if it's done it can't be sent back."
3. **No regression verb exists.** If the PR fails (CI red, closed-unmerged, a new defect uncovered), the workflow's own pre-done `feedback-to: implementation` wiring is *inert* (it's scoped to pre-done rejections). The pr-merge mod's answer to a closed PR is literally "report to the captain and wait"; the only physical reverse routes are frontmatter surgery or a second gate cycle that would have to dance around the consumed record. Prior live precedent (`worktree-clear-guard-misfire-on-backward-route`) shows even bookkeeping of a backward route trips the same merge-hook guard.

The pure-`merge: local` workflow hides this; any PR-gated workflow exposes it. Three review rounds missed it because they scoped to the API of gates, not the failure mode of delivery-after-consume.

## Problem

{Ideation fills this in. Seed: the consumption semantic treats "terminal-stage approval" and "terminal delivery" as one event; they are two, and the second can fail. The model needs a semantic where the approval-authority, the stage status, and the terminal-fields cohere in BOTH directions of that boundary.}

## Proposed approach

{Ideation fills this in — solution discovery explicitly left to ideation. Candidate shapes from the incident discussion (not authorized, for ideation's consideration): (a) spend-deferred — approval stays `approved-pending` and the consume only marks authority spent when the merge hook *proves* delivery; (b) a regression verb (`gate regress` / `merge halt --send-back`) recording a reversal with routing into `feedback-to`; (c) a stage split between closing and delivered (`done` reserved for delivered work) so the merge-hook ceremony lives in a non-terminal stage. Filing names the problem and semantic inconsistency only.}

## Out of scope

{Ideation fills this in. Seeded exclusions: changing the pr-merge hook's PR body creation mechanics, changing `merge: local` workflow semantics (which hide the hole), and revisiting the shipped fo-boot-install-hint work (its fallout is the incident, not the target of this task).}

## Acceptance criteria

Each AC names a property of the finished entity. Seeded; ideation refines.

**AC-1 - A PR-gated task whose delivery fails after validation-gate approval has a structured, recorded path back to a working stage, without frontmatter surgery and without pretending the consumed approval was re-bought.**
Verified by: {ideation names the test — a fixture that drives a gated terminal-target entity to a post-approval delivery failure and asserts the send-back path records a regression and routes the work (vs the current "report and wait").}

**AC-2 - A single canonical "is this entity really finished?" answer is derivable from frontmatter without composing two status classes, for all four terminal-target hook classes in the repo (no hook, `merge: local`, `merge: pr`, any future hook).**
Verified by: {ideation names the test — a fixture or status-output golden that fails on the current `done-but-undelivered` shape and passes on the fixed shape.}

**AC-3 - The fix does not regress the clean-mixed workflow (`merge: local`) case: a merge-hook-free task reaches done when the last gate approval is consumed (as today).**
Verified by: {ideation names the test — regression fixture for the no-hook path.}

## Test plan

{Ideation fills this in: given the incident, the tests should exercise the delivery-failure-after-consume path as a first-class state, not via report-and-wait; include the worktree-clear-guard-misfire case as a fixture sibling.}
