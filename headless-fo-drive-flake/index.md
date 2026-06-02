---
id: hdr9dd8qywdc42cq700vefy4
title: Headless FO intermittently doesn't drive the cycle (-p / teams) — root-cause the stall
status: ideation
source: "FO (2026-06-02): live-e2e flakiness — under spacedock claude -p with EXPERIMENTAL_AGENT_TEAMS=1 the FO intermittently stalls (one model per run): on 38's PR sonnet booted, loaded the contract (Skill->Read), then went silent (~70s quiet trip) while opus drove the full cycle. The deferred FO-runtime-await question, now confirmed real. Investigation; gated on the transcript-artifact follow-up for data."
started: 2026-06-02T07:57:40Z
completed:
verdict:
score: "0.34"
worktree:
issue:
---

The live-e2e flakiness root: under headless `spacedock claude -p` with teams enabled, the FO (first-officer) intermittently does NOT drive the cycle to done — one model per run stalls (failures alternated sonnet/opus across the 0.19.3 + 38 PRs). On 38's PR (#261) the new `streamWatcher` caught it cleanly: sonnet booted, started loading the operating contract (`Skill → Read`), then went silent (the ~60s quiet budget tripped at ~70s) while opus drove the full cycle green. So the binary + the watcher are sound; the FO MODEL sometimes stalls headless.

This is the FO-runtime-await / headless-drive question deferred from 38 — now confirmed real. It is an INVESTIGATION: root-cause WHY the FO stalls headless, then decide a fix.

## Dependency
GATED on `live-e2e-transcript-artifact` (3g — merge first): the investigation needs the FAILING run's transcript, which today is lost (gh truncation + binary-only artifact). Ideation may plan the approach now; data-gathering + root-cause execution follow the artifact landing + a captured stall.

## The contract contradiction this investigation centers on (ideation finding)

Reading the live launch against the FO contract surfaces a concrete, code-grounded reason the headless FO is intermittently free to take a stall-prone path. Two INDEPENDENT mode-selection axes collide in the live launch, and the contract never states their precedence:

1. **Team availability** (`skills/first-officer/references/claude-first-officer-runtime.md:9-12`): when `ToolSearch(select:TeamCreate)` finds TeamCreate — which it does, because the CI live job sets `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` (`.github/workflows/runtime-live-e2e.yml:91`) — `TeamCreate` MUST be the first team-mode call, and dispatch runs in teams mode (async `Agent(team_name=…)`, completions return via `task_notification`/inbox poll).
2. **Single-entity mode** (`skills/first-officer/references/first-officer-shared-core.md:75`): activates when the session is non-interactive (`-p`) AND the prompt names a specific entity. The live launch is exactly this — `live_test.go:95-105` shells `spacedock claude … -- -p … <task>` where `task` (`live_test.go:77-79`) names the specific backlog entity `make-it-work.md` and says "by dispatching an ensign for it." Single-entity mode's runtime adapter (`claude-first-officer-runtime.md:26`) says: *"skip team creation. Use bare-mode dispatch … the Agent tool without `team_name` blocks until the subagent completes, which prevents premature session termination in `-p` mode."*

So the contract simultaneously tells the headless FO (a) teams is available → TeamCreate first, async dispatch, and (b) this is single-entity `-p` → skip team creation, bare blocking dispatch. The precedence is **undocumented**, so the model reconciles it per-run. When it takes the teams branch, it can async-dispatch and end its `-p` turn awaiting a `task_notification` the headless loop may never re-wake it for; when it takes the single-entity/bare branch, the blocking `Agent()` keeps the turn alive and it drives to done. This is the mechanism that makes the intermittent stall plausible WITHOUT positing pure model nondeterminism — and it is the hypothesis the captured transcript must confirm or refute.

A second, aggravating observation: the binary's bare-mode advisory (`internal/dispatch/build.go:120-127` → `internal/claudeteam/claudeteam.go:73-79`) WARNS a bare_mode dispatch that lacks recent TeamCreate evidence and nudges "run ToolSearch select:TeamCreate and TeamCreate first" — i.e., it pushes the FO TOWARD the teams branch, which under single-entity `-p` is the stall-prone path. This is not the root cause (it fires only after a bare dispatch is already chosen) but it is design context the fix stage must weigh.

