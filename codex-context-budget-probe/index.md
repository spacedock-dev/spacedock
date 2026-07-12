---
title: Bind Codex context budget to per-thread token telemetry
status: implementation
source: "Captain request 2026-07-11, with 2026-07-12 native-index and read-only WAL evidence from live Codex 0.144.1."
started: 2026-07-11T05:50:29Z
completed:
verdict:
score: 0.9
worktree: .worktrees/spacedock-ensign-codex-context-budget-probe-jsonl
issue:
id: bct1zbqhbkatrwmqf6s6qd9v
---

## Problem

The Codex first-officer adapter needs a direct-child identity before it can decide whether to reuse a worker. `$CODEX_THREAD_ID` provides the parent identity, but `list_agents` exposes only a task path and status. A task path is inventory, not a thread identity.

The current implementation solves that gap by walking and fully parsing every session JSONL. On this host that means 978 files and about 1.1 GiB before it identifies the one child log, taking about 30 seconds. That all-history replay is not acceptable: it conflates child discovery with active-token projection, can grow without bound, and is unnecessary because Codex maintains its own local thread inventory.

The rejected observer experiment remains rejected: a second initialized app-server client could not resume the producer-owned rollout and received no usable token event. This redesign is limited to a synchronous local-index lookup plus selected-log projection. It does not add an observer, socket, sidecar, capability file, telemetry cache, companion process, PR/mod behavior, generic roster behavior, or a Codex launch wrapper.

## Spike findings

The spikes read only thread IDs, task paths, file paths, status, and allowlisted `token_count` fields. They did not read, retain, or print prompts, responses, reasoning, tool arguments, or lifetime token totals.

- Codex 0.144.1's local state DB contains `thread_spawn_edges(parent_thread_id, child_thread_id, status)` and `threads(id, rollout_path, agent_path, ...)`. A read-only query joining those tables on the current parent and exact worker uses `idx_thread_spawn_edges_parent_status`, then the `threads` primary-key index.
- A live query opened `/Users/clkao/.codex/state_5.sqlite` with SQLite's `-readonly` mode and URI `mode=ro&cache=private`. It reported `journal_mode=wal`, returned exactly one binding for parent `019f4fa1-3e99-7593-b79c-1d816e9bd467` plus `/root/bc_session_live_probe2`, and resolved child `019f53a3-ea0c-7e10-b307-47d3da3bf6e9` to its canonical rollout JSONL. The edge status was `open`; it is not a running-state signal and must not filter completed-but-addressable children.
- Codex's documented precedence is top-level `sqlite_home` in config, then `CODEX_SQLITE_HOME`, then `CODEX_HOME`; relative SQLite-home paths resolve from the current working directory. The config option identifies the directory holding the SQLite-backed state DB, not a fixed `~/.codex/state_5.sqlite` pathname.
- A selected session's metadata still revalidates the exact child ID, both parent fields, and byte-exact worker path. Its final projected `token_count` remains the only source of active-window tokens; the preserved compaction fixture is 18,893 -> 4,097 in a 353,400-token window. `threads.tokens_used` is lifetime state and is not a budget input.
- The current parent has had more than one historical child for a task path after replacement. The index query must therefore return at most one exact binding; it must not choose newest, highest counter, modified time, or an edge status.

## Proposed direction

Ship no observer, socket, sidecar, capability file, telemetry cache, or new `spacedock codex` launch mode. The public surface remains the synchronous command:

```text
spacedock dispatch codex-context-budget --worker <task-path>
```

The caller obtains `<task-path>` from the live roster, but the command treats it only as an exact lookup key. `$CODEX_THREAD_ID` is the required direct-parent anchor.

**Storage resolution and read-only index.** Resolve `CODEX_HOME` as today. Resolve the SQLite directory using Codex's documented precedence: user-level `sqlite_home` in `$CODEX_HOME/config.toml`, then `CODEX_SQLITE_HOME`, then `CODEX_HOME`; resolve a relative value from the current working directory. Parse TOML rather than grepping it. Canonicalize and validate that directory, then recognize only regular, non-symlinked `state_*.sqlite` candidates beneath it. Open a uniquely schema-qualified candidate with a pure-Go SQLite driver through a file URI whose query is exactly `mode=ro&cache=private`; use one short-lived connection and normal SQLite locking. Never enable `immutable=1` or `nolock=1`, never issue a mutating pragma or migration, and close the query and DB promptly. A small SQLite driver and TOML decoder are justified dependencies here; hand-parsing either format would weaken the safety boundary.

