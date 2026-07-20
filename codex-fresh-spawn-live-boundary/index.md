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

**AC-1 (VALUE) — An instruction-driven Codex First Officer cannot execute a fresh Spacedock worker spawn that inherits parent turns.** Verify this with one integrated live record joining the generated dispatch artifact, each exact FO-issued spawn input, the boundary decision, and child-visible context. Omitted, `"all"`, and numeric `fork_turns` inputs must be denied before child creation; exact `"none"` must be allowed and its child must lack the parent canary. Disabling the guard must allow an unsafe spawn and expose the canary, turning the proof red.

**AC-2 — The guard decides only from the isolation field and never rewrites tool input.** Exact `"none"` is allowed with the FO-issued argument object unchanged. Missing, `"all"`, numeric, non-string, or malformed input is denied. Repeat allowed and denied cases with `model`, `reasoning_effort`, `service_tier`, and arbitrary unknown fields; those fields must not affect the decision and must reach the child unchanged on an allowed call.

**AC-3 — Deliberate continuity remains exclusively `followup_task` on the existing worker handle.** Verify this in the AC-1 live journey: the fresh child lacks the parent canary, then `followup_task` reaches that exact child and the child recalls its own first-turn marker. A second spawn is not continuity.

**AC-4 — Adapter maps, contract wording, hook structure, and raw host probes remain supporting evidence.** They cannot substitute for the AC-1 join. The supported launch must load and trust the guard, malformed input must exit `2` and block, and release validation must positively observe the guard decision. A missing, disabled, untrusted, or failed hook invalidates the release claim instead of narrowing the acceptance boundary.

**AC-5 — v0.25.2 ships the fix on the stable line without rewinding `next`.** Verify the exact release candidate SHA with required Go/full/race gates and AC-1 through AC-3 live evidence, cut annotated `v0.25.2` from `main`, and retain the invariant on `next` through the documented propagation path.

## Scope guard

Do not build a new general agent harness or fork-mode configuration surface. Reuse the narrowest existing live capture path that can observe the real FO tool call. Do not expose `fork_turns` to workflow authors. Do not change stage reuse policy. Do not fold unrelated model/effort routing into this patch.

## Riskiest-mechanism spike and decision

The archive already contains the two halves of the problem, but not the required join. Top-level session `codex:019f7d9a-5b06-75a0-a04a-02b0b2ccd6a2` records the escaped FO call with `fork_turns: "all"` and its retry with `"none"`; the first call failed only on an incompatible agent type. Session `codex:019f79e5-3974-75e2-ab25-d6d07836cc72`, corroborated by parent session `codex:019f7007-8fba-7503-8c44-5ebf9a7cc945`, records a raw-host isolated child and successful same-handle follow-up. No archived run joins generated dispatch bytes, the exact FO-issued arguments, absence of a parent canary, and same-child continuity.

A first disposable captain probe exercised input rewriting in Codex CLI 0.144.6. A trusted `PreToolUse` hook observed canonical tool name `agentsspawn_agent` and returned documented `updatedInput`, but the child still inherited the parent canary. That proves rewriting is not the viable seam in this runtime.

A second disposable probe exercised the cheaper fail-closed decision on that same canonical tool path. The hook denied omission, `"all"`, and numeric `"1"`; all three calls returned `Tool call blocked by PreToolUse hook` and created no child. It allowed exact `"none"`; child `/root/allow_none` reported the parent canary absent and marker `CHILD-MARKER-7KQ2-N9V4`, then `followup_task` on that same handle recovered the marker. The hook log recorded the four exact inputs and decisions. Evidence is in parent session `codex:019f7dd1-f012-7e32-9b14-c1e9278390c5` and child session `codex:019f7dd2-590e-7cd0-a6ea-f9729b8aefaf`.

**Decision: PROCEED with a deny-only guard.** Namespaced input rewriting remains unsupported, but it is unnecessary: denial on the same live seam is effective. The guard makes unsafe calls impossible to execute and leaves safe calls untouched. Adapter hardening and FO prose remain defense in depth. Both disposable probe directories were removed; this ideation adds no permanent harness, recorder, fixture, parser, or test infrastructure.

## Implementation design