## Hypotheses to test (in priority order)

- **H1 — Single-entity-vs-teams precedence is unresolved → the FO takes the teams branch under `-p` and ends its turn mid-await.** The strongest hypothesis given the contradiction above. Distinguishing evidence in the transcript: a `TeamCreate` `tool_use` followed by an async `Agent(team_name=…)` `tool_use`, then the FO's turn ENDS with no further assistant turn while an ensign is still in-flight (no `task_notification` re-wake) — the quiet budget then trips. The 38 watcher even bakes in this assumption: its first watched beat is `expect(isTeamCreate)` (38 entity body), i.e. the harness EXPECTS teams mode under `-p`.
- **H2 — The stall is at contract-load, before any dispatch decision (genuine model nondeterminism mid-`Skill→Read`).** The confirmed-real symptom on #261 was sonnet going silent WHILE loading the operating contract (`Skill → Read`), before reaching TeamCreate at all. Distinguishing evidence: the last streamed activity is the contract-load `Skill`/`Read` tool_use (or its result) with NO subsequent TeamCreate/dispatch tool_use before the quiet trip. If the transcript shows this, the cause is model-side, not the precedence contradiction (H1), and the disposition path of AC-2 is the right answer.
- **H3 — The `-p` turn ends mid-await even in a non-team path (headless loop never re-wakes the FO for any teammate signal).** A harness/runtime variant of H1 where the await semantics, not the mode choice, are the defect. Distinguishing evidence: the FO reached a blocking point that SHOULD have re-woken it (a tool_result or notification arrived) but the `-p` turn did not resume. Lower prior — bare `Agent()` is documented to block synchronously — but the transcript must rule it out rather than assume it away.

The three are mutually distinguishable by WHERE the last streamed activity sits relative to TeamCreate and the dispatch: H2 = before TeamCreate (contract-load); H1 = after an async team dispatch, turn-end mid-await; H3 = after a blocking point that failed to resume. The 38 watcher's labelled `stepTimeout` (`"TeamCreate"` / `"dispatch close"` / `"expect_exit"`) already localizes which of the three regions the stall fell in — that label plus the transcript tail is the primary discriminator.

## Data-gathering plan (executes AFTER 3g lands + a stall is captured)

1. **Capture N≥3 failing-run transcripts.** 3g (`live-e2e-transcript-artifact`) makes the streamed FO transcript survive a failed/killed CI-E2E run as an uploaded artifact (today it is lost to gh log truncation + binary-only upload). Re-run the live matrix (sonnet on CI-E2E, opus on CI-E2E-OPUS) until ≥3 stalls are captured — the flake is per-run and alternates models, so a handful of matrix runs should yield several. Each captured artifact carries the watcher's labelled `stepTimeout`/`stepFailure` + the stream-json tail.
2. **Classify each transcript against H1/H2/H3** using the discriminator above: read the watcher's step label and the last 3-5 stream-json entries before the quiet trip. Tabulate which hypothesis each stall matches.
3. **Distinguish runtime/harness cause from model nondeterminism.** A runtime/harness cause (H1/H3) shows a STRUCTURAL pattern — the same stall point correlated with the teams-vs-single-entity branch the FO took (visible in whether a TeamCreate + async Agent appears). Model nondeterminism (H2) shows the stall scattered at contract-load with no dispatch-decision dependence. If ≥2 of 3 captures match H1's teams-branch-under-`-p` signature, the contradiction is the root cause and the candidate fix (force single-entity bare dispatch under `-p`) is justified. If they scatter at contract-load (H2), the disposition path is justified instead.
4. **Cross-check against the upstream Python net.** The upstream harness (`~/git/spacedock/scripts/test_lib.py` FOStreamWatcher, referenced by 38) drove the same FO contract reliably — compare whether upstream forced bare/single-entity dispatch under `-p` (a no-`team_name` Agent) where v1 leaves it to the contract. A divergence there is direct evidence for H1.

