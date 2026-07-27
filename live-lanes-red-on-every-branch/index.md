---
id: se0v37bt7mhsrmhta1nyns0r
title: "Runtime Live E2E has failed on every branch for three days, so no merge since 2026-07-24 carries a live signal"
status: ideation
source: "Found attributing PR #571's four red live lanes, 2026-07-27. The change was exonerated by baseline comparison and local reproduction; the lanes were already red on every branch."
started: 2026-07-27T17:05:38Z
completed:
verdict:
score: 0.85
worktree:
issue:
gates:
    version: 1
    current:
        gate: gate:se0v37bt7mhsrmhta1nyns0r:ideation
    records:
        - id: gate:se0v37bt7mhsrmhta1nyns0r:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:se0v37bt7mhsrmhta1nyns0r-backlog-1
              briefing:
                id: briefing:se0v37bt7mhsrmhta1nyns0r:backlog:attempt-1:revision-1
                digest: sha256:a19bdffd9042fe57c50bc9d86b6f84cbea1c947d4b3929bd171cb44e3c99d875
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:se0v37bt7mhsrmhta1nyns0r:backlog:1
                briefing: briefing:se0v37bt7mhsrmhta1nyns0r:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-07-27T17:03:46.899149Z"
                decision: approve
                reason: Captain approved the release-blocking ideation plan after reviewing the bound failure classification and local-proof boundary.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:se0v37bt7mhsrmhta1nyns0r:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:se0v37bt7mhsrmhta1nyns0r-ideation-1
              briefing:
                id: briefing:se0v37bt7mhsrmhta1nyns0r:ideation:attempt-1:revision-1
                digest: sha256:05243ce76ee0e635de8309e9582fb548dc47aba260cc8c8a960fa6bc9ce967a2
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-1
---

Restore a trustworthy live signal without deleting the behavior the signal was
created to protect. The latest retained run, `30257280066`, supplies enough
transcript, command-log, entity, and git evidence to rule every current red as
conduct, oracle, or harness before implementation starts.

## Binding authority and corrected baseline

The archived 6y record, `first-officer-gate-command-lifecycle`, is the authority
for the recorded-gate journey. Its approved contract requires one semantic
root-assistant gate review after the selected Briefing commit and before the
decision, an open bound stop when no authority exists, one consumed decision,
one successor dispatch, and one later durable effect on every supported host.
The recorded-gate scenario is therefore an approved cross-host obligation, not
a candidate for withdrawal. Cycle 37's Pi review-text deferral conflicts with
that contract and is superseded here.

Run `30097092217` remains the last wholly green Runtime Live E2E run, but it is
not proof that today's gate-guardrail oracle was green: `live_gate_stop_test.go`
was substantially rewritten by `9577380d` and strengthened by `75b4aaca` after
that run. It also predates the recorded-gate scenario. Baseline status can
attribute provenance; only the current grader history and retained artifacts
can classify the current failures.

Run `30257280066` is the comparison point:

| lane | failing journey | ruling |
| --- | --- | --- |
| Claude Sonnet | default gate-stop | oracle |
| Claude Sonnet | gate-guardrail | oracle |
| Claude Opus | default gate-stop | oracle |
| Claude Opus | gate-guardrail | oracle |
| Claude Opus | recorded-gate-lifecycle | harness |
| Codex | gate-guardrail | oracle |
| Codex | rejection-flow | harness fixture |
| Codex | keep-moving-posture | oracle |
| Pi | recorded-gate-lifecycle | conduct |

## Failure rulings and bounded repairs

### Gate guardrail: oracle

Evidence: the Claude default and shared journeys and the Codex shared journey
committed the expected Briefing id and digest and left the selected attempt
unresolved. Their shared grader nevertheless requires the literal YAML
fragment `state: open`. Production `gates.Attempt` has no `state` field;
`gates.Read` represents open as `Resolution == nil && Application == nil`, and
`internal/gates/gates_test.go` rejects `attempt.state` as an unknown field. The
deterministic fixture passed only because it injected state that the real CLI
cannot serialize.

