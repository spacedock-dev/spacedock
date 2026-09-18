# H1 producer evidence

Candidate: `1f9220634fe198388ac041b54f308032e1c8e5c2`. Diff from rejected 9a614522: one file: 18 insertions / 3 deletions (+15 net), claude_runtime_helpers_test.go. Existing runner unchanged.

- Red command: `go test -tags live ./internal/ensigncycle -run '^TestSameStageCommittedReviewHold/(not_provided|unrelated_absence)$' -count=1 -v`; exit 1, 3.915s. Both exact reviewer cases failed before helper correction.
- Normal: `go test -tags live ./internal/ensigncycle -run '^TestSameStage' -count=1 -v`; exit 0, 17.309s.
- Race: `go test -race -tags live ./internal/ensigncycle -run '^TestSameStage' -count=1 -v`; exit 0, 18.609s.
- gofmt on changed file and git diff --check pass. Between normal and race the literal newline/tab regex was rendered with equivalent escaped notation for readability; pattern semantics unchanged. Race covers final bytes.

Paragraph boundaries are blank lines (including horizontal whitespace). Wrapped lines stay in the same note. This bounds textual association and does not understand arbitrary prose or resolve unrelated clauses within one paragraph. All existing durable controls remain strict. No model, broad suite, push, CI, rebase, new agent or authority mutation. Combined checks and native acceptance remain pending.

Original briefing and first four review-log entries are byte-preserved. Own producer proposal and correction entries append actual actions. FO owns the one proposed round 2 publication and Cycle projection.
