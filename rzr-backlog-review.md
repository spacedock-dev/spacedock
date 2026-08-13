# BACKLOG GATE — Conflict-owner-handoff XFAIL grading (`rzr`)

Recommendation: **APPROVE and dispatch ideation.**

## Capability and value

The `owned-conflict-owner-handoff` journey is XFAIL-bound for claude-opus and
pi. Its shared assert reads the worker-produced marker via `readFile(t, ...)`,
which calls `t.Fatal` when the file is missing. When the FO produces no marker
(the handoff did not happen), the test hard-fails with `--- FAIL` and bypasses
`gradeLive` — so the XFAIL binding can never grade the no-handoff mode as XFAIL.
This task routes the missing-marker failure through `gradeLive` by switching to
`readFileAllowMissing`, so a missing marker returns an error (XFAIL for a bound
target) instead of aborting the lane.

## Binding boundaries

- General live-test grading fix, runtime-neutral. Both claude-opus and pi
  bindings benefit; no binding is removed.
- Only `assertConflictOwnerHandoff`'s marker read changes. `fixture.entity`,
  all `git` reads, and the present-but-wrong marker path are unchanged.
- Does not touch any XFAIL binding, any other journey's assert, or any CI lane.
- Preserves `repair-pi-owner-handoff` (fe7) and `repair-opus-owner-handoff` (bqy)
  AC-3 ("preserve the owner-handoff assertion").

## Proof direction

Ideation confirms the grading path: `readFileAllowMissing` returns `""` for a
missing marker, the assert returns `fmt.Errorf("worker marker missing, want
runtime-worker-owner")`, and `gradeLive` grades it XFAIL for a bound target.
Implementation proves the offline unit tests (`ConflictOwner|GradeLive|LiveGrade`)
pass and the present-but-wrong path is byte-unchanged. The edit is small and
already verified offline by the FO before filing.

## Decision ask

Approve this grading-correctness fix for ideation, or revise/hold with a
concrete boundary.
