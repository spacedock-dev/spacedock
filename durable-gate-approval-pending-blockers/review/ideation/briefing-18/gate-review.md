# Gate recorder ideation — approved fixture fold

Decision status: the captain approved revision 17 (`lgtm`) with one required fixture. This
revision adds that fixture and changes no other ruling: unreleased v1 in place, one semantic
`gate record` verb plus read-only `gate validate`, exact Result consumption, verified
full-package association before normalization, Git-owned rebind audit history, gates-only
writes, and no migration.

## Design/spec amendment

For `spacedock gate record ENTITY --briefing FILE`, the target logical gate is derived from
the entity's **current workflow stage**, not from `gates.current`. The recorder looks up that
logical gate independently under the entity lock. `gates.current` is the output selection,
never a prerequisite for finding the target.

Therefore, when workflow status has returned to ideation while a closed validation gate is
still globally selected, binding a new ideation Briefing finds the existing ideation gate,
observes its last attempt is closed, appends the successor attempt, and selects it globally.
It does not mutate either prior closure. The approved lifecycle vocabulary remains internal:
this is post-closure **supersede/re-entry**, derived by the binary from semantic input.

The exact command is:

```text
spacedock gate record 3k --briefing revision-18.json --workflow-dir docs/dev
```

The complete byte-anchored before/after projection is
[`cross-logical-gate-reentry.md`](cross-logical-gate-reentry.md). It starts from the real
dogfood state at state commit `71c61fbc`: nine closed ideation attempts, one closed
validation `revise` globally selected, and `status: ideation`. Success appends minimal
ideation attempt 10 and changes the global/legacy-current selection only; both closed
attempt byte ranges remain identical.

## Behavioral proof amendment

The first implementation pass copies the frozen source into `internal/gates/testdata/`
(test discovery-safe), drives the public CLI with `revision-18.json`, and asserts exit 0,
ten ideation attempts, exact new binding/selection, gates-only mutation, and byte-identical
ideation-9 plus validation-1 closures. This serves AC-10, AC-12, and AC-14 without changing
their end-value wording.

The mandatory adversarial edit reinstates the current recorder's global-current
prerequisite for target gate lookup. The full suite must fail in this fixture with the
observed wrong behavior: pointer conflict against validation. That mutant survived any
single-gate lifecycle, so the two-logical-gate state is the independent contrast.

The cycle-23 expected surface remains ~220-360 production LOC touched, ~300-500
test/fixture LOC, and ~80-150 documentation lines at 2x tolerance; this fixture is inside
that test budget. Product HEAD stays `9d279b87` until implementation dispatch.

## Immutable references

- Approved cleaned design/spec and exact command/write examples:
  [`../briefing-17/gate-review.md`](../briefing-17/gate-review.md)
- Exact late Review v1 negative fixture:
  [`../briefing-17/exact-review-v1-result.json`](../briefing-17/exact-review-v1-result.json)
- Exact new complete Briefing input: [`revision-18.json`](revision-18.json)
- Exact cross-gate fixture: [`cross-logical-gate-reentry.md`](cross-logical-gate-reentry.md)
- Full entity reference: [`../../../index.md`](../../../index.md)
