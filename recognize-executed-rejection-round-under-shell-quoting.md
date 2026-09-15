---
title: "Recognize successful rejection-round recording under nested shell quoting"
status: ideation
source: "Captain filing request; pre3 and 0.27.3 live release failure triage"
started: 2026-09-15T04:02:58Z
completed: ""
verdict: ""
score: "0.9"
worktree: ""
issue: ""
pr: ""
mod-block: ""
id: 82de96hqv8cfwxqagj6ah1k6
gates:
    version: 1
    records:
        - id: gate:82de96hqv8cfwxqagj6ah1k6:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:82de96hqv8cfwxqagj6ah1k6-backlog-1
              briefing:
                id: briefing:82de96hqv8cfwxqagj6ah1k6:backlog:attempt-1:revision-1
                digest: sha256:879d692755c3f9593adae22cf3e3004f7a33e30f047b5ea0e1e0b1dd9940fa63
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:82de96hqv8cfwxqagj6ah1k6:backlog:1
                briefing: briefing:82de96hqv8cfwxqagj6ah1k6:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-15T04:02:39.792196Z"
                decision: approve
                reason: Captain requested dispatch of the Codex live failure tasks, targeted local verification, and stacked PRs.
              application:
                target-stage: ideation
                state: consumed
---

Recognize successful rejection-round recording under nested shell quoting.

## Problem

The 0.27.3 Codex rejection-flow grader reports zero recorder calls although gate record --round executed successfully. This differs from the completed repair-codex-rejection-round-recording task, which addressed missing gate preparation and inline durability.

Evidence: https://github.com/spacedock-dev/spacedock/actions/runs/34923156550, job 104235722711, artifact 10379038870, rejection-flow. A nested /bin/bash -lc command with a Python heredoc and quoted SPACEDOCK_BIN runs gate record rejection-task --round validation/1 and exits 0. Its output includes: round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4.

The plain equivalent in pre3 run https://github.com/spacedock-dev/spacedock/actions/runs/34923154318 is recognized. A one-off Go replay using the exact directRoundLauncher regexp returned pre3 exit=0 recognized=true and 0273 exit=0 recognized=false. Local supporting replay: /tmp/spacedock-round-regex-replay.go and .json. Root surface: internal/ensigncycle/shared_round_recording_test.go (directRoundLauncher, commandRecordsRejectionRound, codexRejectionRoundPublications).

## Proposed approach

For the Codex rejection journey, install the existing `writeRecordedGateLoggingShim` through `withStubPATH` before running the FO. Keep its observer log outside the workflow and retain `command.log` in the scenario artifact directory, following the existing recorded-gate journeys. Read successful `gate record rejection-task --round` executions from this log, after shell quoting has been resolved; use each recorder's exit code rather than the enclosing shell's exit code. Count every successful round id and keep the existing exactly-one-`validation/1` assertion. Reuse the existing entity/round-room validation and complete-log checks without weakening them.

Use a small fixture-specific scan of the existing `exit=0` log rows and argument fields (`--round value` and `--round=value`), not a shell grammar or success-output matcher. Wire only the Codex branch; keep Claude/Pi behavior and native worker-topology grading unchanged. This serves AC-1/2: the simpler output-only alternative admits echo/quoted text and cannot distinguish a failed recorder hidden by a later successful shell command. Durable state alone also misses repeated same-round invocations. Existing logging supplies that missing execution fact without new instrumentation or product changes.

## Risk evidence

Throwaway spike at `/tmp/spacedock-82de96-ideation-spike`, exported from `origin/main` `2a7b87198`, ran `go test ./internal/ensigncycle -run '^TestIdeationRoundLogSpike$' -v -count=1` successfully. It executed both retained commands from `/tmp/spacedock-round-regex-replay.json` through the existing logger and freshly built binary against `writeRejectionWorkflow` fixtures, replacing only the captured fixture-root path. Results: pre3 old recognizer=true, 0.27.3=false; both actual execution counts=1 and existing durable-round oracle=valid. Echo/quoted example generated no logger call; failed recorder followed by `true` added no successful publication; removed round room and duplicate publication list were rejected. This is deterministic mechanism evidence, not a completed live AC.

