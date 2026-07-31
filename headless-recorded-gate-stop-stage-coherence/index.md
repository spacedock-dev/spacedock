---
title: Make the headless recorded-gate stop fixture stage-coherent
status: backlog
source: "PR #580 run 30591046287, Sonnet job 91033369022"
started:
completed:
verdict:
score: 0.9
worktree:
issue:
milestone: 0.27.0
id: 26nk8qd48zknqnn4kc123sez
gates:
    version: 1
    current:
        gate: gate:26nk8qd48zknqnn4kc123sez:backlog
    records:
        - id: gate:26nk8qd48zknqnn4kc123sez:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:26nk8qd48zknqnn4kc123sez-backlog-1
              briefing:
                id: briefing:26nk8qd48zknqnn4kc123sez:backlog:attempt-1:revision-1
                digest: sha256:fea869611abb6a21b3bdf569d264e8c7dbc6166b5869203beec12d8aec962afb
                digest-domain: canonical-bytes
                request-digest: sha256:c6dd2c6b17d18deb57e14686317e8a856fb17c96ae5f6072c601fd0beba9b649
                room-ref: ./review/backlog/briefing-1
---

The required Sonnet live lane for PR #580 failed `TestLiveDefaultHeadlessStopsAtGate/default-headless-recorded-gate-stop`. The fixture changes `recorded-gate-task` from `status: validation` to `status: implementation` and adds an implementation stage definition, but leaves the entity body’s `## Stage Report: validation` in place. Sonnet interpreted that contradictory episode as stale authority: it withdrew the implementation attempt and manually set `status=validation`, crossing the committed no-authority boundary.

The product candidate’s `gate withdraw` did not mutate status; the FO authored the stage transition. Opus does not exercise this scenario, so its green lane does not settle the Sonnet finding.

Ideation must determine whether the supported fixture is invalid or whether the withdrawal contract lacks a real no-authority guard. Prefer the smallest behaviorally falsifiable correction. Do not add transcript/model/provider grammar or weaken the test merely to accept the observed mutation.

## Acceptance criteria

**AC-1 — The live journey starts from one coherent workflow episode.** The fixture’s status, selected gate, Stage Report, and stage definition describe the same pre-gate state; the intended implementation dispatch and validation stop can be followed without inventing a stage repair.

**AC-2 — Headless no-authority remains load-bearing.** A supported Sonnet run binds and commits the validation Briefing, presents it, and stops open without recording a decision, consuming, withdrawing current authority, manually changing status, or dispatching a successor.

**AC-3 — The fix targets the owned defect.** If fixture inconsistency is the cause, change only fixture/setup and its deterministic controls; if a supported coherent fixture reproduces the withdrawal/status mutation, route a distinct product-contract correction. No transcript/provider observer is added.
