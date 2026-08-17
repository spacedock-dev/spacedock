## Runtime Live CI

The live lanes prove runtime behavior, not text shape. Static grep checks over workflow YAML or skill prose are not a substitute for launching the real host front door, observing its output, and checking the resulting workflow state.

A runtime regression is proved by one of the 17 exported `TestLiveCommon...`
functions registered in [`runtime-live-ci-registry.md`](runtime-live-ci-registry.md).
Each declaration has an adjacent `liveJourney(...)` call that binds its stable
journey ID, fixture builder, target-scoped TODO or strict-XFAIL owner, runtime-
neutral exercise, and durable assertion. A TODO skips only when the target
cannot run. An XFAIL runs the target and accepts its typed semantic failures.
There is no scenario table or runtime runner registry.

The helper selects only the Claude, Codex, or Pi transport from
`SPACEDOCK_LIVE_RUNTIME`. The selected transport launches the current checkout;
the exercise and durable grade are shared. Runtime-specific substrate proofs stay
separate because they verify host boundaries rather than workflow semantics.

### Registry reconciliation

`TestRuntimeLiveRegistryReconciliation` parses the real Go declarations and calls,
the immediately adjacent journey and fixture annotations, the desired registry,
and the executable workflow. It fails on missing or duplicate IDs, unclassified
live tests, malformed TODO ownership, builder/assertion drift, orphan fixtures,
and an incorrect common-suite selector.

Run it after changes to `internal/ensigncycle/`, `internal/livescenario/`, the
registry, or `.github/workflows/runtime-live-e2e.yml`:

```bash
go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$'
```

The state checkout changes independently from a code commit. Run the mutable
owner join during sprint close and before a release:

```bash
SPACEDOCK_LIVE_STATE_DIR=docs/dev/.spacedock-state \
  go test ./internal/contractlint -run '^TestRuntimeLiveTODOOwnersAreActive$'
```

This check fails when a TODO or XFAIL names an inactive entity. Stable code CI
does not fetch mutable workflow state.

Live records use `pass`, `xfail`, `xpass`, or `fail`. After infrastructure
succeeds, the grade runs the durable semantic assertions. One or more typed
semantic failures produce XFAIL for an XFAIL target. The metric keeps all
observed semantic codes. An empty semantic set is XPASS. XPASS keeps the lane
green only so the complete lane can finish, and emits an alert with the target
and owner. XPASS is not a terminal green registry state. Before archiving the
owner, remove the source binding and matching reconciliation expectation, then
run the unchanged candidate without the binding and require PASS. Run the
active-owner join at that terminal gate. Authentication, launch, timeout,
fixture, parsing, state-read, and metric failures remain ordinary failures.

### Local live execution

Build the binary and export the resolution hooks once:

```bash
go build -o ./spacedock ./cmd/spacedock
export SPACEDOCK_BIN="$PWD/spacedock"
export SPACEDOCK_REPO_ROOT="$PWD"
```

Run the common journeys by selecting one transport. Claude and Codex run at most
three common journeys at one time. Pi runs at most two. Codex setup artifacts are
isolated under `codex-shared-scenarios/_setup/<journey-id>/`. The Claude and Pi
commands keep `-failfast`, but Go can start queued parallel journeys after a
failure. The suite timeouts remain loose runaway backstops:

```bash
SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 90m -run '^TestLiveCommon' -parallel 3 ./internal/ensigncycle -v
```

Run all three current Claude substrate proofs with one 20-minute backstop:

```bash
go test -tags live -count=1 -timeout 20m -run 'TestLiveMergedTeamModeDispatch|TestLiveBareReachable|TestLiveBreakGlassShimRecovery' ./internal/ensigncycle -v
```

For Codex, install and authenticate the CLI (or set `OPENAI_API_KEY`), then run:

```bash
SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveCommon' -parallel 3 ./internal/ensigncycle -v
```

Leave `SPACEDOCK_CODEX_LIVE_REQUIRED` unset for this local path. When no `OPENAI_API_KEY` is set, the harness copies `~/.codex/auth.json` into an isolated `CODEX_HOME`; if the variable is already set, run `unset SPACEDOCK_CODEX_LIVE_REQUIRED` first.

Run the Pi live proofs locally with the same package versions pinned in CI:

```bash
npm install -g @earendil-works/pi-coding-agent@0.80.10
npm install --prefix "$HOME/.pi/agent/npm" pi-subagents@0.35.1 pi-intercom@0.6.0
export PI_SUBAGENTS_PACKAGE_ROOT="$HOME/.pi/agent/npm/node_modules/pi-subagents"
export PI_INTERCOM_PACKAGE_ROOT="$HOME/.pi/agent/npm/node_modules/pi-intercom"
```

Authenticate with `pi login` or `OPENAI_API_KEY`. Configure the child model with the exact provider-qualified spelling for that authentication path. Custom providers (e.g. `lunaroute`) declare their models in `~/.pi/agent/models.json`, not `auth.json`; the harness mirrors `models.json` (alongside `auth.json`, both `0o600`) into the isolated Pi home so a custom-provider model resolves. If the child reports `Model ... not found`, the `models.json` mirror is missing or does not declare that provider/model. Use the `provider/model:thinking` form to request a specific thinking effort:

