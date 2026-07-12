---
title: Codex post-compaction contract reload
status: ideation
source: Absorbed from task njr36mfyhbafy8zx9ydks8ep in another workflow; canonical handoff /tmp/first-officer-compaction-rehydration.md; captain directed repo-local absorption 2026-07-11
started: 2026-07-11T04:15:29Z
completed:
verdict:
score: 0.95
worktree:
issue:
id: c60nzb396vgf0f8a9v0sggwm
---

Conversation compaction can preserve the story of a workflow while dropping the obligation that drives its next action. In the originating v9 replay, the compacted summary mentioned Codex wait semantics, but the first officer dispatched a critical-path feedback repair and returned before it consumed the completion signal, verified the committed report, and re-ran validation. The worker remained live; the first officer lost the continuation.

This entity absorbs the design from cross-workflow task `njr36mfyhbafy8zx9ydks8ep`. That ID is provenance only. This workflow owns the new identity `c60nzb396vgf0f8a9v0sggwm` and all subsequent lifecycle.

## Problem

A launch-time bootstrap proves only initial loading. After manual or automatic compaction, Codex may continue from a summary that omits or distorts:

- the current first-officer contract and runtime adapter;
- an unresolved worker and its assignment epoch;
- the report baseline that distinguishes old evidence from new evidence;
- the next required continuation: wait, verify, route feedback, re-review, present a gate, or terminalize.

The current contracts already state the right behavior:

- `skills/first-officer/references/codex-first-officer-runtime.md` requires `wait_agent(timeout_ms)` whenever an unresolved Codex worker is the only remaining work, and requires wait to be reinstalled after operator interruption.
- `skills/feedback-rejection-flow/SKILL.md` requires critical-path follow-up to complete through `«completion-signal»`, followed by durable verification and reviewer re-run.
- `skills/first-officer/references/first-officer-shared-core.md` forbids stopping after a completed non-gated, non-terminal stage; the FO must advance and dispatch the next stage.
- `skills/first-officer/references/fo-dispatch-core.md` owns dispatch identity, assignment epochs, reuse, and completion attribution.

Generated summaries are advisory. Repo instructions can require a re-read, but they do not reconstruct live assignments. Durable workflow state proves stage and committed reports, but does not presently distinguish `await_completion` from `verify_report`, `rerun_review`, or `present_gate`.

External evidence sharpens the boundary. ECC's strategic-compaction hook uses actual transcript context size as its primary pressure signal, tool-call count only as a fallback, and rate-limits reminders by context-growth buckets. Its own guidance says to compact at logical phase boundaries rather than when a threshold fires. ECC's PreCompact summary is deliberately lossy and failure-open: it samples at most 25 recent turns and 7,000 characters, invokes an LLM, and falls back to logging when generation fails. Its SessionStart loader marks prior summaries as historical-only and skips prior-summary injection for `resume`, `clear`, and `compact` modes. These are useful risk reductions, but none proves same-session continuation or replaces a typed obligation record.

The invariant is: after dispatching or reusing a critical-path worker, the FO cannot produce a completion response until it consumes the matching completion signal, verifies the durable report, records the next continuation, and reaches a gate, terminal state, or explicit blocker.

## Live Codex 0.144.1 schema findings

The experimental schema was generated from the live `codex-cli 0.144.1` and retained at `/tmp/codex-c6-schema.SEuqNa`. It proves names and payload shapes, not runtime ordering or enforcement coverage.

### Schema-proven surfaces

