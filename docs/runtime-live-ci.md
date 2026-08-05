## Runtime Live CI

The live lanes prove runtime behavior, not text shape. Static grep checks over workflow YAML or skill prose are not a substitute for launching the real host front door, observing its output, and checking the resulting workflow state.

A runtime regression should be caught once per user journey and then exercised by EACH supported host. The shared runtime scenarios make that real: one host-neutral scenario table, one common runner map, and Claude, Codex, and Pi transport adapters implementing or accounting for the same scenario IDs. A parity guard fails if a desired scenario is absent from any shared surface.

### Shared runtime scenarios

The scenario surface lives in `internal/ensigncycle` and splits into four host-neutral layers plus one host-specific layer:

| Layer | File | Host-neutral? |
|-------|------|---------------|
| Scenario table | `shared_scenarios_test.go` (`sharedRuntimeScenarios()`) | Yes |
| Fixtures + prompts | `shared_fixtures_test.go` | Yes |
| Assertions | `gate_assert_impl_test.go`, `shared_assertions_impl_test.go` | Yes |
| Runner adapter | `codex_live_runner_test.go`, `claude_live_runner_test.go`, `pi_shared_live_runner_test.go` | No — one transport adapter per host |

Each runner adapter turns a shared scenario into a real launch and returns `(before, after, observed)` for the shared assertions:

| Concern | Codex runner | Claude runner |
|---------|--------------|---------------|
| Auth / HOME isolation | isolated `CODEX_HOME` + minimal `config.toml` plus copied `auth.json` / `OPENAI_API_KEY` | clean `HOME` + OAuth benchmark-token / `ANTHROPIC_API_KEY` (`isolatedClaudeEnv`) |
| Plugin install | `spacedock codex --plugin-dir <checkout>` consumes the checkout before `--` | `spacedock claude --plugin-dir <checkout> --skip-compat-check` |
| Launch | `spacedock codex <task> -- exec --json --enable multi_agent_v2 --output-last-message <file>` | `spacedock claude -- -p <prompt> --output-format stream-json` |
| `observed` extract | durable workflow state; final message only where the scenario promises user-facing text | durable workflow state; final message only where the scenario promises user-facing text |
| Artifacts | jsonl / final-message / stderr | stream jsonl / final-message |

The shared scenarios reuse the old shared Claude/Codex Python journey overlap (`tests/test_gate_guardrail.py`, `tests/test_rejection_flow.py`, `tests/test_merge_hook_guardrail.py`):

- `gate-guardrail`: starts at a human gate and asserts the first officer presents the gate instead of self-approving, mutating, or archiving the entity.
- `rejection-flow`: drives a two-cycle rejection trajectory — route the concrete finding back through implementation, re-implement, and re-validate a second cycle, reusing the kept-alive reviewer when the host exposes an addressable-worker route and otherwise fresh-dispatching a separate reviewer.
- `feedback-3-cycle-escalation`: starts from two prior rejection cycles at a third REJECTED validation and asserts the first officer escalates to the human on the third cycle instead of auto-bouncing a fourth time.
- `merge-hook-guardrail`: attempts terminalization while a merge hook is registered and asserts the guard refuses bypass without `mod-block`, PR, or force.

Assertions prefer durable workflow state over transcript phrasing: entity frontmatter (status / completed / verdict), archive-vs-no-archive, the exact fix marker and a second stage report, and only the durable user-facing final-message obligations (a gate review and a decision prompt). `extractClaudeFinalMessage` surfaces a stale-credential `is_error`/`401` `result` event as a LOUD launch failure, distinct from a scenario-assertion failure, so a credential problem is never misread as a runtime regression.

Keep-moving completion is provider-independent. Each expected task must have its own
worker proof: either a dispatch with `started` followed by a later entity-file-only Stage
Report, or one child adding `started` and the new report over a parent with neither. A dispatch
or atomic worker child may add only gate-room files newly bound below that revision's exact-slug
`room-ref`; the first terminal signal after the worker report must already include terminal
status, `completed`, and `verdict`. That complete transition may be a separate commit or the
final entity-owned archive.
All expected tasks must engage before any one terminalizes, and that final archive must retain
terminal fields. A questioned hold requires a stage transition and ticket-file-only Stage Report
without any historical terminal status or field.
Dispatch commits may start two or more expected tickets, but credit only those valid starts;
omitted tickets still need their own dispatch. A split batch may carry `questioned` only through
its nonterminal review-to-ideation start. Terminal batches require the complete expected set;
neither batch form admits a foreign ticket, and reports and archives remain individually attributed.
When a report persists a previously set start after another ticket terminalizes, canonical
timestamps count only if two earlier ticket engagements corroborate a frontier no earlier than
every expected start, and every expected start strictly predates the earliest completion.
Same-slug sidecars are allowed only there or at a corrected-held boundary; foreign slugs reject.
Transcript JSONL, command text, provider events, and model narration remain diagnostic only; the commissioned-task fallback uses the same durable oracle.

