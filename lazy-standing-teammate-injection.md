---
title: "Standing teammates inject lazily at first routed use, not eagerly at first dispatch"
status: backlog
source: "Principle-derived contract review, 0250 Commander session 2026-07-07 (captain-requested filing). Live exhibit: this session's comm-officer — spawned at first dispatch per fo-dispatch-core's standing-injection rule, delivered its online greeting, then received zero polish requests across an entire multi-member sprint drive while costing its spawn prompt, residency, and multiple idle-ping wake-ups of the FO's own context. Certain every-session token waste with near-zero efficacy offset."
started:
completed:
verdict:
score: 0.4
worktree:
issue:
id: s8g111hd7e4vg4tc5k7c1qb0
---

fo-dispatch-core mandates standing-teammate injection before the first worker dispatch, unconditionally — machinery-first where the contract's own smallest-sufficient-mechanism principle (zm) says need-first. Direction: injection defers to the teammate's first routed use (e.g. the first polish request routes through a just-in-time spawn), or the mod declares `injection: eager|lazy` with lazy the default; commissioned semantics and zm's scope clause are unaffected (the teammate is still workflow-declared — only its spawn TIMING changes). Acceptance sketch: value — a session that never routes to the standing teammate spawns zero standing agents (baseline: one wasted live agent + idle pings per session, observed 2026-07-06/07); a session that does route gets identical service one spawn-latency later; mechanism — the dispatch-core clause + spawn-standing-all surface updated together. Coordinate wording with the merged zm blockquote to avoid duplicated resident bytes.