- `PluginDetail` exposes plugin `hooks` as `PluginHookSummary`; the static hook configuration schema enumerates `SessionStart`, `PostCompact`, `SubagentStart`, `SubagentStop`, and `Stop`; runtime hook notifications spell those events `sessionStart`, `postCompact`, `subagentStart`, `subagentStop`, and `stop`. `HookRunSummary` can identify `source: plugin`, `sync` or `async` execution, and `context`, `stop`, `feedback`, `warning`, or `error` output kinds. Evidence: `PluginDetail`, `PluginHookSummary`, `ManagedHooksRequirements`, `v2/HookStartedNotification.json`, and `v2/HookCompletedNotification.json`.
- `item/completed` can carry a `ThreadItem` whose type is `contextCompaction`. Evidence: `v2/ItemCompletedNotification.json`.
- Deprecated `thread/compacted` remains present as a compatibility notification. Evidence: `v2/ContextCompactedNotification.json` and the method union in `codex_app_server_protocol.v2.schemas.json`.
- `thread/compact/start` accepts a `threadId`. Evidence: `v2/ThreadCompactStartParams.json`.
- `thread/inject_items` accepts a `threadId` and raw Responses API items described as appended to the thread's model-visible history. Evidence: `v2/ThreadInjectItemsParams.json`.
- `turn/completed`, `hook/started`, and `hook/completed` notifications are present. Evidence: their v2 schemas and the protocol method union.

### Claims that require the first spike

The schema does **not** prove that:

- an installed plugin can register and run the proposed `SessionStart`, `PostCompact`, `SubagentStart`, `SubagentStop`, and `Stop` configuration in the required scope;
- `SessionStart` or `PostCompact` can re-inject the fixed FO wait contract as developer context before the model's first relevant action;
- `SubagentStart` and `SubagentStop` have complete worker coverage, or that their `PLUGIN_DATA` survives and is visible across the lifecycle boundary with the assignment identity needed for attribution;
- a `Stop` handler accepts `decision: "block"`, prevents an inappropriate final stop, and delivers its fixed “reconcile agents and call `wait_agent`” directive without racing an operator turn;
- an explicit captain `do not wait` opt-out reaches the hook, suppresses that block only for its recorded scope, and cannot be confused with a missing or stale worker record;
- `thread/inject_items` called after compaction is ordered before that action;
- `preToolUse` observes or can block collaboration lifecycle calls such as spawn, follow-up, wait, or interrupt;
- duplicate hook/item/legacy notifications have a stable order or one-to-one relationship;
- a `turn/completed` watchdog may safely start another turn without racing operator input.

These remain hypotheses. Implementation must not claim the stronger Codex enforcement path until the riskiest-path spike proves them.

## Options considered

1. **Instruction-only re-read.** Add “after compaction, re-read the first-officer contract” to `AGENTS.md` or the runtime adapter. This is cheap and legible, but the summary may omit the trigger or the obligation. It cannot recover a worker identity or assignment epoch.
2. **Durable continuation fence.** Record unresolved critical-path obligations outside model context and reconcile them at turn entry. This survives summary omission and works without a compaction event. It is the minimum correctness contract and the recommended shared design.
3. **Lifecycle injection only.** Use Codex hook or app-server notifications to inject a reload instruction. This offers a precise trigger if ordering is proven, but injected prose alone does not prove that the reload or state reconciliation occurred.

Adopt option 2. Use option 3 only as a Codex enforcement adapter after the spike. Keep option 1 as a readable statement of the invariant, not the safety mechanism.

## Proposed design

### 1. Launcher-owned continuation checkpoint

The launcher creates one session-scoped checkpoint and passes its absolute path as `${SPACEDOCK_FO_CHECKPOINT}`. The default lives below the user's Spacedock state directory, never in a workflow checkout or code worktree. A launcher-resolved contract root is recorded alongside it so rehydration never searches for an installed sibling.

The checkpoint is a versioned ledger, not a prose recap. Its top level contains:

- schema version, random session nonce, host thread/session identity, and monotonic generation;
- canonical launcher and contract-root paths selected at launch;
- `compaction_epoch`, `needs_rehydrate`, the last successful rehydrate epoch, and source digests;
- zero or more obligations keyed by workflow identity, entity full ID, stage, feedback cycle, and assignment epoch.

Each obligation contains:

- canonical workflow directory, entity reference, stage, cycle, assignment epoch, and worker identity when bound;
- entity/state commit and Stage Report count observed before dispatch;
- one continuation state: `await_completion`, `verify_report`, `route_feedback`, `rerun_review`, `present_gate`, or `terminalize`;
- the next required durable boundary and the last verified boundary;
- a monotonic obligation revision.