Codex deliberately follows the Spacedock front door before handing off host arguments. The
runner passes `--plugin-dir` and `--skip-compat-check` before `--`; the front door installs the
current checkout and appends the fixed first-officer bootstrap to the task. After `--`, the
runner passes Codex's `exec` command, explicitly enables `multi_agent_v2`, and uses
`--dangerously-bypass-approvals-and-sandbox` to match Claude's live `bypassPermissions` posture.
The isolated home receives only the three-key `features.multi_agent_v2` fragment and, for local
OAuth, `auth.json`; it never copies the operator's full config, plugin cache, or other
credentials. Codex has no Claude `--agent` flag or equivalent stream result event, so its
bootstrap is positional and its final message comes from `--output-last-message`. CI pins the
Codex model through the existing non-recursive `codex exec` shim; the runner does not duplicate
that pin.

**To add a shared runtime scenario:**

1. Add the desired journey and fixture IDs to [`runtime-live-ci-registry.md`](runtime-live-ci-registry.md), the normative desired-state registry.
2. Add a `sharedRuntimeScenario` entry to `sharedRuntimeScenarios()` with a unique `name`, its old Python provenance, and behavior `intent`. Keep it host-neutral.
3. Add a fixture writer (README + entity + any `_mods/`) and a prompt. Bind the record and builder with `//spacedock:live-journey` and `//spacedock:live-fixture`.
4. Add a host-neutral assertion over `(before, after, observed)` strings (or reuse an existing one) and at least one offline negative case that builds the broken end-state and proves the assertion goes red.
4. Add one entry to `sharedScenarioRunners()`. Runtime launch, authentication, output parsing, artifacts, and liveness remain behind the Claude, Codex, and Pi adapters.

The shared coverage test enforces parity in both directions across exactly 16 registry records and the sole common runner map.

### Local live execution

Build the binary and export the resolution hooks once:

```bash
go build -o ./spacedock ./cmd/spacedock
export SPACEDOCK_BIN="$PWD/spacedock"
export SPACEDOCK_REPO_ROOT="$PWD"
```

Run the same 16-journey selector on each runtime. Only `SPACEDOCK_LIVE_RUNTIME` changes:

```bash
SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 40m -run '^TestLiveSharedScenarios$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveSharedScenarios$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=pi SPACEDOCK_PI_LIVE_REQUIRED=1 go test -tags live -count=1 -timeout 40m -run '^TestLiveSharedScenarios$' ./internal/ensigncycle -v
```

`TestLiveSharedScenarios` owns one subtest identity for each common journey. Each lane selects a runtime adapter through `SPACEDOCK_LIVE_RUNTIME`. Every adapter preserves its existing launch, liveness, artifact, and metric contract. Pi gives each journey a 12-minute process deadline, stops after the first failure, and never retries automatically.

Run all three current Claude substrate proofs with one 20-minute backstop:

```bash
go test -tags live -count=1 -timeout 20m -run 'TestLiveMergedTeamModeDispatch|TestLiveBareReachable|TestLiveBreakGlassShimRecovery' ./internal/ensigncycle -v
```

Run the Codex shared suite locally (`npm install -g @openai/codex` then `codex login`, or set `OPENAI_API_KEY`). Local runs may authenticate either through an existing Codex login at `~/.codex/auth.json` or through `OPENAI_API_KEY`. The test seeds only the minimal `features.multi_agent_v2` fragment and copies only `auth.json` for the local subscription path; it does not copy local plugin state, other credentials, or the rest of the operator's Codex config. CI does not use local subscription auth.

Each Codex shared scenario launches one `spacedock codex` front-door process, which launches one `codex exec`. The shared stream watcher applies a 60-second quiet budget to each Codex scenario. Each complete JSONL line resets the budget.

On stream silence, the runner kills the process and reports the last event and artifact directory. It preserves JSONL, stderr, the process result, and post-run durable entity/Git evidence, and it does not retry. A failed keep-moving run prints and retains its native Git root; a passing run removes it. The suite-wide `-timeout 40m` remains the runaway backstop.

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

`TestLivePiFrontDoorSmoke` remains the Pi-specific substrate smoke. The common Pi lane is the complete 16-journey selector above, not a coverage row or Pi-only substitute. Its artifacts retain root output, root session JSONL, model, process status, duration, and durable workflow state. Each shared journey derives its `runtime=pi`, `host=pi`, provider/model, token usage, and cost record from that archived Pi root session; stdout and stderr are diagnostic output, not Pi tool or usage evidence. The planning estimate is $2.46, with a $3.08 approval ceiling; recorded actual cost is evidence, not a pass criterion.