1. Add a plugin `PreToolUse` binding for the runtime-verified canonical `agentsspawn_agent` name. Route it to a small public `spacedock dispatch guard-codex-spawn` command rather than a plugin-private script. Add aliases only when a live host proves a different canonical name.
2. Have the command parse the hook envelope and require an object `tool_input` whose `fork_turns` value is the exact string `"none"`. Emit documented `permissionDecision: "allow"` without `updatedInput` for that single valid case. Emit `permissionDecision: "deny"` for omission, `"all"`, numeric strings, non-string values, and every other value. Exit `2` with a concise reason when the envelope itself cannot be parsed so malformed input also blocks.
3. Leave the existing Codex adapter and First Officer request at explicit `fork_turns: "none"`. The guard is unconditional: it must not depend on whether model or effort overrides are present, and it must never repair or retry a rejected call.
4. Keep worker continuity unchanged: only `followup_task` may address an existing handle. No workflow field, CLI flag, helper option, or host-specific fork choice is introduced.

This is the smallest behavior-owning design because it decides immediately before the live spawn executes and uses the runtime behavior the spike proved. A Go command gives the plugin a stable, testable command surface; a private script is cheaper but recreates the plugin-script migration problem. A deny-only predicate avoids the unsupported rewrite mechanism and cannot drop `e3g` fields; reconstructing or mutating the input is more complex and already failed live. Explicit denial is cheaper than a general agent harness and stronger than adapter/prose checks. Failing closed on malformed input protects AC-1; trusting host defaults, a hook error, or FO self-correction is cheaper but reproduces the escaped defect.

## Coordination with `per-host-stage-model-override` (`e3g`)

The two changes compose at the final input object. `e3g` may supply `model`, `reasoning_effort`, and `service_tier`; this entity's guard inspects only `fork_turns` and never rewrites the object. `e3g` must not add a fork-mode surface, make isolation conditional on an override, or restore omission/`"all"`/numeric values. Conversely, this work must not select, validate, or default model/effort values. An allowed exact-`"none"` call carrying all `e3g` fields unchanged is the shared integration contract.

## Test plan

- Add table tests for omitted, `"all"`, numeric strings, JSON numbers, non-string values, arbitrary strings, and exact `"none"`. Assert only exact `"none"` returns allow; all other valid envelopes return deny. Add malformed JSON and non-object-input cases that exercise exit-code-2 blocking. Mutating any unsafe case to allow must fail.
- For exact `"none"`, assert the command emits no `updatedInput` and does not inspect or reconstruct unrelated fields. Repeat with `model`, `reasoning_effort`, `service_tier`, and unknown keys. For unsafe values, prove the decision remains deny with the same unrelated fields present.
- Add a narrow hook/configuration smoke test proving the plugin invokes the public binary command for exact canonical `agentsspawn_agent`. Add a process-level fixture proving allow/deny JSON and exit behavior. Treat these only as routing/mechanism evidence, not AC-1 evidence.
- Run the shortest disposable integrated captain probe below. Retain only its review artifact and JSONL as release evidence; do not productize the probe. The enabled run must join generated artifact, exact arguments, guard decisions, absence of children for denied calls, allowed-child canary absence, and same-handle follow-up. With the guard disabled, an unsafe call must execute and expose the parent canary.
- On the release-candidate SHA run focused tests, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; then perform the existing `main` tag and `next` propagation checks for AC-5.

## Shortest captain re-probe

Launch one disposable trusted-hook session matching canonical `agentsspawn_agent` and seed `PARENT_ONLY_CANARY_<random>`. Make the FO consume generated dispatch artifacts for omission, `"all"`, numeric, and exact `"none"`, then call `followup_task` on the allowed child's handle. Pass only if one record joins each artifact to its exact FO call and guard decision; unsafe calls are denied with no child; exact `"none"` is allowed unchanged; the child reports the parent canary absent; and that same child recalls its first-turn marker. Repeat the unsafe case once with the guard disabled; it must execute and the child must see the canary. Any missing join fails AC-1.

## Documentation delta

No public workflow or release documentation changes: fresh isolation is already the promised behavior, and no user-selectable setting is added. Update only `skills/first-officer/references/codex-first-officer-runtime.md` to say that the FO requests exact `fork_turns: "none"`, the plugin boundary guard rejects every other or missing value immediately before execution without rewriting other arguments, and continuity uses `followup_task` exclusively. Do not describe the guard as shipped until AC-1 passes on the release candidate.

