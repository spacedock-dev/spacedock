---
id: dyxqywnwb4c3zwb3pka1p6s0
title: Survey detects spacedock as the incumbent scaffold when genuinely in use (not just its own self-invocation)
status: backlog
source: "Captain observation 2026-06-13, reviewing the 0202 5wv survey-output redesign: the survey ran on the spacedock repo itself, yet the SCAFFOLD section named only `superpowers: 7` and never spacedock — the dominant scaffold here. Root cause: the `scaffold-usage` query deliberately excludes `WHERE family <> 'spacedock'` to stop the survey's own invocation making every repo read as a spacedock user. Pre-existing (#319/#317.1); surfaced by, not caused by, 5wv."
started:
completed:
verdict:
score:
worktree:
issue:
group: survey
---

The `spacedock:survey` SCAFFOLD section does not name spacedock even when spacedock is the genuine incumbent scaffold driving the repo, because the `scaffold-usage` query deliberately excludes the spacedock family.

## Problem

`references/queries.sql` `scaffold-usage` ends with `WHERE family <> 'spacedock'`, and SKILL.md step 3 ("Recognize an incumbent scaffold") documents the rationale: counting spacedock-family skill calls would let the survey's OWN invocation (`spacedock:survey`, plus dispatched ensign/FO calls) dominate, making every surveyed repo falsely read as "uses spacedock" — a false positive.

That blunt exclusion produces a FALSE NEGATIVE on a genuine spacedock repo. Surveyed on spacedock itself (2026-06-13, 68-session corpus), the SCAFFOLD section named `superpowers: 7` and omitted spacedock — the actual dominant scaffold (the whole first-officer/ensign/dispatch workflow). The survey hides the one scaffold that matters most there.

## Proposed approach (ideation to firm)

Distinguish a GENUINE spacedock incumbent from the survey's own self-invocation, so the SCAFFOLD section names spacedock when it really drives the repo while keeping the exclusion where only the survey self-call appears. Candidate signals to reconcile (mirror step 3's existing "file probe + behavioral tally" join):

- **File probe** — spacedock workflow presence on disk: `.spacedock-state/` (split-root state), a workflow README with spacedock frontmatter (`commissioned-by:` / `stages:`), `_mods/`, a `docs/*/README.md` workflow definition.
- **Behavioral** — non-survey spacedock invocations in MAIN (non-subagent) sessions: `spacedock:first-officer` / `spacedock:ensign` distinct from the `spacedock:survey` self-call that triggered the run. (Most ensign calls live in EXCLUDED dispatched-subagent sessions, so the FO-in-main signal plus the existing `dispatch-fact` count are the visible behavioral evidence.)
- **Fallback** — keep the current exclusion for a repo with only a `spacedock:survey` self-invocation and no spacedock workflow on disk.

## Out of scope

The survey report STRUCTURE (0202 5wv shipped that). This is scaffold-DETECTION accuracy only.

## Acceptance criteria (sketch — ideation firms)

**AC-1 (sketch) — a genuine spacedock repo names spacedock in SCAFFOLD; a survey-self-only repo does not.**
Verified by: a fixture/probe test where a fixture with a spacedock workflow on disk (`.spacedock-state/` + a spacedock-frontmatter README) renders spacedock in SCAFFOLD, while a fixture with only a `spacedock:survey` self-invocation and no workflow on disk does not — the expected outcome deriving from the fixture's on-disk state, never survey prose.

## Notes

Not a 5wv regression — 5wv restructured the report and added the two new queries; the `WHERE family <> 'spacedock'` exclusion predates it. Captain confirmed the "manual" track classification on the same run is accurate (sprints are driven manually), so the keep-signal is intact — this seed is scaffold detection only.
