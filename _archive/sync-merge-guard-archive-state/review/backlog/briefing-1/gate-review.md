# BACKLOG GATE — Durable split-root merge finalization (`rd`)

Recommendation: **APPROVE and dispatch ideation as a prerelease blocker.**

## Capability and value

When `merge guard` reports a split-root task finalized, its path-scoped archive commit must be durably visible from another host. An interruption after the local archive commit must have a supported binary-owned resume, without raw Git repair.

## Evidence

Roborev 2146 traced the current sequence: the `pr-merge` sentinel is committed and pushed, but `merge guard` commits the archive move locally only. The remote can therefore retain an active terminal sentinel after the local checkout has only the archive. Current output says “Next: push” in worktree cleanup prose but exposes no restart-safe state-publication operation for an archived slug.

## Boundary

Reuse the existing state-sync conflict discipline, preserve path-scoped sibling isolation, and add no gate/provider/PR-host semantics. This is separate from 6y’s nine-file lifecycle correction but blocks the durable-decisions prerelease.

## Decision ask

Approve ideation to spike the exact interruption seam and design the smallest supported durable publication path.
