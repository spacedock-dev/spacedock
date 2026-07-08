---
title: "codex-live rejection-flow oracle accepts the current multi_agent_v2 spawn-evidence family (lane red on unchanged main)"
status: validation
sprint:
source: "0250 Commander session 2026-07-06/07. Timeline: Jul-4 main run green; 2026-07-06 k7 branch run, branch rerun, AND a baseline rerun of the Jul-4 green main run (same code) all red at codex_live_runner_test.go:42 'no validation spawn_agent found — the FO never created a cycle-1 reviewer to reuse'. codex-cli 0.142.5 identical across green and red runs, so not a CLI version change: service-side behavior drift under --enable multi_agent_v2 (the harness opts in deliberately and asserts the flag, codex_live_runner_test.go:283,399-400). Transcript evidence (artifacts, run 28806717191 + run 28693890413 attempt 2): the FO ran a REAL two-worker flow — separate implementation worker and validation reviewer, followup_task cycle-2 reuse, populated agents_states on collab wait calls, correct end-state markers — but zero spawn_agent-shaped events; its spawn mechanism now emits evidence the oracle does not recognize. Blocks every 0250 member's merge (all-lanes-green is the DoD's proof policy)."
started: 2026-07-06T23:10:29Z
completed:
verdict:
score: 0.6
worktree: .worktrees/spacedock-ensign-codex-live-spawn-oracle-v2-family
issue:
id: 3v5cd2fcx8y5bk1dcxt7swzf
mod-block:
pr: "#483"
---

The codex-live shared-scenario runner proves worker separation by grepping the exec stream for spawn_agent-shaped events; under current multi_agent_v2 service behavior a genuinely-separated flow emits a different spawn/collab evidence family, so the lane is red on unchanged main and cannot go green for any diff. Direction: harden the spawn-evidence detection (internal/ensigncycle codex runner, with internal/dispatch/codex_v2_adapter.go's v2 modeling) to accept the current v2 family as reviewer-existence proof while REMAINING a real separation oracle — it must still go red when no separate reviewer existed (its whole point is catching an FO that validates inline; narration must never satisfy it). Evidence anchors for ideation: green-run artifacts (run 28693890413 attempt 1, Jul-4) vs red-run artifacts (attempt 2 + run 28806717191 attempts 1-2) — diff what spawn evidence looked like green vs red before writing the new oracle. High-stakes surface (a live lane's own tests → CI machinery): detached adversarial audit required, including a constructed no-separate-reviewer transcript the hardened oracle must reject. Acceptance sketch: value — codex-live returns green on unchanged main scenarios with the hardened oracle (baseline: red today), and 0250 members' merges unblock; mechanism — the oracle change ships with the anti-tautology negative case.

## Ideation notes (paused 2026-07-07)

PAUSED by captain: multi_agent_v2 compatibility is deferred out of the current sprint. This section parks the perishable evidence (the green-vs-red spawn-evidence diff) so a resumer does not have to re-fetch artifacts that may be purged. The fuller reframe + provisional ACs below (Spike / Problem statement / Proposed approach / Acceptance criteria) are unreviewed and paused — treat as a resumption seed, not an approved plan.

Perishable diff facts (from the Jul-4 attempt-1 green vs the red runs; artifacts were on `/tmp` and are transient):

- Corpus: 13 codex transcripts — green run 28693890413 attempt-1 (Jul-4, PASSED); red run 28693890413 attempt-2 (main baseline); red run 28806717191 (k7 branch fo-boot-engage-split) — each across the six shared scenarios.
- Spawn/collab surface, ALL 13 transcripts: 0 `spawn_agent`, 0 `followup_task`, 0 `send_input`, 0 populated `agents_states`. The only `collab_tool_call` that appears is `tool:"wait"` with empty `receiver_thread_ids` and empty `agents_states`. codex-cli 0.142.5 emits `item.type:"error"` "Under-development features enabled: multi_agent_v2" at stream start. There is NO distinct v2 spawn-evidence family — green and red are structurally identical.
- rejection-flow `assertCodexReviewerReuse` trace (reproduced by reimplementing its branches):
  - GREEN 28693890413 a1: addressableWorkerAbsent=false, validation spawn_agent count=0, advance-fallback initV=1 advV=1 implA=1 kept=true waits=8 → PASS (via `codexReviewerReuseViaAdvanceCommands`, never via spawn detection).
  - RED 28693890413 a2: spawn=0, initV=0 advV=1 implA=2 kept=false waits=8 → RED "no validation spawn_agent found".
  - RED 28806717191 k7: spawn=0, initV=0 advV=1 implA=2 kept=false waits=8 → RED (identical).
