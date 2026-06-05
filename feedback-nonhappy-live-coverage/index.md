---
id: gq9g4vrz03kgd8w46cvf09k7
title: Live-scenario coverage for the non-happy feedback-rejection paths
status: validation
source: "captain (2026-06-04) — a9 detached audit surfaced that the feedback-rejection non-happy-path guarantees are guarded only by review + the single-cycle live scenario. Investigation: the old tests/test_rejection_flow.py drove 2 full cycles + reviewer-reuse; the current Go rejection-flow scenario simplified to a single route-back; and NEITHER era ever drove the 3rd-cycle escalation or the budget-probe fail-safe. Use the existing prose-based shared-scenario runner to exercise these."
score: "0.30"
started: 2026-06-04T20:04:39Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-feedback-nonhappy-live-coverage
issue:
mod-block: merge:pr-merge
---

The feedback-rejection procedure carries non-happy-path behavioral guarantees that no standing test exercises. The a9 detached audit (commit 98629283) demonstrated each is gut-able while the suite stays green, because the static oracles assert only substring-presence of the prose, not the behavior. Add live-scenario coverage via the existing prose-based shared-scenario runner.

## The gaps (concrete)

- **3-cycle escalation.** The contract says: route findings back to implementation on cycles 1 and 2, and **on cycle 3 escalate to the human instead of auto-bouncing a 4th time**. Drift → the FO loops forever (reject → re-implement → reject …), burning tokens. NEVER tested: the live `rejection-flow` scenario drives one cycle; the old `tests/test_rejection_flow.py` drove two (impl → val REJECTED → feedback → cycle-2 impl-fix → cycle-2 reval via reviewer-reuse) but never a third. Highest blast-radius (runaway cost).
- **Budget-probe fail-safe (reuse condition 0).** Before reusing the kept-alive ensign for rework, the FO must consult the context-budget probe and **fresh-dispatch if the ensign is over budget or the probe is unavailable**. Drift → reuse of a blown-context ensign → degraded/failed work. Never tested either era.
- **Coverage regression in the Go port.** The current Go `rejection-flow` scenario *starts* from an already-rejected report and asserts a single route-back — it dropped the live 2-cycle trajectory + the SendMessage reviewer-reuse the Python test actually exercised. Restoring that multi-cycle trajectory is part of this task.

## Why this, why now

This generalizes across the 0.19.6 decomposition line (t3/a9/wm/p2): moving contract prose into lazy skills leaves these behavioral guarantees guarded only by review + sparse live scenarios. a9 did not regress it (the prose was never standing-behaviorally-tested) — but the guarantees are real, with real failure modes, so the lever is expanding live-scenario coverage, prioritized by blast-radius.

## Direction (for ideation to flesh out)

- Add a `feedback-3-cycle-escalation` shared/live scenario using the existing prose runner (`Use $spacedock:first-officer`): a fixture seeded with two prior rejections so the next is the 3rd; assert the durable end-state is an **escalation to the human** (an escalation marker / no 4th auto-bounce), graded on durable state, never transcript phrasing (per scenario-testing-principles).
- Evaluate restoring the full multi-cycle trajectory + reviewer-reuse the Go port simplified away.
- Decide whether the budget-probe fail-safe warrants its own scenario or is better served by the binary-gate task (see `feedback-guarantee-binary-gate`) — some of these guarantees may be cheaper to enforce in the binary than to drive live.
- Triage by blast-radius; not every non-happy path earns a (costly) live scenario.

## Notes

Sibling: `feedback-guarantee-binary-gate` (the "promote prose → binary code gate" lever — the stronger fix where a guarantee is mechanizable). Provenance: a9 (`feedback-rejection-flow-skill-extraction`) detached audit, 2026-06-04.

---

## Problem statement

The Feedback Rejection Flow (`skills/first-officer/references/first-officer-shared-core.md:158-170`) carries three durable-state guarantees that no standing test exercises:

1. **3-cycle escalation** — on the 3rd consecutive REJECTED validation the FO must escalate to the human instead of auto-bouncing a 4th time (line 164). Drift → an infinite reject→re-implement→reject loop burning tokens. **Highest blast-radius** (runaway cost).
2. **Multi-cycle trajectory + reviewer-reuse** — the old `tests/test_rejection_flow.py` (recovered at `fcb70def`) drove TWO full cycles and asserted the FO reuses the kept-alive validation reviewer via `SendMessage` for re-review (the `#141` keepalive contract). The Go port's `rejection-flow` simplified to a single route-back from an already-rejected report, dropping both the 2nd cycle and the reviewer-reuse signal.
3. **Budget-probe fail-safe (reuse condition 0)** — before reusing a kept-alive ensign the FO must consult `spacedock dispatch context-budget --name {ensign}` and fresh-dispatch if `reuse_ok:false` or the probe is unavailable (lines 134, 165; `claude-first-officer-runtime.md:128`).

