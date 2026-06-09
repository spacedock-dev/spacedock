# Your first launch

Once the `spacedock` launcher and the host plugin are installed (see [Install](install.md)), a single command starts a session. Spacedock launches your first officer — the orchestrator agent — inside your coding harness.

## Launch the first officer

Point Spacedock at a project you already have and let it survey:

```bash
spacedock claude "/spacedock:survey"
```

This starts the first officer in Claude Code and runs the survey skill. Using Codex or Pi instead? Swap the subcommand:

```bash
spacedock codex "/spacedock:survey"   # or: spacedock pi "/spacedock:survey"
```

The first launch sets up the plugin for you, so this single line is enough. When a `.safehouse` profile is present in the working directory, the launch runs sandboxed.

## What the front door does

The launch command is the front door. Its grammar is:

```bash
spacedock claude "task" [--safehouse…] [-- host-flags…]
```

- The task comes first. It's handed to the first officer as the launch prompt.
- Anything after `--` forwards verbatim to the host (`claude` / `codex` / `pi`) — `--model`, `--resume`, `--plugin-dir`, and the like.
- `--safehouse` forces the launch through the sandbox; a `.safehouse` profile in the working directory does the same automatically.

## Survey what you already built

`/spacedock:survey` reads your project's existing agent history and reports three things:

- the workflow you've been running without naming it,
- how you've been calling work done,
- the decisions still open and waiting on you.

It then offers to commission a Spacedock workflow from what it found — which takes you to [your first workflow](first-workflow.md).
