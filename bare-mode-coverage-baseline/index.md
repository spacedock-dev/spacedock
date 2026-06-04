---
id: e3z4yjk7mfna7mktyxw7xbw7
title: "Model bare mode as a coverage variant + establish a live baseline (contract simplification may have broken it)"
status: ideation
source: "captain (2026-06-04) — bare mode (no-team FO dispatch / Degraded-Mode fallback) has NO live coverage and is not a declared variant axis (the spec's `mode` = codified/llm-live evidence path, not team-vs-bare). This session's operating-contract simplification (zd #291 extracted `## Team Creation` — which INCLUDES the bare-mode entry 'if ToolSearch returns no match, enter bare mode' — into the lazily-loaded using-claude-team skill) may have ORPHANED the bare-mode path: bare mode is the case where team setup is unavailable, so the FO may never load the skill that tells it to go bare. Design how bare mode should be modeled AND do an exploratory baseline run to see where we actually are."
score: "0.34"
worktree:
started: 2026-06-04T07:33:00Z
completed:
verdict:
issue:
---

Bare mode is the FO's no-team fallback: sequential blocking `Agent()` (no `team_name`), no `SendMessage` reuse, Degraded-Mode semantics. It activates on team-infra failure and is conceptually the teams-off path. Today it is offline-tested only (`dispatch build` emits the right bare shape) and has ZERO live coverage; worse, it is not a coverage dimension at all. AND the session's contract decomposition may have broken its bootstrapping.

## Problem

1. **No live coverage.** Every live lane runs `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` (team mode). No live test drives an FO through a workflow in bare mode — we have no evidence an FO completes a cycle without a team.
2. **Not modeled.** The scenario spec's variant axes are `runtime {claude,codex} × mode {codified, llm-live}`, where `mode` = evidence path. Team-vs-bare is not a declared dimension, so the gap is invisible to the coverage matrix.
3. **Possible regression from the contract simplification (the riskiest unknown).** zd (#291) moved `## Team Creation` — including the bare-mode entry rule and Degraded Mode — into the lazy `using-claude-team` skill, loaded via `Skill(skill=...)` AT the team-creation step. Bare mode is the case where team setup is unavailable/fails. If the FO only learns "enter bare mode" from a skill it loads while engaging the team path, a genuinely teams-off environment may never surface the fallback. The decomposition was faithfulness-audited for CONTENT, but the bare-mode BOOTSTRAPPING ordering (does the FO reach the bare-mode instruction without a working team?) is exactly the kind of seam a content diff would not catch.

## Proposed approach (ideation)

1. **Design the variant model.** Add team-vs-bare as a first-class coverage dimension — likely a `dispatch-mode {team, bare}` axis distinct from the spec's evidence-path `mode {codified, llm-live}`. Define how it composes with the existing `runtime × mode` matrix and what the cost ledger / coverage meta-tests should require. Decide whether bare is a per-scenario variant or a dedicated lane.
2. **EXPLORATORY BASELINE RUN (the spike — RUN at ideation).** Actually drive an FO in bare mode live and record where we are. Concretely: launch the FO with `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` UNSET (or otherwise force the no-team path) on a small real workflow, and observe: (a) does the FO correctly DETECT and ENTER bare mode given the post-zd contract structure (the bootstrapping question — does it surface the fallback without a working team?); (b) does it complete a dispatch→gate→(merge) cycle sequentially; (c) where exactly does it break or degrade if it does. Record the baseline verdict (works / partial / broken + the precise failure point) in the entity body. This is the riskiest unknown — exercise it FIRST; the rest of the design depends on what the baseline shows.
3. **Branch on the baseline.** If bare mode WORKS: design the standing live bare lane (AC). If bare mode is BROKEN by the simplification: this entity's deliverable becomes the FIX (re-home the bare-mode entry so it is reachable without a team — e.g. keep the bare-mode detection in the always-on skeleton, with only the team-success path deferred to the skill) PLUS the live lane that would have caught it. Either way, end with a live bare-mode check so the fallback stops being invisible.

## Acceptance criteria (seed — firm at ideation after the baseline)

- **AC-1 (seed):** A documented bare-mode coverage model — team-vs-bare as a declared variant dimension in `docs/specs/scenario-testing-principles.md`, composing with runtime × mode, with the coverage/cost shape stated.
- **AC-2 (seed):** A recorded live baseline — an actual no-team FO run with its verdict (works / partial / broken + failure point), proving the exploratory spike ran (a `Verified by: live <ref>` citation under p4's gate).
- **AC-3 (seed):** A standing live bare-mode check (a bare variant of a shared scenario or a dedicated lane) that grades the bare-mode FO on durable outcomes — OR, if the baseline shows a break, the fix that restores the bare-mode bootstrapping plus the check that pins it.

## Out of scope

- Reverting zd's decomposition (if a break is found, re-home the bare-mode ENTRY to the always-on skeleton; do not undo the team-lifecycle extraction wholesale).

## Test plan (seed)

- The baseline is a live exploratory run (p4's `livescenario` primitive is the natural authoring surface). Local live setup: build, export SPACEDOCK_BIN/SPACEDOCK_REPO_ROOT, rotate ~/.claude/benchmark-token (ping FO on 401), force the no-team path. Offline: the variant-model meta-test + any fix's unit coverage.

## Notes

Does NOT block 0.19.5 (bare mode has been live-untested for many releases — not a release regression). Connects to the live-verification line (p4 shipped the citation gate + primitive; n1a is hardening the team-mode live cycle) and the scenario roadmap (`docs/specs/scenario-testing-principles.md`). Sibling context: zd `extract-team-orchestration-skill` (#291) is the decomposition under suspicion; n1a's 1b bare-mode dispatch-path fallback is offline-tested only — consistent with this hole.
