# Your first launch

One command starts your first session and orients you in a project you already have:

```bash
spacedock claude "/spacedock:survey"
```

This launches the [first officer](../concepts/operating-model.md#three-roles) (the orchestrator agent that runs a Spacedock workflow) inside Claude Code and hands it the survey task. The first launch sets up the plugin for you; there is no separate setup step.

Run it from inside a project that already has some agent history, such as a repo you have been coding in with Claude Code. Survey reads that history; an empty directory has nothing to report on.

## What survey reports

In one read-only pass (it never edits your files), survey reconstructs what the agents in this project have been doing. Under a one-line headline (project, sessions, date range, decision and interruption counts), the report leads with **the decisions still open and waiting on you**. Then it names the [workflow](../concepts/workflows-and-entities.md) you have been running without naming it, the distinct workstreams, and how often you have had to step in. If the project has no agent history, survey says so and stops.

For the full section-by-section breakdown and how survey reads your history, see [what survey reports](../running-workflows/survey.md#what-it-reports).

## What comes next

Survey ends with an offer, not an action. It asks whether you want it to [commission](../running-workflows/commission.md) a real Spacedock workflow built from what it found, turning the open decisions into [approval gates](../concepts/gates-and-decisions.md) and the workstreams into work items. You can say yes and let it scaffold the workflow, or say no and keep the survey as a standalone orientation. Either way, nothing has changed in your project until you choose.

To go straight to defining a workflow yourself instead, see [your first workflow](first-workflow.md). If `spacedock` is not yet installed, start with [installing Spacedock](install.md).

## The command grammar

The front door is one shape, and every launch uses it:

```bash
spacedock claude "task" [--safehouse…] [-- host-flags…]
```

- **The task comes first.** It is handed to the first officer as the launch prompt. Here the task is `/spacedock:survey`, a skill the first officer runs. It could just as well be a plain sentence describing work.
- **`--safehouse` forces the launch through the sandbox.** A `.safehouse` profile in the working directory does the same automatically, so you only pass the flag when you want to force it.
- **Anything after `--` forwards verbatim to the host** (`claude` itself), including flags like `--resume`, `--model`, and `--plugin-dir`.

`spacedock codex "task"` and `spacedock pi "task"` take the same shape for the Codex and Pi harnesses. Claude Code is the primary surface; the examples here use it.
