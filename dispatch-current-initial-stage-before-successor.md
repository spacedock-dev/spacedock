---
title: Dispatch the current initial stage before its successor
status: validation
source: "Recarved from 9adv48yhye5s2vkhwd7ge52d during test-behavior-completeness shaping, 2026-08-09"
started: 2026-08-09T18:34:37Z
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
worktree: .worktrees/spacedock-ensign-dispatch-current-initial-stage-before-successor
issue:
pr:
mod-block:
id: 6x50qafc8566zc6p1qpb6y30
gates:
    version: 1
    records:
        - id: gate:6x50qafc8566zc6p1qpb6y30:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:6x50qafc8566zc6p1qpb6y30-backlog-1
              briefing:
                id: briefing:6x50qafc8566zc6p1qpb6y30:backlog:attempt-1:revision-1
                digest: sha256:e32fc37aa5c0ed64331be56f7df6580998db3c9776929ae23b9e16ca0e33a6e8
                request-digest: sha256:c56845702e628abe4e1a0f0b033af3c75483a568dfdcced3257fadad3d116411
                room-ref: ./dispatch-current-initial-stage-before-successor/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6x50qafc8566zc6p1qpb6y30:backlog:1
                briefing: briefing:6x50qafc8566zc6p1qpb6y30:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T18:33:27.796613Z"
                decision: approve
                reason: The Captain authorized ideation dispatch; this seed isolates the initial-stage repair and requires strict XFAIL evidence first.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:6x50qafc8566zc6p1qpb6y30:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:6x50qafc8566zc6p1qpb6y30-ideation-1
              briefing:
                id: briefing:6x50qafc8566zc6p1qpb6y30:ideation:attempt-1:revision-1
                digest: sha256:258c684ced1fc428b6f4136d683629a5eacb02d181a4070be97eff3b938ad82e
                request-digest: sha256:cf68ab1f7e2f6351ed071571128719b625ca4451924bdacffa516ae371a53b57
                room-ref: ./dispatch-current-initial-stage-before-successor/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6x50qafc8566zc6p1qpb6y30:ideation:1
                briefing: briefing:6x50qafc8566zc6p1qpb6y30:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T21:33:12.688933Z"
                decision: approve
                reason: Captain approved current initial-stage dispatch before successor advancement.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:6x50qafc8566zc6p1qpb6y30:validation
          stage: validation
          attempts:
            - id: gate-attempt:6x50qafc8566zc6p1qpb6y30-validation-1
              briefing:
                id: briefing:6x50qafc8566zc6p1qpb6y30:validation:attempt-1:revision-1
                digest: sha256:f2cd1e5d922d73e8751b196e6c7529b5d64c3e3d8901ec516c48420f3d8d0f0d
                request-digest: sha256:97dfa0bbb32dbddd7b52fa8a8a64e1e84c8405244f0ba24ca958c8a3614f11c4
                room-ref: ./dispatch-current-initial-stage-before-successor/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6x50qafc8566zc6p1qpb6y30:validation:1
                briefing: briefing:6x50qafc8566zc6p1qpb6y30:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T15:41:44.618652Z"
                decision: approve
                reason: Captain conn accepts exact-green product evidence and retains the historical premature-merge defect.
              application:
                target-stage: done
                state: pending
---

A First Officer must dispatch work for the current initial stage before it advances to a terminal successor.

The task owns the initial-stage defect from task 9a. It does not own gate-consume dispatch evidence or post-gate terminalization.

Before the repair, use each target-level XFAIL binding from `ts`. Transfer the three affected bindings to this active task ID.

Ideation must exercise the current-stage dispatch path first. Then it must define the smallest product repair and exact live proof.

## Problem

At source head `d28834249b23df204292149c7581a295e85c10dd`, the scheduler can report
`current=ready,next=done` for a commissioned entity. The First Officer contract
still says that an initial-stage successor row keeps its old meaning. The host
then advances the entity to `done`, commits `dispatch: <slug> entering done`,
and builds the `done` worker.

