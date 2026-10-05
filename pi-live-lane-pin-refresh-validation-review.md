# Validation gate — pi-live-lane-pin-refresh (mh)

## Outcome

The layer is validated after the captain's floor change, and it is ready for the merge ceremony.

## Candidate

- Tip `5ead9b85c` on `spacedock-ensign/pi-live-lane-pin-refresh`, pushed.
- That tip adds one change to the tip validated earlier: `PiCodingAgentFloor` is now `1.0.0`.
- The stack above it: the observer fix `9755d96d1`, the recorded-gate fix `e5907d00b`, the lane
  models `238a2d59a`.

## Evidence

- **The successor validation** verified at this tip: the floor constant is `1.0.0`, no comment
  carries a version string, the pi-subagents floor and node floor are unchanged, the tests that
  assumed the old floor were corrected rather than deleted, AC-4 holds, and the code halves of
  AC-1, AC-2 and AC-3 hold. The four criteria amended to the shipped mechanism still record their
  lost guarantees, and none was silently strengthened.
- **The lane run at the chain top**, `37338782614` at `238a2d59a`, which contains this layer:
  `offline` green. Its `pi-live` job was still running when this gate was prepared, so the lane
  halves of AC-1, AC-2 and AC-3 are **pending**, not met.
- **Race, local, uncached**: 21/21 packages passed at the previous tip `5a2ad16e9`. This is one
  commit stale, and the successor change is a three-line constant edit.

## Amendment

AC-1 to AC-4 are restated against the shipped mechanism, each recording the guarantee it loses,
and AC-5's Pi-floor clause is amended to the captain-directed raise. This layer claims the pinned
family, the manifest-resolved checks, the integrity comparison, and the single pin source. It does
not claim a standalone pin-revert falsifier, an independent structural guard, an independent
integrity falsifier, or a standing doc-versus-pin check.

## Residuals recorded, not hidden

- AC-2, AC-3, AC-7, AC-8 and AC-10 rest on demonstrations authored by the same change; AC-6 holds
  the one independent live artifact.
- Five journey bindings engaged and two are stale.
- `default-headless-gate-stop` engaged on a fault its own comment says it does not cover.
- The lane halves of AC-1, AC-2 and AC-3 await the run above.

## Question

Do you accept the validation of this layer, `pi-live-lane-pin-refresh`?

This gate approves the validation only. It approves no merge, and no merge runs when it is approved.
The approval stays pending, because the gate's target stage is terminal. The tool `merge guard`
consumes that pending approval later, and only after it sees a merged pull request.
