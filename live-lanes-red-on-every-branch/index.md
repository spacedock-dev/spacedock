---
id: se0v37bt7mhsrmhta1nyns0r
title: "Runtime Live E2E has failed on every branch for three days, so no merge since 2026-07-24 carries a live signal"
status: validation
source: "Found attributing PR #571's four red live lanes, 2026-07-27. The change was exonerated by baseline comparison and local reproduction; the lanes were already red on every branch."
started: 2026-07-27T17:05:38Z
completed:
verdict:
score: 0.85
worktree: .worktrees/spacedock-ensign-live-lanes-red-on-every-branch
issue:
gates:
    version: 1
    current:
        gate: gate:se0v37bt7mhsrmhta1nyns0r:validation
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
                id: briefing:se0v37bt7mhsrmhta1nyns0r:ideation:attempt-1:revision-2
                digest: sha256:7b0cf07c49f5286b085dbbd884f45aafb247dddd7a2075d7344010a640b7aad2
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:se0v37bt7mhsrmhta1nyns0r:ideation:1
                briefing: briefing:se0v37bt7mhsrmhta1nyns0r:ideation:attempt-1:revision-2
                by: agent:first-officer
                at: "2026-07-27T17:44:37.747284Z"
                decision: approve
                reason: 'Sprint conn exercised after independent staff approval: nine failures are evidence-classified, the 12-file recovery preserves every journey, all five ACs have falsifiable same-tip local proof, and no material finding remains.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:se0v37bt7mhsrmhta1nyns0r:validation
          stage: validation
          attempts:
            - id: gate-attempt:se0v37bt7mhsrmhta1nyns0r-validation-1
              briefing:
                id: briefing:se0v37bt7mhsrmhta1nyns0r:validation:attempt-1:revision-1
                digest: sha256:9f5025e5b0119f3895e26f7cb072e9eeac8772b99e85413c1668561aa09494d8
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:se0v37bt7mhsrmhta1nyns0r:validation:1
                briefing: briefing:se0v37bt7mhsrmhta1nyns0r:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-29T04:32:30.299796Z"
                decision: revise
                reason: 'Run 30421227237 is 3/4: Sonnet, Opus, and Pi passed; Codex durably consumed the approved gate, dispatched implementation, and archived done, but the keep-moving extractor recognizes only legacy status --set advancement. Revise only that successful-consume oracle and preserve failed-consume/no-advance controls.'
              application:
                action: feedback
                target-stage: implementation
                state: pending
review-round:
    id: round:se0v37bt7mhsrmhta1nyns0r:validation:8
    stage: validation
    cycle: 8
    briefing:
        id: briefing:se0v37bt7mhsrmhta1nyns0r:validation:round-8
        digest: sha256:25020abb7c59295cc831b0a55feeefe1e9ec78b41cdac793ba7515f9540dc7f2
        digest-domain: canonical-bytes
        room-ref: ./review/validation/round-8
pr: "#572"
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

All scenarios remain registered and unweakened. Two defects have narrow,
model-scoped TODOs. `TestLiveDefaultHeadlessStopsAtGate` remains skipped for
Sonnet and Opus under
[q3vpb8hes1b3k3f1jps1kvpk](../gate-record-stage-coherence-guard.md) until
cross-stage gate recording is rejected. The shared
`recorded-gate-lifecycle` journey is skipped only for Sonnet under
[w5bfnrvpcphw857nzz93340c](../reliable-exact-digest-in-gate-review.md) until
Sonnet reliably renders the exact selected Briefing digest. Opus and every
other host/scenario remain active. A Runtime Live E2E job may be green with
these named skips; that green proves neither default gate-stop capability nor
Sonnet exact-digest presentation. There is no compatibility layer, new
standing harness, prompt coaching, fixture prompt change, or CI trigger.
Existing live runners, transcript formats, command log, production gate
model/validator, and negative-test files are sufficient.

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

AC-1 through AC-4 each serve AC-5. Each mechanism criterion re-anchors to the
same-tip value result in AC-5 and is not sufficient when AC-5 is unmet.