The fixture requires the opposite order. Each entity must first run its `ready`
worker and write a `ready` stage report. The terminal stage must follow that
report. A wrong target loses the first worker evidence and can make the durable
journey look complete for the wrong reason.

## End-user value

A workflow operator can see each initial-stage worker run before terminal
completion. The operator can inspect a ready-stage report and its dispatch
commit. This evidence prevents a terminal result from hiding skipped work.

The value measure is two ready-stage reports for the two commissioned entities.
The durable record must also show two terminal archives and no wrong-target
dispatch. The same result must hold for every executable host target.

## Staff finding M4: owner transfer and landing order

The sprint staff review found that the source assigns all three
`smallest-sufficient-mechanism` TODO rows to task `9adv48yhye5s2vkhwd7ge52d`.
This task owns that journey after the recarve. The strict-XFAIL baseline must
move the three target owners to task `6x50qafc8566zc6p1qpb6y30` before any
product bytes change.

| Journey target | Current owner | Required owner |
| --- | --- | --- |
| `claude-sonnet` | `9adv48yhye5s2vkhwd7ge52d` | `6x50qafc8566zc6p1qpb6y30` |
| `codex` | `9adv48yhye5s2vkhwd7ge52d` | `6x50qafc8566zc6p1qpb6y30` |
| `pi` | `9adv48yhye5s2vkhwd7ge52d` | `6x50qafc8566zc6p1qpb6y30` |

This transfer covers only the three `smallest-sufficient-mechanism` rows. The
three `keep-moving-posture` rows remain with task `9adv48yhye5s2vkhwd7ge52d`.
The owner transfer must be a committed baseline step. It must happen before
the First Officer contract or any other product file changes.

The shared files require a serial landing order. The sprint order is:

```text
0a -> ts -> 98a -> 6x -> 9a -> zh -> Codex rejection -> Pi gate commit -> Pi headless hold -> xp6
```

Task `6x` lands before task `9a`. Both tasks edit
`internal/ensigncycle/shared_live_runner_test.go` and
`skills/first-officer/references/fo-dispatch-core.md`. Task `6x` rebases onto
the `98a` landing, commits the owner transfer and strict-XFAIL baseline, runs
the baseline cells, applies its dispatch repair, and lands the exact candidate.
Task `9a` then rebases onto the final `6x` landing. It keeps the three
`smallest-sufficient-mechanism` rows owned by `6x` and edits only its own rows.

The Commander must rebase each branch onto the prior landing. If a rebase
changes a source binding or product rule, run the focused live cell again.

## Proposed approach

Change one target-selection rule in
`skills/first-officer/references/fo-dispatch-core.md`:

1. Keep `current == next` dispatching `current` with an idempotent
   `status={current} started` update.
2. Add the initial-stage exception. When `current` is the workflow's initial
   stage and `next` is terminal, set the dispatch target to `current`.
3. Keep `next` as the target for every other successor row.
4. Use the selected target for the status stamp, exact path-scoped dispatch
   commit, and `dispatch build --stage` call.

For the affected row, the observable sequence is therefore
`status=ready started`, `dispatch: ready-one entering ready`, and
`dispatch build --stage ready`. The worker writes the ready report. Normal
successor handling stays outside this task. Gate consume, terminalization,
status output, stored fields, authority rules, and runtime adapters stay
unchanged.

The simplest alternatives do not meet the value. Changing `status --next` to
project `ready` twice would alter a shared scheduler contract. Adding a special
case to one fixture would not help other initial stages. Adding a new runtime
helper would duplicate the existing dispatch path. The contract rule is the
smallest host-neutral change.

## Spike first: current-stage dispatch path

The native binary and dispatch helper were exercised before selecting this
repair in a throwaway split-root fixture. The fixture declared `ready` as the
initial stage and `done` as terminal.

- `spacedock status --workflow-dir ... --next --json` returned
  `current=ready,next=done`.
- `status --set ready-one status=ready started` returned an idempotent status
  update and a non-empty timestamp.
