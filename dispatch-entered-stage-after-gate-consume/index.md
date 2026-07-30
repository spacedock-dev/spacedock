---
title: Dispatch the entered working stage after gate consumption
status: done
source: "Codex live run 30197794474 on 2026-07-26: rejection-flow began at implementation, but the FO advanced directly to validation without an implementation worker/report; the strict two-cycle provenance assertion correctly failed."
started: 2026-07-26T10:57:18Z
completed:
verdict:
score: 1.0
worktree: .worktrees/spacedock-ensign-dispatch-entered-stage-after-gate-consume
issue:
sprint: durable-decisions
id: gqsw81ghf48hr2n3jg6k7nx8
gates:
    version: 1
    current:
        gate: gate:gqsw81ghf48hr2n3jg6k7nx8:validation
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
        - id: gate:gqsw81ghf48hr2n3jg6k7nx8:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:gqsw81ghf48hr2n3jg6k7nx8-ideation-1
              briefing:
                id: briefing:gqsw81ghf48hr2n3jg6k7nx8:ideation:attempt-1:revision-1
                digest: sha256:2ce71b3cd6b65d989346a0ab81180b6e43bd5a42c3a82ffdda5ba880187e8a00
                digest-domain: canonical-bytes
                request-digest: sha256:303aa4c963d300d0a06d86d1d5fad040cffd12b29c8661a8c8b3ffc3f4cfffae
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:gqsw81ghf48hr2n3jg6k7nx8:ideation:1
                briefing: briefing:gqsw81ghf48hr2n3jg6k7nx8:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T12:17:23.586877Z"
                decision: approve
                reason: AC-1 through AC-4 are evidenced and staff corrections are closed. Record approval now; apply after s4 lands to avoid overlapping shared lifecycle proof surfaces.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:gqsw81ghf48hr2n3jg6k7nx8:validation
          stage: validation
          attempts:
            - id: gate-attempt:gqsw81ghf48hr2n3jg6k7nx8-validation-1
              briefing:
                id: briefing:gqsw81ghf48hr2n3jg6k7nx8:validation:attempt-1:revision-1
                digest: sha256:8dac2e7f8bd99cfb9d84fac2446fac12fc318fc571c22c6c4ee7baad2486d203
                digest-domain: canonical-bytes
                request-digest: sha256:9e8cde28a37107216caafc407f253db2014a9d827a871c344d75c6330b78040a
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:gqsw81ghf48hr2n3jg6k7nx8:validation:1
                briefing: briefing:gqsw81ghf48hr2n3jg6k7nx8:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T14:12:33.606877Z"
                decision: revise
                reason: 'Captain narrowed pre2 to the provider-free journey: retain worktree-safe entered-stage actionability, revert merge-finalization changes, and remove the cold checklist-omission promise.'
              application:
                action: feedback
                target-stage: implementation
                state: superseded
            - id: gate-attempt:gqsw81ghf48hr2n3jg6k7nx8-validation-2
              briefing:
                id: briefing:gqsw81ghf48hr2n3jg6k7nx8:validation:attempt-2:revision-1
                digest: sha256:a0344d898b3119b92251a00735af627ba1b2b8d5808fa25b554407b1b8617c77
                digest-domain: canonical-bytes
                request-digest: sha256:3790c53f7931f353153cdec7d9a5b8290eecc0381a2d83c4e7f30a09d3b5524c
                room-ref: ./review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:gqsw81ghf48hr2n3jg6k7nx8:validation:2
                briefing: briefing:gqsw81ghf48hr2n3jg6k7nx8:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-07-30T14:48:03.888476Z"
                decision: approve
                reason: Exact-tip cb01129b passes the narrowed provider-free post-consume journey, focused/full/race checks, and conflict-drift inspection with unrelated semantics absent.
              application:
                action: advance
                target-stage: done
                state: consumed
                blockers: []
