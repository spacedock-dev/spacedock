---
title: Replace reviewer-reuse inference with structured evidence
status: implementation
source: "Test-infrastructure audit 2026-07-14."
started: 2026-07-15T06:52:43Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-replace-reviewer-reuse-oracle-with-structured-evidence
issue:
milestone: 0.26.0
id: 6vn56z3423xk3ej3wvk54z6r
---

## Problem

The reviewer-reuse oracle (`internal/ensigncycle`, test infrastructure) claims that the FO reused the kept-alive cycle-1 validation reviewer for the cycle-2 re-review. On Codex it does so through an ordered fallback chain that reconstructs identity from signals that cannot carry it:

- `codexAddressableWorkerAbsent` + `codexReuseRouteAbsencePatterns` — a bank of 7 regexes over free-form FO **narration** (`agent_message` text).
- `codexReviewerReuseViaAssignmentSurfaces` (+ `codexDispatchBuildFromCommand`, `codexDispatchStageFromFlags`, `jsonObjectsIn`) — counts `dispatch build --stage validation` **command strings** and `wait` **counts** (`waitCalls >= 2`).
- `assertCodexReviewerReuseWithDurableState` + `codexReviewerDurableValidationOutput` — reads the **durable entity report** (`cycle 1: rejected` + `cycle 2: passed`).

The oracle's own comment admits the gap (`shared_reviewer_reuse_test.go:270`: "Codex exec does not surface enough current multi_agent_v2 metadata to prove a distinct reviewer process"), yet the fallbacks pass anyway. A durable report proves a re-review *occurred*; it cannot prove *who* performed it. More parsing cannot create the missing fact. On Claude the branch is mostly structured but still accepts a `strings.Contains(to, "validation")` **name substring** rather than correlating to the specific spawn's handle. A test (`TestCodexReviewerReuseDocumentsExecObservabilityBoundary`) even asserts on the oracle's own source text (a self-source assertion).

## Required outcome

Claim reviewer reuse only from structured native spawn/follow-up handle correlation. When a host does not expose that surface, report identity as an explicit unsupported result; durable reports may prove that re-review occurred, but not who performed it.

## Proposed approach

**One structured identity correlator per host, replacing the fallback chain.** Each host gets a single function that reads the transcript once and returns an explicit outcome — `reuse`, `fresh`, `no_reuse`, or `unsupported` — with a short detail. The existing reuse assertions become thin wrappers mapping that outcome to pass/fail per their own contract. `unsupported` is surfaced as a distinct sentinel error `errReviewerIdentityUnsupported` (checkable via `errors.Is`) and is never a reuse pass. This is a **consolidation**: the existing single-pass `collab_tool_call` / stream-json scan is kept and the *alternative* parsers are deleted — net fewer parsers, no new dialect parser or identity state machine (AC-4).

**Codex** (`assertCodexReviewerReuse`) — keep the structured PRESENT branch; drive every outcome from `receiver_thread_ids`:
- bind each validation `spawn_agent` (`item.completed`, prompt names validation) to its `receiver_thread_ids`;
- exactly one bound validation thread + a `followup_task`/`send_input` to that same thread id → `reuse`;
- two or more distinct validation threads → `fresh` (distinct reviewers, structurally proven — the contract-correct behavior when `«addressable-worker»` is ABSENT, per `codex-first-officer-runtime.md`);
- a follow-up to a non-validation (impl/uncorrelated) thread while a validation reviewer exists → `no_reuse` (fail);
- no thread-bearing `collab_tool_call` at all → `unsupported`.
- **Delete** `codexAddressableWorkerAbsent`, `codexNarrationNegatesReuseRoute`, `codexReuseRouteAbsencePatterns`, `assertCodexFreshValidationWhenAddressableAbsent` (narration — the fresh case is now structural); `codexReviewerReuseViaAssignmentSurfaces` + `codexDispatchBuildFromCommand` + `codexDispatchStageFromFlags` + `jsonObjectsIn` (command-string/wait inference); `assertCodexReviewerReuseWithDurableState` + `codexReviewerDurableValidationOutput` + the then-unused `validationReport` regex (durable-report identity inference).

