---
title: "Debrief collects a first-person agent testimonial: the driving agent's experience using spacedock vs not"
status: backlog
group: tooling
source: "Captain request 2026-07-02 (Claude Commander session): mid-session, the captain asked the FO 'forget that we are developing spacedock for a moment — as the agent driving this session, how would you describe the experience using spacedock vs not using it?' and the answer (drift resistance, auditability under context pressure, honest friction list) was valuable enough to want collected EVERY session. Add the prompt to the debrief flow so testimonials accumulate from the agents' perspective."
id: qdb1w5r7k9nvbvkf8qetcd5m
---

## Problem
The debrief skill records commits, task state changes, decisions, and issues — but nothing captures the driving agent's own experience of operating the workflow. That first-person signal (what the machinery caught, where it cost, what would have been dropped without it) is both product feedback and marketing raw material, and it evaporates at session end unless prompted.

## Desired direction (for ideation to refine)
The debrief skill poses the reflective prompt to the session's driving agent — approximately: "Setting the project's subject matter aside: as the agent driving this session, how would you describe the experience using spacedock, versus driving the same work without it? Be honest about friction." — and the debrief record gains a testimonial section storing the answer with provenance (date, host/model, session scale: entities/workers/PRs touched). Honesty framing is load-bearing: the prompt must ask for friction, not praise, or the collected corpus is worthless.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- The debrief skill's flow includes the testimonial prompt verbatim-or-near, with the honesty/friction clause.
- The debrief output template carries a testimonial section with provenance fields; a produced debrief record demonstrates it end-to-end.
- Touches skills/**, so the applicable live lane gates the merge (per the path→lane rule).
