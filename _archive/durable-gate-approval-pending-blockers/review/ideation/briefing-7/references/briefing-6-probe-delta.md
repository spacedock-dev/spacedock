# Briefing 6 Probe delta

This is a frozen presentation summary of the provider-owned exact-Briefing result and
comparison. The provider's canonical records remain outside the Briefing package and
are joined by Briefing id.

- Probe: Can this probe-and-reanswer mechanism operate without Spacedock, and where
  are the persistent Probe and each Briefing-bound ProbeResult serialized?
- Comparison: `changed`; the established answer remains supported.
- Still holds: the mechanism is provider-owned and works with ordinary Review & Gate
  Briefings without Spacedock.
- New evidence: this review-room instance stores its Probe history at
  `../probes.jsonl` relative to Briefing 6.
- Narrowed limitation: the instance path is now concrete, but the provider contract
  still does not mandate one universal filesystem layout.

The Briefing 6 Subspace presentation did not visibly surface the fresh ProbeResult or
this comparison. That is an observed Subspace product gap, not a 3k UI requirement.
