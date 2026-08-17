---
id: hz2ankag6fk379ssabpv4ckc
title: Repair Codex rejection-round recording in the live rejection flow
status: validation
source: "Captain directive 2026-08-16 after two same-day codex failures (runs 31915540750 and 31922268382, both FAIL /rejection-flow observed=[rejection-round-missing]) on a journey whose XFAIL c6a336a33 retired on one unbound pass; old owner continue-codex-rejection-after-first-validation is archived done and fixed a different mode"
started: 2026-08-16T03:40:08Z
completed:
verdict:
score: "0.90"
worktree: .worktrees/spacedock-ensign-repair-codex-rejection-round-recording
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
        - id: gate:hz2ankag6fk379ssabpv4ckc:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:hz2ankag6fk379ssabpv4ckc-ideation-1
              briefing:
                id: briefing:hz2ankag6fk379ssabpv4ckc:ideation:attempt-1:revision-1
                digest: sha256:8f390134fea3622188df7a425e5be3aac025e34a23d7533bf33d83a9022ff281
                request-digest: sha256:1e620185acd5174e05487ba1a4a77b6bc7ab4a6fc45191f4e28992144458d922
                room-ref: ./repair-codex-rejection-round-recording/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:hz2ankag6fk379ssabpv4ckc:ideation:1
                briefing: briefing:hz2ankag6fk379ssabpv4ckc:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T03:40:01.113881Z"
                decision: approve
                reason: 'Captain 2026-08-16 (dispatch): ideation approved into implementation'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:hz2ankag6fk379ssabpv4ckc:validation
          stage: validation
          attempts:
            - id: gate-attempt:hz2ankag6fk379ssabpv4ckc-validation-1
              briefing:
                id: briefing:hz2ankag6fk379ssabpv4ckc:validation:attempt-1:revision-1
                digest: sha256:3f8b45b6c85479d02f167cdfde464fb564d71f68d351bfcfe4e1f52caac3325a
                request-digest: sha256:7c707d260e92d367d81caba9225e00d358d5a5401466cd7bac538a5542ed9611
                room-ref: ./repair-codex-rejection-round-recording/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:hz2ankag6fk379ssabpv4ckc:validation:1
                briefing: briefing:hz2ankag6fk379ssabpv4ckc:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-17T02:56:44.315111Z"
                decision: approve
              application:
                target-stage: done
                state: pending
pr: local-merge:61dd8e435
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

## Stage Report: implementation

- DONE: Execute the gated design exactly: the binary fix (state commit becomes a real commit on inline workflows), the three grading-honesty repairs, nothing absorbed from the two separately-named defects
  `runStateCommit`'s inline branch commits the entity's commit unit through the existing path-scoped seam and never pushes; outside a Git work tree it stays an exit-0 no-op naming that reason. The prepared-gate condition gets its own `rejection-gate-not-prepared` code; `liveGrade` retains each finding's message; no XFAIL binding is added (`git grep liveXFail` shows only the pre-existing `pi` binding on `rejection-flow`). Neither named defect is touched: `journey_metrics_live_test.go` is unmodified and the assertion stays codex-only. Commits 7627cbf84, 762de7bf3, 571017df3, 3826cdd63 — branch head 3826cdd63 as pushed.
