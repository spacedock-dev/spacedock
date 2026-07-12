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

## Decision

Adopt a launcher-owned, typed continuation ledger plus a guarded rehydration/action gateway. Strategic compaction reduces exposure, but correctness comes only from durable obligations or standalone action intents, source reloading, generation/revision checks, and crash reconciliation. Generated summaries and plugin lifecycle data are advisory.

Instruction-only reload is readable but cannot recover assignment identity or typed continuation. Lifecycle injection is useful only where live ordering and interception are proved. The first Codex hook-ordering spike therefore remains blocking for every strong host-enforcement claim; unproved raw host actions stay outside the guarantee.

This decision deliberately separates prevention, recovery, and enforcement. `compaction_ready` prevents unsafe voluntary boundaries; the ledger and reconciliation protocol recover detected or already-durable work; the gateway enforces only action classes it owns or the live client proves interceptable. None of those layers may borrow evidence from another: readiness does not prove recovery, a hook notification does not satisfy an obligation, and a generated summary does not authorize an effect. That separation is the smallest design that preserves continuity without claiming omniscience over the host.

## Problem and enforcement boundary

Compaction can preserve a workflow's story while dropping its next obligation. In the originating replay, the summary mentioned waiting, but the FO dispatched a feedback repair and returned before consuming its completion, verifying the newer committed report, and re-running review. The worker remained live; the continuation was lost.

The central invariant is: after dispatch or reuse, the FO cannot claim completion until it consumes the matching completion signal, verifies a newer committed Stage Report against the stored baseline, durably records the typed successor, and reaches a gate, terminal state, or explicit blocker. Assignment epoch N cannot satisfy N+1. This invariant applies equally after compaction, interruption, timeout, resume, and feedback repair.

Summaries are historical reference only. They may omit obligations, contain stale executable-looking instructions, or fail to generate. No transition, permit, worker identity, completion attribution, or recovery decision may derive from summary prose.

The enforceable boundary is evidence-limited. Observed lifecycle/item events can set `needs_rehydrate`; an already-durable subject can fence later actions; and a proved gateway can intercept its action classes. If automatic compaction emits no observable event, no hook runs, the checkpoint has no active subject, and a raw host action bypasses the gateway, Spacedock has no fact from which to infer compaction. That path is unsupported, not a successful fallback. Likewise, collaboration calls, assistant prose, and final responses are enforced only after the live spike proves a client interceptor for those classes.

## Evidence and open spike claims

ECC main commit `40927950c49f6e742d341e20ff7b9b7e1e7bfff5` supplies the applicable strategic-compaction evidence:

| Evidence | Observed behavior | c6 consequence |
| --- | --- | --- |
| Pressure signal | `suggest-compact.js` uses latest transcript context tokens as the primary signal, tool count as fallback, and context-growth buckets to rate-limit non-blocking reminders. | Thresholds may suggest compaction; they never authorize it. Suggest only at a logical boundary after `compaction_ready`. |
| Lossy summary | `pre-compact.js` invokes an LLM and fails open to logging; `llm-summary.js` samples at most 25 recent turns and 7,000 characters. Session replay labels summaries historical-only. | Summary content is never authoritative and cannot replace the ledger or fence. |
| Startup/compact distinction | `session-start.js` distinguishes `startup`, `resume`, `clear`, and `compact`, and skips prior-summary injection for every non-startup mode. | Prove all four sources separately; same-session compact cannot inherit startup proof. |

The initial Codex 0.144.1 schema was generated with `codex app-server generate-json-schema --experimental --out <dir>`. Schema shape proves: plugin hook declarations for `SessionStart`, `PostCompact`, `SubagentStart`, `SubagentStop`, and `Stop`; corresponding runtime hook notifications; plugin-sourced hook runs and output kinds; `contextCompaction` items; legacy `thread/compacted`; `thread/compact/start`; `thread/inject_items`; and turn/hook completion notifications.

Implementation must preserve this evidence in a versioned extraction fixture containing the command, `codex --version`, selected file SHA-256 digests, and extracted assertions. A scripted regeneration test reports schema drift. Schema does not prove plugin discovery, callback delivery/order, developer-context visibility, `PLUGIN_DATA` lifetime or worker coverage, `Stop decision:block`, opt-out propagation, injection ordering, collaboration-tool interception, duplicate-event order, or watchdog safety. The live plugin/client spike must prove each claimed class or leave it unsupported.