These are guarded today only by review plus the single-cycle `rejection-flow` scenario. `internal/ensigncycle/feedback_test.go:28-30` already records the diagnosis: the deterministic *dispatch-side* reflow seam (`is_feedback_reflow` + `feedback_context`, `build.go` Rule 5 / Section 6) is byte-tested, but "the FO's PASSED/REJECTED gate decision and the 3-cycle `### Feedback Cycles` escalation has no in-process Go seam (FO-LLM prose, no internal/status parser)." Confirmed independently: `grep` of `internal/status/` finds no feedback-cycle field. So escalation is gradeable today ONLY by a live producer run against durable entity-body state — which is the lever this task owns.

## Riskiest-unknown spike (done, throwaway)

**Question that would invalidate the rest of the plan if false:** can a host-neutral pure-function assertion grade "3rd-cycle escalation vs 4th auto-bounce" on durable entity state ALONE (no transcript), given a fixture seeded with two prior `### Feedback Cycles` entries sitting at a 3rd REJECTED report?

**Exercised:** wrote a throwaway `TestSpikeEscalationGrading` in `internal/ensigncycle/` driving a candidate `assertThirdCycleEscalation(entity)` over three constructed end-states, ran it green, deleted it.

**Result — PASS.** The durable signals are sufficient and separable:
- GOOD (escalated): `### Feedback Cycles` carries ≥3 cycle entries, a fixture-instructed escalation marker line is present, and NO 2nd `## Stage Report: implementation` was appended → assertion green.
- BAD (4th auto-bounce): cycle-3 entry routes back to implementation AND a 2nd implementation report appears, no marker → assertion red on both the marker-presence and the no-4th-report checks.
- BAD (stalled at cycle 2): only 2 entries, no marker → red on the cycle-count check.

The escalation marker is durable only because the **fixture README instructs** the FO to record an escalation entry on the 3rd rejection (mirroring how the existing rejection fixture instructs an exact `shared-rejection-fix: applied` marker line). The marker is then a fixture-driven on-disk obligation graded as state, NOT transcript phrasing — satisfying scenario-testing-principles. This seeds the implementation's first codified-executor negative case.

## Decisions

**(a) Restore the multi-cycle trajectory + reviewer-reuse the Go port dropped — YES, folded into `rejection-flow`, not a new scenario.** The existing `rejection-flow` scenario IS the regression: it should drive the 2-cycle trajectory the Python test drove. Restoring it means evolving the existing fixture/assertion (the entity starts at cycle 1, the FO drives cycle-1 route-back → cycle-2 re-implement → cycle-2 re-validation), and asserting the durable 2-cycle end-state plus — on Claude teams — the reviewer-reuse `SendMessage` signal. Reviewer-reuse is host-specific (Codex uses `send_input`, bare mode has no reuse), so it is graded by the host runner, not the shared assertion. This keeps the shared assertion host-neutral over durable state and pushes the host-specific producer signal into the adapter, matching the existing split.

**(b) Budget-probe fail-safe — DEFER to `feedback-guarantee-binary-gate`, do NOT give it a live scenario here.** The probe is already a binary command (`spacedock dispatch context-budget`), and the sibling task is explicitly chartered to decide whether the dispatch/reuse path should *force* the consult (refuse to reuse an over-budget member) rather than instructing the FO in prose. A binary guard over a real reading is a deterministic, always-on, model-free oracle — strictly stronger and cheaper than a ~4-min live run that can only sample the behavior. Driving "the FO fresh-dispatched because the probe said over-budget" live is also brittle: it depends on manufacturing an over-budget transcript for the probe to read. Recorded deferral, no silent cap: if the sibling concludes the budget-probe is NOT mechanizable (stays prose), a `feedback-budget-fresh-dispatch` live scenario returns to THIS task's queue. Tracked in the Deferred section below.

## Coverage matrix (scenario × executor/mode/runtime)

