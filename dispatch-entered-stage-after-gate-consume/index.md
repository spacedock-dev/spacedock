---
title: Dispatch the entered working stage after gate consumption
status: ideation
source: "Codex live run 30197794474 on 2026-07-26: rejection-flow began at implementation, but the FO advanced directly to validation without an implementation worker/report; the strict two-cycle provenance assertion correctly failed."
started: 2026-07-26T10:57:18Z
completed:
verdict:
score: 1.0
worktree:
issue:
sprint: durable-decisions
id: gqsw81ghf48hr2n3jg6k7nx8
gates:
    version: 1
    current:
        gate: gate:gqsw81ghf48hr2n3jg6k7nx8:backlog
    records:
        - id: gate:gqsw81ghf48hr2n3jg6k7nx8:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:gqsw81ghf48hr2n3jg6k7nx8-backlog-1
              briefing:
                id: briefing:gqsw81ghf48hr2n3jg6k7nx8:backlog:attempt-1:revision-1
                digest: sha256:4a74b1208239ddf0168759fe2f50fa4bb2a02e1329133d29a83aa7e455a7ed47
                digest-domain: canonical-bytes
                request-digest: sha256:2e8c911fa6c4c2dbea158e9a20cee9bfa16054eada3d2b9774b92e76bba2d053
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:gqsw81ghf48hr2n3jg6k7nx8:backlog:1
                briefing: briefing:gqsw81ghf48hr2n3jg6k7nx8:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T10:56:42.99005Z"
                decision: approve
                reason: The supported live journey lost its first implementation round because stage entry was credited without a worker/report; this must be shaped before the sprint's assembled walkthrough.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

A gate application can atomically move a ticket into a non-gated working stage, but the First Officer must still dispatch that entered stage before advancing again. In Codex run `30197794474`, `status --boot --identify --json` exposed the fixture as `current=implementation,next=validation`; the FO followed that projection, ran `status ... status=validation started`, and built validation without an implementation worker. The final entity carried the first REJECTED validation, one rework implementation report, and the second PASSED validation, but no original implementation report. The strict two-cycle assertion correctly failed with “left 1 implementation reports, want at least 2.”

The relevant FO lifecycle and Codex adapter prose at source commit `8c9aa160` already said that every ready entity, including one advanced by approval, needs an observed worker spawn and that `dispatch build` is not dispatch. Those files are unchanged through the current checkout. The failure is therefore not missing exhortation: the one scheduling projection named the successor when the current working stage had never run.

## Boundary to shape

Ship one host-neutral first-entry invariant in the existing status scheduler:

> A declared non-initial, non-gated, non-terminal current stage that has not produced a durable, conforming, semantically complete current-stage report is the dispatch target itself. No `status --set` may change its status away from that stage. Once completion is proven, ordinary successor projection resumes.

`status --next` and boot continue to use the same `dispatchAnalysis` source. In the entered/unreported case their existing five-field row becomes `current=<stage>,next=<same-stage>` and uses that stage's existing concurrency/worktree settings. The normal dispatch recipe can therefore set the same status idempotently, build that stage, and spawn it without a new field or host branch. A durable, structurally complete current-stage report changes the row back to `current=<stage>,next=<successor>`, subject to the FO's semantic completion veto below.

The one predicate has two mechanical halves:

1. **Report shape and durability.** Select only the latest report whose stage token is the exact current stage. A heading alone, an empty section, a blank checklist item, an item with no evidence/rationale, a `FAILED` item, a missing/empty `### Summary`, a wrong-stage report, an older valid report masked by a later malformed current-stage section, or entity bytes not present cleanly in the local Git commit is incomplete. The parser reuses the shipped report selector and checklist ranges; it adds no Markdown implementation. Path-scoped Git cleanliness means unrelated sibling dirt does not block recovery.
2. **Transition guard.** Every `status --set` containing a status value different from the current stage shares that readiness result and refuses before mutation. Direction is irrelevant: successor hops, backward hops, terminal jumps, repeated/chained `status=` updates, and `--force` all refuse byte-clean. A same-stage dispatch mutation such as `status=<current> started` remains allowed; unrelated non-status updates are outside this guard.

Structural validity is necessary but not the whole completion judgment. The First Officer contract must distinguish these two cases:

- **Live completion:** the FO has an active dispatch epoch, so it MUST observe that host adapter's `«completion-signal»` and then verify a durable current-stage report against the retained dispatched checklist. Every item is accounted for with non-empty evidence or rationale, no item is `FAILED`, and the summary is substantive. A report without the runtime signal cannot complete a live worker.
- **Cold recovery:** a crash may erase the ephemeral runtime signal and worker handle. On boot, a conforming, semantically complete current-stage report already committed under the worker's body/report write scope is sufficient recovery evidence. The FO reconstructs the stage checklist from the entity and stage definition, confirms the report-bearing commit is within the ensign's assigned body/report/artifact scope, performs the same semantic review, and advances without dispatching that already-completed stage again. If the report is absent, dirty, malformed, partial, failed, or stale, completion is vetoed and the current stage is repaired/dispatched once. Report verification may veto the scheduler's row; it never selects a different stage.

This is deliberately first-entry-only. Current workflow topology enters an unvisited working stage after a gate; feedback recovery dispatches its target directly and does not use `status --next`. A report from a prior visit therefore cannot satisfy this invariant in the supported path. A same-stage re-entry epoch is explicitly deferred until either a gate can target a previously visited stage or feedback recovery begins using `status --next`; that future capability needs an epoch/visit identity rather than pretending the latest historical report is fresh. Initial seed stages preserve legacy successor projection, gate and terminal stages retain suppression, and a set `worktree` retains existing in-flight suppression before this rule is considered.

### Spike result

The exact archived `runtime-live-e2e-codex-live` artifact was replayed. Its first mutation was `implementation -> validation`; no initial implementation dispatch preceded it, and the final entity contained exactly one implementation report (`cycle 2`). A throwaway real-CLI test over `writeRecordedGateFixture` then consumed validation into handoff and observed:

- before a handoff report, both `status --next --json` and `--boot --identify --json` returned exactly `current=handoff,next=done`;
- `status --read ... --stage handoff --checklist --json` alone failed with “no Stage Report”;
- after a conforming handoff report, the same status/boot projection was byte-identical while the report read succeeded.

The throwaway file was removed. This falsifies “consume stdout plus current prose is sufficient”: stdout is ephemeral, and a heading-only report read is too weak to credit completion. The status projection lacks a durable cold-boot readiness signal and must own the structural/durability half; the FO completion path remains responsible for host signal provenance and checklist semantics.

### Proposed change

Add one shared read-only predicate for “entered current stage awaiting first completion proof.” It reads the existing latest-stage-report/checklist primitives plus a literal path-scoped Git cleanliness check. `dispatchAnalysis` uses it to select the current stage instead of its successor; `runSet` uses the same predicate, before any force-bypass guards, to reject every away-status update. The refusal must leave the entity byte-identical and emit no success stdout.

Amend `first-officer-shared-core.md` at Completion and Gates with the live-signal-versus-cold-report ruling and the semantic report checklist. Amend `fo-dispatch-core.md` so `current==next` means dispatch the entered stage, while a cold recovered `current!=next` row is advanced only after report verification; neither core may manufacture a completion signal or report. No host runtime adapter changes are required because each adapter already binds its own completion signal and all hosts share the same durable report verification. Add the behavior assertion before changing this FO command text, per the repository's skill-smoke rule.

Do not change `gate consume`, gate schema, entity frontmatter, dispatch-build packages, runtime adapters, ensign protocol, feedback routing, or the strict rejection assertion.

Cheaper alternatives considered:

- More lifecycle prose or relying on `gate consume target-stage=...` cannot recover after a crash and was already present in the failing trace.
- Treating any matching heading as completion would duplicate the original bug with an empty, partial, dirty, failed, or wrong-cycle report. Report shape and path durability must be part of the one scheduler predicate.
- A durable dispatch ledger/lease could model every crash window, but the requested consume-before-spawn and report-present recovery needs no new state; existing worktree/roster behavior remains responsible for already-in-flight workers.
- Persisting runtime completion signals or checklist digests would widen worker/FO write scope and create a second scheduler state. The existing host signal plus committed report and deterministic FO checklist reconstruction are sufficient for this first-entry boundary.

## Acceptance criteria

