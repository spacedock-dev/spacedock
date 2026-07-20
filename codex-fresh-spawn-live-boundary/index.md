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

**AC-1 (VALUE) — An instruction-driven Codex First Officer cannot execute a fresh Spacedock worker spawn that inherits parent turns.** Verify this with one integrated live record joining the helper-generated dispatch artifact, its exact `"none"` FO spawn, the boundary decision, and child-visible context, plus controlled direct FO attempts with omitted, `"all"`, and numeric `fork_turns`. Unsafe attempts must be denied before child creation; exact `"none"` must be allowed and its child must lack the parent canary. Compare with the archived ineffective-guard baseline where `"all"` executed and exposed the canary, so the value can move red.

**AC-2 — The guard decides only from the isolation field and never rewrites tool input.** Exact `"none"` is allowed with the FO-issued argument object unchanged. Missing, `"all"`, numeric, non-string, or malformed input is denied. Repeat allowed and denied cases with `model`, `reasoning_effort`, `service_tier`, and arbitrary unknown fields; those fields must not affect the decision and must reach the child unchanged on an allowed call.

**AC-3 — Deliberate continuity remains exclusively `followup_task` on the existing worker handle.** Verify this in the AC-1 live journey: the fresh child lacks the parent canary, then `followup_task` reaches that exact child and the child recalls its own first-turn marker. A second spawn is not continuity.

**AC-4 — Adapter maps, contract wording, hook structure, and raw host probes remain supporting evidence.** They cannot substitute for the AC-1 join. The supported launch must load and trust the guard; malformed input or an unavailable JSON interpreter must exit `2` and block; release validation must positively observe the guard decision. A missing, disabled, untrusted, or otherwise failed hook invalidates the release claim instead of narrowing the acceptance boundary.

**AC-5 — v0.25.2 ships the fix on the stable line without rewinding `next`.** Verify the exact release candidate SHA with required Go/full/race gates and AC-1 through AC-3 live evidence, cut annotated `v0.25.2` from `main`, and retain the invariant on `next` through the documented propagation path.

## Scope guard

Do not build a new general agent harness or fork-mode configuration surface. Reuse the narrowest existing live capture path that can observe the real FO tool call. Do not expose `fork_turns` to workflow authors. Do not change stage reuse policy. Do not fold unrelated model/effort routing into this patch.

## Riskiest-mechanism spike and decision

The archive already contains the two halves of the problem, but not the required join. Top-level session `codex:019f7d9a-5b06-75a0-a04a-02b0b2ccd6a2` records the escaped FO call with `fork_turns: "all"` and its retry with `"none"`; the first call failed only on an incompatible agent type. Session `codex:019f79e5-3974-75e2-ab25-d6d07836cc72`, corroborated by parent session `codex:019f7007-8fba-7503-8c44-5ebf9a7cc945`, records a raw-host isolated child and successful same-handle follow-up. No archived run joins generated dispatch bytes, the exact FO-issued arguments, absence of a parent canary, and same-child continuity.

A first disposable captain probe exercised input rewriting in Codex CLI 0.144.6. A trusted `PreToolUse` hook observed canonical tool name `agentsspawn_agent` and returned documented `updatedInput`, but the child still inherited the parent canary. That proves rewriting is not the viable seam in this runtime.

A second disposable probe exercised the cheaper fail-closed decision on that same canonical tool path. The hook denied omission, `"all"`, and numeric `"1"`; all three calls returned `Tool call blocked by PreToolUse hook` and created no child. It allowed exact `"none"`; child `/root/allow_none` reported the parent canary absent and marker `CHILD-MARKER-7KQ2-N9V4`, then `followup_task` on that same handle recovered the marker. The hook log recorded the four exact inputs and decisions. Evidence is in parent session `codex:019f7dd1-f012-7e32-9b14-c1e9278390c5` and child session `codex:019f7dd2-590e-7cd0-a6ea-f9729b8aefaf`.