Schema qualification requires the tables and columns used by the lookup: `thread_spawn_edges(parent_thread_id, child_thread_id, status)` and `threads(id, rollout_path, agent_path)`. It must not hardcode the operator's `~/.codex/state_5.sqlite` path or a schema version. Missing, symlinked, locked, unreadable, schema-mismatched, or multiply-qualified state DBs are unavailable evidence. The query is bounded to two rows:

```sql
SELECT e.child_thread_id, t.rollout_path, t.agent_path
FROM thread_spawn_edges AS e
JOIN threads AS t ON t.id = e.child_thread_id
WHERE e.parent_thread_id = ? AND t.agent_path = ?
LIMIT 2;
```

It deliberately has no `status` predicate. Zero rows, a query error, or index/schema inconsistency yields `unavailable: binding`; two rows yield `unavailable: ambiguous`. Neither condition falls back to a full JSONL replay. V1 intentionally ships no header-scan fallback; a future dependency-free fallback would need a separately approved, fixed-byte header-only scan and could never rescue an existing index failure.

**Selected-log validation and token projection.** The single `rollout_path` returned by the index is not trusted blindly. It must be an absolute, canonical, regular non-symlink `.jsonl` below the canonical `$CODEX_HOME/sessions` root, with an identity check before and after open. Stream only that one file and recognize only:

1. `session_meta.payload.id`, `session_meta.payload.parent_thread_id`, and `session_meta.payload.source.subagent.thread_spawn.{parent_thread_id,agent_path}`.
2. `event_msg` records whose payload type is `token_count`, limited to the event timestamp, `info.last_token_usage.total_tokens`, and `info.model_context_window`.

The selected metadata must contain exactly one well-formed binding whose child ID equals the index row, whose two parent fields equal `$CODEX_THREAD_ID`, and whose worker path equals `--worker` byte-for-byte. Any metadata mismatch, duplicate target metadata, target-relevant malformed known record, unsafe path, or missing final record fails closed. A direct child that is completed remains eligible when this exact evidence is fresh because edge status does not determine addressability.

The command accepts a lower later `last_token_usage` value after compaction. It compares active-window values with integer arithmetic: `active_tokens * 100 <= context_window * 60`. Equality reuses; the next token fresh-dispatches. The parser has no field for `total_token_usage`, and it never serializes or persists raw JSONL records or content-bearing event types.

On valid evidence, stdout is stable JSON in this order: `worker`, `thread_id`, `active_tokens`, `context_window`, `remaining_tokens`, `remaining_pct`, `usage_pct`, `threshold_pct`, `source`, `observed_at`, `reuse_ok`. `source` is `session-jsonl`; an over-threshold result is valid JSON with `reuse_ok: false`. On unavailable evidence, the command exits non-zero, writes no JSON result, and emits one stable reason on stderr.

This is an automatic, conditional reuse gate rather than diagnostic-only output. Its safety boundary is deliberately narrow: the Codex runtime treats every non-zero result as PRESENT-but-unavailable and fresh-dispatches; it never treats a failed JSONL lookup as ABSENT. The command does not choose the newest candidate, a highest counter, or a modified file as a replacement tie-breaker.

| Evidence condition | Result |
| --- | --- |
| Missing or malformed `$CODEX_THREAD_ID`; unavailable/locked/schema-mismatched native DB; no index row; index/metadata parent-path-child disagreement | Non-zero `unavailable: binding`; fresh-dispatch. |
| More than one index row or duplicate matching metadata after a replacement/concurrent spawn | Non-zero `unavailable: ambiguous`; fresh-dispatch. |
| Missing final `token_count`; malformed timestamp or fields; zero window; negative or over-window active count | Non-zero `unavailable: record`; fresh-dispatch. |
| Timestamp older than 120 seconds or later than `now` | Non-zero `unavailable: stale`; fresh-dispatch. |
| State DB or selected rollout escapes its canonical root, is a symlink, or is not a regular file | Non-zero `unavailable: unsafe-path`; fresh-dispatch. |
| A unique fresh, valid snapshot | Exit zero; `reuse_ok` decides follow-up versus a new child. |

The reader never walks the sessions tree for discovery. It drops unknown selected-log records without logging their payloads. A malformed known target record invalidates the selected candidate. A stale record is never a fallback for a fresh record, and a tie is never resolved by ordering. These rules contain replacement, compaction, path, missing-record, index, and mismatch failures without retaining session content.

## Acceptance criteria