Repair `assertGateHeld` in
`internal/ensigncycle/gate_assert_impl_test.go` to decode frontmatter with
YAML `KnownFields` into the production `gates.Document`, call
`gates.Validate`, and assert the selected attempt, expected Briefing, nil
resolution, and nil application. Replace the impossible positive fixture and
retain malformed/wrong-Briefing/resolved/applied counterexamples in
`internal/ensigncycle/gate_assert_test.go` and
`internal/ensigncycle/shared_scenarios_negative_test.go`.

Falsification: a production parse failure, a missing/wrong Briefing, or any
resolution/application must remain red. If a host's retained post-state fails
those structural checks, reclassify that host as conduct; do not broaden the
oracle.

### Claude Opus recorded gate: harness

Evidence: the root stream contains a truthful review after bind with entity,
stage, exact Briefing id/digest, recommendation, and decision ask, and the
close/consume/successor lifecycle completed. Opus had changed into a nested
entity directory. The logging shim then ran `git -C .spacedock-state` relative
to that mutable cwd, emitted `fatal: cannot change to '.spacedock-state'`, and
omitted the `state-head` marker which authorizes the review extractor to start.
The extractor discarded a valid review and reported that its facts were
missing.

Repair only `writeRecordedGateLoggingShim` and its controls in
`internal/ensigncycle/recorded_gate_lifecycle_test.go`: derive the fixture's
absolute state root from the already-absolute command-log path when writing the
shim, embed that root, and compute `state-head` there. Retain the root-message,
successful-commit, and event-order filters.

Falsification: add a nested-cwd successful-commit case that must extract the
review, plus uncommitted, child/tool-only, and review-before-bind controls that
must remain red. If an absolute-root trace still contains no qualifying root
review, reclassify it as conduct.

### Pi recorded gate: conduct

Evidence: Pi bound and committed the Briefing and completed the durable
close/consume/successor lifecycle. Its review appeared only inside a shell
`cat` tool call/tool-result; the root session has no assistant text block before
the decision. `present-gate` expressly says child/tool output does not count,
and archived AC3 applies that root-review obligation to Pi.

Repair the contradictory final paragraph of
`skills/first-officer/references/pi-first-officer-runtime.md`: ordinary worker
proof remains durable-state based, while gate presentation is explicitly a
root-session assistant text event before the decision mutation. Shell output,
tool results, child output, and later summaries do not qualify. Pin that
adapter rule in
`internal/contractlint/fo_function_reference_invariant_test.go`.

This remains tool-agnostic: it names the semantic output channel and ordering
boundary only. It does not prescribe a shell command, exact wording,
recommendation, decision, or answer; the existing host-neutral gate template
still owns the facts. The Pi extractor is not weakened.

Falsification: the current tool-result-only retained session must stay red,
while a root assistant text event with the required facts after bind and before
decision passes. Missing facts, wrong order, child text, and tool-result text
remain red.

### Codex rejection flow: harness fixture

Evidence: the fixture says it starts “before first implementation” but seeds
`status: implementation`. Normal FO routing treats status as the completed
stage and `status --next` correctly returns validation, so Codex began with the
reviewer. It subsequently performed the rejection, implementation correction,
and reviewer rerun correctly, leaving one real implementation report because
the fixture had skipped the first implementation by construction.

Seed `status: backlog` in
`internal/ensigncycle/shared_fixtures_test.go`, so normal routing dispatches the
first implementation. Extend the existing fixture-negative control in
`internal/ensigncycle/shared_scenarios_negative_test.go` to require
`backlog -> implementation` and zero pre-seeded implementation reports. Keep
the oracle's two implementation reports, two validations, rejection record,
correction, and rerun unchanged.

Falsification: after the seed repair, one report or one validation is a conduct
failure. A prompt change that talks the FO around invalid seed state does not
satisfy this repair.

### Codex keep-moving: oracle

