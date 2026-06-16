---
title: Build a terminal/pty live harness for team-mode e2e (residency + teardown)
status: ideation
source: "FO + captain (2026-06-16, during 2yf): headless `claude -p` cannot sustain team mode — anthropics/claude-code 2.1.178 dropped the native TeamCreate/TeamDelete tools from headless sessions (anthropics/claude-code#68721), and even with tools present the SDK/headless session lifecycle races to end_turn before teammates finish (anthropics/claude-code-action#1124). Per 7e's recorded steer (headless `-p` goes bare), the two forced-team `-p` live tests (TestLiveEnsignCycleTeamTeardown, TestLiveStandingResidencyInjectsCommOfficer) were RETIRED in 2yf because they cannot work headless. Team-mode MECHANISMS stay covered offline (internal/dispatch/spawn_standing_all_test.go for comm-officer injection; internal/ensigncycle/teardown_grade_watcher_test.go + testdata/sonnet_teamdelete_*.jsonl for the bounded-teardown marker grading). The GAP this task closes: live end-to-end team-mode coverage (FO creates a real team, injects the comm-officer standing teammate into the roster, dispatches through the team, terminalizes, and runs the bounded teardown emitting TERMINAL_TEARDOWN_BOUNDED). That requires an INTERACTIVE (pseudo-terminal) harness where team tools are present and the session stays alive — the current live suite is entirely `claude -p`."
started: 2026-06-16T15:15:39Z
completed:
verdict:
score:
worktree:
issue:
id: m40mphxan8phr3t3tp03gk89
sprint: 0204-structured-reads
mod-block:
pr:
sprint-readiness: in-progress
---

## Problem

Team mode is interactive-only: it needs a live tty session where the native team tools (TeamCreate/TeamDelete) are exposed and the agent loop stays alive while teammates work. The entire live-e2e suite drives `claude -p` (headless), which cannot do either reliably. So the live end-to-end team-mode path — comm-officer roster injection and the bounded terminal teardown — has no working harness; it has only offline mechanism coverage.

## What's needed

A live harness that drives a real INTERACTIVE Claude session via a pseudo-terminal (pty), boots the FO, and exercises:
- standing-teammate (comm-officer) injection landing in the team `config.json` roster;
- the bounded terminal teardown emitting the `TERMINAL_TEARDOWN_BOUNDED` marker after a real TeamCreate→dispatch→terminalize→teardown cycle.

## Notes / prior art

- Retirement of the two forced-team `-p` tests + the go-bare alignment happened in 2yf (shared-merge-dispatch-contract). 7e (headless-dispatch-mode-intent, #381) determined headless `-p` goes bare. #271 (headless-fo-drive-flake) first surfaced the headless silent-await stall and deferred the runtime-await question.
- Mechanism coverage that already exists (do not duplicate): `internal/dispatch/spawn_standing_all_test.go`, `standing_parity_test.go`; `internal/ensigncycle/teardown_grade_watcher_test.go`, `teardown_grade_test.go`, `streamwatch_test.go` + `testdata/sonnet_teamdelete_*.jsonl`.
- Upstream constraints to watch: anthropics/claude-code#68721 (2.1.178 headless team-tool regression — may resolve and restore headless tools), anthropics/claude-code-action#1124 (SDK/headless team lifecycle).
