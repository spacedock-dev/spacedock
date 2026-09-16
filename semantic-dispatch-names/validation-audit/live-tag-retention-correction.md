# Combined-stack live-tag correction

FO authorized the naming-owned call adaptation after a read-only reproduction. Code commit: `2a8491808`; parent: `70165e9b5`.

The lower retention helper requires `(codexHome, publicStream, artifactDir)`. The semantic proof retained its old two-argument call. The correction adds `result.artifactDir` to that call and preserves its existing `native-lifecycle.jsonl` write. One existing code file changed, +1/-1. The helper also retains `codex-native-lifecycle.jsonl` before isolated-home cleanup.

## Red evidence

`go test -tags live ./internal/ensigncycle -run '^$'` exited 1 before correction:

```text
internal/ensigncycle/semantic_names_live_test.go:53:62: not enough arguments in call to codexNativeLifecycleStream
    have (string, string)
    want (string, string, string)
FAIL github.com/spacedock-dev/spacedock/internal/ensigncycle [build failed]
```

Log: `/tmp/semantic-restack-live-compile-failure.log`.

## Green evidence

- The same live-tag compile command exited 0: `ok github.com/spacedock-dev/spacedock/internal/ensigncycle 0.335s [no tests to run]`. Log: `/tmp/semantic-live-compile-fix-green.log`.
- `go test -tags live ./internal/ensigncycle -run '^(TestCodexNativeLifecycleUsesCorrelatedSessionHandle|TestCodexNativeLifecycleParentRolloutLookupFailsClosed)$' -count=1` exited 0 in 0.174s. These existing controls check retained bytes after home cleanup, refusal of retention-write failure, correlation, and missing/ambiguous parent rollout refusal. Log: `/tmp/semantic-retention-fix-green.log`.
- `gofmt` ran on the changed file; `git diff --check` passed.

No live host, broad normal/race suite, push, restack, or CI run was performed. The changed file is excluded from normal/race builds. Existing approved report, gate, and frontmatter remain unchanged. This compile and fixture evidence does not claim a new native run.