Evidence: the trace advances the approved entity, builds three dispatches,
waits, reads three durable implementation reports, and records exact successful
outputs `finalized: <entity> -> done` for all three entities. The extractor
can reconstruct the required build, wait, and durable-report phases, while its
merge helper requires the literal entity name in the command and misses the
valid `for slug ...; merge guard "$slug"` confirmation. Staff review also found
the converse false positive: `codexKeepMovingTrace` directly treats a literal
`merge guard <entity>` command as dispatch even when no build, wait, or durable
report exists. Both sides must be repaired together so the oracle is neither
blind to the retained trace nor green on a planted command.

In `internal/ensigncycle/shared_keep_moving_test.go`, remove
`kmMergeGuardTerminalizes` as direct positive dispatch evidence for the
approved and independent entities. Keep it as a negative forward-drive signal
for the questioned entity: attempting to merge that entity before its
correction folds must still fail. Existing explicit `spawn_agent`,
standing-worker `status=done`, and phase-bound dispatch evidence remain
accepted. The retained batched trace qualifies through the existing
`dispatchEvidence.stageReport` assignments after build, wait, and named report;
merge output is terminal corroboration, never standalone dispatch proof.

In `internal/ensigncycle/shared_keep_moving_negative_test.go`, convert the
existing `mergeGuardSurface` at lines 252–260 and its preceding positive
comment into the promised host-neutral unphased negative: keep its literal
merge commands and otherwise pass-shaped final message, but require
`assertCodexKeepMoving` to fail without dispatch build, wait, and durable
reports. The adjacent `doneSurface` positive remains unchanged, so the
standing-worker `status=done` dialect is still protected while raw merge
syntax loses positive credit.

In `internal/ensigncycle/codex_dispatch_evidence_test.go`, allow an exact
successful `finalized: <entity> -> done` line from a batched merge-guard command
to attach to that entity only after its state machine has observed successful
dispatch build, completed wait, and a durable stage report. Add a
retained-shape positive and missing-build, missing-wait, missing-report,
failed-command, missing-entity, literal-command-only, and planted-output
controls in
`internal/ensigncycle/codex_dispatch_evidence_regression_test.go`.

Falsification: an otherwise pass-shaped keep-moving transcript containing
literal `merge guard approved-gate`, `merge guard ready-one`, and `merge guard
ready-two` but no dispatch build, wait, or durable reports must fail the
host-neutral assertion. Removing any one phase from the retained-shape positive
must also fail; a finalization-looking string from a failed/non-merge command
remains red. If the fully phase-bound retained trace still finds no report for
an entity, that entity's failure is conduct.

## Implementation surface and boundary

Exactly these twelve existing files are planned:

| failure | files | hand-written estimate |
| --- | --- | ---: |
| gate oracle | `internal/ensigncycle/gate_assert_impl_test.go`, `internal/ensigncycle/gate_assert_test.go`, `internal/ensigncycle/shared_scenarios_negative_test.go` | +46/-18 |
| Opus harness | `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +22/-4 |
| rejection fixture | `internal/ensigncycle/shared_fixtures_test.go`, shared negative file above | +9/-2 |
| keep-moving oracle | `internal/ensigncycle/shared_keep_moving_test.go`, `internal/ensigncycle/shared_keep_moving_negative_test.go`, `internal/ensigncycle/codex_dispatch_evidence_test.go`, `internal/ensigncycle/codex_dispatch_evidence_regression_test.go` | +73/-21 |
| Pi conduct | `skills/first-officer/references/pi-first-officer-runtime.md`, `internal/contractlint/fo_function_reference_invariant_test.go` | +14/-1 |
| operator docs | `docs/runtime-live-ci.md` | +3/-3 |

Expected hand-written total is +167/-49, 216 changed lines across 12 existing
files. Generated/golden impact is exactly 0 files and 0 lines; live artifacts
remain uncommitted. Implementation must stop for re-approval above 14 files or
292 changed hand-written lines (two files or 35 percent above the estimate),
or if any production Go package becomes necessary.

No scenario is withdrawn, quarantined, skipped, or weakened. There is no
compatibility layer, new standing harness, prompt coaching, fixture prompt
change, or CI trigger. Existing live runners, transcript formats, command log,
production gate model/validator, and negative-test files are sufficient.

### Exact Pi behavior and operator documentation delta

In `skills/first-officer/references/pi-first-officer-runtime.md`, replace this
exact current paragraph:

```text
The durable proof for Pi support is not transcript phrasing. A valid live proof dispatches a Pi ensign against a temp split-root workflow and verifies process exit, state checkout file changes, git log, and stage report content.
```

with:

```text
Ordinary Pi worker proof is durable-state based: dispatch a Pi ensign against a temp split-root workflow and verify process exit, state checkout file changes, git log, and stage report content.

