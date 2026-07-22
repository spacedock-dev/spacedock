# Entity-label drive fixtures

Throwaway workflows for the `entity-label-localization` task's validation drive.
Each declares a different `entity-label` so a live first-officer drive over both
proves the FO's *generated* captain-facing prose tracks the declared label rather
than a hardcoded or test-supplied noun:

- `ticket-workflow/` — declares `entity-label: ticket` / `entity-label-plural: tickets`
- `experiment-workflow/` — declares `entity-label: experiment` / `entity-label-plural: experiments`

Each has one entity sitting at a gated stage, carrying a `## Stage Report` and an
`## Acceptance criteria` section, so a FO pointed at the workflow can present a gate.

The expected noun ("ticket" / "experiment") lives ONLY in each fixture's README
`entity-label` field — an independent source the FO must read and resolve at boot.
The validator points the FO at a fixture path (never handing it the label) and
observes "ticket(s)" / "experiment(s)" in the FO's generated gate/status prose. A FO
that ignores the convention emits "entity" and the drive fails. The two-fixture
differential (ticket vs experiment) kills the tautology: two READMEs yielding two
different generated nouns proves resolution, not parroting.
