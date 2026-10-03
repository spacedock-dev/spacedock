---
title: Reconcile the live Pi XFAIL bindings with the owners and evidence that exist today
status: backlog
score: 0.7
source: "pi-ux carve review, 2026-10-03: live Pi bindings name an archived owner, one binding cannot XPASS in code, and one real Pi FAIL carries no binding."
id: d525n1p5zgnz99hmtjq16z57
sprint: pi-ux
group: tooling
sprint-readiness: ready
---

The live Pi XFAIL registry does not state the truth. Some bindings name owners that no
longer exist, one journey still grades a Pi stream with a Claude-dialect extractor, and
one real Pi conduct gap carries no binding at all.

## Problem

Three concrete facts, verified in the code on 2026-10-03:

1. The `rejection-flow` Pi binding names owner `p17swb3375rt525fn7f8xt7e`, which resolves
   to `finish-pi-rejection-flow` and is **archived** with verdict PASSED.
2. The `smallest-sufficient-mechanism` Pi journey still grades through
   `claudeMechanismTrace` (`pi_shared_live_runner_test.go:47`), a Claude stream-json
   extractor, so its binding cannot XPASS as written.
3. `recorded-gate-lifecycle` carries no Pi binding, and its owner describes an ordinary
   Pi FAIL with a real conduct gap.

Bindings that cannot fail, and owners that no longer exist, make every Pi pass rate and
every XFAIL claim unreliable.

## Scope boundary

This task reconciles bindings, owner status, and evidence only. It must not absorb the
`pi-live-completeness` journey repairs, must not weaken their assertions, and must not
delete a binding without exact passing evidence.

## Expected surface and tolerance

Estimate net LOC change: +25, across 2 files (the shared live-runner registry and the Pi
runner). Insertions ~+35, deletions ~-10. Tolerance: +/-15 net LOC, +/-1 file.

Declared semantic changes: which journeys carry a Pi binding, and which owner each names.
This task must NOT change the journeys themselves or the Codex and Claude bindings.

## Acceptance criteria

**AC-1 - Every live Pi binding names an active owner or is cleared by evidence.**
Verified by: a Go test that reads the registry and asserts each Pi binding's owner id
resolves to a non-archived entity, or the binding is absent. Independent baseline that can
move the wrong way: today the `rejection-flow` binding names an archived owner.
Falsifying edit: restore the archived owner id; the test must turn RED.

**AC-2 - No Pi journey is graded through a Claude-dialect extractor.**
Verified by: a Go test asserting each Pi journey's trace extractor is Pi-specific, with
`smallest-sufficient-mechanism` as the named case. Falsifying edit: point the Pi driver back
at `claudeMechanismTrace`; the test must turn RED.

**AC-3 - Every cleared or kept binding is recorded with its reason.**
Verified by: the task body states, per Pi binding, whether it was cleared, kept, or
re-anchored, and the evidence that decided it. Falsifying edit: remove one disposition line;
the record is incomplete.

**AC-4 (no-regression) -** `go test ./internal/ensigncycle/...` is green, and no
`pi-live-completeness` assertion is weakened.

## Test plan

Primary proof owner: the shared live-runner test beside `shared_live_runner_test.go`.
Deterministic Go tests. One live pi cadence run confirms the resulting binding states.