The FO records `await_completion` before spawning or routing work, then binds the returned worker identity. A crash between those writes leaves an unbound unresolved dispatch; it never implies completion. A completion from epoch N cannot satisfy epoch N+1. Clearing requires the matching assignment epoch, a newer committed report, checklist verification, and a durably recorded successor state.

### 2. Checkpoint trust, filesystem safety, and lifecycle

The checkpoint defends against accidental corruption and path attacks, not a malicious same-UID process with unrestricted filesystem access.

- Create a private `0700` session directory and `0600` files. Reject symlinks, non-regular files, owner mismatch, path traversal, contract roots outside the launcher-selected root, unknown schema versions, and nonce/thread mismatch.
- Write through a launcher-owned API using a lock plus compare-and-swap generation. Persist an atomic same-directory temporary file, `fsync` it, rename it, and `fsync` the directory. Readers reject partial or generation-regressing records.
- Treat all strings in the checkpoint as data. Never execute a stored command or inject checkpoint text verbatim into model context. The launcher renders a fixed directive containing only validated identifiers and the checkpoint path.
- Store no credentials, prompts, transcript text, or worker output. The record contains identifiers, revisions, paths, digests, and typed states only.
- Create the record after launcher/contract-root resolution. On clean session exit with zero obligations, remove the session directory. On crash or exit with obligations, retain it for explicit resume. Resume adopts it only after nonce/thread/workflow validation; mismatched or stale records are quarantined and reported, never silently reused or deleted.
- A maintenance command may list or garbage-collect terminal/quarantined checkpoints. Age alone never marks an unresolved obligation complete.

The model-facing `spacedock dispatch rehydrate` command may request transitions, but the launcher/checkpoint service validates the state machine and performs writes. This prevents an accidental shell edit from becoming authoritative. A hostile same-user model or process remains outside the threat model and must not be implied secure.

### 3. Concurrent obligation ledger

A single scalar “current worker” is insufficient: the FO may have multiple independent workers, a queued completion for one, and a gate for another. The ledger therefore supports multiple obligations.

At each fence, the reconciler snapshots the ledger generation, roster/mailbox evidence, and workflow state, then computes typed next actions for every unresolved obligation. It commits transitions with compare-and-swap; on generation conflict it re-reads and recomputes rather than overwriting another completion.

Deterministic priority is:

1. verify matching completion evidence already received;
2. finish an in-progress feedback/re-review chain;
3. present ready gates or terminalize verified work;
4. restore waits for still-running workers;
5. dispatch newly ready independent work only after the existing critical obligations are accounted for.

One worker completion advances only its matching obligation. The fence clears globally only when every obligation is reconciled and the contract digests match the current launcher-selected sources.

### 4. Rehydration fence

`spacedock dispatch rehydrate --checkpoint "$SPACEDOCK_FO_CHECKPOINT" --json` is the proposed integration boundary. It returns authoritative sources, digest expectations, reconciliation results, and typed next actions; it never infers completion from a generated summary.

After `needs_rehydrate` is set, the FO must pass this fence before its first worker lifecycle, workflow mutation, gate resolution, merge, or user-facing completion action. Every new FO turn must also pass the fence while any obligation remains.

The fence:

1. Re-reads `skills/first-officer/SKILL.md`, all of its eager imports, and the active runtime adapter from the launcher-selected contract root.
2. Re-reads `references/fo-dispatch-core.md` for unresolved dispatches and `skills/feedback-rejection-flow/SKILL.md` for feedback/re-review continuations. Other deferred modules load only when their normal triggers fire.
3. Reads the workflow stage definition, current entity frontmatter, latest Stage Report, and `### Feedback Cycles` when applicable.
4. Verifies source digests and reconciles every assignment epoch against mailbox/roster evidence and durable report baselines.
5. Returns a typed next action: verify, route, re-review, present, terminalize, or wait. It clears `needs_rehydrate` only after the successful digest and ledger transition is committed.