Per scenario-testing-principles, the baseline row is `runtime {claude, codex} × mode {codified, llm-live}`. Pi is a separate smoke lane (`TestLivePiSubagentEnsignSmoke`), NOT in the shared parity set (`TestSharedScenarioRunnerCoverage` checks only `codexScenarioRunners()` + `claudeScenarioRunners()`), so a shared scenario costs exactly two live runners.

| scenario | codified (offline, always-on) | llm-live claude | llm-live codex | blast-radius |
|---|---|---|---|---|
| `feedback-3-cycle-escalation` (new) | negative case: 4th-bounce + stalled end-states go red | real escalation run | real escalation run | **highest** (infinite loop) |
| `rejection-flow` (evolve to 2-cycle) | negative cases: un-routed + status-not-back | 2-cycle + reviewer-reuse | 2-cycle + send_input reuse | medium |
| budget-probe fail-safe | — (deferred to binary-gate sibling) | — | — | medium (deferred) |

Priority order (by blast-radius): the new `feedback-3-cycle-escalation` first (it is the un-tested infinite-loop guarantee); the `rejection-flow` 2-cycle restoration second (a coverage regression, but the single-cycle scenario still proves route-back today).

## Acceptance criteria

**AC-1 — A `feedback-3-cycle-escalation` shared scenario exists and is host-parity-enforced.** A `sharedRuntimeScenario` entry named `feedback-3-cycle-escalation` is registered, and both `codexScenarioRunners()` and `claudeScenarioRunners()` carry a runner for it.
Verified by: `TestSharedScenarioRunnerCoverage` and `TestSharedRuntimeScenarioDefinitions` (the parity + definition guards) pass with the new ID; a deliberate removal of either host's runner makes `TestSharedScenarioRunnerCoverage` red.

**AC-2 — The escalation assertion grades durable state, never transcript, and goes red on the 4th-auto-bounce end-state.** A host-neutral `assertThirdCycleEscalation(entity)` requires ≥3 `### Feedback Cycles` entries, the fixture-instructed escalation marker, and the ABSENCE of a post-cycle-3 implementation report.
Verified by: a new offline negative case in `shared_scenarios_negative_test.go` that builds the 4th-auto-bounce broken end-state (cycle-3 routes back + a new implementation report) from the real fixture and asserts `assertThirdCycleEscalation` returns an error; plus a stalled-at-cycle-2 case. (Spike already proved this grading is sound and separable.)

**AC-3 — The live escalation producer reaches a 3rd-cycle decision and parks for the human.** A fixture seeded with two prior `### Feedback Cycles` entries at a 3rd REJECTED validation report drives the real FO; the durable post-run entity shows the escalation marker recorded and status NOT routed back to implementation a 4th time.
Verified by: `TestLiveClaudeSharedScenarios/feedback-3-cycle-escalation` and `TestLiveCodexSharedScenarios/feedback-3-cycle-escalation` (a cited LLM-executor run, per the citation gate — the codified executor alone does not satisfy a runtime-observable AC). Cost: ~4-min timeout class per host (mirrors `rejection-flow`), two hosts → ~8 model-minutes per CI live pass.

**AC-4 — `rejection-flow` drives a 2-cycle trajectory with reviewer-reuse, restoring the Go-port regression.** The `rejection-flow` fixture starts at cycle 1; the FO routes back, re-implements, re-validates a second cycle, and the durable end-state carries two implementation reports and two recorded cycles; the Claude runner additionally observes the reviewer-reuse `SendMessage` (Codex: `send_input`).
Verified by: the evolved `assertRejectionFlow` over the 2-cycle durable end-state (offline negative case proving a single-cycle end-state now goes red), plus `TestLive{Claude,Codex}SharedScenarios/rejection-flow` with the host-specific reuse-signal check in the runner. Cost: the existing `rejection-flow` live budget likely needs a bump from 4 min (it now drives two full cycles); the implementation stage sizes this against a real run.

**AC-5 — The seed-scenarios doc block stays bound to the code table.** `docs/specs/scenario-testing-principles.md`'s `<!-- seed-scenarios -->` block and `docs/dev/README.md`'s shared-scenario list name `feedback-3-cycle-escalation`.
Verified by: the existing doc-lock test that binds the doc seed IDs to `sharedRuntimeScenarios()` (reds on drift in either direction); `TestSharedScenarioDocsContract` for the README contract clauses.

## Deferred (no silent caps)

