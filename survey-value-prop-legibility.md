---
id: 9hbm8yqsd1c4dwjrdv8kka2y
title: Survey value-prop legibility — lead with plain "helps you do X" + concrete numbers, not abstract jargon
status: backlog
source: "Live partner-meeting user run, 2026-06-12 (~200 sessions / ~2200 logs). Records only spacedock:survey behavioral feedback; the user's corpus (projects, tools, counts) stays in the uncommitted meeting notes, per the survey anonymization discipline."
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

## Acceptance criteria

(Ideation firms. Each verified by the survey's rendered OUTPUT on a constructed fixture — never a prose-grep of the skill.)

**AC-1 (sketch) — the survey output opens with plain-language value + a concrete number, not abstract track jargon.** Verified by: a survey run over a fixture renders a lede that states a concrete value ("reduce ~N interruptions" / "X hanging threads") in plain language before any "gated automation / thread bookkeeping" vocabulary — observed in rendered output, the number derived from the fixture's session rows (independent source).

**AC-2 (sketch) — repetitive-but-substantive tracks are labeled "manual", not "mechanical".** Verified by: a survey run over a fixture with a repetitive non-trivial track renders it as "manual" (or equivalent), with "mechanical" reserved for trivial edits — observed in rendered output.

## Test plan

(Ideation/implementation firms.) Fixture-driven survey renders: a fixture whose session rows yield a countable interruption/hanging-thread figure (AC-1) and one with a repetitive substantive track (AC-2). Per the survey discipline, a grep over SKILL.md never satisfies the behavioral AC.

## Notes

- **Two-author keep-signal grows:** this is now a third independent run validating the core inference accuracy (with `zby`#3 / `zwg`#0 / `5xw`). Do not regress it.
- **Corroborates `zrc` (0.20.1):** the user independently asked for a non-sandbox auto-accept (shift-tab auto-mode) default — exactly `non-sandboxed-launch-auto-mode`. Cross-sprint validation.
- **Feeds 0.21.x (decision abstraction):** his most substantive critique — "historic decisions about things you did NOT choose pollute the context window; I want the steady-state of what's there now, not how I got there; short-term vs long-term memory" — is a strategic input to the decision-event model / ready-room. The "threads to pull, not threads you've lost" reframe is the survey-local expression; the principle generalizes.
- **Cross-cutting (not this seed):** landing-page CTA (run-the-scan-first) → spacedock-landing; bundle the agent-conversation-index dependency so it auto-installs → a productization/dependency task. Surfaced to the captain for filing.
- **Follow-up owed by CL:** send the feedback-giver his portion of the transcript (redacted) — he offered to refine it.
