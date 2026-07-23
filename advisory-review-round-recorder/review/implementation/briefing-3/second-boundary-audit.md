# Advisory-round second-boundary audit

Classification: material mechanism drift requiring a narrow design reset; value acceptance criteria remain unchanged.

The prior 500–520 estimate undercounted necessary shared code. The current 670 pre-CLI LOC comprise approximately:

| Mechanism | LOC |
| --- | ---: |
| Shared entity/frontmatter I/O | 186 |
| Canonical review parser and triage | 147 |
| Round loading, paths, artifacts, pointer, projection, and publication | 288 |
| Model and operation routing | 49 |

The corrected two-step implementation is no longer broadly duplicative. Its credible compacted floor is 610–630 pre-CLI and 665–685 total, so accepting it would require hard ceilings near 650/710.

The smallest value-preserving design is one-shot completed-round publication:

1. Accept either reviewer approval with no findings, or a complete reviewer-findings plus reviewer-Resolution plus authorized-worker-dispositions and worker-Resolution log.
2. Create the retained room once and treat it as immutable.
3. Make exact replay a no-op and reject any divergent existing room or log.
4. Publish the entity pointer and optional Feedback Cycles projection together.
5. Remove interim pending persistence, strict-prefix append, existing-log replacement, stale-log CAS, and extended-log restoration.

This preserves the 3j replay, the no-findings/all-declines distinction, exact replay, occupied/divergent target refusal, entity CAS, new-room rollback, digest refusal, lock cleanliness, and all no-gate/no-status invariants. Strict-prefix extension was a proposed mechanism, not a value acceptance criterion, and the caller already records only after triage.

Expected surface is 550–575 pre-CLI and 605–630 total. Binding hard stops are 580 pre-CLI and 640 total.
