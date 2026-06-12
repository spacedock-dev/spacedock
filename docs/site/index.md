# Spacedock

You're the human, and agents from a dozen sessions ping you all day — a design call, a one-line approval, "ship without full coverage?" — none of it on a schedule you can plan. Generation got cheap; your judgment is now the bottleneck, and every task still ends in a decision someone has to own.

**Spacedock is a multi-agent orchestrator where nothing ships without a decision.** It lives within your existing harness — Claude Code, Codex, or Pi. It breaks work into stages and surfaces the decisions each stage needs, batched for you. Each decision arrives with evidence measured against a predefined bar for what good looks like. You approve, send back, or escalate — or delegate the call to an agent. Either way, the decision is recorded with its evidence and reason.

You'll meet two terms right away: a **workflow** is the directory of work that defines the stages and the gates, and a **gate** is the point where the workflow pauses for your decision instead of advancing on its own. Entities (the individual work items) and stages are introduced in [Concepts](concepts/operating-model.md), where you first need them.

Three roles run a workflow:

| Role | Who |
|------|-----|
| **Captain** | You. You define the mission and make the calls at gates unless you delegate them. |
| **First Officer** | The orchestrator agent that runs the workflow and reports to you at gates. |
| **Ensign** | The worker agent that moves one entity through one stage. |

## What's different

- **The agent doesn't get to judge its own work.** Review runs as a separate stage with fresh context, no access to the maker's reasoning. It pushes back on thin evidence and work that looks busy without proving its claim.
- **Every decision leaves a trail.** Each gate carries a stage report: findings, verdicts, artifacts, anomalies. You decide on evidence, not the transcript, and the record outlives the reviewer.
- **The bar sharpens as you use it.** Each stage declares what good means and the agent works to that line. When a standard turns out fuzzy in practice, the agent proposes an edit to the written criteria for your approval.
- **Batch the work; decide as it flows back.** Queue many entities at once. Agents advance each through its stages, and you handle gates as they surface, not one session at a time.
- **Work survives the context limit.** When an agent runs out of context, a successor carries forward what's in flight.

## Where to go next

- **[Get started](get-started/install.md)**: install the `spacedock` launcher and the host plugin, then make your [first launch](get-started/first-launch.md) and build your [first workflow](get-started/first-workflow.md).
- **[Concepts](concepts/operating-model.md)** covers the operating model, workflows and entities, the stage lifecycle, gates and decisions, and a worked example.
- **[Running workflows](running-workflows/commission.md)** walks through commissioning a workflow, surveying an existing project, operating a running workflow, and debriefing and refitting between sessions.
- **[Contributing](contributing/development-workflow.md)** covers the development workflow, agent development, the proof policy, and releasing.

New here? Start with [Install](get-started/install.md). It walks a fresh install end to end and names the output to expect at each step.

## For agents using Spacedock

Spacedock's docs are read by agents too. A user's first officer parsing these docs is itself an agent. The build emits a curated `llms.txt` index of the docs at the site root for product-using agents. (Repo-development guidance for an agent working ON Spacedock lives under Contributing → [Agent development](contributing/agent-development.md).)