- **AC-1:** Against a fixture with 978 session JSONLs (the measured current-store baseline), each successful budget decision opens and parses exactly one selected JSONL and no decoy JSONL. An instrumented file-access test records the count, and a CLI fixture proves the same resolved child ID and output as the native index row. This is the measurable end value: one selected-log read rather than an all-history 978-file replay.
- **AC-2:** Every successful `codex-context-budget` result maps one exact `(parent thread ID, roster worker path)` pair through the native thread index to one child, then revalidates its two JSONL parent fields, child ID, and worker path. SQLite fixtures cover zero, one, and two query rows; wrong parent/path; stale index-versus-log disagreement; duplicate metadata; and a completed-but-addressable edge. A live-gated replay verifies a current worker's command result against its independently recorded child ID.
- **AC-3:** The state index is read-only and WAL-safe: `sqlite_home` overrides `CODEX_SQLITE_HOME`, which overrides `CODEX_HOME`; a normal-locking `mode=ro&cache=private` reader sees a writer's committed WAL row; and a write attempt through the reader is rejected. Absent, locked, unreadable, symlinked, schema-mismatched, or ambiguous index state produces non-zero/no-JSON and a fresh dispatch. Tests create a temporary WAL DB and prove these outcomes without `immutable=1` or `nolock=1`.
- **AC-4:** Only the index-selected path may supply active tokens. The selected path must remain canonical, beneath the sessions root, regular, and non-symlinked, and the parser must use only `last_token_usage.total_tokens` plus `model_context_window`. Fixtures cover path escape, symlink/non-regular paths, malformed known records, the 18,893-to-4,097 compaction pair, 60.0%, and the next-token boundary; source/output/artifacts omit lifetime and transcript-content sentinels.
- **AC-5:** Missing, stale, future-dated, malformed, ambiguous, replacement, locked, or mismatched evidence never causes a follow-up. CLI and adapter decision tests assert non-zero/no-JSON plus exactly one fresh dispatch for every unavailable condition, while a valid below-budget result retains the verified child ID and a valid above-budget result receives a distinct ID.
- **AC-6:** The production command works using only a native read-only state query and the selected JSONL. An integration fixture deliberately supplies no app-server, socket, sidecar, cache, capability file, remote UI, or companion executable and still exercises valid reuse and unavailable-to-fresh behavior.

## Test plan

1. Add storage-resolution tests using TOML fixtures plus `CODEX_HOME`, `CODEX_SQLITE_HOME`, and current-working-directory fixtures. Assert documented precedence and relative-path behavior; recognize exactly one regular, non-symlinked schema-qualified `state_*.sqlite` file and fail closed for absent, invalid, multiple, or lock-contended DBs. Estimated cost: medium.
2. Add a native-index package using a pure-Go SQLite driver. Its table fixtures create the required schema and exact parent/worker rows. Assert the `LIMIT 2` lookup returns zero/one/two correctly, ignores edge `status` as a running signal, and exposes no row content beyond child ID/rollout path/agent path. Estimated cost: medium.
3. Add a read-only/WAL test: a separate writer uses WAL and commits a matching row while it remains open; the `mode=ro&cache=private` lookup sees the commit. A write through the production read-only connection must return SQLite read-only, and lock/query/schema errors must reach the existing non-zero/no-JSON fresh path. Do not set `immutable=1`, `nolock=1`, migrations, or mutating pragmas in the production reader. Estimated cost: medium.
4. Refactor the session reader to accept the single index-resolved rollout. Inject an open-counting filesystem seam, place 977 content-sentinel decoys beside the selected fixture, and assert only the selected file opens. Retain selected-file metadata revalidation, path identity checks, compaction, threshold, stale/future, malformed-token, and content-sentinel coverage. Estimated cost: medium.
5. Update command goldens and adapter decision tests. Preserve Claude's existing `context-budget` bytes; distinguish `reuse_ok:false` from unavailable output; prove below-budget retains the exact verified child ID and above-budget/unavailable invokes one fresh child with a distinct ID. Estimated cost: small.
6. Add an opt-in live Codex replay with an explicitly recorded expected child ID. It queries the live native index through the built binary, verifies the selected session's fresh active-window result, then proves exact-path follow-up or fresh dispatch through entity markers, state-checkout git log, and clean status. It records DB mode/schema/query-plan metadata only, never transcript content. Estimated cost: medium.

The riskiest mechanism was native state-index lookup against a live WAL DB. The current read-only indexed join exercised it before implementation. No observer spike remains.

## Documentation change

The implementation updates `docs/runtime-support.md` and the Codex first-officer runtime binding. The public-facing addition is:

