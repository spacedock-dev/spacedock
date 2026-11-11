# Proposal: per-harness, per-machine model and reasoning-effort preference

This proposes that a workflow README keep declaring a portable model *intent* per stage, that a
machine-local mapping turn that intent into a concrete model and reasoning effort for each harness,
and that `dispatch build` flags (open PR #845) sit on top as a one-dispatch override.

Status: proposal for discussion (2026-11-11). Nothing is implemented. Facts below were read on
`main` at `5785dbed7`.

## Terms

"Host" has two meanings below and this doc never uses the bare word: `--host` names the harness;
"hosta"/"hostb" name machines.

- **Harness** — the agent runtime a worker is spawned in: `--host claude|codex|pi` on `spacedock dispatch build`.
- **Machine** — the computer and account a harness runs on. Two machines can reach the same model
  through different provider routes or profiles (a default LLM gateway on one, native provider routing on another).
- **FO** — the first officer, the orchestrating agent that runs `dispatch build` and spawns workers.
- **Envelope** — the JSON `dispatch build` prints (`buildOutput`); the FO copies its fields into the spawn call.

## Problem

A stage can ask for a model only on Claude: on Codex and Pi the declared model is dropped and every
worker inherits the FO's model (`«worker-identity»` in each runtime adapter), so a cheap stage runs on
an expensive tier. There
is no persistent place to say "for stage X use model M at effort E" per harness. A concrete id
committed in the README cannot be right on every machine, because the same model is reached through a
different route per machine (see the `live_models.txt` pair below). Reasoning effort has no
representation in non-test `internal/dispatch` code or in `skills/first-officer/references`.

## Current behaviour on main

| Harness | What a README `model:` does | Effort | Where it is decided |
|---|---|---|---|
| claude | stage > `stages.defaults` > null; value must be in `modelEnum` (sonnet, opus, haiku, fable), else exit 1; FO passes it as `Agent(model=…)`, omitted when null | No slot in the `Agent()` template of `claude-fo-dispatch.md` | `internal/dispatch/build.go`, block after `stageMeta.Model()` |
| codex | Ignored with stderr note `[build] declared model '…' ignored on host`; envelope `model: null`; worker inherits the thread's model | None. `CodexMultiAgentV2Spawn.ToolArgs` emits only `task_name`, `message`, `fork_turns`; the "future overrides" case of `TestCodexMultiAgentV2SpawnInputAlwaysIsolatesFreshDispatch` pins that `model`/`reasoning_effort` are dropped | same block; `codex_v2_adapter.go` |
| pi | Ignored with the same note; envelope `model: null` (`TestBuildPiHostIgnoresModelWithNote`). The FO stamps the parent's live model on null (roadmap `0223`, member `pi-dispatch-model-stamping`, PR #405 — a null-handling fix, not a passthrough) | No dispatch slot. The live test harness's child model takes it inside the id as `provider/model:thinking` (`docs/runtime-live-ci.md`; `internal/release/live_models.txt`); not checked for a `subagent(...)` spawn | same block; `pi-first-officer-runtime.md` `«worker-identity»` |

Two further facts shape the design:

- **Pi ids are machine-specific in this repo's own test config.** `internal/release/live_models.txt`
  lists the same model twice by auth path: `pi.oauth=openai-codex/gpt-6-luna:max` and
  `pi.api-key=openai/gpt-6-luna:max`. Custom providers (for example `lunaroute`) resolve only if the
  machine's `~/.pi/agent/models.json` declares them (`docs/runtime-live-ci.md`).
- **The stage parser is scalar-only.** By reading (not executed): `scalarMap` in
  `internal/status/stages.go` stores `""` for a mapping value, so a nested `model: {codex: …}` would
  reach the `modelEnum` check as `''` and fail `dispatch build --host claude` with exit 1. A new sibling
  key is dropped by `ParseStagesWithDefaultsData`, which keeps only `feedback-to`, `agent`, `fresh`,
  `model`; no other validator was checked.

Spacedock has no per-user or per-machine config file of its own as far as checked (env lookups in
`internal/`, `internal/cli/local_plugin.go`, `docs/specs/extension-install-contract.md`). The only
machine-local model config found belongs to the harnesses themselves (`~/.pi/agent/models.json`, `CODEX_HOME`).

## Requirements

- **R1 Persisted, per stage, before dispatch.** "the problem is persisting that preference. which is what the FM tries to do" (CL).
- **R2 Model and reasoning level.** Set "the model and the reasoning level of a dispatched worker by rule before dispatch, per role/stage" (Kent).
- **R3 Generalize the frontmatter path, allow override.** "the claude path reads from readme frontmatter. we should probably generalize it and allow override" (CL).
- **R4 Per harness.** "we don't have a good way to represent per-harness preference" (CL).
- **R5 Per machine.** "on pi for different host, it might have different profile config to the model"; "hosta … default llm gateway … hostb uses native provider routing (lunaroute/deepseek-4.1-flash)" (CL).
- **R6 An invoker exists.** "you still need an entry point (by prose or by a invoker somewhere)" (CL). The persisted preference is that entry point; the flags are only the mechanism.
- **R7 Existing READMEs keep working.** `docs/dev/README.md` declares `model: opus` on two stages today.

## Candidate shapes

**A. Per-harness concrete values in the README** (sibling key, because nesting under `model:` breaks today's binary):

```yaml
- name: implementation
  model: opus                       # claude, unchanged
  harness-models:
    codex: {model: gpt-6-luna, effort: high}
    pi:    {model: openai-codex/gpt-6-luna:max}   # `:max` on a spawn: unverified
```

Committed: everything. Fails R5: the `pi` id is right on one machine only. Old binaries drop the key
silently; `model: opus` is unchanged.

**B. Portable intent in the README, machine-local mapping per harness:**

```yaml
# README — unchanged
- name: implementation
  model: opus
```
```yaml
# machine-local mapping on hostb (location: open question)
codex:
  opus:   {model: gpt-6-luna, effort: high}
  sonnet: {model: gpt-6-luna, effort: low}
pi:
  opus:   {model: "lunaroute/deepseek-4.1-flash"}
```

Committed: the intent only. hosta writes `pi.opus` with its gateway route; hostb writes the lunaroute
id; the README is identical on both (R5). On claude with no entry, the intent is the model, as today
(R7). A machine with no mapping behaves exactly as today: `model: null` plus the note. The cost: the
Claude names become the cross-harness tier vocabulary unless a new key is added.

**C. Hybrid: README defaults per harness (A), machine-local override (B, keyed by harness and stage).**
Committed: a per-harness default plus the intent. Meets R5 only where the machine overrides; a
missing override silently uses a committed id that may not resolve on that machine. Two places to
read when debugging which model ran.

## Precedence (stated once)

The intent is resolved as today: stage `model:` > `stages.defaults.model` > none. The concrete value
for one dispatch is then:

1. `dispatch build --model / --reasoning-effort` (the #845 flags), for that dispatch only;
2. the machine-local entry for (harness, intent) — on claude, pending open question 3;
3. on claude only, the intent itself, still validated against `modelEnum`;
4. otherwise `model: null` with today's stderr note, and the FO's existing per-harness null handling
   (Codex thread inheritance, Pi parent-live-model stamp, Claude omits `model`).

Effort is emitted as `reasoning_effort` only when a layer above resolves it (omitted otherwise, as #845
already does). `buildAdvanceOutput.Model`, the value `«reuse.model-match»` compares, would carry the
resolved model, not the intent.

## Recommendation

Shape B. Reasons:

1. It is the only shape that commits nothing machine-specific, so it answers hosta/hostb (R5) and not only the per-harness case.
2. It needs no README change: `model: opus` on Claude is byte-identical, and a machine without a mapping runs as today.
3. The #845 flags become layer 1 of the precedence unchanged, and the mapping is their missing invoker (R6).

What gets worse: until a machine has a mapping, Codex and Pi still inherit the FO's tier, and the only
sign is the existing stderr note. Each machine carries one more file that no Spacedock installer
creates. Reusing Claude names as tiers reads oddly on Codex and Pi.

## What happens to #845

Survives as layer 1: the `--model` / `--reasoning-effort` flags, their refusals (empty, with
`--advance`), the `reasoning_effort` omitempty envelope field, and forwarding in
`CodexMultiAgentV2Spawn.ToolArgs`. On main grep finds no non-test Go caller of that adapter (only
`conflict_owner_fixture_test.go` and the contractlint test); the FO calls `spawn_agent` from
`codex-first-officer-runtime.md`, so the runtime-doc change is the operative part. Superseded: the
Codex-only restriction (`«worker-identity»` in `pi-first-officer-runtime.md` already has the FO pass
stage-declared models on Pi), and suppressing the
"ignored" note only for `--model`, which becomes "suppress when any layer resolves". The branch is 151
commits behind main and needs a rebase.

## Open questions for CL

1. Intent vocabulary: reuse `model:` values (sonnet/opus/haiku/fable) as cross-harness tiers, or add a new stage key — and if new, its name?
2. Mapping location: a fixed per-user path, an env var naming a file, or both?
3. May a mapping entry also override Claude (pin `opus` to another enum value on one machine)?
4. Mapped model and effort values on codex/pi: validated by Spacedock (per-harness enum), or opaque strings passed through?
5. Pi: does effort stay folded in the id (`:max`), or does Spacedock emit it separately?
6. When a mapping file exists but lacks the entry, stderr note (as today) or exit 1?
7. On reuse, does a model or effort mismatch force a fresh dispatch, given the flags are refused with `--advance` so a flag-overridden worker will not match the resolved value?
8. Are standing teammates in scope (`standing.go` validates their `## Hook: startup` model against `modelEnum`)?
9. On Pi, `«worker-identity»` in `pi-first-officer-runtime.md` says stage-declared models "are passed unchanged", but `TestBuildPiHostIgnoresModelWithNote` pins `model: null` — which is intended?

Unverified in this repository (needs a live check, not a ruling): whether Codex `spawn_agent` accepts
`model` and `reasoning_effort`; whether a Codex worker also inherits the FO's reasoning effort;
whether Claude's `Agent()` accepts any effort value; whether Pi's `subagent(...)` takes a thinking
level at all, in the id or as a separate parameter.
