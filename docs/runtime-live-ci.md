## Runtime Live CI

The live lanes prove runtime behavior, not text shape. Static grep checks over workflow YAML or skill prose are not a substitute for launching the real host front door, observing its output, and checking the resulting workflow state.

A runtime regression should be caught once per user journey and then exercised by EACH supported host. The shared runtime scenarios make that real: one host-neutral scenario table, per-host runner adapters (Claude and Codex today, with Pi tracked through an explicit live/codified/gap coverage map until its shared runners are live-safe) implementing or accounting for the same scenario IDs, and a parity guard that fails if a scenario exists for one host only.

### Shared runtime scenarios

The scenario surface lives in `internal/ensigncycle` and splits into four host-neutral layers plus one host-specific layer:

| Layer | File | Host-neutral? |
|-------|------|---------------|
| Scenario table | `shared_scenarios_test.go` (`sharedRuntimeScenarios()`) | Yes |
| Fixtures + prompts | `shared_fixtures_test.go` | Yes |
| Assertions | `gate_assert_impl_test.go`, `shared_assertions_impl_test.go` | Yes |
| Runner adapter | `codex_live_runner_test.go`, `claude_live_runner_test.go`, `pi_shared_coverage_test.go` | No — one per host; Pi currently records explicit live/codified/gap status for each shared scenario |

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

1. Add a `sharedRuntimeScenario` entry to `sharedRuntimeScenarios()` with a unique `name`, its old Python provenance, the behavior `intent`, and a live `timeout`. Keep it host-neutral — no launch/auth/plugin field.
2. Add a fixture writer (README + entity + any `_mods/`) and a prompt to `shared_fixtures_test.go`. The prompt must say `Use $spacedock:first-officer`; both hosts honor it. Reuse the existing fixtures verbatim where the journey is the same.
3. Add a host-neutral assertion over `(before, after, observed)` strings (or reuse an existing one) and at least one offline negative case in `shared_scenarios_negative_test.go` that builds the broken end-state and proves the assertion goes red.
4. Add a runner entry for the new `name` to BOTH `codexScenarioRunners()` and `claudeScenarioRunners()`. `TestSharedScenarioRunnerCoverage` fails until both hosts cover it.

The shared coverage meta-test enforces parity in both directions: every shared scenario must have a Claude and Codex runner plus a Pi live/codified/gap coverage entry, and every runner or Pi coverage entry must map to a defined scenario.

### Local live execution

Build the binary and export the resolution hooks once:

```bash
go build -o ./spacedock ./cmd/spacedock
export SPACEDOCK_BIN="$PWD/spacedock"
export SPACEDOCK_REPO_ROOT="$PWD"
```

Run the Claude shared suite locally (skips when no Claude auth is available — set `~/.claude/benchmark-token` for the OAuth path or `ANTHROPIC_API_KEY` for the API-key path; runs against a fresh isolated `HOME`). The `-timeout 40m` is a LOOSE BACKSTOP only — sized above the full 4-scenario serial-suite wall-time (~27m opus). The REAL liveness guard is the per-stage no-progress quiet budget (the shared `streamWatcher`, 60s) in the runners: it resets on every stream line and kills a hang at 60s of stream silence. The 40m ceiling never fires in a healthy run, it only bounds a pathological progressing-but-runaway loop and keeps the suite off Go's too-short default 10m binary timeout:

```bash
go test -tags live -count=1 -timeout 40m -run TestLiveClaudeSharedScenarios ./internal/ensigncycle -v
```

Run the Codex shared suite locally (`npm install -g @openai/codex` then `codex login`, or set `OPENAI_API_KEY`). Local runs may authenticate either through an existing Codex login at `~/.codex/auth.json` or through `OPENAI_API_KEY`. The test seeds only the minimal `features.multi_agent_v2` fragment and copies only `auth.json` for the local subscription path; it does not copy local plugin state, other credentials, or the rest of the operator's Codex config. CI does not use local subscription auth.

