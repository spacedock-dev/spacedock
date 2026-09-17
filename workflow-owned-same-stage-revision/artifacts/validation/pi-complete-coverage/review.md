# Pi coverage correction: independent bounded validation

Recommendation: **PASSED for the fail-fast correction only** at code `7508a773bb7eb42b981016ae9480263941a6047c`; no new material finding. This does not approve the whole task or establish final runtime acceptance. Producer report/evidence was published in state `38fe3050754abf3c5596b0c47dc03bdb52bf00a4`.

## Surface and observation boundary

Exact diff: three files, three insertions and six deletions. `.github/workflows/runtime-live-e2e.yml` removes only Pi common-journey `-failfast`; `docs/runtime-live-ci.md` mirrors that command; `internal/contractlint/live_registry_reconciliation_test.go` includes Pi in the existing no-failfast policy and removes the contrary Pi assertion. Runtime, selector, tags, count, timeout, parallelism, package, gotestsum JSON output and step failure handling remain unchanged. No skip, XFAIL, model, auth, fixture, grader or scheduler change exists in this commit. `candidate.patch` preserves the exact diff.

The existing `TestLiveCommonSameStageRevision` enters `liveJourney`, which enables parallel top-level Pi journeys. Its `runSameStageRevisionJourney` loops serially over six `t.Run` children and does not break on their boolean result. `finishLiveScenario` calls `t.Fatalf` on the failing child's `testing.T`. Source paths/hashes and the exact journey function are retained. This connects command selection to the actual sibling control flow; it is not an instruction-prose search standing in for runtime proof.

## Independent semantic adversarial check

Producer evidence in `../implementation-pi-continue` used `t.Error` in plain and recorded one child with failfast versus six without, both exit1. I independently read its JSON events and scoped source; the retained focused check log reports existing suite-timeout/failfast/registry checks PASS in 1.664s. Those green checks were not rerun.

The unowned claim tested here is whether **fatal**, rather than nonfatal, child failure under a parallel parent still permits sibling execution and preserves the actual gotestsum exit. A temporary standard-library Go file used the candidate's literal six-name inventory, `t.Parallel` on its parent and `t.Fatal` only in plain. Both commands used installed gotestsum with the CI selector, tags, count, timeout and parallel flags. The sole before/after difference was `-failfast`; the target was the temporary file, not the live package. The temporary directory was deleted. This is a bounded Go/command-wrapper exercise, not a modeled Pi runtime or new maintained test framework.

| Observation | With -failfast | Without -failfast |
|---|---|---|
| Children started | plain only | all six, in declared order, once each |
| plain result | FAIL | FAIL |
| Other five children | not started | PASS |
| Parent/package | FAIL / FAIL | FAIL / FAIL |
| Actual gotestsum exit | 1 | 1 |
| Skipped entries | 0 | 0 |

`results.json` retains argv, exits, starts and terminal statuses; before/after JSONL and stdout/stderr retain original observations. `detached-source.txt` preserves the exercised source as evidence. `sources.json` identifies candidate and Go version. Restoring failfast reproduces the omission; removing it retains failure and admits subsequent children. No failure was hidden by an expected failure or an outer successful shell.

The observed original omission was command-level suppression after plain's failure. This correction resolves that mechanism. It does not promise completion after package panic, external cancellation or global timeout, and does not establish successful Pi model outcomes. No such additional failure was exercised or waived here.

## Acceptance and outstanding work

- **AC-1:** historical native correction/open-gate evidence remains unchanged. This command-only review adds no final-tip native evidence.
- **AC-2:** required-review, round and guard assertions are unchanged. Claude's safe hold versus entity-body missing-source recording ambiguity remains under FO HOLD; no missing independent-reviewer obligation is invented. Conventional retained proof remains historical.
- **AC-3:** removal of Pi failfast resolves the observed sibling-coverage omission while retaining the existing routine selector and nonzero failure. All six variants still require actual final-stack-tip native execution on every supported host. No targeted-only exception or coverage exemption was added.

Required final combined `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` remain explicitly FO-deferred pending the independently owned #806 correction/restack. No earlier full-suite result is relabeled as evidence for this candidate. No #806 worktree was read or mutated during this review.

All prior dispositions in report `66bf2a4b0` remain: required-review durable recording location needs a contract decision; roadmap remains a Material exact-target authorization outcome failure; semantic naming / first-outage recovery remains a Material compatibility conflict requiring captain direction. This report neither weakens their graders nor authorizes candidate action or reruns. Earlier run35238049892 red results are historical, not superseded by this detached exercise.

Only the task's report body and this artifact directory were written. No frontmatter or candidate changes, package/full reruns, models, live CI, pushes or new agents occurred. FO owns final gate recommendation, combined checks and live acceptance.
