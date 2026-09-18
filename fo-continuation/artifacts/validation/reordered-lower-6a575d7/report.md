# Reordered lower-stack verification: 6a575d7

Exact final head `6a575d7ae140f3bbf8dd99b59da37f6ac981ae82`, base `5f6e978f9`, naming head `1bcee3400`, durability head `0dee00c21`. Initial/final worktree clean; every required command recorded the same HEAD before and after. Verification completed; this is not an all-green result, fresh native acceptance, or blanket merge clearance.

## Completion checklist

- DONE: Verify the reordered lower stack main→802→803→804 with required normal/race/formatting checks and live-tag compile, preserving actual failures.
  All four commands ran once, sequentially. Normal and race exit1 solely on the known cached-pre0 resolver failure; every other package passes. No data-race diagnostic. Formatting and live-tag no-test compile exit0. Exact command/timing/HEAD JSON and complete logs accompany this report.
- DONE: Record exact candidate and scope equivalence, including the moved Codex retention prerequisite, without claiming fresh native acceptance.
  #803's four patches and #804's continuation patch match their backup ranges; the suffix patch differs only in import context. All eight original #801 commits are absent from lower ancestry. #800 is inherited through main. No models, new agents, push, rebase, CI or approved entity/gate mutation.

## Results

| Check | argv | Exit | Wall time | UTC start → end |
| --- | --- | --- | --- | --- |
| normal | `go test ./... -p 1 -timeout 30m` | 1 | 404.133s | 2026-09-18T19:55:22.110292+00:00 → 2026-09-18T20:02:06.243121+00:00 |
| race | `go test ./... -race -p 1 -timeout 30m` | 1 | 592.702s | 2026-09-18T20:02:46.162244+00:00 → 2026-09-18T20:12:38.863784+00:00 |
| format | `gofmt -w ./cmd ./internal` | 0 | 0.128s | 2026-09-18T20:13:29.854869+00:00 → 2026-09-18T20:13:29.983319+00:00 |
| live-compile | `go test -tags live ./internal/ensigncycle -run ^$` | 0 | 2.207s | 2026-09-18T20:13:42.352738+00:00 → 2026-09-18T20:13:44.559536+00:00 |

Both broad commands report exactly:

```text
--- FAIL: TestCodexResolveManifestAgainstInstalledHost
    codex_resolve_test.go:44: spacedock@spacedock not installed in codex, but resolver returned "/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json"
```

This is the previously declined installed-host baseline, not a green suite. The first frozen candidate3952164 had an additional registry failure; its red evidence remains in ../reordered-lower-3952164. On this corrected candidate, full contractlint passes normal0.621s and race3.371s; ensigncycle passes normal197.655s and race210.927s. The existing registry owner's focused positives and five negative mutations are reused from semantic-dispatch-names/artifacts/implementation/reorder-802-registry/report.md rather than repeated.

`gofmt -w ./cmd ./internal` changed only the existing Default/Options alignment in internal/release/runtime_live_evidence_workflow_test.go, identical to the authorized baseline at pi-native-completion-evidence/artifacts/validation/final-tip-1a28/format-diff.txt. Pre-command bytes are retained in format-baseline-before.txt; format-diff.txt records the delta, and format.json records hashes. The validator restored only that exact formatting delta, returning the file to its original bytes. No new formatting drift or product change occurred; repository-wide formatting cleanliness is not claimed.

The live-tag command passed0.394s with `[no tests to run]`. It proves compilation only. Final diff check passed and final-source.json records clean status.

## Patch and dependency equivalence

range-803.txt compares `e0a250460..backup/reorder-803-20260918` against `1bcee3400..0dee00c21`: all four entries are `=`. range-804.txt compares `56d77c77f..backup/reorder-804-20260918` against `0dee00c21..6a575d7`: continuation is `=`; suffix's `!` is solely surrounding import context changing from gitsource to gates. Its literal and comment delta are unchanged.

The retention prerequisite6419b10bf remains +15 net across two existing files: helper/controls from65b0e9464 plus the shared artifact-directory argument. Earlier retained comparison in ../reordered-lower-3952164 shows identical helper patch identity and the deliberately excluded state.bundle addition. The only subsequent prerequisite is1bcee3400, one checker expectation line requiring exactly one scheduled TestLiveSemanticNamesCodex. No #801 coverage rewrite or acceptance weakening was imported; original focused proof and mutation results are in the producer artifacts.

ancestry.json checks each of the eight original #801 commits from backup/reorder-801-20260918 and records ancestor exit1 for all. The inherited #800 commit4ce49f1ea records ancestor exit0 through main. lower-log.txt makes the exact current stack visible. Upper branches were outside this validator's write scope and were not changed.

## Summary

The reordered lower stack completes required offline verification with the single known resolver failure preserved. The newly discovered registry prerequisite is corrected and passes in both broad suites; exact owned patch equivalence and final clean head are retained. No fresh runtime outcome or merge approval is inferred from these checks.
