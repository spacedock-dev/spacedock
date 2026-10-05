# Final independent validation: frozen combined tip 0997b6b7

Recommendation: **PASSED for the deterministic #806 correction and shared-runner reconciliation**, with actual red broad exits and pending live obligations below. Exact candidate `0997b6b7590a5a89ba979194238e93f7f47d94e4` atop `c572d62f05ae57d643dbaab5ecfb7de4a2ccff31`. Candidate HEAD stayed exact before/after every executed command and the final worktree is clean. This is not an all-green stack or final live acceptance recommendation.

## Scope and authorization

P1/P2 Material AC-1 evidence fixes and the sole shared-runner conflict reconciliation were explicitly FO-authorized. Validator owned only evidence and entity-body reporting. No candidate edit, frontmatter/gate action, rebase, push, PR or CI action occurred. No native/model test ran. Independent probes ran on a disposable detached checkout of this exact candidate under the assigned worktree, then it was removed. Prior green checks were reused by owner; their source hashes are in reused-evidence.json.

## AC-1: completion identity and parent order

PASS locally. Independent tests use the full original tip35149242496 streams, not only implementation projections. Default sync produces exact spawn65/done66 and validation spawn80/done81; async same-stage produces spawn112/done116 on the exact same fresh worker before gate135. The original implementation advance remains73. Removing report evidence fails the full default lifecycle. Omitting or replacing the async artifact root fails completion grading; artifact propagation therefore remains load-bearing.

Eight additional detached variants cover previously unowned combinations: altered retained-parent cwd, Unicode byte change in original task, nonterminal child tail, stale child replay after a second fresh sync dispatch, cross-worker call and run evidence, completion moved after advancement while keeping its earlier timestamp, and valid EOF without newline. Invalid variants fail both existing lifecycle and route graders; EOF normalization passes. Exactly ten subtests across two independent tests pass (`independent.log`, source retained in detached_pi_independent_test.go). Candidate code was never altered for these probes.

Reused owned captured red→green and focused controls establish explicit success/exit0, unique result, call/run/owner/agent/task/cwd/epoch attribution, successful child terminal stop and duplicate/missing/error/order rejection. Accepting a missing exit as zero, using directory UUID as run identity, matching completion prose, or crediting a prior worker would falsify those controls. A single shared child verifier serves both observers; no synthetic event, second Git parser, runtime controller or new acceptance rule was introduced.

## AC-2: outcome guards and other hosts

PASS locally for observer separation. Existing whole-journey absent-gate and uncommitted-report controls remain owned and green in reconciliation logs and the completed serial package suites. Their expected outcome remains rejection; bypassing either durable check fails the retained control. Prior original auto absent-gate outcome is not reclassified as success. Captured historical Pi status/wait, native async and Claude/Codex controls remain in the current passing ensigncycle suite.

#801 independent native-identity/recorder report and falsifiers at `workflow-owned-same-stage-revision/artifacts/validation-native-identity/` were reused, not commissioned or rerun. Exact diff versus c572d62 confirms the reconciled same-stage runner retains `countRouteEvents(routes, routeSpawn) > 1` rather than counting turns and adds Pi routeErr without replacing that predicate. Rejection flow retains non-Pi actual recorder command-log/exit interception; Pi alone uses its retained artifact-root route observer and error propagation. Thus #801's true native worker identity and recorder exit, not display names or enclosing shell success, remain the authority. Runner diff is retained here.

## AC-3: actual combined checks and explicit limits

**Original required normal run failed with more than the known resolver defect.** `go test ./...` exited1 after614.91s. CLI and ensigncycle reached the default10m package timer while active leaf tests were only14s and13s old; stacks wait in ordinary Git clone/commit operations. It also reproduced the exact installed-manifest error. Original normal.log and normal-result.json remain untouched and are not described as resolver-only or passed.

FO initially held race and further runs. Read-only comparison found all CLI code and both timed-out ensign fixture/helper files unchanged from the earlier completed ab23 run on the same Go1.26.7. Several unrelated package durations had inflated roughly3x. Distinct FO authorization allowed exactly two sequential isolated leaf probes, unchanged candidate/default timeout: CLI TestStateCommitFromRootResolves passed18.62s; ensign TestDurableTaskJourneys/atomic_first_worker passed18.98s. They refute persistent deadlock at the observed leaves, support cumulative host/Git slowness, and do not prove its global cause or complete the timed-out broad run.

FO then explicitly authorized one full normal rerun followed by one race run with `-p 1 -timeout 30m`, based on measured slowdown. This changed command-level concurrency/observation budget only; all packages and original assertions remained selected, with no test/code/CI policy change or outcome threshold relaxation. Both completed:

| Command | Exit | Evidence |
| --- | --- | --- |
| `go test ./...` | 1 | Original resolver failure plus CLI/ensigncycle default10m package timeouts. |
| `go test ./... -p 1 -timeout 30m` | 1 | Only exact resolver failure; all other packages pass, ensigncycle280.959s; no timeout. Wall1080.283s. |
| `go test ./... -race -p 1 -timeout 30m` | 1 | Only exact resolver failure; all other packages pass, ensigncycle264.114s; no timeout/data-race diagnostic. Wall708.361s. |

Every command's complete stdout/stderr, exit, UTC times and exact HEAD before/after is retained. No further retry or budget increase was made. The infrastructure cause remains qualified; successful serialized completion resolves the incomplete deterministic package observation without erasing the initial failure.

Fresh normal and race logs each contain exactly the previously FO-declined installed-manifest error:

```text
codex_resolve_test.go:44: spacedock@spacedock not installed in codex, but resolver returned "/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json"
```

That limited disposition is applicable only after checking these exact logs; it does not turn either exit1 green. No new assertion failure or data race was reported by the completed serial suites. Read-only formatting still reports only previously documented baseline `internal/release/runtime_live_evidence_workflow_test.go`; diff check is clean. Required write-formatting and current-head live-tag compile-only were performed by owner and reused, not duplicated; compile-only selected no tests.

## Findings, live obligations and recommendation limits

P1/P2 deterministic evidence fixes are validated. No new candidate-specific defect was found. The new timeout verification finding, FO holds, bounded probes, further run authorization and actual results are retained in timeout-finding.md plus logs; global host/Git cause is supported but not conclusively established. No baseline repair or waiver is implied.

**AC-3 final native CI remains pending.** All six same-stage variants must execute on each host at the final tip: plain, review-required, separate-review-required, round-required, round-missing, cycle-limit. Earlier Pi plain failed and five others were unstarted; unstarted/skipped/cancelled is not passed. Current offline evidence does not relabel old CI results green or supply current native coverage.

The separate observed Claude roadmap-authorization outcome failure remains HOLD/route for decision. In run35148797301 the exact-target prompt authorized creating/committing roadmap-strategy.md, but the FO refused under entity-fit rules. Retained authority: `fo-continuation/artifacts/validation/tip35148797301/report.md`, state commit698cf989d. Pi's same case passed in35149242496; that is not a waiver of the Claude failure. The unchanged strict case must run at the next full tip CI. No speculative prompt patch, new permission requirement or grader relaxation was made.

Canonical checklist: 3 DONE, 0 SKIPPED, 0 FAILED for the three assigned validation work items; this counts completed review/execution/reporting, not green test commands. Broad exits remain1 and final native acceptance remains pending. Recommendation PASSED is bounded to deterministic #806/#801 reconciliation evidence above, not release approval.
