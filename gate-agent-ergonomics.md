---
title: Make recorded gate operation self-guiding for First Officers
status: validation
source: "Durable-decisions sprint dogfood: manual 0c and xb gate/round operation, 2026-07-24."
score: 1.0
id: skwchfe30ac6ntr63j1g0txj
sprint: durable-decisions
started: 2026-07-26T12:36:08Z
gates:
    version: 1
    records:
        - id: gate:skwchfe30ac6ntr63j1g0txj:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:skwchfe30ac6ntr63j1g0txj-backlog-1
              briefing:
                id: briefing:skwchfe30ac6ntr63j1g0txj:backlog:attempt-1:revision-1
                digest: sha256:49d1b8e3720bc54f0f9ef8db12d8ad33f7a277bdf156a9957b4a33d3a4134ddb
                request-digest: sha256:688d21e89d03f92cf3db7cdcbe1c4d46b4866c540fe022d41439b27afc348392
                room-ref: ./gate-agent-ergonomics/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:skwchfe30ac6ntr63j1g0txj:backlog:1
                briefing: briefing:skwchfe30ac6ntr63j1g0txj:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T12:36:24.13334Z"
                decision: approve
                reason: Roborev 2826 established a supported cold-restart stall after stage completion. Promote only the readiness projection and engage routing slice; reconcile prior scaffold/command work against landed s4 and add no second scheduler.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:skwchfe30ac6ntr63j1g0txj:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:skwchfe30ac6ntr63j1g0txj-ideation-1
              briefing:
                id: briefing:skwchfe30ac6ntr63j1g0txj:ideation:attempt-1:revision-1
                digest: sha256:587fa107d3df8935fd4bfc212273e3fce57cb45e0f2834c5b5d786363fa7b4c2
                request-digest: sha256:65a2e61692f30c33753f208d78a71547ac0ba3e93666b302d589b18d31b328fb
                room-ref: ./gate-agent-ergonomics/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:skwchfe30ac6ntr63j1g0txj:ideation:1
                briefing: briefing:skwchfe30ac6ntr63j1g0txj:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T13:26:20.52081Z"
                decision: approve
                reason: Independent staff re-review found all four prior material findings closed; the design delivers one readiness projection with no new scheduler or state machinery. Approval is recorded now; application waits for s4, gqs, and 0m6 to land.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:skwchfe30ac6ntr63j1g0txj:validation
          stage: validation
          attempts:
            - id: gate-attempt:skwchfe30ac6ntr63j1g0txj-validation-1
              briefing:
                id: briefing:skwchfe30ac6ntr63j1g0txj:validation:attempt-1:revision-1
                digest: sha256:a29d71f9cc1bdd4664ecd6c58325791cd4f455dbdb4a53b3a3370aafb7b52c3f
                request-digest: sha256:2a013593bd7763a9633da55c7444fd925fc71191c08e6490a70c6a4e0a0d7712
                room-ref: ./gate-agent-ergonomics/review/validation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-04T15:32:17.489446Z"
                reason: validation candidate changed after clean rebase onto origin/main 1947aacb; prior Briefing and evidence are stale
            - id: gate-attempt:skwchfe30ac6ntr63j1g0txj-validation-2
              briefing:
                id: briefing:skwchfe30ac6ntr63j1g0txj:validation:attempt-2:revision-1
                digest: sha256:0aded4a5ce61b1f3220af35d9148c7253f8ecf5684ba24941c95b5e560bfe3c4
                request-digest: sha256:670753cbdd2797055e642fb835eef65db08310ed9130ea42c9c99856f87193ef
                room-ref: ./gate-agent-ergonomics/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:skwchfe30ac6ntr63j1g0txj:validation:2
                briefing: briefing:skwchfe30ac6ntr63j1g0txj:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-04T16:35:46.115735Z"
                decision: approve
                reason: Captain explicitly directed opening the PR. Deterministic exact-head validation and all AC evidence pass; required live lanes are deferred to exact-head CI and remain a merge prerequisite.
              application:
                target-stage: done
                state: superseded
worktree: .worktrees/spacedock-ensign-gate-agent-ergonomics
mod-block:
pr:
---

## Problem

The original brief bundled package scaffolding, path normalization, command guidance,
round identity, validation diagnostics, and readiness. The promoted defect is narrower:
Roborev 2826 proved that a gated current stage with a mechanically complete, committed
current-stage report and no selected current-stage attempt is invisible to both sides
of the First Officer scheduler. Gate stages are suppressed from `dispatchable`, while
the gate reducer calls the state `validating` and omits it from `ready_gates`.

```text
mechanically complete committed current-stage report
  + absent or prior-stage gates.current
  -> gate-readiness=validating
  -> ready_gates=[]
  + dispatchable=[]
  -> cold First Officer declares quiescence
```

Roborev job 2826 recorded this after `829308d8` in AgentsView session
`codex:019f9a58-f3c0-7653-8af9-e0dcd026ecb1` (`#0-2761 @2716-2728`). The parent
session confirmed the defect and assigned it separately from s4 in
`codex:019f855d-8be8-7112-b467-7197851a2dc5` (`#7016-7169 @7145-7152`).
`prepare-provider-neutral-gate-room/index.md` retains the source disposition at state
commits `6e6204bf`, `ae881c47`, and `27473b4b`.

This design adds one scheduler projection and its First Officer route. It does not own
gate preparation, report parsing, withdrawal, decision, consumption, provider
operation, or worker transport.

## Spike result

A removed real-CLI spike replayed two split-root states:

1. `status: validation`, no `gates:` block, and a committed complete
   `## Stage Report: validation`;
2. the same report and status, with `gates.current` still selecting ideation.

Both returned `gate-readiness="validating"`, `ready_gates=[]`, and
`dispatchable=[]`. After a current-stage Briefing was selected, both returned
`gate-readiness="awaiting-captain"` and exactly:

```json
[{"id":"task","slug":"task","current":"validation","readiness":"awaiting-captain"}]
```

`TestSpikeCompletedGatedReportColdRestartIsInvisible` passed both cases in 0.680s
before the throwaway test was deleted. AgentsView session
`codex:019f9e6d-d524-7720-a4e3-235f5093fe05` (`#0-28 @21-24`) retains the run.
The spike isolates report-to-selection discovery; it is not evidence for a new
mutation path.

## Prerequisites and owner map

Implementation may start only after these approved sibling designs have landed and
their exact implementation commits are pinned in the implementation Stage Report:

| Dependency | Reused authority | This task's required treatment |
|---|---|---|
| s4, `prepare-provider-neutral-gate-room` | `gate prepare`, ordinary human judgment inputs, request/Briefing/room publication, selection, exact replay, rollback, and emitted `room=` authority | Invoke it; make no gate preparation or recorder change. |
| gqs, `dispatch-entered-stage-after-gate-consume` | Latest exact-stage report selection, structural checklist/Summary checks, failed-item rejection, and path-scoped committed-byte cleanliness | Reuse its result as mechanical evidence; add no parser or Git probe. |
| 0m6, `withdraw-stale-open-gate-attempt` | Open/withdrawn/closed classifier, `withdrawn-awaiting-prepare`, `gate withdraw`, and s4 preparation of attempt N+1 | Land first. Pin its reducer behavior as a zero-delta dependency and add only integration controls. |

On the current checkout those three entity designs are approved for implementation but
their full surfaces are not yet present. This entity must not implement substitutes in
advance. If any landed owner cannot supply the behavior named above, or if its semantics
drift, this design returns to ideation.

The mutation owner remains s4. The report-evidence owner remains gqs. The withdrawal
owner remains 0m6. The scheduler remains the existing status projection plus the
prose `«dispatch.next-action»()` loop. The existing registered Claude, Codex, and Pi
runners remain the behavioral proof surface.

## Proposed approach

Extend the canonical status-owned gate readiness reducer with one value:
`needs-preparation`. It means only:

> The current gate stage has no current-stage open/withdrawn/closed authority that
> already owns the route, and gqs mechanically found a latest exact-current-stage
> report with a nonempty conforming checklist, evidence/rationale, no `FAILED` item,
> a nonempty Summary, and path-clean committed entity bytes.

It does not mean “semantically complete,” “approved,” or “safe to mutate.” The First
Officer owns that judgment at engage. A heading, status string, prior-stage report, or
older valid report cannot produce the candidate.

### Canonical projection

