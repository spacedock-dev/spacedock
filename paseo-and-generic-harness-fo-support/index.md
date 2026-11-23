---
id: xdhbscjmpv0prk3k5yz1ngpw
title: Paseo-wide FO runtime support and generic unsupported-harness fallback hints
status: backlog
source: captain (2026-11-23)
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:xdhbscjmpv0prk3k5yz1ngpw:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:xdhbscjmpv0prk3k5yz1ngpw-backlog-1
              briefing:
                id: briefing:xdhbscjmpv0prk3k5yz1ngpw:backlog:attempt-1:revision-1
                digest: sha256:69a37a507dadd7f2fba03e3f3b7bdaea1f1a17f05c5036e44db8d3a34a79d253
                room-ref: '@review/backlog/briefing-1'
---

Make any Paseo-hosted first officer boot and run cleanly regardless of inner harness (fixing `F4` and Paseo-vs-native dispatch selection), audit whether any runtime-specific constructs remain when running under Paseo, and give unsupported harnesses clear guidance to map Spacedock's worker capabilities onto their own Claude/Codex-like native features.

## Problem

1. **Non-{Claude,Codex,Pi} Paseo harnesses abort on Boot Step 1 (`F4`):** `internal/cli/cli.go:913-930` (`printVersion`) returns early after 3 lines — omitting `Runtime:`, `Sandbox:`, and `contract 3` — whenever `runtimehost.Detect` sees none of `CLAUDECODE`, `CODEX_THREAD_ID`, or `PI_CODING_AGENT_DIR`. Any FO running inside Paseo on Antigravity (`agy`), Goose, OpenCode, Copilot, Cursor, etc. immediately aborts on `first-officer-shared-core.md:9`'s `contract 3` binary gate.
2. **Paseo vs. inner-host priority & sandboxed fallback:** Once `paseo-safehouse` `tv` preserves `PASEO_AGENT_ID` for all sandboxed agents, a normal sandboxed Claude/Codex/Pi FO inside Paseo (`paseo-access=workspace`, no daemon access) will have `PASEO_AGENT_ID` set too. `skills/first-officer/SKILL.md:38` currently selects `paseo-fo-dispatch.md` unconditionally on `PASEO_AGENT_ID` without checking whether the `paseo` daemon is reachable (and says "at first dispatch" instead of before resume reconciliation, per `pre-merge-audit.md` Section 4.1).
3. **Runtime-specific residue audit under Paseo:** Nearly everything in `claude-first-officer-runtime.md`, `codex-first-officer-runtime.md`, and `pi-first-officer-runtime.md` is native worker transport (spawn, wait, roster, shutdown); once Paseo owns worker dispatch, confirm whether any remaining host-specific constructs are actually required under Paseo or can stay cleanly decoupled.
4. **No fallback guidance for unsupported harnesses:** When an FO runs on a harness without a dedicated adapter (e.g. Antigravity outside Paseo, Cursor, OpenCode, Goose, etc., or a Paseo session without daemon access on an unsupported harness), Spacedock emits no hint explaining how to map Spacedock's generic capabilities (`spawn` fresh worker, async `wait`, `roster`, `shutdown`, direct captain text) onto that harness's own native features resembling Claude (`Agent`/`SendMessage`) or Codex (`spawn_agent`/`wait_agent`/`list_agents`).

## Proposed approach

1. **Binary `--version` support for Paseo (`F4`, `internal/cli/cli.go:printVersion`):**
   - Keep `internal/runtimehost.Detect` untouched (no ambiguity on Claude/Codex/Pi).
   - In `printVersion`, when `host == ""` and `PASEO_AGENT_ID != ""`, emit `Runtime: paseo (PASEO_AGENT_ID)`, `Sandbox: ...`, and `contract 3` (leaving human terminal output at 3 lines when neither a runtime marker nor `PASEO_AGENT_ID` is set).
2. **Paseo runtime priority & graceful fallback (`skills/first-officer/SKILL.md`):**
   - When `PASEO_AGENT_ID` is set and the `paseo` daemon is reachable (or captain explicitly selects Paseo), Paseo takes priority over Claude/Codex/Pi native subagent transports (and loads before resume reconciliation, not only at first dispatch).
   - When `PASEO_AGENT_ID` is set on a sandboxed Claude/Codex/Pi FO without daemon access (`paseo-access=workspace`), fall back cleanly to the inner host's native adapter instead of failing with `EPERM`.
3. **Unsupported-harness fallback hint (`skills/first-officer/SKILL.md` and/or `--version` / `dispatch build`):**
   - When running on an unsupported harness (neither Claude/Codex/Pi nor reachable Paseo), emit a concise, actionable hint instructing the FO to use the harness's own native equivalents resembling Claude (`Agent` / background subagent) or Codex (`spawn_agent` / `wait_agent` / `list_agents`) for fresh worker dispatch and async monitoring, while keeping durable Git commits and `## Stage Report` as the authoritative completion gate and communicating with the captain in direct conversation text.
4. **Runtime-construct audit:** Verify in ideation whether any remaining runtime-specific assumptions in `first-officer-shared-core.md` or `fo-dispatch-core.md` leak into Paseo-hosted or generic-harness runs, and document/remove them if unnecessary.
