# Reordered lower-stack validation: 3952164

Candidate `3952164f97272b295206c95bec9c99ccebce5571`, main base `5f6e978f9`, naming head `6419b10bf4263acd88af278afe37d594b0240481`. Initial and final worktree clean; HEAD unchanged across the command. This candidate is not verified green. No native/model claim or merge clearance.

## Required checks and disposition

- FAILED: `go test ./... -p 1 -timeout 30m` exited 1 after 420.818s. Start 2026-09-18T19:46:19.741555+00:00; end 2026-09-18T19:53:20.560014+00:00. Exact argv, heads and epoch timestamps are in normal.json; full output is normal.log.
  Two failed tests: `TestCodexResolveManifestAgainstInstalledHost` at codex_resolve_test.go:44 reports `spacedock@spacedock not installed in codex` but resolves `/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json`; this matches the known baseline exactly. NEW `TestRuntimeLiveRegistryReconciliation` at live_registry_reconciliation_test.go:255 rejects `unexpected scheduled callables: map[TestLiveSemanticNamesCodex:1]`. All other packages pass; ensigncycle passed in 152.754s.
- SKIPPED: race suite, `gofmt -w ./cmd ./internal`, and live-tag compile. After the new failure was reported, FO explicitly directed completing only the active normal command and stopping checks until the registry prerequisite is corrected. None of these three commands was started. No formatting changes or restoration occurred.
- DONE: Record patch equivalence and dependency scope without fresh native acceptance.
  range-803.txt maps all four old patches to identical reordered patches. range-804.txt maps the continuation patch identically and shows only import context changing from gitsource to gates around the unchanged suffix patch; the owned literal/comment edit is unchanged.

## Dependency boundary

The moved prerequisite `6419b10bf` is +21/-6 (+15 net), two existing files. Its lifecycle helper/control patch has the same stable patch ID as the original `65b0e9464` helper change: `6534b2f73ecc550a607522332d0795797653cda8`. The shared call gains result.artifactDir; the original unrelated state.bundle addition remains outside the lower stack. Original and moved patches are retained alongside the producer's already-green focused evidence in semantic-dispatch-names/artifacts/implementation/reorder-802/report.md; those tests were not duplicated.

ancestry.json verifies none of the eight original #801 commits from backup/reorder-801-20260918 is an ancestor of this lower head (all merge-base ancestry checks return1). #800 commit4ce49f1ea is an ancestor through merged main; it is not #801 leakage. lower-log.txt lists only the naming layer and moved prerequisite, four durability patches, and two continuation patches above main.

## Finding and ownership

New registry failure is an evidence/dependency defect: normal supported lower-stack verification encounters a scheduled naming callable omitted by the checker expectation inherited below #801. Observable harm is a red registry check after restack. Authority: contract[internal/contractlint/live_registry_reconciliation_test.go#TestRuntimeLiveRegistryReconciliation] scheduled callables must reconcile with the declared live registry without losing missing/duplicate/unexpected coverage. Trigger is this exact normal-suite output. FO acknowledged and routed the prerequisite to #802's owner. Validator performed no repair or rerun; forthcoming correction requires a new exact-head result.

## Summary

Normal-suite verification completed honestly with the known resolver failure and one new registry dependency failure. Owned #803/#804 patches preserve their behavior-bearing changes, but final lower-stack verification remains incomplete pending the separately routed correction. No code, approved entity/frontmatter, gate, push, CI or upper-branch mutation occurred; only this evidence directory is committed.