| Durable current-stage state | Mechanical report evidence | Readiness/visibility | Existing owner and route |
|---|---|---|---|
| No current-stage gate authority | Absent, dirty, malformed, partial, failed, wrong-stage, or later malformed current report | `validating`, omitted | Report/stage work remains incomplete. |
| No current-stage gate authority | gqs mechanical proof passes | `needs-preparation`, ready | Semantic gate below; if accepted, exactly one s4 prepare. |
| `gates.current` names a prior stage; no current-stage authority exists | gqs mechanical proof passes | `needs-preparation`, ready | Semantic gate, then exactly one s4 prepare for current stage. |
| Prior pointer hides an exact current-stage open s4 replay candidate | gqs mechanical proof passes | `needs-preparation`, ready | Exactly one s4 prepare reuses its selection/replay idempotency; no new authority mechanism. |
| Selected current-stage open attempt | Irrelevant | `awaiting-captain`, ready | Present/resume that attempt; zero prepare calls. |
| Selected or prior-pointer-hidden current-stage withdrawn attempt | Irrelevant | `withdrawn-awaiting-prepare`, ready | 0m6 route; s4 prepare appends N+1. Never `needs-preparation`. |
| Selected current-stage approve/advance pending, nonterminal target | Irrelevant | `approved-awaiting-advance`, ready | Existing consume/commit/dispatch route. |
| Selected current-stage approve/advance pending, terminal target | Irrelevant | `approved-awaiting-merge`, ready | Existing consume/commit/merge route. |
| Open attempt is semantically stale | Irrelevant | `awaiting-captain` until withdrawal | 0m6 `gate withdraw` first, commit, then `withdrawn-awaiting-prepare`. |
| Closed attempt is stale | Irrelevant | Existing closed/pending state | Existing consume/supersede owner first; never direct preparation. |
| Revise, hold, blocked, feedback-pending, consumed, superseded, or not-applicable | Irrelevant | Existing value or omission | Existing feedback/stop/successor/merge owner; never prepare. |
| Terminal or archived entity | Irrelevant | Existing omission | Existing terminal/archive owner. |
| Malformed gate document, malformed selected attempt, mismatched closed/current authority, or unknown topology | Any | `invalid` or existing fail-closed omission | Diagnostic/repair owner; zero preparation. |

The reducer must inspect current-stage authority before using report evidence, including
when a prior pointer hides it. That ordering prevents an existing withdrawn, closed,
pending, consumed, revise, hold, blocked, or malformed current-stage attempt from being
misclassified as “absent.” `approved-awaiting-advance` and
`approved-awaiting-merge` keep their current meanings. No durable value or schema is
added.

### Status surfaces and one scheduler envelope

- `status --boot --identify --json` keeps four-key `ready_gates` rows and includes
  `needs-preparation` in the actionable set. Ordinary `--boot` remains unchanged.
- `status --next --json` returns one object containing the unchanged `dispatchable`
  array and the same ordered `ready_gates` array. Human `--next` remains unchanged.
- `gate-readiness` through existing field projections uses the same reducer.
- `needs-preparation` never enters `dispatchable`; gqs keeps ownership of non-gated
  entered-stage dispatch.

Every `status --next --json` read in `«dispatch.next-action»()` consumes that same
two-array envelope. This includes the initial read and the second read after the idle
hook and roster reconciliation. On every read, process the first `ready_gates` row
before any `dispatchable` row. Only when both arrays are empty may the initial read
invoke the existing idle sequence; only when both arrays remain empty on the second
read may the iteration stop.

Therefore an idle hook that makes a gate ready cannot be followed by a false idle
report. If the post-idle envelope contains both a ready gate and a dispatchable entity,
the gate route wins exactly as it does on the initial read. This changes no scheduler,
hook, adapter, or runner.

### `needs-preparation` engage route

For one selected `needs-preparation` row:

1. Load s4's gate lifecycle ownership before gate action, then re-read the entity,
   latest exact-stage report, reconstructed stage checklist, and report-bearing commit
   under the existing FO read/write-scope rules.
2. Decide semantic completeness. The binary row is only a mechanical candidate.
3. If any obligation, evidence claim, Summary conclusion, or worker-scope claim is
   semantically insufficient, perform zero mutation and stop this engage exactly once
   with `report-incomplete: <concrete unmet obligation or evidence defect>`. Do not run
   `gate prepare`, state commit, presentation, idle hooks, or another scheduler
   iteration. Add no persisted or ephemeral veto cache; a later externally initiated
   engage simply re-evaluates current durable state.
4. If the semantic review passes, choose the ordinary human inputs s4 requires:
   question, one committed Markdown Artifact, its concise human-written summary, and
   any committed References. Invoke exactly once:

   ```text
   ${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY \
     --question QUESTION --artifact ARTIFACT --summary SUMMARY \
     [--reference REFERENCE ...] --workflow-dir WORKFLOW_DIR
   ```

   The FO supplies judgment and file choices only. It does not author JSON, ids,
   digests, locators, or rooms.
5. On nonzero prepare, stop with its concrete error and do not retry. On success,
   require s4's `room`, `briefing`, `digest`, and `state=open` output, commit the
   entity-owned room/binding through the existing state owner, and re-read
   `status --next --json`.
6. Require exactly one same-slug `awaiting-captain` row from the same
   `dispatchable+ready_gates` envelope before presentation. Any other result stops
   fail-closed. Present the emitted binding through the existing gate lifecycle.

There are zero positive-path or recovery uses of legacy
`gate record --briefing`. This task does not modify that compatibility surface; it
simply never routes new preparation through it.

S4 owns exact replay and prior-pointer selection idempotency. A crash after successful
prepare is recovered by the next scheduler read seeing `awaiting-captain`, not by a
blind second invocation. A prior-pointer replay candidate is the only case where the
one authorized invocation may exercise s4's exact replay/selection behavior. This task
adds no attempt counter, retry token, cache, or transaction.

### Semantic-veto control

The required counterexample is a path-clean committed current-stage report that passes
every gqs structural check but whose prose omits one stage obligation. Status must emit
exactly one `needs-preparation` row. The First Officer must stop exactly once with, for
example, `report-incomplete: no evidence addresses the rollback acceptance check`.
The command log must contain zero `gate prepare`, zero `gate record --briefing`, zero
state mutation/commit, zero presentation/provider call, and zero repeated
`status --next` or idle-loop attempt after the veto. Entity and room trees remain
byte-identical. This is a deliberate two-layer result, not a request for the binary to
parse meaning.

## Guard matrix

The focused fixture must exercise all rows below against the production reducer and
existing FO lifecycle:

| Control | Expected projection/route | Required falsifier |
|---|---|---|
| Absent gate + mechanically complete clean current report | One `needs-preparation`; semantic pass invokes one s4 prepare | Removing promotion returns empty ready rows and no durable open room. |
| Prior-stage pointer + no current-stage authority + complete report | One `needs-preparation`; one current-stage s4 prepare | Current-stage binding/selection absent after engage. |
| Prior pointer hides exact open replay candidate | One `needs-preparation`; one s4 replay/selection invocation, unchanged attempt cardinality | Duplicate attempt or wrong selected stage. |
| Selected current-stage open | `awaiting-captain`; zero prepare | Selected-open cardinality mutant prepares once or changes attempt count and must fail. |
| Selected withdrawn and prior-pointer-hidden withdrawn | `withdrawn-awaiting-prepare`; s4 prepare N+1 under 0m6 | Either becomes `needs-preparation`, reuses N, or bypasses 0m6. |
| Stale open | Withdraw and commit first; then 0m6 readiness and s4 N+1 | Direct prepare or fabricated Resolution. |
| Stale closed | Existing consume/supersede route | Direct `needs-preparation` or prepare. |
| Approved pending nonterminal/terminal | Existing advance/merge rows and owner | Any prepare or wrong successor route. |
| Revise, hold, blocked, feedback-pending | Existing owner/stop | Any prepare or application bypass. |
| Consumed, superseded, not-applicable | Existing successor/merge/omission | Repeat consume or prepare. |
| Terminal and archived | Existing omission | Any ready or dispatch row. |
| Malformed gate or mismatched current closed authority | `invalid`/fail closed | Any prepare or entity mutation. |
| Prior-stage report only | `validating`, no ready row | Matching by positional-latest report. |
| Older valid current report followed by later malformed current report | `validating`, no ready row | Falling back to the older valid section. |
| Dirty target entity/report | `validating`, no ready row | Reading committed bytes while worktree differs. |
| Clean target with dirty unrelated sibling | Positive `needs-preparation` | Repository-wide dirt suppresses the target. |
| Structurally valid, semantically incomplete report | Mechanical row, then one concrete zero-side-effect veto | Prepare, mutation, vague stop, cache, or loop. |
| Ordinary non-gate entered stage | Unchanged gqs dispatch projection | Gate readiness changes non-gate bytes. |
| First envelope empty; idle hook creates ready gate | Second envelope routes ready gate before dispatch/idle stop | Post-idle gate is ignored or a dispatchable wins. |

