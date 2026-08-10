---
title: Commit the Pi gate before presentation
status: validation
source: "Staff review M3 for test-behavior-completeness, 2026-08-09"
started: 2026-08-09T20:36:19Z
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
worktree: .worktrees/spacedock-ensign-commit-pi-gate-prepare-before-presentation
issue:
pr: pr-merge:656
mod-block: merge:pr-merge
id: 2e4fe65gy9vcr4xck6akzmdd
gates:
    version: 1
    records:
        - id: gate:2e4fe65gy9vcr4xck6akzmdd:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:2e4fe65gy9vcr4xck6akzmdd-backlog-1
              briefing:
                id: briefing:2e4fe65gy9vcr4xck6akzmdd:backlog:attempt-1:revision-1
                digest: sha256:4f4d2ebcbf4b234f26a87c87c1435e3bf527ee89713ffcc1a4f64e53b7357948
                request-digest: sha256:8afb36d8c635124e272551c8ee60e2fdbee8481cc110c3271097a8ea111834be
                room-ref: ./commit-pi-gate-prepare-before-presentation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:2e4fe65gy9vcr4xck6akzmdd:backlog:1
                briefing: briefing:2e4fe65gy9vcr4xck6akzmdd:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T20:35:49.543597Z"
                decision: approve
                reason: The Captain authorized shaping and requires end-user value; this task owns a Pi review bound to committed state.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:2e4fe65gy9vcr4xck6akzmdd:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:2e4fe65gy9vcr4xck6akzmdd-ideation-1
              briefing:
                id: briefing:2e4fe65gy9vcr4xck6akzmdd:ideation:attempt-1:revision-1
                digest: sha256:da4249eb26fa723f8276aaf0d18142ca5f38f2436d3d118bac69acf6108363cc
                request-digest: sha256:58cd28b5f1933706111e40bd0a8dbfdc043db980f11fe0bfbdd16e951dbab739
                room-ref: ./commit-pi-gate-prepare-before-presentation/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-09T21:28:30.615994Z"
                reason: The ideation AC labels are not readable by the workflow AC extractor.
            - id: gate-attempt:2e4fe65gy9vcr4xck6akzmdd-ideation-2
              briefing:
                id: briefing:2e4fe65gy9vcr4xck6akzmdd:ideation:attempt-2:revision-1
                digest: sha256:51af8e7b2cd1722194960cde64039dfb663e2b2a205cc131b3a487e8bd43a89d
                request-digest: sha256:771cc01da140a4ca71f0bdff88e7af0e216611981b452406d90338c9906163bc
                room-ref: ./commit-pi-gate-prepare-before-presentation/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:2e4fe65gy9vcr4xck6akzmdd:ideation:2
                briefing: briefing:2e4fe65gy9vcr4xck6akzmdd:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-09T21:33:23.549678Z"
                decision: approve
                reason: Captain approved committed and reread Pi gate state before presentation.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:2e4fe65gy9vcr4xck6akzmdd:validation
          stage: validation
          attempts:
            - id: gate-attempt:2e4fe65gy9vcr4xck6akzmdd-validation-1
              briefing:
                id: briefing:2e4fe65gy9vcr4xck6akzmdd:validation:attempt-1:revision-1
                digest: sha256:a4e4223e0b3707927e0a1a0f0145366f82e8e41f62eebdeb9b9561bc27d60f00
                request-digest: sha256:eece0af7910df93f58e62723cae6bdab2143c8962544d9a9bbaa82e01a6ac515
                room-ref: ./commit-pi-gate-prepare-before-presentation/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:2e4fe65gy9vcr4xck6akzmdd:validation:1
                briefing: briefing:2e4fe65gy9vcr4xck6akzmdd:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T06:17:36.833725Z"
                decision: approve
                reason: Exact candidate f7a86be96 satisfies AC-1 through AC-5; independent validation found no Material finding, and the Captain approved the 2e4-first merge-order exception.
              application:
                target-stage: done
                state: pending
---

## Problem

A Pi operator must receive a gate review that is bound to committed state.

The real `gate-guardrail` journey prepares a gate, but Pi can present before a
successful state commit and reread. The command log then has no committed
binding for the review. The operator sees a review that is not durable.

The task starts from target-level XFAIL under this task ID. The repaired journey must prepare, commit,
reread, and present in that order.

This is a journey repair. A command helper alone does not satisfy it. A prose
change alone does not satisfy it. The real Pi `gate-guardrail` cell must pass.

Reuse the current gate room and commands. Do not change gate storage or command
grammar.

## Spike result

