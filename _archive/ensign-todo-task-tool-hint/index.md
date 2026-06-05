---
id: 3c0bcrn8p60wptvc12fsv5x2
title: Hint ensigns to use todo/task tools when available — spike the placement + measure empirical impact
status: done
source: "captain (2026-06-04) — empirically ensigns never convert the helper-built `### Completion checklist` into internal todos (TodoWrite=0 across all ensign transcripts; neither the ensign contract nor the dispatch file hints it). Hint ensigns to use todo/task tools IF the runtime has them — but spike the placement and show the empirical performance impact before committing."
score: "0.27"
started: 2026-06-04T22:27:44Z
completed: 2026-06-05T01:11:11Z
verdict: rejected
worktree:
issue:
archived: 2026-06-05T01:11:11Z
---

The `spacedock dispatch build` helper assembles a `### Completion checklist` into every ensign dispatch, but ensigns consume it as prose and never externalize it as a tracked todo/task list (verified this session: `TodoWrite`=0 across all ensign transcripts; `skills/ensign/` and the dispatch file carry no todo/task guidance). A todo/task hint could improve mid-stage captain visibility (Shift+Down into the ensign pane shows live progress) and possibly ensign completeness/focus — but only where the runtime actually has such tools, and only if it earns its token/wallclock cost. Hint the ensign to use them **if available**.

## Problem

- Ensigns don't use todo/task tools; the checklist stays prose-only and progress is invisible until the stage report.
- Tool availability is **runtime-specific**: Claude has `TodoWrite` (internal) + team `TaskCreate`; Codex / Pi may have different or no equivalents. A blanket hint that names a Claude-only tool in the shared core would break host-neutrality (the ensign portability oracle).

## Placement options to spike (the design fork)

- **(a) Shared ensign contract** (`skills/ensign/references/ensign-shared-core.md`) — one host-neutral hint ("if your runtime offers a todo/task tool, track the checklist in it; otherwise account for it in the stage report as today"). Pro: one place, loaded every dispatch. Con: generic phrasing only (cannot name `TodoWrite` — portability oracle); may nudge runtimes that lack the tool.
- **(b) Runtime-specific contracts** (`claude-ensign-runtime.md` / codex / pi) — name the actual tool per host. Pro: precise, host-neutral by construction. Con: duplicated across runtimes; loads even when not wanted.
- **(c) Dispatch-builder injection** — `spacedock dispatch build` appends a todo-hint line (exactly as it already appends the standing-teammates routing block via `show-standing`) ONLY when the target host has the tool. Pro: conditional on host AND checklist presence, most precise, zero always-on contract cost. Con: adds dispatch-builder logic + a host→tool capability map.

## Riskiest unknown — spike in ideation (captain: show empirical performance impact)

The unverified question is NOT the placement mechanics — it is **whether the hint actually changes ensign behavior, and at what cost.** Ideation MUST run an empirical comparison before committing, not assert the hint is good:

- Dispatch a real ensign **with** the todo-hint vs **without**, on a comparable stage, and measure from the ensign transcripts using the per-turn jsonl token+wallclock parse already exercised this session (the boot-analysis method):
  - does the ensign actually invoke the tool (`TodoWrite`/task count > 0)?
  - stage-report completeness / checklist accounting (does every item get DONE/SKIPPED/FAILED; fewer missed items)?
  - token + wallclock delta (does externalizing the checklist cost or save)?
- Record the measured numbers in the entity body. If the hint shows no measurable benefit (or a net cost), say so — but per dev policy a decision-only outcome belongs in the roadmap, so this task should land a hint in the winning placement justified by the measurement, or be re-scoped.

## Spike results (ideation, measured 2026-06-04)

### What the premise got wrong

The seed claim "ensigns never convert the checklist into todos (`TodoWrite`=0)" is true for `TodoWrite` but **false in spirit**. The two real opus ensigns this session ran with NO hint (k0py7yvv team) split:

- **wm implementation ensign** (`fo-auto-continues-after-stage-completion`, `agentId a00fe2018e673e1a0`): spontaneously externalized its checklist into the team task tool — 5 `TaskCreate` + 10 `TaskUpdate`, full in_progress→completed lifecycle, no prompting. `TodoWrite`=0, `task`=15.
- **gq implementation ensign** (`feedback-nonhappy-live-coverage`, `agentId ad7f1a4e1bfa3c4e5`): used zero task/todo tools; kept the checklist in-context. `TodoWrite`=0, `task`=0.