If a source cannot be resolved, a digest changes during the read, the checkpoint is invalid, or evidence is ambiguous, the fence fails closed: no spawn, follow-up, state mutation, gate resolution, merge, or completion claim. It records a typed blocker and reports the exact recovery action to the captain.

With no compaction and no unresolved obligation, the command is a no-op: no roster reconciliation, contract streaming, or checkpoint write.

### 5. Computed compaction readiness

Context pressure is a suggestion signal, never permission to compact mid-obligation. The launcher computes `compaction_ready` from durable state. It is true only when every active obligation has been written to the checkpoint, every dispatched obligation has its worker identity and assignment epoch bound, every report baseline is stored, and every obligation has one typed next action. A context-size or tool-count threshold may suggest compaction only at a logical phase boundary and only when this predicate is true.

For launcher/client-initiated compaction, `compaction_ready=false` refuses the request and reports the missing durable fields; it does not synthesize them from transcript text or a generated summary. An automatic host compaction that bypasses this prevention layer still sets `needs_rehydrate` and must pass the durable fence. Thus strategic compaction reduces exposure while the ledger and fence preserve correctness when timing control fails. With zero obligations, readiness evaluation is a cheap in-memory predicate and adds no checkpoint write, contract stream, roster query, or recurring ceremony.

### 6. Codex plugin lifecycle adapter, conditional on the spike

A bundled Codex plugin is a stronger lifecycle aid than skill prose, but it is never the source of truth. The durable checkpoint and rehydration fence remain authoritative for worker state, completion attribution, and every guarded action.

The intended plugin policy is:

- `SessionStart` and `PostCompact` target re-injection of a fixed FO wait-contract developer-context block. It identifies the launcher-selected contract and tells the FO to reconcile unresolved agents and call `wait_agent`; it never copies checkpoint contents, worker output, or transcript text into context. Until a live trace proves role, visibility, and ordering, this is only a candidate context directive, not a claim that developer context arrived before a guarded action.
- `SubagentStart` and `SubagentStop` record observed worker lifecycle facts in plugin-owned `PLUGIN_DATA`, keyed by validated thread/session and assignment identity. This is an observation cache for prompt repair and stop handling, not a second ledger: it cannot satisfy, clear, or invent a durable obligation.
- `Stop` consults the unresolved state derived from the checkpoint/fence. While an obligation remains and no applicable captain opt-out is recorded, its planned result is `decision: "block"` plus the same fixed reconcile-and-`wait_agent` directive. A block asks Codex to continue reconciliation; it does not declare any worker complete, spawn a replacement, or itself start a watchdog turn.
- A captain's explicit `do not wait` is recorded as typed, scoped opt-out state with a checkpoint generation and explicit reversal. It may be mirrored into `PLUGIN_DATA`, but the checkpoint is authoritative. For matching obligations, the plugin omits the wait directive and `Stop` does not block, so it never fights that instruction; the opt-out does not clear the obligation or waive the rehydration fence for mutations, gates, merges, or completion claims.

Lifecycle handling distinguishes fresh `startup`, persisted-session `resume`, context `clear`, and same-session `compact`. Each source must set or preserve the correct checkpoint identity and fence state; none may treat a prior generated summary as live instruction. In particular, `compact` requires an independent PostCompact/turn-entry proof because ECC's analogous SessionStart path intentionally skips prior-summary injection for every non-startup source. `clear` may discard model context but cannot discard an unresolved durable obligation, and `resume` adopts retained state only through the validation rules above.

The event-capable client still deduplicates `postCompact`, `item/completed` with `contextCompaction`, and legacy `thread/compacted` into one compaction epoch. Duplicate signals mark the same epoch; they do not trigger repeated reloads. If the spike proves plugin context ordering, `PostCompact` marks `needs_rehydrate` through the checkpoint service and contributes the fixed directive. If plugin context is unavailable but injection ordering is proven, the client uses `thread/inject_items` for that directive. If `preToolUse` coverage is proven for a tool class, it blocks that class while the fence is outstanding.

Evidence is deliberately tiered:

- **Schema-proven now:** the live v2 schema names the lifecycle hooks, recognizes plugin-sourced runs, and describes hook modes and output kinds.
- **Fixture-proven after implementation:** plugin configuration rendering, fixed-context redaction, epoch de-duplication, `PLUGIN_DATA` record serialization, and the pure unresolved/opt-out predicate can be tested without claiming host execution.
- **Live-plugin-spike required:** plugin discovery and registration, callback delivery and order, developer-context role/visibility, `PLUGIN_DATA` existence/lifetime and worker coverage, the exact `decision: "block"` wire behavior, the effect on a final stop, and opt-out propagation. A failed or ambiguous claim disables that enforcement claim; it cannot weaken the fence.

No hook is authoritative. Hook failure, timeout, malformed output, absent notification, late context, stale plugin data, or uncovered collaboration call leaves `needs_rehydrate` set. The launcher/client-side fence remains the authority. Where the client cannot intercept a call, the design claims only the host-neutral turn-entry contract until a live probe proves stronger coverage.

### 7. Bounded continuation

Automatic enforcement must not create an injection or turn-start loop.

- At most one directive injection and one watchdog-started continuation turn are allowed per `(compaction_epoch, checkpoint_generation)` pair.
- A `Stop` `decision: "block"` is one stop refusal, not a new watchdog turn or durable progress. Repeated stop events cannot recursively inject or restart the FO.
- A watchdog continuation is allowed only when the previous turn completed while an obligation required a non-wait next action and no operator turn is already pending.
- Durable progress means a generation advance tied to a verified source read, report boundary, or typed continuation transition. Assistant prose and repeated hook notifications are not progress.
- If the automatic continuation makes no durable progress, fails the fence, or returns to the same state, the watchdog stops. It preserves the obligation, records the blocker, and waits for an external turn or operator recovery; it never recursively injects.
- Foreground waits remain individually bounded by the host timeout. A timeout leaves the same worker and epoch unresolved. The next externally entered idle turn reinstalls wait; timeout never causes redispatch or completion.

This preserves the keep-moving invariant without an unbounded autonomous loop.

## Acceptance criteria

### Offline

**AC-1: Authoritative post-compaction reload.** Before the first post-compaction lifecycle, mutation, gate, merge, or completion action, evidence contains a successful rehydrate epoch whose source digests equal the current launcher-selected first-officer entry skill, eager imports, active runtime adapter, applicable dispatch/feedback modules, workflow stage definition, and entity evidence.

**AC-2: Critical-path continuity.** Across manual and automatic compaction, a summary omitting every wait instruction, queued completion, and operator interruption, the replay records zero premature shutdowns, replacement dispatches, completion-only final responses, or state advances before the matching committed Stage Report. After the report arrives, the entity reaches its next gate or terminal continuation in the same drive.

**AC-3: Durable recovery without summary state.** Reconstructing from only the checkpoint, launcher-selected contracts, workflow checkout, roster fixture, and mailbox fixture recovers the same entity, stage, cycle, assignment epoch, report baseline, continuation state, and next typed action as before compaction. The test replaces any generated summary with stale executable-looking prose and proves that no action is derived from it.

**AC-4: Event detection and fallback.** Event-capable fixtures deduplicate `postCompact`, `contextCompaction`, and legacy `thread/compacted` into one epoch. Eventless and failed-hook fixtures still fence the first subsequent FO turn and every active turn before guarded actions.

**AC-5: Completion attribution.** A stale completion from assignment epoch N cannot satisfy N+1; a matching completion cannot clear its obligation until a newer committed Stage Report passes checklist verification. Stale, duplicate, uncommitted, malformed, and cross-entity evidence remain unresolved.

**AC-6: Concurrent obligations.** Two or more workers can complete in either order without lost updates, cross-attribution, or a global fence clearing early. A simulated generation conflict retries from the new ledger rather than overwriting it.

**AC-7: Checkpoint safety and lifecycle.** Symlink, owner, permission, nonce/thread, schema, traversal, partial-write, generation-regression, and stale-resume fixtures fail closed. Clean empty exit removes the checkpoint; crash with obligations preserves it; resume requires explicit validated adoption.

