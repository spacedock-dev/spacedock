---
title: Post-compaction contract reload
status: validation
source: Absorbed from task njr36mfyhbafy8zx9ydks8ep in another workflow; canonical handoff /tmp/first-officer-compaction-rehydration.md; captain directed repo-local absorption 2026-07-11
started: 2026-07-11T04:15:29Z
completed:
verdict:
score: 0.95
worktree: .worktrees/spacedock-ensign-codex-post-compaction-contract-reload
issue:
id: c60nzb396vgf0f8a9v0sggwm
milestone: 0.26.0
---

## Decision

Ship two host-neutral advisory rules and no continuation controller:

1. When context pressure is apparent, the first officer may suggest compaction only after current workflow state is already durable at a clean boundary.
2. After the host compacts, the first officer rereads the authoritative `spacedock:first-officer` contract and reconciles durable workflow and live worker state before the next workflow effect.

Both rules are first-officer judgment rules at the shared-core level, inherited by every runtime adapter; neither depends on a host mechanism. What IS host-specific is DELIVERY of the post-compaction reminder — abstracted as a per-host binding `«post-compact-notice»` that mirrors how the contract already treats `«context-budget»` (PRESENT/ABSENT per host). The binding is failure-open: when it cannot reach the model, a manual captain cue triggers Rule 2 and nothing is blocked.

These rules do not authorize actions, block actions, reconstruct a session, or make summaries authoritative. The existing workflow files, committed Stage Reports, state-checkout history, live worker roster, and first-officer contract remain the sources the FO normally reconciles. The superseded ledger, authorization grant, permits, action gateway, crash-replay controller, interception matrix, and watchdog design is preserved at [artifacts/superseded-controller-design.md](artifacts/superseded-controller-design.md), not repaired.

## Problem

Host compaction can preserve enough narrative to continue while dropping an operating detail such as the current wait contract or the need to consume a worker's durable report. The useful intervention is a timely reminder at each side of the boundary, not a second workflow state machine. This is true on every host, not only Codex.

A pre-compaction suggestion is safe only when the current boundary can be recovered without relying on conversational memory. For this task, that means:

- workflow/entity/report changes already made by the FO are committed in the state checkout;
- assigned code or report work already claimed complete is committed in its owning worktree or state path;
- no received completion, gate decision, state transition, archive, merge, or other FO-owned effect is half-applied or awaiting reconciliation; and
- any unresolved workers can be rediscovered from the live roster and their durable entity stage/cycle, rather than only from prose in the conversation.

If any term is false, the FO finishes that durability or reconciliation work first and does not recommend compaction yet. This is a judgment rule in the FO contract, not a new persisted `compaction_ready` field.

After compaction, the generated summary is useful history but not the contract. The FO must reread the active plugin's authoritative first-officer contract and use the existing boot/status/roster/report reconciliation flow before taking another workflow action. The **authoritative contract to reread is identical on every host**: the `spacedock:first-officer` `SKILL.md`, its eager import `@references/first-officer-shared-core.md`, and the active host runtime adapter (`claude-`/`codex-`/`pi-first-officer-runtime.md`). Rereading those three re-triggers the normal boot, which reloads the deferred dispatch/write/merge modules on demand. A captain cue such as "we compacted" or "reread the FO contract" is sufficient to trigger this behavior even when no host delivery binding is available.

## Proposed approach

### Rule 1: safe-to-compact (host-neutral FO judgment rule)

Add a short first-officer rule at the shared-core level:

> When context pressure is apparent, suggest compaction only after the current workflow boundary is durable and recoverable from committed workflow/report state plus the live roster. If a completion, gate, state mutation, archive, or merge still needs reconciliation, finish it first. At a safe boundary, tell the captain: "Context is getting tight. Current Spacedock state is durable at a clean boundary; now is a safe time to compact." The suggestion is optional and non-blocking.

