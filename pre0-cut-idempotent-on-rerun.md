---
id: 83tx3a7zwnvz92d7kq0cwf2a
title: Make the auto-pre0 cut idempotent under a workflow re-run
status: ideation
source: "Validation of `collapse-duplicate-edge-marketplace-routes`, cycle 2, 2026-08-18. Recorded there as a deferred risk with a captain-decision note; the captain elected to file it rather than spend a third feedback cycle. Reproduced verbatim by that validator: `fatal: tag 'v0.27.0-pre0' already exists`, rc=128, red under `set -euo pipefail`."
started: 2026-08-18T02:36:44Z
completed:
verdict:
score:
worktree:
issue:
---

Re-running a stable release's `edge-advance` job after its pre0 tag already reached origin makes the job die on the tag it created. The retired mechanism prevented this; its replacement does not.

## Problem

`release.yml`'s `edge-advance` job auto-cuts a `vX.(Y+1).0-pre0` tag after a stable release. The decision that gates it now compares the tag against the highest known version derived from git tag history.

That scan **excludes the ref being decided**. On a re-run, the pre0 tag the previous attempt already pushed is therefore invisible to the comparison, the decision advances a second time, and `git tag -a` dies:

```
fatal: tag 'v0.27.0-pre0' already exists
```

rc=128, red under `set -euo pipefail`.

This is a regression against the mechanism it replaced. Fed the same moment, the retired `next`-manifest read decided `skip`, because its stable path stamped `next` to `X.(Y+1).0-pre1` — one notch above the pre0 it had just cut. The retired step's own comment named "no colliding pre0 auto-tag" as something it prevented. Nothing supplies that notch on a re-run now.

The trigger is real rather than theoretical. The pre0 step's verify-or-fail poll exits 1 **after** the tag is pushed, so re-running is the natural response to that failure — and three `release.yml` runs in this repository already carry `run_attempt > 1` (v0.23.0, v0.19.5, v0.19.4).

It was classified as a deferred risk rather than material because no value AC fails, nothing wrong publishes, and it fails loudly. It is filed because the promote condition is one observed re-run away.

## Proposed approach

{Ideation confirms. The fix the implementer weighed and set aside is roughly two lines: skip when the pre0 tag the run would form already exists, `git rev-parse -q --verify "refs/tags/v$PRE0_VERSION"`. Confirm that this composes with the existing notch rather than replacing it — the notch handles the old-line patch case, this handles the re-run case, and they are different inputs.}

## Out of scope

The notch mechanism itself (`HighestBareStableVersion`, `HighestKnownEdgeVersion`, `EdgeAdvanceDecision`) — validated across every named boundary input plus a 16-evaluation replay of real release history. Do not re-open it.

The two silent-disable mutations recorded alongside this risk: hardcoding the decision step to `advance=false`, and stripping `fetch-depth: 0` / `fetch-tags: true` from the edge-advance checkout. Both still pass the suite. They are a separate concern, and closing them may need a kind of check this project has been deliberately reducing.

## Expected surface and tolerance

Estimate net LOC change: +25 across 2 files (`.github/workflows/release.yml` and one test file). Tolerance: net +25 ± 20, files 2 ± 1. Report insertions and deletions separately; do not declare a gross tolerance. Semantics changed: the auto-pre0 step becomes a no-op when its target tag exists, instead of failing the job.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - Re-running a stable release's edge-advance job after its pre0 tag exists leaves the job green and the tag untouched.**
This is the measuring AC: the re-run's exit status must be 0 where it is currently 128, and the tag's target commit must be unchanged. Verified by replaying the recorded reproduction — the same tag state that produced `fatal: tag 'v0.27.0-pre0' already exists` must now decide skip. Fails if the job still dies, or if it "succeeds" by moving or force-replacing an existing tag.

**AC-2 - The old-line patch protection still holds.**
Verified by re-running the real `v0.25.1` replay, which must still decide skip for its own reason (the notch), independently of the new tag-exists check. Fails if the new condition masks or replaces the notch rather than composing with it — the regression that would trade one guard for another.
