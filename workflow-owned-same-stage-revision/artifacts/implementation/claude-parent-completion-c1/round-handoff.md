# C1 handoff to existing producer

Current canonical pointer validation/2 and Cycle2; next proposed publication validation/3. Confirm current state before publication. Original reviewer package is four entries: reviewer finding/revise and transcribed FO FIX/revise. Existing schema matches prior H1 package; no new schema or machinery.

Producer should copy this package into its owned artifacts and append actual actor:ensign proposal/Resolution and eventual correction/evidence closure. Reference `annotation:zz1yqc2w2k-c1` and `annotation:zz1yqc2w2k-c1-fo-fix`. Do not invent completion before implementation. FO alone projects Cycle3 and records complete round once; retained reviewer correction recheck does not publish again.

Exact reproducer: copy `wrong-worker-completion.jsonl` over `internal/dispatch/testdata/claude-completion-route.jsonl` on detached candidate; run `go test ./internal/dispatch -run '^TestBuildMergedModeCompletionSignal$' -count=1 -v`. Current result is false pass. `baseline.jsonl` must pass; `missing-completion.jsonl` must fail. Report pinned by immutable state commit and SHA in briefing. Reviewer does not alter candidate/frontmatter/gate/round state.