- Why red flaps (phrasing variance, not regression): the red FO built the initial validation dispatch via stdin-JSON (`printf '{…"stage":"validation"…}' | dispatch build`) so the command has no `--stage validation` substring → initV=0; and narrated "kept validation reviewer" / "kept cycle-1 reviewer" without the literal token "kept-alive" → kept=false. Green happened to use the flag form and said "kept-alive".
- Durable end-state (`assertRejectionFlow`: fix marker, ≥2 implementation reports, ≥2 `### Feedback Cycles`, output surfaces reject+implementation) PASSES in green AND both reds. The workflow completes correctly in red; only the phrasing-keyed producer-signal fails.
- Indirect separation evidence, identical green-to-red: the FO's main-thread `file_change`s to the entity are limited to the two feedback-cycle bookkeeping writes (green: item_68 Cycle 1 REJECTED, item_94 Cycle 2 PASSED). The implementation and validation stage reports land in the entity with no main-thread `file_change` and no shell redirection — authored by worker contexts on threads codex exec `--json` does not stream. So real separation is indirectly evidenced and unchanged between green and red; the redness is a brittle-heuristic flap in `codexReviewerReuseViaAdvanceCommands`, keyed on exact FO command-phrasing and narration wording.

Bottom line for resumption: there is nothing to "accept" — no v2 spawn family exists on the codex exec surface. When multi_agent_v2 compat is picked back up, the real question is whether to de-brittle the codex producer-signal and honestly scope what codex can machine-verify (Direction 1 below), or wait for a codex CLI that surfaces spawn/thread/agents_states metadata in exec mode.

## Spike: green-vs-red evidence diff (riskiest mechanism, run first)

The riskiest path was the seed's premise itself — "a genuinely-separated flow now emits a DIFFERENT spawn/collab evidence family the oracle must learn to accept." I diffed the Jul-4 GREEN rejection-flow transcript (run 28693890413 attempt-1) against both RED transcripts (run 28693890413 attempt-2 main baseline; run 28806717191 k7 branch), reimplementing the `assertCodexReviewerReuse` branches to see exactly why each verdict fell. The premise did not survive.

Findings (all from observed events, not guessed):

- **No spawn-evidence family exists — green and red are structurally identical.** Across all 13 available codex transcripts (green + both reds × six scenarios) there are ZERO `spawn_agent`, `followup_task`, `send_input` events and ZERO populated `agents_states`. The only `collab_tool_call` that ever appears is `tool:"wait"` with empty `receiver_thread_ids` and empty `agents_states`. codex-cli 0.142.5 `exec --json --enable multi_agent_v2` emits an `item.type:"error"` "Under-development features enabled: multi_agent_v2" at stream start and surfaces no multi-agent spawn/thread/collab metadata in exec mode.
- **The green pass never touched spawn detection.** Green had 0 validation `spawn_agent`s; it passed solely via the `codexReviewerReuseViaAdvanceCommands` fallback (initial `dispatch build --stage validation` flag-form ×1, `--stage validation --advance` ×1, `--stage implementation --advance` ×1, "kept-alive validation reviewer" reuse narration present, waits ≥2).
- **Red fails on FO phrasing variance, not a behavior regression.** Both red runs did the same workflow but (a) built the initial validation dispatch via stdin-JSON (`printf '{…"stage":"validation"…}' | dispatch build`), so the command string has no `--stage validation` substring → `initV=0`; and (b) narrated "kept validation reviewer" / "kept cycle-1 reviewer" without the literal token "kept-alive" → `kept=false`. The fallback therefore misses them and `assertCodexReviewerReuse` errors at "no validation spawn_agent found — the FO never created a cycle-1 reviewer to reuse."
- **The durable end-state passes in red.** `assertRejectionFlow` (host-neutral: fix marker applied, ≥2 implementation reports, ≥2 `### Feedback Cycles` entries, output surfaces reject+implementation) passes in both green and red. The workflow completes correctly in red; only the producer-signal flaps.
- **Real separation is indirectly evidenced, and identical green-to-red.** In both green and red the FO's main-thread `file_change`s to the entity are limited to the two feedback-cycle bookkeeping writes (green item_68 Cycle 1 REJECTED, item_94 Cycle 2 PASSED). The implementation and validation stage reports land in the entity with no main-thread `file_change` and no shell redirection — authored by separate worker contexts on threads codex exec does not stream. The redness is not a separation regression; it is a brittle-heuristic flap.

Conclusion: there is nothing to "accept" — no v2 spawn-evidence family. The lane is red because `codexReviewerReuseViaAdvanceCommands` grades exact FO command-phrasing and narration wording, both non-deterministic under an LLM FO. This is the anti-pattern the codebase already names in `assertThirdCycleEscalation`: "Grading on these on-disk facts, never on transcript phrasing, is what keeps this a behavioral oracle rather than a tautology." Fixtures captured: green 28693890413 a1, red 28693890413 a2, red 28806717191 (rejection-flow `codex-exec.jsonl`) — these seed the implementation's first tests.

## Problem statement (corrected by the spike)