review-round:
    id: round:gqsw81ghf48hr2n3jg6k7nx8:implementation:1
    stage: implementation
    cycle: 1
    briefing:
        id: briefing:gqsw81ghf48hr2n3jg6k7nx8:implementation:round-1
        digest: sha256:2008c595794b500ed2275828817fccfd6a0fc4a4a47649d586d971f273560e4f
        digest-domain: canonical-bytes
        room-ref: ./review/implementation/round-1
---

A gate application can atomically move a ticket into a non-gated working stage, but the First Officer must still dispatch that entered stage before advancing again. In Codex run `30197794474`, `status --boot --identify --json` exposed the fixture as `current=implementation,next=validation`; the FO followed that projection, ran `status ... status=validation started`, and built validation without an implementation worker. The final entity carried the first REJECTED validation, one rework implementation report, and the second PASSED validation, but no original implementation report. The strict two-cycle assertion correctly failed with “left 1 implementation reports, want at least 2.”

The relevant FO lifecycle and Codex adapter prose at source commit `8c9aa160` already said that every ready entity, including one advanced by approval, needs an observed worker spawn and that `dispatch build` is not dispatch. Those files are unchanged through the current checkout. The failure is therefore not missing exhortation: the one scheduling projection named the successor when the current working stage had never run.

## Boundary to shape

Ship one host-neutral first-entry invariant in the existing status scheduler:

> A declared non-initial, non-gated, non-terminal current stage that has not produced a durable, conforming current-stage report is the dispatch target itself. Once its normal same-stage dispatch owns a worktree, no `status --set` may change its status away from that stage until the report is durable. Once completion is proven, ordinary successor projection resumes.

`status --next` and boot continue to use the same `dispatchAnalysis` source. In the entered/unreported case their existing five-field row becomes `current=<stage>,next=<same-stage>` and uses that stage's existing concurrency/worktree settings. The normal dispatch recipe can therefore set the same status idempotently, build that stage, and spawn it without a new field or host branch. A durable, structurally complete current-stage report changes the row back to `current=<stage>,next=<successor>`; the active live worker still completes through the existing First Officer signal-and-report contract below.

The one predicate has two mechanical halves:

1. **Report shape and durability.** Select only the latest report whose stage token is the exact current stage. A heading alone, an empty section, a blank checklist item, an item with no evidence/rationale, a `FAILED` item, a missing/empty `### Summary`, a wrong-stage report, an older valid report masked by a later malformed current-stage section, or entity bytes not present cleanly in the local Git commit is incomplete. The parser reuses the shipped report selector and checklist ranges; it adds no Markdown implementation. Path-scoped Git cleanliness means unrelated sibling dirt does not block recovery.
2. **Post-dispatch transition guard.** Once the same-stage dispatch has set `worktree`, every `status --set` containing a status value different from the current stage shares that readiness result and refuses before mutation. Direction is irrelevant: successor hops, backward hops, terminal jumps, repeated/chained `status=` updates, and `--force` all refuse byte-clean. The same-stage dispatch mutation that sets the worktree remains allowed; unrelated non-status updates are outside this guard.

Structural validity is the scheduler's recovery boundary. The First Officer still owns the live completion judgment:

- **Live completion:** the FO has an active dispatch epoch, so it MUST observe that host adapter's `«completion-signal»` and verify the durable current-stage report against the retained dispatched checklist. A report without the runtime signal cannot complete a live worker.
- **Cold recovery:** after a restart, the status-owned structural and durability predicate alone controls whether ordinary successor projection resumes. This narrowed journey makes no additional semantic checklist-reconstruction promise and adds no cold runtime harness.

This is deliberately first-entry-only. Current workflow topology enters an unvisited working stage after a gate; feedback recovery dispatches its target directly and does not use `status --next`. A report from a prior visit therefore cannot satisfy this invariant in the supported path. A same-stage re-entry epoch is explicitly deferred until either a gate can target a previously visited stage or feedback recovery begins using `status --next`; that future capability needs an epoch/visit identity rather than pretending the latest historical report is fresh. Initial seed stages preserve legacy successor projection, gate and terminal stages retain suppression, and a set `worktree` retains existing in-flight suppression before this rule is considered.

