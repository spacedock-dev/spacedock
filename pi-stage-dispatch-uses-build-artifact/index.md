---
title: Pi stage dispatches should use dispatch build artifacts, not hand-rolled prompts
status: backlog
source: captain (2026-06-04) — FO manually composed Pi subagent task prompts for fc/d2 instead of routing the canonical spacedock dispatch build artifact that carries entity slug/stage context
score: "0.30"
started:
completed:
verdict:
worktree:
issue:
id: z68h8vwxeetp011b1484c2jx
---

Spacedock stage dispatches should be driven by the canonical `spacedock dispatch build` artifact. That artifact is responsible for carrying the entity slug, entity path, workflow directory, target stage, stage definition, worktree path, checklist, and host/runtime constraints. During recent Pi work, the FO manually composed `subagent(...)` task prompts for `launcher-binary-path-passthrough` and `pi-stage-dispatch-fresh-context`, adding slug/stage details by hand. That was sufficient for the moment but bypasses the dispatch-builder contract and makes omissions likely.

## Problem

Manual Pi subagent prompts can drift from the Spacedock stage contract. They may omit entity slug/stage, use stale stage definitions, miss split-root state paths, omit worktree constraints, or phrase completion requirements differently from the host-neutral dispatch builder. This weakens parity with Claude/Codex dispatch and makes it harder to test the FO's behavior.

## Desired outcome

Pi FO stage dispatches should wrap the canonical dispatch-builder output, not replace it. The Pi-specific wrapper may add runtime transport parameters such as `context: "fresh"`, display `phase`/`label`, and the no-`acceptance` rule, but the assignment content should come from `spacedock dispatch build` or its emitted dispatch file.

## Acceptance criteria

**AC-1 - Pi stage dispatch instructions require use of `spacedock dispatch build`.**
Verified by: a skill/integration invariant test over Pi first-officer runtime guidance that fails unless the Dispatch section tells the FO to build the assignment with `spacedock dispatch build` and forward the emitted dispatch-file prompt/content to the subagent.

**AC-2 - Pi subagent wrapper constraints are additive, not a replacement assignment.**
Verified by: an invariant test or runtime fixture that requires `context: "fresh"`, no `acceptance`, and optional `phase`/`label` to be applied around a dispatch-builder artifact while preserving entity slug/stage/workflow/worktree/checklist fields from the artifact.

**AC-3 - Generated Pi dispatch artifacts include the fields needed to name and audit the worker.**
Verified by: `internal/dispatch` tests for `host: "pi"` asserting the artifact/prompt contains entity slug, target stage, entity path, workflow directory, and worktree path when applicable.

**AC-4 - FO manual-prompt fallbacks are limited and explicit.**
Verified by: docs/invariant tests that allow hand-written subagent prompts only for debugging or when `dispatch build` is unavailable, and require the fallback prompt to state why the canonical artifact could not be used.

## Notes

This is distinct from `pi-stage-dispatch-fresh-context`: that task controls the subagent context boundary. This task controls the assignment source of truth.