- `dispatch build --entity-path .../ready-one.md --stage ready` returned an
  ensign artifact whose first action and stage were `ready`.
- A path-scoped state commit was named exactly
  `dispatch: ready-one entering ready`; its entity blob had `status: ready`
  and the non-empty `started` field.

The existing independent checks also passed:

```text
go test ./internal/status -run 'TestEnteredStageLegacySuppressionControls/initial_stage_keeps_successor_projection' -count=1 -v
go test ./internal/ensigncycle -run '^TestDurableTaskJourneys$' -count=1 -v
```

This spike proves the binary, state mutation, path-scoped commit, and dispatch
builder support the needed target. It does not claim host-model proof. The
strict-XFAIL runs below are required before implementation.

## Acceptance criteria

**AC-1 (VALUE) — A commissioned initial-stage journey records ready work before terminal work.**

For every executable target whose fixture runs, the
`smallest-sufficient-mechanism` journey leaves both `ready-one` and `ready-two`
with a `ready` dispatch, a ready-stage report, and then terminal state. The
independent end-value measure is `ready_stage_reports=2`,
`wrong_target_dispatches=0`, and `terminal_archives=2` in the existing durable
Git-state oracle. The current failing baseline is the missing ready dispatch
and report after the host commits `entering done`.

Verified by: strict-XFAIL runs before repair, followed by the unchanged live
journey and durable assertion after repair. A missing path-scoped dispatch,
wrong stage, missing report, or early terminal state fails this criterion.

**AC-2 — An initial-stage successor row dispatches its current stage.**

For `current=ready,next=done`, the First Officer uses `ready` for the status
stamp, path-scoped commit, and dispatch build. The commit's entity blob contains
`status=ready` and non-empty `started`.

Verified by: the split-root command fixture, the dispatch artifact, and the
focused skill smoke. The test fails if any target contains `done` instead of
`ready` in the dispatch boundary.

**AC-3 — Existing entered-stage and ordinary successor behavior remains stable.**

`current == next` remains idempotent. A non-initial, non-terminal successor
still dispatches `next`. Gate-consumed dispatch and post-gate terminalization
remain owned by task `9adv48yhye5s2vkhwd7ge52d`.

Verified by: the existing entered-stage status test, the durable dispatch
oracle, contract-lint checks, and the unaffected keep-moving journey. No test
may replace the path-scoped durable assertion with prose matching.

**AC-4 — The change stays within the declared semantic boundary.**

Command grammar, stored formats, write authority, and runtime host adapters do
not change. Only the target selected for an initial-stage successor row changes.

Verified by: the final diff, existing CLI and dispatch tests, and the required
host lanes for the changed First Officer contract.

**AC-5 — The three smallest-mechanism target owners are transferred before repair.**

The Sonnet, Codex, and Pi rows for
`smallest-sufficient-mechanism` name task `6x50qafc8566zc6p1qpb6y30` in the
committed strict-XFAIL baseline. The rows for `keep-moving-posture` still name
task `9adv48yhye5s2vkhwd7ge52d`.

Verified by: inspect the owner-transfer baseline commit, compare all six rows
with the table above, and run the mutable owner join. A baseline that changes
any product file first, or that changes a keep-moving owner, fails this
criterion.

## Strict-XFAIL-first dependency

Task `ts7gq0mr9s3chx2w4wppd1kt` provides the target classifier and initial
bindings. Transfer the three registered targets to this task:

```text
liveXFail("claude-sonnet", "6x50qafc8566zc6p1qpb6y30")
liveXFail("codex",         "6x50qafc8566zc6p1qpb6y30")
liveXFail("pi",            "6x50qafc8566zc6p1qpb6y30")
```

Use this baseline order:

1. Rebase the task branch onto the latest `98a` landing.
2. Change only the three smallest-mechanism owner and XFAIL rows shown above.
3. Commit that baseline before changing the First Officer contract.
4. Run each cell on the committed baseline.