- **AC-1** At one clean exact tip, seven focused executions pass: shared gate on
   Sonnet; shared gate and recorded gate on Opus; gate, recorded gate,
   rejection, and keep-moving on Codex. Sonnet `recorded-gate-lifecycle`
   reports only
   [w5bfnrvpcphw857nzz93340c](../reliable-exact-digest-in-gate-review.md) and is
   excluded from the passing count; Pi `TestLivePiRecordedGateLifecycle` reports
   only [9w59t6m1qc46hccd54p04z2j](../pi-delegated-gate-continuation-reliability.md)
   and is likewise excluded because its delegated presentation-to-application/
   dispatch capability remains unreliable. The Sonnet and Opus executions of
   `TestLiveDefaultHeadlessStopsAtGate` are likewise excluded under
   [q3vpb8hes1b3k3f1jps1kvpk](../gate-record-stage-coherence-guard.md) until
   cross-stage gate recording is rejected before re-enabling. Any other skip
   fails AC-1. These named skips may leave a CI job green but do not prove
   any quarantined capability.
- **AC-2** Every new deterministic positive has a named counterexample: structural
   closed/wrong gate, nested-cwd uncommitted review, tool-result-only Pi review,
   pre-completed rejection seed, and unphased/planted keep-moving output all
   remain red. This prevents a green signal obtained by oracle weakening.
- **AC-3** The full registered non-TODO shared scenario suites pass locally on
   Sonnet, Opus, and Codex, and Pi front-door passes locally. Pi recorded gate
   reports only its 9w skip. Complete Sonnet may report only its
   `recorded-gate-lifecycle` w5 skip; complete Sonnet and Opus may report only
   their top-level q3 gate-stop skips; complete Pi may report only 9w. Model,
   exact git tip, command log/session JSONL, entity post-state, and Go JSON are
   retained under distinct local artifact roots.
- **AC-4** The exact ID+digest recorded-gate oracle remains unchanged. Same-tip
   live evidence proves one qualifying root review, one consumed decision, one
   successor dispatch, and one later durable effect on Opus and Codex. Sonnet
   exact-digest presentation remains explicitly unproven while w5 is open, and
   Pi delegated presentation-to-application/dispatch remains explicitly
   unproven while 9w is open; neither skip can count as evidence. Rejection
   still proves two
   implementation reports and two validations. Keep-moving still proves
   advance and dispatch for every ready entity; Codex's transcript-only path
   requires successful build, completed wait, and durable report, and any
   claimed finalization must be exact, successful, and attached to that
   phase-bound entity.
- **AC-5 (VALUE)** The measurable release result moves from 0/4 green lane jobs in retained run
   `30257280066` to 4/4 green jobs in exactly one manually dispatched Runtime
   Live E2E run at the same locally proven tip. Prerelease remains blocked until
   that same-tip run is 4/4; CI is confirmation only, never the iteration loop.
   While the linked TODOs remain, 4/4 may include the named q3, w5, zbc, and 9w
   skips and is not evidence that any quarantined capability works.

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
output and artifacts. Local auth must be present. Any skip other than the
Sonnet and Opus executions of `TestLiveDefaultHeadlessStopsAtGate`, the Sonnet
`recorded-gate-lifecycle` selection, and Pi `TestLivePiRecordedGateLifecycle`
fails AC-1. Those four skip events must report their linked TODOs and are
excluded from the passing count:

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

After the seven required focused executions pass, the Sonnet recorded-gate
selection reports only w5, Pi recorded gate reports only 9w, and the two named
gate-stop invocations report only q3, run the complete affected local live proof:

```bash
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/sonnet-complete" SPACEDOCK_LIVE_MODEL=sonnet go test -json -tags live -count=1 -timeout 40m -run '^(TestLiveDefaultHeadlessStopsAtGate|TestLiveClaudeSharedScenarios)$' ./internal/ensigncycle > live-artifacts/local-proof/sonnet-complete.jsonl
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/opus-complete" SPACEDOCK_LIVE_MODEL=claude-opus-4-8 go test -json -tags live -count=1 -timeout 40m -run '^(TestLiveDefaultHeadlessStopsAtGate|TestLiveClaudeSharedScenarios)$' ./internal/ensigncycle > live-artifacts/local-proof/opus-complete.jsonl
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/codex-complete" SPACEDOCK_CODEX_LIVE_REQUIRED=1 go test -json -tags live -count=1 -timeout 40m -run '^TestLiveCodexSharedScenarios$' ./internal/ensigncycle > live-artifacts/local-proof/codex-complete.jsonl
SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/live-artifacts/local-proof/pi-complete" SPACEDOCK_PI_LIVE_REQUIRED=1 go test -json -tags live -count=1 -timeout 40m -run '^(TestLivePiFrontDoorSmoke|TestLivePiRecordedGateLifecycle)$' ./internal/ensigncycle > live-artifacts/local-proof/pi-complete.jsonl
```

