# Gate decision: keep historical schema warnings actionable

## Recommendation

Approve ideation. Validation that is permanently red on immutable history cannot distinguish current corruption, so retired gate-schema fields should warn without determining the exit code.

## Why this matters for 0.27.x

- Historical entities contain fields that the current gate schema no longer accepts.
- Those immutable records make `status --validate` fail forever on otherwise healthy workflows.
- Downgrading only known retired fields restores the command's ability to signal current defects.

## Approved scope

- Define an explicit, reviewable set of retired gate-schema fields.
- Emit visible warnings for those fields.
- Exit zero when retired-field warnings are the only findings.
- Preserve nonzero exit for current-schema corruption and all unrelated validation errors.

## Exclusions

- No migration or rewrite of historical entities.
- No change to existing flat-room conversion warnings.
- No broad weakening of unknown-field validation.

## Proof owed at the next gate

Ideation must explain how a field qualifies as retired, specify warning and exit-code behavior, refine the expected surface, and design fixtures where retired-only history passes while one real current-schema defect still fails.

