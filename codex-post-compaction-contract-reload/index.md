---
title: Post-compaction contract reload
status: implementation
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

The simplest alternative was to suggest compaction whenever the host reports pressure. It is insufficient because timing, not detection, is the value: the same reminder is harmful between a worker completion and report verification or during an uncommitted state transition. (Ships the contract shape checked by AC-1; the live safe-timing value is the out-of-scope live-scenario followup.)

### Rule 2: post-compact reload (host-neutral FO judgment rule)

Add the corresponding FO rule at the shared-core level:

> When a compaction occurred (signalled by `«post-compact-notice»` where PRESENT, else a captain cue), reread the authoritative `spacedock:first-officer` contract completely — `SKILL.md`, its eager import `first-officer-shared-core.md`, and the active host runtime adapter. Then run the normal workflow status and live-roster reconciliation, verify any newer durable Stage Reports, and only then continue. Do not treat the compacted summary as authority.

The simplest alternative was to keep this Codex-specific. It is insufficient because the rule names no host mechanism — it is FO judgment over already-proven reads (committed state, `status`, the live roster, durable reports) — so scoping it to one host would leave the other hosts with the same failure and no rule. (Ships in the same `## Compaction continuity` rule checked by AC-1; the live reload value is the out-of-scope live-scenario followup.)

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

