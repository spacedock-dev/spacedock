---
id: h8qkmx04azabr6vbkgke4d6c
title: Journey-delta PR comment is a terse, unreadable wall of numbers — no baseline identity, no highlighted signal
status: backlog
source: Captain review of the real posted comment on PR #470 (2026-07-03, run 28645201074) — first genuine live journey-delta comment produced by boot-metrics-record-and-pr-delta. Captain's verbatim reaction: "it's very terse and very hard to read. it doesn't tell me what it is delta against, and it should highlight the most important thing."
started:
completed:
verdict:
score:
worktree:
issue:
---

The journey-delta PR comment (`internal/release/journeydelta.go`'s `RenderJourneyDeltaComment`, shipped in boot-metrics-record-and-pr-delta / PR #470) works mechanically — real deltas, real numbers — but the actual posted output is hard for a human reviewer to use. Captured verbatim from PR #470's real comment (run 28645201074, 14 scenario/model rows):

```
<!-- spacedock:journey-delta -->
### Journey cost delta

| Scenario | Runtime | Model | Turns Δ | Cache Read Δ | Cache Creation Δ | Tokens Δ (total) | Cost Δ (USD) | Duration Δ (ms) | Baseline |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| feedback-3-cycle-escalation | claude | claude-opus-4-8 | -2 | -86973 | -187 | -85734 | -0.0289 | +2823 | latest published |
| ... (14 rows total, every "Baseline" cell reads "latest published") |
```

## Problem

{Ideation to flesh out. Seed framing from the captain's reaction — two concrete gaps:
1. **No baseline identity.** Every row's "Baseline" column reads the literal string "latest published" — no release version, no tag, no date, no link to what's actually being compared against. A reviewer cannot tell if the baseline is last week's release or six months old. `journeydelta.go`'s existing `BaselineRunURL` field is per-observation (points at the source CI run, not the release), so even that doesn't answer "which release." The release VERSION the baseline ledger came from is available at render time (it's the file just downloaded) but is not threaded into `RenderJourneyDeltaComment` at all today.
2. **No highlighted signal.** All 14 rows render with equal visual weight in one flat table, sorted alphabetically by scenario name. A reviewer scanning this has to manually compare 9 numeric columns across 14 rows to notice that `rejection-flow` (sonnet) moved by +26 turns and +$1.15 — the actual standout in this run — while `filing`/`merge-hook-guardrail` moved by routine single-digit amounts. Nothing surfaces the biggest mover(s) or flags a scenario that crossed some threshold.
}

## Proposed approach

{Explicitly NOT decided here — ideation's job. Candidate directions to evaluate, not commit to:
- Render the release version/tag (and ideally a link to the release page) once, near the top of the comment, instead of repeating an unhelpful string per row.
- A short "highlights" line or sub-section above the full table calling out the largest delta(s) by some measure (cost? tokens? a normalized magnitude?) — needs a definition of "most important," which is a real design question, not just formatting.
- Whether the full per-row table stays as supporting detail below the highlight, or gets collapsed/reorganized (e.g. sorted by magnitude instead of alphabetically, grouped by direction of change).
Ideation should look at the actual rendering code (`RenderJourneyDeltaComment` in internal/release/journeydelta.go) and the real posted comment before proposing a shape.}

## Out of scope

{Ideation to determine — likely candidates: changing what's measured/computed (this is a rendering/readability problem, not a data problem); the underlying delta-computation logic in `ComputeJourneyDeltas` should not need to change.}

## Acceptance criteria

{Ideation to flesh out with concrete Verified-by evidence. Seed ACs: (1) the comment names the specific baseline release (version/tag), not just the literal string "latest published"; (2) the comment surfaces the largest/most notable delta(s) distinctly from the full row-by-row table, with the "most notable" criterion explicitly defined and testable.}

## Test plan

{Ideation to fill in — likely a fixture test on the rendered comment body's structure/content, mirroring the existing `TestRenderJourneyDeltaCommentIncludesExactDeltasAndMarker` / `TestRenderJourneyDeltaCommentShowsTokenClassBreakdown` pattern in internal/release/journeydelta_test.go.}