- **Budget-probe fail-safe live scenario** → owned by `feedback-guarantee-binary-gate` (decision (b)). Returns to this queue as `feedback-budget-fresh-dispatch` ONLY if that sibling rules it non-mechanizable.
- **Pi runtime parity for these scenarios** → out of scope. Pi runs a separate smoke lane, not the shared parity set; extending Pi parity is a standalone matrix expansion, not part of the feedback-guarantee coverage.

## Test plan

- **Codified executor (offline, always-on, no model):** AC-2 and the AC-4 single-cycle-now-red negative case live in `shared_scenarios_negative_test.go` as pure functions over `(entity[, observed])` strings. Cost: milliseconds, every CI. This is the spike's grading, promoted.
- **LLM executor (live, gated, cited):** AC-3 (both hosts) and AC-4's live legs. Cost: ~4 min/host for escalation; the 2-cycle `rejection-flow` likely needs a budget bump measured during implementation. Run via `go test -tags live -count=1 -run TestLive{Claude,Codex}SharedScenarios ./internal/ensigncycle -v`.
- **Static parity/doc guards (offline):** AC-1 and AC-5 ride the existing meta + doc-lock tests; no new test machinery, just new expected IDs.
- **No spike needed beyond the one done:** the escalation-grading mechanism (the one unverified path) is exercised above. Everything else composes already-proven machinery — the shared-scenario table, the runner adapters, and the offline negative-case pattern are all shipped and green.

## Stage Report: ideation

- DONE: Design the `feedback-3-cycle-escalation` live scenario on the existing prose-based shared-scenario runner (fixture + durable-outcome assertion graded on durable state, never transcript)
  Fleshed out in the body: fixture seeded with two `### Feedback Cycles` entries at a 3rd REJECTED report; assertion `assertThirdCycleEscalation` over ≥3 cycle entries + escalation marker + no 4th impl report. Per scenario-testing-principles; AC-1..AC-3.
- DONE: Exercise the riskiest unknown — run the smallest probe that the runner + a 2-prior-rejection fixture can reach a 3rd-cycle decision before committing to the plan
  Throwaway `TestSpikeEscalationGrading` proved the durable signals separate escalated/4th-bounce/stalled end-states (green, then deleted). Recorded under "Riskiest-unknown spike". Confirmed no `internal/status` feedback-cycle seam, so durable entity-body state is the only oracle.
- DONE: Scope decision (a) — restore the multi-cycle trajectory + SendMessage reviewer-reuse the Go `rejection-flow` scenario dropped
  YES; folded into the existing `rejection-flow` (it IS the regression), reviewer-reuse graded by the host runner not the shared assertion. AC-4.
- DONE: Scope decision (b) — does the budget-probe fail-safe get its own scenario or defer to `feedback-guarantee-binary-gate`
  DEFER to the binary-gate sibling (the probe is already a binary command; a guard is stronger + cheaper than a live run). Recorded deferral with a return-condition; no silent cap.
- DONE: Name the coverage matrix (scenario × executor/mode/runtime) and prioritize by blast-radius
  Matrix table in body: `runtime {claude, codex} × mode {codified, llm-live}`; Pi is a separate smoke lane outside parity. Priority: 3-cycle escalation first (infinite-loop, highest), 2-cycle `rejection-flow` second.
- DONE: Write ACs each backed by an outside-body gate (a failing live/offline scenario assertion), with cost estimates; note deferred paths (no silent caps)
  AC-1..AC-5 each cite a named test that can red (parity guard, offline negative case, live cited run, doc-lock). Costs: offline ms; live ~4 min/host. Deferred paths listed explicitly.

### Summary

Designed `feedback-3-cycle-escalation` on the existing prose-based shared-scenario runner and settled both scope questions: restore the 2-cycle trajectory + reviewer-reuse into `rejection-flow` (decision a), and defer the budget-probe fail-safe to the mechanizable binary-gate sibling (decision b). The riskiest unknown — that 3rd-cycle escalation vs 4th-auto-bounce is gradeable on durable entity-body state alone (no `internal/status` seam exists for feedback cycles) — was exercised with a throwaway green-then-deleted spike, seeding the implementation's first codified-executor negative case. Five ACs each cite an outside-body gate (parity/doc-lock tests, offline negative cases, and cited live runs); the budget-probe and Pi-parity deferrals are recorded with return-conditions, no silent caps.

## Stage Report: implementation