Whole-entity and room-tree bytes, exact JSON, command order/cardinality, attempt ids and
counts, and state-commit ancestry are the offline oracles. Transcript claims are not.

## Mechanism-to-value trace

| New mechanism | Value served | Simpler alternative | Why insufficient |
|---|---|---|---|
| `needs-preparation` from gqs mechanical proof after existing authority classification | AC-1 | More s4 prose | Roborev 2826 has no scheduled slug on which that prose can run. |
| Same ready rows in every machine `--next` envelope | AC-1, AC-2 | Greeting-only boot rows | Split-root readiness can change after greet and during idle hooks. |
| Ready-gate-first handling on both initial and post-idle reads | AC-2 | Apply priority only on the first read | The existing second read can otherwise reproduce the same false quiescence. |
| Semantic veto before one s4 prepare | AC-1, AC-3 | Let structure authorize mutation | Structural conformance cannot establish that evidence satisfies the stage. |
| Pin 0m6 and reuse s4/gqs | AC-3, AC-4 | Add local withdrawal/preparation/report logic | That duplicates three approved owners and creates conflicting state transitions. |

## Acceptance criteria

**AC-1 (VALUE) — A mechanically complete cold gate becomes discoverable without manufacturing semantic completion.**
Absent and prior-pointer cold states each emit
exactly one `needs-preparation` row instead of the current empty scheduler. A semantic
pass performs exactly one s4 `gate prepare`, zero `gate record --briefing`, commits the
open binding, and re-reads the same slug as `awaiting-captain`. A structurally valid
but semantically incomplete control emits the candidate yet stops once with a concrete
`report-incomplete:` reason and zero side effects. **Verified by:** the two-case
split-root real-CLI fixture, exact command/tree oracles, and the existing recorded-gate
live scenario with operation coaching removed. Deleting the promotion restores the
observed 0/2 discovery baseline.

**AC-2 (SCHEDULER) — Every machine next read uses one envelope and one priority rule.**
Boot identify and machine next expose the same ordered ready rows. Both the initial and
post-idle `status --next --json` reads consume `dispatchable+ready_gates`, always
routing a ready gate first. Existing dispatch row bytes and human output remain
unchanged. **Verified by:** exact JSON parity, mixed ready/dispatch controls on both
reads, and an idle-hook control that creates a ready gate only before the second read.
A first-read-only mutant falsely stops idle or dispatches the lower-priority row.

**AC-3 (OWNERSHIP/INTEGRITY) — Preparation occurs only in the two authorized recovery shapes and reuses landed owners.**
Selected open, approved pending, revise, hold,
blocked, feedback, consumed, superseded, not-applicable, stale, withdrawn, terminal,
archived, malformed, and mismatched-owner states never route directly through
`needs-preparation`. Selected and prior-pointer-hidden withdrawn attempts retain
0m6's `withdrawn-awaiting-prepare -> s4 prepare N+1` route; stale open withdraws first,
and stale closed uses consume/supersede. **Verified by:** the full guard matrix,
selected-open cardinality mutant, whole-tree byte controls, and pinned 0m6/s4 replay
tests. Any duplicate attempt, legacy bind, direct stale/withdrawn preparation route, or
unauthorized mutation fails.

**AC-4 (SCOPE) — The correction adds no second owner or execution surface.** Gqs is the
only report shape/durability classifier; s4 is the only preparation/room publisher; 0m6
is the withdrawal/replacement-state owner; status is the one readiness projection;
the prose event loop is the existing scheduler; existing live runners grade behavior.
There is no new parser, store, command, schema, adapter, provider branch, runner, or
harness. **Verified by:** dependency-tip audit, production import/call trace, existing
gqs/0m6/s4 suites, full/race tests, path-to-lane checks, and registered Claude/Codex/Pi
jobs. Any semantic expansion named above requires design re-entry regardless of diff
size.

## Test plan

Implementation begins by pinning landed s4, gqs, and 0m6 commits and adding the red
status/FO controls before changing reducer or skill text. Reuse `buildSplitRoot`, gqs's
report fixtures, the existing gate lifecycle fixture, and existing registered runners.
Add no fixture-only production branch or new harness.

Focused offline:

```bash
go test ./internal/status -run 'TestBoot.*ReadyGate|Test.*GateReadiness|TestEnteredStage' -count=1
go test ./internal/gates -run 'TestCurrentStageReadiness|TestPrepare|TestWithdraw' -count=1
go test ./internal/ensigncycle -run 'TestRecordedGateLifecycle|TestRecordedGateReview' -count=1
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

The matrix includes absent gate, prior pointer, prior-pointer-hidden open,
selected open, withdrawn selected/hidden, stale open/closed, pending advance/merge,
revise, hold, blocked, feedback, consumed, superseded, not-applicable, terminal,
archived, malformed and mismatched authority, prior-stage report mismatch, later
malformed masking older valid, dirty target, dirty sibling positive, structural-pass
semantic-veto, ordinary non-gate, mixed ready/dispatch priority, and post-idle
ready-gate priority.

After offline green, run the existing recorded-gate lanes without adding a runner:

```bash
SPACEDOCK_LIVE_MODEL=sonnet go test -tags live -count=1 -timeout 40m \
  -run 'TestLiveClaudeSharedScenarios/recorded-gate-lifecycle$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_MODEL=claude-opus-4-8 go test -tags live -count=1 -timeout 40m \
  -run 'TestLiveClaudeSharedScenarios/recorded-gate-lifecycle$' ./internal/ensigncycle -v
go test -tags live -count=1 -timeout 40m \
  -run 'TestLiveCodexSharedScenarios/recorded-gate-lifecycle$' ./internal/ensigncycle -v
go test -tags live -count=1 -timeout 40m \
  -run 'TestSharedScenarioRunnerCoverage|TestPiSharedScenarioCoverage|TestLivePiRecordedGateLifecycle' \
  ./internal/ensigncycle -v
