---
title: "Recognize successful rejection-round recording under nested shell quoting"
status: validation
source: "Captain filing request; pre3 and 0.27.3 live release failure triage"
started: 2026-09-15T04:02:58Z
completed: ""
verdict: ""
score: "0.9"
worktree: .worktrees/spacedock-ensign-recognize-executed-rejection-round-under-shell-quoting
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
        - id: gate:82de96hqv8cfwxqagj6ah1k6:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:82de96hqv8cfwxqagj6ah1k6-ideation-1
              briefing:
                id: briefing:82de96hqv8cfwxqagj6ah1k6:ideation:attempt-1:revision-1
                digest: sha256:539e67ba626efb17ef5c4e1eb43dd9abc2aedacb48af53f511239ae019caaa66
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:82de96hqv8cfwxqagj6ah1k6:ideation:1
                briefing: briefing:82de96hqv8cfwxqagj6ah1k6:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-09-15T04:12:12.560253Z"
                decision: approve
                reason: Existing logger plus durable round validation passed both exact command forms and negative controls; narrow test-only correction is ready.
                conn:
                    quote: dispatch codex live test failure tasks, verify locally for targeted failure, and open PR as stack, then trigger codex ci on stack tip.
                    source: active thread goal supplied by captain
              application:
                target-stage: implementation
                state: consumed
        - id: gate:82de96hqv8cfwxqagj6ah1k6:validation
          stage: validation
          attempts:
            - id: gate-attempt:82de96hqv8cfwxqagj6ah1k6-validation-1
              briefing:
                id: briefing:82de96hqv8cfwxqagj6ah1k6:validation:attempt-1:revision-1
                digest: sha256:c6f29d9099e118e60210d7c14b13e1c7d8460a111557332fee3c63624f358a39
                room-ref: '@review/validation/briefing-1'
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

## Stage Report: implementation

- DONE: Use the existing command logger and durable round oracle to recognize real nested-shell recording without adding command-source parsing.
  Commit 5ff85f00b wires Codex rejection-flow to recorder exit rows and retains command.log; durable room, complete log, gate, and topology checks remain unchanged.
- DONE: Commit focused positive and nonexecution, failure, missing-round, duplicate regressions and prepare the exact targeted live rejection-flow run.
  TestRejectionRoundPublicationExecution runs minimal retained pre3/0.27.3 shell forms against a built recorder; restoring source recognition fails both positives. Echo, masked recorder failure, duplicates, validation/2, and removed room fail acceptance; dropping exit/count/durable checks breaks the corresponding controls.
- DONE: Focused offline verification and stack placement.
  `go test ./internal/ensigncycle -run 'TestRejection(FlowRound|RoundPublication)' -count=1` passed, also with `-tags live` to compile live wiring; layer is based on da50d61d, 3 files +96/-17 (net +79), formatted and clean.
- SKIPPED: Targeted live AC-3 and stack-wide checks.
  FO explicitly held live/full suite/PR work pending the stacked focused commit; no live pass is claimed. Prepared command: `env -u SPACEDOCK_BIN SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_CODEX_LIVE_REQUIRED=1 go test -tags live ./internal/ensigncycle -run '^TestLiveCommonRejectionFlow$' -v -count=1` from the assigned worktree.

### Summary

The Codex publication oracle now uses actual successful recorder invocations after shell quoting resolves, preserving repeated calls and recorder failures hidden by a later successful shell command. Focused proof and implementation are committed on stack layer 3; targeted live verification remains explicitly held by the FO.

## Review-finding disposition

- Observation: `TestRejectionRoundPublicationCounter` retains a local Codex transcript-counter closure after the active Codex counter moved to execution logs (5ff85f00b).
- Released user and normal workflow: maintainers of the Codex rejection-flow test; the live path uses the new logger-backed counter.
- Observable harm: redundant legacy test maintenance; no demonstrated incorrect acceptance or rejection in the active path.
- Authority: none: the closure is not AC-1/AC-2 evidence and masks no defect in the separate real-recorder controls.
- Trigger evidence: the legacy closure runs only in the older table test; `TestRejectionRoundPublicationExecution` calls the active counter with real logger output.
- Proposal: polish, evidence-maintenance concern; owned by this test surface; decline cleanup for this narrow layer. FO explicitly authorized decline in the validation worker mailbox. No candidate edits made.

## Stage Report: validation

- DONE: Independently verify rejection layer5ff85f00 relative da50d61d uses actual recorder execution and durable state, preserves failure/cardinality/topology, and has meaningful negative controls.
  Reviewed all three changed files at exact HEAD 5ff85f00bac795ad298cddaca09fc51f27484716; clean code worktree, +96/-17 (net +79), within approved tolerance.