### Spike result

The exact archived `runtime-live-e2e-codex-live` artifact was replayed. Its first mutation was `implementation -> validation`; no initial implementation dispatch preceded it, and the final entity contained exactly one implementation report (`cycle 2`). A throwaway real-CLI test over `writeRecordedGateFixture` then consumed validation into handoff and observed:

- before a handoff report, both `status --next --json` and `--boot --identify --json` returned exactly `current=handoff,next=done`;
- `status --read ... --stage handoff --checklist --json` alone failed with “no Stage Report”;
- after a conforming handoff report, the same status/boot projection was byte-identical while the report read succeeded.

The throwaway file was removed. This falsifies “consume stdout plus current prose is sufficient”: stdout is ephemeral, and a heading-only report read is too weak to credit completion. The status projection lacks a durable cold-boot readiness signal and must own the structural/durability half; the FO completion path remains responsible for host signal provenance and checklist semantics.

### Proposed change

Add one shared read-only predicate for “entered current stage awaiting first completion proof.” It reads the existing latest-stage-report/checklist primitives plus a literal path-scoped Git cleanliness check. `dispatchAnalysis` uses it to select the current stage instead of its successor; after same-stage worktree dispatch, `runSet` uses the same predicate before any force-bypass guards to reject every away-status update. The refusal must leave the entity byte-identical and emit no success stdout.

Amend `fo-dispatch-core.md` so `current==next` means dispatch the entered stage; neither the FO nor a helper may manufacture a completion signal or report. No host runtime adapter changes are required because each adapter already binds its own completion signal. Add the behavior assertion before changing this FO command text, per the repository's skill-smoke rule.

Do not change `gate consume`, gate schema, entity frontmatter, dispatch-build packages, runtime adapters, ensign protocol, feedback routing, or the strict rejection assertion.

Cheaper alternatives considered:

- More lifecycle prose or relying on `gate consume target-stage=...` cannot recover after a crash and was already present in the failing trace.
- Treating any matching heading as completion would duplicate the original bug with an empty, partial, dirty, failed, or wrong-cycle report. Report shape and path durability must be part of the one scheduler predicate.
- A durable dispatch ledger/lease could model every crash window, but the requested consume-before-spawn and report-present recovery needs no new state; existing worktree/roster behavior remains responsible for already-in-flight workers.
- Persisting runtime completion signals or checklist digests would widen worker/FO write scope and create a second scheduler state. The existing host signal plus the committed structural report are sufficient for this narrowed first-entry boundary.

## Acceptance criteria

**AC-1 (VALUE)** In a real two-cycle rejection journey, the durable ticket contains the original implementation report, first REJECTED validation, rework implementation report, and second PASSED validation in order; the existing strict rejection assertion remains unchanged and passes. **Test:** run only the real Codex `rejection-flow` scenario and inspect its archived entity/JSONL; removing current-stage projection or the mutation guard reproduces run `30197794474` and fails on one implementation report.

**AC-2 (DISPATCH)** After gate consume or equivalent first entry reaches a non-initial working stage, `status --next` and boot name that same stage until completion proof is durable. The same-stage worktree dispatch mutation is allowed; after it sets `worktree`, every away-status mutation—including successor, backward, terminal, repeated/chained update, and `--force` forms—is refused byte-clean until the report is durable. An observed worker spawn, the host completion signal, and an FO-validated committed report precede the live transition. **Test:** a real-CLI table performs the normal same-stage worktree dispatch, then snapshots bytes around each direction/force/chained mutation; the existing Claude and Codex `recorded-gate-lifecycle`/`rejection-flow` runners and Pi `TestLivePiRecordedGateLifecycle` order entered-stage dispatch/completion before the next transition.

