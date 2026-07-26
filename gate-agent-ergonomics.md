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
    current:
        gate: gate:skwchfe30ac6ntr63j1g0txj:backlog
    records:
        - id: gate:skwchfe30ac6ntr63j1g0txj:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:skwchfe30ac6ntr63j1g0txj-backlog-1
              briefing:
                id: briefing:skwchfe30ac6ntr63j1g0txj:backlog:attempt-1:revision-1
                digest: sha256:49d1b8e3720bc54f0f9ef8db12d8ad33f7a277bdf156a9957b4a33d3a4134ddb
                digest-domain: canonical-bytes
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
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

## Problem

The original brief bundled package scaffolding, path normalization, lifecycle command
guidance, round identity, validation diagnostics, and readiness. Landed s4 now owns the
recorded gate procedure and its explicit bind, decision, commit, and consume barriers.
The provider-neutral room task owns room preparation. The only promoted defect in this
cycle is the readiness gap demonstrated by Roborev job 2826; the other ergonomics remain
deferred.

The supported cold-restart state is:

- the entity's current declared stage is `gate: true`;
- the latest current-stage report is complete, has evidence and a substantive Summary,
  and is clean in a commit;
- `gates.current` is absent or still names a prior-stage gate; and
- no current-stage Briefing has yet been selected.

The landed reducer returns `validating` for that state. `ready_gates` deliberately omits
`validating`, while ordinary `status --next` deliberately suppresses gate stages. The
First Officer therefore sees neither a ready gate nor a dispatchable entity and an
unscoped headless boot can stop even though gate preparation is its next action:

```text
complete committed gated-stage report
  + absent/prior-stage gates.current
  -> gate-readiness=validating
  -> ready_gates=[]
  + gate suppression from dispatchable
  -> no slug reaches engage or the gate lifecycle
```

Roborev job 2826 recorded the finding after commit `829308d8` in AgentsView session
`codex:019f9a58-f3c0-7653-8af9-e0dcd026ecb1` (`#0-2761 @2716-2728`), and the
top-level session confirmed and assigned the narrow fix separately from s4
(`codex:019f855d-8be8-7112-b467-7197851a2dc5`, `#7016-7169 @7145-7152`).
The retained source state and disposition also live in
`prepare-provider-neutral-gate-room/index.md` at state commits `6e6204bf`,
`ae881c47`, and `27473b4b`.

## Spike result

A removed real-CLI spike replayed both supported variants using the existing split-root
status fixture:

1. `status: validation`, no `gates:` block, and a committed complete
   `## Stage Report: validation`;
2. the same report and status, with `gates.current` still selecting an ideation gate.

For both, the current launcher returned `gate-readiness="validating"`,
`ready_gates=[]`, and `dispatchable=[]`. After recording the existing validation
Briefing, it returned `gate-readiness="awaiting-captain"` and exactly:

```json
[{"id":"task","slug":"task","current":"validation","readiness":"awaiting-captain"}]
```

`go test ./internal/status -run
TestSpikeCompletedGatedReportColdRestartIsInvisible -count=1 -v` passed both cases in
0.680s; the throwaway test was then deleted. AgentsView retains the exact exercise in
session `codex:019f9e6d-d524-7720-a4e3-235f5093fe05` (`#0-28 @21-24`).
This proves the report-to-selection boundary is the missing mechanism. It does not
justify a new gate command, room generator, or workflow engine.

## Proposed approach

Extend the one status-owned scheduler projection with one value:
`needs-preparation`. Reuse gqs's landed current-stage report shape/durability
classifier; do not add another report parser, epoch, ledger, frontmatter field, or
scheduler. The classifier supplies the mechanical candidate. At engage, the First
Officer still performs the existing semantic checklist, evidence, AC, and write-scope
review before any gate mutation. That judgment may veto a candidate and leave it
`validating`; it never manufactures completion.

The canonical reducer remains fail-closed:

| Current gated-stage evidence | `gate-readiness` | Scheduled? | Route |
|---|---|---:|---|
| No complete, path-clean committed current-stage report and no selected current-stage attempt | `validating` | no | Validation remains in progress. |
| Complete committed report; current-stage gate absent | `needs-preparation` | yes | Verify report semantics, then enter the existing prepare/bind owner. |
| Complete committed report; `gates.current` names a prior stage and no conflicting closed current-stage attempt exists | `needs-preparation` | yes | Verify, then bind/reselect the current-stage attempt. |
| Selected current-stage latest attempt is open with a valid Briefing | `awaiting-captain` | yes | Present/resume; never prepare a second attempt. |
| Selected current-stage approval is `advance/pending`, unblocked, and targets a nonterminal successor | `approved-awaiting-advance` | yes | Consume, commit, then hand the successor to ordinary scheduling. |
| The same approval targets a terminal successor | `approved-awaiting-merge` | yes | Consume, commit, then enter the merge ceremony. |
| Current-stage application is consumed, superseded, held, blocked, feedback-pending, or not-applicable | existing landed value | no | Existing resume/stop owner; never prepare. |
| Malformed state, a mismatched pointer hiding a closed/pending or consumed current-stage attempt, dirty/partial/failed report, or unknown topology | `invalid` or `validating` as appropriate | no | Fail closed without mutation. |

`approved-awaiting-advance` and `approved-awaiting-merge` are the two existing
renderings of the conceptual `approved-awaiting-application` state. No durable state
is added. Readiness derives from the current stage, the gqs report proof, the canonical
selected attempt, Resolution/application, and declared target taxonomy—not from the
status string or Stage Report heading alone.

### Status surfaces

- `status --boot --identify --json` keeps its four-key `ready_gates` rows and adds
  `needs-preparation` to the actionable set. Ordinary `--boot` remains unchanged.
- Machine-readable `status --next --json` appends the same `ready_gates` array after
  its existing `dispatchable` array. The existing five dispatch fields and ordering
  remain byte-identical, and human `--next` remains the dispatch table. This gives
  engage a fresh post-`state ready` scheduler read rather than relying on a stale
  greeting snapshot.
- `gate-readiness` through `--fields`/`--all-fields` uses the same reducer. The default
  status table remains unchanged.
- Gate rows use existing status ordering. `needs-preparation` never enters
  `dispatchable`; gqs continues to own only non-gated entered working-stage dispatch.

### Engage routing

After `state ready`, `state sweep`, and the startup hook, `«dispatch.next-action»`
reads one `status --next --json` envelope. It handles the first actionable
`ready_gates` row before ordinary dispatchables:

1. `needs-preparation`: load the existing gate lifecycle, read the latest report and
   reconstruct its checklist, confirm semantic completeness and the report-bearing
   commit's worker write scope, then use the existing preparation owner and
   `gate record --briefing`. Commit, re-read the scheduler, and require the same slug
   to be `awaiting-captain` before presentation.
2. `awaiting-captain`: present or resume the selected open Briefing.
3. `approved-awaiting-advance` / `approved-awaiting-merge`: preserve the existing
   commit-ancestry check and consume/application route.

A gate presentation remains an engage stopping condition. Interactive greet only names
the row. Headless boot follows the same engage route. An empty ready-gate array falls
through to ordinary dispatch and the existing idle double-check.

### Guards and idempotency

The readiness value is the routing guard. The First Officer must not decide preparation
from `status: <gate-stage>`, a report heading, a retained room, or “some gate exists.”
Before the bind it re-reads the entity/report; if the row no longer matches, it discards
the stale action and reruns the scheduler. Every refused route leaves entity bytes
unchanged.

| State when the scheduler is re-read | Preparation behavior |
|---|---|
| Current-stage gate absent + `needs-preparation` | One existing prepare/bind call creates and selects attempt 1; replaying the same Briefing is a byte no-op. |
| Prior-stage selection + current-stage open same-Briefing attempt + `needs-preparation` | Existing recorder repairs only `gates.current`; attempt ids/counts and bytes outside `gates` stay fixed; replay is a byte no-op. |
| Selected current-stage open attempt | No prepare call. Present the exact selected Briefing; a routing mutant that binds again fails. |
| Selected current-stage closed pending approval | No prepare call. Route consume; preparing here would supersede valid authority and is a blocking defect. |
| Consumed current-stage application | No prepare or repeat consume. The row is omitted; successor scheduling/merge owns the post-consume state. |
| Malformed, blocked, held, feedback, dirty/partial/failed report, or conflicting mismatched closed attempt | No prepare call and no entity mutation. Surface the landed diagnostic/owner. |