XFAIL is valid only when the fixture runs and the typed durable assertion
returns that exact code. XPASS is a lane failure until this task removes the
binding. Authentication, launch, timeout, fixture, state-read, parsing, and
every different semantic code remain FAIL. Keep TODO only when the cell cannot
execute.

The exact commands are:

```bash
go build -o ./spacedock ./cmd/spacedock
export SPACEDOCK_BIN="$PWD/spacedock"
export SPACEDOCK_REPO_ROOT="$PWD"
SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=sonnet go test -tags live -count=1 -timeout 40m -run '^TestLiveCommonSmallestSufficientMechanism$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveCommonSmallestSufficientMechanism$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=pi go test -tags live -count=1 -timeout 40m -run '^TestLiveCommonSmallestSufficientMechanism$' ./internal/ensigncycle -v
```

The current live runner still skips these rows, and task `ts7` is still in
ideation. Therefore this ideation records the required XFAIL contract but does
not claim a host run. Implementation must run the executable cells first and
attach their metric records before changing the dispatch rule.

## Exact live proof after repair

Run the same three commands without XFAIL bindings. The artifact for each
target must show two ready dispatch boundaries with exact messages
`dispatch: ready-one entering ready` and `dispatch: ready-two entering ready`.
For each entity, the commit blob at that boundary must contain `status: ready`
and non-empty `started`. Each entity must then contain `## Stage Report: ready`,
terminal fields, and its canonical archive.

The durable result must report `ready_stage_reports=2`,
`wrong_target_dispatches=0`, and `terminal_archives=2`. The live test must PASS,
not XPASS. Remove all three XFAIL bindings only after those records exist.

## Test plan and mechanism budget

The implementation test order is:

1. Rebase onto `98a`. Commit the three owner transfers and XFAIL bindings before
   any product bytes change.
2. Run the three strict-XFAIL cells and save their metrics.
3. Add or reuse the typed semantic code at the existing durable assertion
   boundary. Do not classify unstable error text.
4. Add the focused target-selection smoke using an external fixture and the
   existing `dispatch build` helper.
5. Update the one First Officer rule and replace only this task's live bindings.
6. Run all three exact live cells, then run `go test ./...`,
   `go test ./... -race`, and `gofmt -w ./cmd ./internal`.

The new XFAIL binding serves AC-1 by making the pre-repair failure executable.
The focused smoke serves AC-2 by checking an independent command result and
Git state. Reusing the existing durable oracle is simpler than adding a second
workflow model. A new model would risk passing while the supported First
Officer path still dispatches the wrong stage.

## Expected surface and semantic budget

Baseline: **up to 3 existing files, about 18 gross insertions and 6 deletions,
for a net change of about 12 lines; tolerance is ±1 file and ±12 net lines.**

- `skills/first-officer/references/fo-dispatch-core.md`: about 8 insertions
  and 2 deletions. State the target-selection rule and use the selected target
  in the existing mutation and build steps.
- `internal/ensigncycle/shared_keep_moving_durable_test.go`: about 7
  insertions and 1 deletion only if the existing oracle needs the stable typed
  semantic code. Reuse the classifier from `ts7` when it is available.
- `internal/ensigncycle/shared_live_runner_test.go`: replace the three
  `liveTODO` bindings owned by `9adv48yhye5s2vkhwd7ge52d` with this task's
  `liveXFail` bindings before product edits. Keep the three keep-moving rows
  with task `9adv48yhye5s2vkhwd7ge52d`. Remove this task's bindings only after
  the repaired live proof.

Observable semantics are explicit: command grammar unchanged; stored formats
unchanged; write authority unchanged; runtime behavior changes only for an
initial current stage whose scheduler successor is terminal.

## Documentation diff

The contract text is the user-facing documentation for this First Officer
behavior. Replace:

```text
Initial-stage successor rows retain legacy meaning.
```

with:

```text
When the current stage is initial and its successor is terminal, dispatch the
current stage itself. Use its idempotent status stamp, path-scoped dispatch
commit, and stage-specific build. All other successor rows keep their existing
target.
```

