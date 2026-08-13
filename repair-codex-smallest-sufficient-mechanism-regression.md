---
title: Repair Codex smallest-sufficient mechanism regression
status: ideation
source: PR #679 run 31728107636, Codex job 94541783359
sprint: test-behavior-completeness
sprint-readiness: ready
score: 0.95
id: bfmczd31ydpp4stqjstf6xwx
gates:
    version: 1
    records:
        - id: gate:bfmczd31ydpp4stqjstf6xwx:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:bfmczd31ydpp4stqjstf6xwx-backlog-1
              briefing:
                id: briefing:bfmczd31ydpp4stqjstf6xwx:backlog:attempt-1:revision-1
                digest: sha256:fbe0c77dc6680f348a746c335c9bb00b50dcf4bac9c42f0d6731a50433a01316
                request-digest: sha256:56165dd9b5981ef952d2a5020d59ce18bbcaff0320dfe6946db4028d4964232c
                room-ref: ./repair-codex-smallest-sufficient-mechanism-regression/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:bfmczd31ydpp4stqjstf6xwx:backlog:1
                briefing: briefing:bfmczd31ydpp4stqjstf6xwx:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-13T22:05:22.190377Z"
                decision: approve
                reason: Captain approved the scoped direction for ideation.
              application:
                target-stage: ideation
                state: consumed
started: 2026-08-13T22:06:12Z
---

Codex completed both ready tasks but also created `roadmap-strategy.md`, violating the smallest-sufficient mechanism contract.

## Acceptance criteria

- **AC-1 (VALUE):** Codex completes the two ready tasks without creating unrelated files or orchestration.
- **AC-2:** The exact PR #679 artifact is classified as product behavior, not hidden by an evidence-only change.
- **AC-3:** A bound target-level XFAIL is permitted only if the repair cannot land immediately; XPASS remains a green alert.
- **AC-4:** The final exact targeted local Codex journey passes normally after any binding removal.

## Required evidence

- Artifact from run `31728107636`, job `94541783359`, observed code `smallest-mechanism-violation`.
- Read-only diagnosis identifies the instruction or behavior that caused `roadmap-strategy.md`.
- Focused, full, race, and exact targeted local Codex checks cover final bytes.

## Stage Report: backlog

- DONE: Record the released workflow and observable harm.
- DONE: Preserve exact run, job, code, and artifact identity.
- DONE: Classify this as a product repair.
- DONE: Define the XFAIL-or-fix decision boundary.

### Summary

Restore smallest-sufficient Codex behavior and keep any temporary expected failure explicit.
