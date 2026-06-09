# Your first launch

One command starts your first session and orients you in a project you already have:

```bash
spacedock claude "/spacedock:survey"
```

This launches the first officer — the orchestrator agent that runs a Spacedock workflow — inside Claude Code, and hands it the survey task. The first launch sets up the plugin for you, so this single line is enough; you do not need a separate setup step. When a `.safehouse` profile is present in the working directory, the launch runs sandboxed automatically.

Run it from inside a project that already has some agent history — a repo you have been coding in with Claude Code. Survey reads that history; an empty directory has nothing to report on.

## The command grammar

The front door is one shape, and every launch uses it:

```bash
spacedock claude "task" [--safehouse…] [-- host-flags…]
```

- **The task comes first.** It is handed to the first officer as the launch prompt. Here the task is `/spacedock:survey`, a skill the first officer runs. It could just as well be a plain sentence describing work.
- **`--safehouse` forces the launch through the sandbox.** A `.safehouse` profile in the working directory does the same automatically, so you only pass the flag when you want to force it.
- **Anything after `--` forwards verbatim to the host** — `claude` itself — including flags like `--resume`, `--model`, and `--plugin-dir`.

`spacedock codex "task"` and `spacedock pi "task"` take the same shape for the Codex and Pi harnesses. Claude Code is the primary surface; the examples here use it.

## What survey reports

Survey reconstructs what the agents in this project have implicitly been doing, then reports it back — read-only. It never edits your files. It reads your recorded agent session history (through the `agentsview` tool, which it offers to install if it is missing), and it shows you, in one pass:

- **The inferred workflow** — the loop you have been running without naming it, written out as a stage chain.
- **The workstreams** — the distinct tracks of work, each labeled as mechanical (a routine loop) or exploration (human-driven creative work).
- **The decisions still open and waiting on you** — the forks that were raised but never resolved. Survey leads with these, because they are the work that is actually blocked on you.
- **How often you had to step in** — a count of the interruptions and decision points across those sessions, so you can see where your attention has been going.

A typical report opens with a one-line headline naming the project, the number of sessions, the date range, and the decision and interruption counts, then lays out each section below it. If the project has no agent history, survey says so plainly and stops — there is nothing to misreport.

## What comes next

Survey ends with an offer, not an action. After the report, it asks whether you want it to commission a real Spacedock workflow built from what it found — turning the open decisions into approval gates and the workstreams into work items. You can say yes and let it scaffold the workflow, or say no and keep the survey as a standalone orientation. Either way, nothing has changed in your project until you choose.

To go straight to defining a workflow yourself instead, see [your first workflow](first-workflow.md). If `spacedock` is not yet installed, start with [installing Spacedock](install.md).
