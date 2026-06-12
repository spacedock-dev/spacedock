# Operating a workflow

Operating a workflow is a loop: see what is ready, dispatch the first officer to move it, and make a decision when work reaches a gate. The captain drives that loop; the first officer does the orchestration and the ensigns do the stage work.

## The day-to-day loop

You run the same three steps each session:

1. **See what is ready.** Query workflow state to find the entities that can move (the dispatchable set) and where everything sits.
2. **Dispatch the first officer.** Launch a session and let it pull a dispatchable entity through its next stage.
3. **Handle gates.** When a stage is gated, the first officer stops and presents the result. You approve, reject, or send it back.

The loop ends a session when nothing is dispatchable, or when every dispatchable entity is waiting on a gate decision from you.

## See what is ready

`spacedock status` reads the workflow state and prints it. Run it against the workflow directory, the one holding the commissioned `README.md`:

```bash
spacedock status --workflow-dir docs/dev
```

Prints the status table: one row per active entity, with its ID, title, and current stage. This is the full picture.

To list only the entities ready to dispatch, add `--next`:

```bash
spacedock status --workflow-dir docs/dev --next
```

Prints the dispatchable set: the entities whose next stage can run now, given concurrency limits and what is already in flight. This is the query the first officer runs each loop. When nothing is ready, the result is empty; there is nothing to dispatch.

To filter the table by a frontmatter field, use `--where "field=value"`. The filter takes `=` (equals) or `!=` (not equals):

```bash
spacedock status --workflow-dir docs/dev --where "status=ideation"
spacedock status --workflow-dir docs/dev --where "verdict!="
```

The first prints every entity in the `ideation` stage. The second prints entities whose `verdict` field is set (the `!=` against an empty value). Use `--where` to answer targeted questions: what is in a given stage, which entities carry an external `issue`, which already have a `verdict`.

Two more queries are worth knowing:

- **`--validate`** checks every entity against the workflow's contract and reports problems (a missing or malformed ID, a duplicate ID, a stage name that breaks the naming rule). Run it when the table looks wrong.
- **`--resolve REF`** looks up one entity by slug, full ID, or ID prefix, so you can name it unambiguously before acting on it.

All status queries are read-only. They print state, they do not change it. For the full flag list and the `--set` and `--archive` mutation forms, see the [Command reference](../reference/command-reference.md).

## Dispatch

Hand the first officer the workflow and let it run the dispatch cycle. Launch with your harness subcommand and a task that names the work; the first officer takes the workflow directory from the path you give it in the task, or runs `spacedock status --discover` to find it:

```bash
spacedock claude "/spacedock:first-officer operate the workflow in docs/dev"
```

The first officer reads the workflow `README.md`, runs its own `status --next`, and for each dispatchable entity it dispatches an ensign to move the entity through its next stage. The ensign does the stage work (write the design, produce the deliverable, run the validation), commits, and files a stage report. The first officer reads the report, checks it against the stage's outputs and the entity's acceptance criteria, and advances the entity. A completed non-gated, non-terminal stage is not a stopping point: the first officer advances it and dispatches the next stage on its own, without waiting for you.

It stops and returns to you only at a gate, at a terminal entity's merge ceremony, on a blocker, or when nothing is left to dispatch.

## Handle gate decisions

A gate stops the loop and hands you the call: the first officer presents the stage report and its review, then waits. It never self-approves. You make one of three calls: **approve**, **redo with feedback**, or **reject**. What each one does, what the gate report carries, and the feedback-cycle cap are covered in full in [the three calls](../concepts/gates-and-decisions.md#the-three-calls).

When you approve a terminal stage, the entity is closed: the first officer records the merge, sets the `completed` timestamp and `verdict`, clears the worktree, and tears the worker down. At that point the loop returns to the top: run `status --next` and see what moved into reach.
