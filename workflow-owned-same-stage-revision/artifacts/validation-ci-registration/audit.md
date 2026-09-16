# Routine live coverage validation

Candidate: bf64ebeca949667c2d60e40a9c595ebdf62375b2. Detached worktree audit; no candidate edits, native/model calls, duplicate broad tests, code push, PR or CI.

## Preserved behavior and intended scheduling

All moved same-stage grader statements compare identically to c7568815b after whitespace normalization and the bound assertion-parameter rename. The conventional rejection body from commandLog through its final grading is byte-identical. Builder code changes only its fixture annotation. Six variants retain exact frozen attempt/plan bytes, selected Git identity/digest, native worker completion, independent-review identity, round cardinality and cycle limit checks. Earlier independent raw/state replay remains valid; this registration change does not replace it.

TestLiveCommonSameStageRevision binds the actual builder/assertion through liveJourney. Its six child variants remain serial; no gap/skip/XFAIL is added. Current Claude/Codex scheduler has one common callable row; Pi selector matches the common entry. Runtime adapters and scheduler parallel suppression remain intact. Offline execution of the actual Codex scheduler with logging-only callable stubs observes one execution of the new journey. Expanded required-host outcomes remain pending updated stack-tip CI; neither prior CI nor local retained Codex proof establishes them.

Three explicit experiments are byte-identical: TestLiveHaikuLoopSpike, TestLiveHaikuLoopSpikeN, TestLiveCodexWaitMatrixFromShippedAdapter. Targeted implementation proofs is removed. Upper naming-layer work was not absorbed.

## Checker mutation matrix

Run checker-audit.py from the detached candidate root; it restores source bytes after each case. Each case executes `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$' -count=1`.

| Mutation | Actual exit | Expected |
|---|---:|---:|
| unchanged candidate | 0 | 0 |
| first scheduler selector changed to Other | 1 | 1 |
| same-stage scheduler row removed | 1 | 1 |
| same-stage row filtered out for Codex during jobs construction | 0 | 1 |
| actual first Claude gotestsum command commented with # | 0 | 1 |
| journey moved into targeted/unclassified section | 1 | 1 |
| first explicit exemption reason removed | 1 | 1 |

The two false greens are real enforcement defects, not proof the unchanged candidate currently omits the journey. No attack matrix beyond these demonstrated holes was added.

### V2a: declared row counted despite actual runtime omission

Mutation: inside `for _, row := range rows`, insert `if runtime == "codex" && row.name == "TestLiveCommonSameStageRevision" { continue }`. Reconciliation passes because its AST scan counts four-field row literals as selected for both hosts without observing executable selection. To prove the behavioral difference without models, replace only each common job's callable with a logger, then execute `SPACEDOCK_LIVE_RUNTIME=codex go test -tags live ./internal/ensigncycle -run '^TestLiveScheduled$' -count=1 -v`. Baseline logs AUDIT_SELECTED TestLiveCommonSameStageRevision once; filtered scheduler logs it zero times. Both complete normally. Exact output is retained.

### V2b: shell comment counted as executed selector

Prefix the actual Claude gotestsum scheduler command at workflow line288 with `#`. Reconciliation still passes: liveTestCoverageError scans raw workflow lines and treats the commented selector as lane execution. The corresponding shell command cannot execute.

## Repair proposal and limits

Both findings concern captain-required omission enforcement, not changed workflow output. Use the existing checker/actual scheduler observation boundary: observe selected callback counts by executing current scheduling with non-native stubs, rather than interpret arbitrary Go. Existing workflow_trunk_test.go already uses yaml.v3 and mappingValue; extract jobs.steps.run, then exclude shell comments in script scalars. Producer proposes those bounded owners; no repair is authorized by this review.

Current correction is +141 net/six changed files; full task layer +393 net/seven. Retained actual normal and race logs each fail only the known installed-host resolver, with ensigncycle passing both. Those failures remain explicit and FO-declined; no wholly green suite claim. No new broad/model run was performed. New tip CI must still execute the promoted journey.