```diff
+ Codex context-budget checks use Codex's local read-only thread index to bind a live
+ roster task path to a unique direct child of the current first officer. They then read
+ fresh active-window metadata only from that selected session JSONL. Missing, locked,
+ stale, ambiguous, inconsistent, or unsafe evidence causes a fresh dispatch; Spacedock
+ never modifies Codex state or reads transcript content for this decision.
```

The runtime binding continues to call the same command and failure rule. It does not mention app-server fields or imply a background observer; it says that the index is discovery-only and token values come from the selected session JSONL.

The earlier stage reports below are historical evidence for the superseded direct-all-history reader, not validation for this revision.

## Stage Report: ideation

- DONE: AC-1 worker/task-path identity. Live spike: parent subAgentActivity mapped /root/alpha to child thread 019f4fc2-547f-7720-8ed9-a6183fe51e21; child thread/read confirmed its parent and task path.
- DONE: AC-4 active-window accounting. Compaction changed last.totalTokens 18,893 -> 4,097; lifetime total.totalTokens stayed 18,893 (modelContextWindow 353,400).
- SKIPPED: AC-2, AC-3, AC-5, AC-6, and AC-7. Intentionally deferred to implementation/validation: AC-2 golden CLI and routed follow-up (3-4); AC-3 threshold/adapter (2-3); AC-5 ordering/alpha-beta (2,4); AC-6 sentinel/mode/symlink (2); AC-7 snapshot/live threshold (4).

### Summary

Only AC-1 and AC-4 have live-spike evidence. The other criteria remain proof-planned; the two-client observer is the first implementation gate.

## Feedback Cycles

### Cycle 1 — captain-directed re-ideation (2026-07-11)

- The implementation bridge gate rejected the app-server observer-first direction: a second initialized client received `thread/started` but could not resume the producer-owned thread after its rollout existed (`-32603`), so it received no `thread/tokenUsage/updated` event.
- Preserve the existing uncommitted worktree experiment at `.worktrees/spacedock-ensign-codex-context-budget-probe` as non-shipping evidence; it must not become a partial public binding.
- Rework the v1 direction around Codex's own worker inventory and session JSONL: map `list_agents` task paths through parent/child metadata to fresh `token_count` records, remain fail-closed for missing, stale, or ambiguous evidence, and never parse or persist prompt/response content.
- Required ideation proof: a live current-worker mapping to a child session record, active-window (`last`) accounting across compaction, and explicit stale/mismatch safety. Decide whether the result can support automatic reuse or remains diagnostic-only.

### Cycle 2 — validation rejection (2026-07-11)

- Validation report `deafe49d7745104c767009459db7a45f3f197854` recommended REJECTED. Keep the direct JSONL-only design and rework in the existing `.worktrees/spacedock-ensign-codex-context-budget-probe-jsonl` worktree.
- Add deterministic evidence that a roster-selected worker retains its child thread when `reuse_ok:true` and receives a distinct child when the budget is above threshold or unavailable; do not substitute adapter-prose inspection.
- Make the opt-in live replay accept an independently supplied expected child thread ID, assert exact equality, and execute it when the host has sufficient Codex quota and temporary-disk capacity.
- Reconcile `skills/ensign/references/codex-ensign-runtime.md` with the FO's PRESENT context-budget binding. The validator's independent race rerun could not start because the Go temp filesystem had 136 MiB free; preserve artifacts and rerun it once capacity is available.

## Stage Report: ideation (cycle 1)

- DONE: Reframe the v1 observation mechanism around Codex inventory and session JSONL.
  Replaced the observer, socket, capability, and cache proposal with a direct-parent `$CODEX_THREAD_ID` plus allowlisted session-JSONL lookup.
- DONE: Prove the smallest current-worker mapping and active-window evidence without retaining content.
  Live roster path mapped to child `019f50f7-8f5b-7950-8a9b-102082cdcd3c`; at 2026-07-11T11:56:50Z its final projected `last_token_usage` was 120,352 / 353,400, 25 seconds old.
- DONE: Specify fail-closed eligibility, revised acceptance criteria, and a validation plan.
  Unique fresh binding supports conditional automatic reuse; duplicate, missing, stale, malformed, mismatched, and unsafe-path evidence exits non-zero and fresh-dispatches.
- DONE: AC-1 — Partial ideation evidence only: the live current-worker mapping and fresh active-window snapshot are recorded; retained-follow-up versus new-child behavior remains future validation.
- DONE: AC-2 — Partial ideation evidence only: one child session's two parent fields matched the current parent; command-wide mapping, fixtures, and completed-worker follow-up remain future validation.
- DONE: AC-4 — Partial ideation evidence only: the preserved compaction pair proves active-window semantics; parser field exclusion and threshold-boundary behavior remain future validation.
- SKIPPED: AC-3 — Production failure-path behavior is not proven at ideation; non-zero/no-JSON and fresh-dispatch tests remain planned.
- SKIPPED: AC-5 — Secure reader, path-safety, and content-sentinel behavior are not proven at ideation; temp-tree validation remains planned.
- SKIPPED: AC-6 — The production JSONL-only command and launcher path are not exercised at ideation; command and launcher validation remain planned.

