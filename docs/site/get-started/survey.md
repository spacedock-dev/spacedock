# Survey your project

When you use the `spacedock:survey` skill, it looks at your existing agent
conversation logs on local disk (through [agentsview](https://agentsview.io/),
an open source session-history tool). It is read-only; if `agentsview` is
missing, it asks before installing it.

## What it reports

- **The repeated manual behavior**: the loop you have been driving by hand,
  run after run, without naming it.
- **The steering and interruptions observed**: how often your agents needed
  you to step in, and for what.
- **Your current workstreams**: what is in flight, clustered into tracks.
- **What is still undecided**: forks raised but never resolved, cross-checked
  against the repo first so work that already shipped is dropped.

## The commission offer

Survey ends with an offer, not an action. It asks whether you want it to
[commission](../running-workflows/commission.md) a real Spacedock
[workflow](../concepts/workflows-and-entities.md) built from what it found,
turning the open decisions into
[approval gates](../concepts/gates-and-decisions.md) and the workstreams into
work items. Nothing has changed in your project until you say yes.

The offer matches the work it saw, citing real numbers from the scan:

- **Routine loops** (issue → worktree → PR) get an automation offer: a
  workflow that gates the crucial decisions and lets the agent drive between
  gates.
- **Exploration** (creative or design work where your steering is the point)
  gets a book-keeping offer: structure for the parallel threads and their
  state. There is no automate-the-human-out pitch; the involvement is the work.

On a **yes**, the workflow is built from what survey found: the inferred loop
becomes the proposed stages, the workstreams become the seed work items, and
the open forks become the gates. On a **no**, it stops; the survey stands on
its own as orientation.

To define a workflow yourself instead, see
[your first workflow](first-workflow.md).
