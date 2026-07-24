# Advisory-round implementation drift audit

Reviewed WIP commit: `0e9a313fdc3a736a648638e954e6b2c604bac7a6`

Classification: genuine mechanism reset; value ACs unchanged.

The draft reached 699 net production LOC before CLI wiring, above the declared 600 target and the binding 680 stop. The excess is not only size:

- `internal/gates/io.go` adds a second entity rebuild and multi-file publication path beside 3k's existing atomic writer.
- The new publisher does not compare retained Briefing/log room bytes against their preflight values before writing, so stale-room CAS is missing.
- `internal/gates/operation.go` adds a second partial Review & Gate entry/parser and duplicates Resolution validation instead of reusing the canonical vocabulary.
- Round completion is inferred from a count of two Resolutions rather than a validated worker-triage Resolution, so two reviewer Resolutions can wrongly publish the Feedback Cycles line.
- Record and validate duplicate the same Briefing/log load, digest, artifact verification, and parsing pipeline.
- The WIP lacks stale-room CAS, injected write-failure, and rollback coverage for the mechanism that consumed most of the excess surface.

The focused round tests pass, so the commit is retained as a counterexample rather than discarded. It is not a merge candidate.

Reset direction:

1. one shared 3k top-level entity mutation/CAS/atomic-replace helper;
2. one room publication primitive that accepts expected room bytes and refuses staleness;
3. the existing Resolution validator plus one complete Annotation parser;
4. one shared `load/validateRound` pipeline for record and validate;
5. a narrow splice for the one owned Feedback Cycles section.

Credible revised surface is about 470-500 total production LOC including CLI, without changing the approved operation or ACs. Re-ideation must prove this boundary before implementation resumes.
