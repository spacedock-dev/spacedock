# Proposal: per-harness, per-machine model and reasoning-effort preference

This proposes that `dispatch build`, which already knows `--host`, emit model and effort values
ready for that harness's spawn call, so the first officer copies them verbatim and translates
nothing. It also proposes that a workflow README declare a vendor-neutral tier and an optional effort
per stage, that a machine-local mapping resolve the tier to a concrete model per harness, and that
the `dispatch build` flags of open PR #845 override one dispatch.

Status: proposal for discussion (2026-11-11). Nothing is implemented. Repository facts were read on
`main` at `5785dbed7`.

## Terms

"Host" has two meanings, so this doc avoids the bare word: `--host` names the harness, and
"hosta"/"hostb" name machines.

- **Harness**: the agent runtime a worker is spawned in, `--host claude|codex|pi` on `spacedock dispatch build`.
- **Substrate**: the spawn tool a harness offers. Pi has two in `pi-first-officer-runtime.md`
  (`pi-subagents`, `pi-agent-teams`).
- **Machine**: the computer and account a harness runs on. Two machines can reach the same model
  through different provider routes or profiles.
- **FO**: the first officer. It is an LLM that runs `dispatch build` and copies envelope fields into
  the spawn call by following prose in `skills/first-officer/references/*-runtime.md`.
- **Envelope**: the JSON that `dispatch build` prints (`buildOutput`).
- **Tier**: a vendor-neutral label for how capable a stage's model should be, such as `deep`.

## Problem

A stage can request a model only on Claude (the `host == "claude"` branch in `build.go`). On Codex and Pi the declared model is dropped, and each
worker inherits the FO's model (`«worker-identity»` in each runtime adapter), so a cheap stage runs
on an expensive tier. No persistent place says "for stage X use model M at effort E" for each
harness. A concrete id committed in a README depends on the auth path and provider configuration in
use, so one README can be wrong for a given machine. Effort has no representation in non-test
`internal/dispatch` code or in `skills/first-officer/references`. The copy step in the FO is where
the harnesses differ, and it is done by an LLM reading prose.

## Current behaviour on main

| Harness | What a README `model:` does | Effort | Decided in |
|---|---|---|---|
| claude | stage > `stages.defaults` > null. Both values are checked against `modelEnum` (sonnet, opus, haiku, fable); either invalid gives exit 1. Passed as `Agent(model=…)`, omitted when null | No slot in the `Agent()` template in `claude-fo-dispatch.md` | `internal/dispatch/build.go`, block after `stageMeta.Model()` |
| codex | Ignored with stderr note `[build] declared model '…' ignored on host`; `model: null` (`TestBuildCodexHostIgnoresModelWithNote`); worker inherits the thread's model | `CodexMultiAgentV2Spawn.ToolArgs` emits only `task_name`, `message`, `fork_turns`; the "future overrides" case in `TestCodexMultiAgentV2SpawnInputAlwaysIsolatesFreshDispatch` pins that both fields are dropped | same block; `codex_v2_adapter.go` |
| pi | Ignored with the same note; `model: null` (`TestBuildPiHostIgnoresModelWithNote`). On null the FO stamps the parent's live model (PR #405, `pi-dispatch-model-stamping`, a null-handling fix rather than a passthrough) | No dispatch slot. `pi-subagents` takes it as an id suffix (below) | same block; `pi-first-officer-runtime.md` `«worker-identity»` |

Further facts that shape the design:

- **Pi ids depend on auth path and provider configuration.** `internal/release/live_models.txt` lists one
  model under two ids: `pi.oauth=openai-codex/gpt-6-luna:max` and `pi.api-key=openai/gpt-6-luna:max`.
  Custom providers such as `lunaroute` resolve only where `~/.pi/agent/models.json` declares them
  (`docs/runtime-live-ci.md`). One machine can therefore use either id, so a mapping keyed only by
  (harness, tier) assumes one auth path per machine.
- **The stage parser is scalar-only.** Read, not executed: `scalarMap` in `internal/status/stages.go`
  stores `""` for a mapping value, so a nested `model: {codex: …}` would fail the `modelEnum` check on
  claude with exit 1. The optional-field allowlist in `ParseStagesWithDefaultsData` is `feedback-to`,
  `agent`, `fresh`, `model`, so today's binary drops a new stage key.
- **Native machine config.** Spacedock has no per-machine config file of its own, as far as checked.
  For Codex, this repo shows `--profile` only as a launch argument passed through `spacedock codex --`
  (`docs/roadmap/0250-fo-behavioral-discipline/post-sprint-audit.md`), so it covers a whole FO session.