Recorded-gate presentation is the explicit exception. On Pi, the captain-facing review is one root-session assistant text block after the selected Briefing commit and before the decision mutation; shell output, tool results, child output, and later summaries do not qualify. Recorded-gate lifecycle proof combines that root event with the durable command, state, git, and successor-effect evidence.
```

This is tool-agnostic host behavior: it binds channel and order only, without
prescribing a command, exact wording, recommendation, decision, or answer.
`internal/contractlint/fo_function_reference_invariant_test.go` pins those
semantic invariants and fails if the Pi adapter again says transcript events
never matter or permits child/tool output.

In `docs/runtime-live-ci.md`, replace the current local Pi paragraph and command:

```text
Run the Pi front-door smoke locally (`npm install -g pi-coding-agent`, `pi install npm:pi-subagents`, and either `pi login` or `OPENAI_API_KEY`). The smoke loads the current checkout's Spacedock first-officer and ensign skills plus the local pi-subagents extension/skill explicitly; it verifies durable state in the split-root state checkout rather than transcript wording alone.

go test -tags live -count=1 -run TestLivePiFrontDoorSmoke ./internal/ensigncycle -v
```

with:

```text
Run the Pi live proofs locally (`npm install -g pi-coding-agent`, `pi install npm:pi-subagents`, and either `pi login` or `OPENAI_API_KEY`). `TestLivePiFrontDoorSmoke` loads the current checkout's Spacedock first-officer and ensign skills plus the local pi-subagents extension/skill and verifies durable split-root worker state. `TestLivePiRecordedGateLifecycle` loads the same current-checkout skills through the Pi front door, drives the shared recorded-gate fixture, and requires the root-session assistant review after the committed Briefing and before the decision in addition to durable command, state, git, and successor evidence.

