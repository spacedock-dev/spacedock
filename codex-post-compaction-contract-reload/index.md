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

There is also an enforcement boundary the design must state honestly. If automatic compaction emits no observable lifecycle/item notification, no hook runs, and the checkpoint has zero obligations, durable state contains no fact from which Spacedock can infer that compaction occurred. A later raw collaboration call or user-facing final response is not intercepted merely because a skill says to rehydrate. Correctness claims therefore apply only when compaction is observed, an obligation was already durable, or the guarded action travels through a concrete Spacedock/client interceptor whose coverage has been proved. Eventless zero-obligation compaction followed by an uncovered host action remains an explicit unprotected path, not a fallback success.

The guard itself must preserve the dispatch protocol it protects. Because `await_completion` is written before spawn, “any unresolved obligation blocks dispatch” deadlocks the required first spawn. A guard that checks state and releases its lock before recording intent permits compaction or another ledger write to invalidate the decision. Holding a filesystem lock across an external effect removes that race but does not make the effect and successor ledger write crash-atomic: the process can die after a worker, mutation, archive, or merge becomes visible and before consuming its permit. Finally, an unbounded or indefinitely lock-blocked “small read” is not an inactive-cost contract. The design therefore needs typed action authorization, a durable pre-effect intent, restart reconciliation without blind replay, and quantitative checkpoint/lock limits.

External evidence sharpens the boundary. ECC's strategic-compaction hook uses actual transcript context size as its primary pressure signal, tool-call count only as a fallback, and rate-limits reminders by context-growth buckets. Its own guidance says to compact at logical phase boundaries rather than when a threshold fires. ECC's PreCompact summary is deliberately lossy and failure-open: it samples at most 25 recent turns and 7,000 characters, invokes an LLM, and falls back to logging when generation fails. Its SessionStart loader marks prior summaries as historical-only and skips prior-summary injection for `resume`, `clear`, and `compact` modes. These are useful risk reductions, but none proves same-session continuation or replaces a typed obligation record.

The invariant is: after dispatching or reusing a critical-path worker, the FO cannot produce a completion response until it consumes the matching completion signal, verifies the durable report, records the next continuation, and reaches a gate, terminal state, or explicit blocker.

## Codex 0.144.1 schema findings

The initial experimental schema was generated from live `codex-cli 0.144.1` with `codex app-server generate-json-schema --experimental --out <dir>` and inspected at `/tmp/codex-c6-schema.SEuqNa`. The temporary directory is discovery evidence only. Implementation must add a versioned extraction fixture containing the exact command, `codex --version`, selected schema filenames and SHA-256 digests, and the extracted event/method/type assertions below. A scripted test regenerates into a temporary directory and compares the extraction fixture; a version or digest mismatch requires explicit fixture review. The fixture makes names and payload-shape evidence durable and reproducible, but still cannot prove runtime ordering or enforcement coverage.

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
2. **Durable continuation fence.** Record unresolved critical-path obligations outside model context and reconcile them at guarded action/turn entry. This survives summary omission and works without a compaction event for obligations that were already recorded. It is the minimum correctness contract and the recommended shared design, but it cannot infer an unobserved zero-obligation compaction.
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
- one typed next action with exact action class and validated target digest; an unbound pre-spawn `await_completion` obligation uses `spawn_initial`, while a bound one uses `wait`;
- the next required durable boundary and the last verified boundary;
- a monotonic obligation revision.

The FO records `await_completion` before spawning or routing work, with an unbound worker, a new assignment epoch, the report baseline, and typed next action `spawn_initial`. The reconciler may authorize exactly that initial unbound transition; after the spawn returns, the launcher binds the returned worker identity and changes the typed next action to `wait`. If the process dies before completing that transition, the crashed process records nothing further: the already-fsynced executing permit is the recovery fact. A later launcher detects it as indeterminate, reconciles external evidence, and never blindly redispatches. A completion from epoch N cannot satisfy epoch N+1. Clearing requires the matching assignment epoch, a newer committed report, checklist verification, and a durably recorded successor state.

### 2. Checkpoint trust, filesystem safety, and lifecycle

The checkpoint defends against accidental corruption and path attacks, not a malicious same-UID process with unrestricted filesystem access.

- Create a private `0700` session directory and `0600` files. Reject symlinks, non-regular files, owner mismatch, path traversal, contract roots outside the launcher-selected root, unknown schema versions, nonce/thread mismatch, and checkpoint files larger than 1 MiB.
- Prefix the checkpoint with a checksummed fixed 4 KiB guard header containing magic, schema, total/body length, generation, `needs_rehydrate`, compaction/rehydrate epochs, nonce and host identity, bounded canonical launcher/contract-root paths, aggregate source digest, obligation count, outstanding-permit count, and body digest. Header paths are capped at 1,024 bytes each. Cap the whole checkpoint at 1 MiB including the header (body at most 1 MiB minus 4 KiB), obligations at 128, outstanding permits at one, and any body string at 4 KiB. When both counts are zero, the atomic writer compacts the body to length zero, so the header is the complete inactive checkpoint. The inactive path performs one `stat` plus one `pread` of at most 4 KiB, decodes zero body records, and uses at most 16 KiB transient allocation. Bad magic/checksum/length/count/digest, truncation, excess size, or zero counts with a nonzero body fails closed.
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

