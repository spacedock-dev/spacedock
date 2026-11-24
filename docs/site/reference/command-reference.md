# Command reference

Run `spacedock <command> --help` for the full options. Run `spacedock` for the command list.

## Version

`spacedock --version` reports the version, platform, and release channel. Inside a supported or Paseo-hosted agent session, it also reports the runtime, sandbox, and contract level.

An unstamped source build includes a seven-character commit ID when Go provides VCS information: `0.29.0-pre0+dev.5210d36`.
A modified checkout adds `.dirty`. Without VCS information, the version ends in `+dev`. Linker-stamped versions stay unchanged.

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

Once per UTC day, an interactive launch asks `spacedock.md` for the newest Spacedock release, following its redirect to GitHub's releases API, and prints an update hint on stdout when it is newer.
Neither request carries an identifier in its URL or headers, though both servers see the connecting IP address like any request.
Use `--no-update-check` before `--`, or set `SPACEDOCK_NO_UPDATE_CHECK=1`, to skip the check.
Machine modes, non-TTY stdout, help, version, completions and agent-facing subcommands such as `status` and `dispatch` never check.
The day is recorded under `${XDG_STATE_HOME:-$HOME/.local/state}/spacedock/version-checks/`.
A slow or failed check stays silent and does not hold up the launch.

## Setup

| Command | Purpose |
| --- | --- |
| `spacedock install` | Install or refresh the agent integration |
| `spacedock doctor` | Check installation and compatibility |
| `spacedock update` | Update Spacedock, Review, Compass and their plugins to the newest release on this binary's channel |

`install` and `doctor` accept `--host claude`, `--host codex` or `--host pi`. The default is Claude Code. Follow the doctor output to resolve setup problems.

`update` prints each step it runs to stderr. `update --verbose` also prefixes every line with the time since the update started and reports how long each step took.

See [Installation](../get-started/install.md) for the manual setup path.

## Workflow

Your agent uses these commands to operate the workflow.

| Command | Purpose |
| --- | --- |
| `spacedock status` | Read or update work item state |
| `spacedock new` | Create a work item |
| `spacedock dispatch` | Prepare worker assignments |
| `spacedock gate` | Prepare reviews, record decisions and read gate history |
| `spacedock merge guard` | Complete delivery and record the outcome |
| `spacedock state` | Manage the workflow state checkout |
| `spacedock completion` | Generate shell completions |

`spacedock dispatch build --dispatch-dir DIR` selects the assignment directory. The flag takes precedence over `SPACEDOCK_DISPATCH_DIR`.
The default is `/tmp/spacedock-dispatch`. Relative paths resolve from the current directory. Empty or invalid values cause an error.
Assignment filenames include an eight-character content hash. Equal assignments reuse the file. Conflicting content causes an error without replacement.
If workers use a private `/tmp`, select a directory that both processes can read.

`spacedock gate read <ref-or-path> --json` returns the stored `gates` object, including every attempt and its application state.
References and paths use the same resolution as `status --read`, including archived entities. The command changes no files.
Records and attempts retain their stored order. Scalar values retain their stored values.
An entity without gate history returns `{"records":[],"version":1}`.
Invalid gate data exits nonzero without JSON output.

```bash
spacedock gate read gate-history-json-reader --json --workflow-dir docs/dev
spacedock gate read archive:finished-task --json --workflow-dir docs/dev
```

See [Workflow operation](../running-workflows/operating.md) for the conversation you use to direct the work.

<a id="-version"></a>
<a id="state-commit"></a>
