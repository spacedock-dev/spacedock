# Implementation correction gate: advisory-round shared composition

## Capability and reviewed change

The cycle-2 checkpoint replaces the duplicate architecture and passes its red-first CAS, rollback, replay, parser, actor, and projection controls. It stopped before CLI at 683 net production LOC.

## Evidence

- Entity mutation, room publication, and round loading are now genuinely shared.
- The original 365-LOC pre-CLI estimate was invalid; an independent mechanism map supports 500-520 pre-CLI and 555-575 total.
- The current 683-LOC checkpoint still contains concrete canonical-validation, multiple-triage, section-splice, URI, duplication, and whole-operation failure-test defects.
- ACs remain unchanged and no new writer, parser, loader, field, journal, or lifecycle mechanism is authorized.

## Recommendation and decision

Recommendation: **revise in implementation**. Correct the named defects under hard ceilings of 540 production LOC before CLI and 600 total. Stop again if meeting those limits would weaken CAS, rollback, canonical validation, exact projection, or byte-clean refusal.

Decision requested: revise within the measured bounds, approve the 683-LOC checkpoint, or hold for another design reset.
