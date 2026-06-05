---
id: 19fhrfae24d221wzgqm4zarn
title: Bring workflow-discovery into spacedock — orient on a project's implicit agent workflow as the front-door to commission
status: backlog
source: "captain (2026-06-05) — 'check the prototype orient skill in ~/.claude, this is the workflow discovery thing we should bring in.' Prototype at ~/.claude/skills/orient/SKILL.md."
score: "0.33"
started:
completed:
verdict:
worktree:
issue:
---

There is a working prototype `orient` skill (`~/.claude/skills/orient/SKILL.md`) that reconstructs a project's IMPLICIT workflow from its AI-agent session history: the inferred loop, the workstreams, the recent decisions, and — load-bearing — the OPEN decisions (abandoned/unanswered forks) plus an interruption count (how often the human had to step in). It reads agentsview's multi-agent session DB (`~/.agentsview/sessions.db`) with plain `sqlite3` queries. Its closing pitch is exactly the spacedock thesis: "spacedock turns these [interruptions] into gates so the agent advances on its own between your calls."

That makes orient the natural FRONT-DOOR to commission: discover what a brownfield project's agents have implicitly been doing → surface the workstreams + the OPEN decision frontier → propose commissioning it as a real spacedock workflow with gates.

## Problem

Spacedock today starts from a commissioned workflow; it has no way to DISCOVER a workflow from a project that already has ad-hoc multi-agent history. The bridge from "agents have been working here unsupervised" to "a spacedock workflow with gates" is missing — orient is the prototype of that bridge.

## Direction (for ideation)

- Read the prototype `~/.claude/skills/orient/SKILL.md` (the agentsview sqlite scan + the synthesis format) as the reference; decide what comes into spacedock and in what shape — a spacedock-owned `orient`/discover skill, a `spacedock` command, or a commission pre-step.
- The discovery → commission bridge: the OPEN decisions and the interruption stats are the raw material for proposing gates; the workstreams map to candidate stages/entities. Design how discovery output feeds the commission skill.
- Dependency: agentsview (the session DB this reads). Decide how spacedock handles its absence (the prototype asks consent + installs on a yes). Do NOT make discovery hard-require agentsview without a graceful path.
- Spike-first risk: confirm the agentsview DB shape + the sqlite queries actually reconstruct a useful workflow on a real project's history (the prototype already does this — exercise it on this repo's own history as the spike) before committing to the integration shape.

## Out of scope

Building agentsview itself (it's an external tool). Re-deriving the scan — reuse the prototype's queries.

## Acceptance criteria

**AC-1 — spacedock can discover a project's implicit workflow + OPEN-decision frontier from agent-session history, and bridge it toward commission.**
Verified by: running the discovery surface on a real project with agent history produces the workstreams + OPEN decisions + interruption stats (observable output), and the discovery → commission hand-off is exercised (a discovered workflow feeds the commission flow) — a live/command exercise on real history, not prose review. (Ideation fixes the exact surface + assertion + the agentsview-absent graceful path.)

## Test plan

Spike on this repo's own agentsview history first (the prototype's queries). Then the integration's exercise: discovery output on a history-bearing project + the commission hand-off. Ideation sizes whether a live-history fixture or a real run is needed.
