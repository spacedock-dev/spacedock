---
id: 3c0bcrn8p60wptvc12fsv5x2
title: Hint ensigns to use todo/task tools when available — spike the placement + measure empirical impact
status: backlog
source: "captain (2026-06-04) — empirically ensigns never convert the helper-built `### Completion checklist` into internal todos (TodoWrite=0 across all ensign transcripts; neither the ensign contract nor the dispatch file hints it). Hint ensigns to use todo/task tools IF the runtime has them — but spike the placement and show the empirical performance impact before committing."
score: "0.27"
started:
completed:
verdict:
worktree:
issue:
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

## Acceptance criteria (preliminary — ideation formalizes)

**AC-1 — The chosen placement carries a runtime-aware todo/task-tool hint** (host-neutral where shared; tool-named where runtime-specific or builder-injected), and for option (c) the hint appears only when the host supports the tool.
Verified by: a presence oracle over the chosen surface, or a `dispatch build` golden/test showing the hint conditional on host capability.

**AC-2 — The empirical impact is measured and recorded** — a with-hint vs without-hint ensign comparison (tool-use, report completeness, token+wallclock), reproducible from the two ensign transcripts.
Verified by: the recorded measurement in the entity body citing the two transcripts + the parse; if the shipped behavior is runtime-observable, a `live <ci-run:|session:>` citation per the live-AC policy.

**AC-3 — Host-neutrality preserved** — no Claude-only tool name leaks into the shared ensign core.
Verified by: the ensign portability / host-neutrality oracle stays green.

## Test plan

- **Spike (ideation):** the with/without empirical comparison — needs live ensign runs (a credential) + the jsonl per-turn parse. This is the riskiest unknown and gates the placement choice.
- **Presence oracle / dispatch-build golden** for the chosen placement.
- **Portability oracle** for host-neutrality.

## Notes

Provenance: captain observation 2026-06-04 (this session's ensign-todo investigation — `TodoWrite`=0 across transcripts; not instructed in contract or dispatch file). Option (c) precedent: `dispatch build` already conditionally appends the standing-teammates routing block via `spacedock dispatch show-standing`, so a capability-gated hint line is a known mechanism, not a new one.