**Spawn-tool schemas: agent self-reports of their own tool schema in one session each on one machine
on 2026-11-11, not repository facts and not a guarantee for other installs.**

| Spawn tool | Model | Effort | Notes |
|---|---|---|---|
| Claude Code `Agent` | `model` enum: sonnet, opus, haiku, fable | separate `effort`: low, medium, high, xhigh, max | |
| Codex CLI 0.162.1 `spawn_agent` | optional `model: string` | optional `reasoning_effort: string` | Overrides require `fork_turns: "none"`; a full-history fork inherits model and effort. Listed effort values vary by model (`gpt-6-luna` lists no `ultra`). Behaviour on an unknown value is not specified and was not tried |
| omp 18.2.0 (a Pi-family harness) `task` | none | none | Exposes no `subagent` or `intercom` tool |
| stock Pi 0.73.1, no extension | none | none | Exposes only read, bash, edit, write: no spawn tool |

**`pi-subagents` 0.88.0, read from the published package source on 2026-11-11, not exercised in a live
spawn; this repo records only `subagents 0.75.0` as installed (`docs/roadmap/pi-ux/staff-review-2.md`).**

- `subagent` schema (`src/extension/schemas.js`): `model` is an optional string, "Child model
  provider/id … Suffix :off/minimal/low/medium/high/xhigh/max overrides agent thinking default". The
  separate `thinking` parameter is "watchdog.configure only … Dispatch ignores this; use model suffix."
- Machine-local settings (`docs/models.md`): `~/.pi/agent/settings.json`, or project `.pi/settings.json`
  (project wins), carries `subagents.defaultModel`, `defaultProvider`, `defaultThinking`,
  `agentOverrides.<agent>.model / .thinking`, and `agentOverridesByProvider.<provider>.<agent>`, keyed
  by the active parent session's provider. `subagents.disableThinking` serves providers that reject
  suffixed ids. Saved role-to-model profiles live under `~/.pi/agent/profiles/pi-subagents/`.
- Precedence as documented: "Per-run model overrides and `agentOverrides.<name>.model` win over
  frontmatter and the global default." The order between a per-run `model` and `agentOverrides` is
  not stated.

Consequence: the harness alone does not determine what a spawn can carry. Two Pi-family substrates
differ, and stock Pi has no spawn tool.

## Requirements

- **R1 Persisted, per stage, before dispatch.** "the problem is persisting that preference. which is what the FM tries to do" (CL).
- **R2 Model and reasoning level.** "I want to have a rule that I can setup the reasoning level before dispatch, because agent will just use it's own level" (Kent).
- **R3 Generalize the frontmatter path, allow override.** "the claude path reads from readme frontmatter. we should probably generalize it and allow override" (CL).
- **R4 Per harness.** "we don't have a good way to represent per-harness preference" (CL).
- **R5 Per machine.** "on pi for different host, it might have different profile config to the model"; "hosta … default llm gateway … hostb uses native provider routing (lunaroute/deepseek-4.1-flash)" (CL).
- **R6 An invoker exists.** "you still need an entry point (by prose or by a invoker somewhere)" (CL). The persisted preference is that entry point; the flags are the mechanism.
- **R7 Existing READMEs keep working.** `docs/dev/README.md` declares `model: opus` on two stages today.

## Core rule: harness-ready envelope

`dispatch build` resolves the value, encodes it for the harness, and emits it. The FO prose for each
harness reduces to "pass these fields verbatim". When the harness cannot carry the requested value,
`dispatch build` decides; the FO does not.

| Harness | Model encoding | Effort encoding | Requested value cannot be carried |
|---|---|---|---|
| claude | `model`, a `modelEnum` value | new `effort` field (self-report above); `claude-fo-dispatch.md` needs the slot added | exit non-zero (value outside the enum) |
| codex | `model`, an opaque string | `reasoning_effort`, an opaque string; `fork_turns: "none"` is already fixed by `CodexMultiAgentV2Spawn.ToolArgs` | Spacedock holds no per-model effort list, so an incompatible pair reaches the harness (Proposed decisions) |
| pi (`pi-subagents`) | `model`, a `provider/id` string | suffix in the id (`provider/id:high`); a separate field is ignored on dispatch | `pi-agent-teams` encoding unverified; a substrate with no model parameter, such as omp's `task` (Proposed decisions) |

## Candidate shapes