**AC-3 (RECOVERY)** A cold boot after consume but before spawn exposes exactly one row with `current=<target>,next=<target>`. A heading-only, empty, structurally partial, failed, wrong-stage, later-malformed-over-older-valid, or uncommitted report does not change that row. A structurally complete, path-clean committed report restores ordinary successor projection. **Test:** extend the recorded-gate real-CLI fixture with committed malformed-report and dirty-valid-report mutants plus a committed-complete control; each exact JSON assertion fails if structural report validation is weakened.

**AC-4 (SCOPE)** The correction composes with consume, `status --next`, boot parity, initial/gate/terminal/worktree suppression, standing dispatch, fresh/reuse, and direct feedback routing without new scheduler state, report epochs, frontmatter, or fixture/runtime branches. Same-stage re-entry remains explicitly out of scope until topology or feedback scheduling can reach it. **Test:** table controls keep initial-stage successor projection, gate/terminal suppression, worktree-set suppression, terminal consume, and report-present direct feedback behavior byte-identical; full and race suites plus the README path→lane gate require both `claude-live` model legs, `codex-live`, and `pi-live` green.

## Test plan

Implementation starts with the status and live-scenario behavior assertions, before editing the FO references, then adds the shared predicate, projection, guard, and contract text. The existing `TestRejectionFlowNegativeSingleCycle` and `assertRejectionFlow` are preserved byte-for-byte; they already reject the exact one-implementation-report end state.

Focused offline commands:

```bash
go test ./internal/status -run 'TestEnteredStage|TestBootJSONDispatchableMirrorsNext' -count=1
go test ./internal/ensigncycle -run 'TestRecordedGateLifecycleRealCLIReplay|TestRecordedGateLifecycleTerminalConsumeHasNoDispatchableSuccessor|TestRejectionFlowNegativeSingleCycle' -count=1
go test ./...
go test ./... -race
```

Focused live preflight (existing runners/scenarios only; no new harness):

```bash
SPACEDOCK_LIVE_MODEL=sonnet go test -tags live -count=1 -timeout 40m -run 'TestLiveClaudeSharedScenarios/(recorded-gate-lifecycle|rejection-flow)$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_MODEL=claude-opus-4-8 go test -tags live -count=1 -timeout 40m -run 'TestLiveClaudeSharedScenarios/(recorded-gate-lifecycle|rejection-flow)$' ./internal/ensigncycle -v
go test -tags live -count=1 -timeout 40m -run 'TestLiveCodexSharedScenarios/(rejection-flow|recorded-gate-lifecycle)$' ./internal/ensigncycle -v
go test -tags live -count=1 -run 'TestSharedScenarioRunnerCoverage|TestPiSharedScenarioCoverage' ./internal/ensigncycle -v
go test -tags live -count=1 -run 'TestLivePiFrontDoorSmoke|TestLivePiRecordedGateLifecycle' ./internal/ensigncycle -v
```

The preflight is not the merge gate. Both planned reference edits are shared and host-neutral, so `docs/dev/README.md` requires every host lane that loads them. The implementation remains at validation until these exact registered jobs are genuinely green:

