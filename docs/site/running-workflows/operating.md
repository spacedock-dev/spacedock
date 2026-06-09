# Operating a workflow

Driving a workflow is a loop: see what's ready, dispatch it, handle the gates that surface. The first officer runs the loop; you make the calls.

## The loop

The first officer reads the workflow README, checks which entities are ready to advance, and dispatches ensigns. Stages that need isolation run in their own git worktree; lightweight stages run inline. As work flows back, the first officer presents each gate for your decision — approve, redo with feedback, or reject.

## Check workflow state

Read the current state of a workflow with the launcher:

```bash
spacedock status --workflow-dir docs/dev
```

To list the entities ready for dispatch — the query the first officer runs each loop:

```bash
spacedock status --workflow-dir docs/dev --next
```

To filter entities by a frontmatter field, use `--where`:

```bash
spacedock status --workflow-dir docs/dev --where sprint=0200-flip
```

See the [Command reference](../reference/command-reference.md) for the full `status` surface.

## Handling gates

When the first officer pauses at a gate, it presents the stage report. Decide on the evidence in the report, not the transcript:

- **Approve** to advance the entity.
- **Redo with feedback** to bounce it back to an earlier stage with concrete direction.
- **Reject** when it does not meet the bar.

For the mechanics of gates, feedback cycles, and the loop cap, see [Gates & decisions](../concepts/gates-and-decisions.md).
