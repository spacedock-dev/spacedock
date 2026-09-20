# C1 producer correction

Candidate `1b39a96d928739c362d40837a74d8e6c825c0312`; production route remains from a34274756. Existing test now retains delivery by exact send ID, after a matching successful result under the same parent. A later completed notification must name that parent before the send is credited. Recipient identity alone never joins two workers.

The two reviewer negative files were copied exactly into dispatch testdata; hashes are in negative-provenance.json. The original successful native capture remains unchanged. Three files changed, 83 insertions / 54 deletions = +29 net; most diff movement wraps the existing predicate in a local replay closure for the three existing evidence cases. No new runtime observer or parser framework.

- Before correlation fix: `go test ./internal/dispatch -run '^TestBuildMergedModeCompletionSignal$' -count=1 -v`, exit 1 / 0.567s. Original and missing completion expectations pass; exact wrong-worker/no-success case returns true instead of false.
- After fix: same focused command, exit 0 / 0.657s; all three cases pass.
- `go test ./internal/dispatch`: exit 0 / 73.898s.
- `go test ./internal/dispatch -race`: exit 0 / 74.732s.
- Changed Go file formatted, git diff --check passes.

One initial command accidentally ran formatting/focused test in root before the correct worktree red run. Root target file remained clean; its passing result is excluded from this proof. All candidate edits and counted checks use assigned worktree.

Full normal/race at a34274756 remain historical production-code coverage with the known resolver failure, not green checks for this new test head. FO explicitly authorized focused package verification for this test-only correction. No model, broad-suite repeat, push, CI, frontmatter or gate mutation. Final corrected-tip Claude cycle-limit and all six routine variants remain native CI obligations. The original four reviewer/FO entries and briefing are byte-preserved; only actual producer proposal and closure append. FO owns one round3 publication and Cycle/captain-escalation handling.
