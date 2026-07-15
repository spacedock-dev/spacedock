---
title: Replace reviewer-reuse inference with structured evidence
status: validation
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

## Stage Report: implementation

- DONE: Implement the consolidated per-host structured identity correlator returning reuse/fresh/no_reuse/unsupported; Codex thread-id + followup_task/send_input correlation; Claude agentId + input.name → SendMessage.to (replace the strings.Contains substring); surface errReviewerIdentityUnsupported (errors.Is-checkable, never a reuse pass); DELETE the narration regex bank, command-string/wait inference, durable-report matcher, and self-source test — net fewer parsers
  `codexReviewerIdentity` + `claudeReviewerIdentity` in `shared_reviewer_reuse_test.go`; the three `assert*ReviewerReuse`/`assertClaudeSingleEntityRejectionFlow` are thin outcome-mapping wrappers; deleted `codexAddressableWorkerAbsent`/`codexNarrationNegatesReuseRoute`/`codexReuseRouteAbsencePatterns`/`assertCodexFreshValidationWhenAddressableAbsent` (narration), `codexReviewerReuseViaAssignmentSurfaces`/`codexDispatchBuildFromCommand`/`codexDispatchStageFromFlags`/`jsonObjectsIn` (command-string/wait), `assertCodexReviewerReuseWithDurableState`/`codexReviewerDurableValidationOutput`/`validationReport` (durable-report), and `TestCodexReviewerReuseDocumentsExecObservabilityBoundary` (self-source). Commit 2a171761.
- DONE: Rework fixtures per the 5-part test plan — structured same-handle → pass; wrong-reviewer + ambiguous → no_reuse/unsupported; convert inference tests to errors.Is(errReviewerIdentityUnsupported); add input.name to Claude validation-spawn fixtures; add the AC-1 meta table
  Codex table reworked (two-distinct-spawn now passes as `fresh`; wrong-thread → no_reuse; command/wait, narration, empty → unsupported sentinel via a `wantUnsupported` column). Claude tables gained `input.name`. Inference tests converted: `TestAssertCodexReviewerReuseRetiresCommandStringWaitInference` + `...RetiresNarrationAndDurableInference`. AC-1 meta table: `TestReviewerReuseAC1InferenceRetirement`. Watchdog caller (`codex_collab_wait_watchdog_impl_test.go:355`) maps the oracle: fail on no_reuse, accept reuse/fresh, record unsupported (new `reviewerIdentityUnprovable` field surfaced by the live runner).
- DONE: gofmt -l clean; go vet; go test ./... and go test ./... -race exit 0; commit the deliverable; report the AC-1 value measurement and net-fewer-parsers delta
  `gofmt -l internal/ensigncycle/` clean; `go vet ./...` clean; `go test ./...` and `go test ./... -race` both exit 0 (ensigncycle 19.0s under -race). Deliverable committed on branch `spacedock-ensign/replace-reviewer-reuse-oracle-with-structured-evidence`. (Note: a pre-existing gofmt flag on the untouched `internal/release/journeydelta.go` is not part of this deliverable — `git diff HEAD` on it is empty.)

### Summary

Consolidated the reviewer-reuse oracle into one structured identity correlator per host; reviewer identity now comes only from native spawn/follow-up handles (Codex `receiver_thread_ids` + `followup_task`/`send_input`; Claude `agentId`/`input.name` + `SendMessage.to`), with `errReviewerIdentityUnsupported` when no handle exists. AC-1 VALUE measurement: the inference-only reuse-pass count moved from a measured baseline of 5 under the previous oracle (the 3 advance-mode command/wait fixtures, the current-v2 assignment fixture, and the durable-report fixture — 4 confirmed by a throwaway harness before deletion plus the live-advance-narration fixture proven by its pre-existing green test) to **0**, while every structured fixture still passes (`TestReviewerReuseAC1InferenceRetirement`). Net-fewer-parsers delta: three inference parser families deleted (narration regex bank, command-string/wait counter, durable-report matcher — 11 symbols) and two consolidated correlators added, reusing the existing single-pass `collab_tool_call`/stream-json scan; net diff -240 lines (431 insertions / 671 deletions), no new transcript-dialect parser or identity state machine (AC-4). Note for the gate: the deliberate coverage trade-off from ideation stands — on a host emitting no identity handle the oracle abstains on WHO (durable `assertRejectionFlow` still proves the re-review OCCURRED) rather than manufacturing a false pass.

## Stage Report: validation