Complete Opus must pass every shared scenario and may skip only the named
top-level q3 test. Complete Sonnet must pass every non-TODO shared scenario and
may skip only its shared recorded-gate w5 selection plus the named top-level q3
test. Complete Pi must pass its front-door smoke and may skip only its
recorded-gate 9w selection. Their green exits do not prove any quarantined
capability while the linked TODOs remain.

There is no separate spike: retained sessions already prove production gate
serialization, nested-cwd Opus review capability, Pi root assistant text
capability, normal `status --next`, and Codex batched finalization. Only after
all deterministic tests, seven required focused passes, the four exact named
TODO skips, and complete non-TODO local scenarios pass at the same tip may one
Runtime Live E2E workflow be dispatched for 4/4 confirmation.

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

## Stage Report: ideation (cycle 4)

- DONE: Label all five acceptance criteria as `**AC-1**` through `**AC-5 (VALUE)**` and explicitly state that AC-1 through AC-4 serve AC-5 so mechanism/value re-anchoring is machine-visible and reviewable without changing semantics.
  Read-only `status --read … --stage ideation --ac-scan --json` enumerates exactly AC-1, AC-2, AC-3, AC-4, and AC-5 with evidence; deleting any bold heading drops that ID, while the serve relation re-anchors all four mechanisms to the same-tip 4/4 value.

### Summary

Cycle 4 makes the five approved criteria visible to the deterministic gate
scanner and states their existing mechanism-to-value relationship explicitly.
No classification, implementation surface, LOC budget, command, boundary,
product, live-host, or CI behavior changed.

### Feedback Cycles

- Cycle 1: REJECTED — Roborev panel job 3188; surface 13 files/246 LOC vs estimate 12 files/216 LOC (114%); AC unchanged
- Cycle 2: REJECTED — Roborev panel job 3191; surface 13 files/285 LOC vs estimate 12 files/216 LOC (132%); AC unchanged
- Cycle 3: REJECTED — Roborev panel job 3194; surface 14 files/287 LOC vs estimate 12 files/216 LOC (133%); AC unchanged
- Cycle 4: REJECTED — Roborev panel job 3197; surface 14 files/293 LOC vs estimate 12 files/216 LOC (136%); AC unchanged
- Cycle 5: REJECTED — Roborev panel job 3218; surface 16 files/369 LOC vs estimate 12 files/216 LOC (171%); AC unchanged
- Cycle 6: REJECTED — Roborev panel job 35 and captain-approved design reset; surface 18 files/811 LOC vs estimate 12 files/216 LOC (375%); AC unchanged
- Cycle 7: REJECTED — Roborev panel job 85 and captain-approved narrow reset; surface 18 files/799 LOC vs estimate 12 files/216 LOC (370%); AC amended: 9w is named non-proof
- Cycle 8: REJECTED — exact-tip PR run 30378538074 and delegated sprint-conn design reset; surface 19 files/803 LOC vs estimate 12 files/216 LOC (372%); AC unchanged

## Stage Report: implementation

- DONE: Restore structural gate, recorded-gate, rejection, and keep-moving
  grading without weakening their deterministic counterexamples.
- DONE: Retain exact-tip live evidence for every non-TODO focused journey and
  quarantine Pi delegated gate continuation under
  [9w59t6m1qc46hccd54p04z2j](../pi-delegated-gate-continuation-reliability.md)
  after two classified exact-tip failures.
  `final-ce4ac943-pi-recorded-gate-pinned` reached durable dispatch but presented
  the wrong Briefing digest in 203.39s;
  `final-ce4ac943-pi-recorded-gate-pinned-retry` presented the exact ID/digest
  but stopped before application/dispatch in 93.94s. Prior passing classification
  history remains retained in `se0-b25-focused/pi-recorded-gate-final-tree-retry`
  (201.02s) and `cycle2c-pi-recorded` (352.52s). The 9w skip permits the lane to
  run but does not prove the quarantined capability.
