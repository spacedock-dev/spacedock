---
title: Enforce Codex fresh-spawn isolation at the live FO boundary
status: ideation
source: "Captain-directed v0.25.2 follow-up, 2026-07-20; live FO escape after v0.25.1 / archived rt8 / PR #532"
score: "1.0"
milestone: 0.25.2
started: 2026-07-20T04:01:57Z
completed:
verdict:
worktree:
issue:
id: 6cc3rvfd44y6x3352hh21v8b
---

v0.25.1 hardened `internal/dispatch/codex_v2_adapter.go` so its generated map always carries `fork_turns: "none"`, and aligned the Codex FO prose. That did not close the live boundary: in a 2026-07-20 first-officer session, the FO directly invoked `spawn_agent` with `fork_turns: "all"`. The call was rejected only because it also carried an incompatible explicit agent type; the Spacedock architecture itself did not make inherited-turn spawning impossible.

This is the escaped value defect from archived `codex-fresh-dispatch-context-isolation` (`rt8`, PR #532), not a reopening of its historical release record. Its first validation correctly reported that adapter output and a one-off host probe did not prove the instruction-driven FO invocation. The task later narrowed the acceptance boundary to the adapter map and shipped v0.25.1, leaving the original live claim unenforced.

The 0.25.2 fix must bind the actual worker-creation seam used by an assumed Spacedock FO. In current architecture, `fork_turns` is never a selectable continuity mechanism: every fresh spawn is isolated, and deliberate continuity exclusively uses `followup_task` on the existing handle. `"all"`, numeric forks, omission that defaults to full history, and helper/runtime override channels are invalid Spacedock behavior.

Ideation must identify the smallest executable enforcement seam. If Spacedock cannot enforce this at the live tool boundary, stop with a proven upstream/runtime blocker rather than relabeling adapter or prose evidence as the fix. Coordinate with `per-host-stage-model-override` (`e3g`) so model/effort work cannot restore or conditionalize the isolation invariant.

## Acceptance criteria

**AC-1 (VALUE) — An instruction-driven Codex First Officer cannot issue a fresh Spacedock worker spawn that inherits parent turns.** Verify this with one integrated live record joining the generated dispatch artifact, the exact FO-issued spawn input before and after enforcement, and child-visible context. Seed the artifact with `"all"`, a numeric fork, and omission in separate cases; every executed spawn must contain exact `fork_turns: "none"`, and disabling the enforcement must make the parent canary visible to the child.

**AC-2 — Enforcement changes only the isolation field.** For every AC-1 case, compare the pre- and post-enforcement input objects: `fork_turns` changes to exact `"none"`; `model`, `reasoning_effort`, `service_tier`, and arbitrary unknown fields retain their original JSON values. A reconstruction that drops or rewrites any unrelated field fails.

**AC-3 — Deliberate continuity remains exclusively `followup_task` on the existing worker handle.** Verify this in the AC-1 live journey: the fresh child lacks the parent canary, then `followup_task` reaches that exact child and the child recalls its own first-turn marker. A second spawn is not continuity.

**AC-4 — Adapter maps, contract wording, hook structure, and raw host probes remain supporting evidence.** They cannot substitute for the AC-1 join. If the installed runtime cannot apply an input rewrite to the actual namespaced collaboration spawn, the gate is BLOCKED with that upstream dependency named; it is not passed by narrowing the boundary.

**AC-5 — v0.25.2 ships the fix on the stable line without rewinding `next`.** Verify the exact release candidate SHA with required Go/full/race gates and AC-1 through AC-3 live evidence, cut annotated `v0.25.2` from `main`, and retain the invariant on `next` through the documented propagation path.

## Scope guard

Do not build a new general agent harness or fork-mode configuration surface. Reuse the narrowest existing live capture path that can observe the real FO tool call. Do not expose `fork_turns` to workflow authors. Do not change stage reuse policy. Do not fold unrelated model/effort routing into this patch.

## Riskiest-mechanism spike and decision

The archive already contains the two halves of the problem, but not the required join. Top-level session `codex:019f7d9a-5b06-75a0-a04a-02b0b2ccd6a2` records the escaped FO call with `fork_turns: "all"` and its retry with `"none"`; the first call failed only on an incompatible agent type. Session `codex:019f79e5-3974-75e2-ab25-d6d07836cc72`, corroborated by parent session `codex:019f7007-8fba-7503-8c44-5ebf9a7cc945`, records a raw-host isolated child and successful same-handle follow-up. No archived run joins generated dispatch bytes, the exact FO-issued arguments, absence of a parent canary, and same-child continuity.

A disposable captain probe exercised the smallest plausible live seam in Codex CLI 0.144.6: a trusted `PreToolUse` hook on the actual collaboration tool call. The hook observed canonical tool name `agentsspawn_agent` and input containing `fork_turns: "all"`. It returned the documented `hookSpecificOutput.updatedInput`, preserving every input key and replacing only `fork_turns` with `"none"`. The child nevertheless saw the parent canary. A Bash control proved session hook loading worked, so the blocker is narrower: this runtime observes the namespaced collaboration call but does not apply its input rewrite before execution.

**Decision: BLOCKED upstream at Codex namespaced collaboration-tool input rewriting.** There is no Spacedock-owned executable seam that can currently satisfy AC-1. Adapter hardening and FO prose remain useful defense in depth, but neither may be relabeled as live enforcement. The disposable probe files were removed; this ideation adds no harness, recorder, fixture, parser, or permanent test infrastructure.

## Conditional implementation design

Implement the following only after the captain re-probe proves that the installed Codex runtime applies `updatedInput` to the live namespaced spawn:

1. Add a plugin `PreToolUse` binding for the runtime-verified collaboration spawn name, retaining compatibility aliases only when a live probe verifies them. Route it to a small public `spacedock dispatch normalize-codex-spawn` command rather than a plugin-private script.
2. Have the command parse the hook envelope, require an object `tool_input`, shallow-copy the whole input object, overwrite only `fork_turns` with exact `"none"`, and emit the supported `updatedInput` envelope. Reject malformed input rather than allowing an unnormalized spawn.
3. Leave the existing Codex adapter and First Officer request at explicit `fork_turns: "none"`. The boundary normalizer is unconditional: it must not depend on whether model or effort overrides are present.
4. Keep worker continuity unchanged: only `followup_task` may address an existing handle. No workflow field, CLI flag, helper option, or host-specific fork choice is introduced.

This is the smallest behavior-owning design because it acts immediately before the live spawn executes. A Go command gives the plugin a stable command surface and testable transformation; a private script would recreate the migration problem this repository exists to remove. A shallow copy preserves host evolution and `e3g` overrides; rebuilding a whitelist is cheaper but would silently discard model, effort, service-tier, or future fields. A one-off integrated captain proof is enough to validate value; a permanent general agent harness would add machinery without improving the live boundary. Failing closed on malformed or unsupported rewriting protects AC-1; trusting host defaults or prose is cheaper but reproduces the shipped defect.

## Coordination with `per-host-stage-model-override` (`e3g`)

The two changes compose at the final input object. `e3g` may supply `model`, `reasoning_effort`, and `service_tier`; this entity's normalizer preserves those values and overwrites only `fork_turns`. `e3g` must not add a fork-mode surface, make isolation conditional on an override, or restore omission/`"all"`/numeric values. Conversely, this work must not select, validate, or default model/effort values. The normalizer's preservation test is the shared integration contract.

## Test plan

- Add table tests for absent, `"all"`, numeric, and already-`"none"` fork inputs. Assert exact `"none"` after normalization and equality of every unrelated JSON value, including arbitrary unknown fields. Add malformed-envelope and non-object-input rejection cases. Mutating away the overwrite or reconstructing a field whitelist must fail these tests.
- Add a narrow hook/configuration smoke test proving the plugin invokes the public binary command for the runtime-verified canonical tool name. Treat this only as routing evidence, not AC-1 evidence.
- Run the shortest disposable captain probe below. Retain only its review artifact and JSONL as release evidence; do not productize the probe. The enabled run must join artifact, pre/post exact arguments, child canary absence, and same-handle follow-up. The disabled/identity-normalizer baseline must expose the parent canary and turn the proof red.
- On the release-candidate SHA run focused tests, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; then perform the existing `main` tag and `next` propagation checks for AC-5.

## Shortest captain re-probe

On a Codex build that claims namespaced collaboration-tool rewrite support, launch one disposable trusted-hook session matching the observed canonical `agentsspawn_agent`. Seed `PARENT_ONLY_CANARY_<random>`, make the FO consume a generated dispatch artifact that requests `fork_turns: "all"`, and then call `followup_task` on the returned child handle. Pass only if one record shows: artifact input `"all"`; hook input `"all"`; emitted and executed input exact `"none"` with all other fields preserved; child reports the parent canary absent; and the same child recalls a marker from its first turn. Repeat once with the normalizer disabled or identity-mutated; the child must see the parent canary. Any failure leaves this entity BLOCKED.

## Documentation delta

No public workflow or release documentation changes: fresh isolation is already the promised behavior, and no user-selectable setting is added. After the re-probe is green, update only `skills/first-officer/references/codex-first-officer-runtime.md` to say that the FO requests exact `fork_turns: "none"` and the plugin boundary normalizer reasserts it immediately before execution while preserving model/effort/service-tier and unknown fields; state that continuity uses `followup_task` exclusively. Do not document the conditional design as available before AC-1 passes.

## Stage Report: ideation

- DONE: Identify/exercise smallest actual live FO-to-Codex enforcement seam; use existing logs or short captain probe; stop on proven upstream blocker.
  Archive evidence established the escaped live FO call and the separate successful raw-host isolation/continuity probe. A disposable probe then exercised the live `agentsspawn_agent` `PreToolUse` seam: it observed the call and emitted a correct rewrite, but Codex CLI 0.144.6 executed inherited context. The entity stops at the named upstream rewrite blocker and the scratch probe was removed.
- DONE: Concrete behavior-first design, ACs, test plan preserving integrated artifact -> exact FO spawn args -> child-visible context, continuity followup_task only.
  AC-1 retains the full artifact-to-arguments-to-child join and requires a red inherited-context baseline. AC-2 protects unrelated fields, AC-3 proves same-handle `followup_task` continuity, and the test plan separates unit/routing support from the required disposable integrated live proof.
- DONE: Coordinate with per-host-stage-model-override e3g, name every mechanism’s value and rejected cheaper alt, exact minimal docs, no model/effort or fork-mode surface.
  The conditional normalizer shallow-copies all input fields and overwrites only `fork_turns`, so `e3g` owns model/effort values without weakening isolation. The design rejects adapter/prose substitution, private scripts, field whitelists, permanent harnesses, host defaults, and any fork selector; only the Codex runtime binding note changes after live proof passes.

### Summary

The ideation found the correct narrow boundary but proved it is not presently enforceable: Codex CLI 0.144.6 observes the namespaced collaboration spawn in `PreToolUse` yet does not apply `updatedInput`. Implementation is therefore blocked on upstream live-tool rewrite support; once that support exists, the specified shallow-copy normalizer and integrated captain probe are the minimum path to AC-1 without colliding with `e3g`.
