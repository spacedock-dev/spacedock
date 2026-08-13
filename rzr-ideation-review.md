# IDEATION GATE — Conflict-owner-handoff XFAIL grading (`rzr`)

Recommendation: **APPROVE and dispatch implementation.**

## Proposed approach

One-line, runtime-neutral change to `assertConflictOwnerHandoff` in
`internal/ensigncycle/conflict_owner_handoff_live_test.go`:

- Replace `readFile(t, fixture.marker)` with `readFileAllowMissing(fixture.marker)`
  (helper already in-package at `shallow_boot_assert_test.go:126`, returns `""`
  on missing instead of `t.Fatal`).
- Add a `got == ""` branch returning
  `fmt.Errorf("worker marker missing, want runtime-worker-owner")` so a missing
  marker routes through `gradeLive` and grades XFAIL for a bound target.
- Keep the present-but-wrong path (`worker marker = %q, want runtime-worker-owner`)
  byte-unchanged. Keep `readFile(t, fixture.entity)` (a fixture file the test
  wrote; missing there is a real setup failure, not a semantic failure).

## Why this is sound

- `readFile` calling `t.Fatal` is the root cause: it aborts the test before
  `gradeLive` can classify the failure. `readFileAllowMissing` returns `""`, the
  assert returns an error, and `liveJourney`'s exercise wrapper feeds it to
  `gradeLive(xfail=true, err)` → status `xfail` (green for a bound target).
- Runtime-neutral: both `liveXFail("claude-opus",…)` and `liveXFail("pi",…)`
  bindings benefit. No binding removed (fe7/bqy AC-3 preserved).
- Smallest sufficient mechanism: one helper swap + one error branch. No new
  abstraction, no oracle change, no retry.

## Proof plan

- Offline: `gofmt`, `go vet -tags live ./internal/ensigncycle`,
  `go build -tags live`, and `go test -run 'ConflictOwner|GradeLive|LiveGrade|KeepMovingDurable' ./internal/ensigncycle`.
- The `live_grade_unit_test.go` table already covers `xfail` with one/many
  semantic codes and empty-set `xpass`; the assert change lets the missing-marker
  mode reach that path.
- No live run required for this grading fix (it's offline-testable); a later
  live run on a bound target confirms XFAIL-green, but is not a dispatch gate.

## Risks

- None material. The present-but-wrong path is unchanged; the `fixture.entity`
  read stays fatal-on-missing (correct — setup failure).

## Decision ask

Approve to dispatch implementation (worktree) for this one-line grading fix, or
revise/hold with a concrete boundary.
