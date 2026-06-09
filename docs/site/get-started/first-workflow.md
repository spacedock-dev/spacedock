# Your first workflow

A workflow is a directory of plain-text work item files plus a README that defines the stages, the schema, and the gates. You commission one by describing what you want — Spacedock generates the scaffolding and starts driving it.

## Commission a workflow

Describe the workflow you want:

```bash
spacedock claude "/spacedock:commission Dev task workflow: design -> plan ->
implement -> review, with the design and implementation plan inlined in each work
item, implementation on isolated worktrees with strict TDD, design and review
gated for approval."
```

The first officer commissions the workflow and opens a worktree for the implementation stage. It pauses at the design and review gates for your call.

The same shape drives non-dev work. This example triages a Gmail inbox (it requires a Gmail integration set up before you run it):

```bash
spacedock claude "/spacedock:commission Email triage: fetch, categorize, and act
on my Gmail inbox. Entity: a batch of up to 50 emails. Stages: intake (triage
in:inbox, categorize, propose an action per email as a table) -> approval
(Captain reviews the proposal) -> execute (carry out approved actions). Walk me
through Gmail setup if needed."
```

## What happens next

The first officer reads the workflow README, checks which items are ready to advance, and dispatches ensigns. Stages that need isolation run in their own git worktree; lightweight stages run inline. At a gate, the first officer pauses and presents the stage report for a decision: approve, redo with feedback, or reject. Some gates wait on you; others resolve through a delegated agent review. Rejected work bounces back to an earlier stage for revision, and a hard cap prevents loops.

## Keep going

- To understand the roles and the lifecycle behind this, read [Concepts](../concepts/operating-model.md).
- To drive a workflow day to day — dispatch, gates, status queries — see [Running workflows](../running-workflows/operating.md).
- When you end a session, `/spacedock:debrief` captures what happened so the next session picks up where you left off.
