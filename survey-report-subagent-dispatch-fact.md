---
id: zajryfzgmnbzb5vjw5hven0z
title: Survey reports the FACT of subagent dispatch (so orchestrated repos don't look idle)
status: ideation
source: "captain (2026-06-08) — follow-up surfaced during 47rx's F spike. The survey scopes to the Claude session set and EXCLUDES dispatched-subagent sessions, so in a spacedock/agent-orchestrated repo most work is invisible (e.g. 7468/7680 internal edits here landed in excluded subagent sessions) → the survey under-reports and the repo looks idle. Captain scoping: we do NOT need subagent CONTENT for now — just the FACT that the session dispatched subagents."
started: 2026-06-13T04:49:28Z
completed:
verdict:
score:
worktree:
issue:
group: survey
sprint-readiness: ready
sprint: 0202-survey-improvements
---

When a session dispatched subagents, the survey should say so — surface the fact of orchestration so a spacedock/agent-orchestrated repo isn't reported as having done little.

## Problem

The survey's body (work-by-area, workstreams, interruptions) is built from the parent Claude sessions and EXCLUDES dispatched-subagent sessions. In a repo where the real work happens inside dispatched subagents (a spacedock-orchestrated repo), that work is invisible — the survey under-reports and the project reads as idle. 47rx's F spike measured this directly here: 7468/7680 internal edits were in excluded subagent sessions.

The fix is NOT to fold in subagent CONTENT (their edits/work-by-area) — that's deliberately out of scope. It is to report the FACT: this session dispatched subagents (orchestration happened), so the reader understands the parent-session view is only part of the picture.

## Proposed approach

Ideation firms. Sketch: from the session data already read, detect that a parent session dispatched subagents and surface a fact line — e.g. a count of dispatched subagents (or sessions-that-orchestrated), rendered in the survey so an orchestrated repo's body carries "N subagents dispatched" rather than appearing empty. Content of those subagent sessions stays excluded.

## Out of scope

- Subagent session CONTENT (edits, work-by-area, decisions) — explicitly deferred; this surfaces only the dispatch fact.
- Changing the Claude-session scoping of the body itself.

## Acceptance criteria

Ideation/implementation fills in. Sketch:

- The rendered survey, run over a repo whose sessions dispatched subagents, surfaces the dispatch fact (a count / presence of orchestration), distinct from the parent-session body. Verified by: a live drive / query-smoke over a fixture with dispatched-subagent sessions — the rendered output names the dispatch fact; the expected count comes from the fixture session rows (independent source), not skill prose.

## Test plan

Ideation/implementation fills in. Query-smoke over a fixture (sessions with dispatched-subagent markers) + a live-drive of the rendered survey. Per the survey discipline, a grep over SKILL.md never satisfies the behavioral AC.
