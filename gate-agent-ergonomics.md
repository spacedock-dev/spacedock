---
title: Make recorded gate operation self-guiding for First Officers
status: ideation
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
                state: pending
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
