# Verification chronology

1. Captured-handoff test added first: `go test ./internal/dispatch -run '^TestBuildMergedModeCompletionSignal$' -count=1` exits1 (0.500s). Generated team-lead has no successful captured delivery and parent completion. red.log.
2. After binding correction: `go test ./internal/dispatch ./internal/contractlint` passes (30.782s /2.081s); `go test -tags live ./internal/ensigncycle -run '^TestSameStage' -count=1` passes18.031s. Existing wrong/missing native identity, report and gate controls retained.
3. `gofmt -w ./cmd ./internal` run; restored only unrelated existing release-test field alignment. Changed files formatted, git diff --check clean.
4. First `go test ./...` exits1: known TestCodexResolveManifestAgainstInstalledHost plus stale completion-literal assertions TestEnsignCycleMechanicalOutputs and TestEnsignCycleGoesRedOnBrokenOutput. normal.log. New failures reported separately to FO, who authorized only four expected/mutated recipient replacements in cycle_test.go.
5. First `go test ./... -race` interrupted by explicit FO instruction before owner alignment; Ctrl-C exits1. race.log is incomplete, not a passing result.
6. After four literal replacements: existing mechanical positive and broken-call/prose-trap negative pass0.586s (`go test ./internal/ensigncycle -run '^TestEnsignCycle(MechanicalOutputs|GoesRedOnBrokenOutput)$' -count=1`). cycle.log.

Final normal and race results will be recorded separately. No local model run, push or CI.

## Final candidate results

Candidate `a34274756ff34fab2cf91c0c32cd701dbac92a8b`. Final `go test ./...` exits 1 only for the known installed-host resolver failure; ensigncycle passes 344.096s, dispatch 85.349s, contractlint 5.693s. Final `go test ./... -race` exits 1 only for that same baseline failure; ensigncycle passes 427.043s, dispatch 125.327s, contractlint 25.456s. All other packages pass or reuse cached passing results. No race diagnostic appears. Full logs are final-normal.log and final-race.log. These are performed required checks with a known limitation, not green-suite claims. Prior FO decline keeps the unrelated resolver finding out of this correction; its disposition is resolved by that decline, its observed failure remains recorded.

Change: 125 insertions / 63 deletions, +62 net across 21 files. Fifteen are mechanical dispatch golden updates; one contains five exact captured events; five are existing generator, test and Claude adapter owners. Scope includes separately authorized four recipient-literal alignments in existing cycle_test.go; no assertion removed. Code worktree is clean after commit.

Remaining acceptance: independent review, then actual corrected-tip CI must demonstrate completion of the Claude cycle-limit leaf, with committed correction/report, actual correlated parent completion, cycle-limit escalation and held gate. No local model drive or live success is claimed. Pi, roadmap refusal and other reported issues remain outside scope; all six variants remain routine CI obligations.