The riskiest unverified mechanism was the real Pi ordering. I exercised it
before selecting a repair.

Durable evidence in
`docs/dev/.spacedock-state/restore-live-evidence-after-completed-repairs.md`
records a real Pi `gate-guardrail` probe at source commit `a929fcb60`. The
probe ran the live fixture. It did not skip. Its semantic result was exactly:
`state commit missing or before the successful gate prepare`. The same evidence
separates this failure from Pi `default-headless-gate-stop`, whose code is
`gated entity is not held at its open validation boundary`.

I also ran the exact target in a detached checkout after removing only the Pi
TODO row. Two low-cost Pi model runs reached the fixture and stopped without a
prepared room. A Mistral Pi run completed the same journey in 66.20 seconds.
The model-dependent results show that the live fixture and command recorder can
run. They also show why a passing model run is not enough without the exact
command-log assertion. OpenRouter credit and stale OAuth failures were setup
failures. They are not XFAIL evidence.

The selected baseline is the durable semantic result above. It is the narrow
failure that this task owns.

## Value

After the repair, a Pi operator sees one review whose room, Briefing, and gate
binding were written and reread from committed state. A later operator can
resume from that state without relying on an uncommitted worktree.

The independent baseline is the recorded command log from the pre-repair live
journey. It has a successful prepare but no later successful commit. The value
measure is the same log after repair: one successful prepare, one successful
state commit after it, and one state-head reread before the root presents.

## Proposed approach

### 1. Make the shared gate sequence executable

Update `skills/fo-gate-lifecycle/SKILL.md` at the existing prepared-gate
boundary. State the literal sequence in one block:

1. Run the existing `gate prepare` command.
2. Require a successful open room and the emitted Briefing.
3. Run the existing `state commit ENTITY` command once.
4. Run the existing `status --read ENTITY --json` command and reread the same
   room, Briefing digest, and open state.
5. Load `spacedock:present-gate` and present once from that reread state.

The text must state that prepare output is not presentation input. A nonzero
commit or a mismatch in the reread stops the run. The root must not present,
decide, consume, dispatch, withdraw, or archive between these commands.

Value served: AC-1, AC-2, AC-3, and AC-4. The shared skill names the one
authoritative lifecycle for every host.

Smallest alternative considered: keep the current sentence that says to commit
and then present. It is insufficient. The live Pi probe already had that
instruction and still emitted no commit. The literal order and stop condition
must be adjacent to the point where Pi acts.

### 2. Bind the Pi root session to the shared sequence

Update `skills/first-officer/references/pi-first-officer-runtime.md` in the
existing Pi runtime binding. Add a `gate.lifecycle` rule that binds the shared
sequence to one root-session boundary:

- After `gate prepare` returns `state=open`, the root runs `state commit`.
- The root reads the same entity and checks the same Briefing digest.
- Only then can the root load and emit the one captain-facing review block.
- A failed commit or reread ends the run at the gate. It cannot be repaired by
  a summary, child output, or a second presentation.

Keep the existing Pi rule that only one root-session assistant text block is a
captain-facing review. This binding adds the commit and reread prerequisite. It
does not add a Pi command or a second gate protocol.

Value served: AC-1, AC-3, and AC-4. It closes the host-specific gap exposed by
the Pi transcript while keeping the common lifecycle authoritative.

Smallest alternative considered: edit only the shared gate skill. It is
insufficient because the current shared skill already describes commit and
presentation order, yet the Pi adapter crossed that boundary. Edit only the Pi
adapter is also insufficient because other hosts would lose the shared source
of truth. Both layers are needed for the host binding.

### 3. Use the Pi target-level XFAIL binding

In `internal/ensigncycle/shared_live_runner_test.go`, use the Pi
`gate-guardrail` target-level XFAIL binding owned by this entity.

Use the existing `runGateStopScenario`, `assertGateHeld`, and command-log
classifier. Do not add a new fixture, command helper, or alternate scenario.
The binding has this lifecycle:

1. Run the exact Pi target and record the sole semantic XFAIL.
2. Apply the skill and Pi adapter repair while the binding remains.
3. Run the exact target again. XPASS must fail the test.
4. Remove the binding in the same repair change.
5. Run the exact target again and require a normal pass.

Value served: AC-1, AC-2, and AC-3. The binding makes the real Pi cell the
release gate. It prevents a helper-only or fixture-only pass.

Smallest alternative considered: leave the Pi row as TODO. It is insufficient
because a skipped cell cannot measure the committed-state value. A broad
expected-error matcher is also insufficient because it could hide the separate
open-boundary failure.