**AC-8: Hook trust and bounded continuation.** Hook absence, failure, timeout, duplicate events, late injection, and uncovered tools never clear the fence. Each epoch/generation produces at most one injection and one watchdog continuation; no-progress stops with the obligation intact and a typed blocker.

**AC-9: Wait and feedback guarantees.** Timeout or operator interruption returns control without failing, closing, or redispatching the worker. The next idle action reinstalls wait; a rejected fix proceeds through repair completion, durable verification, reviewer re-run, and normal gate flow.

**AC-10: Inactive cost.** With no observed compaction requiring reload and no unresolved obligation, readiness and rehydration perform no roster reconciliation, contract streaming, checkpoint write, or repeated suggestion work beyond the existing pressure-bucket check. A trace fixture asserts zero such operations.

**AC-11: Compaction readiness prevents unsafe voluntary compaction.** `compaction_ready` is false whenever any active obligation lacks durable persistence, bound dispatch identity/epoch, stored report baseline, or typed next action; launcher/client-initiated compaction is refused with the exact missing fields. Readiness becomes true only when all obligations satisfy that predicate, and compaction is allowed only when readiness is true and the current point is a declared logical boundary. Table tests toggle every term independently, and an automatic-compaction fixture proves that bypassing this prevention layer still leaves `needs_rehydrate` set and the fence mandatory.

**AC-12: Distinct lifecycle sources preserve continuity.** Fresh `startup`, validated `resume`, context `clear`, and same-session `compact` each produce the expected checkpoint identity, `needs_rehydrate` transition, and typed next action without executing prior-summary text. Mode-specific fixture traces cover all four sources, and a live Codex spike separately proves same-session `compact` ordering rather than inferring it from startup behavior.

### Interactive

**AC-13: Live Codex demonstration.** During a delayed worker repair, CL triggers `/compact`, interrupts one foreground wait, and sends an unrelated status question. The FO answers from durable state, restores wait, consumes the final notification, verifies the committed report, and continues to re-review or the next gate without a reminder. Evidence includes app-server notifications, plugin hook and `PLUGIN_DATA` traces, checkpoint generations, workflow git log, and clean state.

**AC-14: Plugin lifecycle assistance is evidence-gated and opt-out-safe.** Before the plugin claims enforcement on a target Codex version, a live trace establishes plugin discovery plus all five callback traces; `SessionStart`/`PostCompact` developer-context delivery before the first guarded action; two-worker `SubagentStart`/`SubagentStop` accounting in `PLUGIN_DATA` under duplicate and out-of-order events; and `Stop` `decision: "block"` acceptance for an unresolved fixture obligation while a resolved obligation may stop. The same trace proves that a recorded explicit captain `do not wait` opt-out suppresses the matching wait directive and stop block only. Absent, late, malformed, stale, or unproven plugin behavior leaves the durable fence as the sole guard and never treats work as complete.

## Test plan

The riskiest unproved mechanism remains Codex ordering and coverage. Invalidate that mechanism before implementing the full protocol.