**Claude** (`assertClaudeReviewerReuse` and the reuse half of `assertClaudeSingleEntityRejectionFlow`) — keep the structured branch; correlate to the specific spawn's handles:
- bind the single validation `Agent`/`Task` spawn to the `agentId` returned in its `tool_result` AND the teammate handle declared in its `input.name`;
- a `SendMessage.to` equal to one of those correlated handles → `reuse`; two validation spawns → `fresh`;
- a `SendMessage` to a different (impl/uncorrelated) handle while a validation reviewer exists → `no_reuse`;
- a validation spawn exposing neither `agentId` nor `input.name`, or no validation spawn → `unsupported`.
- Replace `strings.Contains(to, "validation")` with exact correlation to the spawn's declared handle. Spike-confirmed: real recorded streams carry `input.name` on every spawn and `SendMessage.to` matches it exactly.

**Caller change** — `codex_collab_wait_watchdog_impl_test.go:344`: replace `assertCodexReviewerReuseWithDurableState(...)` with the reuse oracle mapped as: fail on `no_reuse`; accept `reuse`/`fresh`; treat `unsupported` as identity-not-provable (record, do not fail). `assertRejectionFlow` (line 341, unchanged) still proves the two-cycle re-review *occurred*.

**Coverage trade-off (deliberate, gated):** on a host whose transcript carries no identity handle, the oracle no longer manufactures a reviewer-identity pass from durable reports; it honestly abstains on WHO while still proving the re-review happened. The prior fallbacks gave false confidence — they "passed" without proving identity. If current Codex exec *does* emit `receiver_thread_ids` (as the collab-wait watchdog machinery and the runtime contract indicate), the live path lands on `reuse`/`fresh` with full identity coverage; if it does not, the honest outcome is `unsupported`.

## Acceptance criteria

- **AC-1 (VALUE):** The reuse oracle returns a passing reuse claim for exactly the structured-same-handle transcripts. Measured against the current-oracle baseline: the count of committed inference-only fixtures (command-string/wait-count, durable-report, narration) that yield a reuse pass moves from >0 today to **0**, while every structured fixture (Codex thread-id `followup_task`/`send_input`; Claude `agentId`/`input.name` `SendMessage`) still passes. Baseline can move the wrong way: a surviving fallback keeps the count >0. Verified by `go test ./internal/ensigncycle/`.
- **AC-2:** An absent or ambiguous identity surface yields the explicit `errReviewerIdentityUnsupported` sentinel — never a reuse pass — and identity is never derived from command strings, wait counts, durable reports, or narration. Verified by the converted inference tests asserting `errors.Is(err, errReviewerIdentityUnsupported)`.
- **AC-3:** The fixture tables reject wrong-reviewer (follow-up/`SendMessage` to a non-validation handle while a validation reviewer exists) and ambiguous-transcript (no structured handle) cases, and no test asserts on the oracle's own source text or on free-form prose/narration (`TestCodexReviewerReuseDocumentsExecObservabilityBoundary` removed; the narration regex bank deleted).
- **AC-4:** Both supported-host routing shapes stay behaviorally covered (Codex thread-id reuse + fresh; Claude agentId/name reuse + fresh; the `//go:build live` runners unchanged) with net **fewer** parsers — the command-string parser, narration regex bank, and durable-report matcher are deleted and no new transcript-dialect parser or identity state machine is added.

## Test plan

Offline Go table tests in `internal/ensigncycle` under the DEFAULT build tags (no model spend, milliseconds); the `//go:build live` runners keep feeding the same assertions real transcripts. Commands: `gofmt -l` (then `gofmt -w`), `go test ./...`, `go test ./... -race`.