An unresolved obligation is not a blanket prohibition. The reconciler compares the requested action with the deterministic priority order and the obligation's typed next action. If they match, it may issue a single-use action permit containing a random permit nonce, action class, exact target/argument digest, checkpoint generation, obligation key and revision, assignment epoch, and permitted successor transition. Issuance is itself an atomic checkpoint mutation: it stores the permit, increments generation, and stamps the permit with that resulting generation. A different action class/target, a changed generation or obligation revision, a second use, or a revoked permit is rejected before side effects. Independent dispatch is permitted only after higher-priority obligations have typed actions that are already executing, waiting, or otherwise accounted for.

### 4. Rehydration fence

`spacedock dispatch rehydrate --checkpoint "$SPACEDOCK_FO_CHECKPOINT" --json` is the proposed integration boundary. It returns authoritative sources, digest expectations, reconciliation results, and typed next actions; it never infers completion from a generated summary.

After `needs_rehydrate` is set, the FO must pass this fence before its first guarded action. The launcher/checkpoint service concretely guards Spacedock-owned action paths—initial dispatch assembly, workflow mutation/archive, and merge—by checking the checkpoint before the command executes. An event-capable client may guard collaboration lifecycle calls and final-stop handling only after the live spike proves interception coverage for those classes. While an obligation remains, the host-neutral contract still requires a fence on every entered FO turn, but this is behavioral rather than an enforcement claim for raw host actions that bypass the guarded command/client paths.

Concretely, the guarded commands read `${SPACEDOCK_FO_CHECKPOINT}` and ask one launcher-owned action gateway for a permit before emitting a dispatch spec or performing any filesystem/git mutation. The initial action classes are `spawn_initial`, `spawn_replacement`, `state_mutation`, `archive`, and `merge`. `needs_rehydrate=true` always refuses issuance until rehydration succeeds. Otherwise unresolved obligations are reconciled: the gateway issues a permit only when the requested action and target digest equal the currently authorized typed next action. This explicitly allows the pre-recorded unbound `spawn_initial` transition instead of deadlocking on its own `await_completion` obligation.

Permit issuance and execution are separate. The action gateway acquires the checkpoint lock, rechecks permit nonce, action/target digest, generation, obligation revision, assignment epoch, issued status, and `needs_rehydrate`, then changes the permit to `executing`, records the action-specific reconciliation key and expected before/after evidence, increments generation, and atomically persists plus `fsync`s the checkpoint file and directory. That durable executing transition is the action's linearization point. The gateway then releases the lock before invoking the external effect. A concurrent context-loss or ledger writer therefore serializes before the executing commit and stales/refuses the permit, or after the action has durably entered recovery-visible execution.

If the effect returns, the gateway reacquires the lock, verifies the same executing permit and the observed result, commits the typed successor, marks the permit consumed, increments generation, and fsyncs before returning. Only an error that conclusively occurred before any external visibility may consume the permit as no-effect and authorize reconciliation toward retry. A timeout, lost response, partial result, or any other ambiguous error keeps the executing intent and enters the same indeterminate reconciliation path as a crash. A crash cannot perform any post-effect write; after restart, the persisted `executing` permit is interpreted as `indeterminate`, blocks all conflicting actions, and is never replayed merely because the successor is absent. Raw host calls that the gateway cannot own remain outside the enforcement claim.

Action-specific reconciliation is normative:

- **`spawn_initial`:** the permit nonce/assignment epoch is the stable dispatch key passed through the proved client gateway. Exactly one durable spawn acknowledgement plus matching roster/mailbox identity binds that worker and commits `wait` without spawning again. A gateway/idempotency query that conclusively proves the key was never accepted may revoke the old permit and authorize a new one. Missing, multiple, mismatched, or unqueryable evidence records a typed `indeterminate_spawn` blocker; absence alone is never permission to redispatch.
- **`state_mutation`:** the intent stores canonical path, preimage hash, exact postimage hash/value, and any expected path-scoped commit marker. Exact postimage/commit evidence commits the successor without rewriting. Exact preimage plus proof that no matching commit/effect exists may authorize a new compare-and-set permit. Mixed fields, an unexpected commit, or any third state records `indeterminate_state_mutation`.
- **`archive`:** the intent stores source/destination paths, source hash, expected destination hash, and commit marker. Source absent plus exact destination/commit evidence commits the successor without moving again. Exact source present plus destination absent proves no effect and may authorize a new permit. Both paths, neither path without matching commit evidence, or hash/commit mismatch records `indeterminate_archive`.
- **`merge`:** the intent stores repository/ref identity, pre-merge head, expected merged object or PR identity, and terminal-state marker. Durable git/remote/PR evidence that the expected merge landed exactly once plus matching terminal state commits the successor without merging again. Exact pre-merge state plus an authoritative remote/PR query proving no merge may authorize a new permit. Partial ancestry, changed heads, ambiguous API results, or mismatched terminal state records `indeterminate_merge`.