| | README commits | Machine holds | Granularity | hosta/hostb | Old binary reading new README |
|---|---|---|---|---|---|
| A | concrete id per harness (`harness-models:` key) | nothing | per stage | fails: committed id fits one auth path | key dropped; today's behaviour |
| B | today's `model:` as intent | (harness, intent) → model, effort | per intent | works | unchanged |
| C | A's ids as defaults | (harness, stage) overrides | per stage | works only where overridden | key dropped |
| D | tier plus optional `effort:` per stage; Claude names kept as tier aliases | (harness, tier) → model | model per tier, effort per stage | works | `model: deep` exits 1 on claude; `effort:` dropped |
| E | nothing new | the harness's own config (Pi: `pi-subagents` `agentOverrides`) | per agent role | works on Pi, by provider too; Codex not shown | unchanged |

```yaml
# D, README
- name: implementation
  model: deep        # tier; `opus` remains an alias for `deep`
  effort: high       # optional, vendor-neutral
```
```yaml
# D, machine-local mapping on hostb (location: Needs CL's call)
codex: {deep: gpt-6-luna, fast: gpt-6-luna}
pi:    {deep: "lunaroute/deepseek-4.1-flash"}
```

B needs `opus → {gpt-6-luna, high}` and `sonnet → {gpt-6-luna, low}`, so a Claude name stands in for
an effort level, and `model: opus` reads as "use Claude Opus". D separates the two; each tier has a
built-in Claude alias, so Claude needs no mapping.

E exists on Pi: `«worker.spawn»` passes the artifact's `agent` as the `subagent` `agent`, so
`agentOverrides.worker` applies; this may be what R5's "different profile config" refers to. Limits:
it is per role, and `piSpawnAgent` is `worker` for every stage without `agent:`; it is Pi only; and
the FO's null-model stamp sends a per-run `model` whose order against `agentOverrides` is undocumented.

## Precedence (stated once)

Model and effort resolve independently, field by field. Each takes the highest layer that sets it:

1. the `dispatch build` flags (`--model`, `--reasoning-effort`), for that dispatch only;
2. the stage's README value, and for effort the README `effort:`; the machine mapping resolves a tier
   to a model at this layer;
3. `stages.defaults`;
4. none: `model` is null (or the effort field is omitted), as today.

So `--model N` keeps the README's effort, and an effort-only flag keeps the resolved model. On
`pi-subagents` the id is built as `model:effort` after both resolve. A model/effort pair the
harness rejects is passed through (Proposed decisions).
With no mapping file at all, codex/pi behave as today (`model: null` plus the note) and claude is
unchanged. Other failures are in Proposed decisions.

## Recommendation

Shape D, resolved through a machine-local mapping, with the harness-ready envelope. Reasons:

1. It commits nothing machine-specific (R5), and it keeps effort apart from model choice (R2) without per-stage keys in a machine file.
2. Existing READMEs keep working: `model: opus` is an alias, and claude output is unchanged (R7).
3. The #845 flags become layer 1 unchanged, and the mapping is their invoker (R6).

D is preferred over B because B's example has to use model names to express effort. D is preferred
over deferring to native config (E) wherever it exists, because on Pi that config is per agent role
rather than per stage (R1), it exists on one harness, and its precedence against the FO's stamp is
undocumented. What gets worse under D: the `ParseStagesWithDefaultsData` allowlist and `modelEnum` in
`build.go` change; an old binary exits 1 on claude for a README with a new tier
name; each machine carries a file that no Spacedock installer creates; two stages with the same tier cannot get different models; and a Pi user who
already configured `agentOverrides` configures the route a second time. Under E instead, Codex and
Claude would still need D's mapping, so two mechanisms would coexist.

## What happens to #845

Survives as layer 1: the flags, their refusals (empty, with `--advance`), the `reasoning_effort`
omitempty field, and forwarding in `CodexMultiAgentV2Spawn.ToolArgs`, which has no non-test Go caller
on main, so the runtime-doc change is the operative part. Superseded: the Codex-only restriction, and suppressing the "ignored" note only for
`--model`. The branch is 151 commits behind main.

## Facts that settle former questions

