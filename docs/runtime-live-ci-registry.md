# Desired runtime live-test registry

This registry defines the runtime behavior that live CI must prove. It is the
desired-state component for `docs/runtime-live-ci.md`. That operating guide
remains the single normative entry point for runtime live CI.

The operating guide links this registry and names the executable reconciliation.

A journey can be registered before its test, runner, or fixture exists. Missing
implementation is a reconciliation result. It does not weaken the desired state
recorded here.

## Semantics

- A **common journey** is required on every supported runtime target by default.
- A common journey lists an exception only when the behavior is genuinely not
  applicable to that runtime. Missing support, cost, quarantine, and an unwired
  selector are not exceptions.
- Each common journey has one canonical exported `TestLiveCommon...` entry
  point. Runtime-specific launch, authentication, output, and liveness
  behavior stays behind the runtime adapter.
- Each journey references one or more stable fixture IDs and describes their
  semantic setup. Source annotations, not this registry, bind fixture IDs to
  concrete builders.
- Multiple fixtures can prove the same journey when they preserve the same
  intention and assertion while varying topology or starting state. Multiple
  journeys can share a fixture. An annotated fixture ID referenced by no
  registered journey or runtime-specific proof is orphaned.
- Runtime-specific proofs cover only behavior unique to one host substrate. They
  do not duplicate common workflow semantics.
- A live test intentionally excluded from CI is not release evidence and must be
  listed under **Non-gating live experiments** with a reason.

## Supported runtime targets

| Target | Live lane |
|---|---|
| Claude Sonnet | `claude-live` matrix role `sonnet` |
| Claude Opus | `claude-live` matrix role `opus` |
| Codex | `codex-live` |
| Pi | `pi-live` |

The common-journey entries below omit coverage lists. Unless an entry contains an
explicit **Exceptions** field, all four targets are required.

## How to use this registry

### Add or change a common journey

1. Add or update one entry under **Common journeys** before relying on the test as
   release evidence. State the observable outcome, not the implementation plan.
2. Give the journey one stable ID. Its canonical executable path is
   the exported `TestLiveCommon...` function on every supported runtime target.
3. Reference every fixture variant that must prove the same outcome. For each
   fixture, state its stable ID and distinguishing semantic setup.
4. Add the matching source bindings to the shared scenario and concrete fixture
   builder. Builder symbols and paths belong in source, not in this registry.
5. Make every supported live lane invoke the common suite.
6. If one runtime cannot run the journey, leave the desired registry unchanged.
   Reconciliation must report the missing implementation or invocation.

### Decide whether behavior is common or runtime-specific

Put behavior in **Common journeys** when its required outcome is meaningful for
more than one runtime. This rule applies even if only one runtime implements the
behavior today. Adapters absorb differences in transport, authentication, tool
names, and output formats.

Put behavior in **Runtime-specific live proofs** only when the outcome is about a
unique host boundary. Examples include Claude's merged Agent shape and Pi's
package and child-agent substrate. Name the one lane that owns the proof. Explain
the host-specific outcome.

### Add or reuse a fixture

- Reuse an existing fixture ID when its setup contract is already sufficient.
- Add another fixture to the same journey when it varies topology or starting
  state but preserves the journey's intention and assertion.
- Create a separate journey when the required outcome or grader changes.
- Give every fixture builder a `//spacedock:live-fixture id=<fixture-id>` source
  binding. A bound fixture referenced by no registry entry is orphaned. Renaming
  or moving the builder keeps the stable fixture ID and does not require copying
  its new symbol or path into this file.

### Keep a live test outside CI

An intentionally unselected live test belongs under **Non-gating live
experiments**. Record why a negative result is useful data rather than a release
failure. Otherwise promote it to a registered journey or runtime proof, move its
deterministic coverage to the default suite, or delete it.

## Common journeys

### `full-ensign-cycle`

