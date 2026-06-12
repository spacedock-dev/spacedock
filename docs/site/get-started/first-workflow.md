# Your first workflow

Your first workflow comes from one of two places: [survey](first-launch.md)
offers to build one from what it found in your project, or you describe one from
scratch to the `/spacedock:commission` skill. Either way you land here. A
workflow is a directory of plain-text work items plus a README that defines the
stages they move through, the schema each item carries, and the gates where you
make a call. This page walks the commission end to end: the questions it asks,
the design and review gates it sets up, and what happens once the workflow starts
running.

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
gates; the first officer is the orchestrator agent that runs the workflow, and
the ensign is the worker that moves one entity through one stage. The
[operating model](../concepts/operating-model.md#three-roles) covers the three
roles in full.

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
on your seed items). In the design phase it asks, one question at a time: what
the workflow is for and what each entity is (a "design idea" becomes an `idea`),
the stages an entity moves through, which of those stages are gated and where a
rejected entity bounces back to, and the quality bar for each stage. Last come the
entity-ID style and two or three seed items to start on. You confirm or adjust
each proposal. The
[commission reference](../running-workflows/commission.md#the-four-things-you-name)
covers every decision in full.

You do not have to get every answer right. After the questions, the skill
presents the full design as a summary (stages, gates, seed items, where the
files will live) and waits. **Nothing is generated until you accept.** Tell it
what to change and it re-presents.

## What gets generated

Once you accept, the skill writes the workflow into a new directory under `docs/`:
a `README.md` that is the workflow's living spec (mission, schema, and a section
per stage with its `Good:` / `Bad:` bar) and one file per seed entity. The
per-stage prose is a best-guess starting point. Tighten it before any work runs,
because an agent dispatched against a vague bar is expensive to correct. See
[what gets generated](../running-workflows/commission.md#what-gets-generated) for
the full file layout and the `review stages` walk.

## The design and review gates

This is the line Spacedock draws: work flows through the stages, but
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
seed entities: it takes the first-officer role itself, reads the README, and
dispatches ensigns to move ready entities through their stages until the workflow
goes idle or reaches a gate. From there you are running the workflow: approving,
sending back, and resuming in later sessions. [Operating a workflow](../running-workflows/operating.md)
covers the day-to-day loop and how to resume.