- **The Pi runtime sentence is stale.** "stage-declared model overrides are passed unchanged"
  (`pi-first-officer-runtime.md`) dates from `24222955e` (2026-06-20), when every host got the declared
  model. `1af455784` (#463, 2026-07-02) made codex/pi ignore it and left that file untouched.
- **No Go code found selects the Pi substrate.** Grep finds no flag, env var or envelope field; the runtime doc
  says the launch/test harness selects it. A `--host pi` envelope carries the `pi-subagents` `agent` /
  `skill` pair (`buildOutput`). `piruntime.TeamsAction`, Spacedock's `pi-agent-teams` subset, has no
  model field.
- **Reuse today.** `«reuse.model-match»` is prose, skipped when `effective_model` is null. Claude reads
  the team `config.json`; Codex reuses via `followup_task` and records a model only on null; Pi defers
  reuse-advance (`«addressable-worker»`). No envelope field or `«worker-identity»` binding records effort.
- **Standing teammates are Claude-shaped.** `buildSpawnSpec` takes no host and emits an `Agent()` spec.
- **Config-path precedent.** Spacedock uses `XDG_DATA_HOME` (`install.sh`, `internal/cli/versions.go`)
  and `XDG_STATE_HOME` (`internal/cli/review_mode.go`), not `XDG_CONFIG_HOME`. Its non-test `SPACEDOCK_*`
  variables include `SPACEDOCK_BIN`, `SPACEDOCK_REPO_ROOT`, `SPACEDOCK_PLUGINS` and `SPACEDOCK_CHANNEL`.
- **Value-space precedent.** Spacedock validates the Claude enum (`modelEnum`) and treats other hosts' model
  spaces as theirs (the "dispatch-settable model space" comment in `build.go`; "the Pi model space is
  OWNED BY" in `pi-first-officer-runtime.md`).
- **Effort vocabularies.** `pi-subagents` 0.88.0: off, minimal, low, medium, high, xhigh, max. From
  the self-reports: Claude low through max; Codex low through max, plus `ultra` on some models. The five
  levels low, medium, high, xhigh, max are common to all three.

## Proposed decisions

Shape-deciding first.

| Decision | Proposed | Because | Changes if |
|---|---|---|---|
| Granularity | Model per tier, effort per stage (D) | R1 asks per stage; separate effort covers the cases in B's example without stage keys on the machine | Two stages need one tier but different models |
| Vocabulary | Tier in `model:`, Claude names as aliases | On claude an old binary exits 1 on an unknown `model:` value; a new key is dropped silently (`ParseStagesWithDefaultsData` allowlist) | CL prefers a new key and accepts the silent drop on old binaries |
| Pi: defer to native or map | Map, like the other harnesses | `agentOverrides` is per role, and `piSpawnAgent` is `worker` for every stage without `agent:` | Per-role granularity is acceptable on Pi |
| Effort values | Validate README `effort:` against low, medium, high, xhigh, max; pass flag values through opaque | Those five are common to all three vocabularies above | A harness drops one of the five |
| Model/effort pair the harness may reject | Pass through, no per-model list | Value-space precedent above; Codex lists vary per model | A live check shows a silent fallback on an unknown pair |
| Mapping file present, entry missing | Exit non-zero | The FO, not a person, runs `dispatch build`; a stderr note on exit 0 is unseen | CL wants a softer rollout |
| Pi substrate with no or unverified model parameter | Exit non-zero when a model or effort was requested. This needs a substrate signal `dispatch build` does not have today (next row); until one exists it cannot tell the substrates apart | Same reason; `piruntime.TeamsAction` has no model field | `pi-agent-teams` is shown to carry a model |
| Pi substrate encoding | Encode for `pi-subagents` | A `--host pi` envelope already has the `pi-subagents` shape | A substrate signal is added to `dispatch build` |
| Reuse | Envelope and `buildAdvanceOutput` carry effort; `«worker-identity»` records the emitted model and effort, flag values included; a mismatch in either forces a fresh dispatch | Effort is not recorded today, and Codex reuses through `followup_task` | — |
| Standing teammates | Out of scope; `buildSpawnSpec` keeps the four Claude values | It takes no host and emits an `Agent()` spec | Standing teammates are spawned on Codex or Pi |
| Pi runtime sentence | Delete "passed unchanged" | Stale, per the history above | — |

## Needs CL's call

No evidence favours either option here; each has the default we would implement.

1. Mapping location. Default: `${XDG_CONFIG_HOME:-~/.config}/spacedock/models.yaml`, plus a `SPACEDOCK_*` env var naming another file.
2. Mapping keyed also by provider, as `agentOverridesByProvider` is? Default: no; key by (harness, tier).
3. May the mapping override a Claude alias on one machine? Default: no.

Checks owed during implementation (not questions): `pi-subagents` in a live spawn at the installed
version; `pi-agent-teams` model parameters; a Codex per-spawn native config; behaviour on unknown
values; whether `Agent()` and `spawn_agent` apply the self-reported fields.