The inactive cost claim applies only to the uncontended, quiescent guard evaluation when no action is requested: acquire the checkpoint lock with a 250 ms maximum, read only the 4 KiB header, and release. Lock timeout returns typed `checkpoint_busy` and authorizes no action; no read/allocation latency bound is claimed for the contended path. Every requested external-effect action, even with zero obligations, uses a persisted permit and executing intent; there is no non-persisted side-effect fast path. Any active count/flag requires the bounded full decode and typed reconciliation. Oversized, corrupt, inconsistent, or generation-regressing input fails closed before permit issuance.

The fence:

1. Re-reads `skills/first-officer/SKILL.md`, all of its eager imports, and the active runtime adapter from the launcher-selected contract root.
2. Re-reads `references/fo-dispatch-core.md` for unresolved dispatches and `skills/feedback-rejection-flow/SKILL.md` for feedback/re-review continuations. Other deferred modules load only when their normal triggers fire.
3. Reads the workflow stage definition, current entity frontmatter, latest Stage Report, and `### Feedback Cycles` when applicable.
4. Verifies source digests and reconciles every assignment epoch against mailbox/roster evidence and durable report baselines.
5. Returns a typed next action: spawn initial/replacement, verify, route, re-review, present, terminalize, or wait. It clears `needs_rehydrate` only after the successful digest and ledger transition is committed.

If a source cannot be resolved, a digest changes during the read, the checkpoint is invalid, or evidence is ambiguous, the fence fails closed: no spawn, follow-up, state mutation, gate resolution, merge, or completion claim. It records a typed blocker and reports the exact recovery action to the captain.

With no compaction and no unresolved obligation, the command is a no-op: no roster reconciliation, contract streaming, or checkpoint write.

### 5. Computed compaction readiness

Context pressure is a suggestion signal, never permission to compact mid-obligation. The launcher computes `compaction_ready` from durable state. It is true only when every active obligation has been written to the checkpoint, every dispatched obligation has its worker identity and assignment epoch bound, every report baseline is stored, every obligation has one typed next action, and no permit is issued or executing. A context-size or tool-count threshold may suggest compaction only at a logical phase boundary and only when this predicate is true.

For launcher/client-initiated compaction, `compaction_ready=false` refuses the request and reports the missing durable fields; it does not synthesize them from transcript text or a generated summary. An automatic host compaction that bypasses this prevention layer sets `needs_rehydrate` only when a lifecycle/item signal is observed; an already-active obligation independently keeps Spacedock-owned guarded actions fenced. If neither fact exists, no automatic-compaction guarantee is claimed. Thus strategic compaction reduces exposure while the ledger and fence preserve correctness on detectable/interceptable paths. With zero obligations and no observed context-loss event, readiness evaluation is a cheap in-memory predicate and adds no checkpoint write, contract stream, roster query, or recurring ceremony.

### 6. Codex plugin lifecycle adapter, conditional on the spike

A bundled Codex plugin is a stronger lifecycle aid than skill prose, but it is never the source of truth. The durable checkpoint and rehydration fence remain authoritative for worker state, completion attribution, and every guarded action.

The intended plugin policy is:

- `SessionStart` and `PostCompact` target re-injection of a fixed FO wait-contract developer-context block. It identifies the launcher-selected contract and tells the FO to reconcile unresolved agents and call `wait_agent`; it never copies checkpoint contents, worker output, or transcript text into context. Until a live trace proves role, visibility, and ordering, this is only a candidate context directive, not a claim that developer context arrived before a guarded action.
- `SubagentStart` and `SubagentStop` record observed worker lifecycle facts in plugin-owned `PLUGIN_DATA`, keyed by validated thread/session and assignment identity. This is an observation cache for prompt repair and stop handling, not a second ledger: it cannot satisfy, clear, or invent a durable obligation.
- `Stop` consults the unresolved state derived from the checkpoint/fence. While an obligation remains and no applicable captain opt-out is recorded, its planned result is `decision: "block"` plus the same fixed reconcile-and-`wait_agent` directive. A block asks Codex to continue reconciliation; it does not declare any worker complete, spawn a replacement, or itself start a watchdog turn.
- A captain's explicit `do not wait` is recorded as typed, scoped opt-out state with a checkpoint generation and explicit reversal. It may be mirrored into `PLUGIN_DATA`, but the checkpoint is authoritative. For matching obligations, the plugin omits the wait directive and `Stop` does not block, so it never fights that instruction; the opt-out does not clear the obligation or waive the rehydration fence for mutations, gates, merges, or completion claims.