- DONE: `feedback-3-cycle-escalation` registered in `sharedRuntimeScenarios()` with both a Claude and a Codex runner so `TestSharedScenarioRunnerCoverage` + `TestSharedRuntimeScenarioDefinitions` stay green
  Added the table entry + `runClaude/CodexFeedback3CycleEscalationScenario` in both runner maps; host-neutral `assertThirdCycleEscalation(entity)` grades ≥3 `### Feedback Cycles` entries + the fixture-instructed `feedback-escalation: human-review-required` marker + no post-cycle-3 implementation report; offline negatives `TestThirdCycleEscalationNegativeAutoBounce` (4th-bounce + stalled-at-cycle-2) and `TestAssertThirdCycleEscalation` red the broken end-states. Parity falsifiability confirmed: removing the Codex runner reds `TestSharedScenarioRunnerCoverage` with "no Codex runner". Commit 7e03d5ea.
- DONE: `rejection-flow` evolved to drive the 2-cycle trajectory + the host-specific reviewer-reuse signal (Claude SendMessage / Codex send_input), restoring the Go-port regression, with an offline negative that reds on a single-cycle end-state (AC-4)
  Fixture seeds Cycle 1 + instructs `### Feedback Cycles` per validation round; `assertRejectionFlow` now requires ≥2 implementation reports + ≥2 cycle entries (dropped the single-cycle `status: implementation` check); `assertClaude/CodexReviewerReuse` parse the SendMessage/send_input tool-call shape targeting the validation reviewer. `TestRejectionFlowNegativeSingleCycle` reds a single-cycle end-state on the second-cycle check.
- DONE: Doc-lock holds — the `<!-- seed-scenarios -->` block + dev README shared-scenario list both name `feedback-3-cycle-escalation` (the doc-lock test + `TestSharedScenarioDocsContract` green); secret-free `go test ./...` green
  Authored `TestSeedScenariosDocLock` binding the seed block IDs to `sharedRuntimeScenarios()` (no such doc-lock test existed before — the doc only promised one); drift falsifiability confirmed (renaming the doc ID reds it). `go test ./...` = 1112 passed / 15 packages; live lane compiles + the 4 no-model live meta-tests pass. Live legs (AC-3 both hosts, AC-4 live cycle) carry their `live <ci-run:|session:>` citation at validation/terminal per the live-AC policy.
- SKIPPED: budget-probe fail-safe scenario
  Out of scope per the dispatch and ideation decision (b): deferred to `feedback-guarantee-binary-gate` (non-mechanizable ruling would return it as `feedback-budget-fresh-dispatch`). No budget scenario added here.

### Summary

Authored the `feedback-3-cycle-escalation` shared scenario (host-neutral durable-state escalation oracle + both host runners + offline negatives) and evolved `rejection-flow` to drive the full 2-cycle trajectory with host-specific reviewer-reuse, restoring the regression the Go port dropped. Bound the seed-scenarios doc block to the code table with a new doc-lock test (none existed; the spec assumed an "existing" one). Design decision worth flagging: a "cycle" is graded as one validation round, so both fixtures record `- Cycle N:` per round consistently, and the evolved `assertRejectionFlow` drops the brittle `status: implementation` check because the correct 2-cycle end-state sits at validation-passed, not implementation. The OFFLINE secret-free lane is fully green (1112/15 pkgs); the LIVE legs (AC-3 both hosts, AC-4 live cycle) are authored, compile under `-tags live`, and are cited at validation/terminal per the live-AC policy — not burned during implementation.

## Stage Report: validation

- DONE: AC-1 parity — `TestSharedScenarioRunnerCoverage` reds on a removed host runner in BOTH directions; `TestSharedRuntimeScenarioDefinitions` green
  Both meta-tests green under `-tags live` (no-model). Falsifiability driven: removing the Codex runner for `feedback-3-cycle-escalation` reds with `"feedback-3-cycle-escalation" has no Codex runner`; removing the Claude runner reds with `… has no Claude runner`. Both edits reverted; worktree clean at 7e03d5ea.
- DONE: AC-2 offline negatives — `TestThirdCycleEscalationNegativeAutoBounce` + `TestAssertThirdCycleEscalation` red the 4th-auto-bounce and stalled-at-cycle-2 end-states
  Both green. Gutting `assertThirdCycleEscalation` to `return nil` reds BOTH negatives (`negative_test.go:116`, `shared_assertions_test.go:61`) — proven non-tautological behavioral oracles, not transcript spelling checks. Reverted.
