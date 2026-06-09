# Command reference

The `spacedock` launcher groups its commands into Launch, Setup, and Workflow. Run `spacedock <command> --help` for the per-command detail.

## Launch

Start a coding host as your Spacedock first officer.

| Command | Description |
|---------|-------------|
| `spacedock claude [task] [-- claude-flags]` | Start Claude Code as your Spacedock first officer |
| `spacedock codex [task] [-- codex-flags]` | Start Codex as your Spacedock first officer |
| `spacedock pi [task] [-- pi-flags]` | Start Pi as your Spacedock first officer |

The optional `task` is the launch prompt. Everything after `--` forwards verbatim to the host. A `--plugin-dir` launch loads a local plugin checkout and relaxes the contract gate, so it does not require a prior `spacedock install`.

## Setup

| Command | Description |
|---------|-------------|
| `spacedock install [--host claude\|codex\|pi]` | Install the Spacedock plugin for a host, then check it |
| `spacedock doctor [--host claude\|codex\|pi]` | Check the installed plugin and this binary are compatible |

## Workflow

| Command | Description |
|---------|-------------|
| `spacedock status [args]` | Show or update workflow state |
| `spacedock new [--folder] SLUG` | Create an entity from a stdin body (auto-discovers the workflow) |
| `spacedock state init` | Initialize a cloned split-root workflow's state checkout |
| `spacedock completion bash\|zsh` | Print a bash or zsh completion script |
| `spacedock dispatch build \| show-stage-def` | Build worker dispatch artifacts |

## Other

| Command | Description |
|---------|-------------|
| `spacedock --version` | Print the installed version |

## Status queries

`spacedock status` reads workflow state. Common forms:

```bash
spacedock status --workflow-dir docs/dev            # the status table
spacedock status --workflow-dir docs/dev --next     # entities ready to dispatch
spacedock status --workflow-dir docs/dev --where sprint=0200-flip   # filter by frontmatter
spacedock status --workflow-dir docs/dev --validate # entity-contract validation
```