### AC evidence map

No AC is fully implementation-proven at ideation. `922c7b9337ac715cb391040905b7ddceb5d583b6` is the durable body/evidence commit; its `## Test plan` records future validation, not passing results.

- **AC-1 — PARTIAL live evidence:** `922c7b9337ac715cb391040905b7ddceb5d583b6`, `## Spike findings`, bullets 1–2, proves one current task-path-to-child binding and a fresh active-window snapshot. The retained-follow-up versus new-child decision remains planned for `## Test plan` items 3–4.
- **AC-2 — PARTIAL live evidence:** `922c7b9337ac715cb391040905b7ddceb5d583b6`, `## Spike findings`, bullets 1–2, proves both child parent fields match the current parent for one worker. This stage did not exercise successful command results, fixture tables, or a completed-worker follow-up; `## Test plan` items 1, 3, and 4 own those proofs.
- **AC-3 — PLANNED validation:** `922c7b9337ac715cb391040905b7ddceb5d583b6`, `## Proposed direction`, defines the unavailable-result matrix, but the stage did not exercise any failure path. `## Test plan` items 1–3 must prove non-zero/no-JSON results and fresh dispatch.
- **AC-4 — PARTIAL semantic evidence:** `922c7b9337ac715cb391040905b7ddceb5d583b6`, `## Spike findings`, bullet 4, records the 18,893-to-4,097 active-window compaction pair. Parser field exclusion and the 60% boundary remain planned in `## Test plan` items 1 and 3.
- **AC-5 — PLANNED validation:** `922c7b9337ac715cb391040905b7ddceb5d583b6`, `## Spike findings`, records a manually restricted projection, not a secure reader implementation. `## Test plan` items 1–2 must prove canonical-root, symlink, non-regular-file, and content-sentinel behavior.
- **AC-6 — PLANNED validation:** `## Feedback Cycles` records the rejected observer bridge, and `922c7b9337ac715cb391040905b7ddceb5d583b6`, `## Proposed direction`, excludes its replacement. This stage did not exercise a production command or launcher; `## Test plan` item 3 must prove the JSONL-only path.

### Summary

The revised v1 is an automatic, conditional JSONL gate, not a diagnostic or observer bridge. Body commit `922c7b9337ac715cb391040905b7ddceb5d583b6` records the live mapping, compaction boundary, fail-closed matrix, revised ACs, and JSONL-only validation plan.

## Stage Report: implementation

- DONE: AC-1 / AC-2 direct-child reader and exact identity binding in code commit `914e0584bede8d9e577d412e17e9a7b6ff37d238`.
  `internal/codexsession.ReadBudget` requires a well-formed `$CODEX_THREAD_ID`, both metadata parent fields, and a byte-exact roster task path; its fixtures cover zero, one, duplicate, mismatched, replacement, and completed-status-independent candidates. Equality at 60% reuses and the next token produces `reuse_ok: false`.
- DONE: AC-3 fail-closed command and adapter routing in `914e0584bede8d9e577d412e17e9a7b6ff37d238`.
  `spacedock dispatch codex-context-budget --worker <task-path>` emits no JSON and exits 1 for `binding`, `ambiguous`, `record`, `stale`, and `unsafe-path`; `CodexContextBudgetRouteForProbe` sends every non-zero, malformed, missing, or false result to fresh dispatch.
- DONE: AC-4 active-window-only accounting in `internal/codexsession/budget.go` and its deterministic fixtures.
  The typed projection reads only `last_token_usage.total_tokens` and `model_context_window`; it accepts the 18,893-to-4,097 compaction drop and uses integer comparison for the 60% boundary.
- DONE: AC-5 canonical-path and transcript-exclusion coverage in `internal/codexsession/budget_test.go`.
  Tests reject symlinked roots/candidates and non-regular `.jsonl` paths, latch malformed known metadata or token records even when a later record is valid, and assert a prompt/response sentinel never reaches result JSON or errors.
- DONE: AC-6 JSONL-only product path.
  The implementation adds `internal/codexsession` plus the synchronous dispatch command only; it adds no observer, socket, sidecar, UI route, capability file, or telemetry cache. The first-officer binding now invokes the command and treats unavailable evidence as fresh dispatch.