Each Codex shared scenario launches one `spacedock codex` front-door process, which launches one `codex exec`. A fixed 15-minute wall-clock process limit is its only scenario-level liveness guard; JSONL activity, `wait_agent` events, and durable writes do not extend the deadline, and the runner does not retry. The runner preserves JSONL, stderr, the process result, and post-run durable entity/Git evidence, then requires exit 0 and grades the existing workflow assertions. A failed keep-moving run prints and retains its native Git root; a passing run removes it. The suite-wide `-timeout 40m` remains a loose outer backstop.

```bash
go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v
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

`TestLivePiFrontDoorSmoke` remains active: it loads the current checkout's Spacedock first-officer and ensign skills plus the local pi-subagents extension/skill and verifies durable split-root worker state. `TestLivePiRecordedGateLifecycle` remains selected but emits `TODO(9w59t6m1qc46hccd54p04z2j)` while delegated gate presentation-to-application/dispatch is quarantined; a green command does not prove that capability.

```bash
go test -tags live -count=1 -timeout 15m -run '^(TestLivePiFrontDoorSmoke|TestLivePiRecordedGateLifecycle)$' ./internal/ensigncycle -v
```

The parity and definition guards run with no model spend — useful before paying for a live run:

```bash
go test -tags live -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions|TestPiSharedScenarioCoverage' ./internal/ensigncycle -v
```

Without auth, the respective live suite skips locally (Claude/Codex/Pi), except in CI where the lane requires it.

### GitHub setup

Workflow: `.github/workflows/runtime-live-e2e.yml`. The offline gate job (`go test ./...`, no secrets) must pass before either live lane burns its environment approval.

- `claude-live` (matrix: `sonnet` on `CI-E2E`, `claude-opus-4-8` on `CI-E2E-OPUS`): secret `ANTHROPIC_API_KEY`. Runs `TestLiveEnsignCycle` (the full-cycle smoke), `TestLiveClaudeSharedScenarios` (the shared suite over the headless `-p` transport), and the pty/tmux team-mode harness (`TestLivePtyStandingResidencyInjectsCommOfficer` + `TestLivePtyEnsignCycleTeamTeardown`) — which drives a real interactive session where team mode is exposed (tmux is installed for this). Artifacts under `live-artifacts/claude/<model>/` plus the session jsonl under `$CLAUDE_CONFIG_DIR`.
  For local Spacedock task `sonnet-gate-guardrail-no-authority` (`3zzpdw704df1g8pg1x9thzmw`), only the Claude Sonnet `gate-guardrail` case is temporarily non-evidence, based on run `30708727845`, job `91392375253`, artifact `8821429777` (resolved model `claude-sonnet-5`, head `57489d491`). Its narrow runner-boundary `TODO(3zzpdw704df1g8pg1x9thzmw)` skip does not disable the other Sonnet scenarios or add a skip to any Opus, Codex, or Pi case. Promote it back to evidence only after the defect is fixed and a fresh approved Sonnet live run passes the unchanged strict gate oracle; then remove the skip.
- `codex-live` (environment `CI-E2E-CODEX`): secret `OPENAI_API_KEY`, `SPACEDOCK_CODEX_LIVE_REQUIRED=1` so a missing key fails clearly after approval. Runs `TestLiveCodexSharedScenarios`. Artifacts under `live-artifacts/codex/`.
- `pi-live` (environment `CI-E2E-PI`): secret `OPENAI_API_KEY`, `SPACEDOCK_PI_LIVE_REQUIRED=1` so missing Pi/OpenAI prerequisites fail clearly after approval. Installs `pi-coding-agent`, `pi-subagents`, and `pi-intercom`, runs the Pi shared coverage guard and active `TestLivePiFrontDoorSmoke`, and keeps `TestLivePiRecordedGateLifecycle` selected so its `TODO(9w59t6m1qc46hccd54p04z2j)` quarantine is visible. A green job does not prove delegated gate continuation while that skip remains. Artifacts upload under `live-artifacts/pi/`.

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