`assertCodexReviewerReuse` claims to prove reviewer separation/reuse, but on the current codex `exec --json --enable multi_agent_v2` surface there is no observable spawn/collab evidence to grade. Its `spawn_agent` path is dead (0 spawns in every run) and its command+narration fallback is brittle: it red-flags correct separated flows whenever the FO phrases dispatch commands via stdin-JSON or narrates "kept" instead of "kept-alive". This reds the codex lane on unchanged main and blocks every 0250 member (all-lanes-green DoD). The fix must restore stable green on correct behavior WITHOUT collapsing into a pure narration tautology, and must be honest about what codex exec can and cannot machine-verify.

## Proposed approach (Direction 1, recommended)

Rework the codex rejection-flow producer-signal to grade what the surface reliably provides:

1. Keep the durable, host-neutral end-state (`assertRejectionFlow`) as the primary lane grade — it passes today and is what codex can actually prove.
2. Replace `codexReviewerReuseViaAdvanceCommands` with a phrasing-robust assignment-separation check: the FO built a DISTINCT validation-stage dispatch assignment, separate from the implementation-stage dispatch, accepting BOTH command forms (stdin-JSON payload `"stage":"validation"` and flag `--stage validation`) by parsing the dispatch stage from payload/flags, not from prose. Drop the "kept-alive" narration substring requirement entirely.
3. Honest scoping: document (reworked-assertion code comment + this task) that codex exec cannot machine-verify a distinct reviewer PROCESS (no spawn/thread/agents_states surface); the codex lane grades assignment-separation + durable end-state, and the process-level separation oracle lives in the Claude runner (real `Agent`/`SendMessage` evidence). Precedent: shallow-boot's host-neutral-vacuous checks and the `assertThirdCycleEscalation` "on-disk facts, not phrasing" principle.

Residual gap, stated plainly: a transcript that builds a separate validation dispatch but then fabricates the verdict inline without running a worker is NOT distinguishable on codex exec. The Claude runner covers that. Direction 1 catches the collapse that IS observable — validation performed without a distinct validation assignment.

## Alternative (Direction 2, spike-gated, NOT recommended as primary)

Key on the indirect signal the spike found: in a genuine flow the FO's main-thread entity writes are limited to feedback bookkeeping and the stage reports appear with no main-thread `file_change`; RED a transcript whose main-thread `file_change`s include stage-report authorship. Risk: `file_change` items carry path+kind only (no diff), so the oracle can only COUNT entity writes, not classify them, and an inline FO could self-write via shell redirection (no `file_change` at all). Buildable only if a spike first confirms codex reliably emits distinguishable evidence for FO self-writes vs worker writes. If that spike fails, Direction 2 is unbuildable.

## Acceptance criteria

- **AC-1 (VALUE, measured vs baseline):** The codex-live rejection-flow lane returns GREEN on unchanged main. Baseline that can move wrong: RED today — `assertCodexReviewerReuse` errors at "no validation spawn_agent found" on run 28693890413 a2 and run 28806717191. End-state: replaying the three captured transcripts (green a1, red a2, red 28806717191 — all genuine separated flows with correct durable end-state) through the reworked assertion yields 0 failures, versus 2 failures at baseline.
  - Test: offline table test (default build tags, no model spend) feeding the three captured JSONL fixtures; assert all three pass. Plus one live `-tags live` codex rejection-flow run confirming green end-to-end (~4 min).
- **AC-2 (ANTI-TAUTOLOGY negative, named):** A constructed no-separate-reviewer transcript goes RED. It performs validation WITHOUT a distinct validation-stage dispatch (the implementation dispatch is advanced to emit the PASSED verdict, and/or the FO authors the validation stage report via a main-thread `file_change`), while narrating "separate validation reviewer" and calling `wait`. The reworked assertion returns a non-nil error.
  - Test: offline negative fixture in `shared_reviewer_reuse_table_test.go` (the negative-fixture pattern already lives there); assert the exact error message (pristine-output capture).
- **AC-3 (phrasing robustness):** The separation verdict is independent of FO command-form and narration wording. Fixtures building the validation dispatch via stdin-JSON and via `--stage validation` flag both PASS; a fixture identical except narration says "kept" vs "kept-alive" does not change the verdict.
  - Test: offline table fixtures over both command forms and both narration phrasings; verdict stable across all.
- **AC-4 (honest-scope doc, on-disk):** The reworked assertion's code comment and this task record that codex exec `--enable multi_agent_v2` (0.142.5) surfaces no spawn/thread/collab metadata, so the codex lane grades assignment-separation + durable end-state while the Claude runner owns process-level separation.
  - Test: a doc/comment-presence check (contractlint-style or a focused unit assertion) that the scoping note is present.

## Test plan

