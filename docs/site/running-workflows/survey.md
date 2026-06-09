# Survey an existing project

Survey reads a brownfield project's own agent history and reports the implicit workflow you've been running — before you commit to any structure. Use it when you arrive at a project with prior agent work and want the lay of the land.

## Run the survey skill

```bash
spacedock claude "/spacedock:survey"
```

Point it at a project you already have. It reads your own agent session history (read-only) and reports three things:

- the workflow you've been running without naming it,
- how you've been calling work done,
- the decisions still open and waiting on you.

## From survey to a workflow

The survey is read-only — it changes nothing. After reporting, it offers to commission a Spacedock workflow from what it found, so the implicit structure becomes an explicit one you can drive. From there, see [Commission a workflow](commission.md) and [Operating a workflow](operating.md).
