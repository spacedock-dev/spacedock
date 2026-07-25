# Independent staff review — rq ideation cycle 4

Verdict: **REVISE**.

## Material

### 1. Planning estimates are incorrectly promoted to implementation authority

Cycle 4 calls the Spacedock estimate a hard cap of 12 files / 1,358 changed LOC and
the Subspace estimate a hard cap of 24 files / 3,203 changed LOC. Those limits are
not derived from a value, safety, authority, or compatibility boundary. They are
planning estimates, but the current wording lets raw file or line count reject a
semantically correct implementation.

That is Material because it can force implementation to optimize for the estimate,
hide a useful helper split, or reopen ideation even when the Captain's contract is
unchanged. Conversely, staying below either cap would not make an authority expansion
acceptable. File count and changed LOC are evidence for drift review, not semantic
authority.

**Required correction:** retain both named-file tables and expected deltas as
estimates, but delete “hard cap” and any statement that tolerance can authorize or
reject implementation. Require implementation to declare its intended files/LOC
before editing and reconcile actual variance afterward. Variance should trigger an
explanation and reviewer scrutiny, never rejection by count alone.

Return to ideation only for a semantic expansion: another public argument, caller-
selected authority, inability to derive mechanics from the bound room, a third
authoritative prepared metadata file, copied source payloads, remote acquisition,
`association.json`, a compatibility path, weaker Briefing binding, or a changed
provider/validator ownership boundary. Tests and gates must not assert repository
file/LOC totals.

## Direction that stands

The sole public entry is exactly `/subspace:r gate <room>`. The model supplies no
entity, workflow, actor, approver, Briefing, destination, manifest, provider, or
terminal coordinate. Fixed code derives and validates those mechanics from the bound
room before provider effects; the integration-private room-only materializer does not
create a second public composition surface.

The folder and flat room layouts match corrected s4. Both converge on
`<entity-root>/<slug>/review/<stage>/briefing-N`, while entity resolution distinguishes
`<slug>/index.md` from `<slug>.md`. rq consumes the landed resolver and correctly
refuses to compensate for a changed s4 contract with public flags.

The manifest correctly binds the canonical Briefing id and full canonical SHA-256.
Subspace reopens the unchanged Briefing, recomputes both identity facts, and rejects a
mismatch before model installation or display. Keeping the exact summary only in the
canonical Briefing avoids a second source of truth.

The prepared-room invariant is also correct: `request.json` plus its located canonical
Briefing are the two authoritative metadata files immediately after s4 preparation.
Later `ROOM/provider` evidence is permitted and does not weaken that invariant.
No `association.json`, rewritten Briefing, copied prepared-room payload, or
compatibility layer is introduced.

The closed manifest parser, provider-owned ephemeral payload child, catchable cleanup,
honest hard-kill residue, canonical validator profile, supervisor ordering, and
terminal-safe summary rendering are proportionate to the existing provider boundary.
They enforce the stated identity, evidence, and display guarantees rather than invent
another authority.

## Dependency and re-review bar

rq implementation must still wait for corrected s4 to land. If its room resolver,
folder/flat placement, request/root-map contract, two-file preparation state, or exact
summary/Briefing binding differs, rq returns to joint ideation.

The next review needs only confirm that the two count limits became non-authoritative
estimates and that semantic reset triggers remain explicit. No implementation,
provider run, compatibility mechanism, or additional design machinery is required for
this correction.
