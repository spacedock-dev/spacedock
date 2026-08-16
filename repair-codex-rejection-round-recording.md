---
id: hz2ankag6fk379ssabpv4ckc
title: Repair Codex rejection-round recording in the live rejection flow
status: ideation
source: "Captain directive 2026-08-16 after two same-day codex failures (runs 31915540750 and 31922268382, both FAIL /rejection-flow observed=[rejection-round-missing]) on a journey whose XFAIL c6a336a33 retired on one unbound pass; old owner continue-codex-rejection-after-first-validation is archived done and fixed a different mode"
started:
completed:
verdict:
score: "0.90"
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:hz2ankag6fk379ssabpv4ckc:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:hz2ankag6fk379ssabpv4ckc-backlog-1
              briefing:
                id: briefing:hz2ankag6fk379ssabpv4ckc:backlog:attempt-1:revision-1
                digest: sha256:1a03a0685d3c3ef04d8aa2aae0f3435c94e2d1bbac51d81ffed388bbedbd63cb
                request-digest: sha256:8b3615051af7a7b0d18b2bf2f35a0ec367018010efe709fad15b61515ab89bf6
                room-ref: ./repair-codex-rejection-round-recording/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:hz2ankag6fk379ssabpv4ckc:backlog:1
                briefing: briefing:hz2ankag6fk379ssabpv4ckc:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T03:15:18.533316Z"
                decision: approve
                reason: 'Captain directive 2026-08-16: file it and dispatch on top of the current stack; local live run for the targeted flake, then PR onto the stack'
              application:
                target-stage: ideation
                state: consumed
---

Codex intermittently completes the live rejection flow without recording the rejection round: the FO-side flow reaches the feedback stage but `gate record --round` never runs, so the journey assertion finds no round record (`rejection-round-missing`). Fail-pass-fail across the last three runs proves the c6a336a33 retirement premature for this journey.

Single deliverable (captain correction 2026-08-16, superseding the original two-part
framing): diagnose and repair the behavior at its owning surface, prove with a TARGETED
LOCAL live run of the rejection-flow journey (SPACEDOCK_LIVE_RUNTIME=codex, -run scoped to
the one journey, repeated until the mechanism - not luck - explains the green), and land
as the stack layer whose full matrix proves the journey green with NO XFAIL binding
present. The XFAIL restore is NOT a deliverable: binding-with-owner is the lane-honesty
mechanism for debt that sits unfixed, and this task IS the fix, so adding and removing a
binding within one branch is ceremony.

Contingency, not deliverable: if the fix is parked or fails validation, the fallback is
the registry policy's XFAIL bound to this entity id (hz2ankag6fk379ssabpv4ckc) on the
`rejection-flow` journey.