- DONE: AC-4 offline — `TestRejectionFlowNegativeSingleCycle` reds a single-cycle end-state
  Green. Weakening the second-cycle check from `< 2` to `< 1` cycle entries reds it at `negative_test.go:87` (the second-cycle check) plus `TestAssertRejectionFlow` at `shared_assertions_test.go:28` — confirms the dropped `status: implementation` check did NOT weaken it; the ≥2-cycle-entry check is the load-bearing replacement and the correct 2-cycle end-state sits at validation-passed (not implementation), so the old status check would have been wrong. Reverted.
- DONE: AC-5 doc-lock — `TestSeedScenariosDocLock` + `TestSharedScenarioDocsContract`
  Both green. Renaming the doc seed ID reds `TestSeedScenariosDocLock` (`docs_test.go:58`, set-equality so symmetric on code-side drift too); breaking a README live-command clause reds `TestSharedScenarioDocsContract` (`docs_test.go:106`). Reverted.
- DONE: Full offline `go test ./...` green
  1112 passed / 15 packages, twice (start and after all falsifiability edits reverted). `go vet -tags live ./internal/ensigncycle` clean.
- DONE: Confirm assertions grade DURABLE entity-body state, never transcript phrasing
  `assertThirdCycleEscalation(entity)` takes only the entity string (cycle-entry count + escalation marker + impl-report count — all on-disk). `assertRejectionFlow`'s load-bearing checks are over `entity`; `observed` only does weak `reject`/`implementation` substring checks. A validation-only probe confirmed the seeded fixtures do NOT pre-satisfy their assertions (seeded escalation entity = 2 cycles / no marker / reds; seeded rejection entity reds), so a live pass requires the real producer to drive the durable state — no free pass.
- DONE: Confirm reviewer-reuse assertion requires a real SendMessage/send_input targeting the VALIDATION reviewer, not just any tool call
  A live-tagged probe drove `assertClaudeReviewerReuse`/`assertCodexReviewerReuse` over crafted transcripts: both RED on narration-only ("validation" in prose), RED on a different tool (Task/dispatch) targeting validation, the Claude one REDs on a SendMessage to a non-validation target, and both PASS only on a real `SendMessage`/`send_input` tool-call event referencing the validation reviewer. Probe removed; worktree clean.
- DONE: AC-3 + AC-4 live legs — assess the `live <ci-run:|session:>` citation path; recommend drive-now vs close-post-merge
  RECOMMENDATION: close-post-merge. The dev workflow README declares no `require-external-proof`, so `resolveExternalProofPolicy` = OFF and the `live-run guard` in `handlers.go:251` does not gate terminalization. Additionally the AC-3/AC-4 `Verified by:` clauses lead with test names, not `live`, so `classifyLiveACs` would not flag them even if the guard were on (no enforced `pending-live-run` state). The offline halves are sound + falsifiable (above) and the seeded fixtures don't pre-pass, so the live legs genuinely exercise the producer when driven. Live harness is runnable in this session (local `~/.claude/benchmark-token` + `~/.codex/auth.json` present; auth-decision units green → would launch, not skip) but ~32 wall-min × real model across 4 multi-cycle runs is disproportionate inside validation. Matches the established convention: sibling `fo-auto-continues` (commit 8931b90a) passed validation with the live half explicitly recommended post-merge for the same not-opted-in reason.

### Summary

PASSED. Every AC's OUTSIDE-body evidence was reproduced by exercising the behavior, not re-reading: the full offline suite is green (1112/15) and each named oracle was driven to RED on the exact broken end-state it guards (gutted/weakened assertion, removed host runner, drifted doc/README clause), proving they are behavioral oracles rather than substring tautologies. The two dispatch-flagged risks are clear: dropping the `status: implementation` check did NOT weaken `assertRejectionFlow` (the ≥2-cycle-entry check reds the single-cycle regression, and the correct 2-cycle end-state sits at validation-passed), and the reviewer-reuse assertions require a real SendMessage/send_input tool-call targeting the validation reviewer (red on prose, red on other tools, red on wrong target). No AC has self-referential-only proof. AC-3/AC-4 are runtime-observable; their offline halves pass and the live legs are recommended close-post-merge because this workflow does not opt into the live-run guard — the same path the sibling `fo-auto-continues` validation took. Polish-only nits (non-blocking): the spec body's line 55 says "three" host-neutral scenarios while the seed block now lists four; and stale AC-number labels in test comments (`shared_coverage_meta_test.go` "AC-2/AC-3", `shared_scenarios_docs_test.go` "AC-6") don't match this entity's AC-1/AC-5 numbering — cosmetic, no behavioral gap. Awaiting the detached adversarial auditor's Material findings, which join this gate.

