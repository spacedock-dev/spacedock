---
title: "Recognize successful rejection-round recording under nested shell quoting"
status: "backlog"
source: "Captain filing request; pre3 and 0.27.3 live release failure triage"
started: ""
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
                state: pending
---

Recognize successful rejection-round recording under nested shell quoting.

## Problem

The 0.27.3 Codex rejection-flow grader reports zero recorder calls although gate record --round executed successfully. This differs from the completed repair-codex-rejection-round-recording task, which addressed missing gate preparation and inline durability.

Evidence: https://github.com/spacedock-dev/spacedock/actions/runs/34923156550, job 104235722711, artifact 10379038870, rejection-flow. A nested /bin/bash -lc command with a Python heredoc and quoted SPACEDOCK_BIN runs gate record rejection-task --round validation/1 and exits 0. Its output includes: round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4.

The plain equivalent in pre3 run https://github.com/spacedock-dev/spacedock/actions/runs/34923154318 is recognized. A one-off Go replay using the exact directRoundLauncher regexp returned pre3 exit=0 recognized=true and 0273 exit=0 recognized=false. Local supporting replay: /tmp/spacedock-round-regex-replay.go and .json. Root surface: internal/ensigncycle/shared_round_recording_test.go (directRoundLauncher, commandRecordsRejectionRound, codexRejectionRoundPublications).

## Proposed approach

Grade actual recorder execution and its successful result through the existing command log or correlated execution/result and durable round evidence. Prefer existing authoritative evidence over extending shell-quoting regexes. Preserve checks for exactly the required publication and reject quoted examples, echo-only text, failed commands, and duplicates. Do not change the recorder product or re-open the separate acceptance-scan mismatch.

## Risk evidence

The captured release artifacts above establish the failure trigger. Replay them before implementation; a live runtime claim requires the targeted live AC below, not a prose-presence check.

## Out of scope

Unrelated release failures and broader test-harness redesign.

## Expected surface and tolerance

Existing round-publication grader and adjacent tests in internal/ensigncycle; estimate net +20 to +80 LOC across 2–3 files. No shell parser framework, new instrumentation, product recorder changes, or new CI lane.

## Acceptance criteria

**AC-1 — Both successful command forms are recognized exactly once.**
Verified by: retained pre3 and 0.27.3 evidence replay yields one successful publication for plain and nested-shell execution. The captured 0.27.3 case fails the old oracle.

**AC-2 — Nonexecution and invalid publication remain failures.**
Verified by: behavioral controls for echo-only or quoted example text, nonzero exit, absent durable round, and duplicate publication are rejected. Counting output text alone must fail these controls.

**AC-3 — The Codex rejection-flow journey passes live.**
Verified by: a targeted local rejection-flow run observes the required round publication and gate lifecycle without weakening either assertion.

## Test plan

Add focused failing behavioral cases first. Reuse the existing test owner and captured artifacts, then run the targeted live AC. Before completion run go test ./..., go test ./... -race, and gofmt -w ./cmd ./internal. Report unavailable live execution separately from passing evidence. No prose-grep proof.

### Feedback Cycles

