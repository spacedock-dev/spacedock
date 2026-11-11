# Proposal: per-stage model and reasoning effort through each harness's own agent roles

Reading to correct: model identifiers live in each harness's own configuration on each machine
(a harness is the agent runtime: Claude Code, Codex or Pi). The workflow README names a role per
stage and nothing machine- or vendor-specific. Status: proposal, nothing implemented (`main` at `5785dbed7`).

## Proposed frontmatter

```yaml
stages:
  states:
    - name: ideation
      agent: planner      # a role name; each machine's harness config gives it a model and effort
    - name: implementation
      agent: builder
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

**Pi (pi-subagents)**: an agent file, `~/.pi/agent/agents/builder.md`. The model and thinking may
instead sit in each machine's `~/.pi/agent/settings.json` as `subagents.agentOverrides.builder`, so
hosta and hostb share the file and differ only in settings. A role setting no model needs `model: inherit`.

```yaml
---
name: builder
description: Spacedock worker for build stages
skills: ensign
model: lunaroute/deepseek-4.1-flash
thinking: high
---
```

**Codex**: Codex has no per-role sub-agent config; its `spawn_agent` takes `model` and
`reasoning_effort` but no role. Smallest native fit: a Codex profile file named after the role,
`$CODEX_HOME/builder.config.toml` (the same file `codex --profile builder` layers). `dispatch build`
reads its two keys into the spawn fields.

```toml
model = "gpt-6-luna"
model_reasoning_effort = "high"
```

## Default and fallback

**Default.** A stage with no `agent:` behaves as today: Claude honours the stage `model:`; Codex and
Pi use the model of the first officer (the LLM that dispatches workers).

**Fallback.** When a role does not resolve here (the harness refuses the spawn before the worker starts,
for an unknown agent or model, or Codex has no profile file), the default worker runs and the first officer
copies one line from the `dispatch build` output to the user (`«reuse.model-match»` has such a line in prose): `role builder not configured for pi on this machine; ran the default worker. Add ~/.pi/agent/agents/builder.md`.

The role generalises `model:`; the flags are the override. A per-run model beats a role's own config,
so on Claude and Pi Spacedock sends one only for a `--model` or `--reasoning-effort` flag (PR #845, related;
Pi as `provider/id:high`, so effort needs a model) or a Claude `model:` on a stage without `agent:`
(with `agent:`, it feeds only the fallback). Pi's first officer stops passing its own model on role stages.

## Open questions (each with the default we would implement)

1. Use `agent:` as the role hook, or a new `role:` key? Default: `agent:`; a new key keeps the fallback away from purpose-built agents.
2. On Codex, read the role from `$CODEX_HOME/<role>.config.toml`? Default: yes, a two-key reader, no new TOML dependency.

What gets worse: a missing purpose-built agent runs as the plain worker, announced but not refused; a
fresh clone prints the line on every role stage until configured; Spacedock reads one Codex file format.

## Appendix: evidence

- `internal/dispatch/build.go` on `main`: Rule 6 sets `subagent_type` from `stageMeta.Agent()`; on Pi an override replaces `piSpawnAgent` ("worker") and omits `piSpawnSkill` ("ensign"), hence `skills: ensign` in the Pi example. Stage `model:` is checked against `modelEnum` and dropped on codex/pi with a stderr note (`TestBuildCodexHostIgnoresModelWithNote`, `TestBuildPiHostIgnoresModelWithNote`). `CodexMultiAgentV2Spawn.ToolArgs` emits no model or effort today. `docs/dev/README.md` declares `model: opus` on two stages without `agent:`, so it is unaffected.
- Pi stamp: `«worker-identity»` in `pi-first-officer-runtime.md` says "stamp the parent's live model". The roadmap (`docs/roadmap/0223-pi-dispatch-contract/`, Q2 and Q12 in `dispatch-sprint-execution.md`) shows it is passed per run because a null model resolved to `settings.json` `defaultModel`; hence the stamp stays for role-less stages.
- Claude Code docs, "Create custom subagents": frontmatter `model` and `effort` (low through max); model order is per-invocation parameter, then definition, then `CLAUDE_CODE_SUBAGENT_MODEL`, then main conversation.
- pi-subagents 0.88.0 published source, `docs/models.md`: "per-run override → provider-scoped role override → `agentOverrides.<name>.model` → agent frontmatter `model` → `subagents.defaultModel` → the parent session model"; an unknown name fails with `Unknown agent '<name>'`. Spacedock records only 0.75.0 as installed.
- Codex CLI 0.162.1: `--profile` "Layer $CODEX_HOME/<name>.config.toml on top of the base user config" (`codex --help`); `spawn_agent` parameters from the tool's own schema, one machine.
- Not tested: that any harness applies these values to a live worker, accepts a role file with an empty body, or how Claude Code answers an unknown `subagent_type`.