1. **Codex correlator table** (`TestAssertCodexReviewerReuse`, reworked): structured `followup_task`/`send_input` to the bound validation thread → pass (`reuse`); two distinct validation `spawn_agent`s **with no narration line** → pass (`fresh`, proving narration is no longer load-bearing); follow-up to the impl thread while a validation reviewer exists → fail (`no_reuse`); follow-up to an uncorrelated thread → fail; command-string+wait-only → `errReviewerIdentityUnsupported`; narration-only → `errReviewerIdentityUnsupported`; empty → unsupported. Remove all `liveAbsentTwoFresh`/`falseAbsentReuse`/`absent*` narration cases.
2. **Codex inference-retirement:** convert `TestAssertCodexReviewerReuseAcceptsAdvanceModeValidationReroute`, `...AcceptsCIAdvanceModeWithoutImplementationFeedbackReflow`, `...AcceptsLiveAdvanceModeReuseNarration`, `...AcceptsCurrentV2AssignmentSurfaces` from pass expectations to `errors.Is(err, errReviewerIdentityUnsupported)`. Delete `TestAssertCodexReviewerReuseWithDurableState*` (durable-identity path gone) and `TestCodexReviewerReuseDocumentsExecObservabilityBoundary` (self-source, AC-3). Watchdog callers (`codex_collab_wait_watchdog_test.go:325`, `codexReviewerReuseJSONL`) stay green on the structured shape.
3. **Claude correlator table** (`TestAssertClaudeReviewerReuse` + `TestClaudeSingleEntityRejectionFlow`, reworked): add `input.name` to the validation-spawn fixtures (faithful to real dispatch); reuse-by-agentId and reuse-by-correlated-name → pass; `SendMessage` to impl handle while a validation reviewer exists → fail; uncorrelated agentId → fail; fresh cycle-2 spawn + name message → fail; validation spawn with neither name nor agentId → `unsupported`; empty → unsupported. Bare-mode two-fresh-spawns and impl-as-validator (shutdown-exempt) cases unchanged behaviorally.
4. **Negative scenarios** (`shared_scenarios_negative_test.go:106-111`): the no-reuse Claude/Codex fixtures still non-pass (now `no_reuse`/`unsupported`).
5. **AC-1 value measurement:** a meta table enumerating the committed inference-only fixtures asserts the count yielding a reuse pass is 0 (baseline >0), and the structured-fixture pass count is unchanged.

Cost/complexity: low. Deletions + one consolidated correlator per host + fixture `input.name` additions; no new files; the mechanism is fully covered by offline fixtures (no new live-only test needed). Estimate < 1 day.

## Documentation

No user-visible behavior change — CLI output, banners, and host integration are unchanged, and the reviewer-reuse *contract* in `codex-first-officer-runtime.md` (`«addressable-worker»` PRESENT → reuse, ABSENT → fresh) is unchanged; only the test oracle that grades it changes. No doc-site diff. The removed in-code self-admission comment (`shared_reviewer_reuse_test.go:270`) and its self-source test are internal test infrastructure, not documentation.

## Spike evidence

Riskiest unverified mechanism: does the chosen structured handle actually appear in a real (or committed-fixture) host transcript, and does structured-only correlation cleanly separate reuse / fresh / wrong-reviewer / unsupported once the inference fallbacks are gone?

- **Structured-correlator exercise (run):** a throwaway Go program reimplemented the proposed per-host correlator and fed it the committed fixture byte-shapes. All cases matched expectations: Codex `followup_task`→bound validation thread = `reuse`; two distinct validation spawns = `fresh`; follow-up to impl thread (validation reviewer present) = `no_reuse`; command-string+wait-only and narration-only = `unsupported`; narration-absent + two fresh spawns = `fresh` (structure overrides narration). Claude: agentId reuse and correlated-name reuse = `reuse`; two spawns = `fresh`; `SendMessage` to impl handle = `no_reuse`; uncorrelated agentId = `no_reuse`; spawn with no name/agentId or no validation spawn = `unsupported`. `CODEX_ALL_PASS=true CLAUDE_ALL_PASS=true`.
- **Claude handle appears in real data (verified):** across all `internal/ensigncycle/testdata/*.stream.jsonl`, every `Agent`/`Task` spawn's input key-set is `[description, name, prompt, subagent_type, team_name]`, and the set of `SendMessage.to` values equals the set of spawn `input.name` values (100% overlap). Structured name-correlation works on recorded streams; the current fixtures merely omitted `input.name`.
- **Codex handle (committed-fixture + contract):** the structured shape lives in committed fixtures (`codexReviewerReuseJSONL`, table `realReuse`/`realReuseV2`) and the existing structured branch already passes on them (green baseline: `go test -run 'ReviewerReuse|RejectionFlow|SingleEntity'` = ok). No real Codex `.jsonl` artifact exists in-repo; the runtime contract names `followup_task(target,message)` + `receiver_thread_ids` as the real reuse surface. Whether current live Codex emits it is exactly what the `unsupported` branch handles honestly, so the design is correct either way.
- **Baseline:** `go test ./internal/ensigncycle/ -run 'ReviewerReuse|RejectionFlow|SingleEntity'` passes today; deletion scope confirmed (each inference helper is used only by the paths being removed).

