---
title: Continue Codex rejection after the first validation
status: backlog
source: "Staff review M2 for test-behavior-completeness, 2026-08-09"
started:
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
worktree:
issue:
pr:
mod-block:
id: dvddbpsf4tdt3yjw1yjyp14k
gates:
    version: 1
    records:
        - id: gate:dvddbpsf4tdt3yjw1yjyp14k:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:dvddbpsf4tdt3yjw1yjyp14k-backlog-1
              briefing:
                id: briefing:dvddbpsf4tdt3yjw1yjyp14k:backlog:attempt-1:revision-1
                digest: sha256:f8211760d0bf561280fa0671e01ccd07238ebcc2d984f118fe07f8b7a955ce67
                request-digest: sha256:aa4c58b93bc1e78c48359d493dcd20922bda0a810de46ab3d7dcc5308b93bd62
                room-ref: ./continue-codex-rejection-after-first-validation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:dvddbpsf4tdt3yjw1yjyp14k:backlog:1
                briefing: briefing:dvddbpsf4tdt3yjw1yjyp14k:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T20:35:44.798466Z"
                decision: approve
                reason: The Captain authorized shaping and requires end-user value; this task owns the complete Codex correction journey.
              application:
                target-stage: ideation
                state: pending
---

A Codex operator must receive a complete correction journey after the first rejected candidate.

The task must start from strict XFAIL code `rejection-flow-not-completed`. It must then reach correction, validation/2, and one fresh final gate.

A recorder helper or feedback-flow component does not satisfy this task. The complete Codex journey is the required value.

Do not change recorder bytes, round format, or the Pi recorder-order repair.