```

The complete registered `claude-live` matrix for both model legs, `codex-live`, and
`pi-live` must pass. Missing authorization, skip, timeout, or red is not green. The
oracle grades command order and durable state.

## Expected implementation surface

The likely integration surface is:

- the landed gqs report-evidence helper and its tests;
- the existing gate readiness reducer and tests, after 0m6 lands;
- status discovery/JSON/actionable-row formatting and existing boot/next tests;
- the existing recorded-gate lifecycle scenario and its host-neutral oracle;
- `fo-gate-lifecycle`, `first-officer-shared-core`, and `fo-dispatch-core`;
- gate concept and command-reference documentation.

Actual files and changed lines are reconciled in implementation review; LOC counts and
tolerances have no authority. Diff size alone neither approves nor rejects the work.
Semantic expansion is the reset trigger: a new parser, store, command, schema,
adapter, provider path, runner, harness, scheduler, retry/cache, or lifecycle mutation
requires design re-entry. A dependency moving code without changing the named contract
requires only a pinned-tip audit and zero-delta integration evidence.

No implementation may edit s4 preparation semantics, 0m6 withdrawal/reducer semantics,
gate recorder/application rules, provider code, host adapters, workflow schema, or
room scaffolding under this entity.

## Documentation diff

In `docs/site/concepts/gates-and-decisions.md`, replace the validating-only recovery
description with:

> A gate with no current-stage authority remains `validating` until gqs's mechanical
> report checks pass. It then appears as `needs-preparation` on boot and every machine
> scheduler read. Engage performs semantic report review. A concrete
> `report-incomplete:` veto stops without mutation; otherwise the First Officer calls
> `gate prepare` exactly once with its question, Artifact, summary, and References,
> commits the emitted binding, re-reads `awaiting-captain`, and presents it. Open,
> withdrawn, stale, closed, and spent attempts retain their existing lifecycle routes.

In `docs/site/reference/command-reference.md`, document that JSON `--next` contains
both unchanged `dispatchable` and canonical `ready_gates`, including the mechanical
`needs-preparation` candidate and 0m6's `withdrawn-awaiting-prepare`. State that the
First Officer applies ready-gate-first priority on every read and that preparation uses
the existing `gate prepare` command. Do not document legacy `record --briefing` as the
recovery route.

## Out of scope

- Gate-room/package scaffolding, path normalization, Briefing/request identity,
  provider behavior, materialization, and advisory-round preparation.
- Withdrawal or stale-state transitions, recorder/application eligibility, decision
  authority, merge behavior, feedback routing, and report parsing.
- New validation-error persistence, veto cache, retry loop, report epoch, lease,
  ledger, queue, compatibility layer, driver command, scheduler, runtime adapter, live
  runner, or test harness.
- Expanding gqs's non-gated entered-stage contract or weakening s4/0m6 authority.

## Stage Report: ideation

- DONE: Route `needs-preparation` through exactly one s4 `gate prepare` with ordinary human judgment inputs after loading lifecycle ownership; use zero legacy `gate record --briefing`, commit the emitted binding, re-read `awaiting-captain`, and reuse s4 replay/selection idempotency.
  AC-1 and AC-3 evidence: the engage procedure fixes exact inputs, command cardinality, output/commit barriers, same-envelope re-read, recovery behavior, and the selected-open mutant.
- DONE: Define `needs-preparation` as a mechanical candidate and make semantic-report veto a concrete one-stop, zero-side-effect result without a cache or loop.
  AC-1 evidence: the structural-pass/semantic-fail control requires one `report-incomplete:` stop and zero prepare, legacy bind, commit, presentation, mutation, idle, or repeat-next effects.
- DONE: Preserve 0m6 withdrawn and stale ownership, including selected and prior-pointer-hidden withdrawn controls, `withdrawn-awaiting-prepare -> s4 prepare N+1`, withdraw-first stale-open handling, and consume/supersede stale-closed handling.
  AC-3 evidence: 0m6 is a landing prerequisite and pinned zero-delta dependency; the projection and guard matrices forbid direct `needs-preparation` for every withdrawn or stale state.
- DONE: Make every existing machine scheduler read consume the same `dispatchable+ready_gates` envelope with ready-gate-first priority, including the post-idle double-check.
  AC-2 evidence: initial/post-idle mixed controls and the idle-hook-created ready gate fail if either read ignores ready rows, chooses dispatch first, or reports idle.
- DONE: Complete the fail-closed guard matrix for report mismatch/masking/cleanliness, every existing gate owner, and the selected-open cardinality mutant.
  AC-2 and AC-3 evidence: exact JSON, whole-tree bytes, command counts, attempt identity, dirty-target/dirty-sibling controls, and all open/pending/consumed/revise/hold/blocked/terminal/archived/malformed/mismatched routes are specified.
- DONE: Remove LOC-count authority and retain only semantic expansion as design re-entry, with s4/gqs/0m6 prerequisites, existing fixtures/runners, four acceptance criteria, and final evidence accounting.
  AC-4 evidence: the owner map, implementation surface, live/offline plan, exclusions, and reset triggers add no parser, store, command, schema, adapter, provider, runner, harness, scheduler, or mutation owner.

### Summary

Ideation now isolates one missing scheduler projection. Mechanical report evidence
makes the cold gate discoverable; human semantic review decides whether the existing
s4 preparation operation may run. All neighboring states retain their gqs, s4, 0m6,
application, feedback, merge, terminal, and archive owners.

## Stage Report: implementation

- DONE: Implemented the canonical status-owned `needs-preparation` projection. The
  reducer classifies existing current-stage authority first, then promotes only a
  mechanically complete, path-clean, committed exact-stage report. Absent authority
  and prior-stage authority promote; selected/open, withdrawn, closed, malformed, and
  empty-record documents retain their existing owner or fail closed. No parser, store,
  schema, command, adapter, provider, runner, scheduler binary, retry, cache, or
  lifecycle mutation was added.
- DONE: Extended machine `status --next --json` to the ordered
  `command`/`dispatchable`/`ready_gates` envelope. Boot identify and next use the same
  four-key ready-gate rows; `needs-preparation` is actionable but never dispatchable.
  Human `--next` remains unchanged. Updated the canonical JSON golden to include the
  empty `ready_gates` array.
- DONE: Added the First Officer route in `fo-gate-lifecycle` and `fo-dispatch-core`:
  re-read the exact report and commit, issue one concrete `report-incomplete:` veto
  with zero side effects when semantic review fails, or invoke existing `gate prepare`
  exactly once on a semantic pass, commit/re-read `awaiting-captain`, and present once.
  The scheduler consumes the same envelope on the initial and post-idle reads and
  applies ready-gate-first priority. Existing open/withdrawn/pending/revise/hold/
  consumed/terminal routes remain owned by s4, gqs, or 0m6.
- DONE: Pinned and audited dependency tips: s4
  `acae980fc145e624d9e04e7ec9f7fdb599585f6e`, gqs
  `cb01129b6325f7af646363785c24ef69e8bd16bd`, and 0m6
  `9881639697d1af391133c9ecf4111fd1673f537c`; all are ancestors of implementation
  HEAD `ae1233a8b`.
- DONE: Exact code/docs scope at `ae1233a8b` (13 files, 227 insertions, 22 deletions):
  `internal/gates/model.go`, `internal/gates/gates_test.go`,
  `internal/status/discover.go`, `internal/status/format.go`,
  `internal/status/handlers.go`, `internal/status/json_commands.go`,
  `internal/status/gate_readiness_needs_preparation_test.go`,
  `internal/status/testdata/golden/seq-next.json`,
  `docs/site/concepts/gates-and-decisions.md`,
  `docs/site/reference/command-reference.md`,
  `skills/fo-gate-lifecycle/SKILL.md`,
  `skills/first-officer/references/fo-dispatch-core.md`, and
  `skills/first-officer/references/first-officer-shared-core.md`.
- DONE: Formatting and static checks: `gofmt -w ./cmd ./internal` completed with no
  remaining diff; `git diff --check` passed; `go test ./internal/contractlint -count=1`
  passed.
- DONE: Focused functional evidence:
  `go test ./internal/status -run
  'TestGateReadiness|TestBootIdentifyReadyGates|TestStatusProjectsSharedGateReadinessReducer|TestEnteredStage' -count=1`
  passed, including absent/prior authority promotion and dirty/malformed rejection;
  the full `go test ./internal/status -count=1` package passed. The focused gates
  readiness/prepare/withdraw tests, recorded-gate lifecycle/review tests, and JSON
  golden tests passed after the envelope update.
- DONE: Focused race evidence:
  `go test -race` passed for the focused status, gates, contractlint, and ensigncycle
  selections. Full race compile-only coverage (`go test -race ./... -run '^$' -count=1`)
  passed for every package.
- AC-1 evidence: `TestGateReadinessPromotesMechanicallyCompleteColdReports` covers
  absent and prior-stage authority, and asserts one `needs-preparation` row, an empty
  dispatchable array, the shared next envelope, and the `gate-readiness` projection;
  `TestGateReadinessRejectsDirtyAndMalformedColdReports` covers the fail-closed
  cleanliness/malformed controls. The semantic-veto and one-prepare engage behavior is
  encoded in the FO route text and awaits the registered live-lane validation.
- AC-2 evidence: `nextJSON` emits the ordered two-array envelope, the updated
  `seq-next.json` golden passes, and `TestBootIdentifyReadyGates` plus the focused
  status suite preserve boot/next dispatch parity and existing human output. The
  post-idle ready-gate-first rule is documented in `fo-dispatch-core`; it needs the
  live scheduler lane for behavioral execution evidence.
- AC-3 evidence: `CurrentStageReadinessWithReport` promotes only absent authority and
  leaves selected open/withdrawn/closed or malformed documents to existing reducers;
  focused `internal/gates` prepare/withdraw and `internal/ensigncycle` lifecycle/review
  tests pass. No preparation, withdrawal, recorder, or application implementation was
  changed.
- AC-4 evidence: the dependency ancestry audit and 13-file scope show one status
  projection plus declarative FO/docs changes; forbidden-scope search found no new
  parser/store/schema/adapter/provider/runner/scheduler/retry/cache implementation.
- BLOCKED: Full `go test ./...` and full `go test ./... -race` reached all packages but
  fail only in `internal/gates/TestV1PilotManifestReadsAndValidates`: six manifest
  paths named by the repository's v1 pilot fixture are absent from the shared state
  checkout (`bind-post-rework-briefing-at-rejection-regate.md`,
  `collapse-gate-approval-ceremony/index.md`,
  `minimize-v1-gate-application-schema/index.md`,
  `shared-git-scaffold-helper.md`, `status-pagination-and-default-sorting.md`, and
  `status-where-robust-and-discoverable.md`). The status JSON golden mismatch seen on
  the first full run was corrected in `ae1233a8b`; the full status package then passed.
  No live Claude/Codex/Pi runner lane was started in this implementation worktree, so
  there is no live authorization evidence to claim.
- NEXT: Keep the implementation entity at `status: implementation` pending Captain
  review/validation. The validation stage should first restore or mount the six shared
  state manifests, then rerun the complete offline/race matrix and the registered
  recorded-gate live lanes. Do not prepare, consume, or mutate a validation gate from
  this implementation report.

### Completion checklist accounting

The three dispatched checklist items are accounted for explicitly:

1. **DONE — implementation and acceptance criteria.** The approved projection,
   scheduler envelope, owner routing, dependency pins, and AC-1 through AC-4 evidence
   are recorded above at code HEAD `ae1233a8b`.
2. **FAILED — full-suite evidence; SKIPPED — live-lane evidence.** Both
   `go test ./...` and `go test ./... -race` reached every package but are red only at
   `internal/gates/TestV1PilotManifestReadsAndValidates` because six required manifest
   paths are absent from the shared state checkout. Registered Claude/Codex/Pi live
   lanes were skipped because this worktree had no live-lane authorization; no live
   result is claimed.
3. **DONE — report, scope, formatting, and adversarial accounting.** The exact
   13-file scope, dependency ancestry, `gofmt`, `git diff --check`, focused tests,
   malformed/dirty controls, and next-stage recommendation are recorded above.

Checklist totals: 3 top-level items — 2 DONE, 1 with FAILED full-suite evidence and
SKIPPED live-lane evidence; subcheck totals are 2 DONE, 1 FAILED, and 1 SKIPPED.

### Summary

Implementation evidence is complete for the recorded-gate readiness projection and
First Officer routing, with focused tests, race compile coverage, dependency pins, and
scope checks recorded above. Full offline and race suites remain blocked only by the
six absent shared-state manifests (`bind-post-rework-briefing-at-rejection-regate.md`,
`collapse-gate-approval-ceremony/index.md`,
`minimize-v1-gate-application-schema/index.md`, `shared-git-scaffold-helper.md`,
`status-pagination-and-default-sorting.md`, and
`status-where-robust-and-discoverable.md`); registered Claude/Codex/Pi live lanes were
skipped because this worktree had no live-lane authorization.

## Stage Report: validation

- DONE: Re-validated the implementation at exact code HEAD
  `f35436e0a7526ebf5aeeb9e3adaa1b9da4c951ea` (`status: surface cold gated report
  preparation`, parent `11308c81662b8399731811370207e858e3279552`). The code worktree
  was clean before and after validation. The implementation patch remains the same
  13-file, +227/-22 surface described above. The required dependency tips are all
  ancestors of this HEAD: s4 `acae980fc145e624d9e04e7ec9f7fdb599585f6e`, gqs
  `cb01129b6325f7af646363785c24ef69e8bd16bd`, and 0m6
  `9881639697d1af391133c9ecf4111fd1673f537c`. `git merge-base HEAD origin/main` is
  `20cfc809a7342b2428509c76ddd6e91423db39b7`.

- DONE: Audited changed-surface drift after the implementation report. The exact
  `origin/main...HEAD` range is 16 files, +801/-27: the same 13 implementation files
  (+227/-22) plus three pre-existing registry/design documents from the rebased
  parents (`docs/roadmap/durable-decisions/index.md`,
  `docs/roadmap/live-test-truth/index.md`, and `docs/runtime-live-ci-registry.md`).
  Those three files came from `11308c816`/`46e78e3fd`; validation added no code,
  command, schema, parser, scheduler, runner, or owner surface.

- DONE: Re-ran the focused acceptance and owner suites. All passed:
  `go test ./internal/status -run 'TestGateReadiness|TestBootIdentifyReadyGates|TestStatusProjectsSharedGateReadinessReducer|TestEnteredStage' -count=1`,
  `go test ./internal/gates -run 'TestCurrentStageReadiness|TestPrepare|TestWithdraw' -count=1`,
  `go test ./internal/ensigncycle -run 'TestRecordedGateLifecycle|TestRecordedGateReview' -count=1`,
  `go test ./internal/status -count=1`, and
  `go test ./internal/contractlint -count=1`. The focused status and gates runs
  cover absent and prior-stage authority promotion, path-clean committed reports,
  dirty/malformed rejection, selected open/withdrawn/closed authority, approved
  pending advance/merge, feedback/hold/consumed/superseded/not-applicable, ordinary
  and terminal stages, and invalid topology.

- DONE: Re-ran the semantic adversarial matrix against production reducers and the
  existing lifecycle oracles. These all passed: status matrix (`TestGate*`,
  `TestBootIdentifyReadyGates`, `TestStatusProjectsSharedGateReadinessReducer`,
  `TestEnteredStage`, and default-stage controls); gates matrix
  (`TestCurrentStageReadiness*`, withdrawal/prepare/replay/refusal, eligibility,
  record/consume, round, duplicate/prototype/digest fail-closed controls); and
  recorded-gate matrix (`TestRecordedGateLifecycleRealCLIReplay`,
  `EnteredStageRecoveryMatrix`, `AC5RefusalMatrix`, `AC7ResumeMatrix`,
  `ProvenanceMutants`, `MissingEventControls`, `TerminalConsumeHasNoDispatchableSuccessor`,
  `CollapsedConsumeTrace`, and gate assertion controls). The structural-pass,
  semantic-veto obligation remains judgment-owned by the FO route; the tests and
  skill contract require one concrete `report-incomplete:` stop with zero prepare,
  legacy bind, mutation, presentation, idle, or repeat-next effects. No live claim
  is made for a semantic-veto prose judgment that the existing lanes do not execute.

- DONE: Re-validated the machine envelope. The cold split-root fixture's exact
  `--next --json` contract is one object with ordered keys
  `command`, `dispatchable`, `ready_gates`; absent/prior candidates have an empty
  dispatchable array and one four-key `needs-preparation` row. The retained Claude
  boot read from the real recorded-gate fixture emitted this exact candidate row
  (the boot envelope also carries definition/stage metadata):

  ```json
  {"id":"recorded-gate-task","slug":"recorded-gate-task","current":"validation","readiness":"needs-preparation"}
  ```

  The retained Codex post-consume scheduler read emitted the corresponding exact
  next envelope (the lower-priority non-gate row remains in `dispatchable` and the
  gate array is empty after consumption):

  ```json
  {"command":"next","dispatchable":[{"id":"recorded-gate-task","slug":"recorded-gate-task","current":"handoff","next":"handoff","worktree":"no"}],"ready_gates":[]}
  ```

  The direct current-workflow read also returned the same ordered envelope with an
  empty `dispatchable` array and the ordered canonical `ready_gates` rows; no human
  `--next` output changed.

- FAILED: Complete offline and race suites are red only at the pre-existing shared
  state manifest oracle. `go test ./...` and `go test ./... -race` reached every
  package; all packages pass except
  `internal/gates/TestV1PilotManifestReadsAndValidates`. The six historical paths
  named in the implementation report remain absent:
  `bind-post-rework-briefing-at-rejection-regate.md`,
  `collapse-gate-approval-ceremony/index.md`,
  `minimize-v1-gate-application-schema/index.md`,
  `shared-git-scaffold-helper.md`, `status-pagination-and-default-sorting.md`, and
  `status-where-robust-and-discoverable.md`. The current run also reports
  `cut-workflow-specific-round-recorder-from-v1/index.md`, which was archived by
  pre-existing state commit `2987a085c` at 2026-08-03 22:48:45 +0800, after the
  implementation HEAD and its six-path report. All seven are absent from the
  active state root but present under `_archive`; no unrelated fixture was edited.

- DONE: Formatting and contract checks passed after the full runs:
  `gofmt -w ./cmd ./internal` (no code diff), `git diff --check`, and
  `go test ./internal/contractlint -count=1`. The final code worktree remained
  clean.

- DONE: Ran the registered live validation lanes and preserved retained evidence.
  Runtime access was available; no lane was labeled green when it skipped:

  - Claude Sonnet: exit 0, **SKIPPED** by
    `TODO(w5bfnrvpcphw857nzz93340c): Sonnet must reliably render the exact selected
    Briefing digest before re-enabling this journey`.
  - Claude Opus (`SPACEDOCK_LIVE_MODEL=claude-opus-4-8`): **PASS** on the initial
    run (372.58s) and on a retained-artifact rerun (415.43s).
  - Codex: **PASS** on the initial run (228.28s) and on a retained-artifact rerun
    (223.57s). The retained JSONL setup records source HEAD
    `f35436e0a7526ebf5aeeb9e3adaa1b9da4c951ea`.
  - Pi recorded-gate lifecycle: exit 0, **SKIPPED** by
    `TODO(9w59t6m1qc46hccd54p04z2j): delegated gate presentation-to-application/dispatch
    is unreliable`; `TestPiSharedScenarioCoverage` and
    `TestSharedScenarioRunnerCoverage` passed. No lane was unauthorized or
    falsely promoted to green.

  Retained artifacts are under
  `/tmp/spacedock-gate-agent-ergonomics-validation-live-20260804/` (outside both
  worktrees). Claude hashes: `claude-stream.jsonl`
  `d433b8e2b45f2f5c034c5d523a380f1f55e12f33123fdc800a81eb2460953cf4`,
  `command.log` `799ab8e82eec1482e7739fe18353a584c6904ae746d7cec28bc06d765f4b347d`,
  and `claude-final-message.txt`
  `06048e85f611b967d0742f5f968d1af9722e12f418544445c4e65550454f4aff`.
  Codex hashes: `codex-exec.jsonl`
  `be4aa028fcc047358907495faa748f63022646cf8de5421b911d1f9bdb239039`,
  `codex-exec.stderr.txt`
  `29d3053d2febe6c9fd03e6b014411b0fc67e2955964bec88842a6423397562dd`,
  `codex-final-message.txt`
  `6bce962d9aeb351b251d7cb8ef2371f9762a6da22bbc2f192412a9d7328cdff0`, and
  `codex-process-result.txt`
  `9b914c6b416c41418dd833bc079d8b02b398ed0baa492bfe13cfaddc65a74782`.
  The live artifacts show one cold `needs-preparation` boot row, exactly one
  `gate prepare`, one binding commit, one approval/consume path, and one committed
  successor dispatch for each green host; the semantic-veto and post-idle-created
  gate controls remain offline/prose-owned as stated above.

### Completion checklist accounting

1. **DONE — exact-head revalidation.** Acceptance criteria, focused status/gates
   suites, exact machine envelopes, semantic matrix, dependency ancestry, and
   changed-surface drift are recorded above for `f35436e0a`.
2. **FAILED — complete offline/race suites; DONE — formatting and contract checks.**
   The only failure is the seven absent active-state manifest paths described above;
   no implementation or unrelated state fixture changed.
3. **DONE — registered live lanes executed and classified.** Opus and Codex passed
   with retained artifacts; Sonnet and Pi were explicit TODO skips, not green proof.

Checklist totals: 3 top-level items — 2 DONE and 1 FAILED due to the pre-existing
shared-state manifest environment; live subcheck totals are 2 PASS, 2 SKIPPED, and
0 unauthorized.

### Summary

Validation confirms the exact-head readiness projection, canonical scheduler
envelope, fail-closed owner matrix, and green Opus/Codex recorded-gate lifecycle.
The complete offline and race commands remain red only because the shared state
checkout lacks six historical active manifests plus one entity archived after the
implementation report; this validation did not restore or edit those fixtures.
Sonnet and Pi remain explicit quarantined skips, while retained Opus/Codex artifacts
provide durable live evidence at the path and hashes above. The entity remains at the
Captain's validation gate for final review.

## Stage Report: validation (fresh exact-tip cycle)

- DONE: Re-validated the assigned candidate without editing implementation files.
  Code HEAD was exactly `94a66df8c937b58d7ffbebed3467ba135fdfb633`, the assigned
  clean-rebase base was `origin/main` `1947aacb0d7c3481c18a846f3566645fd2cb89ee`,
  and `git merge-base HEAD origin/main` at the target was that same `1947aacb0`.
  The code worktree was clean before and after every check. No gate prepare,
  record, consume, state mutation, or implementation edit was made.

- FAILED: The shared `origin/main` ref advanced during this validation to
  `db7f1e84aef5df2daf20fb02deac440df4ae1af1` (PR #614). The candidate is now stale:
  `git merge-base HEAD origin/main` remains `1947aacb0`,
  `git rev-list --left-right --count origin/main...HEAD` is `13 3`, and
  `git diff --stat HEAD..origin/main` is 17 files, +255/-260. That current-main
  range removes the candidate's readiness implementation and focused test files
  (`internal/gates/model.go`, `internal/gates/gates_test.go`,
  `internal/status/discover.go`, `internal/status/format.go`,
  `internal/status/handlers.go`, `internal/status/json_commands.go`, and
  `internal/status/gate_readiness_needs_preparation_test.go`) while also changing
  the gate/docs surface. This validator did not rebase or reconcile a changed
  candidate. Exact-tip evidence below is therefore not a current-main validation
  result and cannot clear the gate.

- DONE: Audited the exact assigned diff against its target base. The range
  `1947aacb0...94a66df8` is 16 files, +801/-27: the same 13-file implementation
  surface (+227/-22) from the implementation report plus the three parent-supplied
  roadmap/registry documents (`docs/roadmap/durable-decisions/index.md`,
  `docs/roadmap/live-test-truth/index.md`, and
  `docs/runtime-live-ci-registry.md`). Validation added no code, parser, store,
  schema, command, scheduler, runner, provider, adapter, or owner surface.
  The required s4/gqs/0m6 tips remain ancestors of exact HEAD:
  `acae980fc145e624d9e04e7ec9f7fdb599585f6e`,
  `cb01129b6325f7af646363785c24ef69e8bd16bd`, and
  `9881639697d1af391133c9ecf4111fd1673f537c`.

- DONE: Focused exact-tip functional checks all passed:

  - `go test ./internal/status -run 'TestGateReadiness|TestBootIdentifyReadyGates|TestStatusProjectsSharedGateReadinessReducer|TestEnteredStage' -count=1`
  - `go test ./internal/gates -run 'TestCurrentStageReadiness|TestPrepare|TestWithdraw' -count=1`
  - `go test ./internal/ensigncycle -run 'TestRecordedGateLifecycle|TestRecordedGateReview' -count=1`
  - `go test ./internal/contractlint -count=1`

  The full status package also passed. The focused status/gates controls independently
  reproduce absent and prior-stage authority promotion, committed/path-clean report
  proof, dirty/malformed rejection, selected open/withdrawn/closed authority,
  approved advance/merge, feedback/hold/consumed/superseded/not-applicable,
  ordinary/terminal stages, and invalid topology.

- DONE: The PR #612 shared-state manifest blocker is fixed on the exact target.
  `go test ./internal/gates -run 'TestV1PilotManifestReadsAndValidates$' -count=1 -v`
  passed every listed path, including the seven `_archive/` entries that were
  missing in the prior cycle (`bind-post-rework-briefing-at-rejection-regate.md`,
  `collapse-gate-approval-ceremony/index.md`,
  `cut-workflow-specific-round-recorder-from-v1/index.md`,
  `minimize-v1-gate-application-schema/index.md`,
  `shared-git-scaffold-helper.md`, `status-pagination-and-default-sorting.md`,
  and `status-where-robust-and-discoverable.md`).

- DONE: Complete exact-tip offline and race suites passed. `go test ./...` passed
  all packages, including `internal/gates`, `internal/status`,
  `internal/ensigncycle`, `internal/dispatch`, `internal/cli`, and
  `skills/integration`. `go test ./... -race` passed the same package set. The
  previous manifest-only red is absent at `94a66df8`; these results do not waive
  the later current-main drift above.

- DONE: Formatting and contract checks passed without changing code:
  `gofmt -d ./cmd ./internal` produced no output, `git diff --check` passed, and
  `go test ./internal/contractlint -count=1` passed. The read-only `gofmt -d`
  check was used because validation is forbidden to edit implementation files.

- DONE: Re-ran the required semantic adversarial matrix against production
  reducers and existing lifecycle oracles. The status matrix passed
  `TestGateReadinessPromotesMechanicallyCompleteColdReports`,
  `TestGateReadinessRejectsDirtyAndMalformedColdReports`,
  `TestBootIdentifyReadyGates`, `TestStatusProjectsSharedGateReadinessReducer`,
  `TestEnteredStageProjectionRequiresCommittedCompleteReport` (including wrong
  stage, later-malformed masking, dirty target, and dirty sibling controls),
  `TestEnteredStageMutationControls`, `TestEnteredStageLegacySuppressionControls`,
  default-stage controls, and boot/next parity. The gates matrix passed
  `TestCurrentStageReadinessFailClosedTable`,
  `TestCurrentStageReadinessWithReportPromotesOnlyAbsentAuthority`, all prepare/
  withdraw/replay/refusal controls, eligibility, record/consume, round,
  duplicate/prototype/unknown-shape, and digest-domain fail-closed controls. The
  recorded-gate matrix passed `TestRecordedGateLifecycleRealCLIReplay`,
  withdrawn cold-boot replacement, terminal consume, entered-stage recovery,
  AC5 refusal, AC7 resume, provenance mutants, collapsed consume, committed-before-
  dispatch, missing-event, phase-help, and gate-guardrail controls. Exact output,
  ordering, cardinality, byte-clean refusal, authority identity, and terminal-state
  assertions all passed.

- DONE: AC-1 exact-tip mechanical evidence is independently reproduced. The
  split-root status fixture emits exactly one four-key `needs-preparation` row for
  both absent and prior-stage authority, an empty `dispatchable` array, and the
  same ordered two-array `--next --json` envelope; field projection reports the
  same readiness. Dirty target and malformed gate controls emit no candidate. The
  structural-pass/semantic-veto path remains judgment-owned by the First Officer
  contract: `skills/fo-gate-lifecycle/SKILL.md` and
  `skills/first-officer/references/fo-dispatch-core.md` require one concrete
  `report-incomplete:` stop with zero prepare, legacy bind, commit, presentation,
  idle, or repeat-next effects. No live lane exercised that prose-only semantic
  veto, so no behavioral veto claim is made.

- DONE: AC-2 exact-tip machine-envelope behavior is independently reproduced by
  status tests and the updated JSON golden. `--next --json` emits ordered
  `command`, `dispatchable`, `ready_gates`; ready rows retain exactly
  `id`, `slug`, `current`, `readiness`, and `needs-preparation` never enters
  `dispatchable`. Human `--next` behavior remains unchanged. The ready-gate-first
  initial/post-idle event-loop rule is present in the existing declarative FO
  contract, but no executable cold-candidate/post-idle live control exists; this
  is an evidence limitation, not an inferred pass.

- DONE: AC-3 exact-tip ownership and integrity matrix passes. Existing open,
  withdrawn, stale, closed, pending, revise, hold, blocked, feedback, consumed,
  superseded, not-applicable, terminal, archived, malformed, and mismatched
  authorities remain on their existing reducers; prepare/withdraw/replay and
  whole-tree byte controls pass. No new recorder, withdrawal, application,
  provider, parser, scheduler binary, retry, cache, or mutation owner appears in
  the exact diff.

- DONE: AC-4 exact-tip scope audit passes. Status is the only new readiness
  projection; s4 remains preparation/room owner, gqs remains report-evidence
  owner, and 0m6 remains withdrawal/replacement owner. The exact diff contains no
  forbidden execution surface. The current-main drift means this audit must be
  repeated after a clean rebase before final approval.

- FAILED: Required registered live coverage is incomplete, and no complete PASS is
  claimed. Exact-tip results were:

  - Claude Sonnet:
    `SPACEDOCK_LIVE_MODEL=sonnet go test -tags live -count=1 -timeout 40m -run 'TestLiveClaudeSharedScenarios/recorded-gate-lifecycle$' ./internal/ensigncycle -v`
    exited 0 but **SKIPPED** by
    `TODO(w5bfnrvpcphw857nzz93340c): Sonnet must reliably render the exact selected Briefing digest before re-enabling this journey`.
  - Claude Opus (`SPACEDOCK_LIVE_MODEL=claude-opus-4-8`) passed the same lane in
    439.80s.
  - Codex first exact-tip run exited 1 after 275.19s: the live oracle observed
    two successful `dispatch build` attempts instead of the required one (`2/2`,
    want `1/1`). An immediate clean rerun passed in 244.70s. This is recorded as
    a transient host/model reproducibility finding; it was not hidden or promoted
    to green by the first failure, and no implementation change was authorized.
  - Pi:
    `go test -tags live -count=1 -timeout 40m -run 'TestSharedScenarioRunnerCoverage|TestPiSharedScenarioCoverage|TestLivePiRecordedGateLifecycle' ./internal/ensigncycle -v`
    passed both coverage meta-tests, while the lifecycle lane was **SKIPPED** by
    `TODO(9w59t6m1qc46hccd54p04z2j): delegated gate presentation-to-application/dispatch is unreliable`.

  The Sonnet/Pi skips are existing quarantines, not authorization failures, but
  the stage contract explicitly says a skipped or red required lane is not green.
  The exact-tip live set therefore remains incomplete.

- DEFERRED RISK: The Codex duplicate-dispatch observation is a supported-host
  reproducibility risk, not a proven candidate code defect: the only source delta
  from the prior validated candidate is the PR #612 manifest-path refresh, the
  offline dispatch cardinality oracles pass, and the immediate rerun passed. Promote
  this to a material outcome finding if two consecutive clean runs reproduce
  duplicate builds, or if an offline command-log control fails; do not fix it in
  this entity because live-test-truth sprint members are out of scope.

- MATERIAL EVIDENCE FINDING: The candidate is stale against current `origin/main`
  (`db7f1e84`), and the required live matrix is incomplete (Sonnet/Pi skips plus
  one Codex red run). These are evidence/release-scope blockers, not a finding that
  the exact 94a implementation fails its offline value behavior. A clean rebase or
  explicit current-main reconciliation, followed by the full focused/full/race/
  semantic/live matrix, is required before a new validation recommendation.

### Completion checklist accounting

1. **DONE for the assigned exact tip; BLOCKED for current-main release evidence.**
   Focused, full, race, manifest, formatting, dependency, acceptance, and semantic
   matrix checks were reproduced at `94a66df8`; the candidate then became stale when
   shared `origin/main` advanced to `db7f1e84`.
2. **DONE for exact-tip offline/race; FAILED for complete live coverage.** Both
   complete offline commands pass at the assigned tip. Opus and one Codex rerun
   pass; Sonnet and Pi are explicit quarantined skips, and the first Codex run is
   a recorded red duplicate-dispatch result.
3. **DONE — report, scope, findings, and no-mutation accounting.** Only this
   state-entity report is being changed; no implementation or gate bytes changed.

Checklist totals: exact-tip offline checks 2 PASS; required live subchecks 2 PASS,
2 SKIPPED, and 1 transient RED before rerun. Current-main freshness is a material
evidence blocker. **Recommendation: REJECTED/HOLD** pending clean rebase or
reconciliation onto current `origin/main` and complete supported live coverage.

### Summary

At the assigned exact HEAD `94a66df8`, the full offline/race suites, PR #612
manifest oracle, focused ownership suites, and adversarial matrix all pass, and
Claude Opus plus a Codex rerun pass the recorded-gate lane. Validation cannot
recommend a complete PASS: shared `origin/main` advanced to `db7f1e84` and now
removes the candidate surface, Sonnet and Pi remain explicitly skipped, and the
first Codex run recorded a duplicate dispatch before its clean rerun. No candidate
or gate mutation was made; rebase/reconciliation and fresh complete live evidence
are required.

## Stage Report: validation (fresh exact-head cycle after PR #612)

- DONE: Revalidated the candidate at exact code HEAD
  `cd76f52abc3b3f00c0344566ad039f62586936d2`, which is the assigned clean-rebase
  target and is currently exactly `origin/main`
  (`db7f1e84aef5df2daf20fb02deac440df4ae1af1`). `git merge-base HEAD origin/main`
  equals `db7f1e84`; `git rev-list --left-right --count origin/main...HEAD` is
  `0 1`. The code worktree was clean before and after every check. Validation made
  no implementation, gate, room, or state mutation.

  Validation records the rebase lesson explicitly: a rebase invalidates prior
  evidence because the candidate bytes changed; it does not automatically imply a
  behavior defect. This cycle therefore reran the local focused, full, race, and
  semantic checks at the exact head, while the registered live matrix is handed off
  as one exact-head CI run.

- DONE: Audited the exact diff against current main. `git diff --stat
  origin/main...HEAD` is the same 13-file, +227/-22 implementation surface from
  the implementation report: the status/gates projection and tests, JSON golden,
  First Officer contract text, and gate documentation. Required s4, gqs, and 0m6
  dependency tips remain ancestors of this HEAD:
  `acae980fc145e624d9e04e7ec9f7fdb599585f6e`,
  `cb01129b6325f7af646363785c24ef69e8bd16bd`, and
  `9881639697d1af391133c9ecf4111fd1673f537c`. No new parser, store, schema,
  command, scheduler binary, runner, provider, adapter, retry, cache, or lifecycle
  owner was introduced.

- DONE: Focused functional checks all passed at the exact head:

  - `go test ./internal/status -run 'TestGateReadiness|TestBootIdentifyReadyGates|TestStatusProjectsSharedGateReadinessReducer|TestEnteredStage' -count=1`
  - `go test ./internal/gates -run 'TestCurrentStageReadiness|TestPrepare|TestWithdraw' -count=1`
  - `go test ./internal/ensigncycle -run 'TestRecordedGateLifecycle|TestRecordedGateReview' -count=1`
  - `go test ./internal/contractlint -count=1`

  The full status package also passed. These controls reproduce absent and
  prior-stage authority promotion, committed/path-clean exact-stage report proof,
  dirty/malformed rejection, selected open/withdrawn/closed authority, pending
  advance/merge, feedback/hold/consumed/superseded/not-applicable, ordinary and
  terminal stages, and invalid topology.

- DONE: The PR #612 shared-state manifest blocker is fixed on this exact target.
  `go test ./internal/gates -run 'TestV1PilotManifestReadsAndValidates$' -count=1 -v`
  passed every active and `_archive/` manifest entry, including the seven paths
  that were absent in the prior cycle (`bind-post-rework-briefing-at-rejection-
  regate.md`, `collapse-gate-approval-ceremony/index.md`,
  `cut-workflow-specific-round-recorder-from-v1/index.md`,
  `minimize-v1-gate-application-schema/index.md`,
  `shared-git-scaffold-helper.md`, `status-pagination-and-default-sorting.md`,
  and `status-where-robust-and-discoverable.md`).

- DONE: Complete offline and race suites passed without exceptions:
  `go test ./...` and `go test ./... -race`. All packages, including
  `internal/gates`, `internal/status`, `internal/ensigncycle`, `internal/dispatch`,
  `internal/cli`, and `skills/integration`, reported `ok`.

- DONE: Formatting and contract checks are clean. `gofmt -d ./cmd ./internal`
  produced no output, `git diff --check` passed, and
  `go test ./internal/contractlint -count=1` passed. The read-only `gofmt -d`
  check was used because validation is forbidden to edit implementation files.

- DONE: Re-ran the semantic adversarial matrix against production reducers and
  existing lifecycle oracles. All targeted status, gates, and recorded-gate
  controls passed:

  - status: `TestGateReadiness*`, `TestBootIdentifyReadyGates`,
    `TestStatusProjectsSharedGateReadinessReducer`,
    `TestEnteredStageProjectionRequiresCommittedCompleteReport`,
    `TestEnteredStageMutationControls`, `TestEnteredStageLegacySuppressionControls`,
    default-stage controls, and boot/next JSON parity;
  - gates: `TestCurrentStageReadiness*`, all prepare/withdraw/replay/refusal,
    eligibility, record/consume, round, duplicate/prototype/unknown-shape, digest,
    terminal-delivery, and application-shape controls;
  - recorded-gate lifecycle: `TestRecordedGateLifecycle*`, expected-actor,
    committed-before-dispatch, missing-event, phase-detection, gate-held, and
    guardrail controls.

  Exact output/order/cardinality, authority identity, path-clean byte refusal,
  terminal state, stale/withdrawn ownership, later-malformed masking, dirty target
  versus dirty sibling, and no-direct-preparation invariants all passed. The
  structural-pass/semantic-veto obligation remains First Officer judgment-owned;
  the declarative route requires one concrete `report-incomplete:` stop with zero
  prepare, legacy bind, mutation, presentation, idle, or repeat-next effects.

- DONE: AC-1 mechanical evidence is independently reproduced. The split-root
  status controls emit exactly one four-key `needs-preparation` row for both absent
  and prior-stage authority, an empty `dispatchable` array, and the ordered
  `command`/`dispatchable`/`ready_gates` next envelope; dirty and malformed reports
  emit no candidate. The registered Claude Opus recorded-gate lane also passed on
  this exact head (382.88s), providing one live cold-gate lifecycle proof. Sonnet
  and the remaining host lanes are not promoted to green below, so no complete
  multi-host live claim is made; the semantic-veto prose has no executable live
  control in this repository.

- DONE: AC-2 machine-envelope evidence passes in status tests and the JSON golden.
  `--next --json` has ordered `command`, `dispatchable`, and `ready_gates` keys;
  ready rows retain exactly `id`, `slug`, `current`, and `readiness`; and
  `needs-preparation` never enters `dispatchable`. Human `--next` output is
  unchanged. Ready-gate-first handling on initial and post-idle reads is present
  in the existing First Officer contract; no executable post-idle-created-gate
  lane exists, so that prose-only boundary is explicitly not over-claimed.

- DONE: AC-3 ownership/integrity controls pass. Open, withdrawn, stale, closed,
  pending, revise, hold, blocked, feedback, consumed, superseded, not-applicable,
  terminal, archived, malformed, and mismatched authorities remain on their
  existing reducers; prepare/withdraw/replay and whole-tree byte controls pass.
  No direct preparation, duplicate attempt, recorder mutation, or new owner is
  observable in the exact diff or tests.

- DONE: AC-4 scope audit passes. Status is the only readiness projection; s4 remains
  preparation/room owner, gqs remains report-evidence owner, 0m6 remains
  withdrawal/replacement owner, and the event loop remains declarative. The
  implementation surface contains no forbidden execution mechanism.

- FAILED (release evidence): Required local live coverage is intentionally not
  complete. The Sonnet command exited 0 as an explicit quarantine skip:
  `TODO(w5bfnrvpcphw857nzz93340c): Sonnet must reliably render the exact selected
  Briefing digest before re-enabling this journey`. Claude Opus passed as recorded
  above. A Codex lane had been started immediately before the First Officer's scope
  correction and was interrupted at ~32s (exit 1) when instructed not to launch
  additional local Codex/Pi lanes. Pi was not launched. Per current FO direction,
  Codex and Pi coverage is handed off to the registered exact-head CI lanes; skipped,
  aborted, or deferred lanes are not green evidence. No local PASS is claimed.

- MATERIAL EVIDENCE FINDING: The candidate's deterministic value behavior is green,
  but the required supported-host live matrix is incomplete in this validation
  environment (Sonnet quarantined skip, Codex intentionally aborted at scope
  correction, Pi deferred to CI). This blocks a complete validation PASS without
  indicating an offline outcome defect. Reopen/recommend PASS only after the
  registered CI lanes execute at this exact head and classify every skip/red result;
  live-test-truth sprint members remain out of scope for this entity.

### Completion checklist accounting

1. **DONE — exact-head checks and acceptance evidence.** Focused, full, race,
   manifest, formatting, dependency, acceptance, and adversarial matrix checks all
   pass at `cd76f52ab`; current `origin/main` is the same commit.
2. **FAILED — complete supported live coverage.** Opus passes; Sonnet is an
   explicit quarantine skip; Codex was aborted when the FO narrowed scope to CI;
   Pi was not launched. The required live result set is deferred to exact-head CI.
3. **DONE — report, scope, findings, and no-mutation accounting.** Only this state
   entity report is changed; the code worktree remains clean and no gate bytes were
   mutated.

Checklist totals: offline focused/full/race/manifest checks PASS; live subchecks 1
PASS, 1 SKIPPED, 1 ABORTED, and 1 DEFERRED-to-CI. **Recommendation: REJECTED/HOLD**
pending the registered CI live lanes; do not consume or present this validation
gate as a complete PASS.

### Summary

Fresh exact-head validation is green for the implementation's deterministic status,
gate ownership, scheduler envelope, manifest, full/race, and adversarial behavior.
Claude Opus also passes the recorded-gate lifecycle. The gate remains on HOLD because
the supported live matrix is intentionally incomplete: Sonnet is quarantined,
Codex was stopped when scope moved to CI, and Pi is deferred to CI. No implementation
or gate mutation was made; the next authority is the registered exact-head CI lanes.

## Stage Report: implementation (cycle 2)

- DONE: Reproduce both exact-head Codex lifecycle failures and submit the required four-field finding proposal before changing candidate bytes.
  Archived run `30930766742` attempt 1/job `92064915737` reproduced successor dispatch builds `2/1`; attempt 2/job `92071656858` reproduced trace `[prepare]`; candidate HEAD stayed `cd76f52a` until the First Officer authorized `fix`.
- DONE: After explicit First Officer disposition, make only the authorized correction that restores one-attempt gate and successor-dispatch behavior.
  Commit `a29401884513619aac8f3920e772adc18e357a9d` changes only the existing FO gate/dispatch prose: conn presentation continues immediately to record/consume, and successor build uses the bound adapter shape once.
- DONE: Run focused lifecycle controls and preserve the validated head scope; report any required scope expansion as Needs decision.
  Recorded-gate/contract controls, `go test ./...`, `go test ./... -race`, and the exact Codex live lifecycle all pass; the live oracle would fail on either a second build attempt or a prepare-only trace, and no scope expansion occurred.

### Authorized finding disposition

- Released user and normal workflow: Codex First Officer headless single-entity recorded-gate operation under an explicit delegated conn.
- Observable harm: one trace probed forbidden Codex bare mode before named dispatch; the retry stopped after presentation without recording, consuming, or dispatching.
- Affected value: `value-ac[AC-1]` requires the operation-coaching-free recorded-gate scenario to proceed exactly once without manufacturing semantic completion.
- Trigger evidence: exact candidate `cd76f52a`, run `30930766742`, attempt 1/job `92064915737` and attempt 2/job `92071656858`, with archived JSONL command traces.
- Materiality: Material.
- Ownership: this task's existing First Officer SK surface.
- Disposition: First Officer authorized `fix`; no binary, runner, parser, adapter, retry loop, lifecycle owner, accepted value, AC, or scope changed.

### Summary

The correction makes delegated recorded-gate operation self-guiding on Codex while preserving the existing lifecycle and runtime owners. Exact committed-head live evidence now reaches prepare, decision, consume, and one named successor dispatch in one pass.