Captain-directed process: work on top of the current stack (branch from stack #717's tip, spacedock-ensign/rule-superseded-verdict-vocabulary at 4edc82f07); iterate with the targeted local live run; PR onto the stack as the next layer for the full-matrix run.

## Problem

The seed premise is wrong in two ways, and correcting it changes the whole design.

**The round IS recorded. The label lies.** I re-ran the harness's own stream oracle
(`codexRecordedRejectionRound`, copied verbatim from `shared_round_recording_test.go`)
against all three codex `codex-exec.jsonl` streams. It returns `true` for both failing
runs. In both, `spacedock gate record rejection-task --round validation/1 ...` ran and
exited 0 with the complete four-entry round summary, and `--round validation/2` after it.
Nothing about round recording is broken.

What is actually missing is the **prepared gate**. `runClaudeRejectionFlowScenario`
(`internal/ensigncycle/claude_live_runner_test.go:397-403`) carries a codex-only clause:

    if _, _, err := gates.Read(entityPath); err != nil {
        recordedRound = false
    }

`gates.Read` returns `ErrNoGateRecord` ("entity has no gates record") whenever the entity
frontmatter has no `gates:` key — which is exactly the state of an entity whose gate was
never prepared. So "the FO never prepared the gate" is reported under the
`rejection-round-missing` code. The sibling oracle disagrees with it:
`assertRejectionRoundGateBoundary` explicitly *tolerates* `ErrNoGateRecord` and returns
nil. One journey holds two contradictory positions on the same condition, and the
contradiction is invisible because both surface under one code.

**The seed's evidence base is also wrong.** Run 31915540750 has TWO attempts at the same
sha `e9c408bb`: attempt 1 FAILED, attempt 2 PASSED. Run 31922268382 FAILED. So the
observed record is 2 failures and 1 pass at effectively one commit — a stochastic
recovery failure, not a regression, and not bisectable.

### The real mechanism: a dirty-tree dead end, on every host

After the neutral round recorder writes the entity and the round room, the entity is
uncommitted. Two binary surfaces then refuse, and neither refusal is recoverable by an
FO whose durability verb is a no-op:

1. `status --next --json` drops the `needs-preparation` row when the entity is dirty.
   This is deliberate — `TestGateReadinessRejectsDirtyAndMalformedColdReports` pins
   "dirty target bytes are never mechanical proof". The FO sees
   `{"dispatchable":[],"ready_gates":[]}` — an empty scheduler with no error and no
   explanation — and concludes the run is finished.
2. `gate prepare --artifact <entity>` exits 1 with `selected source differs from its
   committed Git object; commit the exact source before preparation`.

And the verb the FO contract names for durability cannot clear either. `«state.commit»(slug)`
is declared in `first-officer-shared-core.md:126` as "record an entity's change durably",
shipped as `spacedock state commit <slug>`. On an **inline** workflow it prints
`Inline workflow — entities live beside the README; nothing to commit to a state checkout.`
and exits 0, leaving the tree dirty (`internal/cli/state_sync.go:64-69`). The command
reference makes the trap explicit: the `--round` form of `gate record` "still needs an
explicit `state commit` afterward" — and on inline that explicit `state commit` does
nothing. Writer and reader disagree: the gate surfaces require a committed entity, and
the only verb authorized to commit it declines.

The two failing runs are the two exits from that dead end:

- **Run 31915540750 attempt 1** — recorded both rounds, ran `state commit` three times
  (each a no-op), never attempted `gate prepare` at all, read the empty scheduler, and
  reported success claiming "with its review gate open". No gate was ever prepared.
- **Run 31922268382** — recorded both rounds, ran `gate prepare` while dirty, got the
  committed-object refusal, ran `state commit` (no-op), then committed with raw git —
  and never retried the prepare. Its final message reports `needs-preparation` as a
  finished state.

**This is not a Codex deficiency.** In the claude-sonnet lane of the same run
31915540750, the FO hit the *identical* `gate prepare` refusal, then ran
`git add && git commit` and retried, and passed. Every host walks into the same trap;
passing is a matter of whether the model happens to retry with raw git. That is the
value chain this entity exists to repair.

### Verified, not asserted

- Rebuilt the codex round oracle from the repo's own regexes and ran it over all three
  streams: `true`, `true`, `true`.
- Reproduced the entire chain locally from a clean inline fixture against the stack-tip
  binary: clean entity → `ready_gates:[{readiness:"needs-preparation"}]`; dirty entity →
  `ready_gates:[]`; dirty entity + `gate prepare` → exit 1 with the exact CI error;
  `state commit` on the dirty inline entity → exit 0, tree still dirty.
- Probed the git seam the fix depends on: `git -C <workflow-subdir> add/commit -- <pathspec>`
  commits path-scoped when the workflow dir is nested inside a repo (exit 0, tree clean);
  a non-repo directory returns exit 128, so the fix needs an explicit work-tree check to
  keep an honest no-op there.

## Proposed approach

Four changes. Only the first changes product behavior; the rest make the lane and its
failures honest.

**1. Make `state commit` durable on inline workflows (the value fix).**
In `runStateCommit`, the `StateInline` branch resolves the entity under the workflow dir
and calls the existing `commitEntityPathsScoped` seam, with no publish — an inline
workflow repo is the user's code repo and must never be pushed by this verb. A workflow
dir not inside a git work tree keeps an honest no-op at exit 0 with a reason naming that.

*Value AC served:* AC-1 and AC-2. *Simplest alternative considered:* add prose telling
the FO to run `git add && git commit` when the workflow is inline. *Why insufficient:*
it hands a mechanical guard back to model discipline at precisely the point where model
discipline is already proven to fail — twice on codex, once on claude, at the same sha.
It must be repeated in `feedback-rejection-flow`, `fo-gate-lifecycle`, and
`fo-dispatch-core`, and it contradicts the project's own priority that the binary owns
mutation guards. No new mechanism is introduced here: no command, flag, readiness value,
or schema field. An existing capability stops lying about having performed its job.

*Alternative also considered and rejected:* teach `status --next` to emit a distinct
"candidate suppressed because dirty" signal so the silence becomes legible. It does not
close the loop — the FO still needs a working verb to clean the tree — and it widens a
machine envelope consumed by three contract documents.

**2. Give the missing-gate condition its own code.**
The codex-only clause reports under a new `rejection-gate-not-prepared` code instead of
`rejection-round-missing`. Deliberately kept codex-only: the contract requires the
prepared gate on every host, but widening the assertion to claude in a repair entity
would enlarge the blast radius on n=1 evidence. Widening is a declarable follow-up.

**3. Carry graded findings' messages into the lane failure.**
`liveGrade` gains a `details` field populated by `gradeLive` from each `gradedErr.msg`,
rendered alongside the codes in `finishLiveScenario`'s fatal. Today those messages are
constructed and discarded, so a CI failure prints a code and nothing else — which is why
diagnosing this took three artifact downloads and a rebuilt oracle.

*Value AC served:* AC-3. *Simplest alternative:* leave it and read artifacts.
*Why insufficient:* the artifacts do not contain the durable end state (it lives in a
`t.TempDir`), so the failing sub-check is not recoverable from them at all.

**4. No XFAIL binding is added.** Per the captain's correction, the lane ships unbound and
the full matrix on the stack PR is what proves the journey green. `c6a336a33` already
removed the codex binding, the reconciliation expectation, and the registry row, so this
requires no change at all: `shared_live_runner_test.go`,
`live_registry_reconciliation_test.go`, and `docs/runtime-live-ci-registry.md` are all
untouched by this entity.

*Risk this accepts, stated plainly:* an unbound lane means that if the repair is
incomplete, the codex rejection-flow failure fails the stack PR's matrix rather than
being absorbed as an expected failure. That is the intended trade — the branch is the fix,
so a red matrix is the correct signal. The contingency if the fix is parked or fails
validation is the registry policy's XFAIL bound to this entity id, applied then and not
before.

### Stack position

The dispatch names the stack tip as `spacedock-ensign/rule-superseded-verdict-vocabulary`
at `4edc82f07`. That branch has since been restacked: its tip is now `64c037b56` with the
same commit subject, and `4edc82f07` is no longer an ancestor. Implementation branches
from the branch tip as it stands at implementation time and lands as the next stack layer.

## Out of scope

Pi's rejection flow (p1 owns it). The registry amendment policy itself.

Two defects found during diagnosis, to be filed separately rather than absorbed here:

- **Journey metrics record an unbound failure as a pass.** `scenarioBehaviorResult`
  (`journey_metrics_live_test.go:181`) sets `Passed: true` and only attaches the real
  grade when the journey carries an XFAIL binding. Both failing runs shipped a metrics
  artifact reading `"outcome": {"status": "passed"}`. This spans every journey and host,
  so it is a metrics-schema entity, not a rejection-flow one.
- **Widening the prepared-gate assertion to all hosts**, per change 2's note.

## Expected surface and tolerance

Five files, roughly 105 insertions, tolerance ±35 and ±2 files. The seed estimated
"contract prose or the journey runner's completion conditions"; the diagnosis puts the
owning surface in the binary instead, which is the material scope change this gate is
being asked to approve.

| File | Change |
|---|---|
| `internal/cli/state_sync.go` | inline branch of `runStateCommit` (~25) |
| `internal/cli/state_commit_test.go` | deterministic chain test (~60) |
| `internal/ensigncycle/claude_live_runner_test.go` | new code + fatal detail (~6) |
| `internal/ensigncycle/claude_runtime_helpers_test.go` | `liveGrade.details` (~6) |
| `docs/site/reference/command-reference.md` | `state commit` doc diff (~8) |

No live-registry file is touched: `shared_live_runner_test.go`,
`live_registry_reconciliation_test.go`, and `docs/runtime-live-ci-registry.md` stay as
`c6a336a33` left them.

**Declared semantic changes.** Command grammar: none. Stored formats: none. Authority:
none — `state commit` gains no new write target beyond the entity paths it already
commits in split-root, and never pushes an inline repo. Runtime behavior: `spacedock
state commit <slug>` in an **inline** workflow changes from an unconditional exit-0 no-op
to a path-scoped commit of that entity in the workflow repo; a non-git workflow dir keeps
the exit-0 no-op with a new reason string. Split-root behavior is untouched. Live-lane
semantic vocabulary gains `rejection-gate-not-prepared`.

## Documentation diff

`docs/site/reference/command-reference.md`, `### state commit`, first sentence:

- Before: "`spacedock state commit <slug>` commits and synchronizes one active or clean archived entity from a split-root state checkout."
- After: "`spacedock state commit <slug>` commits and synchronizes one active or clean archived entity. In a split-root workflow it commits and pushes in the state checkout. In an inline workflow it commits the same entity paths in the workflow's own repository and pushes nothing; if the workflow directory is not inside a git work tree, it reports that and does nothing."

Same section, appended after the `--round` sentence: "In an inline workflow that
`state commit` is what makes the entity's recorded round visible to `gate prepare` and to
the `needs-preparation` scheduler row; both read committed bytes only."

## Acceptance criteria

**AC-1 — In an inline workflow, `state commit` leaves the entity committed and the
blocked gate surfaces reachable.**
Verified by: a Go test driving the real binary over an inline fixture, asserting the whole
chain — after `state commit`, `git status --porcelain` for the entity path is EMPTY;
`status --next --json` yields exactly one `needs-preparation` row; `gate prepare
--artifact <entity>` exits 0 emitting `state=open`. The independent baseline that can move
the wrong way is the porcelain byte count, which is nonzero on today's binary and must
reach zero. Running this test against the pre-fix binary fails at the first assertion.

**AC-2 — The codex rejection-flow journey is green on the stack PR's full matrix run with
no XFAIL binding present for codex on this journey.**
Verified by: the stack full run's `TestLiveCommonRejectionFlow` codex result is a PASS
(not XFAIL, not XPASS), with `git grep liveXFail` over `shared_live_runner_test.go`
showing no codex binding on `rejection-flow` at the commit that run tested. Before that
run, the local repair loop must have reached 5 consecutive targeted greens, each
confirmed by its exec stream showing a `gate prepare` that emitted `state=open` — that is
the mechanism-named evidence, not the pass count alone. The independent baseline that can
move the wrong way is the observed pre-fix rate of 1 pass in 3 runs at one sha; 5
consecutive greens at that rate is ~0.4% by chance.

**AC-3 — A rejection-flow failure names what actually failed.**
Verified by: a unit test feeding an entity with no `gates:` block plus a stream that DID
record the round, asserting the emitted code is `rejection-gate-not-prepared` and NOT
`rejection-round-missing` — it fails today, where that input yields the round code; and a
`gradeLive` unit test asserting each graded finding's message survives into
`liveGrade.details`, which fails if the messages are dropped as they are now.

**AC-4 — The suite stays green.**
Verified by: `go test ./...` plain and `-race`, plus the contractlint live-registry
reconciliation still green with no registry change (this entity adds no binding, so the
reconciliation expectations must be untouched and still pass).

## Test plan

Two layers, because the mechanism is deterministic and the outcome is stochastic.

**Offline (deterministic, cheap, standing).** One Go test in `internal/cli` drives the
real binary across the full causal chain on an inline fixture — dirty entity → `state
commit` → clean tree → `needs-preparation` row → successful `gate prepare`. This is the
regression guard: it fails if any link breaks, and it fails today. Two small unit tests
cover the semantic code and the retained grade details. No new harness, no new fixture
directory.

**Live (measurement, throwaway).** The repair loop is
`SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -run TestLiveCommonRejectionFlow
./internal/ensigncycle/` — one journey, ~6 minutes per run. Run until 5 consecutive
greens or a failure resets the count. Each green is confirmed by reading its exec stream
for the `state=open` prepare, not by the pass alone. The loop is the gating evidence for
proposing the layer; the stack PR's full matrix, unbound, is AC-2's terminal proof. No
live test is added to the standing suite.

If the loop cannot reach 5 consecutive greens, the layer is not proposed as a fix.
Implementation stops and reports, and the contingency XFAIL bound to this entity id is
applied at that point — not preemptively.

**No spike needed beyond what ideation already ran:** the three mechanisms this design
rests on were each exercised above — the round oracle against the real streams, the
dirty/clean scheduler and prepare behavior against the stack-tip binary, and the
path-scoped git commit seam in both the nested-repo and non-repo cases.

## Stage Report: ideation

- DONE: Diagnose the rejection-round-missing mechanism from BOTH failing runs' codex exec streams (gh run download 31915540750 and 31922268382 - the codex artifacts carry codex-exec jsonl); name where the round recording is skipped and why, before designing the fix
  Round recording is NOT skipped: the repo's own `codexRecordedRejectionRound` oracle, rebuilt verbatim, returns true on both failing streams. The skip is `gate prepare`, mislabeled by the codex-only `gates.Read` clause at claude_live_runner_test.go:397-403. Run 31915540750 needed attempt 1's artifact (id 9254966513) — attempt 2 passed and overwrites the default download.
- DONE: Design the two-part deliverable: XFAIL restored bound to this entity id (hz2ankag6fk379ssabpv4ckc) for lane honesty, plus the behavioral fix at the owning surface; the repair loop uses TARGETED local codex live runs (one journey, -run scoped) - design the loop and the N-consecutive-green bar
  Partially superseded: the captain withdrew the XFAIL-restore half mid-stage (binding-with-owner covers debt that sits unfixed; this task IS the fix). Delivered half — owning surface is `runStateCommit`'s inline no-op, not contract prose; N=5, justified against the measured 1-in-3 pre-fix pass rate (~0.4% by chance). Body deliverable and AC-2 amended: AC-2 is now the journey green on the stack full run with no binding present, XFAIL kept as a one-line contingency only.
- DONE: Design against the stack tip (spacedock-ensign/rule-superseded-verdict-vocabulary, 4edc82f07); implementation lands as the next stack layer for the full-matrix run
  That branch was restacked: tip is now 64c037b56, same subject, and 4edc82f07 is no longer an ancestor. Design and all spikes ran against the current tip's binary; recorded in the body.

### Summary

The seed premise was wrong twice: the round IS recorded (proven by re-running the
harness's own oracle over all three streams), and the two "same-day failures" are two
attempts plus one separate run at effectively one sha — stochastic, not a regression. The
real mechanism is a host-neutral dirty-tree dead end: `gate prepare` and the
`needs-preparation` scheduler row both require a committed entity, while `state commit`,
the only durability verb the FO contract names, is an exit-0 no-op on inline workflows.
Claude hit the identical refusal in the same CI run and passed only because it retried
with raw git. Fix moves to the binary; three oracle-honesty repairs ride along; two
further defects are named for separate filing rather than absorbed. Per the captain's
mid-stage correction the XFAIL restore was dropped as a deliverable and kept as a
one-line contingency, which removes three files from the expected surface and makes the
unbound stack full run AC-2's terminal proof — with the accepted risk that an incomplete
repair reddens that matrix rather than being absorbed.
