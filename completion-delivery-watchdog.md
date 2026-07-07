---
title: "Completion-delivery watchdog — a dispatched worker's finished state is detected within a bounded window even when the inbox stalls"
status: backlog
source: "0250 Commander session 2026-07-07, quantified worst instance of the known delivery-lag family: the zr implementation ensign SENT its Done to team-lead at 09:41-42 UTC (transcript + stage-report state commit 15c2cc5e at 17:41:19+08), but delivery into the lead's context stalled until ~12:29 UTC (the lead's advance state commit 29033ad5 at 20:29:57+08) — a ~2h48m stall with ZERO intermediate wake-ups, so the lead's existing check-durable-state-on-wake posture had no wake to act on. The captain observed 'ensigns stopped reporting' and manually prodded a validator that had in fact already sent (prod arrived 3 minutes after the unprompted send). Ensign-side contract verified intact across all session transcripts; the defect is harness inbox delivery."
started:
completed:
verdict:
score: 0.5
worktree:
issue:
id: dbmpqm2sg9n869457z6sqjk8
---

The FO's completion detection has exactly one channel — inbox delivery — and no deadline on it: a stalled delivery parks the whole pipeline silently for hours with the worker finished and the lead unaware. Direction: after each dispatch, the FO arms a bounded background watchdog (the external-wait watcher pattern, cf. 3j — e.g. a background `sleep <envelope>; probe` whose exit wakes the lead) that checks DURABLE state (the entity's stage report / state-checkout commits — the dispatch contract already names the stage report as "the gate in every case") and treats verified-finished work as complete without the message; delivery-stall instances get logged for the harness-side fix. Scope options for ideation: contract prose in the dispatch core's Awaiting-Completion sibling (runtime-neutral, per the 3j convention) and/or a shipped `spacedock dispatch await --name X --timeout N` verb that polls the durable surfaces (code gate > prose). Coordinate with 3j (external-wait watcher) and eh (helper completion contract) — three seeds, one family: signals need senders, contracts, AND deadlines. Acceptance sketch: value — a seeded delivery-stall drive reaches gate handling within the watchdog envelope instead of parking indefinitely (baseline: the 2026-07-07 2h48m stall); mechanism — the watchdog clause/verb ships with a live-observed arming.