## Normative design

The launcher creates a private session checkpoint and exposes its absolute path as `${SPACEDOCK_FO_CHECKPOINT}`. It records schema, random nonce, host identity, launcher/contract roots, source digest, monotonic generation, compaction/rehydrate epochs, `needs_rehydrate`, active subjects, and at most one permit. Strings are data, never commands or model context. The launcher/checkpoint service alone validates transitions and performs writes.

### Subjects

Three independent axes are used everywhere below. `continuation` is the worker obligation's workflow phase. `next_action` is the exact authorization: `spawn_initial|spawn_replacement|wait|verify_report|route_feedback|rerun_review|present_gate|terminalize` for obligations, or the named non-worker action for standalone subjects. `execution` is the external-effect record lifecycle: `planned -> issued -> executing -> consumed|blocked`. Restart treats persisted `execution=executing` as indeterminate evidence requiring reconciliation, but `indeterminate` is not a fourth axis or an authorization to replay. Neither `next_action` nor `execution` may be described as a continuation.

| Subject | Creation and identity | `continuation` | `next_action` and permit binding |
| --- | --- | --- | --- |
| Worker obligation | Before spawn/reuse/replacement, persist workflow/entity/stage/cycle, assignment epoch, report baseline, revision, and both typed axes. | `await_completion`, `verify_report`, `route_feedback`, `rerun_review`, `present_gate`, or `terminalize`. | Permit binds obligation key/revision and epoch. Initial dispatch uses `next_action=spawn_initial`; replacement first advances to a new assignment epoch/target digest and uses `next_action=spawn_replacement`. Success binds the returned worker, keeps `continuation=await_completion`, changes `next_action` to `wait`, and independently consumes execution. |
| Standalone session action | With zero unresolved obligations, under lock require `needs_rehydrate=false`, no active standalone action, and no permit; persist action ID, session identity, class/target digest, generation, revision 1, before/after evidence, successor, and empty result. | Not applicable; it is not worker workflow state. | `next_action` is exactly `state_mutation`, `archive`, or `merge`. Permit binds session identity plus action ID/revision and MUST omit worker, obligation, and assignment fields. Its separate execution record begins `planned`. |

Standalone creation is a real state transition, not a shortcut around the ledger. The gateway validates the canonical target and expected before/after evidence, writes `execution=planned` at subject revision 1, increments checkpoint generation, and fsyncs before it may issue a permit. Each later execution transition increments both revisions. Success fills result and successor fields; ambiguity fills the typed blocker and observed evidence. Terminal standalone records remain for the session so an old nonce is still recognized as consumed/blocked, then clean-session teardown removes them.

The ledger supports multiple obligations. Reconciliation snapshots checkpoint generation, roster/mailbox evidence, workflow state, entity/report baseline, and applicable contracts. It uses compare-and-swap: conflict means reread and recompute. Deterministic priority is matching completion verification; active feedback/re-review; ready gate/terminalization; waits for running workers; then newly ready independent dispatch. One completion advances only its matching epoch; independent work starts only after higher-priority subjects are accounted for.

`spacedock dispatch rehydrate --checkpoint "$SPACEDOCK_FO_CHECKPOINT" --json` is the integration boundary. It reloads the FO entry skill/eager imports, active runtime adapter, dispatch and feedback modules when applicable, stage definition, entity/frontmatter, latest report, and feedback cycles; verifies source digests; reconciles durable evidence; and returns exact typed next actions. It clears `needs_rehydrate` only after committing the successful epoch. Missing sources, changing digests, corrupt state, or ambiguous evidence leave the flag set and return a typed blocker.

Resume never silently repairs or deletes uncertainty. A mismatched nonce/host/workflow/schema/owner is quarantined; `execution=executing` enters indeterminate reconciliation; a crash before a dispatch effect retains `continuation=await_completion` with `next_action=spawn_initial|spawn_replacement`; a crash after a visible effect requires the action-specific evidence table. Age cannot mark a subject complete. A maintenance command may list or garbage-collect terminal/quarantined records, but unresolved state needs explicit recovery.

