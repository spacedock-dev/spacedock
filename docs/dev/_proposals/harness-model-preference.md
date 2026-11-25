# Proposal: per-stage model and reasoning effort through each harness's own agent roles

Reading to correct: model identifiers live in each harness's own configuration on each machine (a harness is the agent runtime: Claude Code, Codex or Pi). The workflow README names a role per stage and nothing machine- or vendor-specific. Status: proposal, nothing implemented (`main` at `5785dbed7`).

## Proposed frontmatter

```yaml
stages:
  states:
    - name: implementation
      agent: builder      # a role name; each machine's harness config gives it a model and effort
    - name: validation    # no agent: the default worker, exactly as today
```

- Why `agent:`: it already exists and already reaches two harnesses: on Claude it becomes the spawn's `subagent_type`, on Pi the spawn's `agent`. No new README key.
- A role file holds almost no Spacedock content: the dispatch prompt loads the ensign contract (the Spacedock worker instructions) on every harness (`dispatchPointerPrompt`, `firstActionBlock`); Pi adds `skills: ensign` because an override drops the skill preload.
- The same README works on hosta and hostb; reasoning effort sits next to the model, set before dispatch.

## Where the user configures a role, per harness

**Claude Code**: a subagent file, `~/.claude/agents/builder.md`.

```yaml
---
name: builder
description: Spacedock worker for build stages
model: opus
effort: high
---
```

**Pi (pi-subagents)**: an agent file, `~/.pi/agent/agents/builder.md`. The model and thinking may instead sit in each machine's `~/.pi/agent/settings.json` as `subagents.agentOverrides.builder` (it beats the file), so hosta and hostb share the file and differ only in settings. A role without `thinking` runs at `medium`, not the parent's level; a role without `model` takes `subagents.defaultModel`, else the parent's.

```yaml
---
name: builder
description: Spacedock worker for build stages
skills: ensign
model: lunaroute/deepseek-4.1-flash
thinking: high
---
```

**Codex**: Codex has no per-role sub-agent config; `spawn_agent` takes `model` and `reasoning_effort` but no role (its session records carry an `agent_role` field, always null here). Smallest native fit: a Codex profile file named after the role, `$CODEX_HOME/builder.config.toml` (the file `codex --profile builder` layers); `dispatch build` reads its two keys into the spawn fields.

```toml
model = "gpt-6-luna"
model_reasoning_effort = "high"
```

**Reasoning effort**, per harness: Codex `reasoning_effort` applied (verified); Pi `thinking:` applied as the model suffix (verified); Claude `effort:` accepted without error, but nothing records the level a worker ran at (not verifiable).

## Default and fallback

**Default.** A stage with no `agent:` behaves as today: Claude honours the stage `model:`; Codex and Pi use the model of the first officer (the LLM that dispatches workers).

**Fallback.** On every harness an unknown role or model fails at the spawn call with a tool error naming what is available, and no worker starts. The first officer then spawns once from a `fallback` object that `dispatch build` adds beside the role (the default worker's spawn fields), and tells the user one line quoting the harness's error: `builder fell back to the default worker: Unknown agent: builder`. `dispatch build` cannot see Claude or Pi roles in advance. On Codex it reads the profile file itself, so a missing file is known before spawn: it emits only the default fields and the line `role builder: no builder.config.toml; default worker`, which the first officer repeats.

The role generalises `model:`; the flags are the override. A per-run model beats a role's own config, so on Claude and Pi Spacedock sends one only for a `--model` or `--reasoning-effort` flag (PR #845, related; Pi as `provider/id:high`, so effort needs a model) or a Claude `model:` on a stage without `agent:` (with `agent:`, it feeds only the fallback). Pi's first officer stops passing its own model on role stages.

## Open questions (each with the default we would implement)

1. Use `agent:` as the role hook, or a new `role:` key? Default: `agent:`; a new key keeps the fallback away from purpose-built agents.
2. On Codex, read the role from `$CODEX_HOME/<role>.config.toml`? Default: yes, a two-key reader, no new TOML dependency.

What gets worse: a missing purpose-built agent runs as the plain worker, announced but not refused; a fresh clone prints the line on every role stage until configured; Spacedock reads one Codex file format.

## Appendix: evidence

- `internal/dispatch/build.go` on `main`: Rule 6 sets `subagent_type` from `stageMeta.Agent()`; on Pi an override replaces `piSpawnAgent` ("worker") and omits `piSpawnSkill` ("ensign"), hence `skills: ensign`. Stage `model:` is dropped on codex/pi with a stderr note (`TestBuildPiHostIgnoresModelWithNote`). `docs/dev/README.md` sets `model: opus` on two stages without `agent:`, so it is unaffected.
- The Pi stamp is passed per run: `docs/roadmap/0223-pi-dispatch-contract/dispatch-sprint-execution.md` Q2 and Q12 (a null model had resolved to `settings.json` `defaultModel`); hence it stays for role-less stages.
- Spawn tools called directly, one machine, 2026-11-11, one-word child task; results are what each harness recorded.

| Harness | Claim | Result | Record |
|---|---|---|---|
| Claude Code 2.1.322 | role `model: haiku` applies; per-invocation `model` beats it; no `model` takes the parent's | confirmed, all three | child stream `message.model` |
| Claude Code | role file needs no Spacedock content | confirmed: worker ran `Skill(skill="spacedock:ensign")` | child output |
| Claude Code | unknown role | spawn error `Agent type 'no-such-role' not found. Available agents: …`; no child | tool result |
| Claude Code | `effort: "high"` applied | accepted; level not recorded | none |
| Codex CLI 0.162.1 | `model` and `reasoning_effort` with `fork_turns: "none"` apply; without them both are inherited | confirmed, both | child rollout `turn_context` `model`, `effort` |
| Codex CLI | unknown model | spawn error ``Unknown model `not-a-real-model` for spawn_agent. Available models: …``; no child | tool result |
| Pi 1.5.2, pi-subagents 0.88.0 | role `model` and `thinking: high` apply; `agentOverrides.builder.model` beats the file | confirmed, both | run `*_meta.json` `model` |
| Pi | per-run model beats role, unsuffixed keeps role thinking; no `model` or `model: inherit` gives the parent's | confirmed, both | same |
| Pi | `skills:` reaches the worker | recorded and in catalog; reading not shown | `*_meta.json` `skills` |
| Pi | unknown role | spawn error `Unknown agent: no-such-role`, dirs consulted, `Available agents:`; no child | tool result |
| Pi | role without `thinking` | `medium` while parent ran `low` | transcript `thinkingLevel` |

Pi caveats observed: pi-subagents 0.88.0 does not load on Pi 0.73.1 (needs `@earendil-works/pi-*`), and `subagents_enable({})` had to run before the `subagent` tool existed.

Still not tested: Claude effort actually applied; the Codex profile-file reader (it does not exist); Pi 1.0.0 with subagents 0.75.0, the versions Spacedock records (`docs/roadmap/pi-ux/staff-review-2.md`); `pi-agent-teams`; anything through `dispatch build` and a first officer.