Neither received the harness's "task tools haven't been used recently" nudge (confirmed: 0 occurrences across all 37 agent transcripts of the live session) — so the wm ensign's task use was model initiative, not a host nudge. The real tool reached is `TaskCreate`/`TaskUpdate` (the team task system), NOT `TodoWrite`. A hint that names `TodoWrite` would target the wrong tool.

### Controlled A/B (the captain's explicit ask)

A clean isolation of the hint variable: identical 10-item Go-coding stage staged in a throwaway git repo, headless `claude -p --output-format stream-json` on **opus** (the baseline ensign model), 3 trials/arm, the ONLY difference being one host-neutral hint line ("If your runtime offers a todo or task tracking tool, externalize this checklist into it…"). Metrics parsed per-turn from the stream jsonl (same `usage`/tool-use parse as the boot analysis). Harness, runbook, and parser: `/tmp/ensign-todo-spike/ab/` (`run_arm.sh`, `parse_stream.py`).

| metric | WITHOUT hint (n=3) | WITH hint (n=3) | delta |
|---|---|---|---|
| externalized checklist (task/todo>0) | **0/3** | **2/3** | hint raises rate |
| mean wallclock | 110.7 s | 131.1 s | **+18%** |
| mean cost | $0.59 | $0.87 | **+48%** |
| mean turns | 15.7 | 32.7 | **+108%** |
| mean output tokens | 1213 | 1656 | +37% |
| mean report accounting (of 11) | 10.7 | 10.7 | **no change** |
| work correctness (`go test ./...`) | 3/3 pass | 3/3 pass | **no change** |

Per-trial: with-1 task=22/$1.06, with-2 task=25/$1.10, with-3 task=0/$0.45; without-1/2/3 task=0, $0.49–$0.73. All six landed complete, correct, test-passing code with an honest stage report.

### The decisive finding — the hint shows a net cost, no measured benefit

Splitting by behavior rather than arm isolates the cost OF externalizing from the hint itself: trials that externalized averaged **$1.08**, trials that stayed inline **$0.56** — a **+94% cost premium** for externalizing, with **identical** report completeness (10.7/11) and **identical** work quality (all passed). On the three metrics the captain named — tool-invocation count, checklist-accounting completeness, token+wallclock — the hint **buys more tool calls and more cost and zero completeness/quality gain.** It does not make ensigns more complete or more correct; it makes them externalize, and externalizing is the expense.

The one benefit the A/B *cannot* measure is the captain's stated motive: **mid-stage visibility** (Shift+Down into a live ensign's pane shows tracked-task progress vs. a silent prose run). That benefit exists ONLY in Claude team mode (the pane + `TaskCreate` both require it) and is invisible to a transcript metric. So the honest trade is: the hint spends ~$0.50 and ~20s and doubles tool calls per externalizing ensign to buy live progress visibility — a Claude-team-only, human-UX benefit — and buys nothing on completeness or correctness.

### What the data justifies (placement)

Because the only benefit is Claude-team-mode visibility and the cost is real, the cost must fall ONLY where the benefit exists. That rules out the always-on surfaces:

- **(a) shared ensign core** — loads on every dispatch including bare-mode and codex/pi, where there is no Shift+Down pane and (for the no-team case) no `TaskCreate`. It would impose the measured cost on contexts with zero benefit, and can only use generic wording (the portability/leakage oracles forbid naming `TaskCreate`). Rejected: pays the cost everywhere, benefits nowhere outside Claude teams.
- **(b) runtime adapters** — names the tool per host but still loads for every Claude dispatch including bare-mode (no team, no pane). Better than (a) but still over-broad, and duplicates the line across adapters. Rejected.
- **(c) dispatch-builder injection, host-AND-team-conditional** — `dispatch build` appends the hint to the dispatch body ONLY when `host == "claude"` AND `team_name != ""` (team mode), the exact context where the pane-visibility benefit and `TaskCreate` both exist. Bare-mode and codex/pi dispatches never carry it, so the measured cost lands only where the benefit does. The shared core is untouched → AC-3 holds by construction. The mechanism is already shipped: `internal/dispatch/build.go:245` resolves `host`, and the codex/pi body branches (`firstActionBlock`, `completionSignalBlock`) prove host drives body content today (`TestBuildCodexHostPromptShape`/`TestBuildPiHost*` green). **Chosen.**

