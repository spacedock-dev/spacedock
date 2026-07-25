# Gate review: First Officer recorded gate lifecycle — cycle 32

## Decision

Should delegated chat decisions stop persisting unauthenticated Captain prose and retain
only the asserted First Officer renderer, required evidence reason, canonical decision,
application, and Git route?

## Capability change

Delete both `--directive` and the rejected `--directive-file` proposal. A delegated
chat close uses `--actor agent:first-officer --reason ...`; direct Captain decisions
remain `person:captain`, and provider decisions remain frozen `--room` Results.

## Why

The retained Claude journey lost only a final period while completing every authorized
lifecycle effect. The recorder authenticated neither the inline text nor a proposed
file against a Captain message; any caller could supply arbitrary bytes. Exact copying
was ceremony, not an enforceable authority boundary.

## Artifact summary

- **Gate review:** this concise decision, mechanism, evidence, and recommendation.
- **Staff review:** independent approval that the cut removes false provenance without
  weakening an enforceable recorder boundary.
- **Entity:** the complete cycle-32 design, corrected ACs/test plan, historical audit,
  nine-file `~+72/-79` surface, prior counterexamples, and stage report.

## Evidence

- History recovery found no explicit Captain ruling that chat grant bytes must be
  persisted exactly; exactness entered through implementation ideation.
- Cycle 32 reports 4 DONE, 1 SKIPPED, 0 FAILED.
- Independent staff review found no material issue.
- Exact and missing-period legacy directive forms become equal byte-clean rejection
  controls; the canonical positive retains renderer, reason, route, and no new
  `adoption-note`.

## Audit boundary

New durable chat state does not claim to prove grant wording, author, message, session,
supersession, or revocation. Those require a future authenticated host-turn producer.
The FO skill and granted/no-conn journeys continue to own the behavioral authority
boundary, including instructions issued later in the active conversation.

## Recommendation

Approve the rigorous cut and return implementation to the existing clean worktree.
