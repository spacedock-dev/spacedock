# Spacedock

Spacedock runs agent work through defined stages, so you can delegate in batches and make only the calls that need your judgment.

The first officer coordinates the flow: it dispatches workers to advance each work item and surfaces approval-worthy decisions to you, the captain, so batches move forward without pulling you into every session.

**You want Spacedock if:**

- **You're a human tired of context-switching** between agent sessions to make approval decisions. Spacedock batches the decisions an agent wants to hand back to you and presents each with evidence, so you approve or redirect without re-loading context.
- **You're an agent delegating repeatable work** and want a structured place to queue up approval-worthy decisions for your human without interrupting them for every tiny step.

## What's Different

- **Approval gates with structured evidence.** Every gate comes with a stage report: findings, verdicts, artifacts, anomalies. You approve, redirect, or bounce back faster than sifting through raw output or a sprawling log.
- **Adversarial review gates.** Review stages can be configured to push back rather than rubber-stamp. They target sycophancy, thin evidence, and work that looks busy without proving its claim. Work clears the gate when it survives the challenge.
- **Plan in batches, decide as work flows back.** Queue multiple work items at once; agents advance each through its stages independently while you handle approvals as they surface.
- **The workflow learns with you.** The first officer helps you adjust it when patterns emerge: a stage that never fires, a gate that keeps bouncing the same issue back, a schema field that always ends up empty.
- **Isolation when needed.** Stages that touch shared state run in their own git worktree; lightweight stages run inline. You declare which is which, and the first officer enforces it.
- **Work doesn't die at the context limit.** When an agent runs out of context, Spacedock swaps in a successor that carries forward what's in flight. Nothing gets lost in the handoff.

## Quick Start

