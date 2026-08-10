---
title: Recognize Codex atomic task filing
status: backlog
score: "0.90"
source: "PR #664 Codex filing evidence failure, 2026-08-10"
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
id: n513fy4c3c9mtqkq38bfh3fh
gates:
    version: 1
    records:
        - id: gate:n513fy4c3c9mtqkq38bfh3fh:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:n513fy4c3c9mtqkq38bfh3fh-backlog-1
              briefing:
                id: briefing:n513fy4c3c9mtqkq38bfh3fh:backlog:attempt-1:revision-1
                digest: sha256:126f8c887e3849becc7f8fbfedc0d71c5f7e7ff1b3aca212332cf0d267e1961a
                request-digest: sha256:e83c57568bf3cd06862e4d971bc04d97b767c22ae8fef6b57b41e5ca4a81c81b
                room-ref: ./recognize-codex-atomic-task-filing/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:n513fy4c3c9mtqkq38bfh3fh:backlog:1
                briefing: briefing:n513fy4c3c9mtqkq38bfh3fh:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T20:28:10.455087Z"
                decision: approve
                reason: Captain directed immediate ideation for the exact Codex filing evidence repair.
              application:
                target-stage: ideation
                state: pending
---
## Problem

Codex created `wire-the-thing` through the blessed atomic `spacedock new` command. The filing oracle did not recognize the launcher-variable command and reported no invocation.

## Value

A task created through actual `spacedock new` command evidence passes the filing journey. Manual or non-atomic creation remains red.

## Scope

- Repair only Codex filing command evidence.
- Use PR #664 artifact 9078284956 as the exact baseline.
- Preserve the filed task state and existing product command.
- Do not change dvd, n28, Pi, XFAIL policy, or unrelated journeys.
- Do not add a permanent XFAIL.
- Use local Codex subscription authentication before required PR CI.

## Acceptance criteria

- AC-1: Exact local Codex `TestLiveCommonFiling` passes normally and retains artifacts.
- AC-2: The oracle recognizes the actual launcher variable command that invokes `spacedock new wire-the-thing`.
- AC-3: Manual file creation, `--next-id` plus a write, wrong slug, failed command, and narration-only evidence remain red.
- AC-4: Full, race, format, registry, active-owner, and required exact PR checks pass. Pi remains skipped.

## Baseline evidence

- Released user and workflow: Codex atomic seed-task filing.
- Observable harm: a successful atomic filing is falsely rejected by the command oracle.
- Value authority: the filing journey requires the blessed `spacedock new` path.
- Trigger: run 31427144353, job 93581749681, artifact 9078284956. The command sets `sd_bin=${SPACEDOCK_BIN:-spacedock}` and invokes `$sd_bin new wire-the-thing`; the oracle reports no invocation.

## Ideation requirements

- Name exact files and gross/net estimate before product edits.
- Define one actual-command positive and adjacent non-atomic negative controls.
- Keep local failing baseline, repaired normal PASS, validation, PR, and merge flow.
