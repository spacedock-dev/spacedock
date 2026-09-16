# Combined stack verification: ab23daf

Exact candidate: `ab23dafeb1bcd5aa6021218e950383f67ff20c0b`, tree `2b9d2c77a69b72659c4fe3cf2298f9913c25a70a`, branch `spacedock-ensign/pi-native-completion-evidence`. Initial and final worktree status are clean; HEAD was checked before and after each test command and stayed unchanged.

This bounded follow-up verifies this combined candidate only. `342edfa8903e0e822ec545d2ee45401352eb6c60` (#801 coverage guard) and `7535aad706e1d0ae88cad8bd45f5ebf4eadf96b8` (#802 promotion/branch correction) are confirmed ancestors. Exact log and source metadata are in provenance.json and final-source.json. Proposed #804 constant changes are absent from this claim; future changed code needs its own applicable verification.

## Results

| Command | Exit | Observed result |
| --- | --- | --- |
| `go test ./...` | 1 | Only TestCodexResolveManifestAgainstInstalledHost fails; every other package passes. ensigncycle passes in 218.530s. |
| `go test ./... -race` | 1 | Same sole failure; every other package passes. ensigncycle passes in 377.281s. No data race diagnostic. |
| `go test -tags live ./internal/ensigncycle -run '^$'` | 0 | Compile-only check passes, `[no tests to run]`; no native/model test selected. |
| `gofmt -l ./cmd ./internal` | 0 | Reports only pre-existing `internal/release/runtime_live_evidence_workflow_test.go`; formatting is not claimed globally clean. No file changed. |
| `git diff --check` | 0 | No diff/whitespace errors. |

Broad commands ran sequentially: normal 2026-09-16 19:35:51–19:39:34 UTC; race 19:40:35–19:47:02 UTC. Compile-only ran afterward, 19:47:43–19:47:45 UTC. Exact command, timing, HEAD and exit codes are in the corresponding result JSON files; full output is retained in normal.log, race.log and live-compile.log. Required write-formatting was previously performed by implementation owners; this dispatch authorized only the read-only formatting check.

## Known failure retained exactly

Both normal and race runs report:

```text
--- FAIL: TestCodexResolveManifestAgainstInstalledHost
    codex_resolve_test.go:44: spacedock@spacedock not installed in codex, but resolver returned "/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json"
```

The same installed-manifest defect was previously explicitly declined by FO for this task. That disposition does not turn either exit-1 suite green. No new failure appeared, so no correction was attempted. The pre-existing formatting output is the same release-test file already recorded in implementation evidence.

## Scope and recommendation

Combined offline verification is complete with the known red resolver result preserved. This is not an all-green check or a live acceptance claim. Native/model tests, host lanes and final tip CI were not run; independent green adversarial checks were not rerun. Candidate, entity and frontmatter were not edited, and no push, CI, PR or rebase was performed. Evidence is confined to this directory and committed locally for FO synchronization.
