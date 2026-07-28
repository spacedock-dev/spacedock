# Implementation round 6 evidence

- Candidate: `47020422ba54e7ec61d6f6d10ef2ecad62930d18`
- Reviewer: Roborev `branch_final` panel job 35
- Actual surface: 18 files and 811 changed lines from the `origin/main` merge base
- Prior estimate and stop: 12 files / 216 lines estimated; 14 files / 292 lines maximum
- Captain-approved design reset: exactly 18 files / 811 changed lines, no new files and no net growth
- Acceptance criteria: amended after exact-tip proof; q3, w5, zbc, and 9w are the only accepted model-scoped skips, and none proves its quarantined capability

## Final verification and review

- Final code candidate: `55b3b13316571f1a84215f489c205e66327faeca`
- Final surface: 19 files / 803 changed lines from the `main` merge base,
  within the captain-approved 19-file / 811-line reset.
- Green deterministic proof: `gofmt -w ./cmd ./internal`; `git diff --check`;
  focused Pi coverage/definition/workflow guards; focused gates, contractlint,
  and ensigncycle packages; `go test ./...`; `go test ./... -race`; and
  `go vet ./...`.
- The live-tagged recorded-gate selection emitted only
  `TODO(9w59t6m1qc46hccd54p04z2j)` in 0.00s; no Pi model session ran.
- Roborev `branch_final` panel job 85 found the false ledger and operator-guide
  capability claims. The captain classified them Material, retained the
  workflow selection as an explicit visible quarantine, and authorized the
  narrow reset recorded below. Re-panel job 90 returned `No issues found.`
- No push, CI dispatch, or further live retry was performed.

## Exact-tip Pi classification history

- `final-ce4ac943-pi-recorded-gate-pinned`: failed after 203.39s; the lifecycle
  dispatched and persisted the handoff, but the root review invented the bound
  Briefing digest.
- `final-ce4ac943-pi-recorded-gate-pinned-retry`: failed after 93.94s; the root
  review presented the exact Briefing ID/digest, then stopped before delegated
  application/dispatch and produced zero child sessions.
- Prior green evidence remains retained at
  `se0-b25-focused/pi-recorded-gate-final-tree-retry` (201.02s) and
  `cycle2c-pi-recorded` (352.52s). It records history, not current reliability.
- Disposition: TODO
  [9w59t6m1qc46hccd54p04z2j](../pi-delegated-gate-continuation-reliability.md)
  permits only `TestLivePiRecordedGateLifecycle` to skip. Pi front-door remains
  active, and this skip is not capability evidence.

## Finding dispositions

### Panel job 85 false ledger and operator claims — Material, task-owned

- Released user and workflow: a release maintainer reads the Pi coverage ledger
  and operator guide while running the required Runtime Live E2E job.
- Observable harm: the registered 9w skip could leave a green job while the
  ledger and guide still claimed recorded-gate live capability evidence.
- Affected boundary: AC-4 and AC-5 require every permitted skip to remain
  explicit non-proof so the restored signal is trustworthy.
- Trigger evidence: every Pi job selects the skipped test while 9w is open; the
  stale `mode: live` row and two guide claims were normal current surfaces.
- Ruling: change only the existing row to `gap` with the explicit TODO,
  registered/selected quarantine, and non-proof reason; correct both existing
  guide claims. Keep the workflow selection so CI surfaces the named skip.
- Captain reset: at most 19 changed files / 811 changed lines; no new mode,
  machinery, oracle weakening, live retry, push, or CI.

### Actionable-decision matcher false-green — Material, task-owned

- Released user and workflow: a release maintainer relies on the required Runtime Live E2E recorded-gate journey in the normal same-tip release workflow.
- Observable harm: the lane can accept a root review that lists a decision option and later mentions an unrelated downstream noun without connecting the option to a consequence.
- Affected boundary: AC-4's qualifying root gate review and the non-negotiable no-oracle-weakening boundary.
- Trigger evidence: supported hosts generate free-form decision prose, and the exact counterexample `Decision: approve, reject, or hold before dispatch.` passes the candidate matcher.
- Disposition: fix only this matcher and its exact counterexample in the next correction pass, within the reset ceiling.

### Named q3/w5/zbc/9w skips — declined

- Released user and workflow: a release maintainer sees green required jobs with the named model-scoped TODO skips.
- Observable harm: none to the accepted value boundary; the task already says those green jobs do not prove the skipped capabilities.
- Affected boundary: AC-1, AC-3, and AC-5 explicitly permit only q3, w5, zbc, and 9w and exclude them from capability evidence.
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
