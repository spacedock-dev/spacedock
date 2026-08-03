---
title: Remove the standalone gate eligibility ceremony
status: ideation
source: "Captain simplification review, 2026-08-03: eligibility is a necessary atomic safety predicate but an unnecessary operator preflight and readiness vocabulary."
score: 0.9
sprint: durable-decisions
group: gate-operator-ux
milestone: 0.27.0
sprint-readiness: ready
id: bv3hhbqr5spt1wn4557qyp8c
gates:
    version: 1
    records:
        - id: gate:bv3hhbqr5spt1wn4557qyp8c:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:bv3hhbqr5spt1wn4557qyp8c-backlog-1
              briefing:
                id: briefing:bv3hhbqr5spt1wn4557qyp8c:backlog:attempt-1:revision-1
                digest: sha256:4b936d38cb7ea50d8823bb73c50ef257b510d64736bf786950a32ff37c97ed58
                request-digest: sha256:098a9a069ace77b1c9afdc1e0166bb33788841d7d2ea8f37e8b14f43f41def03
                room-ref: ./remove-standalone-gate-eligibility/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:bv3hhbqr5spt1wn4557qyp8c:backlog:1
                briefing: briefing:bv3hhbqr5spt1wn4557qyp8c:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-03T23:30:19.040582Z"
                decision: approve
                reason: Captain directed dispatch of the next ideation; latest recorded SO concurrence places bv after g3 and preserves the private fail-closed predicate while removing only the public eligibility ceremony.
              application:
                target-stage: ideation
                state: consumed
---

## Outcome

A First Officer uses `status` to see the next action and uses `gate consume` to apply an approval. The operator does not run a separate eligibility ceremony.

`gate consume` and terminal delivery must retain the internal fail-closed predicate. This predicate prevents stale approval, wrong-stage advancement, wrong-target advancement, and duplicate spending.

## Problem

The public `gate eligibility` command exposes a read-only preview of the predicate that `gate consume` recomputes under the entity lock. The normal First Officer path does not require this preview. The skill calls it an optional diagnostic.

The extra command creates another operator step and another set of condition labels. These labels can look like scheduler state even though they are derived diagnostics.

The safety mechanism and the command are different things. Removing the command must not remove or weaken the mutation guard.

## Proposed direction

Remove `gate eligibility` from the public CLI, help, command reference, completion, and First Officer procedure. Let `status` project the next action from durable facts.

Keep one internal predicate at each authority-spending boundary. `gate consume` must evaluate it while it holds the entity lock. Terminal delivery and rework must use the same authority checks.

On refusal, the acting command must report the concrete reason. Do not expose the predicate as a separate workflow phase or stored state class.

This is an unreleased-v1 cut. Do not add compatibility aliases, deprecation machinery, a replacement preflight command, or a new readiness class.

## Acceptance criteria

**AC-1 (VALUE) — The normal gate journey has no eligibility preflight.**
Verified by: a real CLI lifecycle runs `status`, close-and-consume, successor dispatch or terminal routing, and no `gate eligibility` command.

**AC-2 — Authority spending remains fail-closed and atomic.**
Verified by: stale input, wrong stage, wrong target, consumed authority, superseded authority, and malformed binding controls refuse without an unauthorized status transition.

**AC-3 — The operator receives an actionable refusal from the acting command.**
Verified by: `gate consume` and terminal delivery controls report the applicable refusal reason without requiring a prior diagnostic command.

**AC-4 — The public surface contains no standalone eligibility ceremony.**
Verified by: CLI help, command reference, completion, and First Officer instructions omit `gate eligibility`; source retains the private predicate only at mutation boundaries.

## Expected surface

Expected changes are limited to gate CLI routing and help, First Officer gate instructions, public documentation, completion, and focused lifecycle tests. The internal predicate and its table tests remain.

## Test plan

Start with the real operator lifecycle and negative controls. Remove the public command only after these tests prove that `consume` owns the same checks.

Run focused gate, CLI, status, terminal-delivery, and recorded-lifecycle tests. Then run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.

## Stage-specific test gates

- Ideation must trace every production caller of the public command and the private predicate. It must identify any consumer that cannot use `status` or the acting command.
- Implementation must remove no safety branch from the private predicate.
- Validation must use a detached mutation that weakens one consume guard and make sure that a behavioral test fails.

