# AC-1 live loop — rejection-flow in team mode (zqb683j8jth0tyr2eme231e2, implementation)

`TestLiveCommonRejectionFlow` on the composed tree (571017df3) plus this layer, under
the CI codex shim, 2026-08-16.

**AC-1 is met on both runtimes.** A run counts as a conforming green only when the
focused test exits 0 AND it persisted a topology digest whose rows are that run's
branch chain in order — exit 0 alone never counted a chain, which is how two
self-reviewing single-worker runs graded green before this layer.

| Runtime | Conforming greens | Runs used | Budget | Branch |
|---|---|---|---|---|
| codex | 3 consecutive (runs 1-3) | 3 | 8 | `reuse` |
| claude-sonnet | 3 consecutive (runs 3-5) | 5 | 8 | `fresh` |

Baseline this moves off: **0 consecutive conforming greens on Codex** (the retained
loop ledger showed 2 passes in 4 runs, and both passes certified a single-worker chain
that re-reviewed its own fix), and ungraded on Claude, where no topology check existed.

Files:

- `ledger-{runtime}.tsv` — one row per run: verdict, digest check, branch, the full
  observed chain, the running consecutive count, and wallclock seconds.
- `topology-{runtime}-run-N.tsv` — that run's persisted digest, the per-run evidence
  AC-1 requires. Every codex run carries the reuse chain (`spawn`/`spawn`/`reuse`/
  `reuse`); every claude run carries the fail-safe fresh chain (four `spawn`s), which
  is the conforming shape there because the `dispatch context-budget` probe cannot
  succeed on headless `-p`.
- `claude-run-2-red.go-test.log` — the one red, kept because it is the informative one
  (below).
- `loop.sh` — the driver, including the conforming-green rule it applies.

## The two FO adherence observations

**Mode 1 — ends without `gate prepare`.** Recurred on the first codex run of this
layer, before any hardening. Hardening round 1 was applied at the exact instruction
the run's own stream showed the FO following into it: the scenario prompt's goal line
said "stop after … the second validation passes; leave the entity nonterminal", and
the FO reported "Completed the requested two-cycle rejection journey … Ticket remains
nonterminal at `validation`" without preparing the gate — while having loaded
`feedback-rejection-flow/SKILL.md`, whose step 5 says "Done when exactly one fresh open
gate has been prepared and presented". The harness prompt was handing the FO a stop
condition that contradicted the contract it was also reading, and the prompt won. The
goal line now names the prepared gate as the bounded stop. The mode did not recur in
any of the 8 subsequent runs across both runtimes. The grade was not loosened.

**Mode 3 (new) — records the round before the correction lands.** Claude run 2 only.
The FO ran `gate record --round validation/1` while the review log still held the
reviewer's 2 entries, so the recorder reported `entries=2` and the pre-existing
`rejectionRoundSuccess` recognizer (which pins `entries=4`) redded it. That recognizer
is correct: `feedback-rejection-flow/SKILL.md` step 1 is done only when "the target
worker's own entries clos[e] this round's review log", and step 2 is done only when the
recorder reports "the complete round summary, counting every entry this round
accumulated" — correction first, then record. One occurrence in 8 runs; no hardening
round was spent on it, and the recognizer was deliberately NOT relaxed.

**This contradicts the entity's own determined-shape chain**, which orders it
`… REJECTED report → Cycle 1 line → gate record --round validation/1 → state commit →
rework routing → fix Done: …` and cites `feedback-rejection-flow/SKILL.md` steps 1-5
for it. The skill puts the correction first. The entity text is wrong on this one
ordering; the implementation does not depend on it, because the branch-keyed topology
grader orders only worker routing events (`spawn`/`reuse`/`done`) and never places
`gate record` within them. Flagged for the captain rather than resolved here.

## Reproduce

Build the worktree binary, put `ci-codex-shim-replica.sh` first on PATH as `codex`,
then `loop.sh codex 3 8` / `loop.sh claude 3 8`. Raw streams (~130 KB public per run,
plus the native rollout) stayed in the job scratch dir and are not durable; the
digests and the one red's test log are what was kept.