## Mechanism/value trace

- **Mechanism: consolidated structured correlator returning explicit `reuse`/`fresh`/`no_reuse`/`unsupported`.** Serves AC-1 (only structured same-handle correlation passes) and AC-2 (explicit unsupported sentinel). Simplest alternative considered: keep the fallbacks but downgrade them to warnings. Insufficient — a warning that still returns a *pass* leaves the false reuse claim in place, which is the exact defect. Second alternative: delete the fallbacks and leave the structured branch as-is. Insufficient — the Claude branch keeps the `strings.Contains(to,"validation")` prose substring (AC-3), there is no explicit unsupported result (AC-2), and `assertClaudeSingleEntityRejectionFlow` keeps its own substring, duplicating identity logic (AC-4).
- **Value:** a reviewer-reuse claim that means what it says. Structured native handles (Codex `receiver_thread_ids`, Claude `agentId`/`input.name`) can prove same-vs-distinct reviewer; command strings, wait counts, durable reports, and narration cannot. Where the runtime omits the handle, the claim stays honestly unproven rather than reconstructed.

## Stage Report: ideation

- DONE: Design the per-host structured reviewer-identity correlation naming the concrete native handle on each host and the honest unsupported result (AC-1, AC-2)
  Proposed approach names Codex `spawn_agent`→`receiver_thread_ids` + `followup_task`/`send_input` correlation and Claude `Agent` spawn→`agentId`+`input.name` + `SendMessage.to` correlation; `errReviewerIdentityUnsupported` sentinel for absent surfaces; deletes narration/command-string/wait/durable-report fallbacks (names each helper).
- DONE: Test plan — fixtures that pass only on structured same-handle correlation, reject wrong-reviewer + ambiguous cases, remove self-source/prose assertions (AC-3), keep supported-host routing covered without a new parser/state machine (AC-4); names gofmt, go test ./..., go test ./... -race
  Five-part offline table plan in the body; consolidation deletes 3 inference parsers (net fewer), removes `TestCodexReviewerReuseDocumentsExecObservabilityBoundary`.
- DONE: Spike the riskiest path — structured handle appears in a real/committed-fixture transcript before the gate
  Ran a throwaway structured-correlator against committed fixture byte-shapes (CODEX_ALL_PASS=true CLAUDE_ALL_PASS=true); verified real testdata streams carry `input.name` with 100% `SendMessage.to` overlap; confirmed green baseline `go test -run 'ReviewerReuse|RejectionFlow|SingleEntity'`. Evidence recorded in `## Spike evidence`.

### Summary

The reviewer-reuse oracle is pure test infrastructure in `internal/ensigncycle`. The redesign keeps the already-structured branches (Codex thread-id, Claude agentId/name), replaces Claude's `strings.Contains(to,"validation")` prose substring with exact spawn-handle correlation, and deletes the three Codex inference fallbacks (narration regexes, command-string/wait counts, durable-report identity) plus the self-source test — surfacing an explicit `errReviewerIdentityUnsupported` when no structured handle exists. The load-bearing decision is the deliberate coverage trade-off: on a host emitting no identity handle the oracle abstains on WHO (relying on `assertRejectionFlow` for re-review-occurred) rather than manufacturing a false pass — worth a staff review at the gate. The riskiest mechanism was exercised, not asserted: the structured correlator separates all four outcomes on committed fixtures, and real recorded Claude streams confirm the `input.name`↔`SendMessage.to` handle.
