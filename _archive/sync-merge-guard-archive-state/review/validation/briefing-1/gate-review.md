# Validation gate — durable split-root merge finalization

Recommendation: approve candidate `fb1c9dbbd06e52f86d9fa3d073651ef9c0f4066a` for merge and terminalization.

## Capability delivered

- `merge guard` publishes the exact archive commit before reporting remote durability, so a fresh host observes one archived terminal task instead of the stale active sentinel.
- `state commit <slug>` resumes a clean interrupted archive publication idempotently without creating a second archive commit.
- The shared standard-library publisher preserves linear peer history, halts with recoverable conflict evidence, and never force-pushes or autostashes.
- Archive moves and rollback treat Git-looking task names as literal paths; the real `:(glob)*` counterexample no longer sweeps sibling dirt.
- Inline and no-origin workflows report truthful local outcomes without moving an unrelated code remote.

## Validation evidence

Cycle-2 validation reports 8 DONE, 0 SKIPPED, and 0 FAILED. AC-1 through AC-6 all have independently reproduced code, command, and state evidence.

Focused split-root, archive-resume, peer/conflict, inline/no-origin, identity-collision, wrong-root/branch, and rebase-ownership tests passed. `go test ./...`, `go test ./... -race`, formatting, and diff checks also passed on the corrected tip.

The first validation cycle found one material AC-4 defect: raw Git pathspec magic could include a dirty sibling in the archive commit. Commit `fb1c9dbb` literalizes forward and rollback pathspecs. The same validator reproduced the original counterexample, verified exact two-path archive deltas, and proved the rollback preserves the sibling's staged binary diff byte-for-byte.

## Deferred risk

Unrelated tracked sibling dirt combined with a simultaneous non-fast-forward remains outside the promised direct-push path. The supported path is proven, and recovery succeeds after the dirt settles; this becomes material only if publication through both conditions is promised.

## Decision requested

Approve to merge and terminalize this exact candidate, or revise only if a material acceptance or durability defect remains.
