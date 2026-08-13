---
title: Repair Codex smallest-sufficient mechanism regression
status: backlog
source: PR #679 run 31728107636, Codex job 94541783359
sprint: test-behavior-completeness
sprint-readiness: ready
score: 0.95
id: bfmczd31ydpp4stqjstf6xwx
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