**Decision: PROCEED with a deny-only guard.** Namespaced input rewriting remains unsupported, but it is unnecessary: denial on the same live seam is effective. The guard makes unsafe calls impossible to execute and leaves safe calls untouched. Adapter hardening and FO prose remain defense in depth. Both disposable probe directories were removed; this ideation adds no permanent harness, recorder, fixture, parser, or test infrastructure.

## Implementation design

1. Extend the existing root `hooks.json` with one `PreToolUse` matcher group for the live-verified canonical `^agentsspawn_agent$`. Like the shipped PostCompact hook, its single command is a plugin-root-absolute executable: `${PLUGIN_ROOT}/hooks/codex_fresh_spawn_guard.sh`.
2. Add that one POSIX shell hook. Before reading stdin, it checks `python3` availability and exits `2` with a concise blocking reason when unavailable. It then `exec`s one embedded Python-stdlib JSON predicate over the original stdin: require an object envelope and object `tool_input`; emit documented `permissionDecision: "allow"` with no `updatedInput` only when `fork_turns` is the exact string `"none"`; emit `permissionDecision: "deny"` for missing or every other value; exit `2` for malformed JSON or invalid envelope shape.
3. Do not change the helper, Codex adapter, CLI routing, or any Go package. They already request exact `fork_turns: "none"`; the escaped defect was a direct FO call. The guard is unconditional, never repairs or retries a rejected call, and never reads model/effort fields.
4. Keep worker continuity unchanged: only `followup_task` may address an existing handle. No workflow field, CLI flag, helper option, or host-specific fork choice is introduced.

This is the smallest behavior-owning design because it adds one matcher and one script on the live seam already proved. POSIX `grep`/`sed` matching is fewer lines but cannot safely distinguish JSON keys, types, escapes, duplicates, and nesting. `jq` adds an undeclared binary dependency. A standalone Python shebang cannot turn a missing interpreter into the documented blocking exit, while the shell wrapper can. A public `spacedock` subcommand, generalized policy engine, process fixture, or configuration-test framework adds a subsystem for a one-field predicate. Rewriting input is both larger and proven ineffective. The deny-only hook serves AC-1 directly and leaves `e3g` fields untouched for AC-2.

## Gross changed-LOC budget before implementation

| File | Gross changed LOC | Essential line categories |
| --- | ---: | --- |
| `hooks.json` | ~11 | One event group, exact canonical matcher, one plugin-root-absolute command. |
| `hooks/codex_fresh_spawn_guard.sh` | ~30 | Shebang/comments; Python availability check with exit `2`; JSON envelope/type validation; exact-`"none"` predicate; allow/deny serialization. |
| `skills/first-officer/references/codex-first-officer-runtime.md` | ~2 | Replace one runtime-binding line with the exact guard promise below. |
| **Total** | **~43** | No helper, adapter, CLI, Go package, fixture, harness, or public docs changes. |

The executable mode bit for the new hook is required but is not LOC. If implementation materially exceeds this budget, stop and return to design rather than growing a policy subsystem.

## Coordination with `per-host-stage-model-override` (`e3g`)

The two changes compose at the final input object. `e3g` may supply `model`, `reasoning_effort`, and `service_tier`; this entity's guard inspects only `fork_turns` and never rewrites the object. `e3g` must not add a fork-mode surface, make isolation conditional on an override, or restore omission/`"all"`/numeric values. Conversely, this work must not select, validate, or default model/effort values. An allowed exact-`"none"` call carrying all `e3g` fields unchanged is the shared integration contract.

## Test plan

- Before implementation completion, directly pipe a disposable stdin matrix into the shipped hook from an unrelated cwd; do not commit the driver. Exact `"none"` must exit `0` with allow and no `updatedInput`. Missing, `"all"`, numeric strings, JSON numbers, null, and arbitrary strings must exit `0` with deny. Malformed JSON, non-object envelopes, non-object `tool_input`, and a `PATH` without `python3` must exit `2`. Repeat safe and unsafe inputs carrying `model`, `reasoning_effort`, `service_tier`, and unknown keys; decisions must depend only on `fork_turns`.
- Reuse cycle 2's live probe as the mechanism spike: session `019f7dd1-f012-7e32-9b14-c1e9278390c5` already proves the canonical matcher denies omission/`"all"`/numeric and allows exact `"none"`; child session `019f7dd2-590e-7cd0-a6ea-f9729b8aefaf` proves canary absence and same-handle continuity. Do not add a committed harness, process fixture, generalized policy engine, or configuration smoke-test framework.
- Run the manual release journey below against the installed candidate plugin. This is the only new live proof and must join helper artifact, exact FO calls, guard decisions, child context, and same-handle continuity. Compare it with the archived ineffective-guard `"all"` baseline rather than disabling the candidate guard.
- On the release-candidate SHA run focused tests, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; then perform the existing `main` tag and `next` propagation checks for AC-5.