Spacedock-owned dispatch, state mutation/archive, and merge paths go through the action gateway before effects. Replacement dispatch is never an escape hatch: it first records the new assignment epoch/target and `next_action=spawn_replacement`, then follows the same permit/execution/reconciliation fence. `needs_rehydrate=true` blocks subject creation, permit issue, and permit execute. Raw host actions are covered only where a client interceptor is live-proved.

### Permit lifecycle

Permits are a discriminated `obligation|standalone` union containing nonce, action class, exact target/argument digest, resulting checkpoint generation, selected subject revision, and successor. Only obligation permits carry assignment epoch.

| `execution` transition | Durable requirements | Result and refusal rules |
| --- | --- | --- |
| `planned` -> `issued` | Under lock reconcile exact `next_action`; store one permit; increment subject/execution revision and checkpoint generation; atomic write plus file/directory fsync. | Wrong subject/action/target, active permit, `needs_rehydrate`, or priority conflict refuses with zero effects. `continuation` does not change. |
| `issued` -> `executing` | Under lock recheck nonce, generation, subject identity/revision, action/target, epoch if applicable, and flag; store reconciliation key and expected evidence; increment revisions; atomic write and fsync. | This fsynced intent is the external action's linearization point. A preceding writer stales the permit; a later writer sees recovery-visible execution. |
| `executing` -> `consumed` | Release lock for effect. On exact returned result, reacquire; verify subject/permit; fill result/successor; consume execution; increment revisions; fsync. | Preserve any later `needs_rehydrate`. Either spawn action binds worker and changes only `next_action` from the spawn action to `wait`; `continuation=await_completion` remains. |
| `executing` -> `blocked` | Crash, timeout, lost response, or partial result enters indeterminate reconciliation. Restart never invokes it directly; inspect action-specific durable evidence. | Exact success consumes without replay. Conclusive no-effect terminalizes the old execution and may create a new `planned` record. Missing/partial/multiple/mismatched evidence blocks. |

The crashed process records nothing after death; the already-fsynced `executing` intent is the recovery fact. A conclusively pre-effect live error may consume as no-effect. An ambiguous error remains executing/indeterminate. Consumed/blocked standalone records stay durable until clean session exit so old permit nonces remain rejectable.

Automatic continuation is bounded: at most one directive injection and one watchdog-started turn per `(compaction_epoch, checkpoint_generation)`. Stop refusal is not durable progress. A watchdog runs only when `next_action != wait` and no operator turn is pending; repeated state, failed fence, or no generation advance records a blocker and stops. Wait timeouts preserve the same worker/epoch; they never imply completion or redispatch.

## Crash recovery by action class

| Action | Exact-success evidence | Conclusive no-effect evidence | Ambiguous/partial result |
| --- | --- | --- | --- |
| `spawn_initial` / `spawn_replacement` | Exactly one durable acknowledgement keyed by permit nonce, action, new assignment epoch, and target digest plus matching roster/mailbox identity: bind it, set `next_action=wait`, and never spawn again. | Only an authoritative gateway/idempotency query proving that exact dispatch key was never accepted may authorize a new permit. Absence alone is insufficient. | Missing/unqueryable, multiple, wrong action/target, mismatched, or cross-epoch evidence -> `indeterminate_spawn`. |
| `state_mutation` | Canonical path equals exact postimage/value and expected path-scoped commit marker: commit successor without rewriting. | Exact preimage plus proof no matching commit/effect exists may authorize a new compare-and-set permit. | Mixed fields, unexpected commit, or third state -> `indeterminate_state_mutation`. |
| `archive` | Source absent, destination exact hash, and expected commit marker: commit successor without moving again. | Exact source present and destination absent may authorize a new permit. | Both/neither paths without proof, hash mismatch, or commit mismatch -> `indeterminate_archive`. |
| `merge` | Expected merged object/PR landed exactly once and terminal marker matches: commit successor without merging again. | Exact pre-merge head plus authoritative remote/PR proof of no merge may authorize a new permit. | Partial ancestry, changed heads, ambiguous API, or terminal mismatch -> `indeterminate_merge`. |

Every class is kill-tested after executing-intent fsync/before effect, after external visibility, and after result observation/before successor commit. Restart effect count must remain at most one.

## Lifecycle and compaction readiness