The **trigger is equally weak on all hosts today**: none exposes the FO's own session pressure programmatically (the Claude/Codex context-budget probes measure a dispatched *worker's* window, not the FO session). So this rule leans on a captain cue or apparent host pressure on every host. This task does not invent token thresholds or inspect transcripts. The separate context-budget work may later provide a better signal without changing this safety rule.

The simplest alternative was to suggest compaction whenever the host reports pressure. It is insufficient because timing, not detection, is the value: the same reminder is harmful between a worker completion and report verification or during an uncommitted state transition. (Serves AC-1.)

### Rule 2: post-compact reload (host-neutral FO judgment rule)

Add the corresponding FO rule at the shared-core level:

> When a compaction occurred (signalled by `«post-compact-notice»` where PRESENT, else a captain cue), reread the authoritative `spacedock:first-officer` contract completely — `SKILL.md`, its eager import `first-officer-shared-core.md`, and the active host runtime adapter. Then run the normal workflow status and live-roster reconciliation, verify any newer durable Stage Reports, and only then continue. Do not treat the compacted summary as authority.

The simplest alternative was to keep this Codex-specific. It is insufficient because the rule names no host mechanism — it is FO judgment over already-proven reads (committed state, `status`, the live roster, durable reports) — so scoping it to one host would leave the other hosts with the same failure and no rule. (Serves AC-2.)

### `«post-compact-notice»`: per-host delivery binding

DELIVERY of the Rule 2 reminder is abstracted as a per-host binding that mirrors `«context-budget»`'s PRESENT/ABSENT discipline, not a host-specific procedure baked into the rule:

