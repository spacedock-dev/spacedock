---
id: ecn07f3hwp5wgs8xf14h59sj
title: CI log hygiene — live-runner stream jsonl belongs in the artifact, not stdout
status: backlog
source: "captain (2026-06-14) — observed while debugging #368's opus `gate-guardrail` no-progress failure: the live-runner jsonl-to-stdout dump (`internal/ensigncycle/claude_live_runner_test.go:365`) bloated the CI log (~143KB for ~80 lines on one failed step) and buried the actual failure line."
started:
completed:
verdict:
score: "0.30"
worktree:
issue:
sprint: 0203-fo-efficiency
---

The live shared-scenario runner dumps the full host jsonl stream to CI stdout — `internal/ensigncycle/claude_live_runner_test.go:365` `t.Logf`s every parsed stream line — so the whole conversation lands in the test log. That bloats CI output and buries the actual test failure. The jsonl is ALREADY uploaded as a per-scenario artifact (`claude-stream.jsonl`), so the stdout dump is pure redundancy. Keep the stream in the artifact; keep stdout to clean test-framework output (FAIL lines, assertions) — or only dump the stream to stdout on failure.

## Problem

- `claude_live_runner_test.go:365` logs each parsed stream line to the test log → the entire conversation lands in CI stdout.
- Debugging #368's opus failure (2026-06-14) meant fighting ~143KB of jsonl to find one `claude_live_runner_test.go:101` "no-progress quiet budget" line.
- The jsonl is already captured to an artifact, so duplicating it on stdout adds log cost and reviewer friction with no benefit.
- Likely mirrored in the Codex / Pi runners — confirm at ideation.

## Out of scope

{Ideation fills — e.g. drop the stdout stream log entirely vs only-on-failure; whether to apply the same fix to the Codex/Pi runners in this task or a follow-on.}

## Acceptance criteria

{Ideation fills. The proof is checkable outside this task body: e.g. a scenario run's captured test stdout contains NO raw stream jsonl lines (no `"type":"assistant"`/`"type":"user"` lines; a bounded line count) WHILE the jsonl remains present in the uploaded `claude-stream.jsonl` artifact. Never a string/regex match over an instruction file.}

## Test plan

{Ideation fills.}