| Source | Identity and epoch | Fence transition | Failure |
| --- | --- | --- | --- |
| `startup` | New host identity, nonce, checkpoint; epoch 0; never adopt prior state implicitly. | Complete bootstrap/source validation before dispatch; false afterward with no active subject. | Abort launch/dispatch on identity, root, or source failure. |
| `resume` | Explicit retained checkpoint; validate nonce, resume identity, workflow IDs, schema, owner; preserve epoch. | True if retained flag or active subject; clear only after committed reconciliation. | Quarantine mismatch/stale/corrupt state; never reconstruct from summary. |
| `clear` | Observed clear for current validated identity; compaction epoch unchanged. | Set true before returning control, even with zero subjects; verify and clear at same epoch. | Leave true; guarded actions fail closed. Uncovered raw actions remain unsupported. |
| `compact` | Observed PostCompact/contextCompaction/legacy signal; deduplicate matching signals and increment once. | Set true before control returns, even with zero subjects; verify and clear at new epoch. | Leave true. Unobservable signals cause no transition and no eventless claim. |

`compaction_ready` is prevention, not recovery. It is true only at a declared logical phase boundary when every obligation is durable, every dispatched obligation is bound to worker/epoch, every report baseline and `next_action` exists, no standalone action is nonterminal, and no `execution` is issued/executing. Launcher/client compaction refuses false readiness and names missing fields. Observed automatic compaction still fences; unobserved zero-subject automatic compaction remains outside the guarantee. Pressure thresholds only suggest checking readiness.

Hooks assist but never authorize. SessionStart/PostCompact may inject a fixed reload/wait directive; SubagentStart/SubagentStop may cache observations in `PLUGIN_DATA`; Stop may block while a durable obligation remains; a typed captain `do not wait` opt-out suppresses only its matching wait directive/block. No hook or plugin data can create, satisfy, clear, or supersede a ledger subject.

## Safety and operating bounds

| Property | Exact bound/behavior |
| --- | --- |
| Filesystem | Session directory `0700`; files `0600`; reject symlinks, non-regular files, owner mismatch, traversal, nonce/identity mismatch, unknown schema, partial/regressing records. Canonical contract roots must remain inside the launcher-selected root after realpath resolution; reject outside paths and symlink escapes. Threat model excludes malicious unrestricted same-UID processes. |
| Header/checkpoint | Checksummed fixed 4 KiB header; total checkpoint <= 1 MiB; canonical header paths <= 1,024 bytes each; body strings <= 4 KiB. Header contains lengths, generations, flags, counts, identities, roots, source/body digests. |
| Records | <= 128 obligations; <= 128 standalone records with at most one active; exactly one outstanding permit. All-zero counts require zero body length. Excess/corrupt/inconsistent input fails closed before authorization/effect. |
| Data minimization | Serialization MUST NOT store credentials, prompts, transcript text, or worker output. Allowed content is bounded identities, paths, digests, revisions, epochs, typed states/actions, baselines, evidence keys/results, and blockers. |
| Quiescent probe | Uncontended no-action path only: one `stat`, one `pread <= 4096`, zero body decodes, <= 16 KiB transient allocation, no roster/mailbox query, contract stream, or write. |
| Lock | 250 ms acquisition ceiling. Timeout returns typed `checkpoint_busy` with zero authorization/effects; no latency/read bound is claimed for contended execution. External effects use persisted executing intent, not a no-write fast path. |
| Persistence | Lock + CAS generation; same-directory temp, file fsync, rename, directory fsync. Clean exit removes checkpoint only with no obligation/nonterminal standalone/permit; otherwise retain for explicit validated resume. |

## Acceptance and proof matrix