- DONE: AC-1 (VALUE) — inference-only reuse-pass count 0 vs the >0 baseline, structured fixtures still pass; baseline confirmed real, not asserted
  Reproduced the baseline out-of-tree: ran the OLD oracle in a throwaway `git worktree` at `main` — the 3 advance-mode command/wait fixtures + the v2-assignment fixture (via `assertCodexReviewerReuse`) and the durable-report fixture (via `assertCodexReviewerReuseWithDurableState`) each graded reuse-pass (`err==nil`), 5/5. New code: `TestReviewerReuseAC1InferenceRetirement` PASS (uncached, -race) — the same fixture byte-shapes yield `reviewerUnsupported` (0 reuse passes) while all 3 structured fixtures pass and grade `reviewerReuse`.
- DONE: AC-2 — absent/ambiguous surface → errReviewerIdentityUnsupported via errors.Is, never a reuse pass; identity never from commands/waits/reports/narration
  `TestAssertCodexReviewerReuseRetiresCommandStringWaitInference` + `...RetiresNarrationAndDurableInference` assert `errors.Is(err, errReviewerIdentityUnsupported)` on the exact retired byte-shapes; table `wantUnsupported` columns confirm empty/uncorrelated/narration → sentinel. Correlators read identity only from `spawn_agent`/`Agent` handles.
- DONE: AC-3 — fixtures reject wrong-reviewer + ambiguous; self-source test and narration regex bank gone
  Codex `send_input`→impl thread = `no_reuse` (fail), uncorrelated thread = `no_reuse`; Claude impl-agentId / non-validation `SendMessage` = unsupported; ambiguous (no handle) = unsupported. `grep` confirms `TestCodexReviewerReuseDocumentsExecObservabilityBoundary`, `codexReuseRouteAbsencePatterns` and the other 10 inference symbols are absent from the tree.
- DONE: AC-4 — both host routings covered, net FEWER parsers (~ -240 lines), no new dialect parser / identity state machine
  `git diff main...HEAD --shortstat internal/ensigncycle/` = 431 insertions / 671 deletions (-240). Three inference parser families deleted; `codexCollabItem` pre-existed on `main` (same file) so the correlators reuse the existing `collab_tool_call`/stream-json scan — straight-line two-pass classification, not a state machine.
- DONE: Semantic adversarial pass — handle specificity, no caller maps unsupported→pass, honest coverage trade-off, journeydelta flag not from this diff
  Correlation is exact map lookup (`validationHandles[to]`, `validationThreads[tid]`) — the only substrings are stage classification of the spawn, not identity routing. Callers: `runCodexRejectionFlowWithRetry` runs `assertRejectionFlow` FIRST (proves re-review OCCURRED from durable state) then fails on `no_reuse`, accepts reuse/fresh, records `unsupported` in `reviewerIdentityUnprovable` (logged, never converted to a reuse claim); negative-scenario callers require `err != nil`; Claude live runner fails on unsupported (stricter). `git diff HEAD -- internal/release/journeydelta.go` is empty — its gofmt flag is pre-existing, NOT attributable to this diff.
- DONE: gofmt clean on changed files; go test ./... and go test ./... -race exit 0
  `gofmt -l` clean on all 6 changed files; `go test ./...` exit 0; `go test ./... -race` exit 0; full `internal/ensigncycle` pkg uncached under -race = ok (11.0s).

### Summary

PASSED. All four ACs have valid, independently reproduced evidence. AC-1's value claim was verified out-of-tree: the deleted oracle at `main` really graded 5 inference-only fixtures as reuse-passes, and the new correlators map those same byte-shapes to `errReviewerIdentityUnsupported` while every structured same-handle fixture still passes — so "5 → 0" is a reproduced behavioral change, not a comment. The consolidated per-host correlators bind identity to specific native handles (Codex `receiver_thread_ids`, Claude `agentId`/`input.name`) via exact lookup, with no substring/command/wait/durable/narration inference resurrected. The deliberate coverage trade-off is honestly implemented: `assertRejectionFlow` independently proves the re-review occurred, so abstaining on WHO (unsupported) never manufactures a reuse claim. This is a low-blast-radius, test-only oracle change fully covered by offline fixtures; a full detached adversarial audit was not warranted, and the baseline reproduction was itself run on a throwaway `main` checkout. No material findings; no deferred risks that fail a value AC under the supported workflow. The pre-existing gofmt flag on `internal/release/journeydelta.go` is out of scope for this diff.
