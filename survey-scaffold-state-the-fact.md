---
id: 1p27mtdt6vxt956f5gn00353
title: survey SCAFFOLD section — state the observed fact, drop the recovered-vs-skill-installation classification
status: backlog
source: "captain (2026-06-08) — follow-up to the survey signal-correction (xn). The scaffold reporting carries a recovered-vs-installed taxonomy; just state the observed scaffold fact (what was used + whether it's checked in)."
score: "0.2"
started:
completed:
verdict:
worktree:
issue:
---

Follow-up to xn (survey-signal-correction): simplify the survey's SCAFFOLD reporting.

## Problem

The survey's SCAFFOLD section frames scaffold/skill usage as a recovered-vs-skill-installation classification. That taxonomy is more than the reader needs — just state the observed fact.

## Proposed approach (ideation firms)

Report the observed scaffold fact directly: which skills/scaffolds were used (invocation counts) and whether they are checked in. Desired shape:

```
SCAFFOLD
superpowers was recovered from behavior, not files: 186 skill invocations, but no checked-in .claude/skills/superpowers. Other recovered one-offs: plan-writing, using-git-worktrees, systematic-debug, simplify, debugging.
```

Drop the recovered-vs-installation classification framing; keep the factual statement (used N times, checked in or not). The "recovered from behavior" phrasing in the example IS the fact to state — the change is dropping the taxonomy overhead, not the observation.

## Acceptance criteria (sketch)

- A live survey drive produces a SCAFFOLD section that states the observed fact (skill + invocation count + checked-in presence/absence) without the recovered-vs-installation taxonomy — verified by the live-drive output (the survey's proof bar) plus the SKILL.md / queries change.

## Notes

Survey-skill refinement; same area as xn. Candidate for a survey-followup / 0.19.8 sprint.