1. Build a temporary target-Codex plugin/client probe with `SessionStart`, `PostCompact`, `SubagentStart`, `SubagentStop`, `Stop`, a context marker, `thread/inject_items`, and a blocking `preToolUse` hook. First fixture-test the rendered hook map, fixed context-directive text, `PLUGIN_DATA` serialization, duplicate-epoch handling, and the pure unresolved/opt-out predicate; those fixtures do not certify host execution. Then confirm plugin discovery and capture distinct `startup`, `resume`, `clear`, and `compact` traces, retaining callback stdin/env/stdout, `item/completed`, `hook/started`, `hook/completed`, `PLUGIN_DATA` transitions, context/injection acknowledgements, `turn/completed`, the first post-compaction tool event, and the final stop decision. Start two delayed subagents, exercise duplicate and out-of-order start/stop events, invoke `thread/compact/start`, and test manual and automatic compaction, unresolved and resolved stop, and an explicit captain `do not wait` opt-out. Fail the stronger path if any callback is absent, context is absent or late, same-session compact is inferred only from startup, duplicate events create multiple epochs, `PLUGIN_DATA` is unavailable or cannot attribute both fixture workers, the stop result is not demonstrably blocking or a resolved obligation cannot stop, the opt-out still injects or blocks, operator input races the watchdog, or any collaboration lifecycle call bypasses the claimed hook. Record exact coverage; do not generalize from shell-tool coverage.
2. Add checkpoint serialization, private-path validation, atomic-write, compare-and-swap, incomplete-bind, digest, lifecycle, and stale-resume tests under a proposed `internal/rehydrate/` package. Add table tests that toggle each `compaction_ready` term independently, require a logical boundary, report missing fields, refuse unsafe launcher/client compaction, and prove automatic bypass still requires the fence.
3. Add concurrent-ledger and state-machine tests: out-of-order completions, two workflows, stale epochs, generation conflicts, matching and malformed reports, hook failure, no-progress watchdog, bounded wait timeout, and all four lifecycle sources. Inject stale executable-looking summary text and assert that transitions depend only on typed durable inputs.
4. Add a split-root replay fixture under `internal/ensigncycle/testdata/compaction-rehydration/`: validation rejects to implementation, a delayed fixer commits a report, and validation must run again before the gate.
5. Add `internal/ensigncycle/compaction_rehydration_test.go`. Replace context with a summary that omits the obligation, deliver completion before and after the next turn, interrupt wait, verify the committed report, and assert the next gate or terminal state.
6. Extend the Codex idle-notification integration region to prove wait restoration, durable file verification, eventless turn-entry fallback, and zero inactive roster/contract/checkpoint operations. Assertions must use tool/event traces and workflow state, not transcript phrasing.
7. Run AC-13 and AC-14 live and retain app-server JSONL, plugin hook log, `PLUGIN_DATA` audit, checkpoint/readiness transitions, and the temp workflow git log as CI artifacts.

## Out of scope

Implementing the protocol in ideation; changing Codex compaction algorithms; persisting conversation history; treating generated summaries as authoritative; changing Claude or Pi behavior; weakening existing wait, feedback, gate, or merge contracts; using workflow mods as context-lifecycle hooks; defending against a malicious same-UID process with unrestricted filesystem access; and fixing the separate `spacedock codex -- resume` bootstrap-prompt bug.

## Evidence reviewed for this revision

- ECC `strategic-compact`, `suggest-compact.js`, `pre-compact.js`, `session-start.js`, and `llm-summary.js` at `affaan-m/ECC` main commit `40927950c49f6e742d341e20ff7b9b7e1e7bfff5`.

## Stage Report: ideation

- DONE: Replace the thin c6 body with the full ideation artifact at /tmp/first-officer-compaction-rehydration.md, preserving c6's repo-local full ID and recording njr36mfyhbafy8zx9ydks8ep as source provenance rather than identity.
  Absorbed the original problem, alternatives, durable continuation fence, Codex adapter, acceptance evidence, and riskiest-path-first test plan under `c60nzb396vgf0f8a9v0sggwm`; `njr36mfyhbafy8zx9ydks8ep` appears only as source provenance.
- DONE: Validate every Codex app-server claim against the live 0.144.1 experimental schema at /tmp/codex-c6-schema.SEuqNa and current repo contracts; distinguish schema-proven names from ordering/coverage claims that still require a spike.
  The design records the schema-proven hook events, modes, output kinds, compaction item and compatibility notification, compact/injection methods, and completion notifications. Registration, ordering, model visibility, collaboration interception, deduplication, and watchdog safety are explicitly hypotheses for the first spike.
- DONE: Close the design gaps for checkpoint security/lifecycle, concurrent obligations, hook trust/failure fallback, and bounded continuation; keep the riskiest-path-first proof and avoid expanding into unsupported cross-runtime implementation.
  Added a private atomic checkpoint lifecycle and threat boundary, a compare-and-swap multi-obligation ledger, fail-closed hook fallback, and one-injection/one-continuation bounds. The first test still attempts to invalidate Codex ordering and coverage before product implementation; Claude and Pi changes remain out of scope.