- Offline table tests (default tags, no spend, minutes): 3 replay fixtures (green + 2 red) all pass; ≥2 anti-tautology negatives all red; phrasing-robustness fixtures. The captured rejection-flow JSONL become checked-in testdata.
- Live smoke (`-tags live`, needs `OPENAI_API_KEY` + codex on PATH): 1–2 codex rejection-flow runs → green (~4 min each).
- Riskiest mechanism already exercised (this spike); implementation's first test is the green-replay-passes / negative-reds pair.

## Doc / coordination note (0250 merge order)

The fix is test-only inside `internal/ensigncycle` (no product-surface, no CLI/doc-site change), so it carries no doc-diff dependency and can merge ahead of the behavioral 0250 members. Recommend merging the oracle repair FIRST; once the codex lane is green again, k7/z25/zm/vcm/tv unblock in their existing wave order under the all-lanes-green DoD.

## Flag to FO/captain

This ideation reframes the seed. The seed asked to "accept the current v2 spawn-evidence family"; the spike shows no such family exists (green and red are spawn-evidence-free and identical) and the redness is a brittle-heuristic phrasing flap. The entity is captain-PAUSED as of 2026-07-07, so the corrected direction (de-brittle + honestly scope; do NOT chase a nonexistent spawn family) needs captain confirmation before implementation. Strange things are afoot at the Circle K: the fastest lane repair (Direction 1) makes the codex separation oracle deliberately weaker than the seed implies, leaning on the Claude runner for real process-level separation proof — that trade should be an explicit captain decision, not a silent one.

## Stage Report: ideation

- DONE: Green-vs-red evidence diff FIRST (the riskiest mechanism): pull run 28693890413 attempt-1 (Jul-4 green) codex rejection-flow artifacts and diff their spawn/collab event shapes against the red attempts (28693890413 attempt-2, 28806717191); the new oracle's accepted family must be derived from observed events, not guessed
  Reimplemented `assertCodexReviewerReuse`'s branches over all three transcripts: green, red-a2, red-28806717191 all have 0 spawn_agent/followup_task/send_input and empty agents_states. The diff overturned the premise — there is no distinct v2 spawn family; green passed on the brittle command+narration fallback, red fails on stdin-JSON dispatch phrasing + "kept" vs "kept-alive" narration.
- DONE: Design keeps the oracle falsifiable: accepted v2 evidence family proves a SEPARATE reviewer existed, and a constructed no-separate-reviewer transcript (FO validates inline, narrates workers) still goes RED — the anti-tautology negative is a named AC with a test plan
  AC-2 names the anti-tautology negative (validation without a distinct validation-stage dispatch → RED) with an offline table-fixture test plan. Honestly bounded: codex exec cannot distinguish "built a validation dispatch then faked the verdict inline"; that residual is the Claude runner's job, documented in AC-4 and the approach.
- DONE: ACs include the lane-level value proof (codex-live green on unchanged main scenarios, baseline red today) and the doc/coordination note for 0250 members' merge order once the lane repairs
  AC-1 measures 0 vs 2 assertion failures on the three replayed transcripts against the RED-today baseline plus a live green smoke; the doc/coordination note records the test-only-in-internal/ensigncycle merge-first ordering that unblocks k7/z25/zm/vcm/tv.

### Summary

The riskiest-path spike (green-vs-red evidence diff) invalidated the seed premise: codex `exec --json --enable multi_agent_v2` (0.142.5) surfaces no spawn/collab evidence in any of the 13 available transcripts, so green and red are structurally identical and the Jul-4 green never passed via spawn detection — it rode the brittle `codexReviewerReuseViaAdvanceCommands` fallback that red flaps on purely by FO phrasing. Reframed the task to de-brittle the codex producer-signal (assignment-separation + durable end-state, phrasing-robust) and honestly scope codex's limits, with a falsifiable anti-tautology negative and a value AC measured against the RED-today baseline. Flagged to captain: the corrected direction makes the codex oracle deliberately weaker than the seed implied and the entity is PAUSED, so it needs captain confirmation before implementation.

### Feedback Cycles