> **`«post-compact-notice»`: deliver the post-compaction reload reminder**
>
> - **block:** failure-open. The binding automates the Rule 2 trigger; when it is ABSENT, UI-only, disabled, untrusted, or fails, a manual captain cue ("we compacted" / "reload the FO contract") triggers Rule 2 and nothing is blocked. No Spacedock state file, marker, background process, stop block, or automatic turn is created on any host.
> - → **Codex:** PRESENT as a captain-facing UI warning only — a bundled `PostCompact` command hook (`manual|auto`) emitting a `systemMessage`. The live 0.144.4 probe proved this warning does NOT reach model context, so for the model it degrades to the manual-cue fallback while the visible warning prompts the captain to give that cue. Proof: [artifacts/codex-0.144.4-hook-probe.md](artifacts/codex-0.144.4-hook-probe.md).
> - → **Claude:** model-context delivery is possible via `SessionStart(compact)` (Claude docs' "re-inject context after compaction"), which would make Rule 2 fire automatically rather than via a captain cue. This binding is SPLIT to task `claude-post-compaction-contract-reload` (id `cdbhzxc`) and is NOT designed here; until it ships, Claude uses the manual-cue fallback.
> - → **Pi:** ABSENT — no known compaction-hook surface; manual-cue only.

The Codex binding's `systemMessage` is exactly this reminder and the hook performs no writes:

> Spacedock: compaction completed. Before continuing, ask the first officer to reread the authoritative `spacedock:first-officer` contract and reconcile durable workflow and live worker state.

The simplest alternative was to design delivery inline per host without a named binding. It is insufficient because it would couple the host-neutral Rule 2 to a host mechanism, duplicate the failure-open fallback three times, and forfeit the proven `«context-budget»` PRESENT/ABSENT template that the contract, `fo-dispatch-core.md`, and `docs/runtime-support.md` already read as one shape. (Serves AC-3 for the Codex leg; the Claude leg's value AC lives in `cdbhzxc`.)

### Shared prerequisite: hook-shipping surface

Both the Codex binding (this task) and the Claude binding (`cdbhzxc`) need a plugin hook-shipping surface the plugin does not have today: neither `.claude-plugin/plugin.json` nor `.codex-plugin/plugin.json` declares a `hooks` key (verified 2026-07-18 — a `grep -rl hooks` over both manifests returns nothing). This task's Codex binding establishes the `.codex-plugin` hooks surface for its `PostCompact` hook; the split task owns the `.claude-plugin` surface for `SessionStart(compact)`. The prerequisite's shape is shared even though each host ships its own hook.

### Spike status

- **Codex binding delivery behavior — spiked and proven.** The riskiest delivery claim (does the post-compaction reminder reach the model?) was exercised first in a live Codex 0.144.4 TUI on 2026-07-17: manual `/compact` fired `PreCompact`/`PostCompact`; a `PostCompact` `systemMessage` appeared as one visible warning but the next model turn answered `NONE`, and `SessionStart(compact)` did not fire same-session. On this client, `systemMessage` is a captain/UI reminder, not developer context. Full probe: [artifacts/codex-0.144.4-hook-probe.md](artifacts/codex-0.144.4-hook-probe.md). The design therefore makes the Codex binding failure-open with a manual-cue fallback and does not claim automatic model instruction.
- **Codex binding remaining unverified path — implementation's first spike.** The probe used an inline, trusted hook (`--dangerously-bypass-hook-trust`), not a plugin-bundled one. Whether `.codex-plugin/plugin.json` can declare a `hooks` key that Codex loads and trusts is the Codex binding implementation's first spike, ahead of the fixture matrix.
- **Host-neutral Rules 1 and 2 — no spike needed.** They rely only on already-proven FO mechanisms: committed state reads, `spacedock status`, live-roster reconciliation, and durable Stage Report reads. No new mechanism underlies the judgment rules themselves.

## Acceptance criteria

- **AC-1 — safe timing (host-neutral value AC):** Given the same context-pressure cue in two live or fixture-backed FO scenarios, a fully durable clean-boundary scenario produces exactly one non-blocking safe-to-compact suggestion, while a scenario with an uncommitted state/report change or an unconsumed worker completion produces zero suggestions until that work is durably reconciled. Evidence includes the relevant state/worktree commit OIDs before the suggestion. This measures the end value against a baseline that can move the wrong way — a suggestion fired at an unrecoverable boundary — not merely that the rule text shipped.
- **AC-2 — reload before the next effect (host-neutral):** In a split-root FO replay, after a compaction cue (delivered by the host binding where PRESENT, else a manual captain cue), the next workflow effect occurs only after observable reads of the authoritative contract (`SKILL.md`, `first-officer-shared-core.md`, and the active host runtime adapter), a fresh workflow `status` query, live roster reconciliation where available, and verification of any newer committed Stage Report. A stale summary that says "continue directly" must not skip those observations. Proof uses captured reads/tool calls plus workflow and state-checkout OIDs, not an assertion that response prose mentions reloading.
- **AC-3 — Codex delivery binding, visible warning:** On the target Codex client with the bundled hook trusted, one manual `/compact` produces exactly one visible warning containing the required reread-and-reconcile instruction after compaction completes. The hook configuration matches both `manual` and `auto`; a command-level fixture drives both event payloads and asserts one valid JSON `systemMessage` per event. The test does not claim the warning enters model context. This AC proves only the Codex binding; the Claude auto-inject AC lives in `claude-post-compaction-contract-reload` (`cdbhzxc`).
- **AC-4 — harmless absence (host-neutral):** With the delivery binding disabled, untrusted, unavailable, or returning non-zero on any host, compaction and the next captain turn continue without a Spacedock-created state file, background process, blocked stop, automatic turn, or workflow mutation. A manual captain cue still exercises AC-2.

## Test plan

1. Add a small Codex hook fixture test that parses the shipped hook configuration, drives `manual` and `auto`, validates the exact JSON `systemMessage`, and asserts the command performs no filesystem writes. Run an absent/disabled/failing-handler matrix and assert normal exit/continuation behavior. Cost: small; serves AC-3 and AC-4.
2. Add a host-neutral first-officer integration fixture with paired safe and unsafe pressure cases exercising the shared-core Rule 1. Record state/worktree commit OIDs, pending completion/gate state, and emitted captain messages; assert one suggestion only after the unsafe case is reconciled and committed. Do not accept a grep for the instruction text as proof. Cost: medium; serves AC-1.
3. Add a split-root post-compaction reload replay. Supply a misleading compacted summary, then a compaction cue (manual — the host-neutral path). Capture the active `SKILL.md`/`first-officer-shared-core.md`/host-adapter reads, `spacedock status`, roster reconciliation when present, report OID checks, and the first later workflow mutation. Assert ordering and clean Git state. Cost: medium; serves AC-2 and AC-4.
4. Keep an opt-in live Codex TUI probe for manual `/compact`. Trust the test hook explicitly, assert the warning appears once after `Context compacted`, and ask the next model turn whether the warning was in its context. Until Codex behavior changes, the expected answer is `NONE`; this guards against accidentally upgrading a UI reminder into an unsupported automatic-reload claim. Cost: medium/live; serves AC-3 and the Codex host boundary. The Claude live auto-inject probe is NOT here — it belongs to `cdbhzxc`.

## Documentation change

Implementation updates the shared first-officer core, the dispatch-adapter capability list, and `docs/runtime-support.md` with this concrete host-neutral addition (representative unified diff):

```diff
# skills/first-officer/references/first-officer-shared-core.md — new "Compaction continuity" rule
+ **Compaction continuity (host-neutral).** When context pressure is apparent, suggest
+ compaction only after the current workflow boundary is durable and recoverable from
+ committed workflow/report state plus the live roster; if a completion, gate, state
+ mutation, archive, or merge still needs reconciliation, finish it first. After a
+ compaction (signalled by «post-compact-notice» where PRESENT, else a captain cue),
+ reread the authoritative contract — this SKILL.md, its eager import
+ first-officer-shared-core.md, and the active host runtime adapter — then run the normal
+ status + live-roster reconciliation and verify any newer committed Stage Report before
+ the next workflow effect. Do not treat the compacted summary as authority.

# skills/first-officer/references/fo-dispatch-core.md — Dispatch Adapter capability list
+ ## «post-compact-notice»: deliver the post-compaction reload reminder
+
+ - **block:** failure-open — an ABSENT/UI-only/disabled/failed binding falls back to a
+   manual captain cue; the reload rule still holds and nothing is blocked.
+ - → **Codex:** PRESENT (UI-only) — bundled `PostCompact` `systemMessage`; not model
+   context (probe). · **Claude:** `SessionStart(compact)` model-context injection — owned
+   by task cdbhzxc. · **Pi:** ABSENT — manual cue only.

# docs/runtime-support.md — ### Runtime binding-block shape list
+ - `«post-compact-notice»` -> post-compaction reload delivery binding, or manual-cue-only.
```

## Out of scope

Durable authorization or action ledgers; permits; checkpoints; continuation controllers; crash replay; interception or enforcement matrices; stop blocking; automatic turns; watchdogs; summary generation; transcript parsing; token-threshold design; and recovering effects that were not already made durable by existing Spacedock workflow rules.

Also out of scope for this entity: the Claude `SessionStart(compact)` model-context delivery binding and its live spike — owned by `claude-post-compaction-contract-reload` (id `cdbhzxc`); and any Pi post-compaction delivery beyond the manual captain cue. The host-neutral rules here cover all hosts; only the Claude auto-inject DELIVERY and its proof move to the split task.

## Feedback Cycles

**Cycle 1 (captain feedback, 2026-07-14).** A compacted FO session lost wait/reconciliation discipline and acted on unrelated ready work. The first revision responded with a typed authorization grant and durable continuation machinery.

**Cycle 2 (independent ideation review).** Review found gaps in that controller's action vocabulary, gate binding, interception proof, and report identity. These findings apply only to the superseded controller design.

**Cycle 3 (captain scope reset, 2026-07-17).** The controller design was rejected as unnecessary. The intended product is exactly two hints: recommend compaction only at an already-durable boundary, then remind the post-compaction FO to reread its authoritative contract and reconcile durable state. This body replaces the rejected design rather than repairing it.

**Cycle 4 (captain host-neutral reframe, 2026-07-18).** The two hints were reframed from Codex-specific to host-neutral FO rules at the shared-core level, inherited by every runtime adapter. Post-compaction DELIVERY became a per-host binding `«post-compact-notice»` mirroring `«context-budget»`'s PRESENT/ABSENT discipline: Codex = captain-facing UI warning + manual cue (proven via the 0.144.4 probe); Pi = manual cue. The Claude `SessionStart(compact)` model-context delivery was SPLIT to `claude-post-compaction-contract-reload` (id `cdbhzxc`), referenced here but not designed. Two hints, no controller; the superseded-controller artifact reference is preserved.

## Stage Report: ideation (cycle 3)

- DONE: Replace the continuation controller with exactly two bounded hints.
  The current body contains one safe-boundary compaction suggestion and one post-compaction contract-reload reminder. It removes the authorization ledger, permits, checkpoint, gateway, crash replay, interception matrix, monitor, watchdog, and their derived acceptance criteria.
- DONE: Exercise and honestly bound the current Codex lifecycle surface.
  Live Codex 0.144.4 manual `/compact` fired `PreCompact` and `PostCompact`; the shipped-shape `PostCompact` `systemMessage` appeared as a warning but was absent from the next model's context. `SessionStart(compact)` did not fire. The proposed hook is therefore captain-facing and failure-open, with a manual cue fallback.
- DONE: Rewrite the problem, proposed approach, acceptance criteria, test plan, and documentation change around observable hint timing and reload behavior.
  AC-1 pairs safe and unsafe pressure scenarios, AC-2 proves the visible warning without claiming model injection, AC-3 proves actual contract/status/roster/report reads before the next effect, and AC-4 proves harmless absence. The full rejected design remains beside the entity as an artifact.

### Summary

c6 is now a small continuity aid: the FO suggests compaction only when existing state is recoverable, and Codex warns the captain after compaction to trigger a real contract reload and state reconciliation. The live host does not inject that warning into model context, so the design says so and falls back to a manual captain cue instead of adding lifecycle machinery.

## Stage Report: ideation (cycle 4)

- DONE: State the two hints as host-neutral FO rules at the shared-core level (not Codex-only); the reload rule names the host-neutral reread target (`SKILL.md` + `@first-officer-shared-core.md` + the active host runtime adapter) and the existing boot/status/roster/report reconciliation flow.
  Decision + Rule 1/Rule 2 in Proposed approach now read as shared-core FO judgment rules inherited by every adapter; the reread target and its identical-on-every-host framing are stated in Problem and Rule 2. Trigger weakness is framed host-neutrally.
- DONE: Make post-compaction DELIVERY an abstracted per-host binding mirroring `«context-budget»` PRESENT/ABSENT — Codex = UI warning + manual cue with the probe artifact as proof; Claude = SPLIT to `cdbhzxc`, referenced not designed; Pi = manual cue.
  Added the `«post-compact-notice»` binding block with per-host `→` lines and a failure-open `block:`; the Claude `SessionStart(compact)` delivery is referenced and explicitly not designed here. Shared hook-surface prerequisite recorded (no `hooks` key in either manifest, verified).
- DONE: Rewrite Decision, Problem, Proposed approach, Acceptance criteria, Test plan, Documentation change, and Out of scope to match, and add a `## Feedback Cycles` Cycle 4 entry recording the reframe and the split.
  ACs renumbered host-neutral-first (AC-1 safe timing value AC, AC-2 reload-before-effect, AC-3 Codex binding proof, AC-4 harmless absence); test plan, doc diff, and Out of scope updated; the blanket "Claude or Pi changes" exclusion replaced by the split-task + Pi-manual-cue carve-out while the controller/ledger/permits/gateway/crash-replay/watchdog exclusions and the superseded-controller artifact reference are preserved.

### Summary

The entity is now host-neutral: two FO judgment rules live at the shared-core level and apply to every host, while post-compaction DELIVERY is the per-host binding `«post-compact-notice»` (Codex UI-warning+cue proven, Pi manual cue), with the Claude `SessionStart(compact)` auto-inject delivery split to `cdbhzxc` and referenced but not designed. No controller; the superseded design and Codex probe remain as artifacts, and the shared hook-shipping prerequisite is flagged.

## Stage Report: implementation

- DONE: Shared-core host-neutral Rules 1 & 2 plus the «post-compact-notice» capability (fo-dispatch-core.md + docs/runtime-support.md) shipped as designed; Rule 2 names the reread target (SKILL.md + first-officer-shared-core.md + active host runtime adapter) and the boot/status/roster/report reconciliation flow. (AC-2, AC-4)
  "## Compaction continuity" added to first-officer-shared-core.md (before-compaction durable-boundary suggest + after-compaction reread-and-reconcile), the «post-compact-notice» binding block in fo-dispatch-core.md, and the shape line in docs/runtime-support.md; all contractlint invariants green (capability define/call balance, per-host Claude/Codex/Pi → coverage, size ceilings). Code commit d124f6f6.
- DONE: Codex «post-compact-notice» binding bundled in .codex-plugin (establishing the plugin hooks surface) emitting the exact reread-and-reconcile systemMessage, failure-open; the first spike confirms .codex-plugin/plugin.json can declare a hooks key that Codex loads and trusts. (AC-3)
  .codex-plugin/plugin.json gains `"hooks": "./hooks.json"`; hooks.json declares PostCompact (matcher manual|auto) → ./hooks/codex_post_compact_notice.sh, a bundled script emitting the exact systemMessage. Live Codex 0.144.4 spike proved the hooks key LOADS + TRUSTS + FIRES reproducibly and is inert/harmless untrusted: artifacts/codex-0.144.4-plugin-hooks-spike.md.
- DONE: Offline fixtures prove the value ACs: Rule-1 safe/unsafe pressure-timing with commit-OID evidence (AC-1), split-root post-compaction reload-before-effect replay (AC-2), Codex hook JSON manual|auto + absent/disabled/failing matrix (AC-3, AC-4).
  internal/ensigncycle: compaction_timing_test.go (AC-1 real git OIDs, safe→one/unsafe→zero, + reject a suggestion at an unrecoverable boundary), compaction_reload_test.go (AC-2 ordered oracle over a replay + reject stale-summary skip / missing adapter reread), codex_post_compact_hook_test.go (AC-3 manual|auto → one valid systemMessage + AC-4 absent/disabled/failing no-writes matrix). All green.

### Summary

Shipped the two host-neutral FO compaction-continuity judgment rules (boot-resident shared core) plus the per-host `«post-compact-notice»` delivery binding, and the Codex leg: a bundled `.codex-plugin` `PostCompact` hook whose script emits the reread-and-reconcile `systemMessage`, failure-open. The riskiest claim — that a plugin-declared hooks key loads and trusts on Codex — was spiked FIRST against a live 0.144.4 CLI and proven reproducibly (form: `"hooks": "./hooks.json"` + a `hooks.json` + a bundled script; the vendor `plugin-creator` linter's "hooks unsupported" claim is stale, and curated OpenAI plugins ship the same mechanism). Key nuance in the spike artifact: hook commands run with no shell (whitespace-split, no quote-stripping), so the reminder must be a bundled script, not an inline command. Contract additions were trimmed to fit under the shared-core (<26755B) and FO-surface (<122400B) ratchet ceilings. Full non-live module suite is green except the pre-existing environmental `internal/cli` `TestCodexResolveManifestAgainstInstalledHost`, which fails only because the sandbox cannot read the real `~/.codex/config.toml` (it skips cleanly with a readable CODEX_HOME) — unrelated to this change.