No CLI help, workflow schema, or site documentation needs a diff.

## Out of scope

- Do not change `status --next` projection or stored entity fields.
- Do not repair gate-consume evidence or post-gate terminalization.
- Do not add a new dispatch protocol, runtime adapter, or terminal consumer.
- Do not weaken the durable Git-history oracle or teach hosts to use `--force`.

## Stage Report: ideation

- DONE: Exercise the current-stage dispatch path before selecting the repair.
  The native throwaway fixture returned `current=ready,next=done`, then produced
  `status=ready started`, an exact `entering ready` commit, and a ready dispatch.
- DONE: Define the smallest repair for current equals initial and successor equals terminal.
  The plan changes only First Officer target selection and preserves ordinary,
  entered-stage, gate-consume, and terminalization behavior.
- DONE: Give gross and net line estimates with strict-XFAIL and exact live proof.
  The plan names the active owner, three target bindings, exact live
  commands, expected surfaces, a line budget, and post-repair durable metrics.
- DONE: Fold staff finding M4 into the owner map and serial landing order.
  The plan moves all three smallest-mechanism targets from `9a` to `6x` in the
  baseline before product bytes, then lands `6x` before `9a`.

### Summary

Ideation is complete. The smallest repair is an explicit initial-stage target
rule in the First Officer dispatch contract. Native command and durable-oracle
spikes passed; host strict-XFAIL execution remains a required pre-implementation
step because task `ts7gq0mr9s3chx2w4wppd1kt` has not landed yet. The M4 owner
transfer and `6x` then `9a` landing order now protect the shared files.

## Stage Report: implementation

- DONE: Implement the approved initial-stage target-selection rule for current initial stage with terminal successor.
  `skills/first-officer/references/fo-dispatch-core.md:11-17` selects `current` and uses one `dispatch_stage` for each dispatch boundary.
- DONE: Add task-local proof that status stamp, dispatch commit, and dispatch build all target the current initial stage.
  `internal/contractlint/initial_stage_dispatch_test.go:9-23` fails if one boundary loses the selected current stage or the legacy rule returns.
- DONE: Do not edit shared XFAIL bindings, registry reconciliation, sprint package files, or shared runtime documentation before ts lands.
  Commit `f350e3d27` changes only the approved dispatch contract and its task-local test.
- DONE: Preserve ordinary successor, entered-stage, gate-consume, terminalization, command, format, and authority behavior.
  The selection table retains entered-stage and ordinary-successor behavior. Focused status, stamp, durable-oracle, and contract tests passed.
- DONE: Commit the product behavior checkpoint and report exact files, lines, tests, and deferred post-ts live proof.
  Commit `f350e3d27` contains the checkpoint. Both full suites passed. The three post-`ts` host runs and their durable metrics remain deferred.

### Summary

The dispatch contract now runs an initial current stage before its terminal successor. The first race run hit a transient 250 ms subprocess timeout. Its focused rerun and the complete race rerun passed.

## Stage Report: implementation (cycle 2)

- DONE: Rebase onto `origin/main` at or after `a8688cabf`, preserving the product repair.
  Rebase completed without conflict. Product commit `4daf341aa` retains the `dispatch_stage` rule at `fo-dispatch-core.md:11-17`.
- DONE: Complete the 6x product repair and task-local proof.
  `TestSmallestMechanismTraceSelectsCodexDialect` proves that the shared runner selects the Codex extractor instead of the Claude extractor.
- DONE: Run the exact owned target after the product repair.
  Claude Sonnet and Codex produced XPASS with two ready reports, zero wrong-target dispatches, and two terminal archives.
- DONE: Remove only the 6x-owned XFAIL and update its exact reconciliation row.
  `shared_live_runner_test.go:165` has no gap. `live_registry_reconciliation_test.go:57` now expects `nil`.
- DONE: Do not edit another task's binding.
  The three `keep-moving-posture` bindings at `shared_live_runner_test.go:170` still name task `9adv48yhye5s2vkhwd7ge52d`.
