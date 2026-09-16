# Detached durability audit

Candidate: `9454cf3d2`, scheduling base: `4ce49f1ea`. Recommendation: **PASSED**. No material, deferred-risk, or polish findings identified. No candidate changes, code push, CI, broad suite rerun, or live agent lane.

The candidate was checked out at detached HEAD in an independently registered Git worktree under the assigned code path. Two audit-only test files added there used existing fixture setup and new assertions; they were not committed to the candidate. Their exact sources are retained beside this report. Copy `cli_audit_test.go` to `internal/cli/durability_validation_test.go` and `status_audit_test.go` to `internal/status/durability_validation_test.go` in a detached candidate checkout to reproduce.

Commands, each exited 0:

```sh
go test ./internal/cli -run '^TestValidationRetirementRollbackExactCompanions$' -count=1 -v
go test ./internal/status -run '^TestValidationDeliveryAncestryConflictAndRecovery$' -count=1 -v
go test ./internal/statesync -run '^TestReadCheckPreservesRebaseWhilePreflightHalts$' -count=1 -v
```

| Boundary | Adjacent variants and exact invariant | Evidence |
| --- | --- | --- |
| Delivery proof | Missing sentinel, old reachable sentinel with undelivered worktree HEAD, task SHA absent from trunk, real conflict, relative worktree path, missing worktree, missing trunk; refusal preserves conflicting index. Real merge enables proof. | New `status_audit_test.go`, `delivery.log` (1.431s package total). |
| Authority | Real prepare/record/consume, pending before successful proof, consumed after delivery, repeat cannot spend/commit again; state-only delivery versus retirement retaining pending and validation status. | Retained `TestLocalDeliveryProofBeforeTerminalSpend` and `TestStateOnlyLocalDeliveryAndRetirementKeepDistinctAuthority`; complete normal/race runs execute these. |
| Complete retirement transaction | Flat/folder; binary companion bytes; EOF without newline; tracked-deleted companion; staged sibling with additional unstaged content; injected commit failure restores entity/artifact, refs, sibling index and worktree; retry commits active deletion and leaves deleted artifact absent. | New `cli_audit_test.go`, `archive.log` (3.632s package total). Removing rollback, adding unrelated index contents or omitting source paths breaks assertions. |
| Publication | Flat/folder × local-only/remote; fresh clone; failed push leaves local archive, rc=1; archived `state commit` recovery; peer conflict retains local and remote histories and rc=3. | Retained `TestRetirementCommitsCompleteMove`, `TestRetirementFailureRecovery`, `TestRetirementSameEntityRebaseHalts`. |
| Storage refusal | Missing, ordinary directory/parent discovery, wrong branch, detached HEAD, valid-empty; status default/JSON/boot and both filing paths; init/new recovery and definition/discovery exceptions. | Retained `TestStateStorageRefusalAndRecovery` plus extracted storage code inspection: no mutation during checkout check. |
| Rebase safety | A real conflicted rebase has detached HEAD. Read check refuses with index/HEAD/rebase intact; mutation preflight aborts and HALTs before ordinary branch validation. | Detached existing narrow test, `rebase.log` (0.843s package total). |

Trace review: delivery validation precedes both finalization routes; retirement invokes raw archive without consuming gate authority, commits literal source/destination pathspecs, rolls back on commit failure, and retains committed state on publication failure. Pure storage validation is separated from rebase-aware mutation preflight. The new Git ancestry/storage checks make a fixed number of subprocess calls per command and add no unbounded loop or allocation; no scaling test was warranted. Existing archive snapshots and Git operations keep their prior artifact-size behavior.

Scope measured from the exact base: 22 files, 631 insertions, 44 deletions, +587 net; within amended 22-file / +800-net cap (captain binding `binding-1789516510589003000`). The additional fixture changes repair positive-fixture setup rather than widen observable semantics. All three ACs remain unchanged.

Reused required checks: implementation `normal.log` and `race.log` both exit 1 solely at `TestCodexResolveManifestAgainstInstalledHost` because installed host discovery and the installed local manifest disagree. Every other package/test passed; no data-race report. FO explicitly declined the environment-dependent resolver fix as out of scope. This is a limitation, not a wholly green suite claim. Implementation records required gofmt completion and strict docs build success. No independent format-changing pass was run on the producer-owned candidate.
