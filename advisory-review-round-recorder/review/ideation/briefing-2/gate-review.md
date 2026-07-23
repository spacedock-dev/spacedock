# Re-ideation gate review: shared 3k round composition

## Capability and change

Replace the rejected 699-LOC advisory-round mechanism with a function-level composition of landed 3k primitives. The revised plan requires one shared entity mutation/CAS path, one expected-room-byte publisher, one canonical Annotation/Resolution parser, one shared load/validate pipeline, and one narrow recorder-owned Feedback Cycles splice.

The advisory value, CLI intent, and no-gate/no-status-effect ACs are unchanged.

## Evidence and reviewed snapshot

- WIP `0e9a313f` remains the counterexample and rewrite seed.
- Entity expectation is now a mandatory shared-helper parameter: full bytes for rounds, landed gates-node/status expectations for existing writers.
- Only `actor:ensign` plus the backward disposition graph authorizes worker triage. A different reviewer actor cannot complete the round.
- No-findings passes `project=false`; authorized all-declines passes `project=true` exactly once.
- Commit 1 must prove stale-room CAS, entity-write failure, rollback, replay, graph falsifiers, and unchanged gates/status/product bytes before any CLI wiring.
- Hard stops: 365 net production LOC before CLI and 500 total. Any journal, second writer/parser, optional expectation, new actor field, or weakened AC requires another ruling.

The exact revised entity, landed contract, and prior drift audit are identified by URI and SHA in `briefing.json`.

## Findings

- The prior duplicate writer/parser/load paths are explicitly deleted rather than compacted in place.
- The three gate-review clarifications—central CAS expectation, worker authority, and projection semantics—are now resolved with red controls.
- Remaining implementation risk is measurable at the pre-CLI checkpoint; it is not deferred to final review.

## Recommendation and decision

Recommendation: **approve re-ideation and proceed to fresh implementation** on the preserved branch. Require the 365-LOC checkpoint evidence before CLI work and Roborev at the end.

Decision requested: approve, revise with a concrete remaining seam, or hold for a named prerequisite.
