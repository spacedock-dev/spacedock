---
id: 9hbm8yqsd1c4dwjrdv8kka2y
title: Survey value-prop legibility — lead with plain "helps you do X" + concrete numbers, not abstract jargon
status: backlog
source: "Live partner-meeting user run, 2026-06-12 (~200 sessions / ~2200 logs). Records only spacedock:survey behavioral feedback; the user's corpus (projects, tools, counts) stays in the uncommitted meeting notes, per the survey anonymization discipline."
started: 2026-06-13T04:49:28Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0202-survey-improvements
group: survey
sprint-readiness: defer
---

A live user ran the survey on his own machine and reacted in real time. The inference accuracy landed hard ("this is very accurate, this stuff is hanging" — it caught his plan→worktree→implement pattern, his manual recurring work-streams, and done-vs-hanging work). His strongest critique was that the OUTPUT does not explain its own value to a newcomer.

## Feedback (the asks)

**1. (Keep — accuracy lands.)** The workflow inference and the done/half-done/"hanging" + decision-needed-but-never-closed detection were validated live as accurate. Preserve.

**2. Value-prop first-run clarity (his strongest point).** "I had to do a lot of reading to get to this point. There's not enough explanation of the value — what does this actually do? I'm a new user; I create a SpaceDock workflow — what am I actually getting?" The output leads with abstract jargon ("gated automation for the mechanical tracks, thread bookkeeping for exploration") — "I couldn't tell you what these six words mean." Lead instead with plain "this will help you do X," and anchor the value in concrete numbers from the user's own data (e.g. "~N steering interruptions that could be reduced"), not abstractions.

**3. Terminology: "mechanical" → "manual".** "Mechanical" read to him as trivial (find-and-replace). His work was effortful-but-repetitive — "manual." Reserve "mechanical" for genuinely trivial changes; use "manual" for the repetitive-but-substantive tracks.

**4. Output framing: "threads to pull, not threads you've lost".** Lead with the steady-state — where you are now + a few unresolved/hanging things — and proactively prompt ("have you thought about X, Y, Z?"), rather than dwelling on decision history / the "why" of past choices. (The deeper decision-history-vs-context-pollution principle behind this is a 0.21.x decision-abstraction input — see Notes; this seed's scope is the survey OUTPUT'S framing.)

## Proposed direction (ideation fills in)

