# PR802 dependency extraction for reordered stack

Candidate: `6419b10bf4263acd88af278afe37d594b0240481`, on rebased naming head `6eee6ab4a`, based on main `5f6e978f9`. No rebase, push, CI, model invocation or approved entity/gate mutation occurred in this task.

## Completion checklist

- DONE: Make PR802 compile and preserve its native evidence checks without depending on PR801, using only the existing lifecycle-retention prerequisite.
  Extracted the lifecycle-retention helper, existing positive/negative tests and shared call-site argument from PR801 commit `65b0e9464d33958fc07f509b70b34fbc8cade8b0`. The unchanged naming caller in rebased `b463fbf58` (original `71e1992c8`) retains its artifact-directory argument. Live-tag compilation changed from exit 1 to exit 0.
- DONE: Record the exact dependency move and focused verification without weakening live acceptance or mutating approved gate state.
  Focused normal and race checks pass; changed-file gofmt and diff check pass. This artifact is the only state surface changed; entity body/frontmatter and gate packages remain unchanged.

## Exact move and scope

- `internal/ensigncycle/claude_runtime_helpers_test.go`: original retention implementation and controls from `65b0e9464`, +20/-5. Retains correlated parent lifecycle JSONL before home cleanup; errors on retention failure. Tests verify exact retained bytes after cleanup, reject unwritable/missing artifact destination, and preserve missing/ambiguous parent rollout and lifecycle identity negatives.
- `internal/ensigncycle/claude_live_runner_test.go`: only the existing shared helper call now passes `result.artifactDir`, +1/-1. The unrelated rejection-scenario `state.bundle` addition from `65b0e9464` was not moved.

Correction total: two existing files, +21/-6 (+15 net). Whole naming layer relative to `5f6e978f9`: 58 files, +1020/-360 (+660 net). No native semantic acceptance assertion or scenario was removed or relaxed. The existing naming artifact write remains unchanged.

## Verification

1. Red: `go test -tags live ./internal/ensigncycle -run '^$'` exited 1 on starting head: `semantic_names_live_test.go:45:76: too many arguments in call to codexNativeLifecycleStream`, have three strings, want two (`red.log`).
2. Green: same live-tag no-test compile exited 0, 0.590s (`compile.log`). No models invoked.
3. `go test ./internal/ensigncycle -run '^(TestCodexNativeLifecycle.*|TestSemantic.*|TestConflictOwnerStampedIdentity)$' -count=1` exited 0, 1.194s (`normal.log`). The selector executes the two existing lifecycle tests and offline stamped-owner test; there are no default-tag `TestSemantic` functions, so no semantic live execution is claimed.
4. `go test ./internal/ensigncycle -race -run '^(TestCodexNativeLifecycle.*|TestSemantic.*|TestConflictOwnerStampedIdentity)$' -count=1` exited 0, 2.513s, no race diagnostic (`race.log`).
5. `gofmt -w internal/ensigncycle/claude_runtime_helpers_test.go internal/ensigncycle/claude_live_runner_test.go` and `git diff --check` completed successfully.

No other PR801 dependency surfaced in live compilation or these focused checks. This is bounded evidence, not a full-suite or native-runtime claim. FO owns final combined normal/race verification and further stack reconciliation. Native proof still requires its original exact semantic spawn, same-handle followup, two marker-bearing commits and retained raw lifecycle evidence.
