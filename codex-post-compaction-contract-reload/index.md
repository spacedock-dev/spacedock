---
title: Post-compaction contract reload
status: ideation
source: Absorbed from task njr36mfyhbafy8zx9ydks8ep in another workflow; canonical handoff /tmp/first-officer-compaction-rehydration.md; captain directed repo-local absorption 2026-07-11
started: 2026-07-11T04:15:29Z
completed:
verdict:
score: 0.95
worktree:
issue:
id: c60nzb396vgf0f8a9v0sggwm
milestone: 0.26.0
---

## Decision

Ship two advisory hints and no continuation controller:

1. When context pressure is apparent, the first officer may suggest compaction only after current workflow state is already durable at a clean boundary.
2. After Codex compacts, remind the captain to tell the first officer to reread the authoritative `spacedock:first-officer` contract and reconcile durable workflow and live worker state before continuing.

These hints do not authorize actions, block actions, reconstruct a session, or make summaries authoritative. The existing workflow files, committed Stage Reports, state-checkout history, live worker roster, and first-officer contract remain the sources the FO normally reconciles. The superseded ledger, authorization grant, permits, action gateway, crash-replay controller, interception matrix, and watchdog design is preserved at [artifacts/superseded-controller-design.md](artifacts/superseded-controller-design.md), not repaired.

## Problem

Codex compaction can preserve enough narrative to continue while dropping an operating detail such as the current wait contract or the need to consume a worker's durable report. The useful intervention is a timely reminder at each side of the boundary, not a second workflow state machine.

A pre-compaction suggestion is safe only when the current boundary can be recovered without relying on conversational memory. For this task, that means:

- workflow/entity/report changes already made by the FO are committed in the state checkout;
- assigned code or report work already claimed complete is committed in its owning worktree or state path;
- no received completion, gate decision, state transition, archive, merge, or other FO-owned effect is half-applied or awaiting reconciliation; and
- any unresolved workers can be rediscovered from the live roster and their durable entity stage/cycle, rather than only from prose in the conversation.

If any term is false, the FO finishes that durability or reconciliation work first and does not recommend compaction yet. This is a judgment rule in the FO contract, not a new persisted `compaction_ready` field.

After compaction, the generated summary is useful history but not the contract. The FO must reread the active plugin's authoritative first-officer skill and its required eager imports, then use the existing boot/status/roster/report reconciliation flow before taking another workflow action. A captain cue such as “we compacted” or “reread the FO contract” is sufficient to trigger this behavior even when no host hook is available.

## Host spike: Codex 0.144.4

The riskiest mechanism was exercised first in a live Codex TUI on 2026-07-17. The exact probe and observations are recorded in [artifacts/codex-0.144.4-hook-probe.md](artifacts/codex-0.144.4-hook-probe.md).

- Manual `/compact` ran configured `PreCompact(manual)` and `PostCompact(manual)` command hooks.
- A `PostCompact` hook returning JSON `systemMessage` produced one visible `PostCompact hook (completed) warning` after `Context compacted`.
- The next model turn, asked without tools to repeat the warning or answer `NONE`, answered `NONE`. On this client, `systemMessage` is a captain/UI reminder, not developer context for the FO.
- A configured `SessionStart` hook matched only on `compact` did not fire during the same-session manual `/compact` probe. The design therefore does not rely on `SessionStart(source=compact)`.

This matches the current Codex hook documentation: `systemMessage` is surfaced as a UI/event-stream warning, while `SessionStart` is the documented hook whose `additionalContext` becomes developer context. Hook definitions also require trust and can be disabled. Therefore no supported live hook in this probe automatically instructs the post-compact model. The safe fallback is a visible, non-blocking captain reminder plus the FO's manual-cue rule.

## Proposed approach

### Hint 1: safe-to-compact

Add a short Codex first-officer rule:

> When context pressure is apparent, suggest compaction only after the current workflow boundary is durable and recoverable from committed workflow/report state plus the live roster. If a completion, gate, state mutation, archive, or merge still needs reconciliation, finish it first. At a safe boundary, tell the captain: “Context is getting tight. Current Spacedock state is durable at a clean boundary; now is a safe time to compact.” The suggestion is optional and non-blocking.

The FO uses an available host pressure indication or captain cue. This task does not invent token thresholds or inspect transcripts. The separate context-budget work may later provide a better signal without changing this safety rule.

The simplest alternative was to suggest compaction whenever the host reports pressure. It is insufficient because timing, not detection, is the value: the same reminder is harmful between a worker completion and report verification or during an uncommitted state transition.

### Hint 2: post-compact reload

Bundle one Codex `PostCompact` command hook matching `manual|auto`. It emits exactly this `systemMessage` and performs no writes:

> Spacedock: compaction completed. Before continuing, ask the first officer to reread the authoritative `spacedock:first-officer` contract and reconcile durable workflow and live worker state.

Add the corresponding FO rule:

> When the captain says compaction occurred or asks for a contract reload, reread the active `spacedock:first-officer` `SKILL.md` completely and its required eager imports. Then run the normal workflow status and live-roster reconciliation, verify any newer durable Stage Reports, and only then continue. Do not treat the compacted summary as authority.

The hook is deliberately a reminder. If hooks are unsupported, disabled, untrusted, or fail, Codex continues normally and the captain can give the same cue manually. Do not create a marker file, checkpoint, monitor, permit, background process, automatic follow-up turn, stop block, or action fence.

The simplest alternative was a `SessionStart(compact)` developer-context injection. The live TUI did not deliver that event for same-session manual compaction, so it cannot be the v1 dependency. A `PostCompact` warning is weaker but real and harmless.

## Acceptance criteria

- **AC-1 — safe timing:** Given the same context-pressure cue in two live or fixture-backed FO scenarios, a fully durable clean-boundary scenario produces exactly one non-blocking safe-to-compact suggestion, while a scenario with an uncommitted state/report change or an unconsumed worker completion produces zero suggestions until that work is durably reconciled. Evidence includes the relevant state/worktree commit OIDs before the suggestion. This measures the end value: suggestions occur at recoverable boundaries, not merely whenever pressure is mentioned.
- **AC-2 — visible host reminder:** On the target Codex client with the bundled hook trusted, one manual `/compact` produces exactly one visible warning containing the required reread-and-reconcile instruction after compaction completes. The hook configuration matches both `manual` and `auto`; a command-level fixture drives both event payloads and asserts one valid JSON response per event. The test does not claim the warning enters model context.
- **AC-3 — reload before action:** In a split-root FO replay, after a captain compaction cue, the next workflow effect occurs only after observable reads of the active first-officer contract and required eager imports, a fresh workflow status query, live roster reconciliation when available, and verification of any newer committed Stage Report. A stale summary that says “continue directly” must not skip those observations. Proof uses captured reads/tool calls plus workflow and state-checkout OIDs, not an assertion that response prose mentions reloading.
- **AC-4 — harmless absence:** With hooks disabled, untrusted, unavailable, or returning non-zero, Codex compaction and the next captain turn continue without a Spacedock-created state file, background process, blocked stop, automatic turn, or workflow mutation. A manual captain cue still exercises AC-3.

## Test plan