| AC | Required outcome | Proof |
| --- | --- | --- |
| AC-1 | On observed-compaction/already-durable paths, reload all authoritative source digests before covered action; execute only through matching generation/revision permit. | Source-omission/digest fixtures plus covered gateway traces; raw classes join only after live interceptor proof. |
| AC-2 | Launcher compaction, observed automatic compaction, and already-durable obligations produce zero premature stop, redispatch, completion, or state advance; matching newer report reaches successor in the same drive. | Split-root replay with omitted-summary/interruption variants; explicitly exclude unobserved zero-subject raw-host path. |
| AC-3 | Checkpoint + contracts + workflow + roster/mailbox reconstruct entity, stage, cycle, epoch, baseline, `continuation`, `next_action`, and `execution` without summary state. | Replace summary with stale executable prose and compare all reconstructed axes. |
| AC-4 | Deduplicate all observed compact signals; failed hook cannot clear an observed fence; eventless durable subjects allow only exact typed action; eventless zero-subject raw action is unsupported. | Event matrix asserts epoch/flag/action results, including no-transition limit. |
| AC-5 | Epoch N completion never satisfies N+1; completion clears only after matching signal, newer committed report, checklist, and successor. | Stale, duplicate, malformed, uncommitted, and cross-entity fixtures. |
| AC-6 | Multiple obligations complete either order without lost update/cross-attribution/early global clear; generation conflict invalidates old permits. | CAS/concurrency tests across two workflows and reversed completions. |
| AC-7 | Invalid paths/ownership/schema/header/digest/count/size/resume fail closed; contract roots cannot escape the launcher-selected root; prohibited sensitive/content fields never serialize; crash-retained subjects become indeterminate without replay. | Boundary/one-over/corruption, inside/outside/symlink-root, prohibited-field serialization, and crash-adoption fixtures. |
| AC-8 | Hook failure/timeout/late/duplicate/uncovered calls never clear existing durable state; continuation injection/watchdog is bounded and no-progress blocks. | Hook-failure and duplicate-event fixtures with injection/turn counters. |
| AC-9 | Initial/replacement dispatch persists `continuation=await_completion` and exact `next_action`; exactly one permitted spawn binds the returned worker, changes `next_action` to `wait`, and consumes execution independently. | Full `spawn_initial` and new-epoch/target-bound `spawn_replacement` paths plus wait/interruption and rejected-fix replay. |
| AC-10 | Quiescent no-action probe meets exact read/allocation/lock bounds; contention returns `checkpoint_busy` and zero effects. | Instrumented IO, `-benchmem`, and held-lock fixture. |
| AC-11 | Voluntary compaction is allowed only at logical boundary with all durable/bound fields and no active standalone/permit; observed automatic bypass sets fence. | Predicate table toggles each term and observed/unobserved automatic cases. |
| AC-12 | Startup/resume/observed clear/observed compact match the normative identity, epoch, flag, success, and failure table. | Table-driven lifecycle fixtures plus separate live same-session compact ordering. |
| AC-13 | Obligation and standalone permits bind the correct subject fields; wrong/stale/duplicate use has zero effects; executing-intent ordering preserves later flag. | Misuse table and deterministic writer-before/writer-after barriers. |
| AC-14 | Every external effect has fsynced executing intent; restart recognizes exact success, proves no-effect, or blocks ambiguity without duplicate replay. | Three kill points for both spawn actions and all three non-worker classes; durable fake external ledger; effect count <= 1. |
| AC-15 | Zero-obligation mutation/archive/merge create revisioned standalone subjects with explicit before/after/successor/result and no worker fields. | Success, misuse, stale revision, duplicate, crash, and ambiguous recovery table for all three. |
| AC-16 | Live delayed repair survives compact + interrupted wait + unrelated question, then consumes completion, verifies report, and continues without reminder. | App-server/plugin/checkpoint JSONL, workflow git log, and clean state. |
| AC-17 | Target-version plugin proves five callbacks, context ordering, two-worker `PLUGIN_DATA`, stop block/resolved stop, and scoped opt-out; unproved classes remain unsupported. | First live plugin/client spike with exact covered/uncovered class report. |

### Six proof suites

#### Suite 1 — Schema and live-host spike

Generate, version, hash, and extract the schema fixture first. Then probe plugin discovery, all five callbacks, context/injection order, two-worker coverage, `PLUGIN_DATA`, stop refusal and resolved stop, scoped opt-out, duplicate/out-of-order events, watchdog/operator race, and every collaboration/final-stop class claimed intercepted. Record exact covered and uncovered classes; any absent or ambiguous behavior disables only that stronger claim and cannot weaken the durable gateway.

#### Suite 2 — Checkpoint, lifecycle, and bounds

