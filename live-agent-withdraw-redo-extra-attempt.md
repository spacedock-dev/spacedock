---
id: 47gnqfm1ft6f2hcahz98m2jv
title: Live agent chose withdraw-and-redo over commit-and-retry after a gate-prepare CAS rejection, producing an extra attempt
status: backlog
source: "Found via a dedicated classification workflow, 2026-08-03, diagnosing PR #600 (collapse-gate-approval-ceremony) live CI failures. claude-live (opus), PR #600 run 30754109029 (the rerun), job 91518827506: TestLiveDefaultHeadlessStopsAtGate/default-headless-withdrawn-gate-recovery failed -- the recovered entity had 3 gate attempts instead of the expected 2 (internal/ensigncycle/claude_live_runner_test.go:411-413). Root cause: the FO's first gate prepare call for the withdrawn attempt's successor (relative --artifact/--reference paths) failed with a pre-existing, unchanged-from-main CAS guard in internal/gitsource/source.go ('selected source is not the exact committed file; commit the exact source before preparation'). Instead of committing the exact source and retrying identically, the model diagnosed by dropping --reference and using an absolute artifact path, which succeeded but produced an attempt with no References bound; it judged that attempt deficient, withdrew it, and re-prepared a third time with absolute paths for artifact and both references -- attempt 3 is what it ultimately bound/presented. Classified flakiness (live-model nondeterminism), medium-high confidence: on the same commit, a different attempt of the same scenario passed cleanly while a different scenario failed instead, and two commits earlier both Opus jobs passed fully -- the signature of model variance, not a deterministic regression. One unproven but plausible nudge: this PR's own new skills/fo-gate-lifecycle/SKILL.md Resume paragraph added 'stale -> supersede then replace' wording, and the model's own use of 'stale' to describe its incomplete attempt echoes that new wording, so it may have nudged the withdraw-and-redo choice over commit-and-retry -- worth a look, not confirmed."
started:
completed:
verdict:
score: 0.3
worktree:
issue:
sprint: test-behavior-completeness
---

A live Claude (opus) agent, recovering a withdrawn gate attempt, hit a pre-existing CAS guard (`internal/gitsource/source.go`, unchanged from `main`) rejecting an uncommitted source on its first `gate prepare` attempt. Rather than committing the exact source and retrying the identical command, it worked around the rejection by dropping `--reference` and using an absolute artifact path -- succeeding, but producing an attempt with no References bound. It then judged that attempt deficient on its own, withdrew it, and re-prepared a third time with absolute paths for everything. The scenario's test expects exactly 2 attempts (the original + one legitimate recovery); 3 attempts fails it.

## Why this isn't (yet) classified as a product defect

The CAS guard itself fired correctly and is untouched by collapse-gate-approval-ceremony. The model's workaround was a *legal* sequence of real commands, not a bug tripped by faulty product logic -- it's a live-agent judgment-call variance. Same-commit reruns show different scenarios failing on different attempts, which is the signature of model nondeterminism rather than a deterministic code regression.

## Open question worth checking

This PR's `skills/fo-gate-lifecycle/SKILL.md` Resume paragraph newly added "stale -> supersede then replace" wording. The model's own final message described its incomplete attempt as "stale" before withdrawing and redoing it -- an unproven but plausible echo of that new wording nudging it toward withdraw-and-redo as a first response to a prepare rejection, rather than the simpler commit-and-retry. If confirmed, tightening that wording (to bias toward commit-and-retry on a CAS rejection specifically, reserving withdraw-and-redo for genuinely stale/superseded state) could reduce this class of extra-attempt noise. Needs ideation to determine if this is worth a scoped wording fix or should stay as accepted live-agent variance.

Before a product change, run this journey with the strict XFAIL behavior from
`ts7gq0mr9s3chx2w4wppd1kt`. Convert the TODO only if repeated evidence produces
one stable semantic failure code. If results vary, keep this task as a behavior
investigation and do not hide the variance as XFAIL.

## Out of scope

Any change to the CAS guard itself (`internal/gitsource/source.go`) -- it is behaving correctly, unrelated to this finding.
