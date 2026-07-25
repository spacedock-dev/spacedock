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
gates:
    version: 1
    current:
        gate: gate:docs-dev:v3:backlog
    records:
        - id: gate:docs-dev:v3:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:v3-backlog-1
              briefing:
                id: briefing:docs-dev:v3:backlog:attempt-1:revision-1
                digest: sha256:2cda6cc97a6f5db6ef432781f25aae29f9781316c73999451d134ac549d26dc6
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:v3:backlog:1
                briefing: briefing:docs-dev:v3:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-25T08:11:16.055672Z"
                decision: approve
                reason: Two retained 6y Claude journeys prove the live harness leaks the Codex runtime marker; the narrow ideation preserves production ambiguity refusal and is required for a truthful host journey.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge.
              application:
                action: advance
                target-stage: ideation
                state: pending
                blockers: []
---

The cross-host live harness must simulate the launched host, not leak the captain host's runtime identity. `isolatedClaudeEnv` removes `CLAUDECODE` before the front door but currently retains `CODEX_THREAD_ID` (and potentially Pi markers). The launched Claude session therefore sees multiple runtime families, making ordinary `spacedock dispatch build` fail as ambiguous and encouraging an explicit-host recovery that a real top-level Claude session would not need.

Keep the production ambiguity guard. Correct only live-runner environment isolation, with the same principle applied symmetrically where existing Codex/Pi runners can inherit foreign runtime markers.

## Acceptance criteria

**AC-1** A Claude live child launched from a Codex parent reaches its session with Claude identity and no Codex/Pi runtime markers; ordinary flag-free host derivation succeeds as Claude.

**AC-2** Existing deterministic environment tests fail if any foreign runtime-family marker leaks through the per-host live environment builder, while preserving required credentials, PATH, HOME, and host-specific state.

**AC-3** The recorded-gate live journey no longer emits the mixed-marker ambiguity or uses an explicit `--host` recovery solely because of harness leakage.

**AC-4** Production mixed-marker detection and its refusal tests remain unchanged.
