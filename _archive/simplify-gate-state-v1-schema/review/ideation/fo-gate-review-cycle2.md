# Simplify the unreleased v1 gate-state schema — ideation correction 2

## Chosen direction

Keep the clean v1 direction: derive the current gate from the unique record matching entity status and the last ordered attempt, with canonical digest semantics and no stored pointer or digest-domain field.

## Evidence

- The repaired entity now has exact `## Acceptance criteria` and `## Test plan` sections.
- AC-1 through AC-5 all pass the independent `status --read jc --ac-scan`, with citations from the correction report.
- The stale-approval adversarial fixture still measures 0 of 2 false approvals versus the 1 of 2 pointer-selected baseline.
- The corrected inventory names direct producers/consumers and bounds the surface at 20–22 files with a 24-file tolerance; the full and race plans remain explicit.

## Recommendation

Present this fresh Briefing for Captain decision; approval would enter implementation under the stated 24-file bound, while any further scope growth returns to ideation.
