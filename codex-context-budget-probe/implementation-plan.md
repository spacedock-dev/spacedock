# Native Codex Thread Index Implementation Plan

**Goal:** Replace the all-history Codex session discovery replay with a read-only
native thread-index lookup, then project tokens from exactly one revalidated
rollout JSONL.

**Scope guard:** Keep `spacedock dispatch codex-context-budget --worker` and its
stable JSON unchanged. No observer, socket, app-server client, cache, sidecar,
or JSONL discovery fallback is permitted. Missing or unsafe index evidence is a
non-zero fresh-dispatch result.

## Files and boundaries

- `go.mod`, `go.sum` — add a pure-Go SQLite driver and TOML decoder; no CGO or
  shell dependency.
- `internal/codexsession/index.go` — resolve the effective SQLite home, qualify
  one safe `state_*.sqlite` file, open it via `file:` URI with
  `mode=ro&cache=private`, and issue the bounded parent/worker lookup.
- `internal/codexsession/index_test.go` — TOML precedence, path/schema/row
  ambiguity, WAL/read-only behavior, query limit, and lock/failure tests.
- `internal/codexsession/budget.go` — replace `WalkDir` discovery with the one
  index-selected rollout reader; preserve selected-log metadata, token,
  compaction, freshness, and path checks.
- `internal/codexsession/budget_test.go` — selected-log mismatch and the
  978-file open-count fixture, with existing token and unsafe-path coverage
  adapted to native-index setup.
- `internal/dispatch/codex_context_budget.go` and
  `internal/dispatch/codex_context_budget_test.go` — pass Codex home into the
  reader and retain no-JSON/non-zero unavailable command behavior.
- `internal/ensigncycle/codex_context_budget_live_test.go` — record index-only
  metadata and retain the expected-child live replay contract.
- `docs/runtime-support.md`, `skills/first-officer/references/codex-first-officer-runtime.md`,
  and `skills/ensign/references/codex-ensign-runtime.md` — describe the native
  discovery index and selected-log-only token projection without adding a new
  runtime surface.

## Test-first order

1. Add failing storage-resolution tests for `sqlite_home` TOML precedence,
   relative CWD resolution, and safe candidate qualification. Run
   `rtk go test ./internal/codexsession -run 'Test(ResolveSQLiteHome|OpenStateIndex)' -count=1`
   and record the expected missing-symbol failure.
2. Add the smallest index/config implementation and pure-Go dependencies; rerun
   the same focused tests green.
3. Add failing native-index tests for zero/one/two rows, exact parent/worker
   binding, required schema, normal-WAL visibility, read-only write rejection,
   and bounded lock failure. Run their focused command red, then implement the
   `mode=ro&cache=private` reader and rerun green.
4. Add a failing selected-rollout test with 977 content-sentinel decoys and an
   open-count seam. Run it red, then refactor `ReadBudget` to query first and
   stream only the selected rollout; rerun the codexsession package green.
5. Add/update failing command and adapter behavior tests for index-unavailable,
   path mismatch, below-budget reuse, and above-budget fresh dispatch; implement
   only the command wiring needed to pass and rerun focused dispatch tests.
6. Update the opt-in live harness and runtime documentation after the behavior
   tests are green. Run the live harness only with explicit safe inputs.
7. Run `rtk go test ./...`, `rtk go test ./... -race`, `rtk gofmt -w ./cmd ./internal`,
   a post-format test pass, and `rtk git diff --check`; record any unrelated
   pre-existing formatting drift separately.

## Dependency choices

Use a pure-Go SQLite driver and TOML decoder rather than hand-parsing either
format. The production connection is short-lived, context-bounded, normal-WAL,
and read-only; it must not set `immutable`, `nolock`, mutating PRAGMAs, or run a
migration. The selected JSONL remains the only source of active-window tokens.