### Summary

The minimum reliable contract is a durable continuation fence, not a better summary. Every active first-officer turn reconstructs unresolved obligations from a launcher-owned typed ledger and authoritative checkout sources. Codex 0.144.1 exposes promising hook, compaction, injection, and completion names, but the schema does not prove their ordering or tool coverage; the first spike must do so before Spacedock claims enforced post-compaction rehydration. The revised design also prevents concurrent-worker loss, fails closed on corrupt or failed checkpoints/hooks, and bounds automatic continuation so recovery cannot become an infinite turn loop.

## Stage Report: ideation

- DONE: Amend c6's design with the missing Codex plugin lifecycle-hook note: SessionStart/PostCompact re-inject the FO wait contract as developer context; SubagentStart/SubagentStop track workers in PLUGIN_DATA; Stop returns decision block while unresolved workers exist; an explicit captain opt-out honors do not wait.
  Added the named lifecycle policy, fixed reconcile-and-`wait_agent` directive, non-authoritative `PLUGIN_DATA` cache, planned `decision: "block"`, and typed captain opt-out without allowing any of them to clear a durable obligation.
- DONE: Keep the current durable rehydration fence as authority and distinguish each lifecycle-hook claim that needs a live plugin spike from schema- or fixture-proven behavior; do not overclaim ordering or interception.
  Re-read the retained Codex 0.144.1 schema: `PluginDetail`/`PluginHookSummary` and the hook enums prove the static surface, while registration, callback delivery/order, developer-context visibility, `PLUGIN_DATA`, stop semantics, and opt-out propagation remain explicitly live-spike claims.
- DONE: Update the relevant acceptance/test plan and append an ideation stage report; make no product implementation changes.
  Added AC-12 and a fixture-versus-live probe plan covering every named hook and the opt-out; only this split-root entity body changed, and `git diff --check` passes.

### Summary

The plugin is now designed as a lifecycle assist that reinjects the FO wait contract, observes worker events, and can prevent an accidental stop while work remains. The checkpoint-backed rehydration fence remains the only authority; no hook behavior is claimed until the target Codex plugin spike captures its real registration, order, data lifetime, and stop semantics. An explicit captain `do not wait` opt-out suppresses only the matching wait directive and stop block, and never waives durable verification.

## Stage Report: ideation

- DONE: Add a computed `compaction_ready` prevention layer without weakening the durable rehydration fence.
  The predicate now requires every active obligation to be durable, every dispatch identity and assignment epoch to be bound, every report baseline to be stored, and every obligation to have a typed next action. Voluntary launcher/client compaction is refused with missing fields until readiness is true and a logical phase boundary is reached; automatic host compaction still sets `needs_rehydrate` and must pass the existing fence.
- DONE: Turn the applicable ECC findings into evidence-backed acceptance criteria and realistic lifecycle tests.
  Recorded transcript context size as the primary pressure signal, tool count as fallback, rate-limited non-blocking suggestions, lossy/failure-open PreCompact summaries, and historical-only replay. AC-10 through AC-12 and test steps 1 through 3 now exercise readiness terms, near-zero inactive cost, stale executable-looking summaries, and distinct startup/resume/clear/compact behavior. The first live hook-ordering spike remains blocking for strong Codex enforcement claims.
- DONE: Keep c6 focused on post-compaction continuity and exclude unrelated workflow-efficiency ideas.
  The entity remains limited to ECC-informed compaction prevention plus the authoritative typed ledger and rehydration fence; unrelated workflow-efficiency mechanisms and evaluation changes are excluded.

### Summary

Strategic compaction now acts as a prevention layer: real context pressure may suggest a logical boundary, but computed durable readiness decides whether voluntary compaction is safe. Generated summaries remain historical and non-authoritative. Correctness still comes only from the typed obligation ledger and rehydration fence, with startup, resume, clear, and same-session compact proved separately and the live Codex hook-ordering spike still blocking any strong enforcement claim.
