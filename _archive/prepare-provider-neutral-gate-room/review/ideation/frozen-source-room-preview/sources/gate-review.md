# Gate review: provider-neutral frozen-source room

## Decision

Should `spacedock gate prepare` freeze every selected Artifact and Reference inside the
retained room, using only clean Briefing-relative `sources/...` URIs and full raw-byte
SHA-256 revisions?

## Capability

The First Officer supplies a decision question, one Markdown gate-review Artifact, and
optional supporting files. Spacedock derives the room, ids, media types, locators,
digests, request authority, and open gate attempt. No caller authors JSON or output
paths.

## Corrected direction

Copy every selected source into `room/sources/` before publishing the room. The
canonical Briefing addresses only those frozen copies. The original main or state
checkout can then move or disappear without changing what the captain reviews.

This preserves Review v1's existing Briefing-relative URI rule and avoids inventing a
Git-address schema, provider transport, compatibility wrapper, or `association.json`.

## Artifact summary

- **`artifact:gate-review` — this file.** Summarizes the capability, corrected design,
  evidence, boundaries, and decision requested from the captain.
- **`artifact:staff-review` — `sources/staff-review.md`.** Records the independent
  adversarial review: no material findings; the frozen-source mechanism is the smallest
  durable correction; large package growth is the only deferred risk.
- **`artifact:entity` — `sources/entity.md`.** Freezes the complete currently recorded
  s4 ticket. Its canonical cycle-5 design landed at state commit `c8eb0f49`; the copy
  also includes the later gate history, problem, command contract, room/request layout,
  AC-1 through AC-5, test order, expected 16-file `+1,090/-161` surface, dependency on
  6y, and cycle-5 report.

## Evidence

- The old cross-checkout URI
  `../../../../../../../docs/specs/gate-resolution-frontmatter-contract.md` stopped
  resolving when main and state were placed under unrelated roots.
- A frozen room-owned copy reopened with the same expected raw SHA and exact bytes after
  the original checkout topology was unavailable.
- Cycle 5 reports 3 DONE, 1 SKIPPED, 0 FAILED.
- Independent staff review returned APPROVE with no material findings.

## Boundaries

- The room contains presentation inputs, not provider transport.
- Provider association remains derived and unstored.
- Implementation remains held until 6y's delegated-authority interface lands and s4
  rechecks its overlapping surface.
- The current recorder cannot yet bind this arbitrary canonical filename; that is the
  behavior s4 exists to implement.

## Recommendation

Approve the frozen-source direction. Do not start implementation until 6y lands and the
overlap check passes.