- DONE: Keep `TestLivePiFrontDoorSmoke` active, preserve the production
  gate/oracle surfaces, mark the registered/selected recorded-gate runner as an
  explicit non-proof gap, and stay at 19 changed code files within the
  captain-approved 811-line ceiling.
  Final code candidate `55b3b13316571f1a84215f489c205e66327faeca` is
  19 files/803 changed lines; all deterministic, full, race, and vet checks
  pass, and Roborev `branch_final` re-panel job 90 reports no issues.

### Summary

The implementation restores trustworthy non-TODO live grading at the exact
candidate tip. Pi front-door support remains active; only delegated
presentation-to-application/dispatch is quarantined under 9w, with both red
exact-tip traces and the prior green history retained for the follow-up.

## Stage Report: validation

- DONE: Independently reproduce the final deterministic, race, vet, formatting,
  and diff checks at candidate 55b3b133.
  At exact `55b3b13316571f1a84215f489c205e66327faeca`, focused gates/contractlint/ensigncycle,
  `go test ./...`, `go test ./... -race`, `go vet ./...`, exact `gofmt -w ./cmd ./internal`,
  and `git diff --check` all passed without changing tracked bytes.
- FAILED: Cross-check AC-1 through AC-5 against retained exact-tip evidence, with
  q3/w5/zbc/9w skips treated as explicit non-evidence and no other skip accepted.
  Exact-tip PR run `30378538074` exposed three non-TODO material evidence defects;
  q3/w5/zbc/9w remain explicit non-evidence and are not used to excuse them.
- FAILED: Verify PR #572's required path-derived CI lanes genuinely run green, the
  19-file/803-line boundary holds, and no task-owned Material finding remains.
  The boundary is exactly 19 existing files/803 lines (713 additions, 90 deletions),
  but jobs Sonnet `90340453203` and Codex `90340453205` failed; only Opus
  `90340453222` and Pi `90340453132` passed, and three Material findings remain.
- FAILED: AC-1 — all required non-TODO focused executions pass at one exact tip.
  Sonnet `rejection-flow` and Codex `recorded-gate-lifecycle`/`keep-moving-posture`
  failed outside the four named quarantine skips.
- DONE: AC-2 — every deterministic positive retains a named red counterexample.
  Focused replay passed structural wrong/closed gate, nested-cwd/unordered review,
  pre-completed rejection, and missing/planted Codex phase controls.
- FAILED: AC-3 — the complete affected exact-tip live suites are green.
  Sonnet and Codex shared suites failed; Opus shared/team and Pi front-door passed.
- FAILED: AC-4 — the exact lifecycle and phase-bound live evidence is validly graded.
  The run produced correct-looking durable evidence, but three grader boundaries
  rejected it, so the required behaviors lack valid exact-tip evidence.
- FAILED: AC-5 (VALUE) — move retained `30257280066` from 0/4 to 4/4 green lanes.
  Run `30378538074` reached only 2/4 green at the exact PR head.
- FAILED: Material evidence defect — Sonnet rejection round cardinality.
  The final message reports two implementations and two validations, but
  `assertRejectionFlow` substitutes rejection-only `### Feedback Cycles` count for
  validation count and false-reds the valid second PASSED validation.
- FAILED: Material evidence defect — Codex recorded-review extraction.
  A failed pre-bind decision invocation is followed by a valid post-bind exact
  ID/digest review and consumed lifecycle, but the extractor returns at the earlier
  invocation because the enclosing shell later exits zero.
- FAILED: Material evidence defect — Codex keep-moving durable-report extraction.
  The exact trace has successful build, completed wait, three durable numbered
  `rg -n` reports, and exact successful finalizations; the parser accepts only an
  unnumbered heading. This satisfies the prior merge-guard finding's recorded
  promotion trigger and makes it Material now.
- SKIPPED: q3, w5, zbc, and 9w quarantined capabilities.
  Their exact TODO skips remain accepted non-proof only; none counts toward an AC.
- SKIPPED: Combined Pi 15-minute timeout risk.
  Still deferred: Pi front-door passed in 132.37s and no supported timeout occurred.
- DONE: Validation recommendation: REJECTED.
  Roborev jobs 35 and 85 were dispositioned and re-panel 90 was clean, but the
  later exact-tip CI artifacts supply new trigger evidence for the three findings.

### Summary

Candidate-local deterministic proof and the 19-file/803-line boundary are clean,
but the release value remains unmet at 2/4 live jobs. Validation rejects on three
task-owned Material evidence defects surfaced by the single exact-tip PR run; no
product or test file was modified during validation.

