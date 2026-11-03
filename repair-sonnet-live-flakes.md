---
title: Repair Sonnet keep-moving and selected-bare live flakes
status: ideation
score: "0.95"
source: "Exact-main Runtime Live E2E run 31894934349 attempts 1 and 2"
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
id: 060xp69y61yhrww23g3wvwqy
gates:
    version: 1
    records:
        - id: gate:060xp69y61yhrww23g3wvwqy:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:060xp69y61yhrww23g3wvwqy-backlog-1
              briefing:
                id: briefing:060xp69y61yhrww23g3wvwqy:backlog:attempt-1:revision-1
                digest: sha256:09c57e5fa0e468308825f688c19ec7539dc7e75e89e655d69e0a5f1416aa6012
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:060xp69y61yhrww23g3wvwqy:backlog:1
                briefing: briefing:060xp69y61yhrww23g3wvwqy:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-11-03T18:26:25.611996Z"
                decision: approve
                reason: 'Captain directed in chat on 2026-11-03: ''for those 2, do we know how to fix it? if not dispatch local probe on top of the worktree.'''
              application:
                target-stage: ideation
                state: consumed
started: 2026-11-03T18:28:41Z
worktree:
mod-block:
pr:
---
## Problem

Exact main `c6a336a33f866536d2b1ea9ef313bbbc6f244dd1` produced opposite outcomes across two Sonnet attempts:

- `TestLiveCommonKeepMovingPosture` passed attempt 1 and failed attempt 2 with `keep-moving-violation` after the First Officer bypassed the formal revise/feedback route.
- `TestLiveBreakGlassShimRecovery/selected-bare` failed attempt 1 because the Agent call did not preserve selected bare mode, then passed attempt 2.

Both behaviors are flaky on unchanged bytes. The live registry must retain active ownership until each journey is repaired and passes without an XFAIL binding.

## Proposed approach

Product repair per mechanism, both already named in the evidence (captain ruling 2026-08-16: an XFAIL owner is a product task, not a tracking stub):

1. Keep-moving: the FO bypassed the formal revise/feedback route. Recover the attempt-2 stream from run 31894934349 (or produce a fresh targeted local red under the CI shim) and quote the exact decision point where the sonnet FO leaves the route. The product surface is the FO instruction text governing that route — harden it at the deviation's entry point, with falsifying evidence, the same method the rejection-flow arc used for its adherence failures.
2. Break-glass selected-bare: the Agent call did not preserve selected bare mode. Same method: quote the deviation from the failing stream; the product surface is the dispatch adapter's bare-mode preservation instruction.
3. Retirement per AC-3: bound failure/XPASS evidence, then unchanged-byte unbound PASS.

## Acceptance criteria

- **AC-1:** Sonnet keep-moving behavior uses the formal revise/feedback route and passes repeatedly without an XFAIL binding.
- **AC-2:** Sonnet selected-bare break-glass recovery preserves blocking bare mode and passes repeatedly without owned expected-failure handling.
- **AC-3:** Each temporary expected-failure marker names this active owner and is removed only after bound failure/XPASS and unchanged-byte unbound PASS evidence.
- **AC-4:** Registry reconciliation, active-owner join, full, race, and exact Sonnet live checks pass.

## Evidence

- Run `31894934349`, attempt 1: common keep-moving PASS; selected-bare FAIL.
- Run `31894934349`, attempt 2: keep-moving FAIL; selected-bare PASS.