## Acceptance criteria

Each criterion has a named test or artifact. The implementation must meet all
criteria.

**AC-1 (VALUE) — Value: committed Pi review**

On the exact candidate, with the XFAIL binding removed,

```text
SPACEDOCK_LIVE_RUNTIME=pi go test -tags live -count=1 -timeout 40m \
  -run '^TestLiveCommonGateGuardrail$' -failfast ./internal/ensigncycle -v
```

must exit zero. The recorded command log must show a successful prepare, a
successful state commit after that prepare, and a state-head reread after the
commit. The Pi transcript must show the captain-facing review only after that
reread. This measures the value against the pre-repair log, which has no
successful commit.

Test: exact Pi live target plus `assertRecordedGateHoldLog` and the saved live
artifact.

**AC-2 — Target baseline and XPASS proof**

Before the repair, the exact target must run as target-level XFAIL. A setup
error or a skipped target must fail the baseline. After the repair,
the same binding must produce XPASS and fail. The binding is then removed, and
the exact target must pass.

Test: the target-XFAIL classifier and the three exact-candidate runs.

**AC-3 — No presentation before commit and reread**

The command log must contain exactly one successful `gate prepare`, followed by
exactly one successful `state commit`, followed by a `state-head` read. No
decision, consume, dispatch, withdraw, status mutation, or archive may occur
after prepare and before presentation. The root review must be one assistant
text block after the read.

Test: `assertRecordedGateHoldLog` plus the Pi transcript and command log from
the exact target.

**AC-4 — The prepared gate remains the same open gate**

After presentation, the fixture entity must retain one open attempt with the
same room reference, Briefing digest, and prepared state. It must have no
Resolution, Application, successor, or archive. Presentation must not mutate
the gate.

Test: existing `assertGateHeld` and the reread JSON artifact.

**AC-5 — Existing semantics remain stable**

The existing `gate prepare`, `state commit`, `status --read`, gate storage, and
authority rules must remain unchanged. Claude and Codex `gate-guardrail` runs,
offline gate tests, and the full unit suite must keep their current results.

Test: focused gate tests, non-Pi common live runs where available, and the
repository verification commands in the test plan.

## Test plan

1. Run the exact Pi baseline for AC-2 before any repair. Save the command log,
   transcript, entity JSON, and exit result. If the run has an infrastructure
   error, record that separately and do not classify it as XFAIL. The smallest
   alternative is an offline fixture test. It is insufficient for AC-1 because
   it cannot prove Pi command order.

2. Run the existing command-log unit tests for AC-2 and AC-3. Add focused
   assertions only if the current classifier cannot reject a missing commit,
   an out-of-order commit, or a second semantic error. The smallest alternative
   is to inspect a free-form transcript. It is insufficient because transcript
   text does not prove durable state order.

3. Apply the shared skill and Pi runtime binding for AC-1 and AC-3. Run skill
   contract checks and the focused gate tests. The smallest alternative is a
   new command wrapper. It is insufficient because it would create a second
   command grammar and would not force the real Pi journey to use it.

4. Run the exact Pi target with the strict binding still present for AC-2.
   Require XPASS failure. Inspect the command log for the exact prepare,
   commit, and reread sequence. A passing test without this artifact is not
   enough.

5. Remove only the repaired Pi XFAIL binding. Run the exact Pi target again for
   AC-1, AC-3, and AC-4. Require the live artifact and entity assertions. Do
   not remove or weaken the `assertRecordedGateHoldLog` checks.

6. Run the Claude and Codex `gate-guardrail` targets, when their configured
   credentials are available, for AC-5. Run offline gate tests even when a live
   host is unavailable. The smallest alternative is Pi-only proof. It is
   insufficient for AC-5 because the shared lifecycle text serves all hosts.

7. Run `go test ./...`, `go test ./... -race`, and
   `gofmt -w ./cmd ./internal` for AC-5. Inspect the diff and the state commit.

## Expected surface and estimate

Expected implementation surface:

- `skills/fo-gate-lifecycle/SKILL.md`: about 10 added lines and 4 removed
  lines. Add the literal prepare, commit, reread, present sequence.
- `skills/first-officer/references/pi-first-officer-runtime.md`: about 7
  added lines and 1 removed line. Add the Pi root-session binding.
- `internal/ensigncycle/shared_live_runner_test.go`: about 5 added lines and 3
  removed lines. Use the Pi target-level XFAIL binding and retain
  the existing scenario and assertion.