## Stage Report: implementation (cycle 8 recovery)

- DONE: Repair Sonnet validation-report cardinality from exact run 30378538074.
  `assertRejectionFlow` now requires two anchored durable validation reports
  instead of two rejection-only Feedback Cycle entries. The retained Sonnet
  stream reconstructs its two implementation and two validation reports and
  passes; the one-validation controls remain red.
- DONE: Repair Codex recovery after the exact failed pre-bind decision batch.
  The extractor skips the completed zero-exit item whose output begins with the
  exact no-gate error and contains a state head, then accepts the later committed
  briefing, exact root review, and final decision. Existing failed-commit,
  ordering, duplicate-review, nested-root, and decision-before-review controls
  remain red.
- DONE: Repair Codex numbered `rg -n` durable Stage Report recognition.
  Anchored headings may carry a numeric line prefix. The retained keep-moving
  trace proves all three entity reports after build and wait; planted,
  unphased, wrong-entity, and missing-finalization controls remain red.
- DONE: Remove Roborev-only hardening and advisory state from rounds 9–16.
  The implementation is six files and 46 changed lines over branch tip
  `55b3b133`; validation round 8 is preserved as the current all-fixed advisory
  record. No Pi, workflow, documentation, keep-moving product, or live-host
  behavior changes remain in this recovery patch.
- DONE: Complete the required local gate on the bounded candidate.
  Exact retained Sonnet and Codex artifact replay, paired focused regressions,
  `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`,
  `go test ./... -race`, and `go vet ./...` all pass.

### Summary

The candidate is reconstructed to the three observed run-30378538074 oracle
repairs only and is ready for independent validation. No CI, push, merge, or
additional advisory panel was triggered after the bounded cut.

## Stage Report: validation (cycle 2)

- DONE: Validate exact candidate `c34941ec6d80dd2b97321815ef5b2ca740db6210`.
  Its sole parent is `55b3b13316571f1a84215f489c205e66327faeca`;
  the bounded diff is six existing `internal/ensigncycle` test/grader files,
  30 additions and 16 deletions. The code worktree has no tracked changes after
  validation; only the pre-existing untracked `live-artifacts/` tree remains.
- DONE: Replay the retained Sonnet run-30378538074 rejection evidence without
  substituting Feedback Cycles for validation rounds.
  The retained stream's durable heading read lists implementation at lines 26
  and 49 and validation at lines 34 and 61, while listing one Feedback Cycles
  section at line 57. The repaired assertion counts the two anchored
  implementation headings and two anchored validation headings. Its
  one-validation and one-implementation controls remain red.
- DONE: Replay the retained Codex failed-pre-bind gate sequence with an exact
  recovery boundary.
  Artifact item 11 is a completed zero-exit decision batch whose output begins
  `Error: entity has no gates record` and includes `state-head`; item 13 then
  binds and commits the exact Briefing, item 14 presents one qualifying root
  review with exact id/digest, and item 15 records and consumes the decision.
  Recovery accepts only that exact failed-pre-bind shape. Failed commit,
  review-before-bind, decision-before-review, nested/child review, and
  duplicate review controls remain red; broader error prose cannot satisfy the
  exact prefix-plus-state-head predicate.
- DONE: Replay the retained Codex keep-moving phases and numbered durable reads.
  Successful dispatch builds at items 27–30 precede completed wait item 32;
  item 33 reports `40:`, `13:`, and `13:` numbered anchored Stage Report
  headings for the three implementation entities; successful item 36
  finalizes those same entities after an intentionally failed item 35.
  The focused phase-bound test accepts this build→wait→reports→finalization
  shape. Report-before-build/wait, failed-build, planted command/output,
  anonymous or wrong/missing entity, failed command, and missing finalization
  controls remain red.
- DONE: Confirm Roborev jobs 106–127 and their advisory state rounds are absent
  from the candidate.
  The exact six-file diff contains neither Roborev text nor those job ids, the
  commit reflog has no intervening round commits, the durable entity remains
  bound to `validation:8`, and the state checkout was clean before this report.
  No Roborev command was run.
- DONE: Run the required local gate at the exact candidate.
  Focused retained-shape replay/regressions, `gofmt -w ./cmd ./internal`,
  `git diff --check`, `go test ./...`, `go test ./... -race`, and
  `go vet ./...` all completed successfully. Gofmt changed no tracked bytes.