## Shortest captain re-probe

Install the release candidate through the supported Codex plugin launch and trust the bundled hook; `/hooks` must show it active. Seed `PARENT_ONLY_CANARY_<random>`. First have the FO consume a normal helper-generated dispatch artifact: its exact spawn must contain `fork_turns: "none"`, the guard must allow it unchanged, and the child must report the canary absent plus a unique marker. Use `followup_task` on that handle and require the marker. Then direct the same FO to attempt omission, `"all"`, and numeric `"1"` exactly once each; the JSONL must show three hook-blocked calls and no child handles. Retain that one JSONL and artifact as release evidence and compare its unsafe outcomes with the archived ineffective-guard baseline where `"all"` executed and exposed the canary. Any missing join fails AC-1.

## Documentation delta

No public workflow or release documentation changes: fresh isolation is already the promised behavior, and no user-selectable setting is added. In `skills/first-officer/references/codex-first-officer-runtime.md`, replace the current `«worker.spawn»` sentence ending `Every spawn is a fresh dispatch; deliberate continuity uses followup_task with the existing handle.` with: `Every spawn is a fresh dispatch. The bundled PreToolUse guard allows worker creation only when the live input carries exact fork_turns="none"; missing or any other value is denied without rewriting model, reasoning_effort, service_tier, or unknown fields. Deliberate continuity uses followup_task with the existing handle.` Keep the surrounding helper-message and task-name instructions unchanged. Do not describe the guard as shipped until AC-1 passes on the release candidate.

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
- Cycle 2: REJECTED by the captain at the ideation gate on 2026-07-20 as overbuilt. The proposed public `spacedock dispatch guard-codex-spawn` command, process fixture, configuration smoke test, and 100+ LOC estimate turn a one-field boundary predicate into a subsystem. Return to ideation and produce the smallest change on the already-proven `PreToolUse` hook path. Reuse existing live output as evidence and provide manual release-test instructions; add no harness or generalized policy surface. Give a gross changed-LOC estimate by file before implementation and stop if the minimal hook cannot parse and deny safely.

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

## Stage Report: ideation (cycle 3)

- DONE: Replace the public-command subsystem design with the smallest safe enforcement on the already-proven PreToolUse boundary.
  The design now adds only one `hooks.json` matcher and one plugin-root shell hook with an embedded Python-stdlib predicate; it changes no helper, adapter, CLI, Go package, or fork-mode surface.
- DONE: Give a gross changed-LOC estimate by file before implementation; identify every line category that is essential.
  The pre-implementation budget is ~43 gross LOC across `hooks.json` (~11), `hooks/codex_fresh_spawn_guard.sh` (~30), and one runtime-reference line replacement (~2), with each parsing, fail-closed, decision, binding, and documentation category named.
- DONE: Reuse existing live probe output and provide manual release-test instructions; add no permanent harness, generalized policy engine, process fixture, or configuration smoke-test framework.
  Cycle 2's parent/child sessions remain the mechanism proof; the release plan adds one disposable installed-plugin journey joining the normal helper artifact, exact safe/unsafe calls, guard decisions, child canary absence, and same-handle marker recall.

### Summary

The overbuilt public-command design has been replaced with the repository's existing direct plugin-hook pattern. A single fail-closed shell hook safely delegates JSON parsing to Python's standard library, blocks if Python or the envelope is unavailable, allows only exact `fork_turns: "none"`, and stays within a ~43-gross-LOC implementation budget with no permanent test infrastructure.
