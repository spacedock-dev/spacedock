---
title: Sprint state projection feed and progress dashboard
status: ideation
source: "Captain-directed visibility prototype, 2026-07-22."
id: 68r2tcyvf9d0v8yv7sz00qt9
---

## Problem

When a packaged execution deputy drives a sprint, the captain and commanding First Officer cannot see member progress without polling transcripts and reconstructing state by hand. The durable workflow state already knows current stages and Git knows transitions, but there is no legible projection that separates observed history from estimated completion.

## Intended value

Provide a small read-only local web app backed by a machine-readable state projection feed. It should make the durable-decisions sprint visible now and establish the smallest useful feed contract for later clients without creating a second workflow database.

## Prototype brief

Design and implement the smallest projection that can show:

- every selected sprint member with stable identity, title, current stage, and current gate/block state;
- observed stage history with timestamps and source state commits;
- current projected completion time, confidence/basis, and the full sequence of ETA revisions;
- a clear visual distinction between durable observations and projections;
- refresh from authoritative workflow/state Git without mutating workflow state.

Prefer a standard-library Go server and dependency-free HTML/CSS/JavaScript unless the ideation proves a smaller existing surface. Do not add this task to `durable-decisions`; it is an adjacent visibility prototype and must not delay that release train.

## Acceptance criteria

- **AC-1 — Real state:** A local run renders the actual durable-decisions member set and current stages from Spacedock workflow state; no member or stage is hardcoded into the UI.
- **AC-2 — Provenance:** Every observed transition in the feed carries its state-commit identity and timestamp, derived from state history rather than transcript text.
- **AC-3 — Honest projection:** Current ETA and confidence/basis are explicitly marked as projections. Missing evidence produces `unknown`, not invented precision.
- **AC-4 — Revision history:** ETA changes append a revision containing `as_of`, prior/new ETA, basis, confidence, and cause; earlier revisions remain visible and unchanged.
- **AC-5 — Read-only boundary:** Serving or refreshing the dashboard does not mutate entity files, state Git, gates, or worktrees.
- **AC-6 — Falsifiable fixture:** A fixture-backed test exercises multiple members, at least three stage transitions, one unknown ETA, and two ETA revisions; a smoke test proves the browser receives the same projection as the JSON feed.
- **AC-7 — Operable prototype:** One documented local command starts the app and names the feed URL; shutdown leaves both main and state checkouts clean.

## Ideation questions

- What is the smallest feed/event schema that keeps observed state and projection revisions separate?
- Which parts can be derived from state Git, and which ETA inputs must be recorded explicitly?
- Should the first prototype expose a new `spacedock` subcommand, a standalone dev command, or a static generator plus server?
- What files/modules and production/test/docs LOC are intended, with tolerance, before implementation begins?