- DONE: Cross-check AC-1 through AC-5 without promoting quarantines.
  AC-2 is satisfied locally by the named red counterexamples. The repaired
  graders validly recognize the retained AC-4 behaviors and are ready for a
  same-tip live confirmation. AC-1, AC-3, the same-tip live portion of AC-4,
  and AC-5 remain pending because no Runtime Live E2E run exists at this
  candidate and retained run `30378538074` remains 2/4. The q3, w5, zbc, and
  9w skips remain explicit non-evidence and satisfy none of those pending
  criteria.
- SKIPPED: Trigger or inspect a new PR CI run, push code, run Roborev, or merge.
  Those actions are outside this bounded independent-validation assignment.
- DONE: Validation recommendation: PASSED for pre-CI readiness.
  No material defect remains in the three bounded replay repairs. The candidate
  is ready for one normal PR Runtime Live E2E confirmation; this verdict does
  not claim the release value or full live-host proof until that run is 4/4.

### Summary

The bounded cycle-8 reconstruction passes independent local validation at
`c34941ec6d80dd2b97321815ef5b2ca740db6210`. Retained Sonnet and Codex evidence
now grades on the exact durable and ordered facts while paired false-positive
controls stay red. Release value remains pending one normal same-tip 4/4 PR CI
run, with all four named quarantines still non-evidence.

## Stage Report: implementation (zbc quarantine extension)

- DONE: Extend only the existing zbc rejection-flow model quarantine to Sonnet.
  PR run `30412397240` completed two Sonnet validations but rebound
  `briefing:rejection-task:validation:round-1` as the final gate; the existing
  `assertReviewRoundRecorded` oracle correctly rejected that advisory round as
  a gate attempt.
- DONE: Keep the quarantine explicit non-evidence under the existing owner.
  `TODO(zbcj98qfwtax61vxdzrf615e)` now covers Claude Opus and Sonnet only;
  Codex, Pi, other Claude scenarios, deterministic round/gate coverage, and
  product semantics remain unchanged.
- DONE: Add focused model-scope coverage before implementation.
  The live-tagged predicate test fails without the zbc helper and proves Sonnet
  and Opus aliases skip while Haiku and unrelated names remain active.
- DONE: Run the required local verification gate.
  Focused live-tagged predicate tests, full tests, race, vet, gofmt, and diff
  checks are recorded with the separate code and state commits.

### Summary

The newly observed Sonnet failure is quarantined under its existing zbc owner
without weakening the oracle or broadening any other scenario. The lane remains
explicit non-evidence until zbc supplies distinct post-rework Briefing binding.

## Stage Report: implementation (Opus coordinated gate decision repair)

- DONE: Replay the exact Opus root review from run `30412397240`, job
  `90451410538`.
  The retained review now fails before the repair and qualifies after it; changing
  the coordinated `record ... and consume` decision clause makes the replay red.
- DONE: Select that review only in its durable ordered interval.
  The existing Claude extractor replay binds a successful `state-head`, presents
  the exact root text, then invokes `gate record --decision`; moving, nesting,
  duplicating, or omitting the review remains covered by existing red controls.
- DONE: Keep negated and unrelated decision prose non-actionable.
  A directly paired `record ... and do not consume` control stays red, alongside
  the existing unrelated downstream-noun and negated-action controls.
- DONE: Limit the repair to the exact supported coordinated clause.
  Commit `f75f482e` changes two existing `internal/ensigncycle` test files by
  35 additions and one deletion; generic grammar, product packages, prompts,
  fixtures, gate semantics, TODOs, and the Sonnet zbc commit are unchanged.
- DONE: Run the required local verification gate.
  Focused tests first reproduced both failures, then passed; gofmt, diff check,
  `go test ./...`, `go test ./... -race`, and `go vet ./...` all passed.
- SKIPPED: Push, trigger CI, run Roborev, or merge.
  Those actions remain outside this task-owned exact-CI repair.

### Summary

The Opus lane had completed the durable gate lifecycle, but the extractor
discarded its truthful pre-decision review because one exact coordinated clause
fell outside the bounded predicate. The retained replay now qualifies without
broadening generic decision grammar and is ready for independent validation.

## Stage Report: validation (cycle 3)

- DONE: Validate exact candidate `f75f482e04a0a79268d6f6bb29a4c4247b3f0622`.
  The two post-run commits change only three existing `internal/ensigncycle`
  test/grader files: 53 additions and 3 deletions over `c34941ec`.
