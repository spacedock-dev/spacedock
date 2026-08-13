---
title: Fix conflict-owner-handoff XFAIL grading for missing marker
status: backlog
sprint: test-behavior-completeness
group: common-evidence
source: "fo write-scope review of direct test edits"
id: rzrx7a00yxk9kkvy4mxcj8ep
---

## Problem

`assertConflictOwnerHandoff` reads the worker-produced `owner-handoff.marker` via `readFile(t, fixture.marker)`, which calls `t.Fatal(err)` when the file is missing. When the FO produces no marker (the handoff did not happen), the test hard-fails with `--- FAIL` and bypasses `gradeLive`. The `liveXFail("claude-opus", ...)` and `liveXFail("pi", ...)` bindings on `TestLiveCommonOwnedConflictOwnerHandoff` can therefore never grade the "no marker" failure mode as XFAIL, even though that mode is the expected semantic failure the bindings anticipate.

## Value

A missing worker marker grades as XFAIL (green for a bound target) instead of a hard lane `--- FAIL`. Both runtimes bound to this journey (claude-opus and pi) get the intended XFAIL grading for the no-handoff mode. The lane stops failing on a model/transport issue that produces no marker, and reports the expected semantic failure instead.

## Acceptance criteria

- AC-1: `assertConflictOwnerHandoff` returns an error (not `t.Fatal`) when `fixture.marker` is missing, so the failure routes through `gradeLive` and grades as XFAIL for a bound target.
- AC-2: A present-but-wrong marker still returns the existing `worker marker = %q, want runtime-worker-owner` error unchanged.
- AC-3: The other reads in `assertConflictOwnerHandoff` (`fixture.entity`, `git rev-parse`, `git log`, `git show`, `git status`, `git worktree list`, `git branch`) are unchanged.
- AC-4: This change is runtime-neutral. It does not alter the `liveXFail("claude-opus", ...)` or `liveXFail("pi", ...)` bindings, and preserves `repair-pi-owner-handoff` (fe7) and `repair-opus-owner-handoff` (bqy) AC-3 ("preserve the owner-handoff assertion").
- AC-5: `gofmt`, `go vet -tags live ./internal/ensigncycle`, `go build -tags live ./internal/ensigncycle`, and the offline `ConflictOwner|GradeLive|LiveGrade|KeepMovingDurable` unit tests pass.

## Notes

- The package already has `readFileAllowMissing(path) string` (returns `""` on missing) in `shallow_boot_assert_test.go`. Use it for the marker read; keep `readFile(t, ...)` for `fixture.entity` (a fixture file the test itself wrote, so missing there is a real setup failure).
- This is a general live-test grading fix, not a Pi repair. Do not bundle it with the Pi runtime-specific repair entities.