- (#2) Restructure the survey's lede/offer to open with concrete, plain-language value tied to the user's own numbers, before any mode/track vocabulary. The abstract track labels, if kept, come after the plain "what you get."
- (#3) Replace "mechanical" with "manual" for repetitive-but-substantive tracks across the rendered output; reserve "mechanical" for trivial edits (couples to zw's mode-aware framing — the mode labels themselves are in scope there).
- (#4) Lead the decision-frontier section with unresolved/hanging "threads to pull" + a steady-state snapshot, de-emphasizing decision-history narration (couples to 5x's recent-window-snapshot framing and the decision-frontier triage that is the 0.21.x wedge).

## Out of scope

- The landing-page CTA ("run the scan on your own sessions" at the top) — that is spacedock-landing, not the survey skill.
- Bundling the agent-conversation index dependency so it auto-installs — a separate productization/dependency task.
- Generalizing "steady-state over decision-history" beyond the survey — 0.21.x decision-abstraction.
- Reproducing the user's surveyed corpus (omitted per discipline).

## Firmed direction

This seed owns the BODY structure of the captain-locked "value & numbers first" mock — the
`WHAT THIS GIVES YOU` lede, the `BY THE NUMBERS` block, the `manual`-not-`mechanical`
terminology rule, and the `threads to pull` actionable section. The subhead above it (title +
`recent-window snapshot · agent logs only · …`) is `5x`'s slice; the mode LABELS the numbers
point at are `zw`'s. These render into ONE shared mock — non-overlapping slices.

Three render-side moves, all from numbers the survey ALREADY computes (no new orchestration,
two-of-three need no new query — see Riskiest mechanism):

- **(#2) A plain-value lede + a numbers block atop the body.** `WHAT THIS GIVES YOU` states the
  value in plain language ("You steer your agents by hand ~N times a month; a SpaceDock
  workflow can run those and stop only where you'd want a say"), and `BY THE NUMBERS` lists the
  three lede figures — hand-steering interruptions, hanging threads, decisions-with-no-follow-up
  — each a real scan number. The mode/track vocabulary ("gated automation … thread
  book-keeping") demotes BELOW this, into `HOW YOU WORK` and the commission offer.
- **(#3) `mechanical` → `manual` for repetitive-but-substantive tracks.** The rendered track
  label and the `HOW YOU WORK` line read "manual" (effortful-but-repetitive); "mechanical" is
  reserved for trivial edits. The classifier's internal `mode` value stays `mechanical` (that is
  `zw`'s concern and the query contract) — this seed changes only the RENDERED word.
- **(#4) `NEEDS YOU` → a "threads to pull" framing.** The OPEN-frontier section leads with the
  steady-state ("where you are now + a few unresolved threads") and a proactive prompt, not
  decision-history narration. The hanging-threads / open-decisions list is the steady-state the
  lede's numbers point at.

## Doc diff (SKILL.md — implementation applies)

**Edit A — Overview, de-jargon the closing-move sentence (line ~13).**

old: `… and the offer is keyed to each track's MODE (automation for mechanical tracks, book-keeping for exploration tracks).`
new: `… and the offer leads with plain value ("this helps you run the repetitive work and stop only at the calls you'd want to make"), keyed to each track's MODE (automation for the manual-but-repetitive tracks, book-keeping for exploration tracks).`

**Edit B — report fence, prepend the value+numbers lede ABOVE the existing body (after the `PROJECT:`/`CODEX`/subhead block, before `SCAFFOLD`).** New block:

```
WHAT THIS GIVES YOU
  You steer your agents by hand ~{interruptions} times over this window. A SpaceDock
  workflow can run the repetitive parts for you and stop only where you'd want a say.

BY THE NUMBERS
  {interruptions}  hand-steering interruptions (decisions + vetoes)
  {open_forks}  hanging threads (started, never closed)
  {no_followup}  decisions you made with no follow-up action
  {if no_followup=0 AND open_forks=0: drop these two lines — say "nothing left hanging"}
```

Every figure is a FILL slot from the step-2 scan; a literal `{slot}` shown is a bug (the
existing fence rule). `{interruptions}` = the existing interruption total (decisions + veto
markers); `{open_forks}` = the post-cross-check OPEN count; `{no_followup}` = the new
no-follow-up count (query below).

**Edit C — rename the rendered `mechanical` track label → `manual` (the WORKSTREAMS `mode`
column render + the `HOW YOU WORK` line + the `INTERRUPTIONS` "mechanical tracks:" line).** The
report renders the word `manual` where `mode=mechanical`; reserve "mechanical" for a trivial
edit. The `mode-classification` QUERY value is unchanged — only the rendered word changes. In
the fence:
- `mechanical tracks: {n} steps across {m} sessions …` → `manual tracks: {n} steps across {m} sessions …`
- `HOW YOU WORK` line: "Mostly manual, repetitive tracks (not trivial — they take real work)."

**Edit D — reframe `NEEDS YOU` as "threads to pull" (fence header + the synthesis lead-in).**

old fence header: `NEEDS YOU   (only if any decision is still OPEN after the repo cross-check)`
new: `THREADS TO PULL   (where you are now + what's still open — only if any decision is still OPEN after the repo cross-check)`

Lead-in prose gains: "Lead with the steady-state — where you are now and the few unresolved
threads — and prompt the next move ('have you thought about X?'), rather than narrating the
history of past choices." (The deeper decision-history principle is 0.21.x — Notes; this is the
survey-local framing only.)

**Edit E — commission-bridge offer (line ~215), demote the jargon behind the plain value.**

old: `Want me to commission a spacedock workflow from this{if both modes: — gated automation for the mechanical tracks, thread book-keeping for the exploration tracks}?`
new: `Want me to commission a spacedock workflow from this — so the repetitive work runs for you and stops only at the calls you'd want to make{if both modes: (gated automation for the manual tracks; thread book-keeping for the exploration tracks)}?`

## New query (implementation adds to references/queries.sql)

`decision-no-followup` — count `done` decisions (answered AskUserQuestion / approved
ExitPlanMode) in repo-scoped Claude sessions that have NO Edit/Write LATER in the same session.
"Later" is the real chronological order: `tool_calls.message_id → messages.ordinal`, NOT the
fragile `tool_calls.id` insertion order. Proven against live data (Riskiest mechanism below).

## Riskiest mechanism — the three lede figures are computable (PROVEN, no spike needed)

Exercised against a real `agentsview sync` of this repo (1740 Claude sessions, 2026-06-13):

- **`{interruptions}`** — PROVEN. The existing "Honest signal accounting" total: AskUserQuestion
  + ExitPlanMode count + veto markers (`[Request interrupted` / `doesn't want to proceed` in
  `messages`). Live: 114 decisions + 19 vetoes = 133. No new query.
- **`{open_forks}`** — PROVEN. The existing `decision-open` OPEN rows (post repo cross-check).
  Live: 10 OPEN. No new query.
- **`{no_followup}`** — PROVEN, but needs the `decision-no-followup` query AND a fixture
  extension. The figure derives from chronological order WITHIN a session. The real
  `tool_calls` table carries `message_id`; the real `messages` table carries `ordinal` (NOT
  NULL) and `timestamp` — so "an Edit/Write after this decision" is a robust ordinal join, NOT
  `tool_calls.id` insertion order. Live: 2 of 102 done decisions had no follow-up.
  **Load-bearing caveat for implementation:** the committed fixture's `tool_calls` has NO
  `message_id` and its `messages` has NO `ordinal` — so the AC-3 fixture MUST be extended to
  carry message ordinals (and `tool_calls.message_id`) for this query to be testable via the
  REAL ordering mechanism. Without that, a fixture test would lean on `tool_calls.id`, which is
  insertion order, not chronology — a false oracle.

Determination: **no spike needed beyond this exercise.** Figures #1/#2 reuse already-proven
queries; figure #3's mechanism (`message_id → ordinal` chronological join) was exercised
end-to-end against live data and confirmed. The only build risk — the fixture lacking ordinals
— is recorded as an explicit implementation requirement, not an open unknown.

## Acceptance criteria

(Each verified by the survey's rendered OUTPUT on a constructed fixture — never a prose-grep of
the skill. The expected figures derive from the FIXTURE's session/tool_call/message rows, an
independent source that can diverge from SKILL.md prose.)

**AC-1 — the survey body opens with a plain-language value lede + a `BY THE NUMBERS` block,
above any mode/track jargon.** Verified by: a survey run over the AC-1 fixture renders, before
any "gated automation / thread book-keeping" vocabulary, a `WHAT THIS GIVES YOU` plain-value
line and a `BY THE NUMBERS` block whose `hand-steering interruptions` figure equals the
fixture's (decisions + veto-marker) count and whose `hanging threads` figure equals the
fixture's OPEN-fork count — both numbers derived from the fixture rows, observed in rendered
output. Non-vacuous: adding one veto-marker message to the fixture increments the rendered
interruptions figure by one.

**AC-2 — the `decision-no-followup` figure is computed by the real chronological join and
rendered in `BY THE NUMBERS`.** Verified by: the `decision-no-followup` query, run against an
ordinal-carrying fixture, returns the count of `done` decisions with no later Edit/Write (by
`message_id → ordinal`), and that figure renders in the `BY THE NUMBERS` block. Non-vacuous:
inserting an Edit at a higher ordinal than a no-follow-up decision's message decrements the
count (proving the chronological join is load-bearing, not a constant).

**AC-3 — repetitive-but-substantive tracks render as "manual", not "mechanical".** Verified by:
a survey run over a fixture with a `mode=mechanical` track renders that track's label and the
`HOW YOU WORK` / `INTERRUPTIONS` lines using the word "manual"; "mechanical" appears only where
the rendered semantics are a trivial edit — observed in rendered output. (The classifier query's
internal `mechanical` value is unchanged; this AC is on the rendered WORD.)

**AC-4 — the open-frontier section renders as "threads to pull" (steady-state-first), not
decision-history narration.** Verified by: a survey run over a fixture with ≥1 OPEN fork renders
the frontier under a "threads to pull"-framed header with the OPEN forks as the steady-state
list — observed in rendered output. The OPEN-fork count rendered equals the fixture's OPEN count
(independent source).

## Test plan

Two layers, mirroring the sibling seeds' proof shape:

1. **Query-smoke (AC-1 numbers, AC-2) — `skills/integration/survey_queries_test.go`.** Add the
   `decision-no-followup` query to the smoke's run-each-labeled-query list and assert its count
   against the fixture rows. This REQUIRES extending `testdata/survey/fixture-sessions.sql`:
   add `message_id` to the `tool_calls` rows and `ordinal` to the `messages` rows, and seed a
   `done` decision with no later Edit (counts) plus a `done` decision with a later Edit (does
   not) — so the expected count derives from the fixture and the non-vacuous mutation (insert a
   higher-ordinal Edit) flips it. The interruption total (AC-1) is already derivable from the
   existing decision + veto-marker fixture rows. Cost: ~1 fixture-extension + 1 smoke sub-test,
   low (the harness exists).
2. **Live-drive render (AC-1 layout, AC-3, AC-4) — the render is SKILL.md prose, so the proof is
   a survey run whose rendered output is observed.** A constructed fixture (or the existing
   live-drive harness used by the sibling ACs) drives the full report; assert (a) `WHAT THIS
   GIVES YOU` + `BY THE NUMBERS` precede any "gated automation"/"thread book-keeping" text, (b)
   the `mode=mechanical` track renders the word "manual", (c) the frontier header is the
   "threads to pull" framing. Per the survey discipline, a grep over SKILL.md NEVER satisfies
   these — the render must be exercised and observed. Cost: folds into the cluster's shared
   live-drive (the six seeds render into one mock), low marginal.

## Notes

- **Two-author keep-signal grows:** this is now a third independent run validating the core inference accuracy (with `zby`#3 / `zwg`#0 / `5xw`). Do not regress it.
- **Corroborates `zrc` (0.20.1):** the user independently asked for a non-sandbox auto-accept (shift-tab auto-mode) default — exactly `non-sandboxed-launch-auto-mode`. Cross-sprint validation.
- **Feeds 0.21.x (decision abstraction):** his most substantive critique — "historic decisions about things you did NOT choose pollute the context window; I want the steady-state of what's there now, not how I got there; short-term vs long-term memory" — is a strategic input to the decision-event model / ready-room. The "threads to pull, not threads you've lost" reframe is the survey-local expression; the principle generalizes.
- **Cross-cutting (not this seed):** landing-page CTA (run-the-scan-first) → spacedock-landing; bundle the agent-conversation-index dependency so it auto-installs → a productization/dependency task. Surfaced to the captain for filing.
- **Follow-up owed by CL:** send the feedback-giver his portion of the transcript (redacted) — he offered to refine it.

## Stage Report: ideation

- DONE: Firm the survey output against the banked "value & numbers first" structure (WHAT THIS GIVES YOU lede, BY THE NUMBERS figures, "manual"-not-"mechanical" rule, "threads to pull" section)
  Firmed this seed's slice (body lede + numbers block + terminology + frontier framing) within the captain-locked structure; subhead is `5x`'s, mode labels are `zw`'s — coherent, non-overlapping. Doc diff Edits A–E in the body.
- DONE: Riskiest-mechanism check FIRST — the concrete lede figures are computable from existing session data
  Exercised against a live `agentsview sync` of this repo (1740 sessions): interruptions 114+19=133, OPEN forks 10, no-follow-up 2/102. Figures #1/#2 reuse proven queries; #3 needs the real `tool_calls.message_id → messages.ordinal` chronological join (proven) + a fixture extension (recorded). "No spike needed beyond this exercise."
- DONE: Propose the doc diff + a test plan whose AC proof is a fixture render of plain value + a number DERIVED FROM fixture rows + "manual" labeling — never a prose-grep
  Doc diff = SKILL.md Edits A–E + new `decision-no-followup` query. ACs 1–4 each verified by rendered output / query count over a fixture, expected values from fixture rows. Test plan: query-smoke (extend fixture-sessions.sql with message ordinals) + live-drive render.

### Summary

Firmed the value-prop-legibility seed into four ACs and a five-edit SKILL.md doc diff, all rendering into the captain-locked "value & numbers first" mock — no structure redesign. The load-bearing finding: all three lede figures are computable, but figure #3 (decisions-with-no-follow-up) requires the real `message_id → ordinal` chronological join, and the committed fixture lacks both columns — so the AC-2 fixture MUST be extended with message ordinals, recorded as an explicit implementation requirement rather than a silent assumption. Mechanism proven live (no spike needed); two of three figures reuse already-proven queries.