1. Add a small hook fixture test that parses the shipped hook configuration, drives `manual` and `auto`, validates the exact JSON `systemMessage`, and asserts the command performs no filesystem writes. Run an absent/disabled/failing-handler matrix and assert normal exit/continuation behavior. Cost: small; serves AC-2 and AC-4.
2. Add a Codex first-officer integration fixture with paired safe and unsafe pressure cases. Record state/worktree commit OIDs, pending completion/gate state, and emitted captain messages; assert one suggestion only after the unsafe case is reconciled and committed. Do not accept a grep for the instruction text as proof. Cost: medium; serves AC-1.
3. Add a split-root post-compaction replay. Supply a misleading compacted summary, then a captain cue. Capture the active skill/eager-import reads, `spacedock status`, `list_agents` when present, report OID checks, and the first later workflow mutation. Assert ordering and clean Git state. Cost: medium; serves AC-3 and AC-4.
4. Keep an opt-in live Codex TUI probe for manual `/compact`. Trust the test hook explicitly, assert the warning appears once after `Context compacted`, and ask the next model turn whether the warning was in its context. Until Codex behavior changes, the expected answer is `NONE`; this guards against accidentally upgrading a UI reminder into an unsupported automatic-reload claim. Cost: medium/live; serves AC-2 and the host boundary.

## Documentation change

Implementation updates the Codex first-officer runtime reference and `docs/runtime-support.md` with this concrete addition:

```diff
+ When Codex context pressure is apparent, the first officer recommends compaction
+ only after current workflow and report state is committed at a clean boundary.
+
+ When plugin hooks are enabled and trusted, Codex shows a Spacedock warning after
+ compaction asking the captain to have the first officer reload the authoritative
+ first-officer contract and reconcile workflow status, durable reports, and the live
+ worker roster. The warning is advisory and is not injected into the model context on
+ the currently validated Codex client. If the hook is unavailable, give the same cue
+ manually; normal operation remains unblocked.
```

## Out of scope

Durable authorization or action ledgers; permits; checkpoints; continuation controllers; crash replay; interception or enforcement matrices; stop blocking; automatic turns; watchdogs; summary generation; transcript parsing; token-threshold design; Claude or Pi changes; and recovering effects that were not already made durable by existing Spacedock workflow rules.

## Feedback Cycles

**Cycle 1 (captain feedback, 2026-07-14).** A compacted FO session lost wait/reconciliation discipline and acted on unrelated ready work. The first revision responded with a typed authorization grant and durable continuation machinery.

**Cycle 2 (independent ideation review).** Review found gaps in that controller's action vocabulary, gate binding, interception proof, and report identity. These findings apply only to the superseded controller design.

**Cycle 3 (captain scope reset, 2026-07-17).** The controller design was rejected as unnecessary. The intended product is exactly two hints: recommend compaction only at an already-durable boundary, then remind the post-compaction FO to reread its authoritative contract and reconcile durable state. This body replaces the rejected design rather than repairing it.

## Stage Report: ideation (cycle 3)

- DONE: Replace the continuation controller with exactly two bounded hints.
  The current body contains one safe-boundary compaction suggestion and one post-compaction contract-reload reminder. It removes the authorization ledger, permits, checkpoint, gateway, crash replay, interception matrix, monitor, watchdog, and their derived acceptance criteria.
- DONE: Exercise and honestly bound the current Codex lifecycle surface.
  Live Codex 0.144.4 manual `/compact` fired `PreCompact` and `PostCompact`; the shipped-shape `PostCompact` `systemMessage` appeared as a warning but was absent from the next model's context. `SessionStart(compact)` did not fire. The proposed hook is therefore captain-facing and failure-open, with a manual cue fallback.
- DONE: Rewrite the problem, proposed approach, acceptance criteria, test plan, and documentation change around observable hint timing and reload behavior.
  AC-1 pairs safe and unsafe pressure scenarios, AC-2 proves the visible warning without claiming model injection, AC-3 proves actual contract/status/roster/report reads before the next effect, and AC-4 proves harmless absence. The full rejected design remains beside the entity as an artifact.

### Summary

c6 is now a small continuity aid: the FO suggests compaction only when existing state is recoverable, and Codex warns the captain after compaction to trigger a real contract reload and state reconciliation. The live host does not inject that warning into model context, so the design says so and falls back to a manual captain cue instead of adding lifecycle machinery.
