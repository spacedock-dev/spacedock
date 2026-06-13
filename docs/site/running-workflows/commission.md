# Commission a workflow

`/spacedock:commission` turns a description of the work you want tracked into a runnable workflow. You make four design decisions in conversation; the skill derives the rest, presents the full design, and writes nothing until you accept. [Your first workflow](../get-started/first-workflow.md) shows the whole flow, including the design summary and pilot run you will see.

Invoke it from a session started with `spacedock claude`. You can pass the mission inline:

```
/spacedock:commission product idea to simulated customer interview
```

Text after the command name becomes the workflow mission, presented for confirmation rather than asked from scratch. With no argument, the skill greets you and asks.

## The four things you decide

1. **The mission and what each work item is.** What the workflow is for, and what one entity represents. The description you give becomes the label the workflow uses everywhere it talks about your work: "a design idea" makes it a workflow of ideas.

2. **The stages.** The skill detects the workflow's shape from your mission (shipping code, testing a hypothesis, or iterating on an artifact) and proposes a stage list for you to confirm, modify, add to, or trim. It pushes back on redundant names: `awaiting_validation` reads as "the entity is in awaiting_validation," so it suggests `validation` instead.

3. **The gated stages.** A gate is a stage where the workflow pauses for your decision before an entity advances. By default one gate sits before the terminal stage, and each gate gets a rejection target, the earlier stage work bounces back to when you reject. Both appear in the design summary in plain language ("If you reject at `review`, it goes back to `draft` for revision").

4. **The per-stage quality bar.** Each stage in the generated README tells a dispatched ensign what "good" means there: what to produce, the quality bar to meet, the anti-patterns to avoid. The skill drafts these from the mission, but they are starting prose, not commitments.

Everything else is derived or asked with a recommendation attached: where the files go, how entities are identified, how rejections route. The design summary shows it all; ask for the tradeoffs on any of it before you accept.

## What gets generated

After you accept, everything lands in `docs/{mission-slug}/` as plain text: the README that is the workflow's living spec, and one file per seed entity. A stage that writes code or produces artifacts runs each entity on an isolated branch, so your main checkout stays clean; if you ship through PR review, the skill also offers the [pr-merge mod](../advanced/mods-and-standing-teammates.md), which manages the PR lifecycle so merging never needs to be a stage.

## Tighten the README before the first dispatch

The README is the workflow's living spec. Before launching, the skill reminds you that the auto-generated per-stage rules are best-guesses and prompts you to tighten them. You have two ways to do it:

- Open `docs/{mission-slug}/README.md` and edit the bullets under each stage heading directly.
- Type `review stages` to have the skill walk you through each stage one at a time, flag the rules that read as generic, and apply your amendments inline.

Editing here costs minutes; un-editing after agents have been dispatched against vague rules costs more.

## The pilot run and after

On accept, commission takes the first-officer role itself and runs the pilot: it dispatches your seed entities and reports back when the workflow goes idle or a gate is waiting on your decision. In every later session, `spacedock claude` picks up where the last one left off; [Operating a workflow](operating.md) covers that loop.
