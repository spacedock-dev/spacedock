# Implementation round 6 evidence

- Candidate: `47020422ba54e7ec61d6f6d10ef2ecad62930d18`
- Reviewer: Roborev `branch_final` panel job 35
- Actual surface: 18 files and 811 changed lines from the `origin/main` merge base
- Prior estimate and stop: 12 files / 216 lines estimated; 14 files / 292 lines maximum
- Captain-approved design reset: exactly 18 files / 811 changed lines, no new files and no net growth
- Acceptance criteria: unchanged; q3, w5, and zbc remain the only accepted model-scoped skips

## Finding dispositions

### Actionable-decision matcher false-green — Material, task-owned

- Released user and workflow: a release maintainer relies on the required Runtime Live E2E recorded-gate journey in the normal same-tip release workflow.
- Observable harm: the lane can accept a root review that lists a decision option and later mentions an unrelated downstream noun without connecting the option to a consequence.
- Affected boundary: AC-4's qualifying root gate review and the non-negotiable no-oracle-weakening boundary.
- Trigger evidence: supported hosts generate free-form decision prose, and the exact counterexample `Decision: approve, reject, or hold before dispatch.` passes the candidate matcher.
- Disposition: fix only this matcher and its exact counterexample in the next correction pass, within the reset ceiling.

### Named q3/w5/zbc skips — declined

- Released user and workflow: a release maintainer sees green required jobs with the named model-scoped TODO skips.
- Observable harm: none to the accepted value boundary; the task already says those green jobs do not prove the skipped capabilities.
- Affected boundary: AC-1, AC-3, and AC-5 explicitly permit only q3, w5, and zbc and exclude them from capability evidence.
- Trigger evidence: the skips are observed and linked to filed followups, but are the captain-approved acceptance boundary.
- Promote when: the captain expands the supported capability boundary or the linked TODO is closed.

### Combined Pi timeout — Deferred risk

- Released user and workflow: a release maintainer runs the combined Pi front-door and recorded-gate package in the normal live job.
- Observable harm: a sufficiently slow valid run could hit the 15-minute package timeout before grading and complete artifact reporting.
- Affected boundary: AC-3 and AC-5's complete Pi proof and retained diagnostics.
- Trigger evidence: the exact retained passing partition is 105.86s plus 352.52s, or 458.38s total, leaving about 441.62s of the 15-minute ceiling; no supported run has timed out.
- Promote when: a supported run times out or measured p95 approaches the 15-minute job ceiling.

### Merge-guard dispatch credit — declined

- Released user and workflow: a release maintainer grades Codex keep-moving evidence in the normal Runtime Live E2E lane.
- Observable harm: the suggested credit would allow terminal merge syntax without the required completed wait and durable stage report to stand in for dispatch proof.
- Affected boundary: AC-4's build, completed-wait, and durable-report invariant and the no-oracle-weakening boundary.
- Trigger evidence: the proposed build-to-merge-without-wait dialect contradicts the binding acceptance criterion; the supported retained trace already supplies build, wait, report, and exact finalization.
- Promote when: a supported trace contains build, completed wait, durable report, and exact entity finalization but the phase-bound grader still rejects it.