| Required lane | Existing registered surface (unchanged) | Cost and liveness bound | Stop/green rule |
|---|---|---|---|
| `offline` | `go test ./...` after build; local plan also runs `go test ./... -race` | No model spend; seconds to low minutes | Must pass before any environment-approved live job starts |
| `claude-live` / `sonnet` / `CI-E2E` | `TestLiveEnsignCycle\|TestLiveDefaultHeadlessStopsAtGate\|TestLiveZeroDiscoverReportsAndStops`; all `TestLiveClaudeSharedScenarios`; pty pair; `TestLiveMergedTeamModeDispatch` | Focus: 2 model journeys in parallel, about 9–15 minutes wall (the README records rejection-flow at 8.98m). Full job: 3 basic + 10 shared + 2 pty + 1 merged registrations. Each Go step has a 40m loose backstop; the shared runner kills 60s of stream silence. | Separate environment approval; recorded-gate stops after exactly one durable handoff dispatch/report, rejection-flow after cycle-2 PASSED validation. Shared/pty/merged steps use `if: !cancelled()`, so an earlier red does not skip their evidence; only cancellation stops them. Any red is rerun serial/isolated to green, never waived. |
| `claude-live` / `claude-opus-4-8` / `CI-E2E-OPUS` | Same four registered test steps as the sonnet leg | Same call count and bounds; independently approved and artifacted | This leg is independently required; a green sonnet leg cannot substitute for it. |
| `codex-live` / `CI-E2E-CODEX` | All 10 `TestLiveCodexSharedScenarios` through the current-checkout local marketplace | Focus: 2 serial `codex exec` calls (previous pair estimate 6–10 minutes). Full lane: 10 calls. Each call has a fixed 15m wall limit, no activity extension and no retry; suite step has a 40m outer backstop. | Missing auth fails after approval. Exit 0 plus durable assertions is required; timeout/red is evidence to diagnose and rerun isolated, not a skip. |
| `pi-live` / `CI-E2E-PI` | `TestSharedScenarioRunnerCoverage\|TestPiSharedScenarioCoverage`, then `TestLivePiFrontDoorSmoke\|TestLivePiRecordedGateLifecycle` | Coverage guard has no model spend. Two serial Pi launches, each with the existing fixed 5m process context (the registered Go step keeps its default 10m outer limit). | Front-door smoke stops after a committed implementation report; recorded-gate stops after exactly one durable handoff dispatch/report. Pi's `rejection-flow` entry remains the documented `gap`; do not add a harness or claim it ran. |

Local missing-auth skips are useful setup diagnostics but do not satisfy this matrix. A CI deployment left waiting/unapproved is likewise not green; validation stops until each approved lane actually runs and passes. The two shared scenarios are the focused change proof, while the full registered jobs are the path→lane merge gate.

Falsifiers are load-bearing:

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
| `internal/status/handlers.go` | +35/-0 | Refuse post-dispatch worktree away-status mutations before force-bypass guards |
| `internal/status/entered_stage_test.go` | +250/-0 | Projection, structure/durability matrix, post-dispatch all-direction guard, boot parity, and unchanged controls |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +85/-5 | Real consume/cold-boot/malformed/committed recovery and skill-trace ordering smoke |
| `skills/first-officer/references/first-officer-shared-core.md` | +0/-0 | Preserve the existing live completion contract |
| `skills/first-officer/references/fo-dispatch-core.md` | +18/-5 | Interpret current=current as dispatching the entered stage |
| `docs/site/concepts/gates-and-decisions.md` | +6/-1 | Document recovery rows, durable report rule, and unbypassable guard |

Expected total: **551 additions, 21 deletions (572 touched LOC)**, tolerance **±25% touched LOC (429–715)**. The FO dispatch reference is required while the shared completion core remains unchanged; `SKILL.md`, ensign instructions, and host runtime adapters are not. Any new production package/file beyond `entered_stage.go`, schema/frontmatter field, persisted signal/checklist/epoch, fixture-only branch, runtime adapter edit, or scheduler/lease state is a design deviation requiring gate reconfirmation.

### Documentation diff

In `docs/site/concepts/gates-and-decisions.md`, replace:

> Approval then uses `gate consume`, which rechecks eligibility and atomically writes the successor stage and consumed mark; the consumed descendant commit lands before ordinary successor dispatch.

with:

> Approval then uses `gate consume`, which rechecks eligibility and atomically writes the successor stage and consumed mark. Until that first-entered working stage has a durable, complete Stage Report, `status --next` and boot name it as both `current` and `next`. Once the same-stage dispatch sets its worktree, every away-status `status --set`—backward or forward, even with `--force`—is refused until the report is durable. The consumed descendant commit therefore lands before one recoverable successor dispatch.

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

## Stage Report: ideation (cycle 3)

