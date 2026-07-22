# Validation gate: folder-form entity commit scope

Recommendation: **APPROVE** commit `d4d39ed616f10021d2737f5f919eb243ba62eae0` for landing.

## Capability

`spacedock state commit <slug>` now treats one canonical top-level entity as the commit unit: a flat entity is exactly `<slug>.md`, while a folder entity is the whole `<slug>/` directory. The implementation includes tracked changes and deletions plus new non-ignored reports and artifacts, retains literal Git pathspec isolation, rejects path-bearing pseudo-slugs, and preserves existing conflict-HALT, no-force, retry, no-origin, JSON/text, and clean no-op behavior.

## Validation evidence

- Fresh validation passed all six acceptance criteria with real-Git fixtures, exact committed-path assertions, residual-dirt checks, two-host disjoint/conflicting writes, and invalid-operand zero-mutation controls.
- Focused, full, race, formatting, and cleanliness checks passed at exact candidate `d4d39ed6`; the final surface is 62 production, 363 test, and 9 documentation lines within the approved 2× tolerance.
- Roborev jobs 534 and 537 are closed by deletion and literal-pathspec fixes. Job 540's flat↔folder conversion request is explicitly declined/deferred because conversion is unsupported; its promotion condition is a supported entity-form conversion workflow.
- Live sprint dogfood used this candidate to commit xb's complete five-file validation room as state commit `38dba458198c8a411e6a63c9f52993f6464bdd7e`; the commit contains exactly all five new files and no sibling path.

## Decision

- `approve`: authorize landing and terminal completion.
- `revise`: return a concrete material finding to implementation.
- `hold`: retain validation for a named prerequisite.