**AC-1 (VALUE)** In a real two-cycle rejection journey, the durable ticket contains the original implementation report, first REJECTED validation, rework implementation report, and second PASSED validation in order; the existing strict rejection assertion remains unchanged and passes. **Test:** run only the real Codex `rejection-flow` scenario and inspect its archived entity/JSONL; removing current-stage projection or the mutation guard reproduces run `30197794474` and fails on one implementation report.

**AC-2 (DISPATCH)** After gate consume or equivalent first entry reaches a non-initial working stage, `status --next` and boot name that same stage until completion proof is durable. Every away-status mutation—including successor, backward, terminal, repeated/chained update, and `--force` forms—is refused byte-clean, while same-stage dispatch mutation is allowed. An observed worker spawn, the host completion signal, and an FO-validated committed report precede a live transition; build output, `started`, narration, an FO-authored report, and a syntactic heading do not satisfy the oracle. **Test:** a real-CLI table snapshots bytes around each direction/force/chained mutation and a focused live trace orders implementation spawn and completion before the first validation transition.

**AC-3 (RECOVERY)** A cold boot after consume but before spawn exposes exactly one row with `current=<target>,next=<target>`. A heading-only, empty, structurally partial, failed, wrong-stage, later-malformed-over-older-valid, or uncommitted report does not change that row and cannot unlock mutation. A conforming, semantically complete, path-clean committed report is sufficient cold-recovery evidence: boot exposes the ordinary successor, the FO revalidates it against the reconstructed checklist and worker write scope, and no duplicate target worker is spawned. A structurally complete report that omits a reconstructed checklist obligation remains vetoed by the FO and is never credited as recovery. **Test:** extend the recorded-gate real-CLI fixture with committed malformed-report and dirty-valid-report mutants plus a committed-complete control; each exact JSON/bytes assertion fails if report validation is weakened.

**AC-4 (SCOPE)** The correction composes with consume, `status --next`, boot parity, initial/gate/terminal/worktree suppression, standing dispatch, fresh/reuse, and direct feedback routing without new scheduler state, report epochs, frontmatter, or fixture/runtime branches. Same-stage re-entry remains explicitly out of scope until topology or feedback scheduling can reach it. **Test:** table controls keep initial-stage successor projection, gate/terminal suppression, worktree-set suppression, terminal consume, and report-present direct feedback behavior byte-identical; full and race suites plus adversarial predicate mutations turn focused tests red.

## Test plan

Implementation starts with the status and live-scenario behavior assertions, before editing the FO references, then adds the shared predicate, projection, guard, and contract text. The existing `TestRejectionFlowNegativeSingleCycle` and `assertRejectionFlow` are preserved byte-for-byte; they already reject the exact one-implementation-report end state.

Focused offline commands:

```bash
go test ./internal/status -run 'TestEnteredStage|TestBootJSONDispatchableMirrorsNext' -count=1
go test ./internal/ensigncycle -run 'TestRecordedGateLifecycleRealCLIReplay|TestRecordedGateLifecycleTerminalConsumeHasNoDispatchableSuccessor|TestRejectionFlowNegativeSingleCycle' -count=1
go test ./...
go test ./... -race
```

Focused live command (two existing scenarios, no new runner or scheduler):

```bash
go test -tags live -count=1 -timeout 20m -run 'TestLiveCodexSharedScenarios/(rejection-flow|recorded-gate-lifecycle)$' ./internal/ensigncycle -v
```

The offline focus is expected to cost seconds and no model calls; the focused live pair is expected to cost 6–10 minutes and two Codex scenario calls. Falsifiers are load-bearing:

- Make the predicate accept only the report heading: the heading-only, empty, missing-evidence/summary, `FAILED`, and dirty-report cases must turn red.
- Skip the literal path cleanliness check: the uncommitted-report case must turn red while the unrelated-dirty-sibling control stays green.
- Restore successor-only comparison or let `--force` bypass: the backward-hop and same-command chained-bypass byte snapshots must turn red.
- Remove live spawn/completion ordering: the recorded-gate live oracle must reject the trace; remove first-stage dispatch and the unchanged rejection oracle must again report one implementation report.

### Expected surface

Exact baseline, measured as added/deleted LOC:

