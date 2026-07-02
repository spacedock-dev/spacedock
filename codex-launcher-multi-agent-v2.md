---
id: 52j8z0hhbtgq2afqybb2482j
title: "spacedock codex launcher enables multi_agent_v2"
status: backlog
group: tooling
source: "Codex session 2026-07-02: operator config has features.multi_agent_v2=true, but spacedock codex launches can run with isolated or alternate CODEX_HOME and should not silently lose the v2 collaboration surface."
started:
completed:
verdict:
score:
worktree:
issue:
---

Ensure `spacedock codex` launches Codex with the intended multi-agent collaboration surface, including the current `multi_agent_v2` rollout knob, without relying on every operator or CI environment to carry the same user config.

## Problem

Current Codex releases document `features.multi_agent` as stable and default-enabled, while this environment also uses an under-development `features.multi_agent_v2` knob. The operator's normal Codex config has `multi_agent_v2 = true`, but Spacedock-launched Codex sessions may use a different or isolated `CODEX_HOME`, especially in live CI. That means the first-officer runtime can appear to work locally and then lose the expected v2 collaboration surface under `spacedock codex` or the live harness.

The launcher should make the intended surface explicit at the process boundary. Skill prose should not be the mechanism, because the tool surface is decided before the agent reads any skill.

## Proposed approach

Teach the Codex front door to add Codex config overrides when it assembles the launched argv:

- enable the stable subagent collaboration feature (`features.multi_agent=true`, or the host-supported equivalent)
- enable the current v2 rollout key (`features.multi_agent_v2=true`)

Prefer default-on semantics where Spacedock's injected overrides come before user passthrough flags, so an operator can still intentionally disable or vary the feature with a later `--disable` or `-c ...=false`. If implementation decides the v2 surface is required for Spacedock's Codex runtime contract, replace default-on with an explicit fail-fast when passthrough disables it, and document that decision in this task before coding.

The Codex live runner should use the same launcher path or pass the same overrides directly, so local `spacedock codex` sessions and isolated-CODEX_HOME CI runs do not diverge.

Also use this task to retire the old foreground-wait "resume/interruption" hint if `multi_agent_v2` proves the newer mailbox behavior. Keep the `codex resume` launch-subcommand guard separate: resumed sessions still should not receive a fresh bootstrap prompt, but that is unrelated to whether the first officer must keep warning the operator that an interrupted `wait_agent` will be reinstalled.

## Out of scope

Do not redesign Codex runtime reuse semantics in this task beyond binding the v2 follow-up surface already exposed by Codex. Do not add new Spacedock workflow stages. Do not depend on a public `multi_agent_v2` documentation guarantee; treat it as a current host knob that may need a compatibility guard.

## Acceptance criteria

**AC-1 - `spacedock codex` makes multi-agent enablement explicit at launch.**
Verified by: a focused `runCodex` unit test with a fake host asserts the launched argv contains a stable multi-agent enablement override and the `features.multi_agent_v2=true` override.

**AC-2 - operator override semantics are intentional and tested.**
Verified by: a front-door test covers passthrough that disables or changes the multi-agent feature. The test documents whether user passthrough wins or Spacedock fails fast, and the implementation behaves accordingly.

**AC-3 - launch-session resume and safehouse paths keep their existing behavior.**
Verified by: existing resume/safehouse front-door tests still pass, plus at least one assertion that the new config override does not reintroduce the bootstrap prompt on `codex resume` or break safehouse wrapping.

**AC-4 - isolated-CODEX_HOME live runs receive the same multi-agent setting.**
Verified by: the Codex live runner either calls the same launcher path or asserts/persists the argv/config evidence showing the stable multi-agent flag and `features.multi_agent_v2=true` are present for `codex exec`.

**AC-5 - runtime capability is proven by behavior, not argv alone.**
Verified by: a live or fixture-backed test records that a Spacedock-launched Codex session exposes the expected collaboration tools for the v2 surface, including spawn, wait, and turn-starting follow-up when available. The evidence should be a tool-surface or transcript artifact, not a prose claim in this task.

**AC-6 - the old foreground-wait resume hint is retired or justified.**
Verified by: running the Codex idle-notification probe, or an equivalent live v2 probe, and then updating `codex-first-officer-runtime.md` plus tests to remove/demote the mandatory "interruption returns control; reinstall wait" operator cue when v2 no longer needs it. If the probe shows the hint is still needed, the stage report must cite the evidence and keep the hint intentionally.

## Test plan

Start with unit tests around `runCodex` argument assembly. Add the minimal launcher change. Then run the focused front-door tests for Codex, followed by the Codex live runner or a targeted live probe that proves the actual collaboration surface under an isolated Codex home. Finally, run the idle-notification probe to decide whether the foreground-wait resume hint can be removed from the Codex first-officer runtime.
