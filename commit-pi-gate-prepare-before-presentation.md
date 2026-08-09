---
title: Commit the Pi gate before presentation
status: ideation
source: "Staff review M3 for test-behavior-completeness, 2026-08-09"
started: 2026-08-09T20:36:19Z
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
worktree:
issue:
pr:
mod-block:
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
---

## Problem

A Pi operator must receive a gate review that is bound to committed state.

The real `gate-guardrail` journey prepares a gate, but Pi can present before a
successful state commit and reread. The command log then has no committed
binding for the review. The operator sees a review that is not durable.

The task starts from strict XFAIL code
`gate-prepare-state-commit-missing`. The repaired journey must prepare, commit,
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

### 3. Replace the Pi TODO with the strict XFAIL binding

In `internal/ensigncycle/shared_live_runner_test.go`, replace the Pi
`gate-guardrail` TODO with the committed strict-XFAIL binding owned by this
entity. The binding must name only
`gate-prepare-state-commit-missing`.

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

### AC-1 — Value: committed Pi review

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

### AC-2 — Strict baseline and XPASS proof

Before the repair, the exact target must run as strict XFAIL with the sole code
`gate-prepare-state-commit-missing`. A setup error, a skipped TODO, a second
semantic code, or a different code must fail the baseline. After the repair,
the same binding must produce XPASS and fail. The binding is then removed, and
the exact target must pass.

Test: the strict-XFAIL classifier and the three exact-candidate runs.

### AC-3 — No presentation before commit and reread

The command log must contain exactly one successful `gate prepare`, followed by
exactly one successful `state commit`, followed by a `state-head` read. No
decision, consume, dispatch, withdraw, status mutation, or archive may occur
after prepare and before presentation. The root review must be one assistant
text block after the read.

Test: `assertRecordedGateHoldLog` plus the Pi transcript and command log from
the exact target.

### AC-4 — The prepared gate remains the same open gate

After presentation, the fixture entity must retain one open attempt with the
same room reference, Briefing digest, and prepared state. It must have no
Resolution, Application, successor, or archive. Presentation must not mutate
the gate.

Test: existing `assertGateHeld` and the reread JSON artifact.

### AC-5 — Existing semantics remain stable

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
  removed lines. Replace the Pi TODO with the strict XFAIL binding and retain
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
