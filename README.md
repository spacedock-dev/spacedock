# Spacedock

Hand an agent a multi-step job and it drifts, skips steps, or invents its own
path. Hand it a workflow and it follows the workflow.

Spacedock runs agent work through stages you define. Each stage gets a fresh
context, an explicit gate, and a structured report. Reviewers are allowed to
push back. State lives outside the agent, so work survives the context limit and
picks up the next session where it left off.

**You want Spacedock if:**

- **You delegate repeatable work to agents** — the same pipeline run over many
  inputs (code reviews, content drafts, outreach batches) — and you want each
  run to finish, not stop halfway. Spacedock dispatches a fresh agent per stage,
  with an approval gate where your judgment actually matters.
- **You don't trust single-agent output** on its own. Spacedock review stages
  push back instead of rubber-stamping, with a 3-strikes escalation so you only
  see what the reviewer couldn't resolve.
- **Your work spans days or weeks** — a blog draft, a launch plan, a multi-week
  benchmark. Spacedock holds artifact state outside the agent, so the next
  session resumes instead of restarts.

## What's different

- **Approval gates with structured evidence.** Every gate comes with a stage
  report: findings, verdicts, artifacts, anomalies. You approve, redirect, or
  bounce back without sifting through raw output.
- **Adversarial review gates.** Review stages can be configured to push back
  rather than rubber-stamp, targeting thin evidence and work that looks busy
  without proving its claim.
- **Plan in batches, decide as work flows back.** Queue many work items at once;
  agents advance each through its stages while you handle approvals as they
  surface.
- **The workflow learns with you.** When a pattern emerges — a stage that never
  fires, a gate that keeps bouncing the same issue — the first officer helps you
  adjust the workflow.
- **Isolation when needed.** Stages that touch shared state run in their own git
  worktree; lightweight stages run inline.
- **Work doesn't die at the context limit.** When an agent runs out of context,
  a successor carries forward what's in flight.

## Install

Spacedock is two pieces: the `spacedock` launcher and a host plugin (the
first-officer and ensign agents) loaded by Claude Code or Codex.

Install the launcher with Homebrew, then add the plugin:

```bash
brew install spacedock-dev/homebrew-tap/spacedock
spacedock install --host claude
```

That installs the launcher, adds the Spacedock plugin to Claude Code, and runs a
compatibility check. Now launch the first officer with a task:

```bash
spacedock claude -- "your task"
```

See [`docs/install-journey.md`](docs/install-journey.md) for the full first-run
walkthrough, the Codex path, and a from-source build for development.

> [agent-safehouse](https://agent-safehouse.dev) is an optional sandbox for
> agent runs — install it separately. A `.safehouse` profile in the working
> directory (or a `--safehouse` flag) wraps the launch through it.

## Quick start

Commission a workflow by describing what you want it to do:

```bash
spacedock claude -- "/commission Email triage: fetch, categorize, and act on my
Gmail inbox. Entity: a batch of up to 50 emails. Stages: intake (triage
in:inbox, categorize, propose an action per email as a table) -> approval
(Captain reviews the proposal) -> execute (carry out approved actions). Walk me
through Gmail setup if needed."
```

The first officer commissions the workflow, dispatches an ensign to gather your
inbox, then pauses with a categorized proposal and waits for your approval
before touching anything.

For a development workflow:

```bash
spacedock claude -- "/commission Dev task workflow: design -> plan -> implement
-> review, with the design and implementation plan inlined in each work item,
implementation on isolated worktrees with strict TDD, design and review gated
for approval."
```

## How it works

A workflow is a directory of markdown work item files plus a README that defines
the stages, the schema, and the gates. There are three roles:

| Role | Who |
|------|-----|
| **Captain** | You. You define the mission and make the calls at approval gates. |
| **First Officer** | The orchestrator agent that runs the workflow and reports to you at gates. |
| **Ensign** | The worker agent that moves one item forward through one stage. |

The first officer reads the workflow README, checks which items are ready to
advance, and dispatches ensigns. Stages that need isolation run in their own git
worktree; lightweight stages run inline. At a gate the first officer pauses and
presents the ensign's stage report: approve, redo with feedback, or reject.
Rejected work bounces back to an earlier stage for revision, with a hard cap so
you never get stuck in a loop.

A work item carries everything in its body:

```yaml
---
id: 054
title: Session debrief command
status: done
---

Problem statement, design notes, acceptance criteria, and stage reports all
live in the body of this file as the work moves through its stages.
```

When you end a session, `/spacedock:debrief` captures what happened — commits,
state changes, decisions, open issues — into a record the next session picks up.
When a new Spacedock release is out, `/spacedock:refit` upgrades your workflow
scaffolding while keeping local modifications.

## Usage

```bash
spacedock claude [host-flags…] [--safehouse…] -- "task"   # launch the first officer in Claude Code
spacedock codex  [host-flags…] [--safehouse…] -- "task"   # launch the first officer in Codex
spacedock doctor                                          # plugin compatibility check
spacedock --version                                       # print the installed version
```

Flags before `--` pass through to the host; the bare text after `--` is the
launch task.

## License

Spacedock is released under the [Apache License 2.0](LICENSE).