## Stage Report: ideation

- DONE: Identify/exercise smallest actual live FO-to-Codex enforcement seam; use existing logs or short captain probe; stop on proven upstream blocker.
  Archive evidence established the escaped live FO call and the separate successful raw-host isolation/continuity probe. A disposable probe then exercised the live `agentsspawn_agent` `PreToolUse` seam: it observed the call and emitted a correct rewrite, but Codex CLI 0.144.6 executed inherited context. The entity stops at the named upstream rewrite blocker and the scratch probe was removed.
- DONE: Concrete behavior-first design, ACs, test plan preserving integrated artifact -> exact FO spawn args -> child-visible context, continuity followup_task only.
  AC-1 retains the full artifact-to-arguments-to-child join and requires a red inherited-context baseline. AC-2 protects unrelated fields, AC-3 proves same-handle `followup_task` continuity, and the test plan separates unit/routing support from the required disposable integrated live proof.
- DONE: Coordinate with per-host-stage-model-override e3g, name every mechanism’s value and rejected cheaper alt, exact minimal docs, no model/effort or fork-mode surface.
  The conditional normalizer shallow-copies all input fields and overwrites only `fork_turns`, so `e3g` owns model/effort values without weakening isolation. The design rejects adapter/prose substitution, private scripts, field whitelists, permanent harnesses, host defaults, and any fork selector; only the Codex runtime binding note changes after live proof passes.

### Summary

The ideation found the correct narrow boundary but proved it is not presently enforceable: Codex CLI 0.144.6 observes the namespaced collaboration spawn in `PreToolUse` yet does not apply `updatedInput`. Implementation is therefore blocked on upstream live-tool rewrite support; once that support exists, the specified shallow-copy normalizer and integrated captain probe are the minimum path to AC-1 without colliding with `e3g`.

### Feedback Cycles

- Cycle 1: REJECTED by the captain at the ideation gate on 2026-07-20. Do not treat ignored input rewriting as the terminal seam result. Exercise the existing fail-closed `PreToolUse` decision on the observed namespaced spawn: reject missing, `"all"`, and numeric `fork_turns`; allow exact `"none"`; verify child context and same-handle `followup_task` with existing logs or one disposable captain probe. Add no permanent test infrastructure. If denial is also ignored, retain the upstream blocker with that stronger evidence.

## Stage Report: ideation (cycle 2)

- DONE: Test whether the existing PreToolUse decision can DENY the observed namespaced collaboration spawn when fork_turns is missing, `all`, or numeric, while allowing exact `none`.
  Disposable Codex CLI 0.144.6 session `019f7dd1-f012-7e32-9b14-c1e9278390c5` recorded deny for omission, `"all"`, and numeric `"1"`, with no child created, and allow for exact `"none"` on canonical `agentsspawn_agent`.
- DONE: If denial works, redesign around a tiny guard that blocks unsafe calls rather than rewriting input; verify an allowed `none` child lacks the parent canary and `followup_task` preserves the same child.
  Child session `019f7dd2-590e-7cd0-a6ea-f9729b8aefaf` reported the parent canary absent and marker `CHILD-MARKER-7KQ2-N9V4`; same-handle `followup_task` recovered that marker. The design now emits allow only for exact `"none"` and deny/exit-2 for every unsafe or malformed case, with no rewrite.
- DONE: Use existing logs or one disposable manual probe; create no permanent test infrastructure.
  One manual probe supplied the missing fail-closed evidence; both disposable probe directories were removed, leaving no repo harness, parser, fixture, recorder, or script.
- DONE: Update ACs/design/test instructions and append a new ideation stage report for this cycle, commit/push state only, then send completion.
  AC-1 now joins artifact, exact call, guard decision, and child context; AC-2 protects `e3g` fields without mutation; the implementation, tests, integrated probe, and minimal documentation delta all use the proven deny-only seam.

### Summary

Cycle 2 overturns the earlier blocker: input rewriting is ignored on the namespaced collaboration path, but `permissionDecision: "deny"` reliably blocks unsafe spawns. The smallest viable design is therefore a public binary-backed, fail-closed guard that allows only exact `fork_turns: "none"`, leaves all other arguments untouched, and preserves continuity solely through `followup_task`.