```bash
go test -tags live -count=1 -timeout 15m -run TestLivePiFrontDoorSmoke ./internal/ensigncycle -v
```

The parity and definition guards run with no model spend — useful before paying for a live run:

```bash
go test -tags live -run 'TestSharedScenarioRunnerCoverageFinal|TestSharedRuntimeScenarioDefinitions|TestSharedLiveRuntimeSelection|TestPromotedCommonJourneyEntrypoints' ./internal/ensigncycle -v
```

Without auth, the respective live suite skips locally (Claude/Codex/Pi), except in CI where the lane requires it.

### GitHub setup

| Selected command | Unique evidence | Measured sample or cost |
|---|---|---|
| Claude `TestLiveSharedScenarios` | All 16 common journeys through the Claude adapter | Journey metrics record duration, tokens, model, and available cost. |
| Claude substrate: `TestLiveMergedTeamModeDispatch`, `TestLiveBareReachable`, `TestLiveBreakGlassShimRecovery` | Merged, bare, and break-glass dispatch | Merged baseline: 127s Sonnet and 144s Opus. Cost was not available. |
| Codex resolver and `TestLiveSharedScenarios` | Current-checkout resolution and all 16 common journeys | PR and release jobs consume Codex metrics. |
| Pi `TestLiveSharedScenarios` and `TestLivePiFrontDoorSmoke` | All 16 common journeys plus the four-part Pi substrate proof | The lane uploads per-journey and root/child evidence. |

The deletion removes 17 seconds of tmux setup and avoids a 172.5-second duplicate Pi smoke per run.

Workflow: `.github/workflows/runtime-live-e2e.yml`. The offline gate job (`go test ./...`, no secrets) must pass before either live lane burns its environment approval.

- `claude-live` runs the common suite plus merged, bare, and break-glass substrate proofs. Its matrix uses `sonnet` and `claude-opus-4-8`.
- `codex-live` runs the resolver and shared proofs. The PR delta and release ledger consume its metrics.
- `pi-live` runs the complete common suite and one front-door smoke. It uploads the grade, root session, child session, per-journey metrics, and diagnostics.

### Registry reconciliation

Run reconciliation whenever `internal/ensigncycle/`, `internal/livescenario/`, or `.github/workflows/runtime-live-e2e.yml` changes. Compare the desired registry with adjacent source annotations, the 16-record scenario table, the sole 16-entry runner map, and the three exact workflow selectors.

The current desired inventory is: 16 bound journeys, 21 bound fixtures (17 common, three runtime-proof fixtures beyond the shared lifecycle fixture, and one experiment), four runtime proofs, one suite, three adapters, and four required runtime targets selected through three lanes. Structural diagnostics are `MISSING=0`, `UNSELECTED=0`, `DUPLICATE=0`, `INVALID=0`, `ORPHAN=0`, `UNACCOUNTED-TEST=0`, and `UNACCOUNTED-BUILDER=0`. Observed source bindings separately report eight target-scoped `MISSING-EVIDENCE` results: Claude Sonnet owns `default-headless-gate-stop` through repair `26nk8qd48zknqnn4kc123sez`; Claude Sonnet and Codex each own `smallest-sufficient-mechanism` and `keep-moving-posture` through repair `9adv48yhye5s2vkhwd7ge52d`; Pi owns `rejection-flow` through repair `zbcj98qfwtax61vxdzrf615e`, because the audited run recorded the rejection round at two entries before the complete four-entry log existed; Codex owns `withdrawn-gate-recovery` through repair `47gnqfm1ft6f2hcahz98m2jv`, because hosted run `31029501075` passed the full-ensign, gate-guardrail, and default-headless checks before the withdrawn-gate recovery check exhausted three attempts; and Codex separately owns `rejection-flow` through repair `zbcj98qfwtax61vxdzrf615e`, because hosted run `31032033236` passed full-ensign-cycle, gate-guardrail, default-headless, and recorded-gate-lifecycle, skipped the owned withdrawn-gate case, then prepared an ordinary gate and stopped instead of routing rejected feedback through rework and cycle-2 validation. Codex has passing evidence for `default-headless-gate-stop`. `auto-continue-after-implementation` has no current owner-failure binding. Pi and Claude Opus remain required and runnable wherever evidence is unverified; an unverified target is missing run evidence, not an exception or an inferred owner failure. TODO bindings do not count as passing evidence. Any other nonzero structural result or changed target/owner binding lists exact journey, target, and owner and fails reconciliation.

Registry reconciliation SHA: `79d4055e8fae70c4e2f74c9edabc3f5f1ab45544`

The SHA guard compares the recorded commit with the watched paths. A stale base must fail and name every changed watched path; the recorded base must pass. Source bindings are the semantic inventory, while the path guard only proves that inventory has not gone stale.

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
