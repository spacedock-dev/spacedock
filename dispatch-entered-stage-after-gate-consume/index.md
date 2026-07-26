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

Ship one host-neutral invariant in the existing status scheduler:

> A declared non-initial, non-gated, non-terminal current stage with no matching `## Stage Report: <current>` is the dispatch target itself; it cannot transition to a later stage. Once its report exists and the FO validates it, ordinary successor projection resumes.

`status --next` and boot continue to use the same `dispatchAnalysis` source. In the entered/unreported case their existing five-field row becomes `current=<stage>,next=<same-stage>` and uses that stage's existing concurrency/worktree settings. The normal dispatch recipe can therefore set the same status idempotently, build that stage, and spawn it without a new field or host branch. `status --set status=<later-stage>` shares the predicate and refuses byte-clean—even with `--force`—until the matching report exists; same-stage dispatch mutations remain allowed. Initial seed stages preserve the legacy successor behavior, gate and terminal stages retain their existing suppression, and a set `worktree` retains the existing in-flight suppression before this rule is considered.

This is deliberately a first-entry invariant. Feedback re-entry continues through the existing feedback route, which already dispatches rework directly and owns cycle identity; inventing report epochs, a dispatch ledger, or a lease here would create the forbidden second scheduler. Worker provenance remains a runtime obligation: a report unlocks the status projection only after the existing FO completion path verifies it against the dispatched checklist and completion signal. The strict rejection-flow assertion stays unchanged as the durable end-to-end oracle.

### Spike result

The exact archived `runtime-live-e2e-codex-live` artifact was replayed. Its first mutation was `implementation -> validation`; no initial implementation dispatch preceded it, and the final entity contained exactly one implementation report (`cycle 2`). A throwaway real-CLI test over `writeRecordedGateFixture` then consumed validation into handoff and observed:

- before a handoff report, both `status --next --json` and `--boot --identify --json` returned exactly `current=handoff,next=done`;
- `status --read ... --stage handoff --checklist --json` alone failed with “no Stage Report”;
- after a conforming handoff report, the same status/boot projection was byte-identical while the report read succeeded.

The throwaway file was removed. This falsifies “consume stdout plus current prose is sufficient”: stdout is ephemeral, and using the separate report read as an FO pre-pass would make a second scheduling decision. The status projection lacks the cold-boot readiness signal and must own it.

### Proposed change

Add one shared read-only predicate for “entered current stage awaiting first report.” `dispatchAnalysis` uses it to select the current stage instead of its successor; `runSet` uses the same predicate to refuse a status change away from that stage. Reuse the shipped latest-stage-report selector rather than adding another Markdown parser. Do not change `gate consume`, gate schema, entity frontmatter, dispatch packages, runtime adapters, FO skills, feedback routing, or the rejection assertion.

Cheaper alternatives considered:

- More lifecycle prose or relying on `gate consume target-stage=...` cannot recover after a crash and was already present in the failing trace.
- A prose pre-pass combining `status --next` with `status --read --checklist` duplicates scheduler choice across two outputs that the spike proved disagree.
- A durable dispatch ledger/lease could model every crash window, but the requested consume-before-spawn and report-present recovery needs no new state; existing worktree/roster behavior remains responsible for already-in-flight workers.

## Acceptance criteria

**AC-1 (VALUE)** In a real two-cycle rejection journey, the durable ticket contains the original implementation report, first REJECTED validation, rework implementation report, and second PASSED validation in order; the existing strict rejection assertion remains unchanged and passes. **Test:** run only the real Codex `rejection-flow` scenario and inspect its archived entity/JSONL; removing current-stage projection or the mutation guard reproduces run `30197794474` and fails on one implementation report.

**AC-2 (DISPATCH)** After a gate consume or equivalent first entry reaches a non-initial working stage, `status --next` names that same stage until its matching report exists, and any attempted mutation to a later stage is refused byte-clean. An observed worker spawn plus the FO-validated Stage Report precede the eventual transition; build output, `started`, narration, and an FO-authored report do not satisfy the live oracle. **Test:** a real-CLI table attempts successor mutation before report, after `started`, and after `dispatch build` and expects the same refusal/bytes; the focused live JSONL must show implementation dispatch evidence before the first validation transition.

