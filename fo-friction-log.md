# FO Friction Log — sprint 0198-pre-flip-hardening

Friction encountered while driving 0198. Newest first.

## 2026-06-08 — rtk-filtered git output was stale, hid a commit (qa validation)

The qa validation ensign reported that `rtk`-filtered `git log`/`git diff` returned
**stale output** that initially hid the late `32ceb73e` comm-officer-polish commit and
gave a bogus "files identical" result. It was caught only by comparing raw blob SHAs.
Impact: a validator could validate a stale HEAD and ship a different commit than it
verified. Mitigation this session: the FO compared blob/commit SHAs directly (raw git,
not rtk-filtered) and steered the validator to confirm against the true HEAD. Worth a
tooling note — rtk's git output caching/filtering is unsafe for commit-identity checks
inside a live worktree that a sibling agent is mutating.

## 2026-06-08 — kept-alive implementation ensign committed during validation (qa)

The qa implementation ensign was kept alive for `feedback-to`, but it had an outstanding
comm-officer polish round-trip when it signaled `Done`. The polish arrived late; the
ensign committed it (`32ceb73e`) into the worktree AFTER the FO had already dispatched the
validation ensign into the SAME worktree — a concurrency race. No harm landed (the polish
was AC-neutral, oracles unchanged, suite green at HEAD), but the FO had to steer both
ensigns to resolve it. Lesson for the remaining member (z9): confirm the implementation
ensign has no pending polish round-trip before advancing to validation.