### Feedback Cycles

**Cycle 1 — validation REJECTED via detached adversarial audit (2026-06-04).** The validator recommended PASSED on the offline suite; the detached audit (separate checkout of `7e03d5ea`) found two material test-strength holes the green suite structurally cannot see. Routed to `implementation` (reused impl ensign, `reuse_ok` 14.3%, same worktree). The live legs stay close-post-merge per the live-AC policy — unaffected by this cycle.

Material (fix both):
- **A — `assertThirdCycleEscalation` report-count check has no isolating negative** (`internal/ensigncycle/shared_assertions_impl_test.go`, the `reports > 1` clause). Mutation: changing it to `reports > 99999` (deleting the no-post-cycle-3-report guard) leaves the whole suite GREEN — every existing negative carrying the stray report ALSO lacks the escalation marker, so each reds on the marker check first and never exercises the report-count clause. Add an isolating negative: marker present + cycle-3 + a stray post-cycle-3 implementation report and no other defect, which ONLY the report-count check rejects (mirror the merge-hook suite's dedicated isolating case at `shared_scenarios_negative_test.go:167`).
- **B — `assertThirdCycleEscalation` does not check the entity stayed non-terminal.** A marker-written-but-terminalized end-state (`status: done`) PASSES — the FO wrote the escalation marker AND terminalized the entity (auto-resolving instead of parking for the human; the escalation prompt says "do not advance to done"). Add a park-not-advance check (status NOT terminal) mirroring `assertGateHeld`/`assertMergeHookGuardHeld`, with its own isolating negative.

Polish (fold in):
- Section-scope the escalation-marker + `- Cycle N:` matches to the `### Feedback Cycles` section (currently matched anywhere in the body).
- Add a cheap OFFLINE table test for `assertClaude/CodexReviewerReuse` (currently exercised only by the model-spending live runners) using the wrong-recipient / unrelated-tool / loose-narration shapes the audit + the validator's live-probe already exercised.
- Validator nits: `docs/specs/scenario-testing-principles.md:55` says "three" host-neutral scenarios while the seed block now lists four; stale AC-number labels in `shared_coverage_meta_test.go` / `shared_scenarios_docs_test.go` comments.

Re-run the validator after the fix; keep offline `go test ./...` green.

## Stage Report: implementation (cycle 2)

- DONE: Material A — isolating negative for `assertThirdCycleEscalation`'s no-post-cycle-3-report check
  Added `markerWithStrayReport` to `TestThirdCycleEscalationNegativeAutoBounce` (`shared_scenarios_negative_test.go`): marker present + 3 cycle entries + a stray post-cycle-3 `## Stage Report: implementation`, no other defect — only the report-count check rejects it. Mutation control: changing `reports > 1` → `reports > 99999` reds EXACTLY this case at `negative_test.go:151`, suite otherwise green (the validator's exact finding, now guarded). Commit b337483c.
- DONE: Material B — `assertThirdCycleEscalation` asserts the entity stayed NON-terminal (park-not-advance)
  Added `terminalStatus` regex + the `status: done` check mirroring `assertGateHeld`/`assertMergeHookGuardHeld`; added `markerButTerminalized` isolating negative (marker + 3 cycles + 1 impl report + `status: done`, no other defect). Mutation control: deleting the non-terminal check reds EXACTLY this case at `negative_test.go:170`, suite otherwise green.
- DONE: Polish — section-scope the `- Cycle N:` + escalation-marker matches to the `### Feedback Cycles` section
  Added `feedbackCyclesSection(entity)` (heading → next heading/EOF) + `nextHeading` regex; both the cycle-count and marker checks now run over the section, not the whole body. Added the out-of-section isolating case to `TestAssertThirdCycleEscalation`. Mutation control: matching the whole body reds EXACTLY that case at `shared_assertions_test.go:96`.
- DONE: Polish — offline table test for `assertClaude/CodexReviewerReuse`
  Moved both reuse asserts to a default-tag file `shared_reviewer_reuse_test.go` (stdlib-only, removed now-unused `encoding/json`/`fmt` imports from both live runners; `strings` retained in the Codex runner) and added `shared_reviewer_reuse_table_test.go`: wrong-recipient / unrelated-tool / loose-narration / empty all RED, real tool-call PASSES. Mutation control: stubbing `assertClaudeReviewerReuse` to `nil` reds all four negative subtests.
- DONE: Polish — validator doc/comment nits
  Dropped the stale "three" count in `scenario-testing-principles.md:55` (the seed block + `TestSeedScenariosDocLock` are the count source of truth — avoids the same drift recurring on the next scenario); realigned stale AC labels: `shared_coverage_meta_test.go` AC-2/AC-3 → AC-1, `shared_scenarios_docs_test.go` AC-6 → AC-5, negatives-file header generalized off the colliding AC-5 prefix.

### Summary

Cycle-2 fixes both Material items the detached audit found, each with its mutation control proven by exercising (not asserting): the no-post-cycle-3-report check and the new park-not-advance (non-terminal) check each now have a dedicated isolating negative that reds ONLY that check, confirmed by mutating each check and watching exactly the one case go red while the rest of the suite stays green. Folded in all polish: section-scoped the cycle/marker matches to the `### Feedback Cycles` section (with its own mutation-controlled isolating case), moved the host reviewer-reuse asserts to the offline lane behind a real table test (wrong-recipient/unrelated-tool/loose-narration all red), and cleared the doc/comment nits. Offline `go test ./...` green (1124/15 pkgs); the live lane compiles and the no-model live meta + reuse tests pass. The live legs (AC-3 both hosts, AC-4 live cycle) are unchanged and stay close-post-merge per the live-AC policy.

## Stage Report: validation (cycle 2)

- DONE: Material A closed — `markerWithStrayReport` isolating negative for the no-post-cycle-3-report check; mutation control
  Re-ran in the worktree at b337483c (did not trust the report). `reports > 1` → `reports > 99999` reds EXACTLY `negative_test.go:151` (markerWithStrayReport), and ONLY `TestThirdCycleEscalationNegativeAutoBounce` — rest of suite green. The 4th-auto-bounce case (no marker) still reds on the marker check first, which is why the isolating case is the one that proves the report-count clause. Reverted.
- DONE: Material B closed — park-not-advance (non-terminal) check + `markerButTerminalized` isolating negative; mutation control
  Deleting the `terminalStatus` check reds EXACTLY `negative_test.go:170` (markerButTerminalized: marker + 3 cycles + 1 report + status: done), and ONLY that case — rest of suite green. Confirms the non-terminal check is independently covered. Reverted.
- DONE: Polish — section-scoped cycle/marker matches with mutation control; offline reviewer-reuse table test with stub-to-nil control; doc/comment nits
  Section-scope: replacing `feedbackCyclesSection` with `return entity` (defeating scoping) reds EXACTLY `shared_assertions_test.go:96`, suite otherwise green. Reviewer-reuse: `assertClaude/CodexReviewerReuse` now live in default-tag `shared_reviewer_reuse_test.go` (live lane still compiles + still calls them at claude_live_runner_test.go:154 / codex_live_runner_test.go:175); stubbing BOTH to `return nil` reds EXACTLY the 4 error-expecting subtests per host (loose-narration / unrelated-tool / wrong-recipient / empty) while the real-tool-call subtest still passes. Doc/comment nits cleared: scenario-testing-principles.md:55 no longer says "three"; comments realigned (`shared_coverage_meta_test.go` → AC-1 parity, `shared_scenarios_docs_test.go` → AC-5 README-contract). All reverted; worktree clean at b337483c.
- DONE: Offline `go test ./...` green
  1124 passed / 15 packages (twice: before and after all mutation edits reverted, zero residual diff). `go vet -tags live ./internal/ensigncycle` clean. AC-1 parity-guard regression re-checked post-rework: removing the Codex runner still reds `TestSharedScenarioRunnerCoverage` at `shared_coverage_meta_test.go:25`.

### Summary

PASSED (cycle 2). Both Material audit findings are closed and independently verified by re-running in the worktree, not by trusting the implementation report: each mutation control reds EXACTLY the single isolating negative the team-lead named (Material A → `negative_test.go:151`; Material B → `negative_test.go:170`; section-scope polish → `shared_assertions_test.go:96`; reviewer-reuse stub-to-nil → the 4 negative subtests per host), with the rest of the suite green in every case — proving the new checks are not green-on-everything tautologies. The reviewer-reuse asserts moved to the offline lane without breaking the live wiring, and both cycle-1 Polish nits I flagged (the stale "three" count and stale AC-number comment labels) are now cleared. Full offline suite green 1124/15; live lane compiles. The AC-3/AC-4 live legs are unchanged and remain close-post-merge per the existing live-AC policy — team-lead handles the run-now-vs-defer call at the gate.