The release artifact URLs above are the durable source of the captured commands. Implementation must embed their minimal command cases in repository tests; the temporary replay files are diagnostic inputs, not a permanent test dependency. The only newly wired mechanism is the existing command logger; the spike proved it survives the exact failing quotation shape.

## Out of scope

Unrelated release failures and broader test-harness redesign.

## Expected surface and tolerance

Estimate net LOC change: +60, across 3 files (approximately 90 insertions and 30 deletions). Tolerance: net +20 to +80, at most 3 files: `internal/ensigncycle/shared_round_recording_test.go`, `internal/ensigncycle/claude_live_runner_test.go`, and one adjacent focused test file if needed. Existing helpers are reused, not extended into a framework.

Permitted semantic change: Codex rejection-flow test grading recognizes actual successful recorder executions independently of shell source quoting. Command grammar, stored formats, authority, recorder runtime behavior, skill text, and other hosts are unchanged. No user-facing documentation diff is needed for this test-only oracle correction.

## Acceptance criteria

**AC-1 — Both successful command forms are recognized exactly once.**
Verified by: retained pre3 and 0.27.3 evidence replay yields one successful publication for plain and nested-shell execution. The captured 0.27.3 case fails the old oracle.

**AC-2 — Nonexecution and invalid publication remain failures.**
Verified by: behavioral controls for echo-only or quoted example text, nonzero exit, absent durable round, and duplicate publication are rejected. Counting output text alone must fail these controls.

**AC-3 — The Codex rejection-flow journey passes live.**
Verified by: a targeted local rejection-flow run observes the required round publication and gate lifecycle without weakening either assertion.

## Test plan

Primary proof owner: `shared_round_recording_test.go` and the existing `TestLiveCommonRejectionFlow` journey. Add focused failing cases before implementation. Embed plain and nested retained command forms, executing them against real fixture state through the existing logger. AC-1 fails if grading returns to shell-source recognition. AC-2 controls cover echo/quoted example, failed recorder followed by successful shell completion, absent durable round, duplicate same-round calls, and a second `validation/2` publication; dropping the logger exit guard, durable check, or invocation count must respectively fail these controls. Reuse existing durable-oracle tests rather than duplicating room validation.

Focused offline command: `go test ./internal/ensigncycle -run 'TestRejection(FlowRound|RoundPublication)' -count=1` (name added cases under these existing prefixes). Cost: seconds, with one local binary build. Targeted live command from the stack layer worktree: `SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_CODEX_LIVE_REQUIRED=1 go test -tags live ./internal/ensigncycle -run '^TestLiveCommonRejectionFlow$' -v -count=1`. Cost: one Codex rejection journey, several minutes; require retained command log, valid round, final prepared gate, and unchanged topology assertions. The harness builds its checkout candidate; do not inject the stale inherited binary override.

Per captain-approved stack strategy, run focused offline checks and the targeted live journey per layer; run `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and full CI once at the completed stack tip. No candidate code or live runtime execution is claimed by this ideation stage.

### Feedback Cycles


## Stage Report: ideation

- DONE: Prove the smallest execution-backed round-publication oracle using existing evidence; preserve echo, failed-command, missing-round, and duplicate negatives.
  Throwaway `TestIdeationRoundLogSpike` passed: exact retained pre3/0.27.3 forms produced one successful logger row plus valid durable state; old recognizer missed 0.27.3; independent controls rejected false success.
- DONE: Record a bounded implementation plan and targeted Codex rejection-flow command for the approved stack.
  Plan reuses the logger and durable owner in at most 3 files, net +20–80 LOC; exact focused/live commands and per-layer versus stack-tip verification are recorded above.

### Summary

The existing command logger provides the execution fact the shell regexp misses, and its recorder-level exit code also excludes failures hidden by a successful outer shell. Both retained command forms were exercised against real fixture state without candidate code edits; implementation and the targeted live AC remain subsequent-stage work.
