# Gate review: Codex launcher guarantees multi-agent v2 surface

## Stage

Validation cycle 7 at candidate `64f2f8773a3ba4d848c15f7c9f71bbc45e55c395`, based on `origin/main` `48a7ea0d97042f0e7aaac258e1b77f16157c5281`.

## Recommendation

REJECTED. Hold z5 and route the material Codex First Officer dispatch defect to the runtime and state-commit owner. Do not change the z5 candidate.

## Scope and estimate

The candidate remains five patch-equivalent commits, nine files, `+380/-42`. The ideation estimate was 205-344 added lines, with a 380-line tolerance and a ten-file limit. The candidate stays within the approved numeric and file limits.

## Checklist result

- DONE: Validate exact final base, scope, and patch equivalence.
- DONE: Run focused, full, race, format, diff, and detached checks.
- DONE: Reproduce AC-1 and AC-4 isolated-home behavior.
- DONE: Reproduce AC-2 exact table and fail-closed boundaries.
- DONE: Reproduce AC-3 identical launcher configuration variants.
- FAILED: Complete the required Codex shared live lane. The durable oracle rejects successor dispatch after consume.
- DONE: Preserve artifacts and classify the live finding without candidate edits.

Assessment: 6 done, 0 skipped, 1 failed.

## Acceptance criteria

The authoritative AC scan reports AC-1, AC-2, AC-3, and AC-4 with two citations each and no unevidenced criteria.

## Material finding

The full Codex lane and one targeted rerun fail `recorded-gate-lifecycle`. The retained JSONL shows a successful dispatch build, First Officer narration, no structured `worker.spawn` event, one `wait` event with empty receivers and agent state, and a committed marker. The durable oracle rejects the required post-consume dispatch ordering and ancestry. This is a material outcome defect in Codex First Officer runtime and state-commit integration, not in the z5 launcher candidate. Preserve artifacts at `/tmp/spacedock-codex-final-live-artifacts/codex-shared-scenarios/recorded-gate-lifecycle`.

## Resolution

Science Officer approves the material classification and routes the finding to the Codex runtime owner. A separate local task, `codex-recorded-gate-successor-dispatch`, records the repair and its test TODO. Re-enter this gate only after that repair and a fresh exact-head Codex live validation pass.
