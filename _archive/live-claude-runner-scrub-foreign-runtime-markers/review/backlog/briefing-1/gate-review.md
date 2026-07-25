# Gate review: scrub foreign runtime markers from live host journeys

## Capability

The live harness launches each supported host with that host's runtime identity and no
runtime-family markers inherited from the captain host. Ordinary flag-free
`spacedock dispatch build` host detection then exercises the same path as a real
top-level session.

## Artifact summary

This review describes the repeated 6y live-test failure mode, the narrow harness-only
correction, its acceptance evidence, and the production boundary that must remain
unchanged.

## Evidence

- Two independent 6y Claude journeys launched from Codex inherited both
  `CODEX_THREAD_ID` and `CLAUDECODE`.
- Ordinary host detection rejected the mixed markers before the agent recovered with
  explicit `--host claude`.
- `isolatedClaudeEnv` currently scrubs `CLAUDECODE` before entering the Claude front
  door but does not scrub the foreign Codex or Pi runtime markers.

## Boundaries

- Preserve production mixed-marker detection and refusal tests.
- Change only per-host live-runner environment isolation and its deterministic tests.
- Do not fold the unrelated 6y review/ancestry oracle correction into this ticket.

## Recommendation

Approve ideation so the worker can identify the smallest symmetric host-environment
contract and exact fixture proof before implementation.

## Decision

Approve to enter ideation; revise to narrow the harness boundary; or hold if the live
journeys should continue with explicit-host recovery.