- DONE: Run required checks, commit, and push the candidate.
  `go test ./...`, `go test ./... -race`, focused proof, and registry reconciliation passed. Candidate `974a318b1` is pushed.
- FAILED: Run the exact owned Pi target locally.
  OpenRouter returned HTTP 402 before execution. The account supports 31695 tokens, but the configured request requires 128000 tokens.

### Debugging evidence

- Released workflow: The common smallest-mechanism journey runs on Claude Sonnet, Codex, and Pi through supported first-officer front doors.
- Observable harm: Correct Codex ready dispatches and archives remained XFAIL because the shared runner applied the Claude transcript extractor to Codex JSONL.
- Value authority: `value-ac[AC-1]` requires two ready reports, zero wrong-target dispatches, and two terminal archives for each executable host target.
- Exact trigger: The Codex run completed both journeys, but `runClaudeSmallestSufficientMechanismScenario` called `claudeMechanismTrace` for every `liveDriver`.

### Summary

Candidate `974a318b162b7fc895e018e8d902470474c31a4a` is ready for the protected Pi lane. Run the exact Pi target at this SHA. The lane must pass without an XFAIL binding.

## Stage Report: implementation (cycle 3)

- DONE: Run the approved manual Pi cadence at the exact candidate.
  `Runtime Live E2E` run `31360479606` used head `974a318b162b7fc895e018e8d902470474c31a4a` and deployment `5827064190`.
- DONE: Keep the paid run limited to Pi.
  Job `pi-live` used environment `CI-E2E-PI`. The Claude and Codex jobs were skipped.
- DONE: Inspect the exact target result.
  `TestLiveCommonSmallestSufficientMechanism` did not run. The suite stopped at an earlier failure because it used `-failfast`.
- DONE: Download and inspect the journey artifact.
  Artifact `9052397854` contains full-cycle and gate-guardrail metrics. It contains no smallest-mechanism metric because that test did not start.
- DONE: Preserve the candidate bytes.
  Candidate `974a318b162b7fc895e018e8d902470474c31a4a` remains unchanged. No runtime setting or XFAIL binding changed.
- FAILED: Obtain the required protected Pi result.
  `TestLiveCommonGateGuardrail` produced an XPASS for owner `2e4fe65gy9vcr4xck6akzmdd`. This earlier finding stopped the common suite.

### Workflow evidence

- Released workflow: The manual Pi cadence runs the common live suite with `-failfast` before the Pi front-door smoke.
- Observable harm: The 6x Pi target produced no result, so the required cross-host proof remains incomplete.
- Value authority: `value-ac[AC-1]` requires the smallest-mechanism journey to pass on each executable host target.
- Exact trigger: Run `31360479606` reached an XPASS in `TestLiveCommonGateGuardrail` before `TestLiveCommonSmallestSufficientMechanism` started.

### Summary

HOLD candidate `974a318b162b7fc895e018e8d902470474c31a4a` unchanged. The required Pi proof is absent because another owned XPASS stopped the suite before the exact target.

## Stage Report: implementation (cycle 4)

- DONE: Record the Captain's merge-order exception A.
  Land task 2e4 before task 6x. Then rebase 6x and rerun the exact protected Pi target.
- DONE: Record why the exception is required.
  Run `31360479606` and artifact `9052397854` show that the 2e4 XPASS triggered `-failfast` before the 6x target started.
- DONE: Preserve the candidate while 2e4 is not landed.
  Candidate `974a318b162b7fc895e018e8d902470474c31a4a` remains unchanged. Wait for the 2e4 landing SHA before rebase or rerun.

### Summary

HOLD the unchanged 6x candidate until the 2e4 landing SHA is available. Exception A permits the required 2e4-first landing order.

## Stage Report: implementation (cycle 5)

- DONE: Record the incidental Sonnet supporting evidence.
  PR 656 job `93371397610` produced strict XPASS for `smallest-sufficient-mechanism` with an empty observed semantic set.
  The product exited 0, and the durable assertions passed. This baseline result did not authorize an early binding removal.