- DONE: Require every registered applicable host lane for the two shared First Officer reference edits.
  The proof matrix now requires offline, both independently approved Claude model legs, Codex, and Pi, with exact registered test surfaces, call counts, liveness bounds, auth behavior, and stop conditions.
- DONE: Reuse existing runners and scenarios without hiding Pi's documented coverage gap.
  Claude and Codex reuse recorded-gate-lifecycle plus rejection-flow; Pi reuses its coverage guard, front-door smoke, and live recorded-gate lifecycle while rejection-flow remains explicitly `gap` and no harness is added.

### Summary

Cycle 3 corrects only the proof matrix. All accepted cycle-2 completion, recovery, guard, and first-entry-deferral semantics remain unchanged, and the full README path-to-lane gate—not a Codex-only focused run—now controls implementation completion.

## Stage Report: ideation (cycle 4)

- DONE: AC-1 — preserve the strict two-cycle report order and rejection oracle.
  The existing unchanged `assertRejectionFlow`/`TestRejectionFlowNegativeSingleCycle` oracle must observe original implementation, first REJECTED validation, rework implementation, and second PASSED validation in that order; omitting the rework report reproduces run `30197794474` and fails.
- DONE: AC-2 — self-project the entered current stage, guard every away mutation byte-clean, and require the live completion signal.
  Before durable completion proof, boot and `status --next` expose `current=<target>,next=<target>`; successor, backward, terminal, chained, and `--force` away mutations leave bytes unchanged, while live advancement additionally requires the worker spawn, host completion signal, and FO-verified committed report.
- DONE: AC-3 — recover cold from a committed semantically complete report and prove both boot sides.
  Before that report, cold boot keeps the target self-projected; after a path-clean committed report satisfies the reconstructed checklist and summary, cold boot exposes the ordinary successor without a duplicate target spawn, while heading-only, partial, failed, wrong-stage, stale, dirty, or checklist-incomplete reports remain vetoed.
- DONE: AC-4 — add no scheduler state and require every applicable registered host lane.
  The design adds no scheduler, ledger, lease, epoch, schema/frontmatter, host branch, or persisted completion signal; because both FO reference edits are shared, completion requires offline plus both Claude model legs, Codex, and Pi under the existing registered runners and documented Pi gap.

### Summary

Cycle 4 is an evidence-only addendum mapping each acceptance criterion to its decisive oracle. It changes no design, checklist, product file, or frontmatter and preserves the accepted first-entry-only scope.

### Feedback Cycles

- Cycle 1: REJECTED — Roborev panel job 2948; surface 22/715 vs estimate 8/572 (125%); AC unchanged
- Cycle 2: REJECTED — validation and captain scope reset; surface 22/715 vs estimate 8/572 (125%); AC narrowed: retain worktree-safe entered-stage actionability, revert merge-finalization changes, and remove the cold checklist-omission promise

## Stage Report: implementation

- DONE: Preserve the strict two-cycle rejection oracle and make the original entered implementation dispatch/report appear before first validation.
  Commits `6558ad8c`, `d6f2578b`, and `bd2e467b` leave the strict rejection assertion unchanged; the final-SHA isolated Codex `rejection-flow` passed in 379 seconds with the original implementation report, first REJECTED validation, rework implementation report, and second PASSED validation durably observed.
- DONE: Self-project an entered, incomplete working stage and reject every away-status mutation byte-clean until durable completion proof; preserve same-stage dispatch and cold recovery.
  `internal/status/entered_stage.go` supplies the shared exact-stage report and path-clean commit predicate used by boot/next projection and the transition guard; focused real-CLI tables prove malformed, failed, wrong-stage, later-masked, dirty, force, chained, backward, terminal, same-stage, worktree, and committed-recovery cases, and the final recorded-gate Codex journey passed in 194 seconds.
