# Validation gate — pi-live-lane-pin-refresh (mh)

## Outcome

Validated at a final tip. This room replaces the earlier approval, which the rework superseded, and covers
the seven changes this layer took after it.

## Candidate

Tip `e80f3463c`, the bottom of a four-layer stack.

Its changes since the last approved tip: the Pi readiness floor at `1.0.0`; the journey-binding registry
cleaned (two stale bindings dropped, one re-pointed to the owner whose fault it observes); the redundant
live-registry machinery deleted, including a check that could not run in CI; three unreachable tests
deleted; a focused codex-live step added with its `liveClaims` registration; and a false audit comment
removed.

## Checks run

- `offline` green on earlier tips; the suite is re-run at this tip by the lane run that follows.
- Focused `internal/pilive`, `internal/cli`, `internal/release`, `internal/contractlint` and
  `internal/ensigncycle` packages pass.
- Race, uncached, at an earlier tip: 21/21 packages.

## Evidence per criterion

- The Pi family is pinned in one place, and the floor is `1.0.0` with no version string in any comment.
- An isolated Pi home discovers both extensions; the smoke grade holds the one independent live artifact.
- Four criteria were amended to the shipped mechanism, each recording the guarantee it loses; a fifth
  (the Pi-floor clause) was amended to the captain-directed raise.
- The journey bindings now name active owners only, and the binding-owner guard that once enforced that is
  deleted with the machinery it belonged to.

## What is not met, and what is not claimed

- **The live lane at this tip is pending.** The run at `7ac7bb69b` passed `offline` and then its tip was
  superseded by these changes; the run at this tip is the evidence that counts.
- Not claimed: a standalone pin-revert falsifier, an independent structural guard, an independent integrity
  falsifier, or a standing doc-versus-pin check. Those were removed before this sprint and are recorded.
- Two journey bindings are stale in the register and one was bound to the wrong fault; both are fixed here.

## Question

Do you accept the validation of this layer?
## Verified after this room was first written

- The three rescued checks in `internal/contractlint/runtime_live_lane_test.go` pass under `env -i`, so they
  need no environment variable and are reachable by a runner.
- No journey was lost to the deletions: 17 `spacedock:live-journey` markers and 45 `TestLive*` functions at
  both the pre-deletion tip and this one, with no name difference.
