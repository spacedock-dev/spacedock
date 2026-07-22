## Stage Report: validation

- DONE: Reproduce that one folder-form state commit includes the index plus tracked deletions/changes and new room artifacts, never false-no-ops on artifact-only changes, and leaves sibling/top-level dirt untouched.
  `TestStateCommitFolderIncludesWholeEntity` passed against real Git; its exact four-path assertion and residual porcelain checks fail on any omitted target path, false no-op, or swept sibling/top-level path.
- DONE: Attack exact scoping with Git pathspec-magic names, invalid/path-bearing operands, flat-form exact-file compatibility, and two-host disjoint/conflicting folder writes without force or discarded bytes.
  Focused real-Git tests passed; their assertions compare exact committed names, HEAD/index/worktree/origin bytes, linear history, exit 3, clean rebase abort, and both hosts' surviving content.
- DONE: Verify the final 62-production/363-test/9-doc surface, Roborev 534/537 fixes and 540 deferred classification, plus focused/full/race/format/cleanliness evidence against the committed branch.
  `07811d27..d4d39ed6` is exactly 62 production, 363 test, and 9 documentation changed lines in the approved three files; gofmt and `git diff --check` left clean HEAD `d4d39ed6`.
- DONE: AC-1 - A folder-form state commit durably includes the entity index and all changed or newly created reports and artifacts below that entity.
  `TestStateCommitFolderIncludesWholeEntity` committed and pushed exactly index, tracked report change, tracked deletion, and new artifact, then observed no target-folder dirt.
- DONE: AC-2 - A folder-form state commit never sweeps dirty sibling entities or unrelated state-checkout paths.
  The same test retained flat sibling, folder sibling, and top-level dirt; `TestStateCommitTreatsSlugAsLiteralGitPathspec` also kept matching tracked/untracked siblings off the commit and origin.
- DONE: AC-3 - Nested artifact dirt prevents a false clean no-op even when `index.md` is unchanged.
  The artifact-only phase advanced HEAD with exactly `folder-task/artifacts/evidence.md`; an unchanged rerun returned the established JSON `no-op` result.
- DONE: AC-4 - Flat-form behavior remains compatible.
  `TestStateCommitIsPathScoped` and `TestStateCommitFlatDeletion` required exactly `first-task.md`; JSON/text HALT, no-origin, retry/rebase, and clean-no-op coverage passed in focused/full suites.
- DONE: AC-5 - Concurrency remains entity-scoped for folder form.
  `TestStateCommitFolderMultiWriterHappyPath` proved disjoint artifacts reach linear origin history; `TestStateCommitFolderConflictHalts` proved named nested conflict, clean abort, no force, and preservation of both writers' bytes.
- DONE: AC-6 - The command rejects noncanonical and path-bearing operands without side effects.
  Eight real-Git cases, including both separators, dot/traversal/absolute forms, and the dogfood pseudo-slug, preserved HEAD, index, worktree bytes/status, and remote while returning the invalid-slug diagnostic.
- DONE: Roborev material findings and deferred risk classification.
  Job 534 is closed by index-only/whole-folder/flat deletion tests; job 537 is closed by literal-pathspec and nonexistent-alias tests. Job 540 remains deferred only for simultaneous flat↔folder conversion, which is unsupported and canonically diagnosed as a conflict; promote when a supported conversion workflow exists, and reject any dual-form staging expansion before then.
- DONE: Focused, full, race, format, and cleanliness verification.
  Focused state-commit matrix and canonical dual-form validator passed; `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and `git diff --check` passed with no worktree dirt.
- DONE: PASSED recommendation.
  All six value ACs have executable Git/state evidence, no material finding remains, and the only deferred risk has an unsupported trigger plus a concrete promotion condition.

### Summary

Fresh validation reproduced every promised folder, flat, operand, deletion, and two-host boundary against real Git at `d4d39ed6`. Recommendation: PASSED; job 540's form-conversion request remains a deferred risk until the product supports conversion, and no hidden dual-form compatibility expansion is accepted.