- DONE: Stay within the accepted first-entry/status-owned boundary, run focused local live proof before CI, and finish with classified Roborev review.
  The final diff is 22 files and 715 touched lines with no scheduler state, ledger, lease, epoch, schema/frontmatter, host branch, or persisted completion signal; `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` pass, no CI was triggered, and advisory implementation round 1 records Roborev job 2948 plus all supported-workflow, harm, boundary, and trigger classifications after fixing the material worktree-overlay and rejected-merge compatibility findings from the earlier panels.

### Summary

Implementation now keeps a gate-entered working stage as its own dispatch target until a durable exact-stage report is structurally complete and committed, and it refuses every direct away-status mutation byte-clean during that first-entry window. The shared First Officer contract distinguishes live host-signal verification from cold committed-report recovery, existing worktree and rejected-merge compatibility paths are preserved, full and race suites pass, and both required final-SHA Codex lifecycle journeys are green.

## Stage Report: validation

- DONE: Reproduce the final-SHA Codex rejection-flow and recorded-gate journeys and verify entered implementation is durably dispatched before validation.
  At source tip `bd2e467b`, isolated-OAuth Codex runs passed in 346s and 215s; the clean rejection entity records original implementation, first REJECTED validation, rework implementation, and second PASSED validation in order, while recorded-gate durably commits and completes handoff after consume.
- DONE: Adversarially verify self-projection, byte-clean away-status guards, cold committed-report recovery, and unchanged supported compatibility paths against AC-1 through AC-4.
  Focused real-CLI tests and an external final-SHA binary matrix passed malformed/failed/wrong-stage/later-masked/dirty/force/chained/backward/terminal/same-stage/rejected-merge/committed-recovery controls, but reproduced the material worktree and passed-merge defects below.
- DONE: Run applicable focused/full/race/format checks, audit Roborev dispositions with the release-scope template, and recommend PASSED or REJECTED without pushing or triggering CI.
  `gofmt -w ./cmd ./internal`, focused status/ensigncycle/merge tests, `go test ./...`, `go test ./... -race`, and `git diff --check` pass; no product write, push, PR, CI run, or schema/state expansion occurred.

### Material Findings

- Outcome defect — normal gate-entered worktree dispatch is supported, and `worktree` is set by the shipped FO recipe; afterward any away-status update bypasses the guard, so unfinished work can reach validation without a report. This violates AC-2's every-away-mutation boundary; the final-SHA external CLI changed `build -> review` with exit 0 and changed bytes.
- Outcome defect — documented direct `merge guard --verdict passed` use can clear an existing merge `mod-block` before the incomplete-report terminal transition fails, leaving a failed command's state partially mutated. This violates the non-negotiable merge state-integrity/byte-clean refusal boundary adjacent to AC-2; the final-SHA external CLI exited 1 after changing `mod-block: merge:local-merge` to empty.
- Evidence defect — AC-3 promises that a cold, structurally complete report missing a reconstructed checklist obligation is vetoed by the FO, but the recovery matrix proves only scheduler structure/durability and both live journeys carry completion signals. No command or durable cold-FO journey proves semantic veto/no-duplicate behavior for that supported crash trigger.

### Deferred Risks and Roborev Recheck

- Job 2948 non-versioned state is deferred: Git-less/untracked/unborn state cannot recover, but AC-3 promises committed path-clean evidence; promote when public workflow support includes non-versioned state.
- Job 2948 same-stage re-entry is deferred by AC-4; promote when a gate targets a visited stage or feedback recovery uses `status --next`.
- Job 2948 Git-process scaling is deferred: five reads over the supported 131-entity workflow measured 0.14-0.17s; promote if that workload exceeds a 1s read budget or a registered latency regression appears.
- Job 2948 merge-frontmatter strictness is deferred; promote when a sanctioned FO trace carries an uncommitted non-`pr`/`mod-block` delta and completion is refused. Gate/terminal exclusion coverage and the hardcoded CLI fixture are polish with no observed supported harm.
- Earlier worktree-overlay deadlock and rejected-merge finalization findings are exercised: rejected merge is fixed, while the overlay deadlock's skip-based correction leaves the broader material worktree bypass above.