- **Entry point:** `TestLiveCommonFullEnsignCycle`
- **Required outcome:** A normal workflow member is dispatched, completed,
  validated, terminalized, and preserved with durable state evidence.
- **Fixtures:**
  - `realistic-lifecycle` — a discoverable multi-stage workflow with one member
    ready for execution.

### `gate-guardrail`

- **Entry point:** `TestLiveCommonGateGuardrail`
- **Required outcome:** The first officer binds and presents the retained review
  package, then stops without deciding, advancing, dispatching, or archiving.
- **Fixtures:**
  - `recorded-gate/held` — a member parked at a human gate with retained gate
    inputs.

### `default-headless-gate-stop`

- **Entry point:** `TestLiveCommonDefaultHeadlessGateStop`
- **Required outcome:** A headless launch without decision authority dispatches
  and completes the preceding-stage worker, then presents the first human gate
  and stops open.
- **Fixtures:**
  - `recorded-gate/pre-gate` — the held-gate workflow with the member starting in
    the preceding implementation stage.

### `withdrawn-gate-recovery`

- **Entry point:** `TestLiveCommonWithdrawnGateRecovery`
- **Required outcome:** The first officer preserves a withdrawn attempt, prepares
  and commits its successor, presents the successor, and stops without decision
  or dispatch.
- **Fixtures:**
  - `recorded-gate/withdrawn` — a prepared gate whose first attempt is withdrawn
    without rewriting its retained room, ready for one open successor attempt.

### `recorded-gate-lifecycle`

- **Entry point:** `TestLiveCommonRecordedGateLifecycle`
- **Required outcome:** Delegated authority is bound, recorded, committed, and
  consumed exactly once before successor dispatch.
- **Fixtures:**
  - `recorded-gate/prepared` — a retained prepared gate with command logging and a
    dispatchable successor.

### `rejection-flow`

- **Entry point:** `TestLiveCommonRejectionFlow`
- **Required outcome:** A rejected candidate is corrected and independently
  checked before a fresh final gate is presented. Rejected authority cannot
  satisfy the final approval.
  The journey runs in team mode: the correction is routed to the producer and
  the re-review to a worker that did not produce the fix, through whichever
  route the host's reuse conditions leave available.
- **Fixtures:**
  - `rejection/before-validation-1` — a candidate entering its first validation
    with a deliberate defect and a two-cycle correction path.

### `feedback-3-cycle-escalation`

- **Entry point:** `TestLiveCommonFeedbackThreeCycleEscalation`
- **Required outcome:** A third consecutive rejection is escalated to the captain
  instead of being routed into a fourth automatic correction cycle.
- **Fixtures:**
  - `rejection/before-validation-3` — a member with two prior correction cycles
    and a third rejected validation report.

### `merge-hook-guardrail`

- **Entry point:** `TestLiveCommonMergeHookGuardrail`
- **Required outcome:** The first officer refuses terminalization while a
  registered merge hook remains unsatisfied.
- **Fixtures:**
  - `merge-hook/blocked` — a terminal-ready member with a registered unsatisfied
    merge hook.

### `filing`

- **Entry point:** `TestLiveCommonFiling`
- **Required outcome:** The first officer creates a seed through the atomic
  supported filing path rather than previewing an ID and hand-writing state.
- **Fixtures:**
  - `filing/empty-workflow` — a commissioned workflow with no members and a
    configured ID style.

### `shallow-boot`

- **Entry point:** `TestLiveCommonShallowBoot`
- **Required outcome:** Startup identifies and reports held workflow state without
  mutation, dispatch, team creation, or eager deferred-module work.
- **Fixtures:**
  - `boot/held-gate` — a discoverable workflow with one member already held at a
    human gate and no engage workload.

### `zero-discovery`

- **Entry point:** `TestLiveCommonZeroDiscovery`
- **Required outcome:** Startup reports that no managed workflow exists and stops
  without broad filesystem discovery or team creation.
- **Fixtures:**
  - `boot/no-workflow` — a Git repository containing no commissioned workflow.

