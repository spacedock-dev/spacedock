# New verification finding: package time budget exhausted

Exact candidate 0997b6b7590a5a89ba979194238e93f7f47d94e4. Normal command ran once, exit1; HEAD unchanged before/after. Full stdout/stderr and stacks are in normal.log; timings/exit in normal-result.json. This finding is separate from the exact previously declined installed-manifest resolver error, which also recurred.

`internal/cli` exceeded the default10m package timer while TestStateCommitFromRootResolves had run14s; stack is the unchanged twoHostStateWorkflow helper's `git clone -q <temp-bare> <hostB>`. `internal/ensigncycle` exceeded the same timer while TestDurableTaskJourneys had run5m24s and atomic_first_worker13s; stack is a path-scoped Git commit in durableJourneyFixture. These stacks do not establish a ten-minute deadlock in either leaf.

Read-only comparison: all internal/cli content, ensigncycle/cycle_test.go and shared_keep_moving_durable_test.go are unchanged from ab23daf. That candidate's completed normal suite recorded CLI148.305s and ensigncycle218.530s on the same Go1.26.7 darwin/arm64. Other current package durations also increased: dispatch42.996→142.850s, release27.403→78.623s, status54.732→170.152s. FO separately observed boot/status operations taking minutes and slow Git binary startup. This supports a cumulative host/Git-slowness hypothesis; it does not prove the cause or make interrupted package checks pass.

Four evidence fields:
- Normal workflow: required once-only combined deterministic validation on the frozen stack.
- Observable harm: CLI and ensigncycle did not finish, so final-tip broad proof is incomplete.
- Authority: value-ac[AC-3] Local deterministic evidence must retain actual outcomes rather than call timeouts green.
- Trigger evidence: normal.log lines7–110 and114–220 plus exact prior/current duration comparison.

Proposed kind: verification/infrastructure evidence failure, candidate-specific defect not established. Proposed materiality: blocks claiming completed combined verification; not a proven product regression. Ownership: FO validation/environment decision, not an authorized candidate fix. Proposed disposition: HOLD race and further broad testing, investigate with minimal bounded probes only after distinct FO authorization. FO explicitly issued HOLD for race and forbade timeout increase, rerun, baseline fix and candidate edits during analysis.

Proposed falsifier: one exact existing CLI test TestStateCommitFromRootResolves on unchanged HEAD, then one exact existing ensign test TestDurableTaskJourneys/atomic_first_worker if the first completes. Use -count=1 -v and original timeout, save actual duration/output; no broad retry. A reproducible isolated hang would oppose a merely cumulative slowdown hypothesis. Passing leaves would refute deadlock for these observed stack locations, but would not finish either broad package. Proposal only until FO disposition.


## Authorized isolated probes

Distinct FO disposition authorized exactly the two sequential isolated tests, default timeout and unchanged candidate, while holding broad retry/race and candidate edits. Both completed with exit0 and HEAD0997 before/after:

- `go test ./internal/cli -run '^TestStateCommitFromRootResolves$' -count=1 -v`: leaf18.62s, command34.993s.
- `go test ./internal/ensigncycle -run '^TestDurableTaskJourneys$/^atomic_first_worker$' -count=1 -v`: leaf18.98s, command26.052s.

Logs and structured exit/duration records are adjacent. These results refute a persistent deadlock in the exact two observed leaf operations. Combined with unchanged code and cross-package duration inflation, they support cumulative host/Git slowness, without proving a global root cause. The timed-out broad run remains exit1 and incomplete. Race and broad retry remain held pending FO disposition; no candidate edit or timeout increase occurred.


## Subsequent distinct FO disposition and outcome

FO retained the timeout as unresolved infrastructure/verification evidence, with persistent observed-leaf deadlock refuted, then explicitly authorized exactly one normal rerun and one race run sequentially: `go test ./... -p 1 -timeout 30m`; `go test ./... -race -p 1 -timeout 30m`. Thirty minutes was tied to measured approximately3x slowdown; command-only accommodation, no code/CI/assertion change. No further retry or budget enlargement authorized.

Both authorized serial commands completed on unchanged0997 with exit1 only for the exact previously declined resolver error. All other packages pass; ensigncycle280.959s normal/264.114s race; no timeout or data-race diagnostic. Full original timeout remains retained and cannot be called green. These results resolve incomplete package observation for deterministic validation under the recorded command conditions; they do not conclusively establish the infrastructure cause.
