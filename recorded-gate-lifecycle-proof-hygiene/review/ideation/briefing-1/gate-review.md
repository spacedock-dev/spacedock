# Ideation gate: recorded-gate lifecycle proof hygiene

## Capability

Make the recorded-gate lifecycle evidence easier to trust and diagnose without
changing product behavior, command semantics, skill instructions, or supported
runtime outcomes.

## What the design changes

- Delete the tautological command-text mutant and its orphaned substring parser.
- Restore the existing `gh pr view` quarantine test to a structural positive
  discriminator instead of a shipped-prose map.
- Check the known archive path before reading the active entity in the existing
  Codex and Claude guardrail journeys, so an archival regression produces the
  diagnostic it names.

## Evidence and scope

- Existing real-CLI, refusal, resume, terminal, merge-guard, and live-journey
  tests remain the executable owners of lifecycle behavior.
- The intended implementation touches exactly four test files and no product or
  instruction files.
- The estimated `+14/-64` delta is a planning checkpoint, not an acceptance band.
  Actual variance must be explained and triaged against the deletion-only intent;
  LOC variance alone cannot force a redesign.
- Focused tests, live-tag compilation, the repository full/race suites, and a
  final Roborev panel form the validation evidence.

## Independent concern resolved

The first ideation report made a numeric LOC range binding. The correction removed
that invented obligation while preserving the exact four-file, no-product boundary.

## Recommendation

Approve the corrected deletion-first design for implementation.

## Decision

Approve to implement this exact test-only cleanup; revise to change the evidence
boundary; or hold it outside the pre-release close.
