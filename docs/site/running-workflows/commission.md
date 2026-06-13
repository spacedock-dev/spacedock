# Commission a workflow

A workflow is designed in conversation: you make four decisions, commission derives the rest, and nothing is written until you accept the presented design. [Your first workflow](../get-started/first-workflow.md) shows the whole flow, including the design summary and pilot run.

Invoke it from a session started with `spacedock claude`. You can pass the mission inline:

```
/spacedock:commission product idea to simulated customer interview
```

Text after the command name becomes the workflow mission; with no argument, commission greets you and asks.

## The four things you decide

1. **The mission and what each work item is.** What the workflow is for, and what one entity represents. The description you give becomes the label the workflow uses everywhere it talks about your work: "a design idea" makes it a workflow of ideas.

2. **The stages.** Commission detects the workflow's shape from your mission (shipping code, testing a hypothesis, or iterating on an artifact) and proposes a stage list for you to confirm, modify, add to, or trim. It pushes back on redundant names: `awaiting_validation` reads as "the entity is in awaiting_validation," so it suggests `validation` instead.

3. **The gated stages.** Which stages pause for your decision. By default one gate sits before the terminal stage, and each gate gets a rejection target: the earlier stage work bounces back to when you reject. Both appear in the design summary in plain language ("If you reject at `review`, it goes back to `draft` for revision").

4. **The per-stage quality bar.** Each stage in the generated README tells a dispatched ensign what "good" means there: what to produce, the quality bar to meet, the anti-patterns to avoid. Commission drafts these from the mission, but they are starting prose, not commitments.

Everything else is derived or asked with a recommendation attached: the directory under `docs/` where everything lands as plain text, how entities are identified, how rejections route. Stages that write code give each entity an isolated worktree, so your main checkout stays clean; if you ship through PR review, commission offers the [pr-merge mod](../advanced/mods-and-standing-teammates.md), which manages the PR lifecycle so merging never needs to be a stage. The design summary shows it all; ask about the tradeoffs before you accept.

## Two ways to tighten the README

The generated per-stage rules are best-guesses; [tighten them before any work runs](../concepts/workflows-and-entities.md#the-readme-is-where-you-set-the-rules). Either:

- open the README and edit the bullets under each stage heading directly, or
- type `review stages` to have commission walk you through each stage, flag the rules that read as generic, and apply your amendments inline.

## The pilot run

On accept, commission takes the first-officer role itself and runs the pilot: it dispatches your seed entities and reports back when the workflow goes idle or a gate needs your decision. From there, [Operate a workflow](operating.md) covers the day-to-day loop.