| File | Expected LOC | Purpose |
|---|---:|---|
| `internal/status/entered_stage.go` | +110/-0 | Shared first-entry readiness, report semantics, and literal path-durability predicate |
| `internal/status/format.go` | +25/-6 | Select current vs successor inside existing `dispatchAnalysis` |
| `internal/status/handlers.go` | +35/-0 | Refuse every away-status mutation before force-bypass guards |
| `internal/status/entered_stage_test.go` | +250/-0 | Projection, semantics/durability matrix, all-direction guard, boot parity, and unchanged controls |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +85/-5 | Real consume/cold-boot/malformed/committed recovery and skill-trace ordering smoke |
| `skills/first-officer/references/first-officer-shared-core.md` | +22/-4 | Bind live completion and cold-report recovery semantics |
| `skills/first-officer/references/fo-dispatch-core.md` | +18/-5 | Interpret current=current dispatch and recovered successor without duplicate spawn |
| `docs/site/concepts/gates-and-decisions.md` | +6/-1 | Document recovery rows, durable report rule, and unbypassable guard |

Expected total: **551 additions, 21 deletions (572 touched LOC)**, tolerance **±25% touched LOC (429–715)**. The two shared FO references are required; `SKILL.md`, ensign instructions, and host runtime adapters are not. Any new production package/file beyond `entered_stage.go`, schema/frontmatter field, persisted signal/checklist/epoch, fixture-only branch, runtime adapter edit, or scheduler/lease state is a design deviation requiring gate reconfirmation.

### Documentation diff

In `docs/site/concepts/gates-and-decisions.md`, replace:

> Approval then uses `gate consume`, which rechecks eligibility and atomically writes the successor stage and consumed mark; the consumed descendant commit lands before ordinary successor dispatch.

with:

> Approval then uses `gate consume`, which rechecks eligibility and atomically writes the successor stage and consumed mark. Until that first-entered working stage has a durable, complete Stage Report, `status --next` and boot name it as both `current` and `next`, and every away-status `status --set`—backward or forward, even with `--force`—is refused. Live completion requires the runtime completion signal plus the verified report; after a cold restart, the committed report alone is sufficient only when the first officer reconstructs the checklist and verifies every item and summary. The consumed descendant commit therefore lands before exactly one recoverable successor dispatch.

## Stage Report: ideation

- DONE: Replay the exact skipped-first-implementation trace and test the cheapest existing lifecycle/status correction first.
  Run 30197794474 was replayed; a removed real-CLI spike proved boot/next stay on the successor while only a separate report read distinguishes pre/post report.
- DONE: Define one host-neutral dispatch-before-next-transition invariant with cold-boot recovery and no duplicate worker.
  The shared current-stage/report predicate drives both existing dispatch projection and a byte-clean transition guard; no new state, scheduler, host branch, or fixture exception is introduced.
- DONE: Declare exact files and LOC plus offline and focused-live falsifiers without a second scheduler or fixture special case.
  Six files, 279 touched LOC ±25%, exact commands, mutation falsifiers, and the concrete site-doc replacement are recorded above.

### Summary

Ideation chooses the smallest status-owned correction because the archived trace already contained the lifecycle prose and the spike proved the durable projection was ambiguous. The design preserves the strict two-cycle assertion, adds no dispatch state, and makes the consumed/entered stage recoverable through the one scheduler the FO already follows.

## Stage Report: ideation (cycle 2)

- DONE: Correct the lifecycle ruling for live completion and cold recovery, including semantic report validation.
  The design now requires a host completion signal plus committed report live, permits a committed semantically complete report after cold boot, and rejects heading-only, empty, partial, failed, dirty, or stale evidence.
- DONE: Guard every status change away from an entered/unreported stage and add byte-clean bypass controls without widening re-entry scope.
  Forward, backward, terminal, repeated/chained, and `--force` status updates share one pre-mutation guard; same-stage dispatch remains allowed and same-stage re-entry epochs are explicitly deferred.
- DONE: Reconcile the exact implementation surface and required First Officer contract changes.
  Eight files and 572 touched LOC ±25% are declared; both shared FO references are required, while runtime adapters, ensign instructions, schemas, ledgers, and the strict rejection assertion remain unchanged.

### Summary

Cycle 2 incorporates the binding staff ruling without adding product implementation. The corrected design keeps status as the sole stage selector, makes recovery proof durable and falsifiable, and gives the First Officer an explicit host-neutral rule for live provenance versus crash recovery.
