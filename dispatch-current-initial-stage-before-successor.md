---
title: Dispatch the current initial stage before its successor
status: implementation
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