The estimate is about 22 gross additions, 8 gross deletions, and 14 net lines
across 3 files. Tolerance is one file and 12 net lines. A new command helper,
new gate storage format, new fixture, or new scenario is outside this estimate
and requires a new task decision.

The staff estimate was about 20 additions and 8 deletions. The spike adds the
small Pi binding text and keeps the estimate within the same small-change
range.

## Observable semantic boundary

- Command grammar: unchanged. Use the current `gate prepare`, `state commit`,
  and `status --read` commands.
- Stored format: unchanged. Use the current gate room, Briefing, digest, and
  entity state fields.
- Authority: unchanged. The first officer still mutates entity state. The Pi
  root only commits the selected gate binding and presents it to the Captain.
- Runtime behavior: Pi pauses at the same open gate, but presentation now
  follows a successful commit and reread. A failed commit or mismatch stops the
  run.
- User-facing documentation: no new command is exposed. The existing
  `docs/runtime-live-ci.md` already describes the retained gate package and
  stop boundary. It needs no wording change for this internal ordering repair.

## Scope

In scope:

- The shared prepared-gate lifecycle wording.
- The Pi runtime binding for commit and reread before presentation.
- The strict Pi `gate-guardrail` XFAIL, XPASS, and pass sequence.
- Focused proof of command order and unchanged open-gate state.

Out of scope:

- The Pi `default-headless-gate-stop` task. Its open-boundary code has a
  different owner.
- Gate storage or Briefing format changes.
- New CLI grammar or a command helper.
- A new fixture, live scenario, or transcript-only assertion.
- A change to who owns entity mutation or Captain decisions.
- Product implementation in this ideation stage.

## Stage Report: ideation

- DONE: Exercise the real Pi gate-guardrail failure before selecting a repair.
  Durable live evidence records the sole `gate-prepare-state-commit-missing` result; detached live runs also exercised the real Pi target.
- DONE: Define the smallest change that presents only committed and reread gate state.
  The plan uses the existing shared lifecycle, a Pi runtime binding, and the real gate-guardrail strict-XFAIL path. It adds no command or storage protocol.
- DONE: Give a visible-value statement and gross and net line estimates.
  The value is a review bound to committed state. The estimate is 22 gross additions, 8 gross deletions, and 14 net lines across 3 files.
- SKIPPED: Product implementation.
  This stage produces the repair plan and proof contract. No product files were changed.

### Summary

The task now has a live, semantic baseline and a bounded repair that requires
prepare, commit, reread, and presentation in order. The exact Pi journey must
show XPASS while its strict binding remains, then pass after that binding is
removed. The state entity is the only artifact changed in this stage.

## Stage Report: ideation (cycle 2)

- DONE: Exercise the real Pi gate-guardrail failure before selecting a repair.
  The durable semantic baseline remains `gate-prepare-state-commit-missing`; no repair semantics changed.
- DONE: Define the smallest change that presents only committed and reread gate state.
  The three-surface plan and its ordered lifecycle remain unchanged.
- DONE: Give a visible-value statement and gross and net line estimates.
  The committed-state review value and 22 gross additions, 8 gross deletions, and 14 net lines remain unchanged.
- DONE: Restore extractor-visible AC labels without changing a criterion.
  Five criterion labels now use the required bold form. AC-1 retains the `(VALUE)` annotation.

The post-correction command was `spacedock status --read 2e4 --ac-scan`.
Its exact output was:

```text
stage=ideation
ac=AC-1 line=189 unevidenced=false citations=2
ac=AC-2 line=207 unevidenced=false citations=1
ac=AC-3 line=217 unevidenced=false citations=1
ac=AC-4 line=228 unevidenced=false citations=1
ac=AC-5 line=237 unevidenced=false citations=1
```

### Summary

The withdrawn gate's structural issue is corrected. The AC scan now discovers
all five criteria, including the real Pi value criterion. The criteria,
semantic boundary, repair scope, and estimates were not changed.

## Stage Report: implementation

- DONE: Implement the approved Pi prepare, commit, reread, then present
  behavior with the existing gate commands and room.
  Commit `9550735b8` adds the shared sequence at
  `skills/fo-gate-lifecycle/SKILL.md:19-32` and the Pi binding at
  `skills/first-officer/references/pi-first-officer-runtime.md:15,29`.
- DONE: Add task-local proof for one successful prepare, later state commit,
  state-head reread, and unchanged open gate.
  `TestAssertGateHeldAcceptsPreparedFixtureBinding` at
  `internal/ensigncycle/gate_assert_test.go:59-108` runs the real commands.
  Missing state or changed gate identity makes it fail.