The riskiest *placement* mechanism — host-conditional body emission — is therefore already proven (no spike needed beyond the existing host-shape tests). The riskiest *value* question — does the hint earn its cost — was the spike, and its answer is "only in Claude team mode, only for visibility, never for completeness."

## Acceptance criteria

Implementation lands option (c): a host-AND-team-conditional todo/task hint line in the `dispatch build` body. The ACs below are entity-level (properties of the finished change), each backed by a gate that lives outside this task body and can fail.

**AC-1 — A Claude team-mode ensign dispatch body carries a todo/task-tool hint.** `dispatch build` with `host="claude"` and a non-empty `team_name` emits, in the dispatch-file body, a line directing the ensign to track the `### Completion checklist` in its task tool. The wording names the tracking behavior (and, in the Claude branch where it is host-correct, may name `TaskCreate`/the task tool).
Verified by: a Go test in `internal/dispatch` modeled on `TestBuildCodexHostPromptShape` — build a `host="claude"`, team-mode dispatch, `readDispatchBody`, assert the hint line is present; a golden capture of the body via the `golden_harness_test.go` `-update` flow freezes it. Reds if the line is dropped or moved out of the body.

**AC-2 — Bare-mode and non-Claude dispatch bodies do NOT carry the hint.** The same builder, invoked with `bare_mode=true` (no team) OR `host` in {`codex`,`pi`}, produces a dispatch body with no todo/task hint line — confining the measured cost to the context with the visibility benefit.
Verified by: the companion `internal/dispatch` test asserting the hint substring is ABSENT in (i) a `host="claude"` bare-mode body and (ii) a `host="codex"` body. Reds if the injection is unconditional. (This is the conditionality gate the A/B's net-cost finding requires.)

**AC-3 — The shared ensign core stays host-neutral; no Claude task-tool name leaks into it.** `skills/ensign/references/ensign-shared-core.md` (and the FO shared core) name none of `TodoWrite`, `TaskCreate`, `TaskUpdate` as an unqualified generic step.
Verified by: a banned-literal entry for the task-tool tokens added to the existing `internal/hostneutrality/ensign_dev_leakage_locks_test.go` sweep (the `devLeakageLiterals` table over `devLeakageCorePaths`), plus the unchanged `skills/integration` portability oracle. Reds if a future edit re-homes a Claude tool name into the universal core. (Option (c) touches no prose, so this AC holds the line against drift rather than gating a new change.)

**AC-4 — The empirical measurement is recorded and reproducible.** The with/without A/B (tool-use rate, report completeness, token+wallclock, cost) and the two without-hint baseline ensign readings are recorded in this body with the per-turn jsonl parse method named.
Verified by: the "Spike results" section above, citing baseline transcripts `agent-a00fe2018e673e1a0.jsonl` / `agent-ad7f1a4e1bfa3c4e5.jsonl` and the A/B harness `/tmp/ensign-todo-spike/ab/{run_arm.sh,parse_stream.py}`; the measurement is re-runnable against a credential. This AC is satisfied at ideation (the spike); it gates that the placement decision rests on data, not assertion.

## Test plan

- **Spike (ideation) — DONE.** The with/without A/B over a live opus `claude -p` drive (benchmark-token credential) + the per-turn jsonl parse. Cost: 6 opus trials, ~11 min wallclock, ~$4.4 total. Result recorded above. This was the riskiest *value* unknown.
- **AC-1 / AC-2 dispatch-build tests** (implementation): two Go tests in `internal/dispatch` (present-in-claude-team-body, absent-in-bare/codex), reusing the `readDispatchBody` + `golden_harness_test.go` substrate. Cheap fixture-level, no live run. This is the conditionality proof the net-cost finding demands.
- **AC-3 host-neutrality** (implementation): one banned-literal addition to the existing `ensign_dev_leakage_locks_test.go` sweep; the portability oracle stays green. Cheap.
- No live workflow test needed for the shipped behavior: the hint is a static body line whose presence/absence is a build-output property, fully gated by the golden/substring tests above. The runtime-behavior question (does the hint change ensign behavior, at what cost) is the ideation spike, already answered — re-running it live in implementation would re-pay the ~$4 with no new claim to prove.

## Notes

Provenance: captain observation 2026-06-04 (this session's ensign-todo investigation — `TodoWrite`=0 across transcripts; not instructed in contract or dispatch file). Option (c) precedent: `dispatch build` already conditionally appends the standing-teammates routing block via `spacedock dispatch show-standing`, so a capability-gated hint line is a known mechanism, not a new one — and the codex/pi body branches in `build.go` already prove the resolved `host` drives body content.

Spike-corrected scope notes:
- The premise's `TodoWrite`=0 holds, but the wm baseline ensign DID externalize via `TaskCreate`/`TaskUpdate` — the real tool is the team task system, not `TodoWrite`. Implementation must word the hint around "task tool" (Claude team) not `TodoWrite`.
- The A/B found a net cost on the captain's three measured axes and no completeness/quality gain; the only benefit is Claude-team Shift+Down visibility, which a transcript metric cannot capture. The host-AND-team-mode gate (AC-2) is therefore load-bearing, not a nicety: an unconditional hint would impose the measured ~+94% externalization cost on bare-mode and codex/pi contexts that get zero visibility benefit.
- Boundary held: this ships ONE hint line in option (c) + its conditionality + neutrality gates + the recorded measurement. It is NOT a todo-workflow redesign and NOT the team `TaskCreate` coordination path — the per-ensign checklist→task hint only.
- Open question for the gate (captain call): the A/B shows no completeness benefit, so if the captain does not value the Claude-team mid-stage visibility enough to accept ~$0.50/+20s per externalizing ensign, the right outcome is to record the decision in the roadmap and NOT ship — the dev policy's "decision-only → roadmap" branch. The FO should surface this trade explicitly at the ideation gate rather than treat "ship option (c)" as foregone.

## Stage Report: ideation

- DONE: Run the riskiest-unknown spike — measure the hint's EMPIRICAL impact (with vs without ensign comparison over tool-invocation count, stage-report checklist-accounting completeness, token + wallclock delta), parsed from transcripts via the per-turn jsonl method; record the measured numbers in the entity body. The placement decision MUST be justified by this data, and if the hint shows no measurable benefit or a net cost, say so.
  Measured: without 0/3 vs with 2/3 externalized; +18% wallclock, +48% cost, +108% turns, identical 10.7/11 completeness; isolated +94% cost premium for externalizing with zero quality/completeness gain. Net cost on the three measured axes recorded in "Spike results"; the only benefit (Claude-team Shift+Down visibility) is un-measurable by transcript and called out as such. Baselines: `agent-a00fe2018e673e1a0.jsonl` (wm, task=15 spontaneous) / `agent-ad7f1a4e1bfa3c4e5.jsonl` (gq, task=0); A/B harness `/tmp/ensign-todo-spike/ab/`.
- DONE: Choose among the three placements grounded in the spike + host-neutrality, and write entity-level ACs each backed by an outside-body gate, keeping the shared ensign core host-neutral.
  Chose (c) dispatch-builder injection, host-AND-team-conditional — the data confines the measured cost to the only context with the (un-measurable) benefit. AC-1 (claude team body carries hint) → `internal/dispatch` body test + golden; AC-2 (bare/codex body omits hint) → companion absence test; AC-3 (no task-tool name in shared core) → banned-literal add to `ensign_dev_leakage_locks_test.go` + portability oracle; AC-4 (measurement recorded) → the spike section. Option (c) touches no prose → shared core stays neutral by construction.

### Summary

The spike overturned the seed's framing: `TodoWrite`=0 is real but the wm baseline ensign already externalized its checklist via the team task tool with no hint, while gq did not — natural model-initiative variance. The controlled opus A/B (3 trials/arm, identical 10-item stage, only the hint differing) showed the hint raises externalization 0/3→2/3 but at +94% cost-per-externalizing-trial with no completeness or correctness gain; the sole benefit is Claude-team mid-stage visibility, invisible to a transcript metric. That asymmetry — real cost, Claude-team-only un-measurable benefit — drives the placement to option (c) gated on `host=claude` AND team mode, so bare-mode/codex/pi never pay the cost. The option (c) mechanism is already shipped (host drives body content today; codex/pi body tests green); implementation is two `dispatch build` body tests + one banned-literal addition, no live run needed. Surfaced for the captain at the gate: if Claude-team visibility isn't worth ~$0.50/+20s per externalizing ensign, the dev "decision-only → roadmap" branch applies and this should not ship — the FO should present that trade rather than assume option (c) ships.

## Spawner-arm spike — pre-created tasks (ideation cycle 2, measured 2026-06-04)

The captain's third arm: instead of the ENSIGN externalizing its checklist (the with-hint arm: +94% cost-per-externalizing-trial, +108% turns), the SPAWNER (team-lead/FO) pre-creates the team tasks (`TaskCreate`, owner=ensign) BEFORE spawning, so the ensign just consumes them (`TaskList`) + `TaskUpdate`s. Hypothesis: the creation cost moves to the cheap FO side, the ensign recovers Shift+Down visibility WITHOUT the externalization premium, and the checklist drops from the prompt.

### Faithful harness (the engineering the team-lead flagged)

A bare `claude -p` is NOT a team member and cannot read a pre-created `TaskList`, so the prompt-only A/B harness cannot capture this arm. The faithful harness: a `claude -p` LAUNCHER (with `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`) that runs the real team flow — `TeamCreate` → 10× `TaskCreate` + 10× `TaskUpdate(owner="worker")` → `Agent(team_name, name="worker")` spawns an in-process member → member reads its `TaskList`, does the work, `TaskUpdate`s each item → teardown. The worker's subagent transcript (under the launcher session's `subagents/`) carries the worker-side metrics. Identical 10-item Go stage, 3 trials. Harness: `/tmp/ensign-todo-spike/spawner/` (`run_spawner.sh`, `launcher_prompt.txt`, `parse_worker.py`, `aggregate.py`). Smoke trial validated the full mechanism end-to-end before the 3-trial run (worker read TaskList, did the work, all tests passed). Mechanism note: every launcher hit the using-claude-team terminal-teardown hang (a `claude -p` launcher cannot self-exit while members linger — the documented #275/#282 behavior); the WORKER finishes early so the measurement is intact, but it burned ~7 min/trial of launcher wallclock on teardown. That hang is itself a finding: this arm's real-FO cost includes the teardown tax unless the FO is a long-lived interactive session.

### Faithfully-comparable metrics (same measurement surface on both arms)

| metric | embedded NO-hint | embedded WITH-hint | SPAWNER pre-created |
|---|---|---|---|
| externalized (worker uses task tool) | 0/3 | 2/3 | **3/3** |
| worker task-ops (TaskList+TaskUpdate) | 0 | 15.7 | **22.0** (20 update + 2 list) |
| worker wallclock (transcript window) | 87 s | 106 s | **137 s** |
| worker prompt size | ~1500 ch (checklist embedded) | ~1500 ch | **570 ch** (checklist in tasks, not prompt) |
| report completeness | 10.7/11 | 10.7/11 | full (10/10 accounted, all trials) |
| work correctness (`go test ./...`) | 3/3 pass | 3/3 pass | 3/3 pass |
| ADDED launcher/FO task-ops | 0 | 0 | **20** (10 TaskCreate + 10 TaskUpdate) |

### The answer — spawner creation does NOT invert the +94%; it relocates and ADDS cost

- **It does not make the worker cheaper.** The worker still does MORE task-ops than the embedded with-hint worker (22 vs 15.7) and runs LONGER (137 s vs 106 s vs 87 s baseline). Consuming + updating pre-created tasks is itself the expense, and pre-creation makes externalization mandatory (3/3 vs 2/3), so the worker pays the update cost every time instead of sometimes. Moving CREATION off the worker did not remove the worker's UPDATE cost.
- **It ADDS a launcher cost on top.** The FO now does 10 `TaskCreate` + 10 `TaskUpdate` (20 ops) it did not do before, PLUS the teardown tax. System-wide task-ops rose to launcher-20 + worker-22 = 42, vs the embedded with-hint worker's 15.7 alone. The cost did not invert; it grew and split across two agents.
- **Cost-in-dollars is NOT faithfully cross-comparable** and I will not fake a number: the embedded arm is a top-level `claude -p` whose authoritative cost is its `result.total_cost_usd` ($0.59 / $0.87); the spawner worker is an in-process SUBAGENT with no `result` cost event, and its transcript's per-turn `output_tokens` are full-turn values while the top-level `-p` stream emits tiny streaming-delta `output_tokens` (max 68 vs 920) — different accounting surfaces. Reconstructing worker cost from summed per-turn cache-read tokens overcounts wildly ($3.80–$5.07, implausible) because per-turn cache_read is cumulative-by-nature. The instrumentation cannot produce a trustworthy cross-arm dollar figure; the comparable axes above (ops, wallclock, prompt size, completeness) carry the verdict instead. (This surface mismatch is the honest limit of the headless harness.)

### What it DID buy

1. **Guaranteed mid-stage visibility (the captain's actual motive).** 3/3 workers drove observable `in_progress`(10) → `completed`(10) transitions on the pre-created tasks from t=0 — full live Shift+Down progress every run, vs the embedded with-hint arm's 2/3. Pre-creation is the only arm that makes visibility reliable.
2. **A 62% smaller worker prompt** (570 vs ~1500 ch): the checklist lives in tasks, not the dispatch prompt — a real prompt-token saving on the worker side.

### Recommendation — does the 3rd arm change the ship-vs-roadmap call

No, it does not flip 3c0b to "ship." The spawner arm delivers the visibility benefit MORE reliably (3/3) and shrinks the worker prompt, but it does so by ADDING net system work (FO 20 ops + teardown tax + a slower, busier worker), not by removing the externalization premium the captain hoped to dodge. It trades "ensign pays to externalize, sometimes" for "FO pays to pre-create AND ensign pays to consume, always." On completeness and correctness it is still a wash (all arms 10.7–11/11, all pass). So the decision is unchanged from cycle 1: the only real benefit across ALL three arms is Claude-team mid-stage visibility, and it is never free. If the captain values that visibility, the spawner arm is the most RELIABLE way to get it (3/3) and the cleanest on prompt size — but it is the heaviest on total system ops and carries the teardown tax, so it should ship ONLY as an FO-side option the captain explicitly opts into, not as a default. If visibility is not worth the added FO ops + teardown cost, the dev "decision-only → roadmap" branch still applies and nothing ships. Either way the recommendation remains: surface the trade at the gate; do not treat shipping as foregone.

## Stage Report: ideation (spawner-arm spike)

- DONE: Measure the SPAWNER-side-tasks arm (FO pre-creates owner-tagged tasks before spawning; ensign consumes + TaskUpdates) faithfully, vs the embedded-checklist arms, on the same 10-item Go stage, 3 trials.
  Built a faithful team harness (a `claude -p` launcher does real TeamCreate → 10× TaskCreate(owner=worker) → Agent(member) → worker reads TaskList + does work + TaskUpdates); a bare `claude -p` can't read a TaskList, so the prompt-only harness was insufficient — engineered the real-team launcher instead. Smoke-validated end-to-end first. Harness + parsers in `/tmp/ensign-todo-spike/spawner/`.
- DONE: Per-turn jsonl parse of worker-side tool-use, turns, wallclock, report completeness, PLUS prompt-tokens-saved and mid-run task-state observability.
  Measured (3/3 trials): worker 22 task-ops (20 update + 2 list) vs embedded with-hint 15.7; worker wall 137s vs 106s/87s; worker prompt 570ch vs ~1500ch (62% smaller); 10 in_progress→10 completed transitions observable mid-run every trial; all 3 trials' work passed `go test`. Faithful cross-arm DOLLAR cost is NOT obtainable (top-level result-cost vs subagent transcript-tokens are different surfaces) — reported as a measurement-limit finding, not faked.
- DONE: Answer whether spawner creation inverts the +94% finding and changes the ship-vs-roadmap call.
  It does NOT invert it: the worker still pays the update cost (more ops, slower) and pre-creation makes it mandatory (3/3), while ADDING 20 FO-side ops + a teardown tax. It buys reliable visibility (3/3 vs 2/3) + a smaller prompt, but by adding net system work, not removing the premium. Recommendation unchanged: only real benefit across all three arms is Claude-team visibility, never free; ship the spawner arm only as an explicit FO-side opt-in if the captain values reliable visibility, else roadmap. Surface the trade at the gate.

### Summary

Built a faithful real-team harness (the team-lead's flagged constraint — a bare `claude -p` is not a team member) and measured the spawner-pre-creates-tasks arm over 3 opus trials against the existing embedded baselines. The arm does not invert the cycle-1 +94% ensign-side cost: the worker still does the full TaskUpdate work (22 ops, more than the with-hint arm's 15.7) and runs longer (137s), pre-creation just makes externalization mandatory (3/3) and ADDS 20 FO-side ops plus the using-claude-team teardown tax. It does deliver two real wins — guaranteed mid-run visibility (10 in_progress→10 completed transitions every trial) and a 62%-smaller worker prompt (checklist in tasks, not prompt). A faithful cross-arm dollar comparison is impossible with headless instrumentation (top-level result-cost vs subagent transcript-tokens) and I reported that limit rather than fabricate a figure. Net: the third arm is the most RELIABLE path to the visibility benefit but the heaviest on total system work, so the ship-vs-roadmap call is unchanged — surface the visibility-vs-cost trade at the gate; ship spawner-side only as an explicit FO opt-in, else roadmap.