- DONE: Confirm the Sonnet run-30412397240 disposition is conduct, not an oracle
  false positive. The retained result records REJECTED validation/1, rework, and
  PASSED validation/2, then fails because the stale validation/1 Briefing became
  the final gate attempt; the semantic oracle correctly stayed red.
- DONE: Confirm `5a2c06e` extends only the existing Claude zbc quarantine.
  The live-tagged scope test proves Sonnet and Opus family aliases skip while
  Haiku and unrelated model names remain active. Other Claude scenarios,
  deterministic assertions, and the pre-existing Codex zbc quarantine are
  unchanged; every zbc skip remains explicit non-evidence.
- DONE: Confirm the existing zbc owner carries the new trigger.
  `bind-post-rework-briefing-at-rejection-regate.md` pre-existed this cycle and
  now names run `30412397240`, Sonnet's stale round-1 binding, and both affected
  Claude families without weakening its distinct-post-rework-Briefing ACs.
- DONE: Confirm the retained Opus review and lifecycle facts.
  The root message at `2026-07-29T01:02:15.257Z` names the task/stage, exact
  Briefing id and digest, approve recommendation, and coordinated decision ask;
  later commands successfully record, commit, consume, dispatch, and retain the
  successor effect.
- DONE: Confirm `f75f482e` recognizes only the supported coordinated clause.
  The generic actionable-option regex is unchanged; one exact normalized prefix
  admits the coordinated record-and-consume clause. Exact ordered replay passes,
  while negated, unrelated, vague, failed, misordered, nested, duplicate, and
  batched controls remain red.
- DONE: Run the exact-tip local verification gate.
  Focused scope/replay tests, gofmt, diff check, full tests, race, and vet all
  completed green; gofmt changed no tracked bytes.
- DONE: Confirm the full PR quarantine and state boundary.
  q3, w5, zbc, and 9w remain explicit non-evidence; no Roborev 106–127
  hardening returned. State changes touch only the existing zbc owner and this
  existing task body, with no new task and no review-round/frontmatter change;
  the durable pointer remains `validation:8`.
- SKIPPED: Modify or push code, run Roborev, trigger CI, or merge.
  These remain outside this independent validation assignment.
- DONE: Validation recommendation: PASSED for pre-CI readiness.
  The candidate is ready for exactly one updated normal PR Runtime Live E2E
  run. Run `30412397240` was at `c34941ec`, so no full live/value claim is made
  for `f75f482e` until that updated run confirms it.

### Summary

Both completed-run dispositions are bounded and falsifiable at the exact
candidate tip. The Sonnet conduct failure is quarantined as non-evidence under
its existing owner, and the Opus false-negative is repaired without broadening
generic grammar or lifecycle semantics.

## Stage Report: implementation (Sonnet live-cycle cleanup)

- DONE: Diagnose the sole Sonnet failure in Runtime Live E2E run `30419214076`
  from retained evidence before changing code.
  Job `90472359633` passed the journey assertions at `03:21:46Z`; its descendant
  transcript continued until `03:23:05Z` and observed cleanup deleting `.git`
  and README before recreating `make-it-work.md`.
- DONE: Classify the cleanup failure as harness-owned or product-lifecycle-owned
  using the four-field journey test.
  The supported journey reached durable terminal state; observable harm was a
  false-red proof/temp writer; the affected boundary is process-exit/clean-state
  evidence; the trigger was launcher-only `SIGKILL` after the durable barrier.
- FAILED: If harness-owned, implement and locally prove only the smallest
  deterministic cleanup correction; if product-owned, preserve evidence and
  stop for a separate decision.
  Commit `90243bf7` drains the existing watcher to zero exit before TempDir
  cleanup, but both focused Sonnet attempts stopped pre-journey on the expired
  local OAuth benchmark token (`401`), so real-live proof remains pending.
- DONE: Preserve the existing se0 CI-greening boundary.
  One existing harness file changed by 13 additions and two deletions; scenario
  assertions, prompts, product packages, launcher semantics, and TODOs are
  unchanged.
- DONE: Run the deterministic verification gate.
  Tagged drain/exit tests, gofmt, diff check, `go test ./...`,
  `go test ./... -race`, and `go vet ./...` all passed.