**Prerequisites:** [Homebrew](https://brew.sh) and a host CLI — Claude Code (primary), Codex CLI, or Pi.

Spacedock installs as two pieces: the `spacedock` binary (the launcher and contract gate) and the host plugin (the `spacedock:first-officer` / `spacedock:ensign` skills and agents).

### Claude Code

1. Install the binary:

   ```bash
   brew install spacedock-dev/homebrew-tap/spacedock
   ```

2. Install the host plugin:

   ```bash
   spacedock install --host claude
   ```

   This adds the plugin from the `spacedock-dev/spacedock` marketplace and runs the contract doctor — expect `OK: binary contract 1 satisfies plugin range >=1,<2.`

3. Commission a workflow with your own mission prompt:

   ```bash
   spacedock claude "/commission <your mission prompt>"
   ```

   `spacedock claude` launches `claude --agent spacedock:first-officer` behind the contract gate. Tokens after `--` forward verbatim to the host, e.g. `spacedock claude "task" -- --model opus`.

4. Or start from one of these example workflows — copy and run:

   **Email triage:**
   ```bash
   spacedock claude "/commission Email triage: fetch, categorize, and act on Gmail inbox. Entity: a batch of up to 50 emails. Stages: intake (use gws-cli, triage in:inbox and read email body if necessary, categorize, propose action per email, output as table) → approval (Captain reviews proposal) -> execute (carry out approved actions, do not mark as read). Use gws-cli (https://github.com/googleworkspace/cli/tree/main/skills/gws-gmail), GOOGLE_WORKSPACE_CLI_CONFIG_DIR=~/.config/gws/<account> for different accounts. Walk me through gws-cli setup if not already done."
   ```

   **[Superpowers](https://github.com/obra/superpowers)-style dev task workflow:**
   ```bash
   spacedock claude "/commission Dev task workflow: superpowers-style design → plan → implement → review with ## Design and ## Implementation Plan inlined in the entity body (no separate spec/plan files), implement on isolated worktrees with strict TDD, design and review gated for approval."
   ```

### Codex CLI and Pi

The same two steps with a different `--host`, then launch with the matching subcommand:

```bash
spacedock install --host codex            # or: --host pi
spacedock codex "/commission <your mission prompt>"   # or: spacedock pi "..."
```

> Codex multi-agent is experimental. The Claude Code path is the primary supported surface.

### Upgrading

```bash
brew upgrade spacedock
spacedock install --host claude
```

Run `spacedock doctor` any time to check that the installed plugin and binary are contract-compatible; rerun `spacedock install` when it reports the plugin is stale. (A plain `claude plugin update` will not pick up new releases — the plugin ships through `spacedock install`.)

### Sandboxing and source builds

[safehouse](https://agent-safehouse.dev) is a separate runtime dependency for sandboxed launches — not installed by brew. A `.safehouse` profile in the working directory (or the `--safehouse` flags) wraps the launch through it.

To build from source and load the repo's own plugin checkout with `--plugin-dir`, see the [fresh-install journey](https://github.com/spacedock-dev/spacedock/blob/next/docs/install-journey.md).

## Supported models

Spacedock's first officer drives a startup protocol (load operating-contract skills, read stage definitions, enforce gate semantics) that requires careful instruction-following. The supported floor for the **bare-mode** first officer is **Claude Sonnet or above**. Haiku has been empirically observed to flake below the protocol-adherence bar in bare mode (cwd drift, gate self-approve, early-terminate failure modes — see archived empirical evidence in `docs/plans/_archive/haiku-bare-fo-startup-protocol-adherence.md`).

In **teams mode**, the first officer can run reliably on Haiku because the teams runtime provides extra protocol scaffolding (skill preloading, agent-discovery hooks). Ensigns may use any model the user chooses; the floor applies to the first officer role specifically.

If you point Spacedock at a sub-floor model for the bare-mode first officer, expect intermittent stage-advancement failures and gate-state corruption. The Spacedock test matrix exercises `sonnet` for bare-mode and `haiku` for teams-mode.

## What a Work Item Looks Like

```yaml
---
id: 054
title: Session debrief command
status: done
---

Problem statement, design notes, acceptance criteria, and stage reports
all live in the body of this file as the work moves through stages.
```

See [a completed example](https://github.com/spacedock-dev/spacedock/blob/main/docs/plans/_archive/session-debrief.md) from Spacedock's own workflow.

## Concepts

| Concept | What it is |
|---------|------------|
| **Mission** | The purpose of the workflow: what it processes and what it delivers. |
| **Work item** | A single markdown file describing one thing being worked on: an email batch, a dev task, a draft. |
| **Workflow** | A directory of work items plus the README that defines stages, schema, and gates. |
| **Stage** | A named step a work item passes through (e.g. design, implement, review). |
| **Gate** | A pause point at a stage boundary where the captain approves, redirects, or bounces the work back. |

*"I am the master of my fate, I am the captain of my soul."* -- William Ernest Henley, *Invictus*

| Role | Who |
|---------|------------|
| **Captain** | You. You define the mission and make the calls at approval gates. |
| **First Officer** | The orchestrator agent that manages the workflow and reports to you at gates. |
| **Ensign** | The worker agent that moves a single item forward through one stage. |

## How It Works

The first officer reads the workflow README, checks work item statuses, and dispatches ensigns for items ready to advance. Stages that need isolation (typically implementation work with commits) run inside their own git worktree; lightweight stages (design, review, triage) run inline. At approval gates the first officer pauses and presents the ensign's stage report for your review: approve, redo with feedback, or reject. Rejected work automatically bounces back for revision in a fresh round of the earlier stage, with a hard cap so you never get stuck in an infinite loop. When you end a session, `/spacedock:debrief` captures what happened (commits, task state changes, decisions, open issues) into a record the next session picks up automatically (see [an example debrief](https://github.com/spacedock-dev/spacedock/blob/main/docs/plans/_debriefs/2026-04-09-01.md) from a real session).

## What Gets Generated

When you run `/spacedock:commission`, the following files are added to your workflow directory:

- **`{dir}/README.md`**: workflow schema, stage definitions, and work item template
- **`{dir}/*.md`**: seed work item files
- **`{dir}/_mods/`**: local modifications carried across refits

**Shipped by the Spacedock plugin:**

- **`spacedock:first-officer`**: the orchestrator agent that reads workflow state and dispatches ensigns
- **`spacedock:ensign`**: the worker agent dispatched to do stage work
- **`skills/commission/bin/status`**: read and advance workflow state without switching to a separate tracking tool

The generated workflow README is the single source of truth. The first officer reads it to know what stages exist, what quality criteria to enforce, and when to pause for your review.

Workflows can extend their own behavior via markdown mod files (`_mods/*.md`) that declare hook handlers for lifecycle events like startup, idle, or merge. For example, the [`pr-merge` mod](docs/plans/_mods/pr-merge.md) opens a pull request automatically when a completed worktree branch is ready to land.

When a new Spacedock release is available, use `/spacedock:refit` to upgrade your workflow scaffolding while keeping local modifications.

## Tips

- **Run Spacedock inside a sandbox.** Recommended: [safehouse](https://agent-safehouse.dev) (macOS, integrates with the `--safehouse` launch flags), [packnplay](https://github.com/obra/packnplay), a devcontainer, or a VM.
- **Talk directly to an ensign.** Claude Code supports agent team chat: while a dispatched ensign is running, you can `Shift+Up` / `Shift+Down` to switch panes and give the ensign feedback directly instead of routing everything through the first officer.

## Use Cases

- **Email triage**: classify and route incoming messages with AI agents, escalate to a human at review gates
- **Dev task workflow**: [superpowers](https://github.com/obra/superpowers)-style design -> plan -> implement -> review with approval gates
- **Content publishing**: manage drafts through editing, review, and publication stages
- **Research workflows**: process papers or data through analysis, synthesis, and validation
- **Dogfooding Spacedock's own development.** Spacedock is self-hosted. Its own development runs on a plain text workflow at [`docs/plans/`](docs/plans/). Run `skills/commission/bin/status --workflow-dir docs/plans` to see the current state.

## License

Spacedock is released under the [Apache License 2.0](LICENSE).
