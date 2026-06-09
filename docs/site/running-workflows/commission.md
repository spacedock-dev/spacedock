# Commission a workflow

Commissioning generates a new workflow from a description. You say what stages you want, what an entity is, and which stages are gated; Spacedock writes the scaffolding and the first officer starts driving it.

## Run the commission skill

```bash
spacedock claude "/spacedock:commission <describe your workflow>"
```

A useful description names four things:

- **The stages** and their order (e.g. `design -> plan -> implement -> review`).
- **The entity** — what one work item is (a task, a batch of emails, a document).
- **Which stages are gated** — where you want to approve before the work advances.
- **Any per-stage rules** — isolated worktrees, strict TDD, inlined design notes.

## Example

```bash
spacedock claude "/spacedock:commission Dev task workflow: design -> plan ->
implement -> review, with the design and implementation plan inlined in each work
item, implementation on isolated worktrees with strict TDD, design and review
gated for approval."
```

The first officer commissions the workflow and opens a worktree for the implementation stage, pausing at the design and review gates for your call.

## After commissioning

Once the workflow exists, you drive it: add entities, dispatch stages, and handle gates. See [Operating a workflow](operating.md) for the day-to-day loop. If you're starting from an existing project rather than a blank slate, [survey it first](survey.md).