Exercise serialization, atomic write/fsync/rename, CAS, adoption, quarantine, and cleanup. Drive the four lifecycle rows and readiness predicate term by term. Test each exact/one-over bound, corruption combination, all-zero-count/nonzero-body rejection, quiescent IO/allocation, and held-lock `checkpoint_busy`. Add canonical contract-root fixtures inside the launcher root plus outside and symlink-escape refusals; attempt credentials, prompts, transcript text, and worker output and prove none serialize. Confirm crashes retain every nonterminal subject/permit while clean exit removes only eligible state.

#### Suite 3 — Subjects, permits, and concurrency

Drive full `spawn_initial` and `spawn_replacement` obligation paths, including a new replacement assignment epoch/target digest, plus zero-obligation standalone creation for all three non-worker actions. Table-test wrong subject kind, action, target, session, action/obligation revision, assignment epoch, generation, and duplicate nonce. Use deterministic barriers for writer-before/writer-after executing intent, two workflows, reverse completion order, priority, stale completion, and generation conflict. Every refusal asserts a zero external-effect counter.

#### Suite 4 — Crash reconciliation

Run both obligation-backed spawn actions and all three standalone actions in subprocesses; kill each at all three barriers. Restart from the on-disk checkpoint plus a durable fake external ledger keyed by permit nonce/action/target. Cover exact success, authoritative no-effect, missing/partial/multiple effects, mismatched epoch/target/identity/hash/commit/ancestry, and ambiguous API response. Reconciliation never calls the effect directly; every run proves total effect count at most one and the exact successor or blocker.

#### Suite 5 — Split-root continuity replay

Use a generated split-root workflow where validation rejects to a delayed implementation repair. Replace context with a summary omitting the obligation and another containing stale executable-looking text; deliver completion before and after the next turn; interrupt a wait; ask an unrelated question; verify the newer committed report; and continue through reviewer re-run and gate. Assertions use checkpoint bytes, event traces, exit behavior, workflow files, and git state—not transcript wording.

#### Suite 6 — Covered integration and live demonstration

Test every Spacedock-owned gateway path before effects, observed/eventless boundaries, `needs_rehydrate` issue/execute refusal, inactive bounds, wait restoration, and bounded continuation. Then run the live AC-16/17 delayed-worker scenario on the target Codex version. Retain app-server JSONL, hook log, plugin audit, checkpoint/permit/reconciliation generations, temporary workflow git log, and clean status as artifacts.

## Out of scope and evidence

Out of scope: implementing this protocol during ideation; changing Codex compaction algorithms; persisting conversation history; treating summaries as authority; changing Claude or Pi; weakening wait, feedback, gate, or merge contracts; using workflow mods as lifecycle hooks; malicious same-UID defense; and the separate `spacedock codex -- resume` bootstrap-prompt bug.

Evidence: ECC files at commit `40927950c49f6e742d341e20ff7b9b7e1e7bfff5`; live Codex 0.144.1 discovery schema (temporary provenance only, to become the versioned extraction fixture); current FO, dispatch, feedback, and runtime contracts. Superseded ideation reports are preserved verbatim in `artifacts/ideation-history.md`.

## Stage Report: ideation

- DONE: Compress within budget.
  Replaced the 9,229-word amendment stack with one current decision/design/proof artifact using the required outline and compact tables; the pre-commit result is 3,781 words and 185 source lines.
- DONE: Preserve invariant and AC coverage.
  Retained the completion/report/successor invariant, summary distrust, epoch/CAS isolation, obligation/standalone subjects, initial and replacement dispatch fencing, fsynced executing intent, no-blind-replay reconciliation, lifecycle/readiness boundary, contract-root containment, data minimization, non-authoritative hooks, bounded continuation, all exact limits, and AC-1 through AC-17. The current design consistently separates `continuation`, `next_action`, and `execution`; successful spawn leaves `continuation=await_completion`, changes `next_action` to `wait`, and consumes execution independently.
- DONE: Move history durably and verify readability structure.
  Moved every superseded ideation Stage Report verbatim to `artifacts/ideation-history.md`; `index.md` now presents one current design and one current Stage Report. Markdown headings/tables and `git diff --check` pass.

### Summary

c6 is now a compact current contract rather than an amendment log. Its audit history remains beside the folder-form entity, while the main body exposes one evidence-limited design, normative state/action tables, crash recovery, exact operating bounds, stable AC mapping, and six proof suites.
