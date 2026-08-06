# Architecture notes

Spacedock is a Go binary plus a set of harness skills, and the two halves divide cleanly: the binary owns state and command behavior, the skills own orchestration prose. This page maps the project shape, the design contracts under `docs/specs/`, and the runtime live CI model that proves the orchestration actually behaves. For the proof discipline these all serve, see [Proof policy](proof-policy.md).

## Project shape

The binary is small Go packages with narrow boundaries; the orchestration lives in markdown skills the host loads.

- **`cmd/spacedock/`** holds the process entry point only (`main.go`). A second entry point, `cmd/spacedock-release/`, drives release cutting.
- **`internal/cli/`** owns command routing, usage text, and exit-code behavior. The front-door verbs (`spacedock claude`, `spacedock codex`, `spacedock pi`, `spacedock doctor`) launch a host with the Spacedock plugin loaded; `init`, `new`, and `status` are the workflow-facing verbs. Path resolution and launch shape live here, not in the skills.
- **`internal/status/`** is the `status` implementation: frontmatter parse, stage enumeration, the read queries (`--next` dispatch, `--resolve`, `--short-id`, `--validate`), mutation (`--set`, `--archive`), and the guards that refuse an unsafe mutation. Output is held stable by golden fixtures under `internal/status/testdata`.
- **`docs/specs/`** holds the design contracts (see below).
- **`skills/`** holds the host-loaded orchestration skills: `commission/`, `survey/`, `debrief/`, `refit/`, `first-officer/`, `ensign/`, `present-gate/`, and `feedback-rejection-flow/`. Each is a `SKILL.md` (some with `references/` and `bin/`). Skill instructions call `spacedock status`, never a plugin-private script path. The binary owns path resolution and mutation guards, and the skills stay declarative.

Other `internal/` packages support these: `internal/contract` and `internal/contractlint` (the shipped contract and its structural lints), `internal/ensigncycle` (the runtime live scenario surface), `internal/dispatch`, `internal/safehouse` (the `.safehouse` sandbox profile), and `internal/release`.

The division is deliberate: a behavior that can be guarded by the binary or a failing test belongs in the binary, not in a sentence in a skill file.

## Design contracts under `docs/specs/`

`docs/specs/` holds the contracts downstream code cites instead of re-deriving. Two are current.

- **`state-behavior-extension.md`** defines the split-root storage profile. A development workflow keeps its README in the main repo and its mutable entities in a per-workflow `.spacedock-state` checkout, so shared issues advance without noisy state commits on the code branch. The README's `state: .spacedock-state` frontmatter field names the checkout, resolved relative to the README directory. The spec fixes the v0 layout (entities directly under `.spacedock-state`, no `entities/` directory; `_archive/` and `_debriefs/` siblings) and the mutation rules: reads compose the main README's stages with the checkout's entities, while `--set` and `--archive` write only inside the checkout.
- **`scenario-testing-principles.md`** sets out the semantic model for scenario testing. Common LLM journeys are the 16 exported `TestLiveCommon...` functions reconciled against the live registry, fixture annotations, and workflow selectors.

## Runtime live CI model

The live lanes prove runtime behavior by launching a real headless host, observing its output, and checking the resulting workflow state. A static grep over workflow YAML or skill prose is not a substitute. This is the LLM-executor side of the scenario contract above.

Each common test binds its stable ID, fixture builder, target-scoped TODO owner, runtime-neutral exercise, and durable assertion in one adjacent `liveJourney(...)` call. The helper selects only the Claude, Codex, or Pi transport; it never redispatches by journey name. Liveness remains the transport's no-progress quiet budget.

CI runs these in `.github/workflows/runtime-live-e2e.yml`. The offline gate job (`go test ./...`, no secrets) must pass before either live lane spends its environment approval:

- **`claude-live`** (matrix `sonnet` and `claude-opus-4-8`): secret `ANTHROPIC_API_KEY`. Runs the full-cycle smoke and the shared suite, loading the current checkout via `spacedock claude --plugin-dir "$GITHUB_WORKSPACE"`.
- **`codex-live`**: secret `OPENAI_API_KEY`. Builds a local marketplace under `$RUNNER_TEMP` and fails if the listing names a remote `github.com`/`ref next` install instead of the local path.
- **`pi-live`**: installs `pi-coding-agent` and runs the Pi coverage guard plus the front-door smoke.

Every live lane tests the current checkout, never a remote `--ref next` install. For the local invocation commands and the full layer-by-layer breakdown of the scenario surface, see [the development workflow](development-workflow.md).

## See also

- [Proof policy](proof-policy.md): why behavior is proven by exercising it, the instruction-file-read quarantine, and the detached adversarial audit.
- [The development workflow](development-workflow.md): the authoritative stage-by-stage rules, the entity field reference, and the live-suite commands.
- [Agent development](agent-development.md): the first-officer/ensign write-scope rules and the durable-state evidence discipline.
