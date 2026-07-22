# Gate review: Approved-pending blockers, execution holds, and dispatch eligibility (h1) — ideation

**What you are looking at:** the application layer — what a recorded decision DOES. Three artifacts: this summary, the full design, and the recorder contract whose application sections this task owns.

**Chosen direction, in two honest halves:**
1. **Ships: the application record.** The one-use advance authorization — pending, consumed exactly once through the EXISTING transition path, superseded on post-closure drift, not-applicable on hold — extending the recorder binary (never a second writer). Its consumer is proven, not assumed: the 0260 Commander consumed pending advances all drive. Its red fixture is real: the recorder's own ideation gate briefly held TWO pending advances (caught at preflight); the extended acceptance criterion now demands at most one pending across all attempts after any supersede.
2. **Declined: the blocker-satisfaction evaluator and hold authoring.** No live consumer exists — every recorded approval to date carried zero declared blockers. The decline is recorded with its promotion condition: a live consumer appearing promotes that half to its own captain-approved design round. What ships still HONORS any present blockers/holds fail-closed — missing or unqueryable state never reads satisfied and never consumes an approval.

**The spike:** exactly-once consumption exercised against the real frontmatter writer before design lock.

**Preflight:** the cross-attempt criterion extension and the decline-fallback rewording were its two findings here — both applied and verified.

**Recommend approve** — this is the sprint's own thesis practiced at design time: build the half with a proven consumer, decline the half without one, on the record.

**Decision:** approve = a pending advance to implementation. Revise = annotate. Hold = discuss.

---
*Companion artifacts: the full design (entity snapshot), the recorder contract.*
