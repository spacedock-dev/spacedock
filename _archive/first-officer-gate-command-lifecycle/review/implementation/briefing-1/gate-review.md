# Implementation gate review: defer the First Officer gate lifecycle

## Capability and reviewed change

The checkpoint implements and tests much of the six-event gate lifecycle, but places 6.2 KB of procedure in the boot-resident First Officer core and still lacks required live-spawn evidence and several adversarial controls.

## Evidence

- Focused deterministic lifecycle tests pass; production Go remains unchanged.
- Full normal/race suites fail the boot-core and all-host prompt-load ratchets by exactly 6,197 bytes.
- The already-gated headless/no-conn path can still present before bind/open validation.
- Codex completed and consumed the six gate commands but emitted no observed successor spawn; Claude/Pi stopped on external authentication.
- WIP `cabdef33` is preserved as a counterexample and test bank, not a merge candidate.

The exact entity, gate contract, and detailed topology audit are identified by URI and SHA in `briefing.json`.

## Recommendation and decision

Recommendation: **revise**. Return to fresh bounded ideation for a deferred non-user-invocable gate-lifecycle skill, one small boot trigger covering every gate entry, unchanged prompt ceiling, and the named missing proof arms. Do not rebaseline the boot core or narrow the spawn AC.

Decision requested: revise to ideation, approve the over-budget resident-core layout, or hold for another prerequisite.
