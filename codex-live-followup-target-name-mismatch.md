---
title: Codex live lifecycle fails on a follow-up target name mismatch
sprint: test-behavior-completeness
source: "Run 37421110550 codex-live, after the CI-E2E-CODEX credential refresh: codex_multi_agent_test.go:143: structured launcher lifecycle: follow-up target \"worker\" != spawned worker \"/root/worker...\"."
id: phwzty950t2de11g4g7h95me
---

## Problem

`TestCodexIsolatedHomeCollaborationLifecycle` fails on a follow-up target naming a different
worker identity than the one spawned. Before the credential refresh this was hidden behind
`workspace routing discovery unauthorized (401)`; a fresh `CI-E2E-CODEX` secret cleared the 401
and exposed this.

## Value

The codex lane is red on every pull request, and this is the reason once auth is out of the way.

## Acceptance criteria

**AC-1** The live test passes in CI, or the mismatch is shown to be environment-specific and the expectation corrected.
**AC-2** Reproduce locally with `SPACEDOCK_CODEX_REAL_BIN`, an isolated `CODEX_HOME`, and `SPACEDOCK_LIVE_CODEX_MULTI_AGENT=1`, and record whether the failure is a product fault or a CI path artifact.

## Verification

Run the focused live target locally, then confirm in the lane.