This task does not alter `gate record --briefing` semantics. Its atomic validator and
entity lock remain the mutation boundary. The fix ensures the FO calls it only on the
two preparation states where its existing idempotency is correct.

### Mechanism-to-value trace

| New mechanism | Value AC served | Simplest alternative | Why the alternative is insufficient |
|---|---|---|---|
| Promote gqs's committed-report proof to `needs-preparation` in the canonical reducer | AC-1 | Add another s4 instruction saying to bind after completion | Cold boot has no entity slug on which to execute that instruction. |
| Append the shared ready rows to `status --next --json` | AC-1, AC-2 | Reuse only the greeting's boot snapshot | `state ready` can refresh split-root state after greet; engage needs a current scheduling read. |
| Route `needs-preparation` before ordinary dispatch with the state matrix above | AC-1, AC-3 | Infer the route from `status: <gate-stage>` | It cannot distinguish unfinished validation, an open Briefing, or approval pending consumption. |
| Reuse gqs's report classifier and s4's lifecycle instead of introducing gate-specific copies | AC-4, serving AC-1 | Add a gate-report parser and a combined lifecycle command | That creates a second readiness authority and duplicates already exercised mutations. |

## Alternatives considered

- **More s4 prose only:** insufficient because the unscoped boot has no slug to which
  it could apply the prose. Roborev 2826 is the counterexample.
- **Treat every gate-stage status as ready:** reintroduces the landed 5/3
  false-positive failure and can present unfinished validation.
- **Parse report prose in a second gate scheduler:** duplicates gqs and creates two
  readiness authorities. Reusing its structural/durability classifier is sufficient.
- **A new `gate next`, `gate prepare`, room scaffold, or lifecycle transaction:**
  broader than the demonstrated defect. Landed gate record/application commands and
  the provider-neutral room owner already cover the mutations.
- **Keep readiness boot-only:** insufficient after `state ready`; engage needs a fresh
  scheduler read and otherwise can repeat the cold stall on remotely refreshed state.

## Acceptance criteria

**AC-1 (VALUE)** — Completed cold gates are discovered and operated, not parked.
For both supported cold states (absent `gates:` and prior-stage `gates.current`), a
complete path-clean committed gated-stage report yields exactly one
`needs-preparation` row: 2/2 true positives versus the current 0/2 baseline. An
unscoped headless run through the existing `recorded-gate-lifecycle` fixture then
records exactly one Briefing bind, one bound-review presentation, one authorized
decision, one consume, and one durable successor effect instead of stopping
quiescent. **Verified by:** the two-case native fixture plus the existing registered
Claude/Codex/Pi recorded-gate scenario with its prompt stripped of gate-operation
coaching; removing the report promotion returns `ready_gates=[]` and the live durable
marker is absent.

**AC-2** — One reducer distinguishes validation, preparation, Captain wait, and
approved application on boot and engage. The same mixed fixture reports
`validating`, `needs-preparation`, `awaiting-captain`,
`approved-awaiting-advance`, and `approved-awaiting-merge` from independent state
variants. `status --boot --identify --json` and `status --next --json` emit identical
ordered ready rows; next's existing `dispatchable` bytes and human output are
unchanged. **Verified by:** raw JSON/table assertions over the shared fixture and a
stage-only mutant, which falsely schedules unfinished validators and fails the exact
population.

**AC-3** — Preparation is guarded and idempotent across adjacent durable states.
Absent and mismatched preparation create/select at most one open attempt; identical
replay is byte-clean. Open, closed-pending, consumed, malformed, held, blocked,
feedback, and incomplete/dirty-report states invoke no preparation. In particular a
pending approval is neither superseded nor replaced, and consumed authority is never
spent twice. **Verified by:** real CLI command-log cardinality, attempt identity/count,
whole-entity before/after bytes, and the existing same-Briefing recorder fixture; a
route mutant that binds on open/closed/consumed state changes bytes or event counts and
fails.