go test -tags live -count=1 -run '^(TestLivePiFrontDoorSmoke|TestLivePiRecordedGateLifecycle)$' ./internal/ensigncycle -v
```

Also replace the current `pi-live` GitHub-setup sentence:

```text
Installs `pi-coding-agent`, `pi-subagents`, and `pi-intercom`, runs the Pi shared coverage guard plus `TestLivePiFrontDoorSmoke`, and uploads artifacts under `live-artifacts/pi/`.
```

with:

```text
Installs `pi-coding-agent`, `pi-subagents`, and `pi-intercom`, runs the Pi shared coverage guard, `TestLivePiFrontDoorSmoke`, and `TestLivePiRecordedGateLifecycle`, and uploads artifacts under `live-artifacts/pi/`; the recorded-gate test additionally grades the ordered root assistant review event.
```

Public site behavior does not change; the host adapter and operator guide are
being brought into conformance with the already-published gate contract.

## Acceptance criteria

1. At one clean exact tip, eleven focused executions pass: default gate-stop,
   shared gate, and recorded gate on Sonnet and Opus; gate, recorded gate,
   rejection, and keep-moving on Codex; recorded gate on Pi. These cover all
   nine current reds plus the already-green Sonnet/Codex recorded-gate controls.
   A skip is not a pass.
2. Every new deterministic positive has a named counterexample: structural
   closed/wrong gate, nested-cwd uncommitted review, tool-result-only Pi review,
   pre-completed rejection seed, and unphased/planted keep-moving output all
   remain red. This prevents a green signal obtained by oracle weakening.
3. The full registered shared scenario suites pass locally on Sonnet, Opus, and
   Codex, and Pi front-door plus recorded-gate tests pass locally, with model,
   exact git tip, command log/session JSONL, entity post-state, and Go JSON
   retained under distinct local artifact roots.
4. Recorded-gate still proves one qualifying root review, one consumed
   decision, one successor dispatch, and one later durable effect on every
   supported host. Rejection still proves two implementation reports and two
   validations. Keep-moving still proves advance and dispatch for every ready
   entity; Codex's transcript-only path requires successful build, completed
   wait, and durable report, and any claimed finalization must be exact,
   successful, and attached to that phase-bound entity.
5. The measurable release result moves from 0/4 green lane jobs in retained run
   `30257280066` to 4/4 green jobs in exactly one manually dispatched Runtime
   Live E2E run at the same locally proven tip. Prerelease remains blocked until
   that same-tip run is 4/4; CI is confirmation only, never the iteration loop.

## Test plan

First run the deterministic counterexample layer:

```bash
gofmt -w ./cmd ./internal
git diff --check
go test ./internal/gates ./internal/contractlint ./internal/ensigncycle -count=1
go test ./...
go test ./... -race
```

Then run the exact focused live commands, retaining each command's `-json`
output and artifacts. Local auth must be present, and any skip fails AC1:

```bash
mkdir -p live-artifacts/local-proof
git rev-parse HEAD > live-artifacts/local-proof/git-tip.txt
set -o pipefail
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/sonnet-focused" SPACEDOCK_LIVE_MODEL=sonnet go test -json -tags live -count=1 -timeout 40m -run 'TestLiveClaudeSharedScenarios/(gate-guardrail|recorded-gate-lifecycle)$' ./internal/ensigncycle | tee live-artifacts/local-proof/sonnet-focused.jsonl
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/sonnet-default-gate" SPACEDOCK_LIVE_MODEL=sonnet go test -json -tags live -count=1 -timeout 40m -run '^TestLiveDefaultHeadlessStopsAtGate$' ./internal/ensigncycle | tee live-artifacts/local-proof/sonnet-default-gate.jsonl
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/opus-focused" SPACEDOCK_LIVE_MODEL=claude-opus-4-8 go test -json -tags live -count=1 -timeout 40m -run 'TestLiveClaudeSharedScenarios/(gate-guardrail|recorded-gate-lifecycle)$' ./internal/ensigncycle | tee live-artifacts/local-proof/opus-focused.jsonl
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/opus-default-gate" SPACEDOCK_LIVE_MODEL=claude-opus-4-8 go test -json -tags live -count=1 -timeout 40m -run '^TestLiveDefaultHeadlessStopsAtGate$' ./internal/ensigncycle | tee live-artifacts/local-proof/opus-default-gate.jsonl
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/codex-focused" SPACEDOCK_CODEX_LIVE_REQUIRED=1 go test -json -tags live -count=1 -timeout 40m -run 'TestLiveCodexSharedScenarios/(gate-guardrail|recorded-gate-lifecycle|rejection-flow|keep-moving-posture)$' ./internal/ensigncycle | tee live-artifacts/local-proof/codex-focused.jsonl
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/pi-focused" SPACEDOCK_PI_LIVE_REQUIRED=1 go test -json -tags live -count=1 -timeout 40m -run '^TestLivePiRecordedGateLifecycle$' ./internal/ensigncycle | tee live-artifacts/local-proof/pi-focused.jsonl
```

After those pass, run the complete affected local live proof:

```bash
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/sonnet-complete" SPACEDOCK_LIVE_MODEL=sonnet go test -json -tags live -count=1 -timeout 40m -run '^(TestLiveDefaultHeadlessStopsAtGate|TestLiveClaudeSharedScenarios)$' ./internal/ensigncycle > live-artifacts/local-proof/sonnet-complete.jsonl
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/opus-complete" SPACEDOCK_LIVE_MODEL=claude-opus-4-8 go test -json -tags live -count=1 -timeout 40m -run '^(TestLiveDefaultHeadlessStopsAtGate|TestLiveClaudeSharedScenarios)$' ./internal/ensigncycle > live-artifacts/local-proof/opus-complete.jsonl
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/codex-complete" SPACEDOCK_CODEX_LIVE_REQUIRED=1 go test -json -tags live -count=1 -timeout 40m -run '^TestLiveCodexSharedScenarios$' ./internal/ensigncycle > live-artifacts/local-proof/codex-complete.jsonl
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/pi-complete" SPACEDOCK_PI_LIVE_REQUIRED=1 go test -json -tags live -count=1 -timeout 40m -run '^(TestLivePiFrontDoorSmoke|TestLivePiRecordedGateLifecycle)$' ./internal/ensigncycle > live-artifacts/local-proof/pi-complete.jsonl
```

There is no separate spike: retained sessions already prove production gate
serialization, nested-cwd Opus review capability, Pi root assistant text
capability, normal `status --next`, and Codex batched finalization. Only after
all deterministic, focused, and complete local commands pass at the same tip
may one Runtime Live E2E workflow be dispatched for 4/4 confirmation.

## Stage Report: ideation

- DONE: Correct the filed newcomer/regression classification using the archived 6y authority, current grader history, and retained current-run artifacts; rule each failure as conduct, oracle, or harness.
- DONE: Define measurable entity-level acceptance criteria and a bounded repair design that keeps the approved gate/rejection/keep-moving journeys, blocks the prerelease, and requires focused plus complete local Claude/Codex/Pi proof before one CI confirmation.
- DONE: Declare the exact expected file/LOC surface, documentation delta, and test plan; preserve the no-compatibility, no prompt-coaching, no quarantine, no new-harness, and no CI-led-iteration boundary.

### Summary

Ideation resolves the red lanes into two oracle repairs, two harness
repairs, and one Pi conduct repair while preserving every approved journey.
Implementation is bounded to existing surfaces and must earn complete local
Claude, Codex, and Pi proof before the single same-tip CI confirmation that
unblocks prerelease.

## Stage Report: ideation (cycle 2)

- DONE: The keep-moving design must include `internal/ensigncycle/shared_keep_moving_test.go` and a falsifiable planted/unphased negative control because lines 255-260 currently credit literal `merge guard <entity>` without dispatch-build/wait/durable-report proof; update exact surface/LOC and stop boundary without weakening the journey.
  The revised 3-file repair removes raw merge commands as positive dispatch evidence, retains them as a corrected-entity violation, and budgets 11 files at +160/-42 LOC with a 13-file/273-line stop.
- DONE: Make every complete-local command independently failure-preserving (`set -o pipefail` in the same block, or remove the pipeline).
  All four complete-local commands now redirect JSON directly instead of piping through `tee`, so each preserves the `go test` exit independently.
- DONE: Provide exact before/after text for the Pi/host behavior and operator documentation, including `docs/runtime-live-ci.md`'s currently missing `TestLivePiRecordedGateLifecycle` description; revise the doc estimate credibly.
  Exact adapter, local-run, and GitHub-lane replacements are specified above; the operator-doc budget is now +3/-3 and generated/golden impact remains zero.

### Summary

Cycle 2 closes all three staff-review findings, with checklist 3/3 DONE and all
nine failure classifications unchanged. The trustworthy-green boundary still
forbids scenario withdrawal, prompt coaching, quarantine, compatibility
layers, a new harness, CI-led iteration, product edits during ideation, and any
live-host or CI run before implementation earns complete local proof.

## Stage Report: ideation (cycle 3)

- DONE: Add `internal/ensigncycle/shared_keep_moving_negative_test.go`, whose existing `mergeGuardSurface` positive at lines 252-260 must be deleted or converted to the promised host-neutral unphased negative so raw literal `merge guard <entity>` cannot remain positive proof and the suite still passes; recalculate exact expected hand-written LOC and the +2-file/35% stop from a 12-file baseline.
  The existing surface becomes a seven-line assertion inversion while `doneSurface` stays positive; the exact plan is 12 files at +167/-49 LOC, with re-approval above 14 files or 292 changed lines.

### Summary

Cycle 3 closes the remaining test-surface omission and leaves every cycle-2
ruling and boundary unchanged. No product, live-host, or CI work was performed.