- DONE: Opt-in live replay harness added at `internal/ensigncycle/codex_context_budget_live_test.go`.
  It runs the built command from an FO session with an explicitly roster-selected child, verifies a fresh content-free snapshot, then checks the child entity marker, state-checkout log, and clean entity status. It compiled and skipped without the four explicit live-run inputs in this implementation session; the earlier ideation spike remains the recorded live mapping evidence.

### Verification

- `rtk go test ./...` — exit 0.
- `rtk go test ./... -race` — exit 0.
- `rtk go test ./internal/codexsession -count=1` — 28 passing tests.
- `rtk go test ./internal/dispatch -run 'TestCodex(ContextBudget|MultiAgentV2)' -count=1` — 22 passing tests.
- `rtk proxy go test -tags live -run TestLiveCodexContextBudgetCurrentWorkerReplay ./internal/ensigncycle -count=1 -v` — compiled; skipped as designed without `SPACEDOCK_CODEX_CONTEXT_BUDGET_{WORKER,ENTITY,STATE_ROOT,MARKER}`.

### Summary

Implemented the approved direct-parent session-JSONL context-budget gate in code commit `914e0584bede8d9e577d412e17e9a7b6ff37d238`. The production path is content-free and fail-closed; validation should execute the opt-in live replay with a named child in a temporary split-root workflow before accepting the gate.

## Stage Report: validation

- DONE: Independently reproduce the direct JSONL identity, threshold, compaction, and unavailable-result behavior; assess every AC from observable evidence.
  `go test ./internal/codexsession -count=1`, `go test ./internal/dispatch -run 'TestCodex(ContextBudget|MultiAgentV2)' -count=1`, and `go test ./...` exited 0; AC-3 through AC-5 have fixture and golden coverage. AC-1 lacks retained-versus-new child-thread evidence, and AC-2 lacks an independently anchored live identity assertion.
- FAILED: Adversarially inspect the JSONL-only command and adapter for fail-open paths, unsafe traversal, observer/cache/capability routes, and transcript-content leakage.
  The shipped reader/command are direct JSONL-only (`internal/codexsession/budget.go:79-158`, `internal/dispatch/codex_context_budget.go:18-37`), but `skills/ensign/references/codex-ensign-runtime.md:12` still declares context-budget unavailable while the FO binding declares it PRESENT at `skills/first-officer/references/codex-first-officer-runtime.md:25`.
- FAILED: Run the relevant focused, repository, race, and feasible live-replay checks; distinguish a skipped live replay from a demonstrated current-worker replay and recommend PASSED or REJECTED.
  The live test skipped without its four required inputs (`internal/ensigncycle/codex_context_budget_live_test.go:37-45`); it accepts only a nonempty `ThreadID` (`:62-63`). The focused race command could not start because the Go temp filesystem had 136 MiB free and returned `no space left on device`.

### Recommendation

REJECTED. Bounce back to implementation:

1. Add deterministic retained-versus-new child thread-ID decision evidence for below/above/unavailable budget outcomes.
2. Make the opt-in live replay take an independently supplied expected named-child thread ID, assert equality, and execute it when infrastructure permits.
3. Reconcile the stale Codex Ensign context-budget-unavailable binding with the new PRESENT FO binding, and adjust guarding tests as appropriate.

### Summary

The fixture and command behavior is promising, but it does not yet prove the actual reuse decision or a live named-child mapping. The skipped live replay and stale runtime binding prevent a PASSED recommendation; the race rerun remains an infrastructure block, not a product failure.

## Stage Report: implementation (cycle 2)

- DONE: Add deterministic behavior evidence that a roster-selected worker retains its child thread below budget and receives a distinct child above budget or when evidence is unavailable. (AC-1, AC-3)
  Code commit `59ba810`; `TestCodexContextBudgetDecisionRetainsOrFreshensChildThread` uses the command-golden envelope, retains the verified child below budget, and invokes the injected fresh boundary exactly once for above-budget, unavailable, partial, malformed, and mismatched evidence.
- DONE: Require exact equality to an independently supplied expected child thread ID in the live replay, and execute it once Codex quota and temporary-disk capacity permit; record a genuine block otherwise. (AC-1, AC-2)
  Parent-FO replay passed in 11.09s: the command's `session-jsonl` identity equalled the named child's independently recorded `CODEX_THREAD_ID`; the child marker, state git log (`75dd88c`), and clean path status passed. Preserved artifact: `/tmp/spacedock-bc-live-replay.7OTBln`.
