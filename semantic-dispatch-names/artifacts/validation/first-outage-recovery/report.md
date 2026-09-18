# Independent first-outage recovery validation

Recommendation: PASSED for the bounded correction, with native first-outage efficacy and final combined checks still pending. No new finding or candidate repair was required.

Exact tested candidate: `bfc97adbda716201f7420fefa3759965788252f3`; base: `f5af7282dda9151771ff9888dd246553e38afa82`. Code worktree was clean before and after verification. Scope independently measured: seven existing files, +186/-4 (+182 net), inside the authorized seven-path +160–190 estimate. Candidate/frontmatter bytes were not changed.

## Independent probe

Built this checkout with `go build -o /tmp/semantic-outage-validation/spacedock ./cmd/spacedock`, then ran `python3 probe.py /tmp/semantic-outage-validation/spacedock`. Both exited 0. The retained script accepts an explicit checkout binary argument and requires Python 3 and Git; it creates only temporary fixtures. `probe.log` contains exact outcomes.

This probe attacks a boundary not established by the producer's status-only checks: relative paths in a split-root workflow, comparing SHA-256 snapshots of every fixture file, including Git metadata and sentinel bytes. An executable forwarding shim refuses every `dispatch build` with exit 73 and delegates other commands to the actual checkout binary. There is no earlier successful envelope. Short and shortened long identity queries succeed between two refused builds and preserve the entire file snapshot. Only after the uncached checks does the probe invoke successful builds directly; both JSON names exactly equal the queried identities.

Missing entity, undeclared stage, malformed workflow, unexpected positional argument, forbidden mutation option and split-root cross-generation collision all refuse with empty stdout and unchanged fixture snapshots. The independent short-name expectation is literal; the long name must remain within 56 characters and match the subsequent build exactly. Existing producer tests own additional invalid-name/stage-budget/worktree-copy cases and independent wrong-name/collision mutants; those were reused, not rerun.

## Preserved boundaries and documentation

The correction changes only dispatch routing, canonical-query validation, focused tests and the two recovery instruction owners plus command documentation. Native runners, selected-mode controls, duplicate-worker refusal and complete committed-report graders are untouched. Producer focused logs show those existing controls green, alongside dispatch/integration tests; retained wrong-name and bypassed-collision overlays fail as intended.

The command reference accurately describes required readable inputs, declared stage, canonical validation and absence of assignment/stamp/commit/worktree effects. Recovery instructions first reuse a validated same-assignment envelope, otherwise query once. They retain suffix/occupancy checks, the selected named call, its completion target and committed-report obligations. They forbid reconstructing a name, retrying the failed build for identity, changing mode, or treating a failed dependency/hold as successful recovery. Complete executable absence remains unavailable, as explicitly documented.

## Acceptance mapping and outstanding evidence

AC-1: independent short output and exact build equivalence preserve semantic public identity; existing branch proof is unchanged. AC-2: independent split-root alias refusal and long-name equivalence preserve conservative canonical identity; producer budget/invalid-name controls and mutants are retained. AC-3: this correction replaces the superseded cached-envelope-or-hold rule under captain-ruling[2026-09-18] for first build-only outage; offline proof establishes uncached safe identity availability without mutation. Existing one-worker, mode and committed-report controls remain enforced, but this probe does not spawn a worker. AC-4: previous native semantic spawn/reuse evidence remains historical; no new native recovery success or live committed report is claimed.

Normal/race/full-tree gofmt on the final combined stack remain FO-owned and pending. Changed-file formatting evidence is reused and current `git diff --check f5af7282..HEAD` passes. Earlier resolver baseline and retained native failures remain visible and are not waived. No models, broad suites, CI, pushes, rebases, new agents or code edits occurred. Final native CI must establish selected-mode first-outage recovery producing exactly one worker and its committed report before claiming that runtime outcome.
