# Validation gate — pin-lane-models-in-one-place

## Outcome

Validated on values and structure, with live and offline clauses pending.

## Candidate

Layer 4, the top of the stack. Five lane model ids live once in `internal/release/live_models.txt`, read by Go
through `go:embed` and by the workflow with bash builtins; the Go command that printed them is removed.

## Evidence

- Simulating the workflow's own loop over the real file yields the five records.
- No model literal remains in the workflow.
- `pi_liveenv.go` aliases the release values for both auth paths.
- The claude pin is retained; focused `internal/release` and `internal/claudeteam` checks pass.

## Pending, unmet, and one correction

- **Live smoke and the offline suite are pending.** There is no green lane at this tip.
- The clause "a test fails when they diverge" is **not met as a test**. My read: it needs none. Go and the
  workflow read the same file, so no second copy exists; the way it breaks is the file's format changing, and
  that surfaces in the lane, whose steps consume those values. Recorded as carried by the lane.
- **Correction:** earlier stage reports in this entity cite three tests that its own later commits removed.
  The code is fine; the report was stale, and the entity now says so.

## Question

Do you accept the validation of this layer? This approves the validation only, not any merge.