- Family tracking note (FO, 2026-07-07, pre-cut audit NB#2 — for the codex-side resumer): two further surface-variance instances joined this family after the pause, recorded in PR comments: VARIANT 2 (PR #478 comment) — shallow-boot red from an FO call-shape fumble on `state sweep` flag placement, plus the earlier keep-moving grader blind spot fixed in-sprint by crediting the `status --set status=done` dispatch surface (zm cycle-1/2); VARIANT 3 (PR #480 comment) — keep-moving "did not dispatch ready-one" false-negative because that run's FO terminalized via `merge guard`, so the status-set-done surface never appeared as an FO command; transcript proved the parallel batch dispatch happened. The unifying lesson for the resumption: codex FOs reach the same durable end-state through MULTIPLE legitimate command surfaces (direct status --set, merge guard, collab threads) — oracles must key on durable state or the full surface family, never one command shape. See also seed sc (live-runner-boot-preamble-hardening) for the adjacent boot-preamble flake class.

## Stage Report: implementation

- DONE: Reworked the Codex live rejection-flow oracle around durable state plus observable assignment/command surfaces instead of spawn_agent-only evidence.
  `assertCodexReviewerReuse` now keeps the real spawn/thread path when present, but when Codex exec emits no spawn metadata it parses `dispatch build` assignment surfaces from both `--stage ...` flags and stdin JSON payloads. The fallback requires a validation-stage assignment before the first review, a validation-stage assignment for the re-review, and an implementation rework assignment, while `assertRejectionFlow` continues to grade the durable two-cycle entity body.
- DONE: Taught keep-moving Codex assertions to credit the merge-guard terminalization surface while preserving false-negative/false-positive negative fixtures.
  `codexKeepMovingTrace` now treats `spacedock merge guard <slug>` as terminal dispatch evidence for the approved and independent entities, and as a forward-drive violation for the corrected entity. Existing permission-question, silent-park, missing-dispatch, and corrected-forward-drive negatives remain in the same table, and the new PR #480 merge-guard positive is covered there.
- DONE: Added replay/fixture tests from the current red run and documented the honest Codex exec multi_agent_v2 observability boundary.
  Added an offline stdin-JSON multi_agent_v2 assignment-surface fixture, a no-validation-assignment anti-tautology negative with an exact error assertion, and a source-comment presence check documenting that Codex exec does not prove a distinct reviewer process. The code comment records the boundary: Codex grades assignment-separation plus durable end-state; Claude remains the process-level separation oracle.

### Summary

Implemented and committed the Codex live oracle repair on `spacedock-ensign/codex-live-spawn-oracle-v2-family` at `eb3b0a77` (`harden codex live oracle surfaces`). Verification passed: `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and targeted live smoke `go test -tags live ./internal/ensigncycle -run 'TestLiveCodexSharedScenarios/rejection-flow' -count=1` (2 passed).

## Stage Report: validation

- DONE: Reproduce the Codex oracle tests, including replay/negative fixtures and keep-moving merge-guard coverage.
  Focused `go test ./internal/ensigncycle -run 'TestAssertCodexReviewerReuse|TestAssertCodexKeepMoving' -count=1 -v` passed 32 tests; `go test ./...` and `go test ./... -race` each passed 2042 tests across 17 packages.
- DONE: Inspect the assertion diff for anti-tautology: durable-state/assignment surfaces pass while no-validation-assignment and forward-drive negatives still fail.
  Offline positives and negatives pass, including the no-validation-assignment exact-error check and the keep-moving corrected-forward-drive negative; the source-comment presence check covers the Codex exec observability boundary.
- FAILED: Run or explicitly account for the targeted live Codex rejection-flow smoke, then append a PASSED/REJECTED validation report.
  `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-validation-codex-live go test -tags live ./internal/ensigncycle -run 'TestLiveCodexSharedScenarios/rejection-flow' -count=1 -v` failed after 450s at `assertCodexReviewerReuse`: durable state reached validation/PASSED with both feedback cycles, but no validation dispatch-build assignment surface was observed.

### Summary

Recommendation: REJECTED. The offline oracle tests are green and preserve the intended negative cases, but the live Codex rejection-flow lane is still red on the repaired assertion. The preserved live artifact shows Codex produced the durable two-cycle end state, yet the only top-level `dispatch build` command was a failed implementation build (`missing required field 'checklist'`), with no validation assignment surface for the new oracle to accept.

### Feedback Cycles

- Cycle 1: REJECTED — targeted live Codex rejection-flow still reds after 450s: durable state reached validation/PASSED with both feedback cycles, but the repaired oracle found no validation dispatch-build assignment surface to accept.
- Cycle 2: REJECTED at PR #483 CI — failed-jobs-only retry still red on codex-live `TestLiveCodexSharedScenarios/rejection-flow` with `Codex foreground-wait watchdog typed stall` (`arm=silent-after-wait`, `durable_progress=false`) and Claude opus `TestLiveClaudeSharedScenarios/{filing,feedback-3-cycle-escalation}` with wrong-root boot outside fixture cwd. Captain directed: send back, diagnose root cause, and get local green before retrying CI.
- Cycle 3: REJECTED at PR #483 CI run 28906653062 on repaired head `749b7db6`. Deterministic lanes, offline, and `pi-live` passed, but `codex-live` failed before scenarios: `codex plugin add spacedock@spacedock` returned `failed to create plugin target directory: Filename too long (os error 36)`. `claude-live (sonnet, CI-E2E)` also failed `Run live ensign cycle` while continuing into shared scenarios; `claude-live (opus)` was still running shared scenarios at escalation time. Per the feedback-rejection-flow cycle-3 rule, FO escalated to captain instead of auto-dispatching another implementation round.
- Cycle 3 override: Captain directed FO to send it back unless reframing is needed and to focus on the Codex side. Scope is the `codex-live` filename-too-long plugin-install failure first; Claude sonnet failures are context only unless the Codex root cause intersects.

## Stage Report: implementation

- DONE: Root-caused the failed targeted live Codex rejection-flow artifact where durable state passes but no validation assignment surface is emitted.
  The cycle-1 implementation still routed the live retry through `assertCodexReviewerReuse(got.result.jsonl)` after `assertRejectionFlow` had already proven the durable two-cycle state. That kept the final gate JSONL-only even when Codex exec omitted spawn/thread/assignment metadata. A preserved local repro attempt was blocked before JSONL creation by a Codex ENOSPC/cache error, but the validator's preserved artifact summary was enough to identify the missing handoff: durable validation worker output was available in `entityAfter`, yet the oracle never received it.
- DONE: Revised the oracle and fixtures so the actual live surface passes while no-reviewer/inline-validation negatives still fail without relying on transcript narration.
  Added `assertCodexReviewerReuseWithDurableState`, used only by the Codex rejection-flow retry wrapper. JSONL spawn/reuse and assignment surfaces still win when present; contradictory JSONL evidence still fails. Only the missing-assignment error can fall back to durable validation-stage worker output, requiring a `## Stage Report: validation` heading plus recorded `Cycle 1: REJECTED` and `Cycle 2: PASSED` feedback. Added a live-surface-style positive with no validation assignment surface and a narration-heavy inline-validation negative that lacks a validation report and stays red.
- DONE: Reran focused offline/full/race checks and the targeted live rejection-flow smoke, then appended this cycle-2 implementation report.
  Verification passed: focused durable fallback/retry tests, `go test ./internal/ensigncycle`, `go test ./...` (2044 passed across 17 packages), `go test ./... -race` (2044 passed across 17 packages), `gofmt -w ./cmd ./internal`, and targeted live smoke `go test -tags live ./internal/ensigncycle -run 'TestLiveCodexSharedScenarios/rejection-flow' -count=1` (2 passed). Code commit: `4e9cd99` (`accept codex durable reviewer evidence`).

### Summary

Cycle 2 repaired the live-only gap by letting the Codex live rejection-flow oracle consume durable validation-stage worker output when Codex exec emits no spawn or validation-assignment surface. The anti-tautology boundary remains: transcript-only tests still reject narration without a validation assignment, and the durable fallback rejects inline validation when no validation stage report exists.

## Stage Report: validation (cycle 2)

- DONE: Reproduce the cycle-2 focused durable-fallback/retry tests, full suite, race suite, and targeted live Codex rejection-flow smoke.
  Focused `go test ./internal/ensigncycle -run 'TestAssertCodexReviewerReuse|TestRunCodexRejectionFlowRetry|TestRejectionFlowRejectsWaitReuseTranscriptWithoutDurableSecondCycle' -count=1 -v` passed 39 tests; final-tree `go test ./...` and `go test ./... -race` each passed 2043 tests across 17 packages; live `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-validation-codex-live-cycle2 go test -tags live ./internal/ensigncycle -run 'TestLiveCodexSharedScenarios/rejection-flow' -count=1 -v` passed in 515.017s.
- DONE: Verify the durable fallback only applies when JSONL lacks usable spawn/assignment evidence and still rejects inline/narration-only validation.
  Focused tests exercised the no-assignment durable positive, inline-without-validation-report negative, narration-without-validation-assignment negative, and retry narrowness cases; live artifact JSONL also showed validation dispatch-build surfaces at lines 70-71 and 107-108, so the live pass did not require narration-only acceptance.
- DONE: Append a cycle-2 validation report with PASSED/REJECTED recommendation and exact evidence.
  Recommendation: PASSED. Oracle implementation commit `4e9cd995` remains intact; required `gofmt -w ./cmd ./internal` produced formatter-only code commit `70697e9b` in `internal/release/journeydelta.go`.

### Summary

Recommendation: PASSED. Cycle-2 validation reproduced the focused durable-fallback and retry coverage, the full and race baselines, and a live Codex rejection-flow smoke that completed successfully with preserved artifacts. The anti-tautology boundary is still covered by executable tests: durable fallback accepts missing JSONL assignment evidence only when durable validation-stage worker output exists, and inline or narration-only validation remains red.

### Addendum: pre-gate cleanup

Removed the validation-authored formatter drift from `internal/release/journeydelta.go` with cleanup commit `63cbf1dd`, leaving the implementation commits' `internal/ensigncycle` behavior intact. Verification: `git diff --name-only 4e9cd995..HEAD` returned no paths; `git diff --name-only origin/next..HEAD -- internal/release/journeydelta.go internal/ensigncycle` listed only `internal/ensigncycle` files and did not list `internal/release/journeydelta.go`; `go test ./internal/ensigncycle -run 'TestAssertCodexReviewerReuse|TestRunCodexRejectionFlowRetry|TestRejectionFlowRejectsWaitReuseTranscriptWithoutDurableSecondCycle' -count=1` passed 39 tests. I did not rerun the expensive live smoke because this cleanup did not change `internal/ensigncycle`.

## Stage Report: implementation

- DONE: Root-caused and repaired the PR #483 retry failures from the current CI surface family.
  The Codex rejection-flow stall was not a reviewer-oracle regression: the live harness created isolated `CODEX_HOME` under system temp, and current Codex refuses to create helper aliases there (`Refusing to create helper binaries under temporary dir "/tmp"`), leaving spawned workers unable to run exec tools while the FO waited. The Claude opus filing/escalation reds were detector false positives: one stream attempted `spacedock new` from the parent directory, failed with `no commissioned workflow found`, then corrected with `--workflow-dir` into the fixture; another read the plugin skill root README outside the fixture, which is not a workflow-root boot.
- DONE: Implemented compatibility fixes in `internal/ensigncycle` only.
  `newCodexLiveIsolatedHome` now creates isolated Codex homes outside system temp, trying a non-temp artifact root, user cache, then repo-local `.spacedock-live-codex` fallback. The wrong-root detector now scans all tool-use blocks, correlates suspected wrong-root commands with tool results, ignores failed off-fixture attempts that recover, treats explicit in-fixture `--workflow-dir` as the target even when the shell cwd is outside the fixture, and excludes plugin skill root README reads.
- DONE: Added focused regression coverage for both CI failure families.
  Added Codex home-parent/candidate tests so temp artifact roots are not reused as `CODEX_HOME`, plus wrong-root detector positives for plugin skill README reads and failed-parent-new-then-corrected-workflow-dir streams. Existing wait-watchdog, rejection-flow retry, and reviewer-reuse tests remain green.
- DONE: Verified local green for the repair and recorded the local live-host boundary.
  Verification passed on commit `749b7db6` (`repair live runtime CI surfaces`): `go test ./internal/ensigncycle -run 'TestCodexLiveHomeParent|TestDetectWrongRootBoot|TestCodexCollabWaitWatchdog|TestRunCodexRejectionFlowRetry|TestAssertCodexReviewerReuse' -count=1` (1.361s), final-tree `go test ./internal/ensigncycle -count=1` (7.271s), `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` with unrelated `internal/release/journeydelta.go` formatter drift removed before commit. Targeted Codex live smoke passed: `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-ci-repair-live-codex-final go test -tags live ./internal/ensigncycle -run 'TestLiveCodexSharedScenarios/rejection-flow$' -count=1 -v` (400.752s). Targeted Claude opus live probe cannot run locally because Claude launches fail before FO work with `apiKeySource:"none"` and `401 Invalid authentication credentials`; the detector changes for the opus CI surfaces are covered by the offline `TestDetectWrongRootBoot` regressions.

### Summary

Implementation commit `749b7db6` repairs the live Codex worker startup environment and de-flakes the Claude wrong-root detector against the PR #483 observed stream shapes. Code is committed locally only; no code remote push was performed. The Codex live rejection-flow smoke is green locally, while Claude opus live execution is blocked here by missing/invalid local Anthropic auth before the workflow starts.

## Stage Report: validation (cycle 3)

- DONE: Reproduce the implementation's local-green evidence for commit `749b7db6`: focused failure-family tests, `go test ./internal/ensigncycle -count=1`, full suite, race suite, and the targeted Codex live rejection-flow smoke if local credentials/tooling permit.
  Focused `go test ./internal/ensigncycle -run 'TestCodexLiveHomeParent|TestDetectWrongRootBoot|TestCodexCollabWaitWatchdog|TestRunCodexRejectionFlowRetry|TestAssertCodexReviewerReuse' -count=1` passed 68 tests; `go test ./internal/ensigncycle -count=1` passed 292 tests; `go test ./...` and `go test ./... -race` each passed 2072 tests across 17 packages; Codex live `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-validation-codex-live-spawn-oracle-v2-family-cycle3 go test -tags live ./internal/ensigncycle -run 'TestLiveCodexSharedScenarios/rejection-flow$' -count=1 -v` passed in 370.359s.
- DONE: Verify root-cause coverage for both PR #483 failure families: Codex temp `CODEX_HOME`/helper alias stall and Claude opus wrong-root detector false positives.
  Codex coverage includes `TestCodexLiveHomeParentUsesUserCacheOutsideSystemTemp`, `TestCodexLiveHomeParentCandidatesIncludeRepoFallbackWhenArtifactsAreTemp`, and the green live artifact; opus detector coverage includes `plugin_skill_readme_outside_fixture_passes`, `failed_parent_new_then_corrected_workflow_dir_passes`, and the pre-existing PR #446 stream replays in the focused detector run.
- DONE: Check anti-tautology boundaries from the previous cycle remain covered; do not accept a fix that weakens inline/narration-only negatives.
  `TestAssertCodexReviewerReuseRejectsNarratedReviewerWithoutValidationAssignment` and `TestAssertCodexReviewerReuseWithDurableStateRejectsInlineValidationWithoutReport` remain in the focused 68-test run alongside the current-v2 and durable-state positives.
- DONE: Inspect the branch diff for scope, confirm code is committed locally but not pushed, and append a PASSED/REJECTED validation report with exact evidence.
  Repair diff `06d75e51..749b7db6` is six `internal/ensigncycle` files; code status is clean at `749b7db6`, `git rev-list --left-right --count origin/spacedock-ensign/codex-live-spawn-oracle-v2-family...HEAD` reports `0 1`, and no code branch push was performed.

### Summary

Recommendation: PASSED. The local repair is good enough to retry PR #483 CI: the targeted Codex live rejection-flow smoke is green, the temp-home and wrong-root false-positive families have deterministic regression coverage, and the prior anti-tautology negatives remain executable. `gofmt -w ./cmd ./internal` still wants to realign `internal/release/journeydelta.go`, but that drift comes from pre-existing parent commit `06d75e51` and was restored after the check; it is not introduced by repair commit `749b7db6`.

## Stage Report: implementation

- DONE: Root-caused the PR #483 cycle-3 Codex `Filename too long` failure with path evidence.
  CI run 28906653062/job 85755161978 showed the job-level Codex setup succeeded from `/home/runner/work/_temp/spacedock-codex-marketplace/plugins/spacedock`, but the Go live harness then added a per-test marketplace at `/tmp/TestLiveCodexSharedScenarios693401381/002` and failed at `codex plugin add spacedock@spacedock` before scenarios. The cycle-3 repair had added a fallback `CODEX_HOME` under the plugin checkout (`<repo>/.spacedock-live-codex`) for hosts where user cache creation fails. Codex plugin install copies the checkout into `CODEX_HOME/plugins/cache/spacedock/spacedock/0.24.0`; putting that cache inside the checkout creates a self-copy/path-length trap. A local forced reproduction with repo-local `CODEX_HOME=/Users/clkao/git/spacedock-research/spacedock-v1/.spacedock-live-codex/repro-codex-home` failed during `codex plugin add` with `File name too long`.
- DONE: Implemented the smallest Codex-side hardening.
  Commit `30f2fb46` (`keep codex live home outside checkout`) changes only `internal/ensigncycle`: the last fallback is now repo-adjacent, not repo-local. For the CI shape, fallback moves from `/home/runner/work/spacedock/spacedock/.spacedock-live-codex` to `/home/runner/work/spacedock/.spacedock-live-codex/spacedock`, so Codex can copy the checkout without copying its own plugin cache. The temp artifact-root and user-cache candidate ordering stays unchanged.
- DONE: Added a failing regression first, then verified it green.
  The TDD red test initially failed because candidates were `["/blocked/cache/spacedock-live-codex", "/Users/clkao/git/spacedock/.spacedock-live-codex"]` while the expected fallback was outside the checkout. The final test, `TestCodexLiveHomeParentCandidatesKeepFallbackOutsidePluginCheckoutWhenArtifactsAreTemp`, uses the CI checkout shape and asserts no `CODEX_HOME` parent is equal to or under the plugin checkout.
- DONE: Local green, including a CI-like plugin-add reproduction.
  CI-like reproduction passed with marketplace root `/tmp/TestLiveCodexSharedScenarios693401381/002`, plugin source pointing at the worktree, and repo-adjacent `CODEX_HOME=/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/.spacedock-live-codex/spacedock-ensign-codex-live-spawn-oracle-v2-family/repro-codex-home`; `codex plugin add spacedock@spacedock` installed to `.../plugins/cache/spacedock/spacedock/0.24.0`. Verification passed: `go test ./internal/ensigncycle -run 'TestCodexLiveHomeParent|TestDetectWrongRootBoot|TestCodexCollabWaitWatchdog|TestRunCodexRejectionFlowRetry|TestAssertCodexReviewerReuse' -count=1` (1.326s), `go test ./internal/ensigncycle -count=1` (7.190s), `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` with the unrelated `internal/release/journeydelta.go` drift removed. Targeted Codex live smoke passed: `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-ci-repair-live-codex-path-final go test -tags live ./internal/ensigncycle -run 'TestLiveCodexSharedScenarios/rejection-flow$' -count=1 -v` (393.228s).

### Summary

Recommendation: retry PR #483 after validation. The Codex cycle-3 failure was a live-harness setup path bug introduced by the repo-local fallback, not a reviewer-oracle weakening. The repair keeps isolated Codex homes outside system temp and outside the checkout being installed as a plugin, with a regression test covering the CI path shape.
