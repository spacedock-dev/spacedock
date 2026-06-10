# Your first workflow

A workflow is a directory of plain-text work items plus a README that defines
the stages they move through, the schema each item carries, and the gates where
you make a call. You create one by describing what you want, in plain language,
to the `/spacedock:commission` skill. This page walks that first commission end
to end: the questions it asks, the design and review gates it sets up, and what
happens once the workflow starts running.

A few terms used below, defined on first use:

- An **entity** is one work item: a single markdown file (the README also calls
  it a "work item"). A bug report, a design idea, a feature: whatever the
  workflow processes, each one is an entity.
- A **stage** is a bucket an entity sits in as it advances, for example
  `ideation` or `implementation`. The first entity starts in the first stage and
  moves toward a terminal one.
- A **gate** is a decision point at the end of a stage where the workflow pauses
  for your call instead of advancing on its own.

You are addressed as the captain, the workflow operator who makes the calls at
gates. The first officer is the orchestrator agent that runs the workflow; the
ensign is the worker agent that moves one entity through one stage.

## Commission a workflow

Run `/spacedock:commission` inside a Spacedock session and describe the work in
the same line:

```bash
spacedock claude "/spacedock:commission Track design ideas through review stages"
```

If you have not launched a session yet, see
[Install Spacedock](install.md) first. You can also start bare
(`/spacedock:commission` with no description) and answer the questions from
scratch.

The skill greets you and walks three phases: **design** (a few questions),
**generate** (it writes the files), and a **pilot run** (it starts the workflow
on your seed items). The design phase asks you three things, one question at a
time:

1. **The mission and the entity.** What the workflow is for, and what each work
   item represents: "a design idea", "a bug report", "a candidate feature". The
   skill derives a short label from your answer (a "design idea" becomes an
   `idea`).
2. **The stages.** It proposes an ordered list from your mission (for a design
   workflow, something like `backlog → ideation → implementation → validation →
   done`), and you confirm, add, remove, or rename them.
3. **Seed entities.** Two or three starting items to run through, each with a
   title and a short description (and an optional score). These become the
   workflow's first work.

From your answers the skill then derives the gates: which stages pause for your
approval, and which earlier stage rejected work bounces back to. By default it
gates the stage before the terminal one.

You do not have to get every answer right. After the questions, the skill
presents the full design as a summary (stages, gates, seed items, where the
files will live) and waits. **Nothing is generated until you accept.** Tell it
what to change and it re-presents.

## What gets generated

Once you accept, the skill writes the workflow into a new directory under
`docs/` and confirms each file it created:

- `README.md`, the workflow's living spec. It holds the mission, the schema each
  entity carries, and a section per stage describing its inputs, outputs, and
  quality bar (`Good:` / `Bad:`).
- One file per seed entity, named from its title, with YAML frontmatter that
  records its `status`, `score`, and other fields.

The per-stage prose in the README is a best-guess starting point, not a
commitment. The skill flags this directly and offers a `review stages` walk that
steps through each stage's expectations so you can tighten the quality bar before
any work runs. Tightening here is cheap; an agent dispatched against a vague bar
is not.

## The design and review gates

A gate is where the workflow stops and hands you a decision instead of advancing
on its own. This is the line Spacedock draws: work flows through the stages, but
**nothing crosses a gate without a recorded decision.** A development workflow
gates the design stage and the review stage among others, so you sign off on
the approach before code is written, and on the result before it ships.

At each gate the first officer pauses and presents a stage report: the chosen
direction, the evidence behind it, and a single recommendation. You make one of
three calls:

- **Approve**, and the entity advances to the next stage.
- **Redo with feedback**: it goes back for revision against the notes you give.
- **Reject**, and it bounces to an earlier stage (the one the design named as the
  rejection target) to be reworked.

You decide on the report and its evidence, not on the agent's transcript. The
decision is recorded with its reason, so a result can later be traced back to the
call that produced it.

## What happens after

When you accept the design, the commission skill launches a pilot run on your
seed entities. It takes on the first-officer role itself for this first run:
it reads the workflow README, checks which entities are ready to advance, and
dispatches ensigns to move them through their stages. Stages that modify the
repo run in their own git worktree; lighter stages run inline.

The run proceeds until the workflow goes idle or reaches a gate. There the first
officer stops and reports what happened: which entities moved, which stages they
passed through, which gate is waiting on you. From that point you are running the
workflow: approve, send back, or reject, and work continues.

To resume the workflow in a later session, launch the first officer with no task:

```bash
spacedock claude
```

It reads the saved workflow state, picks up where you left off, and dispatches
ensigns for any entity ready for its next stage.