- DONE: Reconcile the Ensign and FO Codex context-budget bindings, then rerun the relevant baseline/race checks when capacity permits without replacing behavior proof with contractlint. (AC-3, AC-4, AC-5, AC-6)
  `skills/ensign/references/codex-ensign-runtime.md` now carries the same PRESENT/fail-closed rule; the baseline/race reruns retain the existing active-window, path-safety, transcript-exclusion, and JSONL-only command evidence. `internal/contractlint/codex_multi_agent_v2_contract_test.go` removes only its obsolete unavailable-string assertion and is not behavior evidence.

### Verification

- `rtk proxy go test ./... -count=1` — completed successfully after the schema fix.
- `rtk proxy go test ./... -race -count=1` — completed successfully after the schema fix.
- `rtk go test ./internal/codexsession ./internal/dispatch ./internal/contractlint ./internal/ensigncycle -count=1` — 780 passing tests.
- `rtk proxy go test -tags=live -run 'TestLiveCodexContextBudgetCurrentWorkerReplay$' ./internal/ensigncycle -count=1 -v` — parent-FO exact-ID replay passed in 11.09s.
- `gofmt -d` on the five changed Go files and `git diff --check` — clean.
- The comm officer's separately running full suite had not completed when this report was written; it is not included as green evidence.

### Summary

The second implementation cycle makes child identity an observable reuse decision, accepts the live session's safe embedded-parent metadata shape without relaxing target-relevant failure handling, and proves the exact live mapping. No observer, sidecar, cache, or transcript-content path was added.

## Stage Report: validation (cycle 2)

- DONE: Independently exercise the feedback decision boundary: only a complete verified envelope may retain the exact child ID; all other results invoke fresh exactly once with a distinct ID. (AC-1, AC-3)
  Fresh `go test ./internal/dispatch -run 'TestCodexContextBudget(CommandGoldens|DecisionRetainsOrFreshensChildThread|RouteFailsFreshOnEveryUnavailableProbe|RequiresWorker)$' -count=1 -v` passed (30 tests): a complete below-budget envelope retains its ID; above-budget, unavailable, partial, malformed, and wrong-worker results call fresh exactly once with the distinct fixture ID.
- DONE: Review the embedded-parent JSONL classification and preserved parent-session live replay artifact; verify exact child-ID, entity marker, state-log, and clean-state evidence without reading transcript content. (AC-2, AC-4, AC-5)
  The fresh codexsession suite exercised the non-target embedded-parent allow, target-relevant malformed-record rejection, active-window compaction/boundary, and canonical-root/content-sentinel safety. Preserved replay `/tmp/spacedock-bc-live-replay.7OTBln` has matching marker `019f51df-2041-75e1-b86a-9d9e1c9a3585`, state commit `75dd88c`, and clean entity path; no JSONL was opened. Current proof commit `442710b` records a completed-and-still-addressable child, fresh content-free exact-ID projection, then exact-path follow-up state commit `c5e7395` with the same ID.
- DONE: Run fresh baseline, race, focused, and live-harness checks as capacity allows; verify the Ensign/FO binding reconciliation and no observer/cache/capability route, then recommend PASSED or REJECTED. (AC-6)
  Fresh baseline and race runs over `./internal/codexsession ./internal/dispatch ./internal/ensigncycle` passed 687 tests each; focused command coverage passed 30 and contract compatibility 93 (not behavioral proof). The validator's parent-only live harness correctly skipped without supplied inputs; the current completed-addressable proof instead executed the parent command with exit 0, `source:"session-jsonl"`, a fresh 21,179/353,400 snapshot, and `reuse_ok:true`. The Ensign/FO bindings agree on the fail-closed command; the synchronous command golden exercises the JSONL-only route. A broad `go test ./...` host run remains blocked at unrelated real-Claude install behavior and is not counted as green evidence.

### Recommendation

PASSED. Commit `59ba8101` closes the prior decision, live-identity, and binding findings; all six ACs now have executable or durable state evidence, including the formerly missing completed-but-addressable path.

### Summary

Validation found one real AC-2 gap, held the verdict, and closed it with a new parent-FO completed-addressable replay rather than prose or contractlint. The proof retains only content-free output and durable state; the reader's JSONL input was never opened by the validator.

## Feedback Cycles

### Cycle 3 — captain-directed native-index re-ideation (2026-07-12)