- DONE: The local repair loop is the gate: 5 consecutive targeted codex rejection-flow greens with the mechanism - not luck - explaining them; if the loop cannot reach 5, STOP and report with the contingency XFAIL applied, never propose the layer
  Superseded by FO disposition, recorded here because the checklist text still reads the original bar. I stopped at max 1 consecutive green and reported rather than proposing the layer; the FO then re-scoped the bar to the in-stream mechanism chain and explicitly forbade the contingency XFAIL, because the blocker was two out-of-scope skill defects (routed to task #19), not this layer. No XFAIL was applied. Mechanism chain captured in one CI-config stream: `gate prepare` exit 1 "selected source differs from its committed Git object", then `state commit` exit 0 "Committed rejection-task in the inline workflow repository; nothing pushed", then `gate prepare` exit 0 `state=open`.
- DONE: Base is origin/main (post-stack, post-pre6-stamp); land as a stack PR per the ratified Stacked mode for the unbound full-matrix run that is AC-2's terminal proof
  Order reversed by captain mid-stage: this layer now sits ON TOP of the skill-rewrite layer. Rebased onto `spacedock-ensign/simplify-feedback-rejection-flow` (ba91b333f); the quoting fix was re-applied semantically to the peer's extracted `rejectionRoundEntity` and to `rejectionRoundArtifactArg`, not merged textually. PR #719, stack #720 (#718 bottom, #719 top).

### Test evidence

- `TestStateCommitInlineUnblocksGatePreparationChain` drives the real command surface across the whole chain on an inline fixture. Reverting `runStateCommit`'s inline branch to its no-op reds it at the `state commit` assertion — verified by running the test against pre-fix source, where it failed exactly there while its two precondition assertions passed, proving those refusals are observed rather than tautological.
- `TestStateCommitInlineScopeAndPublishBoundary` pins that the commit is path-scoped and that nothing is pushed. Making the inline branch publish, or widening its pathspec to sweep a dirty sibling, reds it.
- `TestStateCommitInlineCommitsWhateverTheWorkflowDepth` drives both repo shapes, workflow-at-root and workflow-nested. It asserts an EMPTY porcelain rather than only a moved HEAD, because the nested defect left the entity staged and a HEAD-moved or commit-count check would have walked past it. Removing `--relative` from the diff reds the nested case and leaves the root case green — the red-before/green-after layer 4's validator proved, now pinned.
- `TestRejectionUnpreparedGateReportsItsOwnCode` feeds a recorded round plus an entity with no `gates:` block and asserts the code is `rejection-gate-not-prepared` and not `rejection-round-missing`. Routing the unprepared gate back through the round oracle reds it; I confirmed the pre-fix behavior by replaying the old clause against that input, where it emitted `rejection-round-missing`.
- `TestGradeLiveRetainsFindingMessages` asserts each finding's message survives into `liveGrade.details`, paired with its own code. Dropping the msg field or mispairing code and message reds it.
- `TestRejectionRoundRecognizerAcceptsNestedShellQuoting` carries a live-captured command verbatim and covers a quote run on each of the four operands. Narrowing any site's quote run back to `?` reds its case; the wrong-file, wrong-entity and wrong-round controls confirm the relaxation did not widen what counts.
- `go test ./...` plain and `-race` green on the composed tree.

A caution reached me mid-stage that AC-1's `needs-preparation` assertion might depend on the report's review verdict, was then contested, and was finally retracted by its author. Rather than pick a side of a prose dispute I read the predicate: `internal/status/entered_stage.go` gates the row on `hasCompleteStageReport(data, stage) && entityPathCleanInHEAD(path)`, where completeness means every checklist bullet's status token is `DONE` or `SKIPPED` with non-empty text and a non-blank evidence line, plus a Summary. There are no review semantics in it. The retraction was correct, AC-1 is unaffected, and this entity's fixture satisfies the predicate by construction — which its passing test demonstrates rather than assumes.

### Live evidence

Composed-tree loop, 4 runs under CI's pinned model and reasoning effort (`--model gpt-5.6-luna`, `model_reasoning_effort=max`, via a local replica of the CI PATH shim): 2 green. `state commit` reported the inline commit in 4 of 4 runs and the pre-fix no-op string appeared in 0 of 4. The two reds were single non-overlapping codes outside this layer — one `rejection-gate-not-prepared` where the FO ended without preparing, one `implementation-worker-not-dispatched` from the worker-lifecycle selector. `rejection-round-missing` did not fire in any composed-tree run.

Per the FO's revised attribution rules the two reds are signal, not noise, so both are stated rather than filtered: the never-prepares mode still occurs on the composed tree at roughly 1 in 4 even with the operand defect fixed under this layer.

An earlier loop measured the wrong thing: without the shim the live runner resolves a bare `codex` from PATH, a different model and effort from CI's. Those runs are retained as a labeled dataset, not counted as evidence. One of them shelled a stale pre-fix binary and passed anyway, its stream showing 4 pre-fix no-op strings and a `git add && git commit` chained into the same shell call as `gate prepare` — an independent live reproduction of the trap, and the reason the loop now also requires the inline-commit line rather than trusting `state=open` alone.

**Reusable fact for anyone running a local live loop.** The live runner resolves codex with `exec.LookPath("codex")`, and CI does not run the bare CLI: `.github/workflows/runtime-live-e2e.yml` (the codex shim step) puts a wrapper named `codex` first on PATH that injects `--model gpt-5.6-luna -c model_reasoning_effort="max"` into every `codex exec` front door. A local loop without that shim silently measures a different model at a different reasoning effort. Same trap on the binary side: `spacedockBinary` falls back to a `spacedock` on PATH when `SPACEDOCK_BIN` is unset, which is how one loop here drove a stale pre-fix build. A replica shim and the loop driver are preserved beside this entity.

### Durable evidence

The exec streams lived in the agent job's tmp dir, which does not survive cleanup, so the load-bearing ones are copied to `repair-codex-rejection-round-recording/evidence/` in this state checkout: the pre-fix trap reproduction, the mechanism chain, and a composed-tree green (gzipped streams), all four run ledgers, the loop driver, and the CI shim replica. `evidence/README.md` extracts each cited command sequence from those same artifacts so every claim here can be re-derived rather than taken on trust; it names every verb a command carries, because codex chains several into one shell call and a naive reading would misattribute the raw-git escape.

### Declared deviations

- **Surface past tolerance.** Ideation declared 5 files and ~105 insertions, tolerance ±2 files and ±35 LOC. Actual at branch head 3826cdd63 is 7 files and 594 insertions against the stack base `ba91b333f`: file count at the boundary, LOC 4.3x over the ceiling. Product plus docs is 96; test code is 498. The overrun is acceptance-criteria-mandated tests the surface table allocated no lines for — AC-3's two unit tests were folded into two ~6-line rows, and AC-1's chain test was estimated at 60 — plus the routed nested-workflow repair below, which added 11 product and 73 test lines after the figures first reported. It was 512 insertions when the FO ruled keep-all-guards and no trim was requested; the growth since is that one authorized repair, no new files.
- **Scope added by FO authorization.** The recognizer quote-run fix was not in the approved design. It was found during implementation, held unchanged, proposed through the review-finding checkpoint, and applied only after distinct FO authorization. Authorized scope was `['"]?` to `['"]*` in `rejectionRoundArtifactArg` and its two sibling patterns in `commandRecordsRejectionRound`, plus a regression case carrying the captured command bytes verbatim. On the composed tree one sibling had been extracted into the var `rejectionRoundEntity`, so the fix was re-applied to that var rather than merged textually; the peer had deliberately left both siblings at `['"]?` for this edit. **Correction to the authorization's framing:** it was granted on the understanding that CI-model prevalence was unmeasured. It is now measured — the three-character quote run appeared under the shimmed CI config too, in the pre-quoting-fix CI run that graded `rejection-round-missing` on a round recorded with a complete four-entry summary. The defect is in any case a static property of the recognizer and independent of which model produced the command.
- **Bar re-scoped by FO.** The 5-consecutive-greens gate was replaced with the mechanism-chain proof; the contingency XFAIL was explicitly withheld, on the FO's ruling that a binding here would be the ceremony the captain already rejected for a journey whose remaining defects are owned by a live peer.
- **Loop scoring re-scoped by FO.** A green counts only with the in-stream inline-commit line and zero pre-fix no-op strings — `state=open` alone is insufficient, because the pre-fix baseline reached `state=open` through the model's raw-git escape. The validator should re-run the loop under that same bar; `evidence/repair-loop.sh` encodes it.
- **Base reversed by captain.** Landing on top of #718 rather than directly on origin/main.
- **Nested-workflow repair, routed in from layer 4's validation** (implementation-discovered-truth class, FO-authorized, captain ratifies at this gate). The inline fix as first written committed only when the workflow dir happened to BE the repo root. In the ordinary nested shape — a workflow at `docs/dev` inside a larger repo — `git -C <workflow> diff --cached --name-only` answers with repo-root-relative names while the pathspec built from `entityPaths` is workflow-relative, so nothing matched, nothing was selected, and the verb reported "Nothing to commit … already up to date". Its own `git add` had already staged the change, so this was strictly worse than the exit-0 no-op the entity set out to remove: it mutated the index and still claimed there was nothing to do, and it did so in the shape that matters, since a workflow at the repo root is the unusual case. Fix is `--relative` on that diff, which makes git answer in the same terms the pathspecs are written in and that `git -C checkout commit -- <path>` resolves them in. Split-root is untouched by construction — there the checkout IS the worktree root. Commit 3826cdd63.

### Findings routed out, not absorbed

- The skill documented the round recorder without its entity operand, so an FO following it verbatim got `Error: unknown gate flag: validation/1`; the same step said to invoke once and hold on failure, turning a usage error into a terminal hold. Observed in 4 of 4 runs before the skill rewrite, under both codex configurations — near-deterministic, not a flake. Provenance witness: `c355fbe44` (2026-07-23) documents the correct `gate record <entity> --round <stage>/<cycle> ...` form; cited alone, because the intervening commit carries no recorder command line in that file and so is not a prior-correct witness. A naive grep for `gate record` in that skill is misleading — the file's fixture examples do name the entity, and the step line an FO actually follows is the one that lost it. Routed to task #19 and fixed there; the composed skill now names the operand and calls its omission a usage error.
- `implementation-worker-not-dispatched` fires on this journey when the lifecycle selector matches two `Stage Report: implementation` headings, which a two-cycle rejection legitimately produces. Not investigated further; named for separate filing.

### Summary

The value fix landed and is proven live: an inline `state commit` now commits the
entity, which is what lets `gate prepare` and the `needs-preparation` scheduler row
see it, and one captured stream shows the refusal, the commit, and the successful
retry in sequence. The two honesty repairs do their job on the composed tree — a
run that records its round is no longer told it did not, and the one real remaining
condition is named as an unprepared gate instead of hiding under the round code.
One correction to my own work, found by layer 4's validation rather than by me:
the inline commit as first written only worked when the workflow dir was the repo
root, and in the ordinary nested shape it staged the change and then reported
nothing to do. That is the same lying-no-op class this entity exists to kill,
reintroduced by the fix for it, and my live loop could not have caught it because
the harness builds workflows in temp dirs outside any enclosing repo. It is fixed
and pinned against both repo shapes.

This layer is necessary but not sufficient for the journey, and the report should
not be read as claiming otherwise. Three independent defects sat on the path:
the skill documented the recorder without its entity operand, the skill's
invoke-once language turned that usage error into a terminal hold, and the
dirty-tree dead end. The first two are #19's and are fixed under this layer; only
the third was ever mine. Journey-level green therefore belongs to the composed
tree at the stack tip, which is the shipping vehicle since the stack merges
atomically — not to this layer alone. Two further findings are named for separate
filing and no XFAIL binding was added.

## Stage Report: validation

- DONE: Source of truth is the entity's stage report and its amendments, never message threads. Validate on the stack tip 52b94f9aa; this layer is #719 including the nested-state-commit fix (3826cdd63).
  Throwaway checkout at 52b94f9aa; 3826cdd63 confirmed an ancestor (9 layers above it); PR #719 open, head 3826cdd63, base `spacedock-ensign/simplify-feedback-rejection-flow`.
- DONE: AC-1: verify the needs-preparation predicate FROM SOURCE at the tip — note a layer above split it into gatePreparable (initial-stage exception); confirm the non-initial path still requires the complete committed exact-stage report with bullet tokens, which is what this entity's AC names.
  `gatePreparable` (entered_stage.go:28) keys the exception on `stage.initial` ONLY; the non-initial path requires `hasCompleteCommittedStageReport` — latest exact-stage token (gate_extract.go:110, last match wins), every bullet DONE/SKIPPED with non-empty text and evidence line, non-empty Summary, entity clean in HEAD — and the needs-preparation row is emitted only through it (discover.go:254). Chain test green at tip; reverting the inline branch to the no-op reds it exactly at the state-commit assertion (state_commit_test.go:1073) with both precondition refusals passing first.
- DONE: AC-2: journey-level green on the unbound full matrix at the tip. Evidence: tip CI run 31986034831 PLUS the composed local ledgers. Do not run new live loops.
  Run 31986034831 (headSha 52b94f9aa, concluded 2026-08-17 ~02:10Z): codex `TestLiveCommonRejectionFlow` is a real PASS (424s, 17 tests 1 failure — the failure is ac-reanchor), with only the pre-existing `pi` binding on rejection-flow at the tip (shared_live_runner_test.go:109). Its exec stream shows the repaired mechanism: 0 pre-fix no-op strings, round recorded entries=4, "Committed rejection-task in the inline workflow repository; nothing pushed", prepare `state=open`. Composed local ledgers: this entity's evidence/ 2-of-4 green with the inline-commit line in 4/4 and prefix_noop 0/4; `_evidence/zq-team-mode-ac1` codex 3/3, claude-sonnet 3-consecutive-in-5. No new live loops were run.
- DONE: AC-3: honest labels — rejection-gate-not-prepared fires as its own code, liveGrade retains finding messages, rejection-round-missing does not fire on recorded rounds. Spot-verify by the falsifying edits the report claims.
  `TestRejectionUnpreparedGateReportsItsOwnCode` and `TestGradeLiveRetainsFindingMessages` green at tip; both carry clear-either-way controls. Narrowing the artifact-arg quote run to `?` reds the verbatim-bytes regression (shared_round_recording_test.go:660). Details rendering proven in production: the tip run's claude rejection-flow failure prints each code's message (gradeLive at claude_runtime_helpers_test.go:598 → finishLiveScenario fatal). rejection-round-missing fired in 0 of 4 composed-tree runs; a layer above widened the gate-prepared assert host-neutral, which strengthens rather than contradicts the split.
- DONE: Nested-fix deviation (routed from layer 4's validator, FO-authorized, captain-ratification pending at this gate): verify the --relative fix and the porcelain-empty regression both directions, and that the stage report's final-SHA amendment records it.
  Dropping `--relative` reds the nested subtest at state_commit_test.go:1173 with the lying `"no-op … already up to date"` while the repo-root subtest stays green; both subtests green at the tip, porcelain-empty asserted. The implementation report's commit list and the deviation bullet both record final SHA 3826cdd63.
- DONE: Assess every declared deviation: surface overrun keep-all-guards, recognizer quoting fix with MEASURED CI prevalence, loop re-scopes (no XFAIL applied), shim confound datasets preserved, evidence/ dir README re-derivations intact (spot-check one cited stream extraction).
  Surface re-measured `git diff --numstat ba91b333f..3826cdd63`: 7 files / 594 insertions, product+docs 96, tests 498 — matches the declaration digit-for-digit; the overrun is AC-mandated falsifiable tests, keep-all-guards defensible. Quoting fix confined to the authorized sites; wrong-file/entity/round controls hold. `git grep liveXFail` at tip: no codex binding anywhere. Bare-default confound ledger preserved and labeled non-evidence. README re-derivations reproduced from the raw gz streams: mechanism chain byte-exact (refusal → inline commit → state=open, no raw git), trap-stream counts exact (4 pre-fix no-ops, 2 raw-git commits). AC-4: plain suite green at tip incl. contractlint reconciliation; `-race` green with one caveat below.
- DONE: Validation stage report: per-AC verdict with evidence, PASSED or REJECTED recommendation, path-scoped commit, push, signal the FO, stop. No gate preparation or resolution.
  This section; no gate preparation or resolution performed.

### Findings for disposition (recommendation only)

- Material at STACK level, outside this entity's value ACs: run 31986034831's matrix is red — claude-sonnet rejection-flow FAILED, observed [rejection-cycle-line rejection-round-missing] ("`### Feedback Cycles` holds 0 `- Cycle N:` entries"; "selected validation gate does not contain exactly one attempt"). The same stream shows this layer's mechanism WORKING on claude (inline commits 4x, 0 pre-fix strings, round recorded, `state=open`); the FO withdrew the stale round-1 gate, re-prepared fresh, and never appended the Cycle line — FO-discipline/grading-shape conditions in the lane this entity's AC-2 deliberately scoped out. Route to the claude-lane owner. The codex lane's red is ac-reanchor, a different journey.
- Deferred risk: `rejection-round-missing` still labels an attempt-cardinality condition on claude (round was recorded); details now name it, so triage is no longer blind. Promote if it misleads a future triage; candidate follow-up is its own code.
- Polish: `internal/ensigncycle` under `-race` blew the default 10m package timeout once on a loaded machine (plain 572s; race ok at 388s on retry with -timeout 25m). Polish: the pre-quoting CI-config stream that sourced the recognizer test's verbatim capture is not among the preserved gz streams (ledger row + in-test bytes remain).

### Summary

All four ACs verified at the stack tip with reproduced, falsifiable evidence: the predicate from source, the codex journey's real PASS on the unbound full matrix with the repaired chain visible in that run's own stream, the honest-label repairs redding under every claimed falsifying edit, and the nested `--relative` fix red/green in both repo shapes. Recommendation: PASSED. The stack PR's matrix is nonetheless red from two findings outside this entity's surface (claude-lane FO-discipline codes on rejection-flow, codex ac-reanchor), routed above; captain ratification is pending at this gate for the declared deviations (nested repair, surface overrun, FO bar re-scopes with XFAIL withheld).
