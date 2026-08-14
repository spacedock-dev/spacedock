---
title: Fix conflict-owner-handoff XFAIL grading for missing marker
status: validation
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
        - id: gate:rzrx7a00yxk9kkvy4mxcj8ep:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:rzrx7a00yxk9kkvy4mxcj8ep-ideation-1
              briefing:
                id: briefing:rzrx7a00yxk9kkvy4mxcj8ep:ideation:attempt-1:revision-1
                digest: sha256:084c62d6330371ecfaea2b39f47d880218695a675132524dd65133f1b3aae892
                request-digest: sha256:fa188aa9ab7a0b1ae3bcf22039291ee1f71fc33e55b716ab208084f6eaf83361
                room-ref: ./fix-conflict-owner-handoff-xfail-grading/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:rzrx7a00yxk9kkvy4mxcj8ep:ideation:1
                briefing: briefing:rzrx7a00yxk9kkvy4mxcj8ep:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-13T22:54:23.992116Z"
                decision: approve
                reason: Captain approved ideation gate; dispatch implementation for the one-line XFAIL grading fix.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:rzrx7a00yxk9kkvy4mxcj8ep:validation
          stage: validation
          attempts:
            - id: gate-attempt:rzrx7a00yxk9kkvy4mxcj8ep-validation-1
              briefing:
                id: briefing:rzrx7a00yxk9kkvy4mxcj8ep:validation:attempt-1:revision-1
                digest: sha256:9e606e91a2427f4cf721b38c6dbc3c1fd0611586215759feda1356d3e4844982
                request-digest: sha256:8b9ad4c9e9ded0402b0de34860ac28589f87ac9e57fae5fe538d0bb91e85c977
                room-ref: ./fix-conflict-owner-handoff-xfail-grading/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:rzrx7a00yxk9kkvy4mxcj8ep:validation:1
                briefing: briefing:rzrx7a00yxk9kkvy4mxcj8ep:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-14T02:50:39.863599Z"
                decision: approve
                reason: Captain approved validation; terminalize the XFAIL grading fix.
              application:
                target-stage: done
                state: pending
started: 2026-08-13T22:55:16Z
worktree: .worktrees/spacedock-ensign-fix-conflict-owner-handoff-xfail-grading
mod-block: merge:pr-merge
pr: pr-merge:684
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

## Stage Report: implementation

- DONE: swapped `readFile(t, fixture.marker)` to `readFileAllowMissing(fixture.marker)` in `assertConflictOwnerHandoff` (`internal/ensigncycle/conflict_owner_handoff_live_test.go`), adding a `got == ""` branch returning `fmt.Errorf("worker marker missing, want runtime-worker-owner")`; the present-but-wrong `worker marker = %q` path is byte-unchanged and `readFile(t, fixture.entity)` is untouched.
  `internal/ensigncycle/conflict_owner_handoff_live_test.go:77-78` — `if got := strings.TrimSpace(readFileAllowMissing(fixture.marker)); got == "" { return fmt.Errorf("worker marker missing, want runtime-worker-owner") }`; commit 6bebae7c8 (`1 file changed, 3 insertions(+), 1 deletion(-)`) touches only this assert.
- DONE: ran `gofmt -w internal/ensigncycle`, `go vet -tags live ./internal/ensigncycle`, `go build -tags live ./internal/ensigncycle`, and `go test -run 'ConflictOwner|GradeLive|LiveGrade|KeepMovingDurable' ./internal/ensigncycle`; all pass.
  `gofmt -l` empty; `go vet -tags live` silent; `go build -tags live` silent; `go test -run …` → `ok\tgithub.com/spacedock-dev/spacedock/internal/ensigncycle\t(cached)`.
- DONE: did not touch any `liveXFail(...)` binding, any other assert, or any CI lane; fe7/bqy AC-3 preserved.
  `git diff HEAD~1 HEAD -- internal/ensigncycle/ | grep liveXFail` empty; only `conflict_owner_handoff_live_test.go` changed (1 file, 3+/1-).

### Summary

One-line change to the shared assert: missing `owner-handoff.marker` now returns an error (`worker marker missing, want runtime-worker-owner`) via `readFileAllowMissing` instead of `t.Fatal`, so `gradeLive` can grade the no-handoff mode as XFAIL on bound targets. Offline checks (`gofmt`, `go vet -tags live`, `go build -tags live`, and the `ConflictOwner|GradeLive|LiveGrade|KeepMovingDurable` tests) all pass.

## Stage Report: validation

- DONE: Independently verified the deliverable against the entity spec; did NOT edit the deliverable.
  `git show HEAD~1:.../conflict_owner_handoff_live_test.go` vs current: `readFile(t, fixture.entity)` untouched at line 74; `readFileAllowMissing(fixture.marker)` at line 77; missing-marker branch `return fmt.Errorf("worker marker missing, want runtime-worker-owner")` at line 78; present-but-wrong `worker marker = %q, want runtime-worker-owner` byte-identical (base line 78 → current line 80, same bytes).
- DONE: Confirmed no `liveXFail(...)` binding, no other assert, and no CI lane touched; fe7/bqy AC-3 preserved.
  `git diff HEAD~1 HEAD --name-only` → only `internal/ensigncycle/conflict_owner_handoff_live_test.go`; `git diff HEAD~1 HEAD | grep liveXFail` → empty (exit 1); no `.github`/`ci`/`workflow`/`lane` files changed.
- DONE: Re-ran the offline checks in the worktree and report actual output.
  `gofmt -l internal/ensigncycle` → empty (exit 0); `go vet -tags live ./internal/ensigncycle` → silent (exit 0); `go build -tags live ./internal/ensigncycle` → silent (exit 0); `go test -count=1 -run 'ConflictOwner|GradeLive|LiveGrade|KeepMovingDurable' ./internal/ensigncycle` → `ok  github.com/spacedock-dev/spacedock/internal/ensigncycle  0.424s` (exit 0).
- DONE: Appended this `## Stage Report: validation` with per-item DONE evidence and a Summary recommending PASSED.

### Summary

Recommending PASSED. The worktree deliverable matches the entity spec exactly: the shared assert reads the marker via `readFileAllowMissing`, returns `worker marker missing, want runtime-worker-owner` on absent file, keeps the present-but-wrong `worker marker = %q` path byte-unchanged, and leaves `readFile(t, fixture.entity)` and all `liveXFail` bindings untouched. Fresh (non-cached) offline checks — gofmt, go vet -tags live, go build -tags live, and the `ConflictOwner|GradeLive|LiveGrade|KeepMovingDurable` tests — all pass with exit 0.
