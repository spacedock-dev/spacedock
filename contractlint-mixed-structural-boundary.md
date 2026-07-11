---
title: "Separate legitimate contractlint structure checks from mixed behavioral prose claims"
status: backlog
source: "Contractlint antipattern sweep, 2026-07-11: mixed checks in boot_resident_closure, first_officer_eager_references, fn_consolidation_structure, legacy_teamcreate_layering, dispatch_recovery_value_binding, reconcile_class_binding, and structural_checks."
score: 0.30
id: b71p9hwyyscckemmer2v3mbb
---

## Problem

Several contractlint tests mix valid closure/topology/schema assertions with claims that a phrase guarantees live behavior (for example eager loading preventing a hunt, a nonempty reference proving capability, or a prose classifier proving an exemption). The valid structural signal is obscured by the tautological behavioral portion.

## Proposed approach

Split each mixed test at its assertion boundary. Preserve structural checks such as reference closure, frontmatter/schema validity, deduplication, block locality, and absence constraints. Move executable behavior claims to owning fixtures/live tests; delete claims with no executable owner instead of widening prose matching.

## Out of scope

Rewriting unrelated contract instructions or adding a meta-grep to police greps.

## Acceptance criteria

**AC-1 (VALUE) - The listed mixed tests retain their independent structural signal without claiming that instruction prose proves behavior.**
Verified by: focused contractlint tests that demonstrate the retained structural failure mode using a fixture or generated temporary document.

**AC-2 - Every removed semantic assertion is either covered by an owning behavior test or recorded as non-enforceable guidance.**
Verified by: an assertion-to-owner inventory in the implementation report and focused behavior-test output for executable cases.

## Test plan

Handle one test file at a time with a before/after assertion inventory. Run focused contractlint tests after each split, then the full Go suite; run host live lanes only for behavior assertions moved into those lanes.