- SKIPPED: Run Roborev, push code, trigger CI, merge, or create a task.
  Those actions remain outside this bounded implementation cycle.

### Summary

The sole red is a harness cleanup race, not a released lifecycle defect. The
smallest correction now keeps the fixture alive until the resident launcher and
Claude exit cleanly; independent live validation is still required because the
local isolated credential expired before either focused journey could start.

## Stage Report: implementation (Sonnet live-cycle cleanup proof)

- DONE: Prove commit `90243bf7` with the focused real Sonnet live cycle after
  refreshing the isolated OAuth snapshot.
  The repository-documented Keychain refresh stayed silent; the exact live test
  passed on `claude-sonnet-5` in 156.99s at tip `90243bf7`.
- DONE: Preserve the exact focused journey evidence.
  Session `75ba325e-40fb-4623-bb3a-66c206a6b3ba` archived `make-it-work.md`,
  then the new drain observed FO exit code zero; either a surviving process or
  nonzero exit would have failed before TempDir cleanup.
- DONE: Report readiness for independent validation.
  `TestLiveEnsignCycle` passed without the prior `TempDir RemoveAll` failure, so
  the clean-exit correction now has real Sonnet evidence.
- DONE: Keep this continuation proof-only.
  Code remained exactly at `90243bf7`; the live result did not falsify the fix,
  so no implementation, scenario, product, or launcher semantics changed.
- SKIPPED: Run Roborev, push code, trigger CI, merge, file tasks, or widen scope.
  These actions remain explicitly outside this continuation.

### Summary

Fresh isolated OAuth allowed the previously blocked real Sonnet proof to run.
The supported journey terminalized, archived, and exited cleanly, confirming the
harness correction is ready for independent validation.

## Stage Report: validation (cycle 4)

- DONE: Verify exact candidate `90243bf7` fixes the Sonnet TempDir cleanup race
  without changing journey or product semantics.
  The sole-parent diff over pushed tip `f75f482e` is one existing harness file,
  `internal/ensigncycle/live_test.go`, at +13/-2; prompts, journey assertions,
  product packages, launcher behavior, and TODO boundaries are unchanged.
- DONE: Verify the retained real Sonnet `TestLiveEnsignCycle` pass proves clean
  launcher and Claude exit at the exact candidate.
  The committed implementation record identifies tip `90243bf7`, model
  `claude-sonnet-5`, session `75ba325e-40fb-4623-bb3a-66c206a6b3ba`, archived
  `make-it-work.md`, exit code zero, and the 156.99s passing result.
- DONE: Inspect the exact process-exit assertion adversarially.
  After the durable terminal barrier, the existing watcher drains until process
  exit or kills on three minutes of silence; the new poll then requires exactly
  `(exited=true, code=0)` before post-state reads and TempDir cleanup.
- DONE: Run focused deterministic tests plus gofmt/diff check, `go test ./...`,
  `go test ./... -race`, and `go vet ./...`.
  Tagged exit/drain progress, clean-exit, stall-kill, and transcript-retention
  controls passed; the focused gates/contractlint/ensigncycle layer, full suite,
  race suite, and vet passed; formatting and diff checks changed no tracked bytes.
- DONE: Classify the supported journey and affected proof boundary.
  The supported lifecycle terminalizes and archives correctly; the observed harm
  was a false-red Sonnet lane after success, affecting AC-3/AC-5 process-exit and
  clean-fixture evidence when the durable barrier preceded Claude's terminal
  ceremony. This is a material evidence defect with the exact observed CI trigger,
  not an outcome defect.
- SKIPPED: Repeat the owned expensive real Sonnet live run.
  The retained exact-tip pass establishes model, journey terminal state, archive,
  clean process exit, and absence of the prior TempDir failure, so another
  token-spending run would add no missing proof.
- DONE: Recommend PASSED without modifying the candidate.
  No material, deferred-risk, or polish finding remains in this bounded cleanup;
  code stayed exactly at `90243bf7` and only the required validation report changed.
- SKIPPED: Run Roborev, push or modify code, trigger CI, merge, file tasks, or
  widen scope.
  These actions remain explicitly outside this validation assignment.

### Summary

PASSED. Exact candidate `90243bf7` waits for the resident launcher and Claude to
exit cleanly before fixture cleanup, while preserving the supported journey and
all product semantics; retained real Sonnet evidence and fresh deterministic,
race, formatting, diff, and vet checks establish the bounded correction.
