# Design gate: warn on exact retired gate-schema fields

## Recommendation

Approve implementation. The revised design restores useful validation without migration or permissive unknown-field handling, and independent review found no remaining material issue.

## Selected design

- Recognize only two exact historical string-scalar encodings: `gates.current.gate` pointing to an existing record, and Briefing `digest-domain: canonical-bytes`.
- Remove those nodes only from a validation clone; retain the original YAML node and stored bytes unchanged.
- Emit visible warnings on active explicit validation and exit zero when they are the only findings.
- Keep malformed lookalikes, wrong locations, unrelated unknown fields, and current-schema corruption fatal.
- Keep readiness and mutation authority on current status plus ordered attempts, never on retired pointers.

## Proof and safety

- AC-1 measures the retired-only baseline changing from exit 1 to exit 0 with exact text and JSON behavior.
- AC-2 and AC-3 prove real corruption and near-miss encodings remain fatal.
- AC-4 structurally checks both returned source nodes, disk bytes, and status/order authority.
- AC-5 pins the public and internal compatibility wording while preserving silent provider-evidence behavior.

## Surface

Estimate: **+155 net LOC across 6 files**. Tolerance: **±45 net LOC and at most 7 files**. Command grammar, canonical writers, stored formats, archive policy, and gate authority remain unchanged.

## Delivery proof owed

Implementation must run focused gate/status fixtures, full and race suites, formatting, and strict documentation build. Validation must reject any surface overrun or AC without behavioral/state evidence.