```bash
export SPACEDOCK_PI_LIVE_CHILD_MODEL=openrouter/openai/gpt-5.4 # OpenRouter login
# or
export SPACEDOCK_PI_LIVE_CHILD_MODEL=openai/gpt-5.4 # direct OpenAI provider
# or, with an explicit thinking effort (provider/model:thinking):
export SPACEDOCK_PI_LIVE_CHILD_MODEL='lunaroute/glm-5.2-vision-background:max'
```

A slow `:max`-thinking model can take minutes per turn. Raise the per-run cap with `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` (a positive integer, in minutes) so multi-dispatch journeys complete to a graded result instead of timing out. Make the outer `go test -timeout` longer than the per-run cap so the suite backstop does not fire before the individual cap. Defaults: smoke `10m`, common `12m`. Example:

`TestLivePiFrontDoorSmoke` is the only Pi substrate smoke. It checks the front door, child dispatch, durable output, and the boot contract. The grade artifact records both models, durations, and available costs.

```bash
go test -tags live -count=1 -timeout 15m -run TestLivePiFrontDoorSmoke ./internal/ensigncycle -v
```

Run the common Pi journeys separately from that substrate proof. Pi calls `t.Parallel()`; the workflow pins `-parallel` to `SPACEDOCK_PI_LIVE_PARALLEL` (the `live_parallel` dispatch input, default 4), and the local command keeps the same pin:

```bash
PI_LIVE_PARALLEL="${SPACEDOCK_PI_LIVE_PARALLEL:-4}"
SPACEDOCK_LIVE_RUNTIME=pi go test -tags live -count=1 -timeout 40m -run '^TestLiveCommon' -failfast -parallel "$PI_LIVE_PARALLEL" ./internal/ensigncycle -v
```

For a custom slow `:max`-thinking model, lower `PI_LIVE_PARALLEL` (e.g. `PI_LIVE_PARALLEL=2`) and raise `-timeout` above the per-run cap set by `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` (e.g. `-timeout 120m` with `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES=40`).

Without auth, the respective live suite skips locally (Claude/Codex/Pi), except in CI where the lane requires it.

### GitHub setup

#### Subscription OAuth secrets

Both the `CI-E2E-CODEX` and `CI-E2E-PI` Environments accept `CODEX_AUTH_JSON`,
the complete `~/.codex/auth.json` payload. Keep `OPENAI_API_KEY` in each
Environment during the migration; the lane uses it only when its OAuth secret
is absent.

Create or rotate the shared secret from a trusted workstation without printing
it:

    gh secret set CODEX_AUTH_JSON --env CI-E2E-CODEX < "$HOME/.codex/auth.json"
    gh secret set CODEX_AUTH_JSON --env CI-E2E-PI < "$HOME/.codex/auth.json"

For Codex the runner writes this object directly to the isolated
`CODEX_HOME/auth.json`. For Pi it safely transforms `tokens.access_token`,
`tokens.refresh_token`, `tokens.account_id`, and the access-token JWT `exp`
into Pi's `openai-codex` record (`access`, `refresh`, `accountId`, and
`expires`); the complete source payload is never passed to the Pi child.
Refreshes stay in the isolated run home and are never written back to GitHub.
Replace a revoked or expired secret from a trusted workstation. If OAuth is
absent, `OPENAI_API_KEY` is used; a lane fails before launch only when both
credentials are absent.

For OAuth Pi uses `openai-codex/gpt-5.6-luna:max`; the API-key fallback uses
`openai/gpt-5.6-luna:max`. The model ID and `max` thinking level are unchanged.

| Selected command | Unique evidence | Measured sample or cost |
|---|---|---|
| Claude `TestLiveCommon...` | The 17 registered common journeys | Journey metrics record duration, tokens, model, and available cost. |
| Claude substrate: `TestLiveMergedTeamModeDispatch`, `TestLiveBareReachable`, `TestLiveBreakGlassShimRecovery` | Merged and bare dispatch, plus break-glass recovery that preserves the selected bare/team mode and commits the worker report | Merged baseline: 127s Sonnet and 144s Opus. Cost was not available. |
| Codex resolver and `TestLiveCommon...` | Current-checkout resolution and common journeys | Both PR and release jobs consume Codex metrics. |
| Pi `TestLiveCommon...` and `TestLivePiFrontDoorSmoke` | Common journeys plus one four-part substrate proof | The detail artifacts preserve each run. |

The manual Pi lane keeps the lean surface: it does not install tmux and runs one front-door substrate smoke.

The optional journey-delta job uses the newest metrics artifact for each live producer in the run.
If one artifact is unavailable or incomplete, the job warns and skips the comment. The required test result does not change.

Workflow: `.github/workflows/runtime-live-e2e.yml`. The offline gate job (`go test ./...`, no secrets) must pass before a live lane uses an environment approval.

- Pull requests run `claude-sonnet-5` at maximum effort and `gpt-5.6-luna` at maximum effort.
- An explicit `live_cadence=opus-pre-release` dispatch runs offline plus `claude-opus-4-8` at maximum effort. It allocates no Codex or Pi runner and requests only `CI-E2E-OPUS` approval.
- An explicit `live_cadence=pi` dispatch runs the 17 common Pi journeys and the Pi front-door proof with `openai-codex/gpt-5.6-luna` for OAuth or `openai/gpt-5.6-luna` for the API-key fallback, at maximum thinking. It waits only for `CI-E2E-PI` approval and retains Pi logs, diagnostics, journey metrics, and session artifacts. Pull requests still run only Sonnet and Codex; Pi is optional and is not a merge requirement. Local Pi execution remains supported with `pi login` or an API key.

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
