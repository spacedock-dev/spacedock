---
title: "codex-live rejection-flow oracle accepts the current multi_agent_v2 spawn-evidence family (lane red on unchanged main)"
status: implementation
sprint:
source: "0250 Commander session 2026-07-06/07. Timeline: Jul-4 main run green; 2026-07-06 k7 branch run, branch rerun, AND a baseline rerun of the Jul-4 green main run (same code) all red at codex_live_runner_test.go:42 'no validation spawn_agent found — the FO never created a cycle-1 reviewer to reuse'. codex-cli 0.142.5 identical across green and red runs, so not a CLI version change: service-side behavior drift under --enable multi_agent_v2 (the harness opts in deliberately and asserts the flag, codex_live_runner_test.go:283,399-400). Transcript evidence (artifacts, run 28806717191 + run 28693890413 attempt 2): the FO ran a REAL two-worker flow — separate implementation worker and validation reviewer, followup_task cycle-2 reuse, populated agents_states on collab wait calls, correct end-state markers — but zero spawn_agent-shaped events; its spawn mechanism now emits evidence the oracle does not recognize. Blocks every 0250 member's merge (all-lanes-green is the DoD's proof policy)."
started: 2026-07-06T23:10:29Z
completed:
verdict:
score: 0.6
worktree: .worktrees/spacedock-ensign-codex-live-spawn-oracle-v2-family
issue:
id: 3v5cd2fcx8y5bk1dcxt7swzf
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