Lifecycle handling distinguishes fresh `startup`, persisted-session `resume`, context `clear`, and same-session `compact`. Each source must set or preserve the correct checkpoint identity and fence state; none may treat a prior generated summary as live instruction. In particular, `compact` requires an independent PostCompact/turn-entry proof because ECC's analogous SessionStart path intentionally skips prior-summary injection for every non-startup source. `clear` may discard model context but cannot discard an unresolved durable obligation, and `resume` adopts retained state only through the validation rules above.

The transition contract is normative:

| Source | Accepted input and identity | Epoch transition | `needs_rehydrate` before a guarded action | After successful fence | Failure behavior |
| --- | --- | --- | --- | --- | --- |
| `startup` | Launcher-selected contract root plus a new host thread/session identity; create a fresh nonce and checkpoint, never adopt a prior summary or checkpoint implicitly. | Initialize `compaction_epoch=0`. | Initial bootstrap must finish before dispatch; after bootstrap, false with zero obligations. | Record current source digests; remain at epoch 0. | Abort launch/dispatch if identity, contract root, or bootstrap source validation fails. |
| `resume` | Explicit retained checkpoint whose nonce, host resume identity, workflow identities, schema, and ownership validate; summaries are ignored for adoption. | Preserve the retained epoch. | True when retained `needs_rehydrate` is true or any obligation remains; otherwise false. | Record the adopted host identity and successful rehydrate epoch; clear the flag only after reconciliation commits. Active obligations still require the next entered-turn fence contract. | Quarantine mismatch/stale/corrupt state and refuse guarded Spacedock actions; never create a replacement obligation from summary text. |
| `clear` | Current validated live checkpoint and matching host identity from an observed clear lifecycle source. | Preserve `compaction_epoch`; clear is a context-loss boundary, not a compaction event. | Set true before returning control, including when there are zero obligations. | Verify sources and obligations, then clear the flag at the same epoch. | Leave true and fail closed on Spacedock-owned guarded actions; uncovered raw host actions remain outside the enforcement claim. |
| `compact` | Current validated live checkpoint plus one observed `PostCompact`, `contextCompaction`, or legacy compact signal for the matching host identity. | Increment exactly once after de-duplicating signals for the same compaction. | Set true before returning control, including when there are zero obligations. | Verify sources and obligations, record the successful new epoch, then clear the flag. | Leave true and fail closed on Spacedock-owned guarded actions; absent/unobservable signals cause no transition and no eventless guarantee. |

The event-capable client still deduplicates `postCompact`, `item/completed` with `contextCompaction`, and legacy `thread/compacted` into one compaction epoch. Duplicate signals mark the same epoch; they do not trigger repeated reloads. If the spike proves plugin context ordering, `PostCompact` marks `needs_rehydrate` through the checkpoint service and contributes the fixed directive. If plugin context is unavailable but injection ordering is proven, the client uses `thread/inject_items` for that directive. If `preToolUse` coverage is proven for a tool class, it blocks that class while the fence is outstanding.

Evidence is deliberately tiered:

- **Schema-proven now:** the live v2 schema names the lifecycle hooks, recognizes plugin-sourced runs, and describes hook modes and output kinds.
- **Fixture-proven after implementation:** plugin configuration rendering, fixed-context redaction, epoch de-duplication, `PLUGIN_DATA` record serialization, and the pure unresolved/opt-out predicate can be tested without claiming host execution.
- **Live-plugin-spike required:** plugin discovery and registration, callback delivery and order, developer-context role/visibility, `PLUGIN_DATA` existence/lifetime and worker coverage, the exact `decision: "block"` wire behavior, the effect on a final stop, and opt-out propagation. A failed or ambiguous claim disables that enforcement claim; it cannot weaken the fence.

No hook is authoritative. Hook failure, timeout, malformed output, late context, stale plugin data, or an uncovered collaboration call cannot clear an already-set `needs_rehydrate` flag or durable obligation. An absent compaction notification may provide no trigger at all when the ledger is empty; that path is explicitly unsupported. The launcher/checkpoint fence remains authoritative for Spacedock-owned guarded commands. Where the client cannot intercept a raw host call, the design claims only the host-neutral behavioral contract until a live probe proves stronger coverage.

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

**AC-1: Authoritative post-compaction reload.** On every observed-compaction or already-obligated path, before the first Spacedock-owned guarded dispatch, mutation/archive, or merge action, evidence contains a successful rehydrate epoch whose source digests equal the current launcher-selected first-officer entry skill, eager imports, active runtime adapter, applicable dispatch/feedback modules, workflow stage definition, and entity evidence; the action then executes only through a matching generation/revision-bound permit. Collaboration/final-stop actions join this criterion only for tool classes whose client interception the live spike proves.

