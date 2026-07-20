---
id: ve7f64jv0gcywa18p3m9ttv5
title: Feedback rejection routes to design-reset when findings indict the mechanism or the diff grows across cycles
status: backlog
source: "0260 shaping — agent-derail forensics audit, 2026-07-19."
score: "0.85"
sprint: 0260-proportionality
group: reframe
---

A rejection whose findings indict the mechanism's architecture — or a repair cycle that grows the diff — must halt at a design-reset decision (park / re-scope / escalate) instead of dispatching another repair. Cycle-counting alone provably fails: every observed runaway loop was contract-legal round by round. Pairs with `bw` (binary-owned cycle count) as the code-gate substrate: dispatch build refuses a repair dispatch at cycle ≥2 with diff growth unless a recorded reframe decision exists.

## Problem

{Ideation fills in. Evidence: e6j 2-defect fix → 10 roborev cycles, 26 files/+3,373, PR closed (codex:019f5fe6); dp one-paragraph fix → 4-cycle ladder, discarded (6d175b2f, ab6c437e); 7h harness repaired twice before park (codex:019f5160:499); 419-line synthetic proof under "AC-1/AC-2 remain unproven" pressure (bef9653f:251-509). The 7h postmortem's "feedback-loop exception" drafted this rule in-session and it never left the session.}