**AC-4** — The correction composes with s4 and gqs without another engine. Gqs remains
the sole owner of non-gated entered-stage report readiness, s4 remains the sole gate
lifecycle owner, and this change adds no schema/frontmatter, command, provider
behavior, room scaffold, compatibility path, report epoch, ledger, lease, runtime
adapter, or harness. **Verified by:** existing entered-stage controls, recorded-gate
offline mutants, full/race suites, and the README path-to-lane gate. Adding a second
report parser/state store or changing non-gated dispatch bytes is a design failure even
if focused gate tests pass.

## Test plan

Implementation starts with the current red cold-restart fixture and the live prompt
de-coaching, before changing readiness or FO instructions.

Focused offline:

```bash
go test ./internal/status -run 'TestBoot.*ReadyGate|Test.*GateReadiness|TestEnteredStage' -count=1
go test ./internal/gates -run 'TestCurrentStageReadiness|TestSameBriefingBind' -count=1
go test ./internal/ensigncycle -run 'TestRecordedGateLifecycle|TestRecordedGateReview' -count=1
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

The status matrix reuses `buildSplitRoot`, the gqs report classifier fixtures, and the
recorded-gate entity/room. It covers absent, prior-stage selection, selected open,
closed pending to both target types, consumed, malformed, blocked, held, feedback,
partial/failed/dirty report, archived, terminal, and ordinary non-gate controls. It
asserts exact JSON, event order/cardinality, attempt identity, and whole-file bytes;
there is no new harness.

Because the diff changes the host-neutral FO contract and gate skill, all registered
live lanes are required after offline green:

```bash
SPACEDOCK_LIVE_MODEL=sonnet go test -tags live -count=1 -timeout 40m \
  -run 'TestLiveClaudeSharedScenarios/recorded-gate-lifecycle$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_MODEL=claude-opus-4-8 go test -tags live -count=1 -timeout 40m \
  -run 'TestLiveClaudeSharedScenarios/recorded-gate-lifecycle$' ./internal/ensigncycle -v
go test -tags live -count=1 -timeout 40m \
  -run 'TestLiveCodexSharedScenarios/recorded-gate-lifecycle$' ./internal/ensigncycle -v
go test -tags live -count=1 \
  -run 'TestSharedScenarioRunnerCoverage|TestPiSharedScenarioCoverage|TestLivePiRecordedGateLifecycle' \
  ./internal/ensigncycle -v
