# The stage lifecycle

An entity moves through a fixed chain of stages, one at a time, and each stage declares the work it owns and the proof it must produce. The dev workflow's chain is `backlog → ideation → implementation → validation → done`; your own workflow names its own stages, but the mechanics are the same. The first officer advances an entity stage by stage, dispatching one ensign per stage and pausing at the gates you declared.

The stage order, names, and per-stage properties live in the workflow README's frontmatter under `stages.states`. This page uses the dev workflow (`docs/site/contributing/development-workflow.md`) as the running example; read that page for the full per-stage Inputs/Outputs/Good/Bad detail.

## What a stage declares

Each entry under `stages.states` is a stage name plus a set of boolean or string properties. The first officer reads these to decide how to dispatch and when to stop. A `stages.defaults` block sets the baseline; a stage entry overrides it. The properties that change behavior:

| Property | Effect |
|----------|--------|
| `initial: true` | The stage an entity starts in. The dev workflow marks `backlog`. |
| `terminal: true` | The stage an entity ends in. Reaching it runs the merge and cleanup ceremony, not another dispatch. The dev workflow marks `done`. |
| `gate: true` | The first officer presents a stage report and waits for your decision instead of advancing on its own. |
| `worktree: true` | The stage's work runs in an isolated git worktree. Absent or `false`, it runs inline. |
| `fresh: true` | The stage always gets a freshly dispatched ensign, never a worker reused from the prior stage. |
| `feedback-to: {stage}` | On rejection, work routes back to the named stage rather than failing outright. |
| `concurrency: N` | How many entities may sit in this stage at once. |
| `agent: {name}` | Which worker skill the first officer dispatches. Defaults to `ensign`. |

Beyond these properties, the prose of each stage's `###` subsection in the README is the stage definition: its Inputs, Outputs, and the Good/Bad bar. The first officer copies that subsection verbatim into the ensign's assignment, so what a stage declares in prose is exactly what the worker is told to do.

## The stages

Read the chain as a pipeline: each stage takes the prior stage's output as its input, and the bar rises from "is this clear?" to "is this proven?".

- **`backlog`, the seed.** An entity enters here when first proposed: a title, a source, a brief description, and the test gates future stages must satisfy. No design work yet. `initial: true`, so this is where every new entity starts. The dev workflow also marks it `gate: true`, so the first officer presents a new entity for your go-ahead before it advances.
- **`ideation`, the design.** A worker clarifies the problem, explores approaches, and produces a fleshed-out body: problem statement, proposed approach, acceptance criteria, and a test plan. Each acceptance criterion names how it will be checked. This stage is `gate: true` in the dev workflow, so the first officer presents the design for your approval before any code is written.
- **`implementation`, the deliverable.** A worker produces the change the entity describes: code, fixtures, instruction text, on-disk state. This stage is `worktree: true`, so the work happens in an isolated checkout. Implementation completing is not a stopping point. A completed, non-gated stage routes straight on to validation.
- **`validation`, the independent check.** A worker verifies the deliverable against the acceptance criteria. It checks what was produced; it does not produce the deliverable itself. It reproduces the evidence each `AC-N` cites and returns a `PASSED` or `REJECTED` recommendation. This stage is `worktree: true`, `fresh: true`, `feedback-to: implementation`, and `gate: true`.
- **`done`, the verdict.** Validation is complete and you approve the result. The entity is closed with a `verdict` of `PASSED` or `REJECTED` and a `completed` timestamp. `terminal: true`.

## Fresh context at validation

Validation declares `fresh: true` because the reviewer must not be the maker. The first officer normally reuses a live worker across consecutive stages to save context, but a `fresh: true` stage forces a new dispatch every time. The validator arrives without the implementer's reasoning in its context, sees only the entity body and the deliverable, and pushes back on thin evidence. This is the mechanism behind the README's claim that "the agent doesn't get to judge its own work."

When validation recommends `REJECTED`, `feedback-to: implementation` routes the concrete finding back to the implementation stage for rework rather than closing the entity. The entity re-enters implementation, the finding is addressed, and a fresh validator checks it again. A hard cap on feedback cycles prevents an endless bounce; on the third cycle the first officer escalates to you.

## Worktree vs. inline

A stage runs in an isolated git worktree when it declares `worktree: true`, and inline at the repo root otherwise. This is the "isolation when it matters" tradeoff: stages that mutate shared state (implementation, validation) get their own checkout so concurrent entities don't collide; lighter stages that only edit the entity body (backlog, ideation) run inline.

The mechanics, run by the first officer:

- **On first dispatch to a worktree stage,** the first officer creates a worktree at `.worktrees/{worker_key}-{slug}` on branch `{worker_key}/{slug}` and records the path in the entity's `worktree` frontmatter field.
- **Inside a worktree-backed stage,** the ensign keeps all reads, writes, and commits under that worktree. The deliverable is isolated there until the entity terminalizes.
- **In a split-root workflow** (the README declares a `state:` checkout, e.g. `state: .spacedock-state`), the entity body and stage report live in the state checkout, not in the worktree; the worktree isolates only the deliverable. The first officer's dispatch hands the ensign the correct state-checkout path for the entity; the worker trusts that path rather than writing entity state into the worktree.
- **At the terminal stage,** the first officer merges the worktree branch, clears the `worktree` field, removes the worktree, and deletes the local branch.

To see where each entity sits and which are ready to advance, read the workflow state:

```bash
spacedock status --workflow-dir docs/dev
```

```bash
spacedock status --workflow-dir docs/dev --next
```

`--next` lists the entities ready for dispatch, the query the first officer runs each loop. The `worktree` column shows the isolated checkout path for any entity currently mid-stage in a worktree-backed stage.

## Where to go next

- The roles that drive this pipeline (captain, first officer, ensign) are in [the operating model](operating-model.md).
- The decision points at stage boundaries are covered in [gates and decisions](gates-and-decisions.md).
- The entity frontmatter these stages update is in the [frontmatter contract](../reference/frontmatter-contract.md).
