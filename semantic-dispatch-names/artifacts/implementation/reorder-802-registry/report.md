# PR802 registry dependency correction

Candidate `1bcee3400` on `6419b10bf`; lifecycle-retention extraction is preserved. One existing file changes: `internal/contractlint/live_registry_reconciliation_test.go`, +1/-1 (zero net). The pre801 checker now requires exactly one scheduled `TestLiveSemanticNamesCodex`, using its existing supported-callable predicate. No PR801 coverage rewrite, exemption, scenario removal or new harness was imported.

## Completion checklist

- DONE: Make the pre801 registry reconcile the802 scheduled naming proof while still rejecting missing, duplicate, or unexpected scheduled callables.
  Actual checker passes with the current scheduler; temporary missing/duplicate/unexpected mutations each exit 1 with the corresponding diagnostic. Mismatch and unregistered negatives also exit 1. All temporary source/document mutations were restored byte-for-byte before commit.
- DONE: Retain the observed red and focused positive/negative evidence for the minimum dependency correction.
  Logs accompany this report. Original checker fails on the observed unexpected semantic callable; corrected normal/race and live-tag no-test compilation pass. Entity body, frontmatter and gate packages are untouched.

## Commands and exact results

- Before correction: `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$' -count=1` exits 1: `live_registry_reconciliation_test.go:255: unexpected scheduled callables: map[TestLiveSemanticNamesCodex:1]` (`red.log`).
- After correction: same command exits 0, 0.367s (`green.log`).
- `go test ./internal/contractlint -race -run '^TestRuntimeLiveRegistryReconciliation$' -count=1` exits 0, 1.814s, no race diagnostic (`race.log`).
- `go test -tags live ./internal/ensigncycle -run '^$'` exits 0, cached no-test compilation (`compile.log`); no model invocation.
- `gofmt -w internal/contractlint/live_registry_reconciliation_test.go` and `git diff --check` pass.

Each negative reruns the exact normal command against the actual changed checker, without prose assertions or checker substitution:

1. Remove the semantic jobs append statement: exit 1, semantic callable appears 0 times, want 1 (`missing.log`).
2. Duplicate that statement: exit 1, semantic callable appears 2 times, want 1 (`duplicate.log`).
3. Retain its name but substitute `TestLiveBareReachable` as callable: exit 1, scheduled name/callable mismatch (`mismatch.log`).
4. Append an additional row named/calling `TestLiveUnexpected`: exit 1, unexpected scheduled callable (`unexpected.log`). The checker parses source; this proves its unexpected-entry refusal without compiling an undefined live test.
5. Temporarily replace the registry's semantic test name with `TestLiveUnregisteredProof`: exit 1, semantic test is unregistered and replacement has no declaration (`unregistered.log`).

The first four mutations affect only `internal/ensigncycle/scheduled_live_test.go`; the fifth affects only `docs/runtime-live-ci-registry.md`. Python try/finally restored their exact original bytes. Only the authorized one-line checker change remains in the product commit.

No other dependency surfaced in these bounded checks. Full normal/race remain FO-coordinated on the final upper tip. This report makes no full-suite or live efficacy claim. No push, rebase, CI, models, body/frontmatter or gate edits occurred.