**AC-2: Critical-path continuity.** Across launcher/client compaction, observed automatic compaction, and automatic compaction while an obligation is already durable, a summary omitting every wait instruction, queued completion, and operator interruption, the replay records zero premature shutdowns, replacement dispatches, completion-only final responses, or state advances within the proved interceptor coverage before the matching committed Stage Report. After the report arrives, the entity reaches its next gate or terminal continuation in the same drive. Unobserved zero-obligation automatic compaction followed by an uncovered raw host action is explicitly excluded.

**AC-3: Durable recovery without summary state.** Reconstructing from only the checkpoint, launcher-selected contracts, workflow checkout, roster fixture, and mailbox fixture recovers the same entity, stage, cycle, assignment epoch, report baseline, continuation state, and next typed action as before compaction. The test replaces any generated summary with stale executable-looking prose and proves that no action is derived from it.

**AC-4: Detection and interceptor coverage.** Event-capable fixtures deduplicate `postCompact`, `contextCompaction`, and legacy `thread/compacted` into one epoch and set `needs_rehydrate` before returning control. Failed-hook fixtures with an observed event remain fenced. Eventless fixtures with an existing obligation are reconciled from the durable ledger: only its exact typed next action can receive a permit; wrong actions are refused without relying on summary prose. A zero-obligation eventless fixture proves the honest limit: no epoch/flag transition occurs, an uncovered raw host action is not claimed safe, and the result is reported as unsupported rather than passing a fallback assertion.

**AC-5: Completion attribution.** A stale completion from assignment epoch N cannot satisfy N+1; a matching completion cannot clear its obligation until a newer committed Stage Report passes checklist verification. Stale, duplicate, uncommitted, malformed, and cross-entity evidence remain unresolved.

**AC-6: Concurrent obligations.** Two or more workers can complete in either order without lost updates, cross-attribution, or a global fence clearing early. A simulated generation conflict retries from the new ledger rather than overwriting it, and invalidates any permit issued against the prior generation.

**AC-7: Checkpoint safety and lifecycle.** Symlink, owner, permission, nonce/thread, schema, traversal, partial-write, generation-regression, stale-resume, bad header checksum/digest, excess count/string length, and checkpoint size over 1 MiB fixtures fail closed. Clean empty exit removes the checkpoint; crash with obligations or an executing permit preserves it; resume requires explicit validated adoption and converts executing permits to reconciliation-required indeterminate state without replay.

**AC-8: Hook trust and bounded continuation.** Hook absence, failure, timeout, duplicate events, late injection, and uncovered tools never clear an already-set fence or durable obligation. Each epoch/generation produces at most one injection and one watchdog continuation; no-progress stops with the obligation intact and a typed blocker. No assertion infers compaction solely from hook absence.

**AC-9: Wait and feedback guarantees.** The pre-spawn protocol records an unbound `await_completion` obligation with typed `spawn_initial`, executes one permitted spawn, binds its returned worker identity, and advances to `wait`. Timeout or operator interruption returns control without failing, closing, or redispatching the worker. The next idle action reinstalls wait; a rejected fix proceeds through repair completion, durable verification, reviewer re-run, and normal gate flow.

**AC-10: Inactive cost.** On the uncontended, quiescent no-action probe with no observed compaction, unresolved obligation, or outstanding permit, guard evaluation acquires/releases the lock within its 250 ms ceiling, performs one `stat` and one `pread` of at most 4 KiB, decodes zero body records, uses at most 16 KiB transient allocation, and performs no roster/mailbox query, contract stream, or checkpoint write. Contention beyond 250 ms returns `checkpoint_busy` with zero authorization/effect; no latency/read bound is claimed for that contended path. Oversized, corrupt, or inconsistent headers fail closed. Trace, allocation, and held-lock fixtures assert these bounds.

**AC-11: Compaction readiness prevents unsafe voluntary compaction.** `compaction_ready` is false whenever any active obligation lacks durable persistence, bound dispatch identity/epoch, stored report baseline, or typed next action, or any permit is issued/executing; launcher/client-initiated compaction is refused with the exact missing fields. Readiness becomes true only when all obligations satisfy that predicate, no permit is active, and the current point is a declared logical boundary. Table tests toggle every term independently, and an observed automatic-compaction fixture proves that bypassing this prevention layer sets `needs_rehydrate`; an unobserved fixture is classified under AC-4's limit.

**AC-12: Distinct lifecycle sources preserve continuity.** Fresh `startup`, validated `resume`, observed context `clear`, and observed same-session `compact` each match the normative transition table's identity adoption, epoch, pre/post `needs_rehydrate`, and failure result without executing prior-summary text. Mode-specific table fixtures cover all four sources, and a live Codex spike separately proves same-session `compact` ordering rather than inferring it from startup behavior.