### Summary

Recommendation: REJECTED. AC-1's strict ordered journey and AC-4's no-new-state/scope boundary are proven, and the structural half of AC-3 passes, but two supported mutation paths violate AC-2/state integrity and the semantic cold-recovery half of AC-3 lacks behavioral evidence. Correct the worktree-aware proof path and merge preflight atomically, then add a cold checklist-omission journey before revalidation.

## Stage Report: implementation (cycle 2)

- DONE: Keep exactly one actionable post-consume entered-stage projection, including safe worktree dispatch.
  Commit `fd657883d` lets the normal same-stage mutation set `worktree`, then refuses every away-status direction byte-clean until the current-stage report is committed; removing that worktree guard makes the focused real-CLI journey fail.
- DONE: Revert unrelated merge-finalization behavior and remove the cold checklist-omission promise and proof.
  Merge production/tests and the shared completion core are byte-identical to the feature base; the entity, FO dispatch reference, and site prose now promise only structural cold projection, with no checklist-reconstruction claim or new live harness.
- DONE: Run focused real-CLI, full Go, race, and format checks on the reduced branch.
  The entered-stage/merge status tests and recorded-gate/rejection-flow real-CLI tests pass; `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check` all pass at `fd657883d`.

### Summary

Cycle 2 narrows the candidate to the provider-free post-consume journey and leaves merge finalization outside this change. The reduced branch self-projects the entered stage, safely dispatches it into its worktree, and prevents premature transition until durable completion proof exists.

## Stage Report: validation (cycle 2)

- DONE: Reproduce post-consume self-projection, same-stage worktree dispatch, and byte-clean refusal before committed completion.
  At final SHA `fd657883d`, the real-CLI consume replay returned `current=handoff,next=handoff`; the normal same-stage worktree mutation succeeded, all successor/backward/terminal/force/chained away mutations exited 1 with empty stdout and identical bytes, and a committed complete report unlocked the successor.
- DONE: Verify merge production/tests are unchanged and no cold checklist-omission semantic or harness remains.
  Path-scoped comparison with feature base `50f8d1fb7` was byte-identical for merge production/tests and `first-officer-shared-core.md`; the reduced diff contains no former checklist-reconstruction production path or cold runtime harness.
- DONE: Run focused real-CLI, full Go, race, format, and diff checks; classify only journey-relevant findings.
  Focused status and recorded-gate/rejection tests passed, the critical worktree journey passed 10 repeated runs, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and both diff checks passed; no CI or provider-backed run was triggered.

### Summary

Recommendation: PASSED for the exact provider-free post-consume journey at `fd657883d`. No journey-relevant material, evidence, or deferred finding remains; merge-finalization and cold checklist-omission semantics were confirmed absent and were not re-expanded.

## Stage Report: validation (cycle 3)

- DONE: Verify exact-tip patch semantics against the previously passed candidate.
  At rebased SHA `cb01129b` on release-stamped base `5e7f1ffa`, the entered-stage predicate, projection, guard, and focused test source are byte-identical to passed tip `fd657883d`; the release-base patch remains 13 files and 641 touched lines.
- DONE: Inspect conflict resolutions for drift without expanding scope.
  The gate-site and FO-dispatch resolutions preserve upstream retained-request/rejection-routing prose while adding only current=current dispatch semantics; the recorded-gate resolution preserves the upstream oracle and adds the provider-free pre/post report projection assertions, with merge and shared completion surfaces unchanged from `5e7f1ffa`.
- DONE: Re-run exact-tip focused, full, race, and diff checks.
  Real-CLI consume/recovery, worktree byte-clean refusal, projection/boot parity, and negative rejection tests passed; SHA-guarded `go test ./...`, `go test ./... -race`, and `git diff --check` passed at `cb01129b`; no CI was triggered.

### Summary

Recommendation: PASSED for exact-tip integration at `cb01129b`. No conflict-resolution drift or critical-journey material/evidence finding was observed.