## Acceptance criteria

**AC-1 — The headless-FO stall is root-caused from captured evidence.** At least 3 failing-run transcripts (from 3g's uploaded CI-E2E artifact) are captured and each classified against H1/H2/H3 by its watcher step-label + stream-json tail; the analysis names the dominant stall region (contract-load / TeamCreate / dispatch-await / exit) and states whether the cause is the single-entity-vs-teams precedence contradiction (runtime/harness) or model nondeterminism.
Verified by: the ≥3 captured transcript artifacts attached to this entity (or referenced by CI run ID + artifact name) + a classification table in the entity body mapping each capture to a hypothesis and stall region; an empty or single-capture analysis fails this AC (the per-run flake demands multiple samples to separate structure from noise).

**AC-2 — A fix lands with a test, or a justified disposition is recorded with a guard.** EITHER (fix path) a runtime/contract change makes the headless FO drive reliably — candidate: the contract states single-entity `-p` mode takes precedence over team availability and forces bare blocking `Agent()` dispatch (no `team_name`) under `-p`, per the contract's own "prevents premature session termination in `-p` mode" guidance — AND the change is exercised by a test that fails without it; OR (disposition path) a documented decision to accept model nondeterminism backed by 38's fast-localized detection + a bounded automatic retry of the live cycle, AND the retry is implemented + tested.
Verified by — fix path: a test (offline unit over the dispatch/mode-selection seam, OR the AC-1 classification re-run green on a follow-up CI-E2E matrix showing the stall no longer reproduces across ≥3 runs) that goes red without the change and green with it. Disposition path: the bounded-retry code + an offline test asserting the retry fires on a first-attempt `stepTimeout` and the run passes when a retry drives to done; PLUS the recorded decision in the entity body naming why nondeterminism (H2) is accepted over a forced-mode fix.

**AC-3 — No speculative runtime change precedes the evidence.** This entity's PR contains NO change to `skills/first-officer/references/` or `internal/cli/frontdoor.go` whose justification is not one of the AC-1 captured transcripts. The investigation (AC-1) is committed before any fix (AC-2); a fix commit that predates the captured-transcript analysis is a process violation.
Verified by: the PR's commit order — the AC-1 transcript-classification commit precedes any FO-runtime/frontdoor change commit; if AC-2 takes the disposition path, there is no FO-runtime/frontdoor change at all.

## Spike (riskiest unverified mechanism)

The riskiest unknown for this investigation is whether 3g's artifact actually yields a transcript rich enough to discriminate H1/H2/H3 — i.e., whether the streamed stream-json names the FO's TeamCreate/Agent tool_use and turn boundaries with enough fidelity to tell "ended turn mid-await" from "stalled at contract-load." **This spike CANNOT run before 3g lands and a stall is captured** — it is the literal first step of the data-gathering plan (capture one transcript, confirm it discriminates). Recorded determination: the smallest end-to-end exercise of the riskiest path is "capture ONE failing transcript via 3g and confirm its step-label + tail places the stall in exactly one of the three regions" — that is AC-1's first capture, and it gates whether the rest of AC-1 (N≥3, classification) is even possible. If the first capture does NOT discriminate (e.g., the tail is too coarse), the investigation bounces to 3g to enrich the artifact before proceeding. This is why this entity is GATED on 3g and is investigation-first.

## Test plan

| Layer | What it proves | Cost | Fixture/CLI/live |
|---|---|---|---|
| Live capture (CI-E2E, gated on 3g) | AC-1: ≥3 failing transcripts captured + classified; the dominant stall region + cause named | model spend, several matrix runs | CI-E2E sonnet + CI-E2E-OPUS, transcript artifact from 3g |
| Offline unit (fix path) | AC-2 fix: the dispatch/mode-selection seam forces bare blocking dispatch under `-p` single-entity (test red without the contract/runtime change, green with it) | seconds, no model | dispatch-build input asserting `bare_mode`/no-`team_name` under the single-entity `-p` shape |
| Offline unit (disposition path) | AC-2 disposition: a first-attempt `stepTimeout` triggers exactly one bounded retry and a retried run that drives to done passes | seconds, no model | synthetic watcher result + retry harness (extends 38's `streamWatcher` unit split) |
| Re-run reliability (fix path, optional binding proof) | AC-2 fix: the stall no longer reproduces across ≥3 follow-up CI-E2E runs | model spend | CI-E2E matrix re-run |

Cost/complexity: the investigation (AC-1) is read-and-classify over captured artifacts — cheap once 3g lands. The fix-or-disposition (AC-2) is small either way: the fix is a contract-precedence clause + the seam test; the disposition is a bounded-retry wrapper + its unit test. Neither is started before AC-1's evidence.

## Notes
- Touches the FO runtime (`skills/first-officer/references/` + `internal/cli/frontdoor.go`) — high-stakes; investigation FIRST, no speculative runtime change (AC-3).
- GATED on `3g` (`live-e2e-transcript-artifact`) — merge first; the data-gathering + root-cause EXECUTION follow 3g landing + a captured stall.
- The candidate fix (single-entity `-p` precedence over team availability → forced bare blocking dispatch) is grounded in the contract's existing "prevents premature session termination in `-p` mode" guidance (`claude-first-officer-runtime.md:26`); it is a CANDIDATE for AC-2, not a foregone conclusion — the captured transcript decides between it and the disposition path.

## Stage Report: ideation

- DONE: Plan the root-cause investigation: lay out the hypotheses and the data-gathering plan
  Three prioritized, mutually-distinguishable hypotheses (H1 single-entity-vs-teams precedence unresolved → teams branch under `-p` ends turn mid-await; H2 contract-load model nondeterminism; H3 `-p` await never re-wakes) added with per-hypothesis transcript discriminators tied to the 38 watcher's step labels; data-gathering plan captures N≥3 failing transcripts via 3g's artifact, classifies each, and distinguishes structural (runtime/harness) from scattered (model) cause + an upstream-Python cross-check.
- DONE: AC: root-cause with evidence PLUS a fix-or-justified-disposition
  AC-1 binds root-cause to ≥3 captured transcripts + a classification table (empty/single-capture fails); AC-2 offers the fix path (contract precedence forcing bare blocking dispatch under single-entity `-p`, test red-without/green-with) OR the disposition path (accept nondeterminism + a tested bounded retry); each AC's "Verified by" names a checkable artifact/test, not prose review.
- DONE: Scope + dependency: high-stakes FO runtime — investigation FIRST, no speculative runtime change; GATED on 3g
  AC-3 enforces investigation-before-fix via PR commit order (no FO-runtime/frontdoor change justified by anything but a captured transcript); Spike + Dependency + Notes record the GATE on 3g and that the riskiest mechanism (does 3g's artifact discriminate H1/H2/H3?) is the literal first capture and cannot run before 3g lands.

### Summary

Hardened the provisional spec into an investigation-first ideation grounded in the actual code. The load-bearing finding: the live launch (`live_test.go` shells `spacedock claude -- -p <task-naming-make-it-work.md>` with `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` set in CI) puts the FO under TWO colliding mode-selection axes — team availability (TeamCreate-first) vs single-entity `-p` mode (skip team creation, bare blocking dispatch) — whose precedence the contract never states, so the model reconciles per-run; the teams branch under `-p` is the stall-prone path that ends the turn mid-await. That contradiction becomes H1 (the strongest hypothesis), with H2 (contract-load nondeterminism, the confirmed #261 symptom) and H3 (await-never-rewakes) as the alternatives, all three distinguishable by where the last streamed activity sits relative to TeamCreate/dispatch — which the 38 watcher's labelled stepTimeout already localizes. The riskiest unknown (whether 3g's artifact is rich enough to discriminate) is honestly recorded as un-spikeable until 3g lands, so the entity stays gated on 3g and execution follows a captured stall. ACs name checkable artifacts/tests; AC-3 enforces no speculative FO-runtime change via PR commit order.
