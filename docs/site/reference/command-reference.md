# Command reference

The `spacedock` binary groups its subcommands into Launch, Setup, and Workflow, plus a top-level `spacedock --version` (the binary version and contract level). For the exact flags of any command, run `spacedock <command> --help`, the always-current source of truth; `spacedock` with no arguments prints the grouped help.

## Launch

`spacedock claude`, `spacedock codex`, and `spacedock pi` start a host with the first officer loaded. Claude Code is the primary surface; Codex and Pi are experimental. The grammar is the same for all three:

```bash
spacedock claude [task] [spacedock-flags] [-- host-flags]
```

The task comes first and becomes the launch prompt. Anything after `--` forwards verbatim to the host (`--model`, `--resume`, and the like). When no plugin is installed, the launcher auto-installs it and launches, so the single command yields a working session; a contract mismatch fails fast. The sandbox flags (`--safehouse` and its knobs) and the contract-gate flags are listed by `spacedock claude --help`.

## Setup

| Command | What it does |
|---------|--------------|
| `spacedock install` | Install the per-host plugin, then run the compatibility check |
| `spacedock doctor` | Run the compatibility check alone |

Both take `--host claude|codex|pi` (default `claude`). When `doctor` reports the plugin is out of date, refresh it with `spacedock install`. See [Install Spacedock](../get-started/install.md) for the full setup path.

## Workflow

The first officer runs these against workflow state as it moves entities; you operate through it, not by hand. They are documented here for completeness and for the rare direct use (scripting, debugging, restoring a state checkout on a fresh clone).

| Command | What it does |
|---------|--------------|
| `spacedock status` | Read or mutate the state: the entity table, `--next`, `--where`, `--set`, `--validate`, `--boot` |
| `spacedock new` | Create an entity (`new [--folder] SLUG`) from a body on stdin |
| `spacedock dispatch` | Build the worker dispatch artifacts (`dispatch build`, `dispatch show-stage-def`) |
| `spacedock state` | Manage a [split-root workflow](../advanced/split-root-state.md)'s state checkout (`state init` resumes one on a fresh clone, `state new` births one) |
| `spacedock completion` | Print a bash or zsh completion script |

[Operate a workflow](../running-workflows/operating.md) covers how the first officer uses `status` on your behalf. Run `spacedock status --help` (and the same for each command) for the full flag list, the mutation guards, and the exit codes.