**AC-13: Typed permits are action-scoped, single-use, and race-free.** A permit authorizes only its exact action class, target digest, assignment epoch, checkpoint generation, and obligation revision. Wrong-action, wrong-target, stale-generation/revision, and duplicate-use attempts cause zero side effects. A writer that changes generation or sets `needs_rehydrate` before the under-lock `executing` commit makes the permit fail; a writer after that fsynced commit observes an already-linearized action and cannot erase its indeterminate/reconciliation state. Successor commit preserves any later `needs_rehydrate` flag.

**AC-14: External effects recover without duplicate replay.** Before any spawn, state mutation, archive, or merge effect, its exact permit is durably `executing` and fsynced. Restart never invokes an executing/indeterminate action directly. For each action class, exact durable success evidence commits the successor and consumes the permit without repeating the effect; conclusive no-effect evidence may authorize a new permit only under that action's rules; partial, missing, multiple, mismatched, or ambiguous evidence remains blocked. Kill-point fixtures before effect, after externally visible effect, and immediately before successor/consumption commit observe zero duplicate effects.

### Interactive

**AC-15: Live Codex demonstration.** During a delayed worker repair, CL triggers `/compact`, interrupts one foreground wait, and sends an unrelated status question. The FO answers from durable state, restores wait, consumes the final notification, verifies the committed report, and continues to re-review or the next gate without a reminder. Evidence includes app-server notifications, plugin hook and `PLUGIN_DATA` traces, checkpoint generations, workflow git log, and clean state.

**AC-16: Plugin lifecycle assistance is evidence-gated and opt-out-safe.** Before the plugin claims enforcement on a target Codex version, a live trace establishes plugin discovery plus all five callback traces; `SessionStart`/`PostCompact` developer-context delivery before the first guarded action; two-worker `SubagentStart`/`SubagentStop` accounting in `PLUGIN_DATA` under duplicate and out-of-order events; and `Stop` `decision: "block"` acceptance for an unresolved fixture obligation while a resolved obligation may stop. The same trace proves that a recorded explicit captain `do not wait` opt-out suppresses the matching wait directive and stop block only. Absent, late, malformed, stale, or unproven plugin behavior cannot clear an existing durable fence; uncovered action classes remain outside the enforcement claim.

## Test plan

The riskiest unproved mechanism remains Codex ordering and coverage. Invalidate that mechanism before implementing the full protocol.

