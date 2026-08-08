## Runtime Live CI

The live lanes prove runtime behavior, not text shape. Static grep checks over workflow YAML or skill prose are not a substitute for launching the real host front door, observing its output, and checking the resulting workflow state.

A runtime regression is proved by one of the 16 exported `TestLiveCommon...`
functions registered in [`runtime-live-ci-registry.md`](runtime-live-ci-registry.md).
Each declaration has an adjacent `liveJourney(...)` call that binds its stable
journey ID, fixture builder, target-scoped TODO owner, runtime-neutral exercise,
and durable assertion. There is no scenario table or runtime runner registry.

The helper selects only the Claude, Codex, or Pi transport from
`SPACEDOCK_LIVE_RUNTIME`. The selected transport launches the current checkout;
the exercise and durable grade are shared. Runtime-specific substrate proofs stay
separate because they verify host boundaries rather than workflow semantics.

### Registry reconciliation

`TestRuntimeLiveRegistryReconciliation` parses the real Go declarations and calls,
the immediately adjacent journey and fixture annotations, the desired registry,
and the executable workflow. It fails on missing or duplicate IDs, malformed TODO
ownership, builder/assertion drift, orphan fixtures, and any lane that does not use
the exact `-run '^TestLiveCommon' -failfast` selector.

Run it after changes to `internal/ensigncycle/`, `internal/livescenario/`, the
registry, or `.github/workflows/runtime-live-e2e.yml`:

```bash
go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$'
```

### Local live execution

Build the binary and export the resolution hooks once:

```bash
go build -o ./spacedock ./cmd/spacedock
export SPACEDOCK_BIN="$PWD/spacedock"
export SPACEDOCK_REPO_ROOT="$PWD"
```

Run the common journeys by selecting one transport. Each command uses the same
16 exported sequential tests and stops at the first non-TODO failure. Claude's
90-minute timeout is a loose suite-wide runaway backstop; Codex and Pi retain
the 40-minute backstop:

```bash
SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 90m -run '^TestLiveCommon' -failfast ./internal/ensigncycle -v
```

Run all three current Claude substrate proofs with one 20-minute backstop:

```bash
go test -tags live -count=1 -timeout 20m -run 'TestLiveMergedTeamModeDispatch|TestLiveBareReachable|TestLiveBreakGlassShimRecovery' ./internal/ensigncycle -v
```

For Codex, install and authenticate the CLI (or set `OPENAI_API_KEY`), then run:

```bash
SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveCommon' -failfast ./internal/ensigncycle -v
```

Run the Pi live proofs locally with the same package versions pinned in CI:

```bash
npm install -g @earendil-works/pi-coding-agent@0.80.10
npm install --prefix "$HOME/.pi/agent/npm" pi-subagents@0.35.1 pi-intercom@0.6.0
export PI_SUBAGENTS_PACKAGE_ROOT="$HOME/.pi/agent/npm/node_modules/pi-subagents"
export PI_INTERCOM_PACKAGE_ROOT="$HOME/.pi/agent/npm/node_modules/pi-intercom"
```

Authenticate with `pi login` or `OPENAI_API_KEY`. Configure the child model with the exact provider-qualified spelling for that authentication path:

```bash
export SPACEDOCK_PI_LIVE_CHILD_MODEL=openrouter/openai/gpt-5.4 # OpenRouter login
# or
export SPACEDOCK_PI_LIVE_CHILD_MODEL=openai/gpt-5.4 # direct OpenAI provider
```

`TestLivePiFrontDoorSmoke` is the only Pi substrate smoke. It checks the front door, child dispatch, durable output, and the boot contract. The grade artifact records both models, durations, and available costs.

```bash
go test -tags live -count=1 -timeout 15m -run TestLivePiFrontDoorSmoke ./internal/ensigncycle -v
```

Run the common Pi journeys separately from that substrate proof:

```bash
SPACEDOCK_LIVE_RUNTIME=pi go test -tags live -count=1 -timeout 40m -run '^TestLiveCommon' -failfast ./internal/ensigncycle -v
```

Without auth, the respective live suite skips locally (Claude/Codex/Pi), except in CI where the lane requires it.

### GitHub setup

| Selected command | Unique evidence | Measured sample or cost |
|---|---|---|
| Claude `TestLiveCommon...` | The 16 registered common journeys | Journey metrics record duration, tokens, model, and available cost. |
| Claude substrate: `TestLiveMergedTeamModeDispatch`, `TestLiveBareReachable`, `TestLiveBreakGlassShimRecovery` | Merged, bare, and break-glass dispatch | Merged baseline: 127s Sonnet and 144s Opus. Cost was not available. |
| Codex resolver and `TestLiveCommon...` | Current-checkout resolution and common journeys | Both PR and release jobs consume Codex metrics. |
| Pi `TestLiveCommon...` and `TestLivePiFrontDoorSmoke` | Common journeys plus one four-part substrate proof | The detail artifacts preserve each run. |

The deletion removes 17 seconds of tmux setup and avoids a 172.5-second duplicate Pi smoke per run.

Workflow: `.github/workflows/runtime-live-e2e.yml`. The offline gate job (`go test ./...`, no secrets) must pass before either live lane burns its environment approval.

- `claude-live` runs the core, shared, merged, bare, and break-glass proofs. Its matrix uses `sonnet` and `claude-opus-4-8`.
- `codex-live` runs the resolver and shared proofs. The PR delta and release ledger consume its metrics.
- `pi-live` runs the coverage guards and one front-door smoke. It uploads the grade, root session, child session, and diagnostics.

All live lanes must test the current checkout, not a remote `--ref next` install. The Codex lane generates a local marketplace under `$RUNNER_TEMP`:

```text
.agents/plugins/marketplace.json
plugins/spacedock -> $GITHUB_WORKSPACE
```

The marketplace manifest uses `source: local` and `path: ./plugins/spacedock`. The workflow's
host-install smoke runs `codex plugin marketplace add`, `codex plugin add spacedock@spacedock`,
and `codex plugin list`, and fails if the listing names `github.com` or `ref next` instead of
the local path. `go test ./internal/cli -run TestCodexResolveManifestAgainstInstalledHost -v`
then confirms Spacedock resolves the installed Codex manifest. The shared Codex runner also
uses `spacedock codex --plugin-dir` so its isolated home receives the same current-checkout
install through the production front door. The Claude lane loads the current checkout directly
via `spacedock claude --plugin-dir "$GITHUB_WORKSPACE"`. Both paths therefore exercise the
current checkout rather than a remote `next` install, while retaining their host-native plugin
and output differences.

A one-off host-only smoke is not enough for either lane: it can prove plugin/login plumbing while missing shared runtime regressions in gate handling, rejection routing, or merge-hook guards. The shared scenarios run real headless hosts, observe output, and check resulting workflow state; jsonl, stderr, and final-message artifacts upload for debugging.
