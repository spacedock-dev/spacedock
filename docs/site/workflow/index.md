# What is Spacedock Workflow

You hand your agents whole chunks of work, and every result comes back to you for a call. Most of those calls repeat a judgment you've already made; only a few really need you. **Workflow is where you write down once what good looks like at each step. Work that misses that bar goes back to the agent to fix before it reaches you, so what does reach you is worth your judgment, with the evidence attached.**

## How it works

Workflow runs your work as a series of stages. **Nothing crosses a gate without a decision you own.**

A gate is a checkpoint where the workflow pauses and puts the question to you: ship this, or not? You approve it, send it back, or escalate. You can also delegate the call to an agent. Either way, the decision is recorded with its evidence and its reason. That is the whole idea. Everything else is detail.

You are the captain. You set the bar and make the calls; the agents do the rest. The bar starts rough and sharpens every time you reject, so calls that once needed you become ones you can hand off with confidence. See [the operating model](../concepts/operating-model.md) for how the three roles divide the work.

## What you get

- **The agent doesn't get to judge its own work.** A separate review stage checks it with fresh context and no access to the maker's reasoning. It pushes back on thin evidence and work that looks busy without proving its claim.
- **Rework stops before it loops.** Work that fails review goes back up to three times. After the third failed round, the call comes to you instead, with every round on the record. See [rejections](../concepts/gates-and-decisions.md#rejections).
- **Every decision leaves a trail.** Each gate carries a stage report: findings, verdicts, artifacts, anomalies. You decide on evidence, not the transcript, and the record outlives the reviewer.
- **The bar sharpens as you use it.** Each stage declares what good means and the agent works to that line. When a standard turns out fuzzy in practice, the agent proposes an edit to the written criteria for your approval.
- **Batch the work; decide as it flows back.** Queue many work items at once. Agents advance each through its stages, and you handle gates as they surface, not one session at a time.
- **Work survives the context limit.** When an agent runs out of context, a successor carries forward what's in flight.

## Next

- **Start from what you already do.** [Install Workflow](../get-started/install.md), then [survey your project](../get-started/survey.md): it reads your past agent sessions and names the workflow you're already running. Or [start a fresh workflow](../get-started/first-workflow.md) from a common shape like development or research.
- **Decide its calls where you read best.** When Workflow holds a call for you, you can decide it in [Spacedock Review](../review/index.md) instead of in chat.
- **Learn the model** in [Concepts](../concepts/operating-model.md): the operating model, workflows and entities, the stage lifecycle, and gates and decisions.

## For agents using Spacedock

Agents read these docs too. Start from [`llms.txt`](/llms.txt), the curated index of these pages.