1. Add a versioned schema-extraction fixture during implementation, not ideation. Its generator runs `codex --version` and `codex app-server generate-json-schema --experimental --out <tempdir>`, hashes the selected source files, and extracts the exact hook events, output kinds, notification methods, `contextCompaction` item, compact-start input, and injection input asserted in this design. A test regenerates the manifest and fails with a reviewable diff on version, digest, or extracted-value drift. The fixture proves schema shape only.
2. Build a temporary target-Codex plugin/client probe with `SessionStart`, `PostCompact`, `SubagentStart`, `SubagentStop`, `Stop`, a context marker, `thread/inject_items`, and a blocking `preToolUse` hook. First fixture-test the rendered hook map, fixed context-directive text, `PLUGIN_DATA` serialization, duplicate-epoch handling, and the pure unresolved/opt-out predicate; those fixtures do not certify host execution. Then confirm plugin discovery and capture distinct `startup`, `resume`, `clear`, and `compact` traces, retaining callback stdin/env/stdout, `item/completed`, `hook/started`, `hook/completed`, `PLUGIN_DATA` transitions, context/injection acknowledgements, `turn/completed`, the first post-compaction tool event, and the final stop decision. Start two delayed subagents, exercise duplicate and out-of-order start/stop events, invoke `thread/compact/start`, and test manual and automatic compaction, unresolved and resolved stop, and an explicit captain `do not wait` opt-out. Fail the stronger path if any callback is absent, context is absent or late, same-session compact is inferred only from startup, duplicate events create multiple epochs, `PLUGIN_DATA` is unavailable or cannot attribute both fixture workers, the stop result is not demonstrably blocking or a resolved obligation cannot stop, the opt-out still injects or blocks, operator input races the watchdog, or any collaboration lifecycle call bypasses the claimed hook. Record exact covered and uncovered action classes; do not generalize from shell-tool coverage.
3. Add checkpoint serialization, private-path validation, atomic-write, compare-and-swap, incomplete-bind, digest, lifecycle, and stale-resume tests under a proposed `internal/rehydrate/` package. Generate exact-boundary and one-byte/one-record-over fixtures for the 4 KiB header, 1 MiB total checkpoint, 1,024-byte header paths, 128 obligations, one outstanding permit, and 4 KiB body strings; corrupt magic/checksum/digest/length/count and zero-count/nonzero-body fixtures fail before allocation or action. On an uncontended no-action probe, instrument one `stat`, one `pread <= 4096`, zero body decodes, and `-benchmem <= 16,384 B/op`. Hold the lock beyond 250 ms and assert typed `checkpoint_busy`, zero authorization/effect, and no fast-path latency claim. Add table tests that toggle each `compaction_ready` term independently, including an issued/executing permit, require a logical boundary, report missing fields, and refuse unsafe launcher/client compaction.
4. Drive the normative lifecycle table as data: for `startup`, `resume`, observed `clear`, and observed `compact`, assert identity adoption, epoch before/after, `needs_rehydrate` before/after, and each failure result. Add separate observed-event, failed-hook-with-event, eventless-with-obligation, and eventless-zero-obligation fixtures. The last must assert no transition and an unsupported/uncovered classification, not successful fallback.
5. Add permit/state-machine tests for the full pre-spawn path: persist unbound `await_completion` plus report baseline, assignment epoch, obligation revision, and `spawn_initial`; reconcile and issue the exact permit; under lock transition it to `executing` and verify file/directory fsync before the fake spawn counter can increment; execute spawn; bind the returned worker; commit `wait`; consume the permit. Assert zero effects for wrong action class, wrong target digest, stale generation, stale obligation revision, wrong assignment epoch, and duplicate use. A conclusively pre-effect error may consume as no-effect and reconcile toward retry; timeout/lost-response/partial errors remain executing/indeterminate. A killed process is not expected to write anything after death.
6. Add deterministic generation/compaction race tests around the executing-intent linearization point. If a writer increments generation and/or sets `needs_rehydrate` before the gateway's under-lock executing commit, the permit is rejected and the effect counter remains zero. If the executing commit fsyncs first, the writer sees the executing state, may set `needs_rehydrate` afterward, and cannot erase it; the effect result/successor commit preserves that flag. Cover two workflows, out-of-order completions, stale epochs, generation conflicts, and matching/malformed reports.
7. For each of `spawn_initial`, `state_mutation`, `archive`, and `merge`, run a subprocess fixture with kill barriers (a) after executing-intent fsync but before effect invocation, (b) after the effect becomes externally visible, and (c) after result observation but before successor/permit-consumption commit. Restart reconciliation from disk and a durable fake external ledger keyed by permit nonce. Assert no automatic replay and total effect count <= 1. Exercise exact-success recognition, conclusive no-effect recovery where the action rules permit it, and partial/missing/multiple/mismatched/ambiguous evidence producing the named typed blocker. Spawn absence without a conclusive idempotency query must block, not redispatch.
8. Add remaining concurrent-ledger and state-machine tests for hook failure, no-progress watchdog, and bounded wait timeout. Inject stale executable-looking summary text and assert that transitions depend only on typed durable inputs.
9. Add a split-root replay fixture under `internal/ensigncycle/testdata/compaction-rehydration/`: validation rejects to implementation, a delayed fixer commits a report, and validation must run again before the gate. Add `internal/ensigncycle/compaction_rehydration_test.go` to replace context with a summary that omits the obligation, deliver completion before and after the next turn, interrupt wait, verify the committed report, and assert the next gate or terminal state.
10. Extend integration tests for each Spacedock-owned guarded command to prove an observed flag blocks before effects, a matching typed permit writes/fsyncs executing intent before effect, successful rehydrate unlocks permit issuance, and the no-action/no-event probe meets the quantitative uncontended bounds. For raw collaboration and final-stop actions, assert only the action classes the live client spike proves; record uncovered classes rather than calling them eventless fallback. Assertions use tool/event traces and workflow state, not transcript phrasing.
11. Run AC-15 and AC-16 live and retain app-server JSONL, plugin hook log, `PLUGIN_DATA` audit, checkpoint/readiness/permit/reconciliation transitions, and the temp workflow git log as CI artifacts.

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

## Stage Report: ideation

- DONE: HIGH — Repair the impossible eventless fallback claim with a concrete guarded-action boundary.
  Added launcher-owned `guardCheckpoint(action_class)` before side effects on Spacedock-owned dispatch, state mutation, archive, and merge commands, including stable blocked output and a single bounded inactive metadata read. AC-1, AC-2, AC-4, AC-8, AC-10, AC-11, AC-14, and test steps 2, 4, and 7 now distinguish observed compaction, an already-durable obligation, proved client interception, and uncovered raw host actions. Zero-obligation unobserved automatic compaction now explicitly causes no transition and carries no correctness claim. This supersedes the prior report's absolute statement that every automatic host compaction sets `needs_rehydrate`.
- DONE: MEDIUM — Add a normative lifecycle transition table.
  The design now fixes accepted input/identity adoption, epoch behavior, `needs_rehydrate` before and after the fence, and failure behavior for `startup`, `resume`, observed `clear`, and observed `compact`. AC-12 and test step 4 drive these rows directly, including the unobservable-signal limit, so implementations cannot choose transitions that merely make their own tests pass.
