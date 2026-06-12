# Command reference

The `spacedock` binary has ten subcommands in three groups, plus a top-level `--version`. This page is the map: what each command is for and when you reach for it. For the exact flags of any command, run `spacedock <command> --help`, which is the always-current source of truth.

| Command | Group | What it does |
|---------|-------|--------------|
| [`spacedock claude`](#launch) | Launch | Start Claude Code with the first officer loaded |
| [`spacedock codex`](#launch) | Launch | Start Codex (experimental) with the first officer |
| [`spacedock pi`](#launch) | Launch | Start Pi (experimental) with the first officer |
| [`spacedock install`](#setup) | Setup | Install the per-host plugin, then run the compatibility check |
| [`spacedock doctor`](#setup) | Setup | Run the compatibility check alone |
| [`spacedock status`](#workflow) | Workflow | Read or mutate workflow state: the table, `--next`, `--where`, `--set` |
| [`spacedock new`](#workflow) | Workflow | Create an entity from stdin |
| [`spacedock state`](#workflow) | Workflow | Manage a split-root workflow's state checkout |
| [`spacedock completion`](#workflow) | Workflow | Print a bash or zsh completion script |
| [`spacedock dispatch`](#workflow) | Workflow | Build the dispatch artifacts the first officer hands an ensign |
| `spacedock --version` | | Print the binary version and the contract level |

Run `spacedock` with no arguments for the grouped help, and `spacedock <command> --help` for a command's own flags.

## Launch

`spacedock claude`, `spacedock codex`, and `spacedock pi` start the named host with the first officer loaded. Claude Code is the primary surface; Codex and Pi are experimental. The grammar is the same for all three:

```bash
spacedock claude [task] [spacedock-flags] [-- host-flags]
```

The task comes first and becomes the launch prompt. Anything after `--` forwards verbatim to the host (`--model`, `--resume`, and the like). When no plugin is installed, the launcher auto-installs it and launches, so the single command yields a working session; a contract mismatch fails fast. The sandbox flags (`--safehouse` and its knobs) and the contract-gate flags are listed by `spacedock claude --help`.

## Setup

`spacedock install` installs the per-host plugin, then runs the compatibility check; `spacedock doctor` runs the check alone. Both take `--host claude|codex|pi` (default `claude`). When `doctor` reports the plugin is out of date, refresh it with `spacedock install`. See [Install Spacedock](../get-started/install.md) for the full setup path.

## Workflow

These commands read and mutate workflow state.

- **`status`** is the main one: with no flag it prints the entity table, `--next` lists what is ready to dispatch, `--where` filters, `--set` mutates frontmatter, `--validate` checks the workflow, and `--boot` prints the first-officer boot view. It resolves the workflow from `--workflow-dir`, then `PIPELINE_DIR`, then by walking up to the enclosing workflow. The day-to-day reads are covered in [Operating a workflow](../running-workflows/operating.md).
- **`new`** creates an entity from stdin.
- **`state`** initializes or creates the state checkout for a [split-root workflow](../advanced/split-root-state.md).
- **`completion`** prints a bash or zsh completion script.
- **`dispatch`** builds the worker dispatch artifacts the first officer hands an ensign.

Run `spacedock status --help` (and the same for each command) for the full flag list, the mutation guards, and the exit codes.