### `auto-continue-after-implementation`

- **Entry point:** `TestLiveCommonAutoContinueAfterImplementation`
- **Required outcome:** After observing a completed implementation report, the
  first officer advances to validation, dispatches a fresh validator, and
  presents the validation gate, leaving it open. The fixture's validation stage
  is `gate: true` and the runbook grants no conn, so a run that reaches `done`
  resolved a human gate nobody approved: that is a failure, not a success.
- **Fixtures:**
  - `auto-continue/single-root` — a single-root workflow parked at completed
    implementation with worktree-backed validation.
  - `auto-continue/split-root` — the same invariant in a separate state checkout
    with non-worktree validation.
- **Evidence:** Every fixture variant, on every runtime, must pass the shared
  `assertAutoContinue` durable-state assertion — end state `status: validation`
  only, with `done` red under code `human-gate-bypassed` — and the unconditional
  dispatch-evidence check: an open validation gate, a committed
  `## Stage Report: validation`, and one fresh validator dispatched and completed
  before `gate prepare`.

### `self-evidence-merge-triage`

- **Entry point:** `TestLiveCommonSelfEvidenceMergeTriage`
- **Required outcome:** The first officer refuses unsupported terminalization and
  diagnoses the current run's evidence instead of trusting an inherited label.
- **Fixtures:**
  - `merge-triage/unapproved-live-evidence` — a terminal candidate whose required
    live lane is unapproved and whose current failure differs from the inherited
    diagnosis.

### `smallest-sufficient-mechanism`

- **Entry point:** `TestLiveCommonSmallestSufficientMechanism`
- **Required outcome:** The first officer performs directly authorized work
  directly. The first officer dispatches commissioned ready work without an
  unnecessary workflow, worker, PR, or per-member justification.
- **Fixtures:**
  - `mechanism-choice/mixed-authority` — commissioned ready members alongside
    directly writable deterministic work.

### `keep-moving-posture`

- **Entry point:** `TestLiveCommonKeepMovingPosture`
- **Required outcome:** Approval triggers immediate advancement. Independent work
  proceeds concurrently. Async dispatch does not end the turn prematurely. A
  correction pauses only the affected member.
- **Fixtures:**
  - `keep-moving/mixed-events` — independent ready members, a gate approval, and a
    questioned correction path.

### `ac-value-reanchor`

- **Entry point:** `TestLiveCommonACValueReanchor`
- **Required outcome:** A gate rejects mechanism-only success when the value
  criterion that mechanism serves has regressed.
- **Fixtures:**
  - `ac-reanchor/means-pass-value-regressed` — a gated candidate whose mechanism
    criterion passes while its paired value criterion regresses.

### `owned-conflict-owner-handoff`

- **Entry point:** `TestLiveCommonOwnedConflictOwnerHandoff`
- **Required outcome:** After aborting an owned moving-target conflict, the first
  officer routes reconciliation to the workflow owner of the existing registered
  checkout. Live reuse and ordinary fresh fallback preserve the owner tuple from
  the initial stamped dispatch without changing workflow authority.
- **Fixtures:**
  - `conflict-owner/stamped-checkout` — an initial stamped dispatch creates the
    registered implementation checkout before its branch and `main` diverge into
    a real rebase conflict.

## Runtime-specific live proofs

These proofs are intentionally separate from common journeys. Each must remain
limited to the named runtime boundary.

### `claude-merged-agent-dispatch`

- **Entry point:** `TestLiveMergedTeamModeDispatch`
- **Lane:** `claude-live`
- **Required outcome:** Current Claude dispatches a named background Agent without
  legacy `TeamCreate`, records the host-native member identity, and completes the
  durable workflow.
- **Fixture:** `realistic-lifecycle`, shared with `full-ensign-cycle`.

### `claude-bare-dispatch`

