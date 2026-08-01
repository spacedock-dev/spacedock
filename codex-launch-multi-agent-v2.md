---
id: z5gwwz2748sg6vxr0g3kdsar
title: "Codex launcher guarantees multi-agent v2 surface"
status: ideation
source: "Captain feedback, 2026-07-02: ordinary Codex config currently enables multi_agent_v2; make Spacedock-launched Codex enable or prove the same surface instead of relying on ambient session setup."
started: 2026-08-01T14:29:25Z
completed:
verdict:
score: 1.0
worktree:
issue:
sprint: durable-decisions
---

Spacedock's Codex front door should launch a first-officer session with the same multi-agent v2 control surface an operator expects from a normal Codex session. The guarantee should live in launcher behavior, config setup, or an explicit preflight, not only in prompt text or an assumed user-level config file.

## Problem

The Codex launcher owns plugin installation, optional sandbox wrapping, and the first-officer bootstrap prompt, but the current launch path does not visibly force or verify the Codex multi-agent v2 feature. That leaves several fragile cases:

- a clean or isolated Codex home may not carry the operator's feature config;
- a sandboxed launch may preserve the binary invocation but miss the config needed to expose the v2 tools;
- live CI or fresh user machines may silently start with a weaker multi-agent surface, changing dispatch/reuse behavior after Spacedock has already begun operating.

The user-visible failure mode is not "a flag is missing"; it is that Spacedock thinks it launched an interactive multi-agent cockpit, but the first officer sees a reduced or older tool surface.

## Desired direction

Ideation should confirm the exact Codex-supported way to enable the feature at launch time. Prefer a launcher-owned guarantee: pass a supported config override, create/merge a minimal temporary config for isolated launches, or preserve the operator's configured Codex home through the wrapper intentionally. If Codex cannot expose the v2 surface, the launcher should fail early with a clear diagnostic or the first officer should record the downgraded surface before dispatch.

Do not solve this by adding stronger prose to the bootstrap prompt. The proof has to observe effective launch behavior: argv/env/config seen by Codex, the resulting tool surface exposed to the session, or durable live-run evidence.

## Out of scope

- Changing Codex itself.
- Redesigning Spacedock's dispatch or reuse contract.
- Making multi-agent v1 behavior equivalent to v2.
- Depending on a particular developer's local config as the only proof.

## Acceptance criteria

**AC-1 - Spacedock-launched Codex starts with the multi-agent v2 surface when the host Codex supports it.**
Verified by: a launcher or integration test that starts from an isolated Codex home with no preexisting user feature config and proves the effective launched session receives the v2 multi-agent surface.

**AC-2 - The launcher does not silently downgrade to a weaker multi-agent surface.**
Verified by: a negative-path test where the feature cannot be enabled or observed; the launch fails with an actionable diagnostic, or the first-officer boot record explicitly marks the downgrade before any worker dispatch can proceed.

**AC-3 - The guarantee covers the normal Codex launch variants Spacedock owns.**
Verified by: tests for plain Codex launch, local-plugin launch, and sandbox-wrapped launch showing the same feature-enabling path or preflight behavior is applied.

**AC-4 - Runtime CI and operator documentation align with the launcher-owned guarantee.**
Verified by: runtime-live setup and any affected user-facing docs stop relying on a developer's personal Codex config as the only source of the feature, and instead use the same launcher-backed path or document the required Codex support boundary.

## Test plan

Ideation should first spike the smallest observable launch path: run a fake or instrumented Codex binary under Spacedock's Codex front door and record which argv, environment, and config inputs are actually delivered in plain and wrapped launches. Use that spike to choose between config override, temporary config, preserved Codex home, or explicit preflight.

Implementation should add focused launcher tests before changing behavior, then run the relevant front-door tests and the full Go suite. If the chosen mechanism affects live runtime setup, finish with a small Codex live smoke that proves the first officer receives the v2 multi-agent control surface in a clean launch context.
