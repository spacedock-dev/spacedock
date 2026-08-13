---
title: Repair Codex filing command-ledger observation
status: ideation
source: PR #679 run 31728107636, Codex job 94541783359
sprint: test-behavior-completeness
sprint-readiness: ready
score: 0.9
id: 6ker7h25hj86983e5ef71ahm
gates:
    version: 1
    records:
        - id: gate:6ker7h25hj86983e5ef71ahm:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:6ker7h25hj86983e5ef71ahm-backlog-1
              briefing:
                id: briefing:6ker7h25hj86983e5ef71ahm:backlog:attempt-1:revision-1
                digest: sha256:97242c2367a7b4e96084790b7765170b01ec932b6a9873d73c5233c694dc79cc
                request-digest: sha256:96f4def00a71e0db7e1e41b3c0d0e4405e2deed271501dcdb2b4ad62429f82b3
                room-ref: ./repair-codex-filing-command-ledger-observation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6ker7h25hj86983e5ef71ahm:backlog:1
                briefing: briefing:6ker7h25hj86983e5ef71ahm:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-13T22:09:09.620643Z"
                decision: approve
                reason: Captain approved the scoped direction for ideation.
              application:
                target-stage: ideation
                state: consumed
---

Codex created `wire-the-thing` correctly, but the filing oracle reported no `spacedock new` invocation.

## Acceptance criteria

- **AC-1 (VALUE):** The exact PR #679 Codex filing artifact is graded from the real atomic filing command evidence.
- **AC-2:** Manual, non-atomic, failed, wrong-slug, and missing-command cases remain failures.
- **AC-3:** The correction changes test observation only unless diagnosis proves a separate product defect.

## Required evidence

- Artifact from run `31728107636`, job `94541783359`, `TestLiveCommonFiling`.
- A focused fixture reproduces the observed command shape before any product change.
- Exact targeted local Codex filing passes on the final candidate.

## Stage Report: backlog

- DONE: Record the released workflow and observable harm.
- DONE: Preserve exact run, job, test, and artifact identity.
- DONE: State the product/test boundary.
- DONE: Define falsifiable positive and negative evidence.

### Summary

Repair the filing evidence boundary without changing successful filing behavior.