- The current all-history JSONL discovery scans 978 files / about 1.1 GiB and is too slow for a context-budget gate.
- Replace discovery with Codex's native local `thread_spawn_edges` plus `threads` index, anchored by `$CODEX_THREAD_ID` and exact worker path. The selected rollout JSONL remains the sole active-token source.
- Open the state DB only through an explicit read-only `mode=ro&cache=private` URI, preserve normal SQLite/WAL locking, and never use `immutable=1` or `nolock=1`.
- Resolve the SQLite directory through effective Codex configuration rather than hardcoding `~/.codex/state_5.sqlite`; absent, locked, inconsistent, ambiguous, unsafe, or schema-mismatched evidence must fresh-dispatch, never trigger a full-transcript fallback.

## Stage Report: ideation (cycle 3)

- DONE: Reframed child discovery around the native indexed thread inventory rather than a full session-tree walk. The design's measurable target is one selected JSONL open against the 978-file current-store baseline.
- DONE: Exercised the riskiest live mechanism read-only. A `mode=ro&cache=private` query against the WAL state DB used the parent-edge and thread primary-key indexes and returned one exact parent/worker-to-rollout binding; no transcript content was read.
- DONE: Made read-only/WAL behavior, configuration precedence, schema qualification, selected-path validation, and fail-closed outcomes explicit in the problem, direction, ACs, and test plan.
- DONE: Preserved the existing public command and the selected-log active-window contract. No implementation, observer, cache, or fallback was added during ideation.

### Summary

The next implementation replaces all-history JSONL discovery with a short-lived native read-only index lookup, revalidates one selected rollout, and measures one file open rather than a 978-file replay. SQLite failure is safety evidence for fresh dispatch, not a reason to reintroduce an observer or broad scan.

## Stage Report: implementation (cycle 3)

- DONE: Write and commit a task-local implementation plan before production edits; include exact files and red-green test order.
  State commit `a105c38` records the file map, dependency choices, and focused red/green commands; code commit `b083ec6` adds the tested TOML storage-precedence resolver.
- FAILED: Prove the native index opens read-only and WAL-safe, then parse exactly one selected rollout JSONL.
  The WAL/read-only test went red as intended, but `modernc.org/sqlite@v1.34.5` cannot be fetched here: `proxy.golang.org` and `modernc.org` DNS time out; the one bounded GitHub-mirror check requires persistent replacement handling for transitive `modernc.org/*` modules, which the captain declined.
- SKIPPED: Prove all unavailable index/path/token evidence fresh-dispatches without a full JSONL replay.
  Stopped at the approved dependency boundary; no alternate driver, cgo binding, SQLite CLI, observer, or fallback was introduced.

### Summary

The storage-resolution substep is committed and its three focused tests pass (`rtk go test ./internal/codexsession -run '^TestResolveSQLiteHome' -count=1`). The actual native-index delivery is blocked on an obtainable approved pure-Go SQLite dependency; the command, JSONL reader, and runtime bindings remain unchanged.

## Stage Report: implementation (cycle 4)

- DONE: Write and commit a task-local implementation plan before production edits; include exact files and red-green test order.
  The prior task-local plan remains at state commit `a105c38`, and the worktree remained clean for this retry.
- FAILED: Prove the native index opens read-only and WAL-safe, then parse exactly one selected rollout JSONL.
  Fresh approved fetch failed: `rtk go mod download -json modernc.org/sqlite@v1.34.5` exited 1 with `dial tcp: lookup proxy.golang.org: i/o timeout` for the module `.info` URL.
- SKIPPED: Prove all unavailable index/path/token evidence fresh-dispatches without a full JSONL replay.
  Per captain direction, no retry workaround, replacement, cgo binding, CLI, observer, or broad replay was attempted.

### Summary

This retry changed no code or dependency files; the assigned worktree is clean. The same network/module-resolution blocker prevents the required pure-Go SQLite test and implementation from proceeding.

## Stage Report: implementation (cycle 5)

- DONE: Write and commit a task-local implementation plan before production edits; include exact files and red-green test order.
  The task-local plan and the prior green TOML-resolution substep remain unchanged; this retry began from a clean assigned worktree.
- FAILED: Prove the native index opens read-only and WAL-safe, then parse exactly one selected rollout JSONL.
  The required normal fetch still fails: `rtk go mod download -json modernc.org/sqlite@v1.34.5` exited 1 with `dial tcp: lookup proxy.golang.org: i/o timeout` for the `.info` URL.
- SKIPPED: Prove all unavailable index/path/token evidence fresh-dispatches without a full JSONL replay.
  The captain prohibited mirrors, replacements, cgo, external SQLite, sidecars, observers, and full-replay fallbacks; none was added.

### Summary

Product code and dependency files are unchanged by this retry. This report is committed locally only, per the current no-push instruction.
