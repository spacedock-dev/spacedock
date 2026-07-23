# Gate review: Expose ready-gate entities in boot identify JSON — validation

Chosen direction: reject stage-derived readiness and re-ideate the projection around durable gate lifecycle state.

Recommend **reject: current-stage taxonomy cannot distinguish validating, awaiting-captain, and approved-awaiting-merge states**.

## Checklist

- DONE: Emit current gate-stage identities in boot JSON.
- DONE: Preserve dispatchable and ordinary boot compatibility.
- DONE: Pass focused, full, race, and compile checks.
- DONE: Record implementation and review evidence.

## Material finding

Five live tickets all display `validation`, but only three have complete validation reports awaiting gate/merge. The current selector emits all five because it checks only `gate: true`. One of the three complete tickets still points `gates.current` at its old ideation gate, so no durable state advertises its pending validation decision.

The useful boundary is:

- completion verification prepares the canonical Briefing and current-stage open attempt;
- boot/status project `validating`, `awaiting captain`, or `approved, awaiting merge`;
- engage prioritizes awaiting-captain attempts and presents the retained Briefing;
- approval is consumed immediately into merge and terminalization.

Assessment: 4 done, 0 skipped, 0 failed; end-value acceptance fails.

Decision: revise to ideation, preserve commit `c5a96678` as the stage-derived counterexample, and coordinate the projection with `6y` lifecycle wiring and the existing presentation provider.
