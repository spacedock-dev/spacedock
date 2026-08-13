---
title: Repair Codex filing command-ledger observation
status: backlog
source: PR #679 run 31728107636, Codex job 94541783359
sprint: test-behavior-completeness
sprint-readiness: ready
score: 0.9
id: 6ker7h25hj86983e5ef71ahm
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
