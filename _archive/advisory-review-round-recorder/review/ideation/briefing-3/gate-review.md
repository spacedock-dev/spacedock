# Gate review: Extend 3k's recorder to persist advisory review rounds — ideation

Chosen direction: publish only complete advisory rounds into new immutable rooms; exact replay is a no-op and any divergent existing room is refused.

Recommend **approve**.

## Checklist

- DONE: Replace mutable prefix append with complete-round publication.
- DONE: Bound the risk-first implementation at 580/640 LOC.

Assessment: 2 done, 0 skipped, 0 failed.

## Evidence

- AC-1 through AC-5 each have explicit evidence in the cycle-3 Stage Report.
- The independent boundary audit found the two-step design's credible floor was 610–630 LOC before CLI, while one-shot publication preserves the value within an expected 550–575.
- The design retains 3k's parser, lock, entity writer, full-byte expectation, immutable Briefing validation, atomic pointer/projection write, and new-room rollback.
- It removes only interim pending persistence, strict-prefix append, existing-log replacement, stale-log CAS, and extended-log restoration.
- The preserved 670-LOC checkpoint remains counterexample evidence; no acceptance criterion is narrowed.

Decision: approve to resume implementation in worktree `.worktrees/spacedock-ensign-advisory-review-round-recorder` under the 580 pre-CLI and 640 total hard stops.
