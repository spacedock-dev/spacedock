# Command reference

Run `spacedock <command> --help` for the full options. Run `spacedock` for the command list.

## Version

`spacedock --version` reports the version, platform and release channel. Inside an agent session, it also reports the runtime and sandbox.

## Launch

| Command | Agent |
| --- | --- |
| `spacedock claude` | Claude Code |
| `spacedock codex` | Codex |
| `spacedock pi` | Pi |

Claude Code is the primary Workflow host. Codex and Pi support is experimental.

Launch from your project directory. You can include a task:

```bash
spacedock claude "continue the workflow"
```

Put agent-specific options after `--`. See [Supported sandboxes](sandbox.md) for sandbox options and [Local plugin development](../advanced/plugin-development.md) for local checkouts.

## Setup

| Command | Purpose |
| --- | --- |
| `spacedock install` | Install or refresh the agent integration |
| `spacedock doctor` | Check installation and compatibility |

Both accept `--host claude`, `--host codex` or `--host pi`. The default is Claude Code. Follow the doctor output to resolve setup problems.

See [Installation](../get-started/install.md) for the manual setup path.

## Workflow

Your agent uses these commands to operate the workflow.

| Command | Purpose |
| --- | --- |
| `spacedock status` | Read or update work item state |
| `spacedock new` | Create a work item |
| `spacedock dispatch` | Prepare worker assignments |
| `spacedock gate` | Prepare reviews and record decisions |
| `spacedock merge guard` | Complete delivery and record the outcome |
| `spacedock state` | Manage the workflow state checkout |
| `spacedock completion` | Generate shell completions |

See [Workflow operation](../running-workflows/operating.md) for the conversation you use to direct the work.

<a id="-version"></a>
<a id="state-commit"></a>