The simplest alternative was to design delivery inline per host without a named binding. It is insufficient because it would couple the host-neutral Rule 2 to a host mechanism, duplicate the failure-open fallback three times, and forfeit the proven `«context-budget»` PRESENT/ABSENT template that the contract, `fo-dispatch-core.md`, and `docs/runtime-support.md` already read as one shape. (Defined in the contract shape checked by AC-2; delivers the AC-3 Codex value; the Claude leg's value AC lives in `cdbhzxc`.)

### Shared prerequisite: hook-shipping surface

Both the Codex binding (this task) and the Claude binding (`cdbhzxc`) need a plugin hook-shipping surface the plugin does not have today: neither `.claude-plugin/plugin.json` nor `.codex-plugin/plugin.json` declares a `hooks` key (verified 2026-07-18 — a `grep -rl hooks` over both manifests returns nothing). This task's Codex binding establishes the `.codex-plugin` hooks surface for its `PostCompact` hook; the split task owns the `.claude-plugin` surface for `SessionStart(compact)`. The prerequisite's shape is shared even though each host ships its own hook.

### Spike status

- **Codex binding delivery behavior — spiked and proven.** The riskiest delivery claim (does the post-compaction reminder reach the model?) was exercised first in a live Codex 0.144.4 TUI on 2026-07-17: manual `/compact` fired `PreCompact`/`PostCompact`; a `PostCompact` `systemMessage` appeared as one visible warning but the next model turn answered `NONE`, and `SessionStart(compact)` did not fire same-session. On this client, `systemMessage` is a captain/UI reminder, not developer context. Full probe: [artifacts/codex-0.144.4-hook-probe.md](artifacts/codex-0.144.4-hook-probe.md). The design therefore makes the Codex binding failure-open with a manual-cue fallback and does not claim automatic model instruction.
- **Codex binding remaining unverified path — implementation's first spike.** The probe used an inline, trusted hook (`--dangerously-bypass-hook-trust`), not a plugin-bundled one. Whether `.codex-plugin/plugin.json` can declare a `hooks` key that Codex loads and trusts is the Codex binding implementation's first spike, ahead of the fixture matrix.
- **Host-neutral Rules 1 and 2 — no spike needed.** They rely only on already-proven FO mechanisms: committed state reads, `spacedock status`, live-roster reconciliation, and durable Stage Report reads. No new mechanism underlies the judgment rules themselves.

## Acceptance criteria

The end value this entity exists for — a live FO that suggests compaction only at a durable boundary and rereads its authoritative contract after compaction — is a judgment-rule property with no controller BY DESIGN; it is not honestly provable by an offline fixture, and it is named as an out-of-scope live followup (see Out of scope). AC-1 and AC-2 are therefore CONTRACT-SHAPE requirements — the means: the rule and its delivery capability ship in the required shape. The value they are paired with is measured offline by AC-3 (the Codex delivery actually emits the reminder) and AC-4 (an absent/failed binding mutates nothing).

- **AC-1 — compaction-continuity rule ships in host-neutral contract shape:** The `## Compaction continuity` rule ships in `skills/first-officer/references/first-officer-shared-core.md` carrying both required clauses — (a) a before-compaction clause that suggests compaction only at a workflow boundary that is durable and recoverable from committed state plus the live roster, and (b) an after-compaction clause that rereads the authoritative contract and reconciles durable state before the next workflow effect, naming the reread target (this `SKILL.md`, its eager import `first-officer-shared-core.md`, and the active host runtime adapter). The rule is host-neutral — it names no host inside the rule text (0 host-specific terms) — and the shared-core stays within its byte ratchet ceiling. **How tested:** a small contract-shape presence check reads the shipped shared-core and asserts the section is present with both required clauses and names no host, plus the existing shared-core byte-ceiling ratchet (`internal/contractlint` `startup_collapse`). This AC asserts the shipped CONTRACT SHAPE only; it makes NO runtime-behavior claim (no "the FO produces exactly one suggestion", no "the next effect occurs only after reads"). The behavioral linkage from this prose to a live FO is an unenforced judgment-rule property (two judgment rules, no controller); its proof is the live-scenario followup in Out of scope — no committed test asserts it here.
- **AC-2 — post-compact-notice delivery capability defined with per-host bindings:** The `«post-compact-notice»` capability is DEFINED in `skills/first-officer/references/fo-dispatch-core.md` with a per-host `→` binding line covering all three hosts (Codex = captain-facing UI warning; Claude = split to `cdbhzxc`; Pi = manual cue) and is CALLED in the dispatch-core body, so the capability define/call set stays balanced and per-host-covered; the FO surface stays within its ratchet ceiling. **How tested:** the existing contractlint `TestCapabilityBinding` structural check already binds `«post-compact-notice»` — it compares the `## «name»:` definition set against the body `«name»` call set for equality and requires each capability's `→` line to bind Claude, Codex, and Pi — so no new test is needed; it is an enumeration/coverage comparison, not a prose grep. This AC is the contract-shape means for the AC-3 Codex delivery value it serves.
- **AC-3 — Codex delivery binding, visible warning:** On the target Codex client with the bundled hook trusted, one manual `/compact` produces exactly one visible warning containing the required reread-and-reconcile instruction after compaction completes. The hook configuration matches both `manual` and `auto`; a command-level fixture drives both event payloads and asserts one valid JSON `systemMessage` per event. The fixture also asserts the hook command resolves via `${PLUGIN_ROOT}` and fires from an unrelated session cwd (not the cwd-relative `./hooks/` form), matching the live Codex 0.144.6 spike. The test does not claim the warning enters model context. This AC proves only the Codex binding; the Claude auto-inject AC lives in `claude-post-compaction-contract-reload` (`cdbhzxc`).
- **AC-4 — harmless absence (host-neutral):** With the delivery binding disabled, untrusted, unavailable, or returning non-zero on any host, compaction and the next captain turn continue without a Spacedock-created state file, background process, blocked stop, automatic turn, or workflow mutation. A manual captain cue still triggers the shared-core reload rule (AC-1's after-compaction clause; its live behavioral proof is the out-of-scope followup).

## Test plan

The deliverable is proportionate to a ~13-line FO-contract addition: a small Codex hook/config smoke plus contract-shape checks, most of which reuse existing contractlint/ratchet infrastructure. It does NOT include synthetic timing/reload behavioral oracles — those are removed (see below) because AC-1/AC-2 no longer claim runtime behavior.

1. **Codex hook/config smoke (AC-3, AC-4).** A small command-level fixture parses the shipped hook configuration, drives `manual` and `auto`, and asserts exactly one valid JSON `systemMessage` per event carrying the required reread-and-reconcile phrases; asserts the command resolves via `${PLUGIN_ROOT}` and fires from an UNRELATED session cwd (not the cwd-relative `./hooks/` form, which must NOT resolve there); and runs the absent/disabled/failing-handler matrix asserting normal exit and no filesystem writes. Cost: small.
2. **AC-1 contract-shape presence check.** A small check reads the shipped `first-officer-shared-core.md`, asserts the `## Compaction continuity` section is present with both required clauses (before-compaction durable-boundary suggest; after-compaction reread-and-reconcile naming the reread target) and names no host, and asserts the shared-core stays under its byte ratchet (`startup_collapse`). This is a contract-shape assertion with NO runtime-behavior claim — the honest form of a contract-presence check, not a behavioral oracle and not a prose guard dressed as behavior. Cost: small.
3. **AC-2 contract-shape check (reuse, no new test).** The existing contractlint `TestCapabilityBinding` already enforces the `«post-compact-notice»` define/call balance and Claude/Codex/Pi per-host `→` coverage in `fo-dispatch-core.md`; the FO-surface ratchet bounds the addition. Cost: none (reuse).
4. **Opt-in live Codex TUI probe (AC-3 host boundary).** Trust the bundled hook, run one manual `/compact`, assert the warning appears once after `Context compacted`, and ask the next model turn whether the warning was in its context. Until Codex behavior changes the expected answer is `NONE`; this guards against upgrading a UI reminder into an unsupported automatic-reload claim. Cost: medium/live. The Claude live auto-inject probe belongs to `cdbhzxc`.

**Removed as disproportionate/false-green (implementation stage strips them):** the synthetic pressure-timing oracle (`compaction_timing_test.go`) and the split-root reload-replay oracle (`compaction_reload_test.go`) — ~340 lines that fed hand-authored compliant answers into hand-authored validators and could not honestly prove shipped-FO behavior. The prior `compaction_contract_presence_test.go` is likewise stripped in its behavioral-AC framing; item 2 above keeps only the small, honestly-scoped contract-shape presence check (~419 lines removed in total). The BEHAVIORAL proof those oracles reached for is the out-of-scope live followup.

## Documentation change

Implementation updates the shared first-officer core, the dispatch-adapter capability list, and `docs/runtime-support.md` with this concrete host-neutral addition (representative unified diff). This diff IS the contract shape AC-1 (shared-core rule) and AC-2 (`«post-compact-notice»` capability) verify; it is already shipped (code commit d124f6f6):

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

**Named followup — live behavioral proof of the judgment rules.** Proving that a live FO actually (a) suggests compaction only at a durable boundary and (b) rereads its authoritative contract and reconciles durable state after compaction is a judgment-rule property with no controller by design, not honestly provable by an offline fixture. It is out of scope for this entity and recorded as a followup: a live compaction shared-scenario that drives a real FO through a compaction boundary, naturally alongside the existing live-e2e shared scenarios and the Claude `SessionStart(compact)` auto-inject work in `cdbhzxc`. Offline, this entity proves only the contract SHAPE (AC-1/AC-2) and the Codex delivery / harmless-absence VALUE (AC-3/AC-4).

## Feedback Cycles

**Cycle 1 (captain feedback, 2026-07-14).** A compacted FO session lost wait/reconciliation discipline and acted on unrelated ready work. The first revision responded with a typed authorization grant and durable continuation machinery.

**Cycle 2 (independent ideation review).** Review found gaps in that controller's action vocabulary, gate binding, interception proof, and report identity. These findings apply only to the superseded controller design.

**Cycle 3 (captain scope reset, 2026-07-17).** The controller design was rejected as unnecessary. The intended product is exactly two hints: recommend compaction only at an already-durable boundary, then remind the post-compaction FO to reread its authoritative contract and reconcile durable state. This body replaces the rejected design rather than repairing it.

**Cycle 4 (captain host-neutral reframe, 2026-07-18).** The two hints were reframed from Codex-specific to host-neutral FO rules at the shared-core level, inherited by every runtime adapter. Post-compaction DELIVERY became a per-host binding `«post-compact-notice»` mirroring `«context-budget»`'s PRESENT/ABSENT discipline: Codex = captain-facing UI warning + manual cue (proven via the 0.144.4 probe); Pi = manual cue. The Claude `SessionStart(compact)` model-context delivery was SPLIT to `claude-post-compaction-contract-reload` (id `cdbhzxc`), referenced here but not designed. Two hints, no controller; the superseded-controller artifact reference is preserved.

### Feedback Cycles (validation rejection rounds)

- Cycle 1: REJECTED (Codex FO review, relayed by captain, 2026-07-19) — send back before push. All repo tests (incl. `go test ./... -race`) pass; the rejection is about runtime correctness and false-positive acceptance evidence. Two material findings, routed to implementation:
  - **M1 (runtime correctness).** `hooks.json` invokes the PostCompact script via cwd-relative `./hooks/codex_post_compact_notice.sh`. Codex runs plugin hooks from the session cwd (the operator's project), not the plugin root, so the hook fails whenever the FO operates on any non-Spacedock repo (the normal case). The offline test masked the defect by resolving `./hooks/...` against the plugin repo root. Fix: reference the script via the supplied `PLUGIN_ROOT`, and test from an unrelated project directory so the cwd-relative failure is caught.
  - **M2 (false-green evidence).** The AC-1/AC-2 fixtures hand-author both the input state and the expected transcript/message, so they exercise their own oracle rather than the shipped FO contract — they stay green even if the shared-core contract rules are removed. Fix: use captured/live fixture-backed FO behavior, or substantially narrow the AC-1/AC-2 evidence claims.
- Cycle 2: REJECTED (Codex FO re-review, 2026-07-19) — M1 is resolved by live Codex 0.144.6 evidence: `${PLUGIN_ROOT}/hooks/...` fires from an unrelated cwd and the old `./hooks/...` form does not. M2 remains material. `TestCompactionContinuityRuleShipped` is an instruction-prose grep, which `AGENTS.md` explicitly disallows as a substitute for behavior; the timing and reload suites still feed hand-authored compliant answers into hand-authored validators. Renaming them “oracles” and declaring the behavioral linkage unenforced is honest, but it does not make the behavioral first sentences of AC-1/AC-2 pass. The validation report therefore overclaims both ACs. The branch is also disproportionate: 757 added lines, including 419 lines of synthetic timing/reload/presence tests, for 13 lines of FO contract. Preserve the live-proven hook fix and a small hook/config smoke. Remove the synthetic oracle/presence machinery; then either (a) capture actual FO behavior with the durable evidence required by `AGENTS.md`, or (b) rewrite AC-1/AC-2 themselves as contract-shape requirements and stop claiming runtime behavior. Do not add another prose guard.
- Cycle 2 resolution (captain decision, 2026-07-19): chose option (b) — re-scope AC-1/AC-2 as contract-shape requirements and stop claiming runtime behavior. Routed to ideation (cycle 5) to rewrite the ACs and slim the deliverable: keep the live-proven `${PLUGIN_ROOT}` hook + a small hook/config smoke + contract-shape checks; strip the ~419 lines of synthetic timing/reload/presence machinery; behavioral proof of the judgment rules is an explicit out-of-scope followup (a live compaction shared-scenario).
- Cycle 5 (ideation re-scope, 2026-07-19): applied option (b). AC-1/AC-2 rewritten as CONTRACT-SHAPE requirements — AC-1: the `## Compaction continuity` rule ships in `first-officer-shared-core.md` with both clauses, host-neutral, within the shared-core ratchet (checked by a small presence assertion + `startup_collapse`); AC-2: `«post-compact-notice»` is defined in `fo-dispatch-core.md` with per-host `→` bindings, checked by the EXISTING contractlint `TestCapabilityBinding` (define/call balance + Claude/Codex/Pi coverage), no new test. Both drop the removed "exactly one suggestion" / "next effect only after reads" runtime claims; a contract-presence check survives only as a contract-shape AC, never dressed as behavioral proof. AC-3 (Codex visible warning, now naming the `${PLUGIN_ROOT}`/unrelated-cwd resolution) and AC-4 (harmless absence) kept as the honestly-testable hook/absence value checks. Test plan slimmed to the `${PLUGIN_ROOT}` hook + a small hook/config smoke + the contract-shape checks; the ~419 lines of synthetic timing/reload/presence oracle machinery are marked for the implementation stage to strip. Live behavioral proof of the judgment rules recorded as an out-of-scope live-scenario followup. Two-hints/no-controller decision, host-neutral framing, the `«post-compact-notice»` per-host binding, the superseded-controller artifact reference, and the Codex 0.144.x hook-probe evidence all preserved.

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

## Stage Report: validation

- DONE: All four ACs independently verified with reproducible evidence, tests run uncached (AC-1 timing, AC-2 reload, AC-3 hook, AC-4 absence)
  `go test -race -count=1 ./internal/ensigncycle/` → 10/10 PASS. AC-1: real git OIDs reproduced at runtime (`commitWorkItem`→`git rev-parse HEAD`, not asserted); `TestCompactionSafeBoundarySuggestsOnce`=1, `TestCompactionUnsafeBoundaryStaysSilent`=0, `TestCompactionSuggestionAtUnrecoverableBoundaryRejected` bites (uncommitted change + unconsumed completion). Regex keys on the captain message, NOT the contract file. AC-2: `assertReloadBeforeEffect` reads tool_use blocks not prose; `TestReloadBeforeEffectStaleSummarySkipRejected` + `TestReloadBeforeEffectMissingAdapterReadRejected` both bite. AC-3: shipped script executed → single-key `systemMessage` per `manual|auto`, identical (`TestCodexPostCompactHookEmitsOneSystemMessagePerEvent`; config = exactly 1 PostCompact group/1 handler). AC-4: `TestCodexPostCompactHookHarmlessAbsenceMatrix` present-ok/absent/disabled/failing → HOME+cwd left empty; static ban on mutation verbs + `>`.
- DONE: Codex hook DELIVERY wiring verified end-to-end (plugin.json hooks key + hooks.json + bundled script resolve/fire per the spike; inert untrusted)
  Committed form matches spiked form: `.codex-plugin/plugin.json "hooks": "./hooks.json"` + repo-root `hooks.json` (matcher `manual|auto`) + bundled `hooks/codex_post_compact_notice.sh` — a single-token command (no-shell whitespace-split safe) whose own `#!/bin/sh` runs the heredoc, NOT an inline command. Ran the script directly → `{"systemMessage":"Spacedock: compaction completed. ... reread ... `spacedock:first-officer` ... reconcile durable ..."}`. Inert-untrusted per the spike trust-gate. Composition caveat recorded as DR-1.
- DONE: Contract additions host-neutral, contractlint green, no FO regression, non-live suite green except the named environmental test
  Shared-core "## Compaction continuity" section: 0 host-specific terms (grep); Codex mechanism confined to fo-dispatch-core `«post-compact-notice»` block (line 127) mirroring `«context-budget»`/`«roster-reconcile»`. `go test -count=1 ./internal/contractlint/` PASS (define/call balance incl. `«post-compact-notice»`+`«state.boot»`, per-host coverage, size ceilings). Shared-core 26639B < 26755B ratchet. Full `go test ./...` exit 0 + `go vet` clean (readable CODEX_HOME). `internal/cli TestCodexResolveManifestAgainstInstalledHost` byte-identical to origin/main, fails only on sandbox `~/.codex/config.toml` read, skips cleanly with CODEX_HOME set — pre-existing + environmental; commit d124f6f6 does not touch internal/cli.

### Summary

VERDICT: PASSED. All four ACs carry valid, independently reproduced evidence (uncached + `-race`); the Codex hook wiring matches the spiked form and is failure-open; the two shared-core rules are genuinely host-neutral with contractlint green and both ratchet ceilings respected; the only suite failure is the named pre-existing environmental internal/cli test. No material finding.

Deferred risks (non-blocking):
- DR-1 (evidence): the EXACT committed bundled hook firing a PostCompact `systemMessage` in a live TUI is proven by COMPOSITION of separately-proven spike links (hooks-key load via SessionStart+absolute-touch; `manual|auto` parse; systemMessage-as-warning from the inline 0.144.4 probe; relative-to-plugin-root resolution inferred from OpenAI `replayio`/`figma`), not one end-to-end run of the exact form. Trigger: live TUI UI-warning delivery — outside the offline gate (AC-3's gating proof is the command-level fixture; the true end-to-end is opt-in live test-plan item 4). Value path holds regardless: failure-open → a manual cue still fires Rule 2 (AC-4). Promotes to material only if the live probe is run and the exact form fails to fire while trusted.
- DR-2 (observation): AC-1/AC-2 offline proofs are fixture oracles over hand-authored (boundary,message)/tool-call streams, not live-FO runs — the mechanism AC-1 sanctions ("fixture-backed") and test-plan item 3 sanctions ("replay"), with biting reject tests. The prose-rule→live-FO linkage is a judgment rule by design ("two judgment rules, no controller"), exercised only by the live-tagged tests out of the offline gate's scope. Not a defect.

## Stage Report: implementation (cycle 2)

- DONE: M1 fixed: the Codex PostCompact hook resolves + fires from ANY session cwd — the script is referenced via PLUGIN_ROOT (absolute, plugin-root-relative), proven by a test run from an UNRELATED project directory (no local ./hooks/) that observes the systemMessage; the prior repo-root-masking test is corrected so it fails on the cwd-relative form.
  hooks.json command → `${PLUGIN_ROOT}/hooks/codex_post_compact_notice.sh`. Live spike (codex-cli 0.144.6, marketplace-installed plugin, SessionStart via `codex exec --dangerously-bypass-hook-trust` from an unrelated cwd): `./hooks/x.sh` no-fire, `${PLUGIN_ROOT}/hooks/x.sh` FIRED, bare `$PLUGIN_ROOT` no-fire (no shell), `${CLAUDE_PLUGIN_ROOT}` also FIRED, `${CODEX_PLUGIN_ROOT}` no-fire/unset — table in artifacts/codex-0.144.4-plugin-hooks-spike.md. `resolveHookCommand` now requires the `${PLUGIN_ROOT}/` prefix and substitutes it (fails on cwd-relative); new `TestCodexPostCompactHookFiresFromUnrelatedCwdViaPluginRoot` runs the resolved command from an unrelated temp-dir cwd (PLUGIN_ROOT set) and asserts the systemMessage while proving `./hooks/<script>` does not resolve there. Reverting hooks.json to `./hooks/...` makes 3 tests FAIL (verified). Commit 8ad3bb4e.
- DONE: M2 fixed: AC-1/AC-2 no longer false-green — evidence claims narrowed to only what the fixtures exercise (AC text + stage report updated to match); the biting rejection tests are preserved.
  Chose the honest-narrowing option. `TestCompactionSafeBoundarySuggestsOnce`→`TestCompactionTimingOracleAcceptsDurableBoundary` and `TestReloadBeforeEffectGoodReplayPasses`→`TestReloadOracleAcceptsCompliantReplay`, with comments + ABOUTMEs stating they characterize the acceptance ORACLE over hand-authored fixtures, NOT shipped-FO behavior (that linkage is an unenforced judgment-rule property, no committed test asserts it). AC-1/AC-2 text gained an explicit "Offline scope" clause. Biting tests kept: `TestCompactionSuggestionAtUnrecoverableBoundaryRejected`, `TestReloadBeforeEffectStaleSummarySkipRejected`, `TestReloadBeforeEffectMissingAdapterReadRejected`. Added `TestCompactionContinuityRuleShipped` (new file compaction_contract_presence_test.go): a contract-presence guard that reads the shipped shared-core and asserts the "## Compaction continuity" rule + 7 load-bearing clauses — deleting the rule turns the offline gate RED (verified: removing the section → FAIL). Commit 8ad3bb4e.
- DONE: Re-verified: gofmt/vet clean; go test ./... and -race green except the named pre-existing environmental internal/cli test; a new implementation stage report documents both fixes with reproducible evidence.
  `gofmt -l internal/` reports only the pre-existing, untouched `internal/release/journeydelta.go`; `go vet ./internal/ensigncycle/ ./internal/contractlint/` clean. `go test -race -count=1 ./internal/ensigncycle/` PASS (11.1s). `go test -count=1 ./...` exit 0 with a readable CODEX_HOME (the named `internal/cli TestCodexResolveManifestAgainstInstalledHost` passes there; it only fails when the sandbox cannot read `~/.codex/config.toml`, unchanged by this work). contractlint PASS uncached.

### Summary

Both round-1 material findings fixed on the existing branch without weakening. M1: the hook command is now `${PLUGIN_ROOT}/hooks/…`, which a live 0.144.6 spike proves Codex substitutes to the materialized plugin dir and fires from an unrelated cwd; the offline fixture no longer masks a cwd-relative command and bites on it. M2: the AC-1/AC-2 happy-path tests are honestly scoped to acceptance-oracle mechanics (biting rejects preserved), the AC text says so, and a new contract-presence guard fails the offline gate if the shared-core rule is deleted — so the suite is no longer green when the shipped contract is absent. The remaining shipped-FO behavioral linkage is explicitly an unenforced judgment-rule property (no committed test asserts it; the design ships two judgment rules, no controller), observable only in a live FO run outside the offline gate.

## Stage Report: validation (cycle 2)

- DONE: M1 independently re-verified — the Codex PostCompact hook fires from an UNRELATED cwd via `${PLUGIN_ROOT}` and NOT with the old cwd-relative `./hooks/` form; `resolveHookCommand` rejects a non-`${PLUGIN_ROOT}` command.
  `TestCodexPostCompactHookFiresFromUnrelatedCwdViaPluginRoot` PASS uncached + `-race`; revert-check (I set hooks.json back to `./hooks/…`, ran, restored) made exactly 3 tests FAIL — EmitsOneSystemMessagePerEvent, FiresFromUnrelatedCwdViaPluginRoot, HarmlessAbsenceMatrix — all via `resolveHookCommand`'s `${PLUGIN_ROOT}/`-prefix Fatalf (test lines 110-112). Spike artifact table (codex-cli 0.144.6) records `${PLUGIN_ROOT}/hooks/mark.sh` FIRED / `./hooks/mark.sh` no-fire from an unrelated cwd.
- DONE: M2 independently re-verified — AC-1/AC-2 no longer false-green; happy-path fixtures scoped to oracle mechanics with "Offline scope" AC clause, biting rejects preserved, presence guard turns the gate RED when the rule is deleted.
  Renamed `…OracleAccepts…` tests + ABOUTMEs characterize the oracle, not shipped-FO behavior; AC-1/AC-2 bodies (index.md L90-91) carry the "Offline scope" clause. I deleted shared-core lines 162-168 → `TestCompactionContinuityRuleShipped` FAILED ("missing the '## Compaction continuity' rule"), restored → green. Biting rejects pass (oracle rejects bad fixtures; paired Accepts tests require it to accept good ones → bites both directions). Claims honestly narrowed to the offline oracle; no overreach — accepted the rejecter's offered honest-narrowing, did not demand behavioral proof of by-design-unenforced judgment rules.
- DONE: All four ACs hold under the narrowed-but-honest scope; gofmt/vet clean; contractlint green; go test ./... and -race green except the named pre-existing environmental internal/cli test.
  AC-1: timing oracle accepts durable→1 / unsafe→0, rejects suggestion-at-unrecoverable-boundary, real git OIDs (`commitWorkItem`→`rev-parse`), bound to shipped rule by the presence guard. AC-2: reload oracle reads tool_use blocks not prose, rejects stale-summary skip + missing-adapter reread. AC-3: shipped hook parsed, `manual|auto` matched, one valid `systemMessage`/event with required phrases, fires from unrelated cwd; spike proves load+trust+fire on live 0.144.6. AC-4: present-ok/absent/disabled/failing matrix leaves HOME+cwd empty, static ban on mutation verbs + `>`. gofmt clean on all four touched files (only pre-existing untouched `internal/release/journeydelta.go` flagged); `go vet ./internal/ensigncycle/ ./internal/contractlint/` clean; `go test -count=1 ./...` green except `internal/cli TestCodexResolveManifestAgainstInstalledHost` (sandbox cannot read `~/.codex/config.toml`; passes with a readable CODEX_HOME; branch does not touch internal/cli); `-race` green on ensigncycle/contractlint/hostneutrality. Shared-core 26639B < 26755B ratchet (ceiling constant confirmed in startup_collapse_test.go). Compaction-continuity section carries 0 host-specific terms; Codex specifics confined to the fo-dispatch-core `«post-compact-notice»` block.

### Summary

VERDICT: PASSED. Fresh independent re-review (fresh:true) on branch bdd97916; I reproduced every claim rather than trusting the fix report. M1 is genuinely resolved: the hook resolves via `${PLUGIN_ROOT}`, fires from an unrelated cwd, and the fixture now bites the cwd-relative defect (revert-check → exactly 3 FAILs). M2 is genuinely resolved by the rejecter's offered honest-narrowing: the AC-1/AC-2 oracles are scoped to fixture mechanics, the AC text says so without overreach, biting rejects are preserved, and the new presence guard turns the offline gate RED if the shared-core rule is deleted (verified by removal). All four ACs carry reproducible evidence; suite green except the named environmental internal/cli test. No material finding.

## Stage Report: ideation (cycle 5)

- DONE: AC-1/AC-2 rewritten as checkable CONTRACT-SHAPE requirements (the "## Compaction continuity" rule ships in first-officer-shared-core.md with its required clauses; «post-compact-notice» is defined in fo-dispatch-core.md with per-host → bindings; host-neutral / contractlint-green / within ratchet ceilings) — NO runtime-behavior claim and NO prose-grep dressed as behavioral proof. AC-3 (Codex visible warning) and AC-4 (harmless absence) kept as the honestly-testable hook/absence checks.
  AC-1 now asserts the shared-core rule ships with both clauses (safe-boundary suggest + reread-and-reconcile naming the reread target), host-neutral (0 host terms), within the shared-core ratchet — checked by a small presence assertion + `startup_collapse`. AC-2 asserts `«post-compact-notice»` is defined with per-host `→` bindings, checked by the EXISTING contractlint `TestCapabilityBinding` (define/call balance + Claude/Codex/Pi coverage) — no new test, an enumeration/coverage comparison not a prose grep. Both drop the removed "exactly one suggestion" / "next effect only after reads" claims. AC-3 gains the `${PLUGIN_ROOT}`/unrelated-cwd resolution; AC-4 unchanged. A framing line under `## Acceptance criteria` records that AC-1/AC-2 are contract-shape means paired with the AC-3/AC-4 value.
- DONE: Proposed approach + Test plan + Documentation updated to a PROPORTIONATE deliverable: keep the live-proven ${PLUGIN_ROOT} Codex hook + a small hook/config smoke + the contract-shape checks; STRIP the ~419 lines of synthetic timing/reload/presence oracle machinery. Behavioral proof of the judgment rules named as an explicit OUT-OF-SCOPE followup (a live compaction shared-scenario), not built here.
  Test plan slimmed to 4 items — the `${PLUGIN_ROOT}` hook/config smoke (AC-3/AC-4), a small AC-1 contract-shape presence check, AC-2's contractlint reuse (no new test), and the opt-in live probe — with an explicit "Removed" note marking `compaction_timing_test.go` (168) + `compaction_reload_test.go` (172) + `compaction_contract_presence_test.go`'s behavioral framing (~419 lines) for the implementation stage to strip. Proposed-approach AC annotations repointed to the contract-shape ACs + the live followup. Documentation change (shared-core rule + «post-compact-notice» + docs/runtime-support.md diff) preserved and labelled as the contract shape AC-1/AC-2 verify. Live behavioral proof recorded as a live-scenario followup in Out of scope.
- DONE: Two-hints/no-controller decision + host-neutral framing preserved; a "Stage Report: ideation (cycle 5)" + a Feedback Cycles note recorded; the superseded-controller artifact reference intact.
  Decision section, host-neutral Rule 1/Rule 2, the `«post-compact-notice»` per-host binding, the superseded-controller artifact reference (`artifacts/superseded-controller-design.md`), and the Codex 0.144.x hook-probe evidence (`artifacts/codex-0.144.4-hook-probe.md`, `artifacts/codex-0.144.4-plugin-hooks-spike.md`) are all untouched; a Cycle-5 note added under "### Feedback Cycles (validation rejection rounds)".

### Summary

Cycle-5 ideation applied the captain's option (b): AC-1/AC-2 are now contract-shape requirements (the rule ships with its clauses, host-neutral, within ratchet; the delivery capability is defined with per-host bindings) proven by a small presence check + the existing contractlint capability-binding, with NO runtime-behavior claim; AC-3/AC-4 stay as the honestly-testable Codex hook / harmless-absence value checks. The deliverable is slimmed to the live-proven `${PLUGIN_ROOT}` hook + a small hook/config smoke + contract-shape checks; the ~419 lines of synthetic timing/reload/presence oracle machinery are marked for removal, and the live behavioral proof of the judgment rules is an explicit out-of-scope live-scenario followup that the implementation stage does not build.
