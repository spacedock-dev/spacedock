# Spacedock

**Spacedock is a multi-agent orchestrator where nothing ships without a decision.** It lives within your existing harness — Claude Code, Codex, or Pi. It breaks work into stages and surfaces the decisions each stage needs, batched for you. Each decision arrives with evidence measured against a predefined bar for what good looks like. You approve, send back, or escalate. Or you delegate the call to an agent. Either way, the decision is recorded with its evidence and reason.

## Why

- You're the human, and agents from a dozen sessions ping you all day: a design call, a one-line approval, "ship without full coverage?", none of it on a schedule you can plan around. You stopped deciding the work and started just answering the agent.
- You're the agent, and you stall every few steps waiting on a human who isn't watching, with no clear scope and no clear bar for done, so you ask and then you wait.
- Generation got cheap. Your attention and judgment are now the bottleneck. Every task still ends in a decision, and no one owns it, so it falls on whoever's around.

## What's different

- **The agent doesn't get to judge its own work.** Review runs as a separate stage with fresh context, no access to the maker's reasoning. It pushes back on thin evidence and work that looks busy without proving its claim.
- **Every decision leaves a trail.** Each gate — the decision point at the end of a stage — carries a stage report: findings, verdicts, artifacts, anomalies. You decide on evidence, not the transcript. The record outlives the reviewer, so you can trace a bad result back to the call that caused it.
- **The bar sharpens as you use it.** Each stage declares what good means, and the agent works to that line on its own. When a standard turns out fuzzy in practice, the agent proposes an edit to the stage's written criteria for your approval.
- **Batch the work; decide as it flows back.** Queue many items at once. Agents advance each through its stages. You handle gates as they surface, not one session at a time.
- **Isolation when it matters.** Stages that touch shared state run in their own git worktree. Lighter stages run inline.

## Start with what you already built

Point Spacedock at a project you vibe-coded into spaghetti and run `/spacedock:survey`. It reads your own agent history and shows you three things: the workflow you've been running without naming it, how you've been calling work done, and the decisions still open and waiting on you.

## Where to go next

- New to Spacedock? Start with [Get started](get-started/install.md) — install, your first launch, and your first workflow.
- Want the mental model? Read [Concepts](concepts/operating-model.md) — the operating model, the stage lifecycle, and how gates work.
- Operating a workflow? See [Running workflows](running-workflows/commission.md).
- Contributing? Read [the development workflow](contributing/development-workflow.md) and the [voice & tone guide](contributing/voice-and-tone.md).

## For agents

Spacedock's docs are read by agents too — a user's first officer parsing these docs is itself an agent. The build emits an agent-readable surface: an `llms.txt` index of the docs at the site root (discoverable from each page's `<head>` via `<link rel="alternate" type="text/markdown">`) and the repo's [agent instructions](agents/index.md).

## License

Spacedock is released under the Apache License 2.0.
