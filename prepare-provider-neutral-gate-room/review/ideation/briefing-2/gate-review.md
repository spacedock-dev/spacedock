# Gate review: Prepare provider-neutral gate rooms — ideation cycle 5

## Capability

`spacedock gate prepare` will turn a First Officer's decision question and selected
files into one validated, bound gate room without caller-authored JSON, ids, digests,
locators, request authority, or provider paths.

## Corrected direction

Every selected Artifact and Reference is frozen under `room/sources/`. Its Briefing URI
is clean, contained, and Briefing-relative; its full raw SHA-256 revision verifies the
retained bytes. This replaces the approved design's checkout-escaping `..` paths and
adds no Git-address schema or provider-specific transport.

## Evidence

- The independent-checkout spike made the old
  `../../../../../../../docs/specs/gate-resolution-frontmatter-contract.md` locator
  disappear when state and main were placed under unrelated roots.
- The frozen room copy reopened with the exact expected SHA and bytes after the original
  checkout topology was unavailable.
- Cycle 5 reports 3 DONE, 1 SKIPPED, 0 FAILED.
- Independent staff review found no material issue and approved the revised direction.

## Acceptance coverage

- AC-1 measures zero caller-authored metadata, a complete retained room, exact stdout,
  contained source URIs, and independent-checkout reopen after state commit.
- AC-2 exercises one frozen arbitrary Briefing locator through bind, record, validation,
  and eligibility, with traversal/substitution/digest failures byte-clean.
- AC-3 rejects conflicting duplicate members in request, Briefing, Result, and inventory.
- AC-4 retains provider neutrality and reuses the existing Claude/Codex/Pi lifecycle
  observation rather than adding an override simulator.
- AC-5 recomputes one association from retained inputs and stores no association file.

## Dependency

Approval records the corrected s4 direction but does not authorize implementation yet.
The application remains pending until 6y's exact delegated-authority interface lands;
s4 must re-read that tip and reset again if the overlapping surface differs.

## Recommendation

Approve the frozen-source design and retain its implementation application pending
behind 6y.
