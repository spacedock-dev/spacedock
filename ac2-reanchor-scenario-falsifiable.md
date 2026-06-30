---
title: Strengthen the #441 AC-2 re-anchor live scenario — it cannot fail on the regression it polices
status: backlog
score: 0.4
source: "pre-tag spot-audit of #441/#442, 2026-06-30. Real major, verified non-ship-blocker (the gate-on-end-value prose shipped fine; only its behavioral PROOF is vacuous). Test-strength follow-up."
id: w07c8xz91q2yze26b7ft1vka
---

`internal/livescenario/ac2_reanchor.go` (#441's sole behavioral proof of the AC-cross-check re-anchor rule) is non-falsifiable for the change it ships: the runbook (line 30) hands the FO the rule AND the verdict, so deleting the actual deliverable clause (`first-officer-shared-core.md:105`) leaves the test green; the durable-state clauses (before==after, status: ideation) are satisfied identically by a correct REJECT and an incorrect APPROVE (a gate-held FO with no conn writes no verdict either way); and the only discriminating clauses grade on transcript phrasing (the banned steerable-by-narration class per the README proof policy).

## Proposed approach
Remove the rule/verdict from the runbook (don't supply the answer), and grade on a DURABLE divergence — e.g. run WITH the conn so a wrong APPROVE produces a verdict/advance the assertion can catch, while a correct REJECT routes back. The two graded values must come from independent sources that genuinely diverge between correct and incorrect FO behavior.

## Acceptance criteria
- **AC-1** — the scenario goes RED when the re-anchor clause is removed from `first-officer-shared-core.md` and GREEN when present (divergeable against the actual deliverable, not the runbook's restatement).
- **AC-2** — a correct reject and an incorrect approve produce DIFFERENT durable on-disk state the assertion distinguishes (no reliance on transcript phrasing).

## Test plan
- Drive the scenario both ways (correct/incorrect FO behavior) and assert divergent durable state; prove RED-first against the clause-removed deliverable.
