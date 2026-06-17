---
title: Pin CI live-e2e to a known-good Claude Code version (team tools regress at 2.1.178+)
status: backlog
source: captain (2026-06-17) — team-tool availability regression
score: 0.5
id: 61an320jkjgeyqq55rfx7nw0
---

The live-e2e CI lane (`.github/workflows/runtime-live-e2e.yml`) floats the Claude Code version: its `claude_version` workflow_dispatch input defaults to empty, which falls through to the installer's latest. Latest is now 2.1.178+, where headless `claude -p` lost the team tools (TeamCreate/SendMessage), and 2.1.179 interactive regressed too (established by `m4`, live-team-mode-terminal-harness). Team-mode FO drives and the standing level-3-judge residency both need a version where team tools work; 2.1.177 is the last known-good (team tools confirmed present in its registry).

Pin CI on 2.1.177 so the live lanes run against a deterministic, team-tool-capable Claude Code instead of whatever floats to latest.

Brief scope (ideation fleshes and verifies):
- Make 2.1.177 the effective default for the live-e2e lane's Claude Code install, and audit whether any other CI install point (release.yml, install-e2e.yml, scheduled triggers) also floats the version and needs the same pin.
- Produce a checkable proof that runs actually resolve to 2.1.177 (e.g. the lane asserts `claude --version` reports 2.1.177), so the pin is enforced, not just declared.
- Record WHY 2.1.177 (team-tool regression at 2.1.178+) in or beside the workflow so the version is not a mystery magic number to a future reader.

This is a CI-machinery change — a high-stakes surface under the proof policy: implementation in a worktree, validation runs the detached adversarial audit.

Linkage: pinning CI to a team-tool-capable version unblocks the `72` (fo-tier-delegation) residency decision — the standing level-3-judge residency becomes viable in the CI env once the pin lands.
