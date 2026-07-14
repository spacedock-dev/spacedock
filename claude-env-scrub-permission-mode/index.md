---
id: d1sah62r0xckeysjet5j6zzk
title: Preserve Claude first-officer permissions with subprocess environment scrubbing
status: backlog
source: GitHub issue spacedock-dev/spacedock#504, reported by jesserobbins 2026-07-14
started:
completed:
verdict:
score:
worktree:
issue: spacedock-dev/spacedock#504
---

`CLAUDE_CODE_SUBPROCESS_ENV_SCRUB=1` should remain usable with `spacedock claude` without silently forcing a dispatched first officer back to Claude Code's default, prompt-on-everything permission mode.

## Problem

Issue #504 reports that a user who enables Claude Code's recommended subprocess credential scrubbing globally sees every `spacedock claude` launch warn that permission mode was forced to `default`. Spacedock launches Claude Code as a subprocess, and the scrub hardening requires an explicit allowed-tool declaration for the child. The current interaction makes users choose between credential protection and an operable first officer.

Reproduction seed from the issue: put `"CLAUDE_CODE_SUBPROCESS_ENV_SCRUB": "1"` in the `env` block of `~/.claude/settings.json`, launch `spacedock claude` from a commissioned workflow, and observe Claude Code's warning plus default permission mode.

## Proposed approach

Ideation must first exercise the current Claude Code launch path and verify the supported mechanism and configuration boundary rather than assuming the warning's `allowedTools` hint maps directly to Spacedock's CLI invocation. Design the smallest compatibility-first change that preserves environment scrubbing and the intended first-officer tool authority. If the host cannot accept an explicit safe declaration at launch, define a precise doctor diagnostic and user-facing remediation instead of weakening or clearing the user's hardening setting.

## Out of scope

- Disabling, clearing, or recommending removal of `CLAUDE_CODE_SUBPROCESS_ENV_SCRUB` as the product fix.
- Broad redesign of Claude permission modes or policy outside the spawned first-officer path.
- Changes to Codex or Pi launch behavior unless shared launcher code requires a no-regression assertion.

## Acceptance criteria

Ideation will refine these into end-state properties after the live mechanism spike.

**AC-1 - A Claude first officer launched with subprocess environment scrubbing enabled retains the intended dispatched tool authority without exposing scrubbed credentials.**
Verified by: a live or faithful fixture-backed launch journey that observes the permission mode/tool behavior and the child environment boundary.

**AC-2 - Users receive an actionable, stable diagnostic when their installed Claude host cannot support the safe launch configuration.**
Verified by: exact-output doctor or launcher tests covering supported and unsupported host behavior without advising users to disable scrubbing.

## Test plan

Ideation must record the live host/version mechanism spike, identify the owning launcher/configuration surface, and propose focused command/fixture coverage plus a protected Claude live journey when runtime behavior is the claim. Implementation and validation retain the repository baseline gates from `docs/dev/README.md`.
