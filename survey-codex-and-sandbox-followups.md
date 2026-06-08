---
id: 47rx3x8a809wx35vx6rbqqhv
title: Survey Codex body-surfacing + sandbox honesty (0.19.9 follow-ups from vh)
status: ideation
source: "captain (2026-06-08) — deferred from vh (survey-skill-correctness-pass, 0.19.8). vh's live drive proved the codex-presence COUNT fires (61), but Codex sessions are not surfaced in the survey body, and a sandbox-denied ~/.codex yields a silent confident 0. Captain: file these as a bundled 0.19.9 candidate."
score: "0.25"
started: 2026-06-08T19:37:02Z
completed:
verdict:
worktree:
issue:
group: survey
sprint-readiness: ready
sprint: 0199-pre-flip-mechanics
---

Three survey follow-ups deferred from vh (0.19.8) to 0.19.9. They share the survey skill (`skills/survey/SKILL.md` + `references/queries.sql`) and the agentsview Codex/sandbox theme. Ideation should confirm whether they ship as one task or split — D1 (body-surfacing) is the largest and has an upstream dependency; D2/D3 are small honesty fixes. **0.19.9 candidate per captain (2026-06-08); roadmap/sprint assignment is the captain's.**

## Problem

vh shipped the corrected git-root-basename model and a flagged Codex presence count, but left three gaps the live drive surfaced:

- **Codex sessions are absent from the survey BODY.** The report shows a flagged Codex presence COUNT only (matched by `project` name, caveated "may include a same-named sibling"); the decisions / SCAFFOLD / work-by-area / workstreams body is Claude-only. Root reason: agentsview does not record Codex cwd (every Codex session lands `cwd=''`), so a Codex session cannot be cwd-prefix-scoped — and the cwd-prefix is the only thing that excludes a **same-basename sibling repo** (`project` is git-root-basename, which collides: `spacedock`→3 roots, `workspace`→6). Folding Codex into the body naively would corrupt the counts with sibling sessions. This exclusion is correct today (vh AC-2 asserts codex-not-folded-into-scope), but it means a Codex-heavy repo's real work (e.g. this repo's codex-adapter sessions) never reaches the body.
- **Silent-0 under a denied Codex source.** When `~/.codex` is unreadable (Seatbelt sandbox) or absent, `agentsview sync` silently succeeds with 0 Codex sessions, so the survey prints a *confident* "0 Codex" indistinguishable from genuinely-zero (the exact confusion the captain hit before allowing `~/.codex` access).
- **No source-health signal.** The survey cannot detect that a source was denied/disabled vs. genuinely empty.

## Proposed approach (sketch — ideation firms)

- **D1 — surface Codex sessions in the body.** Needs ONE of: (a) upstream agentsview persisting Codex cwd (then Codex joins the cwd-scope cleanly — an upstream `kenn-io/agentsview` change, a dependency not a deliverable here), or (b) an attribution heuristic distinguishing this repo's Codex from a same-basename sibling's (e.g. session-content fingerprinting) so Codex can be safely folded into decisions/scaffold/work-by-area. Largest item; gated on (a) or (b).
- **D2 — silent-0 caveat (small, ~2 lines).** When codex-presence = 0, render a one-line honest caveat ("0 matched — a sandboxed run where ~/.codex is unreadable also shows 0; agentsview ingests Codex from ~/.codex/sessions") instead of a bare confident 0. In-theme with vh's AC-5 and §2 honest-accounting. The rendered line is live-drive-confirmed.
- **D3 — hard readability detector.** Distinguish genuinely-0 from a denied/disabled Codex source. Not cheaply achievable today (the agentsview subprocess is itself sandboxed; an agent-side `[ -r ~/.codex ]` faces the same denial that forced vh's AC-5 `command -v`→`agentsview --version` swap; agentsview exposes no per-source count / source-health signal). Likely needs an upstream source-health signal or a probe through the binary.

## Out of scope

- Changing the cwd-scoped Claude body mechanism (it is correct — D1 adds Codex without altering it).
- The upstream agentsview Codex-cwd change itself (a dependency for D1(a), not this task's deliverable unless explicitly adopted).

## Notes

Provenance: vh feedback cycle 1 + the captain's 2026-06-08 live-drive findings. The implementation ensign's recommended shapes (happy-path 61, the no-all-agents-flag fact, the D2 caveat wording, the D3 detector difficulty) are recorded in vh's archived cycle-1 stage report (`_archive/survey-skill-correctness-pass.md`). Per the survey discipline, a grep over SKILL.md never satisfies a behavioral AC — D2/D1 bottom out on a live drive.