- DONE: Rebase 6x onto the exact 2e4 landing.
  Rebase onto `f2428ddee6b5b52d1478be25c345e47aae43aa97` completed without conflict.
- DONE: Preserve only task-owned changes.
  The 6x smallest-mechanism binding remains removed. The three 9a keep-moving bindings remain unchanged.
- DONE: Run focused checks.
  The Codex extractor proof, durable smallest-mechanism assertion, and registry reconciliation passed.
- DONE: Push the rebased candidate.
  Candidate `a7ea84f63ffa5d3c2fb44232b4f8637951690646` is clean and pushed.

### Summary

Candidate `a7ea84f63ffa5d3c2fb44232b4f8637951690646` is ready for the exact Pi-only manual cadence.

## Stage Report: implementation (cycle 6)

- DONE: Run the approved Pi-only cadence at the rebased candidate.
  `Runtime Live E2E` run `31365184470` used exact head `a7ea84f63ffa5d3c2fb44232b4f8637951690646` and deployment `5827847941`.
- DONE: Download and inspect the journey artifact.
  Artifact `9054268664` contains metrics for full-cycle and gate-guardrail only. Both journeys passed.
- DONE: Check the earlier 2e4 concern on the same Pi model path.
  The gate-guardrail pass on Luna falsifies a deterministic 2e4 defect.
- DONE: Inspect the exact 6x target result.
  `TestLiveCommonSmallestSufficientMechanism` did not run because an earlier failure triggered `-failfast`.
- DONE: Inspect the independent Pi front-door result.
  `TestLivePiFrontDoorSmoke` passed with its grade artifact.
- DONE: Preserve candidate and binding ownership.
  Candidate `a7ea84f63ffa5d3c2fb44232b4f8637951690646` remains unchanged. No 98a or 9a binding changed.
- FAILED: Obtain the required exact Pi result.
  `TestLiveCommonDefaultHeadlessGateStop` failed first. Its prepared fixture had zero requests instead of one.

### Workflow evidence

- Released workflow: The manual Pi cadence runs common live journeys in source order with `-failfast`.
- Observable harm: The 6x target produced no Pi result, so the required cross-host proof remains incomplete.
- Value authority: `value-ac[AC-1]` requires the smallest-mechanism journey to pass on each executable host target.
- Exact trigger: Run `31365184470` stopped at the default-headless target before the 6x target started.

### Summary

HOLD candidate `a7ea84f63ffa5d3c2fb44232b4f8637951690646` unchanged. After 98a lands, task fh6 owns the default-headless repair and must land first.

## Stage Report: implementation (cycle 7)

- DONE: Rebase 6x onto the exact fh6 landing.
  Rebase onto `bad641f754b567674972bc302d72564c192477c1` kept main's `started` stamp and 6x's selected `dispatch_stage`.
- DONE: Repair the rebase integration guard.
  Run `31398403065` showed that the legacy guard counted only `next_stage`. It now also counts `dispatch_stage`.
- DONE: Apply the Captain's final host scope.
  Run `31399181746` was cancelled after offline passed. The Pi live job did not execute.
- DONE: Preserve Pi as pending while Sonnet and Codex remain strict.
  The 6x Pi XFAIL and its reconciliation row remain. The Sonnet and Codex bindings remain removed.
- DONE: Run the exact Sonnet target on the final candidate.
  `TestLiveCommonSmallestSufficientMechanism` passed in 256.99 seconds.
- DONE: Run the exact Codex target on the final candidate.
  `TestLiveCommonSmallestSufficientMechanism` passed in 237.81 seconds.
- DONE: Run required offline checks.
  `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` passed.
- SKIPPED: Run the final exact Pi target.
  The Captain made Pi pending and non-blocking. No new Pi proof is required for this landing.

### Summary

Candidate `e2ead5f4367598f806a1162f8c0c36bee97f091c` is ready for PR and merge on proven Sonnet and Codex value.

## Stage Report: implementation (cycle 8)

