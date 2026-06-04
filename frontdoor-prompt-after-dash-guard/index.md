---
id: nz2aae5kfk1rb4e4z62vsnv1
title: Front-door silently swallows a positional prompt placed after `--` — bootstrap never prepends
status: ideation
source: captain (2026-06-03, session 12) — launched `spacedock claude --safehouse-enable=… --plugin-dir "$(pwd)" -- --model … --effort ultracode '@/tmp/handoff.md'` and the bootstrap/default prompt never prepended; the `@file` after `--` was treated as host passthrough so hasTask=false
started: 2026-06-04T06:30:55Z
completed:
verdict:
score: 0.28
worktree:
issue:
---

The front-door grammar (`internal/cli/frontdoor.go` `parseFrontDoorArgs`) splits args at `--`: tokens before `--` become the spacedock *task* (`fd.task`, `hasTask=true`), tokens after `--` forward verbatim to the host. The bootstrap prompt is only combined with the operator's prompt when the prompt is the task — `launchPrompt` returns `base + " " + task` only when `hasTask`. When the operator places their positional prompt AFTER `--` (intuitive, since "the prompt goes to claude"), it is classified as host passthrough, `hasTask` stays false, and the bare bootstrap is appended as a separate trailing positional after the operator's prompt. The host then receives two positionals and the spacedock bootstrap is effectively lost — the operator silently loses the launch-and-go preamble with no warning.

The grammar inversion (host flags after `--`, task before) was a deliberate fix for a value-taking host flag swallowing the prompt; this task does NOT propose reverting it. It proposes making the failure mode legible instead of silent.

## Problem

A positional (non-flag) token after `--` is almost always an operator who meant it as the launch prompt. Today it silently degrades to host passthrough, the bootstrap is misplaced, and there is no feedback. The captain hit this directly in session 12.

## Proposed approach

{Ideation fleshes this out. Candidate directions to evaluate:}
- Detect a bare positional (non-flag, non-flag-value) among the after-`--` passthrough and emit a warning to stderr naming the corrected form (task before `--`), without changing launch behavior.
- Or refuse with an actionable error when a positional appears after `--` and no task was given before it.
- Decide the boundary against legitimate host positionals (e.g. claude `-p <prompt>`, codex `exec <prompt>`) so the guard does not false-positive on a real host positional.

## Out of scope

- Reverting the host-flags-after-`--` grammar inversion.
- `$0` / binary-path propagation (that is `fc` launcher-binary-path-passthrough).

## Acceptance criteria

{Ideation defines these with a test for each. The behavioral claim — "a positional prompt after `--` no longer silently vanishes" — must be proven by a front-door parse/launch test asserting the warning or error and the assembled argv, not by prose.}

## Test plan

{Ideation. Likely `internal/cli` front-door parse tests over the assembled inner argv, in the style of `frontdoor_parse_test.go` / `plugin_dir_frontdoor_test.go`.}
