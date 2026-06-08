# Spacedock

**Spacedock is a multi-agent orchestrator where nothing ships without a decision.**
It lives within your existing harness: Claude Code, Codex, or Pi. It breaks work
into stages and surfaces the decisions each stage needs, batched for you. Each
decision arrives with evidence measured against a predefined bar for what good
looks like. You approve, send back, or escalate. Or you delegate the call to an
agent. Either way, the decision is recorded with its evidence and reason.

**Why?**

- You're the human, and agents from a dozen sessions ping you all day: a design
  call, a one-line approval, "ship without full coverage?", none of it on a
  schedule you can plan around. You stopped deciding the work and started just
  answering the agent.
- You're the agent, and you stall every few steps waiting on a human who isn't
  watching, with no clear scope and no clear bar for done, so you ask and then
  you wait.
- Generation got cheap. Your attention and judgment are now the bottleneck.
  Every task still ends in a decision, and no one owns it, so it falls on
  whoever's around.

**Start with what you already built.** Point Spacedock at a project you
vibe-coded into spaghetti and run `/spacedock:survey`. It reads your own agent
history and shows you three things: the workflow you've been running without
naming it, how you've been calling work done, and the decisions still open and
waiting on you.

## What's different

- **The agent doesn't get to judge its own work.** Review runs as a separate
  stage with fresh context, no access to the maker's reasoning. It pushes back on
  thin evidence and work that looks busy without proving its claim.
- **Every decision leaves a trail.** Each gate, the decision point at the end of
  a stage, carries a stage report: findings, verdicts, artifacts, anomalies. You
  decide on evidence, not the transcript.
  The record outlives the reviewer, so you can trace a bad result back to the
  call that caused it.
- **The bar sharpens as you use it.** Each stage declares what good means, and
  the agent works to that line on its own. When a standard turns out fuzzy in
  practice, the agent proposes an edit to the stage's written criteria for your
  approval. `/spacedock:debrief` captures each session's learnings so the next
  one starts from them.
- **Batch the work; decide as it flows back.** Queue many items at once. Agents
  advance each through its stages. You handle gates as they surface, not one
  session at a time.
- **Isolation when it matters.** Stages that touch shared state run in their own
  git worktree. Lighter stages run inline.
- **Sandboxing built in.** Drop a profile in the project and Spacedock runs the
  agent sandboxed.
- **Work survives the context limit.** When an agent runs out of context, a
  successor carries forward what's in flight.

## Install

Prerequisite: a coding agent harness. Claude Code, Codex, and Pi are tier-1
supported; through skill systems it also runs in most other harnesses, including
Hermes-class agents.

Install with Homebrew:

```bash
brew tap spacedock-dev/homebrew-tap
brew install spacedock
```

Then launch. The first launch sets up the plugin for you, so a single line gets
you a working session. Point it at a project you already have and let it survey:

```bash
spacedock claude "/spacedock:survey"
```

Using Codex or Pi instead? Swap the subcommand: `spacedock codex "/spacedock:survey"`
or `spacedock pi "/spacedock:survey"`.

See [`docs/install-journey.md`](docs/install-journey.md) for the full first-run
walkthrough, the Codex and Pi paths, and a from-source build for development.

## Quick start

Commission a workflow by describing what you want:

```bash
spacedock claude "/spacedock:commission Dev task workflow: design -> plan ->
implement -> review, with the design and implementation plan inlined in each work
item, implementation on isolated worktrees with strict TDD, design and review
gated for approval."
```

The first officer commissions the workflow and opens a worktree for the
implementation stage. It pauses at the design and review gates for your call.

The same shape drives non-dev work. This example triages a Gmail inbox. It
requires a Gmail integration set up before you run it:

```bash
spacedock claude "/spacedock:commission Email triage: fetch, categorize, and act
on my Gmail inbox. Entity: a batch of up to 50 emails. Stages: intake (triage
in:inbox, categorize, propose an action per email as a table) -> approval
(Captain reviews the proposal) -> execute (carry out approved actions). Walk me
through Gmail setup if needed."
```

## How it works

A workflow is a directory of plain-text work item files plus a README that
defines the stages, the schema, and the gates. Everything about a work item lives
in the file itself: the problem, the design notes, the bar for done, the stage
reports. State survives a session; the next one picks up where you left off.
Three roles:

| Role | Who |
|------|-----|
| **Captain** | You. You define the mission and make the calls at approval gates unless delegated. |
| **First Officer** | The orchestrator agent that runs the workflow and reports to you at gates. |
| **Ensign** | The worker agent that moves one item forward through one stage. |

The first officer reads the workflow README, checks which items are ready to
advance, and dispatches ensigns. Stages that need isolation run in their own git
worktree; lightweight stages run inline. At a gate, the first officer pauses and
presents the stage report for a decision: approve, redo with feedback, or
reject. Some gates wait on you; others resolve through a delegated agent review.
Rejected work bounces back to an earlier stage for revision. A hard cap prevents
loops.

When you end a session, `/spacedock:debrief` captures what happened: commits,
state changes, decisions, open issues, all in a record the next session picks up.
When a new Spacedock release is out, `/spacedock:refit` upgrades your workflow
scaffolding while keeping local modifications.

## Usage

```bash
spacedock claude "task" [--safehouse…] [-- host-flags…]   # launch the first officer in Claude Code
spacedock codex  "task" [--safehouse…] [-- host-flags…]   # launch the first officer in Codex
spacedock pi     "task" [--safehouse…] [-- host-flags…]   # launch the first officer in Pi
spacedock doctor                                          # plugin compatibility check
spacedock --version                                       # print the installed version
```

## License

Spacedock is released under the [Apache License 2.0](LICENSE).
