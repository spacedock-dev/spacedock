---
title: Drop the non-meaningful activeSessionFile characterization subtest
status: implementation
source: "Captain review (CL) 2026-07-22 — mp (#548) gave four tautological tests teeth; three earn their keep, but #1 tests a test-local port not production behavior. Delete it."
id: sf2d5zawymn5n86g8nt597nh
worktree: .worktrees/spacedock-ensign-drop-activesessionfile-characterization-subtest
started: 2026-07-22T05:44:43Z
---

Delete the `activeSessionFile_would_flip_to_teammate` subtest from TestFOSessionPinning — it characterizes a test-local port's mod-time tiebreak, not the shipped session-resolution behavior, so it is not worth keeping even with teeth.

## Problem
mp #1 (`internal/ensigncycle/pty_session_test.go`, the `activeSessionFile_would_flip_to_teammate` subtest) asserts that `activeSessionFile` picks the teammate file. But `activeSessionFile` is a TEST-LOCAL port — a reference-copy of production session-resolution logic (flagged at `pty_live_driver_test.go:570`) — so it guards a fixture's mod-time tiebreak, not production. The captain judged it a delete, not a give-it-teeth.

## Proposed approach
Remove only the `t.Run("activeSessionFile_would_flip_to_teammate", ...)` subtest and any import it alone introduced (the `time` import, if now unused). Leave `activeSessionFile` (used by pty_live_driver_test.go) and every other TestFOSessionPinning subtest intact.

## Acceptance criteria
**AC-1 - the subtest is gone with no other coverage lost.**
Verified by: the subtest is absent from pty_session_test.go; `go test ./internal/ensigncycle` green; a diff review shows only that subtest (and its os.Chtimes/time additions) removed, nothing else touched.

## Test plan
gofmt; go build ./...; go test ./internal/ensigncycle; go test ./... -race. No behavior change; test-file-only.