- DONE: Do not edit shared XFAIL bindings, registry reconciliation, sprint
  package files, or shared runtime documentation before ts lands.
  `git show --stat 9550735b8` lists only the two task-owned product files and
  the task-local gate test.
- DONE: Preserve command grammar, stored formats, authority, fixtures, and other runtime behavior.
  The candidate uses existing commands and fixtures. Focused gate, CLI, and contract tests passed.
- DONE: Commit the product behavior checkpoint and report exact files, lines,
  tests, and deferred post-ts live proof.
  Product commit `9550735b8` contains 26 additions and 7 deletions across 3
  files. The exact normal and race suites passed.
- SKIPPED: Run the post-ts exact Pi XFAIL, XPASS, and normal-pass sequence.
  The Captain deferred this live proof until ts lands. This checkpoint does not change the shared Pi XFAIL binding.

### Summary

The shared gate lifecycle and Pi adapter now require prepare, commit, reread,
and presentation in order. The focused test proves the durable-state sequence
and unchanged open gate. Post-ts live Pi proof remains deferred.

## Stage Report: implementation (cycle 2)

- DONE: Rebase onto `origin/main` at or after `a8688cabf` while preserving the
  product work from `9550735b8`.
  Rebase completed without conflict. The preserved product change is now
  commit `1b529add2` above base `a8688cabf`.
- DONE: Complete the 2e4 Pi product repair and task-local proof.
  `skills/fo-gate-lifecycle/SKILL.md:19-32` and the Pi runtime binding enforce
  commit and reread before presentation. The real-command test remains at
  `internal/ensigncycle/gate_assert_test.go:59-108`.
- DONE: Do not require baseline-XFAIL ceremony.
  No pre-repair baseline run occurred in this cycle. The current strict binding
  produced XPASS after the product repair.
- DONE: Run the exact owned target and remove only the 2e4-owned XFAIL after
  the product repair passes.
  The bound Pi target produced XPASS in 102.45 seconds. The unbound target then
  passed in 99.91 seconds with `openrouter/moonshotai/kimi-k2.5`.
- DONE: Update the exact reconciliation row without editing another task's
  binding.
  Commit `f7a86be96` changes only the gate-guardrail binding and its desired
  reconciliation row. `TestRuntimeLiveRegistryReconciliation` passes.
- DONE: Run required checks, update the durable implementation report, commit,
  and push.
  `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`
  pass. The branch is pushed through commit `f7a86be96`.

### Summary

The Pi gate-guardrail journey now passes without an XFAIL. The review follows
the committed-state reread and leaves the same gate open. The product branch
and this durable report contain the completed transition.

## Stage Report: validation

- DONE: Inspect exact candidate `f7a86be96` and the implementation report. Do not change candidate bytes.
  The worktree stayed clean at `f7a86be96`. The diff has 28 additions and 9 deletions across five files.
- DONE: Verify every acceptance criterion with independent evidence, including the exact Pi target pass and committed-state reread behavior.
  AC-1 passed in the unbound 99.91-second Pi run. The focused real-command check proved the commit and reread boundary.
  AC-2 used the bound 102.45-second XPASS and the unbound normal pass.
  AC-3 passed the order and forbidden-action controls in `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle`.
  AC-4 passed the gate-identity and terminal-state controls in `TestAssertGateHeldAcceptsPreparedFixtureBinding`.
  AC-5 uses the recorded normal and race suite results. The focused gate and registry checks also passed.
- DONE: Inspect the owned XFAIL removal and exact registry reconciliation. Confirm no other task binding changed.
  The candidate removes only `xfail/pi/2e4fe65gy9vcr4xck6akzmdd` and its matching reconciliation row.
- DONE: Perform the required semantic adversarial pass on gate identity, commit order, authority, and failure cleanup.
  Mutants for a missing commit, duplicate prepare, wrong attempt, wrong digest, mutation, and successor dispatch all failed.
  The sequence keeps one open gate. It adds no Resolution, Application, successor, archive, or status mutation.
- DONE: Use existing implementation full and race results without duplicate owned reruns. Run only independent focused checks that can falsify the claims.
  The focused gate checks and `TestRuntimeLiveRegistryReconciliation` passed. I did not repeat the full or race suites.
- DONE: Classify each finding with the workflow four evidence fields and recommend PASSED or REJECTED.
  No finding exists. I recommend PASSED for exact candidate `f7a86be96`.

### Summary

Exact candidate `f7a86be96` meets AC-1 through AC-5. The independent focused checks found no Material finding.
Three additional Pi attempts had external setup failures. They did not reach a product result and do not change the PASSED recommendation.
