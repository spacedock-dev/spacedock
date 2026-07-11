---
title: Bind Codex context budget to per-thread token telemetry
status: implementation
source: "Captain request 2026-07-11 plus live Codex 0.144.1 schema and rollout evidence gathered while reusing 88t workers."
started: 2026-07-11T05:50:29Z
completed:
verdict:
score: 0.9
worktree: .worktrees/spacedock-ensign-codex-context-budget-probe-jsonl
issue:
id: bct1zbqhbkatrwmqf6s6qd9v
---

## Problem

The Codex first-officer adapter declares `«context-budget»` ABSENT, so reuse-condition 0 always passes. Codex exposes enough local metadata to make a bounded decision, but `list_agents` exposes only a task path and status. A task path is inventory, not a thread identity.

The rejected observer experiment proved why v1 cannot depend on an app-server sidecar: a second initialized client saw `thread/started`, but it could not resume the producer-owned rollout and received no usable token event. The retained worktree is evidence for that rejection only. It must never become a bridge, a UI scraper, or a rollout-tail substitute.

This work is limited to the Codex budget binding: identity, fresh active-window telemetry, fail-safe reuse, its command surface, and tests. It does not add PR, mod, generic roster behavior, or a Codex launch wrapper.

## Spike findings

The replacement spike read only session identity fields and `token_count` fields from the live Codex session store. It did not read, retain, or print prompts, responses, reasoning, tool arguments, or lifetime token totals.

- The live roster contained `/root/spacedock_ensign_codex_context_budget_probe_reideation`. Its child session metadata mapped that exact path to `019f50f7-8f5b-7950-8a9b-102082cdcd3c`, with both parent fields equal to `019f4fa1-3e99-7593-b79c-1d816e9bd467`. This worker's `CODEX_THREAD_ID` equalled the child metadata ID.
- At `2026-07-11T11:44:54Z`, the child JSONL's last projected `token_count` was ten seconds old: event timestamp `2026-07-11T11:44:44.685Z`, `last_token_usage.total_tokens` 79,822, and `model_context_window` 353,400. This proves a current worker can be bound to a fresh active-window record without an observer.
- The current parent also has two child metadata records for each of `/root/spacedock_ensign_codex_context_budget_probe_implementation` and `/root/spacedock_ensign_codex_context_budget_probe_ideation`. A task path plus parent is therefore not sufficient after replacement. Selecting the newest file would be an unsafe guess.
- The preserved read-only compaction spike remains the semantic fixture: a compaction moved the active `last` count from 18,893 to 4,097 in a 353,400-token window. The JSONL reader uses the corresponding `last_token_usage` field and permits a lower later value. It never reads lifetime accounting.

## Proposed direction

Ship no observer, socket, sidecar, capability file, or new `spacedock codex` launch mode. The smallest v1 surface is a synchronous, Codex-specific command:

```text
spacedock dispatch codex-context-budget --worker <task-path>
```

The caller obtains `<task-path>` from the live roster, but the command treats it only as a lookup key. `$CODEX_THREAD_ID` is the required direct-parent anchor. The command scans the canonical `$CODEX_HOME/sessions` root and recognizes only these fields:

1. `session_meta.payload.id`, `session_meta.payload.parent_thread_id`, and `session_meta.payload.source.subagent.thread_spawn.{parent_thread_id,agent_path}`.
2. `event_msg` records whose payload type is `token_count`, limited to the event timestamp, `info.last_token_usage.total_tokens`, and `info.model_context_window`.

It accepts one candidate only when the child ID is well formed; both parent fields equal `$CODEX_THREAD_ID`; the metadata path equals `--worker` byte-for-byte; and exactly one matching candidate has a final, well-formed token record whose event timestamp is in `[now - 120 seconds, now]`. A direct child that is `completed` remains eligible when it meets those checks because the roster status does not change its session identity.

The command accepts a lower later `last_token_usage` value after compaction. It compares active-window values with integer arithmetic: `active_tokens * 100 <= context_window * 60`. Equality reuses; the next token fresh-dispatches. The parser has no field for `total_token_usage`, and it never serializes or persists raw JSONL records or content-bearing event types.

On valid evidence, stdout is stable JSON in this order: `worker`, `thread_id`, `active_tokens`, `context_window`, `remaining_tokens`, `remaining_pct`, `usage_pct`, `threshold_pct`, `source`, `observed_at`, `reuse_ok`. `source` is `session-jsonl`; an over-threshold result is valid JSON with `reuse_ok: false`. On unavailable evidence, the command exits non-zero, writes no JSON result, and emits one stable reason on stderr.

This is an automatic, conditional reuse gate rather than diagnostic-only output. Its safety boundary is deliberately narrow: the Codex runtime treats every non-zero result as PRESENT-but-unavailable and fresh-dispatches; it never treats a failed JSONL lookup as ABSENT. The command does not choose the newest candidate, a highest counter, or a modified file as a replacement tie-breaker.