- DONE: Assess AC-1 — Both successful command forms are recognized exactly once.
  Existing green `TestRejectionRoundPublicationExecution` executes retained plain/nested shapes through a fresh checkout binary and the logger; both require exactly `[validation/1]`. Returning to transcript/source recognition breaks the positives. No duplicate deterministic rerun per FO instruction.
- DONE: Assess AC-2 — Nonexecution and invalid publication remain failures.
  Same executed fixture suite rejects echo-only, recorder failure followed by `true`, duplicate calls, and validation/2; removing exit/count guards breaks those controls. Removing the canonical room after successful execution must fail the durable oracle.
- DONE: Perform semantic adversarial review of the observation boundary.
  Traced shell resolution → shim real argv → recorder exit → publication list → independent canonical room/gate/topology checks; empty/failure/repeated/second-round/missing-room variants remain rejected. The shim captures recorder status before outer-shell completion and does not collapse repeat invocations.
- DONE: Confirm durable validation retains ownership of room correctness.
  `assertRejectionRecordedRound` still calls `gates.ValidateRoundFile`, checks identity, entry completeness/advisory authority, exactly one room, exact canonical bytes and preserved workflow/candidate bytes; gate preparation and native Codex worker topology remain separate unchanged live assertions.
- SKIPPED: After FO grants serialized live slot, run targeted local Codex rejection-flow on exact candidate and assess all ACs with retained artifacts, without candidate edits.
  AC-3 remains pending: FO directed a review-only report while keep-moving/headless occupy the serialized lane; no local live success is claimed. Planned artifact root: /tmp/spacedock-stack-rejection-live.
- SKIPPED: Repeat already-green deterministic suites and stack-wide checks.
  FO owns one normal/race/formatted verification at stack tip 9029decfc; no material concern justified redundant runs on this unchanged layer.

### Summary

Independent review found no material outcome or evidence defect in the changed execution-backed oracle. The only polish finding was explicitly declined by the FO; AC-1/AC-2 have implementation-stage behavioral evidence, while AC-3 remains unproven until the serialized exact-checkout live run. Recommendation: keep validation pending live evidence; do not present this report as a complete PASSED stage.

## Stage Report: validation (live completion)

- DONE: Independently verify rejection layer5ff85f00 relative da50d61d uses actual recorder execution and durable state, preserves failure/cardinality/topology, and has meaningful negative controls.
  Review-only report above remains applicable; exact code HEAD 5ff85f00bac795ad298cddaca09fc51f27484716 stayed clean and unchanged through the live run.
- DONE: After FO grants serialized live slot, run targeted local Codex rejection-flow on exact candidate and assess all ACs with retained artifacts, without candidate edits.
  `go test -tags live ./internal/ensigncycle -run '^TestLiveCommonRejectionFlow$' -v -count=1` passed (333.65s scenario, exit 0); fresh candidate binary `/tmp/spacedock-stack-rejection-live/bin/spacedock`, exact HEAD retained in `_setup/rejection-flow/source-head.txt`.
- DONE: AC-1 — Both successful command forms are recognized exactly once.
  Previously green real-recorder plain/nested fixture evidence remains unchanged; this live run additionally retained exactly one successful `gate record rejection-task --round validation/1` row in `codex-shared-scenarios/rejection-flow/command.log`.
- DONE: AC-2 — Nonexecution and invalid publication remain failures.
  Existing executed negatives and canonical-room controls remain as reviewed; the live grader kept publication cardinality, canonical round validation, final gate preparation, cycle-line, and native topology assertions enabled and all passed.
- DONE: AC-3 — The Codex rejection-flow journey passes live.
  `/tmp/spacedock-stack-rejection-live/test.log` records PASS; `rejection-topology.tsv` records implementation and validation spawns, then reuse/completion of those same workers for the correction cycle. Removing publication, durable room, gate preparation, or reuse would fail its corresponding unchanged grader.
- DONE: Retain live evidence before harness cleanup and release the serialized slot.
  Artifact root `/tmp/spacedock-stack-rejection-live` contains stdout/stderr/process result/final message/command log/topology plus `retained-temp/TestLiveCommonRejectionFlow3645076169/002` fixture/git/round state and three native rollouts under `retained-codex`; no auth.json remains in the artifact tree.
- DONE: Use the proven local host setup without global credential changes.
  Removed inherited CODEX markers, API/OAuth override variables and browser-trust marker; used CI shim with `/opt/homebrew/bin/codex` (Luna/max), local OAuth copied by the harness into its isolated home. Reproduction wrapper retained as `run.py`; no retries or candidate edits.
- SKIPPED: Repeat already-green deterministic suites and stack-wide checks.
  Per FO, normal/race/format verification belongs to the completed stack tip; this layer required only the serialized targeted live proof after its prior focused checks.

### Summary

Recommendation: PASSED. All three ACs now have behavioral evidence, including an independently run live journey against the exact candidate; no material findings remain, and the legacy test-closure polish disposition remains FO-authorized decline. The live slot was released immediately after the pass.