**AC-3 (RECOVERY)** A cold boot after consume but before spawn exposes exactly one row with `current=<target>,next=<target>`; after the durable report, boot exposes the ordinary successor and never the target again. **Test:** extend the recorded-gate real-CLI fixture with consume-then-boot and report-then-boot controls; changing either projection makes the corresponding exact JSON assertion fail.

**AC-4 (SCOPE)** The correction composes with consume, `status --next`, boot parity, initial/gate/terminal/worktree suppression, standing dispatch, fresh/reuse, and feedback routing without new scheduler state, frontmatter, or fixture/runtime branches. **Test:** table controls keep initial-stage successor projection, gate/terminal suppression, worktree-set suppression, terminal consume, and report-present feedback behavior unchanged; full and race suites plus an adversarial mutation of the shared predicate must turn the focused tests red.

## Test plan

Implementation starts with the status tests, then the shared predicate, projection, and guard. The existing `TestRejectionFlowNegativeSingleCycle` and `assertRejectionFlow` are preserved byte-for-byte; they already reject the exact one-implementation-report end state.

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

The offline focus is expected to cost seconds and no model calls; the focused live pair is expected to cost 6–10 minutes and two Codex scenario calls. Because this touches native status projection and mutation guarding, validation also requires the workflow's detached adversarial audit: make the predicate ignore report absence or let `--force` bypass the guard and confirm the new tests fail.

### Expected surface

Exact baseline, measured as added/deleted LOC:

| File | Expected LOC | Purpose |
|---|---:|---|
| `internal/status/entered_stage.go` | +45/-0 | One shared current-stage/report predicate; no persisted state |
| `internal/status/format.go` | +18/-7 | Select current vs successor inside existing `dispatchAnalysis` |
| `internal/status/handlers.go` | +18/-0 | Refuse later-stage `--set` before the current report |
| `internal/status/entered_stage_test.go` | +155/-0 | Projection, guard, boot parity, and unchanged-control table |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +32/-0 | Real consume/cold-boot/report/no-duplicate replay |
| `docs/site/concepts/gates-and-decisions.md` | +3/-1 | Document the visible current=current recovery row and guard |

Expected total: **271 additions, 8 deletions (279 touched LOC)**, tolerance **±25% touched LOC (209–349)**. Any new production file, schema/frontmatter field, FO/runtime skill edit, fixture-only branch, or scheduler/lease state is a design deviation requiring gate reconfirmation.

### Documentation diff

In `docs/site/concepts/gates-and-decisions.md`, replace:

> Approval then uses `gate consume`, which rechecks eligibility and atomically writes the successor stage and consumed mark; the consumed descendant commit lands before ordinary successor dispatch.

with:

> Approval then uses `gate consume`, which rechecks eligibility and atomically writes the successor stage and consumed mark. Until that entered working stage has a matching Stage Report, `status --next` and boot name it as both `current` and `next`, and a later-stage `status --set` is refused; after the report is verified, ordinary successor projection resumes. The consumed descendant commit therefore lands before exactly one recoverable successor dispatch.

## Stage Report: ideation

- DONE: Replay the exact skipped-first-implementation trace and test the cheapest existing lifecycle/status correction first.
  Run 30197794474 was replayed; a removed real-CLI spike proved boot/next stay on the successor while only a separate report read distinguishes pre/post report.
- DONE: Define one host-neutral dispatch-before-next-transition invariant with cold-boot recovery and no duplicate worker.
  The shared current-stage/report predicate drives both existing dispatch projection and a byte-clean transition guard; no new state, scheduler, host branch, or fixture exception is introduced.
- DONE: Declare exact files and LOC plus offline and focused-live falsifiers without a second scheduler or fixture special case.
  Six files, 279 touched LOC ±25%, exact commands, mutation falsifiers, and the concrete site-doc replacement are recorded above.

### Summary

Ideation chooses the smallest status-owned correction because the archived trace already contained the lifecycle prose and the spike proved the durable projection was ambiguous. The design preserves the strict two-cycle assertion, adds no dispatch state, and makes the consumed/entered stage recoverable through the one scheduler the FO already follows.