| Evidence condition | Result |
| --- | --- |
| Missing or malformed `$CODEX_THREAD_ID`; no matching direct child; parent/path/child-ID disagreement | Non-zero `unavailable: binding`; fresh-dispatch. |
| More than one fresh matching child after a replacement or concurrent spawn | Non-zero `unavailable: ambiguous`; fresh-dispatch. |
| Missing final `token_count`; malformed timestamp or fields; zero window; negative or over-window active count | Non-zero `unavailable: record`; fresh-dispatch. |
| Timestamp older than 120 seconds or later than `now` | Non-zero `unavailable: stale`; fresh-dispatch. |
| Session root or candidate path escapes the canonical root, is a symlink, or is not a regular file | Non-zero `unavailable: unsafe-path`; fresh-dispatch. |
| A unique fresh, valid snapshot | Exit zero; `reuse_ok` decides follow-up versus a new child. |

The reader resolves the sessions root before walking it, does not follow symlinked entries, and accepts only regular `.jsonl` files that remain beneath that root. It drops unknown records without logging their payloads. A malformed known record invalidates its candidate. A stale record is never a fallback for a fresh record, and a fresh tie is never resolved by ordering. These rules contain replacement, compaction, symlink/path, missing-record, and mismatch failures without retaining session content.

## Acceptance criteria

- **AC-1:** A uniquely bound, fresh direct Codex child receives a follow-up at or below 60% active-window use and a new child above 60%. A deterministic decision test measures the retained versus new child thread IDs, and a live-gated replay verifies one current worker's command result against its session JSONL snapshot.
- **AC-2:** Every successful `codex-context-budget` result maps one roster task path to one child whose two parent fields match the invoking FO's `$CODEX_THREAD_ID`. Fixture tables cover zero, one, and multiple candidates; wrong parent/path; duplicate child metadata; and a completed-but-addressable worker. The live replay proves the current-worker mapping recorded above.
- **AC-3:** Missing, stale, future-dated, malformed, ambiguous, replacement, or mismatched evidence never causes a follow-up. Table and CLI tests assert a non-zero, no-JSON result with the stable reason; adapter tests assert a fresh dispatch for every unavailable result instead of an ABSENT fallback.
- **AC-4:** The threshold uses only `last_token_usage.total_tokens` and `model_context_window`, accepts a lower value after compaction, and returns the documented integer-boundary result. Fixtures include the 18,893-to-4,097 compaction pair, 60.0%, and the next-token boundary; source and output fixtures omit all lifetime fields.
- **AC-5:** The reader accepts only canonical-root regular JSONL files and never emits or persists content-bearing records. Temp-tree tests cover valid data, malformed known records, wrong root, symlink escape, non-regular files, and prompt/response sentinels. Output, stderr, and any test artifact omit the sentinels.
- **AC-6:** The production Codex path has no app-server observer, remote UI, socket, capability file, or telemetry cache. Command and launcher tests exercise the JSONL-only surface and assert that an unavailable lookup fresh-dispatches without launching a companion process.

## Test plan

1. Add a small JSONL reader package with an injected clock and sessions root. Fixture records contain only allowed `session_meta` and `token_count` fields plus opaque content sentinels. Table-test exact parent/path reduction, no candidate, duplicate fresh candidate, stale prior candidate, malformed known record, lower post-compaction count, future timestamp, and threshold arithmetic. Estimated cost: small.
2. Add filesystem fixtures for canonical-root validation. Exercise missing roots, symlinked root entries, symlinked and non-regular candidates, paths outside the canonical root, and a malformed final token record. Assert no content sentinel reaches a result, diagnostic, or persisted file. Estimated cost: small.
3. Add command goldens for valid under-budget, valid over-budget, and every unavailable reason. Preserve Claude's existing `context-budget` bytes. Add adapter decision tests that distinguish valid `reuse_ok: false` from an unavailable non-zero result and fresh-dispatch for both. Estimated cost: small.
4. Add an opt-in live Codex replay with a temp workflow. Capture the FO `CODEX_THREAD_ID`, dispatch a named direct child, use the live roster only to choose its task path, and require one exact child JSONL mapping plus a fresh `last_token_usage`/window result. A fixture, not an observer, covers duplicate replacement and compaction. The replay verifies process exit, entity body, state-checkout git log, and clean status. Estimated cost: medium.

No app-server spike remains. The riskiest mechanism is the live parent-anchor to child-session mapping, and the current-worker spike above exercised it before implementation.

## Documentation change

The implementation updates `docs/runtime-support.md` and the Codex first-officer runtime binding. The public-facing addition is:

```diff
+ Codex context-budget checks bind a live roster task path to a unique direct child
+ session of the current first officer. They read only fresh active-window token
+ metadata from the local session JSONL. Missing, stale, ambiguous, or unsafe evidence
+ causes a fresh dispatch; Spacedock never reads transcript content for this decision.
```

The runtime binding replaces `«context-budget»: ABSENT` with the JSONL command and its failure rule. It does not mention app-server fields or imply a background observer.

The earlier ideation stage report below records the rejected observer proposal. It remains historical evidence, not validation for this revision.

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