- **Entry point:** `TestLiveBareReachable`
- **Lane:** `claude-live`
- **Required outcome:** Explicit bare mode uses the supported bare dispatch shape
  without loading break-glass recovery, and the assigned durable result is committed.
- **Fixture:** `dispatch-recovery/base` — one dispatchable member under a workflow
  that exercises the supported Claude dispatch boundary.

### `claude-dispatch-build-break-glass`

- **Entry point:** `TestLiveBreakGlassShimRecovery`
- **Lane:** `claude-live`
- **Required outcome:** A real `dispatch build` failure reaches recovery, preserves
  the selected blocking-bare or named-background-team dispatch shape, and commits the
  complete worker result and parsed Stage Report in a path-scoped clean commit.
- **Expected failure:** `selected-bare` is flaky on `claude-live`; owner
  `060xp69y61yhrww23g3wvwqy`.
- **Fixtures:**
  - `dispatch-recovery/base` — the shared dispatchable workflow.
  - `dispatch-recovery/failing-build` — the same workflow with `dispatch build`
    failing while every other launcher command remains available.

### `pi-front-door-subagent-dispatch`

- **Entry point:** `TestLivePiFrontDoorSmoke`
- **Lane:** `pi-live`
- **Required outcome:** The current checkout loads through the supported Pi front
  door. It dispatches an ensign through the supported child-agent substrate. It
  leaves durable worker output and boot-contract evidence.
- **Fixture:** `pi/split-root-smoke` — a current-checkout Pi environment and
  split-root workflow with one child-dispatchable member.

## Non-gating live experiments

These tests are intentionally not release evidence and are not selected by a
live CI lane. Each remains live-tagged only for its stated experiment.

### `TestLiveHaikuLoopSpike`

- **Purpose:** Run one exploratory Haiku-only workflow loop and preserve its
  durable and stream observations.
- **Fixture:** `haiku-loop/experimental` — a low-cost multi-stage workflow used
  only for the loop experiment.
- **Reason unselected:** It measures whether an experimental low-cost loop holds.
  A failure is research output, not a release-gating regression.

### `TestLiveHaikuLoopSpikeN`

- **Purpose:** Repeat the same exploratory loop at least three times and report
  which signals break across repetitions.
- **Fixture:** `haiku-loop/experimental`, repeated into isolated artifact roots.
- **Reason unselected:** Its any-break matrix is an experiment whose negative
  result is expected data, not a failing release assertion.

### `TestLiveCodexWaitMatrixFromShippedAdapter`

- **Purpose:** Characterize Codex wait behavior for active, complete, failed,
  and absent worker states through the shipped adapter.
- **Fixture:** `codex/wait-matrix` — four isolated worker-state variants.
- **Reason unselected:** The four real Codex runs measure host behavior. They do
  not prove a common user journey or block a release.

## Source binding convention

Each exported `TestLiveCommon...` declaration carries an immediately adjacent
`//spacedock:live-journey id=<id> fixture=<fixture-ids>` annotation and a single
`liveJourney(...)` call. That call binds real builder, TODO-owner, exercise, and
assertion symbols. Fixture builders carry adjacent `//spacedock:live-fixture`
annotations; runtime-only proofs use `//spacedock:live-proof`. The executable
workflow remains the authority for lane selection.

## Reconciliation boundary

The compact real-repository reconciliation joins this registry to actual common
test declarations and calls, fixture annotations/builders, TODO ownership, and
the three executable workflow selectors. It derives findings from current source;
it does not use a copied gap oracle, mutation laboratory, or recorded-SHA guard.

## Amendment discipline

A change to a live journey, fixture ID, runtime proof, source binding, builder, or
live-lane selector must run reconciliation in the same pull request.

XPASS requires an amendment before its owner is archived: remove the target's
source binding and matching reconciliation expectation, then run the unchanged
candidate without the binding and require PASS. The desired journey remains in
this registry.

Update this registry only when desired state changes. A builder move that keeps
its fixture ID does not require a registry change.
