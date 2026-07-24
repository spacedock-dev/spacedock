# Implementation gate review: advisory-round mechanism reset

## Capability and reviewed change

The first implementation did not produce an acceptable advisory-round recorder. It preserved a focused-test-green WIP counterexample at `0e9a313f`, but crossed the approved size and architecture stops before CLI wiring.

## Evidence

- Actual production surface: 699 net LOC before CLI versus 300 estimated, 600 declared target, and 680 binding stop.
- Material findings: second writer/transaction path, second partial Review & Gate parser, missing retained-room CAS, incorrect Resolution-count completion, and missing risky-path failure coverage.
- The independent audit estimates a 470-500 total-LOC implementation is credible only after re-centering on shared 3k primitives.
- ACs remain unchanged. No behavior or safety condition is narrowed to retain the draft.

The exact entity, recorder contract, and detailed drift audit are identified by URI and SHA in `briefing.json`.

## Recommendation and decision

Recommendation: **revise**. Return to fresh ideation, keep the WIP only as a counterexample, and redesign the mechanism around one shared entity writer, expected-room-byte CAS, canonical Review & Gate parsing, one load/validate pipeline, and a narrow projection splice.

Decision requested: revise to ideation, approve the over-bound draft, or hold for another prerequisite.
