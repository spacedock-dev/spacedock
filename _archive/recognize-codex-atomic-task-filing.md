---
title: Recognize Codex atomic task filing
status: done
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
        - id: gate:n513fy4c3c9mtqkq38bfh3fh:validation
          stage: validation
          attempts:
            - id: gate-attempt:n513fy4c3c9mtqkq38bfh3fh-validation-1
              briefing:
                id: briefing:n513fy4c3c9mtqkq38bfh3fh:validation:attempt-1:revision-1
                digest: sha256:33e8f9b2a88842fcbe8cef42b416aa33fcfbda1bbe5a2aa86140404865f71efe
                request-digest: sha256:832cb2778f7b0481283f6ca5c663498153b2339bec07561b88efd1a7c9be224e
                room-ref: ./recognize-codex-atomic-task-filing/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:n513fy4c3c9mtqkq38bfh3fh:validation:1
                briefing: briefing:n513fy4c3c9mtqkq38bfh3fh:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T21:28:48.175613Z"
                decision: approve
                reason: Exact local filing value passes; positional slug identity and all adjacent negative controls are independently verified; full/race and the approved scope pass.
              application:
                target-stage: done
                state: consumed
started: 2026-08-10T20:28:23Z
worktree: .worktrees/spacedock-ensign-recognize-codex-atomic-task-filing
mod-block:
pr: pr-merge:668
verdict: PASSED
completed: 2026-08-10T22:24:41Z
archived: 2026-08-10T22:24:41Z
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

## Review-finding disposition

### Validation finding: wrong slug can pass through unrelated text

- Reviewer observation: A wrong `new` slug passes when the requested slug occurs only in `--workflow-dir`.
- Defect kind: Evidence defect.
- Release scope: Material.
- Released user and normal workflow: The Codex filing oracle grades atomic task filing from successful command evidence.
- Observable harm: The oracle accepts a command that creates `other-slug`, not the requested `wire-the-thing` task.
- Affected value: value-ac[AC-3] Wrong-slug evidence must remain red.
- Trigger evidence: `spacedock new other-slug --workflow-dir /tmp/wire-the-thing` passed the candidate oracle in a detached focused test.
- Recommended disposition: Fix the slug identity boundary, then run this focused falsifier again.

## Stage Report: validation

- FAILED: Independently verify the exact three-file 105-gross evidence repair recognizes only successful blessed atomic Codex filing.
  Commit `95e7f715d` has three files and 105 gross lines. The wrong-slug falsifier proves that recognition is not exclusive.
- FAILED: Attack the manual, preview-plus-write, wrong-slug, failed-command, and narration-only controls without duplicating owned full/race/live runs.
  Six focused controls stayed red. A seventh control passed when the requested slug occurred only in `--workflow-dir`.
- DONE: Confirm exact candidate/remote, retained local Codex PASS artifacts, required check evidence, and no product/binding/Pi change.
  The clean candidate is local commit `95e7f715d`. The code branch has no remote ref or PR.
  The retained metric records `TestLiveCommonFiling` PASS at `/tmp/spacedock-codex-filing-implementation.Iskzdl`.
  The implementation transcript records format, registry, active-owner, full, and final race PASS results. Pi did not run.
  The diff changes only three `internal/ensigncycle` test files. It changes no product, binding, registry, or Pi file.
- DONE: AC-1 evidence.
  The retained metric reports `status: passed`. The retained transcript records the exact live test PASS in 50.05 seconds.
- DONE: AC-2 evidence.
  The artifact command is the positive regression control, and the retained exact live test passes with that matcher.
- FAILED: AC-3 evidence.
  The detached matrix fails only the contextual wrong-slug case. This material evidence defect violates exact task identity.
- SKIPPED: AC-4 required exact PR checks.
  No remote code branch or PR exists. Local full, race, format, registry, and active-owner evidence is present.
- FAILED: Validation recommendation.
  REJECTED. Fix the wrong-slug identity boundary before another validation run.

### Summary

The candidate recognizes the retained successful Codex command and rejects the basic negative controls. It also accepts one wrong-slug command.
This material AC-3 defect blocks delivery. The First Officer must authorize the candidate correction before implementation changes occur.

## Stage Report: implementation (cycle 3)

- DONE: Recognize the exact successful blessed Codex filing command while every manual, preview-plus-write, wrong-slug, failed, and narration-only control remains red.
  Commit `90d41cad4` binds the requested slug to the positional token after `new`. The exact contextual wrong-slug control stays red.
- DONE: Prove the repaired exact local Codex filing journey passes normally with retained artifacts.
  The retained exact-command evidence remains at `/tmp/spacedock-codex-filing-implementation.Iskzdl`. The authorized correction changed no product bytes.
- DONE: Keep the approved three-file, 105-gross evidence-only boundary and complete required offline checks.
  The candidate remains three files and 105 gross lines. Focused controls, `go test ./...`, and `go test ./... -race` passed.

### Summary

The matcher now requires the requested slug immediately after direct and captured `new` forms. Later path text cannot satisfy the slug check.
The live-runner file remains byte-identical to `95e7f715d`. The retained successful launcher-variable command still passes the focused control.

## Stage Report: validation (cycle 2)

- DONE: Revalidate exact corrected filing head `90d41cad4de0863814674c359b5763bae8e21ca3`.
  The clean local head and the remote branch resolve to the exact requested commit.
- DONE: Inspect only the cap-neutral positional-slug correction and exact contextual wrong-slug negative.
  The correction binds the slug directly after `new`. The contextual wrong-slug focused control passed.
- DONE: Confirm candidate scope and remote identity.
  The range changes three test files with 75 insertions and 30 deletions, for 105 gross lines.
  Commit `90d41cad4` matches `origin/spacedock-ensign/recognize-codex-atomic-task-filing` with no local difference.
- DONE: Confirm retained exact local Codex PASS evidence.
  The retained metric reports PASS at `/tmp/spacedock-codex-filing-implementation.Iskzdl`.
  The focused control also runs the exact artifact command assertion, which passed after the correction.
- DONE: Confirm implementation-owned full and race evidence without reruns.
  The implementation transcript records `go test ./...` PASS after the correction.
  It also records `go test ./... -race` PASS after the correction.
- DONE: Confirm no product, binding, registry, or Pi change.
  The range changes only three `internal/ensigncycle` test files. The live-runner file is unchanged by the correction.
- DONE: AC-1 and AC-2 evidence.
  The retained live metric is green, and the exact launcher-variable artifact assertion passes at the corrected head.
- DONE: AC-3 evidence.
  The contextual wrong-slug command stays red. The correction binds exact slug identity to the positional argument.
- DONE: AC-4 local evidence.
  Retained full, race, format, registry, and active-owner results pass. Pi remains skipped.
- DONE: Validation recommendation.
  PASSED. No material finding or deferred risk remains in the inspected correction surface.

### Summary

The corrected matcher requires the requested slug in the correct command position. The exact successful Codex command still passes.
The candidate remains within the approved three-file, 105-gross boundary. Validation recommends PASSED.
