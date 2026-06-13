---
id: 5wvrtfjvjz78fy9xg55p6pjg
title: Survey output redesign — one coherent "value & numbers first" rewrite folding all six feedback seeds
status: backlog
source: "Coalesces the survey-improvement cluster (six user-sourced feedback seeds, 2026-06-08→06-12) into one redesign of the spacedock:survey output. Captain call 2026-06-13: six seeds all rewrite one skill = one coherent change, not six worktrees colliding on SKILL.md."
started:
completed:
verdict:
score:
worktree:
issue:
sprint: 0202-survey-improvements
group: survey
sprint-readiness: ready
---

One coherent rewrite of the `spacedock:survey` output into the captain-locked **"value & numbers first"** structure (canonical mock in `docs/roadmap/0202-survey-improvements/index.md`). The six feedback seeds below are the requirement record; this task is their single execution. All proof is the survey's rendered output over constructed fixtures — never a prose-grep of SKILL.md (the survey discipline). Do NOT regress the validated inference accuracy or the decision-frontier triage (three+ users confirmed these — the keep-signal).

## Requirement bands (each cites its source seed; all render INTO the locked structure)

- **R1 — value-prop legibility (`9h` survey-value-prop-legibility).** Plain "what this gives you" lede + concrete numbers from the user's own data; rename **"mechanical" → "manual"** for repetitive-substantive tracks (reserve "mechanical" for trivial); a **"threads to pull"** actionable section (unresolved/hanging + proactive suggestions), not decision-history narration.
- **R2 — honest lens (`5x` survey-lens-honesty).** Recent-window snapshot framing (surface evolution where the signal exists, else say "recent window, not whole history") + ONE clear partial-lens caveat naming what the corpus can't see.
- **R3 — subagent-dispatch fact (`za` survey-report-subagent-dispatch-fact).** Report the FACT that a session dispatched subagents so an orchestrated repo isn't read as idle (count/presence only; subagent CONTENT stays excluded).
- **R4 — output hygiene (`h5` survey-output-polish).** Collapse the always-empty Codex `workstreams: (unlabeled) — N` breakdown to one honest line; strip the model scratch-reasoning preamble (`I have everything I need. Let me…`) — cross-check findings stay as report content.
- **R5 — knowledge-work archetype + Codex count (`zb` survey-knowledge-work-archetype).** Name a third archetype (knowledge-work loop) beside mechanical/manual-drive and exploration-steering; lead the Codex section with the workdir-attributed count (name-match demoted to caveat) — **verify against the current skill; zb#1 may already be shipped.**
- **R6 — mode-aware framing (`zw` survey-iteration-framing-and-branch-attribution).** Exploration mode leads with iterate/steer (not gates-as-the-whole-model); manual/mechanical mode keeps gate-drive framing. Branch-aware work-by-area as a **detect-branch-and-merge + caveat** (conditional — it did not bite two of the runs).

**Cross-revision rule (from `h5`'s 2nd round):** verify each ask against the CURRENT skill before building — some may be partly shipped (zb#1) or need narrowing (zw#1/#2). **Pending input:** the author-followup reply (with the captain) folds into R5/R6 — ideate-and-revise, do not block.

## Out of scope

- Landing-page CTA + bundling the agent-conversation-index dependency (separate, non-survey tasks — see 0202 cross-cutting).
- Generalizing the decision-frontier / partial-lens beyond the survey (0.21.x decision-abstraction).
- Reproducing any surveyed corpus content (anonymization discipline).

## Acceptance criteria

Ideation firms the full set — one AC per requirement band, each verified by a survey run over a constructed fixture rendering the corrected output, the expected value DERIVED FROM the fixture's session rows (independent source), never a prose-grep of SKILL.md. Sketch:

**AC-1 (sketch, R1) — output opens with plain value + a concrete number, "manual" not "mechanical".** Verified by: fixture render shows the value lede + a count derived from fixture rows; repetitive-substantive track labeled "manual".

**AC-2 (sketch, R2) — recent-window framing + one partial-lens caveat.** Verified by: early-vs-late-divergent fixture renders evolution/snapshot framing; off-corpus-gap fixture renders one caveat naming gaps from fixture rows.

(R3–R6 ACs firmed in ideation, same fixture-render discipline.)

## Test plan

Ideation/implementation firms. One fixture suite over constructed session sets exercising each band's render; a grep over SKILL.md never satisfies a behavioral AC. Riskiest mechanism to exercise first: that the lede's concrete figures (interruptions, hanging threads, dispatch count) are computable from the existing session data.