```

Those focused runs are preflight, not the merge gate. The complete registered
`claude-live` matrix (both model legs), `codex-live`, and `pi-live` jobs must each run
and pass; an unapproved, skipped, missing-auth, timed-out, or red job is not green.
The live oracle grades command order and durable state, never a transcript claim.

## Expected surface and estimates

These LOC estimates are non-authoritative planning aids; the exact file list and
behavioral boundaries are authoritative. Expected total is about 460 touched LOC,
with a ±30% review tolerance (322–598). Crossing it requires a design check, not an
AC waiver.

| File | Estimated touched LOC | Purpose |
|---|---:|---|
| `internal/gates/model.go` | 33 | Extend the pure reducer with the preparation candidate while preserving fail-closed lifecycle states. |
| `internal/gates/gates_test.go` | 75 | Exact reducer and conflicting-mismatch matrix. |
| `internal/status/entered_stage.go` | 25 | Reuse/factor gqs's report shape and path-durability classifier; no second parser. |
| `internal/status/discover.go` | 13 | Feed report readiness to the canonical gate reducer. |
| `internal/status/format.go` | 7 | Include `needs-preparation` in the existing ordered actionable set. |
| `internal/status/json_commands.go` | 15 | Append shared ready rows to machine-readable `--next`. |
| `internal/status/boot_identify_test.go` | 125 | Two-case replay, boot/next parity, negative matrix, unchanged dispatch bytes. |
| `internal/status/gates_coexist_test.go` | 23 | Human/JSON readiness vocabulary. |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | 60 | Remove operation coaching and require discovery plus exact route cardinality. |
| `skills/fo-gate-lifecycle/SKILL.md` | 28 | Add preparation-required guards and resume matrix. |
| `skills/first-officer/references/first-officer-shared-core.md` | 13 | Route fresh engage scheduling through ready gates first. |
| `skills/first-officer/references/fo-dispatch-core.md` | 18 | Consume the combined machine scheduler envelope. |
| `docs/site/concepts/gates-and-decisions.md` | 17 | Explain completed-report recovery and guarded preparation. |
| `docs/site/reference/command-reference.md` | 8 | Document `needs-preparation` and `--next --json` ready rows. |

No other file is intended. In particular, do not edit gate CLI/recorder/application
semantics, provider code, host adapters, workflow schema, room scaffolding, or add a
new test runner. The implementation must start after gqs lands so
`internal/status/entered_stage.go` is a real integration target; if its landed shape
cannot supply the classifier, return to the gate rather than duplicate it.

## Documentation diff

In `docs/site/concepts/gates-and-decisions.md`, replace:

> After completion verification, a gate with no selected attempt remains
> `validating`. The first officer binds and commits the retained Briefing before
> presenting anything.

with:

> A gate with no selected attempt remains `validating` until its current-stage report
> is complete and committed. It then becomes `needs-preparation` on both boot and the
> machine scheduler read. Engage verifies the report semantically, binds and commits
> the retained Briefing exactly once, re-reads `awaiting-captain`, and only then
> presents. Open, pending-approved, and consumed gates never re-enter preparation.

In `docs/site/reference/command-reference.md`, change the status row:

```diff
- `--next`, ... `--boot` (with `--identify`, ... canonical ready-gate scheduling rows ...)
+ `--next` (JSON appends the same canonical `ready_gates` rows after the unchanged
+ `dispatchable` array), ... `--boot` (with `--identify`, ... canonical ready-gate
+ scheduling rows, including `needs-preparation` only for a complete committed
+ gated-stage report ...)
```

## Out of scope

- Gate-room/package scaffolding, path normalization, new Briefing/request identity
  derivation, provider behavior, and advisory-round preparation.
- New validation-error UX, lifecycle command generation, or descendant-commit fields.
- Changing Briefing/decision/consume authority separation, application eligibility,
  merge behavior, feedback routing, or s4's review/presentation ownership.
- Expanding gqs beyond non-gated entered stages or adding same-stage epochs.
- A generic workflow driver, scheduler service, durable queue, compatibility layer, or
  provider-specific branch.

## Stage Report: ideation

- DONE: Replay the exact Roborev 2826 cold-restart state: current gated stage has a semantically complete committed report, while gates.current is absent or names a prior stage; show why current boot/engage cannot discover it.
  AC-1 evidence: the removed two-case real-CLI spike returned validating + empty ready_gates/dispatchable before bind and awaiting-captain afterward; AgentsView sessions and retained state commits identify the original finding and assignment.
- DONE: Reconcile the original gate-agent-ergonomics brief against landed/current s4: treat mechanical preparation and room scaffolding as owned elsewhere and scope this cycle to first-class readiness projection plus routing only.
  Proposed scope adds only needs-preparation and engage routing; gate commands, rooms, provider, rounds, diagnostics, and the rest of the original brief are deferred.
- DONE: Define the smallest single-scheduler state model that distinguishes validating, awaiting Captain, and approved-awaiting-application without inferring from status text alone.
  AC-2 evidence: one reducer combines gqs report proof with canonical gate/application state; existing approved target variants remain, and no durable state or second scheduler is added.
- DONE: Specify status --boot/--next/engage behavior, exact byte-clean guards, and preparation idempotency for absent, mismatched, open, closed-pending, and consumed current-stage gates.
  AC-3 evidence: the status, engage, and six-row guard matrices define exact routes, replay behavior, mutation exclusions, and whole-byte/event-count falsifiers.
- DONE: Use existing fixtures and registered live lanes; add no new harness, compatibility layer, provider behavior, or generic workflow engine.
  AC-4 evidence: the plan reuses buildSplitRoot, gqs controls, recorded-gate lifecycle, and every applicable registered host lane.
- DONE: Declare exact intended files and non-authoritative estimates, append a complete ideation Stage Report with AC evidence, and commit the state path only.
  Fourteen intended files, 460 touched LOC ±30%, exact docs changes, four testable ACs, and the path-scoped state-only handoff are recorded.

### Summary

Ideation reduces the broad ergonomics seed to the one demonstrated restart failure.
A single new readiness value makes a complete committed gated-stage report
discoverable, while existing s4/gqs owners retain every mutation and completion rule.
