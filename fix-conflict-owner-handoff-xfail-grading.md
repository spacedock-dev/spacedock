---
title: Fix conflict-owner-handoff XFAIL grading for missing marker
status: ideation
source: "FO write-scope review, 2026-08-13: direct test edit routed through filing"
score: 0.85
sprint: live-evidence-followups
sprint-readiness: ready
group: common-evidence
id: rzrx7a00yxk9kkvy4mxcj8ep
gates:
    version: 1
    records:
        - id: gate:rzrx7a00yxk9kkvy4mxcj8ep:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:rzrx7a00yxk9kkvy4mxcj8ep-backlog-1
              briefing:
                id: briefing:rzrx7a00yxk9kkvy4mxcj8ep:backlog:attempt-1:revision-1
                digest: sha256:d1824570ffe158b5c49568a68c6ee86694d16436a73112df5f4481c8fdb724d5
                request-digest: sha256:39c1da392d035c198eb09647ab13e7ef3db540e6fabae2b802e672763979bf15
                room-ref: ./fix-conflict-owner-handoff-xfail-grading/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:rzrx7a00yxk9kkvy4mxcj8ep:backlog:1
                briefing: briefing:rzrx7a00yxk9kkvy4mxcj8ep:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-13T22:21:28.847522Z"
                decision: approve
                reason: Captain approved backlog gate; advance to ideation for the XFAIL grading fix.
              application:
                target-stage: ideation
                state: consumed
---

## Problem

`assertConflictOwnerHandoff` reads the worker-produced `owner-handoff.marker` via
`readFile(t, fixture.marker)`, which calls `t.Fatal(err)` when the file is missing.
When the FO produces no marker (the handoff did not happen), the test hard-fails
with `--- FAIL` and bypasses `gradeLive`. The `liveXFail("claude-opus", ...)` and
`liveXFail("pi", ...)` bindings on `TestLiveCommonOwnedConflictOwnerHandoff` can
therefore never grade the "no marker" failure mode as XFAIL, even though a missing
marker is the expected semantic failure those bindings anticipate.

This is a grading-correctness gap, not a product repair. The assert turns a typed
semantic failure (no handoff) into an untyped hard failure, which is the opposite
of the `ts` XFAIL policy: typed semantic failures must grade as XFAIL for a bound
target.

## Visible value

A missing worker marker grades as XFAIL (green for a bound target) instead of a
hard lane `--- FAIL`. Both runtimes bound to this journey (claude-opus and pi)
get the intended XFAIL grading for the no-handoff mode. The lane stops failing on
a model or transport issue that produces no marker, and reports the expected
semantic failure instead. Measured against baseline: before this fix, a
no-marker run on a bound target is a red lane `--- FAIL`; after, the same run is
XFAIL-green with the observed code recorded.

## Out of scope

- Removing either XFAIL binding. The bindings stay until their runtime-specific
  repair entities produce XPASS and then a normal PASS (`repair-pi-owner-handoff`
  fe7, `repair-opus-owner-handoff` bqy own binding removal).
- Product behavior of the owner-handoff journey. This task changes only the
  shared assert's missing-file failure mode.
- Sonnet, Codex, or any other journey's assert. Only `assertConflictOwnerHandoff`.
- A new gate command, stored format, authority source, or CI lane.

## Acceptance criteria

**AC-1 (VALUE) — A missing marker grades as XFAIL on a bound target, not a lane FAIL.**

Verified by: a no-marker run of `TestLiveCommonOwnedConflictOwnerHandoff` on a
target with an active XFAIL binding reports XFAIL-green with an observed code,
not a `--- FAIL`. The control is a fixture or run that produces no marker file;
the baseline is the current `readFile` `t.Fatal` behavior.

**AC-2 — The missing-marker failure returns an error, not `t.Fatal`.**

Verified by: `assertConflictOwnerHandoff` reads `fixture.marker` via
`readFileAllowMissing` and returns `fmt.Errorf("worker marker missing, want
runtime-worker-owner")` when the file is absent, so the failure routes through
`gradeLive` instead of aborting the test.

**AC-3 — A present-but-wrong marker is unchanged.**

Verified by: a marker whose content is not `runtime-worker-owner` still returns
the existing `worker marker = %q, want runtime-worker-owner` error; the
`fixture.entity`, `git rev-parse`, `git log`, `git show`, `git status`,
`git worktree list`, and `git branch` reads are unchanged.

**AC-4 — Runtime-neutral; preserves bindings and the assert for repair entities.**

Verified by: no change to the `liveXFail("claude-opus", ...)` or
`liveXFail("pi", ...)` bindings, and `repair-pi-owner-handoff` (fe7) /
`repair-opus-owner-handoff` (bqy) AC-3 ("preserve the owner-handoff assertion")
still holds.

**AC-5 — Offline checks pass.**

Verified by: `gofmt`, `go vet -tags live ./internal/ensigncycle`,
`go build -tags live ./internal/ensigncycle`, and the offline
`ConflictOwner|GradeLive|LiveGrade|KeepMovingDurable` unit tests pass.

## Test plan

Use the offline `ConflictOwner|GradeLive|LiveGrade|KeepMovingDurable` unit tests
first. The `readFileAllowMissing` helper already exists in
`shallow_boot_assert_test.go`. Keep `readFile(t, ...)` for `fixture.entity` (a
fixture file the test itself wrote, so missing there is a real setup failure).

## Notes

- General live-test grading fix, not a Pi repair. Do not bundle with the
  runtime-specific repair entities.
- Filed by the FO after a write-scope review caught direct test edits; the
  worker re-derives the edit from this spec.
