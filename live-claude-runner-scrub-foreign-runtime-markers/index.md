---
title: Scrub foreign runtime markers from live Claude journeys
status: backlog
source: "Repeated 6y live evidence on 2026-07-25: nested Claude inherited CODEX_THREAD_ID from the Codex captain host and first failed dispatch host detection before recovering with explicit --host claude"
started:
completed:
verdict:
score: 0.7
worktree:
issue:
sprint: durable-decisions
id: v3vt8gp2yffmn62r8p95gkph
---

The cross-host live harness must simulate the launched host, not leak the captain host's runtime identity. `isolatedClaudeEnv` removes `CLAUDECODE` before the front door but currently retains `CODEX_THREAD_ID` (and potentially Pi markers). The launched Claude session therefore sees multiple runtime families, making ordinary `spacedock dispatch build` fail as ambiguous and encouraging an explicit-host recovery that a real top-level Claude session would not need.

Keep the production ambiguity guard. Correct only live-runner environment isolation, with the same principle applied symmetrically where existing Codex/Pi runners can inherit foreign runtime markers.

## Acceptance criteria

**AC-1** A Claude live child launched from a Codex parent reaches its session with Claude identity and no Codex/Pi runtime markers; ordinary flag-free host derivation succeeds as Claude.

**AC-2** Existing deterministic environment tests fail if any foreign runtime-family marker leaks through the per-host live environment builder, while preserving required credentials, PATH, HOME, and host-specific state.

**AC-3** The recorded-gate live journey no longer emits the mixed-marker ambiguity or uses an explicit `--host` recovery solely because of harness leakage.

**AC-4** Production mixed-marker detection and its refusal tests remain unchanged.