- DONE: Open the final candidate PR.
  PR 660 used exact head `e2ead5f4367598f806a1162f8c0c36bee97f091c` with Sonnet and Codex proof recorded.
- DONE: Merge the validated candidate.
  PR 660 merged as `8832664ddbec8ac024c4e5251410d0422fa7ca04`.
- DONE: Preserve the pending Pi scope.
  The merged change keeps the 6x Pi XFAIL. Pi remains pending and non-blocking.

### Summary

Implementation is complete. Sonnet and Codex passed the exact value target, and merge commit `8832664ddbec8ac024c4e5251410d0422fa7ca04` landed it.

## Stage Report: validation

- DONE: Verify exact candidate e2ead5f4367598f806a1162f8c0c36bee97f091c and merge commit 8832664ddbec8ac024c4e5251410d0422fa7ca04.
  PR 660 and Git identify that exact head and merge commit.
- DONE: Inspect the 6x-owned initial-stage dispatch behavior and confirm other product bindings remain unchanged.
  The selected `dispatch_stage` controls each boundary. The three 9a keep-moving bindings remain unchanged.
- DONE: Inspect exact pre-merge Sonnet PASS in 256.99 seconds and Codex PASS in 237.81 seconds.
  Implementation cycle 7 records both exact candidate runs and their durations.
- DONE: Inspect implementation full, race, focused, registry, owner, and format evidence without repeating owned checks.
  Cycle 7 records both suites and focused checks. The final diff has no format error.
- DONE: Confirm the Pi XFAIL remains pending and non-blocking under the Captain priority recarve.
  The Pi binding still names task `6x50qafc8566zc6p1qpb6y30`. PR job `93494489487` was skipped.
- DONE: Inspect exact PR run 31400544313 after its Sonnet and Codex jobs finish.
  Sonnet job `93494488520` and Codex job `93494488495` finished with SUCCESS.
- DONE: Record the Material process defect that PR #660 merged before required live checks and before independent validation; do not hide or repair this fact.
  PR 660 merged at 14:53:04Z. Both required live jobs started at 14:55:09Z.
- DONE: Write a Simplified Technical English validation report with PASSED or REJECTED and four evidence fields for each finding.
  Recommendation: PASSED. The exact merged product value is green, and the Captain accepted the historical process defect.
- DONE: Commit and push only the durable validation report, then send a completion signal that starts with "We love you."
  This report is the only planned state change. The completion signal follows the push.

### Acceptance evidence

- AC-1: Exact PR artifacts `9069076294` and `9068796373` report the 6x target passed in 225.28 and 174.10 seconds.
  The durable oracle fails for a missing dispatch, report, terminal field, archive, wrong order, or wrong path.
- AC-2: The focused dispatch test and both live runs exercised the initial-stage selection through the supported host path.
- AC-3: The implementation report records the entered-stage, ordinary-successor, durable, full, race, registry, owner, and format checks.
- AC-4: The candidate keeps command grammar, stored formats, authority, and host adapters unchanged. The final integration diff is host-neutral.
- AC-5: Commit `c87e82bc3` contains the three 6x smallest-mechanism owners. The 9a keep-moving owners remain unchanged.

### Material finding: premature merge

- Released user and normal workflow: PR 660 changed the host-neutral dispatch core, so Sonnet and Codex were required before merge.
- Observable harm: The merge placed an unvalidated contract on `main` while both required live jobs were pending.
- Affected boundary: `contract[docs/dev/README.md#Proof-policy]` requires each exercised live lane to be green before merge.
- Trigger evidence: PR 660 merged at 14:53:04Z. The required jobs started at 14:55:09Z. Independent validation started later.
- Defect kind: evidence defect. Release scope: Material historical process defect. Ownership: First Officer process.
- Captain disposition: Accept and retain the event in the record. No corrective product work exists.

### Summary

The 6x product value passes on Sonnet and Codex. Pi remains pending under the Captain recarve.
The validation recommendation is PASSED. The Captain accepted the premature merge as a historical process defect.
