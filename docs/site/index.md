# Spacedock

Spacedock runs your work as a series of stages. **Nothing crosses a gate without a decision you own.**

A gate is a checkpoint where the workflow pauses and asks you the question the work has reached: ship this, or not? You approve it, send it back, or escalate. You can also delegate the call to an agent. Either way, the decision is recorded with its evidence and its reason. That is the whole idea. Everything else is detail.

You are the captain. You set the bar and make the calls; the agents do the rest. The bar starts rough and sharpens every time you reject, so the calls that once needed you become ones you can hand off with confidence. See [the operating model](concepts/operating-model.md) for how the three roles divide the work.

## What's different

- **The agent doesn't get to judge its own work.** Review runs as a separate stage with fresh context, no access to the maker's reasoning. It pushes back on thin evidence and work that looks busy without proving its claim.
- **Every decision leaves a trail.** Each gate carries a stage report: findings, verdicts, artifacts, anomalies. You decide on evidence, not the transcript, and the record outlives the reviewer.
- **The bar sharpens as you use it.** Each stage declares what good means and the agent works to that line. When a standard turns out fuzzy in practice, the agent proposes an edit to the written criteria for your approval.
- **Batch the work; decide as it flows back.** Queue many work items at once. Agents advance each through its stages, and you handle gates as they surface, not one session at a time.
- **Work survives the context limit.** When an agent runs out of context, a successor carries forward what's in flight.

## Where to go next

- **[Get started](get-started/install.md)**: install the `spacedock` launcher and the host plugin, then make your [first launch](get-started/first-launch.md) and build your [first workflow](get-started/first-workflow.md).
- **[Concepts](concepts/operating-model.md)** covers the operating model, workflows and entities, the stage lifecycle, gates and decisions, and a worked example.
- **[Running workflows](running-workflows/commission.md)** walks through commissioning a workflow, surveying an existing project, operating a running workflow, and debriefing and refitting between sessions.
- **[Contributing](contributing/development-workflow.md)** covers the development workflow, agent development, the proof policy, and releasing.

New here? Start with [Install](get-started/install.md). It walks a fresh install end to end and names the output to expect at each step.

## For agents using Spacedock

Spacedock's docs are read by agents too. A user's first officer parsing these docs is itself an agent. The build emits a curated `llms.txt` index of the docs at the site root for product-using agents. (Repo-development guidance for an agent working ON Spacedock lives under Contributing → [Agent development](contributing/agent-development.md).)
