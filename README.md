# Spacedock

**Spacedock is a multi-agent orchestrator where nothing ships without a decision.**
It lives inside your agent. Each stage ends at a gate. The work arrives with its
evidence, measured against a predefined bar for what good looks like. The human
partner or a delegated agent says yes, sends it back, or escalates. Either way
the decision gets recorded, with its evidence and its reason.

**Why?**

- You're the human, and agents from a dozen sessions ping you all day: a design
  call, a one-line approval, "ship without full coverage?", none of it on a
  schedule you can plan around. You stopped deciding the work and started just
  answering the agent.
- You're the agent, and you stall every few steps waiting on a human who isn't
  watching, with no clear scope and no clear bar for done, so you ask and then
  you wait.
- It's one broken loop on both ends, because generation got cheap and
  verification didn't. Every task still ends in a decision, and right now nobody
  owns it. It lands on whoever's awake.

**Start with what you already built.** Point Spacedock at a project you
vibe-coded into spaghetti and run `/spacedock:survey`. It reads your own agent
history and shows you three things: the workflow you've been running without
naming it, how you've been calling work done, and the decisions still open and
waiting on you.

## What's different

- **The maker doesn't judge itself.** Review stages are adversarial by default.
  They push back on thin evidence and work that looks busy without proving its
  claim. A 3-strikes escalation means you only see what the reviewer couldn't
  settle.
- **Every decision leaves a trail.** Each gate carries a stage report: findings,
  verdicts, artifacts, anomalies. You decide on evidence, not the transcript.
  The record outlives the reviewer, so you can trace a bad result back to the
  call that caused it.
- **You set the bar; the agent works to it.** Each stage declares what done
  means. The agent runs to that line on its own. You show up at the gate, not in
  the loop.
- **The bar sharpens as you go.** When a stage never fires, or a gate keeps
  bouncing the same issue, the first officer helps you adjust the workflow. Good
  is discovered, not declared up front.
- **Batch the work; decide as it flows back.** Queue many items at once. Agents
  advance each through its stages. You handle gates as they surface, not one
  session at a time.
- **Isolation when it matters.** Stages that touch shared state run in their own
  git worktree. Lighter stages run inline.
- **Work survives the context limit.** When an agent runs out of context, a
  successor carries forward what's in flight.

## Install

Spacedock is two pieces: the `spacedock` launcher and a host plugin (the
first-officer and ensign agents) loaded by Claude Code, Codex, or Pi. Install the
launcher with Homebrew:

```bash
brew install spacedock-dev/homebrew-tap/spacedock
```

Then launch. The first command installs the plugin if it's missing, so a single
line gets you a working session:

```bash
spacedock claude "your task"
```

Using Codex or Pi instead? Swap the subcommand: `spacedock codex "your task"` or
`spacedock pi "your task"`.

To stay current: `brew upgrade spacedock` updates the launcher; `spacedock
install` refreshes the plugin. Run `spacedock doctor` any time to confirm the
launcher and plugin are compatible. A plain `claude plugin update` will not pick
up new releases. The plugin ships through `spacedock`.

See [`docs/install-journey.md`](docs/install-journey.md) for the full first-run
walkthrough, the Codex and Pi paths, and a from-source build for development.

> [safehouse](https://agent-safehouse.dev) is an optional sandbox for agent
> runs. Install it separately. A `.safehouse` profile in the working directory
> (or a `--safehouse` flag) wraps the launch through it.

## Quick start

Commission a workflow by describing what you want:

```bash
spacedock claude "/commission Email triage: fetch, categorize, and act on my
Gmail inbox. Entity: a batch of up to 50 emails. Stages: intake (triage
in:inbox, categorize, propose an action per email as a table) -> approval
(Captain reviews the proposal) -> execute (carry out approved actions). Walk me
through Gmail setup if needed."
```

The first officer commissions the workflow and dispatches an ensign to gather
your inbox. It then pauses with a categorized proposal and waits for your
approval before touching anything.

For a development workflow:

```bash
spacedock claude "/commission Dev task workflow: design -> plan -> implement
-> review, with the design and implementation plan inlined in each work item,
implementation on isolated worktrees with strict TDD, design and review gated
for approval."
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

The task goes first. Anything after `--` forwards verbatim to the host, e.g.
`spacedock claude "task" -- --model opus`.

## License

Spacedock is released under the [Apache License 2.0](LICENSE).