- DONE: LOW — Make Codex schema evidence durable and reproducible without implementing product files during ideation.
  Replaced `/tmp` as the claimed durable basis with a specified versioned extraction fixture for implementation: exact generator command, Codex version, selected schema hashes, and extracted assertions, regenerated by a scripted drift test. The original temporary schema remains discovery provenance only; test step 1 owns the future fixture and explicitly does not claim runtime ordering.

### Summary

c6 now makes only enforceable continuity claims. Spacedock-owned guarded commands consult the durable checkpoint before side effects; raw host actions are covered only when the live spike proves a client interceptor, and eventless zero-obligation compaction is named as unsupported. A normative four-source transition table and reproducible schema-extraction fixture make the remaining implementation and validation targets non-circular while preserving the ECC-only readiness and durable-fence design.

## Stage Report: ideation

- DONE: HIGH — Replace blanket unresolved-obligation blocking with typed, single-use action permits.
  The ledger now records an exact typed next action, and permit issuance binds action class, target digest, checkpoint generation, obligation revision, assignment epoch, and successor transition. The required pre-spawn sequence is explicit: persist unbound `await_completion` plus `spawn_initial`, issue one matching permit, execute through the gateway, bind the returned worker identity, consume the permit, and commit `wait`. Wrong action/target, stale generation/revision/epoch, duplicate use, and side-effect failure paths are specified in AC-9/AC-13 and test step 5. This supersedes the prior repair's blanket `guardCheckpoint` predicate.
- DONE: MEDIUM — Eliminate the guard-to-side-effect race.
  Permit execution rechecks all bindings plus `needs_rehydrate` under the checkpoint lock and holds that lock through the side effect and successor ledger transition. Competing writers therefore linearize before execution and stale the permit, or after the committed effect. AC-6/AC-13 and deterministic barrier test step 6 prove both orderings with a zero-side-effect counter on rejection.
- DONE: MEDIUM — Define and test a quantitative inactive-read and corruption bound.
  Added a checksummed fixed 4 KiB header, 1 MiB total checkpoint cap, 128-obligation and one-outstanding-permit caps, 4 KiB string cap, one-stat/one-`pread <= 4096` inactive path, zero body decodes, and at most 16,384 B/op. The inactive gateway holds the checkpoint lock from header read through effect. Oversized, corrupt, inconsistent, truncated, or generation-regressing checkpoints fail closed before permits or effects. AC-7/AC-10 and test step 3 assert the exact boundaries and one-over failures.

### Summary

The fence now authorizes progress instead of deadlocking it. Typed next actions produce generation/revision-bound single-use permits, including the initial unbound spawn-to-bind transition, and the gateway makes validation plus effect linearizable under the checkpoint lock. A fixed bounded header gives the inactive path a measurable cost ceiling while oversized or corrupt state fails closed. ECC-only scope and the existing lifecycle/order evidence gates remain unchanged.

## Stage Report: ideation

- DONE: HIGH — Make external-effect permits crash-recoverable without blind replay.
  Permit execution now fsyncs a durable `executing` intent and action-specific reconciliation key before releasing the lock and invoking any external effect. After restart, executing is treated as indeterminate; the new process—not the crashed one—examines durable evidence for `spawn_initial`, `state_mutation`, `archive`, or `merge`. Exact success commits the successor without replay, conclusive no-effect evidence may authorize a new permit only under the action's rules, and partial/ambiguous evidence stays blocked. AC-7/AC-14 and kill-point test step 7 cover before-effect, after-visible-effect, and before-successor-commit deaths with total effect count at most one.
- DONE: MEDIUM — Preserve the prior race fix with a write-ahead linearization point.
  The gateway rechecks under lock, persists/fsyncs `executing`, then releases the lock for the external call. A generation or context-loss writer before that commit stales the permit; one after it observes recovery-visible execution and cannot erase the intent. The successor commit preserves later `needs_rehydrate`. AC-13 and test step 6 cover both orderings.
- DONE: LOW — Qualify inactive cost and bound lock contention.
  AC-10 now applies only to the uncontended, quiescent no-action probe: one lock acquisition capped at 250 ms, one `stat`, one `pread <= 4096`, zero body decodes, and at most 16,384 B/op. Lock timeout returns `checkpoint_busy` with no authorization or effect and makes no contended-path latency/read claim. Every external-effect action uses a persisted executing intent, even with zero obligations.

### Summary

External effects are no longer described as crash-atomic with the checkpoint. The durable guarantee is a pre-effect executing intent plus restart reconciliation that never blindly replays spawn, mutation, archive, or merge. Generation races remain ordered at that fsynced intent, and inactive cost is now explicitly an uncontended no-action bound with fail-closed lock timeout. All prior typed-permit, lifecycle, ECC readiness, and scope limits remain intact.
