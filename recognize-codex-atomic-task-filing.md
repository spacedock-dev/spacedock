---
title: Recognize Codex atomic task filing
status: implementation
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
                state: consumed
        - id: gate:n513fy4c3c9mtqkq38bfh3fh:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:n513fy4c3c9mtqkq38bfh3fh-ideation-1
              briefing:
                id: briefing:n513fy4c3c9mtqkq38bfh3fh:ideation:attempt-1:revision-1
                digest: sha256:82916450a3bb3dd212ed85a0cae8ef3c858e60d07a86cbf3a8145920df074ae1
                request-digest: sha256:2a38cac752f0fbd21d6219838aead3238547a1bd4d7125c907ab1ca1d4fdb4c4
                room-ref: ./recognize-codex-atomic-task-filing/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:n513fy4c3c9mtqkq38bfh3fh:ideation:1
                briefing: briefing:n513fy4c3c9mtqkq38bfh3fh:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T20:34:02.044323Z"
                decision: approve
                reason: The exact artifact proves the missed blessed command, the three-file design is bounded, and every adjacent non-atomic or unsuccessful path has a falsifier.
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-10T20:28:23Z
worktree: .worktrees/spacedock-ensign-recognize-codex-atomic-task-filing
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

## Stage Report: ideation

- DONE: Read the complete task and artifact `9078284956` at `/tmp/dvd-final-codex.XdHiIa`.
  The filing event creates and reads `wire-the-thing`. The completed command has `exit_code: 0`.
- DONE: Inspect the filing oracle and the exact launcher command.
  The command uses `/bin/bash -lc`, captures `sd_bin`, and invokes the captured launcher after a pipeline.
  The recorded display form starts that segment with `\""'$sd_bin" new wire-the-thing`.
- DONE: Define the smallest evidence repair.
  Extend only the captured-launcher recognizer for the exact `/bin/bash -lc` display form.
  Accept only completed Codex command records with exit code zero.
- DONE: Define the actual-command positive and adjacent negative controls.
  Copy the complete artifact command into one positive control.
  Keep manual creation, `--next-id` plus a write, wrong slug, failed command, and narration-only controls red.
- DONE: Name the exact files and estimate the change before product edits.
  Change three existing test files. Estimate +90/-15 lines, 105 gross lines, and +75 net lines.
- DONE: Preserve unrelated behavior.
  Do not change dvd, n28, Pi, XFAIL policy, filed state, runtime instructions, registry rows, bindings, or other journeys.
- DONE: Define the evidence sequence.
  Run the focused control as a failing baseline before the repair.
  Then run the focused controls, exact local Codex filing journey, full suite, race suite, format, registry, and owner checks.

### Exact file design

1. `internal/ensigncycle/shared_filing_test.go`
   - Add Codex command status and exit-code fields to the existing event shape.
   - Add one shared extractor for completed commands with exit code zero.
   - Use this extractor in the Codex filing assertion.
   - Recognize the exact display-quoted captured launcher only inside `/bin/bash -lc`.
   - Require the captured variable, `new`, and the requested slug in the same command segment.
2. `internal/ensigncycle/codex_live_runner_test.go`
   - Reuse the shared successful-command extractor in `codexObservedCommands`.
   - Ignore started records and completed records with a nonzero or missing exit code.
3. `internal/ensigncycle/shared_filing_negative_test.go`
   - Build valid completed Codex records with JSON encoding.
   - Add the complete artifact command as the positive regression control.
   - Add direct controls for manual creation, preview plus write, wrong slug, failed command, and narration only.

### Control boundary

The positive control contains the exact successful command from artifact `9078284956`.
It includes the `/bin/bash -lc` wrapper, `sd_bin` assignment, pipeline, quoted launcher, slug, and follow-up reads.

The matcher does not accept a general `$VAR new` pattern.
It also does not infer success from a started record or command text alone.
Thus, a failed atomic-looking command stays red.

### Summary

The design repairs only the Codex filing evidence boundary.
It recognizes the successful atomic command and keeps each adjacent non-atomic or unsuccessful command red.

## Stage Report: implementation

- DONE: Recognize the exact successful blessed Codex filing command while every manual, preview-plus-write, wrong-slug, failed, and narration-only control remains red.
  Commit `95e7f715d` adds successful-command filtering and the narrow `/bin/bash -lc` display-form matcher. The focused control fails without either change.
- DONE: Prove the repaired exact local Codex filing journey passes normally with retained artifacts.
  `TestLiveCommonFiling` passed with local subscription authentication. The artifacts remain at `/tmp/spacedock-codex-filing-implementation.Iskzdl`.
- FAILED: Keep the approved three-file, 105-gross evidence-only boundary and complete required offline checks.
  The commit has three files and 105 gross lines. Format, registry, active-owner, focused, and focused-race checks passed.
  `go test ./...` and `go test ./... -race` stopped when the full system volume reported `no space left on device`.

### Summary

The oracle now grades only completed Codex commands with exit code zero. It recognizes the exact artifact command and keeps all required controls red.
The exact local journey passed and retained its artifacts. The full offline suites need a machine with sufficient temporary disk space.

## Stage Report: implementation (cycle 2)

- DONE: Recognize the exact successful blessed Codex filing command while every manual, preview-plus-write, wrong-slug, failed, and narration-only control remains red.
  Commit `95e7f715d` remains unchanged. The focused controls pass.
- DONE: Prove the repaired exact local Codex filing journey passes normally with retained artifacts.
  `TestLiveCommonFiling` passed. The artifacts remain at `/tmp/spacedock-codex-filing-implementation.Iskzdl`.
- DONE: Keep the approved three-file, 105-gross evidence-only boundary and complete required offline checks.
  `go test ./...` passed. `go test ./... -race` passed on the final broad run.
  The first race run had one transient 250 ms timing failure. Its focused rerun passed before the final broad run.

### Summary

The candidate remains at commit `95e7f715d` with three files and 105 gross lines. All required local checks now pass.
The retained live artifacts and the earlier disk-failure record remain unchanged.
