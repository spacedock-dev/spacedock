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
| pi | Ignored with the same note; `model: null` (`TestBuildPiHostIgnoresModelWithNote`). On null the FO stamps the parent's live model (PR #405, `pi-dispatch-model-stamping`, a null-handling fix rather than a passthrough) | No dispatch slot. The live test child model carries effort in the id as `provider/model:thinking` (`docs/runtime-live-ci.md`); not checked for a `subagent(...)` spawn | same block; `pi-first-officer-runtime.md` `«worker-identity»` |

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
- **Native machine config.** Spacedock has no per-user or per-machine config file of its own, as far as
  checked (env lookups in `internal/`, `internal/cli/local_plugin.go`,
  `docs/specs/extension-install-contract.md`). For Pi, `~/.pi/agent/models.json` declares provider
  models. For Codex, this repo shows `--profile` only as a launch argument passed through
  `spacedock codex --` (`docs/roadmap/0250-fo-behavioral-discipline/post-sprint-audit.md`), which
  applies to the whole FO session and not to one spawn.

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
| codex | `model`, an opaque string | `reasoning_effort`, an opaque string; `fork_turns: "none"` is already fixed by `CodexMultiAgentV2Spawn.ToolArgs` | Spacedock holds no per-model effort list, so an incompatible pair reaches the harness (open question 8) |
| pi (`pi-subagents`) | `model`, a `provider/id` string | suffix in the id (`provider/id:high`); a separate field is ignored on dispatch | `pi-agent-teams` encoding unverified; a substrate with no model parameter, such as omp's `task` (open question 7) |

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
# D, machine-local mapping on hostb (location: open question 4)
codex: {deep: gpt-6-luna, fast: gpt-6-luna}
pi:    {deep: "lunaroute/deepseek-4.1-flash"}
```

B needs `opus → {gpt-6-luna, high}` and `sonnet → {gpt-6-luna, low}`, so a Claude name stands in for
an effort level, and `model: opus` reads as "use Claude Opus". D separates the two; each tier has a
built-in Claude alias, so Claude needs no mapping.

E on Pi exists today. `«worker.spawn»` in `pi-first-officer-runtime.md` passes the artifact's `agent`
field as the `subagent` `agent` parameter, so `agentOverrides.worker` would apply. This may be what R5's
"different profile config to the model" refers to. Limits: it is keyed by agent role, and
`dispatch build` emits `agent: worker` (`piSpawnAgent`) for every stage without an `agent:` override,
so stages are not distinguished. It is Pi only; this repo shows no Codex per-spawn equivalent. And the
FO's null-model stamp (`«worker-identity»`) passes an explicit per-run `model`, whose order against
`agentOverrides` is not documented, so deferring to it means dropping that stamp.

## Precedence (stated once)

Model and effort resolve independently, field by field. Each takes the highest layer that sets it:

1. the `dispatch build` flags (`--model`, `--reasoning-effort`), for that dispatch only;
2. the stage's README value, and for effort the README `effort:`; the machine mapping resolves a tier
   to a model at this layer;
3. `stages.defaults`;
4. none: `model` is null (or the effort field is omitted), as today.

So `--model N` keeps the README's effort, and an effort-only flag keeps the resolved model. On
`pi-subagents` the id is built as `model:effort` after both resolve. A model/effort pair the
harness rejects is open question 8.

Reuse compares model and effort together. Today `buildAdvanceOutput` carries `Model` only. Under D
on Codex, two stages with one model and different efforts would look reusable, while on Pi, with
effort folded into the id, they would not. `buildAdvanceOutput` therefore gains the effort field.

## Failure defaults

The FO, not a person, runs `dispatch build`. A stderr note on exit 0 is effectively unseen.
Proposed:

- No mapping file at all: today's behaviour (codex/pi: `model: null` plus the note; claude unchanged).
- Mapping file present, tier entry missing: exit non-zero.
- Requested value the harness cannot carry: exit non-zero (core rule table).

## Recommendation

Shape D, resolved through a machine-local mapping, with the harness-ready envelope. Reasons:

1. It commits nothing machine-specific (R5), and it keeps effort apart from model choice (R2) without per-stage keys in a machine file.
2. Existing READMEs keep working: `model: opus` is an alias, and claude output is unchanged (R7).
3. The #845 flags become layer 1 unchanged, and the mapping is their invoker (R6).

D is preferred over B because B's example has to use model names to express effort. D is preferred
over deferring to native config (E) wherever it exists, because on Pi that config is per agent role
rather than per stage (R1), it exists on one harness, and its precedence against the FO's stamp is
undocumented. What gets worse under D: the `ParseStagesWithDefaultsData` allowlist and `modelEnum` (in
`build.go` and `standing.go`) change; an old binary exits 1 on claude for a README with a new tier
name; each machine carries a file that no Spacedock installer creates; two stages with the same tier cannot get different models; and a Pi user who
already configured `agentOverrides` configures the route a second time. Under E instead, Codex and
Claude would still need D's mapping, so two mechanisms would coexist.

## What happens to #845

Survives as layer 1: the `--model` / `--reasoning-effort` flags, their refusals (empty, with
`--advance`), the `reasoning_effort` omitempty field, and forwarding in
`CodexMultiAgentV2Spawn.ToolArgs`. Grep finds no non-test Go caller of that adapter on main (only
`conflict_owner_fixture_test.go` and the contractlint test), so its runtime-doc change is the
operative part. Superseded: the Codex-only restriction, and suppressing the "ignored" note only for
`--model`. The branch is 151 commits behind main.

## Open questions for CL

1. Granularity: is the preference per tier (B, D), or per stage (A, C)?
2. Vocabulary: a vendor-neutral tier in `model:` with Claude names as aliases, a new stage key, or Claude names as the tiers?
3. On Pi, should Spacedock emit no model and defer to `pi-subagents` settings (dropping the null-model stamp), or resolve through the same mapping as other harnesses?
4. If a new file: a fixed per-user path, an env var naming a file, or both?
5. Should mapping keys also include the auth path or provider, as `agentOverridesByProvider` does, since one machine can use both `pi.oauth` and `pi.api-key` ids?
6. Effort values: a portable enum Spacedock validates, or an opaque per-harness string?
7. For a Pi substrate with no model parameter (omp's `task`) or an unverified one (`pi-agent-teams`): omit with a note, or exit non-zero?
8. Model/effort pair the harness may reject (such as `gpt-6-luna` with `ultra`): pass through, or validate against a list Spacedock maintains?
9. Missing mapping entry when the file exists: exit non-zero as proposed, or the stderr note?
10. May a mapping entry override Claude's built-in alias on one machine?
11. On reuse, does an effort mismatch force a fresh dispatch, and does a flag-set worker (flags are refused with `--advance`) record its values for comparison?
12. Are standing teammates in scope (`buildSpawnSpec` in `standing.go` validates their model against `modelEnum`)?
13. `«worker-identity»` in `pi-first-officer-runtime.md` says stage-declared models "are passed unchanged", but `TestBuildPiHostIgnoresModelWithNote` pins `model: null`. Which one is intended?
14. How does `dispatch build` learn the substrate, given that `--host pi` alone does not determine it?

Unverified, needs a live check: `pi-subagents` behaviour at the installed version in a live spawn;
`pi-agent-teams` parameters; a Codex per-spawn native config; harness behaviour on an unknown model
or effort value.
