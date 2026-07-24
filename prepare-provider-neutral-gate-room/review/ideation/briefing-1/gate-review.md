# IDEATION GATE — Provider-neutral gate-room preparation (`s4`)

Recommendation: **APPROVE the corrected design; keep implementation pending until 6y lands.**

## Chosen direction

Add one provider-neutral `gate prepare` operation that derives an open room and stable identity output from the First Officer’s question, Markdown gate review, and selected existing references. Bind any readable canonical Briefing by frozen locator, id, and digest; reject duplicate members before typed decoding; derive the Result association without storing `association.json`.

## Evidence

- The riskiest-first spike proves identical Briefing bytes fail only at xb’s hardcoded basename checks.
- Four correction cycles removed provider-transport simulation, moved lifecycle ownership to 6y’s `fo-gate-lifecycle`, added the missed eligibility reader, and made stdout/URI/media-type semantics executable.
- The final plan reuses existing gate-room and CLI fixtures, one focused preparation test file, and all three existing recorded-gate live lanes.
- Expected surface is 16 files and 1,147 changed LOC, capped at 18 files and 1,434 changed LOC.
- Independent staff review returned APPROVE after all material findings were corrected.

## Boundary

No compatibility layer, `association.json`, Subspace executable/probe/transport, fake provider harness, `present-gate` lifecycle ownership, or broader gate ergonomics. Subspace q0 owns `/subspace:r gate <room>`. Implementation must start from 6y’s final xb-rebased landing or return for a surface reset.

## Decision ask

Approve the design now; retain the application pending until 6y lands, then enter implementation in `.worktrees/spacedock-ensign-prepare-provider-neutral-gate-room`.
