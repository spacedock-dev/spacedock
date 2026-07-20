# Ideation stage-report history

The following superseded stage reports were moved from `index.md` during the compression revision. Their text is preserved verbatim.

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

## Stage Report: ideation

- DONE: MEDIUM — Give zero-obligation non-worker effects a real permit subject.
  Added durable standalone session action intents for quiescent `state_mutation`, `archive`, and `merge`. Creation under lock is normative and records session/action identity, class/target digest, checkpoint generation, independent monotonic action revision, explicit before/after and successor/result fields, and state without any worker, obligation, or assignment epoch. Permits are now a discriminated `obligation|standalone` union; both use issued -> fsynced executing -> exact-success/conclusive-no-effect/ambiguous-blocker reconciliation and single-use semantics. `spawn_initial` remains obligation-backed.
- DONE: Update checkpoint/readiness/lifecycle criteria for standalone records.
  The header now counts up to 128 standalone records with at most one active; clean/crash retention, `compaction_ready`, inactive eligibility, subject validation, and executing recovery all account for them. AC-7, AC-10, AC-11, AC-13, AC-14, and new AC-15 distinguish standalone action revision/session binding from worker assignment binding.
- DONE: Add zero-obligation success, misuse, crash, and ambiguity tests for every non-worker action class.
  Test step 6 drives standalone creation through consumption for mutation/archive/merge and asserts wrong action/target, stale action revision/generation, duplicate use, and ambiguous blocking. Steps 7 and 8 exercise both subject kinds and all three kill boundaries with total effect count at most one; integration step 11 accepts either exact obligation or standalone permits.

### Summary

Quiescent mutation, archive, and merge no longer depend on fabricated worker metadata. Each receives a session-scoped, independently revisioned durable action record and a subject-typed permit, while spawn retains its assignment-backed obligation. The prior write-ahead execution, crash reconciliation, race ordering, inactive bounds, and ECC-only scope are preserved.
